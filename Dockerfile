# The builder stage is the pinned toolchain from docker/builder.Dockerfile (clang 22, bpftool;
# stock clang 14 emits a guard that overflows the verifier, issue #45). `make deploy` builds it
# first. It runs `build-host`, not `build`: `build` starts its own container, and there is no
# Docker inside an image build. The eBPF still targets this host: RUN steps see the host's
# /sys/kernel/btf/vmlinux.
ARG BUILDER_IMAGE=app-listener-builder:local
FROM ${BUILDER_IMAGE} AS builder

WORKDIR /app/app-listener

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN make build-host

FROM debian:bookworm-slim AS runner

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
    ca-certificates \
    dumb-init \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
ENV PATH="/app:${PATH}"

COPY --from=builder /app/app-listener/build/linux/app-listener /app/app-listener

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
# dumb-init only reaps and forwards signals to its child and exits with it, so it needs one
# that stays up: `docker exec -it app-listener app-listener <mode>`.
CMD ["sleep", "infinity"]
