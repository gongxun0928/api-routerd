// SPDX-License-Identifier: Apache-2.0

package main

import (
	"flag"
	"os"
	"path"
	"runtime"

	"github.com/RestGW/api-routerd/cmd/share"

	"github.com/RestGW/api-routerd/cmd/router"

	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
)

// App Version
const (
	Version  = "0.1"
	ConfPath = "/etc/api-routerd"
	ConfFile = "api-routerd.toml"
	TLSCert  = "tls/server.crt"
	TLSKey   = "tls/server.key"
)

// flag
var IPFlag string
var PortFlag string

type tomlConfig struct {
	Server Network `toml:"Network"`
}

type Network struct {
	IPAddress string
	Port      string
}

func init() {
	const (
		defaultIP   = "0.0.0.0"
		defaultPort = "8080"
	)

	flag.StringVar(&IPFlag, "ip", defaultIP, "The server IP address.")
	flag.StringVar(&PortFlag, "port", defaultPort, "The server port.")
}

func InitConf() (tomlConfig, error) {
	var conf tomlConfig

	confFile := path.Join(ConfPath, ConfFile)
	_, err := toml.DecodeFile(confFile, &conf)
	if err != nil {
		log.Errorf("Fail to read conf file '%s': %v", ConfPath, err)
		return conf, err
	}

	_, err = share.ParseIP(conf.Server.IPAddress)
	if err != nil {
		log.Errorf("Failed to parse IPAddress=%s, %s", conf.Server.IPAddress, conf.Server.Port)
		return conf, err
	}

	_, err = share.ParsePort(conf.Server.Port)
	if err != nil {
		log.Errorf("Failed to parse Conf file Port=%s", conf.Server.Port)
		return conf, err
	}

	log.Debugf("Conf file: Parsed IPAddress=%s and Port=%s", conf.Server.IPAddress, conf.Server.Port)

	return conf, nil
}

func main() {
	share.InitLog()

	conf, err := InitConf()
	if err != nil {
		flag.Parse()
	} else {
		IPFlag = conf.Server.IPAddress
		PortFlag = conf.Server.Port
	}

	log.Infof("api-routerd: v%s (built %s)", Version, runtime.Version())
	log.Infof("Start Server at %s:%s", IPFlag, PortFlag)

	err = router.StartRouter(IPFlag, PortFlag, path.Join(ConfPath, TLSCert), path.Join(ConfPath, TLSKey))
	if err != nil {
		log.Fatalf("Failed to init api-routerd: %s", err)
		os.Exit(1)
	}
}
