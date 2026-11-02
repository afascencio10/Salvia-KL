#!/usr/bin/env node
/**
 * WP-2 · T5 · REG-08 — Pregunta nueva «Tipo de discapacidad» (DP_51), con 7 opciones, que solo aparece
 * cuando la persona responde «Sí» a «¿Tiene dificultades físicas, mentales y/o sensoriales…?» (DP_2).
 *
 * Orden (qa-process.md): 1) backend — catálogo, migración SQL (idempotente), guardado, columna y BD;
 * 2) frontend — solo lo que el backend no cubre: cuándo aparece la pregunta, sus opciones y que se limpie
 * al pasar DP_2 a «No».
 *
 * Uso (con la app corriendo y la migración aplicada):  node cloudDevelopment/e2e/wp2-reg08-tipo-discapacidad.js
 * Casos C1…C8 y U1…U3: ver DocsMD/Screens/Registro V3/QA/E12-guardar-seccion-qa-cases.md (sección REG-08).
 */
const path = require('path');
const { execFileSync } = require('child_process');
const { chromium } = require('playwright');
const { sql, deleteCases, check, finish, loginAs, api, loadCatalog, validAnswers, openCaseUI, openModule, question } = require('./wp2-lib');

const MIGRATION = path.join(__dirname, '..', '..', 'src', 'migrations', '20261102_wp2_registro_v3_catalogos_y_columnas.sql');
const CAT = 'victim_case_form2_disability_type';
const LABELS = ['Discapacidad física', 'Discapacidad auditiva', 'Discapacidad visual', 'Sordoceguera', 'Discapacidad intelectual', 'Discapacidad psicosocial (mental)', 'Discapacidad múltiple'];
const DP2 = '¿Tiene dificultades físicas, mentales y/o sensoriales en su vida diaria?';
const DP51 = 'Tipo de discapacidad';
const created = [];
const rnd = () => String(900000000 + Math.floor(Math.random() * 99999999));

