// SPDX-License-Identifier: Apache-2.0

package hostname

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/RestGW/api-routerd/cmd/share"

	log "github.com/sirupsen/logrus"
)

const (
	dbusInterface = "org.freedesktop.hostname1"
	dbusPath      = "/org/freedesktop/hostname1"
)

var HostNameInfo = map[string]string{
	"Hostname":                  "",
	"StaticHostname":            "",
	"PrettyHostname":            "",
	"IconName":                  "",
	"Chassis":                   "",
	"Deployment":                "",
	"Location":                  "",
	"KernelName":                "",
	"KernelRelease":             "",
	"KernelVersion":             "",
	"OperatingSystemPrettyName": "",
	"OperatingSystemCPEName":    "",
	"HomeURL":                   "",
}

var HostMethodInfo = map[string]string{
	"SetHostname":       "",
	"SetStaticHostname": "",
	"SetPrettyHostname": "",
	"SetIconName":       "",
	"SetChassis":        "",
	"SetDeployment":     "",
	"SetLocation":       "",
}

type Hostname struct {
	Property string `json:"property"`
	Value    string `json:"value"`
}

func (hostname *Hostname) SetHostname() error {
	conn, err := share.GetSystemBusPrivateConn()
	if err != nil {
		log.Error("Failed to get systemd bus connection: ", err)
		return err
	}
	defer conn.Close()

	_, k := HostMethodInfo[hostname.Property]
	if !k {
		return fmt.Errorf("Failed to set hostname property: %s not found", hostname.Property)
	}

	h := conn.Object(dbusInterface, dbusPath)
	r := h.Call(dbusInterface+"."+hostname.Property, 0, hostname.Value, false).Err
	if r != nil {
		log.Errorf("Failed to set hostname: %s", r)
		return errors.New("Failed to set hostname")
	}

	return nil
}

func GetHostname(rw http.ResponseWriter, property string) error {
	conn, err := share.GetSystemBusPrivateConn()
	if err != nil {
		log.Error("Failed to get dbus connection: ", err)
		return err
	}
	defer conn.Close()

	h := conn.Object(dbusInterface, dbusPath)
	for k, _ := range HostNameInfo {
		p, perr := h.GetProperty(dbusInterface + "." + k)
		if perr != nil {
			log.Errorf("Failed to get org.freedesktop.hostname1.%s", k)
			continue
		}

		hv, b := p.Value().(string)
		if !b {
			log.Error("Received unexpected type as value, expected string got :", property, hv)
			continue
		}

		HostNameInfo[k] = hv
	}

	if property == "" {
		return share.JsonResponse(HostNameInfo, rw)
	}

	host := Hostname{
		Property: property,
		Value: HostNameInfo[property],
	}

	return share.JsonResponse(host, rw)
}
