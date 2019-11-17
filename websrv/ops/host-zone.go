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

package ops

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hooto/hlog4g/hlog"
	"github.com/hooto/iam/iamapi"
	"github.com/lessos/lessgo/crypto/idhash"
	"github.com/lessos/lessgo/types"

	iam_db "github.com/hooto/iam/store"
	"github.com/sysinner/incore/data"
	"github.com/sysinner/incore/inapi"
	"github.com/sysinner/incore/status"
)

func (c Host) ZoneListAction() {

	ls := inapi.GeneralObjectList{}
	defer c.RenderJson(&ls)

	//
	rs := data.DataGlobal.NewReader(nil).KeyRangeSet(
		inapi.NsGlobalSysZone(""), inapi.NsGlobalSysZone("")).
		LimitNumSet(100).Query()

	for _, v := range rs.Items {

		var zone inapi.ResZone

		if err := v.Decode(&zone); err != nil || zone.Meta.Id == "" {
			continue
		}

		if c.Params.Get("fields") == "cells" {

			rs2 := data.DataGlobal.NewReader(nil).KeyRangeSet(
				inapi.NsGlobalSysCell(zone.Meta.Id, ""), inapi.NsGlobalSysCell(zone.Meta.Id, "")).
				LimitNumSet(100).Query()

			for _, v2 := range rs2.Items {

				var cell inapi.ResCell

				if err := v2.Decode(&cell); err == nil {
					zone.Cells = append(zone.Cells, &cell)
				}
			}
		}

		ls.Items = append(ls.Items, zone)
	}

	ls.Kind = "HostZoneList"
}

func (c Host) ZoneEntryAction() {

	var set struct {
		inapi.GeneralObject
		inapi.ResZone
	}

	defer c.RenderJson(&set)

	if obj := data.DataGlobal.NewReader(inapi.NsGlobalSysZone(c.Params.Get("id"))).Query(); obj.OK() {

		if err := obj.Decode(&set.ResZone); err != nil {
			set.Error = &types.ErrorMeta{"400", err.Error()}
		} else {

			if c.Params.Get("fields") == "cells" {

				rs2 := data.DataGlobal.NewReader(nil).KeyRangeSet(
					inapi.NsGlobalSysCell(set.Meta.Id, ""), inapi.NsGlobalSysCell(set.Meta.Id, "")).
					LimitNumSet(100).Query()

				for _, v2 := range rs2.Items {

					var cell inapi.ResCell

					if err := v2.Decode(&cell); err == nil {
						set.Cells = append(set.Cells, &cell)
					}
				}
			}

		}
	}

	if set.Meta.Id != "" {
		set.Kind = "HostZone"
	} else {
		set.Error = &types.ErrorMeta{"404", "Item Not Found"}
	}
}

