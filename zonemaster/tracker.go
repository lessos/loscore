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
	"strings"
	"time"

	"github.com/hooto/hlog4g/hlog"
	"github.com/lynkdb/iomix/skv"

	"github.com/sysinner/incore/data"
	"github.com/sysinner/incore/inapi"
	"github.com/sysinner/incore/status"

	iam_api "github.com/hooto/iam/iamapi"
	iam_db "github.com/hooto/iam/store"
)

var (
	zm_leader_refreshed int64 = 0
	zm_pod_statuses           = map[string]uint32{}
)

func zone_tracker() {

	var (
		force_refresh = false
	)

	if status.Host.Operate == nil {
		hlog.Printf("info", "status.Host.Operate")
		return
	}

	// is zone-master
	if !status.IsZoneMaster() {
		hlog.Printf("debug", "status.IsZoneMaster SKIP")
		return
	}

	if status.ZoneId == "" {
		hlog.Printf("error", "config.json host/zone_id Not Found")
		return
	}

	// if leader active
	leader_path := inapi.NsZoneSysMasterLeader(status.Host.Operate.ZoneId)
	if rs := data.ZoneMaster.PvGet(leader_path); rs.NotFound() {

		if !inapi.ResSysHostIdReg.MatchString(status.Host.Meta.Id) {
			return
		}

		if rs2 := data.ZoneMaster.PvNew(
			leader_path,
			status.Host.Meta.Id,
			&skv.ProgWriteOptions{
				Expired: uint64(time.Now().Add(12e9).UnixNano()),
			},
		); rs2.OK() {
			status.ZoneMasterList.Leader = status.Host.Meta.Id
			force_refresh = true
			hlog.Printf("warn", "new zone-master/leader %s", status.Host.Meta.Id)
		} else {
			hlog.Printf("info", "status.Host.Operate")
			return
		}

	} else if rs.OK() {
		node_id := rs.String()
		if inapi.ResSysHostIdReg.MatchString(node_id) &&
			status.ZoneMasterList.Leader != node_id {
			status.ZoneMasterList.Leader = node_id
			force_refresh = true
		}
	} else {
		hlog.Printf("warn", "refresh zone-master leader active failed")
		return
	}

	// refresh host keys
	if len(status.ZoneHostSecretKeys) == 0 {
		if rs := data.ZoneMaster.PvScan(
			inapi.NsZoneSysHostSecretKey(status.Host.Operate.ZoneId, ""), "", "", 1000); rs.OK() {

			rs.KvEach(func(v *skv.ResultEntry) int {
				status.ZoneHostSecretKeys.Set(string(v.Key), v.Bytex().String())
				return 0
			})

		} else {
			hlog.Printf("warn", "refresh host key list failed")
		}
	}

	// is zone-master leader
	if !status.IsZoneMasterLeader() {
		return
	}

	// refresh zone-master leader ttl
	// pv := skv.NewValueObject(status.Host.Meta.Id)
	if rs := data.ZoneMaster.PvPut(
		leader_path,
		status.Host.Meta.Id,
		&skv.ProgWriteOptions{
			// PrevSum: pv.Crc32(), // TODO BUG
			Expired: uint64(time.Now().Add(12e9).UnixNano()),
		},
	); !rs.OK() {
		hlog.Printf("warn", "refresh zone-master leader ttl failed "+rs.String())
		return
	}

	// refresh pod's statuses
	for _, v := range status.ZonePodList {

		if v.Spec == nil || v.Spec.Zone != status.Host.Operate.ZoneId {
			continue
		}

		if v.Operate.ReplicaCap < 1 {
			continue
		}

		action := status.ZonePodRepMergeOperateAction(v.Meta.ID, v.Operate.ReplicaCap)
		if action < 1 {
			continue
		}

		if v2, ok := zm_pod_statuses[v.Meta.ID]; ok && v2 == action {
			continue
		}

		pst := inapi.PodStatus{
			Action: action,
		}
		rs := data.ZoneMaster.PvPut(inapi.NsZonePodStatus(status.Host.Operate.ZoneId, v.Meta.ID), pst, nil)
		if !rs.OK() {
			continue
		}
		rs = data.GlobalMaster.PvPut(inapi.NsGlobalPodStatus(status.Host.Operate.ZoneId, v.Meta.ID), pst, nil)
		if !rs.OK() {
			continue
		}
		zm_pod_statuses[v.Meta.ID] = action
	}

	if !force_refresh &&
		time.Now().UTC().Unix()-zm_leader_refreshed < 60 {
		hlog.Printf("debug", "zone-master/refresh SKIP")
		return
	}
	zm_leader_refreshed = time.Now().UTC().Unix()

	//
	if status.Zone == nil {

		if rs := data.ZoneMaster.PvGet(inapi.NsZoneSysInfo(status.Host.Operate.ZoneId)); rs.OK() {

			var zone inapi.ResZone
			if err := rs.Decode(&zone); err != nil {
				hlog.Printf("error", "No ZoneInfo Setup in db")
				return
			}

			if zone.Meta != nil && zone.Meta.Id == status.Host.Operate.ZoneId {
				status.Zone = &zone
			}
		}

		if status.Zone == nil {
			hlog.Printf("error", "No ZoneInfo Setup in status")
			return
		}

		init_if := iam_api.AccessKey{
			User: "sysadmin",
			Bounds: []iam_api.AccessKeyBound{{
				Name: "sys/zm/" + status.Host.Operate.ZoneId,
			}},
			Description: "ZoneMaster AccCharge",
		}
		if v, ok := status.Zone.OptionGet("iam/acc_charge/access_key"); ok {
			init_if.AccessKey = v
		}
		if v, ok := status.Zone.OptionGet("iam/acc_charge/secret_key"); ok {
			init_if.SecretKey = v
		}
		if init_if.AccessKey != "" {
			iam_db.AccessKeyInitData(init_if)
		}
	}

	// TODO
	if rs := data.ZoneMaster.PvScan(inapi.NsZoneSysCell(status.Zone.Meta.Id, ""), "", "", 100); rs.OK() {

		rss := rs.KvList()
		for _, v := range rss {

			var cell inapi.ResCell
			if err := v.Decode(&cell); err == nil {
				status.Zone.SyncCell(cell)
			}
		}
	}

	// refresh globa zones
	if rs := data.GlobalMaster.PvScan(
		inapi.NsGlobalSysZone(""), "", "", 50); !rs.OK() {
		hlog.Printf("warn", "refresh global-zones failed")
		return
	} else {

		rss := rs.KvList()

		for _, v := range rss {

			var o inapi.ResZone
			if err := v.Decode(&o); err == nil {
				found := false
				for i, pv := range status.GlobalZones {
					if pv.Meta.Id != o.Meta.Id {
						continue
					}
					found = true

					if pv.Meta.Updated < o.Meta.Updated {
						status.GlobalZones[i] = o
					}
				}
				if !found {
					status.GlobalZones = append(status.GlobalZones, o)
				}
			}
		}
	}

	// refresh zone-master list
	if rs := data.ZoneMaster.PvScan(
		inapi.NsZoneSysMasterNode(status.Host.Operate.ZoneId, ""), "", "", 100); !rs.OK() {
		hlog.Printf("warn", "refresh zone-master list failed")
		return
	} else {

		rss := rs.KvList()
		zms := inapi.ResZoneMasterList{
			Leader: status.Host.Meta.Id,
		}

		for _, v := range rss {

			var o inapi.ResZoneMasterNode
			if err := v.Decode(&o); err == nil {
				zms.Sync(o)
			}
		}

		if len(zms.Items) > 0 {
			status.ZoneMasterList.SyncList(zms)
		}
	}

	// refresh host keys
	if rs := data.ZoneMaster.PvScan(
		inapi.NsZoneSysHostSecretKey(status.Host.Operate.ZoneId, ""), "", "", 1000); rs.OK() {

		rs.KvEach(func(v *skv.ResultEntry) int {
			status.ZoneHostSecretKeys.Set(string(v.Key), v.Bytex().String())
			return 0
		})

	} else {
		hlog.Printf("warn", "refresh host key list failed")
	}

	// refresh host list
	if rs := data.ZoneMaster.PvScan(
		inapi.NsZoneSysHost(status.Host.Operate.ZoneId, ""), "", "", 1000); !rs.OK() {
		hlog.Printf("warn", "refresh host list failed")
		return
	} else {

		rss := rs.KvList()
		cell_counter := map[string]int32{}

		for _, v := range rss {

			var o inapi.ResHost
			if err := v.Decode(&o); err == nil {

				if o.Operate == nil {
					o.Operate = &inapi.ResHostOperate{}
				} else {
					// o.Operate.PortUsed.Clean()
				}
				cell_counter[o.Operate.CellId]++

				status.ZoneHostList.Sync(o)
				if gn := status.GlobalHostList.Item(o.Meta.Id); gn == nil {
					if rs := data.GlobalMaster.PvGet(inapi.NsGlobalSysHost(o.Operate.ZoneId, o.Meta.Id)); rs.NotFound() {
						data.GlobalMaster.PvPut(inapi.NsGlobalSysHost(o.Operate.ZoneId, o.Meta.Id), o, nil)
					}
				}
				status.GlobalHostList.Sync(o)

				// hlog.Printf("info", "refresh host refresh %s", o.Meta.Id)

			} else {
				// TODO
				hlog.Printf("error", "refresh host list ### %s", v.Value)
			}

			// hlog.Printf("info", "refresh host refresh %d", len(rss))
		}

		if !status.ZoneHostListImported {

			if len(cell_counter) > 0 {
				if rs := data.GlobalMaster.PvScan(inapi.NsGlobalSysCell(status.Host.Operate.ZoneId, ""), "", "", 1000); rs.OK() {
					rss := rs.KvList()
					for _, v := range rss {
						var cell inapi.ResCell
						if err := v.Decode(&cell); err == nil {
							if n, ok := cell_counter[cell.Meta.Id]; ok && n != cell.NodeNum {
								cell.NodeNum = n

								if rs := data.GlobalMaster.PvPut(inapi.NsGlobalSysCell(status.Host.Operate.ZoneId, cell.Meta.Id), cell, nil); rs.OK() {
									data.ZoneMaster.PvPut(inapi.NsZoneSysCell(status.Host.Operate.ZoneId, cell.Meta.Id), cell, nil)
								}
							}
						}
					}
				}
			}

			if rs := data.ZoneMaster.PvScan(
				inapi.NsZonePodInstance(status.Host.Operate.ZoneId, ""), "", "", 10000,
			); rs.OK() {

				rss := rs.KvList()

				for _, v := range rss {

					var pod inapi.Pod
					if err := v.Decode(&pod); err != nil {
						continue
					}

					status.ZonePodList.Set(&pod)

					for _, opRep := range pod.Operate.Replicas {

						host := status.ZoneHostList.Item(opRep.Node)
						if host == nil {
							continue
						}

						for _, v := range opRep.Ports {

							if v.HostPort == 0 {
								continue
							}

							if host.OpPortAlloc(v.HostPort) == 0 {
								continue
							}

							hlog.Printf("warn", "zone-master/host:%s.operate.ports refreshed", host.Meta.Id)

							data.ZoneMaster.PvPut(
								inapi.NsZoneSysHost(status.Host.Operate.ZoneId, host.Meta.Id),
								host,
								nil,
							)
						}
					}
				}
			}
		}

		if len(status.ZoneHostList.Items) > 0 {
			status.ZoneHostListImported = true
		}

		// hlog.Printf("info", "zone-master/host-list %d refreshed", len(status.ZoneHostList.Items))
	}

	// TODO
	if rs := data.ZoneMaster.PvScan(inapi.NsZonePodServiceMap(""), "", "", 10000); rs.OK() {

		nszs := []*inapi.NsPodServiceMap{}

		rs.KvEach(func(v *skv.ResultEntry) int {

			var nsz inapi.NsPodServiceMap
			if err := v.Decode(&nsz); err == nil {

				if pod := status.ZonePodList.Get(inapi.NsZonePodOpRepKey(nsz.Id, 0)); pod != nil {

					if inapi.OpActionAllow(pod.Operate.Action, inapi.OpActionDestroy) {
						return 0
					}

					if len(pod.Operate.Replicas) < 1 {
						return 0
					}

					nsz_sync := false
					nsz.Id = string(v.Key)

					for _, rep := range pod.Operate.Replicas {

						lan_addr := ""
						if host := status.ZoneHostList.Item(rep.Node); host != nil {
							lan_addr = host.Spec.PeerLanAddr
							if i := strings.IndexByte(lan_addr, ':'); i > 0 {
								lan_addr = lan_addr[:i]
							}
						}

						if lan_addr == "" {
							continue
						}

						for _, sport := range rep.Ports {

							pse := nsz.Get(uint16(sport.BoxPort))
							if pse == nil {
								continue
							}
							psh := inapi.NsPodServiceHostSliceGet(pse.Items, uint32(rep.Id))
							if psh == nil {
								continue
							}

							if psh.Ip != lan_addr {
								psh.Ip = lan_addr
							}

							nsz_sync = true
						}
					}

					if nsz_sync {
						data.ZoneMaster.PvPut(inapi.NsZonePodServiceMap(nsz.Id), nsz, nil)
						nszs, _ = inapi.NsPodServiceMapSliceSync(nszs, &nsz)
					}
				}
			}

			return 0
		})

		if !inapi.NsPodServiceMapSliceEqual(status.ZonePodServiceMaps, nszs) {
			status.ZonePodServiceMaps = nszs
			hlog.Printf("info", "zone-master/pod-service-maps refreshed %d", len(nszs))
		}
	}

	//
	if force_refresh {

		prs_key := inapi.NsZonePodReplicaStatus(status.Zone.Meta.Id, "", 0)

		rs := data.ZoneMaster.PvScan(prs_key, "", "", 100000)
		if !rs.OK() {
			hlog.Printf("warn", "refresh pod-replica-status list failed")
			return
		}

		rss := rs.KvList()
		for _, v := range rss {

			var prs inapi.PbPodRepStatus
			if err := v.Decode(&prs); err != nil {
				continue
			}

			status.ZonePodRepStatusSets, _ = inapi.PbPodRepStatusSliceSync(
				status.ZonePodRepStatusSets, &prs)
		}

		hlog.Printf("info", "zone-master/pod-replica-status refreshed %d", len(status.ZonePodRepStatusSets))
	}

	// hlog.Printf("debug", "zone-master/status refreshed")
}
