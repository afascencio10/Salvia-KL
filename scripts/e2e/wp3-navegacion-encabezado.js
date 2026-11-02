#!/usr/bin/env node
/**
 * WP-3 · REG-06, REG-04 (a y b), REG-03, REG-01 y REG-02 — textos, estados, encabezado, navegación y avance de
 * Registro V3.
 *
 *  REG-06  El usuario final no ve «V3» (título, avisos y título de la pestaña).
 *  REG-04  El estado del Caso se ve con su nombre (Borrador, Iniciado…) en el Cabezote y en el encabezado.
 *  REG-03  El encabezado muestra el nombre y el documento reales (nunca «Pendiente Pendiente»).
 *  REG-01  Menú vertical construido con modulos[], flechas Anterior/Siguiente y avance al guardar.
 *  REG-02  Porcentaje de avance por módulo (menú) y global (encabezado), en vivo.
 *
 * Orden (qa-process.md): 1) API (estados del Cabezote) 2) pantalla, solo lo que la API no cubre.
 * Uso (app corriendo con ./run_dev.sh):  node scripts/e2e/wp3-navegacion-encabezado.js
 * Casos: DocsMD/Screens/Registro V3/QA/E10-E22-E23-navegacion-wp3-qa-cases.md y E04-encabezado-estados-wp3-qa-cases.md
 */
const { chromium } = require('playwright');
const L = require('./wp3-lib');
const { sql, deleteCases, check, finish, loginAs, api, loadCatalog, openRegistroV3, newDraftUI, openCaseUI, BASE_URL } = L;

const created = [];
const noV3 = (t) => !/\bV3\b/.test(t);

