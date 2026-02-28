#!/usr/bin/env bash
# reinstall.sh — Reinstall dependencies for all example projects.
# Usage: bash examples/reinstall.sh

# Get current Node and Bun paths and add them to PATH
NODE_BIN_DIR=$(dirname "$(which node 2>/dev/null || echo '/usr/local/bin')")
BUN_BIN_DIR="$HOME/.bun/bin"
DENO_BIN_DIR="$HOME/.deno/bin"
LOCAL_BIN_DIR="$HOME/.local/bin"

export PATH="$NODE_BIN_DIR:$BUN_BIN_DIR:$DENO_BIN_DIR:$LOCAL_BIN_DIR:$PATH"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

GREEN='\033[0;32m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

step() { echo -e "\n${CYAN}━━━ $1 ━━━${NC}"; }
ok()   { echo -e "${GREEN}  ✓ $1${NC}"; }
warn() { echo -e "${YELLOW}  ⚠ $1${NC}"; }
fail() { 
  echo -e "${RED}  ✗ $1${NC}"
  echo -e "${YELLOW}      PATH atual: $PATH${NC}"
  # Check if npm/pnpm/bun exist in the PATH
  command -v npm >/dev/null || echo -e "${RED}      npm não encontrado no PATH${NC}"
  command -v pnpm >/dev/null || echo -e "${RED}      pnpm não encontrado no PATH${NC}"
  command -v bun >/dev/null || echo -e "${RED}      bun não encontrado no PATH${NC}"
}

# ─────────────────────────────────────────────
# 1. npm projects
# ─────────────────────────────────────────────
if [ -d "node-npm" ]; then
  step "node-npm (npm install)"
  if (cd node-npm && npm install 2>&1); then
    ok "done ($(du -sh node-npm/node_modules 2>/dev/null | cut -f1))"
  else
    fail "npm install failed"
  fi
fi

if [ -d "nextjs-project" ]; then
  step "nextjs-project (npm install + next build)"
  if (cd nextjs-project && npm install 2>&1 && npx next build 2>&1); then
    ok "done (node_modules: $(du -sh nextjs-project/node_modules 2>/dev/null | cut -f1), .next: $(du -sh nextjs-project/.next 2>/dev/null | cut -f1))"
  else
    fail "nextjs install/build failed"
  fi
fi

# ─────────────────────────────────────────────
# 2. pnpm projects
# ─────────────────────────────────────────────
if [ -d "node-pnpm" ]; then
  step "node-pnpm (pnpm install)"
  if (cd node-pnpm && pnpm install 2>&1); then
    ok "done ($(du -sh node-pnpm/node_modules 2>/dev/null | cut -f1))"
  else
    fail "pnpm install failed"
  fi
fi

# ─────────────────────────────────────────────
# 3. bun projects
# ─────────────────────────────────────────────
if [ -d "node-bun" ]; then
  step "node-bun (bun install)"
  if (cd node-bun && bun install 2>&1); then
    ok "done ($(du -sh node-bun/node_modules 2>/dev/null | cut -f1))"
  else
    fail "bun install failed"
  fi
fi

# ─────────────────────────────────────────────
# 4. deno projects
# ─────────────────────────────────────────────
if [ -d "deno-project" ]; then
  step "deno-project (deno install)"
  if (cd deno-project && deno install 2>&1); then
    ok "done"
  else
    fail "deno install failed"
  fi
fi

# ─────────────────────────────────────────────
# 5. rust projects
# ─────────────────────────────────────────────
if [ -d "rust-project" ]; then
  step "rust-project (cargo build)"
  if (cd rust-project && cargo build 2>&1); then
    ok "done ($(du -sh rust-project/target 2>/dev/null | cut -f1))"
  else
    fail "cargo build failed"
  fi
fi

# ─────────────────────────────────────────────
# 6. go projects
# ─────────────────────────────────────────────
if [ -d "go-project" ]; then
  step "go-project (go build)"
  if (cd go-project && go build -o go-project . 2>&1); then
    ok "done"
  else
    fail "go build failed"
  fi
fi

# ─────────────────────────────────────────────
# 7. python projects
# ─────────────────────────────────────────────
if [ -d "python-project" ]; then
  step "python-project (venv + pip install)"
  if (cd python-project && python3 -m venv .venv && .venv/bin/pip install -r requirements.txt 2>&1); then
    ok "done ($(du -sh python-project/.venv 2>/dev/null | cut -f1))"
  else
    fail "pip install failed"
  fi
fi

echo ""
echo -e "${GREEN}━━━ All dependencies reinstalled! ━━━${NC}"
