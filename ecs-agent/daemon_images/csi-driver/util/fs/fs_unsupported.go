//go:build !linux && !darwin && !windows
// +build !linux,!darwin,!windows

package fs

import (
	"fmt"
)

type UsageInfo struct {
	Bytes  int64
	Inodes int64
}

// Info unsupported returns 0 values for available and capacity and an error.
func Info(path string) (int64, int64, int64, int64, int64, int64, error) {
	return 0, 0, 0, 0, 0, 0, fmt.Errorf("fsinfo not supported for this build")
}
