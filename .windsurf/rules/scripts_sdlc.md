---
trigger: always_on
---

# SDLC Scripts Convention

Use the scripts folder for high-level SDLC automation. Prefer this approach over Makefiles or similar build tools.

## Core Scripts Structure

Maintain these three primary scripts:

### s`cripts/setup.sh`
First-time repository setup:
- Check prerequisites (Go, Node.js, etc.)
- Install dependencies (go mod download, npm install)
- Run initial build/test to verify setup
- Display quick reference and next steps

### `scripts/build.sh`
Build the project without testing:
- Build binaries and tools
- Generate code from schemas
- Compile packages
- Pure build operations only

### `scripts/test.sh`
Primary development command - run full test suite:
- Include build steps (build what's needed)
- Run tests for all languages/components
- Verify constraints and invariants
- This should be the main command developers run

## Script Standards

All scripts must:
- Use bash with `#!/bin/bash` shebang
- Include `set -e` for proper error handling
- Use emoji prefixes for output (🚀, 📦, 🧪, ✅, ❌)
- Provide clear, user-friendly messages
- Be executable (`chmod +x`)

## Additional Scripts

Other scripts can live in `scripts/` if they:
- Are high-level SDLC operations (deploy, release, lint, format, etc.)
- Are frequently used by developers
- Naturally fit into the software development lifecycle
- Provide clear value over one-off commands

## Why Not Makefiles

This project uses bash scripts instead of Makefiles for:
- Better cross-platform compatibility
- More readable syntax
- Easier error handling and user feedback
- Direct integration with shell tools and commands