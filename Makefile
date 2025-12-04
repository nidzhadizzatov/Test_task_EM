# Makefile

# Makefile for Subscription Service

.PHONY: all build run test migrate clean

all: build

build:
	go build -o bin/subscription-service ./cmd/server

run: build
	./bin/subscription-service

test:
	go test ./...

migrate:
	./scripts/migrate.sh

clean:
	go clean
	rm -rf bin/*