// Copyright 2015 Eryx <evorui аt gmаil dοt cοm>, All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package zonemaster

import (
	"errors"
	"fmt"
	"plugin"
	"strings"
	"time"

	"github.com/hooto/hlog4g/hlog"
	"github.com/hooto/iam/iamapi"
	"github.com/hooto/iam/iamclient"
	"github.com/lessos/lessgo/crypto/idhash"
	"github.com/lessos/lessgo/types"

	inCfg "github.com/sysinner/incore/config"
	"github.com/sysinner/incore/data"
	"github.com/sysinner/incore/inapi"
	typeScheduler "github.com/sysinner/incore/inapi/scheduler"
	pScheduler "github.com/sysinner/incore/plugin/scheduler"
	"github.com/sysinner/incore/status"
)

var (
	Scheduler         typeScheduler.Scheduler
	zonePodSpecPlans  inapi.PodSpecPlans
	podResFreeTimeMin uint32 = 60
	podInQueue        types.ArrayString
	hostResUsages     = map[string]*hostUsageItem{}
	scheduleHostList  typeScheduler.ScheduleHostList
	errServerError    = errors.New("server error")
)

type destResReplica struct {
	Ports       inapi.ServicePorts `json:"ports,omitempty"`
	ResCpu      int32              `json:"res_cpu,omitempty"` // 1 = .1 cores
	ResMem      int32              `json:"res_mem,omitempty"` // MB
	VolSys      int32              `json:"vol_sys,omitempty"` // GB
	VolSysAttrs uint32             `json:"vol_sys_attrs,omitempty"`
}

type hostUsageItem struct {
	cpu   int32 // 1 = .1 cores
	mem   int32 // MB
	vols  typeScheduler.ScheduleHostVolumes
	ports types.ArrayUint32
	box   int32
}

func SetupScheduler() error {

	if inCfg.Config.ZoneMasterSchedulerPlugin != "" {

		p, err := plugin.Open(inCfg.Prefix + "/plugin/" + inCfg.Config.ZoneMasterSchedulerPlugin)
		if err != nil {
			return err
		}

		nc, err := p.Lookup("NewConnector")
		if err != nil {
			return err
		}

		fn, ok := nc.(func() (typeScheduler.Scheduler, error))
		if !ok {
			return fmt.Errorf("No Plugin/Method (%s) Found", "NewConnector")
		}

		cn, err := fn()
		if err != nil {
			return err
		}

		Scheduler = cn

	} else {
		Scheduler = pScheduler.NewScheduler()
	}

	return nil
}

func scheduleAction() error {

	if status.ZoneId == "" ||
		status.Zone == nil ||
		!status.ZoneHostListImported {
		return errors.New("Zone Status Not Ready")
	}

	if Scheduler == nil {
		return errors.New("No Scheduler Found")
	}

	if err := schedulePodSpecPlanListRefresh(); err != nil {
		return err
	}

	hostResUsages = map[string]*hostUsageItem{}

	if err := schedulePodListRefresh(); err != nil {
		return err
	}

	if err := scheduleHostListRefresh(); err != nil {
		return err
	}

	//
	podInQueue.Clean()

	//
	for _, cell := range status.Zone.Cells {
		schedulePodListQueue(cell.Meta.Id)
	}

	schedulePodListBound()

	return nil
}

func schedulePodSpecPlanListRefresh() error {

	rs := data.GlobalMaster.PvScan(inapi.NsGlobalPodSpec("plan", ""), "", "", 1000)
	if !rs.OK() {
		return errors.New("gm/db error")
	}

	zonePodSpecPlans = inapi.PodSpecPlans{}
	rss := rs.KvList()
	for _, v := range rss {
		var specPlan inapi.PodSpecPlan
		if err := v.Decode(&specPlan); err == nil {
			specPlan.ChargeFix()
			zonePodSpecPlans = append(zonePodSpecPlans, &specPlan)
		}
	}

	if len(zonePodSpecPlans) > 0 {
		return nil
	}

	return errors.New("No PodSpecPlan Found")
}

