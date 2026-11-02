#!/usr/bin/env node
/**
 * WP-2 · T6 · REG-12 — El Tamizaje de riesgo se presenta como tablas «Pregunta | Sí | No» y las preguntas
 * generales (TZ_1..TZ_8) se GUARDAN.
 *
 * Antes las preguntas generales eran bloques sueltos de radios y SaveTamizajeBranch solo escribía las 24
 * de la rama (pareja / no pareja): TZ_1..TZ_8 nunca llegaban a la BD, así que al reabrir el Caso todas las
 * generales aparecían como «No» aunque se hubieran respondido «Sí». Además el guardado no era atómico
 * (marcaba las claves como respondidas antes del UPDATE) y un caso inexistente respondía 200.
 *
 * Orden (qa-process.md): 1) backend — POST /tamizaje, GET del Caso y verificación directa en la BD;
 * 2) frontend — solo lo que el backend no cubre: la estructura de las tablas y el ciclo guardar/reabrir.
 *
 * Uso (con la app corriendo):  node cloudDevelopment/e2e/wp2-reg12-tamizaje-tablas.js
 * Casos: DocsMD/Screens/Registro V3/QA/E11-actualizar-respuesta-y-tamizaje-qa-cases.md (E11-C13…C20; comprobaciones C1…C5, U1…U3) y E21-guardado-o-carga-fallida-qa-cases.md (E21-C7; U4).
 */
const { chromium } = require('playwright');
const { sql, deleteCases, check, finish, loginAs, api, loadCatalog, validAnswers, openCaseUI, openModule } = require('./wp2-lib');

const created = [];
const rnd = () => String(900000000 + Math.floor(Math.random() * 99999999));
const GENERALES = ['TZ_1', 'TZ_2', 'TZ_3', 'TZ_4', 'TZ_5', 'TZ_6', 'TZ_7', 'TZ_8'];
const G_COLS = ['g1', 'g2', 'g3', 'g4', 'g5', 'g6', 'g7', 'vpp'].map((c) => 'victim_case_form2_tzv3_' + c);
const branch = (pref, trues) => Object.fromEntries(Array.from({ length: 24 }, (_, i) => [pref + (i + 1), trues.includes(i + 1) ? 'true' : 'false']));
const tf = (v) => (v ? 'true' : 'false');

