# Development

`CGO_ENABLED=1` is required for every build: the daemon links `google/fscrypt`, which needs cgo.
Contributor-oriented notes (architecture, the BPF verifier budget, test conventions) live in
[`CLAUDE.md`](../CLAUDE.md).

## Makefile targets

| Target | Description |
|--------|-------------|
| `make build` | Dockerized build (rootful Docker): dump `vmlinux.h` from the host BTF, regenerate BPF bindings, build `build/linux/app-listener`. `GUI=1` links the desktop GUI |
| `make build-host` | Same on the host; needs a recent clang/LLVM, `bpftool`, Go 1.26+, GCC. Also takes `GUI=1` |
| `make build-linux` | Go compile only, assuming generated files exist. `GUI=1` links the desktop GUI |
| `make generate` | Regenerate BPF bindings (needs `bpftool` and a readable `/sys/kernel/btf/vmlinux`) |
| `make check-compatibility` | Static host check: can it run app-listener? |
| `make test` | Unit tests |
| `make test-integration` | Docker integration and bypass suite (rootful Docker) |
| `make lint` | golangci-lint |
| `make vuln` | govulncheck over the module (`GUI=1` scans the GUI build; run `make generate` first) |
| `make pprof` | Save CPU, heap, allocation and goroutine profiles of a daemon started with `--pprof` |
| `make clean` | Remove build artifacts |

The eBPF is compiled against the build host's BTF, so builds are native only: build arm64 on an
arm64 host.

## Profiling

`make pprof` sees only the daemon's Go side. The guard hooks run inside every process that reads,
writes or stats a file, and the BPF maps are kernel memory, so neither shows up there.
`sudo bash scripts/kernel-prof.sh [seconds]` records that part while you keep using the machine:
per-program BPF run time, map memory and `guard_inodes` occupancy, and a `perf` profile of the
daemon with kernel stacks. Output goes to `build/pprof/kernel-<timestamp>/`. Profile on the
weakest hardware you target: hook cost per call grows on slow CPUs, map memory does not shrink
with RAM.

## Tests

- `make test` runs the unit tests. They never need root.
- `make test-integration` builds the C exploit corpus in
  [`integrationtests/exploits`](../integrationtests/exploits/README.md) (deliberate bypass attempts:
  io_uring, `open_by_handle_at`, `copy_file_range`, `process_vm_readv`, raw block device, mount,
  `SCM_RIGHTS` fd passing, ...), then runs the suite in privileged containers. It needs rootful
  Docker and skips when `docker` is not on `PATH`.

## Docker

```bash
sudo make deploy
docker exec -it app-listener app-listener monitor /tmp
```

`make deploy` builds the binary inside Docker, on the toolchain image (`docker/builder.Dockerfile`),
so nothing is compiled or written on the host, and starts an idle container to `docker exec` into
(`/app` is the working directory and on `PATH`). It needs rootful Docker: a rootless container can't
load eBPF, so `deploy` refuses it (`sudo` drops a rootless `DOCKER_HOST`). The container is
privileged and gets the host's `/sys/kernel/security` read-only, so the guard modes still refuse to
start on a kernel without BPF-LSM. `$MONITOR_DIR` (default `./data`) is mounted read-only at
`/data`. `make deploy-down` stops it.

## Demo GIFs

The README's media lives in [`media/`](../media): terminal recordings in `media/tui/`, desktop
GUI screenshots and its walkthrough GIF in `media/gui/`. Keep each GIF short (under a minute) and
small (a few hundred KB). For a terminal recording, scale to about 900 px wide at 6 fps with a
reduced palette, for example:

```bash
ffmpeg -i recording.mp4 -vf "fps=6,scale=880:-1:flags=lanczos,split[a][b];[a]palettegen=max_colors=64:stats_mode=diff[p];[b][p]paletteuse=dither=none:diff_mode=rectangle" media/tui/<mode>.gif
```

The GUI media is rendered off-screen from fake engines, at the floating window's default
1280x820, by the GUI end-to-end tests (`-tags gui,e2e`, kept out of `make test` for their run
time; they need neither root nor a display). `make test-e2e` renders into a temp dir;
`SNAPSHOTS=<dir>` keeps one PNG per page and state, `DEMO=<dir>` the walkthrough's frames and a
concat list with each frame's duration:

```bash
mkdir -p /tmp/snaps /tmp/reel
make test-e2e SNAPSHOTS=/tmp/snaps DEMO=/tmp/reel
cd /tmp/reel && ffmpeg -f concat -i frames.txt -fps_mode vfr -vf "split[a][b];[a]palettegen=max_colors=128:stats_mode=diff[p];[b][p]paletteuse=dither=none:diff_mode=rectangle" demo.gif
```

Keep the GUI GIF at its native width: it is sharper and smaller than a scaled one.
