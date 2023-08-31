package csiclient

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/aws/amazon-ecs-agent/ecs-agent/logger"
	"github.com/aws/amazon-ecs-agent/ecs-agent/logger/field"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	PROTOCOL = "unix"
)

type csiClient struct {
	socket     string
	nodeClient csi.NodeClient
}

func NewCSIClient(socketIn string) (csiClient, error) {
	nodeClient, err := newCSINodeClient(socketIn)
	if err != nil {
		logger.Error("Failed to create the CSI node client", logger.Fields{
			field.Error: err,
		})
		return csiClient{}, err
	}

	csiClient := csiClient{
		socket:     socketIn,
		nodeClient: nodeClient,
	}
	return csiClient, nil
}

func newCSINodeClient(socketIn string) (csi.NodeClient, error) {
	dialer := func(addr string, t time.Duration) (net.Conn, error) {
		return net.Dial(PROTOCOL, addr)
	}
	// Set up a connection to the server
	conn, err := grpc.Dial(
		socketIn,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDialer(dialer),
	)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := csi.NewNodeClient(conn)
	return client, nil
}

// Used for testing and integration
// TODO update with stats packing
func (cc *csiClient) GetVolumeMetrics(volumeId string, hostMountPath string) (int64, int64, error) {
	var usedBytes, totalBytes int64

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := cc.nodeClient.NodeGetVolumeStats(ctx, &csi.NodeGetVolumeStatsRequest{
		VolumeId:   volumeId,
		VolumePath: hostMountPath,
	})
	if err != nil {
		logger.Error("Could not get stats", logger.Fields{
			field.Error:  err,
			"VolumeId":   volumeId,
			"VolumePath": hostMountPath,
		})
		return usedBytes, totalBytes, err
	}

	usages := resp.GetUsage()
	// TODO update return type and values to match TCS payload
	if usages == nil {
		return usedBytes, totalBytes, fmt.Errorf("failed to get usage from response: usage is nil")
	}

	for _, usage := range usages {
		unit := usage.GetUnit()
		switch unit {
		case csi.VolumeUsage_BYTES:
			usedBytes = usage.GetUsed()
			totalBytes = usage.GetTotal()
		default:
			logger.Warn("Found missing key in volume usage")
		}
	}
	return usedBytes, totalBytes, nil
}
