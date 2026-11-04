// Pantalla · T7 (REG-15): «Plan de acción Salvia» (nombre, orden a–i, sin títulos A–D) y formulario de activación.
const { abrir, elegirEntidad } = require('./ui-lib');
const OUT = process.env.SHOTS || '/tmp/wp4-shots';
require('fs').mkdirSync(OUT, { recursive: true });
let ok = 0, fail = 0;
const check = (d, a, e) => { if (JSON.stringify(a) === JSON.stringify(e)) { ok++; console.log('  ✅', d); } else { fail++; console.log('  ❌', d, '\n      esperado:', JSON.stringify(e), '\n      obtenido:', JSON.stringify(a)); } };

(async () => {
  const t = await abrir({
    hook: async ({ api, icode }) => {
      // Una activación y una acudida ya guardadas (la acudida solo debe verse en la Ruta de atención integral).
      await api('POST', `${icode}/entidades`, { kind: 'activacion', entityBranchId: 102, pendingActions: [], otherInfo: 'obj act' });
      await api('POST', `${icode}/entidades`, { kind: 'acudida', entityBranchId: 9217, otherInfo: 'obj acu' });
    },
  });
  const { page, icode, sql } = t;

  console.log('== 1. Nombre de la pestaña');
  const tabs = (await page.locator('ul.nav-tabs a').allInnerTexts()).map((x) => x.replace(/🔒/, '').trim());
  check('la pestaña se llama «Plan de acción Salvia» y ya no «Gestión Salvia»', [tabs.includes('Plan de acción Salvia'), tabs.includes('Gestión Salvia')], [true, false]);
  await page.getByRole('link', { name: /Plan de acción Salvia/ }).click();

  console.log('== 2. Orden en el DOM: a b c d e f g h i y § E');
  const orden = await page.evaluate(() => {
    const raiz = document.querySelector('main .container');
    const textos = (re) => { const w = document.createTreeWalker(raiz, NodeFilter.SHOW_TEXT); const out = []; let n; while ((n = w.nextNode())) { if (re.test(n.nodeValue)) out.push(n); } return out; };
    const pos = (re) => { const ns = textos(re); if (!ns.length) return -1; const r = document.createRange(); r.selectNode(ns[0]); return r.getBoundingClientRect().top + window.scrollY; };
    return {
      sintesis: pos(/Párrafo síntesis/), gs1: pos(/denuncia fácil ante la FGN/), gs15: pos(/Está de acuerdo con que SALVIA realice la activación/),
      gs16: pos(/debe activar rutas en el marco de la debida diligencia\?/), gs3: pos(/acciones previstas, las entidades involucradas/),
      remision: pos(/Requiere remisión|¿A qué equipo|Justificación de la remisión/), gs20: pos(/elementos psicosociales integrados/),
      gs4: pos(/acompañamiento realizado/), alerta: pos(/E\. Alerta Inminente/),
    };
  });
  console.log('   posiciones:', JSON.stringify(orden));
  const claves = ['gs1', 'gs15', 'gs16', 'gs3', 'gs20', 'gs4', 'alerta'];
  const posiciones = claves.map((k) => orden[k]);
  check('todos los elementos existen', posiciones.every((p) => p >= 0), true);
  check('el orden vertical es b, c, d, e, h, i y § E', posiciones.every((p, i) => i === 0 || p > posiciones[i - 1]), true);
  check('el párrafo síntesis no aparece sin datos (se conserva la lógica)', orden.sintesis === -1 || orden.sintesis < orden.gs1, true);

  console.log('== 3. Sin títulos A–D; GS_14 y GS_17 salen del Plan de acción');
  const titulos = (await page.locator('main h5').allInnerTexts()).map((x) => x.trim());
  check('solo queda el título «E. Alerta Inminente»', titulos, ['E. Alerta Inminente']);
  check('ya no está «Denuncia y gestión general»', (await page.locator('main').innerText()).includes('Denuncia y gestión general'), false);
  check('GS_14 y GS_17 no están en el Plan de acción', [await page.locator('input[name="GS_14"]').count(), await page.locator('input[name="GS_17"]').count()], [0, 0]);
  check('GS_2 (gestiones interinstitucionales) no está', (await page.locator('main').innerText()).includes('gestiones interinstitucionales'), false);

  console.log('== 4. GS_15 siempre visible; GS_16 despliega el formulario');
  check('GS_15 visible sin responder GS_14 ni GS_16', await page.locator('input[name="GS_15"]').count(), 2);
  check('sin GS_16 en «Sí» no hay formulario de activación', await page.locator('.card.card-body.mb-3').filter({ hasText: 'Siguiente entidad' }).count(), 0);
  check('las activaciones guardadas se ven siempre; la acudida no', (await page.locator('.card.card-body.mb-2').allInnerTexts()).length, 1);
  await page.locator('input[name="GS_16"]').first().check();
  const form = page.locator('.card.card-body.mb-3').filter({ hasText: 'Siguiente entidad' });
  check('con GS_16 = Sí aparece «Activación de ruta»', await form.count(), 1);
  check('el formulario no tiene selector de tipo', await form.locator('label.form-label:has-text("Tipo")').count(), 0);

  console.log('== 5. Formulario de activación: campos y obligatoriedad');
  await elegirEntidad(page, form, { sector: 'justicia', texto: 'usaqu' });
  const etiquetas = (await form.locator('label.form-label').allInnerTexts()).map((x) => x.replace(/\s*\*$/, '').trim());
  check('pide solicitud, radicado y objetivo', ['Solicitud a la entidad', '¿Cuenta con un N° de radicado, NUNC, MP, etc.?', 'Objetivo con la entidad'].every((x) => etiquetas.includes(x)), true);
  check('«+ Siguiente entidad» deshabilitado sin objetivo', await form.getByRole('button', { name: '+ Siguiente entidad' }).isDisabled(), true);
  const chips = form.locator('label.form-label:has-text("Acciones institucionales") ~ div.d-flex button');
  check('hay chips de acciones institucionales (solo pendientes)', (await chips.count()) > 0 && (await form.locator('label:has-text("Completadas")').count()) === 0, true);
  await chips.nth(0).click(); await chips.nth(1).click();
  await form.locator('textarea[placeholder="¿Qué acciones institucionales se requieren?"]').fill('Solicito atención prioritaria');
  await form.locator('input[placeholder="N° de radicado, NUNC, MP, etc."]').fill('NUNC-123');
  await form.locator('textarea[placeholder="¿Qué debe gestionar la víctima en esta entidad?"]').fill('Que radique la denuncia');
  check('con el objetivo se habilita', await form.getByRole('button', { name: '+ Siguiente entidad' }).isEnabled(), true);
  await page.screenshot({ path: `${OUT}/t7-1-activacion.png` });
  await form.getByRole('button', { name: '+ Siguiente entidad' }).click();
  await page.waitForTimeout(900);

  console.log('== 6. Lo guardado: borrador JSON y tarjeta con sus acciones');
  const fila = sql(`select (victim_case_v3_draft_data->'entities'->2->>'kind')||'|'||(victim_case_v3_draft_data->'entities'->2->>'solicitud')||'|'||(victim_case_v3_draft_data->'entities'->2->>'fileNumber')||'|'||(victim_case_v3_draft_data->'entities'->2->>'objetivo')||'|'||jsonb_array_length(victim_case_v3_draft_data->'entities'->2->'pending') from salvia.victim_case where victim_case_i_code='${icode}'`);
  check('kind, solicitud, radicado, objetivo y 2 acciones pendientes', fila, 'activacion|Solicito atención prioritaria|NUNC-123|Que radique la denuncia|2');
  check('no hay filas en entity_case (sigue siendo borrador)', sql(`select count(*) from salvia.entity_case where case_id='${icode}'`), '0');
  await page.waitForTimeout(500);
  const tarjetas = await page.locator('.card.card-body.mb-2').allInnerTexts();
  check('ahora hay 2 tarjetas de activación', tarjetas.length, 2);
  check('la tarjeta nueva lista sus 2 acciones como «Pendiente: …»', (tarjetas.join('\n').match(/Pendiente: /g) || []).length, 2);

  console.log('== 7. Remisión sin cambios en borrador y aviso');
  check('aviso de borrador visible', (await page.locator('main .alert-warning').allInnerTexts()).some((x) => x.includes('Borrador:')), true);
  check('los campos de remisión están habilitados en borrador', await page.locator('input[name="GS_5"]').first().isEnabled(), true);

  console.log('== 8. Errores de consola');
  check('sin errores de JS', t.errores, []);
  await t.browser.close();
  console.log(`\nRESULTADO: ${ok} ok / ${fail} fallos`);
  process.exit(fail ? 1 : 0);
})().catch((e) => { console.error(e); process.exit(2); });
