.PHONY: build build-linux install-linter install-deps run lint test tidy clean

BINARY_NAME = app-listener
OUTPUT_DIR  = build/linux

build: build-linux

build-linux:
	GOOS=linux GOARCH=amd64 go build -o $(OUTPUT_DIR)/$(BINARY_NAME) .
.PHONY: build-linux

install-linter:
	@GOPATH=$$(go env GOPATH); \
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $$GOPATH/bin v2.9.0
.PHONY: install-linter

install-deps:
	go mod download
	go mod verify
.PHONY: install-deps

run:
	go run ./...
.PHONY: run

lint:
	@golangci-lint run
.PHONY: lint

test:
	go test ./... --count=1 -p 1
.PHONY: test

tidy:
	go mod tidy
.PHONY: tidy

clean:
	rm -rf $(OUTPUT_DIR)
.PHONY: clean
