// SPDX-License-Identifier: Apache-2.0

package kmod

import (
	"fmt"
	"net/http"
	"os/exec"

	"github.com/RestGW/api-routerd/cmd/proc"
	"github.com/RestGW/api-routerd/cmd/share"
	log "github.com/sirupsen/logrus"
)

type KModules struct {
	Name string `json:"name"`
	Args string `json:"args"`
}

func LsMod(w http.ResponseWriter) error {
	return proc.GetModules(w)
}

func (r *KModules) ModProbe() error {
	err := share.CheckBinaryExists("modprobe")
	if err != nil {
		return err
	}

	cmd := exec.Command("/usr/sbin/modprobe", r.Name, r.Args)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("Failed to load module %s: %s", r.Name, stdout)
		return fmt.Errorf("Failed to load module '%s': %s", r.Name, stdout)
	}

	return nil
}
