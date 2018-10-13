package hostlet

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hooto/hlog4g/hlog"
	"github.com/lessos/lessgo/encoding/json"
	"github.com/lessos/lessgo/types"
	ps_disk "github.com/shirou/gopsutil/disk"

	"github.com/sysinner/incore/config"
	"github.com/sysinner/incore/hostlet/nstatus"
	"github.com/sysinner/incore/inapi"
)

type QuotaConfig struct {
	mu            sync.Mutex
	path          string
	sync_maps_sum string
	sync_ids_sum  string
	Items         []*QuotaProject `json:"items,omitempty"`
	Updated       int64           `json:"updated"`
	IdOffset      int             `json:"id_offset"`
}

type QuotaProject struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Soft int64  `json:"soft"`
	Hard int64  `json:"hard"`
	Used int64  `json:"used"`
}

func (it *QuotaConfig) Fetch(name string) *QuotaProject {
	it.mu.Lock()
	defer it.mu.Unlock()

	for _, v := range it.Items {
		if name == v.Name {
			return v
		}
	}

	return nil
}

func (it *QuotaConfig) FetchById(id int) *QuotaProject {
	it.mu.Lock()
	defer it.mu.Unlock()

	for _, v := range it.Items {
		if id == v.Id {
			return v
		}
	}

	return nil
}

func (it *QuotaConfig) FetchOrCreate(name string) *QuotaProject {

	it.mu.Lock()
	defer it.mu.Unlock()

	for _, v := range it.Items {
		if name == v.Name {
			return v
		}
	}

	var (
		bind_id = 0
	)

	if it.IdOffset < 100 || it.IdOffset >= 100000 {
		it.IdOffset = 100
	}

	for i := it.IdOffset; i <= 110000; i++ {
		hit := false
		for _, v := range it.Items {
			if v.Id == i {
				hit = true
				break
			}
		}
		if !hit {
			bind_id = i
			break
		}
	}

	if bind_id == 0 {
		return nil
	}

	p := &QuotaProject{
		Id:   bind_id,
		Name: name,
	}

	it.Items = append(it.Items, p)

	it.IdOffset = bind_id + 1

	return p
}

func (it *QuotaConfig) Remove(name string) {

	it.mu.Lock()
	defer it.mu.Unlock()

	for i, v := range it.Items {
		if name == v.Name {
			it.Items = append(it.Items[:i], it.Items[i+1:]...)
			return
		}
	}
}

func (it *QuotaConfig) Sync() error {
	it.Updated = time.Now().Unix()
	return json.EncodeToFile(it, it.path, "  ")
}

func (it *QuotaConfig) SyncVendor() error {

	//
	maps := ""
	for _, v := range it.Items {
		if v.Id < 1 {
			continue
		}
		maps += fmt.Sprintf("%d:%s\n", v.Id, filepath.Clean(config.Config.PodHomeDir+"/"+v.Name))
	}

	maps_sum := fmt.Sprintf("%x", sha1.Sum([]byte(maps)))
	if maps_sum != it.sync_maps_sum {
		if err := put_file("/etc/projects", maps); err != nil {
			return err
		}
		it.sync_maps_sum = maps_sum
	}

	return nil
}

func put_file(path, data string) error {

	fp, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer fp.Close()

	fp.Seek(0, 0)
	fp.Truncate(0)

	_, err = fp.WriteString(data)
	return err
}

var (
	quotaInited           = false
	quotaRefreshed  int64 = 0
	quotaCmd              = "xfs_quota"
	quotaMountpoint       = ""
	regMultiSpace         = regexp.MustCompile("\\ \\ +")
	quotaConfig     QuotaConfig
	err             error
)

func QuotaKeeperInit() error {

	if quotaInited {
		return nil
	}

	defer func() {
		quotaInited = true
	}()

	_, err := exec.Command(quotaCmd, "-V").Output()
	if err != nil {
		return errors.New("command " + quotaCmd + " not found")
	}

	devs, _ := ps_disk.Partitions(false)

	sort.Slice(devs, func(i, j int) bool {
		if strings.Compare(devs[i].Mountpoint, devs[j].Mountpoint) > 0 {
			return true
		}
		return false
	})

	for _, d := range devs {
		if !strings.HasPrefix(config.Config.PodHomeDir, d.Mountpoint) {
			continue
		}
		if d.Fstype != "xfs" {
			return errors.New("invalid fstype (" + d.Fstype + ") to enable quota")
		}
		if !strings.Contains(d.Opts, "prjquota") {
			return errors.New("the option:prjquota required on mountpoint of " + d.Mountpoint)
		}
		quotaMountpoint = d.Mountpoint
		break
	}

	if quotaMountpoint == "" {
		return errors.New("no quota path found")
	}

	_, err = exec.Command(quotaCmd, "-x", "-c", "report", quotaMountpoint).Output()
	if err != nil {
		return err
	}

	//
	cfgpath := "/etc/sysinner_vol_quota.json"
	if err := json.DecodeFile(cfgpath, &quotaConfig); err != nil && !os.IsNotExist(err) {
		return err
	}
	quotaConfig.path = cfgpath

	quotaRefreshed = time.Now().Unix()

	return nil
}