func schedulePodListRefresh() error {

	//
	var (
		tn            = uint32(time.Now().Unix())
		chg           = false
		podServiceChg = false
	)

	// local zone scheduled pods
	rs := data.ZoneMaster.PvScan(
		inapi.NsZonePodInstance(status.ZoneId, ""), "", "", 10000)
	if !rs.OK() {
		return errors.New("zm/db err")
	}

	rss := rs.KvList()
	for _, v := range rss {

		var srcPod inapi.Pod
		if err := v.Decode(&srcPod); err != nil {
			hlog.Printf("warn", "zm/pod data/struct err %s", err.Error())
			continue
		}

		pod := status.ZonePodList.Items.Get(srcPod.Meta.ID)
		if pod == nil ||
			srcPod.Operate.Version > pod.Operate.Version {
			//
			hlog.Printf("info", "zm/pod set %s", srcPod.Meta.ID)
			status.ZonePodList.Items.Set(&srcPod)
			pod = &srcPod
		}

		// destroy
		if inapi.OpActionAllow(pod.Operate.Action, inapi.OpActionDestroy) {

			//
			if (pod.Operate.Operated + inapi.PodDestroyTTL) < tn {

				if rs := data.ZoneMaster.KvPut(inapi.NsKvZonePodInstanceDestroy(status.ZoneId, pod.Meta.ID), pod, nil); rs.OK() {
					rs = data.ZoneMaster.PvDel(inapi.NsZonePodInstance(status.ZoneId, pod.Meta.ID), nil)
					//
					status.ZonePodList.Items.Del(pod.Meta.ID)
					hlog.Printf("warn", "zm/scheduler pod %s, remove", pod.Meta.ID)
				}
			}
			continue
		}

		// TODO
		if inapi.OpActionAllow(pod.Operate.Action, inapi.OpActionStop) &&
			len(pod.Operate.RepMigrates) < 1 &&
			!inapi.OpActionAllow(pod.Operate.Action, inapi.OpActionResFree) &&
			(pod.Operate.Operated+podResFreeTimeMin) < tn {

			pod.Operate.Action = pod.Operate.Action | inapi.OpActionResFree
			//
			if rs := data.ZoneMaster.PvPut(inapi.NsZonePodInstance(status.ZoneId, pod.Meta.ID), pod, nil); !rs.OK() {
				hlog.Printf("info", "zm/scheduler pod %s, db err %s", pod.Meta.ID, rs.Bytex().String())
				continue
			}

			pod.Operate.OpLog, _ = inapi.PbOpLogEntrySliceSync(
				pod.Operate.OpLog,
				inapi.NewPbOpLogEntry(inapi.OpLogNsZoneMasterPodScheduleResFree, inapi.PbOpLogOK, "free CPU/RAM resources on host"),
			)

			hlog.Printf("warn", "zm/scheduler pod %s, stop and res-free", pod.Meta.ID)
		}

		specRes := pod.Spec.ResComputeBound()

		if specRes == nil || pod.Spec.VolSys == nil {
			hlog.Printf("warn", "zm/scheduler pod %s, Spec.ResCompute or Spec.Volume(system) not found", pod.Meta.ID)
			continue
		}

		for _, rp := range pod.Operate.Replicas {

			if rp.Node == "" {
				continue
			}

			if inapi.OpActionAllow(rp.Action, inapi.OpActionDestroy) &&
				!inapi.OpActionAllow(rp.Action, inapi.OpActionMigrate) {
				continue
			}

			hostRes, ok := hostResUsages[rp.Node]
			if !ok {
				hostRes = &hostUsageItem{}
				hostResUsages[rp.Node] = hostRes
			}

			if !inapi.OpActionAllow(pod.Operate.Action, inapi.OpActionResFree) {
				hostRes.cpu += specRes.CpuLimit
				hostRes.mem += specRes.MemLimit
				hostRes.box += 1
			}

			if rp.VolSysMnt == "" {
				rp.VolSysMnt = "/opt"
			}

			hostResVol := hostRes.vols.Get(rp.VolSysMnt)
			if hostResVol == nil {
				hostResVol = &typeScheduler.ScheduleHostVolume{
					Name: rp.VolSysMnt,
				}
				hostRes.vols.Sync(hostResVol)
			}

			hostResVol.Used += pod.Spec.VolSys.Size

			for _, rpp := range rp.Ports {

				if rpp.HostPort == 0 {
					continue
				}

				hostRes.ports.Set(uint32(rpp.HostPort))
			}
		}

		// refresh pod's service endpoints
		zmPodService := inapi.AppServicePodSliceGet(status.ZonePodServices.Items, pod.Meta.ID)
		if zmPodService == nil {
			zmPodService = &inapi.AppServicePod{
				PodId:   pod.Meta.ID,
				Updated: inapi.TimeNowMs(),
			}
			status.ZonePodServices.Items, _ = inapi.AppServicePodSliceSync(status.ZonePodServices.Items, zmPodService)
			hlog.Printf("info", "zm/pod %s, service init", pod.Meta.ID)
			podServiceChg = true
		}
		//
		for _, app := range pod.Apps {

			for _, appSpecPort := range app.Spec.ServicePorts {

				appSpecId := app.Spec.Meta.ID
				if appSpecPort.AppSpec != "" {
					appSpecId = appSpecPort.AppSpec
				}

				asp := inapi.AppServicePortSliceGet(zmPodService.Ports, uint32(appSpecPort.BoxPort))
				if asp == nil {
					asp = &inapi.AppServicePort{
						Spec:    appSpecId,
						Port:    uint32(appSpecPort.BoxPort),
						Updated: inapi.TimeNowMs(),
						Name:    appSpecPort.Name,
					}
					zmPodService.Ports, _ = inapi.AppServicePortSliceSync(zmPodService.Ports, asp)
					hlog.Printf("info", "zm/scheduler pod %s, app-spec %s, service-port %d, init",
						pod.Meta.ID, app.Spec.Meta.ID, appSpecPort.BoxPort)
					podServiceChg = true
				}
				if asp.Name != appSpecPort.Name {
					asp.Name = appSpecPort.Name
				}

				for _, rp := range pod.Operate.Replicas {

					if rp.Node == "" {
						continue
					}

					repIp := status.ZoneHostIp(rp.Node)
					if repIp == "" {
						continue
					}

					for _, rpp := range rp.Ports {

						if rpp.HostPort == 0 || rpp.BoxPort != appSpecPort.BoxPort {
							continue
						}

						aspRep := &inapi.AppServiceReplica{
							Rep:  rp.RepId,
							Port: uint32(rpp.HostPort),
							Ip:   repIp,
						}
						if asp.Endpoints, chg = inapi.AppServiceReplicaSliceSync(asp.Endpoints, aspRep); chg {
							asp.Updated = inapi.TimeNowMs()
							podServiceChg = true
							hlog.Printf("info", "zm/scheduler pod %s, rep %d, addr %s:%d app-service-replica refresh",
								pod.Meta.ID, rp.RepId, aspRep.Ip, aspRep.Port)
						}
					}
				}
			}
		}
	}

	if podServiceChg {
		status.ZonePodServicesFlush()
	}

	// inapi.ObjPrint("ns", status.ZonePodServices.Items)

	return nil
}

func scheduleHostListRefresh() error {

	scheduleHostList.Items = []*typeScheduler.ScheduleHostItem{}

	var (
		cellStatuses = map[string]*inapi.ResCellStatus{}
		tn           = uint32(time.Now().Unix())
	)

	//
	for _, host := range status.ZoneHostList.Items {

		if host.Operate == nil ||
			host.Spec == nil || host.Spec.Capacity == nil {
			continue
		}

		var (
			res, ok         = hostResUsages[host.Meta.Id]
			sync            = false
			cellStatus, cok = cellStatuses[host.Operate.CellId]
			vols            = typeScheduler.ScheduleHostVolumes{}
		)
		if !cok {
			cellStatus = &inapi.ResCellStatus{
				Updated: tn,
			}
			cellStatuses[host.Operate.CellId] = cellStatus
		}

		if !ok {
			res = &hostUsageItem{}
			hostResUsages[host.Meta.Id] = res
		}

		for _, scv := range host.Spec.Capacity.Vols {

			cellStatus.VolCap += scv.Value

			vol := &typeScheduler.ScheduleHostVolume{
				Name:  scv.Name,
				Total: scv.Value,
				Attrs: scv.Attrs,
			}

			if res != nil {
				if pv := res.vols.Get(scv.Name); pv != nil {
					vol.Used = pv.Used
				}
			}

			vols.Sync(vol)
		}

		if !ok {

			if host.Operate.CpuUsed > 0 {
				host.Operate.CpuUsed, sync = 0, true
			}

			if host.Operate.MemUsed > 0 {
				host.Operate.MemUsed, sync = 0, true
			}

			if len(host.Operate.VolUsed) > 0 {
				host.Operate.VolUsed, sync = []*inapi.ResVolValue{}, true
			}

			if len(host.Operate.PortUsed) > 0 {
				host.Operate.PortUsed, sync = []uint32{}, true
			}

			if host.Operate.BoxNum > 0 {
				host.Operate.BoxNum, sync = 0, true
			}

		} else {

			if host.Operate.CpuUsed != res.cpu {
				host.Operate.CpuUsed, sync = res.cpu, true
			}

			if host.Operate.MemUsed != int64(res.mem) {
				host.Operate.MemUsed, sync = int64(res.mem), true
			}

			if !res.vols.Equal(vols) {
				volNews := []*inapi.ResVolValue{}
				for _, v := range res.vols {
					volNews = append(volNews, &inapi.ResVolValue{
						Name:  v.Name,
						Value: v.Used,
					})
				}
				host.Operate.VolUsed, sync = volNews, true
			}

			if !res.ports.Equal(host.Operate.PortUsed) {
				host.Operate.PortUsed, sync = res.ports, true
			}

			if host.Operate.BoxNum != res.box {
				host.Operate.BoxNum, sync = res.box, true
			}
		}

		// inapi.ObjPrint("AA", vols)

		cellStatus.CpuCap += int64(host.Spec.Capacity.Cpu)
		cellStatus.MemCap += host.Spec.Capacity.Mem

		cellStatus.CpuUsed += int64(host.Operate.CpuUsed)
		cellStatus.MemUsed += host.Operate.MemUsed

		for _, vol := range host.Operate.VolUsed {
			cellStatus.VolUsed += vol.Value
		}

		cellStatus.HostCap += 1
		if host.Operate.Action == 1 {
			cellStatus.HostIn += 1
		}

		if sync {

			host.OpPortSort()

			if rs := data.ZoneMaster.PvPut(
				inapi.NsZoneSysHost(status.ZoneId, host.Meta.Id), host, nil,
			); !rs.OK() {
				return fmt.Errorf("host %s sync changes failed %s", host.Meta.Id, rs.Bytex().String())
			}
		}

		if pv := vols.Get("/opt"); pv != nil {
			if vols.Get("/") != nil {
				vols.Del("/")
			}
		}

		scheduleHostList.Items = append(scheduleHostList.Items, &typeScheduler.ScheduleHostItem{
			Id:               host.Meta.Id,
			CellId:           host.Operate.CellId,
			OpAction:         host.Operate.Action,
			CpuTotal:         host.Spec.Capacity.Cpu,
			CpuUsed:          host.Operate.CpuUsed,
			MemTotal:         int32(host.Spec.Capacity.Mem),
			MemUsed:          int32(host.Operate.MemUsed),
			Volumes:          vols,
			BoxDockerVersion: host.Spec.ExpDockerVersion,
			BoxPouchVersion:  host.Spec.ExpPouchVersion,
		})
	}

	for id, v := range cellStatuses {
		cell := status.Zone.Cell(id)
		if cell == nil {
			continue
		}
		if cell.Status != nil && (cell.Status.Updated+600) > v.Updated {
			continue
		}

		cell.Status = v
		if rs := data.GlobalMaster.PvPut(inapi.NsGlobalSysCell(status.Host.Operate.ZoneId, id), cell, nil); rs.OK() {
			data.ZoneMaster.PvPut(inapi.NsZoneSysCell(status.Zone.Meta.Id, id), cell, nil)
		}
	}

	return nil
}

