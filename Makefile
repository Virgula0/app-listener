.PHONY: build build-linux install-linter install-deps generate run lint test tidy clean

BINARY_NAME = app-listener
OUTPUT_DIR  = build/linux

build: generate build-linux

build-linux:
	GOOS=linux GOARCH=amd64 go build -o $(OUTPUT_DIR)/$(BINARY_NAME) .
.PHONY: build-linux

GEN_DIR = build/generated

generate:
	@mkdir -p $(GEN_DIR) internal/infrastructure/ebpf/embeds
	GOPACKAGE=ebpf GOOS=linux GOARCH=amd64 go run github.com/cilium/ebpf/cmd/bpf2go \
		-cc clang \
		-cflags "-O2 -g -Wall -Wno-visibility -Wno-attributes -D__TARGET_ARCH_x86" \
		-target bpf \
		-output-dir $(GEN_DIR) \
		Monitor ./internal/infrastructure/ebpf/bpf/monitor.bpf.c
	@mv $(GEN_DIR)/monitor_bpf.go internal/infrastructure/ebpf/monitor_bpf.go
	@mv $(GEN_DIR)/monitor_bpf.o internal/infrastructure/ebpf/embeds/monitor_bpf.o
	@sed -i 's|monitor_bpf\.o|embeds/monitor_bpf.o|' internal/infrastructure/ebpf/monitor_bpf.go
	@rm -rf $(GEN_DIR)
	@echo "BPF generation complete"
.PHONY: generate

bpftool-headers:
	bpftool btf dump file /sys/kernel/btf/vmlinux format c 2>/dev/null > internal/infrastructure/ebpf/bpf/vmlinux.h
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
	go run ./... monitor $(ARGS)
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
	rm -rf $(OUTPUT_DIR) \
		internal/infrastructure/ebpf/monitor_bpf.go \
		internal/infrastructure/ebpf/embeds/ \
		internal/infrastructure/ebpf/bpf/vmlinux.h
.PHONY: clean
