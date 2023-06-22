//go:build linux
// +build linux

package driver

import (
	mountutils "k8s.io/mount-utils"
)

func (m *NodeMounter) PathExists(path string) (bool, error) {
	return mountutils.PathExists(path)
}
