package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	goodPath = "/healthz"
)

func TestMain_makeRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, goodPath, req.URL.String())
		_, _ = rw.Write([]byte(`OK`))
	}))
	defer server.Close()

	t.Run("success", func(t *testing.T) {
		u := fmt.Sprintf("%s%s", server.URL, goodPath)
		err := makeRequest(&u)
		assert.Nil(t, err)
	})

	t.Run("failure", func(t *testing.T) {
		u := fmt.Sprintf("%s%s", "http://localhost", "/")
		err := makeRequest(&u)
		assert.NotNil(t, err)
	})
}
