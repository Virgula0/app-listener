.PHONY: build build-host build-image build-linux install-linter install-deps generate generate-monitor generate-guard generate-networkmonitor run run-guard run-networkmonitor lint test test-integration check-compatibility deploy deploy-down tidy clean

BINARY_NAME = app-listener
OUTPUT_DIR  = build/linux
GEN_DIR     = build/generated
BUILD_IMAGE = app-listener-builder:local

# VERSION is injected into the binary by the release workflow
# (pre-<date>-<sha>); when empty the embedded constants.Version default
# (dev marker) is kept.
VERSION ?=

.PHONY: require-kernel51
require-kernel51:
	@kernel_ver=$$(uname -r | cut -d. -f1); \
	if [ -z "$$kernel_ver" ] || [ "$$kernel_ver" -lt 5 ] 2>/dev/null; then \
		echo "ERROR: kernel $$(uname -r) is too old: kernel 5.x or newer is required (build host detected $${kernel_ver}.x)"; \
		exit 1; \
	fi

# Isolated build: runs the whole pipeline (vmlinux.h dump + BPF bindings +
# Go build) inside a rootful Docker container. The host's BTF vmlinux is
# mounted read-only (the CO-RE programs must target the host kernel) and the
# repo — including the output directory — is mounted read-write. The
# container runs as the calling host user, so every artifact (binary, embeds,
# vmlinux.h) ends up owned by the host user. Requires a rootful Docker daemon.
build:
	@if ! command -v docker >/dev/null 2>&1; then \
		echo "ERROR: docker not found — 'make build' needs a rootful Docker daemon"; \
		echo "       (install Docker, or install clang/LLVM, bpftool, Go and GCC and run 'make build-host')"; \
		exit 1; \
	fi
	@docker info >/dev/null 2>&1 || { \
		echo "ERROR: docker daemon not reachable — is it running (rootful)?"; \
		exit 1; \
	}
	@if [ ! -r /sys/kernel/btf/vmlinux ]; then \
		echo "ERROR: /sys/kernel/btf/vmlinux not readable — run 'make check-compatibility' (kernel needs CONFIG_DEBUG_INFO_BTF)"; \
		exit 1; \
	fi
	@if [ -e $(OUTPUT_DIR) ] && [ ! -w $(OUTPUT_DIR) ]; then \
		echo "ERROR: $(OUTPUT_DIR) is not writable by user $$(id -u) — fix ownership with:"; \
		echo "       sudo chown -R $$(id -u):$$(id -g) $(OUTPUT_DIR)"; \
		exit 1; \
	fi
	@mkdir -p $(OUTPUT_DIR)
	$(MAKE) build-image
	docker run --rm \
		-v /sys/kernel/btf/vmlinux:/sys/kernel/btf/vmlinux:ro \
		-v "$$PWD:/app/app-listener:rw" \
		--user "$$(id -u):$$(id -g)" \
		-e HOME=/tmp \
		-e GOCACHE=/tmp/.gocache \
		-e GOPATH=/tmp/gopath \
		-e GOMODCACHE=/tmp/gopath/pkg/mod \
		-e VERSION="$(VERSION)" \
		-w /app/app-listener \
		$(BUILD_IMAGE) \
		make bpftool-headers generate build-linux
.PHONY: build

# On-host build: requires clang/LLVM, bpftool, Go and GCC installed locally.
build-host: require-kernel51 bpftool-headers generate build-linux

build-image:
	docker build -t $(BUILD_IMAGE) -f docker/builder.Dockerfile .
.PHONY: build-image

build-linux:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build $(if $(VERSION),-ldflags "-X github.com/Virgula0/app-listener/internal/constants.Version=$(VERSION)",) -o $(OUTPUT_DIR)/$(BINARY_NAME) .
.PHONY: build-linux

# Ensure the shared vmlinux.h is dumped before any BPF module is compiled
generate: bpftool-headers generate-monitor generate-guard generate-networkmonitor generate-networkguard

generate-monitor: bpftool-headers
	@mkdir -p $(GEN_DIR) internal/monitor/embeds
	GOPACKAGE=monitor GOOS=linux GOARCH=amd64 go run github.com/cilium/ebpf/cmd/bpf2go \
		-cc clang \
		-cflags "-O2 -g -Wall -Wno-visibility -Wno-attributes -D__TARGET_ARCH_x86 -I internal/bpf -I/usr/include/x86_64-linux-gnu" \
		-target bpf \
		-output-dir $(GEN_DIR) \
		Monitor ./internal/monitor/bpf/monitor.bpf.c
	@mv $(GEN_DIR)/monitor_bpf.go internal/monitor/monitor_bpf.go
	@mv $(GEN_DIR)/monitor_bpf.o internal/monitor/embeds/monitor_bpf.o
	@sed -i 's|monitor_bpf\.o|embeds/monitor_bpf.o|' internal/monitor/monitor_bpf.go
	@rm -rf $(GEN_DIR)
	@echo "Monitor BPF generation complete"
.PHONY: generate-monitor

