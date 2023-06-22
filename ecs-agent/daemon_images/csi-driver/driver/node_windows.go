//go:build windows
// +build windows

package driver

import (
	"fmt"

	"github.com/aws/amazon-ecs-agent/ecs-agent/daemon_images/csi-driver/mounter"
)

// IsBlockDevice checks if the given path is a block device
func (d *nodeService) IsBlockDevice(fullPath string) (bool, error) {
	return false, nil
}

// getBlockSizeBytes gets the size of the disk in bytes
func (d *nodeService) getBlockSizeBytes(devicePath string) (int64, error) {
	proxyMounter, ok := (d.mounter.(*NodeMounter)).SafeFormatAndMount.Interface.(*mounter.CSIProxyMounter)
	if !ok {
		return -1, fmt.Errorf("failed to cast mounter to csi proxy mounter")
	}

	sizeInBytes, err := proxyMounter.GetDeviceSize(devicePath)
	if err != nil {
		return -1, err
	}

	return sizeInBytes, nil
}
