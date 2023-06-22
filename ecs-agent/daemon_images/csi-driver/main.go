package main

import (
	"flag"

	"k8s.io/klog/v2"

	"github.com/aws/amazon-ecs-agent/ecs-agent/daemon_images/csi-driver/driver"
)

func main() {
	fs := flag.NewFlagSet("csi-driver", flag.ExitOnError)
	klog.InitFlags(fs)
	srvOptions, err := GetServerOptions(fs)
	if err != nil {
		klog.ErrorS(err, "Failed to get the server options")
		klog.FlushAndExit(klog.ExitFlushTimeout, 1)
	}

	klog.V(4).InfoS("Server Options are provided", "ServerOptions", srvOptions)

	drv, err := driver.NewDriver(
		driver.WithEndpoint(srvOptions.Endpoint),
	)
	if err != nil {
		klog.ErrorS(err, "Failed to create driver")
		klog.FlushAndExit(klog.ExitFlushTimeout, 1)
	}

	if err := drv.Run(); err != nil {
		klog.ErrorS(err, "Failed to run driver")
		klog.FlushAndExit(klog.ExitFlushTimeout, 1)
	}
}
