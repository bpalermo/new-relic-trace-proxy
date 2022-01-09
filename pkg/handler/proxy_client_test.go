package handler

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func getHopHeaders() http.Header {
	headers := http.Header{}
	for _, h := range hopHeaders {
		headers.Add(h, "")
	}
	return headers
}

func TestDelHopHeaders(t *testing.T) {
	const headerValue = "test"

	expected := http.Header{}
	expected.Add(apiKeyHeaderName, headerValue)

	actual := getHopHeaders()
	actual.Add(apiKeyHeaderName, headerValue)

	delHopHeaders(actual)
	assert.Equal(t, expected, actual)
}

func TestCopyHeaderWithoutOverwrite(t *testing.T) {
	const headerValue = "test"

	expected := http.Header{}
	expected.Add(apiKeyHeaderName, headerValue)

	actual := http.Header{}
	copyHeaderWithoutOverwrite(actual, expected)

	assert.Equal(t, expected, actual)
}

func TestAddNewRelicHeaders(t *testing.T) {
	const apiKey = "test"

	client := ProxyClient{
		ApiKey: apiKey,
	}

	actual := http.Header{}
	client.addNewRelicHeaders(actual)

	assert.Equal(t, apiKey, actual.Get(apiKeyHeaderName))
	assert.Equal(t, dataFormatHeaderValue, actual.Get(dataFormatHeaderName))
	assert.Equal(t, dataFormatVersionHeaderValue, actual.Get(dataFormatVersionHeaderName))
}
