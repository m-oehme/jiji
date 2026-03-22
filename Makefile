VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE    ?= $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
LDFLAGS := -X github.com/m-oehme/jiji/internal/version.Version=$(VERSION) \
           -X github.com/m-oehme/jiji/internal/version.Commit=$(COMMIT) \
           -X github.com/m-oehme/jiji/internal/version.Date=$(DATE)

.PHONY: build run lint test clean watch

build:
	go build -ldflags "$(LDFLAGS)" -o bin/jiji ./cmd/jiji

run:
	$(if $(wildcard .env),set -a && . ./.env && set +a &&,) go run -ldflags "$(LDFLAGS)" ./cmd/jiji --debug

lint:
	golangci-lint run

test:
	go test ./...

watch:
	watchexec -e go -r make run

clean:
	rm -rf bin/
