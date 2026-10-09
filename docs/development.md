# Development

`CGO_ENABLED=1` is required for every build: the daemon links `google/fscrypt`, which needs cgo.
Contributor-oriented notes (architecture, the BPF verifier budget, test conventions) live in
[`CLAUDE.md`](../CLAUDE.md).

## Makefile targets

| Target | Description |
|--------|-------------|
| `make build` | Dockerized build (rootful Docker): dump `vmlinux.h` from the host BTF, regenerate BPF bindings, build `build/linux/app-listener` |
| `make build-host` | Same on the host; needs a recent clang/LLVM, `bpftool`, Go 1.26+, GCC |
| `make build-linux` | Go compile only, assuming generated files exist. `GUI=1` links the desktop GUI |
| `make generate` | Regenerate BPF bindings (needs `bpftool` and a readable `/sys/kernel/btf/vmlinux`) |
| `make check-compatibility` | Static host check: can it run app-listener? |
| `make test` | Unit tests |
| `make test-integration` | Docker integration and bypass suite (rootful Docker) |
| `make lint` | golangci-lint |
| `make vuln` | govulncheck over the module (`GUI=1` scans the GUI build; run `make generate` first) |
| `make clean` | Remove build artifacts |

The eBPF is compiled against the build host's BTF, so builds are native only: build arm64 on an
arm64 host.

## Tests

- `make test` runs the unit tests. They never need root.
- `make test-integration` builds the C exploit corpus in
  [`integrationtests/exploits`](../integrationtests/exploits/README.md) (deliberate bypass attempts:
  io_uring, `open_by_handle_at`, `copy_file_range`, `process_vm_readv`, raw block device, mount,
  `SCM_RIGHTS` fd passing, ...), then runs the suite in privileged containers. It needs rootful
  Docker and skips when `docker` is not on `PATH`.

## Docker

```bash
docker compose build
docker compose run --rm app-listener monitor -w /tmp
```

Multi-stage build on a Debian-slim runner. eBPF needs the host kernel: run privileged, or with
`/sys/kernel/btf/` mounted and `CAP_BPF`. `make build` uses a separate toolchain image
(`docker/builder.Dockerfile`), not this one.

## Demo GIFs

The per-mode GIFs in [`media/`](../media) are referenced by the README. Keep each one short
(under a minute) and small (a few hundred KB): scale to about 900 px wide at 6 fps with a reduced
palette, for example:

```bash
ffmpeg -i recording.mp4 -vf "fps=6,scale=880:-1:flags=lanczos,split[a][b];[a]palettegen=max_colors=64:stats_mode=diff[p];[b][p]paletteuse=dither=none:diff_mode=rectangle" media/<mode>.gif
```
