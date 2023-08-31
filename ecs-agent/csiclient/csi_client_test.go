//go:build unit
// +build unit

package csiclient

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

type mockNodeClient struct {
}

func (c mockNodeClient) NodeStageVolume(ctx context.Context,
	in *csi.NodeStageVolumeRequest,
	opts ...grpc.CallOption) (*csi.NodeStageVolumeResponse, error) {
	// TODO - mock
	return nil, nil
}
func (c mockNodeClient) NodeUnstageVolume(ctx context.Context,
	in *csi.NodeUnstageVolumeRequest,
	opts ...grpc.CallOption) (*csi.NodeUnstageVolumeResponse, error) {
	// TODO - mock
	return nil, nil
}
func (c mockNodeClient) NodePublishVolume(ctx context.Context,
	in *csi.NodePublishVolumeRequest,
	opts ...grpc.CallOption) (*csi.NodePublishVolumeResponse, error) {
	// TODO - mock
	return nil, nil
}
func (c mockNodeClient) NodeUnpublishVolume(ctx context.Context,
	in *csi.NodeUnpublishVolumeRequest,
	opts ...grpc.CallOption) (*csi.NodeUnpublishVolumeResponse, error) {
	// TODO - mock
	return nil, nil
}
func (c mockNodeClient) NodeGetVolumeStats(ctx context.Context,
	in *csi.NodeGetVolumeStatsRequest,
	opts ...grpc.CallOption) (*csi.NodeGetVolumeStatsResponse, error) {
	return &csi.NodeGetVolumeStatsResponse{
		Usage: []*csi.VolumeUsage{
			&csi.VolumeUsage{
				Available: 100,
				Total:     250,
				Used:      150,
				Unit:      csi.VolumeUsage_BYTES,
			},
		},
	}, nil
}
func (c mockNodeClient) NodeExpandVolume(ctx context.Context,
	in *csi.NodeExpandVolumeRequest,
	opts ...grpc.CallOption) (*csi.NodeExpandVolumeResponse, error) {
	// TODO - mock
	return nil, nil
}
func (c mockNodeClient) NodeGetCapabilities(ctx context.Context,
	in *csi.NodeGetCapabilitiesRequest,
	opts ...grpc.CallOption) (*csi.NodeGetCapabilitiesResponse, error) {
	// TODO - mock
	return nil, nil
}
func (c mockNodeClient) NodeGetInfo(ctx context.Context,
	in *csi.NodeGetInfoRequest,
	opts ...grpc.CallOption) (*csi.NodeGetInfoResponse, error) {
	// TODO - mock
	return nil, nil
}

func TestGetVolumeMetrics(t *testing.T) {
	csiClient := csiClient{
		socket:     "",
		nodeClient: mockNodeClient{},
	}

	used, total, err := csiClient.GetVolumeMetrics("vol-1", "/")
	require.NoError(t, err)

	assert.Equal(t, int64(150), used)
	assert.Equal(t, int64(250), total)
}
