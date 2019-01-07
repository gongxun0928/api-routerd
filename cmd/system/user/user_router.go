// SPDX-License-Identifier: Apache-2.0

package user

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func RouterUserAdd(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":

		g := new(User)
		err := json.NewDecoder(r.Body).Decode(&g)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		err = g.UserAdd()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	}

	rw.WriteHeader(http.StatusOK)
}

func RouterUserModify(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":

		g := new(User)
		err := json.NewDecoder(r.Body).Decode(&g)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		err = g.UserModify()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	}

	rw.WriteHeader(http.StatusOK)
}

func RouterUserDel(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "DELETE":

		g := new(User)
		err := json.NewDecoder(r.Body).Decode(&g)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		err = g.UserDel()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	}

	rw.WriteHeader(http.StatusOK)
}

func RegisterRouterUser(router *mux.Router) {
	s := router.PathPrefix("/user").Subrouter().StrictSlash(false)

	s.HandleFunc("/add", RouterUserAdd)
	s.HandleFunc("/delete", RouterUserDel)
	s.HandleFunc("/modify", RouterUserModify)
}
