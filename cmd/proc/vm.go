// SPDX-License-Identifier: Apache-2.0

package proc

import (
	"net/http"
	"path"

	"github.com/RestGW/api-routerd/cmd/share"
)

const (
	vmPath = "/proc/sys/vm"
)

//VM JSON message
type VM struct {
	Property string `json:"property"`
	Value    string `json:"value"`
}

//GetVM read proc vm property
func (req *VM) GetVM(rw http.ResponseWriter) error {
	line, err := share.ReadOneLineFile(path.Join(vmPath, req.Property))
	if err != nil {
		return err
	}

	vmProperty := VM{
		Property: req.Property,
		Value:    line,
	}

	return share.JSONResponse(vmProperty, rw)
}

//SetVM write a value to VM
func (req *VM) SetVM(rw http.ResponseWriter) error {
	err := share.WriteOneLineFile(path.Join(vmPath, req.Property), req.Value)
	if err != nil {
		return err
	}

	line, err := share.ReadOneLineFile(path.Join(vmPath, req.Property))
	if err != nil {
		return err
	}

	vmProperty := VM{
		Property: req.Property,
		Value:    line,
	}

	return share.JSONResponse(vmProperty, rw)
}
