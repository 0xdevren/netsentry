#!/usr/bin/env bash
set -euo pipefail

BINARY_NAME="netsentry"
CMD_PATH="./cmd/netsentry"
MODULE="github.com/0xdevren/netsentry"
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS="-s -w -X ${MODULE}/internal/app.Version=${VERSION} -X ${MODULE}/internal/app.Commit=${COMMIT} -X ${MODULE}/internal/app.BuildDate=${BUILD_DATE}"

mkdir -p bin

echo "Building ${BINARY_NAME} ${VERSION} ..."
go build -trimpath -ldflags "${LDFLAGS}" -o "bin/${BINARY_NAME}" "${CMD_PATH}"
echo "Binary: bin/${BINARY_NAME}"

echo "Cross-compiling (linux/amd64) ..."
GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "${LDFLAGS}" -o "bin/${BINARY_NAME}-linux-amd64" "${CMD_PATH}"

echo "Cross-compiling (darwin/arm64) ..."
GOOS=darwin GOARCH=arm64 go build -trimpath -ldflags "${LDFLAGS}" -o "bin/${BINARY_NAME}-darwin-arm64" "${CMD_PATH}"

echo "Build complete."