func schedulePodListQueue(cellId string) {

	// TODO pager
	var (
		offset = inapi.NsKvGlobalSetQueuePod(status.ZoneId, cellId, "")
		cutset = inapi.NsKvGlobalSetQueuePod(status.ZoneId, cellId, "")
	)

	rss := data.GlobalMaster.KvScan(offset, cutset, 10000).KvList()
	if len(rss) == 0 {
		return
	}

	for _, v := range rss {

		var podq inapi.Pod

		if err := v.Decode(&podq); err != nil {
			hlog.Printf("error", "invalid data struct: %s", err.Error())
			continue
		}

		if podq.Spec == nil || podq.Spec.Cell != cellId {
			hlog.Printf("error", "invalid data struct: no spec/cell found")
			continue
		}

		var (
			pod = status.ZonePodList.Items.Get(podq.Meta.ID)
		)

		if rs := data.ZoneMaster.PvGet(
			inapi.NsZonePodInstance(status.ZoneId, podq.Meta.ID),
		); rs.OK() {

			var prev inapi.Pod
			if err := rs.Decode(&prev); err != nil {
				hlog.Printf("error", "bad prev podq %s instance, err %s", podq.Meta.ID, err.Error())
				continue
			}

			if pod == nil || pod.Operate.Version < prev.Operate.Version {

				hlog.Printf("info", "zm/pod set %s", prev.Meta.ID)
				status.ZonePodList.Items.Set(&prev)
				pod = &prev
			}

		} else if !rs.NotFound() {
			hlog.Printf("warn", "zm/scheduler db err")
			continue
		}

		if pod == nil {
			hlog.Printf("info", "zm/pod set %s", podq.Meta.ID)
			status.ZonePodList.Items.Set(&podq)
			pod = &podq
		}

		if podq.Operate.Version > pod.Operate.Version {

			pod.Meta = podq.Meta
			pod.Spec = podq.Spec

			apps := pod.Apps

			pod.Apps = podq.Apps

			for _, appPrev := range apps {

				if len(appPrev.Operate.Services) == 0 {
					continue
				}

				for _, app := range pod.Apps {

					if app.Meta.ID != appPrev.Meta.ID {
						continue
					}

					for _, posp := range appPrev.Operate.Services {

						for _, osp := range app.Operate.Services {

							if posp.Port == osp.Port {
								osp.Endpoints = posp.Endpoints
								break
							}
						}
					}

					break
				}
			}

			if inapi.OpActionAllow(podq.Operate.Action, inapi.OpActionRestart) {
				// podq.Operate.Action = inapi.OpActionRemote(podq.Operate.Action, inapi.OpActionRestart)
				hlog.Printf("info", "Scheduler Pod %s, restart", pod.Meta.ID)
			}

			pod.Operate.Action = podq.Operate.Action
			pod.Operate.Version = podq.Operate.Version
			pod.Operate.Priority = podq.Operate.Priority
			pod.Operate.ReplicaCap = podq.Operate.ReplicaCap
			pod.Operate.Operated = podq.Operate.Operated
			pod.Operate.Access = podq.Operate.Access
			pod.Operate.OpLog = podq.Operate.OpLog
			pod.Operate.BindServices = podq.Operate.BindServices
			pod.Operate.ExpSysState = podq.Operate.ExpSysState

			if len(podq.Operate.RepMigrates) > 0 {
				migrates := types.ArrayUint32{}
				for _, v := range pod.Operate.RepMigrates {
					migrates.Set(v)
				}
				for _, v := range podq.Operate.RepMigrates {
					migrates.Set(v)
				}
				pod.Operate.RepMigrates = migrates
			}
		}

		err := schedulePodItem(pod)
		if err != nil {
			hlog.Printf("warn", "Scheduler Pod %s, ER %s", pod.Meta.ID, err.Error())
		}

		if rs := data.ZoneMaster.PvPut(inapi.NsZonePodInstance(status.ZoneId, pod.Meta.ID), pod, nil); !rs.OK() {
			hlog.Printf("error", "zone/podq saved %s, err (%s)", pod.Meta.ID, rs.Bytex().String())
			continue
		}

		if err == nil {
			data.GlobalMaster.KvDel(inapi.NsKvGlobalSetQueuePod(status.ZoneId, pod.Spec.Cell, pod.Meta.ID), nil)
			hlog.Printf("info", "zone/podq queue/clean %s", pod.Meta.ID)
			podInQueue.Set(pod.Meta.ID)
		}
	}
}

func schedulePodListBound() {

	if len(status.ZonePodList.Items) > 0 {
		hlog.Printf("debug", "zone/pod/list %d", len(status.ZonePodList.Items))
	}

	for _, podq := range status.ZonePodList.Items {

		if podInQueue.Has(podq.Meta.ID) {
			continue
		}

		// hlog.Printf("info", "Scheduler Pod/Migrate N %d", len(podq.Operate.RepMigrates))

		if inapi.OpActionAllow(podq.Operate.Action, inapi.OpActionHang) {
			continue
		}

		if podq.Spec == nil {
			hlog.Printf("error", "invalid data struct: no spec/cell found")
			continue
		}

		err := schedulePodItem(podq)
		if err != nil {
			hlog.Printf("error", "Scheduler Pod %s, ER %s", podq.Meta.ID, err.Error())
		}

		if len(podq.Operate.OpLog) > 0 {
			//
		}

		if err == nil {

			err = schedulePodMigrate(podq)
			if err != nil {
				hlog.Printf("error", "Scheduler Pod/Migrate %s, ER %s", podq.Meta.ID, err.Error())
			}
		}

		if err == nil {
			err = schedulePodFailover(podq)
			if err != nil {
				hlog.Printf("error", "Scheduler Pod/Failover %s, ER %s", podq.Meta.ID, err.Error())
			}
		}

		if rs := data.ZoneMaster.PvPut(inapi.NsZonePodInstance(status.ZoneId, podq.Meta.ID), podq, nil); !rs.OK() {
			hlog.Printf("error", "zone/podq saved %s, err (%s)", podq.Meta.ID, rs.Bytex().String())
			continue
		}
	}
}

