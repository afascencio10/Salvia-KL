#!/usr/bin/env node
/**
 * WP-2 · T8 · REG-11 — Guardar una sección de Registro V3 es atómico, valida lo que se escribe y, si
 * algo falla, la persona lo sabe y puede corregirlo (comentario de Felipe en la planeación: «si algo falla
 * el usuario lo sepa; si es por el formato del dato, pedirle que lo corrija»).
 *
 * Antes: el municipio de más de 8 letras, una hora como «7pm», una fecha o lista vacía o una dirección de más
 * de 150 caracteres daban un 500 genérico; markKeysAnswered corría ANTES del UPDATE, así que el módulo
 * quedaba «respondido» con los datos sin guardar y el Caso pasaba a Iniciado con los hechos en blanco; las
 * fechas se devolvían en formato crudo de Go (el campo quedaba en blanco al reabrir) y los departamentos de
 * los hechos y de residencia no se guardaban.
 *
 * Orden (qa-process.md): 1) backend — endpoint SaveSection/GetCase + verificación directa en la BD + logs;
 * 2) frontend — solo lo que el backend no cubre: lo que ve la persona (aviso en rojo, reapertura).
 *
 * Uso (con la app corriendo y la migración aplicada):
 *   APP_LOG=/ruta/al/log/de/la/app node cloudDevelopment/e2e/wp2-reg11-guardado-seccion.js   (APP_LOG es opcional)
 * Casos: ver DocsMD/Screens/Registro V3/QA/E12-guardar-seccion-qa-cases.md.
 */
const fs = require('fs');
const { chromium } = require('playwright');
const { sql, deleteCases, check, finish, loginAs, api, openRegistroV3, question } = require('./wp2-lib');

const created = [];
const rnd = () => String(900000000 + Math.floor(Math.random() * 99999999));
let rejected400 = 0; // 400 que este script provoca a propósito (se comparan con las líneas del log)

