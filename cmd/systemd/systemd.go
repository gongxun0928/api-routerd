// SPDX-License-Identifier: Apache-2.0

package systemd

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/RestGW/api-routerd/cmd/share"

	sd "github.com/coreos/go-systemd/dbus"
	"github.com/godbus/dbus"
	log "github.com/sirupsen/logrus"
)

type Unit struct {
	Action   string `json:"action"`
	Unit     string `json:"unit"`
	UnitType string `json:"unit_type"`
	Property string `json:"property"`
	Value    string `json:"value"`
}

type Property struct {
	Property string `json:"property"`
	Value    string `json:"value"`
}

type UnitStatus struct {
	Status string `json:"property"`
	Unit   string `json:"unit"`
}

func SystemdProperty(property string) (dbus.Variant, error) {
	conn, err := share.GetSystemBusPrivateConn()
	if err != nil {
		log.Error("Failed to get dbus connection: ", err)
		return dbus.Variant{}, err
	}
	defer conn.Close()

	c := conn.Object("org.freedesktop.systemd1", "/org/freedesktop/systemd1")
	p, perr := c.GetProperty("org.freedesktop.systemd1.Manager." + property)
	if perr != nil {
		log.Errorf("org.freedesktop.systemd1.Manager.%s", property)
		return dbus.Variant{}, errors.New("dbus error")
	}

	if p.Value() == nil {
		return dbus.Variant{}, errors.New("Failed to get property")
	}

	return p, nil
}

func SystemdState(w http.ResponseWriter) error {
	v, err := SystemdProperty("SystemState")
	if err != nil {
		return err
	}

	prop := Property{
		Property: "SystemState",
		Value: v.Value().(string),
	}

	return share.JsonResponse(prop, w)
}

func SystemdVersion(w http.ResponseWriter) error {
	v, err := SystemdProperty("Version")
	if err != nil {
		return err
	}

	prop := Property{
		Property: "Version",
		Value: v.Value().(string),
	}

	return share.JsonResponse(prop, w)
}

func SystemdVirtualization(w http.ResponseWriter) error {
	v, err := SystemdProperty("Virtualization")
	if err != nil {
		return err
	}

	prop := Property{
		Property: "Virtualization",
		Value: v.Value().(string),
	}

	return share.JsonResponse(prop, w)
}

func SystemdArchitecture(w http.ResponseWriter) error {
	v, err := SystemdProperty("Architecture")
	if err != nil {
		return err
	}

	prop := Property{
		Property: "Architecture",
		Value: v.Value().(string),
	}

	return share.JsonResponse(prop, w)
}

func SystemdFeatures(w http.ResponseWriter) error {
	v, err := SystemdProperty("Features")
	if err != nil {
		return err
	}

	prop := Property{
		Property: "Features",
		Value: v.Value().(string),
	}

	return share.JsonResponse(prop, w)
}

func SystemdNFailedUnits(w http.ResponseWriter) error {
	v, err := SystemdProperty("NFailedUnits")
	if err != nil {
		return err
	}

	prop := Property{
		Property: "NFailedUnits",
		Value: fmt.Sprint(v.Value().(uint32)),
	}

	return share.JsonResponse(prop, w)
}

func SystemdNNames(w http.ResponseWriter) error {
	v, err := SystemdProperty("NNames")
	if err != nil {
		return err
	}

	prop := Property{
		Property: "NNames",
		Value: fmt.Sprint(v.Value().(uint32)),
	}

	return share.JsonResponse(prop, w)
}

func ListUnits(w http.ResponseWriter) error {
	conn, err := sd.NewSystemdConnection()
	if err != nil {
		log.Errorf("Failed to get systemd bus connection: %s", err)
		return err
	}
	defer conn.Close()

	units, err := conn.ListUnits()
	if err != nil {
		log.Errorf("Failed ListUnits: %s", err)
		return err
	}

	return share.JsonResponse(units, w)
}

func (u *Unit) StartUnit() error {
	conn, err := sd.NewSystemdConnection()
	if err != nil {
		log.Errorf("Failed to get systemd bus connection: %s", err)
		return err
	}
	defer conn.Close()

	reschan := make(chan string)
	_, err = conn.StartUnit(u.Unit, "replace", reschan)
	if err != nil {
		log.Errorf("Failed to start unit %s: %s", u.Unit, err)
		return err
	}

	return nil
}

