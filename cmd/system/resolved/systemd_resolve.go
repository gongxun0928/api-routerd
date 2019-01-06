// SPDX-License-Identifier: Apache-2.0

package resolved

import (
	"github.com/RestGW/api-routerd/cmd/share"
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/go-ini/ini"
	log "github.com/sirupsen/logrus"
)

const (
	ResolvedConfPath = "/etc/systemd/resolved.conf"
)

type DNSConfig struct {
	DNS         []string `json:"dns"`
	FallbackDNS []string `json:"fallback_dns"`
}

func (d *DNSConfig) WriteResolveConfig() error {
	f, err := os.OpenFile(ResolvedConfPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	conf := "[Resolve]\n"

	dnsConf := "DNS="
	for _, s := range d.DNS {
		dnsConf += s + " "
	}
	conf += dnsConf + "\n"

	fallbackDNS := "FallbackDNS="
	for _, s := range d.FallbackDNS {
		fallbackDNS += s + " "
	}
	conf += fallbackDNS + "\n"

	fmt.Fprintln(w, conf)
	w.Flush()

	return nil
}

func ReadResolveConf() (*DNSConfig, error) {
	cfg, err := ini.Load(ResolvedConfPath)
	if err != nil {
		return nil, err
	}

	conf := new(DNSConfig)

	dns := cfg.Section("Resolve").Key("DNS").String()
	fallbackDNS := cfg.Section("Resolve").Key("FallbackDNS").String()

	conf.DNS = strings.Fields(dns)
	conf.FallbackDNS = strings.Fields(fallbackDNS)

	return conf, nil
}

func GetResolveConf(rw http.ResponseWriter) error {
	conf, err := ReadResolveConf()
	if err != nil {
		return err
	}

	return share.JsonResponse(conf, rw)
}

func UpdateResolveConf(rw http.ResponseWriter, r *http.Request) error {
	dns := DNSConfig{
		DNS:         []string{""},
		FallbackDNS: []string{""},
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("Failed to parse HTTP request: ", err)
		return err
	}

	err = json.Unmarshal([]byte(body), &dns)
	if err != nil {
		log.Error("Failed to Decode HTTP request to json: ", err)
		return err
	}

	conf, err := ReadResolveConf()
	if err != nil {
		return err
	}

	// update DNS
	for _, s := range dns.DNS {
		if share.StringContains(conf.DNS, s) {
			continue
		}

		conf.DNS = append(conf.DNS, s)
	}

	// update fallback
	for _, s := range dns.FallbackDNS {
		if share.StringContains(conf.FallbackDNS, s) {
			continue
		}
		conf.FallbackDNS = append(conf.FallbackDNS, s)
	}

	err = conf.WriteResolveConfig()
	if err != nil {
		log.Errorf("Failed Write to resolv conf: %s", err)
		return err
	}

	return share.JsonResponse(conf, rw)
}

func DeleteResolveConf(rw http.ResponseWriter, r *http.Request) error {
	dns := DNSConfig{
		DNS:         []string{""},
		FallbackDNS: []string{""},
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("Failed to parse HTTP request: ", err)
		return err
	}

	err = json.Unmarshal([]byte(body), &dns)
	if err != nil {
		log.Error("Failed to Decode HTTP request to json: ", err)
		return err
	}

	conf, err := ReadResolveConf()
	if err != nil {
		return err
	}

	// update DNS
	for _, s := range dns.DNS {
		if !share.StringContains(conf.DNS, s) {
			continue
		}

		conf.DNS, _ = share.StringDeleteSlice(conf.DNS, s)
	}

	// update Fallback
	for _, s := range dns.FallbackDNS {
		if !share.StringContains(conf.FallbackDNS, s) {
			continue
		}

		conf.FallbackDNS, _ = share.StringDeleteSlice(conf.FallbackDNS, s)
	}

	err = conf.WriteResolveConfig()
	if err != nil {
		log.Errorf("Failed Write to resolv conf: %s", err)
		return err
	}

	return share.JsonResponse(conf, rw)
}
