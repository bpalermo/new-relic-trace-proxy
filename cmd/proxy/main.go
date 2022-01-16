package main

import (
	"fmt"
	"github.com/bpalermo/new-relic-trace-proxy/internal/constants"
	"github.com/bpalermo/new-relic-trace-proxy/pkg/server"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

var (
	logger *logrus.Logger

	debug        = kingpin.Flag("verbose", "Enable additional logging, implies all the logger-* options").Short('v').Bool()
	port         = kingpin.Flag("port", "Port to serve HTTP on").Default(fmt.Sprintf("%d", constants.DefaultPort)).Uint()
	hostOverride = kingpin.Flag("host", "Host to proxy to").Default("").String()
	apiKey       = kingpin.Flag("apiKey", "New Relic API key").String()
)

func init() {
	logger = logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)
}

func main() {
	kingpin.Parse()
	if *debug {
		logger.SetLevel(logrus.DebugLevel)
	}
	logger.Fatal(server.NewServer(*port, apiKey, hostOverride, logger).Start())
}