(async () => {
    const browser = await chromium.launch();
    const logCount = (marca) => (process.env.APP_LOG && fs.existsSync(process.env.APP_LOG))
        ? fs.readFileSync(process.env.APP_LOG, 'utf8').split('\n').filter((l) => l.includes(marca)).length : 0;
    const log0 = { rej: logCount('[RegistroV3] SaveSection rechazado'), fail: logCount('[RegistroV3] SaveSection falló') };
    try {
        const op = await loginAs(browser, 'op');
        const page = op.page;

        // ── utilidades ───────────────────────────────────────────────────────────────────
        const draft = async (doc) => {
            const d = doc || rnd();
            const r = await api(page, 'POST', '/api/v1/victim-cases-v3/draft', { tipo_documento: 'cc', documento: d, telefono: '3009998888' });
            created.push(r.json.icode);
            return { icode: r.json.icode, doc: d };
        };
        const save = async (icode, modulo, answers) => {
            const r = await api(page, 'POST', `/api/v1/victim-cases-v3/${icode}/seccion/${modulo}`, { answers });
            if (r.status === 400) rejected400++;
            return r;
        };
        const getCase = (icode) => api(page, 'GET', `/api/v1/victim-cases-v3/${icode}`);
        const form2 = (icode, col) => sql(`SELECT COALESCE(f.${col}::text, '<NULL>') FROM salvia.victim_case_form2 f
            JOIN salvia.victim_case c ON c.victim_case_id = f.victim_case_form2_victim_case WHERE c.victim_case_i_code = '${icode}'`);
        const vcase = (icode, col) => sql(`SELECT COALESCE(${col}::text, '<NULL>') FROM salvia.victim_case WHERE victim_case_i_code = '${icode}'`);

        const cat = (await api(page, 'GET', '/api/v1/victim-cases-v3/enums')).json;
        const first = (c) => (cat[c] && cat[c][0] && cat[c][0].icode) || '';
        const byCode = (c, code) => ((cat[c] || []).find((o) => o.code === code) || {}).icode || '';
        const valid = (doc) => ({
            CI: { CI_1: byCode('victim_case_form2_pronoun', 'la'), CI_2: 'Ana Perez', CI_3: 'Anita', CI_4: '3001234567', CI_5: byCode('victim_case_form2_victim_doc_type', 'cc'), CI_6: doc, CI_7: 'false', CI_9: 'false' },
            DP: { DP_1: '1990-05-20', DP_2: 'false', DP_5: first('victim_case_form2_nationality'), DP_8: first('victim_case_form2_gender_identity'), DP_9: first('victim_case_form2_sexual_orientation'), DP_11: first('victim_case_form2_specially_protected_population'), DP_12: first('victim_case_form2_ethnic_affiliation'), DP_14: 'false', DP_17: first('victim_case_form2_occupation'), DP_27: 'Atlántico', DP_28: 'Barranquilla', DP_29: 'Calle 1 # 2-3' },
            HV: { HV_1: 'Descripción de prueba de los hechos', HV_2: '2026-10-30', HV_3: '19:30', HV_4: 'Atlántico', HV_5: 'Barranquilla', HV_6: first('victim_case_form2_facts_zone'), HV_7: 'Calle 45 # 12-34', HV_8: first('victim_case_form2_scenario_violence'), HV_9: first('victim_case_form2_type_violence_experienced'), HV_11: first('victim_case_form2_scope_of_violence'), HV_13: 'true', HV_14: first('victim_case_form2_recurrence_aggression') },
            AG: { AG_1: first('victim_case_form2_num_agressors'), AG_2: first('victim_case_form2_proximity_principal_aggressor'), AG_3: first('victim_case_form2_relationship_with_presumed_aggressor_01') },
        });

        // ═══ FASE 1 — BACKEND ═══════════════════════════════════════════════════════════
        // C1 · Hechos completos: fecha/hora/departamento/municipio largo se guardan y se devuelven en el formato de la pantalla.
        const a = await draft();
        const p = valid(a.doc);
        const r1 = await save(a.icode, 'HV', p.HV);
        const g1 = await getCase(a.icode);
        check('C1', r1.status === 200 && form2(a.icode, 'victim_case_form2_facts_date') === '2026-10-30'
            && form2(a.icode, 'victim_case_form2_facts_start_time') === '19:30:00'
            && form2(a.icode, 'victim_case_form2_facts_department') === 'Atlántico'
            && form2(a.icode, 'victim_case_form2_facts_town_code') === 'Barranquilla'
            && g1.json.answers.HV_2 === '2026-10-30' && g1.json.answers.HV_3 === '19:30'
            && g1.json.answers.HV_4 === 'Atlántico' && g1.json.answers.HV_5 === 'Barranquilla',
            `HTTP ${r1.status}; BD fecha=${form2(a.icode, 'victim_case_form2_facts_date')} hora=${form2(a.icode, 'victim_case_form2_facts_start_time')} depto=${form2(a.icode, 'victim_case_form2_facts_department')} municipio=${form2(a.icode, 'victim_case_form2_facts_town_code')}; GET HV_2=${g1.json.answers.HV_2} HV_3=${g1.json.answers.HV_3}`);

        // C2 · Residencia: departamento y municipio largo (también en victim_case).
        const r2 = await save(a.icode, 'DP', p.DP);
        const g2 = await getCase(a.icode);
        check('C2', r2.status === 200 && form2(a.icode, 'victim_case_form2_residence_department') === 'Atlántico'
            && form2(a.icode, 'victim_case_form2_residence_town') === 'Barranquilla'
            && vcase(a.icode, 'victim_case_victim_town_code') === 'Barranquilla'
            && g2.json.answers.DP_1 === '1990-05-20' && g2.json.answers.DP_27 === 'Atlántico' && g2.json.answers.DP_28 === 'Barranquilla',
            `HTTP ${r2.status}; residencia depto=${form2(a.icode, 'victim_case_form2_residence_department')} municipio=${form2(a.icode, 'victim_case_form2_residence_town')}; GET DP_1=${g2.json.answers.DP_1}`);

        // C3 · Datos con formato/tamaño inválido → 400 con el campo y el motivo (antes: 500 genérico).
        const bad = [
            ['HV', 'HV_3', '7pm'], ['HV', 'HV_2', '30/10/2026'], ['HV', 'HV_2', '2026-02-30'],
            ['HV', 'HV_7', 'x'.repeat(151)], ['HV', 'HV_6', 'no-existe'], ['CA', 'CA_2', 'abc'],
            ['CI', 'CI_6', '9'.repeat(33)], ['CI', 'CI_2', 'N'.repeat(33) + ' Perez'],
        ];
        for (const [modulo, key, value] of bad) {
            const b = await draft();
            const r = await save(b.icode, modulo, { [key]: value });
            check(`C3:${key}=${value.length > 12 ? value.slice(0, 5) + '…(' + value.length + ')' : value}`,
                r.status === 400 && r.json && r.json.error === 'valor no válido' && r.json.campo === key && !!r.json.detalle,
                `HTTP ${r.status} ${r.text.slice(0, 150)}`);
        }

        // C4 · Lo que se puede interpretar se acepta: teléfono con espacios, fecha en el formato viejo, hora con segundos.
        const c = await draft();
        const r4a = await save(c.icode, 'CA', { CA_2: '300 123 4567' });
        const r4b = await save(c.icode, 'HV', { HV_2: '2026-10-31 00:00:00 +0000 UTC', HV_3: '08:05:30' });
        const g4 = await getCase(c.icode);
        check('C4', r4a.status === 200 && form2(c.icode, 'victim_case_form2_support_contact_phone') === '3001234567'
            && r4b.status === 200 && form2(c.icode, 'victim_case_form2_facts_date') === '2026-10-31'
            && form2(c.icode, 'victim_case_form2_facts_start_time') === '08:05:30'
            && g4.json.answers.HV_2 === '2026-10-31' && g4.json.answers.HV_3 === '08:05',
            `teléfono→${form2(c.icode, 'victim_case_form2_support_contact_phone')}; fecha vieja→${form2(c.icode, 'victim_case_form2_facts_date')}; hora con seg→${form2(c.icode, 'victim_case_form2_facts_start_time')} (GET ${g4.json.answers.HV_3})`);

        // C5 · ATOMICIDAD: un campo inválido junto a uno válido → no se guarda NADA ni se marca nada como respondido.
        const d = await draft();
        const before = { desc: form2(d.icode, 'victim_case_form2_facts_description'), keys: vcase(d.icode, 'victim_case_v3_answered_keys') };
        const r5 = await save(d.icode, 'HV', { HV_1: 'TEXTO NUEVO QUE NO DEBE QUEDAR', HV_3: '7pm' });
        const g5 = await getCase(d.icode);
        check('C5', r5.status === 400 && form2(d.icode, 'victim_case_form2_facts_description') === before.desc
            && vcase(d.icode, 'victim_case_v3_answered_keys') === before.keys && !g5.json.answeredKeys.some((k) => k === 'HV_1' || k === 'HV_3'),
            `HTTP ${r5.status}; descripción antes="${before.desc}" después="${form2(d.icode, 'victim_case_form2_facts_description')}"; claves respondidas=${JSON.stringify(g5.json.answeredKeys)}`);

        // C6 · El Caso NO pasa a Iniciado con los hechos sin guardar (el caso real de la reproducción).
        const e = await draft();
        const pe = valid(e.doc);
        await save(e.icode, 'CI', pe.CI);
        await save(e.icode, 'DP', pe.DP);
        const r6a = await save(e.icode, 'HV', { ...pe.HV, HV_3: '7pm' });
        const r6b = await save(e.icode, 'AG', pe.AG);
        const g6 = await getCase(e.icode);
        check('C6a', r6a.status === 400 && r6b.status === 200 && r6b.json.triadaCompleta === false && g6.json.triadaCompleta === false
            && vcase(e.icode, 'victim_case_status') === 'bo' && !g6.json.answeredKeys.some((k) => k.startsWith('HV_')),
            `HV inválido→${r6a.status}; AG→${r6b.status} triadaCompleta=${r6b.json.triadaCompleta}; estado=${vcase(e.icode, 'victim_case_status')}; claves HV respondidas=${g6.json.answeredKeys.filter((k) => k.startsWith('HV_')).length}`);
        const r6c = await save(e.icode, 'HV', pe.HV);
        const g6c = await getCase(e.icode);
        check('C6b', r6c.status === 200 && r6c.json.triadaCompleta === true && vcase(e.icode, 'victim_case_status') === 'in'
            && g6c.json.answers.HV_3 === '19:30',
            `HV corregido→${r6c.status} triadaCompleta=${r6c.json.triadaCompleta}; estado=${vcase(e.icode, 'victim_case_status')}`);

        // C7 · Valores vacíos de fecha / hora / lista ya no dan 500 y no borran lo guardado.
        const zoneBefore = form2(a.icode, 'victim_case_form2_facts_zone');
        const r7 = await save(a.icode, 'HV', { HV_2: '', HV_3: '', HV_6: '' });
        check('C7', r7.status === 200 && form2(a.icode, 'victim_case_form2_facts_date') === '2026-10-30'
            && form2(a.icode, 'victim_case_form2_facts_start_time') === '19:30:00' && form2(a.icode, 'victim_case_form2_facts_zone') === zoneBefore,
            `HTTP ${r7.status}; fecha=${form2(a.icode, 'victim_case_form2_facts_date')} hora=${form2(a.icode, 'victim_case_form2_facts_start_time')} zona sin cambio=${form2(a.icode, 'victim_case_form2_facts_zone') === zoneBefore}`);

        // C8 · Sin respuesta ≠ valor marcador: un borrador nuevo no devuelve la fecha 1900-01-01 como respuesta.
        const g8 = await getCase(d.icode);
        check('C8', g8.json.answers.HV_2 === undefined && g8.json.answers.DP_1 === undefined,
            `borrador sin responder: HV_2=${g8.json.answers.HV_2} DP_1=${g8.json.answers.DP_1} (deben ser undefined)`);

        // C9 · Departamento borrado → NULL en la BD y no se devuelve.
        const r9 = await save(a.icode, 'HV', { HV_4: '' });
        const g9 = await getCase(a.icode);
        check('C9', r9.status === 200 && form2(a.icode, 'victim_case_form2_facts_department') === '<NULL>' && g9.json.answers.HV_4 === undefined,
            `HTTP ${r9.status}; BD depto=${form2(a.icode, 'victim_case_form2_facts_department')}; GET HV_4=${g9.json.answers.HV_4}`);

        // C10 · El fan-out del Perfil sigue funcionando después de mover el guardado a una transacción.
        const docShared = rnd();
        const f1 = await draft(docShared);
        const f2 = await draft(docShared);
        const r10 = await save(f1.icode, 'DP', { DP_29: 'Calle del fan-out 123', DP_27: 'Bolívar' });
        check('C10', r10.status === 200 && form2(f2.icode, 'victim_case_form2_residence_address') === 'Calle del fan-out 123'
            && form2(f2.icode, 'victim_case_form2_residence_department') === 'Bolívar',
            `HTTP ${r10.status}; hermano: dirección=${form2(f2.icode, 'victim_case_form2_residence_address')} depto=${form2(f2.icode, 'victim_case_form2_residence_department')}`);

        // C11 · Logs: cada 400 provocado deja su línea [RegistroV3] SaveSection rechazado en el log de la app.
        if (process.env.APP_LOG && fs.existsSync(process.env.APP_LOG)) {
            const rej = logCount('[RegistroV3] SaveSection rechazado') - log0.rej;
            const fail = logCount('[RegistroV3] SaveSection falló') - log0.fail;
            check('C11', rej === rejected400 && fail === 0, `log de esta corrida: ${rej} «rechazado» (400 provocados: ${rejected400}); ${fail} «falló» (500) — debe ser 0`);
        }

        // ═══ FASE 2 — FRONTEND (solo lo que ve la persona) ══════════════════════════════
        const gotoCase = async (pg, doc) => {
            await openRegistroV3(pg);
            await pg.locator('input[placeholder="Número de documento"]').fill(doc);
            await pg.getByRole('button', { name: 'Buscar' }).click();
            await pg.locator('tr', { hasText: `cc-${doc}` }).getByRole('button', { name: 'Abrir' }).click();
            await pg.getByRole('link', { name: 'Hechos de violencia' }).waitFor();
        };
        const goHV = async (pg) => { await pg.getByRole('link', { name: 'Hechos de violencia' }).click(); };
        const pg = page; // misma sesión: un segundo login del mismo usuario invalida la anterior (401 «sesión expirada»)
        const ERR = '#registro-v3-error-guardado';
        const LBL_DIR = '¿Cuál es la dirección donde ocurrieron los hechos (o nombre de la finca si es rural)?';

        // Se vuelve a dejar el caso 'a' con sus datos de hechos completos para la parte de pantalla.
        const rr = await save(a.icode, 'HV', p.HV);
        const gr = await getCase(a.icode);
        const grAns = (gr.json && gr.json.answers) || {};
        check('U0', rr.status === 200 && grAns.HV_4 === 'Atlántico' && ((gr.json && gr.json.answeredKeys) || []).includes('HV_4'),
            `preparación: HTTP ${rr.status} ${rr.text.slice(0, 120)}; GET HTTP ${gr.status} ${gr.status === 200 ? '' : gr.text.slice(0, 160)} HV_4=${grAns.HV_4}`);

        // U1 · Al reabrir, la pantalla muestra lo guardado (antes la fecha salía en blanco).
        await gotoCase(pg, a.doc);
        await goHV(pg);
        const uFecha = await question(pg, '¿En qué fecha ocurrió el hecho?').locator('input').inputValue();
        const uHora = await question(pg, '¿A qué hora ocurrió el hecho?').locator('input').inputValue();
        const uDepto = await question(pg, '¿En qué departamento ocurrieron los hechos?').locator('select').inputValue();
        const uMuni = await question(pg, '¿En qué municipio ocurrieron los hechos?').locator('input').inputValue();
        check('U1', uFecha === '2026-10-30' && uHora === '19:30' && uDepto === 'Atlántico' && uMuni === 'Barranquilla',
            `fecha="${uFecha}" hora="${uHora}" depto="${uDepto}" municipio="${uMuni}"`);

        // U2 · Dato inválido (dirección de 151 caracteres): aviso en rojo que nombra la pregunta; no se guarda; al corregir sí.
        await question(pg, LBL_DIR).locator('input').fill('x'.repeat(151));
        await pg.getByRole('button', { name: 'Guardar sección' }).click();
        await pg.locator(ERR).waitFor({ timeout: 8000 });
        const msg2 = await pg.locator(ERR).innerText();
        const dirDb = form2(a.icode, 'victim_case_form2_facts_address');
        check('U2a', msg2.includes('No se guardó la sección') && msg2.includes(LBL_DIR) && msg2.includes('150') && dirDb === 'Calle 45 # 12-34',
            `aviso: «${msg2.slice(0, 190)}»; BD sigue con «${dirDb}»`);
        await question(pg, LBL_DIR).locator('input').fill('Calle corregida 99');
        await pg.getByRole('button', { name: 'Guardar sección' }).click();
        await pg.locator('text=Sección guardada.').waitFor({ timeout: 8000 });
        check('U2b', (await pg.locator(ERR).count()) === 0 && form2(a.icode, 'victim_case_form2_facts_address') === 'Calle corregida 99',
            `al corregir: aviso en rojo oculto=${(await pg.locator(ERR).count()) === 0}; BD=${form2(a.icode, 'victim_case_form2_facts_address')}`);

        // U3 · Falla del servidor (500): avisa y NO dice «Sección guardada.».
        await pg.route('**/api/v1/victim-cases-v3/*/seccion/HV', (route) => route.fulfill({ status: 500, contentType: 'application/json', body: '{"error":"error al guardar la sección"}' }));
        await pg.getByRole('button', { name: 'Guardar sección' }).click();
        await pg.locator(ERR).waitFor({ timeout: 8000 });
        const msg3 = await pg.locator(ERR).innerText();
        check('U3', msg3.includes('error del servidor') && (await pg.locator('text=Sección guardada.').count()) === 0, `aviso: «${msg3.slice(0, 160)}»`);

        // U4 · Sin conexión: avisa.
        await pg.unroute('**/api/v1/victim-cases-v3/*/seccion/HV');
        await pg.route('**/api/v1/victim-cases-v3/*/seccion/HV', (route) => route.abort('failed'));
        await pg.getByRole('button', { name: 'Guardar sección' }).click();
        await pg.waitForFunction((sel) => { const el = document.querySelector(sel); return el && /no hay conexión/.test(el.innerText); }, ERR, { timeout: 8000 });
        check('U4', true, 'sin conexión: «' + (await pg.locator(ERR).innerText()).slice(0, 120) + '»');

        // U5 · «Finalizar llamada» con el guardado fallando: avisa y deja elegir (cancelar = seguir editando; aceptar = salir).
        let dlg = '';
        pg.once('dialog', (dl) => { dlg = dl.message(); dl.dismiss(); });
        await pg.getByRole('button', { name: 'Finalizar llamada (guardar y salir)' }).click();
        await pg.waitForFunction(() => document.body.innerText.includes('Hechos de violencia'));
        const sigueEnForm = await pg.getByRole('link', { name: 'Hechos de violencia' }).count();
        pg.once('dialog', (dl) => dl.accept());
        await pg.getByRole('button', { name: 'Finalizar llamada (guardar y salir)' }).click();
        await pg.getByRole('button', { name: '+ Nuevo registro' }).waitFor({ timeout: 8000 });
        check('U5', /no hay conexión/.test(dlg) && sigueEnForm === 1, `diálogo: «${dlg.replace(/\n+/g, ' ').slice(0, 150)}»; al cancelar sigue en el formulario=${sigueEnForm === 1}; al aceptar vuelve al Cabezote`);
        await pg.unroute('**/api/v1/victim-cases-v3/*/seccion/HV');

        // U6 · No abre el Caso (500): el Cabezote lo dice.
        await pg.route(`**/api/v1/victim-cases-v3/${a.icode}`, (route) => route.request().method() === 'GET'
            ? route.fulfill({ status: 500, contentType: 'application/json', body: '{"error":"error al cargar el caso"}' }) : route.continue());
        await pg.locator('input[placeholder="Número de documento"]').fill(a.doc);
        await pg.getByRole('button', { name: 'Buscar' }).click();
        await pg.locator('tr', { hasText: `cc-${a.doc}` }).getByRole('button', { name: 'Abrir' }).click();
        await pg.locator('.alert-danger', { hasText: 'No se pudo abrir el registro' }).waitFor({ timeout: 8000 });
        check('U6', (await pg.getByRole('button', { name: '+ Nuevo registro' }).count()) === 1, 'el Cabezote muestra «No se pudo abrir el registro…» y no cambia de pantalla');
        await op.context.close();
    } catch (e) {
        check('ERROR', false, e && e.stack ? e.stack.split('\n').slice(0, 4).join(' | ') : String(e));
    } finally {
        try { deleteCases(created); } catch (e) { console.log('(limpieza falló: ' + e.message + ')'); }
        await browser.close();
    }
    finish('REG-11 · guardado de secciones');
})();
