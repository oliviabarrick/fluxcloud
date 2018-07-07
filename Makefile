.PHONY: test build

test:
	test -z $(shell gofmt -l ./cmd/ ./pkg/)
	go test ./...

build:
	gofmt -w ./cmd ./pkg
	CGO_ENABLED=0 go build -ldflags '-w -s' -installsuffix cgo -o fluxcloud ./cmd/
