# Dockerfile explicitly designed constructing minimal attack surfaces via static 
# compilation mapping directly onto unprivileged Distroless boundaries natively.

# ------------------------------------------------------------------------------
# Phase 1: Compilation Environment Configuration
# ------------------------------------------------------------------------------
FROM golang:1.22-alpine AS builder

# Establish working directory constraints.
WORKDIR /workspace

# Install ca-certificates configuring secure dependency acquisition boundaries.
RUN apk add --no-cache git ca-certificates tzdata

# Copy module definitions resolving dependencies exclusively.
COPY go.mod go.sum ./
RUN go mod download -x

# Ingest application source arrays mapping exactly avoiding external modifications.
COPY cmd/ cmd/
COPY internal/ internal/
COPY pkg/ pkg/

# Compile the static executable aggressively executing completely offline.
# CGO_ENABLED=0 removes dynamic library assignments completely avoiding glibc.
# -ldflags="... -s -w" removes symbol arrays executing smaller binary sizes successfully.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -o netsentry \
    -ldflags="-extldflags '-static' -s -w" \
    ./cmd/netsentry

# ------------------------------------------------------------------------------
# Phase 2: Runtime Execution Boundary
# ------------------------------------------------------------------------------
# Utilize explicitly unprivileged scratch-equivalent bounds maintaining security.
FROM gcr.io/distroless/static-debian11:nonroot

# Transfer compiled executable boundaries avoiding external structural dependencies.
COPY --from=builder /workspace/netsentry /netsentry
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Force execution natively as unprivileged user (1000:1000).
USER nonroot:nonroot

# Explicitly expose HTTP REST API listening boundary port mapping globally.
EXPOSE 8080/tcp

# Instruct container execution logic establishing binary directly.
ENTRYPOINT ["/netsentry"]

# Define standard argument bounds instantiating the HTTP API listener.
CMD ["serve", "--addr", "0.0.0.0:8080"]
