#!/usr/bin/env node
'use strict';
/**
 * E2E de la pantalla "Registro V3" (/salvia/registro-v3) recorrida de inicio a fin por la UI
 * con los datos del seed. Es el camino feliz E2E-1 de DocsMD/Screens/Registro V3/QA/
 * registro-v3-e2e-qa-cases.md: Borrador → Iniciado → Terminado, con remisión, entidad y barrera.
 *
 *   Habeas data → Tipo de contacto → Cabezote (búsqueda + "Nuevo registro") →
 *   CI → DP → CA → HV → AG (la tríada completa pasa el Caso a "Iniciado") →
 *   Tamizaje → Gestión Salvia (denuncia, remisión, entidad, barrera, alerta inminente) →
 *   Finalizar → Cierre ("Terminado")
 *
 * Todo sale del seed: las opciones de los desplegables y casillas vienen de
 * victim_case_form2_enums (seed_reference_enums.sql), la cascada Departamento → Ciudad →
 * Municipio de la geografía DIVIPOLA (seed_reference_data.sql), y las sedes y obligaciones
 * de la entidad de seed_test_data/20261101-registro-v3-entidades-bogota.sql. Al final se
 * verifica en la BD lo que quedó guardado.
 *
 * Los campos se llenan leyendo el propio catálogo de la pantalla (window.CAMPOS_*_V3) y los
 * que solo aparecen por una respuesta anterior (condicionales) se llenan en pasadas sucesivas.
 *
 * Uso (con la plataforma corriendo):  node cloudDevelopment/e2e/registro-v3.js
 *   SCREENSHOTS=/ruta/dir guarda capturas de cada módulo · E2E_VERBOSE=1 traza completa al fallar.
 */
const { BASE_URL, RUN_ID, rows, one, lit, launch, loginAs, shot, check, scenario, finish } = require('./lib');

const aleatorio = (n) => String(Math.floor(10 ** (n - 1) + Math.random() * 9 * 10 ** (n - 1)));
const DOC = aleatorio(10); // cédula de 10 dígitos: cumple la máscara de CI_6 (CC: 6 a 10 dígitos)
const TEL = '31' + aleatorio(8);
const TEL_APOYO = '32' + aleatorio(8);

// ── Utilidades de la pantalla ──────────────────────────────────────────────
const esc = (s) => s.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');

/** El <div class="mb-3"> de un campo genérico (CI/DP/CA/HV/AG/GS-A), ubicado por el texto exacto de su etiqueta. */
const campoWrapper = (page, label) =>
  page.locator('div.mb-3').filter({ has: page.locator('label.form-label').filter({ hasText: new RegExp(`^\\s*${esc(label)}\\s*$`) }) }).first();

/** Campos del catálogo (window.<varName>) que están visibles ahora mismo. */
async function camposVisibles(page, varName) {
  const defs = await page.evaluate((v) => (window[v] || []).map((c) => ({ id: c.id, label: c.label, tipo: c.tipo, req: c.req || '', max: c.max || 0 })), varName);
  const labels = await page.locator('div.mb-3 > label.form-label').evaluateAll((els) => els.map((e) => e.textContent.trim()));
  return defs.filter((d) => labels.includes(d.label));
}

/**
 * Llena un campo según su tipo. `valor` opcional: boolean (radios), string o RegExp (opción de un
 * desplegable), número (cuántas casillas marcar) o array de RegExp (casillas concretas).
 * Devuelve un texto corto con lo que quedó elegido, para el log.
 */
