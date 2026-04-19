.PHONY: build test fmt lint aggregate

build:
	go build -o bin/fnctl ./cmd/fnctl

test:
	go test ./... -count=1

fmt:
	gofmt -w .

lint:
	go vet ./...

aggregate: build
	./bin/fnctl aggregate -c config.yaml