(async () => {
    const browser = await chromium.launch();
    try {
        const op = await loginAs(browser, 'op');
        const page = op.page;
        const pageErrors = [];
        page.on('pageerror', (e) => pageErrors.push(String(e)));
        const catalog = await loadCatalog(page);
        const saveBtn = () => page.getByRole('button', { name: 'Guardar sección' });
        const goTab = (name) => page.getByRole('link', { name }).click();
        const active = async () => (await page.locator('.v3-menu .nav-link.active').innerText()).replace(/\s+/g, ' ').replace(/ \d+ %$| —$/, '').trim();
        const headerText = async () => (await page.locator('.v3-form > .mb-2').first().innerText()).replace(/\s+/g, ' ');
        const saved = () => page.locator('.v3-nav').getByText('Sección guardada.').waitFor();
        const hdrWait = async (re) => {
            await page.waitForFunction((src) => new RegExp(src).test(document.querySelector('.v3-form > .mb-2').innerText.replace(/\s+/g, ' ')), re.source, { timeout: 10000 });
            return headerText();
        };

        // ═══ REG-06 · textos sin «V3» ══════════════════════════════════════════════════════════
        await page.goto(`${BASE_URL}/salvia/registro-v3`);
        await page.getByRole('button', { name: 'Acepto y continuar' }).waitFor();
        let ok = noV3(await page.locator('body').innerText()) && (await page.title()) === 'Registro de víctima' && (await page.locator('h2').first().innerText()) === 'Registro de víctima';
        check('T9a', ok, `etapa de autorización: título "${await page.title()}", h2 "${await page.locator('h2').first().innerText()}"`);
        await page.getByRole('button', { name: 'Acepto y continuar' }).click();
        await page.getByRole('button', { name: 'Víctima - Seguimiento' }).click();
        const avisoTipo = await page.locator('.alert-info').innerText();
        check('T9b', noV3(avisoTipo) && /Registro de víctima/.test(avisoTipo), `aviso de tipo de contacto: "${avisoTipo.slice(0, 80)}…"`);
        const rutaOk = (await api(page, 'GET', '/salvia/registro-v3')).status === 200;
        check('T9c', rutaOk, 'la ruta /salvia/registro-v3 sigue respondiendo 200');

        // ═══ REG-04a · estados en el Cabezote (API crea, SQL fija estados, pantalla verifica) ═══
        const docX = String(700000000 + Math.floor(Math.random() * 99999999));
        const ics = [];
        for (let i = 0; i < 5; i++) ics.push((await L.newDraft(page, created, docX)).icode);
        ['in', 'te', 'cd', 'zz'].forEach((st, i) => sql(`UPDATE salvia.victim_case SET victim_case_status = '${st}' WHERE victim_case_i_code = '${ics[i + 1]}'`));
        await openRegistroV3(page);
        await page.locator('input[placeholder="Número de documento"]').fill(docX);
        await page.getByRole('button', { name: 'Buscar' }).click();
        await page.locator('table tbody tr').first().waitFor();
        const estados = (await page.locator('table tbody tr td:nth-child(4)').allInnerTexts()).map((s) => s.trim()).sort();
        check('T8a', JSON.stringify(estados) === JSON.stringify(['Borrador', 'Cerrado', 'Desconocido', 'Iniciado', 'Terminado']),
            `columna Estado del Cabezote: ${estados.join(', ')} (sin códigos crudos)`);

        // ═══ REG-03 / REG-04b / REG-01 · recorrido de un Caso nuevo por la pantalla ═══════════
        await page.locator('input[placeholder="Número de documento"]').fill('');
        const icode = await newDraftUI(page); created.push(icode);
        const doc = String(600000000 + Math.floor(Math.random() * 99999999));
        const A = L.finalizeAnswers(catalog, doc);
        let h = await headerText();
        check('T7a', /Persona sin identificar/.test(h) && /Documento sin registrar/.test(h) && /Borrador/.test(h) && !/Pendiente/.test(h) && !/\b(bo)\b/.test(h),
            `Borrador sin datos: "${h}"`);
        check('T10a', (await active()) === 'Contacto e identificación', 'el Caso se abre en CI con su ítem activo en el menú');

        // Menú construido con modulos[] (contrato con WP-4) y orden
        const labels = (await page.locator('.v3-menu .nav-link').allInnerTexts()).map((s) => s.replace(/\s+/g, ' ').replace(/ \d+ %$| —$/, '').replace(' 🔒', '').trim());
        const modl = await L.vm(page, "p.modulos.map(function (m) { return m.label; })");
        check('T10b', JSON.stringify(labels) === JSON.stringify(modl) && labels.length === 7, `menú = modulos[] (${labels.length} ítems, mismo orden)`);

        // REG-02: porcentajes iniciales
        const pct0 = (await page.locator('.v3-menu .nav-link').allInnerTexts()).map((s) => ((s.replace(/\s+/g, ' ').trim().match(/(\d+ %|—)$/) || [''])[0]));
        check('T11a', pct0.join('|') === '0 %|0 %|—|0 %|0 %|—|0 %',
            `porcentajes de un Borrador vacío (CI, DP, CA, HV, AG, TZ, GS): ${pct0.join(' | ')}`);
        check('T11b', /Avance 0 %/.test(h), 'avance global 0 % en el encabezado');
        // REG-02: un condicional que aparece sube el total (CI_7 = Sí muestra CI_8)
        await L.setAnswers(page, { ...A.CI, CI_7: 'true' });
        const ci = await L.vm(page, "p.etiquetaPorcentaje('CI')");
        check('T11c', ci === '89 %', `CI_7=Sí muestra CI_8 y el avance de CI baja a ${ci} (8 de 9)`);
        await saveBtn().click();
        const faltantes = page.locator('span.ms-2', { hasText: 'Faltan obligatorias' });
        await faltantes.waitFor();
        const nFalta = (await faltantes.innerText()).replace('Faltan obligatorias de 1er contacto:', '').split(';').filter((x) => x.trim()).length;
        check('T11d', nFalta === 1, `el recordatorio lista ${nFalta} faltante(s), igual a total - respondidas (9 - 8)`);
        // CI completo: guardar -> avanza a DP (REG-01), el encabezado se actualiza (REG-03)
        let posts = 0;
        const onReq = (rq) => { if (rq.method() === 'POST' && /\/seccion\//.test(rq.url())) posts++; };
        await L.setAnswers(page, { CI_7: 'false' });
        await L.unsetAnswers(page, ['CI_8']);
        check('T11e', (await L.vm(page, "p.etiquetaPorcentaje('CI')")) === '100 %', 'CI completo -> 100 %');
        await saveBtn().click();
        await saved();
        h = await hdrWait(/Ana Perez/);
        check('T7b', /Sra\. Ana Perez \(Anita\)/.test(h) && h.includes(`cc-${doc}`) && !/Pendiente/.test(h), `tras guardar CI el encabezado muestra nombre y documento: "${h}"`);
        check('T10c', (await active()) === 'Datos personales', 'guardar CI con éxito avanzó a Datos personales');
        check('T8b1', /Borrador/.test(h), 'la insignia sigue en Borrador con la tríada incompleta');

        // Navegación por flechas: no guarda, saltos sobre lo bloqueado y extremos
        page.on('request', onReq);
        posts = 0;
        await L.setAnswers(page, { CI_2: 'Ana Perez' });
        await goTab('Características del agresor');
        await page.getByRole('button', { name: 'Siguiente →' }).click();
        check('T10d', (await active()) === 'Gestión Salvia', 'con TZ bloqueado, «Siguiente» desde AG salta a Gestión Salvia');
        check('T10e', await page.getByRole('button', { name: 'Siguiente →' }).isDisabled(), '«Siguiente» deshabilitado en el último módulo');
        await page.getByRole('button', { name: '← Anterior' }).click();
        check('T10f', (await active()) === 'Características del agresor', '«Anterior» vuelve al módulo previo no bloqueado');
        await goTab('Contacto e identificación');
        check('T10g', await page.getByRole('button', { name: '← Anterior' }).isDisabled(), '«Anterior» deshabilitado en el primer módulo');
        check('T10h', posts === 0 && (await L.vm(page, "p.answers.CI_2")) === 'Ana Perez', `navegar no guarda (POST a /seccion: ${posts}) y lo escrito se conserva`);
        page.off('request', onReq);

        // Completar la tríada por la pantalla: DP -> (CA) -> HV -> AG; al guardar AG pasa a Tamizaje y el estado cambia
        await goTab('Datos personales');
        await L.setAnswers(page, A.DP); await saveBtn().click(); await saved();
        check('T10i', (await active()) === 'Contacto de apoyo', 'guardar DP avanzó a Contacto de apoyo');
        await page.getByRole('button', { name: 'Siguiente →' }).click();
        await L.setAnswers(page, A.HV); await saveBtn().click(); await saved();
        check('T10j', (await active()) === 'Características del agresor', 'guardar HV avanzó a AG');
        await L.setAnswers(page, { ...A.AG, AG_5: 'false' }); await saveBtn().click(); await saved();
        h = await hdrWait(/Iniciado/);
        check('T10k', (await active()) === 'Tamizaje de riesgo', 'al completar la tríada, guardar AG avanza a Tamizaje (ya desbloqueado)');
        check('T8b2', /Iniciado/.test(h) && !/Borrador/.test(h), `la insignia pasó a Iniciado sin recargar: "${h}"`);
        check('T11f', (await L.vm(page, "p.porcentajeGlobal()")) > 0 && pageErrors.length === 0, `avance global ${await L.vm(page, 'p.porcentajeGlobal()')} % y sin errores de consola (${pageErrors.length})`);

        // Contrato con WP-4: un módulo agregado a modulos[] aparece solo en el menú
        await page.evaluate(() => { document.getElementById('app')._vnode.component.proxy.modulos.push({ id: 'ZZ', label: 'Módulo de prueba', campos: [] }); });
        check('T10l', (await page.locator('.v3-menu .nav-link').count()) === 8 && (await page.getByRole('link', { name: 'Módulo de prueba' }).count()) === 1, 'un módulo agregado a modulos[] aparece solo en el menú');
        await page.evaluate(() => { document.getElementById('app')._vnode.component.proxy.modulos.pop(); });

        // Reabrir desde el Cabezote muestra el mismo encabezado (REG-03)
        await openCaseUI(page, doc);
        h = await headerText();
        check('T7c', /Sra\. Ana Perez/.test(h) && h.includes(`cc-${doc}`) && /Iniciado/.test(h), `al reabrir el Caso: "${h}"`);

        // Maquetación: menú a la izquierda (escritorio) y apilado (móvil)
        const box = async (sel) => page.locator(sel).first().boundingBox();
        await page.setViewportSize({ width: 1280, height: 900 });
        let m = await box('.v3-menu'), n = await box('.v3-nav');
        check('T10m', m.x + m.width <= n.x + 1, `escritorio: menú a la izquierda (menú x=${Math.round(m.x)}..${Math.round(m.x + m.width)}, contenido x=${Math.round(n.x)})`);
        await page.setViewportSize({ width: 375, height: 800 });
        m = await box('.v3-menu'); n = await box('.v3-nav');
        check('T10n', m.y < n.y && Math.abs(m.x - n.x) < 8, `móvil: menú apilado arriba (menú y=${Math.round(m.y)}, contenido y=${Math.round(n.y)})`);
    } catch (e) {
        check('ERROR', false, e && e.stack ? e.stack.split('\n').slice(0, 4).join(' | ') : String(e));
    } finally {
        try { deleteCases(created); } catch (e) { console.log('limpieza:', e.message); }
        await browser.close();
    }
    finish('WP-3 · navegación, encabezado y avance');
})();
