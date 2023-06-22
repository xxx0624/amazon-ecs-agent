package main

import (
	"errors"
	"flag"
	"os"
)

const EmptyCSIEndpoint = ""

type ServerOptions struct {
	// Endpoint is the endpoint that the driver server should listen on.
	Endpoint string
}

func GetServerOptions(fs *flag.FlagSet) (*ServerOptions, error) {
	serverOptions := &ServerOptions{}
	fs.StringVar(&serverOptions.Endpoint, "endpoint", EmptyCSIEndpoint, "Endpoint for the CSI driver server")

	args := os.Args[1:]
	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if serverOptions.Endpoint == EmptyCSIEndpoint {
		return nil, errors.New("no endpoint is provided")
	}

	return serverOptions, nil
}