async function llenarCampo(page, c, valor) {
  const w = campoWrapper(page, c.label);
  await w.waitFor({ state: 'visible', timeout: 8000 });
  switch (c.tipo) {
    case 'text':
      await w.locator('input.form-control').first().fill(String(valor ?? 'Prueba E2E'));
      return String(valor ?? 'Prueba E2E');
    case 'textarea':
      await w.locator('textarea').fill(String(valor ?? 'Texto de prueba E2E del sandbox, con detalle suficiente.'));
      return 'texto';
    case 'date':
      await w.locator('input[type=date]').fill(valor ?? '2026-10-20');
      return valor ?? '2026-10-20';
    case 'boolean':
      await w.locator('input[type=radio]').nth(valor === true ? 0 : 1).check();
      return valor === true ? 'Sí' : 'No';
    case 'dropdown': {
      const sel = w.locator('select');
      const opts = (await sel.locator('option').evaluateAll((os) => os.map((o) => ({ value: o.value, text: o.textContent.trim() })))).filter((o) => o.value);
      if (!opts.length) return '(sin opciones)';
      let elegido = opts[0];
      if (valor instanceof RegExp) elegido = opts.find((o) => valor.test(o.text)) || opts[0];
      else if (typeof valor === 'string') elegido = opts.find((o) => o.text === valor) || opts[0];
      await sel.selectOption(elegido.value);
      return elegido.text.slice(0, 40);
    }
    case 'multiple': {
      const items = w.locator('div.form-check');
      const total = await items.count();
      const marcados = [];
      const marcar = async (i) => { await items.nth(i).locator('input[type=checkbox]').check(); marcados.push((await items.nth(i).innerText()).trim().slice(0, 28)); };
      if (Array.isArray(valor)) {
        for (const re of valor) {
          const textos = await items.evaluateAll((els) => els.map((e) => e.textContent.trim()));
          const i = textos.findIndex((t) => re.test(t));
          if (i >= 0) await marcar(i);
        }
      } else {
        const n = Math.min(typeof valor === 'number' ? valor : 1, c.max || 99, total);
        for (let i = 0; i < n; i++) await marcar(i);
      }
      return marcados.join(' + ') || '(ninguna)';
    }
    default:
      throw new Error(`tipo de campo no soportado: ${c.tipo} (${c.id})`);
  }
}

/** Llena todos los campos visibles de un catálogo, en pasadas: los condicionales aparecen al responder otros. */
async function llenarCatalogo(page, varName, valores = {}) {
  const hechos = new Map();
  for (let pasada = 0; pasada < 6; pasada++) {
    const pendientes = (await camposVisibles(page, varName)).filter((c) => !hechos.has(c.id));
    if (!pendientes.length) break;
    for (const c of pendientes) hechos.set(c.id, await llenarCampo(page, c, valores[c.id]));
  }
  return hechos;
}

const guardarSeccion = async (page) => {
  await page.getByRole('button', { name: 'Guardar sección' }).click();
  const msg = page.locator('span.ms-2').filter({ hasText: /Sección guardada|Faltan obligatorias|Error|no pueden|inválido/ }).first();
  await msg.waitFor({ timeout: 15000 });
  const texto = (await msg.innerText()).trim();
  if (!/Sección guardada/.test(texto)) throw new Error(`no se guardó la sección: ${texto.slice(0, 220)}`);
};
const irAModulo = (page, etiqueta) => page.locator('ul.nav-tabs a.nav-link').filter({ hasText: etiqueta }).first().click();
const radio = (page, nombre, si) => page.locator(`input[type=radio][name="${nombre}"]`).nth(si ? 0 : 1).check();
const casilla = (page, textoRegex) =>
  page.locator('div.form-check').filter({ has: page.locator('label', { hasText: textoRegex }) }).locator('input[type=checkbox]').first();

