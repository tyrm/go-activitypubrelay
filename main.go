package main

import (
	"github.com/juju/loggo"
	"github.com/juju/loggo/loggocolor"
	"os"
	"os/signal"
	"syscall"
)

var logger *loggo.Logger

func main() {
	config := CollectConfig()

	// Init Logging
	newLogger := loggo.GetLogger("main")
	logger = &newLogger
	logger.Infof("Starting LitePub Relay")

	err := loggo.ConfigureLoggers(config.LoggerConfig)
	if err != nil {
		logger.Errorf("Error configuring Logger: %s", err.Error())
		return
	}
	_, err = loggo.ReplaceDefaultWriter(loggocolor.NewWriter(os.Stderr))
	if err != nil {
		logger.Errorf("Error configuring Color Logger: %s", err.Error())
		return
	}

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	nch := make(chan os.Signal)
	signal.Notify(nch, syscall.SIGINT, syscall.SIGTERM)
	logger.Infof("%s", <-nch)
}