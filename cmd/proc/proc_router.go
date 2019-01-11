// SPDX-License-Identifier: Apache-2.0

package proc

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Info struct {
	Path     string `json:"path"`
	Property string `json:"property"`
	Value    string `json:"value"`
}

func RouterGetProcNetDev(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetNetDev(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func RouterGetProcVersion(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetVersion(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func RouterGetProcPlatformInformation(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetPlatformInformation(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func RouterGetProcVirtualization(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetVirtualization(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func RouterGetProcUserStat(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetUserStat(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func RouterGetProcTemperatureStat(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetTemperatureStat(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func RouterGetProcNetStat(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	protocol := vars["protocol"]

	switch r.Method {
	case "GET":
		err := GetNetStat(rw, protocol)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func RouterGetProcPidNetStat(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	protocol := vars["protocol"]
	pid := vars["pid"]

	switch r.Method {
	case "GET":
		err := GetNetStatPid(rw, protocol, pid)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func RouterGetProcInterfaceStat(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetInterfaceStat(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func RouterGetProcProtoCountersStat(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetProtoCountersStat(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func RouterGetProcGetSwapMemoryStat(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetSwapMemoryStat(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func RouterGetProcVirtualMemoryStat(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetVirtualMemoryStat(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func RouterGetProcCPUInfo(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetCPUInfo(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func RouterGetProcCPUTimeStat(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetCPUTimeStat(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func RouterGetProcAvgStat(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetAvgStat(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func ConfigureProcSysVM(rw http.ResponseWriter, r *http.Request) {
	var err error

	vars := mux.Vars(r)
	vm := VM{Property: vars["path"]}

	switch r.Method {
	case "GET":
		err = vm.GetVM(rw)
		break
	case "PUT":

		v := new(Info)

		err = json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		vm.Value = v.Value
		err = vm.SetVM(rw)
		break
	}

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func ConfigureProcSysNet(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	proc := SysNet{Path: vars["path"], Property: vars["conf"], Link: vars["link"]}

	switch r.Method {
	case "GET":
		err := proc.GetSysNet(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	case "PUT":

		v := new(Info)
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		proc.Value = v.Value
		err = proc.SetSysNet(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func RouterGetProcMisc(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetMisc(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func RouterGetProcNetArp(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetNetArp(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func RouterGetProcModules(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetModules(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func RouterGetProcProcess(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid := vars["pid"]
	property := vars["property"]

	switch r.Method {
	case "GET":

		err := GetProcessInfo(rw, pid, property)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func RouterGetPartitions(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetPartitions(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

func RouterGetIOCounters(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetIOCounters(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

func RouterGetDiskUsage(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetDiskUsage(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

func RegisterRouterProc(router *mux.Router) {
	n := router.PathPrefix("/proc").Subrouter().StrictSlash(false)

	n.HandleFunc("/avgstat", RouterGetProcAvgStat)
	n.HandleFunc("/cpuinfo", RouterGetProcCPUInfo)
	n.HandleFunc("/cputimestat", RouterGetProcCPUTimeStat)
	n.HandleFunc("/diskusage", RouterGetDiskUsage)
	n.HandleFunc("/interface-stat", RouterGetProcInterfaceStat)
	n.HandleFunc("/iocounters", RouterGetIOCounters)
	n.HandleFunc("/misc", RouterGetProcMisc)
	n.HandleFunc("/modules", RouterGetProcModules)
	n.HandleFunc("/net/arp", RouterGetProcNetArp)
	n.HandleFunc("/netdev", RouterGetProcNetDev)
	n.HandleFunc("/netstat/{protocol}", RouterGetProcNetStat)
	n.HandleFunc("/partitions", RouterGetPartitions)
	n.HandleFunc("/platform", RouterGetProcPlatformInformation)
	n.HandleFunc("/process/{pid}/{property}/", RouterGetProcProcess)
	n.HandleFunc("/proto-counter-stat", RouterGetProcProtoCountersStat)
	n.HandleFunc("/proto-pid-stat/{pid}/{protocol}", RouterGetProcPidNetStat)
	n.HandleFunc("/swap-memory", RouterGetProcGetSwapMemoryStat)
	n.HandleFunc("/sys/net/{path}/{link}/{conf}", ConfigureProcSysNet)
	n.HandleFunc("/sys/vm/{path}", ConfigureProcSysVM)
	n.HandleFunc("/temperaturestat", RouterGetProcTemperatureStat)
	n.HandleFunc("/userstat", RouterGetProcUserStat)
	n.HandleFunc("/version", RouterGetProcVersion)
	n.HandleFunc("/virtual-memory", RouterGetProcVirtualMemoryStat)
	n.HandleFunc("/virtualization", RouterGetProcVirtualization)
}
