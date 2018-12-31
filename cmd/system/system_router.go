// SPDX-License-Identifier: Apache-2.0

package system

import (
	"api-routerd/cmd/system/resolve"
	"api-routerd/cmd/system/systemdresolved"
	"api-routerd/cmd/system/journal"
	"github.com/gorilla/mux"
	"net/http"
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

func NetworkConfigureResolv(rw http.ResponseWriter, r *http.Request) {
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

func NetworkConfigureSystemdResolved(rw http.ResponseWriter, r *http.Request) {
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

func RegisterRouterSystem(router *mux.Router) {
	n := router.PathPrefix("/system").Subrouter()

	// conf
	n.HandleFunc("/journal/conf", RouterConfigureJournalConf)
	n.HandleFunc("/journal/conf/update", RouterConfigureJournalConf)

	// resolv.conf
	n.HandleFunc("/resolv", NetworkConfigureResolv)
	n.HandleFunc("/resolv/get", NetworkConfigureResolv)
	n.HandleFunc("/resolv/add", NetworkConfigureResolv)
	n.HandleFunc("/resolv/delete", NetworkConfigureResolv)

	// systemd-resolved
	n.HandleFunc("/systemdresolved", NetworkConfigureSystemdResolved)
	n.HandleFunc("/systemdresolved/get", NetworkConfigureSystemdResolved)
	n.HandleFunc("/systemdresolved/add", NetworkConfigureSystemdResolved)
	n.HandleFunc("/systemdresolved/delete", NetworkConfigureSystemdResolved)
}
