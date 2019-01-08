// SPDX-License-Identifier: Apache-2.0

package sysctl

import (
	"fmt"
	"net/http"
	"os/exec"
	"strings"

	"github.com/RestGW/api-routerd/cmd/share"

	log "github.com/sirupsen/logrus"
)

const (
	SysctlPath = "/etc/sysctl.conf"
)

type Sysctl struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Apply string `json:"apply"`
}

// Apply sysctl conf to system
func (s *Sysctl) ApplySysctl() error {
	b, err := share.ParseBool(s.Apply)
	if err != nil {
		return fmt.Errorf("Failed to apply: %s", s.Key)
	}

	if b != true {
		return nil
	}

	path, err := exec.LookPath("sysctl")
	if err != nil {
		return err
	}

	cmd := exec.Command(path, "-p")
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("Failed to load sysctl variable: %s", stdout)
		return fmt.Errorf("Failed to load sysctl variable: %s", stdout)
	}

	return nil
}

// Read sysctl config to a map
func ReadSysctlConfig() (map[string]string, error) {
	lines, err := share.ReadFullFile(SysctlPath)
	if err != nil {
		return nil, err
	}

	sysctl := make(map[string]string)
	for _, line := range lines {
		tokens := strings.Split(line, "=")
		if len(tokens) != 2 {
			log.Errorf("could not parse line %s", line)
			continue
		}

		k := strings.TrimSpace(tokens[0])
		v := strings.TrimSpace(tokens[1])
		sysctl[k] = v
	}

	return sysctl, nil
}

// WriteSysctlConfig: write config to file
func WriteSysctlConfig(sysctl map[string]string) error {
	var lines []string
	var line string

	for k, v := range sysctl {
		line = k + "=" + v
		lines = append(lines, line)
	}

	return share.WriteFullFile(SysctlPath, lines)
}

// GetSysctl read sysctl file
func GetSysctl(rw http.ResponseWriter) error {
	sysctl, err := ReadSysctlConfig()
	if err != nil {
		return err
	}

	return share.JsonResponse(sysctl, rw)
}

// UpdateSysctl update sysctl file
func (s *Sysctl) UpdateSysctl() error {
	sysctl, err := ReadSysctlConfig()
	if err != nil {
		return err
	}

	sysctl[s.Key] = s.Value

	err = WriteSysctlConfig(sysctl)
	if err != nil {
		return err
	}

	return s.ApplySysctl()
}

// DeleteSysctl delete sysctl value in file
func (s *Sysctl) DeleteSysctl() error {
	sysctl, err := ReadSysctlConfig()
	if err != nil {
		return err
	}

	_, ok := sysctl[s.Key]
	if !ok {
		return fmt.Errorf("Failed to delete sysctl parameter '%s'. Key not found", s.Key)
	}

	delete(sysctl, s.Key)
	err = WriteSysctlConfig(sysctl)
	if err != nil {
		return err
	}

	return s.ApplySysctl()
}
