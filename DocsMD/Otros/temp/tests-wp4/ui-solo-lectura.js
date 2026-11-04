// Pantalla · D20: tras materializar (Caso en «Iniciado») entidades, barreras y remisiones son de solo lectura.
const { abrir, elegirEntidad } = require('./ui-lib');
const OUT = process.env.SHOTS || '/tmp/wp4-shots';
require('fs').mkdirSync(OUT, { recursive: true });
let ok = 0, fail = 0;
const check = (d, a, e) => { if (JSON.stringify(a) === JSON.stringify(e)) { ok++; console.log('  ✅', d); } else { fail++; console.log('  ❌', d, '\n      esperado:', JSON.stringify(e), '\n      obtenido:', JSON.stringify(a)); } };

(async () => {
  console.log('== A. Caso ya materializado (Iniciado)');
  const t = await abrir({
    hook: async ({ api, icode, sql }) => {
      const eq = sql(`select victim_case_form2_enums_i_code from salvia.victim_case_form2_enums where victim_case_form2_enums_category='victim_case_form2_referral_team' and victim_case_form2_enums_code='ah'`);
      const crit = sql(`select victim_case_form2_enums_i_code from salvia.victim_case_form2_enums where victim_case_form2_enums_category='victim_case_form2_men_criteria' and victim_case_form2_enums_code='ob'`);
      await api('POST', `${icode}/seccion/GS`, { answers: { GS_5: 'true', GS_6: eq, GS_10: crit, GS_13: 'justificación', GS_15: 'true', GS_16: 'true', GS_17: 'true', GS_14: 'true' } });
      await api('POST', `${icode}/entidades`, { kind: 'activacion', entityBranchId: 102, otherInfo: 'obj', agreement: true, dueDiligence: true, activationChannels: ['notificacion'], managementNotes: 'x', pendingActions: ['4cd330bd-bb75-44a2-9bc1-d2a50da9d63b'] });
      await api('POST', `${icode}/entidades`, { kind: 'acudida', entityBranchId: 9217, fileNumber: 'RAD-1', reportDate: '2026-10-01', otherInfo: 'acu' });
      await api('POST', `${icode}/barreras`, { sectorCode: 'justicia', entityBranchId: 102, description: 'barrera x', management: ['alerta'], date: '2026-10-21' });
      const r = await api('POST', `${icode}/seccion/AG`, { answers: { AG_7: 'Agresor' } });   // completa la tríada: bo -> in y materializa
      if (!r.json || r.json.triadaCompleta !== true) throw new Error('no se completó la tríada');
    },
  });
  const { page, icode, sql } = t;
  check('el Caso quedó en Iniciado y materializado', sql(`select victim_case_status||'/'||(select count(*) from salvia.entity_case where case_id='${icode}') from salvia.victim_case where victim_case_i_code='${icode}'`), 'in/2');

  await page.getByRole('link', { name: /Plan de acción Salvia/ }).click();
  check('aviso de solo lectura (y no el de borrador)', [(await page.locator('main .alert-info').allInnerTexts()).some((x) => x.includes('ya dejó de ser borrador')), await page.locator('main .alert-warning').count()], [true, 0]);
  check('GS_16 = Sí pero no hay formulario de activación', await page.locator('.card.card-body.mb-3').count(), 0);
  check('se ve la tarjeta de la activación (datos de la tabla)', (await page.locator('.card.card-body.mb-2').first().innerText()).includes('Canales de activación: Notificación SALVIA'), true);
  check('la tarjeta lista la acción pendiente (del snapshot)', (await page.locator('.card.card-body.mb-2').first().innerText()).includes('Pendiente: Recibir la denuncia formalmente'), true);
  check('se ve la barrera de la entidad con la gestión como etiqueta', (await page.locator('table.table-sm').first().innerText()).includes('Alerta por barreras'), true);
  check('sin botón «Eliminar entidad», sin «Eliminar» de barrera y sin «+ Agregar barrera» (GS_17 = Sí)', [await page.getByRole('button', { name: 'Eliminar entidad' }).count(), await page.getByRole('button', { name: 'Eliminar', exact: true }).count(), await page.getByRole('button', { name: '+ Agregar barrera a esta entidad' }).count()], [0, 0, 0]);
  check('remisión: GS_5 (ya respondida) y sus detalles deshabilitados', [await page.locator('input[name="GS_5"]').first().isDisabled(), await page.locator('input[name="GS_5"]').nth(1).isDisabled()], [true, true]);
  const detalles = page.locator('div.mb-2 input[type=checkbox], textarea');
  const habilitadosRemision = await page.evaluate(() => Array.from(document.querySelectorAll('main input[type=checkbox]')).filter((i) => !i.disabled && /equipo|criterio|medida/i.test(i.closest('.mb-2, .ms-4')?.innerText || '')).length);
  check('ninguna casilla de equipos, criterios o medidas queda habilitada', habilitadosRemision, 0);
  check('la justificación de la remisión (GS_13) está deshabilitada', await page.locator('textarea:disabled').count() >= 1, true);
  check('GS_15 y GS_16 siguen respondibles (Finalizar exige GS_16)', [await page.locator('input[name="GS_15"]').first().isEnabled(), await page.locator('input[name="GS_16"]').first().isEnabled()], [true, true]);
  await page.screenshot({ path: `${OUT}/solo-lectura-gs.png` });

  await page.getByRole('link', { name: /Ruta de atención integral/ }).click();
  check('Ruta: aviso de solo lectura, tarjeta acudida y sin formulario ni botones', [(await page.locator('main .alert-info').allInnerTexts()).some((x) => x.includes('ya dejó de ser borrador')), await page.locator('.card.card-body.mb-2').count(), await page.locator('.card.card-body.mb-3').count(), await page.getByRole('button', { name: 'Eliminar entidad' }).count()], [true, 1, 0, 0]);
  check('GS_14 y GS_17 siguen respondibles', [await page.locator('input[name="GS_14"]').first().isEnabled(), await page.locator('input[name="GS_17"]').first().isEnabled()], [true, true]);
  const r409 = await t.api('POST', `${icode}/entidades`, { kind: 'acudida', entityBranchId: 101 });
  check('el backend responde 409 con readOnly', [r409.status, r409.json.readOnly], [409, true]);

  console.log('== B. Un Caso que pasa a Iniciado mientras la pantalla seguía abierta (bandera atrasada)');
  const t2 = await abrir({ hook: async ({ api, icode: ic }) => { await api('POST', `${ic}/seccion/GS`, { answers: { GS_16: 'true' } }); } });
  await t2.page.getByRole('link', { name: /Plan de acción Salvia/ }).click();
  const formulario = t2.page.locator('.card.card-body.mb-3').filter({ hasText: 'Siguiente entidad' });
  check('en borrador se ve el formulario y el aviso de borrador', [await formulario.count(), await t2.page.locator('main .alert-warning').count() >= 1], [1, true]);
  await elegirEntidad(t2.page, formulario, { sector: 'justicia', texto: 'usaqu' });
  await formulario.locator('textarea[placeholder="¿Qué debe gestionar la víctima en esta entidad?"]').fill('Objetivo');
  // Otro actor completa la tríada (el Caso se materializa) antes de que se presione «+ Siguiente entidad».
  await t2.api('POST', `${t2.icode}/seccion/AG`, { answers: { AG_7: 'x' } });
  await formulario.getByRole('button', { name: '+ Siguiente entidad' }).click();
  await t2.page.waitForTimeout(1200);
  check('la pantalla pasa a solo lectura (409): sin formulario y con el aviso', [await t2.page.locator('.card.card-body.mb-3').count(), (await t2.page.locator('main .alert-info').allInnerTexts()).some((x) => x.includes('ya dejó de ser borrador'))], [0, true]);
  check('la entidad no se creó', t2.sql(`select (select count(*) from salvia.entity_case where case_id='${t2.icode}')||'/'||jsonb_array_length(victim_case_v3_draft_data->'entities') from salvia.victim_case where victim_case_i_code='${t2.icode}'`), '0/0');

  console.log('== C. Se completa la tríada desde la pantalla: se recarga y queda en solo lectura');
  const t3 = await abrir({ triada: false });
  await t3.page.getByRole('link', { name: /Ruta de atención integral/ }).click();
  check('antes: aviso de borrador', await t3.page.locator('main .alert-warning').count() >= 1, true);
  // Todas las llaves de la tríada ya respondidas salvo las de «Características del agresor»: esas se contestan por la pantalla.
  t3.sql(`update salvia.victim_case set victim_case_v3_answered_keys='CI_1,CI_2,CI_3,CI_4,CI_5,CI_6,CI_7,CI_9,DP_1,DP_2,DP_5,DP_8,DP_9,DP_11,DP_12,DP_14,DP_17,DP_27,DP_28,DP_29,HV_1,HV_2,HV_4,HV_6,HV_8,HV_9,HV_11,HV_13,HV_14' where victim_case_i_code='${t3.icode}'`);
  await t3.page.getByRole('link', { name: /Características del agresor/ }).click();
  const selects = t3.page.locator('main select');
  for (let i = 0; i < 3; i++) await selects.nth(i).selectOption({ index: 1 });
  await t3.page.locator('input[name="AG_5"]').first().check();      // «¿Depende económicamente del agresor?» (obligatoria)
  await t3.page.getByRole('button', { name: 'Guardar sección' }).click();
  await t3.page.waitForTimeout(1500);
  check('el guardado completó la tríada (el Caso pasó a Iniciado)', t3.sql(`select victim_case_status from salvia.victim_case where victim_case_i_code='${t3.icode}'`), 'in');
  await t3.page.getByRole('link', { name: /Ruta de atención integral/ }).click();
  check('después de guardar la tríada: aviso de solo lectura y sin formulario', [(await t3.page.locator('main .alert-info').allInnerTexts()).some((x) => x.includes('ya dejó de ser borrador')), await t3.page.locator('.card.card-body.mb-3').count()], [true, 0]);

  console.log('== D. Errores de consola');
  check('sin errores de JS (salvo los 409 esperados)', [...t.errores, ...t2.errores, ...t3.errores].filter((e) => !e.includes('409')), []);
  await t.browser.close(); await t2.browser.close(); await t3.browser.close();
  console.log(`\nRESULTADO: ${ok} ok / ${fail} fallos`);
  process.exit(fail ? 1 : 0);
})().catch((e) => { console.error(e); process.exit(2); });
