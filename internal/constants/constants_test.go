package constants

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConstants(t *testing.T) {
	var port uint = 9001
	assert.Equal(t, port, DefaultPort)
	assert.Equal(t, "/healthz", HealthCheckPath)
}
