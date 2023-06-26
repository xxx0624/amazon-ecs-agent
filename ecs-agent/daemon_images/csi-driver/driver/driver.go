// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package driver

import (
	"context"
	"net"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc"
	"k8s.io/klog/v2"

	"github.com/aws/amazon-ecs-agent/ecs-agent/daemon_images/csi-driver/util"
	"github.com/aws/amazon-ecs-agent/ecs-agent/daemon_images/csi-driver/version"
)

type Driver struct {
	nodeService

	srv     *grpc.Server
	options *DriverOptions
}

type DriverOptions struct {
	endpoint string
}

func NewDriver(options ...func(*DriverOptions)) (*Driver, error) {
	driverInfo := version.GetVersionInfo()
	klog.InfoS("Driver Information",
		"Driver", "csi-driver",
		"Version", driverInfo.Version,
	)
	klog.V(4).InfoS("Additional driver information",
		"BuildDate", driverInfo.BuildDate,
		"RuntimeGoVersion", driverInfo.GoVersion,
		"RuntimePlatform", driverInfo.Platform,
	)

	driverOptions := DriverOptions{}
	for _, option := range options {
		option(&driverOptions)
	}

	driver := Driver{
		nodeService: newNodeService(),
		options:     &driverOptions,
	}
	return &driver, nil
}

func (d *Driver) Run() error {
	scheme, addr, err := util.ParseEndpoint(d.options.endpoint)
	if err != nil {
		return err
	}

	listener, err := net.Listen(scheme, addr)
	if err != nil {
		return err
	}

	logErr := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			klog.ErrorS(err, "GRPC error")
		}
		return resp, err
	}
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(logErr),
	}
	d.srv = grpc.NewServer(opts...)
	csi.RegisterNodeServer(d.srv, d)

	klog.V(4).InfoS("Listening for connections", "address", listener.Addr())
	return d.srv.Serve(listener)
}

func (d *Driver) Stop() {
	klog.InfoS("Stopping the driver")
	d.srv.Stop()
}

func WithEndpoint(endpoint string) func(*DriverOptions) {
	return func(o *DriverOptions) {
		o.endpoint = endpoint
	}
}
