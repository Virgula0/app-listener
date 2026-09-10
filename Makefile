.PHONY: build build-host build-image build-linux install-linter install-deps generate generate-monitor generate-guard generate-networkmonitor run run-guard run-networkmonitor lint test test-integration check-compatibility deploy deploy-down tidy clean pprof

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
# Go build) inside a Docker container. The host's BTF vmlinux is mounted
# read-only (the CO-RE programs must target the host kernel) and the repo —
# including the output directory — is mounted read-write. Every artifact
# (binary, embeds, vmlinux.h) ends up owned by the invoking host user.
#
# Works with both rootful and rootless Docker: under rootful the container is
# pinned to the caller's uid:gid (--user); under rootless that flag maps to an
# unwritable subordinate uid, so it is dropped and the container's namespaced
# root — which already maps back to the host user — writes the artifacts.
build:
	@if ! command -v docker >/dev/null 2>&1; then \
		echo "ERROR: docker not found — 'make build' needs Docker (rootful or rootless)"; \
		echo "       (or install clang/LLVM, bpftool, Go and GCC and run 'make build-host')"; \
		exit 1; \
	fi
	@docker info >/dev/null 2>&1 || { \
		echo "ERROR: docker daemon not reachable — is it running?"; \
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
	@if docker info -f '{{println .SecurityOptions}}' 2>/dev/null | grep -q rootless; then \
		user_arg=""; \
		echo "make build: rootless Docker detected — running the builder as its namespaced root"; \
		echo "            (artifacts are still written as $$(id -un):$$(id -gn))"; \
	else \
		user_arg="--user $$(id -u):$$(id -g)"; \
	fi; \
	set -x; \
	docker run --rm \
		-v /sys/kernel/btf/vmlinux:/sys/kernel/btf/vmlinux:ro \
		-v "$$PWD:/app/app-listener:rw" \
		$$user_arg \
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

# GUI=1 links the fyne desktop window for `monitor --gui` (~20 MiB larger
# binary, X11/OpenGL deps). Off by default: the daemon and the other
# subcommands share this binary and must not carry a GUI toolkit.
GUI ?=
build-linux:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build $(if $(GUI),-tags gui,) $(if $(VERSION),-ldflags "-X github.com/Virgula0/app-listener/internal/constants.Version=$(VERSION)",) -o $(OUTPUT_DIR)/$(BINARY_NAME) .
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
	go test ./integrationtests/ -v --count=1 -timeout 60m
.PHONY: test-integration

check-compatibility:
	@bash scripts/check-compatibility.sh
.PHONY: check-compatibility

# Collect a full profiling snapshot from a running daemon started with
#   app-listener daemon ... --pprof $(PPROF_ADDR)
# Saves raw *.pb.gz profiles (open later with: go tool pprof -http=: <file>)
# plus ready-to-read -top text reports and the runtime MemStats dump.
#   make pprof                       # 30s CPU sample, default 127.0.0.1:6060
#   make pprof PPROF_SECONDS=60
#   make pprof PPROF_ADDR=127.0.0.1:7070
PPROF_ADDR    ?= 127.0.0.1:6060
PPROF_SECONDS ?= 30
PPROF_URL     := http://$(PPROF_ADDR)/debug/pprof
PPROF_OUT     ?= build/pprof/$(shell date +%Y%m%d-%H%M%S)
pprof:
	@command -v go >/dev/null || { echo "ERROR: go not found"; exit 1; }
	@curl -sf -o /dev/null "$(PPROF_URL)/" 2>/dev/null || { \
		echo "ERROR: no pprof endpoint at $(PPROF_URL)"; \
		echo "  start the daemon with --pprof $(PPROF_ADDR), e.g. via a systemd drop-in:"; \
		echo "    sudo systemctl edit app-listener-daemon"; \
		echo "      [Service]"; \
		echo "      ExecStart="; \
		echo "      ExecStart=/usr/local/sbin/app-listener daemon --headless --blocked-only --pprof $(PPROF_ADDR)"; \
		echo "      Environment=GODEBUG=gctrace=1"; \
		echo "    sudo systemctl restart app-listener-daemon"; \
		exit 1; \
	}
	@mkdir -p "$(PPROF_OUT)"
	@echo ">> profiling $(PPROF_ADDR) -> $(PPROF_OUT)"
	@echo ">> raw profiles (heap, allocs, goroutine, threadcreate)"
	@curl -sf "$(PPROF_URL)/heap"         -o "$(PPROF_OUT)/heap.pb.gz"
	@curl -sf "$(PPROF_URL)/allocs"       -o "$(PPROF_OUT)/allocs.pb.gz"
	@curl -sf "$(PPROF_URL)/goroutine"    -o "$(PPROF_OUT)/goroutine.pb.gz"
	@curl -sf "$(PPROF_URL)/threadcreate" -o "$(PPROF_OUT)/threadcreate.pb.gz" || true
	@echo ">> CPU profile ($(PPROF_SECONDS)s — hold on)"
	@curl -sf "$(PPROF_URL)/profile?seconds=$(PPROF_SECONDS)" -o "$(PPROF_OUT)/cpu.pb.gz"
	@echo ">> text reports"
	@go tool pprof -top -nodecount=40 -inuse_space "$(PPROF_OUT)/heap.pb.gz"   > "$(PPROF_OUT)/heap.inuse_space.txt"   2>/dev/null || true
	@go tool pprof -top -nodecount=40 -inuse_objects "$(PPROF_OUT)/heap.pb.gz" > "$(PPROF_OUT)/heap.inuse_objects.txt" 2>/dev/null || true
	@go tool pprof -top -nodecount=40 -alloc_space "$(PPROF_OUT)/allocs.pb.gz" > "$(PPROF_OUT)/allocs.alloc_space.txt" 2>/dev/null || true
	@go tool pprof -top -nodecount=40 "$(PPROF_OUT)/cpu.pb.gz"                  > "$(PPROF_OUT)/cpu.top.txt"            2>/dev/null || true
	@go tool pprof -tree -nodecount=30 "$(PPROF_OUT)/cpu.pb.gz"                 > "$(PPROF_OUT)/cpu.tree.txt"           2>/dev/null || true
	@echo ">> goroutine summary + runtime MemStats"
	@curl -sf "$(PPROF_URL)/goroutine?debug=1" -o "$(PPROF_OUT)/goroutine.summary.txt" || true
	@curl -sf "$(PPROF_URL)/goroutine?debug=2" -o "$(PPROF_OUT)/goroutine.stacks.txt"  || true
	@curl -sf "$(PPROF_URL)/heap?debug=1"      -o "$(PPROF_OUT)/heap.debug.txt"         || true
	@echo
	@echo "=== CPU top ==="        ; head -25 "$(PPROF_OUT)/cpu.top.txt"          2>/dev/null || true
	@echo "=== heap inuse_space ===" ; head -25 "$(PPROF_OUT)/heap.inuse_space.txt" 2>/dev/null || true
	@echo "=== alloc_space (churn) ===" ; head -25 "$(PPROF_OUT)/allocs.alloc_space.txt" 2>/dev/null || true
	@echo "=== goroutines by count ===" ; grep -E '^goroutine profile|^[0-9]+ @' "$(PPROF_OUT)/goroutine.summary.txt" 2>/dev/null | head -25 || true
	@echo "=== MemStats ===" ; sed -n '/^# runtime.MemStats/,/^# NumGC/p' "$(PPROF_OUT)/heap.debug.txt" 2>/dev/null | head -40 || true
	@echo
	@echo ">> saved to $(PPROF_OUT)/  (interactive: go tool pprof -http=: $(PPROF_OUT)/heap.pb.gz)"

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
	rm -rf $(OUTPUT_DIR) build/test build/pprof \
		internal/bpf/vmlinux.h \
		internal/monitor/monitor_bpf.go internal/monitor/embeds/ \
		internal/guard/guard_bpf.go internal/guard/embeds/ \
		internal/networkmonitor/networkmonitor_bpf.go internal/networkmonitor/embeds/ \
		internal/networkguard/guardnet_bpf.go internal/networkguard/embeds/ \
		internal/monitor/bpf/vmlinux.h internal/guard/bpf/vmlinux.h \
		internal/networkmonitor/bpf/vmlinux.h internal/networkguard/bpf/vmlinux.h
.PHONY: clean
