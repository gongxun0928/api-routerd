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

type KMod struct {
	Name string `json:"name"`
	Args string `json:"args"`
}

func LsMod(w http.ResponseWriter) error {
	return proc.GetModules(w)
}

func (r *KMod) ModProbe() error {
	err := share.CheckBinaryExists("modprobe")
	if err != nil {
		return err
	}

	path, err := exec.LookPath("modprobe")
	if err != nil {
		return err
	}

	cmd := exec.Command(path, r.Name, r.Args)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("Failed to load module %s: %s", r.Name, stdout)
		return fmt.Errorf("Failed to load module '%s': %s", r.Name, stdout)
	}

	return nil
}

func (r *KMod) RmMod() error {
	err := share.CheckBinaryExists("rmmod")
	if err != nil {
		return err
	}

	path, err := exec.LookPath("rmmod")
	if err != nil {
		return err
	}

	cmd := exec.Command(path, r.Name)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("Failed to unload module %s: %s", r.Name, stdout)
		return fmt.Errorf("Failed to unload module '%s': %s", r.Name, stdout)
	}

	return nil
}
