// Pantalla · T3 (REG-20): canal de activación como Seguimiento v4 en el formulario de entidad.
const { abrir, elegirEntidad } = require('./ui-lib');
const OUT = process.env.SHOTS || '/tmp/wp4-shots';
require('fs').mkdirSync(OUT, { recursive: true });
let ok = 0, fail = 0;
const check = (d, a, e) => { if (JSON.stringify(a) === JSON.stringify(e)) { ok++; console.log('  ✅', d); } else { fail++; console.log('  ❌', d, '\n      esperado:', JSON.stringify(e), '\n      obtenido:', JSON.stringify(a)); } };

(async () => {
  const t = await abrir({
    hook: async ({ api, icode, sql }) => {
      // Entidad con el canal en el formato anterior del cliente (códigos no/ll): debe verse con las etiquetas nuevas.
      await api('POST', `${icode}/entidades`, { kind: 'activacion', entityBranchId: 102, activationChannel: 'no,ll', agreement: true, dueDiligence: true, managementNotes: 'x' });
      // Dato histórico ya materializado: solo «Canal de activación: ll» en solicitud, sin meta.
      sql(`insert into salvia.entity_case (id, case_id, entity_branch_id, solicitud, created_by_id, created_at, updated_at) values (gen_random_uuid(), '${icode}', 9217, 'Canal de activación: ll', 'test', now(), now())`);
    },
  });
  const { page, icode, api } = t;
  await page.getByRole('link', { name: /Plan de acción Salvia/ }).click();

  console.log('== 1. Tarjetas: canal en formato anterior y dato histórico');
  const tarjetas = await page.locator('.card.card-body.mb-2').allInnerTexts();
  check('la entidad con «no,ll» muestra las etiquetas de Seguimiento', tarjetas.some((x) => x.includes('Canales de activación: Notificación SALVIA, Articulación institucional telefónica o por mensajería')), true);
  check('el dato histórico «Canal de activación: ll» se sigue viendo', tarjetas.some((x) => x.includes('Canales de activación: Llamada (tarea)')), true);

  console.log('== 2. Formulario: sin GS_15 y GS_16 en «Sí» no se piden canales');
  await page.locator('input[name="GS_16"]').first().check();          // despliega el formulario de activación
  const nueva = page.locator('.card.card-body.mb-3').filter({ hasText: 'Siguiente entidad' });
  await elegirEntidad(page, nueva, { sector: 'justicia', texto: 'usaqu' });
  check('aparece la nota explicativa', (await nueva.innerText()).includes('Los canales de activación se piden cuando la persona está de acuerdo'), true);
  check('no hay casillas de canal', await nueva.locator('.form-check-input[type=checkbox]').filter({ hasNot: page.locator('x') }).count() >= 0, true);
  await nueva.locator('textarea[placeholder="¿Qué debe gestionar la víctima en esta entidad?"]').fill('Objetivo de prueba');
  check('«+ Siguiente entidad» habilitado con el objetivo (sin canales no hay nada más obligatorio)', await nueva.getByRole('button', { name: '+ Siguiente entidad' }).isEnabled(), true);

  console.log('== 3. Con GS_15 y GS_16 en «Sí» aparecen las 5 opciones y las gestiones');
  await page.locator('input[name="GS_15"]').first().check();
  await page.waitForTimeout(300);
  const etiquetasCanal = await nueva.locator('label.form-label:has-text("Seleccione los canales") + .form-check, label.form-label:has-text("Seleccione los canales") ~ .form-check').locator('label').allInnerTexts();
  check('opciones de canal', etiquetasCanal.map((x) => x.trim()), ['Notificación SALVIA', 'Articulación institucional telefónica o por mensajería', 'Denuncia fácil', 'Mi Salvia', 'Ninguna']);
  check('pregunta de gestiones visible', await nueva.locator('label.form-label:has-text("¿Qué gestiones se realizaron desde Salvia para la activación de ruta?")').count(), 1);
  check('«+ Siguiente entidad» deshabilitado sin canales ni gestiones', await nueva.getByRole('button', { name: '+ Siguiente entidad' }).isDisabled(), true);
  await page.screenshot({ path: `${OUT}/t3-1-canales.png` });

  console.log('== 4. «Ninguna» es excluyente');
  const casilla = (txt) => nueva.locator('.form-check', { hasText: txt }).locator('input[type=checkbox]');
  await casilla('Notificación SALVIA').check();
  await casilla('Ninguna').check();
  check('al marcar «Ninguna» se desmarca Notificación', [await casilla('Notificación SALVIA').isChecked(), await casilla('Ninguna').isChecked()], [false, true]);
  await casilla('Mi Salvia').check();
  check('al marcar Mi Salvia se desmarca «Ninguna»', [await casilla('Mi Salvia').isChecked(), await casilla('Ninguna').isChecked()], [true, false]);

  console.log('== 5. Agregar con canales y gestiones');
  await casilla('Notificación SALVIA').check();
  await casilla('Denuncia fácil').check();
  await casilla('Mi Salvia').uncheck();
  check('sigue deshabilitado sin las gestiones', await nueva.getByRole('button', { name: '+ Siguiente entidad' }).isDisabled(), true);
  await nueva.locator('textarea').last().fill('Se llamó a la entidad y se radicó oficio');
  check('se habilita con canales y gestiones', await nueva.getByRole('button', { name: '+ Siguiente entidad' }).isEnabled(), true);
  await page.screenshot({ path: `${OUT}/t3-2-completo.png` });
  await nueva.getByRole('button', { name: '+ Siguiente entidad' }).click();
  await page.waitForTimeout(800);

  console.log('== 6. Lo que llegó al backend');
  const g = (await api('GET', icode)).json;
  const nuevaEnt = g.entities.find((e) => e.entityBranchId !== 102 && e.entityBranchId !== 9217);
  const meta = g.entityMeta[String(nuevaEnt.entityBranchId)];
  check('canales en entityMeta', meta.activationChannels, ['notificacion', 'denuncia']);
  check('acuerdo y debida diligencia copiados de GS_15 y GS_16', [meta.agreement, meta.dueDiligence], [true, true]);
  check('gestiones realizadas', meta.managementNotes, 'Se llamó a la entidad y se radicó oficio');
  check('tipo en entityMeta', meta.kind, 'activacion');
  check('solicitud ya no lleva la nota «Canal de activación»', nuevaEnt.solicitud || '', '');
  const tarjetas2 = await page.locator('.card.card-body.mb-2').allInnerTexts();
  check('la tarjeta nueva muestra «Notificación SALVIA, Denuncia fácil»', tarjetas2.some((x) => x.includes('Canales de activación: Notificación SALVIA, Denuncia fácil')), true);

  console.log('== 7. Con GS_15 = No no se envían canales');
  await page.locator('input[name="GS_15"]').nth(1).check();
  await elegirEntidad(page, nueva, { sector: 'proteccion', texto: 'comisar' });
  await nueva.locator('textarea[placeholder="¿Qué debe gestionar la víctima en esta entidad?"]').fill('Otro objetivo');
  await nueva.getByRole('button', { name: '+ Siguiente entidad' }).click();
  await page.waitForTimeout(800);
  const g2 = (await api('GET', icode)).json;
  const sinCanal = g2.entities.find((e) => ![102, 9217, nuevaEnt.entityBranchId].includes(e.entityBranchId));
  const meta2 = g2.entityMeta[String(sinCanal.entityBranchId)];
  check('sin canales; Agreement = false (GS_15) y DueDiligence = true (GS_16)', [meta2.activationChannels, meta2.agreement, meta2.dueDiligence], [[], false, true]);

  console.log('== 8. Errores de consola');
  check('sin errores de JS', t.errores, []);
  await t.browser.close();
  console.log(`\nRESULTADO: ${ok} ok / ${fail} fallos`);
  process.exit(fail ? 1 : 0);
})().catch((e) => { console.error(e); process.exit(2); });
