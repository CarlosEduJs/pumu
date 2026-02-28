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
# 1. npm
# ─────────────────────────────────────────────
step "node-npm (npm)"
if has npm; then
  mkdir -p node-npm && cd node-npm
  cat > package.json <<'EOF'
{
  "name": "example-npm",
  "version": "1.0.0",
  "private": true,
  "dependencies": {
    "express": "^4.21.0",
    "lodash": "^4.17.21",
    "chalk": "^5.4.0"
  }
}
EOF
  npm install --silent 2>/dev/null
  ok "node-npm created ($(du -sh node_modules 2>/dev/null | cut -f1))"
  cd ..
else
  warn "npm not found, skipping"
fi

# ─────────────────────────────────────────────
# 2. pnpm
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
# 3. bun
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
# 4. deno
# ─────────────────────────────────────────────
step "deno-project (deno)"
if has deno; then
  mkdir -p deno-project && cd deno-project
  cat > deno.json <<'EOF'
{
  "imports": {
    "@std/path": "jsr:@std/path@^1.0.0",
    "@std/fs": "jsr:@std/fs@^1.0.0"
  }
}
EOF
  cat > main.ts <<'EOF'
import { join } from "@std/path";
import { exists } from "@std/fs";

const p = join(".", "hello");
console.log("exists:", await exists(p));
EOF
  deno install --allow-read 2>/dev/null || true
  ok "deno-project created ($(du -sh node_modules 2>/dev/null | cut -f1))"
  cd ..
else
  warn "deno not found, skipping"
fi

# ─────────────────────────────────────────────
# 5. rust
# ─────────────────────────────────────────────
step "rust-project (cargo)"
if has cargo; then
  mkdir -p rust-project && cd rust-project
  cat > Cargo.toml <<'EOF'
[package]
name = "example-rust"
version = "0.1.0"
edition = "2021"

[dependencies]
serde = { version = "1", features = ["derive"] }
serde_json = "1"
tokio = { version = "1", features = ["full"] }
EOF
  mkdir -p src
  cat > src/main.rs <<'EOF'
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug)]
struct Example {
    name: String,
    value: i32,
}

fn main() {
    let e = Example { name: "pumu".into(), value: 42 };
    let json = serde_json::to_string_pretty(&e).unwrap();
    println!("{json}");
}
EOF
  cargo build 2>/dev/null
  ok "rust-project created ($(du -sh target 2>/dev/null | cut -f1))"
  cd ..
else
  warn "cargo not found, skipping"
fi

# ─────────────────────────────────────────────
# 6. go
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
# 7. python
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
# 8. next.js
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
echo ""
echo -e "${GREEN}━━━ All examples ready! ━━━${NC}"
echo ""
echo "Now try:"
echo "  pumu list -p examples/"
echo "  pumu sweep -p examples/"
echo "  pumu prune --dry-run -p examples/"
