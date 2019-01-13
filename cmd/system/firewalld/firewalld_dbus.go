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
	Interfaces  []string `json:"interfaces"`
}

type ServiceSettings struct {
	Version      string            `json:"version"`
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	Ports        [][]interface{}   `json:"ports"`
	Destinations map[string]string `json:"destinations"`
	Protocols    []string          `json:"protocols"`
	SourcePorts  [][]interface{}   `json:"source_ports"`
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
		case 1:
			z.Name = el.(string)
			break
		case 2:
			z.Description = el.(string)
			break
		case 5:
			z.Services = el.([]string)
			break
		case 10:
			z.Interfaces = el.([]string)
		}
	}

	return z, nil
}

func (c *Conn) GetServiceSettings(zone string) (*ServiceSettings, error) {
	out := []interface{}{}

	err := c.object.Call(dbusInterface+".getServiceSettings", 0, zone).Store(&out)
	if err != nil {
		return nil, err
	}

	fmt.Println(out)
	s := new(ServiceSettings)
	for i, el := range out {
		switch i {
		case 1:
			s.Name = el.(string)
			break
		case 2:
			s.Description = el.(string)
			break
		case 3:
			s.Ports = el.([][]interface{})
			break
		case 5:
			s.Destinations = el.(map[string]string)
			break
		case 6:
			s.Protocols = el.([]string)
			break
		case 7:
			s.SourcePorts = el.([][]interface{})
			break
		}
	}

	return s, nil

}

func (c *Conn) AddPort(zone string, port string, protocol string) (string, error) {
	 var r string

	err := c.object.Call(dbusInterface+".zone.addPort", 0, zone, port, protocol, 0).Store(&r)
	if err != nil {
		return r, err
	}

	return r, nil
}

func (c *Conn) RemovePort(zone string, port string, protocol string) (string, error) {
	var r string

	err := c.object.Call(dbusInterface+".zone.removePort", 0, zone, port, protocol).Store(&r)
	if err != nil {
		return r, err
	}

	return r, nil
}