// ── Recorrido ──────────────────────────────────────────────────────────────
(async () => {
  const browser = await launch();
  const { context, page } = await loginAs(browser, 'op');
  const errores = [];
  page.on('pageerror', (e) => errores.push(String(e.message).slice(0, 120)));
  let icode = null;
  let nivel = null;
  console.log(`Registro V3 como cloud.test.op · documento ${DOC} · teléfono ${TEL} · corrida ${RUN_ID}`);

  await scenario('1-3 · Habeas data → tipo de contacto → Cabezote sin coincidencias → Nuevo registro', async () => {
    await page.goto(`${BASE_URL}/salvia/registro-v3`);
    await page.getByRole('heading', { name: 'Autorización de datos personales' }).waitFor();
    await page.getByRole('button', { name: 'Acepto y continuar' }).click();
    await page.getByRole('heading', { name: 'Tipo de contacto' }).waitFor();
    check(true, 'habeas data aceptado; se muestran los 5 tipos de contacto');
    await page.getByRole('button', { name: 'Víctima - Caso Nuevo' }).click();
    await page.getByRole('heading', { name: 'Buscar registro existente' }).waitFor();
    await page.getByPlaceholder('Número de documento').fill(DOC);
    await page.getByPlaceholder('Teléfono', { exact: true }).fill(TEL);
    const [busq] = await Promise.all([
      page.waitForResponse((r) => /victim-cases-v3\/buscar/.test(r.url())),
      page.getByRole('button', { name: 'Buscar', exact: true }).click(),
    ]);
    check(busq.status() === 200 && ((await busq.json()).registros || []).length === 0, 'la búsqueda del documento nuevo no encuentra coincidencias (no hay duplicado)');
    const [draft] = await Promise.all([
      page.waitForResponse((r) => /victim-cases-v3\/draft/.test(r.url()) && r.request().method() === 'POST'),
      page.getByRole('button', { name: '+ Nuevo registro' }).click(),
    ]);
    check(draft.status() === 200, 'POST draft → 200');
    icode = (await draft.json()).icode;
    await page.getByRole('button', { name: /Finalizar llamada/ }).waitFor();
    const estado = one(`SELECT victim_case_status FROM salvia.victim_case WHERE victim_case_i_code = '${lit(icode)}'`)[0];
    check(estado === 'bo', `Borrador creado: victim_case ${icode} en estado "${estado}" (bo = Borrador)`);
    await shot(page, 'v3-01-borrador');
  });

  await scenario('4 · CI — contacto e identificación (discapacidad "Sí" activa CI_8)', async () => {
    const hechos = await llenarCatalogo(page, 'CAMPOS_CI_V3', {
      CI_1: /^Ella$/i, CI_2: 'María Fernanda Pérez Gómez', CI_3: 'Mafe', CI_4: TEL, CI_5: /Cédula de ciudadanía/i, CI_6: DOC, CI_7: true, CI_8: 1, CI_9: false,
    });
    check(hechos.has('CI_8'), 'CI_7 = Sí hizo aparecer CI_8 (ajuste razonable); elegido: ' + hechos.get('CI_8'));
    check(hechos.get('CI_5').startsWith('Cédula'), `tipo de documento del seed: ${hechos.get('CI_5')}`);
    await guardarSeccion(page);
    check(true, 'CI guardada ("Sección guardada.")');
    await shot(page, 'v3-02-ci');
  });

  await scenario('5 · DP — datos personales (nacionalidad extranjera activa DP_6; etnia indígena activa DP_13)', async () => {
    await irAModulo(page, 'Datos personales');
    const hechos = await llenarCatalogo(page, 'CAMPOS_DP_V3', {
      DP_1: '1990-05-15', DP_2: false, DP_5: /Extranjer/i, DP_8: /^Mujer$/i, DP_9: /^Heterosexual$/i, DP_10: /^Mujer$/i, DP_12: /Indígena/i, DP_14: false,
      DP_27: 'Bogotá, D.C.', DP_28: 'Bogotá', DP_29: 'Carrera 7 # 32-16', DP_26: false,
    });
    check(hechos.has('DP_6'), 'nacionalidad extranjera hizo aparecer DP_6 (país); elegido: ' + hechos.get('DP_6'));
    check(hechos.has('DP_13'), 'etnia indígena hizo aparecer DP_13 (pueblo); elegido: ' + hechos.get('DP_13'));
    await guardarSeccion(page);
    check(true, `DP guardada (${hechos.size} campos llenados)`);
    await shot(page, 'v3-03-dp');
  });

  await scenario('5b · CA — contacto de apoyo', async () => {
    await irAModulo(page, 'Contacto de apoyo');
    await llenarCatalogo(page, 'CAMPOS_CA_V3', { CA_1: 'Laura Gómez Ríos', CA_2: TEL_APOYO, CA_3: 'laura.gomez@example.com' });
    await guardarSeccion(page);
    check(true, 'CA guardada');
  });

  await scenario('6 · HV — hechos (2 tipos de violencia activan HV_10; "Ámbito laboral" activa HV_12)', async () => {
    await irAModulo(page, 'Hechos de violencia');
    const hechos = await llenarCatalogo(page, 'CAMPOS_HV_V3', {
      HV_1: 'La persona relata hechos de violencia ocurridos en el último mes, con detalle suficiente para la prueba.', HV_2: '2026-10-20', HV_3: '19:30',
      HV_4: 'Bogotá, D.C.', HV_5: 'Bogotá', HV_7: 'Calle 45 # 10-20', HV_9: 2, HV_10: 2, HV_11: [/laboral/i, /./], HV_13: true,
    });
    check(hechos.has('HV_10'), 'HV_9 con 2 tipos hizo aparecer HV_10 (subtipos); elegido: ' + hechos.get('HV_10'));
    check(hechos.has('HV_12'), '"Ámbito laboral" en HV_11 hizo aparecer HV_12 (sector laboral); ámbitos: ' + hechos.get('HV_11'));
    await guardarSeccion(page);
    check(true, 'HV guardada');
    await shot(page, 'v3-04-hv');
  });

  await scenario('7-8 · AG — agresor (pareja íntima): la tríada queda completa y el Caso pasa a "Iniciado"', async () => {
    await irAModulo(page, 'Características del agresor');
    const hechos = await llenarCatalogo(page, 'CAMPOS_AG_V3', {
      AG_2: /^(?!.*desconoc)/i, AG_3: /pareja/i, AG_5: false, AG_7: 'Carlos Andrés Ruiz', AG_9: '79123456', AG_10: 'Calle 10 # 5-20', AG_11: '3105550000',
    });
    check(/pareja/i.test(hechos.get('AG_3')), `relación con el agresor: "${hechos.get('AG_3')}" (define la rama "pareja" del Tamizaje)`);
    await guardarSeccion(page);
    await page.getByText('La información mínima de contacto, víctima y hecho está completa').waitFor({ timeout: 10000 });
    check(true, 'aparece el aviso "…el Caso pasó a Iniciado"');
    const estado = one(`SELECT victim_case_status FROM salvia.victim_case WHERE victim_case_i_code = '${lit(icode)}'`)[0];
    check(estado === 'in', `BD: victim_case_status = "${estado}" (in = Iniciado)`);
    const bloqueados = await page.locator('ul.nav-tabs a.nav-link').filter({ hasText: '🔒' }).count();
    check(bloqueados === 0, 'Tamizaje y Gestión Salvia quedaron desbloqueados (sin 🔒)');
    await shot(page, 'v3-05-triada-completa');
  });

  await scenario('9-10 · TZ — Tamizaje: rama pareja, cálculo y guardado', async () => {
    await irAModulo(page, 'Tamizaje de riesgo');
    await page.getByText('Se evalúan primero 7 preguntas generales').waitFor();
    await radio(page, 'TZ_1', true); // una de las 5 iniciales en Sí revela el resto
    for (const id of ['TZ_2', 'TZ_3', 'TZ_4', 'TZ_5', 'TZ_6', 'TZ_7', 'TZ_8']) await radio(page, id, false);
    await page.locator('b', { hasText: 'Pareja/expareja' }).waitFor();
    check(true, 'al responder "Sí" en TZ_1 se revela el resto; rama detectada: Pareja/expareja');
    for (let i = 1; i <= 24; i++) await radio(page, `TZP_${i}`, i <= 8); // exactamente 8 "Sí"
    const [res] = await Promise.all([
      page.waitForResponse((r) => /victim-cases-v3\/.+\/tamizaje/.test(r.url())),
      page.getByRole('button', { name: 'Calcular y guardar Tamizaje' }).click(),
    ]);
    const tz = await res.json();
    check(res.status() === 200 && tz.riskType && tz.riskLevel, `POST tamizaje → 200: tipo "${tz.riskType}", nivel "${tz.riskLevel}", puntaje ${tz.score}`);
    nivel = tz.riskLevel;
    await page.getByText(/Tipo de riesgo:/).waitFor();
    await shot(page, 'v3-06-tamizaje');
  });

  await scenario('11-16 · GS — denuncia, remisión psicosocial, entidad acudida, barrera y alerta inminente', async () => {
    await irAModulo(page, 'Gestión Salvia');
    await page.getByRole('heading', { name: /A\. Denuncia y gestión general/ }).waitFor();

    // § A — denuncia fácil, gestiones interinstitucionales y bitácora
    const a = await llenarCatalogo(page, 'CAMPOS_GS_A_V3', { GS_1: true, GS_2: 2, GS_3: 'Gestionar con la Fiscalía y Salud la activación de la ruta de atención.', GS_4: 'Se acompañó a la persona durante la llamada y se activó la ruta de atención.' });
    check(a.has('GS_1') && a.has('GS_2') && a.has('GS_4'), `§ A llenada: interinstitucional = ${a.get('GS_2')}`);

    // § B — remisión a Atención Psicosocial: el obligatorio + criterios que sumen ≥ 3 puntos
    await radio(page, 'GS_5', true);
    await casilla(page, /^Atención Psicosocial$/).check();
    const crit = await page.evaluate(() => window.ENUM_LABELS_OVERRIDE_V3.victim_case_form2_psycho_criteria);
    await casilla(page, new RegExp('^' + esc(crit.ob.slice(0, 40)))).check();
    await casilla(page, new RegExp('^' + esc(crit.cs))).check(); // "Conducta suicida…(3 pts)"
    await page.getByText('Cumple criterios de remisión').first().waitFor();
    check(true, 'remisión a Atención Psicosocial: el obligatorio + "Conducta suicida (3 pts)" → "✓ Cumple criterios de remisión"');
    await page.locator('label', { hasText: 'Justificación de la remisión' }).locator('xpath=following::textarea[1]')
      .fill('Se cumple el criterio obligatorio y conducta suicida derivada de la violencia.');

    // § C — entidad acudida: Bogotá → Salud → sede real del seed, con sus obligaciones
    await radio(page, 'GS_14', true);
    await radio(page, 'GS_16', true);
    await radio(page, 'GS_15', true);
    const sel = (label) => page.locator('div[class*=col-md]').filter({ has: page.locator('label.form-label', { hasText: new RegExp(`^${esc(label)}$`) }) }).locator('select, input').first();
    await sel('Departamento').selectOption({ label: 'BOGOTÁ, D.C.' });
    await sel('Ciudad').selectOption({ label: 'BOGOTÁ, D.C.' });
    await sel('Municipio').selectOption('11001000');
    await sel('Sector').selectOption({ label: 'Salud' });
    await sel('Entidad').waitFor({ state: 'visible', timeout: 10000 }); // aparece cuando llega la lista de sedes
    const sedes = await sel('Entidad').locator('option').evaluateAll((os) => os.map((o) => o.textContent.trim()).filter((t) => t && !/^Seleccione/.test(t)));
    check(sedes.length === 2, `sedes de Salud en Bogotá desde el seed: ${sedes.join(' · ')}`);
    await sel('Entidad').selectOption({ label: 'ZERENIA S.A.S.' });
    await sel('N° radicado').fill(`RAD-${RUN_ID}`);
    await sel('Objetivo con la entidad').fill('Solicitar atención integral en salud para la persona.');
    await page.getByText('Completadas', { exact: true }).waitFor();
    const completadas = page.locator('label.form-label.small:text-is("Completadas") ~ div.form-check input[type=checkbox]');
    const pendientes = page.locator('label.form-label.small:text-is("Pendientes") ~ div.form-check input[type=checkbox]');
    check((await completadas.count()) > 0, `obligaciones de la entidad (seed): ${await completadas.count()} acciones institucionales para marcar`);
    await completadas.nth(0).check();
    if ((await pendientes.count()) > 1) await pendientes.nth(1).check();
    const [ent] = await Promise.all([
      page.waitForResponse((r) => /victim-cases-v3\/.+\/entidades/.test(r.url()) && r.request().method() === 'POST'),
      page.getByRole('button', { name: '+ Siguiente entidad' }).click(),
    ]);
    check(ent.status() === 200, 'POST entidades → 200');
    const tarjeta = page.locator('div.card.card-body').filter({ hasText: 'ZERENIA S.A.S.' }).filter({ hasText: 'Eliminar entidad' }).first();
    await tarjeta.waitFor();
    check(true, 'la tarjeta de la entidad aparece con su sede y ubicación: ' + (await tarjeta.locator('.small.text-muted').first().innerText()));

    // § D — barrera dentro de la entidad (sector Salud: barreras específicas + una transversal)
    await radio(page, 'GS_17', true);
    await tarjeta.getByRole('button', { name: '+ Agregar barrera a esta entidad' }).click();
    const form = tarjeta.locator('div.border.rounded.p-2');
    await form.waitFor();
    const especificas = form.locator('label.form-label:text-is("Barreras específicas del sector") ~ div.form-check input[type=checkbox]');
    check((await especificas.count()) === 12, `barreras específicas de Salud desde el seed: ${await especificas.count()} opciones`);
    await especificas.nth(0).check();
    await especificas.nth(1).check();
    await form.locator('label.form-label:text-is("Barreras institucionales y de talento humano") ~ div.form-check input[type=checkbox]').first().check();
    const competencias = form.locator('label.form-label:text-is("Competencias institucionales de la entidad") ~ div.form-check input[type=checkbox]');
    await competencias.first().waitFor({ timeout: 10000 });
    await competencias.first().check();
    await form.locator('input[type=date]').fill('2026-10-25');
    await form.locator('label', { hasText: 'Funcionario(a) o dependencia' }).locator('xpath=following::input[1]').fill('Oficina de atención al usuario');
    await form.locator('textarea').fill('Demora en la asignación de citas y falta de activación del protocolo de violencia.');
    await form.locator('label.form-label:text-is("Gestión realizada") ~ div.form-check input[type=checkbox]').first().check();
    const [bar] = await Promise.all([
      page.waitForResponse((r) => /victim-cases-v3\/.+\/barreras/.test(r.url()) && r.request().method() === 'POST'),
      form.getByRole('button', { name: '+ Agregar barrera' }).click(),
    ]);
    check(bar.status() === 200, 'POST barreras → 200');
    await tarjeta.locator('table tbody tr').first().waitFor();
    check(true, 'la barrera aparece en la tabla de la entidad');

    // § E — alerta inminente
    await radio(page, 'GS_18', true);
    await page.locator('label', { hasText: 'Justificación de la alerta' }).locator('xpath=following::textarea[1]')
      .fill('Riesgo alto para la vida de la persona: conducta suicida y violencia reciente.');
    await guardarSeccion(page);
    check(true, 'GS guardada (denuncia, remisión, entidad, barrera y alerta inminente)');
    await shot(page, 'v3-07-gs');
  });

  await scenario('17-19 · Finalizar → Cierre "Terminado" y verificación en la BD', async () => {
    const [fin] = await Promise.all([
      page.waitForResponse((r) => /victim-cases-v3\/.+\/finalizar/.test(r.url())),
      page.getByRole('button', { name: 'Finalizar registro y enrutar a Seguimiento' }).click(),
    ]);
    const cuerpo = await fin.json().catch(() => ({}));
    check(fin.status() === 200 && cuerpo.ok === true, `POST finalizar → ${fin.status()} ${JSON.stringify(cuerpo)}`);
    await page.getByText('Registro finalizado.').waitFor();
    const cierre = await page.locator('main, .container').first().innerText();
    check(/Terminado/.test(cierre), 'pantalla de Cierre: estado "Terminado"');
    check(/Se confirma remisión a|Se solicitó la Denuncia fácil/.test(cierre), 'el Cierre informa la denuncia fácil solicitada y la remisión confirmada');
    check(/Párrafo síntesis/.test(cierre) && /María Fernanda/.test(cierre), 'el Cierre muestra el párrafo síntesis con los datos de la persona');
    await shot(page, 'v3-08-cierre');

    // BD
    const vc = one(`SELECT victim_case_status, victim_case_victim_names, victim_case_victim_doc_number FROM salvia.victim_case WHERE victim_case_i_code = '${lit(icode)}'`);
    check(vc[0] === 'te', `BD: victim_case_status = "${vc[0]}" (te = Terminado)`);
    check(vc[2] === DOC, `BD: documento guardado = ${vc[2]}`);
    const cid = one(`SELECT victim_case_id FROM salvia.victim_case WHERE victim_case_i_code = '${lit(icode)}'`)[0];
    const f2 = one(`SELECT victim_case_form2_residence_address, victim_case_form2_victim_phone FROM salvia.victim_case_form2 WHERE victim_case_form2_victim_case = ${cid}`);
    check(f2[0] === 'Carrera 7 # 32-16', `BD: victim_case_form2.residence_address = "${f2[0]}"`);
    const ent = one(`SELECT count(*), max(entity_branch_id) FROM salvia.entity_case WHERE case_id = '${lit(icode)}' AND deleted_at IS NULL`);
    check(ent[0] === '1' && ent[1] === '9218', `BD: 1 entidad en el caso (entity_case → sede ${ent[1]} = ZERENIA S.A.S.)`);
    const bar = one(`SELECT count(*) FROM salvia.barrier_v2 WHERE case_id = '${lit(icode)}' AND deleted_at IS NULL`);
    check(bar[0] === '1', 'BD: 1 barrera en barrier_v2');
  });

  check(errores.length === 0, `sin excepciones de JavaScript en la página${errores.length ? ': ' + errores.join(' | ') : ''}`);
  await context.close();
  await browser.close();
  finish(`Registro V3 (E2E-1, nivel de riesgo "${nivel}")`);
})().catch((e) => { console.error('ERROR FATAL:', e.message); process.exit(2); });
