// SPDX-License-Identifier: Apache-2.0

package network

import (
	"github.com/RestGW/api-routerd/cmd/network/ethtool"
	"github.com/RestGW/api-routerd/cmd/network/netlink"
	"github.com/RestGW/api-routerd/cmd/network/networkd"

	"github.com/gorilla/mux"
)

// RegisterRouterNetwork register with mux
func RegisterRouterNetwork(router *mux.Router) {
	n := router.PathPrefix("/network").Subrouter()

	netlink.RegisterRouterNetlink(n)
	networkd.RegisterRouterNetworkd(n)
	ethtool.RegisterRouterEthtool(n)
}
