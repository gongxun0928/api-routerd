// SPDX-License-Identifier: Apache-2.0

package hostname

import (
	"fmt"
	"net/http"

	"github.com/RestGW/api-routerd/cmd/share"

	log "github.com/sirupsen/logrus"
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

// Hostname commands
type Hostname struct {
	Property string `json:"property"`
	Value    string `json:"value"`
}

// SetHostname set hostname via dbus
func (h *Hostname) SetHostname() error {
	conn, err := NewConn()
	if err != nil {
		log.Errorf("Failed to get systemd bus connection: %v", err)
		return err
	}
	defer conn.Close()

	_, k := hostMethodInfo[h.Property]
	if !k {
		return fmt.Errorf("Failed to set hostname property: %s not found", h.Property)
	}

	return conn.SetHostName(h.Property, h.Value)
}

// GetHostname retrieves properties from hostnamed via dbus
func GetHostname(rw http.ResponseWriter, property string) error {
	conn, err := NewConn()
	if err != nil {
		log.Errorf("Failed to get systemd bus connection: %v", err)
		return err
	}
	defer conn.Close()

	for k := range hostNameInfo {
		p, err := conn.GetHostName(k)
		if err != nil {
			log.Errorf("Failed to get org.freedesktop.hostname1.%s", k)
			continue
		}

		hostNameInfo[k] = p
	}

	if property == "" {
		return share.JSONResponse(hostNameInfo, rw)
	}

	host := Hostname{
		Property: property,
		Value:    hostNameInfo[property],
	}

	return share.JSONResponse(host, rw)
}
