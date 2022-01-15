package handler

import (
	"bytes"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockHttp struct {
	StatusCode int
}

func (c *MockHttp) Do(_ *http.Request) (*http.Response, error) {
	r := ioutil.NopCloser(bytes.NewReader([]byte("fake")))

	h := http.Header{}
	h.Set("x-test", "test")

	return &http.Response{
		StatusCode: c.StatusCode,
		Body:       r,
		Header:     h,
	}, nil
}

type MockBadHttp struct {
}

func (c *MockBadHttp) Do(_ *http.Request) (*http.Response, error) {
	return nil, errors.New("bad request")
}

func newHandler(statusCode int, overrideHost bool) (*http.Request, *httptest.ResponseRecorder, *Handler) {
	const defaultHostOverride = "test.com"
	const apiKey = "fake"

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	var hostOverride string
	if overrideHost {
		hostOverride = defaultHostOverride
	}

	r := httptest.NewRequest(http.MethodPost, "/some-url", nil)
	w := httptest.NewRecorder()
	h := &Handler{
		ProxyClient: &ProxyClient{
			Logger: logger,
			Client: &MockHttp{
				StatusCode: statusCode,
			},
			HostOverride: hostOverride,
			ApiKey:       apiKey,
		},
	}

	return r, w, h
}

func newBadHandler() (*http.Request, *httptest.ResponseRecorder, *Handler) {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	r := httptest.NewRequest(http.MethodGet, "/fail", nil)
	w := httptest.NewRecorder()
	return r, w, &Handler{
		Logger: logger,
		ProxyClient: &ProxyClient{
			Logger: logger,
			Client: &MockBadHttp{},
		},
	}
}

func TestHandler_write(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		w := httptest.NewRecorder()
		h := Handler{}
		h.write(w, http.StatusOK, make([]byte, 0))
	})
}

func TestHandler_ServeHTTP(t *testing.T) {

	t.Run("test 200 proxy response", func(t *testing.T) {
		r, w, h := newHandler(http.StatusOK, true)
		h.ServeHTTP(w, r)
	})

	t.Run("test 400 proxy response", func(t *testing.T) {
		r, w, h := newHandler(http.StatusBadRequest, false)
		h.ServeHTTP(w, r)
	})

	t.Run("bad upstream", func(t *testing.T) {
		r, w, h := newBadHandler()
		h.ServeHTTP(w, r)
	})
}

func TestHandler_New(t *testing.T) {
	var healthy int32 = 1
	k := "fake"
	host := "host"
	h := New(&k, &host, &healthy, logrus.New())
	assert.NotNil(t, h)
}

func TestHandler_healthz(t *testing.T) {
	t.Run("healthy", func(t *testing.T) {
		var healthy int32 = 1
		h, r, w := newHealthyzHandler(&healthy)
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusNoContent, w.Result().StatusCode)
	})

	t.Run("unhealthy", func(t *testing.T) {
		var healthy int32 = 0
		h, r, w := newHealthyzHandler(&healthy)
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusServiceUnavailable, w.Result().StatusCode)
	})
}

func newHealthyzHandler(healthy *int32) (http.Handler, *http.Request, *httptest.ResponseRecorder) {
	h := healthz(healthy)
	r := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	return h, r, w
}
