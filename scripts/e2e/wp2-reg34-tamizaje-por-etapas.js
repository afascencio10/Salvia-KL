#!/usr/bin/env node
/**
 * WP-2 · T7 · REG-34 — El Tamizaje se responde por etapas: primero las 5 preguntas iniciales (TZ_1, 2, 3, 4, 7);
 * si alguna es «Sí» se habilita el resto; si todas son «No» ya se puede calcular y guardar el Tamizaje.
 *
 * Antes el botón «Calcular y guardar Tamizaje» estaba dentro del bloque del resto, así que con las cinco iniciales
 * en «No» nunca aparecía: el Tamizaje no se guardaba y el Caso no se podía finalizar (GAP de flow-E15). Además
 * al reabrir un Caso el resto no se mostraba aunque alguna inicial fuera «Sí».
 *
 * Orden (qa-process.md): 1) backend — POST /tamizaje con solo las iniciales, reinicio del resto, y el portón de
 * Finalizar; 2) frontend — solo lo que el backend no cubre: cuándo aparece el botón, qué se envía y la reapertura.
 *
 * Uso (con la app corriendo):  node scripts/e2e/wp2-reg34-tamizaje-por-etapas.js
 * Casos: DocsMD/Screens/Registro V3/QA/E11-actualizar-respuesta-y-tamizaje-qa-cases.md (E11-C21…C27; comprobaciones C1…C3, U1…U4).
 */
const { chromium } = require('playwright');
const { sql, deleteCases, check, finish, loginAs, api, loadCatalog, validAnswers, openCaseUI, openModule } = require('./wp2-lib');

const created = [];
const rnd = () => String(900000000 + Math.floor(Math.random() * 99999999));
const INICIALES = ['TZ_1', 'TZ_2', 'TZ_3', 'TZ_4', 'TZ_7'];
const branch = (pref, trues) => Object.fromEntries(Array.from({ length: 24 }, (_, i) => [pref + (i + 1), trues.includes(i + 1) ? 'true' : 'false']));
const solo5 = (extra) => ({ TZ_1: 'false', TZ_2: 'false', TZ_3: 'false', TZ_4: 'false', TZ_7: 'false', ...extra });

