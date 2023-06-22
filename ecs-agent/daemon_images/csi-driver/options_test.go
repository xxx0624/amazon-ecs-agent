package main

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetServerOptions(t *testing.T) {
	testFunc := func(t *testing.T, additionalArgs []string) (*ServerOptions, error) {
		flagSet := flag.NewFlagSet("test-flagset", flag.ContinueOnError)

		args := append([]string{
			"/bin/csi-driver",
		}, additionalArgs...)
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()
		os.Args = args

		options, err := GetServerOptions(flagSet)
		return options, err
	}

	testCases := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{
			name: "Only endpoint is given",
			testFunc: func(t *testing.T) {
				opts, err := testFunc(t, []string{"--endpoint=foo"})
				assert.NoError(t, err)
				assert.Equal(t, "foo", opts.Endpoint)
			},
		},
		{
			name: "No argument is given",
			testFunc: func(t *testing.T) {
				_, err := testFunc(t, nil)
				assert.NotNil(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, tc.testFunc)
	}
}
