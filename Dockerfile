FROM golang:1.26-bookworm AS builder

WORKDIR /app/app-listener

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN make build-linux

FROM debian:bookworm-slim AS runner

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
    ca-certificates \
    dumb-init \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/app-listener/build/linux/app-listener /app/app-listener

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["sleep", "infinity"]
