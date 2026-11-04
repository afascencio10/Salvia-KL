'use strict';
/**
 * Utilidades compartidas por los scripts de prueba de WP-2 (scripts/e2e/wp2-*.js).
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

/** Ejecuta un archivo .sql en la BD local con psql -f. Devuelve 0 si salió bien o el código de salida si falló. */
function sqlFile(file) {
    try {
        execFileSync('psql', ['-h', DB.host, '-p', DB.port, '-U', DB.user, '-d', DB.name, '-X', '-q', '-v', 'ON_ERROR_STOP=1', '-f', file],
            { env: { ...process.env, PGPASSWORD: DB.password }, stdio: 'ignore' });
        return 0;
    } catch (e) {
        return e.status || 1;
    }
}

/**
 * Borra casos de prueba por icode y lo que cuelga de ellos, incluido el usuario 'us' (perfil, rol y
 * usuario) que genera cada borrador. Solo para el sandbox. Nunca toca a los cloud.test.*.
 */
function deleteCases(icodes) {
    const list = (icodes || []).filter(Boolean).map((i) => `'${String(i).replace(/'/g, "''")}'`).join(',');
    if (!list) return;
    // Los ids del usuario de cada caso se leen ANTES de borrar el caso: es lo que los une.
    const users = sql(`SELECT gu.general_user_id || ',' || COALESCE(gu.general_user_general_user_profile::text, '')
                         FROM security.general_user gu
                        WHERE gu.general_user_i_code IN (SELECT victim_case_general_user FROM salvia.victim_case WHERE victim_case_i_code IN (${list}))
                          AND gu.general_user_login NOT LIKE 'cloud.test.%'`)
        .split('\n').filter(Boolean).map((r) => r.split(','));
    const userIds = users.map((u) => u[0]);
    const profileIds = users.map((u) => u[1]).filter(Boolean);
    const caseIds = `SELECT victim_case_id FROM salvia.victim_case WHERE victim_case_i_code IN (${list})`;
    const form2Ids = `SELECT victim_case_form2_id FROM salvia.victim_case_form2 WHERE victim_case_form2_victim_case IN (${caseIds})`;
    sql(`DELETE FROM salvia.rel_victim_case_form2_enums_victim_case_form2 WHERE victim_case_form2_id IN (${form2Ids})`);
    sql(`DELETE FROM salvia.victim_case_form2 WHERE victim_case_form2_victim_case IN (${caseIds})`);
    sql(`DELETE FROM salvia.victim_case WHERE victim_case_i_code IN (${list})`);
    if (userIds.length) {
        sql(`DELETE FROM security.rel_role_general_user WHERE general_user_id IN (${userIds.join(',')})`);
        sql(`DELETE FROM security.general_user WHERE general_user_id IN (${userIds.join(',')})`);
    }
    if (profileIds.length) sql(`DELETE FROM security.general_user_profile WHERE general_user_profile_id IN (${profileIds.join(',')})`);
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

// ── Pantalla Registro V3 ─────────────────────────────────────────────────────────────────
/** Lee el catálogo de opciones de Registro V3 y devuelve utilidades para armar respuestas válidas. */
async function loadCatalog(page) {
    const cat = (await api(page, 'GET', '/api/v1/victim-cases-v3/enums')).json || {};
    const first = (c) => (cat[c] && cat[c][0] && cat[c][0].icode) || '';
    const byCode = (c, code) => ((cat[c] || []).find((o) => o.code === code) || {}).icode || '';
    return { cat, first, byCode };
}

/** Respuestas válidas por módulo (CI/DP/HV/AG: las obligatorias de la tríada) para el documento `doc`. */
function validAnswers({ first, byCode }, doc) {
    return {
        CI: { CI_1: byCode('victim_case_form2_pronoun', 'la'), CI_2: 'Ana Perez', CI_3: 'Anita', CI_4: '3001234567', CI_5: byCode('victim_case_form2_victim_doc_type', 'cc'), CI_6: doc, CI_7: 'false', CI_9: 'false' },
        DP: { DP_1: '1990-05-20', DP_2: 'false', DP_5: first('victim_case_form2_nationality'), DP_8: first('victim_case_form2_gender_identity'), DP_9: first('victim_case_form2_sexual_orientation'), DP_11: first('victim_case_form2_specially_protected_population'), DP_12: first('victim_case_form2_ethnic_affiliation'), DP_14: 'false', DP_17: first('victim_case_form2_occupation'), DP_27: 'Atlántico', DP_28: 'Barranquilla', DP_29: 'Calle 1 # 2-3' },
        HV: { HV_1: 'Descripción de prueba de los hechos', HV_2: '2026-10-30', HV_3: '19:30', HV_4: 'Atlántico', HV_5: 'Barranquilla', HV_6: first('victim_case_form2_facts_zone'), HV_7: 'Calle 45 # 12-34', HV_8: first('victim_case_form2_scenario_violence'), HV_9: first('victim_case_form2_type_violence_experienced'), HV_11: first('victim_case_form2_scope_of_violence'), HV_13: 'true', HV_14: first('victim_case_form2_recurrence_aggression') },
        AG: { AG_1: first('victim_case_form2_num_agressors'), AG_2: first('victim_case_form2_proximity_principal_aggressor'), AG_3: first('victim_case_form2_relationship_with_presumed_aggressor_01') },
    };
}

/** Abre /salvia/registro-v3 y pasa la autorización y el tipo de contacto: queda en el Cabezote. */
async function openRegistroV3(page) {
    await page.goto(`${BASE_URL}/salvia/registro-v3`);
    await page.getByRole('button', { name: 'Acepto y continuar' }).click();
    await page.getByRole('button', { name: 'Víctima - Caso Nuevo' }).click();
    await page.getByRole('button', { name: '+ Nuevo registro' }).waitFor();
}

/** Cabezote → «+ Nuevo registro». Devuelve el icode del borrador creado (queda abierto en CI). */
async function newDraftUI(page) {
    const [resp] = await Promise.all([
        page.waitForResponse((r) => r.url().endsWith('/api/v1/victim-cases-v3/draft') && r.request().method() === 'POST'),
        page.getByRole('button', { name: '+ Nuevo registro' }).click(),
    ]);
    const body = await resp.json();
    await page.getByRole('link', { name: 'Datos personales' }).waitFor();
    return body.icode;
}

/** Cabezote → busca por documento y abre el Caso (queda en CI). */
async function openCaseUI(page, doc) {
    await openRegistroV3(page);
    await page.locator('input[placeholder="Número de documento"]').fill(doc);
    await page.getByRole('button', { name: 'Buscar' }).click();
    await page.locator('tr', { hasText: `cc-${doc}` }).getByRole('button', { name: 'Abrir' }).click();
    await page.getByRole('link', { name: 'Contacto e identificación' }).waitFor();
}

/** Abre un módulo por el texto de su pestaña (p. ej. 'Datos personales'). */
async function openModule(page, tabName) {
    await page.getByRole('link', { name: tabName }).click();
}

/** Contenedor (div.mb-3) de una pregunta, por el texto exacto de su etiqueta. */
function question(page, label) {
    return page.locator(`xpath=//div[contains(@class,"mb-3")][./label[normalize-space(.)="${label}"]]`);
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

module.exports = {
    BASE_URL, PASSWORD, DB, sql, sqlFile, deleteCases, check, finish, loginAs, api, bogotaNow, secondsBetween,
    loadCatalog, validAnswers, openRegistroV3, newDraftUI, openCaseUI, openModule, question,
};
