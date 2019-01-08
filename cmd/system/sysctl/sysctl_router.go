// SPDX-License-Identifier: Apache-2.0

package sysctl

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func RouterSysctlGet(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetSysctl(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func RouterSysctlUpdate(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST", "PUT":
		s := new(Sysctl)
		err := json.NewDecoder(r.Body).Decode(&s)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		err = s.UpdateSysctl()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	}

	rw.WriteHeader(http.StatusOK)
}

func RouterSysctlDelete(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "DELETE":
		s := new(Sysctl)
		err := json.NewDecoder(r.Body).Decode(&s)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		err = s.DeleteSysctl()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	}

	rw.WriteHeader(http.StatusOK)
}

func RegisterRouterSysctl(router *mux.Router) {
	s := router.PathPrefix("/sysctl").Subrouter().StrictSlash(false)

	s.HandleFunc("/get", RouterSysctlGet)
	s.HandleFunc("/add", RouterSysctlUpdate)
	s.HandleFunc("/modify", RouterSysctlUpdate)
	s.HandleFunc("/delete", RouterSysctlDelete)
}
