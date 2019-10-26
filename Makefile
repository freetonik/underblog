include version.mk

all: build test targz sha move

build:
	GO111MODULE=on CGO_ENABLED=0 go build -mod=vendor -o underblog ./app/main.go

test:
	go test -v -count 1 -race -cover ./...

bench:
	go test -v -run Bench -bench=. ./...

targz:
	tar -czf underblog_$(version).tar.gz underblog

sha:
	shasum -a 256 underblog_$(version).tar.gz

move:
	mv underblog release/ && mv underblog_$(version).tar.gz release/
