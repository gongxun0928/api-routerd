// SPDX-License-Identifier: Apache-2.0

package proc

import (
	"errors"
	"net/http"
	"path"

	"github.com/RestGW/api-routerd/cmd/share"
)

const (
	sysNetPath     = "/proc/sys/net"
	sysNetPathCore = "core"
	sysNetPathIPv4 = "ipv4"
	sysNetPathIPv6 = "ipv6"
)

//SysNet Json request
type SysNet struct {
	Path     string `json:"path"`
	Property string `json:"property"`
	Value    string `json:"value"`
	Link     string `json:"link"`
}

//getPath read info from proc
func (req *SysNet) getPath() (string, error) {
	var procPath string

	switch req.Path {
	case sysNetPathCore:
		procPath = path.Join(path.Join(sysNetPath, sysNetPathCore), req.Property)
		break
	case sysNetPathIPv4:

		if req.Link != "" {
			procPath = path.Join(path.Join(path.Join(path.Join(sysNetPath, sysNetPathIPv4), "conf"), req.Link), req.Property)
		} else {
			procPath = path.Join(path.Join(sysNetPath, sysNetPathIPv4), req.Property)
		}
		break
	case sysNetPathIPv6:

		if req.Link != "" {
			procPath = path.Join(path.Join(path.Join(path.Join(sysNetPath, sysNetPathIPv6), "conf"), req.Link), req.Property)
		} else {
			procPath = path.Join(path.Join(sysNetPath, sysNetPathIPv6), req.Property)
		}
		break
	default:
		return "", errors.New("Path not found")
	}

	return procPath, nil
}

//GetSysNet read proc value and send response
func (req *SysNet) GetSysNet(rw http.ResponseWriter) error {
	path, err := req.getPath()
	if err != nil {
		return err
	}

	line, err := share.ReadOneLineFile(path)
	if err != nil {
		return err
	}

	property := SysNet{
		Path:     req.Path,
		Property: req.Property,
		Value:    line,
		Link:     req.Link,
	}

	return share.JsonResponse(property, rw)
}

//SetSysNet sets a value to proc
func (req *SysNet) SetSysNet(rw http.ResponseWriter) error {
	path, err := req.getPath()
	if err != nil {
		return err
	}

	return share.WriteOneLineFile(path, req.Value)
}
