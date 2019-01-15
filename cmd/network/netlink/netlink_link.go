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
	Mode    string   `json:"mode"`
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
		log.Errorf("Failed to find bridge link %s: %v", req.Link, err)
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
			log.Errorf("Failed to find slave link %s: %v", n, err)
			continue
		}

		err = netlink.LinkSetMaster(link, br)
		if err != nil {
			log.Errorf("Failed to set link %s master device %s: %v", n, req.Link, err)
		}
	}

	return nil
}

func (req *Link) LinkCreateBridge() error {
	_, err := netlink.LinkByName(req.Link)
	if err == nil {
		log.Infof("Bridge link %s exists. Using the bridge", req.Link)
	} else {

		bridge := &netlink.Bridge{
			LinkAttrs: netlink.LinkAttrs{
				Name: req.Link,
			},
		}
		err = netlink.LinkAdd(bridge)
		if err != nil {
			log.Errorf("Failed to create bridge %s: %v", req.Link, err)
			return err
		}

		log.Debugf("Successfully create bridge link: %s", req.Link)
	}

	return req.LinkSetMasterBridge()
}

func (req *Link) LinkSetMasterBond() error {
	bond, err := netlink.LinkByName(req.Link)
	if err != nil {
		log.Errorf("Failed to find bond link %s: %v", req.Link, err)
		return err
	}

	for _, n := range req.Enslave {
		link, err := netlink.LinkByName(n)
		if err != nil {
			log.Errorf("Failed to find slave link %s: %v", n, err)
			continue
		}

		err = netlink.LinkSetBondSlave(link, &netlink.Bond{LinkAttrs: *bond.Attrs()})
		if err != nil {
			log.Errorf("Failed to set link %s master device %s: %v", n, req.Link, err)
		}
	}

	return nil
}

func (req *Link) LinkCreateBond() error {
	_, err := netlink.LinkByName(req.Link)
	if err == nil {
		log.Infof("Bond link %s exists. Using the bond", req.Link)
	} else {

		bond := netlink.NewLinkBond(
			netlink.LinkAttrs{
				Name: req.Link,
			},
		)

		bond.Mode = netlink.StringToBondModeMap[req.Mode]
		err = netlink.LinkAdd(bond)
		if err != nil {
			log.Errorf("Failed to create bond %s: %v", req.Link, err)
			return err
		}

		log.Debugf("Successfully create bond link: %s", req.Link)
	}

	return req.LinkSetMasterBond()
}

func LinkSetUp(link string) error {
	l, err := netlink.LinkByName(link)
	if err != nil {
		log.Errorf("Failed to find link %s: %v", link, err)
		return err
	}

	err = netlink.LinkSetUp(l)
	if err != nil {
		log.Errorf("Failed to set link %s up: %v", l, err)
		return err
	}

	return nil
}

func LinkSetDown(link string) error {
	l, err := netlink.LinkByName(link)
	if err != nil {
		log.Errorf("Failed to find link %s: %v", link, err)
		return err
	}

	err = netlink.LinkSetDown(l)
	if err != nil {
		log.Errorf("Failed to set link down %s: %v", l, err)
		return err
	}

	return nil
}

func LinkSetMTU(link string, mtu int) error {
	l, err := netlink.LinkByName(link)
	if err != nil {
		log.Errorf("Failed to find link %s: %v", link, err)
		return err
	}

	err = netlink.LinkSetMTU(l, mtu)
	if err != nil {
		log.Errorf("Failed to set link %s MTU %d: %v", link, mtu, err)
		return err
	}

	return nil
}

func SetLink(r *http.Request) error {
	req, err := DecodeLinkJSONRequest(r)
	if err != nil {
		log.Errorf("Failed to decode JSON: %v", err)
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
			log.Errorf("Failed to parse received link %s MTU %s: %v", req.Link, req.MTU, err)
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
			log.Errorf("Failed to find link %s: %v", link, err)
			return err
		}

		return share.JsonResponse(l, rw)

	}

	links, err := netlink.LinkList()
	if err != nil {
		return err
	}

	return share.JsonResponse(links, rw)
}

func DeleteLink(r *http.Request) error {
	req, err := DecodeLinkJSONRequest(r)
	if err != nil {
		log.Errorf("Failed to decode JSON: %v", err)
		return err
	}

	l, err := netlink.LinkByName(req.Link)
	if err != nil {
		log.Errorf("Failed to find link %s: %v", req.Link, err)
		return err
	}

	err = netlink.LinkDel(l)
	if err != nil {
		log.Errorf("Failed to delete link %s up: %v", l, err)
		return err
	}

	return nil
}

func CreateLink(r *http.Request) error {
	req, err := DecodeLinkJSONRequest(r)
	if err != nil {
		log.Errorf("Failed to decode JSON: %v", err)
		return err
	}

	switch req.Action {
	case "add-link-bridge":
		return req.LinkCreateBridge()
	case "add-link-bond":
		return req.LinkCreateBond()
	}

	return nil
}
