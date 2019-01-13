// SPDX-License-Identifier: Apache-2.0

package firewalld

import (
	"fmt"

	"github.com/RestGW/api-routerd/cmd/share"

	"github.com/godbus/dbus"
)

const (
	dbusInterface = "org.fedoraproject.FirewallD1"
	dbusPath      = "/org/fedoraproject/FirewallD1"
)

type Conn struct {
	conn   *dbus.Conn
	object dbus.BusObject
}

type ZoneSettings struct {
	Version     string   `json:"version"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Services    []string `json:"services"`
}

func NewConn() (*Conn, error) {
	c := new(Conn)

	conn, err := share.GetSystemBusPrivateConn()
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to system bus:", err)
	}

	c.conn = conn
	c.object = conn.Object(dbusInterface, dbus.ObjectPath(dbusPath))

	return c, nil
}

func (c *Conn) Close() {
	c.conn.Close()
}

func (c *Conn) ListServices() ([]string, error) {
	var services []string

	err := c.object.Call(dbusInterface+".listServices", 0).Store(&services)
	if err != nil {
		return nil, err
	}

	return services, nil
}

func (c *Conn) GetDefaultZone() (string, error) {
	var zone string

	err := c.object.Call(dbusInterface+".getDefaultZone", 0).Store(&zone)
	if err != nil {
		return "", err
	}

	return zone, nil
}

func (c *Conn) GetZoneSettings(zone string) (*ZoneSettings, error) {
	out := []interface{}{}

	err := c.object.Call(dbusInterface+".getZoneSettings", 0, zone).Store(&out)
	if err != nil {
		return nil, err
	}

	z := new(ZoneSettings)
	for i, el := range out {
		switch i {
		case 0:
			continue
		case 1:
			z.Name = el.(string)
			break
		case 2:
			z.Description = el.(string)
			break
		case 3, 4:
			continue
		case 5:
			z.Services = el.([]string)
		}
	}

	return z, nil
}
