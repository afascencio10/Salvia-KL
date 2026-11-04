// Pantalla · T6 (REG-19): agregar entidad como Gestión institucional (cascada obligatoria, buscar con debounce, chips).
const { abrir } = require('./ui-lib');
const OUT = process.env.SHOTS || '/tmp/wp4-shots';
require('fs').mkdirSync(OUT, { recursive: true });
let ok = 0, fail = 0;
const check = (d, a, e) => { if (JSON.stringify(a) === JSON.stringify(e)) { ok++; console.log('  ✅', d); } else { fail++; console.log('  ❌', d, '\n      esperado:', JSON.stringify(e), '\n      obtenido:', JSON.stringify(a)); } };

(async () => {
  const t = await abrir({ hook: async ({ api, icode }) => { await api('POST', `${icode}/seccion/GS`, { answers: { GS_16: 'true' } }); } });
  const { page, icode, sql } = t;
  await page.getByRole('link', { name: /Plan de acción Salvia/ }).click();
  const peticiones = [];
  page.on('request', (r) => { if (/\/api\/v1\/(entity-branches|locations)/.test(r.url())) peticiones.push(r.url().replace(/^https?:\/\/[^/]+/, '')); });

  const nueva = page.locator('.card.card-body.mb-3').filter({ hasText: 'Siguiente entidad' });
  const sel = (label) => nueva.locator(`label.form-label:has-text("${label}") + select`);
  const buscar = nueva.locator('input[placeholder="Escribe el nombre de la entidad"]');
  const sugerencias = () => nueva.locator('.list-group button').allInnerTexts();

  console.log('== 1. Cascada obligatoria');
  const etiquetas = await nueva.locator('label.form-label').allInnerTexts();
  check('Departamento, Ciudad, Municipio y Sector con asterisco', ['Departamento *', 'Ciudad *', 'Municipio *', 'Sector *'].every((x) => etiquetas.map((e) => e.trim()).includes(x)), true);
  check('Ciudad deshabilitada sin departamento', await sel('Ciudad').isDisabled(), true);
  check('Municipio deshabilitado sin ciudad', await sel('Municipio').isDisabled(), true);
  check('Sector deshabilitado sin municipio', await sel('Sector').isDisabled(), true);
  check('«Buscar entidad» no aparece sin municipio y sector', await buscar.count(), 0);
  await sel('Departamento').selectOption({ label: 'BOGOTÁ, D.C.' });
  await page.waitForTimeout(400);
  await sel('Ciudad').selectOption({ index: 1 });
  await page.waitForTimeout(400);
  await sel('Municipio').selectOption({ index: 1 });
  check('con municipio aún no aparece el buscador (falta el sector)', await buscar.count(), 0);

  console.log('== 2. Una sola petición de sedes al elegir municipio y sector');
  const antes = peticiones.filter((u) => u.startsWith('/api/v1/entity-branches?')).length;
  await sel('Sector').selectOption('justicia');
  await buscar.waitFor();
  const sedes = peticiones.filter((u) => u.startsWith('/api/v1/entity-branches?')).length - antes;
  check('1 petición de sedes', sedes, 1);

  console.log('== 3. Escribir filtra en el navegador, con debounce y sin tildes');
  await buscar.focus();
  check('al enfocar se ofrecen todas las sedes del municipio y sector', (await sugerencias()).length >= 3, true);
  const total = (await sugerencias()).length;
  const antesTeclas = peticiones.length;
  await buscar.pressSequentially('bolivar', { delay: 30 });
  check('justo después de teclear el filtro aún no se aplicó (debounce 250 ms)', (await sugerencias()).length, total);
  await page.waitForTimeout(450);
  const filtradas = await sugerencias();
  check('«bolivar» (sin tilde ni mayúsculas) encuentra «…CIUDAD BOLÍVAR»', filtradas.length === 1 && /BOL[ÍI]VAR/i.test(filtradas[0]), true);
  check('escribir no envió ninguna petición', peticiones.length - antesTeclas, 0);
  await buscar.fill('');
  await buscar.pressSequentially('zzzz', { delay: 20 });
  await page.waitForTimeout(450);
  check('sin coincidencias se avisa', (await nueva.innerText()).includes('Ninguna entidad coincide con lo escrito.'), true);
  await page.screenshot({ path: `${OUT}/t6-1-sin-coincidencias.png` });

  console.log('== 4. Elegir una sugerencia carga las acciones institucionales (chips)');
  await buscar.fill('');
  await buscar.pressSequentially('usaqu', { delay: 20 });
  await page.waitForTimeout(450);
  const antesObl = peticiones.length;
  await nueva.locator('.list-group button').first().dispatchEvent('mousedown');
  await page.waitForTimeout(700);
  check('se pidieron las acciones de esa sede (1 petición)', peticiones.slice(antesObl).filter((u) => u.includes('/obligations')).length, 1);
  const chips = nueva.locator('button[type=button]').filter({ hasText: /Recibir|Asignar|Aplicar|Emitir|Informar/ });
  check('hay chips de acciones institucionales', (await chips.count()) > 0, true);
  check('el campo muestra la entidad elegida', (await buscar.inputValue()).includes('USAQU'), true);
  await page.screenshot({ path: `${OUT}/t6-2-chips.png` });

  console.log('== 5. Marcar chips y agregar');
  await nueva.locator('label.form-label:has-text("Acciones institucionales")').first().scrollIntoViewIfNeeded();
  const pendientes = nueva.locator('label.form-label:has-text("Acciones institucionales") ~ div.d-flex button');
  await pendientes.nth(0).click();
  await pendientes.nth(2).click();
  check('el chip marcado muestra ✓', (await pendientes.nth(0).innerText()).startsWith('✓'), true);
  await nueva.locator('textarea[placeholder="¿Qué debe gestionar la víctima en esta entidad?"]').fill('Objetivo de prueba');
  await nueva.getByRole('button', { name: '+ Siguiente entidad' }).click();
  await page.waitForTimeout(800);
  const json = sql(`select jsonb_array_length(victim_case_v3_draft_data->'entities')||'/'||jsonb_array_length(victim_case_v3_draft_data->'entities'->0->'pending')||'/'||jsonb_array_length(victim_case_v3_draft_data->'entities'->0->'completed') from salvia.victim_case where victim_case_i_code='${icode}'`);
  check('en el borrador: 1 entidad, 2 pendientes y 0 completadas (la activación solo tiene pendientes)', json, '1/2/0');
  check('el formulario se reinició (sin texto de búsqueda)', await nueva.locator('input[placeholder="Escribe el nombre de la entidad"]').count(), 0);

  console.log('== 6. Editar el texto limpia la selección; cambiar el sector limpia la entidad');
  await sel('Departamento').selectOption({ label: 'BOGOTÁ, D.C.' });
  await page.waitForTimeout(400);
  await sel('Ciudad').selectOption({ index: 1 });
  await page.waitForTimeout(400);
  await sel('Municipio').selectOption({ index: 1 });
  await sel('Sector').selectOption('justicia');
  await buscar.waitFor();
  await buscar.focus();
  await buscar.pressSequentially('reaccion', { delay: 20 });
  await page.waitForTimeout(450);
  await nueva.locator('.list-group button').first().dispatchEvent('mousedown');
  await page.waitForTimeout(600);
  check('hay entidad elegida (aparecen las acciones)', await nueva.locator('label.form-label:has-text("Acciones institucionales")').count(), 1);
  await buscar.press('End');
  await buscar.pressSequentially('x', { delay: 20 });
  await page.waitForTimeout(400);
  check('editar el texto limpia la selección (desaparecen las acciones)', await nueva.locator('label.form-label:has-text("Acciones institucionales")').count(), 0);
  await sel('Sector').selectOption('proteccion');
  await page.waitForTimeout(800);
  check('cambiar el sector vacía el campo de búsqueda', await buscar.inputValue(), '');

  console.log('== 7. Una respuesta atrasada no pisa a la nueva');
  await page.route('**/api/v1/entity-branches?*sector=justicia*', async (route) => { await new Promise((r) => setTimeout(r, 1500)); await route.continue(); });
  await sel('Sector').selectOption('justicia');          // petición A (lenta)
  await page.waitForTimeout(100);
  await sel('Sector').selectOption('proteccion');        // petición B (rápida)
  await page.waitForTimeout(2200);                       // llega A después de B
  await buscar.focus();
  const nombres = await sugerencias();
  const esperadoProteccion = sql(`select count(distinct eb.entity_branch_name) from salvia.entity_branch eb join salvia.entity e on e.entity_id=eb.entity_id where eb.entity_branch_town_code='11001000' and e.entity_sector='pt'`);
  check('la lista es la de «protección» (la última), no la de «justicia» que llegó tarde', String(nombres.length), esperadoProteccion);

  console.log('== 8. Errores de consola');
  check('sin errores de JS', t.errores, []);
  await t.browser.close();
  console.log(`\nRESULTADO: ${ok} ok / ${fail} fallos`);
  process.exit(fail ? 1 : 0);
})().catch((e) => { console.error(e); process.exit(2); });
