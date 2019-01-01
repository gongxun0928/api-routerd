// SPDX-License-Identifier: Apache-2.0

package proc

import (
	"api-routerd/cmd/share"
	"encoding/json"
	"net/http"
	"path"
)

const (
	VMPath = "/proc/sys/vm"
)

type VM struct {
	Property string `json:"property"`
	Value    string `json:"value"`
}

func (req *VM) GetVM(rw http.ResponseWriter) error {
	line, r := share.ReadOneLineFile(path.Join(VMPath, req.Property))
	if r != nil {
		return r
	}

	vmProperty := VM{Property: req.Property, Value: line}
	json.NewEncoder(rw).Encode(vmProperty)

	return nil
}

func (req *VM) SetVM(rw http.ResponseWriter) error {
	r := share.WriteOneLineFile(path.Join(VMPath, req.Property), req.Value)
	if r != nil {
		return r
	}

	line, r := share.ReadOneLineFile(path.Join(VMPath, req.Property))
	if r != nil {
		return r
	}

	vmProperty := VM{Property: req.Property, Value: line}
	json.NewEncoder(rw).Encode(vmProperty)

	return nil
}
