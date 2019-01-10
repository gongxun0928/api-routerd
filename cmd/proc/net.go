// SPDX-License-Identifier: Apache-2.0

package proc

import (
	"errors"
	"net/http"
	"path"

	"github.com/RestGW/api-routerd/cmd/share"
)

const (
	SysNetPath     = "/proc/sys/net"
	SysNetPathCore = "core"
	SysNetPathIPv4 = "ipv4"
	SysNetPathIPv6 = "ipv6"
)

type SysNet struct {
	Path     string `json:"path"`
	Property string `json:"property"`
	Value    string `json:"value"`
	Link     string `json:"link"`
}

func (req *SysNet) GetSysNetPath() (string, error) {
	var procPath string

	switch req.Path {
	case SysNetPathCore:
		procPath = path.Join(path.Join(SysNetPath, SysNetPathCore), req.Property)
		break
	case SysNetPathIPv4:

		if req.Link != "" {
			procPath = path.Join(path.Join(path.Join(path.Join(SysNetPath, SysNetPathIPv4), "conf"), req.Link), req.Property)
		} else {
			procPath = path.Join(path.Join(SysNetPath, SysNetPathIPv4), req.Property)
		}
		break
	case SysNetPathIPv6:

		if req.Link != "" {
			procPath = path.Join(path.Join(path.Join(path.Join(SysNetPath, SysNetPathIPv6), "conf"), req.Link), req.Property)
		} else {
			procPath = path.Join(path.Join(SysNetPath, SysNetPathIPv6), req.Property)
		}
		break
	default:
		return "", errors.New("Path not found")
	}

	return procPath, nil
}

func (req *SysNet) GetSysNet(rw http.ResponseWriter) error {
	path, err := req.GetSysNetPath()
	if err != nil {
		return err
	}

	line, err := share.ReadOneLineFile(path)
	if err != nil {
		return err
	}

	property := SysNet{
		Path: req.Path,
		Property: req.Property,
		Value: line,
		Link:req.Link,
	}

	return share.JsonResponse(property, rw)
}

func (req *SysNet) SetSysNet(rw http.ResponseWriter) error {
	path, err := req.GetSysNetPath()
	if err != nil {
		return err
	}

	return share.WriteOneLineFile(path, req.Value)
}
