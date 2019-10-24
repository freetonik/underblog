all: build prepublish move

build:
	go build

move:
	mv underblog release/ && mv underblog.tar.gz release/

prepublish:
	tar -czf underblog.tar.gz underblog

sha:
	shasum -a 256 underblog.tar.gz
