// SPDX-License-Identifier: Apache-2.0

package proc

import (
	"net/http"
	"path"

	"github.com/RestGW/api-routerd/cmd/share"
)

const (
	VMPath = "/proc/sys/vm"
)

type VM struct {
	Property string `json:"property"`
	Value    string `json:"value"`
}

func (req *VM) GetVM(rw http.ResponseWriter) error {
	line, err := share.ReadOneLineFile(path.Join(VMPath, req.Property))
	if err != nil {
		return err
	}

	vmProperty := VM{
		Property: req.Property,
		Value:    line,
	}

	return share.JsonResponse(vmProperty, rw)
}

func (req *VM) SetVM(rw http.ResponseWriter) error {
	err := share.WriteOneLineFile(path.Join(VMPath, req.Property), req.Value)
	if err != nil {
		return err
	}

	line, err := share.ReadOneLineFile(path.Join(VMPath, req.Property))
	if err != nil {
		return err
	}

	vmProperty := VM{
		Property: req.Property,
		Value:    line,
	}

	return share.JsonResponse(vmProperty, rw)
}
