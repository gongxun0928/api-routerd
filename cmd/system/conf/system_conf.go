// SPDX-License-Identifier: Apache-2.0

package conf

import (
	"github.com/RestGW/api-routerd/cmd/share"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	SudoersPath   = "/etc/sudoers"
	SSHConfigFile = "/etc/ssh/sshd_config"
)

type SudoersConf struct {
	Sudoers []string `json:"sudoers"`
}

type SSHdConf struct {
	Value string `json:"value"`
}

// GetSudoers read sudoers file
func GetSudoers(rw http.ResponseWriter) error {
	lines, err := share.ReadFullFile(SudoersPath)
	if err != nil {
		return err
	}

	var sudoers []string
	for _, line := range lines {
		if strings.Contains(line, "%") || strings.Contains(line, "Defaults") || strings.Contains(line, "root") {
			continue
		}

		line = strings.Replace(line, "\t", " ", 5)
		sudoers = append(sudoers, line)
	}

	sshdconf := SudoersConf{Sudoers: sudoers}

	return share.JsonResponse(sshdconf, rw)
}

// SSHConfFileRead read sshd configuration file
func SSHConfFileRead(rw http.ResponseWriter) error {
	lines, err := share.ReadFullFile(SSHConfigFile)
	if err != nil {
		log.Errorf("Failed to read: %s, %s", SSHConfigFile, err)
		return err
	}

	sshdConfMap := make(map[string]string)
	for _, line := range lines {
		fields := strings.Fields(line)
		paramName := fields[0]
		sshdConfMap[paramName] = fields[1]
	}

	return share.JsonResponse(sshdConfMap, rw)
}
