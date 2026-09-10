# Toolchain image used by 'make build' (and 'make build-image'). Only the
# tools needed to regenerate BPF bindings and compile the Go binary live
# here — the source is mounted read-write at build time, so this image is
# intentionally free of any COPY of the repository.
FROM golang:1.26-bookworm

# BPF code generation is highly sensitive to the clang version: Debian
# bookworm's stock clang 14 emits a `guard_path_rename` that overflows the
# BPF verifier's 1,000,000-instruction budget on modern kernels, so the
# container build fails `daemon --check` where an on-host build with a
# current clang passes (issue #45). Pin a current clang from apt.llvm.org so
# `make build`, CI and every release binary match a modern host toolchain
# (clang 22.x) and the guard programs verify.
ARG LLVM_VERSION=22

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        wget \
        gnupg \
        ca-certificates \
        lsb-release \
        software-properties-common \
    && wget -qO /tmp/llvm.sh https://apt.llvm.org/llvm.sh \
    && chmod +x /tmp/llvm.sh \
    && /tmp/llvm.sh ${LLVM_VERSION} \
    && apt-get install -y --no-install-recommends llvm-${LLVM_VERSION} \
    && ln -sf /usr/bin/clang-${LLVM_VERSION}      /usr/bin/clang \
    && ln -sf /usr/bin/clang++-${LLVM_VERSION}    /usr/bin/clang++ \
    && ln -sf /usr/bin/llvm-strip-${LLVM_VERSION} /usr/bin/llvm-strip \
    && ln -sf /usr/bin/llvm-objcopy-${LLVM_VERSION} /usr/bin/llvm-objcopy \
    && apt-get install -y --no-install-recommends \
        bpftool \
        libbpf-dev \
        libc6-dev \
        libc6-dev-i386 \
        make \
        gcc \
        pkg-config \
        libgl1-mesa-dev \
        xorg-dev \
        libwayland-dev \
        libxkbcommon-dev \
        git \
    && rm -rf /var/lib/apt/lists/* /tmp/llvm.sh

# Fail the image build if clang did not end up at the pinned major version.
RUN clang --version | grep -qE "clang version ${LLVM_VERSION}\." \
    || { echo "ERROR: expected clang ${LLVM_VERSION}, got: $(clang --version | head -1)"; exit 1; }
