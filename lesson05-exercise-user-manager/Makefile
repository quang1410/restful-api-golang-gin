APP_NAME := user-manager
CMD_DIR  := ./cmd/api
BIN_DIR  := ./bin
# GOPATH/bin is not always on PATH, so fall back to the absolute path
AIR      := $(shell command -v air 2>/dev/null || echo $(shell go env GOPATH)/bin/air)

.PHONY: run dev install/air build clean tidy test test/run fmt vet

run:
	go run $(CMD_DIR)/main.go

# Hot reload: rebuilds and restarts on every file change
dev:
	@test -x "$(AIR)" || $(MAKE) install/air
	@"$(AIR)"

install/air:
	go install github.com/air-verse/air@latest

build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP_NAME) $(CMD_DIR)/main.go

clean:
	rm -rf $(BIN_DIR)

tidy:
	go mod tidy

fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	go test ./... -v

test/run:
	go test ./... -v -run $(name)
