.DEFAULT_GOAL := build

clean:
	rm -f ./release

install:
	go install -v

build: install
	CGO_ENABLED=0 DEBUG=false GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -a -v -o rattus

release: install
	mkdir -p release
	CGO_ENABLED=0 DEBUG=false GOOS=linux GOARCH=386 go build -ldflags="-s -w" -a -v -o release/rattus-linux-i386
	CGO_ENABLED=0 DEBUG=false GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -a -v -o release/rattus-linux-amd64
	CGO_ENABLED=0 DEBUG=false GOOS=freebsd GOARCH=386 go build -ldflags="-s -w" -a -v -o release/rattus-freebsd-i386
	CGO_ENABLED=0 DEBUG=false GOOS=freebsd GOARCH=amd64 go build -ldflags="-s -w" -a -v -o release/rattus-freebsd-amd64
	CGO_ENABLED=0 DEBUG=false GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -a -v -o release/rattus-darwin-amd64
	CGO_ENABLED=0 DEBUG=false GOOS=windows GOARCH=386 go build -ldflags="-s -w" -a -v -o release/rattus-windows-i386.exe
	CGO_ENABLED=0 DEBUG=false GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -a -v -o release/rattus-windows-amd64.exe
	upx -9 --best --ultra-brute --overlay=strip release/rattus-*
