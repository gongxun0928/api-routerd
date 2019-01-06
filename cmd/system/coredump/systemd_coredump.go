// SPDX-License-Identifier: Apache-2.0

package coredump

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/RestGW/api-routerd/cmd/share"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/go-ini/ini"
	log "github.com/sirupsen/logrus"
)

const (
	CoreDumpConfPath = "/etc/systemd/coredump.conf"
)

type CoreDumpConfig struct {
	Storage         string `json:"Storage"`
	Compress        string `json:"Compress"`
	ProcessSizeMax  string `json:"ProcessSizeMax"`
	ExternalSizeMax string `json:"ExternalSizeMax"`
	JournalSizeMax  string `json:"JournalSizeMax"`
}

func (c *CoreDumpConfig) WriteCoreDumpConfig() error {
	f, err := os.OpenFile(CoreDumpConfPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	conf := "[Coredump]\n"

	if c.Storage != "" {
		conf += "Storage=" + c.Storage + "\n"
	} else {
		conf += "#Storage=" + c.Storage + "\n"
	}

	if c.Compress != "" {
		conf += "Compress=" + c.Compress + "\n"
	} else {
		conf += "#Compress=" + c.Compress + "\n"
	}

	if c.ProcessSizeMax != "" {
		conf += "ProcessSizeMax=" + c.ProcessSizeMax + "\n"
	} else {
		conf += "#ProcessSizeMax=" + c.ProcessSizeMax + "\n"
	}

	if c.ExternalSizeMax != "" {
		conf += "ExternalSizeMax=" + c.ExternalSizeMax + "\n"
	} else {
		conf += "#ExternalSizeMax=" + c.ExternalSizeMax + "\n"
	}

	if c.JournalSizeMax != "" {
		conf += "JournalSizeMax=" + c.JournalSizeMax + "\n"
	} else {
		conf += "#JournalSizeMax=" + c.JournalSizeMax + "\n"
	}

	fmt.Fprintln(w, conf)
	w.Flush()

	return nil
}

func ReadCoreDumpConf() (*CoreDumpConfig, error) {
	cfg, err := ini.Load(CoreDumpConfPath)
	if err != nil {
		return nil, err
	}

	conf := new(CoreDumpConfig)
	conf.Storage = cfg.Section("Coredump").Key("Storage").String()
	conf.Compress = cfg.Section("Coredump").Key("Compress").String()
	conf.JournalSizeMax = cfg.Section("Coredump").Key("JournalSizeMax").String()
	conf.ExternalSizeMax = cfg.Section("Coredump").Key("ExternalSizeMax").String()
	conf.ProcessSizeMax = cfg.Section("Coredump").Key("ProcessSizeMax").String()

	return conf, nil
}

func GetCoreDumpConf(rw http.ResponseWriter) error {
	conf, err := ReadCoreDumpConf()
	if err != nil {
		return err
	}

	return share.JsonResponse(conf, rw)
}

func UpdateCoreDumpConf(rw http.ResponseWriter, r *http.Request) error {
	c := new(CoreDumpConfig)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("Failed to parse HTTP request: ", err)
		return err
	}

	err = json.Unmarshal([]byte(body), &c)
	if err != nil {
		log.Error("Failed to Decode HTTP request to json: ", err)
		return err
	}

	conf, err := ReadCoreDumpConf()
	if err != nil {
		return err
	}

	if c.Storage != "" {
		conf.Storage = c.Storage
	}

	if c.Compress != "" {
		conf.Compress = c.Compress
	}

	if c.JournalSizeMax != "" {
		conf.JournalSizeMax = c.JournalSizeMax
	}

	if c.ExternalSizeMax != "" {
		conf.ExternalSizeMax = c.ExternalSizeMax
	}

	if c.ProcessSizeMax != "" {
		conf.ProcessSizeMax = c.ProcessSizeMax
	}

	err = conf.WriteCoreDumpConfig()
	if err != nil {
		log.Errorf("Failed Write to resolv conf: %s", err)
		return err
	}

	return share.JsonResponse(conf, rw)
}

func DeleteCoreDumpConf(rw http.ResponseWriter, r *http.Request) error {
	c := new(CoreDumpConfig)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("Failed to parse HTTP request: ", err)
		return err
	}

	err = json.Unmarshal([]byte(body), &c)
	if err != nil {
		log.Error("Failed to Decode HTTP request to json: ", err)
		return err
	}

	conf, err := ReadCoreDumpConf()
	if err != nil {
		return err
	}

	if c.Storage == "#" {
		conf.Storage = ""
	}

	if c.Compress == "#" {
		conf.Compress = ""
	}

	if c.JournalSizeMax == "#" {
		conf.JournalSizeMax = ""
	}

	if c.ExternalSizeMax == "#" {
		conf.ExternalSizeMax = ""
	}

	if c.ProcessSizeMax == "#" {
		conf.ProcessSizeMax = ""
	}

	err = conf.WriteCoreDumpConfig()
	if err != nil {
		log.Errorf("Failed Write to coredump conf: %s", err)
		return err
	}

	j, err := json.Marshal(conf)
	if err != nil {
		log.Errorf("Failed to encode json for coredump %s", err)
		return err
	}

	rw.Write(j)

	return nil
}
