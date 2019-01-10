// SPDX-License-Identifier: Apache-2.0

package netlink

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/RestGW/api-routerd/cmd/share"

	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

type Link struct {
	Action  string   `json:"action"`
	Link    string   `json:"link"`
	MTU     string   `json:"mtu"`
	Kind    string   `json:"kind"`
	Enslave []string `json:"enslave"`
}

func DecodeLinkJSONRequest(r *http.Request) (Link, error) {
	link := new(Link)

	err := json.NewDecoder(r.Body).Decode(&link)
	if err != nil {
		return *link, err
	}

	return *link, nil
}

func (req *Link) LinkSetMasterBridge() error {
	bridge, err := netlink.LinkByName(req.Link)
	if err != nil {
		log.Errorf("Failed to find bridge link %s: %s", req.Link, err)
		return err
	}

	br, b := bridge.(*netlink.Bridge)
	if !b {
		log.Errorf("Link is not a bridge: %s", req.Link)
		return errors.New("Link is not a bridge")
	}

	for _, n := range req.Enslave {
		link, err := netlink.LinkByName(n)
		if err != nil {
			log.Errorf("Failed to find slave link %s: %s", n, err)
			continue
		}

		err = netlink.LinkSetMaster(link, br)
		if err != nil {
			log.Errorf("Failed to set link %s master device %s: %s", n, req.Link, err)
		}
	}

	return nil
}

func (req *Link) LinkCreateBridge() error {
	_, err := netlink.LinkByName(req.Link)
	if err == nil {
		log.Infof("Bridge link %s exists. Using the bridge", req.Link)
	} else {

		bridge := &netlink.Bridge{LinkAttrs: netlink.LinkAttrs{Name: req.Link}}
		err = netlink.LinkAdd(bridge)
		if err != nil {
			log.Errorf("Failed to create bride %s: %s", req.Link, err)
			return err
		}

		log.Debugf("Successfully create bridge link: %s", req.Link)
	}

	return req.LinkSetMasterBridge()
}

func LinkSetUp(link string) error {
	l, err := netlink.LinkByName(link)
	if err != nil {
		log.Errorf("Failed to find link %s: %s", link, err)
		return err
	}

	err = netlink.LinkSetUp(l)
	if err != nil {
		log.Errorf("Failed to set link %s up: %s", l, err)
		return err
	}

	return nil
}

func LinkSetDown(link string) error {
	l, err := netlink.LinkByName(link)
	if err != nil {
		log.Errorf("Failed to find link %s: %s", link, err)
		return err
	}

	err = netlink.LinkSetDown(l)
	if err != nil {
		log.Errorf("Failed to set link down %s: %s", l, err)
		return err
	}

	return nil
}

func LinkSetMTU(link string, mtu int) error {
	l, err := netlink.LinkByName(link)
	if err != nil {
		log.Errorf("Failed to find link %s: %s", link, err)
		return err
	}

	err = netlink.LinkSetMTU(l, mtu)
	if err != nil {
		log.Errorf("Failed to set link %s MTU %d: %s", link, mtu, err)
		return err
	}

	return nil
}

func SetLink(r *http.Request) error {
	req, err := DecodeLinkJSONRequest(r)
	if err != nil {
		log.Errorf("Failed to decode JSON: %s", err)
		return err
	}

	switch req.Action {
	case "set-link-up":
		return LinkSetUp(req.Link)
	case "set-link-down":
		return LinkSetDown(req.Link)
	case "set-link-mtu":

		mtu, err := strconv.ParseInt(strings.TrimSpace(req.MTU), 10, 64)
		if err != nil {
			log.Errorf("Failed to parse received link %s MTU %s: %s", req.Link, req.MTU, err)
			return err
		}

		return LinkSetMTU(req.Link, int(mtu))
	}

	return nil
}

func GetLink(rw http.ResponseWriter, r *http.Request, link string) error {
	if link != "" {
		l, err := netlink.LinkByName(link)
		if err != nil {
			log.Errorf("Failed to find link %s: %s", link, err)
			return err
		}

		j, err := json.Marshal(l)
		if err != nil {
			log.Errorf("Failed to encode json linkInfo for link %s: %s", link, err)
			return err
		}

		rw.WriteHeader(http.StatusOK)
		rw.Write(j)

	} else {

		links, err := netlink.LinkList()
		if err != nil {
			return err
		}

		return share.JsonResponse(links, rw)
	}

	return nil
}

func DeleteLink(r *http.Request) error {
	req, err := DecodeLinkJSONRequest(r)
	if err != nil {
		log.Errorf("Failed to decode JSON: %s", err)
		return err
	}

	l, err := netlink.LinkByName(req.Link)
	if err != nil {
		log.Errorf("Failed to find link %s: %s", req.Link, err)
		return err
	}

	err = netlink.LinkDel(l)
	if err != nil {
		log.Errorf("Failed to delete link %s up: %s", l, err)
		return err
	}

	return nil
}

func CreateLink(r *http.Request) error {
	req, err := DecodeLinkJSONRequest(r)
	if err != nil {
		log.Errorf("Failed to decode JSON: %s", err)
		return err
	}

	switch req.Action {
	case "add-link-bridge":
		return req.LinkCreateBridge()
	}

	return nil
}
