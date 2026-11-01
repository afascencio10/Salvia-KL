#!/usr/bin/env node
'use strict';
/**
 * E2E de la pantalla "Registro de Contacto Víctima V2" (el reporte público V2,
 * /salvia/public-v2/nuevo) recorrida de inicio a fin por la UI con los datos del seed:
 *
 *   Autorización → Perfil de ingreso → Formulario específico → Confirmación (radicado)
 *
 * Un escenario por tipo de reportante y rama condicional (Víctima, Tercero anónimo y a
 * nombre propio, Familiar, Hombre). Las listas del asistente (tipo de reporte, ajustes
 * razonables, anonimato, parentesco) se leen del catálogo que entrega la propia pantalla,
 * o sea, de victim_case_form2_enums (seed_reference_enums.sql). Después de cada envío se
 * verifica en la BD que el reporte quedó en victim_contact / victim_contact_form2 /
 * victim_contact_form2_v2 con el tipo de reporte del catálogo.
 *
 * Uso (con la plataforma corriendo):  node cloudDevelopment/e2e/reporte-v2.js
 *   SCREENSHOTS=/ruta/dir  guarda una captura por escenario.
 * Spec: DocsMD/Screens/Registro de Contacto Víctima V2/
 */
const { BASE_URL, RUN_ID, rows, one, lit, launch, bypassCaptcha, loginAs, shot, check, scenario, finish } = require('./lib');

const tel = () => '300' + String(Math.floor(1000000 + Math.random() * 8999999));
const DESCRIPCION = 'Descripción de prueba E2E del sandbox: hechos de violencia reportados para verificar el flujo completo de la pantalla.';

// Un caso por tipo de reportante / rama que cambia los campos del formulario.
const CASOS = [
  {
    id: 'v', titulo: 'Víctima — disposición "Sí" con hora de contacto', tipo: 'v',
    campos: { nombreVictima: 'Ana', apellidoVictima: `Victima-${RUN_ID}`, telefonoVictima: tel(), disposicionLlamada: 'Sí', horaContacto: '10:30', descripcionHechos: DESCRIPCION },
  },
  {
    id: 'r-anonimo', titulo: 'Tercero — anónimo', tipo: 'r',
    campos: { anonimato: 'Anónimo', rolCuidado: 'Sí', victimaConoce: 'No', nombreVictima: 'Beatriz', apellidoVictima: `TerceroAnon-${RUN_ID}`, telefonoVictima: tel(), disposicionLlamada: 'No', descripcionHechos: DESCRIPCION },
  },
  {
    id: 'r-propio', titulo: 'Tercero — a nombre propio', tipo: 'r',
    campos: { anonimato: 'A nombre propio', nombreReportante: 'Rosa', apellidoReportante: `Reporta-${RUN_ID}`, correoReportante: 'rosa.prueba@example.com', rolCuidado: 'No', victimaConoce: 'Sí', nombreVictima: 'Carmen', apellidoVictima: `TerceroProp-${RUN_ID}`, telefonoVictima: tel(), disposicionLlamada: 'No', descripcionHechos: DESCRIPCION },
  },
  {
    id: 'fm', titulo: 'Familiar de víctima de feminicidio', tipo: 'fm',
    campos: { nombreReportante: 'Luis', apellidoReportante: `Familiar-${RUN_ID}`, correoReportante: 'luis.prueba@example.com', parentesco: 'Familia', telefonoReportante: tel(), rolCuidado: 'Sí', disposicionReportante: 'No', nombreVictima: 'Elena', apellidoVictima: `Fallecida-${RUN_ID}`, descripcionHechos: DESCRIPCION },
  },
  {
    id: 'ho', titulo: 'Hombre con posibilidad de ejercer violencia', tipo: 'ho',
    campos: { nombreReportante: 'Mateo', apellidoReportante: `Hombre-${RUN_ID}`, correoReportante: 'mateo.prueba@example.com', telefonoReportante: tel(), disposicionReportante: 'No', descripcionHechos: DESCRIPCION },
  },
];

/** Llena un campo del asistente por su data-field (select por texto visible; el resto con fill). */
async function llenar(page, field, value) {
  const el = page.locator(`[data-field="${field}"]`);
  await el.waitFor({ state: 'visible', timeout: 8000 });
  const tag = await el.evaluate((n) => n.tagName.toLowerCase());
  if (tag === 'select') await el.selectOption({ label: value });
  else await el.fill(value);
}

/** Las listas del asistente: el catálogo JSON que entrega la propia pantalla (viene del seed). */
async function catalogo() {
  const r = await fetch(`${BASE_URL}/salvia/public-v2/nuevo`, { headers: { Accept: 'application/json' } });
  return r.json();
}

