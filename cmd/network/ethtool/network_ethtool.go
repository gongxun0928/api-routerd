// SPDX-License-Identifier: Apache-2.0

package ethtool

import (
	"errors"
	"net/http"
	"unsafe"

	"github.com/RestGW/api-routerd/cmd/share"

	"github.com/safchain/ethtool"
	log "github.com/sirupsen/logrus"
)

type Ethtool struct {
	Action string `json:"action"`
	Link   string `json:"link"`
}

func (req *Ethtool) GetEthTool(rw http.ResponseWriter) error {
	link := share.LinkExists(req.Link)
	if !link {
		log.Errorf("Failed to get link: %s", req.Link)
		return errors.New("Link not found")
	}

	e, err := ethtool.NewEthtool()
	if err != nil {
		log.Errorf("Failed to init ethtool for link %s: %s", err, req.Link)
		return err
	}
	defer e.Close()

	switch req.Action {
	case "get-link-stat":
		stats, err := e.Stats(req.Link)
		if err != nil {
			log.Errorf("Failed to get ethtool statitics for link %s: %s", err, req.Link)
			return err
		}

		return share.JsonResponse(stats, rw)

	case "get-link-features":

		features, err := e.Features(req.Link)
		if err != nil {
			log.Errorf("Failed to get ethtool features for link %s: %s", err, req.Link)
			return err
		}

		return share.JsonResponse(features, rw)

	case "get-link-bus":

		bus, err := e.BusInfo(req.Link)
		if err != nil {
			log.Errorf("Failed to get ethtool bus for link %s: %s", err, req.Link)
			return err
		}

		b := struct {
			Bus string
		}{
			bus,
		}

		return share.JsonResponse(b, rw)

	case "get-link-driver-name":

		driver, err := e.DriverName(req.Link)
		if err != nil {
			log.Errorf("Failed to get ethtool driver name for link %s: %s", err, req.Link)
			return err
		}

		d := struct {
			Driver string
		}{
			driver,
		}

		return share.JsonResponse(d, rw)

	case "get-link-driver-info":

		e, err := NewEthTool()
		if err != nil {
			log.Errorf("Failed to init ethtool for link %s: %s", err, req.Link)
			return err
		}
		defer e.Close()

		drvinfo := EthtoolDrvInfo{
			Cmd: ETHTOOL_GDRVINFO,
		}

		err = e.Ioctl(req.Link, uintptr(unsafe.Pointer(&drvinfo)))
		if err != nil {
			return err
		}

		return share.JsonResponse(drvinfo, rw)

	case "get-link-permaddr":

		permaddr, err := e.PermAddr(req.Link)
		if err != nil {
			log.Errorf("Failed to get ethtool Perm Addr for link %s: %s", err, req.Link)
			return err
		}

		p := struct {
			PermAddr string
		}{
			permaddr,
		}

		return share.JsonResponse(p, rw)

	case "get-link-eeprom":

		eeprom, err := e.ModuleEepromHex(req.Link)
		if err != nil {
			log.Errorf("Failed to get ethtool eeprom for link %s: %s", err, req.Link)
			return err
		}

		e := struct {
			ModuleEeprom string
		}{
			eeprom,
		}

		return share.JsonResponse(e, rw)

	case "get-link-msglvl":

		msglvl, err := e.MsglvlGet(req.Link)
		if err != nil {
			log.Errorf("Failed to get ethtool msglvl for link %s: %s", err, req.Link)
			return err
		}

		g := struct {
			ModuleMsglv uint32
		}{
			msglvl,
		}

		return share.JsonResponse(g, rw)

	case "get-link-mapped":

		mapped, err := e.CmdGetMapped(req.Link)
		if err != nil {
			log.Errorf("Failed to get ethtool msglvl for link %s: %s", err, req.Link)
			return err
		}

		return share.JsonResponse(mapped, rw)
	}

	return nil
}
