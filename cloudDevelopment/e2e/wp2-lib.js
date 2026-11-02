'use strict';
/**
 * Utilidades compartidas por los scripts de prueba de WP-2 (cloudDevelopment/e2e/wp2-*.js).
 *
 * Requisitos (sandbox cloud): la app corriendo en BASE_URL (./run_dev.sh) y el Postgres local
 * (psql en el PATH). Usa el Playwright que ya trae el sandbox: require('playwright') funciona
 * desde cualquier carpeta.
 *
 * Captcha: el backend lo omite cuando el User-Agent es "flutter-client" (mecanismo de la app
 * móvil); aquí se manda ese User-Agent SOLO en el POST de /seguridad/login.
 */
const { execFileSync } = require('child_process');

const BASE_URL = (process.env.BASE_URL || 'http://127.0.0.1:9090').replace(/\/$/, '');
const PASSWORD = process.env.TEST_PASSWORD || 'CloudDev!Test2026';
const DB = {
    host: process.env.DB_HOST || '127.0.0.1',
    port: process.env.DB_PORT || '5432',
    user: process.env.DB_USER || 'postgres',
    password: process.env.DB_PASSWORD || 'postgres',
    name: process.env.DB_NAME || 'salvia_local',
};

/** Ejecuta SQL en la BD local y devuelve la salida sin formato (-At). Lanza si psql falla. */
function sql(query) {
    return execFileSync(
        'psql',
        ['-h', DB.host, '-p', DB.port, '-U', DB.user, '-d', DB.name, '-X', '-q', '-At', '-v', 'ON_ERROR_STOP=1', '-c', query],
        { env: { ...process.env, PGPASSWORD: DB.password }, encoding: 'utf8' }
    ).trim();
}

// ── Registro de resultados ───────────────────────────────────────────────────────────────
const results = [];

/** Registra un resultado. `id` es el caso del MD de QA (C1, C2…) o un id corto. */
function check(id, ok, detail) {
    results.push({ id, ok: !!ok, detail });
    console.log(`${ok ? '✅' : '❌'} ${id} — ${detail}`);
}

/** Imprime el resumen y termina con código 0 (todo ok) o 1 (algo falló). */
function finish(title) {
    const failed = results.filter((r) => !r.ok);
    console.log(`\n${title}: ${results.length - failed.length}/${results.length} comprobaciones ok`);
    if (failed.length) {
        console.log('Fallaron: ' + failed.map((r) => r.id).join(', '));
    }
    process.exit(failed.length ? 1 : 0);
}

// ── Sesión de navegador ──────────────────────────────────────────────────────────────────
/** Abre un contexto nuevo e inicia sesión como cloud.test.<rol>. Devuelve { context, page }. */
async function loginAs(browser, role, contextOptions) {
    const context = await browser.newContext(contextOptions || {});
    const page = await context.newPage();
    await page.route('**/seguridad/login', (route) => {
        const req = route.request();
        if (req.method() !== 'POST') return route.continue();
        return route.continue({ headers: { ...req.headers(), 'user-agent': 'flutter-client' } });
    });
    await page.goto(`${BASE_URL}/seguridad/login`);
    await page.locator('input#login').fill(`cloud.test.${role}`);
    await page.locator('input#pass').fill(PASSWORD);
    await page.locator('input#captcha').fill('0000');
    const [resp] = await Promise.all([
        page.waitForResponse((r) => r.url().endsWith('/seguridad/login') && r.request().method() === 'POST'),
        page.locator('button#login').click(),
    ]);
    if (resp.status() !== 200) {
        throw new Error(`login ${role}: HTTP ${resp.status()}`);
    }
    await page.waitForURL((u) => !u.pathname.startsWith('/seguridad/login'), { timeout: 15000 });
    return { context, page };
}

/**
 * fetch con la sesión del navegador (las cookies no se comparten con context.request).
 * Devuelve { status, json, text }; json es null si el cuerpo no es JSON.
 */
async function api(page, method, path, body) {
    return page.evaluate(
        async (a) => {
            const opt = { method: a.method, headers: {}, credentials: 'same-origin' };
            if (a.body !== undefined) {
                opt.headers['Content-Type'] = 'application/json';
                opt.body = JSON.stringify(a.body);
            }
            const r = await fetch(a.path, opt);
            const text = await r.text();
            let json = null;
            try { json = JSON.parse(text); } catch (e) { /* cuerpo no JSON */ }
            return { status: r.status, json, text };
        },
        { method, path, body }
    );
}

/** "YYYY-MM-DD HH:MM:SS" de AHORA en hora de Bogotá, calculado en el cliente de prueba. */
function bogotaNow() {
    return new Intl.DateTimeFormat('sv-SE', {
        timeZone: 'America/Bogota',
        year: 'numeric', month: '2-digit', day: '2-digit',
        hour: '2-digit', minute: '2-digit', second: '2-digit', hour12: false,
    }).format(new Date());
}

/** Segundos entre dos "YYYY-MM-DD HH:MM:SS" (se leen como hora sin zona). */
function secondsBetween(a, b) {
    const t = (s) => Date.parse(s.replace(' ', 'T').slice(0, 19) + 'Z');
    return Math.abs(t(a) - t(b)) / 1000;
}

module.exports = { BASE_URL, PASSWORD, DB, sql, check, finish, loginAs, api, bogotaNow, secondsBetween };
