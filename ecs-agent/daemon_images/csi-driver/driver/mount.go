package driver

import (
	"github.com/aws/amazon-ecs-agent/ecs-agent/daemon_images/csi-driver/mounter"
	mountutils "k8s.io/mount-utils"
	"os"
	"path/filepath"
)

// Mounter is the interface implemented by NodeMounter.
type Mounter interface {
	PathExists(path string) (bool, error)
}

// NodeMounter implements Mounter.
type NodeMounter struct {
	*mountutils.SafeFormatAndMount
}

func newNodeMounter() (Mounter, error) {
	// mounter.NewSafeMounter returns a SafeFormatAndMount
	safeMounter, err := mounter.NewSafeMounter()
	if err != nil {
		return nil, err
	}
	return &NodeMounter{safeMounter}, nil
}

// DeviceIdentifier is for mocking os io functions used for the driver to
// identify an EBS volume's corresponding device (in Linux, the path under
// /dev; in Windows, the volume number) so that it can mount it. For volumes
// already mounted, see GetDeviceNameFromMount.
// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/nvme-ebs-volumes.html#identify-nvme-ebs-device
type DeviceIdentifier interface {
	Lstat(name string) (os.FileInfo, error)
	EvalSymlinks(path string) (string, error)
}

type nodeDeviceIdentifier struct{}

func newNodeDeviceIdentifier() DeviceIdentifier {
	return &nodeDeviceIdentifier{}
}

func (i *nodeDeviceIdentifier) Lstat(name string) (os.FileInfo, error) {
	return os.Lstat(name)
}

func (i *nodeDeviceIdentifier) EvalSymlinks(path string) (string, error) {
	return filepath.EvalSymlinks(path)
}