func schedulePodItem(podq *inapi.Pod) error {

	// hlog.Printf("error", "exec podq %s instance", podq.Meta.ID)

	if podq.Spec == nil {
		return errors.New("No PodSpec Setup")
	}

	//
	if podq.Spec.VolSys == nil {
		return errors.New("No Spec/VolSys Setup")
	}

	if inapi.OpActionAllow(podq.Operate.Action, inapi.OpActionDestroy) {
		//
	} else if inapi.OpActionAllow(podq.Operate.Action, inapi.OpActionStart) {
		podq.Operate.Action = inapi.OpActionRemove(podq.Operate.Action, inapi.OpActionStop)
		podq.Operate.Action = inapi.OpActionRemove(podq.Operate.Action, inapi.OpActionStopped)
	} else if inapi.OpActionAllow(podq.Operate.Action, inapi.OpActionStop) {
		podq.Operate.Action = inapi.OpActionRemove(podq.Operate.Action, inapi.OpActionStart)
		podq.Operate.Action = inapi.OpActionRemove(podq.Operate.Action, inapi.OpActionRunning)
	}

	// bugfix
	for _, ctrlRep := range podq.Operate.Replicas {

		if inapi.OpActionAllow(podq.Operate.Action, inapi.OpActionDestroy) {
			//
		} else if inapi.OpActionAllow(podq.Operate.Action, inapi.OpActionStart) {
			ctrlRep.Action = inapi.OpActionRemove(ctrlRep.Action, inapi.OpActionStop)
			ctrlRep.Action = inapi.OpActionRemove(ctrlRep.Action, inapi.OpActionStopped)
		} else if inapi.OpActionAllow(podq.Operate.Action, inapi.OpActionStop) {
			ctrlRep.Action = inapi.OpActionRemove(ctrlRep.Action, inapi.OpActionStart)
			ctrlRep.Action = inapi.OpActionRemove(ctrlRep.Action, inapi.OpActionRunning)
		}
	}

	for _, app := range podq.Apps {

		var ports types.ArrayUint32

		// ServicePort of Remotely dependent AppSpec
		for _, sp := range app.Operate.Services {

			//
			if sp.AppId == "" || sp.AppId == app.Meta.ID ||
				sp.PodId == "" || sp.PodId == podq.Meta.ID {
				continue
			}

			zmPodService := inapi.AppServicePodSliceGet(status.ZonePodServices.Items, sp.PodId)
			if zmPodService == nil {
				continue
			}

			srvPort := inapi.AppServicePortSliceGet(zmPodService.Ports, sp.Port)
			if srvPort == nil || srvPort.Updated <= sp.Updated {
				continue
			}

			chg := false
			if sp.Endpoints, chg = inapi.AppServiceReplicaSliceSyncSlice(sp.Endpoints, srvPort.Endpoints); chg {
				hlog.Printf("info", "pod %s, app %s, port %d, service endpoints refreshed ",
					podq.Meta.ID, app.Meta.ID, sp.Port)
			}
			sp.Updated = srvPort.Updated

			ports.Set(sp.Port)
		}

		//
		zmPodService := inapi.AppServicePodSliceGet(status.ZonePodServices.Items, podq.Meta.ID)
		if zmPodService == nil {
			continue
		}

		//
		for _, assp := range app.Spec.ServicePorts {

			if ports.Has(uint32(assp.BoxPort)) {
				continue
			}

			appSpecId := assp.AppSpec
			if appSpecId == "" {
				appSpecId = app.Spec.Meta.ID
			}

			srvPort := inapi.AppServicePortSliceGet(app.Operate.Services, uint32(assp.BoxPort))
			if srvPort == nil {

				hlog.Printf("info", "pod %s, app/depend-spec %s, port %d, service port init",
					podq.Meta.ID, appSpecId, assp.BoxPort)

				srvPort = &inapi.AppServicePort{
					Spec: appSpecId,
					Port: uint32(assp.BoxPort),
					Name: assp.Name,
				}

				app.Operate.Services, _ = inapi.AppServicePortSliceSync(app.Operate.Services, srvPort)

			} else {

				if srvPort.Name != assp.Name {
					srvPort.Name = assp.Name
				}

				if srvPort.Spec != appSpecId {
					srvPort.Spec = appSpecId
				}
			}
		}

		//
		for _, appOpService := range app.Operate.Services {

			if appOpService.PodId != "" {
				continue
			}

			srvPort := inapi.AppServicePortSliceGet(zmPodService.Ports, appOpService.Port)
			if srvPort == nil || srvPort.Updated <= appOpService.Updated {
				continue
			}

			chg := false
			if appOpService.Endpoints, chg = inapi.AppServiceReplicaSliceSyncSlice(appOpService.Endpoints, srvPort.Endpoints); chg {
				hlog.Printf("info", "pod %s, app %s, port %d, service endpoints refreshed ",
					podq.Meta.ID, app.Meta.ID, appOpService.Port)
			}
			appOpService.Updated = srvPort.Updated
		}

		for _, appBindService := range app.Operate.BindServices {

			srvPort := inapi.AppServicePortPodBindSliceGet(podq.Operate.BindServices,
				appBindService.Port, appBindService.PodId)
			if srvPort == nil {
				podq.Operate.BindServices, _ = inapi.AppServicePortPodBindSliceSync(podq.Operate.BindServices, &inapi.AppServicePortPodBind{
					Port:  appBindService.Port,
					PodId: appBindService.PodId,
				})
			}
		}
	}

	for _, v := range podq.Operate.BindServices {

		zmPodService := inapi.AppServicePodSliceGet(status.ZonePodServices.Items, v.PodId)
		if zmPodService == nil {
			continue
		}

		srvPort := inapi.AppServicePortSliceGet(zmPodService.Ports, v.Port)
		if srvPort == nil {
			continue
		}

		if v.Updated >= srvPort.Updated {
			continue
		}

		v.Endpoints = srvPort.Endpoints
		v.Updated = srvPort.Updated
	}

	if podq.OpResScheduleFit() {
		return nil
	}

	var (
		tnStart   = time.Now()
		destRes   *destResReplica
		scaleup   = 0
		scaledown = 0
	)

	// PreChargeValid
	if inapi.OpActionAllow(podq.Operate.Action, inapi.OpActionStart) {

		if err := schedulePodPreChargeValid(podq); err != nil {
			return errors.New("zm/scheduler pod/pre-valid err " + err.Error())
		}

		specRes := podq.Spec.ResComputeBound()
		destRes = &destResReplica{
			ResCpu:      specRes.CpuLimit,
			ResMem:      specRes.MemLimit,
			VolSys:      podq.Spec.VolSys.Size,
			VolSysAttrs: podq.Spec.VolSys.Attrs,
		}

	} else {

		destRes = &destResReplica{}

		if inapi.OpActionAllow(podq.Operate.Action, inapi.OpActionStop) {
			destRes.VolSys = podq.Spec.VolSys.Size
		}
	}

	if len(podq.Operate.Replicas) > 0 {
		podq.Operate.Replicas.Sort()
	}

	for repid := uint32(0); repid < uint32(podq.Operate.ReplicaCap); repid++ {

		oplog := schedulePodRepItem(podq, 0, repid, destRes)
		if oplog.Status == inapi.PbOpLogOK {
			scaleup += 1
			// podq.Operate.OpLog, _ = inapi.PbOpLogEntrySliceSync(
		} else {
			podq.Operate.OpLog, _ = inapi.PbOpLogEntrySliceSync(
				podq.Operate.OpLog, oplog)
		}
	}

	var (
		podStatus = status.ZonePodStatusList.Get(podq.Meta.ID)
		repOuts   = []*inapi.PodOperateReplica{}
	)

	// scaling down
	for _, rep := range podq.Operate.Replicas {

		if rep.RepId < uint32(podq.Operate.ReplicaCap) {
			continue
		}

		scaledown += 1

		if podStatus != nil {
			if repStatus := podStatus.RepGet(rep.RepId); repStatus != nil {
				if inapi.OpActionAllow(rep.Action, inapi.OpActionDestroy) &&
					inapi.OpActionAllow(repStatus.Action, inapi.OpActionDestroyed) {
					scaledown -= 1
					repOuts = append(repOuts, rep)
				}
			}
		}

		if !inapi.OpActionAllow(rep.Action, inapi.OpActionDestroy) {
			hlog.Printf("info", "zm/rep %s:%d destroy", podq.Meta.ID, rep.RepId)
			rep.Action = inapi.OpActionDestroy
			schedulePodRepItem(podq, inapi.OpActionDestroy, rep.RepId, &destResReplica{})
		}
	}

	for _, rep := range repOuts {

		hlog.Printf("info", "zm/pod %s, scaling down rep %d, clean out operate/replica",
			podq.Meta.ID, rep.RepId)

		podq.Operate.Replicas.Del(rep.RepId)
	}

	var (
		opType = inapi.PbOpLogOK
		opMsg  = fmt.Sprintf("schedule %d/%d replicas in %v",
			scaleup, podq.Operate.ReplicaCap, time.Since(tnStart))
	)

	if scaleup != podq.Operate.ReplicaCap ||
		scaledown > 0 {
		opType = inapi.PbOpLogWarn
	}

	podq.Operate.OpLog, _ = inapi.PbOpLogEntrySliceSync(
		podq.Operate.OpLog,
		inapi.NewPbOpLogEntry(inapi.OpLogNsZoneMasterPodScheduleAlloc, opType, opMsg),
	)

	if inapi.OpActionAllow(podq.Operate.Action, inapi.OpActionResFree) &&
		!inapi.OpActionAllow(podq.Operate.Action, inapi.OpActionStop) {
		podq.Operate.Action = inapi.OpActionRemove(podq.Operate.Action, inapi.OpActionResFree)
	}

	return nil
}

