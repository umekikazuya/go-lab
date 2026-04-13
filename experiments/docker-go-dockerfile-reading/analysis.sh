#!/usr/bin/env bash
# analysis.sh — 公式GoコンテナイメージのDockerfile構造解析スクリプト
# Usage: bash analysis.sh
#
# docker-library/golang リポジトリからDockerfileを取得し、構造を解析する。
set -euo pipefail

WORK_DIR="$(mktemp -d)"
REPO_URL="https://github.com/docker-library/golang.git"
trap 'rm -rf "$WORK_DIR"' EXIT

separator() {
  echo ""
  echo "================================================================"
  echo "  $1"
  echo "================================================================"
  echo ""
}

# --- Dockerfileの取得 ---
separator "Cloning docker-library/golang"
git clone --depth 1 "$REPO_URL" "$WORK_DIR/golang"

# --- Dockerfileの一覧 ---
separator "Dockerfile variants"
find "$WORK_DIR/golang" -name "Dockerfile" -type f | sort

# --- bookworm Dockerfileの解析 ---
BOOKWORM_DF="$(find "$WORK_DIR/golang" -path "*/bookworm/Dockerfile" -type f | head -1)"
if [ -z "$BOOKWORM_DF" ]; then
  echo "bookworm Dockerfile not found, searching for default..."
  BOOKWORM_DF="$(find "$WORK_DIR/golang" -name "Dockerfile" -type f | head -1)"
fi

separator "Analyzing: $BOOKWORM_DF"

separator "ENV directives"
grep -n '^ENV\b' "$BOOKWORM_DF" || echo "(no ENV found)"

separator "ARG directives"
grep -n '^ARG\b' "$BOOKWORM_DF" || echo "(no ARG found)"

separator "FROM directives"
grep -n '^FROM\b' "$BOOKWORM_DF" || echo "(no FROM found)"

separator "RUN directives (with line numbers)"
grep -n '^RUN\b' "$BOOKWORM_DF" || echo "(no RUN found)"

separator "Download URLs"
grep -noE 'https?://[^ "]+' "$BOOKWORM_DF" || echo "(no URLs found)"

separator "SHA256 / checksum references"
grep -ni 'sha256\|checksum\|gpg\|signature' "$BOOKWORM_DF" || echo "(no checksum references found)"

separator "Full Dockerfile content"
cat -n "$BOOKWORM_DF"

# --- alpine Dockerfileの解析 ---
ALPINE_DF="$(find "$WORK_DIR/golang" -path "*/alpine/Dockerfile" -type f | head -1)"
if [ -n "$ALPINE_DF" ]; then
  separator "Alpine Dockerfile: ENV directives"
  grep -n '^ENV\b' "$ALPINE_DF" || echo "(no ENV found)"

  separator "Alpine Dockerfile: FROM directives"
  grep -n '^FROM\b' "$ALPINE_DF" || echo "(no FROM found)"

  separator "Alpine Dockerfile: Package installation (apk)"
  grep -n 'apk' "$ALPINE_DF" || echo "(no apk commands found)"

  separator "Diff: bookworm vs alpine (ENV lines)"
  diff <(grep '^ENV' "$BOOKWORM_DF" | sort) <(grep '^ENV' "$ALPINE_DF" | sort) || echo "(differences found above)"
else
  separator "Alpine Dockerfile not found"
fi

separator "Analysis complete."