func (u *Unit) StopUnit() error {
	conn, err := sd.NewSystemdConnection()
	if err != nil {
		log.Errorf("Failed to get systemd bus connection: %s", err)
		return err
	}
	defer conn.Close()

	reschan := make(chan string)
	_, err = conn.StopUnit(u.Unit, "fail", reschan)
	if err != nil {
		log.Errorf("Failed to stop unit %s: %s", u.Unit, err)
		return err
	}

	return nil
}

func (u *Unit) RestartUnit() error {
	conn, err := sd.NewSystemdConnection()
	if err != nil {
		log.Errorf("Failed to get systemd bus connection: %s", err)
		return err
	}
	defer conn.Close()

	reschan := make(chan string)
	_, err = conn.RestartUnit(u.Unit, "replace", reschan)
	if err != nil {
		log.Errorf("Failed to restart unit %s: %s", u.Unit, err)
		return err
	}

	return nil
}

func (u *Unit) ReloadUnit() error {
	conn, err := sd.NewSystemdConnection()
	if err != nil {
		log.Errorf("Failed to get systemd bus connection: %s", err)
		return err
	}
	defer conn.Close()

	err = conn.Reload()
	if err != nil {
		log.Errorf("Failed to reload unit %s: %s", u.Unit, err)
		return err
	}

	return nil
}

func (u *Unit) KillUnit() error {
	conn, err := sd.NewSystemdConnection()
	if err != nil {
		log.Errorf("Failed to get systemd bus connection: %s", err)
		return err
	}
	defer conn.Close()

	signal, err := strconv.ParseInt(u.Value, 10, 64)
	if err != nil {
		log.Errorf("Failed to parse signal number '%s': %s", u.Value, err)
		return err
	}

	conn.KillUnit(u.Unit, int32(signal))

	return nil
}

func (u *Unit) GetUnitStatus(w http.ResponseWriter) error {
	conn, err := sd.NewSystemdConnection()
	if err != nil {
		log.Errorf("Failed to get systemd bus connection: %s", err)
		return err
	}
	defer conn.Close()

	units, err := conn.ListUnitsByNames([]string{u.Unit})
	if err != nil {
		log.Errorf("Failed get unit '%s' status: %s", u.Unit, err)
		return err
	}

	status := UnitStatus{
		Status: units[0].ActiveState,
		Unit: u.Unit,
	}

	json.NewEncoder(w).Encode(status)

	return nil
}

func (u *Unit) GetUnitProperty(w http.ResponseWriter) error {
	conn, err := sd.NewSystemdConnection()
	if err != nil {
		log.Errorf("Failed to get systemd bus connection: %s", err)
		return err
	}
	defer conn.Close()

	if u.Property != "" {
		p, err := conn.GetServiceProperty(u.Unit, u.Property)
		if err != nil {
			log.Errorf("Failed to get service property: %s", err)
			return err
		}

		switch u.Property {
		case "CPUShares", "LimitNOFILE", "LimitNOFILESoft":
			cpu := strconv.FormatUint(p.Value.Value().(uint64), 10)
			prop := Property{Property: p.Name, Value: cpu}

			return share.JsonResponse(prop, w)
		}
	}

	p, err := conn.GetUnitProperties(u.Unit)
	if err != nil {
		log.Errorf("Failed to get service properties: %s", err)
		return err
	}

	return share.JsonResponse(p, w)
}

func (u *Unit) SetUnitProperty(w http.ResponseWriter) error {
	conn, err := sd.NewSystemdConnection()
	if err != nil {
		log.Errorf("Failed to get systemd bus connection: %s", err)
		return err
	}
	defer conn.Close()

	switch u.Property {
	case "CPUShares":
		n, err := strconv.ParseInt(u.Value, 10, 64)
		if err != nil {
			log.Errorf("Failed to parse CPUShares: %s", err)
			return err
		}

		err = conn.SetUnitProperties(u.Unit, true, sd.Property{"CPUShares", dbus.MakeVariant(uint64(n))})
		if err != nil {
			log.Errorf("Failed to set CPUShares %s: %s", u.Value, err)
			return err
		}
		break
	}

	return nil
}

func (u *Unit) GetUnitTypeProperty(w http.ResponseWriter) error {
	conn, err := sd.NewSystemdConnection()
	if err != nil {
		log.Errorf("Failed to get systemd bus connection: %s", err)
		return err
	}
	defer conn.Close()

	p, err := conn.GetUnitTypeProperties(u.Unit, u.UnitType)
	if err != nil {
		log.Errorf("Failed to get unit type properties: %s", err)
		return err
	}

	return share.JsonResponse(p, w)
}
