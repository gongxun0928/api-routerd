// SPDX-License-Identifier: Apache-2.0

package firewalld

import (
	"net/http"

	"github.com/gorilla/mux"
)

func RouterGetFirewalld(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	property := vars["property"]
	value := vars["value"]

	switch r.Method {
	case "GET":

		firewall := Firewall{
			Property: property,
			Value:    value,
		}

		err := firewall.GetFirewalld(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		break
	}
}

func RegisterRouterFirewalld(router *mux.Router) {
	f := router.PathPrefix("/firewalld").Subrouter().StrictSlash(false)

	f.HandleFunc("/get/{property}", RouterGetFirewalld)
	f.HandleFunc("/get/{property}/{value}", RouterGetFirewalld)
}
