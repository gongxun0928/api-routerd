// SPDX-License-Identifier: Apache-2.0

package timedate

import (
	"api-routerd/cmd/share"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

var TimeInfo = map[string]string{
	"Timezone":        "",
	"LocalRTC":        "",
	"CanNTP":          "",
	"NTP":             "",
	"NTPSynchronized": "",
	"TimeUSec":        "",
	"RTCTimeUSec":     "",
}

var TimeDateMethod = map[string]string{
	"SetTime":       "",
	"SetTimezone":   "",
	"SetLocalRTC":   "",
	"SetNTP":        "",
	"ListTimezones": "",
}

type TimeDate struct {
	Property string `json:"property"`
	Value    string `json:"value"`
}

func (t *TimeDate) SetTimeDate() error {
	conn, err := share.GetSystemBusPrivateConn()
	if err != nil {
		log.Error("Failed to get systemd bus connection: ", err)
		return err
	}
	defer conn.Close()

	_, k := TimeDateMethod[t.Property]
	if !k {
		return fmt.Errorf("Failed to set timedate:  %s not found", t.Property)
	}

	h := conn.Object("org.freedesktop.timedate1", "/org/freedesktop/timedate1")

	if t.Value == "SetNTP" {

		b, err := share.ParseBool(t.Value)
		if err != nil {
			return err
		}

		r := h.Call("org.freedesktop.timedate1."+t.Property, 0, b, false).Err
		if r != nil {
			log.Errorf("Failed to set SetNTP: %s", r)
			return r
		}
	} else {

		r := h.Call("org.freedesktop.timedate1."+t.Property, 0, t.Value, false).Err
		if r != nil {
			log.Errorf("Failed to set timedate property: %s", r)
			return r
		}
	}

	return nil
}

func GetTimeDate(rw http.ResponseWriter, property string) error {
	conn, err := share.GetSystemBusPrivateConn()
	if err != nil {
		log.Error("Failed to get dbus connection: ", err)
		return err
	}
	defer conn.Close()

	h := conn.Object("org.freedesktop.timedate1", "/org/freedesktop/timedate1")
	for k := range TimeInfo {
		p, perr := h.GetProperty("org.freedesktop.timedate1." + k)
		if perr != nil {
			log.Errorf("Failed to get org.freedesktop.timedate1.%s", k)
			continue
		}

		switch k {
		case "Timezone":
			v, b := p.Value().(string)
			if !b {
				continue
			}

			TimeInfo[k] = v
			break
		case "LocalRTC":
			v, b := p.Value().(bool)
			if !b {
				continue
			}

			TimeInfo[k] = strconv.FormatBool(v)

			break

		case "CanNTP":
			v, b := p.Value().(bool)
			if !b {
				continue
			}

			TimeInfo[k] = strconv.FormatBool(v)

			break
		case "NTP":
			v, b := p.Value().(bool)
			if !b {
				continue
			}

			TimeInfo[k] = strconv.FormatBool(v)

			break
		case "NTPSynchronized":
			v, b := p.Value().(bool)
			if !b {
				continue
			}

			TimeInfo[k] = strconv.FormatBool(v)

			break
		case "TimeUSec":
			v, b := p.Value().(uint64)
			if !b {
				continue
			}

			t := time.Unix(0, int64(v))
			TimeInfo[k] = t.String()

		case "RTCTimeUSec":
			v, b := p.Value().(uint64)
			if !b {
				continue
			}

			t := time.Unix(0, int64(v))

			TimeInfo[k] = t.String()
			break
		}
	}

	if property == "" {
		b, err := json.Marshal(TimeInfo)
		if err != nil {
			return err
		}

		rw.Write(b)
	} else {
		t := TimeDate{Property: property, Value: TimeInfo[property]}
		b, err := json.Marshal(t)
		if err != nil {
			return err
		}

		rw.Write(b)
	}

	return nil
}
