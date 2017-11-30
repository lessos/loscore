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

package v1

import (
	"fmt"
	"strings"
	"time"

	"github.com/hooto/httpsrv"
	"github.com/hooto/iam/iamapi"
	"github.com/hooto/iam/iamclient"
	"github.com/lessos/lessgo/crypto/idhash"
	"github.com/lessos/lessgo/types"
	"github.com/lynkdb/iomix/skv"
	iox_utils "github.com/lynkdb/iomix/utils"

	"github.com/sysinner/incore/data"
	"github.com/sysinner/incore/inapi"
	zm_status "github.com/sysinner/incore/status"
)

type Pod struct {
	*httpsrv.Controller
	us iamapi.UserSession
}

func (c *Pod) Init() int {

	//
	c.us, _ = iamclient.SessionInstance(c.Session)

	if !c.us.IsLogin() {
		c.Response.Out.WriteHeader(401)
		c.RenderJson(types.NewTypeErrorMeta(iamapi.ErrCodeUnauthorized, "Unauthorized"))
		return 1
	}

	return 0
}

func (c Pod) ListAction() {

	var (
		ls inapi.PodList
	)

	defer c.RenderJson(&ls)

	// TODO pager
	var rs *skv.Result
	if zone_id := c.Params.Get("zone_id"); zone_id != "" {
		rs = data.ZoneMaster.PvRevScan(inapi.NsZonePodInstance(zone_id, ""), "", "", 1000)
	} else {
		rs = data.ZoneMaster.PvRevScan(inapi.NsGlobalPodInstance(""), "", "", 1000)
	}
	rss := rs.KvList()

	var fields types.ArrayPathTree
	if fns := c.Params.Get("fields"); fns != "" {
		fields.Set(fns)
		fields.Sort()
	}

	var exp_app_filter_notin types.ArrayString
	if v := c.Params.Get("exp_app_filter_notin"); v != "" {
		exp_app_filter_notin = types.ArrayString(strings.Split(v, ","))
	}

	action := uint32(c.Params.Uint64("operate_action"))

	for _, v := range rss {

		var pod inapi.Pod

		if err := v.Decode(&pod); err == nil {

			// TOPO
			if pod.Meta.User != c.us.UserName {
				continue
			}

			if c.Params.Int64("destroy_enable") != 1 &&
				inapi.OpActionAllow(pod.Operate.Action, inapi.OpActionDestroy) {
				continue
			}

			if action > 0 && !inapi.OpActionAllow(pod.Operate.Action, action) {
				continue
			}

			if len(exp_app_filter_notin) > 0 {
				found := false
				for _, vpa := range pod.Apps {
					if exp_app_filter_notin.Has(vpa.Spec.Meta.ID) {
						found = true
						break
					}
				}
				if found {
					continue
				}
			}

			if len(fields) > 0 {

				podfs := inapi.Pod{
					Meta: types.InnerObjectMeta{
						ID: pod.Meta.ID,
					},
				}

				if fields.Has("meta/name") {
					podfs.Meta.Name = pod.Meta.Name
				}

				if fields.Has("meta/updated") {
					podfs.Meta.Updated = pod.Meta.Updated
				}

				if fields.Has("spec") && pod.Spec != nil {
					podfs.Spec = &inapi.PodSpecBound{}

					if fields.Has("spec/ref/id") {
						podfs.Spec.Ref.Id = pod.Spec.Ref.Id
					}

					if fields.Has("spec/ref/name") {
						podfs.Spec.Ref.Name = pod.Spec.Ref.Name
					}

					if fields.Has("spec/zone") {
						podfs.Spec.Zone = pod.Spec.Zone
					}

					if fields.Has("spec/cell") {
						podfs.Spec.Cell = pod.Spec.Cell
					}
				}

				if fields.Has("apps") {

					for _, a := range pod.Apps {

						afs := inapi.AppInstance{}

						if fields.Has("apps/meta/id") {
							afs.Meta.ID = a.Meta.ID
						}

						podfs.Apps = append(podfs.Apps, afs)
					}
				}

				if fields.Has("operate") {

					if fields.Has("operate/action") {
						podfs.Operate.Action = pod.Operate.Action | zm_status.ZonePodRepMergeOperateAction(pod.Meta.ID, pod.Operate.ReplicaCap)
					}

					if fields.Has("operate/version") {
						podfs.Operate.Version = pod.Operate.Version
					}

					if fields.Has("operate/replica_cap") {
						podfs.Operate.ReplicaCap = pod.Operate.ReplicaCap
					}

					if fields.Has("operate/replicas") {
						podfs.Operate.Replicas = pod.Operate.Replicas
					}
				}

				ls.Items = append(ls.Items, podfs)
			} else {
				ls.Items = append(ls.Items, pod)
			}
		}
	}

	ls.Kind = "PodList"
}

