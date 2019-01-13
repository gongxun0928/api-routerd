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

	switch f.Property {
	case "list-services":
		services, err := c.ListServices()
		if err != nil {
			return err
		}

		return share.JsonResponse(services, rw)

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
	}

	return nil
}

func Init() error {
	FirewalldMethods = share.NewSet()

	FirewalldMethods.Add("list-services")
	FirewalldMethods.Add("get-default-zone")
	FirewalldMethods.Add("get-zone-settings")

	return nil
}
