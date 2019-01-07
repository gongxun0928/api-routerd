// SPDX-License-Identifier: Apache-2.0

package group

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func RouterGroupAdd(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":

		g := new(Group)
		err := json.NewDecoder(r.Body).Decode(&g)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		err = g.GroupAdd()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	}

	rw.WriteHeader(http.StatusOK)
}

func RouterGroupModify(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":

		g := new(Group)
		err := json.NewDecoder(r.Body).Decode(&g)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		err = g.GroupModify()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	}

	rw.WriteHeader(http.StatusOK)
}

func RouterGroupDel(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "DELETE":

		g := new(Group)
		err := json.NewDecoder(r.Body).Decode(&g)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		err = g.GroupDel()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	}

	rw.WriteHeader(http.StatusOK)
}

func RegisterRouterGroup(router *mux.Router) {
	s := router.PathPrefix("/group").Subrouter().StrictSlash(false)

	s.HandleFunc("/add", RouterGroupAdd)
	s.HandleFunc("/delete", RouterGroupDel)
	s.HandleFunc("/modify", RouterGroupModify)
}