func (c Pod) EntryAction() {

	var (
		id      = c.Params.Get("id")
		fields  = c.Params.Get("fields")
		zone_id = c.Params.Get("zone_id")
		set     inapi.Pod
	)

	defer c.RenderJson(&set)

	if zone_id == "" {
		if rs := data.ZoneMaster.PvGet(inapi.NsGlobalPodInstance(id)); rs.OK() {
			rs.Decode(&set)
			if set.Meta.ID == "" || set.Meta.User != c.us.UserName {
				set = inapi.Pod{}
				set.Error = types.NewErrorMeta("404", "Pod Not Found")
				return
			}
			zone_id = set.Spec.Zone
		}
	}

	if zone_id != "" {

		if rs := data.ZoneMaster.PvGet(inapi.NsZonePodInstance(zone_id, id)); rs.OK() {

			rs.Decode(&set)
			if set.Meta.ID == "" || set.Meta.User != c.us.UserName {
				set = inapi.Pod{}
				set.Error = types.NewErrorMeta("404", "Pod Not Found")
				return
			}
			if fields == "status" {
				if v := pod_status(set.Meta.ID, c.us.UserName); v.Action != 0 {
					set.Status = &v
				} else {
					set.Status = nil
				}
			}

			set.Operate.Action = inapi.OpActionAppend(set.Operate.Action,
				zm_status.ZonePodRepMergeOperateAction(set.Meta.ID, set.Operate.ReplicaCap))
		}
	}

	for _, v := range set.Operate.Replicas {
		if host := zm_status.ZoneHostList.Item(v.Node); host != nil {
			for _, v2 := range v.Ports {
				if i := strings.IndexByte(host.Spec.PeerLanAddr, ':'); i > 0 {
					v2.LanAddr = host.Spec.PeerLanAddr[:i]
				} else {
					v2.LanAddr = host.Spec.PeerLanAddr
				}
				v2.WanAddr = host.Spec.PeerWanAddr
			}
		}
	}

	set.Kind = "Pod"
}