/** Recorre el asistente completo para un caso y devuelve { page, radicado }. */
async function radicar(browser, cat, caso) {
  const context = await browser.newContext({ viewport: { width: 1280, height: 900 } });
  const page = await context.newPage();
  await bypassCaptcha(page, '**/salvia/public-v2'); // captcha: se omite con el User-Agent de la app móvil

  // Paso 1 — Autorización
  await page.goto(`${BASE_URL}/salvia/public-v2/nuevo`);
  await page.getByRole('heading', { name: 'Autorización de datos personales' }).waitFor();
  const continuar = page.getByRole('button', { name: 'Continuar' });
  check(await continuar.isDisabled(), 'paso Autorización: "Continuar" deshabilitado hasta responder');
  await llenar(page, 'autorizacion', 'Sí, autorizo');
  await continuar.click();

  // Paso 2 — Perfil de ingreso (opciones = catálogo del seed)
  const tipo = cat.reportType.find((t) => t.code === caso.tipo);
  const ajuste = cat.adjustmentsGBV.find((a) => a.code === 'nr');
  await llenar(page, 'tipoReporte', tipo.name);
  await llenar(page, 'ajusteRazonable', ajuste.name);
  await page.getByRole('button', { name: 'Continuar' }).click();

  // Paso 3 — Formulario específico del tipo de reportante
  await page.getByRole('heading', { name: new RegExp(`^Datos del Reporte \\(${tipo.name.slice(0, 20)}`) }).waitFor();
  for (const [campo, valor] of Object.entries(caso.campos)) await llenar(page, campo, valor);
  await llenar(page, 'captchaSolution', '0000');
  const enviar = page.getByRole('button', { name: 'Enviar Reporte' });
  check(await enviar.isEnabled(), 'paso Formulario: "Enviar Reporte" se habilita con todo lo obligatorio');
  await shot(page, `v2-${caso.id}-formulario`);
  await enviar.click();

  // Paso 4 — Confirmación con radicado
  await page.getByText('¡Reporte Enviado!').waitFor({ timeout: 20000 });
  const linea = await page.locator('p', { hasText: 'Número de Radicado:' }).innerText();
  const radicado = linea.split(':').pop().trim();
  await shot(page, `v2-${caso.id}-confirmacion`);
  return { page, radicado };
}

/** Verifica en la BD que el reporte quedó donde corresponde. */
function verificarEnBD(caso, radicado) {
  const vc = one(`SELECT victim_contact_id, victim_contact_names, victim_contact_last_names, victim_contact_status, COALESCE(victim_contact_agent_id,'')
                  FROM salvia.victim_contact WHERE victim_contact_i_code = '${lit(radicado)}'`);
  check(vc.length > 0, `victim_contact creado con el radicado ${radicado} (estado "${vc[3]}", agente "${vc[4] || 'sin asignar'}")`);
  if (caso.tipo !== 'ho') {
    check(vc[1] === caso.campos.nombreVictima && vc[2] === caso.campos.apellidoVictima, `víctima guardada: ${vc[1]} ${vc[2]}`);
  } else {
    check(vc[1] === '' && vc[2] === '', 'Hombre: nombre/apellido de la víctima vacíos (los completa el agente)');
  }
  const f2 = one(`SELECT f.victim_contact_form2_id, e.victim_case_form2_enums_code, COALESCE(f.victim_contact_form2_reporter_names,'')
                  FROM salvia.victim_contact_form2 f
                  JOIN salvia.victim_case_form2_enums e ON e.victim_case_form2_enums_id = f.victim_contact_form2_report_type
                  WHERE f.victim_contact_form2_victim_contact = ${vc[0]}`);
  check(f2[1] === caso.tipo, `victim_contact_form2.report_type = "${f2[1]}" (enum del seed) — esperado "${caso.tipo}"`);
  const v2 = one(`SELECT COALESCE(v.victim_contact_form2_v2_reporter_last_names,''), COALESCE(a.victim_case_form2_enums_code,''), COALESCE(p.victim_case_form2_enums_code,'')
                  FROM salvia.victim_contact_form2_v2 v
                  LEFT JOIN salvia.victim_case_form2_enums a ON a.victim_case_form2_enums_id = v.victim_contact_form2_v2_anonimato
                  LEFT JOIN salvia.victim_case_form2_enums p ON p.victim_case_form2_enums_id = v.victim_contact_form2_v2_parentesco
                  WHERE v.victim_contact_form2_v2_victim_contact_form2 = ${f2[0]}`);
  check(v2.length > 0, 'fila de extensión victim_contact_form2_v2 creada');
  if (caso.campos.apellidoReportante) check(v2[0] === caso.campos.apellidoReportante, `apellido del reportante guardado: ${v2[0]}`);
  if (caso.campos.anonimato) check(v2[1] === (caso.campos.anonimato === 'Anónimo' ? 'an' : 'pr'), `anonimato guardado: ${v2[1]}`);
  if (caso.campos.parentesco) check(v2[2] === 'fa', `parentesco guardado: ${v2[2]}`);
}

