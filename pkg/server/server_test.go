package server

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestServer_NewServer(t *testing.T) {
	k := "fake"
	h := "fake"
	adr := ":45566"

	srv := NewServer(&adr, &k, &h, logrus.New())
	assert.NotNil(t, srv)
}

func TestServer_Start(t *testing.T) {
	k := "fake"
	h := "fake"
	adr := ":45566"

	srv := NewServer(&adr, &k, &h, logrus.New())
	go func() {
		time.Sleep(1 * time.Second)
		_ = srv.Shutdown(context.Background())
	}()

	err := srv.Start()
	if err != nil {
		t.Error("unexpected error:", err)
	}
}
