package server

import (
	"fmt"
	"github.com/bpalermo/new-relic-trace-proxy/pkg/handler"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Server struct {
	*http.Server
}

var (
	healthy int32
)

func NewServer(port uint, apiKey *string, hostOverride *string, logger *logrus.Logger) Server {
	h := handler.New(apiKey, hostOverride, &healthy, logger)
	addr := fmt.Sprintf(":%d", port)
	return Server{
		&http.Server{
			Addr:         addr,
			Handler:      h,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Second,
		},
	}
}

func (s Server) Start() error {
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}
