#!/bin/bash
# Playwright CLI (@playwright/cli) listo para usar en el sandbox cloud.
#
# El sandbox ya trae Playwright (librería, usable con require('playwright')) y un
# Chromium en /opt/pw-browsers. El CLI (@playwright/cli) depende de OTRO Playwright
# (alfa) que no encontraría ese Chromium, así que se instala aparte, en
# /opt/playwright-cli y sin descargar navegadores, y se apunta al Chromium que ya
# está. Este script hace todo eso y reenvía los comandos al CLI.
#
# Uso (con la plataforma corriendo: ./run_dev.sh):
#   cloudDevelopment/playwright-cli.sh login ad          # abre la app e inicia sesión como cloud.test.ad
#   cloudDevelopment/playwright-cli.sh snapshot          # estructura de la página con referencias (e12, e15...)
#   cloudDevelopment/playwright-cli.sh click e29
#   cloudDevelopment/playwright-cli.sh goto http://127.0.0.1:9090/salvia/mis-tareas
#   cloudDevelopment/playwright-cli.sh screenshot --filename=/tmp/pantalla.png
#   cloudDevelopment/playwright-cli.sh close             # cierra el navegador (y la sesión)
#   cloudDevelopment/playwright-cli.sh --help            # todos los comandos del CLI
#
# `login <rol>` entra con cloud.test.<rol> (clave CloudDev!Test2026) y deja el
# navegador abierto en la pantalla de aterrizaje del rol, listo para seguir con
# snapshot / click / fill / goto. Para cambiar de rol: `login <otro-rol>` (cierra
# la sesión del navegador anterior).
#
# Captcha: el login tiene captcha de imagen; el backend lo omite si el User-Agent
# es "flutter-client" (mecanismo de la app móvil), así que `login` manda ese
# User-Agent SOLO en el POST de /seguridad/login.
#
# Variables: BASE_URL (http://127.0.0.1:9090), TEST_PASSWORD, PLAYWRIGHT_CLI_DIR
# (/opt/playwright-cli), PLAYWRIGHT_CLI_OUT (/tmp/playwright-cli: capturas y
# snapshots; fuera del repo a propósito).

PFX="${PLAYWRIGHT_CLI_DIR:-/opt/playwright-cli}"
BIN="$PFX/node_modules/.bin/playwright-cli"
BASE="${BASE_URL:-http://127.0.0.1:9090}"
PASS="${TEST_PASSWORD:-CloudDev!Test2026}"
OUT="${PLAYWRIGHT_CLI_OUT:-/tmp/playwright-cli}"

if [ ! -x "$BIN" ]; then
    echo "[playwright-cli] Instalando @playwright/cli en $PFX (una sola vez, sin descargar navegadores)..." >&2
    mkdir -p "$PFX"
    PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD=1 npm install --prefix "$PFX" --silent --no-audit --no-fund @playwright/cli@0.1.23 >&2 || exit 1
fi

export PLAYWRIGHT_MCP_EXECUTABLE_PATH="${PLAYWRIGHT_MCP_EXECUTABLE_PATH:-/opt/pw-browsers/chromium}"
export PLAYWRIGHT_MCP_SANDBOX=false      # corremos como root: Chromium necesita --no-sandbox
export PLAYWRIGHT_MCP_HEADLESS=true
export NO_PROXY="127.0.0.1,localhost${NO_PROXY:+,$NO_PROXY}"
export no_proxy="$NO_PROXY"

# El CLI deja sus artefactos (.playwright-cli/) en el directorio actual: fuera del repo.
mkdir -p "$OUT"
cd "$OUT" || exit 1

if [ "${1:-}" = "login" ]; then
    role="${2:-}"
    if [ -z "$role" ]; then
        echo "Uso: cloudDevelopment/playwright-cli.sh login <rol>   (ad, op, sv, ro, do, no, et, en, an, us, ps, ts, te, mc, sm, em, sc, cr, am)" >&2
        exit 2
    fi
    "$BIN" close >/dev/null 2>&1 || true    # sesión nueva: sin cookies del rol anterior
    "$BIN" open >/dev/null || exit 1
    "$BIN" run-code "async page => {
      await page.route('**/seguridad/login', route => {
        const req = route.request();
        return req.method() === 'POST'
          ? route.continue({ headers: { ...req.headers(), 'user-agent': 'flutter-client' } })
          : route.continue();
      });
      await page.goto('$BASE/seguridad/login');
      await page.locator('input#login').fill('cloud.test.$role');
      await page.locator('input#pass').fill('$PASS');
      await page.locator('input#captcha').fill('0000');
      await Promise.all([
        page.waitForURL(u => !u.pathname.startsWith('/seguridad/login'), { timeout: 15000 }),
        page.locator('button#login').click(),
      ]);
      await page.waitForLoadState('networkidle', { timeout: 15000 }).catch(() => {});
      return page.url() + '  |  ' + (await page.title());
    }"
    exit $?
fi

exec "$BIN" "$@"
