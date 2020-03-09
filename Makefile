install:
	go install -v

build: install
	DEBUG=false GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -a -v -o rattus
	upx -9 --best --ultra-brute rattus

