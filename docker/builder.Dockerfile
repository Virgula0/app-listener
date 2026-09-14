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

# apt.llvm.org (LLVM's own apt repository, not a CDN-backed mirror) and the
# Debian mirror network are occasionally slow or briefly unreachable from CI
# runners; llvm.sh's own preflight (`check_url`, HEAD-probing
# apt.llvm.org/bookworm/) aborts the *whole* script on the first such blip
# with no retry of its own. Observed in CI as this RUN failing outright on a
# push, then succeeding unchanged on a manual re-run (the image is rebuilt
# from scratch every run — nothing here is cached between CI invocations).
# Wrap every network-dependent command below in a short retry/backoff loop
# instead of gambling the whole image build on one shot at these endpoints.
RUN retry() { \
        n=0; max=5; \
        until "$@"; do \
            n=$((n + 1)); \
            if [ "$n" -ge "$max" ]; then \
                echo "retry: giving up after $max attempts: $*" >&2; \
                return 1; \
            fi; \
            echo "retry: attempt $n/$max failed, retrying in $((n * 5))s: $*" >&2; \
            sleep $((n * 5)); \
        done; \
    }; \
    retry apt-get update \
    && retry apt-get install -y --no-install-recommends \
        wget \
        gnupg \
        ca-certificates \
        lsb-release \
        software-properties-common \
    && retry wget -qO /tmp/llvm.sh https://apt.llvm.org/llvm.sh \
    && chmod +x /tmp/llvm.sh \
    && retry /tmp/llvm.sh ${LLVM_VERSION} \
    && retry apt-get install -y --no-install-recommends llvm-${LLVM_VERSION} \
    && ln -sf /usr/bin/clang-${LLVM_VERSION}      /usr/bin/clang \
    && ln -sf /usr/bin/clang++-${LLVM_VERSION}    /usr/bin/clang++ \
    && ln -sf /usr/bin/llvm-strip-${LLVM_VERSION} /usr/bin/llvm-strip \
    && ln -sf /usr/bin/llvm-objcopy-${LLVM_VERSION} /usr/bin/llvm-objcopy \
    && retry apt-get install -y --no-install-recommends \
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