generate-guard: bpftool-headers
	@mkdir -p $(GEN_DIR) internal/guard/embeds
	GOPACKAGE=guard GOOS=linux GOARCH=amd64 go run github.com/cilium/ebpf/cmd/bpf2go \
		-cc clang \
		-cflags "-O2 -g -Wall -Wno-visibility -Wno-attributes -D__TARGET_ARCH_x86 -I internal/bpf -I/usr/include/x86_64-linux-gnu" \
		-target bpf \
		-output-dir $(GEN_DIR) \
		Guard ./internal/guard/bpf/guard.bpf.c
	@mv $(GEN_DIR)/guard_bpf.go internal/guard/guard_bpf.go
	@mv $(GEN_DIR)/guard_bpf.o internal/guard/embeds/guard_bpf.o
	@sed -i 's|guard_bpf\.o|embeds/guard_bpf.o|' internal/guard/guard_bpf.go
	@rm -rf $(GEN_DIR)
	@echo "Guard BPF generation complete"
.PHONY: generate-guard

generate-networkmonitor: bpftool-headers
	@mkdir -p $(GEN_DIR) internal/networkmonitor/embeds
	GOPACKAGE=networkmonitor GOOS=linux GOARCH=amd64 go run github.com/cilium/ebpf/cmd/bpf2go \
		-cc clang \
		-cflags "-O2 -g -Wall -Wno-visibility -Wno-attributes -D__TARGET_ARCH_x86 -I internal/bpf -I/usr/include/x86_64-linux-gnu" \
		-target bpf \
		-output-dir $(GEN_DIR) \
		NetMon ./internal/networkmonitor/bpf/networkmonitor.bpf.c
	@mv $(GEN_DIR)/netmon_bpf.go internal/networkmonitor/networkmonitor_bpf.go
	@mv $(GEN_DIR)/netmon_bpf.o internal/networkmonitor/embeds/networkmonitor_bpf.o
	@sed -i 's|netmon_bpf\.o|embeds/networkmonitor_bpf.o|' internal/networkmonitor/networkmonitor_bpf.go
	@rm -rf $(GEN_DIR)
	@echo "Network monitor BPF generation complete"
.PHONY: generate-networkmonitor

generate-networkguard: bpftool-headers
	@mkdir -p $(GEN_DIR) internal/networkguard/embeds
	GOPACKAGE=networkguard GOOS=linux GOARCH=amd64 go run github.com/cilium/ebpf/cmd/bpf2go \
		-cc clang \
		-cflags "-O2 -g -Wall -Wno-visibility -Wno-attributes -D__TARGET_ARCH_x86 -I internal/bpf -I/usr/include/x86_64-linux-gnu" \
		-target bpf \
		-output-dir $(GEN_DIR) \
		GuardNet ./internal/networkguard/bpf/networkguard.bpf.c
	@mv $(GEN_DIR)/guardnet_bpf.go internal/networkguard/guardnet_bpf.go
	@mv $(GEN_DIR)/guardnet_bpf.o internal/networkguard/embeds/guardnet_bpf.o
	@sed -i 's|guardnet_bpf\.o|embeds/guardnet_bpf.o|' internal/networkguard/guardnet_bpf.go
	@rm -rf $(GEN_DIR)
	@echo "Network guard BPF generation complete"
.PHONY: generate-networkguard

bpftool-headers:
	@if ! command -v bpftool >/dev/null 2>&1; then \
		echo "ERROR: bpftool not found — install it (Ubuntu: apt install linux-tools-generic; Arch: pacman -S bpftool)"; \
		exit 1; \
	fi
	@if [ ! -r /sys/kernel/btf/vmlinux ]; then \
		echo "ERROR: /sys/kernel/btf/vmlinux not readable — run 'make check-compatibility' (kernel needs CONFIG_DEBUG_INFO_BTF)"; \
		exit 1; \
	fi
	@mkdir -p internal/bpf
	bpftool btf dump file /sys/kernel/btf/vmlinux format c \
		> internal/bpf/vmlinux.h
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

run-guard:
	CGO_ENABLED=1 go run ./... guard $(ARGS)
.PHONY: run-guard

run-networkmonitor:
	CGO_ENABLED=1 go run ./... network-monitor $(ARGS)
.PHONY: run-networkmonitor

run-networkguard:
	CGO_ENABLED=1 go run ./... network-guard $(ARGS)
.PHONY: run-networkguard

lint:
	@golangci-lint run
.PHONY: lint

test:
	CGO_ENABLED=1 go test $$(go list ./... | grep -v /integrationtests) --count=1 -p 1
.PHONY: test

test-integration:
	$(MAKE) -C integrationtests/exploits
	go test ./integrationtests/ -v --count=1 -timeout 30m
.PHONY: test-integration

check-compatibility:
	@bash scripts/check-compatibility.sh
.PHONY: check-compatibility

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
		internal/bpf/vmlinux.h \
		internal/monitor/monitor_bpf.go internal/monitor/embeds/ \
		internal/guard/guard_bpf.go internal/guard/embeds/ \
		internal/networkmonitor/networkmonitor_bpf.go internal/networkmonitor/embeds/ \
		internal/networkguard/guardnet_bpf.go internal/networkguard/embeds/ \
		internal/monitor/bpf/vmlinux.h internal/guard/bpf/vmlinux.h \
		internal/networkmonitor/bpf/vmlinux.h internal/networkguard/bpf/vmlinux.h
.PHONY: clean