(async () => {
    const browser = await chromium.launch();
    try {
        const op = await loginAs(browser, 'op');
        const page = op.page;
        const catalog = await loadCatalog(page);

        // ── Caso con la tríada completa (para que el Tamizaje esté habilitado); AG_3 = «Amigo(a)» → rama «No pareja».
        const doc = rnd();
        const d = await api(page, 'POST', '/api/v1/victim-cases-v3/draft', { tipo_documento: 'cc', documento: doc, telefono: '3009998888' });
        const icode = d.json.icode;
        created.push(icode);
        const v = validAnswers(catalog, doc);
        v.AG.AG_3 = catalog.byCode('victim_case_form2_relationship_with_presumed_aggressor_01', 'am');
        for (const m of ['CI', 'DP', 'HV', 'AG']) await api(page, 'POST', `/api/v1/victim-cases-v3/${icode}/seccion/${m}`, { answers: v[m] });

        const colCodes = (cols) => sql(`SELECT concat_ws(',', ${cols.map((c) => `(SELECT e.victim_case_form2_enums_code FROM salvia.victim_case_form2_enums e WHERE e.victim_case_form2_enums_id = f.${c})`).join(', ')})
            FROM salvia.victim_case_form2 f JOIN salvia.victim_case c ON c.victim_case_id = f.victim_case_form2_victim_case WHERE c.victim_case_i_code = '${icode}'`);
        const generalesDb = () => colCodes(G_COLS);
        const tamizaje = (answers, wasPartner, id) => api(page, 'POST', `/api/v1/victim-cases-v3/${id || icode}/tamizaje`, { answers, wasPartner, esOsigd: false, subtipoFeminicidio: false, subtipoHomicidioAbyecto: false });
        const getCase = () => api(page, 'GET', `/api/v1/victim-cases-v3/${icode}`);

        // ═══ FASE 1 — BACKEND ═══════════════════════════════════════════════════════════
        // C1 · Las generales respondidas se guardan en la BD (antes quedaban en el marcador «No»).
        const gen1 = { TZ_1: 'true', TZ_2: 'false', TZ_3: 'true', TZ_4: 'false', TZ_5: 'true', TZ_6: 'false', TZ_7: 'true', TZ_8: 'true' };
        const r1 = await tamizaje({ ...gen1, ...branch('TZNP_', [2, 24]) }, false);
        check('C1', r1.status === 200 && r1.json.score === 7 && r1.json.riskLevel === 'alto' && generalesDb() === 'y,n,y,n,y,n,y,y'
            && colCodes(['victim_case_form2_tzv3_np1', 'victim_case_form2_tzv3_np2', 'victim_case_form2_tzv3_np24']) === 'n,y,y',
            `HTTP ${r1.status} ${JSON.stringify(r1.json)}; generales en BD=${generalesDb()} (esperado y,n,y,n,y,n,y,y); rama np1,np2,np24=${colCodes(['victim_case_form2_tzv3_np1', 'victim_case_form2_tzv3_np2', 'victim_case_form2_tzv3_np24'])}`);

        // C2 · Al reabrir el Caso, las generales vuelven como se respondieron y el resultado está.
        const g2 = await getCase();
        const a2 = g2.json.answers;
        check('C2', GENERALES.every((k) => a2[k] === gen1[k]) && GENERALES.every((k) => g2.json.answeredKeys.includes(k))
            && a2.TZNP_2 === 'true' && g2.json.tzResultado && g2.json.tzResultado.score === 7 && g2.json.tzResultado.riskLevel === 'alto',
            `GET: ${GENERALES.map((k) => k + '=' + a2[k]).join(' ')}; tzResultado=${JSON.stringify(g2.json.tzResultado)}`);

        // C3 · Volver a guardar con otras respuestas actualiza la BD.
        const gen3 = { TZ_1: 'false', TZ_2: 'true', TZ_3: 'true', TZ_4: 'false', TZ_5: 'true', TZ_6: 'false', TZ_7: 'true', TZ_8: 'false' };
        const r3 = await tamizaje({ ...gen3, ...branch('TZNP_', []) }, false);
        const a3 = (await getCase()).json.answers;
        check('C3', r3.status === 200 && generalesDb() === 'n,y,y,n,y,n,y,n' && GENERALES.every((k) => a3[k] === gen3[k]),
            `HTTP ${r3.status} score=${r3.json.score}; BD=${generalesDb()} (esperado n,y,y,n,y,n,y,n)`);

        // C4 · Rama pareja: solo se escriben las columnas de la rama activa; las de la otra quedan como estaban.
        const r4 = await tamizaje({ ...Object.fromEntries(GENERALES.map((k) => [k, 'false'])), ...branch('TZP_', [3, 24]) }, true);
        check('C4', r4.status === 200 && r4.json.score === 2 && r4.json.riskLevel === 'bajo'
            && colCodes(['victim_case_form2_tzv3_p3', 'victim_case_form2_tzv3_p24', 'victim_case_form2_tzv3_p1']) === 'y,y,n',
            `HTTP ${r4.status} ${JSON.stringify(r4.json)}; p3,p24,p1=${colCodes(['victim_case_form2_tzv3_p3', 'victim_case_form2_tzv3_p24', 'victim_case_form2_tzv3_p1'])}`);

        // C5 · Un Caso que no existe ya no responde 200 ni marca nada (antes: 200 con un resultado calculado).
        const r5 = await tamizaje({ TZ_1: 'true', ...branch('TZNP_', []) }, false, 'no-existe-xyz');
        check('C5', r5.status >= 400, `Caso inexistente → HTTP ${r5.status} ${r5.text.slice(0, 80)}`);

        // ═══ FASE 2 — FRONTEND ══════════════════════════════════════════════════════════
        await openCaseUI(page, doc);
        await openModule(page, 'Tamizaje de riesgo');
        await page.locator('table').first().waitFor({ timeout: 10000 });

        // U1 · Estructura: dos tablas «Pregunta | Sí | No», una fila por pregunta, un par de radios por fila.
        const est = await page.evaluate(() => {
            const tablas = [...document.querySelectorAll('table')];
            const heads = tablas.map((t) => [...t.querySelectorAll('thead th')].map((h) => h.innerText.trim()).join('|'));
            const cols = tablas.map((t) => [...t.querySelectorAll('colgroup col')].map((c) => c.style.width).join('/'));
            const filas = tablas.map((t) => [...t.querySelectorAll('tbody tr')]);
            const radiosPorFila = filas.map((fs) => fs.every((tr) => {
                const r = [...tr.querySelectorAll('input[type=radio]')];
                return r.length === 2 && r[0].name === r[1].name && /^Sí — /.test(r[0].getAttribute('aria-label')) && /^No — /.test(r[1].getAttribute('aria-label'));
            }));
            const visibles = filas[0].map((tr) => tr.offsetParent !== null);
            return { n: tablas.length, heads, cols, filas: filas.map((f) => f.length), radiosPorFila, visibles, segundaVisible: tablas[1] ? tablas[1].offsetParent !== null : null };
        });
        check('U1', est.n === 2 && est.heads.every((h) => h === 'Pregunta|Sí|No') && est.cols.every((c) => c === '70%/15%/15%')
            && est.filas[0] === 8 - 1 && est.filas[1] === 25 && est.radiosPorFila.every(Boolean) && est.segundaVisible === false,
            `tablas=${est.n}; encabezados=${JSON.stringify(est.heads)}; anchos=${JSON.stringify(est.cols)}; filas=${JSON.stringify(est.filas)}; radios ok=${JSON.stringify(est.radiosPorFila)}; 2ª tabla visible=${est.segundaVisible}`);

        // U2 · Ciclo de revelado: un «Sí» inicial muestra TZ_5, TZ_6 y la segunda tabla con el botón; todo «No» las oculta.
        const fila = (txt) => page.locator('tr', { hasText: txt });
        const L = { TZ1: '¿Ha aumentado la violencia física en severidad o frecuencia', TZ5: '¿Esta persona le persigue, espía', TZ8: 'motivados por prejuicio' };
        const visTZ5 = async () => fila(L.TZ5).first().isVisible();
        const inicialOcultas = !(await visTZ5());
        await fila(L.TZ1).first().locator('input[type=radio]').nth(0).check();
        const trasSi = { tz5: await visTZ5(), tabla2: await page.locator('table').nth(1).isVisible(), boton: await page.getByRole('button', { name: 'Calcular y guardar Tamizaje' }).isVisible() };
        check('U2', inicialOcultas && trasSi.tz5 && trasSi.tabla2 && trasSi.boton, `antes del «Sí» TZ_5 oculta=${inicialOcultas}; después: TZ_5=${trasSi.tz5} tabla2=${trasSi.tabla2} botón=${trasSi.boton}`);

        // U3 · Guardar desde la pantalla y reabrir: lo respondido se conserva en las tablas.
        await fila(L.TZ8).first().locator('input[type=radio]').nth(0).check();
        await fila('¿La persona agresora está desempleada').first().locator('input[type=radio]').nth(0).check();
        // Se espera la respuesta del POST (el resultado del guardado anterior ya está en pantalla y no sirve de señal).
        await Promise.all([
            page.waitForResponse((r) => r.url().endsWith('/tamizaje') && r.request().method() === 'POST'),
            page.getByRole('button', { name: 'Calcular y guardar Tamizaje' }).click(),
        ]);
        await page.waitForFunction(() => /Puntaje: 3/.test(document.body.innerText), null, { timeout: 8000 }).catch(() => {});
        const resultado = (await page.locator('.alert', { hasText: 'Puntaje:' }).first().innerText()).replace(/\s+/g, ' ');
        const genUi = generalesDb();
        await openCaseUI(page, doc);
        await openModule(page, 'Tamizaje de riesgo');
        const sel = async (txt, i) => fila(txt).first().locator('input[type=radio]').nth(i).isChecked();
        const reabierto = { tz1Si: await sel(L.TZ1, 0), tz2No: await sel('¿Le ha agredido o amenazado utilizando algún arma', 1) };
        check('U3', /Puntaje: 3/.test(resultado) && genUi === 'y,n,n,n,n,n,n,y' && reabierto.tz1Si && reabierto.tz2No,
            `resultado «${resultado.slice(0, 90)}»; generales en BD=${genUi}; al reabrir TZ_1=Sí:${reabierto.tz1Si} TZ_2=No:${reabierto.tz2No}`);

        // U4 · Si el guardado falla, la persona lo sabe (comentario de Felipe en la planeación).
        await page.route('**/api/v1/victim-cases-v3/*/tamizaje', (route) => route.fulfill({ status: 500, contentType: 'application/json', body: '{"error":"error al guardar el tamizaje"}' }));
        // TZ_1 ya viene en «Sí» al reabrir: se pasa por «No» y se vuelve a «Sí» para que el cambio dispare el revelado.
        await fila(L.TZ1).first().locator('input[type=radio]').nth(1).check();
        await fila(L.TZ1).first().locator('input[type=radio]').nth(0).check();
        await page.getByRole('button', { name: 'Calcular y guardar Tamizaje' }).click();
        await page.locator('#registro-v3-error-guardado').waitFor({ timeout: 8000 });
        const aviso = await page.locator('#registro-v3-error-guardado').innerText();
        check('U4', /No se pudo guardar el Tamizaje/.test(aviso), `aviso: «${aviso.slice(0, 120)}»`);
        await op.context.close();
    } catch (e) {
        check('ERROR', false, e && e.stack ? e.stack.split('\n').slice(0, 4).join(' | ') : String(e));
    } finally {
        try { deleteCases(created); } catch (e) { console.log('(limpieza falló: ' + e.message + ')'); }
        await browser.close();
    }
    finish('REG-12 · Tamizaje en tablas');
})();
