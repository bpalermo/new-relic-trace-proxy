package server

import (
	"github.com/bpalermo/new-relic-trace-proxy/pkg/handler"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Server struct {
	*http.Server
}

func NewServer(adr *string, apiKey *string, hostOverride *string, logger *logrus.Logger) Server {
	h := handler.New(apiKey, hostOverride, logger)
	return Server{
		&http.Server{
			Addr:    *adr,
			Handler: h,
		},
	}
}

func (s Server) Start() error {
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}
