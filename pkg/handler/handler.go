package handler

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"sync/atomic"
)

type Handler struct {
	ProxyClient Client
	Logger      *logrus.Logger
	Healthy     *int32
}

func (h *Handler) write(w http.ResponseWriter, status int, body []byte) {
	w.WriteHeader(status)
	_, err := w.Write(body)
	if err != nil {
		h.Logger.WithError(err).Error("could not write response body")
		return
	}
}

func (h *Handler) closeBody(body io.ReadCloser) {
	err := body.Close()
	if err != nil {
		h.Logger.WithError(err).Error("could not close body")
	}
}

func (h *Handler) logError(w http.ResponseWriter, err error, statusCode int, errorMsg string) {
	h.Logger.WithError(err).Error(errorMsg)
	h.write(w, statusCode, []byte(fmt.Sprintf("%v - %v", errorMsg, err.Error())))
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp, err := h.ProxyClient.Do(r)
	if err != nil {
		errorMsg := "unable to proxy request"
		h.logError(w, err, http.StatusBadGateway, errorMsg)
		return
	}
	defer h.closeBody(resp.Body)

	// read response body
	buf := bytes.Buffer{}
	if _, err := io.Copy(&buf, resp.Body); err != nil {
		errorMsg := "error while reading response from upstream"
		h.logError(w, err, http.StatusInternalServerError, errorMsg)
		return
	}

	// copy headers
	for k, vals := range resp.Header {
		for _, v := range vals {
			w.Header().Add(k, v)
		}
	}

	h.write(w, resp.StatusCode, buf.Bytes())
}

func healthz(healthy *int32) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.LoadInt32(healthy) == 1 {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusServiceUnavailable)
	})
}

func New(apiKey *string, hostOverride *string, healthy *int32, logger *logrus.Logger) *http.ServeMux {
	proxy := &Handler{
		Logger: logger,
		ProxyClient: &ProxyClient{
			Logger:       logger,
			Client:       http.DefaultClient,
			HostOverride: *hostOverride,
			ApiKey:       *apiKey,
		},
	}

	router := http.NewServeMux()
	router.Handle("/", proxy)
	router.Handle("/healthz", healthz(healthy))

	return router
}
