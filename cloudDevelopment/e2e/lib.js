'use strict';
/**
 * Utilidades compartidas de los recorridos E2E del sandbox cloud (Playwright, sin
 * test runner). Usan el Playwright que ya trae el sandbox (`require('playwright')`
 * funciona desde cualquier carpeta) y la BD local (`psql`, auth trust) para verificar
 * lo que cada pantalla escribe.
 *
 * Variables: BASE_URL (http://127.0.0.1:9090), TEST_PASSWORD, SCREENSHOTS (carpeta para
 * capturas), RUN_ID (sufijo de los datos creados; por defecto, único por corrida),
 * DB_HOST/DB_PORT/DB_USER/DB_NAME (los de setup-local-postgres.sh), HEADED=1 (ver el
 * navegador; en el sandbox no hay pantalla, así que solo sirve en local), E2E_VERBOSE=1
 * (traza completa y "call log" de Playwright cuando un escenario falla).
 */
const { chromium } = require('playwright');
const { execFileSync } = require('child_process');
const path = require('path');

const BASE_URL = (process.env.BASE_URL || 'http://127.0.0.1:9090').replace(/\/$/, '');
const PASSWORD = process.env.TEST_PASSWORD || 'CloudDev!Test2026';
const SHOTS = process.env.SCREENSHOTS || '';
const RUN_ID = (process.env.RUN_ID || Date.now().toString(36)).toUpperCase();

// ── BD local ───────────────────────────────────────────────────────────────
function psql(query) {
  return execFileSync(
    process.env.PSQL_BIN || 'psql',
    ['-h', process.env.DB_HOST || '127.0.0.1', '-p', process.env.DB_PORT || '5432', '-U', process.env.DB_USER || 'postgres',
      '-d', process.env.DB_NAME || 'salvia_local', '-X', '-q', '-tA', '-F', '|', '-c', query],
    { encoding: 'utf8', env: { ...process.env, PGOPTIONS: '-c client_min_messages=warning' } },
  ).trim();
}
/** Filas como arrays de strings (columnas separadas por '|'). */
const rows = (query) => psql(query).split('\n').filter(Boolean).map((r) => r.split('|'));
/** Primera fila (o [] si no hay). */
const one = (query) => rows(query)[0] || [];
/** Escapa una cadena para usarla dentro de un literal SQL entre comillas simples. */
const lit = (s) => String(s).replace(/'/g, "''");

// ── Navegador y sesión ─────────────────────────────────────────────────────
async function launch() {
  return chromium.launch({ headless: !process.env.HEADED });
}

/**
 * El backend omite el captcha de imagen cuando el User-Agent es "flutter-client" (mecanismo
 * de la app móvil). Se manda ese User-Agent SOLO en los POST a las URLs indicadas; el resto
 * de la navegación va con el del navegador.
 */
async function bypassCaptcha(page, ...globs) {
  for (const glob of globs) {
    await page.route(glob, (route) => {
      const req = route.request();
      if (req.method() !== 'POST') return route.continue();
      return route.continue({ headers: { ...req.headers(), 'user-agent': 'flutter-client' } });
    });
  }
}

/** Inicia sesión por la UI como cloud.test.<rol>. Devuelve { context, page }. */
async function loginAs(browser, role) {
  const context = await browser.newContext({ viewport: { width: 1366, height: 900 } });
  const page = await context.newPage();
  await bypassCaptcha(page, '**/seguridad/login');
  await page.goto(`${BASE_URL}/seguridad/login`);
  await page.locator('input#login').fill(`cloud.test.${role}`);
  await page.locator('input#pass').fill(PASSWORD);
  await page.locator('input#captcha').fill('0000');
  await Promise.all([
    page.waitForURL((u) => !u.pathname.startsWith('/seguridad/login'), { timeout: 20000 }),
    page.locator('button#login').click(),
  ]);
  await page.waitForLoadState('networkidle', { timeout: 20000 }).catch(() => {});
  return { context, page };
}

/** fetch con la sesión del navegador. Devuelve { status, json, text }. */
function api(page, method, url, body) {
  return page.evaluate(async ([m, u, b]) => {
    const r = await fetch(u, { method: m, headers: { 'Content-Type': 'application/json' }, body: b ? JSON.stringify(b) : undefined });
    const text = await r.text();
    let json = null;
    try { json = JSON.parse(text); } catch (e) { /* no era JSON */ }
    return { status: r.status, json, text };
  }, [method, url, body]);
}

// ── Reporte de resultados ──────────────────────────────────────────────────
const results = [];

async function shot(page, name) {
  if (!SHOTS || !page) return;
  await page.screenshot({ path: path.join(SHOTS, `${name}.png`), fullPage: true }).catch(() => {});
}

/** Marca como verificada una condición; si no se cumple, aborta el escenario con ese mensaje. */
function check(cond, msg) {
  if (!cond) throw new Error(msg);
  console.log(`    · ${msg}`);
}

/** Ejecuta un escenario: un fallo lo marca como tal y deja seguir con el siguiente. */
async function scenario(name, fn) {
  const t0 = Date.now();
  console.log(`\n▶ ${name}`);
  try {
    await fn();
    results.push({ name, ok: true });
    console.log(`  ✔ OK (${((Date.now() - t0) / 1000).toFixed(1)} s)`);
  } catch (e) {
    const msg = String(e.message || e).split('\n')[0];
    results.push({ name, ok: false, error: msg });
    console.log(`  ✘ FALLÓ: ${msg}`);
    if (process.env.E2E_VERBOSE) console.log(String(e.stack || e.message).split('\n').slice(0, 14).map((l) => `      ${l}`).join('\n'));
  }
}

/** Imprime el resumen y termina el proceso con 0 (todo OK) o 1. */
function finish(titulo) {
  const fallidos = results.filter((r) => !r.ok);
  console.log(`\n══ ${titulo}: ${results.length - fallidos.length}/${results.length} escenarios OK ══`);
  fallidos.forEach((r) => console.log(`   ✘ ${r.name} — ${r.error}`));
  process.exit(fallidos.length ? 1 : 0);
}

module.exports = { BASE_URL, PASSWORD, RUN_ID, SHOTS, psql, rows, one, lit, launch, bypassCaptcha, loginAs, api, shot, check, scenario, finish };
