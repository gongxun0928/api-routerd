// SPDX-License-Identifier: Apache-2.0

package machine

import (
	"fmt"

	"github.com/RestGW/api-routerd/cmd/share"
	"github.com/godbus/dbus"
)

const (
	dbusInterface = "org.freedesktop.machine1.Manager"
	dbusPath      = "/org/freedesktop/machine1"
)

type Conn struct {
	conn   *dbus.Conn
	object dbus.BusObject
}

func NewMachineConnection() (*Conn, error) {
	c := new(Conn)
	conn, err := share.GetSystemBusPrivateConn()
	if err != nil {
		return nil, err
	}

	c.conn = conn
	c.object = c.conn.Object("org.freedesktop.machine1", dbus.ObjectPath(dbusPath))

	return c, nil
}

func (c *Conn) Close() {
	c.conn.Close()
}

func GetMachineOSRelease(machine string) (dbus.ObjectPath, error) {
	conn, err := NewMachineConnection()
	if err != nil {
		return "", fmt.Errorf("Failed to get systemd bus connection: %s", err)
	}
	defer conn.Close()

	r := conn.object.Call(fmt.Sprintf("%s.%s", dbusInterface, "GetMachineOSRelease"), 0, machine)
	if r != nil {
		fmt.Println(r)
		return "", fmt.Errorf("Failed to get machine release information: %s", r.Err)
	}

	path, typeErr := r.Body[0].(dbus.ObjectPath)
	if !typeErr {
		return "", fmt.Errorf("unable to convert dbus response '%v' to dbus.ObjectPath", r.Body[0])
	}

	return path, nil
}

func CloneImage(image string, newImage string) error {
	conn, err := NewMachineConnection()
	if err != nil {
		return fmt.Errorf("Failed to get systemd bus connection: %s", err)
	}
	defer conn.Close()

	r := conn.object.Call(fmt.Sprintf("%s.%s", dbusInterface, "CloneImage"), 0, image, newImage, false)
	if r != nil {
		fmt.Println(r)
		return fmt.Errorf("Failed to clone image: %s", r.Err)
	}

	return nil
}

func RenameImage(image string, newImage string) error {
	conn, err := NewMachineConnection()
	if err != nil {
		return fmt.Errorf("Failed to get systemd bus connection: %s", err)
	}
	defer conn.Close()

	r := conn.object.Call(fmt.Sprintf("%s.%s", dbusInterface, "RenameImage"), 0, image, newImage)
	if r != nil {
		fmt.Println(r)
		return fmt.Errorf("Failed to Rename Image : %s", r.Err)
	}

	return nil
}

func RemoveImage(image string) error {
	conn, err := NewMachineConnection()
	if err != nil {
		return fmt.Errorf("Failed to get systemd bus connection: %s", err)
	}
	defer conn.Close()

	r := conn.object.Call(fmt.Sprintf("%s.%s", dbusInterface, "RemoveImage"), 0, image)
	if r != nil {
		fmt.Println(r)
		return fmt.Errorf("Failed to Remove Image : %s", r.Err)
	}

	return nil
}