func (c Pod) NewAction() {

	var (
		set       inapi.PodCreate
		spec_plan inapi.PodSpecPlan
	)

	defer c.RenderJson(&set)

	if err := c.Request.JsonDecode(&set); err != nil {
		set.Error = types.NewErrorMeta("400", err.Error())
		return
	}

	set.Name = strings.TrimSpace(set.Name)
	if set.Name == "" {
		set.Error = types.NewErrorMeta("400", "Name can not be null")
		return
	}

	//
	if rs := data.ZoneMaster.PvGet(inapi.NsGlobalPodSpec("plan", set.Plan)); rs.OK() {
		rs.Decode(&spec_plan)
	}
	if spec_plan.Meta.ID == "" || spec_plan.Meta.ID != set.Plan {
		set.Error = types.NewErrorMeta("400", "Spec Not Found")
		return
	}

	//
	if err := set.Valid(spec_plan); err != nil {
		set.Error = types.NewErrorMeta("400", err.Error())
		return
	}
	spec_plan.ChargeFix()

	//
	var zone inapi.ResZone
	if rs := data.ZoneMaster.PvGet(inapi.NsGlobalSysZone(set.Zone)); rs.OK() {
		rs.Decode(&zone)
	}
	if zone.Meta.Id == "" {
		set.Error = types.NewErrorMeta("400", "Zone Not Found")
		return
	}

	//
	var cell inapi.ResCell
	if rs := data.ZoneMaster.PvGet(inapi.NsGlobalSysCell(set.Zone, set.Cell)); rs.OK() {
		rs.Decode(&cell)
	}
	if cell.Meta.Id == "" {
		set.Error = types.NewErrorMeta("400", "Cell Not Found")
		return
	}

	res_vol := spec_plan.ResVolume(set.ResVolume)
	if res_vol == nil {
		set.Error = types.NewErrorMeta("400", "No ResVolume Found")
		return
	}
	if set.ResVolumeSize < res_vol.Request {
		set.ResVolumeSize = res_vol.Request
	} else if set.ResVolumeSize > res_vol.Limit {
		set.ResVolumeSize = res_vol.Limit
	} else {
		if fix := set.ResVolumeSize % res_vol.Step; fix > 0 {
			set.ResVolumeSize += res_vol.Step
		}
	}

	tn := types.MetaTimeNow()
	id := iox_utils.Uint32ToHexString(uint32(tn.Time().Unix())) + idhash.RandHexString(8)

	pod := inapi.Pod{
		Meta: types.InnerObjectMeta{
			ID:      id,
			Name:    set.Name,
			User:    c.us.UserName,
			Created: tn,
			Updated: tn,
		},
		Spec: &inapi.PodSpecBound{
			Ref: inapi.ObjectReference{
				Id:      spec_plan.Meta.ID,
				Name:    spec_plan.Meta.Name,
				Version: spec_plan.Meta.Version,
			},
			Zone:   set.Zone,
			Cell:   set.Cell,
			Labels: spec_plan.Labels,
			Volumes: []inapi.PodSpecResVolumeBound{
				{
					Ref: inapi.ObjectReference{
						Id:   res_vol.RefId,
						Name: res_vol.RefId,
					},
					Name:      "system",
					SizeLimit: set.ResVolumeSize,
				},
			},
		},
		Operate: inapi.PodOperate{
			Action:     inapi.OpActionStart,
			Version:    1,
			ReplicaCap: 1, // TODO
			Operated:   uint32(time.Now().Unix()),
		},
	}

	//
	for _, v := range set.Boxes {

		img := spec_plan.Image(v.Image)
		if img == nil {
			set.Error = types.NewErrorMeta("400", "No Image Found")
			return
		}

		res := spec_plan.ResCompute(v.ResCompute)
		if res == nil {
			set.Error = types.NewErrorMeta("400", "No ResCompute Found")
			return
		}

		pod.Spec.Boxes = append(pod.Spec.Boxes, inapi.PodSpecBoxBound{
			Name:    v.Name,
			Updated: types.MetaTimeNow(),
			Image: inapi.PodSpecBoxImageBound{
				Ref: &inapi.ObjectReference{
					Id:   img.RefId,
					Name: img.RefId,
				},
				Driver:  img.Driver,
				OsDist:  img.OsDist,
				Arch:    img.Arch,
				Options: img.Options,
			},
			Resources: &inapi.PodSpecBoxResComputeBound{
				Ref: &inapi.ObjectReference{
					Id:   res.RefId,
					Name: res.RefId,
				},
				CpuLimit: res.CpuLimit,
				MemLimit: res.MemLimit,
			},
		})
	}

	charge_amount := float64(0)

	// Volumes
	for _, v := range pod.Spec.Volumes {
		charge_amount += iamapi.AccountFloat64Round(
			spec_plan.ResVolumeCharge.CapSize*float64(v.SizeLimit/inapi.ByteMB), 4)
	}

	for _, v := range pod.Spec.Boxes {

		if v.Resources != nil {
			// CPU
			charge_amount += iamapi.AccountFloat64Round(
				spec_plan.ResComputeCharge.Cpu*(float64(v.Resources.CpuLimit)/1000), 4)

			// RAM
			charge_amount += iamapi.AccountFloat64Round(
				spec_plan.ResComputeCharge.Mem*float64(v.Resources.MemLimit/inapi.ByteMB), 4)
		}
	}

	charge_amount = charge_amount * float64(pod.Operate.ReplicaCap)

	charge_cycle_min := float64(3600)
	charge_amount = iamapi.AccountFloat64Round(charge_amount*(charge_cycle_min/3600), 2)
	if charge_amount < 0.01 {
		charge_amount = 0.01
	}

	tnu := uint32(time.Now().Unix())
	if rsp := iamclient.AccountChargePreValid(iamapi.AccountChargePrepay{
		User:      pod.Meta.User,
		Product:   types.NameIdentifier(fmt.Sprintf("pod/%s", pod.Meta.ID)),
		Prepay:    charge_amount,
		TimeStart: tnu,
		TimeClose: tnu + uint32(charge_cycle_min),
	}, zm_status.ZonePodChargeAccessKey()); rsp.Error != nil {
		set.Error = rsp.Error
		return
	} else if rsp.Kind != "AccountCharge" {
		set.Error = types.NewErrorMeta("400", "Network Error")
		return
	}

	//
	if rs := data.ZoneMaster.PvNew(inapi.NsGlobalPodInstance(pod.Meta.ID), pod, nil); !rs.OK() {
		set.Error = types.NewErrorMeta("500", rs.Bytex().String())
		return
	}

	/*
		oplog := inapi.NewOpLog(c.us.UserName)
		if rs := data.ZoneMaster.PvNew(inapi.NsZonePodOperateLog(pod.Spec.Zone, pod.RepId(0)), oplog, nil); !rs.OK() {
			set.Error = types.NewErrorMeta("500", rs.Bytex().String())
			return
		}
	*/

	// Pod Map to Cell Queue
	qmpath := inapi.NsZonePodOpQueue(pod.Spec.Zone, pod.Spec.Cell, pod.Meta.ID)
	data.ZoneMaster.PvNew(qmpath, pod, nil)

	set.Pod = pod.Meta.ID
	set.Kind = "PodInstance"
}

