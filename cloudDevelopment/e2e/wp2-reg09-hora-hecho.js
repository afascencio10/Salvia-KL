#!/usr/bin/env node
/**
 * WP-2 · T4 · REG-09 — «¿A qué hora ocurrió el hecho?» (HV_3) es un selector de hora, no un texto libre.
 *
 * Antes era un campo de texto («Ej: 19:30»): se podía escribir «7pm» o dejarlo vacío y el guardado
 * reventaba con un 500. Ahora es <input type="time">: la persona solo puede elegir una hora válida.
 *
 * Orden (qa-process.md): 1) backend — los bordes de la hora (medianoche, 23:59, fuera de rango) ya que el
 * resto (19:30, hora con segundos, «7pm») lo cubre wp2-reg11; 2) frontend — el control y su ida y vuelta.
 *
 * Uso (con la app corriendo):  node cloudDevelopment/e2e/wp2-reg09-hora-hecho.js
 * Casos: DocsMD/Screens/Registro V3/QA/E12-guardar-seccion-qa-cases.md (E12-C12, C13; comprobaciones C1…C3) y E11-actualizar-respuesta-y-tamizaje-qa-cases.md (E11-C6…C9; C4…C7).
 */
const { chromium } = require('playwright');
const { sql, deleteCases, check, finish, loginAs, api, loadCatalog, validAnswers, openCaseUI, openModule, question } = require('./wp2-lib');

const created = [];
const HORA = '¿A qué hora ocurrió el hecho?';
const rnd = () => String(900000000 + Math.floor(Math.random() * 99999999));

(async () => {
    const browser = await chromium.launch();
    try {
        const op = await loginAs(browser, 'op');
        const page = op.page;
        const catalog = await loadCatalog(page);
        const doc = rnd();
        const d = await api(page, 'POST', '/api/v1/victim-cases-v3/draft', { tipo_documento: 'cc', documento: doc, telefono: '3009998888' });
        const icode = d.json.icode;
        created.push(icode);
        const hv = validAnswers(catalog, doc).HV;
        const horaDb = () => sql(`SELECT f.victim_case_form2_facts_start_time::text FROM salvia.victim_case_form2 f
            JOIN salvia.victim_case c ON c.victim_case_id = f.victim_case_form2_victim_case WHERE c.victim_case_i_code = '${icode}'`);
        const saveHV = (extra) => api(page, 'POST', `/api/v1/victim-cases-v3/${icode}/seccion/HV`, { answers: { ...hv, ...extra } });
        const horaGet = async () => ((await api(page, 'GET', `/api/v1/victim-cases-v3/${icode}`)).json.answers || {}).HV_3;

        // ═══ FASE 1 — BACKEND (bordes de la hora) ═══════════════════════════════════════
        const r1 = await saveHV({ HV_3: '00:00' });
        check('C1', r1.status === 200 && horaDb() === '00:00:00' && (await horaGet()) === '00:00',
            `medianoche → HTTP ${r1.status}; BD=${horaDb()}; GET=${await horaGet()} (00:00 es una respuesta real, no el marcador)`);
        const r2 = await saveHV({ HV_3: '23:59' });
        check('C2', r2.status === 200 && horaDb() === '23:59:00' && (await horaGet()) === '23:59', `23:59 → HTTP ${r2.status}; BD=${horaDb()}; GET=${await horaGet()}`);
        const r3 = await saveHV({ HV_3: '25:00' });
        check('C3', r3.status === 400 && r3.json.campo === 'HV_3' && horaDb() === '23:59:00',
            `25:00 → HTTP ${r3.status} campo=${r3.json && r3.json.campo}; la hora guardada no cambió (${horaDb()})`);

        // ═══ FASE 2 — FRONTEND ══════════════════════════════════════════════════════════
        await saveHV({ HV_3: '19:30' });
        await openCaseUI(page, doc);
        await openModule(page, 'Hechos de violencia');
        const ctl = question(page, HORA).locator('input');
        await ctl.waitFor({ timeout: 10000 });
        const attrs = await ctl.evaluate((e) => ({ type: e.type, step: e.getAttribute('step'), placeholder: e.getAttribute('placeholder'), value: e.value }));
        check('C4', attrs.type === 'time' && attrs.step === '60' && !attrs.placeholder && attrs.value === '19:30',
            `control: type=${attrs.type} step=${attrs.step} placeholder=${attrs.placeholder} valor guardado=${attrs.value}`);

        // C5 · elegir una hora y guardar: queda en la BD y se ve al reabrir.
        await ctl.fill('07:05');
        await page.getByRole('button', { name: 'Guardar sección' }).click();
        await page.locator('text=Sección guardada.').waitFor({ timeout: 8000 });
        const db5 = horaDb();
        await openCaseUI(page, doc);
        await openModule(page, 'Hechos de violencia');
        const re5 = await question(page, HORA).locator('input').inputValue();
        check('C5', db5 === '07:05:00' && re5 === '07:05', `elegir 07:05 → BD=${db5}; al reabrir el control muestra ${re5}`);

        // C6 · medianoche por pantalla: se guarda y al reabrir NO queda en blanco.
        await question(page, HORA).locator('input').fill('00:00');
        await page.getByRole('button', { name: 'Guardar sección' }).click();
        await page.locator('text=Sección guardada.').waitFor({ timeout: 8000 });
        await openCaseUI(page, doc);
        await openModule(page, 'Hechos de violencia');
        const re6 = await question(page, HORA).locator('input').inputValue();
        check('C6', horaDb() === '00:00:00' && re6 === '00:00', `elegir 00:00 → BD=${horaDb()}; al reabrir muestra "${re6}"`);

        // C7 · el control no admite texto libre como «7pm».
        let rechazo = '';
        try { await question(page, HORA).locator('input').fill('7pm', { timeout: 2000 }); } catch (e) { rechazo = String(e.message).split('\n')[0]; }
        check('C7', /Malformed value/i.test(rechazo), `escribir «7pm» en el control → ${rechazo || 'se aceptó (mal)'}`);
        await op.context.close();
    } catch (e) {
        check('ERROR', false, e && e.stack ? e.stack.split('\n').slice(0, 4).join(' | ') : String(e));
    } finally {
        try { deleteCases(created); } catch (e) { console.log('(limpieza falló: ' + e.message + ')'); }
        await browser.close();
    }
    finish('REG-09 · hora del hecho');
})();
