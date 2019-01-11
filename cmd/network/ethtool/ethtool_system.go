// SPDX-License-Identifier: Apache-2.0

package ethtool

import (
	"syscall"
	"unsafe"
)

const (
	IFNAMSIZ = 16
)

const (
	SIOCETHTOOL = 0x8946
)

const (
	ETHTOOL_GDRVINFO = 0x00000003
)

type EthTool struct {
	fd int
}

type ifreq struct {
	ifr_name [IFNAMSIZ]byte
	ifr_data uintptr
}

type EthtoolDrvInfo struct {
	Cmd         uint32   `json:"cmd"`
	Driver      [32]byte `json:"driver"`
	Version     [32]byte `json:"version"`
	FwVersion   [32]byte `json:"fw_version"`
	BusInfo     [32]byte `json:"bus_info"`
	EromVersion [32]byte `json:"erom_version"`
	Reserved2   [12]byte `json:"reserved2"`
	NPrivFlags  uint32   `json:"n_priv_flags"`
	NStats      uint32   `json:"n_stats"`
	TestinfoLen uint32   `json:"testinfo_len"`
	EedumpLen   uint32   `json:"eedump_len"`
	RegdumpLen  uint32   `json:"regdump_len"`
}

var manager *EthTool

// SocketIoctlFd returns a new fd
func (e *EthTool) SocketIoctlFd() error {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM|syscall.SOCK_CLOEXEC, syscall.IPPROTO_IP)
	if err != nil {
		fd, err := syscall.Socket(syscall.AF_NETLINK, syscall.SOCK_RAW|syscall.SOCK_CLOEXEC, syscall.NETLINK_GENERIC)
		if err != nil {
			return err
		}
		e.fd = fd

		return nil
	}

	e.fd = fd

	return nil
}

func (e *EthTool) Ioctl(intf string, data uintptr) error {
	var name [IFNAMSIZ]byte
	copy(name[:], []byte(intf))

	ifr := ifreq{
		ifr_name: name,
		ifr_data: data,
	}

	_, _, ep := syscall.Syscall(syscall.SYS_IOCTL, uintptr(e.fd), SIOCETHTOOL, uintptr(unsafe.Pointer(&ifr)))
	if ep != 0 {
		return syscall.Errno(ep)
	}

	return nil
}

func (e *EthTool) EthtoolConnect() error {
	if e.fd < 1 {
		err := e.SocketIoctlFd()
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *EthTool) Close() {
	syscall.Close(e.fd)
	e.fd = -1
}

func NewEthTool() (*EthTool, error) {
	if manager == nil {
		manager = new(EthTool)
		manager.fd = -1
	}

	err := manager.EthtoolConnect()
	if err != nil {
		return nil, err
	}

	return manager, nil
}
