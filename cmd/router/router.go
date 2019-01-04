// SPDX-License-Identifier: Apache-2.0

package router

import (
	"api-routerd/cmd/network"
	"api-routerd/cmd/proc"
	"api-routerd/cmd/share"
	"api-routerd/cmd/system"
	"api-routerd/cmd/systemd"
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/coreos/go-systemd/activation"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func StartRouter(ip string, port string, tlsCertPath string, tlsKeyPath string) error {
	var srv http.Server

	router := mux.NewRouter()

	// Register services
	network.RegisterRouterNetwork(router)
	proc.RegisterRouterProc(router)
	systemd.RegisterRouterSystemd(router)
	system.RegisterRouterSystem(router)

	// Authenticate users
	amw, err := InitAuthMiddleware()
	if err != nil {
		log.Fatalf("Faild to init auth DB existing: %s", err)
		return fmt.Errorf("Failed to init Auth DB: %s", err)
	}

	router.Use(amw.AuthMiddleware)

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		sig := <-gracefulStop

		log.Printf("Received signal: %+v", sig)
		log.Println("Shutting down api-routerd ...")

		err := srv.Shutdown(context.Background())
		if err != nil {
			log.Errorf("Failed to shutdown server gracefuly: %s", err)
		}

		os.Exit(0)
	}()

	// socket activation
	listeners, err := activation.Listeners()
	if err != nil {
		log.Infof("Failed to retrieve listeners: %s", err)
	}

	if share.PathExists(tlsCertPath) && share.PathExists(tlsKeyPath) {
		cfg := &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: false,
		}
		srv = http.Server{
			Addr:         ip + ":" + port,
			Handler:      router,
			TLSConfig:    cfg,
			TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
		}

		log.Info("Starting api-routerd in TLS mode")

		if len(listeners) <= 0 {
			log.Fatal(srv.ListenAndServeTLS(tlsCertPath, tlsKeyPath))
		} else {
			log.Fatal(srv.ServeTLS(listeners[0], tlsCertPath, tlsKeyPath))
		}
	} else {
		srv = http.Server{
			Addr:    ip + ":" + port,
			Handler: router,
		}
		log.Info("Starting api-routerd in plain text mode")

		if len(listeners) <= 0 {
			log.Fatal(srv.ListenAndServe())
		} else {
			log.Fatal(srv.Serve(listeners[0]))
		}
	}

	return nil
}
