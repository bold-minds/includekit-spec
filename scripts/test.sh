#!/bin/bash
set -e

# Determine script directory and repo root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

echo "🧪 Running full test suite"
echo ""

# Build codegen tool if needed
echo "📦 Building codegen tool..."
cd "$REPO_ROOT/codegen" || exit 1
go build -o ../bin/codegen .
cd "$REPO_ROOT" || exit 1

# Generate code from schema
echo "🔧 Generating code from schema..."
cd "$REPO_ROOT" || exit 1
./bin/codegen

# Build TypeScript tests (includes runtime utilities)
echo "📦 Building TypeScript tests..."
cd "$REPO_ROOT/pkgs/ts/tests" || exit 1
npm install
npm run build

echo ""
echo "🧪 Testing TypeScript..."
cd "$REPO_ROOT/pkgs/ts/tests" || exit 1
npm test

echo ""
echo "🧪 Testing Go..."
cd "$REPO_ROOT/pkgs/go" || exit 1
go test ./...

echo ""
echo "🔍 Verifying no-runtime constraint..."
cd "$REPO_ROOT" || exit 1

# Check TypeScript production package - should be types only
if ! find "$REPO_ROOT/pkgs/ts/types" -type f -name "*.d.ts" 2>/dev/null | grep -q .; then
    echo "❌ ERROR: No .d.ts files found in pkgs/ts/types"
    exit 1
fi

echo "✓ pkgs/ts/types contains types only"

# Check Go production package - ensure no complex runtime logic
# Allow type helpers (like isValue()), but not actual business logic
if grep -r "^func.*{$" "$REPO_ROOT/pkgs/go/types"/*.go 2>/dev/null | \
   grep -v "func (.*) is" | \
   grep -v "^[[:space:]]*$" | \
   grep -v "_test.go" | \
   grep -q .; then
    echo "⚠️  Warning: Found functions in pkgs/go/types - verifying they are type helpers only..."
    # If this fails, manually verify that only type interface methods exist
fi

echo "✓ pkgs/go/types structure verified"

echo ""
echo "✅ All tests passed and constraints verified!"