func podVolQuotaRefresh() error {

	if err := QuotaKeeperInit(); err != nil {
		hlog.Printf("error", "Failed to Enable Disk Quota: %s", err.Error())
		return nil
	}

	tn := time.Now().Unix()

	if quotaRefreshed < 1 || (tn-quotaRefreshed) < 10 {
		return nil
	}

	args := []string{
		"-xc",
		"report",
		quotaMountpoint,
	}
	out, err := exec.Command(quotaCmd, args...).Output()
	if err != nil {
		return err
	}

	var (
		lines       = strings.Split(regMultiSpace.ReplaceAllString(string(out), " "), "\n")
		report_gots = types.ArrayUint32{}
		device_ok   = false
	)
	for _, v := range lines {

		vs := strings.Split(strings.TrimSpace(v), " ")
		if len(vs) < 5 || len(vs[0]) < 2 {
			continue
		}

		if vs[0][0] != '#' {
			continue
		}

		if vs[0] == "#0" {
			device_ok = true
		}

		id, err := strconv.ParseInt(vs[0][1:], 10, 32)
		if err != nil || id < 100 {
			continue
		}

		report_gots.Set(uint32(id))

		proj := quotaConfig.FetchById(int(id))
		if proj == nil {
			continue
		}

		if i64, err := strconv.ParseInt(vs[1], 10, 64); err == nil {
			proj.Used = i64 * 1024
		}

		if i64, err := strconv.ParseInt(vs[2], 10, 64); err == nil {
			proj.Soft = i64 * 1024
		}

		if i64, err := strconv.ParseInt(vs[3], 10, 64); err == nil {
			proj.Hard = i64 * 1024
		}
	}

	// hlog.Printf("info", "get info %d  %d", len(quotaConfig.Items), len(report_gots))

	if device_ok {
		dels := []string{}
		for _, v := range quotaConfig.Items {
			if !report_gots.Has(uint32(v.Id)) {
				dels = append(dels, v.Name)
			}
		}
		if len(dels) > 0 {

			for _, v := range dels {
				quotaConfig.Remove(v)
			}

			if err := quotaConfig.SyncVendor(); err != nil {
				return err
			}
		}
	}

	if err := quotaConfig.Sync(); err != nil {
		return err
	}

	// hlog.Printf("info", "get info %d, pods %d", len(quotaConfig.Items), len(nstatus.PodRepActives))

	nstatus.PodRepActives.Each(func(pod *inapi.Pod) {

		if pod.Spec == nil || pod.Operate.Replica == nil {
			return
		}

		if !inapi.OpActionAllow(pod.Operate.Action, inapi.OpActionStart) {
			return
		}

		spec_vol := pod.Spec.Volume("system")
		if spec_vol == nil {
			return
		}

		name := inapi.NsZonePodOpRepKey(pod.Meta.ID, pod.Operate.Replica.Id)

		//
		path := filepath.Clean(config.Config.PodHomeDir + "/" + name)
		fp, err := os.Open(path)
		if err != nil {
			return
		}
		fp.Close()

		//
		proj := quotaConfig.FetchOrCreate(name)
		if proj == nil {
			hlog.Printf("error", "failed to create quota/project : %s", name)
			return
		}

		if proj.Soft == spec_vol.SizeLimit {
			return
		}

		err = quotaConfig.SyncVendor()
		if err != nil {
			hlog.Printf("warn", "config init %s", err.Error())
			return
		}

		args := []string{
			"-x",
			"-c",
			fmt.Sprintf("\"project -s %d\"", proj.Id),
			quotaMountpoint,
		}

		_, err = exec.Command(quotaCmd, args...).Output()
		if err != nil {
			hlog.Printf("warn", "quota init %s", err.Error())
			return
		}

		//
		args = []string{
			"-x",
			"-c",
			fmt.Sprintf("limit -p bsoft=%d bhard=%d %d",
				spec_vol.SizeLimit, spec_vol.SizeLimit, proj.Id),
			quotaMountpoint,
		}

		if out, err := exec.Command(quotaCmd, args...).Output(); err != nil {
			hlog.Printf("warn", "quota limit %s, {{{%s}}}", err.Error(), string(out))
			fmt.Println(strings.Join(args, " "))
			return
		}
	})

	if err := quotaConfig.Sync(); err != nil {
		return err
	}

	for _, v := range quotaConfig.Items {

		if v.Soft < 1 {
			continue
		}

		pod := nstatus.PodRepActives.Get(v.Name)
		if pod != nil && inapi.OpActionAllow(pod.Operate.Action, inapi.OpActionStart) {
			continue
		}

		//
		args := []string{
			"-x",
			"-c",
			fmt.Sprintf("limit -p bsoft=0 bhard=0 %d", v.Id),
			quotaMountpoint,
		}

		_, err = exec.Command(quotaCmd, args...).Output()
		if err != nil {
			hlog.Printf("warn", "quota clean project %s, err %s", v.Name, err.Error())
			return err
		}

		hlog.Printf("warn", "quota clean project %s, done", v.Name)
	}

	quotaRefreshed = tn

	return nil
}