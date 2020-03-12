.DEFAULT_GOAL := build

install:
	go install -v

build: install
	CGO_ENABLED=0 DEBUG=false GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -a -v -o rattus
	upx -9 --best --ultra-brute --overlay=strip rattus
