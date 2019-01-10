// SPDX-License-Identifier: Apache-2.0

package netlink

import (
	"net/http"

	"github.com/gorilla/mux"
)

func RouterLinkGet(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	link := vars["link"]

	switch r.Method {
	case "GET":
		GetLink(rw, r, link)
		break
	}
}

func RouterLinkAdd(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		err := CreateLink(r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

func RouterLinkDelete(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "DELETE":
		err := DeleteLink(r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

func RouterLinkSet(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":
		err := SetLink(r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

func RouterGetAddress(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	link := vars["link"]

	switch r.Method {
	case "GET":
		GetAddress(rw, link)
		break
	}
}

func RouterAddAddress(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		err := AddAddress(r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func RouterDeleteAddres(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "DELETE":
		err := DelAddress(r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		break
	}
}

func RouterAddRoute(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		err := ConfigureRoutes(r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	case "PUT":
		err := ConfigureRoutes(r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	}
}

func RouterDeleteRoute(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "DELETE":
		err := DeleteGateWay(r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		break
	}
}

func RouterGetRoute(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetRoutes(rw, r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	}
}

func RegisterRouterNetlink(n *mux.Router) {
	// Link
	n.HandleFunc("/link/set", RouterLinkSet)
	n.HandleFunc("/link/add", RouterLinkAdd)
	n.HandleFunc("/link/delete", RouterLinkDelete)
	n.HandleFunc("/link/get/{link}", RouterLinkGet)
	n.HandleFunc("/link/get", RouterLinkGet)

	// Address
	n.HandleFunc("/address/add", RouterAddAddress)
	n.HandleFunc("/address/delete", RouterDeleteAddres)
	n.HandleFunc("/address/get", RouterGetAddress)
	n.HandleFunc("/address/get/{link}", RouterGetAddress)

	// Route
	n.HandleFunc("/route/add", RouterAddRoute)
	n.HandleFunc("/route/del", RouterDeleteRoute)
	n.HandleFunc("/route/get", RouterGetRoute)
}
