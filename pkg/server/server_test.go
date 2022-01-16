package server

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"syscall"
	"testing"
	"time"
)

func TestServer_NewServer(t *testing.T) {
	k := "fake"
	h := "fake"
	var healthy int32 = 1

	srv := NewServer(45566, &k, &h, &healthy, logrus.New())
	assert.NotNil(t, srv)
}

func TestServer_Start(t *testing.T) {
	logger := logrus.New()
	k := "fake"
	h := "fake"
	var healthy int32 = 1

	t.Run("success", func(t *testing.T) {
		srv := NewServer(4566, &k, &h, &healthy, logger)
		go func() {
			time.Sleep(2 * time.Second)
			srv.quit <- syscall.SIGTERM
		}()
		err := srv.Start()
		assert.Nil(t, err)
	})

	t.Run("failure with bad port", func(t *testing.T) {
		srv := NewServer(400566, &k, &h, &healthy, logger)
		err := srv.Start()
		assert.NotNil(t, err)
	})
}
