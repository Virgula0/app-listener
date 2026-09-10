# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

A single Go binary (`app-listener`) that monitors or **denies** filesystem and network
operations with eBPF, to protect credential directories (SSH keys, GPG, browser
profiles, AI-agent tokens) from info-stealers and supply-chain attacks. Runs as root
on Linux only; `linux/amd64` is the released target.

`README.md` is the user-facing reference for every subcommand and flag — consult it
before changing CLI surface. This is a "vibe-coding experiment, not for production".

## Security is the product — do not regress it

This is a **security enforcement module**. `guard` and `daemon` mode are the parts users
trust to *deny* access to secrets; a bug there is a vulnerability, not a defect. Treat
every change to enforcement paths as security-sensitive:

- **Fail closed.** On any error — a hook that won't attach, an unreadable inode map, a
  malformed config, a race during unlock — the correct outcome is "deny / refuse to
  start", never "continue unprotected". Preserve existing fail-closed checks
  (`CheckBPFLSM`, the attach-before-unlock ordering, deferred-binary resolution,
  `/sys/fs/bpf` link pinning, `internal/protected` daemon-alive gates).
- **No TOCTOU windows.** The daemon's invariant is that a resource is never readable
  without a live guard attached. Keep the `attach → unlock → populate → resolve →
  re-sync` order; never move an unlock earlier or a detach earlier.
- **Identity stays inode-based.** Any new policy check must key on `dev:ino` of the
  executable, never on path, name, comm, argv, or a caller-supplied string — those are
  all attacker-controlled.
- **New syscall / kernel path = new bypass surface.** When you add or change an LSM hook
  or a guarded operation, add a matching `integrationtests/exploits/*.c` case that proves
  the operation is denied, and run `make test-integration` locally.
- **Watch the well-known bypass classes** already defended against (see `memory/*` and the
  exploit corpus): io_uring, `open_by_handle_at`, `copy_file_range`/`splice`/`sendfile`,
  `process_vm_readv`, raw block-device reads, bind-mount / rename-over-watchroot,
  `SCM_RIGHTS` fd passing, `ptrace`, stat/metadata leaks. Don't reintroduce a gap a past
  commit closed.
- **The live-edit grant is a real escalation path — keep it narrow.** `edit-protected`
  live mode (`cmd/functions/editprotected` + `cmd/functions/daemon/control.go`)
  authenticates over a `0600` root-only unix socket that *also* checks `SO_PEERCRED`
  uid 0 **and** that the peer's exe inode is the daemon's own binary. The protocol is
  two-phase — `AUTH` (password) must succeed **before** the daemon discloses any
  watch path, then `SELECT <resource>` activates the grant — so an unauthenticated
  caller learns nothing about the protected directories. The grant transiently widens
  **only** that one resource's `GUARD_ALLOW_ROOT` self event-mask (never the
  whitelist), one session at a time, 30-min hard cap, revoked on disconnect / SIGHUP.
  The password hash (`/etc/app-listener/edit-auth.hash`, PBKDF2, `0600`) is self-guarded
  like `fscrypt.key`. Don't loosen any of: the peer-exe check, AUTH-before-disclosure,
  the single-session lock, the auth lockout, the `origin=install` guard on
  `--set-password`, or the mask being restored in `RevokeSelfEditAccess`.
- **The daemon self-protects.** It guards its own `/etc/app-listener` config + fscrypt
  key + edit-auth hash (`selfguards.go`) and drops to `guard.ModeReadOnly` where appropriate — keep those
  guarantees intact.
- Run `/security-review` on diffs that touch `internal/guard`, `internal/networkguard`,
  `internal/usecase/daemon.go`, `internal/fscrypt`, `internal/protected`, or any
  `*.bpf.c`.

## Build

`CGO_ENABLED=1` is mandatory everywhere — the daemon links `google/fscrypt`, which
needs cgo (`mlock`). A `CGO_ENABLED=0` build no longer compiles.

