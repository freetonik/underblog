all: build prepublish move

build:
	GO111MODULE=on CGO_ENABLED=0 go build -mod=vendor -o underblog ./app/main.go

move:
	mv underblog release/ && mv underblog.tar.gz release/

prepublish:
	tar -czf underblog.tar.gz underblog

sha:
	shasum -a 256 underblog.tar.gz

test:
	go test -v -count 1 -race -cover ./...

bench:
	go test -v -run Bench -bench=. ./...