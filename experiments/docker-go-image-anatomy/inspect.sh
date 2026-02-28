#!/usr/bin/env bash
# inspect.sh — 公式Goコンテナイメージの内部調査スクリプト
# Usage: docker run --rm golang:1.26 bash /work/inspect.sh
#   or:  docker build -t go-image-anatomy . && docker run --rm go-image-anatomy
set -euo pipefail

separator() {
  echo ""
  echo "================================================================"
  echo "  $1"
  echo "================================================================"
  echo ""
}

# --- Go toolchain ---
separator "Go Version"
go version

separator "Go Environment (go env)"
go env

separator "Go Tools (go tool)"
go tool -n

separator "Go binaries in GOROOT/bin"
ls -la "$(go env GOROOT)/bin/"

# --- Directory layout ---
separator "GOROOT directory tree (depth=2)"
find "$(go env GOROOT)" -maxdepth 2 -type d | sort

separator "GOPATH directory tree (depth=2)"
find "$(go env GOPATH)" -maxdepth 2 -type d 2>/dev/null | sort || echo "(GOPATH is empty)"

# --- Base OS ---
separator "OS Release"
cat /etc/os-release

separator "Kernel Info"
uname -a

# --- System packages ---
separator "Installed Packages (dpkg)"
dpkg -l 2>/dev/null | tail -n +6 | awk '{print $2, $3}' | sort || echo "(dpkg not available — alpine variant?)"

# --- libc / shared libraries ---
separator "glibc Version"
ldd --version 2>&1 | head -1 || echo "(ldd not available)"

separator "Shared Libraries for Go binary"
ldd "$(go env GOROOT)/bin/go" 2>/dev/null || echo "(statically linked or ldd not available)"

# --- Environment variables ---
separator "Environment Variables"
env | sort

separator "Inspection complete."
