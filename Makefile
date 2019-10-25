all: build prepublish move

build:
	GO111MODULE=on CGO_ENABLED=0 go build -mod=vendor -o underblog ./underblog.go

move:
	mv underblog release/ && mv underblog.tar.gz release/

prepublish:
	tar -czf underblog.tar.gz underblog

sha:
	shasum -a 256 underblog.tar.gz
