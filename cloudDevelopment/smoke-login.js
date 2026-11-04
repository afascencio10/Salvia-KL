#!/usr/bin/env node
/**
 * Smoke test del sandbox cloud: con la plataforma corriendo (./run_dev.sh),
 * inicia sesión por la UI con cada usuario de prueba cloud.test.<rol> y
 * verifica a dónde aterriza. Sirve para comprobar en ~1 minuto que el
 * ambiente quedó sano (Postgres + seed + app + login) antes de desarrollar.
 *
 * Uso:
 *   node cloudDevelopment/smoke-login.js              # los 19 roles
 *   node cloudDevelopment/smoke-login.js ad sv ps     # solo esos roles
 *   SCREENSHOTS=/ruta/dir node cloudDevelopment/smoke-login.js   # PNG por rol
 *   BASE_URL=http://127.0.0.1:9090 TEST_PASSWORD=... node ...    # otros valores
 *
 * Usa el Playwright que ya viene instalado en el sandbox (require('playwright')
 * funciona desde cualquier carpeta). Sale con código 0 si todos los roles
 * inician sesión y aterrizan en una página sana; 1 si alguno falla.
 *
 * Captcha: el login tiene un captcha de imagen. El backend lo omite cuando el
 * User-Agent es "flutter-client" (mecanismo de la app móvil), así que el script
 * manda ese User-Agent SOLO en el POST de /seguridad/login; el resto de la
 * navegación va con el User-Agent normal del navegador.
 */
const path = require('path');
const { chromium } = require('playwright');

const BASE_URL = (process.env.BASE_URL || 'http://127.0.0.1:9090').replace(/\/$/, '');
const PASSWORD = process.env.TEST_PASSWORD || 'CloudDev!Test2026';
const SCREENSHOTS = process.env.SCREENSHOTS || '';
const ALL_ROLES = ['ad', 'op', 'sv', 'ro', 'do', 'no', 'et', 'en', 'an', 'us', 'ps', 'ts', 'te', 'mc', 'sm', 'em', 'sc', 'cr', 'am'];
const roles = process.argv.slice(2).length ? process.argv.slice(2) : ALL_ROLES;
const origin = new URL(BASE_URL).origin;

async function loginAs(browser, role) {
  const login = `cloud.test.${role}`;
  const context = await browser.newContext();
  const page = await context.newPage();
  const serverErrors = [];
  const clientErrors = [];
  const pageErrors = [];

  page.on('pageerror', (e) => pageErrors.push(String(e.message).slice(0, 100)));
  page.on('response', (r) => {
    const u = new URL(r.url());
    if (u.origin !== origin) return;
    const s = r.status();
    const entry = `${s} ${r.request().method()} ${u.pathname}`;
    if (s >= 500) serverErrors.push(entry);
    else if (s >= 400) clientErrors.push(entry);
  });

  // Captcha: User-Agent "flutter-client" solo en el POST de login.
  await page.route('**/seguridad/login', (route) => {
    const req = route.request();
    if (req.method() !== 'POST') return route.continue();
    return route.continue({ headers: { ...req.headers(), 'user-agent': 'flutter-client' } });
  });

  try {
    await page.goto(`${BASE_URL}/seguridad/login`);
    await page.locator('input#login').fill(login);
    await page.locator('input#pass').fill(PASSWORD);
    await page.locator('input#captcha').fill('0000');

    const [loginResp] = await Promise.all([
      page.waitForResponse((r) => r.url().endsWith('/seguridad/login') && r.request().method() === 'POST'),
      page.locator('button#login').click(),
    ]);
    if (loginResp.status() !== 200) {
      const body = (await loginResp.text()).slice(0, 160).replace(/\s+/g, ' ');
      return { role, login, ok: false, landing: '-', title: '-', notes: `login HTTP ${loginResp.status()}: ${body}` };
    }

    await page.waitForURL((u) => !u.pathname.startsWith('/seguridad/login'), { timeout: 15000 });
    await page.waitForLoadState('networkidle', { timeout: 15000 }).catch(() => {});
    const landing = new URL(page.url()).pathname;
    const title = (await page.title()).trim();
    const greeting = (await page.locator('body').innerText().catch(() => '')).match(/Hola,\s*([^\n]+)/);
    if (SCREENSHOTS) await page.screenshot({ path: path.join(SCREENSHOTS, `${role}.png`) });

    const notes = [];
    if (greeting) notes.push(`"Hola, ${greeting[1].trim()}"`);
    if (serverErrors.length) notes.push(`5xx: ${serverErrors.join('; ')}`);
    if (clientErrors.length) notes.push(`4xx: ${clientErrors.join('; ')}`);
    if (pageErrors.length) notes.push(`JS: ${pageErrors.join('; ')}`);
    return { role, login, ok: serverErrors.length === 0, landing, title, notes: notes.join(' | ') };
  } catch (e) {
    return { role, login, ok: false, landing: '-', title: '-', notes: `excepción: ${String(e.message).split('\n')[0]}` };
  } finally {
    await context.close();
  }
}

(async () => {
  const browser = await chromium.launch();
  const results = [];
  for (const role of roles) {
    const r = await loginAs(browser, role);
    results.push(r);
    console.log(`${r.ok ? 'OK  ' : 'FAIL'} ${r.login.padEnd(15)} → ${r.landing.padEnd(34)} ${r.title}${r.notes ? '\n       ' + r.notes : ''}`);
  }
  await browser.close();

  const failed = results.filter((r) => !r.ok);
  console.log(`\n${results.length - failed.length}/${results.length} roles inician sesión y aterrizan sin errores 5xx.`);
  if (failed.length) console.log(`Fallaron: ${failed.map((r) => r.role).join(', ')}`);
  process.exit(failed.length ? 1 : 0);
})();
