// SPDX-License-Identifier: Apache-2.0

package network

import (
	"api-routerd/cmd/network/ethtool"
	"api-routerd/cmd/network/netlink"
	"api-routerd/cmd/network/networkd"
	"net/http"

	"github.com/gorilla/mux"
)

func NetworkLinkGet(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	link := vars["link"]

	switch r.Method {
	case "GET":
		netlink.GetLink(rw, r, link)
		break
	}
}

func NetworkLinkAdd(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		err := netlink.CreateLink(r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

func NetworkLinkDelete(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "DELETE":
		err := netlink.DeleteLink(r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

func NetworkLinkSet(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":
		err := netlink.SetLink(r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

func NetworkGetAddress(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	link := vars["link"]

	switch r.Method {
	case "GET":
		netlink.GetAddress(rw, link)
		break
	}
}

func NetworkAddAddress(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		err := netlink.AddAddress(r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func NetworkDeleteAddres(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "DELETE":
		err := netlink.DelAddress(r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		break
	}
}

func NetworkAddRoute(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		err := netlink.ConfigureRoutes(r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	case "PUT":
		err := netlink.ConfigureRoutes(r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	}
}

func NetworkDeleteRoute(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "DELETE":
		err := netlink.DeleteGateWay(r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		break
	}
}

func NetworkGetRoute(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := netlink.GetRoutes(rw, r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	}
}

func NetworkdConfigureNetwork(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		networkd.ConfigureNetworkFile(rw, r)
		break
	}
}

func NetworkdConfigureNetDev(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		networkd.ConfigureNetDevFile(rw, r)
		break
	}
}

func NetworkdConfigureLink(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		networkd.ConfigureLinkFile(rw, r)
		break
	}
}

func NetworkConfigureEthtool(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	link := vars["link"]
	command := vars["command"]

	ethtool := ethtool.Ethtool{Link: link, Action: command}

	switch r.Method {
	case "GET":

		err := ethtool.GetEthTool(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	}
}

func RegisterRouterNetwork(router *mux.Router) {
	n := router.PathPrefix("/network").Subrouter()

	// Link
	n.HandleFunc("/link/set", NetworkLinkSet)
	n.HandleFunc("/link/add", NetworkLinkAdd)
	n.HandleFunc("/link/delete", NetworkLinkDelete)
	n.HandleFunc("/link/get/{link}", NetworkLinkGet)
	n.HandleFunc("/link/get", NetworkLinkGet)

	// Address
	n.HandleFunc("/address/add", NetworkAddAddress)
	n.HandleFunc("/address/delete", NetworkDeleteAddres)
	n.HandleFunc("/address/get", NetworkGetAddress)
	n.HandleFunc("/address/get/{link}", NetworkGetAddress)

	// Route
	n.HandleFunc("/route/add", NetworkAddRoute)
	n.HandleFunc("/route/del", NetworkDeleteRoute)
	n.HandleFunc("/route/get", NetworkGetRoute)

	// systemd-networkd
	networkd.InitNetworkd()
	n.HandleFunc("/networkd/network", NetworkdConfigureNetwork)
	n.HandleFunc("/networkd/netdev", NetworkdConfigureNetDev)
	n.HandleFunc("/networkd/link", NetworkdConfigureLink)

	// ethtool
	n.HandleFunc("/ethtool/{link}/{command}", NetworkConfigureEthtool)
}
