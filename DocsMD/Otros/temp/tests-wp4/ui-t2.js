// Pantalla · T2 (REG-18): formulario de barrera de Seguimiento v4 dentro de la tarjeta de entidad.
const { abrir } = require('./ui-lib');
const OUT = process.env.SHOTS || '/tmp/wp4-shots';
require('fs').mkdirSync(OUT, { recursive: true });
let ok = 0, fail = 0;
const check = (d, a, e) => { if (JSON.stringify(a) === JSON.stringify(e)) { ok++; console.log('  ✅', d); } else { fail++; console.log('  ❌', d, '\n      esperado:', JSON.stringify(e), '\n      obtenido:', JSON.stringify(a)); } };

(async () => {
  const t = await abrir({
    hook: async ({ api, icode }) => {
      await api('POST', `${icode}/entidades`, { kind: 'activacion', entityBranchId: 101, fileNumber: 'RAD-1', pendingActions: ['4cd330bd-bb75-44a2-9bc1-d2a50da9d63b'], otherInfo: 'obj' });
      await api('POST', `${icode}/seccion/GS`, { answers: { GS_17: 'true' } });
    },
  });
  const { page, icode, api } = t;
  await page.getByRole('link', { name: /Plan de acción Salvia/ }).click();
  // GS_17 = Sí (guardado por la API: la pregunta vive en la Ruta de atención integral) habilita el botón de la tarjeta.
  await page.getByRole('button', { name: '+ Agregar barrera a esta entidad' }).click();

  const etiquetasVisibles = async () => (await page.locator('.border.rounded.p-2.mt-2 label.form-label').allInnerTexts()).map((s) => s.replace(/\s*\*$/, '').trim());

  console.log('== 1. Sin sector solo se ve la pregunta de sector');
  check('solo «Sector…»', await etiquetasVisibles(), ['Sector en el que se presentó la barrera', 'Fecha en la que se presentó la barrera (aproximada)', 'Funcionario/a o dependencia donde se presentó la barrera', 'Departamento', 'Municipio', 'Vereda o zona', 'Descripción de la barrera (tiempo, modo y lugar)', 'Gestión de la barrera', 'Descripción de la gestión realizada con relación a las barreras']);
  check('«+ Agregar barrera» deshabilitado con obligatorias vacías', await page.getByRole('button', { name: '+ Agregar barrera', exact: true }).isDisabled(), true);
  await page.screenshot({ path: `${OUT}/t2-1-sin-sector.png` });

  console.log('== 2. Sector Justicia muestra sus preguntas');
  const form = page.locator('.border.rounded.p-2.mt-2');
  await form.locator('select').first().selectOption('justicia');
  const et = await etiquetasVisibles();
  check('aparecen «Barreras en justicia», institución y «Otra barrera en sector Justicia»', ['Barreras en justicia', '¿A qué institución acudió?', 'Otra barrera en sector Justicia'].every((x) => et.includes(x)), true);
  check('no aparecen las de Salud ni Protección ni las transversales', ['Barreras en salud', 'Barreras en protección', 'Barreras económicas y socioeconómicas'].some((x) => et.includes(x)), false);
  check('aviso informativo de Justicia', (await form.locator('.alert-info').allInnerTexts()).some((x) => x.startsWith('Justicia:')), true);
  check('13 opciones de «Barreras en justicia»', await form.locator('label.form-label:has-text("Barreras en justicia") + div .form-check').count(), 13);
  await page.screenshot({ path: `${OUT}/t2-2-justicia.png` });

  console.log('== 3. Cambiar a Transversales limpia lo anterior y muestra las 4 listas');
  await form.locator('label.form-label:has-text("Barreras en justicia") + div .form-check input').first().check();
  await form.locator('select').first().selectOption('transversales');
  const et2 = await etiquetasVisibles();
  check('4 listas transversales', ['Barreras institucionales y de talento humano', 'Barreras económicas y socioeconómicas', 'Barreras territoriales y geográficas', 'Barreras por ausencia de enfoque diferencial'].every((x) => et2.includes(x)), true);
  check('desaparecen las de Justicia', et2.includes('Barreras en justicia'), false);
  await form.locator('select').first().selectOption('justicia');
  check('al volver a Justicia la selección anterior se limpió', await form.locator('label.form-label:has-text("Barreras en justicia") + div .form-check input:checked').count(), 0);

  console.log('== 4. Departamento → Municipio');
  check('Municipio deshabilitado sin Departamento', await form.locator('label.form-label:has-text("Municipio") + select').isDisabled(), true);
  const deptoSel = form.locator('label.form-label:has-text("Departamento") + select');
  const optsDepto = await deptoSel.locator('option').allInnerTexts();
  check('hay departamentos', optsDepto.length > 1, true);
  await deptoSel.selectOption({ index: 1 });
  await page.waitForTimeout(600);
  check('Municipio se habilita y carga municipios', (await form.locator('label.form-label:has-text("Municipio") + select option').count()) > 1, true);

  console.log('== 5. Gestión «alerta» pide perspectiva profesional');
  let et3 = await etiquetasVisibles();
  check('sin gestión no hay perspectiva', et3.includes('¿Cuál es su perspectiva profesional frente a la barrera?'), false);
  const gestion = form.locator('label.form-label:has-text("Gestión de la barrera") + div');
  await gestion.locator('.form-check', { hasText: 'Alerta por barreras' }).locator('input').check();
  et3 = await etiquetasVisibles();
  check('con alerta aparece la perspectiva', et3.includes('¿Cuál es su perspectiva profesional frente a la barrera?'), true);
  check('con alerta aparece el aviso', (await form.locator('.alert-info').allInnerTexts()).some((x) => x.startsWith('Alerta por barreras notificaría')), true);
  await page.screenshot({ path: `${OUT}/t2-5-alerta.png` });

  console.log('== 6. Completar y agregar');
  await form.locator('label.form-label:has-text("Barreras en justicia") + div .form-check input').nth(2).check();
  await form.locator('label.form-label:has-text("¿A qué institución acudió?") + select').selectOption({ index: 1 });
  await form.locator('input[type=date]').fill('2026-10-20');
  await form.locator('label.form-label:has-text("Funcionario/a") + input').fill('Fiscal Pérez');
  await form.locator('label.form-label:has-text("Municipio") + select').selectOption({ index: 1 });
  await form.locator('label.form-label:has-text("Vereda o zona") + select').selectOption('rural');
  await form.locator('label.form-label:has-text("Descripción de la barrera") + textarea').fill('No recibieron la denuncia');
  await gestion.locator('.form-check', { hasText: 'Activación de ruta' }).locator('input').check();
  check('sigue deshabilitado: falta la perspectiva y la descripción de la gestión', await page.getByRole('button', { name: '+ Agregar barrera', exact: true }).isDisabled(), true);
  await form.locator('label.form-label:has-text("perspectiva profesional") + textarea').fill('Hay riesgo');
  await form.locator('label.form-label:has-text("Descripción de la gestión realizada") + textarea').fill('Se activó la ruta');
  check('habilitado con todo lo obligatorio', await page.getByRole('button', { name: '+ Agregar barrera', exact: true }).isEnabled(), true);
  await page.screenshot({ path: `${OUT}/t2-6-completo.png` });
  await page.getByRole('button', { name: '+ Agregar barrera', exact: true }).click();
  await page.waitForTimeout(800);
  check('el formulario se cierra', await page.locator('.border.rounded.p-2.mt-2').count(), 0);
  check('la tarjeta lista 1 barrera', await page.locator('table.table-sm tbody tr').count(), 1);

  console.log('== 7. Lo que llegó al backend (borrador JSON)');
  const g = (await api('GET', icode)).json;
  const b = g.barriers[0];
  check('id del borrador', b.id.startsWith('d~'), true);
  check('sector / institución / zona / depto', [b.sector, !!b.specificInstitutions, g.barrierMeta[b.id].zone, !!b.departmentId], ['justicia', true, 'rural', true]);
  check('gestión (CSV) y perspectiva', [b.managementActions, g.barrierMeta[b.id].professionalPerspective], ['alerta,activacion', 'Hay riesgo']);
  check('descripción de la gestión', g.barrierMeta[b.id].managementDescription, 'Se activó la ruta');
  check('sin competencias', b.competences, '');

  console.log('== 8. Errores de consola');
  check('sin errores de JS', t.errores, []);
  await t.browser.close();
  console.log(`\nRESULTADO: ${ok} ok / ${fail} fallos`);
  process.exit(fail ? 1 : 0);
})().catch((e) => { console.error(e); process.exit(2); });
