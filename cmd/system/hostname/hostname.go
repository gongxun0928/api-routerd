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

var hostNameInfo = map[string]string{
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

var hostMethodInfo = map[string]string{
	"SetHostname":       "",
	"SetStaticHostname": "",
	"SetPrettyHostname": "",
	"SetIconName":       "",
	"SetChassis":        "",
	"SetDeployment":     "",
	"SetLocation":       "",
}

//Hostname commands
type Hostname struct {
	Property string `json:"property"`
	Value    string `json:"value"`
}

//SetHostname set hostname via dbus
func (hostname *Hostname) SetHostname() error {
	conn, err := share.GetSystemBusPrivateConn()
	if err != nil {
		log.Errorf("Failed to get systemd bus connection: %v", err)
		return err
	}
	defer conn.Close()

	_, k := hostMethodInfo[hostname.Property]
	if !k {
		return fmt.Errorf("Failed to set hostname property: %s not found", hostname.Property)
	}

	h := conn.Object(dbusInterface, dbusPath)
	r := h.Call(dbusInterface+"."+hostname.Property, 0, hostname.Value, false).Err
	if r != nil {
		log.Errorf("Failed to set hostname: %v", r)
		return errors.New("Failed to set hostname")
	}

	return nil
}

//GetHostname retrives properties from hostnamed via dbus
func GetHostname(rw http.ResponseWriter, property string) error {
	conn, err := share.GetSystemBusPrivateConn()
	if err != nil {
		log.Errorf("Failed to get dbus connection: %v", err)
		return err
	}
	defer conn.Close()

	h := conn.Object(dbusInterface, dbusPath)
	for k := range hostNameInfo {
		p, perr := h.GetProperty(dbusInterface + "." + k)
		if perr != nil {
			log.Errorf("Failed to get org.freedesktop.hostname1.%s", k)
			continue
		}

		hv, b := p.Value().(string)
		if !b {
			log.Errorf("Received unexpected type as value, %s expected string got: %s", property, hv)
			continue
		}

		hostNameInfo[k] = hv
	}

	if property == "" {
		return share.JsonResponse(hostNameInfo, rw)
	}

	host := Hostname{
		Property: property,
		Value:    hostNameInfo[property],
	}

	return share.JsonResponse(host, rw)
}
