//go:build windows
// +build windows

package driver

import (
	"fmt"

	"github.com/aws/amazon-ecs-agent/ecs-agent/daemon_images/csi-driver/mounter"
)

func (m *NodeMounter) PathExists(path string) (bool, error) {
	proxyMounter, ok := m.SafeFormatAndMount.Interface.(*mounter.CSIProxyMounter)
	if !ok {
		return false, fmt.Errorf("failed to cast mounter to csi proxy mounter")
	}
	return proxyMounter.ExistsPath(path)
}