func schedulePodRepItem(podq *inapi.Pod, opAction uint32,
	repId uint32, destRes *destResReplica) *inapi.PbOpLogEntry {

	var (
		host        *inapi.ResHost
		hostChanged = false
		opRep       = podq.Operate.Replicas.Get(repId)
		opLogKey    = inapi.OpLogNsZoneMasterPodScheduleRep(repId)
	)

	if opAction == 0 {
		opAction = podq.Operate.Action
	}

	//
	if opRep == nil {

		opRep = &inapi.PodOperateReplica{
			RepId:  repId,
			ResCpu: destRes.ResCpu,
			ResMem: destRes.ResMem,
			VolSys: destRes.VolSys,
		}

		hit, err := Scheduler.ScheduleHost(
			&typeScheduler.SchedulePodSpec{
				CellId:    podq.Spec.Cell,
				BoxDriver: podq.Spec.Box.Image.Driver,
			},
			&typeScheduler.SchedulePodReplica{
				RepId:       opRep.RepId,
				Cpu:         opRep.ResCpu,
				Mem:         opRep.ResMem,
				VolSys:      opRep.VolSys,
				VolSysAttrs: destRes.VolSysAttrs,
			},
			&scheduleHostList,
			nil,
		)

		if err != nil || hit.Host == nil || len(hit.Volumes) < 1 {
			// TODO error log
			return inapi.NewPbOpLogEntry(opLogKey,
				inapi.PbOpLogWarn,
				"no available resources, waiting for allocation")
		}

		host = status.ZoneHostList.Item(hit.HostId)
		if host == nil {
			return inapi.NewPbOpLogEntry(opLogKey,
				inapi.PbOpLogWarn,
				"no available resources, waiting for allocation")
		}

		opRep.Node = hit.HostId
		opRep.Scheduled = uint32(time.Now().Unix())

		// TODO
		opRep.VolSysMnt = hit.Volumes[0].Name

		podq.Operate.Replicas.Set(*opRep)
		podq.Operate.Replicas.Sort()

		host.SyncOpCpu(destRes.ResCpu)
		host.SyncOpMem(destRes.ResMem)

		hit.Host.CpuUsed += destRes.ResCpu
		hit.Host.MemUsed += destRes.ResMem

		for _, hitVol := range hit.Volumes {
			if pv := inapi.ResVolValueSliceGet(host.Operate.VolUsed, hitVol.Name); pv != nil {
				pv.Value += hitVol.Size
			}
			if pv := hit.Host.Volumes.Get(hitVol.Name); pv != nil {
				pv.Used += hitVol.Size
			}
		}

		hostChanged = true

		// TOTK
		if rs := data.ZoneMaster.PvPut(inapi.NsZonePodInstance(status.ZoneId, podq.Meta.ID), podq, nil); !rs.OK() {
			hlog.Printf("error", "zone/podq saved %s, err (%s)", podq.Meta.ID, rs.Bytex().String())
			return inapi.NewPbOpLogEntry("", inapi.PbOpLogWarn, "Data IO error")
		}

		hlog.Printf("info", "schedule rep %s:%d to host %s (new)",
			podq.Meta.ID, opRep.RepId, hit.HostId)

		podq.Operate.OpLog, _ = inapi.PbOpLogEntrySliceSync(podq.Operate.OpLog, inapi.NewPbOpLogEntry(opLogKey,
			inapi.PbOpLogOK, fmt.Sprintf("schedule rep %s:%d to host %s", podq.Meta.ID, opRep.RepId, hit.HostId)))

	} else {

		if inapi.OpActionAllow(opRep.Action, inapi.OpActionMigrate) {
			return inapi.NewPbOpLogEntry(opLogKey,
				inapi.PbOpLogWarn,
				"action in migrating")
		}

		host = status.ZoneHostList.Item(opRep.Node)
		if host == nil {
			return inapi.NewPbOpLogEntry(opLogKey,
				inapi.PbOpLogWarn,
				fmt.Sprintf("host %s not found", opRep.Node))
		}
	}

	if !inapi.OpActionAllow(opRep.Action, opAction) {

		if inapi.OpActionAllow(opAction, inapi.OpActionDestroy) {
			opRep.Action = inapi.OpActionDestroy
		} else if inapi.OpActionAllow(opRep.Action, inapi.OpActionDestroy) ||
			inapi.OpActionAllow(opRep.Action, inapi.OpActionMigrate) {
			//
		} else if inapi.OpActionAllow(opAction, inapi.OpActionStart) {
			opRep.Action = inapi.OpActionRemove(opRep.Action, inapi.OpActionStop) | inapi.OpActionStart
		} else if inapi.OpActionAllow(opAction, inapi.OpActionStop) {
			opRep.Action = inapi.OpActionRemove(opRep.Action, inapi.OpActionStart) | inapi.OpActionStop
		}
	}

	//
	if opRep.ResCpu != destRes.ResCpu ||
		opRep.ResMem != destRes.ResMem ||
		opRep.VolSys != destRes.VolSys {

		hlog.Printf("info", "rep %s-%d spec change", podq.Meta.ID, opRep.RepId)

		if inapi.OpActionAllow(opAction, inapi.OpActionStart) {

			//
			if err := Scheduler.ScheduleHostValid(
				&typeScheduler.ScheduleHostItem{
					Id:       host.Meta.Id,
					CpuTotal: host.Spec.Capacity.Cpu,
					CpuUsed:  host.Operate.CpuUsed,
					MemTotal: int32(host.Spec.Capacity.Mem),
					MemUsed:  int32(host.Operate.MemUsed),
				},
				&typeScheduler.SchedulePodReplica{
					Cpu:         destRes.ResCpu - opRep.ResCpu,
					Mem:         destRes.ResMem - opRep.ResMem,
					VolSys:      destRes.VolSys - opRep.VolSys,
					VolSysAttrs: destRes.VolSysAttrs,
				},
			); err != nil {

				hlog.Printf("warn", "rep %s-%d, err %s",
					podq.Meta.ID, opRep.RepId, err.Error())
				//
				oplog := inapi.NewPbOpLogEntry(opLogKey,
					inapi.PbOpLogWarn,
					"no available resources (spec update), waiting for allocation")
				podq.Operate.OpLog, _ = inapi.PbOpLogEntrySliceSync(
					podq.Operate.OpLog,
					oplog)
				return oplog
			}
		}

		host.SyncOpCpu(destRes.ResCpu - opRep.ResCpu)
		host.SyncOpMem(destRes.ResMem - opRep.ResMem)
		hostChanged = true

		opRep.ResCpu = destRes.ResCpu
		opRep.ResMem = destRes.ResMem
		opRep.VolSys = destRes.VolSys

		// inapi.ObjPrint(podq.Meta.ID, opRep)
	}

	if inapi.OpActionAllow(opAction, inapi.OpActionStart) {
		if ports, chg := schedulePodRepNetPortAlloc(podq, opRep, host); chg {
			opRep.Ports = ports
			hostChanged = true
		}
	}

	if hostChanged {

		hlog.Printf("info", "host %s sync changes", host.Meta.Id)

		host.OpPortSort()

		if rs := data.ZoneMaster.PvPut(
			inapi.NsZoneSysHost(status.ZoneId, host.Meta.Id), host, nil,
		); !rs.OK() {
			hlog.Printf("error", "host %s sync changes failed %s", host.Meta.Id, rs.Bytex().String())
			return inapi.NewPbOpLogEntry("", inapi.PbOpLogWarn,
				fmt.Sprintf("host %s sync changes failed %s", host.Meta.Id, rs.Bytex().String()))
		}

		// TODO
	}

	// inapi.ObjPrint(podq.Meta.ID, podq)

	return inapi.NewPbOpLogEntry("", inapi.PbOpLogOK, "sync to host/"+opRep.Node)
}

