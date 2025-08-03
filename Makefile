# Check to see if we can use ash, in Alpine images, or default to bash.
SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)

# ==============================================================================
# Building containers

# $(shell git rev-parse --short HEAD)
VERSION := 1.0

all: server

server:
	go build -ldflags "-X main.build=$(VERSION)" -o server ./cmd/server

run: server
	./server

run-help: server
	./server --help

# ==============================================================================
# Building containers

build:
	docker build \
		-f Dockerfile.backend \
		-t launchpad-backend:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

# ==============================================================================
# Running from within docker compose

compose-up:
	docker-compose -f deploy/docker-compose.yml up --detach --remove-orphans

compose-down:
	docker-compose -f deploy/docker-compose.yml down --remove-orphans

compose-logs:
	docker-compose -f deploy/docker-compose.yml logs -f

# ==============================================================================
# Modules support

deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

tidy:
	go mod tidy
	go mod vendor

deps-list:
	go list -m -u -mod=readonly all

deps-upgrade:
	go get -u -v ./...
	go mod tidy
	go mod vendor

deps-cleancache:
	go clean -modcache

# ==============================================================================
# Class Stuff

clean:
	rm -f server
	go clean

.PHONY: all server run run-help build compose-up compose-down compose-logs deps-reset tidy deps-list deps-upgrade deps-cleancache clean