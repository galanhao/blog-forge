#!/bin/bash
set -e

echo "=== blog-forge Harness Initialization ==="

# Go version manager
if [ -f /root/.g/env ]; then
  source /root/.g/env
  echo "✅ g version manager loaded"
else
  echo "⚠️  g version manager not found"
fi

# Go version check
GO_VERSION=$(go version 2>/dev/null || echo "NOT FOUND")
echo "Go: $GO_VERSION"

# Run tests
echo ""
echo "=== Running Tests ==="
if go test ./... 2>&1; then
  echo "✅ All tests pass"
else
  echo "❌ Tests FAILED — fix before proceeding"
  exit 1
fi

# Build check
echo ""
echo "=== Build Check ==="
if GOPROXY=https://goproxy.cn,direct GONOSUMCHECK=* go build -o blog-forge ./cmd/blog-forge/ 2>&1; then
  echo "✅ Build succeeds"
else
  echo "❌ Build FAILED — fix before proceeding"
  exit 1
fi

# Site build check
echo ""
echo "=== Site Build Check ==="
if ./blog-forge build 2>&1; then
  echo "✅ Site built to dist/"
else
  echo "❌ Site build FAILED"
  exit 1
fi

echo ""
echo "=== Verification Complete ==="
echo ""
echo "Next steps:"
echo "1. Read feature_list.json to see current feature state"
echo "2. Pick ONE unfinished feature to work on"
echo "3. Implement only that feature"
echo "4. Re-run ./init.sh before claiming done"
