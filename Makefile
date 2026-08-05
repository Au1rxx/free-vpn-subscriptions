.PHONY: build test test-migrations fmt lint aggregate

build:
	go build -o bin/fnctl ./cmd/fnctl

test:
	go test ./... -count=1

test-migrations:
	./scripts/test-migrations.sh "$(CONFIG)"

fmt:
	gofmt -w .

lint:
	go vet ./...

aggregate: build
	./bin/fnctl aggregate -c config.yaml
