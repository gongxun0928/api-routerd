// SPDX-License-Identifier: Apache-2.0

package journal

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/RestGW/api-routerd/cmd/share"

	"github.com/go-ini/ini"
	log "github.com/sirupsen/logrus"
)

const (
	journalConfPath = "/etc/systemd/journald.conf"
)

var journalConfig = map[string]string{
	"Storage":              "",
	"Compress":             "",
	"Seal":                 "",
	"SplitMode":            "",
	"SyncIntervalSec":      "",
	"RateLimitIntervalSec": "",
	"RateLimitBurst":       "",
	"SystemMaxUse":         "",
	"SystemKeepFree":       "",
	"SystemMaxFileSize":    "",
	"SystemMaxFiles":       "",
	"RuntimeMaxUse":        "",
	"RuntimeKeepFree":      "",
	"RuntimeMaxFileSize":   "",
	"RuntimeMaxFiles":      "",
	"MaxRetentionSec":      "",
	"MaxFileSec":           "",
	"ForwardToSyslog":      "",
	"ForwardToKMsg":        "",
	"ForwardToConsole":     "",
	"ForwardToWall":        "",
	"TTYPath":              "",
	"MaxLevelStore":        "",
	"MaxLevelSyslog":       "",
	"MaxLevelKMsg":         "",
	"MaxLevelConsole":      "",
	"MaxLevelWall":         "",
	"LineMax":              "",
	"ReadKMsg":             "",
}

func writeConfig() error {
	f, err := os.OpenFile(journalConfPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	conf := "[Journal]\n"
	for k, v := range journalConfig {
		if v != "" {
			conf += k + "=" + v
		} else {
			conf += "#" + k + "="
		}
		conf += "\n"
	}

	fmt.Fprintln(w, conf)
	w.Flush()

	return nil
}

func readConf() error {
	cfg, err := ini.Load(journalConfPath)
	if err != nil {
		return err
	}

	for k := range journalConfig {
		journalConfig[k] = cfg.Section("Journal").Key(k).String()
	}

	return nil
}

// GetConf Read and send journal conf
func GetConf(rw http.ResponseWriter) error {
	err := readConf()
	if err != nil {
		return err
	}

	return share.JSONResponse(journalConfig, rw)
}

// UpdateConf update the journal conf
func UpdateConf(rw http.ResponseWriter, r *http.Request) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Failed to parse HTTP request: %v", err)
		return err
	}

	conf := make(map[string]string)
	err = json.Unmarshal([]byte(body), &conf)
	if err != nil {
		log.Errorf("Failed to Decode HTTP request to json: %v", err)
		return err
	}

	err = readConf()
	if err != nil {
		return err
	}

	for k, v := range conf {
		_, ok := journalConfig[k]
		if ok {
			journalConfig[k] = v
		}
	}

	err = writeConfig()
	if err != nil {
		log.Errorf("Failed Write to journal conf: %v", err)
		return err
	}

	return share.JSONResponse(journalConfig, rw)
}