func schedulePodMigrate(podq *inapi.Pod) error {

	if len(podq.Operate.RepMigrates) < 1 ||
		inapi.OpActionAllow(podq.Operate.Action, inapi.OpActionDestroy) {
		return nil
	}

	podStatus := status.ZonePodStatusList.Get(podq.Meta.ID)
	if podStatus == nil {
		return nil
	}

	var (
		tn  = uint32(time.Now().Unix())
		mgs = types.ArrayUint32(podq.Operate.RepMigrates)
	)

	for _, repId := range podq.Operate.RepMigrates {

		if repId >= uint32(podq.Operate.ReplicaCap) {
			mgs.Del(repId)
			continue
		}

		hlog.Printf("debug", "sche rep %s:%d",
			podq.Meta.ID, repId,
		)

		ctrRep := podq.Operate.Replicas.Get(repId)
		if ctrRep == nil {
			continue
		}
		if !inapi.OpActionAllow(ctrRep.Action, inapi.OpActionMigrate) {
			continue
		}

		repStatus := podStatus.RepGet(repId)
		if repStatus == nil {
			continue
		}

		if repStatus.OpLog == nil {
			repStatus.OpLog = inapi.NewPbOpLogSets(
				inapi.NsZonePodOpRepKey(repStatus.PodId, repStatus.RepId), 0)
		}

		hlog.Printf("debug", "rep %s:%d, op Action %s, repStatus Action %s",
			podq.Meta.ID, repId,
			strings.Join(inapi.OpActionStrings(ctrRep.Action), "|"),
			strings.Join(inapi.OpActionStrings(repStatus.Action), "|"),
		)

		if ctrRep.Next != nil &&
			inapi.OpActionAllow(ctrRep.Next.Action, inapi.OpActionMigrated) {

			if inapi.OpActionAllow(ctrRep.Action, inapi.OpActionDestroy) {

				if inapi.OpActionAllow(ctrRep.Action, inapi.OpActionDestroyed) {

					ctrRep.PrevNode = ctrRep.Node
					ctrRep.Node = ctrRep.Next.Node
					ctrRep.Action = inapi.OpActionStart
					ctrRep.Next = nil
					ctrRep.Options = nil
					mgs.Del(repId)

					repStatus.Action = 0

					// TODO
					data.ZoneMaster.PvPut(inapi.NsZonePodInstance(status.ZoneId, podq.Meta.ID), podq, nil)

					hlog.Printf("warn", "scheduler rep %s:%d, host %s, set opAction to %s, Migrate DONE",
						podq.Meta.ID, repId, ctrRep.Node,
						strings.Join(inapi.OpActionStrings(ctrRep.Action), "|"),
					)

					repStatus.OpLog.LogSet(
						podq.Operate.Version,
						inapi.NsOpLogZoneRepMigratePrevDestory, inapi.PbOpLogOK,
						fmt.Sprintf("prev box on host %s cleaned up", ctrRep.Node),
					)

					repStatus.OpLog.LogSet(
						podq.Operate.Version,
						inapi.NsOpLogZoneRepMigrateDone, inapi.PbOpLogOK,
						"migrate well done",
					)
				}

			} else {

				ctrRep.Action = inapi.OpActionMigrate | inapi.OpActionDestroy

				repStatus.OpLog.LogSet(
					podq.Operate.Version,
					inapi.NsOpLogZoneRepMigratePrevDestory, inapi.PbOpLogInfo,
					fmt.Sprintf("cleaning the prev box on host %s", ctrRep.Node),
				)

				hlog.Printf("warn", "scheduler rep %s:%d, host %s, set opAction to %s",
					podq.Meta.ID, repId, ctrRep.Node,
					strings.Join(inapi.OpActionStrings(ctrRep.Action), "|"),
				)
			}
		}
	}

	podq.Operate.RepMigrates = mgs

	// schedule new host
	for _, repId := range podq.Operate.RepMigrates {

		if repId >= uint32(podq.Operate.ReplicaCap) {
			continue
		}

		rep := podq.Operate.Replicas.Get(repId)
		if rep == nil {
			continue
		}

		repStatus := podStatus.RepGet(repId)
		if repStatus == nil {
			continue
		}

		if repStatus.OpLog == nil {
			repStatus.OpLog = inapi.NewPbOpLogSets(
				inapi.NsZonePodOpRepKey(repStatus.PodId, repStatus.RepId), 0)
		}

		// hlog.Printf("info", "pod %s", podq.Meta.ID)

		if rep.Next != nil {
			continue
		}
		// hlog.Printf("info", "pod %s", podq.Meta.ID)

		prevHostId := rep.Node
		prevHost := status.ZoneHostList.Item(rep.Node)
		if prevHost == nil {
			continue
		}
		// hlog.Printf("info", "pod %s", podq.Meta.ID)

		if podq.Spec.VolSys == nil {
			continue
		}

		repNext := &inapi.PodOperateReplica{
			RepId:  rep.RepId,
			Action: inapi.OpActionMigrate,
			ResCpu: rep.ResCpu,
			ResMem: rep.ResMem,
			VolSys: rep.VolSys,
		}

		hit, err := Scheduler.ScheduleHost(
			&typeScheduler.SchedulePodSpec{
				CellId:    podq.Spec.Cell,
				BoxDriver: podq.Spec.Box.Image.Driver,
			},
			&typeScheduler.SchedulePodReplica{
				RepId:       repNext.RepId,
				Cpu:         repNext.ResCpu,
				Mem:         repNext.ResMem,
				VolSys:      repNext.VolSys,
				VolSysAttrs: podq.Spec.VolSys.Attrs,
			},
			&scheduleHostList,
			&typeScheduler.ScheduleOptions{
				HostExcludes: []string{rep.Node},
			},
		)

		if err != nil || hit.Host == nil || len(hit.Volumes) < 1 {
			repStatus.OpLog.LogSet(
				podq.Operate.Version,
				inapi.NsOpLogZoneRepMigrateAlloc, inapi.PbOpLogWarn,
				fmt.Sprintf("migrate rep %s:%d, waiting for available resources",
					podq.Meta.ID, repId),
			)
			continue
		}

		host := status.ZoneHostList.Item(hit.HostId)
		if host == nil {
			continue
		}
		repNext.Node = hit.HostId
		repNext.Scheduled = uint32(time.Now().Unix())

		repNext.VolSysMnt = hit.Volumes[0].Name

		if ports, chg := schedulePodRepNetPortAlloc(podq, repNext, host); chg {
			repNext.Ports = ports
			host.OpPortSort()
		}

		hostPeerLan := inapi.HostNodeAddress(prevHost.Spec.PeerLanAddr)

		rep.Options.Set("rsync/host", fmt.Sprintf("%s:%d", hostPeerLan.IP(), hostPeerLan.Port()+5))
		rep.Options.Set("rsync/auth", idhash.RandHexString(30))

		rep.PrevNode = ""
		rep.Next = repNext
		rep.Action = inapi.OpActionStop | inapi.OpActionMigrate
		rep.Updated = tn

		host.SyncOpCpu(repNext.ResCpu)
		host.SyncOpMem(repNext.ResMem)

		hit.Host.CpuUsed += repNext.ResCpu
		hit.Host.MemUsed += repNext.ResMem

		for _, hitVol := range hit.Volumes {
			if pv := inapi.ResVolValueSliceGet(host.Operate.VolUsed, hitVol.Name); pv != nil {
				pv.Value += hitVol.Size
			}
			if pv := hit.Host.Volumes.Get(hitVol.Name); pv != nil {
				pv.Used += hitVol.Size
			}
		}

		// hlog.Printf("info", "host %s sync changes", host.Meta.Id)

		if rs := data.ZoneMaster.PvPut(
			inapi.NsZoneSysHost(status.ZoneId, host.Meta.Id), host, nil,
		); !rs.OK() {
			hlog.Printf("error", "host %s sync changes failed %s", host.Meta.Id, rs.Bytex().String())
			continue
		}

		hlog.Printf("warn", "scheduler rep %s:%d, migrate from host %s to %s",
			podq.Meta.ID, rep.RepId, prevHostId, hit.HostId)

		if rs := data.ZoneMaster.PvPut(inapi.NsZonePodInstance(status.ZoneId, podq.Meta.ID), podq, nil); !rs.OK() {
			hlog.Printf("error", "zone/podq saved %s, err (%s)", podq.Meta.ID, rs.Bytex().String())
		}

		repStatus.OpLog.LogSet(
			podq.Operate.Version,
			inapi.NsOpLogZoneRepMigrateAlloc, inapi.PbOpLogOK,
			fmt.Sprintf("migrate rep from host %s to %s", prevHostId, hit.HostId),
		)
	}

	return nil
}

