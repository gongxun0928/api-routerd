// SPDX-License-Identifier: Apache-2.0

package systemd

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
	SystemConfPath = "/etc/systemd/system.conf"
)

var SystemConfig = map[string]string{
	"LogLevel":                     "",
	"LogTarget":                    "",
	"LogColor":                     "",
	"LogLocation":                  "",
	"DumpCore":                     "",
	"ShowStatus":                   "",
	"CrashChangeVT":                "",
	"CrashShell":                   "",
	"CrashReboot":                  "",
	"CtrlAltDelBurstAction":        "",
	"CPUAffinity":                  "",
	"JoinControllers":              "",
	"RuntimeWatchdogSec":           "",
	"ShutdownWatchdogSec":          "",
	"CapabilityBoundingSe":         "",
	"SystemCallArchitectures":      "",
	"TimerSlackNSec":               "",
	"DefaultTimerAccuracySec":      "",
	"DefaultStandardOutput":        "",
	"DefaultStandardError":         "",
	"DefaultTimeoutStartSec":       "",
	"DefaultTimeoutStopSec":        "",
	"DefaultRestartSec":            "",
	"DefaultStartLimitIntervalSec": "",
	"DefaultStartLimitBurst":       "",
	"DefaultEnvironment":           "",
	"DefaultCPUAccounting":         "",
	"DefaultIOAccounting":          "",
	"DefaultIPAccounting":          "",
	"DefaultBlockIOAccounting":     "",
	"DefaultMemoryAccounting":      "",
	"DefaultTasksAccounting":       "",
	"DefaultTasksMax":              "",
	"DefaultLimitCPU":              "",
	"DefaultLimitFSIZE":            "",
	"DefaultLimitDATA":             "",
	"DefaultLimitSTACK":            "",
	"DefaultLimitCORE":             "",
	"DefaultLimitRSS":              "",
	"DefaultLimitNOFILE":           "",
	"DefaultLimitAS":               "",
	"DefaultLimitNPROC":            "",
	"DefaultLimitMEMLOCK":          "",
	"DefaultLimitLOCKS":            "",
	"DefaultLimitSIGPENDING":       "",
	"DefaultLimitMSGQUEUE":         "",
	"DefaultLimitNICE":             "",
	"DefaultLimitRTPRIO":           "",
	"DefaultLimitRTTIME":           "",
	"IPAddressAllow":               "",
	"IPAddressDeny":                "",
}

func WriteSystemConfig() error {
	f, err := os.OpenFile(SystemConfPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	conf := "[Manager]\n"
	for k, v := range SystemConfig {
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

func ReadSystemConf() error {
	cfg, err := ini.Load(SystemConfPath)
	if err != nil {
		return err
	}

	for k := range SystemConfig {
		SystemConfig[k] = cfg.Section("Manager").Key(k).String()
	}

	return nil
}

func GetSystemConf(rw http.ResponseWriter) error {
	err := ReadSystemConf()
	if err != nil {
		return err
	}

	return share.JSONResponse(SystemConfig, rw)
}

func UpdateSystemConf(rw http.ResponseWriter, r *http.Request) error {
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

	err = ReadSystemConf()
	if err != nil {
		return err
	}

	for k, v := range conf {
		_, ok := SystemConfig[k]
		if ok {
			SystemConfig[k] = v
		}
	}

	err = WriteSystemConfig()
	if err != nil {
		log.Errorf("Failed Write to system conf: %v", err)
		return err
	}

	return share.JSONResponse(SystemConfig, rw)
}
