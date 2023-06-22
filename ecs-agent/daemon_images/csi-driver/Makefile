ARCH?=amd64
OS?=linux
GO111MODULE=on
GOPATH=$(shell go env GOPATH)
GOOS=$(shell go env GOOS)

ROOT := $(shell pwd)
SOURCEDIR=./
CSI_DRIVER_SOURCES=$(shell find $(SOURCEDIR) -name '*.go')
DESTDIR=./

CSI_DRIVER_BINARY=$(ROOT)/bin/csi-driver
$(CSI_DRIVER_BINARY): $(CSI_DRIVER_SOURCES)
	CGO_ENABLED=0 GOOS=linux GOARCH=$(ARCH) go build -o $(CSI_DRIVER_BINARY) $(DESTDIR)

.PHONY: build-csi-driver-image
build-csi-driver-image: $(CSI_DRIVER_BINARY)
	ARCH=${ARCH} ./csi-driver-image

.PHONY: test
test:
	go test -v -race ./...

# please include this `mockgen` into PATH
bin/mockgen: | bin
	go install github.com/golang/mock/mockgen@v1.6.0

.PHONY: mockgen
mockgen: bin/mockgen
	./update-gomock

.PHONY: clean
clean:
	rm -rf bin/