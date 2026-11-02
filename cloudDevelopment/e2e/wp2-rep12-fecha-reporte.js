#!/usr/bin/env node
/**
 * WP-2 · T1 · REP-12 — «Fecha Reporte» del Listado de Reportes y de Mis Reportes.
 *
 * Qué comprueba: GET /api/v1/victim-contacts/list devuelve `creationDateCO` ("YYYY-MM-DD") con el
 * día en que el reporte se creó en Bogotá, y la pantalla lo muestra como DD/MM/AAAA SIN cambiar de
 * día por la zona horaria del navegador (comentario de Felipe en la planeación).
 *
 * Orden (qa-process.md): 1) backend — endpoint + verificación directa en la BD; 2) frontend — solo
 * lo que el backend no cubre (el texto que ve la persona, con el navegador en varias zonas).
 *
 * Uso (con la app corriendo: ./run_dev.sh):
 *   node cloudDevelopment/e2e/wp2-rep12-fecha-reporte.js
 *
 * Crea reportes de prueba (i_code qa-wp2-rep12-*, nombre QAREP12*) directo en la BD y los borra al
 * terminar. Casos: DocsMD/Screens/Listado Reportes/QA/E01-carga-pantalla-qa-cases.md (LR-E01-C14…C20) y Mis Reportes/QA/mis-reportes-qa-cases.md (MR-E01-C10, C11).
 */
const { chromium, request } = require('playwright');
const { BASE_URL, sql, check, finish, loginAs, api, bogotaNow, secondsBetween } = require('./wp2-lib');

const PREFIX = 'qa-wp2-rep12-';
// Hora guardada (ya en Bogotá) → día esperado en el API y texto esperado en pantalla.
const FIXTURES = [
    { code: 'a', name: 'QAREP12A', ts: '2026-10-31 23:59:59', day: '2026-10-31', ui: '31/10/2026', why: 'último segundo del día' },
    { code: 'b', name: 'QAREP12B', ts: '2026-11-01 00:00:00', day: '2026-11-01', ui: '01/11/2026', why: 'primer segundo del día' },
    { code: 'c', name: 'QAREP12C', ts: '2026-11-01 03:30:00', day: '2026-11-01', ui: '01/11/2026', why: 'madrugada: leída como UTC y vista en Bogotá cae el día anterior' },
    { code: 'd', name: 'QAREP12D', ts: '2026-12-31 23:30:00', day: '2026-12-31', ui: '31/12/2026', why: 'fin de año' },
    { code: 'e', name: 'QAREP12E', ts: '2027-01-01 00:00:01', day: '2027-01-01', ui: '01/01/2027', why: 'inicio de año' },
];
const BROWSER_ZONES = ['America/Bogota', 'UTC', 'Asia/Tokyo', 'Pacific/Auckland', 'America/Los_Angeles'];

function cleanup() {
    // Reportes creados por el asistente público (nombre QAREP12W) y las filas que cuelgan de ellos.
    sql(`DELETE FROM salvia.victim_contact_form2_v2 WHERE victim_contact_form2_v2_victim_contact_form2 IN (
            SELECT f.victim_contact_form2_id FROM salvia.victim_contact_form2 f
            JOIN salvia.victim_contact c ON c.victim_contact_id = f.victim_contact_form2_victim_contact
            WHERE c.victim_contact_names = 'QAREP12W')`);
    sql(`DELETE FROM salvia.rel_victim_case_form2_enums_victim_contact_form2 WHERE victim_contact_form2_id IN (
            SELECT f.victim_contact_form2_id FROM salvia.victim_contact_form2 f
            JOIN salvia.victim_contact c ON c.victim_contact_id = f.victim_contact_form2_victim_contact
            WHERE c.victim_contact_names = 'QAREP12W')`);
    sql(`DELETE FROM salvia.victim_contact_form2 WHERE victim_contact_form2_victim_contact IN (
            SELECT victim_contact_id FROM salvia.victim_contact WHERE victim_contact_names = 'QAREP12W')`);
    sql(`DELETE FROM salvia.victim_contact WHERE victim_contact_names = 'QAREP12W' OR victim_contact_i_code LIKE '${PREFIX}%'`);
}

function seed(opAgentICode) {
    cleanup();
    for (const f of FIXTURES) {
        // c y e quedan asignados al operador de prueba para probar "Mis Reportes".
        const agent = f.code === 'c' || f.code === 'e' ? `'${opAgentICode}'` : 'NULL';
        sql(`INSERT INTO salvia.victim_contact
               (victim_contact_i_code, victim_contact_creation_date, victim_contact_update_date,
                victim_contact_general_user, victim_contact_names, victim_contact_last_names,
                victim_contact_status, victim_contact_agent_id)
             VALUES ('${PREFIX}${f.code}', '${f.ts}', '${f.ts}', 'qa', '${f.name}', 'Borde', 'v', ${agent})`);
    }
}

