// Pantalla · T5 (REG-16): GS_20 «elementos psicosociales», «Ninguno» excluyente, guardar y recargar.
const { abrir } = require('./ui-lib');
const OUT = process.env.SHOTS || '/tmp/wp4-shots';
require('fs').mkdirSync(OUT, { recursive: true });
let ok = 0, fail = 0;
const check = (d, a, e) => { if (JSON.stringify(a) === JSON.stringify(e)) { ok++; console.log('  ✅', d); } else { fail++; console.log('  ❌', d, '\n      esperado:', JSON.stringify(e), '\n      obtenido:', JSON.stringify(a)); } };

const ETIQUETAS = [
  'Contención emocional en crisis',
  'Comprensión del ciclo de la violencia',
  'Identificación y fortalecimiento de recursos personales',
  'Despatologización de experiencias y reacciones asociadas a la violencia',
  'Co-construcción del Plan de Seguridad',
  'Identificación de señales tempranas de riesgo y estrategias de respuesta',
  'Orientación para acceso a servicios, programas o apoyos institucionales que fortalecen autonomía',
  'Fortalecimiento de la autonomía económica y vinculación a proyecto de vida',
  'Cuestionamiento de mandatos culturales que sostienen la violencia',
  'Identificación de redes seguras y afirmativas (familiares, comunitarias, institucionales)',
  'Reconocimiento y comprensión de los derechos de las personas víctimas de VBG y VpP',
  'Orientación sobre acciones disponibles para la garantía de derechos y la responsabilidad del Estado',
  'Fortalecimiento de la corresponsabilidad de la persona en la gestión de su seguridad',
  'Sensibilización y orientación sobre servicios de salud mental para acompañar los impactos de la violencia',
  'Ninguno',
];

(async () => {
  const t = await abrir();
  const { page } = t;
  await page.getByRole('link', { name: /Gestión Salvia/ }).click();
  const bloque = page.locator('div.mb-3', { has: page.locator('label.form-label', { hasText: 'elementos psicosociales integrados' }) });
  const opcion = (txt) => bloque.locator('.form-check', { hasText: txt }).locator('input[type=checkbox]');

  console.log('== 1. La pregunta y sus 15 opciones');
  check('texto de la pregunta', (await bloque.locator('label.form-label').first().innerText()).trim(), 'Seleccione los elementos psicosociales integrados en este acompañamiento');
  check('15 opciones con el texto exacto del cliente (y «Ninguno» al final)', (await bloque.locator('.form-check label').allInnerTexts()).map((x) => x.trim()), ETIQUETAS);
  await bloque.scrollIntoViewIfNeeded();
  await page.screenshot({ path: `${OUT}/t5-1-opciones.png` });

  console.log('== 2. «Ninguno» es excluyente en ambos sentidos');
  await opcion('Contención emocional en crisis').check();
  await opcion('Co-construcción del Plan de Seguridad').check();
  check('dos marcadas', await bloque.locator('input:checked').count(), 2);
  await opcion('Ninguno').check();
  check('al marcar «Ninguno» solo queda esa', [await bloque.locator('input:checked').count(), await opcion('Ninguno').isChecked()], [1, true]);
  await opcion('Comprensión del ciclo de la violencia').check();
  check('al marcar otra se quita «Ninguno»', [await bloque.locator('input:checked').count(), await opcion('Ninguno').isChecked(), await opcion('Comprensión del ciclo de la violencia').isChecked()], [1, false, true]);
  await opcion('Sensibilización y orientación sobre servicios de salud mental').check();

  console.log('== 3. Otros campos múltiples no cambian (sin exclusividad)');
  const otrosMultiples = await page.locator('input[type=checkbox]').count();
  check('hay otras casillas en el módulo (remisión, etc.) sin efecto', otrosMultiples >= 15, true);

  console.log('== 4. Guardar la sección y recargar');
  await page.locator('input[name="GS_1"]').first().check();
  await page.locator('div.mb-3', { has: page.locator('label.form-label', { hasText: 'acompañamiento realizado' }) }).locator('textarea').fill('Acompañamiento de prueba');
  await page.getByRole('button', { name: 'Guardar sección' }).click();
  await page.waitForTimeout(900);
  check('mensaje de guardado', (await page.locator('span.ms-2').allInnerTexts()).some((x) => x.includes('Sección guardada')), true);
  const guardado = t.sql(`select count(*) from salvia.rel_victim_case_form2_enums_victim_case_form2 r join salvia.victim_case_form2_enums e on e.victim_case_form2_enums_id=r.victim_case_form2_enums_id join salvia.victim_case_form2 f on f.victim_case_form2_id=r.victim_case_form2_id join salvia.victim_case vc on vc.victim_case_id=f.victim_case_form2_victim_case where vc.victim_case_i_code='${t.icode}' and e.victim_case_form2_enums_category='victim_case_form2_psychosocial_elements'`);
  check('en BD: 2 elementos', guardado, '2');
  await page.reload();
  await page.getByRole('button', { name: 'Acepto y continuar' }).click();
  await page.getByRole('button', { name: 'Víctima - Caso Nuevo' }).click();
  await page.getByPlaceholder('Número de documento').fill(t.doc);
  await page.getByRole('button', { name: 'Buscar' }).click();
  await page.locator('table button', { hasText: 'Abrir' }).first().click();
  await page.waitForSelector('ul.nav-tabs');
  await page.getByRole('link', { name: /Gestión Salvia/ }).click();
  const bloque2 = page.locator('div.mb-3', { has: page.locator('label.form-label', { hasText: 'elementos psicosociales integrados' }) });
  const marcadas = await bloque2.locator('.form-check', { has: page.locator('input:checked') }).locator('label').allInnerTexts();
  check('al reabrir el Caso siguen marcadas las 2', marcadas.map((x) => x.trim()).sort(), ['Comprensión del ciclo de la violencia', 'Sensibilización y orientación sobre servicios de salud mental para acompañar los impactos de la violencia'].sort());

  console.log('== 5. Errores de consola');
  check('sin errores de JS', t.errores, []);
  await t.browser.close();
  console.log(`\nRESULTADO: ${ok} ok / ${fail} fallos`);
  process.exit(fail ? 1 : 0);
})().catch((e) => { console.error(e); process.exit(2); });