(async () => {
  const browser = await launch();
  const cat = await catalogo();
  console.log(`Catálogos del asistente (seed): ${cat.reportType.length} tipos de reporte, ${cat.adjustmentsGBV.length} ajustes, ${cat.anonimato.length} anonimato, ${cat.parentesco.length} parentesco · corrida ${RUN_ID}`);

  await scenario('E-02 · "No autorizo" corta el flujo', async () => {
    const page = await (await browser.newContext()).newPage();
    await page.goto(`${BASE_URL}/salvia/public-v2/nuevo`);
    await llenar(page, 'autorizacion', 'No autorizo');
    await page.getByText('No es posible continuar con el reporte sin la autorización').waitFor();
    check(await page.getByRole('button', { name: 'Continuar' }).isDisabled(), 'aviso de corte visible y "Continuar" deshabilitado');
    await shot(page, 'v2-no-autoriza');
    await page.context().close();
  });

  await scenario('E-06/E-09 · "Atrás" conserva lo diligenciado', async () => {
    const page = await (await browser.newContext()).newPage();
    await page.goto(`${BASE_URL}/salvia/public-v2/nuevo`);
    await llenar(page, 'autorizacion', 'Sí, autorizo');
    await page.getByRole('button', { name: 'Continuar' }).click();
    await llenar(page, 'tipoReporte', cat.reportType.find((t) => t.code === 'v').name);
    await llenar(page, 'ajusteRazonable', cat.adjustmentsGBV.find((a) => a.code === 'ni').name);
    await page.getByRole('button', { name: 'Atrás' }).click();
    const autor = await page.locator('[data-field="autorizacion"]').evaluate((n) => n.options[n.selectedIndex].text);
    check(autor === 'Sí, autorizo', 'al volver a Autorización, la respuesta sigue seleccionada');
    await page.getByRole('button', { name: 'Continuar' }).click();
    const tipoSel = await page.locator('[data-field="tipoReporte"]').evaluate((n) => n.options[n.selectedIndex].text);
    check(tipoSel.startsWith('Usted es la víctima'), 'al volver a Perfil, el tipo de reporte sigue seleccionado');
    await page.context().close();
  });

  for (const caso of CASOS) {
    await scenario(`E-10 · Radicar reporte — ${caso.titulo}`, async () => {
      const { page, radicado } = await radicar(browser, cat, caso);
      check(/^[0-9a-f-]{36}$/.test(radicado), `confirmación con radicado ${radicado}`);
      verificarEnBD(caso, radicado);
      caso.radicado = radicado;
      await page.context().close();
    });
  }

  await scenario('E-11 · "Volver al Inicio" reinicia el asistente', async () => {
    const caso = { id: 'v-reinicio', tipo: 'v', campos: { nombreVictima: 'Ines', apellidoVictima: `Reinicio-${RUN_ID}`, telefonoVictima: tel(), disposicionLlamada: 'No', descripcionHechos: DESCRIPCION } };
    const { page } = await radicar(browser, cat, caso);
    await page.getByRole('button', { name: /Volver al Inicio/ }).click();
    await page.getByRole('heading', { name: 'Autorización de datos personales' }).waitFor();
    const autor = await page.locator('[data-field="autorizacion"]').evaluate((n) => n.selectedIndex);
    check(autor === 0, 'vuelve al paso Autorización con el formulario en blanco');
    await page.context().close();
  });

  // Ciclo de vida completo de un reporte V2 (GR-E2E de DocsMD/Screens/Listado Reportes/QA):
  // radicación pública → supervisor lo ve "Sin asignar" y lo reasigna → agente lo gestiona
  // (3x3) y lo cierra → supervisor ve el resultado.
  await scenario('Ciclo de vida · radicar → reasignar (sv) → gestionar y cerrar (op) → resultado (sv)', async () => {
    const apellido = `Ciclo-${RUN_ID}`;
    const caso = { id: 'ciclo', tipo: 'v', campos: { nombreVictima: 'Paula', apellidoVictima: apellido, telefonoVictima: tel(), disposicionLlamada: 'Sí', horaContacto: '10:00', descripcionHechos: DESCRIPCION } };

    /** Busca por apellido en la tabla de casos-component y devuelve la única fila. */
    const filaDe = async (page) => {
      // La búsqueda tiene debounce y recarga la tabla (y borra la selección): se espera su respuesta
      // antes de tocar la fila, si no se marcaría una fila de la lista sin filtrar.
      await Promise.all([
        page.waitForResponse((r) => /victim-contacts\/list/.test(r.url()) && r.url().includes(apellido), { timeout: 15000 }),
        page.getByRole('textbox', { name: /Buscar por nombres/ }).fill(apellido),
      ]);
      const filas = page.locator('tbody tr').filter({ hasText: apellido });
      await filas.first().waitFor({ timeout: 15000 });
      check((await filas.count()) === 1, `la búsqueda "${apellido}" devuelve una sola fila`);
      return filas.first();
    };
    const textoFila = async (fila) => (await fila.innerText()).replace(/\s+/g, ' ');
    /** Comprueba que el texto contenga todas las piezas (sin distinguir mayúsculas: la tabla las pinta por CSS). */
    const contiene = (txt, ...piezas) => piezas.every((p) => txt.toLowerCase().includes(p.toLowerCase()));

    // C1 — radicación por el formulario público
    const { page: pub, radicado } = await radicar(browser, cat, caso);
    await pub.context().close();
    check(!!radicado, `C1 · reporte radicado ${radicado}`);

    // C2 — supervisor: aparece "Sin asignar" y lo reasigna al agente de prueba
    const sv = await loginAs(browser, 'sv');
    await sv.page.goto(`${BASE_URL}/salvia/reportes`);
    await sv.page.getByRole('heading', { name: 'Reportes' }).waitFor();
    let fila = await filaDe(sv.page);
    let txt = await textoFila(fila);
    check(contiene(txt, `ID: ${radicado}`, 'Víctima', 'Sin intentos', 'Sin asignar'),
      `C2 · fila con el radicado, tipo "Víctima", "Sin intentos" y estado "Sin asignar" — fila: "${txt.slice(0, 160)}"`);
    await fila.getByRole('checkbox').check();
    await sv.page.getByRole('button', { name: /Reasignar Reportes/ }).click();
    await sv.page.getByRole('combobox', { name: 'Nueva persona asignada' }).selectOption({ label: 'Cloud Operador — Sin equipo' });
    await sv.page.getByRole('dialog').getByRole('button', { name: 'Reasignar', exact: true }).click();
    const [reasig] = await Promise.all([
      sv.page.waitForResponse((r) => /victim-contacts\/reasignar-bulk/.test(r.url())),
      sv.page.getByRole('button', { name: 'Sí, reasignar' }).click(),
    ]);
    check(reasig.status() === 200 && (await reasig.json()).updated === 1, 'C2 · POST reasignar-bulk → 200, 1 reporte actualizado');
    const agente = one(`SELECT COALESCE(victim_contact_agent_id,'') FROM salvia.victim_contact WHERE victim_contact_i_code = '${lit(radicado)}'`)[0];
    check(agente === 'user-icode-cloud.test.op', `C2 · BD: victim_contact_agent_id = ${agente}`);
    await sv.context.close();

    // C3 — agente: gestiona (3x3) y cierra por decisión de la persona
    const op = await loginAs(browser, 'op');
    await op.page.goto(`${BASE_URL}/salvia/mis-reportes`);
    await op.page.getByRole('heading', { name: 'Mis Reportes' }).waitFor();
    fila = await filaDe(op.page);
    txt = await textoFila(fila);
    check(contiene(txt, 'Asignado', 'Sin intentos'), `C3 · en Mis Reportes: "Asignado" y "Sin intentos" — fila: "${txt.slice(0, 120)}"`);
    await fila.getByRole('button', { name: 'Ver Reporte' }).click();
    const dlg = op.page.getByRole('dialog').first();
    await dlg.waitFor();
    check((await dlg.innerText()).includes(radicado) && (await dlg.innerText()).includes(apellido), 'C3 · el modal abre con el radicado y el apellido');
    check((await dlg.innerText()).includes('No requiere'), 'C3 · el modal muestra el ajuste razonable que eligió el reportante ("No requiere")');
    await op.page.getByRole('button', { name: 'Procede recontacto' }).click();
    await dlg.getByText('Hoy: 0/3 · Días completos: 0/3').waitFor();
    check(true, 'C3 · esquema 3x3 inicial: "Hoy: 0/3 · Días completos: 0/3"');
    await op.page.getByRole('button', { name: 'Intento Fallido' }).click();
    await dlg.getByText('Hoy: 1/3 · Días completos: 0/3').waitFor();
    check((await dlg.innerText()).includes('Fallido'), 'C3 · intento fallido: "Hoy: 1/3" y una fila "Fallido" en el log');
    await op.page.getByRole('button', { name: 'Intento Exitoso' }).click();
    await dlg.getByRole('heading', { name: 'Contacto establecido' }).waitFor();
    check((await dlg.getByRole('button', { name: /^Cierre por/ }).count()) === 2, 'C3 · intento exitoso: "Contacto establecido" con solo dos formas de cierre');
    await op.page.getByRole('button', { name: 'Cierre por decisión de la persona' }).click();
    const confirmar = op.page.getByRole('button', { name: 'Confirmar cierre' });
    check(await confirmar.isDisabled(), 'C3 · "Confirmar cierre" deshabilitado hasta escribir la manifestación');
    const observacion = `La persona manifiesta que no desea continuar (${RUN_ID})`;
    await op.page.getByRole('textbox', { name: 'Escriba aquí...' }).fill(observacion);
    const [cierre] = await Promise.all([
      op.page.waitForResponse((r) => /victim-contacts\/.+\/cerrar/.test(r.url())),
      confirmar.click(),
    ]);
    const cuerpo = cierre.request().postDataJSON();
    check(cierre.status() === 200 && cuerpo.close_reason === 'dp', `C3 · POST cerrar → 200 con close_reason "${cuerpo.close_reason}"`);
    await dlg.getByText(/cerrado/i).first().waitFor();
    const final = await dlg.innerText();
    check(final.includes(radicado) && /inv[áa]lido/i.test(final) && final.includes(observacion), 'C3 · el modal confirma el cierre con el motivo y la observación');
    await shot(op.page, 'v2-ciclo-modal-cerrado');
    // BD: motivo del enum del seed, observación, fecha de cierre y los dos intentos del agente
    const bd = one(`SELECT e.victim_case_form2_enums_code, vc.victim_contact_close_observation, (vc.victim_contact_closed_date IS NOT NULL)
                    FROM salvia.victim_contact vc
                    JOIN salvia.victim_case_form2_enums e ON e.victim_case_form2_enums_id = vc.victim_contact_close_reason
                    WHERE vc.victim_contact_i_code = '${lit(radicado)}'`);
    check(bd[0] === 'dp' && bd[1] === observacion && bd[2] === 't', `C3 · BD: motivo "${bd[0]}" (enum victim_contact_close_reason del seed), observación y fecha de cierre`);
    const intentos = rows(`SELECT a.victim_contact_attempt_result FROM salvia.victim_contact_attempt a
                           JOIN salvia.victim_contact vc ON vc.victim_contact_i_code = a.victim_contact_attempt_victim_contact
                           WHERE vc.victim_contact_i_code = '${lit(radicado)}' AND a.victim_contact_attempt_deleted_at IS NULL
                           ORDER BY a.victim_contact_attempt_id`).map((r) => r[0]);
    check(intentos.length === 2, `C3 · BD: ${intentos.length} intentos del agente registrados (${intentos.join(', ')})`);
    await op.context.close();

    // C4 — supervisor: ve el resultado, en solo lectura
    const sv2 = await loginAs(browser, 'sv');
    await sv2.page.goto(`${BASE_URL}/salvia/reportes`);
    await sv2.page.getByRole('heading', { name: 'Reportes' }).waitFor();
    fila = await filaDe(sv2.page);
    txt = await textoFila(fila);
    check(contiene(txt, 'Cerrado', 'inválido', 'Cloud Operador'), `C4 · fila del supervisor: "${txt.replace(`ID: ${radicado}`, '').slice(0, 130)}"`);
    await fila.getByRole('button', { name: 'Ver Reporte' }).click();
    const dlg2 = sv2.page.getByRole('dialog').first();
    await dlg2.waitFor();
    check((await dlg2.innerText()).includes(observacion), 'C4 · "Ver Reporte" muestra el mismo motivo y observación que vio el agente');
    check((await dlg2.getByRole('button', { name: /Procede recontacto|Intento/ }).count()) === 0, 'C4 · el supervisor lo ve en solo lectura (sin botones de gestión)');
    await shot(sv2.page, 'v2-ciclo-supervisor-cerrado');
    await sv2.context.close();
  });

  await browser.close();
  finish('Reporte V2 (asistente público + ciclo de vida)');
})().catch((e) => { console.error('ERROR FATAL:', e.message); process.exit(2); });
