#!/bin/bash
set -e

echo "🚀 Setting up includekit-spec repository"
echo ""

# Check prerequisites
echo "📋 Checking prerequisites..."

# Check Go
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go >= 1.22"
    exit 1
fi
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo "✓ Go $GO_VERSION found"

# Check Node.js
if ! command -v node &> /dev/null; then
    echo "❌ Node.js is not installed. Please install Node.js >= 20"
    exit 1
fi
NODE_VERSION=$(node --version)
echo "✓ Node.js $NODE_VERSION found"

# Check npm
if ! command -v npm &> /dev/null; then
    echo "❌ npm is not installed. Please install npm"
    exit 1
fi
NPM_VERSION=$(npm --version)
echo "✓ npm $NPM_VERSION found"

echo ""
echo "📦 Installing dependencies..."

# Install Go dependencies
echo "Installing Go modules..."
go mod download
echo "✓ Go dependencies installed"

# Install Node.js dependencies
echo "Installing npm packages..."
npm install
echo "✓ npm dependencies installed"

echo ""
echo "🧪 Running full test suite (includes build)..."

# Run all tests (includes build + codegen + verification)
./scripts/test.sh

echo ""
echo "✅ Setup complete!"
echo ""
echo "📝 Quick reference:"
echo "  ./scripts/test.sh     - Run full test suite (main command)"
echo "  ./scripts/build.sh    - Just build (without tests)"
echo ""
echo "💡 In normal development, just run: ./scripts/test.sh"
echo "📖 See QUICKSTART.md for more details"
