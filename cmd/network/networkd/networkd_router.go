// SPDX-License-Identifier: Apache-2.0

package networkd

import (
	"net/http"

	"github.com/gorilla/mux"
)

func RouterConfigureNetworkdNetwork(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		ConfigureNetworkFile(rw, r)
		break
	}
}

func RouterConfigureNetworkdNetDev(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		ConfigureNetDevFile(rw, r)
		break
	}
}

func RouterConfigureNetworkdLink(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		ConfigureLinkFile(rw, r)
		break
	}
}

func RegisterRouterNetworkd(router *mux.Router) {
	InitNetworkd()

	n := router.PathPrefix("/networkd").Subrouter().StrictSlash(false)

	// systemd-networkd
	n.HandleFunc("/network", RouterConfigureNetworkdNetwork)
	n.HandleFunc("/netdev", RouterConfigureNetworkdNetDev)
	n.HandleFunc("/link", RouterConfigureNetworkdLink)
}
