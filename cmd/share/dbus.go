// SPDX-License-Identifier: Apache-2.0

package share

import (
	"github.com/godbus/dbus"
)

func GetSystemBusPrivateConn() (*dbus.Conn, error) {
	conn, err := dbus.SystemBusPrivate()
	if err != nil {
		return nil, err
	}

	if err = conn.Auth(nil); err != nil {
		conn.Close()
		conn = nil
		return conn, err
	}

	if err = conn.Hello(); err != nil {
		conn.Close()
		conn = nil
	}

	return conn, nil
}
