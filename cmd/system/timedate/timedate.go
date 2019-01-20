// SPDX-License-Identifier: Apache-2.0

package timedate

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/RestGW/api-routerd/cmd/share"

	log "github.com/sirupsen/logrus"
)

var timeInfo = map[string]string{
	"Timezone":        "",
	"LocalRTC":        "",
	"CanNTP":          "",
	"NTP":             "",
	"NTPSynchronized": "",
	"TimeUSec":        "",
	"RTCTimeUSec":     "",
}

var timeDateMethod = map[string]string{
	"SetTime":       "",
	"SetTimezone":   "",
	"SetLocalRTC":   "",
	"SetNTP":        "",
	"ListTimezones": "",
}

// TimeDate JSON message
type TimeDate struct {
	Property string `json:"property"`
	Value    string `json:"value"`
}

// SetTimeDate set timedate property
func (t *TimeDate) SetTimeDate() error {
	conn, err := NewConn()
	if err != nil {
		log.Errorf("Failed to get systemd bus connection: %v", err)
		return err
	}
	defer conn.Close()

	_, k := timeDateMethod[t.Property]
	if !k {
		return fmt.Errorf("Failed to set timedate:  %s not found", t.Property)
	}

	err = conn.SetTimeDate(t.Property, t.Value)
	if err != nil {
		log.Errorf("Failed to set timedate property: %s", err)
		return err
	}

	return nil
}

// GetTimeDate gets property from timedated
func GetTimeDate(rw http.ResponseWriter, property string) error {
	conn, err := NewConn()
	if err != nil {
		log.Errorf("Failed to get dbus connection: %v", err)
		return err
	}
	defer conn.Close()

	for k := range timeInfo {
		p, err := conn.GetTimeDate(k)
		if err != nil {
			log.Errorf("Failed to get %s", k)
			continue
		}

		switch k {
		case "Timezone":
			v, b := p.Value().(string)
			if !b {
				continue
			}

			timeInfo[k] = v
			break
		case "LocalRTC":
			v, b := p.Value().(bool)
			if !b {
				continue
			}

			timeInfo[k] = strconv.FormatBool(v)

			break
		case "CanNTP":
			v, b := p.Value().(bool)
			if !b {
				continue
			}

			timeInfo[k] = strconv.FormatBool(v)

			break
		case "NTP":
			v, b := p.Value().(bool)
			if !b {
				continue
			}

			timeInfo[k] = strconv.FormatBool(v)

			break
		case "NTPSynchronized":
			v, b := p.Value().(bool)
			if !b {
				continue
			}

			timeInfo[k] = strconv.FormatBool(v)

			break
		case "TimeUSec":
			v, b := p.Value().(uint64)
			if !b {
				continue
			}

			t := time.Unix(0, int64(v))
			timeInfo[k] = t.String()

		case "RTCTimeUSec":
			v, b := p.Value().(uint64)
			if !b {
				continue
			}

			t := time.Unix(0, int64(v))

			timeInfo[k] = t.String()
			break
		}
	}

	if property == "" {
		return share.JSONResponse(timeInfo, rw)
	}

	t := TimeDate{
		Property: property,
		Value:    timeInfo[property],
	}

	return share.JSONResponse(t, rw)
}
