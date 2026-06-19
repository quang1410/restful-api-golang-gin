APP_NAME := project-shopping-cart
CMD_DIR  := ./cmd/api
BIN_DIR  := ./bin

.PHONY: run build clean tidy test

run:
	go run $(CMD_DIR)/main.go

build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP_NAME) $(CMD_DIR)/main.go

clean:
	rm -rf $(BIN_DIR)

tidy:
	go mod tidy

test:
	go test ./... -v

test/run:
	go test ./... -v -run $(name)
