#!/bin/bash
# Build script for imgconvert

set -e

# Setup Go path
export PATH="$HOME/gosdk/bin:$PATH"

echo "Building imgconvert..."
go build -o imgconvert -ldflags "-s -w" ./cmd/imgconvert

echo "✓ Build complete: ./imgconvert"
echo ""
echo "Test it:"
echo "  ./imgconvert --help"
echo "  ./imgconvert --version"
echo ""
echo "Quick test:"
echo "  ./imgconvert test/fixtures/sample.jpg"
