// Pantalla · T8 (REG-17): módulo «Ruta de atención integral» (réplica de «Identificación de entidades» de Seguimiento v4).
const { abrir, elegirEntidad } = require('./ui-lib');
const OUT = process.env.SHOTS || '/tmp/wp4-shots';
require('fs').mkdirSync(OUT, { recursive: true });
let ok = 0, fail = 0;
const check = (d, a, e) => { if (JSON.stringify(a) === JSON.stringify(e)) { ok++; console.log('  ✅', d); } else { fail++; console.log('  ❌', d, '\n      esperado:', JSON.stringify(e), '\n      obtenido:', JSON.stringify(a)); } };

(async () => {
  // Caso en borrador con la tríada INCOMPLETA: la pestaña nueva no debe tener candado (GS sí).
  const t = await abrir({
    triada: false,
    hook: async ({ api, icode }) => {
      await api('POST', `${icode}/entidades`, { kind: 'activacion', entityBranchId: 102, otherInfo: 'obj act' });
      await api('POST', `${icode}/entidades`, { kind: 'acudida', entityBranchId: 9217, fileNumber: 'RAD-9', otherInfo: 'obj acu', reportDate: '2026-10-01' });
    },
  });
  const { page, icode, sql } = t;

  console.log('== 1. Pestaña nueva entre Tamizaje y Plan de acción, sin candado');
  const tabs = (await page.locator('ul.nav-tabs a').allInnerTexts()).map((x) => x.trim());
  const nombres = tabs.map((x) => x.replace(/\s*🔒/, ''));
  const iTZ = nombres.indexOf('Tamizaje de riesgo'), iRA = nombres.indexOf('Ruta de atención integral'), iGS = nombres.indexOf('Plan de acción Salvia');
  check('orden de pestañas: … Tamizaje, Ruta de atención integral, Plan de acción', [iTZ >= 0, iRA === iTZ + 1, iGS === iRA + 1], [true, true, true]);
  check('«Ruta de atención integral» sin candado y «Plan de acción Salvia» bloqueado (tríada incompleta)', [!tabs[iRA].includes('🔒'), tabs[iGS].includes('🔒')], [true, true]);
  await page.getByRole('link', { name: /Ruta de atención integral/ }).click();
  check('el módulo abre con la tríada incompleta', await page.locator('input[name="GS_14"]').count(), 2);

  console.log('== 2. Pregunta y tarjetas: la acudida solo aquí, la activación solo en el Plan de acción');
  check('texto de GS_14 = el de Seguimiento v4', (await page.locator('div.mb-3 > label.form-label').first().innerText()).replace(/\s*\*$/, '').trim(), '¿Ha puesto en conocimiento de alguna entidad o institución los hechos de violencia reportados?');
  check('sin responder GS_14 no hay entidades ni formulario', [await page.locator('.card.card-body.mb-2').count(), await page.locator('.card.card-body.mb-3').count()], [0, 0]);
  await page.locator('input[name="GS_14"]').first().check();
  check('con GS_14 = Sí se ve solo la entidad acudida (1 tarjeta)', await page.locator('.card.card-body.mb-2').count(), 1);
  check('la tarjeta es la de la sede acudida (ZONAMEDICA…)', (await page.locator('.card.card-body.mb-2').first().innerText()).includes('ZONAMEDICA'), true);

  console.log('== 3. Formulario: preguntas y obligatoriedad');
  const form = page.locator('.card.card-body.mb-3').filter({ hasText: 'Entidad a la que ya acudió' });
  check('hay formulario de entidad acudida', await form.count(), 1);
  await elegirEntidad(page, form, { sector: 'justicia', texto: 'usaqu' });
  const etiquetas = (await form.locator('label.form-label').allInnerTexts()).map((x) => x.replace(/\s*\*$/, '').trim());
  check('pide radicado, fecha y otra información asociada a la Ruta', ['Número de radicado (NUNC, medida de protección). Escriba tipo y número', 'Fecha en la que reportó a la entidad', 'Otra información asociada a la Ruta'].every((x) => etiquetas.includes(x)), true);
  check('hay chips de completadas y de pendientes', [await form.locator('label:has-text("Completadas")').count(), await form.locator('label:has-text("Pendientes")').count()], [1, 1]);
  const boton = form.getByRole('button', { name: '+ Siguiente entidad' });
  check('«+ Siguiente entidad» deshabilitado y explica qué falta', [await boton.isDisabled(), (await form.innerText()).includes('Falta:')], [true, true]);
  await form.locator('label.form-label:has-text("Número de radicado") + input').fill('NUNC-555');
  await form.locator('label.form-label:has-text("Fecha en la que reportó") + input[type=date]').fill('2026-10-15');
  await form.locator('label.form-label:has-text("Otra información asociada a la Ruta") + textarea').fill('Acudió a la URI');
  check('sin acciones marcadas sigue deshabilitado', await boton.isDisabled(), true);
  const completadas = form.locator('label:has-text("Completadas") + div button');
  const pendientes = form.locator('label:has-text("Pendientes") + div button');
  await completadas.nth(0).click();
  check('con una acción completada y sin pendientes se habilita (no hay bloque de activación)', [await boton.isEnabled(), await form.locator('label:has-text("¿La persona está de acuerdo")').count()], [true, 0]);

  console.log('== 4. D17: con acciones pendientes se piden acuerdo, debida diligencia, canales y gestiones');
  await pendientes.nth(1).click();
  await pendientes.nth(2).click();
  check('aparece el acuerdo y el botón se deshabilita', [await form.locator('label.form-label:has-text("¿La persona está de acuerdo")').count(), await boton.isDisabled()], [1, true]);
  check('la debida diligencia aún no se ve', await form.locator('label.form-label:has-text("¿Salvia debe activar rutas")').count(), 0);
  await form.locator('input[name="RA_acuerdo"]').first().check();
  check('con acuerdo = Sí aparece la debida diligencia', await form.locator('label.form-label:has-text("¿Salvia debe activar rutas")').count(), 1);
  check('los canales aún no', await form.locator('label.form-label:has-text("Seleccione los canales")').count(), 0);
  await form.locator('input[name="RA_debida"]').first().check();
  check('con ambas en Sí aparecen los canales y las gestiones', [await form.locator('label.form-label:has-text("Seleccione los canales")').count(), await form.locator('label.form-label:has-text("¿Qué gestiones se realizaron")').count()], [1, 1]);
  check('sigue deshabilitado sin canales ni gestiones', await boton.isDisabled(), true);
  await form.locator('.form-check', { hasText: 'Notificación SALVIA' }).locator('input').check();
  await form.locator('.form-check', { hasText: 'Mi Salvia' }).locator('input').check();
  await form.locator('label.form-label:has-text("¿Qué gestiones se realizaron") + textarea').fill('Se envió oficio y se radicó en Mi Salvia');
  check('con todo se habilita', await boton.isEnabled(), true);
  await page.screenshot({ path: `${OUT}/t8-1-formulario.png` });
  await boton.click();
  await page.waitForTimeout(900);

  console.log('== 5. Lo guardado en el borrador');
  const fila = sql(`select e->>'kind'||'|'||(e->>'fileNumber')||'|'||(e->>'reportDate')||'|'||(e->>'objetivo')||'|'||(e->>'agreement')||'|'||(e->>'dueDiligence')||'|'||(e->'activationChannels')::text||'|'||jsonb_array_length(e->'pending')||'|'||jsonb_array_length(e->'completed') from salvia.victim_case, jsonb_array_elements(victim_case_v3_draft_data->'entities') e where victim_case_i_code='${icode}' and e->>'fileNumber'='NUNC-555'`);
  check('acudida con sus datos, canales y acciones', fila, 'acudida|NUNC-555|2026-10-15|Acudió a la URI|true|true|["notificacion", "misalvia"]|2|1');
  check('sigue sin haber filas en entity_case ni tareas (borrador)', sql(`select (select count(*) from salvia.entity_case where case_id='${icode}')||'/'||(select count(*) from salvia.case_task where case_id='${icode}')`), '0/0');
  check('ahora hay 2 tarjetas acudidas', await page.locator('.card.card-body.mb-2').count(), 2);

  console.log('== 6. Barreras: GS_17 = Sí muestra dónde registrarlas y el botón en la tarjeta');
  check('sin GS_17 no hay botón de barrera', await page.getByRole('button', { name: '+ Agregar barrera a esta entidad' }).count(), 0);
  await page.locator('input[name="GS_17"]').first().check();
  check('aviso con dónde registrar', (await page.locator('main').innerText()).includes('Registra cada barrera con el botón'), true);
  const botones = page.getByRole('button', { name: '+ Agregar barrera a esta entidad' });
  check('un botón por tarjeta', await botones.count(), 2);
  await botones.first().click();
  check('abre el formulario de barrera de Seguimiento v4', (await page.locator('.border.rounded.p-2.mt-2 label.form-label').first().innerText()).includes('Sector en el que se presentó la barrera'), true);
  await page.getByRole('button', { name: 'Cancelar' }).click();

  console.log('== 7. Guardar sección guarda GS_14 y GS_17');
  await page.getByRole('button', { name: 'Guardar sección' }).click();
  await page.waitForTimeout(900);
  check('mensaje de guardado', (await page.locator('span.ms-2').allInnerTexts()).some((x) => x.includes('Sección guardada')), true);
  const llaves = sql(`select victim_case_v3_answered_keys from salvia.victim_case where victim_case_i_code='${icode}'`);
  check('las llaves GS_14 y GS_17 quedaron como respondidas', [llaves.includes('GS_14'), llaves.includes('GS_17')], [true, true]);
  check('en BD: reported_to_entity e identifies_barriers respondidas', sql(`select (f.victim_case_form2_reported_to_entity is not null)||'/'||(f.victim_case_form2_identifies_barriers is not null) from salvia.victim_case_form2 f join salvia.victim_case vc on vc.victim_case_id=f.victim_case_form2_victim_case where vc.victim_case_i_code='${icode}'`), 'true/true');

  console.log('== 8. Con GS_17 = Sí y sin entidades el formulario se muestra aunque GS_14 no sea Sí');
  const t2 = await abrir({ triada: false });
  await t2.page.getByRole('link', { name: /Ruta de atención integral/ }).click();
  await t2.page.locator('input[name="GS_14"]').nth(1).check();
  check('con GS_14 = No no hay formulario', await t2.page.locator('.card.card-body.mb-3').count(), 0);
  await t2.page.locator('input[name="GS_17"]').first().check();
  check('con GS_17 = Sí y sin entidades se muestra el formulario y el aviso «Primero agrega la entidad»', [await t2.page.locator('.card.card-body.mb-3').filter({ hasText: 'Entidad a la que ya acudió' }).count(), (await t2.page.locator('main').innerText()).includes('Primero agrega la entidad')], [1, true]);
  await t2.page.getByRole('button', { name: 'Guardar sección' }).click();
  check('Guardar sección responde sin errores', true, true);

  console.log('== 9. Errores de consola');
  check('sin errores de JS', [...t.errores, ...t2.errores], []);
  await t2.browser.close();
  await t.browser.close();
  console.log(`\nRESULTADO: ${ok} ok / ${fail} fallos`);
  process.exit(fail ? 1 : 0);
})().catch((e) => { console.error(e); process.exit(2); });
