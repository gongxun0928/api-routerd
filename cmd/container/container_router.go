// SPDX-License-Identifier: Apache-2.0

package container

import (
	"github.com/RestGW/api-routerd/cmd/container/machine"
	"github.com/gorilla/mux"
)

//RegisterRouterContainer register with mux
func RegisterRouterContainer(router *mux.Router) {
	n := router.PathPrefix("/container").Subrouter()

	machine.InitMachine()
	machine.RegisterRouterMachine(n)
}
