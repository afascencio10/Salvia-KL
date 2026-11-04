#!/usr/bin/env node
/**
 * WP-2 · T2 · REG-07 — «¿En qué zona se encuentra su lugar de residencia?» (DP_30) con las opciones
 * Rural y Urbano.
 *
 * Orden (qa-process.md): 1) backend — catálogo, migración SQL, guardado y BD; 2) frontend — solo lo
 * que el backend no cubre: la lista desplegable que ve la persona.
 *
 * Uso (con la app corriendo y la migración aplicada):
 *   psql ... -f src/migrations/20261102_wp2_registro_v3_catalogos_y_columnas.sql
 *   node scripts/e2e/wp2-reg07-zona-residencia.js
 *
 * Casos: DocsMD/Screens/Registro V3/QA/E12-guardar-seccion-qa-cases.md (E12-C21…C24, API) y E11-actualizar-respuesta-y-tamizaje-qa-cases.md (E11-C4, C5, pantalla).
 */
const path = require('path');
const { chromium } = require('playwright');
const { sql, sqlFile, deleteCases, check, finish, loginAs, api, openRegistroV3, newDraftUI, openModule, question } = require('./wp2-lib');

const MIGRATION = path.join(__dirname, '..', '..', 'src', 'migrations', '20261102_wp2_registro_v3_catalogos_y_columnas.sql');
const ZONE = 'victim_case_form2_residence_zone';
const created = [];

const zoneRowId = (code) => sql(`SELECT victim_case_form2_enums_id FROM salvia.victim_case_form2_enums
                                  WHERE victim_case_form2_enums_category = '${ZONE}' AND victim_case_form2_enums_code = '${code}'`);
const storedZone = (icode) => sql(`SELECT f.victim_case_form2_residence_zone FROM salvia.victim_case_form2 f
                                   JOIN salvia.victim_case c ON c.victim_case_id = f.victim_case_form2_victim_case
                                   WHERE c.victim_case_i_code = '${icode}'`);

(async () => {
    const browser = await chromium.launch();
    try {
        const op = await loginAs(browser, 'op');

        // ═══ FASE 1 — BACKEND ═══════════════════════════════════════════════════════════
        // C1 · el catálogo trae exactamente Rural y Urbano, con texto en español.
        const cat = await api(op.page, 'GET', '/api/v1/victim-cases-v3/enums');
        const zone = (cat.json && cat.json[ZONE]) || [];
        check('C1', cat.status === 200 && JSON.stringify(zone.map((o) => [o.code, o.label])) === JSON.stringify([['ru', 'Rural'], ['ur', 'Urbano']])
            && zone.every((o) => o.icode && o.icode.length === 36),
            `GET /enums → ${zone.length} opciones: ${JSON.stringify(zone.map((o) => o.code + '=' + o.label))}`);

        // C2 · la zona del hecho (facts_zone) no cambió.
        const facts = (cat.json && cat.json['victim_case_form2_facts_zone']) || [];
        check('C2', facts.length === 3 && facts.map((o) => o.code).join() === 'cm,cp,pr'
            && facts.every((o) => !o.label.startsWith('victim_case_form2_')),
            `facts_zone sigue con ${facts.length} opciones: ${facts.map((o) => o.label).join(' | ')}`);

        // C3 · la migración es idempotente: aplicarla de nuevo no duplica ni falla.
        const before = sql(`SELECT count(*) FROM salvia.victim_case_form2_enums WHERE victim_case_form2_enums_category = '${ZONE}'`);
        const rc = sqlFile(MIGRATION);
        const after = sql(`SELECT count(*) FROM salvia.victim_case_form2_enums WHERE victim_case_form2_enums_category = '${ZONE}'`);
        check('C3', rc === 0 && before === '2' && after === '2', `reaplicar la migración: salida=${rc}, filas antes=${before} después=${after}`);

        // C4/C5 · guardar la zona en un borrador y verificarlo en la BD y al reabrir.
        const draft = await api(op.page, 'POST', '/api/v1/victim-cases-v3/draft', { tipo_documento: 'cc', documento: String(900000000 + Math.floor(Math.random() * 99999999)), telefono: '3009998888' });
        const icode = draft.json && draft.json.icode;
        if (icode) created.push(icode);
        const rural = zone.find((o) => o.code === 'ru');
        const urbano = zone.find((o) => o.code === 'ur');
        const s1 = await api(op.page, 'POST', `/api/v1/victim-cases-v3/${icode}/seccion/DP`, { answers: { DP_30: rural.icode } });
        const c1 = await api(op.page, 'GET', `/api/v1/victim-cases-v3/${icode}`);
        check('C4', s1.status === 200 && storedZone(icode) === zoneRowId('ru') && c1.json.answers.DP_30 === rural.icode
            && c1.json.answeredKeys.includes('DP_30'),
            `guardar Rural → HTTP ${s1.status}; BD=${storedZone(icode)} (id Rural=${zoneRowId('ru')}); reabrir DP_30===Rural: ${c1.json.answers.DP_30 === rural.icode}`);
        const s2 = await api(op.page, 'POST', `/api/v1/victim-cases-v3/${icode}/seccion/DP`, { answers: { DP_30: urbano.icode } });
        const c2 = await api(op.page, 'GET', `/api/v1/victim-cases-v3/${icode}`);
        check('C5', s2.status === 200 && storedZone(icode) === zoneRowId('ur') && c2.json.answers.DP_30 === urbano.icode,
            `cambiar a Urbano → HTTP ${s2.status}; BD=${storedZone(icode)} (id Urbano=${zoneRowId('ur')}); reabrir DP_30===Urbano: ${c2.json.answers.DP_30 === urbano.icode}`);

        // ═══ FASE 2 — FRONTEND (solo la lista desplegable) ══════════════════════════════
        await openRegistroV3(op.page);
        const uiCode = await newDraftUI(op.page);
        created.push(uiCode);
        await openModule(op.page, 'Datos personales');
        const sel = question(op.page, '¿En qué zona se encuentra su lugar de residencia?').locator('select');
        await sel.waitFor({ timeout: 10000 });
        const options = await sel.locator('option').allInnerTexts();
        check('C6', JSON.stringify(options) === JSON.stringify(['Seleccione...', 'Rural', 'Urbano']),
            `opciones en pantalla: ${JSON.stringify(options)}`);
        await sel.selectOption({ label: 'Rural' });
        check('C7', (await sel.inputValue()) === rural.icode, 'al elegir «Rural» el valor del campo es el icode de Rural');
        await openModule(op.page, 'Contacto e identificación');
        await openModule(op.page, 'Datos personales');
        check('C8', (await question(op.page, '¿En qué zona se encuentra su lugar de residencia?').locator('select').inputValue()) === rural.icode,
            'al salir del módulo y volver, la selección «Rural» se conserva');
        await op.context.close();
    } catch (e) {
        check('ERROR', false, e && e.stack ? e.stack.split('\n').slice(0, 3).join(' | ') : String(e));
    } finally {
        try { deleteCases(created); } catch (e) { console.log('(limpieza falló: ' + e.message + ')'); }
        await browser.close();
    }
    finish('REG-07 · zona de residencia');
})();
