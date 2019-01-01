// SPDX-License-Identifier: Apache-2.0

package conf

import (
	"api-routerd/cmd/share"
	"encoding/json"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	SudoersPath = "/etc/sudoers"
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
	json.NewEncoder(rw).Encode(sshdconf)

	return nil
}

// SSHConfFileRead read sshd configuration file
func SSHConfFileRead(rw http.ResponseWriter) error {
	lines, err := share.ReadFullFile(SSHConfigFile)
	if err != nil {
		log.Errorf("Failed to read: %s", SSHConfigFile)
		return err
	}

	sshdConfMap := make(map[string]string)
	for _, line := range lines {
		fields := strings.Fields(line)
		paramName := fields[0]
		sshdConfMap[paramName] = fields[1]
	}

	j, err := json.Marshal(sshdConfMap)
	if err != nil {
		log.Error("Failed to encode JSON payload")
		return err
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(j)

	return nil
}
