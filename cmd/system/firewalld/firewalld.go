// SPDX-License-Identifier: Apache-2.0

package firewalld

import (
	"fmt"
	"net/http"

	"github.com/RestGW/api-routerd/cmd/share"

	log "github.com/sirupsen/logrus"
)

type Firewall struct {
	Property string `json:"property"`
	Value    string `json:"value"`

	Zone      string `json:"zone"`
	Port      string `json:"port"`
	Protocol  string `json:"protocol"`
	Permanent bool   `json:"permanent,omitempty"`
}

var FirewalldMethods *share.Set

func (f *Firewall) GetFirewalld(rw http.ResponseWriter) error {
	c, err := NewConn()
	if err != nil {
		return err
	}
	defer c.Close()

	b := FirewalldMethods.Contains(f.Property)
	if !b {
		return fmt.Errorf("Failed to call method firewalld: %s not found", f.Property)
	}

	log.Debugf("Get Firewalld passthrough: %s", f.Property)

	fmt.Println(f)
	switch f.Property {
	case "list-services":
		services, err := c.ListServices()
		if err != nil {
			return err
		}

		return share.JsonResponse(services, rw)

	case "get-zones":
		z, err := c.GetZones()
		if err != nil {
			return err
		}

		return share.JsonResponse(z, rw)

	case "list-all-zones":
		z, err := c.ListAllZones()
		if err != nil {
			return err
		}

		return share.JsonResponse(z, rw)

	case "list-ports":
		z, err := c.ListPorts(f.Value)
		if err != nil {
			return err
		}

		return share.JsonResponse(z, rw)

	case "get-default-zone":
		z, err := c.GetDefaultZone()
		if err != nil {
			return err
		}

		return share.JsonResponse(z, rw)

	case "get-zone-settings":
		z, err := c.GetZoneSettings(f.Value)
		if err != nil {
			return err
		}

		return share.JsonResponse(z, rw)

	case "get-service-settings":
		z, err := c.GetServiceSettings(f.Value)
		if err != nil {
			return err
		}

		return share.JsonResponse(z, rw)
	}

	return nil
}

func (f *Firewall) AddFirewalld(rw http.ResponseWriter) error {
	c, err := NewConn()
	if err != nil {
		return err
	}
	defer c.Close()

	b := FirewalldMethods.Contains(f.Property)
	if !b {
		return fmt.Errorf("Failed to call method firewalld: %s not found", f.Property)
	}

	log.Debugf("Set Firewalld passthrough: %s", f.Property)

	switch f.Property {
	case "add-port":
		var r string

		if f.Permanent == true {
			r, err = c.AddPortPermanent(f.Zone, f.Port, f.Protocol)
		} else {
			r, err = c.AddPort(f.Zone, f.Port, f.Protocol)
		}

		if err != nil {
			return err
		}

		return share.JsonResponse(r, rw)
	}

	return nil
}

func (f *Firewall) DeleteFirewalld(rw http.ResponseWriter) error {
	c, err := NewConn()
	if err != nil {
		return err
	}
	defer c.Close()

	b := FirewalldMethods.Contains(f.Property)
	if !b {
		return fmt.Errorf("Failed to call method firewalld: %s not found", f.Property)
	}

	log.Debugf("Delete Firewalld passthrough: %s", f.Property)

	switch f.Property {
	case "remove-port":
		var r string

		if f.Permanent == true {
			r, err = c.RemovePortPermanent(f.Zone, f.Port, f.Protocol)
		} else {
			r, err = c.RemovePort(f.Zone, f.Port, f.Protocol)
		}

		if err != nil {
			return err
		}

		return share.JsonResponse(r, rw)
	}

	return nil
}

func Init() error {
	FirewalldMethods = share.NewSet()

	FirewalldMethods.Add("list-services")
	FirewalldMethods.Add("get-zones")
	FirewalldMethods.Add("list-all-zones")
	FirewalldMethods.Add("list-ports")
	FirewalldMethods.Add("get-default-zone")
	FirewalldMethods.Add("get-zone-settings")
	FirewalldMethods.Add("get-service-settings")
	FirewalldMethods.Add("add-port")
	FirewalldMethods.Add("remove-port")

	return nil
}
