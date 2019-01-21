// SPDX-License-Identifier: Apache-2.0

package system

import (
	"github.com/RestGW/api-routerd/cmd/system/conf"
	"github.com/RestGW/api-routerd/cmd/system/coredump"
	"github.com/RestGW/api-routerd/cmd/system/firewalld"
	"github.com/RestGW/api-routerd/cmd/system/group"
	"github.com/RestGW/api-routerd/cmd/system/hostname"
	"github.com/RestGW/api-routerd/cmd/system/journal"
	"github.com/RestGW/api-routerd/cmd/system/kmod"
	"github.com/RestGW/api-routerd/cmd/system/login"
	"github.com/RestGW/api-routerd/cmd/system/resolv"
	"github.com/RestGW/api-routerd/cmd/system/resolve"
	"github.com/RestGW/api-routerd/cmd/system/sysctl"
	"github.com/RestGW/api-routerd/cmd/system/timedate"
	"github.com/RestGW/api-routerd/cmd/system/timesync"
	"github.com/RestGW/api-routerd/cmd/system/user"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// RegisterRouterSystem register with mux
func RegisterRouterSystem(router *mux.Router) {
	n := router.PathPrefix("/system").Subrouter()

	// system conf
	conf.RegisterRouterSystemConf(n)

	// coredump
	coredump.RegisterRouterCoreDump(n)

	// firewalld
	err := firewalld.Init()
	if err != nil {
		log.Errorf("Failed to init firewalld: %s", err)
		return
	}

	firewalld.RegisterRouterFirewalld(n)

	// group
	group.RegisterRouterGroup(n)

	// hostname
	hostname.InitHostname()
	hostname.RegisterRouterHostname(n)

	// journald
	journal.InitJournalConf()
	journal.RegisterRouterJournal(n)

	// kmod
	kmod.RegisterRouterKMod(n)

	// login
	login.RegisterRouterLogin(n)

	// sysctl
	sysctl.RegisterRouterSysctl(n)

	// /etc/resolv
	resolv.RegisterRouterResolv(n)

	// resolved
	resolve.RegisterRouterResolve(n)

	// timedate
	timedate.InitTimeDate()
	timedate.RegisterRouterTimeDate(n)

	// timesync
	timesync.RegisterRouterTimeSync(n)

	// user
	user.RegisterRouterUser(n)
}
