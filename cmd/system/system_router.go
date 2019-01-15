// SPDX-License-Identifier: Apache-2.0

package system

import (
	"net/http"

	"github.com/RestGW/api-routerd/cmd/system/conf"
	"github.com/RestGW/api-routerd/cmd/system/coredump"
	"github.com/RestGW/api-routerd/cmd/system/firewalld"
	"github.com/RestGW/api-routerd/cmd/system/group"
	"github.com/RestGW/api-routerd/cmd/system/hostname"
	"github.com/RestGW/api-routerd/cmd/system/journal"
	"github.com/RestGW/api-routerd/cmd/system/kmod"
	"github.com/RestGW/api-routerd/cmd/system/login"
	"github.com/RestGW/api-routerd/cmd/system/resolv"
	"github.com/RestGW/api-routerd/cmd/system/resolved"
	"github.com/RestGW/api-routerd/cmd/system/sysctl"
	"github.com/RestGW/api-routerd/cmd/system/timedate"
	"github.com/RestGW/api-routerd/cmd/system/timesyncd"
	"github.com/RestGW/api-routerd/cmd/system/user"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func routerConfigureJournalConf(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := journal.GetConf(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break

	case "POST":
		err := journal.UpdateConf(rw, r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func configureResolv(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":

		err := resolv.GetConf(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	case "POST":

		err := resolv.UpdateConf(rw, r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	case "DELETE":

		err := resolv.DeleteConf(rw, r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	}
}

func configureSystemdResolved(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":

		err := resolved.GetConf(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	case "POST":

		err := resolved.UpdateConf(rw, r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	case "DELETE":

		err := resolved.DeleteConf(rw, r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	}
}

func configureSystemdTimeSyncd(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":

		err := timesyncd.GetConf(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	case "POST":

		err := timesyncd.UpdateConf(rw, r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	case "DELETE":

		err := timesyncd.DeleteConf(rw, r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	}
}

func configureSystemdCoreDump(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":

		err := coredump.GetConf(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	case "POST":

		err := coredump.UpdateConf(rw, r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	case "DELETE":

		err := coredump.DeleteConf(rw, r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	}
}

func readSudoersConfig(rw http.ResponseWriter, req *http.Request) {
	err := conf.GetSudoers(rw)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func readSSHConfig(rw http.ResponseWriter, req *http.Request) {
	err := conf.SSHConfFileRead(rw)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

//RegisterRouterSystem register with mux
func RegisterRouterSystem(router *mux.Router) {
	n := router.PathPrefix("/system").Subrouter()

	// hostname
	hostname.RegisterRouterHostname(n)

	// timedate
	timedate.RegisterRouterTimeDate(n)

	// kmod
	kmod.RegisterRouterKMod(n)

	// group
	group.RegisterRouterGroup(n)

	// user
	user.RegisterRouterUser(n)

	// sysctl
	sysctl.RegisterRouterSysctl(n)

	// login
	login.RegisterRouterLogin(n)

	// firewalld
	err := firewalld.Init()
	if err != nil {
		log.Errorf("Failed to init firewalld: %s", err)
		return
	}

	firewalld.RegisterRouterFirewalld(n)

	// conf
	n.HandleFunc("/journal/conf", routerConfigureJournalConf)
	n.HandleFunc("/journal/conf/update", routerConfigureJournalConf)

	// resolv.conf
	n.HandleFunc("/resolv", configureResolv)
	n.HandleFunc("/resolv/get", configureResolv)
	n.HandleFunc("/resolv/add", configureResolv)
	n.HandleFunc("/resolv/delete", configureResolv)

	// systemd-resolved
	n.HandleFunc("/systemdresolved", configureSystemdResolved)
	n.HandleFunc("/systemdresolved/get", configureSystemdResolved)
	n.HandleFunc("/systemdresolved/add", configureSystemdResolved)
	n.HandleFunc("/systemdresolved/delete", configureSystemdResolved)

	// systemd-timesyncd
	n.HandleFunc("/systemdtimesyncd", configureSystemdTimeSyncd)
	n.HandleFunc("/systemdtimesyncd/get", configureSystemdTimeSyncd)
	n.HandleFunc("/systemdtimesyncd/add", configureSystemdTimeSyncd)
	n.HandleFunc("/systemdtimesyncd/delete", configureSystemdTimeSyncd)

	// coredump.conf
	n.HandleFunc("/coredump", configureSystemdCoreDump)
	n.HandleFunc("/coredump/get", configureSystemdCoreDump)
	n.HandleFunc("/coredump/add", configureSystemdCoreDump)
	n.HandleFunc("/coredump/delete", configureSystemdCoreDump)

	// Generic system confs
	n.HandleFunc("/conf/sudoers", readSudoersConfig)
	n.HandleFunc("/conf/sshd", readSSHConfig)
}
