// SPDX-License-Identifier: Apache-2.0

package netlink

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/RestGW/api-routerd/cmd/share"

	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

type Address struct {
	Action  string `json:"action"`
	Link    string `json:"link"`
	Address string `json:"address"`
	Label   string `json:"label"`
}

func DecodeAddressJSONRequest(r *http.Request) (Address, error) {
	address := new(Address)

	err := json.NewDecoder(r.Body).Decode(&address)
	if err != nil {
		return *address, err
	}

	return *address, nil
}

func AddAddress(r *http.Request) error {
	address, err := DecodeAddressJSONRequest(r)
	if err != nil {
		log.Errorf("Failed decode Address Json request: %v", err)
		return err
	}

	link, err := netlink.LinkByName(address.Link)
	if err != nil {
		log.Errorf("Failed to find link %v: %s", err, address.Link)
		return err
	}

	addr, err := netlink.ParseAddr(address.Address)
	if err != nil {
		log.Errorf("Failed to parse address %v: %s", err, address.Address)
		return err
	}

	err = netlink.AddrAdd(link, addr)
	if err != nil {
		log.Errorf("Failed to add Address %v to link %s: %s", err, addr, link)
		return err
	}

	return nil
}

func DelAddress(r *http.Request) error {
	address, err := DecodeAddressJSONRequest(r)
	if err != nil {
		log.Errorf("Failed decode Address Json request: %s", err)
		return err
	}

	link, err := netlink.LinkByName(address.Link)
	if err != nil {
		log.Errorf("Failed to find link %v: %s", err, address.Link)
		return err
	}

	addr, err := netlink.ParseAddr(address.Address)
	if err != nil {
		log.Errorf("Failed to parse address %v: %s", err, addr)
		return err
	}

	err = netlink.AddrDel(link, addr)
	if err != nil {
		log.Errorf("Failed to add address %v: %s, %s", err, addr, link)
		return err
	}

	return nil
}

func GetAddress(rw http.ResponseWriter, link string) error {
	l := strings.TrimSpace(link)
	if l != "" {
		link, err := netlink.LinkByName(l)
		if err != nil {
			log.Errorf("Failed to get link %v: %s", l, err)
			return err
		}

		addrs, err := netlink.AddrList(link, netlink.FAMILY_ALL)
		if err != nil {
			log.Errorf("Could not get addresses for link %s: %v", l, err)
			return err
		}

		return share.JsonResponse(addrs, rw)
	}

	addrs, err := netlink.AddrList(nil, netlink.FAMILY_ALL)
	if err != nil {
		log.Errorf("Could not get addresses: %s", err)
		return err
	}

	return share.JsonResponse(addrs, rw)
}
