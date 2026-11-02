'use strict';
/**
 * Utilidades de los scripts de prueba de WP-3 (scripts/e2e/wp3-*.js). Reutilizan wp2-lib.js (login con captcha
 * omitido, fetch con la sesión del navegador, SQL, borrado de casos de prueba) y agregan lo propio de WP-3.
 *
 * Requisitos: igual que wp2-lib (app corriendo en BASE_URL con ./run_dev.sh, psql en el PATH).
 */
const lib = require('./wp2-lib');

/** Respuestas por módulo para poder FINALIZAR un Caso (todas las obligatorias que exige el servidor). */
function finalizeAnswers(catalog, doc) {
    const { first } = catalog;
    const v = lib.validAnswers(catalog, doc);
    return {
        CI: v.CI,
        DP: { ...v.DP, DP_10: first('victim_case_form2_assigned_sex_at_birth'), DP_15: first('victim_case_form2_marital_status') },
        CA: { CA_1: 'Contacto Apoyo', CA_2: '3009876543', CA_3: 'apoyo@example.com', CA_4: first('victim_case_form2_support_contact_kinship') },
        HV: v.HV,
        AG: v.AG,
        TZ5: { TZ_1: 'false', TZ_2: 'false', TZ_3: 'false', TZ_4: 'false', TZ_7: 'false' },
        GS: {
            GS_1: 'false', GS_2: first('victim_case_form2_interinst_action'), GS_4: 'Acompañamiento de prueba',
            GS_5: 'false', GS_14: 'false', GS_16: 'false', GS_17: 'false', GS_18: 'false',
        },
    };
}

/** POST /seccion/:modulo con la sesión del navegador. */
function save(page, icode, modulo, answers) {
    return lib.api(page, 'POST', `/api/v1/victim-cases-v3/${icode}/seccion/${modulo}`, { answers });
}

/** POST /tamizaje con solo las 5 iniciales (REG-34: ya se puede guardar). */
function saveTamizaje(page, icode, answers) {
    return lib.api(page, 'POST', `/api/v1/victim-cases-v3/${icode}/tamizaje`,
        { answers, wasPartner: false, esOsigd: false, subtipoFeminicidio: false, subtipoHomicidioAbyecto: false });
}

/** Crea un Borrador por API con un documento único. Devuelve { icode, doc }. */
async function newDraft(page, created, doc) {
    const d = doc || String(800000000 + Math.floor(Math.random() * 99999999));
    const r = await lib.api(page, 'POST', '/api/v1/victim-cases-v3/draft', { tipo_documento: 'cc', documento: d, telefono: '3009998888' });
    if (created) created.push(r.json.icode);
    return { icode: r.json.icode, doc: d };
}

/** Completa la tríada (CI, DP, HV, AG) por API. */
async function completeTriad(page, icode, A) {
    for (const m of ['CI', 'DP', 'HV', 'AG']) await save(page, icode, m, A[m]);
}

/** Deja el Caso listo para finalizar salvo GS: tríada + DP/CA de 2º contacto + Tamizaje inicial. */
async function readyExceptGS(page, icode, A) {
    await completeTriad(page, icode, A);
    await save(page, icode, 'CA', A.CA);
    await saveTamizaje(page, icode, A.TZ5);
}

/** Fija respuestas en la app Vue (equivale a escribirlas en los campos) sin depender de cada control. */
function setAnswers(page, answers) {
    return page.evaluate((a) => {
        const p = document.getElementById('app')._vnode.component.proxy;
        Object.keys(a).forEach((k) => { p.answers[k] = a[k]; });
    }, answers);
}

/** Quita respuestas de la app Vue (una pregunta sin responder no viaja en el guardado; un texto vacío sí cuenta como respondido). */
function unsetAnswers(page, keys) {
    return page.evaluate((ks) => {
        const p = document.getElementById('app')._vnode.component.proxy;
        ks.forEach((k) => { delete p.answers[k]; });
    }, keys);
}

/** Lee un valor de la app Vue (data o método sin argumentos). */
function vm(page, expr) {
    return page.evaluate((e) => {
        const p = document.getElementById('app')._vnode.component.proxy;
        return new Function('p', `return (${e});`)(p);
    }, expr);
}

/** Estado del Caso en la BD. */
function status(icode) {
    return lib.sql(`SELECT victim_case_status FROM salvia.victim_case WHERE victim_case_i_code = '${icode}'`);
}

/** Luminancia relativa y contraste WCAG entre dos colores "rgb(r, g, b)". */
function contrast(fg, bg) {
    const parse = (s) => (s.match(/\d+(\.\d+)?/g) || [0, 0, 0]).slice(0, 3).map(Number);
    const lum = (c) => {
        const [r, g, b] = c.map((v) => { const x = v / 255; return x <= 0.03928 ? x / 12.92 : Math.pow((x + 0.055) / 1.055, 2.4); });
        return 0.2126 * r + 0.7152 * g + 0.0722 * b;
    };
    const a = lum(parse(fg)), b = lum(parse(bg));
    return (Math.max(a, b) + 0.05) / (Math.min(a, b) + 0.05);
}

module.exports = { ...lib, finalizeAnswers, save, saveTamizaje, newDraft, completeTriad, readyExceptGS, setAnswers, unsetAnswers, vm, status, contrast };
