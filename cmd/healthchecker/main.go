package main

import (
	"flag"
	"fmt"
	"github.com/bpalermo/new-relic-trace-proxy/internal/constants"
	"log"
	"net/http"
	"os"
)

var (
	url = flag.String("url", fmt.Sprintf("http://localhost:%d%s", constants.DefaultPort, constants.HealthCheckPath), "server url to health check")
)

func makeRequest(url *string) error {
	log.Printf("Health checking: %s", *url)
	if _, err := http.Get(*url); err != nil {
		return err
	}
	return nil
}

func main() {
	if err := makeRequest(url); err != nil {
		os.Exit(1)
	}
}
