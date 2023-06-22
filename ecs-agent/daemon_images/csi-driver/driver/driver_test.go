package driver

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWithEndpoint(t *testing.T) {
	value := "endpoint"
	options := &DriverOptions{}

	WithEndpoint(value)(options)

	assert.Equal(t, value, options.endpoint)
}
