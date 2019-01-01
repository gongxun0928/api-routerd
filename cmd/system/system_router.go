// SPDX-License-Identifier: Apache-2.0

package system

import (
	"api-routerd/cmd/system/journal"
	resolv "api-routerd/cmd/system/resolve"
	"api-routerd/cmd/system/systemdresolved"
	"api-routerd/cmd/system/systemdtimesyncd"
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

		err := systemdresolved.GetResolveConf(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	case "POST":

		err := systemdresolved.UpdateResolveConf(rw, r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	case "DELETE":

		err := systemdresolved.DeleteResolveConf(rw, r)
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

		err := systemdtimesyncd.GetTimeSyncConf(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	case "POST":

		err := systemdtimesyncd.UpdateTimeSyncConf(rw, r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	case "DELETE":

		err := systemdtimesyncd.DeleteTimeSyncConf(rw, r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	}
}

func RegisterRouterSystem(router *mux.Router) {
	n := router.PathPrefix("/system").Subrouter()

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
}
