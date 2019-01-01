package systemdtimesyncd

import (
	"api-routerd/cmd/share"
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/go-ini/ini"
	log "github.com/sirupsen/logrus"
)

const (
	TimeSyncdConfPath = "/etc/systemd/timesyncd.conf"
)

type TimeSyncConfig struct {
	NTP                []string `json:"NTP"`
	FallbackNTP        []string `json:"FallbackNTP"`
	RootDistanceMaxSec string   `json:"RootDistanceMaxSec"`
	PollIntervalMinSec string   `json:"PollIntervalMinSec"`
	PollIntervalMaxSec string   `json:"PollIntervalMaxSec"`
}

func (t *TimeSyncConfig) WriteTimeSyncConf() error {
	f, err := os.OpenFile(TimeSyncdConfPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	conf := "[Time]\n"

	ntpConf := "NTP="
	for _, s := range t.NTP {
		ntpConf += s + " "
	}
	conf += ntpConf + "\n"

	FallbackNTP := "FallbackNTP="
	for _, s := range t.FallbackNTP {
		FallbackNTP += s + " "
	}
	conf += FallbackNTP + "\n"

	if t.RootDistanceMaxSec != "" {
		conf += "RootDistanceMaxSec=" + t.RootDistanceMaxSec + "\n"
	}

	if t.PollIntervalMinSec != "" {
		conf += "PollIntervalMinSec=" + t.PollIntervalMinSec + "\n"
	}

	if t.PollIntervalMaxSec != "" {
		conf += "PollIntervalMaxSec=" + t.PollIntervalMaxSec + "\n"
	}

	fmt.Fprintln(w, conf)
	w.Flush()

	return nil
}

func ReadTimeSyncConf() (*TimeSyncConfig, error) {
	cfg, err := ini.Load(TimeSyncdConfPath)
	if err != nil {
		return nil, err
	}

	conf := new(TimeSyncConfig)

	ntp := cfg.Section("Time").Key("NTP").String()
	FallbackNTP := cfg.Section("Time").Key("FallbackNTP").String()

	conf.NTP = strings.Fields(ntp)
	conf.FallbackNTP = strings.Fields(FallbackNTP)

	return conf, nil
}

func GetTimeSyncConf(rw http.ResponseWriter) error {
	conf, err := ReadTimeSyncConf()
	if err != nil {
		return err
	}

	j, err := json.Marshal(conf)
	if err != nil {
		log.Errorf("Failed to encode json for resolv: %s", err)
		return err
	}

	rw.Write(j)

	return nil
}

func UpdateTimeSyncConf(rw http.ResponseWriter, r *http.Request) error {
	t := TimeSyncConfig{
		NTP:         []string{""},
		FallbackNTP: []string{""},
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("Failed to parse HTTP request: ", err)
		return err
	}

	err = json.Unmarshal([]byte(body), &t)
	if err != nil {
		log.Error("Failed to Decode HTTP request to json: ", err)
		return err
	}

	conf, err := ReadTimeSyncConf()
	if err != nil {
		return err
	}

	// update NTP
	for _, s := range t.NTP {
		if share.StringContains(conf.NTP, s) {
			continue
		}

		conf.NTP = append(conf.NTP, s)
	}

	// update fallback
	for _, s := range t.FallbackNTP {
		if share.StringContains(conf.FallbackNTP, s) {
			continue
		}
		conf.FallbackNTP = append(conf.FallbackNTP, s)
	}

	if t.RootDistanceMaxSec != "" {
		conf.RootDistanceMaxSec = t.RootDistanceMaxSec
	}

	if t.PollIntervalMinSec != "" {
		conf.PollIntervalMinSec = t.PollIntervalMinSec
	}

	if t.PollIntervalMaxSec != "" {
		conf.PollIntervalMaxSec = "t.PollIntervalMaxSec"
	}

	err = conf.WriteTimeSyncConf()
	if err != nil {
		log.Errorf("Failed Write to time sync conf: %s", err)
		return err
	}

	j, err := json.Marshal(conf)
	if err != nil {
		log.Errorf("Failed to encode json for resolv %s", err)
		return err
	}

	rw.Write(j)

	return nil
}

func DeleteTimeSyncConf(rw http.ResponseWriter, r *http.Request) error {
	t := TimeSyncConfig{
		NTP:         []string{""},
		FallbackNTP: []string{""},
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("Failed to parse HTTP request: ", err)
		return err
	}

	err = json.Unmarshal([]byte(body), &t)
	if err != nil {
		log.Error("Failed to Decode HTTP request to json: ", err)
		return err
	}

	conf, err := ReadTimeSyncConf()
	if err != nil {
		return err
	}

	// update NTP
	for _, s := range t.NTP {
		if !share.StringContains(conf.NTP, s) {
			continue
		}

		conf.NTP, _ = share.StringDeleteSlice(conf.NTP, s)
	}

	// update Fallback
	for _, s := range t.FallbackNTP {
		if !share.StringContains(conf.FallbackNTP, s) {
			continue
		}

		conf.FallbackNTP, _ = share.StringDeleteSlice(conf.FallbackNTP, s)
	}

	if t.RootDistanceMaxSec != "" {
		conf.RootDistanceMaxSec = t.RootDistanceMaxSec
	}

	if t.PollIntervalMinSec != "" {
		conf.PollIntervalMinSec = t.PollIntervalMinSec
	}

	if t.PollIntervalMaxSec != "" {
		conf.PollIntervalMaxSec = "t.PollIntervalMaxSec"
	}

	err = conf.WriteTimeSyncConf()
	if err != nil {
		log.Errorf("Failed Write to time sync conf: %s", err)
		return err
	}

	j, err := json.Marshal(conf)
	if err != nil {
		log.Errorf("Failed to encode json for resolv %s", err)
		return err
	}

	rw.Write(j)

	return nil
}