func schedulePodRepNetPortAlloc(
	podq *inapi.Pod,
	opRep *inapi.PodOperateReplica,
	host *inapi.ResHost,
) (inapi.ServicePorts, bool) {

	var (
		hostPeerLan  = inapi.HostNodeAddress(host.Spec.PeerLanAddr)
		hostPeerPort = hostPeerLan.Port()
		ports        = podq.AppServicePorts()
		hostChanged  = false
	)

	for i, pv := range ports {

		if pv.HostPort > 0 {

			if pv.HostPort == hostPeerPort {

				ports[i].HostPort = 0

				hlog.Printf("warn", "the host port %s:%d already in use ",
					host.Spec.PeerLanAddr, pv.HostPort)

			} else if host.OpPortHas(pv.HostPort) {

				if ppv := opRep.Ports.Get(pv.BoxPort); ppv == nil ||
					ppv.HostPort != pv.HostPort {

					ports[i].HostPort = 0

					hlog.Printf("warn", "the host port %s:%d is already allocated",
						host.Spec.PeerLanAddr, pv.HostPort)
				}
			} else {
				host.OpPortAlloc(pv.HostPort)
				hostChanged = true
			}

			continue
		}

		for _, pvp := range opRep.Ports {

			if pvp.BoxPort != pv.BoxPort {
				continue
			}

			if pvp.HostPort > 0 {
				ports.Sync(*pvp)
			}

			break
		}
	}

	// clean unused Ports
	for _, p := range opRep.Ports {
		if ports.Get(p.BoxPort) == nil {
			opRep.Ports.Del(p.BoxPort)
		}
	}

	// Assign new Host:Ports
	portsAlloc := []uint16{}
	for i, p := range ports {

		if p.HostPort > 0 {
			continue
		}

		if portAlloc := host.OpPortAlloc(0); portAlloc > 0 {

			ports[i].HostPort = portAlloc
			hostChanged = true

			portsAlloc = append(portsAlloc, portAlloc)

			hlog.Printf("info", "zm new port alloc to rep %s:%d, host %s, port %d",
				podq.Meta.ID, opRep.RepId, host.Meta.Id, portAlloc)

		} else {
			hlog.Printf("warn", "zm host %s res-port out range", host.Meta.Id)
		}
	}

	return ports, hostChanged
}