func (c Pod) OpActionSetAction() {

	set := types.TypeMeta{}
	defer c.RenderJson(&set)

	var (
		pod_id    = c.Params.Get("pod_id")
		op_action = uint32(c.Params.Uint64("op_action"))
	)

	if !inapi.PodIdReg.MatchString(pod_id) {
		set.Error = types.NewErrorMeta("400", "Invalid Pod ID")
		return
	}

	if !inapi.OpActionValid(op_action) {
		set.Error = types.NewErrorMeta("400", "Invalid OpAction")
		return
	}

	var prev inapi.Pod
	if rs := data.ZoneMaster.PvGet(inapi.NsGlobalPodInstance(pod_id)); !rs.OK() {
		set.Error = types.NewErrorMeta("400", "Pod Not Found")
		return
	} else {
		rs.Decode(&prev)
	}

	if prev.Meta.ID != pod_id ||
		prev.Meta.User != c.us.UserName {
		set.Error = types.NewErrorMeta("400", "Pod Not Found or Access Denied")
		return
	}

	if inapi.OpActionAllow(prev.Operate.Action, inapi.OpActionDestroy) {
		set.Error = types.NewErrorMeta("400", "the pod instance has been destroyed")
		return
	}

	//
	prev.Operate.Action = op_action
	prev.Operate.Version++
	prev.Meta.Updated = types.MetaTimeNow()

	tn := uint32(time.Now().Unix())

	// Pod Map to Cell Queue
	qstr := inapi.NsZonePodOpQueue(prev.Spec.Zone, prev.Spec.Cell, prev.Meta.ID)
	if rs := data.ZoneMaster.PvGet(qstr); rs.OK() {
		if tn < prev.Operate.Operated+600 {
			set.Error = types.NewErrorMeta(inapi.ErrCodeBadArgument, "the previous operation is in processing, please try again later")
		}
		return
	}

	prev.Operate.Operated = tn
	if rs := data.ZoneMaster.PvPut(qstr, prev, nil); !rs.OK() {
		set.Error = types.NewErrorMeta(inapi.ErrCodeServerError, "server error, please try again later")
		return
	}
	data.ZoneMaster.PvPut(inapi.NsGlobalPodInstance(prev.Meta.ID), prev, nil)

	set.Kind = "PodInstance"
}

