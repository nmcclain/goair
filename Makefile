#!/usr/bin/make -f

# from https://github.com/davecheney/golang-crosscompile

SHELL=/bin/bash

build: deps
	go build

release: deps golang-crosscompile
	source golang-crosscompile/crosscompile.bash; \
	go-darwin-386 build -o release/goair-Darwin-i386; \
	go-darwin-amd64 build -o release/goair-Darwin-x86_64; \
	go-linux-386 build -o release/goair-Linux-i386; \
	go-linux-386 build -o release/goair-Linux-i686; \
	go-linux-amd64 build -o release/goair-Linux-x86_64; \
	go-linux-arm build -o release/goair-Linux-armv6l; \
	go-linux-arm build -o release/goair-Linux-armv7l; \
	go-freebsd-386 build -o release/goair-FreeBSD-i386; \
	go-freebsd-amd64 build -o release/goair-FreeBSD-amd64; \
	go-windows-386 build -o release/goair.exe; \
	CGO_ENABLED=0 go build -a -ldflags '-s' -o release/goair-Linux-static

golang-crosscompile:
	git clone https://github.com/davecheney/golang-crosscompile.git

deps:
	go clean -i net && go install -tags netgo std