func schedulePodPreChargeValid(podq *inapi.Pod) error {

	specPlan := zonePodSpecPlans.Get(podq.Spec.Ref.Id)
	if specPlan == nil {
		return fmt.Errorf("bad pod.Spec %s", podq.Meta.ID)
	}

	if podq.Spec.VolSys == nil {
		return fmt.Errorf("Pod %s : No PodSpec/VolSys Setup", podq.Meta.ID)
	}

	// Volumes
	chargeAmount := iamapi.AccountFloat64Round(
		specPlan.VolCharge(podq.Spec.VolSys.RefId)*float64(podq.Spec.VolSys.Size), 4)

	if podq.Spec.Box.Resources != nil {
		// CPU
		chargeAmount += iamapi.AccountFloat64Round(
			specPlan.ResComputeCharge.Cpu*(float64(podq.Spec.Box.Resources.CpuLimit)/10), 4)

		// RAM
		chargeAmount += iamapi.AccountFloat64Round(
			specPlan.ResComputeCharge.Mem*float64(podq.Spec.Box.Resources.MemLimit), 4)
	}

	chargeAmount = chargeAmount * float64(podq.Operate.ReplicaCap)

	chargeCycleMin := float64(3600)
	chargeAmount = iamapi.AccountFloat64Round(chargeAmount*(chargeCycleMin/3600), 2)
	if chargeAmount < 0.01 {
		chargeAmount = 0.01
	}

	tnu := uint32(time.Now().Unix())
	if rsp := iamclient.AccountChargePreValid(iamapi.AccountChargePrepay{
		User:      podq.Meta.User,
		Product:   types.NameIdentifier(fmt.Sprintf("pod/%s", podq.Meta.ID)),
		Prepay:    chargeAmount,
		TimeStart: tnu,
		TimeClose: tnu + uint32(chargeCycleMin),
	}, status.ZonePodChargeAccessKey()); rsp.Error != nil {

		status, msg := inapi.PbOpLogWarn, rsp.Error.Message
		if msg == "" {
			msg = "Error Code: " + rsp.Error.Code
		}
		if rsp.Error.Code == iamapi.ErrCodeAccChargeOut {
			status = inapi.PbOpLogError
		}
		podq.Operate.OpLog, _ = inapi.PbOpLogEntrySliceSync(podq.Operate.OpLog,
			inapi.NewPbOpLogEntry(inapi.OpLogNsZoneMasterPodScheduleCharge, status, msg),
		)
		return errors.New(msg)
	} else if rsp.Kind != "AccountCharge" {
		return errors.New("Network Error")
	}

	podq.Operate.OpLog, _ = inapi.PbOpLogEntrySliceSync(podq.Operate.OpLog,
		inapi.NewPbOpLogEntry(inapi.OpLogNsZoneMasterPodScheduleCharge, inapi.PbOpLogOK, "PreValid OK"),
	)

	return nil
}

func schedulePodFailover(podq *inapi.Pod) error {

	if !inapi.OpActionAllow(podq.Operate.Action, inapi.OpActionStart) {
		return nil
	}

	podStatus := status.ZonePodStatusList.Get(podq.Meta.ID)
	if podStatus == nil {
		return nil
	}

	var (
		ft, fnm = podq.FailoverActive()
		fs      = podStatus.HealthFails(ft, podq.IsStateful())
	)

	if len(fs) == 0 {
		return nil
	}

	if len(fs) > fnm {
		hlog.Printf("info", "zm/pod %s, failover active Fails %d, Max %d", len(fs), fnm)
		return nil
	}

	hlog.Printf("info", "zm/pod %s, failover active Time %d, Max %d, Fails %d",
		podq.Meta.ID, ft, fnm, len(fs))

	tn := uint32(time.Now().Unix())

	for _, repId := range fs {

		if repId >= uint32(podq.Operate.ReplicaCap) {
			continue
		}

		rep := podq.Operate.Replicas.Get(repId)
		if rep == nil {
			continue
		}

		repStatus := podStatus.RepGet(repId)
		if repStatus == nil {
			continue
		}

		if inapi.OpActionAllow(rep.Action, inapi.OpActionMigrate) {
			continue
		}

		if (rep.Scheduled + uint32(inapi.HealthFailoverScheduleTimeMin)) > tn {
			hlog.Printf("debug", "zm/pod %s:%d, failover skip schudele time %d",
				podq.Meta.ID, rep.RepId, rep.Scheduled)
			continue
		}

		//
		repNext := &inapi.PodOperateReplica{
			RepId:  rep.RepId,
			Action: inapi.OpActionMigrate,
			ResCpu: rep.ResCpu,
			ResMem: rep.ResMem,
			VolSys: rep.VolSys,
		}

		prevHostId := rep.Node

		hit, err := Scheduler.ScheduleHost(
			&typeScheduler.SchedulePodSpec{
				CellId:    podq.Spec.Cell,
				BoxDriver: podq.Spec.Box.Image.Driver,
			},
			&typeScheduler.SchedulePodReplica{
				RepId:       repNext.RepId,
				Cpu:         repNext.ResCpu,
				Mem:         repNext.ResMem,
				VolSys:      repNext.VolSys,
				VolSysAttrs: podq.Spec.VolSys.Attrs,
			},
			&scheduleHostList,
			&typeScheduler.ScheduleOptions{
				HostExcludes: []string{prevHostId},
			},
		)

		if err != nil || hit.Host == nil || len(hit.Volumes) < 1 {
			hlog.Printf("info", "zm/pod %s:%d, failover schudele unhit",
				podq.Meta.ID, rep.RepId)
			continue
		}

		host := status.ZoneHostList.Item(hit.HostId)
		if host == nil {
			continue
		}

		rep.Node = hit.HostId
		rep.Scheduled = tn
		rep.PrevNode = ""

		rep.VolSysMnt = hit.Volumes[0].Name

		hlog.Printf("info", "zm/pod %s:%d, failover schudele unhit",
			podq.Meta.ID, rep.RepId)

		if ports, chg := schedulePodRepNetPortAlloc(podq, repNext, host); chg {
			rep.Ports = ports
			host.OpPortSort()
		}

		rep.Action = inapi.OpActionStart
		rep.Updated = tn

		host.SyncOpCpu(repNext.ResCpu)
		host.SyncOpMem(repNext.ResMem)

		hit.Host.CpuUsed += rep.ResCpu
		hit.Host.MemUsed += rep.ResMem

		for _, hitVol := range hit.Volumes {
			if pv := inapi.ResVolValueSliceGet(host.Operate.VolUsed, hitVol.Name); pv != nil {
				pv.Value += hitVol.Size
			}
			if pv := hit.Host.Volumes.Get(hitVol.Name); pv != nil {
				pv.Used += hitVol.Size
			}
		}

		if rs := data.ZoneMaster.PvPut(
			inapi.NsZoneSysHost(status.ZoneId, host.Meta.Id), host, nil,
		); !rs.OK() {
			hlog.Printf("error", "host %s sync changes failed %s", host.Meta.Id, rs.Bytex().String())
			continue
		}

		hlog.Printf("warn", "failover rep %s:%d, move from host %s to %s",
			podq.Meta.ID, rep.RepId, prevHostId, hit.HostId)

		if rs := data.ZoneMaster.PvPut(inapi.NsZonePodInstance(status.ZoneId, podq.Meta.ID), podq, nil); !rs.OK() {
			hlog.Printf("error", "zone/podq saved %s, err (%s)", podq.Meta.ID, rs.Bytex().String())
		}

		repStatus.OpLog.LogSet(
			podq.Operate.Version,
			inapi.NsOpLogZoneRepMigrateAlloc, inapi.PbOpLogOK,
			fmt.Sprintf("failover rep from host %s to %s", prevHostId, hit.HostId),
		)
	}

	return nil
}
