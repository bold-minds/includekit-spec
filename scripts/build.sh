#!/bin/bash
set -e

echo "🚀 Building includekit-types"

# Build the Go codegen tool
echo "📦 Building codegen tool..."
go build -o bin/codegen ./codegen

# Run codegen for all languages
echo "🔧 Generating code from schema..."
./bin/codegen -schema schema/v0-1-0.json -output pkgs

# Build TypeScript testkit
echo "📦 Building TypeScript testkit..."
cd pkgs/ts/tests && npm run build && cd ../..

echo "✅ Build complete!"
