// Helpers de Playwright para las pruebas de pantalla de WP-4 (Registro V3).
// Uso: const { abrir } = require('./ui-lib'); const { page, ctx, browser, icode, api, sql } = await abrir({...});
const { chromium } = require('playwright');
const { execSync } = require('child_process');

const BASE = process.env.BASE_URL || 'http://127.0.0.1:9090';
const sql = (q) => execSync(`PGPASSWORD=postgres psql -h 127.0.0.1 -U postgres -d salvia_local -At -F ' | ' -c "${q.replace(/"/g, '\\"')}"`).toString().trim();

const TRIADA = 'CI_1,CI_2,CI_3,CI_4,CI_5,CI_6,CI_7,CI_9,DP_1,DP_2,DP_5,DP_8,DP_9,DP_11,DP_12,DP_14,DP_17,DP_27,DP_28,DP_29,HV_1,HV_2,HV_4,HV_6,HV_8,HV_9,HV_11,HV_13,HV_14,AG_1,AG_2,AG_3';

async function login(browser, role = 'op') {
  const ctx = await browser.newContext({ viewport: { width: 1280, height: 1600 } });
  const page = await ctx.newPage();
  await page.route('**/seguridad/login', (route) => {
    const req = route.request();
    if (req.method() !== 'POST') return route.continue();
    return route.continue({ headers: { ...req.headers(), 'user-agent': 'flutter-client' } });
  });
  await page.goto(`${BASE}/seguridad/login`);
  await page.locator('input#login').fill(`cloud.test.${role}`);
  await page.locator('input#pass').fill('CloudDev!Test2026');
  await page.locator('input#captcha').fill('0000');
  await Promise.all([
    page.waitForResponse((r) => r.url().endsWith('/seguridad/login') && r.request().method() === 'POST'),
    page.locator('button#login').click(),
  ]);
  await page.waitForURL((u) => !u.pathname.startsWith('/seguridad/login'), { timeout: 15000 });
  return { ctx, page };
}

// Crea un Caso en Borrador (API), opcionalmente con la tríada marcada como respondida por SQL
// (GetCase.triadaCompleta = true sin pasar a Iniciado: abre GS en modo borrador) y lo abre por la UI.
async function abrir({ role = 'op', triada = true, hook } = {}) {
  const browser = await chromium.launch();
  const { ctx, page } = await login(browser, role);
  const errores = [];
  page.on('pageerror', (e) => errores.push('pageerror: ' + e.message));
  page.on('console', (m) => { if (m.type() === 'error') errores.push('console: ' + m.text()); });
  // La API se llama desde la página (como la app): el contexto de Node no envía la cookie Secure por http.
  const api = (method, path, body) => page.evaluate(async ([m, p, b]) => {
    const r = await fetch('/api/v1/victim-cases-v3/' + p, {
      method: m, headers: { 'Content-Type': 'application/json' }, body: b ? JSON.stringify(b) : undefined,
    });
    let json = null; try { json = await r.json(); } catch (e) { /* sin cuerpo */ }
    return { status: r.status, json };
  }, [method, path, body || null]);
  const doc = String(Math.floor(900000000 + Math.random() * 99999999));
  const { json } = await api('POST', 'draft', { tipo_documento: 'cc', documento: doc, telefono: '3001112233' });
  const icode = json.icode;
  if (triada) sql(`update salvia.victim_case set victim_case_v3_answered_keys='${TRIADA}' where victim_case_i_code='${icode}'`);
  if (hook) await hook({ api, icode, sql });
  await page.goto(`${BASE}/salvia/registro-v3`);
  await page.getByRole('button', { name: 'Acepto y continuar' }).click();
  await page.getByRole('button', { name: 'Víctima - Caso Nuevo' }).click();
  await page.getByPlaceholder('Número de documento').fill(doc);
  await page.getByRole('button', { name: 'Buscar' }).click();
  await page.locator('table button', { hasText: 'Abrir' }).first().click();
  try {
    await page.waitForSelector('ul.nav-tabs', { timeout: 15000 });
  } catch (e) {
    await page.screenshot({ path: '/tmp/wp4-abrir-fallo.png' });
    console.log('No abrió el formulario. Errores de la página:', errores, '\nTexto visible:', (await page.locator('main').innerText()).slice(0, 600));
    throw e;
  }
  return { browser, ctx, page, icode, api, sql, errores, doc };
}

// Elige una entidad por la UI (cascada Bogotá D.C. + sector + «Buscar entidad»); `nueva` es el locator de la tarjeta.
async function elegirEntidad(page, nueva, { sector = 'justicia', texto = 'usaqu' } = {}) {
  const sel = (label) => nueva.locator(`label.form-label:has-text("${label}") + select`);
  await sel('Departamento').selectOption({ label: 'BOGOTÁ, D.C.' });
  await page.waitForTimeout(400);
  await sel('Ciudad').selectOption({ index: 1 });
  await page.waitForTimeout(400);
  await sel('Municipio').selectOption({ index: 1 });
  await sel('Sector').selectOption(sector);
  const buscar = nueva.locator('input[placeholder="Escribe el nombre de la entidad"]');
  await buscar.waitFor();
  await buscar.focus();
  await buscar.pressSequentially(texto, { delay: 20 });
  await page.waitForTimeout(450);
  await nueva.locator('.list-group button').first().dispatchEvent('mousedown');
  await page.waitForTimeout(600);
}

module.exports = { abrir, sql, BASE, elegirEntidad };