| Command | Use |
|---|---|
| `make build` | Full pipeline in a **rootful Docker** container: dumps `vmlinux.h` from the host BTF, regenerates BPF bindings, builds to `build/linux/app-listener`. Nothing but Docker needed on the host. |
| `make build-host` | Same pipeline on the host — needs a **recent** `clang`/LLVM (`docker/builder.Dockerfile` pins clang 22; clang 14 emits a `guard_path_rename` that overflows the BPF verifier's 1M-insn budget on modern kernels — issue #45), `bpftool`, Go 1.26+, GCC. |
| `make build-linux` | Go compile only (assumes generated files exist). `GUI=1` adds `-tags gui` (fyne desktop window for `monitor --gui`, ~20 MiB larger, X11/OpenGL deps) — off by default because the daemon shares this binary. |
| `make generate` | Regenerate BPF bindings (`bpf2go`). Needs `bpftool` + readable `/sys/kernel/btf/vmlinux`. |
| `make check-compatibility` | Static host check (kernel version, `.config`, BTF, BPF-LSM activation, fscrypt prereqs). Run before building/running on a new host. `--binary <path>` (as root) also loads that binary's guard eBPF into the running verifier (`daemon --check`), catching a prebuilt/kernel mismatch; `scripts/install.sh` runs this against the downloaded release before installing anything. |

BPF C sources live in `internal/<mode>/bpf/*.bpf.c`. `make generate` compiles each with
clang and, via `bpf2go`, emits `internal/<mode>/<name>_bpf.go` plus
`internal/<mode>/embeds/<name>_bpf.o`, which is `//go:embed`ed into the binary. These
generated files and `internal/bpf/vmlinux.h` are gitignored — regenerate, don't edit.

`internal/constants.Version` is overwritten at release time via `-ldflags -X`; leave the
`dev` default in the source.

## Test

- `make test` → `CGO_ENABLED=1 go test $(go list ./... | grep -v /integrationtests) --count=1 -p 1`.
  Package-parallelism is forced to `-p 1` (BPF/global state); keep new tests serial-safe.
- Single unit test: `CGO_ENABLED=1 go test ./internal/guard/ -run TestName -count=1`.
- `make test-integration` → builds the C exploit corpus (`make -C integrationtests/exploits`),
  then `go test ./integrationtests/ -v --count=1 -timeout 60m`. Needs **rootful Docker**;
  spins privileged testcontainers (`CAP_BPF`, `SYS_ADMIN`, host `/sys/kernel/btf` bind-mounted,
  `APPLISTENER_ASSUME_BPF_LSM=1` set inside because containers have no securityfs). Auto-skips
  when `docker` is not on PATH.
- Single integration test: `go test ./integrationtests/ -v -run 'TestIntegrationSuite/TestGuardWhitelist' -count=1`.
- The integration suite compiles the binary under test statically with `-tags ci,osusergo,netgo`
  (`main_test.go`), and runs it against `ubuntu:latest` containers whose glibc is older than a
  typical host. Containers are pooled per test file (`pool_test.go`); add heavy new tests to an
  existing file's pool rather than starting fresh containers.
- `integrationtests/exploits/*.c` are deliberate bypass attempts (io_uring, `open_by_handle_at`,
  `copy_file_range`, `process_vm_readv`, raw block device, mount, `SCM_RIGHTS` fd-passing, …).
  Adding a new BPF hook usually means adding a matching exploit here to prove it's covered.

## Lint

`make lint` → golangci-lint v2.9.0 (`make install-linter`), config in `.golangci.yaml`.
`run.tests: false` — test files are not linted. Notable limits: `funlen` 100 lines,
`gocyclo` 15, `gocognit` 20, `gosec` on, `dupl` on. golangci-lint typechecks the whole
module including the fyne GUI, so CI installs `libgl1-mesa-dev xorg-dev libwayland-dev
libxkbcommon-dev` — a bare box will report GUI typecheck failures that aren't your change.

CI (`.github/workflows/ci.yml`) runs `make lint` and `make test` on non-draft PRs only.
Integration tests are **not** in CI — run them locally when touching BPF or enforcement.

## Architecture

### Layering (ports & adapters)

```
main.go → cmd/root.go (cobra)
  cmd/functions/<mode>/   CLI wiring + TUI/headless/serve plumbing per subcommand
    internal/usecase/     one orchestrator per mode; owns lifecycle & ordering
      internal/repository/ the ports: MonitorRepository, GuardRepository,
                           NetworkMonitorRepository, NetworkGuardRepository
        internal/monitor, internal/guard,
        internal/networkmonitor, internal/networkguard   the eBPF engines (adapters)
```

`internal/infrastructure` (imported as `ebpf`) is the shared kernel below all engines:
event structs (`FileEvent`, `NetEvent`), `EventType`/`NetEventType` enums + parsing,
`BinaryEntry`/`ComputeBinaryEntry` (binary identity), `Target` (watch-path model), and
the **BPF-LSM preflight** (`CheckBPFLSM`) that every enforcing mode runs before attach —
a hook that attaches on a kernel without `bpf` in the active LSM list silently denies
nothing, so the process refuses to start instead.

### eBPF attach model

- **kprobes** (`monitor`) — observe all I/O regardless of syscall path (io_uring, splice,
  sendfile, mmap) plus metadata ops.
- **LSM hooks** (`guard`, `network-guard`, `daemon`) — the only kernel mechanism that can
  **deny**. ~23 hooks; only `file_open` + `file_permission` are mandatory, the rest are
  best-effort (a missing hook logs a warning, enforcement continues).
- **tracepoints / kretprobes** (`network-monitor`) — TCP/UDP/DNS.

**Identity is the executable's inode, never its name or comm** — renaming or comm-spoofing
a binary cannot move it past policy. Guard whitelist/blacklist entries, the network guard,
and the daemon all key on `dev:ino`.

### The daemon (`internal/usecase/daemon.go`, `cmd/functions/daemon`)

Config-driven (`/etc/app-listener/daemon.conf`, parsed by `internal/daemonconfig`), it
runs the guard whitelist engine over many directories plus an fscrypt lifecycle. The
invariant everywhere: **resources are never readable without a live guard attached**.
Ordering is `attach → unlock → populate inodes → resolve deferred binaries → re-sync`;
shutdown deprovisions fscrypt keys in passes *while guards still deny*. LSM links are
pinned under `/sys/fs/bpf` so a `SIGKILL` leaves trees enforced until `ExecStopPost`
re-locks the vaults. `systemctl reload` (SIGHUP) recomputes every binary's inode identity
atomically — new guards attach before old ones detach; a broken config keeps the running
one. `selfguards.go` makes the daemon guard its own `/etc/app-listener` config + key +
edit-auth hash. When `edit-protected` has a password configured, `control.go` opens the
live-edit control socket and `DaemonUseCase.GrantEditAccess` / the guard's
`GrantSelfEditAccess`+`RevokeSelfEditAccess` implement the transient per-resource write
grant. The `GuardRepository` port has daemon-specific methods (`PopulateInodes`,
`ResolvePendingBinaries`, `ReSyncBinaries`, `SweepInodes`, `GrantSelfEditAccess`,
`RevokeSelfEditAccess`) — read their doc comments in
`internal/repository/repository.go` before changing reload/resync behaviour.

### Install / lifecycle packages

- `internal/install` — the TUI install wizard and the built-in **catalog** (`catalog.go`)
  of sensitive directories (SSH, GPG, AI agents, browsers, password stores). The catalog
  narrows watches to the sensitive subtrees so self-updating apps' vault roots stay writable.
- `internal/fscrypt` — encryption lifecycle, xattr policy checks, orphan/migrate handling.
- `internal/protected` — shared "is the daemon running / which dirs are encrypted with the
  master key" checks that gate `install`, `uninstall`, `edit-protected`.
- `internal/systemd` — generates the systemd units, package-manager catalog-refresh hooks
  (pacman `PostTransaction`, apt `DPkg::Post-Invoke`), per-user ssh-agent unit, boot-time
  refresh service.
- `scripts/install.sh` is the `curl … | sudo bash` installer: `check-compatibility` +
  Ed25519 signature / checksum / GitHub asset-digest verification before atomic install.

### Presentation

Every mode has a bubbletea TUI (`internal/tui/*_model.go`), a `--headless` stderr writer
with a stable `EVENT|` / `GUARD|` / `DAEMON …` / `NETEVENT|` line format that the
integration tests parse, and `--serve[=host:port]` which mirrors the event stream
read-only over WebSockets (`internal/tui/serve.go`). `--gui` (build-tag `gui`) is
fyne, monitor-only; `internal/gui/stub.go` is the no-op for non-gui builds.
Verbosity ladder lives in `internal/constants`; `--verbose` requires `--headless`.

## Release

See `memory/release-promote-flow.md`. A release is a *build*, not a tag rename:
`release.yml` cuts a `pre-*` pre-release per push to main; `promote-release.yml` rebuilds
a stable `vX.Y.Z`. `_build-release.yml` is the reusable build job.
