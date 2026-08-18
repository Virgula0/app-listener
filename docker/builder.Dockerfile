# Toolchain image used by 'make build' (and 'make build-image'). Only the
# tools needed to regenerate BPF bindings and compile the Go binary live
# here — the source is mounted read-write at build time, so this image is
# intentionally free of any COPY of the repository.
FROM golang:1.26-bookworm

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        clang \
        llvm \
        lld \
        bpftool \
        libbpf-dev \
        libc6-dev \
        libc6-dev-i386 \
        make \
        gcc \
        libc6-dev \
        pkg-config \
        libgl1-mesa-dev \
        xorg-dev \
        libwayland-dev \
        libxkbcommon-dev \
        git \
        ca-certificates \
    && rm -rf /var/lib/apt/lists/*