(async () => {
    const browser = await chromium.launch();
    try {
        const op = await loginAs(browser, 'op');
        const page = op.page;
        const catalog = await loadCatalog(page);
        const draft = async (doc) => {
            const d = doc || rnd();
            const r = await api(page, 'POST', '/api/v1/victim-cases-v3/draft', { tipo_documento: 'cc', documento: d, telefono: '3009998888' });
            created.push(r.json.icode);
            return { icode: r.json.icode, doc: d };
        };
        const tamizaje = (icode, answers) => api(page, 'POST', `/api/v1/victim-cases-v3/${icode}/tamizaje`, { answers, wasPartner: false, esOsigd: false, subtipoFeminicidio: false, subtipoHomicidioAbyecto: false });
        const getCase = (icode) => api(page, 'GET', `/api/v1/victim-cases-v3/${icode}`);
        const codes = (icode, cols) => sql(`SELECT concat_ws(',', ${cols.map((c) => `(SELECT e.victim_case_form2_enums_code FROM salvia.victim_case_form2_enums e WHERE e.victim_case_form2_enums_id = f.victim_case_form2_tzv3_${c})`).join(', ')})
            FROM salvia.victim_case_form2 f JOIN salvia.victim_case c ON c.victim_case_id = f.victim_case_form2_victim_case WHERE c.victim_case_i_code = '${icode}'`);
        const completo = async () => { // Caso con la tríada completa y todo lo que pide Finalizar, menos el Tamizaje
            const c = await draft();
            const v = validAnswers(catalog, c.doc);
            v.AG.AG_3 = catalog.byCode('victim_case_form2_relationship_with_presumed_aggressor_01', 'am');
            const send = (m, a) => api(page, 'POST', `/api/v1/victim-cases-v3/${c.icode}/seccion/${m}`, { answers: a });
            await send('CI', v.CI);
            await send('DP', { ...v.DP, DP_10: catalog.first('victim_case_form2_assigned_sex_at_birth'), DP_15: catalog.first('victim_case_form2_marital_status') });
            await send('CA', { CA_1: 'Contacto Apoyo', CA_2: '3001112222', CA_3: 'apoyo@example.org', CA_4: catalog.first('victim_case_form2_support_contact_kinship') });
            await send('HV', v.HV);
            await send('AG', v.AG);
            return c;
        };

        // ═══ FASE 1 — BACKEND ═══════════════════════════════════════════════════════════
        // C1 · Solo las cinco iniciales: se guarda, el resto queda en «No» y NO se marca como respondido.
        const a = await draft();
        const r1 = await tamizaje(a.icode, solo5());
        const g1 = await getCase(a.icode);
        const keys1 = g1.json.answeredKeys;
        check('C1', r1.status === 200 && r1.json.score === 0 && r1.json.riskLevel === 'bajo'
            && codes(a.icode, ['g1', 'g2', 'g3', 'g4', 'g7']) === 'n,n,n,n,n'
            && INICIALES.every((k) => keys1.includes(k)) && !['TZ_5', 'TZ_6', 'TZ_8', 'TZNP_1', 'TZP_1'].some((k) => keys1.includes(k)),
            `HTTP ${r1.status} ${JSON.stringify(r1.json)}; respondidas TZ=${keys1.filter((k) => k.startsWith('TZ')).sort().join(',')}`);

        // C2 · Si antes había un resto respondido y luego las iniciales pasan a «No», el resto se reinicia en «No».
        const b = await draft();
        const full = { TZ_1: 'true', TZ_2: 'false', TZ_3: 'false', TZ_4: 'false', TZ_5: 'true', TZ_6: 'false', TZ_7: 'false', TZ_8: 'true', ...branch('TZNP_', [2, 24]) };
        await tamizaje(b.icode, full);
        const antes = codes(b.icode, ['g1', 'g5', 'vpp', 'np2', 'np24']);
        const r2 = await tamizaje(b.icode, solo5());
        const despues = codes(b.icode, ['g1', 'g5', 'vpp', 'np2', 'np24']);
        const a2 = (await getCase(b.icode)).json.answers;
        check('C2', antes === 'y,y,y,y,y' && r2.status === 200 && r2.json.score === 0 && despues === 'n,n,n,n,n' && a2.TZ_5 === 'false' && a2.TZNP_2 === 'false',
            `antes g1,g5,vpp,np2,np24=${antes}; después de guardar solo las iniciales en «No»=${despues}; score=${r2.json.score}`);

        // C3 · Finalizar: sin Tamizaje lo rechaza (módulo TZ); con solo las iniciales en «No» ya pasa a Terminado.
        const f = await completo();
        const fin = (icode) => api(page, 'POST', `/api/v1/victim-cases-v3/${icode}/finalizar`);
        const gs = { GS_1: 'false', GS_2: '', GS_4: 'Bitácora inicial de prueba', GS_5: 'false', GS_14: 'false', GS_16: 'false', GS_17: 'false', GS_18: 'false' };
        await api(page, 'POST', `/api/v1/victim-cases-v3/${f.icode}/seccion/GS`, { answers: gs });
        const sinTz = await fin(f.icode);
        await tamizaje(f.icode, solo5());
        const conTz = await fin(f.icode);
        const estado = sql(`SELECT victim_case_status FROM salvia.victim_case WHERE victim_case_i_code = '${f.icode}'`);
        check('C3', sinTz.status === 400 && (sinTz.json.modulosIncompletos || []).includes('TZ') && conTz.status === 200 && estado === 'te',
            `sin Tamizaje → HTTP ${sinTz.status} módulos incompletos=${JSON.stringify(sinTz.json.modulosIncompletos)}; con las 5 iniciales en «No» → HTTP ${conTz.status}; estado=${estado}`);

        // ═══ FASE 2 — FRONTEND ══════════════════════════════════════════════════════════
        const z = await completo(); // tríada completa, sin Tamizaje
        await openCaseUI(page, z.doc);
        await openModule(page, 'Tamizaje de riesgo');
        const fila = (txt) => page.locator('tr', { hasText: txt }).first();
        const T = {
            TZ1: '¿Ha aumentado la violencia física en severidad o frecuencia', TZ2: '¿Le ha agredido o amenazado utilizando algún arma',
            TZ3: '¿Ha amenazado con matarle', TZ4: '¿Ha intentado asfixiarle', TZ5: '¿Esta persona le persigue, espía', TZ7: '¿La persona agresora tiene acceso a armas',
        };
        const boton = page.getByRole('button', { name: 'Calcular y guardar Tamizaje' });
        const no = (txt) => fila(txt).locator('input[type=radio]').nth(1).check();
        const btn0 = await boton.isVisible();
        for (const k of ['TZ1', 'TZ2', 'TZ3', 'TZ4']) await no(T[k]);
        const btn4 = await boton.isVisible();
        await no(T.TZ7);
        const btn5 = await boton.isVisible();
        const aviso = await page.locator('p.text-muted', { hasText: 'Las cinco preguntas iniciales están en «No»' }).isVisible();
        const restoOculto = !(await fila(T.TZ5).isVisible()) && !(await page.locator('table').nth(1).isVisible());
        check('U1', !btn0 && !btn4 && btn5 && aviso && restoOculto,
            `botón: sin responder=${btn0}, con 4 iniciales=${btn4}, con las 5 en «No»=${btn5}; aviso=${aviso}; resto oculto=${restoOculto}`);

        // U2 · Guardar: solo viajan las cinco iniciales y el resultado es riesgo bajo, puntaje 0.
        const [req] = await Promise.all([
            page.waitForRequest((r) => r.url().endsWith('/tamizaje') && r.method() === 'POST'),
            boton.click(),
        ]);
        const enviadas = Object.keys(JSON.parse(req.postData()).answers).sort();
        await page.locator('.alert', { hasText: 'Puntaje: 0' }).waitFor({ timeout: 8000 });
        const res = (await page.locator('.alert', { hasText: 'Puntaje:' }).first().innerText()).replace(/\s+/g, ' ');
        const keysZ = (await getCase(z.icode)).json.answeredKeys;
        check('U2', JSON.stringify(enviadas) === JSON.stringify([...INICIALES].sort()) && /Nivel: bajo/.test(res)
            && codes(z.icode, ['g1', 'g2', 'g3', 'g4', 'g5', 'g6', 'g7', 'vpp']) === 'n,n,n,n,n,n,n,n' && !keysZ.includes('TZ_5') && !keysZ.includes('TZNP_1'),
            `se enviaron ${enviadas.join(',')}; resultado «${res.slice(0, 80)}»; TZ_5 respondida=${keysZ.includes('TZ_5')}`);

        // U3 · Reabrir: las cinco en «No», botón visible, resto oculto; un «Sí» lo revela.
        await openCaseUI(page, z.doc);
        await openModule(page, 'Tamizaje de riesgo');
        const nos = [];
        for (const k of ['TZ1', 'TZ2', 'TZ3', 'TZ4', 'TZ7']) nos.push(await fila(T[k]).locator('input[type=radio]').nth(1).isChecked());
        const btnRe = await boton.isVisible();
        const restoRe = await fila(T.TZ5).isVisible();
        await fila(T.TZ3).locator('input[type=radio]').nth(0).check();
        const restoTrasSi = await fila(T.TZ5).isVisible() && await page.locator('table').nth(1).isVisible();
        check('U3', nos.every(Boolean) && btnRe && !restoRe && restoTrasSi, `reabierto: cinco en «No»=${nos.every(Boolean)}; botón=${btnRe}; resto oculto=${!restoRe}; tras un «Sí» el resto aparece=${restoTrasSi}`);

        // U4 · Un Caso guardado con una inicial en «Sí» muestra el resto (y el botón) apenas se abre.
        const w = await completo();
        await tamizaje(w.icode, { ...full, TZ_5: 'false' });
        await openCaseUI(page, w.doc);
        await openModule(page, 'Tamizaje de riesgo');
        await boton.waitFor({ state: 'visible', timeout: 5000 }).catch(() => {});
        const vis4 = { tz5: await fila(T.TZ5).isVisible(), tabla2: await page.locator('table').nth(1).isVisible(), boton: await boton.isVisible(), tz1Si: await fila(T.TZ1).locator('input[type=radio]').nth(0).isChecked() };
        check('U4', vis4.tz5 && vis4.tabla2 && vis4.boton && vis4.tz1Si, `al abrir un Caso con TZ_1=Sí: TZ_5 visible=${vis4.tz5}; 2ª tabla=${vis4.tabla2}; botón=${vis4.boton}; TZ_1=Sí marcado=${vis4.tz1Si}`);
        await op.context.close();
    } catch (e) {
        check('ERROR', false, e && e.stack ? e.stack.split('\n').slice(0, 4).join(' | ') : String(e));
    } finally {
        try { deleteCases(created); } catch (e) { console.log('(limpieza falló: ' + e.message + ')'); }
        await browser.close();
    }
    finish('REG-34 · Tamizaje por etapas');
})();