func (c Pod) SetInfoAction() {

	var (
		set  inapi.Pod
		prev inapi.Pod
	)

	defer c.RenderJson(&set)

	if err := c.Request.JsonDecode(&set); err != nil {
		set.Error = types.NewErrorMeta("400", err.Error())
		return
	}

	if set.Operate.Action > 0 && !inapi.OpActionValid(set.Operate.Action) {
		set.Error = types.NewErrorMeta("400", "Invalid OpAction")
		return
	}

	if rs := data.ZoneMaster.PvGet(inapi.NsGlobalPodInstance(set.Meta.ID)); !rs.OK() {
		set.Error = types.NewErrorMeta("400", "Prev Pod Not Found")
		return
	} else {
		rs.Decode(&prev)
	}

	if prev.Meta.ID != set.Meta.ID {
		set.Error = types.NewErrorMeta("400", "Prev Pod Not Found")
		return
	}

	if prev.Meta.User != c.us.UserName {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeAccessDenied, "Access Denied")
		return
	}

	if inapi.OpActionAllow(prev.Operate.Action, inapi.OpActionDestroy) {
		set.Error = types.NewErrorMeta("400", "the pod instance has been destroyed")
		return
	}

	//
	if prev.Meta.Name == set.Meta.Name &&
		prev.Operate.Action == set.Operate.Action {
		set.Kind = "PodInstance"
		return
	}

	prev.Meta.Name = set.Meta.Name
	if set.Operate.Action > 0 {
		prev.Operate.Action = set.Operate.Action
		prev.Operate.Version++
	}

	//
	prev.Meta.Updated = types.MetaTimeNow()

	tn := uint32(time.Now().Unix())

	// Pod Map to Cell Queue
	qstr := inapi.NsZonePodOpQueue(prev.Spec.Zone, prev.Spec.Cell, prev.Meta.ID)
	if rs := data.ZoneMaster.PvGet(qstr); rs.OK() {
		if tn < prev.Operate.Operated+600 {
			set.Error = types.NewErrorMeta(inapi.ErrCodeBadArgument, "the previous operation is in processing, please try again later")
		}
		return
	}

	prev.Operate.Operated = tn
	if rs := data.ZoneMaster.PvPut(qstr, prev, nil); !rs.OK() {
		set.Error = types.NewErrorMeta(inapi.ErrCodeServerError, "server error, please try again later")
		return
	}
	data.ZoneMaster.PvPut(inapi.NsGlobalPodInstance(prev.Meta.ID), prev, nil)

	set.Kind = "PodInstance"
}

func (c Pod) StatusAction() {

	rsp := inapi.PodStatus{}

	defer c.RenderJson(&rsp)

	rsp = pod_status(c.Params.Get("id"), c.us.UserName)
}

func pod_status(pod_id string, user_name string) inapi.PodStatus {

	var (
		pod    inapi.Pod
		status inapi.PodStatus
	)

	//
	if rs := data.ZoneMaster.PvGet(inapi.NsGlobalPodInstance(pod_id)); rs.OK() {
		rs.Decode(&pod)
	}

	if pod.Meta.ID == "" {
		status.Error = types.NewErrorMeta("404.01", "Pod Not Found")
		return status
	}

	if user_name != pod.Meta.User {
		status.Error = types.NewErrorMeta(iamapi.ErrCodeAccessDenied, "Access Denied")
		return status
	}

	//
	if rs := data.ZoneMaster.PvGet(inapi.NsZonePodInstance(pod.Spec.Zone, pod_id)); rs.OK() {
		rs.Decode(&pod)
	}

	if pod.Meta.ID == "" {
		status.Error = types.NewErrorMeta("404.02", "Pod Not Found")
		return status
	}

	if len(pod.Operate.OpLog) > 0 {
		status.OpLog = pod.Operate.OpLog
	}

	for _, rep := range pod.Operate.Replicas {

		if rep.Node == "" {
			continue
		}

		rep_status := inapi.PbPodRepStatusSliceGet(zm_status.ZonePodRepStatusSets, pod.Meta.ID, uint32(rep.Id))
		if rep_status == nil {
			continue
		}

		if rep_status.Action == 0 {
			rep_status.Action = inapi.OpActionPending
		} else {

			if rep_status.Action == inapi.OpActionRunning &&
				(uint32(time.Now().UTC().Unix())-rep_status.Updated) > 600 {
				rep_status.Action = 0
			}

			status.Replicas = append(status.Replicas, rep_status)
		}
	}

	status.Kind = "PodStatus"

	status.Refresh(pod.Operate.ReplicaCap)

	return status
}
