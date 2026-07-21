.PHONY: build build-linux install-linter install-deps generate run lint test test-integration deploy deploy-down tidy clean

BINARY_NAME = app-listener
OUTPUT_DIR  = build/linux

build: generate build-linux

build-linux:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o $(OUTPUT_DIR)/$(BINARY_NAME) .
.PHONY: build-linux

GEN_DIR = build/generated

generate:
	@mkdir -p $(GEN_DIR) internal/infrastructure/embeds
	GOPACKAGE=ebpf GOOS=linux GOARCH=amd64 go run github.com/cilium/ebpf/cmd/bpf2go \
		-cc clang \
		-cflags "-O2 -g -Wall -Wno-visibility -Wno-attributes -D__TARGET_ARCH_x86" \
		-target bpf \
		-output-dir $(GEN_DIR) \
		Monitor ./internal/infrastructure/bpf/monitor.bpf.c
	@mv $(GEN_DIR)/monitor_bpf.go internal/infrastructure/monitor_bpf.go
	@mv $(GEN_DIR)/monitor_bpf.o internal/infrastructure/embeds/monitor_bpf.o
	@sed -i 's|monitor_bpf\.o|embeds/monitor_bpf.o|' internal/infrastructure/monitor_bpf.go
	@rm -rf $(GEN_DIR)
	@echo "BPF generation complete"
.PHONY: generate

bpftool-headers:
	bpftool btf dump file /sys/kernel/btf/vmlinux format c 2>/dev/null > internal/infrastructure/bpf/vmlinux.h
.PHONY: bpftool-headers

install-linter:
	@GOPATH=$$(go env GOPATH); \
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $$GOPATH/bin v2.9.0
.PHONY: install-linter

install-deps:
	go mod download
	go mod verify
.PHONY: install-deps

run:
	CGO_ENABLED=1 go run ./... monitor $(ARGS)
.PHONY: run

lint:
	@golangci-lint run
.PHONY: lint

test:
	CGO_ENABLED=1 go test $$(go list ./... | grep -v /tests) --count=1 -p 1
.PHONY: test

test-integration:
	go test ./tests/ -v --count=1 -timeout 15m
.PHONY: test-integration

deploy:
	docker compose up --build -d
.PHONY: deploy

deploy-down:
	docker compose down
.PHONY: deploy-down

tidy:
	go mod tidy
.PHONY: tidy

clean:
	rm -rf $(OUTPUT_DIR) build/test \
		internal/infrastructure/monitor_bpf.go \
		internal/infrastructure/embeds/ \
		internal/infrastructure/bpf/vmlinux.h
.PHONY: clean
