//go:build linux
// +build linux

package driver

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPathExists(t *testing.T) {
	// set up the full driver and its environment
	dir, err := os.MkdirTemp("", "mount-ebs-csi")
	if err != nil {
		t.Fatalf("error creating directory %v", err)
	}
	defer os.RemoveAll(dir)

	targetPath := filepath.Join(dir, "notafile")

	mountObj, err := newNodeMounter()
	if err != nil {
		t.Fatalf("error creating mounter %v", err)
	}

	exists, err := mountObj.PathExists(targetPath)

	if err != nil {
		t.Fatalf("Expect no error but got: %v", err)
	}

	if exists {
		t.Fatalf("Expected file %s to not exist", targetPath)
	}

}
