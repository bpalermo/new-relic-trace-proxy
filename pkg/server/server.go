package server

import (
	"context"
	"fmt"
	"github.com/bpalermo/new-relic-trace-proxy/pkg/handler"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

// Server proxy server
type Server struct {
	server  *http.Server
	logger  *logrus.Logger
	healthy *int32
	done    chan bool
	quit    chan os.Signal
}

// Start starts the server
func (s *Server) Start() error {
	s.logger.Info("Server is starting...")
	go func() {
		<-s.quit
		s.logger.Info("Server is shutting down...")
		atomic.StoreInt32(s.healthy, 0)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		s.server.SetKeepAlivesEnabled(false)
		if err := s.server.Shutdown(ctx); err != nil {
			s.logger.WithError(err).Fatal("Could not gracefully shutdown the server")
		}
		close(s.done)
	}()

	s.logger.WithField("addr", s.server.Addr).Info("Server is ready to handle requests")
	atomic.StoreInt32(s.healthy, 1)

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.WithField("addr", s.server.Addr).
			WithError(err).
			Error("Could not serve")
		return err
	}

	<-s.done
	s.logger.Println("Server stopped")
	return nil
}

// NewServer creates a server and handler
func NewServer(port uint, apiKey *string, hostOverride *string, healthy *int32, logger *logrus.Logger) Server {
	mux := handler.NewMux(apiKey, hostOverride, healthy, logger)
	addr := fmt.Sprintf(":%d", port)
	s := Server{
		logger: logger,
		server: &http.Server{
			Addr:         addr,
			Handler:      mux,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Second,
		},
		healthy: healthy,
	}

	s.done = make(chan bool)
	s.quit = make(chan os.Signal, 1)

	// register for interupt (Ctrl+C) and SIGTERM (docker)
	signal.Notify(s.quit,
		os.Interrupt,
		syscall.SIGTERM,
	)

	return s
}
