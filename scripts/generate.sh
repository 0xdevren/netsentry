#!/usr/bin/env bash
set -euo pipefail
echo "Running go generate ..."
go generate ./...
echo "Done."
