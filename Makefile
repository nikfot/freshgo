# Build time info
PACKAGE = freshgo
VERSION=$(shell git describe --tags)
BUILD=$(shell git rev-parse HEAD)
DATE=$(shell git show -s --format=%ci ${BUILD})
NAME?=freshgo
MAIN_PATH=cmd/cli

# Binary output file
BINARY  = freshgo

build: clean
	cd ${MAIN_PATH}; \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -v -a -o $(BINARY)

clean:
	rm -f ${MAIN_PATH}/${BINARY}

run: build
	cd cmd/cli && ./freshgo latest
