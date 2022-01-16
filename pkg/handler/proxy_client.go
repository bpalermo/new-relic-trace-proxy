package handler

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
)

const (
	apiKeyHeaderName             = "Api-Key"
	dataFormatHeaderName         = "Data-Format"
	dataFormatHeaderValue        = "zipkin"
	dataFormatVersionHeaderName  = "Data-Format-Version"
	dataFormatVersionHeaderValue = "2"
)

// Hop-by-hop headers. These are removed when sent to the backend.
// http://www.w3.org/Protocols/rfc2616/rfc2616-sec13.html
var hopHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te", // canonicalized version of "TE"
	"Trailers",
	"Transfer-Encoding",
	"Upgrade",
}

// Client proxy server HTTP client
type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

// ProxyClient implements the Client interface
type ProxyClient struct {
	Client       Client
	HostOverride string
	Logger       *logrus.Logger
	APIKey       string
}

func copyHeaderWithoutOverwrite(dst, src http.Header) {
	for k, vv := range src {
		if _, ok := dst[k]; !ok {
			for _, v := range vv {
				dst.Add(k, v)
			}
		}
	}
}

func delHopHeaders(header http.Header) {
	for _, h := range hopHeaders {
		header.Del(h)
	}
}

func (p *ProxyClient) addNewRelicHeaders(header http.Header) {
	header.Add(apiKeyHeaderName, p.APIKey)
	header.Add(dataFormatHeaderName, dataFormatHeaderValue)
	header.Add(dataFormatVersionHeaderName, dataFormatVersionHeaderValue)
}

// Do implement the proxy request handling
func (p *ProxyClient) Do(req *http.Request) (*http.Response, error) {
	proxyURL := *req.URL
	if p.HostOverride != "" {
		proxyURL.Host = p.HostOverride

	} else {
		proxyURL.Host = req.Host
	}
	proxyURL.Scheme = "https"

	if p.Logger.GetLevel() == logrus.DebugLevel {
		initialReqDump, err := httputil.DumpRequest(req, true)
		if err != nil {
			p.Logger.WithError(err).Error("unable to dump request")
		}
		p.Logger.WithField("request", string(initialReqDump)).Debug("Initial request dump:")
	}

	proxyReq, err := http.NewRequest(req.Method, proxyURL.String(), req.Body)
	if err != nil {
		return nil, err
	}

	// Add origin headers after request is signed (no overwrite)
	copyHeaderWithoutOverwrite(proxyReq.Header, req.Header)
	delHopHeaders(req.Header)

	if p.Logger.GetLevel() == logrus.DebugLevel {
		proxyReqDump, err := httputil.DumpRequest(proxyReq, true)
		if err != nil {
			p.Logger.WithError(err).Error("unable to dump request")
		}
		p.Logger.WithField("request", string(proxyReqDump)).Debug("proxying request")
	}

	p.addNewRelicHeaders(proxyReq.Header)

	resp, err := p.Client.Do(proxyReq)
	if err != nil {
		return nil, err
	}
	delHopHeaders(resp.Header)

	if (p.Logger.GetLevel() == logrus.DebugLevel) && resp.StatusCode >= 400 {
		b, _ := ioutil.ReadAll(resp.Body)
		p.Logger.WithField("request", fmt.Sprintf("%s %s", proxyReq.Method, proxyReq.URL)).
			WithField("status_code", resp.StatusCode).
			WithField("message", string(b)).
			Error("error proxying request")

		// Need to "reset" the response body because we consumed the stream above, otherwise caller will
		// get empty body.
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
	}

	return resp, nil
}
