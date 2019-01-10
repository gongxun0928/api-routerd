// SPDX-License-Identifier: Apache-2.0

package network

import (
	"net/http"

	"github.com/RestGW/api-routerd/cmd/network/ethtool"
	"github.com/RestGW/api-routerd/cmd/network/netlink"
	"github.com/RestGW/api-routerd/cmd/network/networkd"

	"github.com/gorilla/mux"
)

func RouterConfigureEthtool(rw http.ResponseWriter, r *http.Request) {
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

	netlink.RegisterRouterNetlink(n)
	networkd.RegisterRouterNetworkd(n)

	// ethtool
	n.HandleFunc("/ethtool/{link}/{command}", RouterConfigureEthtool)
}
