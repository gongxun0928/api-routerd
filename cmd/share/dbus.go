// SPDX-License-Identifier: Apache-2.0

package share

import (
	"os"
	"strconv"

	"github.com/godbus/dbus"
)

func GetSystemBusPrivateConn() (*dbus.Conn, error) {
	conn, err := dbus.SystemBusPrivate()
	if err != nil {
		return nil, err
	}

	methods := []dbus.Auth{dbus.AuthExternal(strconv.Itoa(os.Getuid()))}

	err = conn.Auth(methods)
	if err != nil {
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
