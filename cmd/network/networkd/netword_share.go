// SPDX-License-Identifier: Apache-2.0

package networkd

import (
	"api-routerd/cmd/share"

	log "github.com/sirupsen/logrus"
)

const (
	NetworkdUnitPath = "/etc/systemd/network"
)

type Match struct {
	MAC    string `json:"MAC"`
	Driver string `json:"Driver"`
	Name   string `json:"Name"`
}

func InitNetworkd() error {
	err := share.CreateDirectory(NetworkdUnitPath, 0777)
	if err != nil {
		log.Errorf("Failed create network unit path %s: %s", NetworkdUnitPath, err)
		return err
	}

	return nil
}
