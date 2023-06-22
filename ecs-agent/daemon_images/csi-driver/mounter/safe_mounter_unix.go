//go:build linux || darwin
// +build linux darwin

package mounter

import (
	"k8s.io/mount-utils"
	utilexec "k8s.io/utils/exec"
)

func NewSafeMounter() (*mount.SafeFormatAndMount, error) {
	return &mount.SafeFormatAndMount{
		Interface: mount.New(""),
		Exec:      utilexec.New(),
	}, nil
}