func (c Host) ZoneSetAction() {

	var set struct {
		inapi.GeneralObject
		inapi.ResZone
	}

	defer c.RenderJson(&set)

	if err := c.Request.JsonDecode(&set.ResZone); err != nil {
		set.Error = &types.ErrorMeta{"400", err.Error()}
		return
	}

	set.Meta.Id = strings.ToLower(set.Meta.Id)

	if mat := inapi.ResSysZoneIdReg.MatchString(set.Meta.Id); !mat {
		set.Error = &types.ErrorMeta{
			Code:    "400",
			Message: "Zone Id must consist of letters or numbers, and begin with a letter",
		}
		return
	}

	set.WanApi = strings.Trim(strings.TrimSpace(set.WanApi), "/")
	if len(set.WanApi) > 0 {
		if !strings.HasPrefix(set.WanApi, "http://") &&
			!strings.HasPrefix(set.WanApi, "https://") {
			set.Error = types.NewErrorMeta("400", fmt.Sprintf("Invalid Address (%s)", set.WanApi))
			return
		}
	}

	for _, addr := range set.WanAddrs {

		if !inapi.HostNodeAddress(addr).Valid() {
			set.Error = &types.ErrorMeta{
				Code:    "400",
				Message: fmt.Sprintf("Invalid Address (%s)", addr),
			}
			return
		}
	}

	for _, addr := range set.LanAddrs {

		if !inapi.HostNodeAddress(addr).Valid() {
			set.Error = &types.ErrorMeta{
				Code:    "400",
				Message: fmt.Sprintf("Invalid Address (%s)", addr),
			}
			return
		}
	}

	if len(set.ImageServices) > 5 {
		set.Error = types.NewErrorMeta("400",
			fmt.Sprintf("the number of Image Services cannot be greater than %d", 5))
		return
	}

	for _, v := range set.ImageServices {

		if v.Driver != inapi.PodSpecBoxImageDocker &&
			v.Driver != inapi.PodSpecBoxImagePouch {
			set.Error = types.NewErrorMeta("400", fmt.Sprintf("Invalid ImageService Driver (%s)", v.Driver))
			return
		}

		if _, err := url.ParseRequestURI(v.Url); err != nil {
			set.Error = types.NewErrorMeta("400", fmt.Sprintf("Invalid ImageService URL (%s)", v.Url))
			return
		}
	}

	if obj := data.DataGlobal.NewReader(inapi.NsGlobalSysZone(set.Meta.Id)).Query(); obj.OK() {

		var prev inapi.ResZone
		if err := obj.Decode(&prev); err == nil {

			if prev.Meta.Created > 0 {
				set.Meta.Created = prev.Meta.Created
			}
		}
	}

	if set.Meta.Created == 0 {
		set.Meta.Created = uint64(types.MetaTimeNow())
	}

	set.Meta.Updated = uint64(types.MetaTimeNow())

	data.DataGlobal.NewWriter(inapi.NsGlobalSysZone(set.Meta.Id), set.ResZone).Commit()

	set.Kind = "HostZone"
}

func (c Host) ZoneAccChargeKeyRefreshAction() {

	var set inapi.GeneralObject
	defer c.RenderJson(&set)

	var (
		zone_id = c.Params.Get("zone_id")
		zone    inapi.ResZone
	)

	if zone_id != status.Zone.Meta.Id {
		set.Error = types.NewErrorMeta("400", "Zone Not Found")
		return
	}

	if rs := data.DataGlobal.NewReader(inapi.NsGlobalSysZone(zone_id)).Query(); rs.OK() {
		rs.Decode(&zone)
	}

	if zone.Meta.Id != zone_id {
		set.Error = types.NewErrorMeta("400", "Zone Not Found")
		return
	}

	init_akacc := iamapi.AccessKey{
		User: "sysadmin",
		AccessKey: "00" + idhash.HashToHexString(
			[]byte(fmt.Sprintf("sys/zone/iam_acc_charge/ak/%s", zone_id)), 14),
		SecretKey: idhash.RandBase64String(40),
		Bounds: []iamapi.AccessKeyBound{{
			Name: "sys/zm/" + zone_id,
		}},
		Description: "ZoneMaster AccCharge",
	}
	if err := iam_db.AccessKeyReset(init_akacc); err != nil {
		set.Error = types.NewErrorMeta("500", "database/iam error "+err.Error())
		return
	}

	zone.OptionSet("iam/acc_charge/access_key", init_akacc.AccessKey)
	zone.OptionSet("iam/acc_charge/secret_key", init_akacc.SecretKey)

	if rs := data.DataGlobal.NewWriter(inapi.NsGlobalSysZone(zone_id), zone).Commit(); !rs.OK() {
		set.Error = types.NewErrorMeta("500", "database/global error "+rs.Message)
		return
	}

	if rs := data.DataZone.NewWriter(inapi.NsZoneSysInfo(zone_id), zone).Commit(); !rs.OK() {
		set.Error = types.NewErrorMeta("500", "database/zone error "+rs.Message)
		return
	}

	hlog.Printf("warn", "ops/zone/acc-charge/key/reset %s, %s...",
		init_akacc.AccessKey, init_akacc.SecretKey[:20])

	status.Zone = &zone

	set.Kind = "Zone"
}
