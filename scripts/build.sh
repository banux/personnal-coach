#!/bin/sh
# Build script: compiles frontend then embeds into Go backend binary
set -e

ROOT="$(cd "$(dirname "$0")/.." && pwd)"

echo "==> Building frontend..."
cd "$ROOT/frontend"
npm ci --silent
npm run build

echo "==> Copying dist to backend for embedding..."
rm -rf "$ROOT/backend/dist"
cp -r "$ROOT/frontend/dist" "$ROOT/backend/dist"

echo "==> Building Go backend with embedded frontend..."
cd "$ROOT/backend"
GOROOT=/home/banux/go GOPATH=/home/banux/go go build -ldflags="-s -w" -o "$ROOT/personal-coach" .

echo ""
echo "Build complete! Run: ANTHROPIC_API_KEY=your-key $ROOT/personal-coach"
