// SPDX-License-Identifier: Apache-2.0

package machine

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/RestGW/api-routerd/cmd/share"
	"github.com/coreos/go-systemd/machine1"
)

type Machine struct {
	Path     string `json:"path"`
	Property string `json:"property"`
	Value    string `json:"value"`
}

var MachineMethods *share.Set

func (m *Machine) MachineMethodGet(rw http.ResponseWriter) error {
	c, err := machine1.New()
	if err != nil {
		return err
	}

	b := MachineMethods.Contains(m.Path)
	if !b {
		return fmt.Errorf("Failed to call method machine: %s not found", m.Path)
	}

	switch m.Path {
	case "list-images":
		images, err := c.ListImages()
		if err != nil {
			return err
		}

		return share.JsonResponse(images, rw)

	case "list-machines":
		machines, err := c.ListMachines()
		if err != nil {
			return err
		}

		return share.JsonResponse(machines, rw)

	case "get-machine":
		machine, err := c.GetMachine(m.Property)
		if err != nil {
			return err
		}

		return share.JsonResponse(machine, rw)
	case "describe-machine":
		machine, err := c.GetMachine(m.Property)
		if err != nil {
			return err
		}

		return share.JsonResponse(machine, rw)

	case "get-image":
		image, err := c.GetImage(m.Property)
		if err != nil {
			return err
		}

		return share.JsonResponse(image, rw)
	case "get-machine-by-pid":
		pid, err := strconv.ParseInt(m.Property, 10, 32)
		if err != nil {
			return nil
		}

		p, err := c.GetMachineByPID(uint(pid))
		if err != nil {
			return err
		}

		return share.JsonResponse(p, rw)
	case "get-machine-address":
		addr, err := c.GetMachineAddresses(m.Property)
		if err != nil {
			return err
		}

		return share.JsonResponse(addr, rw)
	}

	return nil
}

func (m *Machine) MachineMethodConfigure(rw http.ResponseWriter) error {
	c, err := machine1.New()
	if err != nil {
		return err
	}

	b := MachineMethods.Contains(m.Path)
	if !b {
		return fmt.Errorf("Failed to call method machine: %s not found", m.Path)
	}

	switch m.Path {
	case "terminate-machine":
		err := c.TerminateMachine(m.Property)
		if err != nil {
			return err
		}

		return nil
	}

	return nil
}

func InitMachine() error {
	MachineMethods = share.NewSet()

	MachineMethods.Add("list-images")
	MachineMethods.Add("list-machines")
	MachineMethods.Add("get-machine")
	MachineMethods.Add("get-image")
	MachineMethods.Add("get-machine-by-pid")
	MachineMethods.Add("get-machine-address")
	MachineMethods.Add("describe-machine")
	MachineMethods.Add("terminate-machine")

	return nil
}
