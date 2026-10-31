#!/bin/bash
# Levanta el backend Go en modo HTTP local (sin TLS) para pruebas de QA/Playwright.
REPO_ROOT="$(cd "$(dirname "$0")" && pwd)"

# En el sandbox cloud de Claude Code (CLAUDE_CODE_REMOTE=true) no hay salida
# de red hacia la BD compartida — levanta Postgres local y exporta DB_*.
# No-op en local: sigue usando config/db_config.json sin tocarlo.
# Ver cloudDevelopment/README.md.
source "$REPO_ROOT/cloudDevelopment/setup-local-postgres.sh"

cd "$REPO_ROOT/src"
export PORT=9090
exec go run -a .
