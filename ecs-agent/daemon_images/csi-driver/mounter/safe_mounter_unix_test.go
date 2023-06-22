//go:build linux
// +build linux

package mounter

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewSafeMounter(t *testing.T) {
	resp, err := NewSafeMounter()
	assert.NotNil(t, resp)
	assert.Nil(t, err)
}
