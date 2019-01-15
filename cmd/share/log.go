// SPDX-License-Identifier: Apache-2.0

package share

import (
	"os"
	"path"

	"github.com/sirupsen/logrus"
)

const (
	defaultLogDir  = "/var/log/api-router"
	defaultLogFile = "api-router.log"
)

//InitLog inits the logger
func InitLog() error {
	log := logrus.New()

	log.Out = os.Stderr
	log.Level = logrus.DebugLevel

	logDir := defaultLogDir

	err := CreateDirectory(logDir, 0644)
	if err != nil {
		log.Errorf("Failed to create log directory. path: %s, err: %s", logDir, err)
		return err
	}

	logFile := path.Join(logDir, defaultLogFile)
	f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Errorf("Failed to create log file. path: %s, err: %s", logFile, err)
		return err
	}

	log.SetOutput(f)
	log.SetReportCaller(true)
	log.Info("Starting API Router")

	return nil
}
