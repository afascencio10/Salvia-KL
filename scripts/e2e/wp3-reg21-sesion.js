#!/usr/bin/env node
/**
 * WP-3 · REG-21 — «Finalizar registro» responde «no autenticado»: reproducción de las causas posibles y verificación
 * de lo que ve la persona.
 *
 * Qué mensaje significa qué (requireRegistroV3Access, victim_case_v3_controller.go):
 *   401 {"error":"no autenticado"}   -> la cookie del navegador NO trae `userData`
 *   401 {"error":"sesión expirada"}  -> la cookie sí lo trae, pero el servidor ya no tiene esa sesión
 *   401 sin cuerpo                   -> sesión válida, el rol no tiene el permiso manage_registro_v3
 *
 * Experimentos (cada uno con el formulario ya abierto, que es como lo vive el agente):
 *   S1  línea base: con sesión vigente «Finalizar» cierra el registro.
 *   S2  H1: se cierra sesión en OTRA pestaña del mismo navegador (la cookie es compartida) y se pulsa «Finalizar».
 *   S3  H1/H3: se borran las cookies del sitio (equivale a cookie vencida) y se pulsa «Finalizar».
 *   S7  se inicia sesión con el MISMO usuario en otro navegador (el login borra la sesión previa del usuario).
 * En S2, S3 y S7 se registra el código y el cuerpo de la primera petición que falla y lo que muestra la pantalla.
 *
 * Uso (app corriendo con ./run_dev.sh):  node scripts/e2e/wp3-reg21-sesion.js
 * Casos: DocsMD/Screens/Registro V3/QA/sesion-perdida-wp3-qa-cases.md
 */
const { chromium } = require('playwright');
const L = require('./wp3-lib');
const { deleteCases, check, finish, loginAs, api, loadCatalog, openCaseUI, BASE_URL } = L;

const created = [];

(async () => {
    const browser = await chromium.launch();
    try {
        const catalog = await loadCatalog((await loginAs(browser, 'op')).page);

        /** Abre un Caso listo para finalizar (GS escrito en pantalla, sin guardar) y devuelve la página y las respuestas. */
        async function formularioListo(ctxUser) {
            const op = ctxUser || (await loginAs(browser, 'op'));
            const page = op.page;
            const d = await L.newDraft(page, created);
            const A = L.finalizeAnswers(catalog, d.doc);
            await L.readyExceptGS(page, d.icode, A);
            await openCaseUI(page, d.doc);
            await page.getByRole('link', { name: 'Gestión Salvia' }).click();
            await L.setAnswers(page, A.GS);
            return { op, page, d, A };
        }
        const finalizar = (page) => page.getByRole('button', { name: 'Finalizar registro y enrutar a Seguimiento' });
        /** Pulsa Finalizar y devuelve la primera respuesta no exitosa { url, status, body } (o null) y el texto de error visible. */
        async function pulsarFinalizar(page) {
            const fallo = new Promise((resolve) => {
                const h = async (r) => {
                    if (/\/api\/v1\/victim-cases-v3\//.test(r.url()) && r.status() >= 400) {
                        page.off('response', h);
                        resolve({ url: r.url().replace(BASE_URL, ''), status: r.status(), body: (await r.text().catch(() => '')).slice(0, 80) });
                    }
                };
                page.on('response', h);
                setTimeout(() => resolve(null), 8000);
            });
            await finalizar(page).click();
            const f = await fallo;
            await page.waitForTimeout(500);
            const alertas = (await page.locator('.alert-danger').allInnerTexts()).map((s) => s.replace(/\s+/g, ' ').trim()).join(' | ');
            return { f, alertas };
        }

        // S1 · línea base
        {
            const { op, page, d } = await formularioListo();
            await finalizar(page).click();
            await page.getByText('Registro finalizado.').waitFor();
            check('S1', L.status(d.icode) === 'te', `sesión vigente: Finalizar cierra el registro (estado=${L.status(d.icode)})`);
            await op.context.close();
        }

        // S2 · cierre de sesión en otra pestaña del mismo navegador
        {
            const { op, page, d, A } = await formularioListo();
            const otra = await op.context.newPage();
            await otra.goto(`${BASE_URL}/salvia/registro-v3`);
            const lo = await api(otra, 'POST', '/seguridad/logout');
            const { f, alertas } = await pulsarFinalizar(page);
            console.log(`   S2 → primera respuesta que falla: ${JSON.stringify(f)} | pantalla: ${alertas.slice(0, 150)}`);
            check('S2', f && f.status === 401 && /no autenticado/.test(f.body) && /sesión venció/.test(alertas),
                `logout en otra pestaña (HTTP ${lo.status}): ${f && f.status} ${f && f.body} -> la pantalla avisa y conserva los datos (GS_4="${await L.vm(page, 'p.answers.GS_4')}")`);
            check('S2b', (await L.vm(page, 'p.answers.GS_4')) === A.GS.GS_4 && L.status(d.icode) === 'in', 'datos intactos en pantalla y el Caso sigue Iniciado');
            await op.context.close();
        }

        // S3 · cookies borradas (equivale a cookie vencida)
        {
            const { op, page, A } = await formularioListo();
            await op.context.clearCookies();
            const { f, alertas } = await pulsarFinalizar(page);
            console.log(`   S3 → primera respuesta que falla: ${JSON.stringify(f)} | pantalla: ${alertas.slice(0, 150)}`);
            check('S3', f && f.status === 401 && /no autenticado/.test(f.body) && /sesión venció/.test(alertas), `cookies borradas: ${f && f.status} ${f && f.body} -> aviso claro`);
            check('S3b', (await L.vm(page, 'p.answers.GS_4')) === A.GS.GS_4, 'lo diligenciado sigue en pantalla');
            await op.context.close();
        }

        // S7 · segundo inicio de sesión del mismo usuario en otro navegador
        {
            const { op, page, A } = await formularioListo();
            const otro = await loginAs(browser, 'op');            // misma cuenta, otro navegador: el login borra la sesión previa
            const { f, alertas } = await pulsarFinalizar(page);
            console.log(`   S7 → primera respuesta que falla: ${JSON.stringify(f)} | pantalla: ${alertas.slice(0, 150)}`);
            check('S7', f && f.status === 401 && /sesión expirada/.test(f.body) && /sesión venció/.test(alertas), `segundo login del mismo usuario: ${f && f.status} ${f && f.body} -> aviso claro (no es "no autenticado")`);
            check('S7b', (await L.vm(page, 'p.answers.GS_4')) === A.GS.GS_4, 'lo diligenciado sigue en pantalla');
            await otro.context.close(); await op.context.close();
        }
    } catch (e) {
        check('ERROR', false, e && e.stack ? e.stack.split('\n').slice(0, 4).join(' | ') : String(e));
    } finally {
        try { deleteCases(created); } catch (e) { console.log('limpieza:', e.message); }
        await browser.close();
    }
    finish('WP-3 · REG-21 sesión');
})();