const listPath = (q) => `/api/v1/victim-contacts/list?page=1&page_size=50&${q}`;

(async () => {
    const browser = await chromium.launch();
    const opICode = 'user-icode-cloud.test.op';
    try {
        seed(opICode);

        // ═══ FASE 1 — BACKEND ═══════════════════════════════════════════════════════════
        const sv = await loginAs(browser, 'sv');

        // C1–C4 · el día del API es el día guardado, en los bordes de día y de año.
        const r1 = await api(sv.page, 'GET', listPath('search=QAREP12&sort=registration_date&order=desc'));
        const rows = (r1.json && r1.json.reports) || [];
        check('C0', r1.status === 200 && rows.length >= FIXTURES.length, `GET list → HTTP ${r1.status}, ${rows.length} filas con QAREP12`);
        for (const f of FIXTURES) {
            const row = rows.find((x) => x.i_code === PREFIX + f.code);
            const label = { a: 'C1', b: 'C2', c: 'C3', d: 'C4a', e: 'C4b' }[f.code];
            check(label, row && row.creationDateCO === f.day,
                `guardado ${f.ts} (${f.why}) → creationDateCO=${row && row.creationDateCO} (esperado ${f.day})`);
        }

        // C6 · el orden por fecha sigue siendo cronológico (desc y asc).
        const codes = (rs) => rs.filter((x) => x.i_code.startsWith(PREFIX)).map((x) => x.i_code.slice(PREFIX.length));
        const expectedAsc = [...FIXTURES].sort((p, q) => p.ts.localeCompare(q.ts)).map((f) => f.code);
        const rAsc = await api(sv.page, 'GET', listPath('search=QAREP12&sort=registration_date&order=asc'));
        check('C6', JSON.stringify(codes(rows)) === JSON.stringify([...expectedAsc].reverse())
            && JSON.stringify(codes(rAsc.json.reports)) === JSON.stringify(expectedAsc),
            `orden desc=${codes(rows).join('')} asc=${codes(rAsc.json.reports).join('')}`);

        // C8 · TODAS las filas traen creationDateCO con formato válido y creationDate (compatibilidad).
        const rAll = await api(sv.page, 'GET', listPath('sort=registration_date&order=desc'));
        const bad = (rAll.json.reports || []).filter((x) => !/^\d{4}-\d{2}-\d{2}$/.test(x.creationDateCO || '') || !x.creationDate);
        check('C8', (rAll.json.reports || []).length > 0 && bad.length === 0,
            `${(rAll.json.reports || []).length} filas, ${bad.length} sin creationDateCO válido o sin creationDate`);

        // C7 · «Mis Reportes» (persona_asignada) trae el campo.
        const rMine = await api(sv.page, 'GET', listPath(`filter_key=persona_asignada&filter_value=${opICode}`));
        const mine = (rMine.json && rMine.json.reports) || [];
        check('C7', mine.length === 2 && mine.every((x) => /^\d{4}-\d{2}-\d{2}$/.test(x.creationDateCO)),
            `persona_asignada=${opICode} → ${mine.length} reportes: ${mine.map((x) => x.creationDateCO).join(', ')}`);

        // C5 · reporte creado AHORA por el asistente público: la BD guarda hora de Bogotá y el API da su día.
        const apiCtx = await request.newContext({ baseURL: BASE_URL, extraHTTPHeaders: { 'User-Agent': 'flutter-client' } });
        const cat = await (await apiCtx.get('/salvia/public-v2/nuevo', { headers: { Accept: 'application/json' } })).json();
        const yes = cat.yesNo.find((e) => e.code === 'y');
        const no = cat.yesNo.find((e) => e.code === 'n');
        const post = await apiCtx.post('/salvia/public-v2', {
            data: {
                victimContact: {
                    names: 'QAREP12W', lastNames: 'Asistente', captchaID: cat.captchaID, captchaSolution: '0000',
                    form2: {
                        reportType: cat.reportType.find((e) => e.code === 'v'), adjustmentsGBV: [cat.adjustmentsGBV[0]],
                        reporterNames: '', reporterPhone: 0, hasCareRole: null, victimAwareOfReport: null,
                        victimColPhone: 3001234567, willReceiveCall: no, bestContactTime: '00:00:00',
                        factsDescription: 'Reporte de prueba automatizada WP-2 REP-12 (fecha del reporte).',
                    },
                    form2V2: {
                        authorizationAnswer: yes, anonimato: null, parentesco: null, reporterLastNames: '',
                        reporterEmail: '', reporterWillReceiveCall: null, reporterBestContactTime: null,
                    },
                },
            },
        });
        const created = await post.json().catch(() => ({}));
        const beforeNow = bogotaNow();
        const stored = created.icode
            ? sql(`SELECT to_char(victim_contact_creation_date, 'YYYY-MM-DD HH24:MI:SS') FROM salvia.victim_contact WHERE victim_contact_i_code = '${created.icode}'`)
            : '';
        check('C5a', post.status() === 200 && !!created.icode, `POST /salvia/public-v2 → HTTP ${post.status()} icode=${created.icode}`);
        check('C5b', stored && secondsBetween(stored, beforeNow) < 180,
            `BD guardó ${stored}; hora de Bogotá ahora ${beforeNow} (si fuera UTC diferirían 5 h)`);
        const rNew = await api(sv.page, 'GET', listPath('search=QAREP12W'));
        const newRow = ((rNew.json && rNew.json.reports) || [])[0];
        check('C5c', newRow && newRow.creationDateCO === beforeNow.slice(0, 10),
            `API creationDateCO=${newRow && newRow.creationDateCO}; día de Bogotá=${beforeNow.slice(0, 10)}`);
        await apiCtx.dispose();
        await sv.context.close();

        // ═══ FASE 2 — FRONTEND (solo lo que el backend no cubre: el texto en pantalla) ═════
        const dateCells = async (page, url) => {
            await page.goto(`${BASE_URL}${url}`);
            const search = page.locator('input[placeholder^="Buscar por nombres"]');
            await search.waitFor({ timeout: 15000 });
            await Promise.all([
                page.waitForResponse((r) => r.url().includes('/api/v1/victim-contacts/list') && r.url().includes('search=QAREP12')),
                search.fill('QAREP12'),
            ]);
            await page.waitForFunction(() => document.querySelectorAll('table tbody tr').length > 0);
            return page.evaluate(() => {
                const headers = [...document.querySelectorAll('table thead th')].map((h) => h.innerText.trim());
                const idx = headers.findIndex((h) => /Fecha Reporte/i.test(h));
                return [...document.querySelectorAll('table tbody tr')].map((tr) => {
                    const tds = [...tr.querySelectorAll('td')];
                    const m = (tr.innerText.match(/ID:\s*(\S+)/) || [])[1];
                    return { icode: m, date: idx >= 0 && tds[idx] ? tds[idx].innerText.trim() : null };
                });
            });
        };

        const seen = {};
        for (const tz of BROWSER_ZONES) {
            const op = await loginAs(browser, 'sv', { timezoneId: tz });
            const cells = await dateCells(op.page, '/salvia/reportes');
            seen[tz] = Object.fromEntries(cells.map((c) => [c.icode, c.date]));
            await op.context.close();
        }
        const expectedUi = Object.fromEntries(FIXTURES.map((f) => [PREFIX + f.code, f.ui]));
        for (const tz of BROWSER_ZONES) {
            const wrong = FIXTURES.filter((f) => seen[tz][PREFIX + f.code] !== f.ui)
                .map((f) => `${f.code}: ${seen[tz][PREFIX + f.code]}≠${f.ui}`);
            check('C9:' + tz, wrong.length === 0, wrong.length ? wrong.join('; ') : 'las 5 filas muestran el día guardado (Listado de Reportes)');
        }
        check('C9:igual', BROWSER_ZONES.every((tz) => JSON.stringify(seen[tz]) === JSON.stringify(seen[BROWSER_ZONES[0]])),
            'mismo texto en las 5 zonas horarias del navegador');

        // C10 · Mis Reportes (solo el rol op): mismas fechas.
        const mineUi = await loginAs(browser, 'op', { timezoneId: 'Pacific/Auckland' });
        const mineCells = await dateCells(mineUi.page, '/salvia/mis-reportes');
        const mineMap = Object.fromEntries(mineCells.map((c) => [c.icode, c.date]));
        check('C10', mineMap[PREFIX + 'c'] === expectedUi[PREFIX + 'c'] && mineMap[PREFIX + 'e'] === expectedUi[PREFIX + 'e']
            && Object.keys(mineMap).length === 2,
            `Mis Reportes (navegador en Pacific/Auckland): ${JSON.stringify(mineMap)}`);
        await mineUi.context.close();
    } catch (e) {
        check('ERROR', false, e && e.stack ? e.stack.split('\n').slice(0, 3).join(' | ') : String(e));
    } finally {
        try { cleanup(); } catch (e) { console.log('(limpieza falló: ' + e.message + ')'); }
        await browser.close();
    }
    finish('REP-12 · fecha del reporte');
})();
