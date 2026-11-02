#!/usr/bin/env node
/**
 * WP-3 · REG-35, REG-13, REG-37, REG-21 (pantalla), REG-05 y REG-36 (verificación) — guardado, bloqueo y cierre de
 * Registro V3.
 *
 *  REG-35  Guardar GS o TZ no vuelve a poner los candados (SaveSection devuelve la tríada real para cualquier módulo).
 *  REG-13  Solo el Tamizaje queda bloqueado; GS se puede abrir desde el Borrador y sigue siendo obligatorio al finalizar.
 *  REG-37  «Finalizar» guarda GS antes de validar; si falta algo dice CUÁLES (camposFaltantes) y no cierra.
 *  REG-21  Un 401 al guardar/finalizar (con o sin cuerpo) avisa con claridad y no pierde lo diligenciado.
 *  REG-05  El recordatorio de obligatorias faltantes sale en rojo con contraste suficiente.
 *  REG-36  Verificación de lo que cubrió REG-11 (WP-2): un 500 ya no dice «Sección guardada.».
 *
 * Orden (qa-process.md): 1) backend por API + verificación en la BD; 2) pantalla, solo lo que el backend no cubre.
 * Uso (app corriendo con ./run_dev.sh):  node scripts/e2e/wp3-guardado-cierre.js
 * Casos: DocsMD/Screens/Registro V3/QA/E12-E15-triada-guardado-y-cierre-wp3-qa-cases.md
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
        const sectionBtn = () => page.getByRole('button', { name: 'Guardar sección' });
        const goTab = (name) => page.getByRole('link', { name }).click();

        // ═══ FASE 1 — BACKEND (API + BD) ═══════════════════════════════════════════════════════
        let d = await L.newDraft(page, created);
        let A = L.finalizeAnswers(catalog, d.doc);
        let r = await L.save(page, d.icode, 'GS', A.GS);
        check('G1', r.status === 200 && r.json.triadaCompleta === false && L.status(d.icode) === 'bo',
            `GS con la tríada incompleta -> ${r.status}, triadaCompleta=${r.json && r.json.triadaCompleta}, estado=${L.status(d.icode)}`);
        r = await L.save(page, d.icode, 'CA', A.CA);
        check('G5', r.json.triadaCompleta === false && L.status(d.icode) === 'bo', 'CA no aporta a la tríada: false y sigue en Borrador');
        for (const m of ['CI', 'DP', 'HV']) r = await L.save(page, d.icode, m, A[m]);
        check('G4a', r.json.triadaCompleta === false && L.status(d.icode) === 'bo', 'CI, DP y HV guardados: la tríada aún no está completa');
        r = await L.save(page, d.icode, 'AG', A.AG);
        check('G4', r.json.triadaCompleta === true && L.status(d.icode) === 'in', `AG completa la tríada -> true y estado ${L.status(d.icode)} (regresión: bo -> in)`);
        r = await L.save(page, d.icode, 'GS', A.GS);
        check('G2', r.status === 200 && r.json.triadaCompleta === true && L.status(d.icode) === 'in',
            `GS con la tríada completa -> triadaCompleta=${r.json && r.json.triadaCompleta}, estado sigue ${L.status(d.icode)} (REG-35)`);
        r = await L.save(page, d.icode, 'TZ', { TZ_1: 'false' });
        check('G3', r.status === 200 && r.json.triadaCompleta === true, `/seccion/TZ con la tríada completa -> triadaCompleta=${r.json && r.json.triadaCompleta}`);

        // Finalizar: GS nunca guardado -> 400 con módulos y campos (REG-13, REG-37)
        d = await L.newDraft(page, created); A = L.finalizeAnswers(catalog, d.doc);
        await L.readyExceptGS(page, d.icode, A);
        r = await api(page, 'POST', `/api/v1/victim-cases-v3/${d.icode}/finalizar`);
        const campos = (r.json && r.json.camposFaltantes) || [];
        const mods = (r.json && r.json.modulosIncompletos) || [];
        const gsEsperadas = ['GS_1', 'GS_14', 'GS_16', 'GS_17', 'GS_18', 'GS_2', 'GS_4', 'GS_5'];
        check('N4', r.status === 400 && mods.join() === 'GS', `finalizar sin GS -> ${r.status}, modulosIncompletos=${mods.join()} (GS sigue siendo obligatorio)`);
        check('C6', JSON.stringify(campos) === JSON.stringify(gsEsperadas) && campos.every((k) => mods.includes(k.split('_')[0])),
            `camposFaltantes ordenado y consistente con los módulos -> ${campos.join(',')}`);
        r = await L.save(page, d.icode, 'GS', A.GS);
        r = await api(page, 'POST', `/api/v1/victim-cases-v3/${d.icode}/finalizar`);
        check('C0', r.status === 200 && L.status(d.icode) === 'te', `con GS guardado, finalizar -> ${r.status}, estado=${L.status(d.icode)}`);

        // Sesión / permisos (REG-21, REG-36): lo que el servidor responde
        const anon = await browser.newContext();
        const anonResp = await anon.request.post(`${BASE_URL}/api/v1/victim-cases-v3/${d.icode}/seccion/CI`, { data: { answers: {} } });
        const anonBody = await anonResp.text();
        check('G6', anonResp.status() === 401 && /no autenticado/.test(anonBody), `sin cookie de sesión -> ${anonResp.status()} ${anonBody}`);
        await anon.close();
        const ro = await loginAs(browser, 'ro');
        const roR = await api(ro.page, 'POST', `/api/v1/victim-cases-v3/${d.icode}/seccion/CI`, { answers: {} });
        check('G7', roR.status === 401 && roR.text === '', `rol sin permiso (ro) -> ${roR.status} con cuerpo ${JSON.stringify(roR.text)} (401 sin cuerpo)`);
        await ro.context.close();

        // ═══ FASE 2 — PANTALLA ═════════════════════════════════════════════════════════════════
        // U1 · REG-13: GS sin candado desde el Borrador; TZ con candado y aviso solo de Tamizaje
        await openRegistroV3(page);
        const icodeUI = await newDraftUI(page); created.push(icodeUI);
        const menuText = async (name) => (await page.getByRole('link', { name }).innerText()).replace(/\s+/g, ' ');
        check('U1a', !(await menuText('Gestión Salvia')).includes('🔒'), `GS sin candado en un Borrador: "${await menuText('Gestión Salvia')}"`);
        check('U1b', (await menuText('Tamizaje de riesgo')).includes('🔒'), 'TZ conserva el candado');
        await goTab('Gestión Salvia');
        await page.getByRole('button', { name: 'Finalizar registro y enrutar a Seguimiento' }).waitFor();
        const bodyGS = await page.locator('body').innerText();
        check('U1c', pageErrors.length === 0 && !/undefined/.test(bodyGS), `GS abierto con CI vacío sin errores de consola ni "undefined" (errores: ${pageErrors.length})`);
        await goTab('Tamizaje de riesgo');
        const aviso = await page.locator('.alert-warning').innerText();
        check('U1d', /Tamizaje/.test(aviso) && !/Gestión Salvia/.test(aviso), `aviso del TZ bloqueado solo habla de Tamizaje: "${aviso.trim()}"`);

        // U2 · REG-35 / REG-05 / REG-36 / REG-21 sobre un Caso con la tríada completa
        d = await L.newDraft(page, created); A = L.finalizeAnswers(catalog, d.doc);
        await L.completeTriad(page, d.icode, A);
        await openCaseUI(page, d.doc);
        await goTab('Gestión Salvia');
        // REG-05: falta GS_4 -> recordatorio en rojo
        await L.setAnswers(page, { GS_1: 'false', GS_2: A.GS.GS_2, GS_4: '' });
        await sectionBtn().click();
        const msg = page.locator('span.ms-2', { hasText: 'Faltan obligatorias de 1er contacto' });
        await msg.waitFor();
        const css = await msg.evaluate((el) => { const s = getComputedStyle(el); return { color: s.color, bg: getComputedStyle(document.body).backgroundColor }; });
        const cr = L.contrast(css.color, 'rgb(255, 255, 255)');
        check('U2a', css.color === 'rgb(176, 42, 55)' && cr >= 4.5, `recordatorio en rojo ${css.color}, contraste ${cr.toFixed(2)}:1 (mínimo 4,5)`);
        // éxito -> sin rojo, sin avanzar (GS es el último módulo) y sin volver a poner candados (REG-35)
        await L.setAnswers(page, { GS_4: A.GS.GS_4 });
        await sectionBtn().click();
        const ok = page.locator('span.ms-2', { hasText: 'Sección guardada.' });
        await ok.waitFor();
        const okColor = await ok.evaluate((el) => getComputedStyle(el).color);
        check('U2b', okColor !== 'rgb(176, 42, 55)', `"Sección guardada." conserva su estilo (color ${okColor})`);
        check('U2c', (await L.vm(page, 'p.triadaCompleta')) === true && !(await menuText('Tamizaje de riesgo')).includes('🔒') && !(await menuText('Gestión Salvia')).includes('🔒')
            && (await page.locator('.alert-success', { hasText: 'Iniciado' }).count()) === 1, 'tras guardar GS: triadaCompleta=true, sin candados y el aviso verde permanece (REG-35)');
        check('U2d', (await L.vm(page, 'p.moduloActivo')) === 'GS', 'en el último módulo (GS) guardar no cambia de módulo');

        // REG-36 (verificación de REG-11): 500, 401 con cuerpo, 401 sin cuerpo
        const errBox = page.locator('#registro-v3-error-guardado');
        const keep = A.GS.GS_4;
        await page.route('**/seccion/GS', (route) => route.fulfill({ status: 500, contentType: 'application/json', body: '{"error":"error al guardar la sección"}' }));
        await sectionBtn().click();
        await errBox.waitFor();
        let t = await errBox.innerText();
        const bodyAfter500 = await page.locator('body').innerText();
        check('U3a', /error del servidor/.test(t) && !/Sección guardada\./.test(bodyAfter500), `500 -> aviso rojo y nunca "Sección guardada." ("${t.slice(0, 60)}…")`);
        await page.unroute('**/seccion/GS');
        await page.route('**/seccion/GS', (route) => route.fulfill({ status: 401, contentType: 'application/json', body: '{"error":"no autenticado"}' }));
        await sectionBtn().click();
        await page.waitForFunction(() => /sesión venció/.test((document.getElementById('registro-v3-error-guardado') || {}).innerText || ''));
        t = await errBox.innerText();
        check('U3b', /sigue en esta pantalla/.test(t) && (await L.vm(page, 'p.answers.GS_4')) === keep, `401 "no autenticado" (REG-21) -> mensaje claro y lo escrito sigue en pantalla ("${t.slice(0, 70)}…")`);
        await page.unroute('**/seccion/GS');
        await page.route('**/seccion/GS', (route) => route.fulfill({ status: 401, body: '' }));
        await sectionBtn().click();
        await page.waitForFunction(() => /sesión venció/.test((document.getElementById('registro-v3-error-guardado') || {}).innerText || ''));
        check('U3c', pageErrors.length === 0, `401 sin cuerpo (permiso) -> mismo aviso y sin excepciones (errores de página: ${pageErrors.length})`);
        await page.unroute('**/seccion/GS');

        // U4 · REG-37 — Finalizar guarda GS antes de validar
        const finalizar = () => page.getByRole('button', { name: 'Finalizar registro y enrutar a Seguimiento' });
        const gsFull = A.GS;
        // U4a: GS § A incompleto -> mensaje rojo y CERO peticiones a /finalizar
        let d4 = await L.newDraft(page, created); let A4 = L.finalizeAnswers(catalog, d4.doc);
        await L.readyExceptGS(page, d4.icode, A4);
        await openCaseUI(page, d4.doc);
        await goTab('Gestión Salvia');
        let finReq = 0;
        const onReq = (rq) => { if (rq.url().endsWith('/finalizar')) finReq++; };
        page.on('request', onReq);
        await L.setAnswers(page, { ...gsFull, GS_4: '' });
        await finalizar().click();
        await page.locator('.alert-danger', { hasText: 'Faltan obligatorias' }).first().waitFor();
        check('U4a', finReq === 0, `GS_4 vacío -> el recordatorio sale junto a Finalizar y NO se llama a /finalizar (${finReq} peticiones)`);
        // U4b: § A completo pero GS_5.. sin responder -> 400 con camposFaltantes y etiquetas
        await L.setAnswers(page, { GS_1: 'false', GS_2: gsFull.GS_2, GS_4: gsFull.GS_4 });
        await L.unsetAnswers(page, ['GS_5', 'GS_14', 'GS_16', 'GS_17', 'GS_18']);
        await finalizar().click();
        await page.locator('.alert-danger', { hasText: 'Gestión Salvia, falta' }).first().waitFor();
        const eb = await page.locator('.alert-danger', { hasText: 'Gestión Salvia, falta' }).first().innerText();
        check('U4b', /Requiere remisión/.test(eb) && /alerta inminente/i.test(eb) && L.status(d4.icode) === 'in', `faltan GS_5.. -> el aviso nombra los campos y el Caso sigue Iniciado ("${eb.slice(0, 110)}…")`);
        // U4c: 401 al finalizar (con y sin cuerpo): aviso claro, datos intactos
        await L.setAnswers(page, gsFull);
        await page.route('**/finalizar', (route) => route.fulfill({ status: 401, contentType: 'application/json', body: '{"error":"no autenticado"}' }));
        await finalizar().click();
        await page.locator('.alert-danger', { hasText: 'sesión venció' }).first().waitFor();
        check('U4c', (await L.vm(page, "p.etapa")) === 'formulario' && (await L.vm(page, 'p.answers.GS_4')) === gsFull.GS_4, '401 "no autenticado" al finalizar (REG-21) -> aviso claro, sigue en el formulario y los datos están intactos');
        await page.unroute('**/finalizar');
        await page.route('**/finalizar', (route) => route.fulfill({ status: 401, body: '' }));
        await finalizar().click();
        await page.waitForFunction(() => /sesión venció/.test((document.getElementById('app') || {}).innerText || ''));
        check('U4d', pageErrors.length === 0, `401 sin cuerpo al finalizar -> sin excepciones (errores de página: ${pageErrors.length})`);
        await page.unroute('**/finalizar');
        // U4e: sin pulsar «Guardar sección», Finalizar guarda GS y cierra; doble clic = una sola petición
        finReq = 0;
        await finalizar().dblclick();
        await page.getByText('Registro finalizado.').waitFor();
        const gs4 = sql(`SELECT f.victim_case_form2_saliva_management_explanation FROM salvia.victim_case_form2 f JOIN salvia.victim_case c ON c.victim_case_id = f.victim_case_form2_victim_case WHERE c.victim_case_i_code = '${d4.icode}'`);
        check('U4e', L.status(d4.icode) === 'te' && gs4 === gsFull.GS_4, `Finalizar sin "Guardar sección" previo: estado=${L.status(d4.icode)}, GS_4 guardado="${gs4}"`);
        check('U4f', finReq === 1, `doble clic -> ${finReq} petición a /finalizar`);
        check('U4g', noV3(await page.locator('body').innerText()) && noV3(await page.title()), 'el cierre y el título de la pestaña no muestran "V3"');
        page.off('request', onReq);
    } catch (e) {
        check('ERROR', false, e && e.stack ? e.stack.split('\n').slice(0, 4).join(' | ') : String(e));
    } finally {
        try { deleteCases(created); } catch (e) { console.log('limpieza:', e.message); }
        await browser.close();
    }
    finish('WP-3 · guardado y cierre');
})();
