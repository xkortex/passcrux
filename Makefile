VERSION := $(shell git describe --always --dirty --tags)

.PHONY: default get test all

default: get
	go build -i -ldflags="-X 'main.Version=${VERSION}'" -o ${GOPATH}/bin/passcrux

get:
	go get

test: default
	bash ./tests/end2end.sh

fmt:
	go fmt ./...

static: get
	CGO_ENABLED=0 go build -i -ldflags="-X 'main.Version=${VERSION}'" -o ${GOPATH}/bin/passcrux

linux_86:
	GOOS=linux GOARCH=386 go build -i -ldflags="-X 'main.Version=${VERSION}'" -o build/passcrux-linux_386

linux_64:
	GOOS=linux GOARCH=amd64 go build -i -ldflags="-X 'main.Version=${VERSION}'" -o build/passcrux-linux_amd64

linux_arm:
	GOOS=linux GOARCH=arm go build -i -ldflags="-X 'main.Version=${VERSION}'" -o build/passcrux-linux_arm

linux_arm64:
	GOOS=linux GOARCH=arm64 go build -i -ldflags="-X 'main.Version=${VERSION}'" -o build/passcrux-linux_arm64

freebsd_64:
	GOOS=freebsd GOARCH=amd64 go build -i -ldflags="-X 'main.Version=${VERSION}'" -o build/passcrux-freebsd_amd64

darwin_64:
	GOOS=darwin GOARCH=amd64 go build -i -ldflags="-X 'main.Version=${VERSION}'" -o build/passcrux-darwin_amd64

windows:
	GOOS=windows GOARCH=amd64 go build -i -ldflags="-X 'main.Version=${VERSION}'" -o build/passcrux-windows_amd64



