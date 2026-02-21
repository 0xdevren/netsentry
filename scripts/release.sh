#!/usr/bin/env bash
set -euo pipefail

VERSION=${1:-$(git describe --tags --always --dirty 2>/dev/null || echo "dev")}
echo "Releasing netsentry ${VERSION} ..."

bash scripts/build.sh

echo "Checksums:"
sha256sum bin/netsentry-* | tee bin/SHA256SUMS

if command -v gh &>/dev/null; then
  gh release create "${VERSION}" bin/netsentry-* bin/SHA256SUMS \
    --title "netsentry ${VERSION}" \
    --notes "Release ${VERSION}"
  echo "GitHub release created: ${VERSION}"
else
  echo "gh CLI not found – skipping GitHub release creation."
  echo "Binaries ready in bin/"
fi