(async () => {
    const browser = await chromium.launch();
    try {
        const op = await loginAs(browser, 'op');
        const page = op.page;
        const catalog = await loadCatalog(page);
        const opts = catalog.cat[CAT] || [];
        const optId = (code) => sql(`SELECT victim_case_form2_enums_id FROM salvia.victim_case_form2_enums WHERE victim_case_form2_enums_category = '${CAT}' AND victim_case_form2_enums_code = '${code}'`);
        const draft = async (doc) => {
            const d = doc || rnd();
            const r = await api(page, 'POST', '/api/v1/victim-cases-v3/draft', { tipo_documento: 'cc', documento: d, telefono: '3009998888' });
            created.push(r.json.icode);
            return { icode: r.json.icode, doc: d };
        };
        const colDb = (icode) => sql(`SELECT COALESCE(f.victim_case_form2_disability_type::text, '<NULL>') FROM salvia.victim_case_form2 f
            JOIN salvia.victim_case c ON c.victim_case_id = f.victim_case_form2_victim_case WHERE c.victim_case_i_code = '${icode}'`);
        const saveDP = (icode, extra) => api(page, 'POST', `/api/v1/victim-cases-v3/${icode}/seccion/DP`, { answers: extra });
        const getCase = (icode) => api(page, 'GET', `/api/v1/victim-cases-v3/${icode}`);

        // ═══ FASE 1 — BACKEND ═══════════════════════════════════════════════════════════
        check('C1', JSON.stringify(opts.map((o) => [o.code, o.label])) === JSON.stringify(LABELS.map((l, i) => ['0' + (i + 1), l])),
            `GET /enums → ${opts.length} opciones: ${opts.map((o) => o.code + '=' + o.label).join(' | ')}`);

        const cntBefore = sql(`SELECT count(*) FROM salvia.victim_case_form2_enums WHERE victim_case_form2_enums_category = '${CAT}'`);
        let rc = 0;
        try { execFileSync('psql', ['-h', '127.0.0.1', '-U', 'postgres', '-d', 'salvia_local', '-X', '-q', '-v', 'ON_ERROR_STOP=1', '-f', MIGRATION], { env: { ...process.env, PGPASSWORD: 'postgres' }, stdio: 'ignore' }); } catch (e) { rc = e.status || 1; }
        const cntAfter = sql(`SELECT count(*) FROM salvia.victim_case_form2_enums WHERE victim_case_form2_enums_category = '${CAT}'`);
        const colN = sql(`SELECT count(*) FROM information_schema.columns WHERE table_schema='salvia' AND table_name='victim_case_form2' AND column_name='victim_case_form2_disability_type'`);
        check('C2', rc === 0 && cntBefore === '7' && cntAfter === '7' && colN === '1', `reaplicar la migración: salida=${rc}; filas ${cntBefore}→${cntAfter}; columna presente ${colN} vez`);

        const a = await draft();
        const g0 = await getCase(a.icode);
        check('C3', colDb(a.icode) === '<NULL>' && g0.json.answers.DP_51 === undefined,
            `borrador nuevo: columna=${colDb(a.icode)} (NULL = sin respuesta, no hay valor marcador); GET DP_51=${g0.json.answers.DP_51}`);

        const fisica = opts.find((o) => o.code === '01');
        const multiple = opts.find((o) => o.code === '07');
        const r4 = await saveDP(a.icode, { DP_2: 'true', DP_51: fisica.icode });
        const g4 = await getCase(a.icode);
        check('C4', r4.status === 200 && colDb(a.icode) === optId('01') && g4.json.answers.DP_51 === fisica.icode && g4.json.answeredKeys.includes('DP_51'),
            `guardar «${fisica.label}» → HTTP ${r4.status}; BD=${colDb(a.icode)} (id 01=${optId('01')}); GET DP_51===icode: ${g4.json.answers.DP_51 === fisica.icode}`);

        const r5 = await saveDP(a.icode, { DP_51: multiple.icode });
        check('C5', r5.status === 200 && colDb(a.icode) === optId('07'), `cambiar a «${multiple.label}» → HTTP ${r5.status}; BD=${colDb(a.icode)} (id 07=${optId('07')})`);

        const r6 = await saveDP(a.icode, { DP_2: 'false', DP_51: '' });
        const g6 = await getCase(a.icode);
        check('C6', r6.status === 200 && colDb(a.icode) === '<NULL>' && g6.json.answers.DP_51 === undefined,
            `limpiar (DP_2=No, DP_51 vacío) → HTTP ${r6.status}; BD=${colDb(a.icode)}; GET DP_51=${g6.json.answers.DP_51}`);

        const r7 = await saveDP(a.icode, { DP_51: 'no-existe' });
        check('C7', r7.status === 400 && r7.json.campo === 'DP_51' && colDb(a.icode) === '<NULL>', `icode inexistente → HTTP ${r7.status} campo=${r7.json && r7.json.campo}; BD=${colDb(a.icode)}`);

        const shared = rnd();
        const s1 = await draft(shared);
        const s2 = await draft(shared);
        const r8 = await saveDP(s1.icode, { DP_2: 'true', DP_51: fisica.icode });
        check('C8', r8.status === 200 && colDb(s2.icode) === optId('01'), `el Perfil se replica al otro Caso del mismo documento: hermano BD=${colDb(s2.icode)} (id 01=${optId('01')})`);

        // ═══ FASE 2 — FRONTEND ══════════════════════════════════════════════════════════
        // Caso con DP completo (DP_2 = No) para que «Guardar sección» pase las obligatorias.
        const b = await draft();
        const v = validAnswers(catalog, b.doc);
        await api(page, 'POST', `/api/v1/victim-cases-v3/${b.icode}/seccion/CI`, { answers: v.CI });
        await saveDP(b.icode, v.DP);
        await openCaseUI(page, b.doc);
        await openModule(page, 'Datos personales');
        const q51 = question(page, DP51);
        const yes = question(page, DP2).locator('.form-check', { hasText: 'Sí' }).locator('input');
        const no = question(page, DP2).locator('.form-check', { hasText: 'No' }).locator('input');
        await question(page, DP2).waitFor({ timeout: 10000 });

        const ocultaAntes = (await q51.count()) === 0;
        await yes.check();
        await q51.waitFor({ timeout: 5000 });
        const texts = await q51.locator('select option').allInnerTexts();
        await no.check();
        const ocultaDespues = (await question(page, DP51).count()) === 0;
        check('U1', ocultaAntes && ocultaDespues && JSON.stringify(texts) === JSON.stringify(['Seleccione...', ...LABELS]),
            `con DP_2=No no aparece (${ocultaAntes}); con Sí aparece con ${JSON.stringify(texts)}; al volver a No se oculta (${ocultaDespues})`);

        await yes.check();
        await question(page, DP51).locator('select').selectOption({ label: 'Sordoceguera' });
        await page.getByRole('button', { name: 'Guardar sección' }).click();
        await page.locator('text=Sección guardada.').waitFor({ timeout: 8000 });
        const saved = colDb(b.icode);
        await openCaseUI(page, b.doc);
        await openModule(page, 'Datos personales');
        const reYes = await question(page, DP2).locator('.form-check', { hasText: 'Sí' }).locator('input').isChecked();
        const reSel = await question(page, DP51).locator('select').inputValue();
        const icode04 = opts.find((o) => o.code === '04').icode;
        check('U2', saved === optId('04') && reYes && reSel === icode04, `elegir «Sordoceguera» y guardar → BD=${saved} (id 04=${optId('04')}); al reabrir DP_2=Sí: ${reYes}, selección===icode: ${reSel === icode04}`);

        await question(page, DP2).locator('.form-check', { hasText: 'No' }).locator('input').check();
        const oculta3 = (await question(page, DP51).count()) === 0;
        await page.getByRole('button', { name: 'Guardar sección' }).click();
        await page.locator('text=Sección guardada.').waitFor({ timeout: 8000 });
        const cleared = colDb(b.icode);
        await openCaseUI(page, b.doc);
        await openModule(page, 'Datos personales');
        const reHidden = (await question(page, DP51).count()) === 0;
        check('U3', oculta3 && cleared === '<NULL>' && reHidden, `pasar DP_2 a No: la pregunta se oculta (${oculta3}); al guardar la BD queda ${cleared}; al reabrir sigue oculta (${reHidden})`);
        await op.context.close();
    } catch (e) {
        check('ERROR', false, e && e.stack ? e.stack.split('\n').slice(0, 4).join(' | ') : String(e));
    } finally {
        try { deleteCases(created); } catch (e) { console.log('(limpieza falló: ' + e.message + ')'); }
        await browser.close();
    }
    finish('REG-08 · tipo de discapacidad');
})();
