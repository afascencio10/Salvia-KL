#!/usr/bin/env node
/**
 * WP-2 · T3 · REG-10 — Los subtipos de violencia (HV_10) se ven con el texto correcto.
 *
 * Bug: etiqueta() buscaba traducciones en TODAS las categorías de ENUM_LABELS_OVERRIDE_V3 cuando el
 * campo tenía subtipoDe, así que un código repetido en otra categoría (tr, al, es, ah, ai, pr) tomaba
 * la etiqueta de esa otra (p. ej. «Trata de personas…» salía como «Transporte»).
 *
 * Orden (qa-process.md): 1) backend — las etiquetas que entrega /enums son las correctas y los códigos
 * de cada categoría son únicos; 2) frontend — lo que dibuja la pantalla es exactamente eso.
 *
 * Uso (con la app corriendo):  node cloudDevelopment/e2e/wp2-reg10-subtipos-violencia.js
 * Casos C1…C5: ver DocsMD/Screens/Registro V3/QA/E11-guardar-seccion-qa-cases.md (sección REG-10).
 */
const { chromium } = require('playwright');
const { deleteCases, check, finish, loginAs, api, openRegistroV3, newDraftUI, openModule, question } = require('./wp2-lib');

const TIPOS = ['fi', 'pl', 'po', 'ps', 're', 'se', 'vi']; // código del tipo (HV_9) = sufijo de la categoría de subtipos
const created = [];
const HV9 = '¿Qué tipo o tipos de violencia ha experimentado?';
const HV10 = '¿Qué subtipo o subtipos de violencia se presentaron?';
const subCat = (suf) => `victim_case_form2_subtype_violence_experienced_${suf}`;

(async () => {
    const browser = await chromium.launch();
    try {
        const op = await loginAs(browser, 'op');

        // ═══ FASE 1 — BACKEND ═══════════════════════════════════════════════════════════
        const cat = (await api(op.page, 'GET', '/api/v1/victim-cases-v3/enums')).json || {};
        const bySuf = Object.fromEntries(TIPOS.map((s) => [s, cat[subCat(s)] || []]));
        const empty = TIPOS.filter((s) => bySuf[s].length === 0);
        const slugs = TIPOS.flatMap((s) => bySuf[s].filter((o) => !o.label || o.label.startsWith('victim_case_form2_')).map((o) => o.code));
        check('C1', empty.length === 0 && slugs.length === 0,
            `/enums: ${TIPOS.map((s) => s + '=' + bySuf[s].length).join(' ')}; categorías vacías=[${empty}]; etiquetas sin traducir=[${slugs}]`);
        const dup = TIPOS.filter((s) => new Set(bySuf[s].map((o) => o.code)).size !== bySuf[s].length);
        check('C2', dup.length === 0, `los códigos de cada categoría de subtipos son únicos (repetidas: [${dup}])`);
        const tipos = cat['victim_case_form2_type_violence_experienced'] || [];
        check('C3', TIPOS.every((s) => tipos.some((o) => o.code === s)),
            `el catálogo de tipos (HV_9) trae los 7 códigos: ${tipos.map((o) => o.code).join(',')}`);

        // ═══ FASE 2 — FRONTEND ══════════════════════════════════════════════════════════
        await openRegistroV3(op.page);
        created.push(await newDraftUI(op.page));
        await openModule(op.page, 'Hechos de violencia');
        const hv9 = question(op.page, HV9);
        await hv9.waitFor({ timeout: 10000 });

        const subtypeLabels = async () => (await question(op.page, HV10).locator('.form-check label').allInnerTexts()).map((t) => t.trim());
        const norm = (arr) => [...arr].map((t) => t.trim()).sort((a, b) => a.localeCompare(b, 'es'));

        // C4 · un solo tipo: Psicológica (ps) → sus 6 subtipos, con los textos del catálogo.
        const psLabel = tipos.find((o) => o.code === 'ps').label.trim();
        await hv9.locator('.form-check', { hasText: psLabel }).locator('input').check();
        const psUi = await subtypeLabels();
        const psExpected = bySuf.ps.map((o) => o.label);
        check('C4', JSON.stringify(norm(psUi)) === JSON.stringify(norm(psExpected)),
            `solo Psicológica → pantalla ${JSON.stringify(norm(psUi))} vs catálogo ${JSON.stringify(norm(psExpected))}`);

        // C5 · los 7 tipos a la vez → todos los subtipos, cada uno con su texto.
        for (const t of tipos) {
            if (t.code === 'ps') continue;
            await hv9.locator('.form-check', { hasText: t.label.trim() }).locator('input').check();
        }
        const allUi = await subtypeLabels();
        const allExpected = TIPOS.flatMap((s) => bySuf[s].map((o) => o.label));
        const missing = norm(allExpected).filter((l) => !norm(allUi).includes(l));
        const unexpected = norm(allUi).filter((l) => !norm(allExpected).includes(l));
        check('C5', allUi.length === allExpected.length && missing.length === 0 && unexpected.length === 0,
            `7 tipos → ${allUi.length} subtipos en pantalla (catálogo: ${allExpected.length}); faltan=[${missing}] sobran/cambiados=[${unexpected}]`);
        await op.context.close();
    } catch (e) {
        check('ERROR', false, e && e.stack ? e.stack.split('\n').slice(0, 3).join(' | ') : String(e));
    } finally {
        try { deleteCases(created); } catch (e) { console.log('(limpieza falló: ' + e.message + ')'); }
        await browser.close();
    }
    finish('REG-10 · subtipos de violencia');
})();
