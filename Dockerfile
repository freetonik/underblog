FROM golang:1.13.4-alpine3.10 as builder
ADD . /build
WORKDIR /build
RUN GO111MODULE=on CGO_ENABLED=0 go build -mod=vendor -o underblog ./app/main.go

FROM scratch
COPY --from=builder /build/underblog /app/underblog
ENV PATH "/app:${PATH}"
WORKDIR /blog
CMD ["underblog"]