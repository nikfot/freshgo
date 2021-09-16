# Build time info
PACKAGE = go-versions
VERSION=$(shell git describe --tags)
BUILD=$(shell git rev-parse HEAD)
DATE=$(shell git show -s --format=%ci ${BUILD})
NAME?=go-versions
MAIN_PATH=cmd/go-versions

# Binary output file
BINARY  = go-versions

# Setup ldflags
LDFLAGS = -ldflags "-X 'alethea-elite/pkg/logger.AppName=$(NAME)' -X '$(PACKAGE)/health.Version=$(VERSION)' -X '$(PACKAGE)/health.Build=$(BUILD)' -X '$(PACKAGE)/health.Date=$(DATE)'"

build: clean
	cd ${MAIN_PATH}; \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -v -a -o $(BINARY)

clean:
	rm -f ${MAIN_PATH}/${BINARY}