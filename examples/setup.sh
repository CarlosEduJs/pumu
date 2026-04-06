#!/usr/bin/env bash
# setup.sh — Bootstrap all example projects for pumu showcase.
# Usage: bash examples/setup.sh

set -euo pipefail

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
fail() { echo -e "${RED}  ✗ $1${NC}"; }

has() { command -v "$1" &>/dev/null; }

# ─────────────────────────────────────────────
# 1. pnpm
# ─────────────────────────────────────────────
step "node-pnpm (pnpm)"
if has pnpm; then
  mkdir -p node-pnpm && cd node-pnpm
  cat > package.json <<'EOF'
{
  "name": "example-pnpm",
  "version": "1.0.0",
  "private": true,
  "dependencies": {
    "fastify": "^5.2.0",
    "zod": "^3.24.0",
    "dayjs": "^1.11.0"
  }
}
EOF
  pnpm install --silent 2>/dev/null
  ok "node-pnpm created ($(du -sh node_modules 2>/dev/null | cut -f1))"
  cd ..
else
  warn "pnpm not found, skipping"
fi

# ─────────────────────────────────────────────
# 2. bun
# ─────────────────────────────────────────────
step "node-bun (bun)"
if has bun; then
  mkdir -p node-bun && cd node-bun
  cat > package.json <<'EOF'
{
  "name": "example-bun",
  "version": "1.0.0",
  "private": true,
  "dependencies": {
    "hono": "^4.7.0",
    "drizzle-orm": "^0.38.0"
  }
}
EOF
  bun install --silent 2>/dev/null
  ok "node-bun created ($(du -sh node_modules 2>/dev/null | cut -f1))"
  cd ..
else
  warn "bun not found, skipping"
fi

# ─────────────────────────────────────────────
# 3. go
# ─────────────────────────────────────────────
step "go-project (go)"
if has go; then
  mkdir -p go-project && cd go-project
  go mod init example-go 2>/dev/null || true
  cat > main.go <<'EOF'
package main

import (
	"fmt"

	"github.com/fatih/color"
)

func main() {
	color.Green("Hello from go-project!")
	fmt.Println("pumu example")
}
EOF
  go get github.com/fatih/color@latest 2>/dev/null
  go build -o go-project . 2>/dev/null
  ok "go-project created"
  cd ..
else
  warn "go not found, skipping"
fi

# ─────────────────────────────────────────────
# 4. python
# ─────────────────────────────────────────────
step "python-project (pip + venv)"
if has python3; then
  mkdir -p python-project && cd python-project
  cat > requirements.txt <<'EOF'
requests==2.32.3
flask==3.1.0
pydantic==2.10.0
EOF
  python3 -m venv .venv
  .venv/bin/pip install -r requirements.txt --quiet 2>/dev/null
  ok "python-project created ($(du -sh .venv 2>/dev/null | cut -f1))"
  cd ..
else
  warn "python3 not found, skipping"
fi

# ─────────────────────────────────────────────
# 5. next.js
# ─────────────────────────────────────────────
step "nextjs-project (next.js)"
if has npx; then
  mkdir -p nextjs-project && cd nextjs-project
  cat > package.json <<'EOF'
{
  "name": "example-nextjs",
  "version": "1.0.0",
  "private": true,
  "scripts": {
    "dev": "next dev",
    "build": "next build"
  },
  "dependencies": {
    "next": "^15.3.0",
    "react": "^19.0.0",
    "react-dom": "^19.0.0"
  }
}
EOF
  mkdir -p app
  cat > app/layout.tsx <<'EOF'
export const metadata = { title: "pumu example" };
export default function RootLayout({ children }: { children: React.ReactNode }) {
  return <html><body>{children}</body></html>;
}
EOF
  cat > app/page.tsx <<'EOF'
export default function Home() {
  return <h1>pumu nextjs example</h1>;
}
EOF
  npm install --silent 2>/dev/null
  npx next build 2>/dev/null || true
  ok "nextjs-project created (node_modules: $(du -sh node_modules 2>/dev/null | cut -f1), .next: $(du -sh .next 2>/dev/null | cut -f1))"
  cd ..
else
  warn "npx not found, skipping"
fi

# ─────────────────────────────────────────────
# 6. mix
# ─────────────────────────────────────────────
step "elixir-mix (mix)"
if has mix; then
  mkdir -p elixir-mix && cd elixir-mix
  if [ ! -f mix.exs ]; then
    warn "mix.exs missing, skipping"
  else
    mix deps.get 2>/dev/null || true
    ok "elixir-mix created ($(du -sh deps 2>/dev/null | cut -f1))"
  fi
  cd ..
else
  warn "mix not found, skipping"
fi

# ─────────────────────────────────────────────
echo ""
echo -e "${GREEN}━━━ All examples ready! ━━━${NC}"
echo ""
echo "Now try:"
echo "  pumu list -p examples/"
echo "  pumu sweep -p examples/"
echo "  pumu prune --dry-run -p examples/"
