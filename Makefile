BINARY := bastion
CMD := ./cmd/bastion
BIN_DIR := ./bin
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -X github.com/taka1156/bastion/internal/app.Version=$(VERSION)

DIST_TARGETS := \
	linux/amd64/tar.gz \
	linux/arm64/tar.gz \
	darwin/amd64/tar.gz \
	darwin/arm64/tar.gz \
	windows/amd64/exe

.PHONY: fmt run build test test-cover clean bin exec dist

fmt:
	go fmt ./...
	golangci-lint run ./...

run:
	go run $(CMD) $(RUN_ARGS)

build:
	mkdir -p $(BIN_DIR)
	rm -r $(BIN_DIR)/${BINARY} || true
	go build -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/$(BINARY) $(CMD)

test:
	go test ./...

test-cover:
	go test -cover ./... -coverprofile=cover.out
	go tool cover -html=cover.out -o cover.html

bin:
	mkdir -p $(BIN_DIR)
	go build -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/$(BINARY) $(CMD)

exec:
	$(BIN_DIR)/bastion -output $(BIN_DIR)/.devcontainer

dist:
	@for target in $(DIST_TARGETS); do \
		GOOS=$$(echo $$target | cut -d/ -f1) \
		GOARCH=$$(echo $$target | cut -d/ -f2) \
		ARCHIVE=$$(echo $$target | cut -d/ -f3) \
		bash scripts/build.sh; \
	done

clean:
	rm -rf $(BIN_DIR) dist tmp
	rm -f cover.out cover.html
