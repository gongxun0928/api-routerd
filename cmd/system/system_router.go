// SPDX-License-Identifier: Apache-2.0

package system

import (
	"api-routerd/cmd/system/conf"
	"api-routerd/cmd/system/coredump"
	"api-routerd/cmd/system/hostname"
	"api-routerd/cmd/system/journal"
	"api-routerd/cmd/system/resolv"
	"api-routerd/cmd/system/resolved"
	"api-routerd/cmd/system/timesyncd"
	"api-routerd/cmd/system/timedate"
	"net/http"

	"github.com/gorilla/mux"
)

func RouterConfigureJournalConf(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := journal.GetJournalConf(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break

	case "POST":
		err := journal.UpdateJournalConf(rw, r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func ConfigureResolv(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":

		err := resolv.GetResolvConf(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	case "POST":

		err := resolv.UpdateResolvConf(rw, r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	case "DELETE":

		err := resolv.DeleteResolvConf(rw, r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	}
}

func ConfigureSystemdResolved(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":

		err := resolved.GetResolveConf(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	case "POST":

		err := resolved.UpdateResolveConf(rw, r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	case "DELETE":

		err := resolved.DeleteResolveConf(rw, r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	}
}

func ConfigureSystemdTimeSyncd(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":

		err := timesyncd.GetTimeSyncConf(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	case "POST":

		err := timesyncd.UpdateTimeSyncConf(rw, r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	case "DELETE":

		err := timesyncd.DeleteTimeSyncConf(rw, r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	}
}

func ConfigureSystemdCoreDump(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":

		err := coredump.GetCoreDumpConf(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	case "POST":

		err := coredump.UpdateCoreDumpConf(rw, r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	case "DELETE":

		err := coredump.DeleteCoreDumpConf(rw, r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	}
}

func ReadSudoersConfig(rw http.ResponseWriter, req *http.Request) {
	err := conf.GetSudoers(rw)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ReadSSHConfig(rw http.ResponseWriter, req *http.Request) {
	err := conf.SSHConfFileRead(rw)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func RegisterRouterSystem(router *mux.Router) {
	n := router.PathPrefix("/system").Subrouter()

	// hostname
	hostname.RegisterRouterHostname(router)

	// conf
	n.HandleFunc("/journal/conf", RouterConfigureJournalConf)
	n.HandleFunc("/journal/conf/update", RouterConfigureJournalConf)

	// resolv.conf
	n.HandleFunc("/resolv", ConfigureResolv)
	n.HandleFunc("/resolv/get", ConfigureResolv)
	n.HandleFunc("/resolv/add", ConfigureResolv)
	n.HandleFunc("/resolv/delete", ConfigureResolv)

	// systemd-resolved
	n.HandleFunc("/systemdresolved", ConfigureSystemdResolved)
	n.HandleFunc("/systemdresolved/get", ConfigureSystemdResolved)
	n.HandleFunc("/systemdresolved/add", ConfigureSystemdResolved)
	n.HandleFunc("/systemdresolved/delete", ConfigureSystemdResolved)

	// systemd-timesyncd
	n.HandleFunc("/systemdtimesyncd", ConfigureSystemdTimeSyncd)
	n.HandleFunc("/systemdtimesyncd/get", ConfigureSystemdTimeSyncd)
	n.HandleFunc("/systemdtimesyncd/add", ConfigureSystemdTimeSyncd)
	n.HandleFunc("/systemdtimesyncd/delete", ConfigureSystemdTimeSyncd)

	// coredump.conf
	n.HandleFunc("/coredump", ConfigureSystemdCoreDump)
	n.HandleFunc("/coredump/get", ConfigureSystemdCoreDump)
	n.HandleFunc("/coredump/add", ConfigureSystemdCoreDump)
	n.HandleFunc("/coredump/delete", ConfigureSystemdCoreDump)

	// Generic system confs
	n.HandleFunc("/conf/sudoers", ReadSudoersConfig)
	n.HandleFunc("/conf/sshd", ReadSSHConfig)

	// timedate
	timedate.RegisterRouterTimeDate(n)
}
