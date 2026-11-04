// Pantalla · T4 (REG-04 c): la celda «Gestión» de la tabla de barreras muestra etiquetas, no códigos.
const { abrir } = require('./ui-lib');
let ok = 0, fail = 0;
const check = (d, a, e) => { if (JSON.stringify(a) === JSON.stringify(e)) { ok++; console.log('  ✅', d); } else { fail++; console.log('  ❌', d, '\n      esperado:', JSON.stringify(e), '\n      obtenido:', JSON.stringify(a)); } };

const casos = [
  // [descripción, management enviado, texto esperado en la celda]
  ['nueva «activacion»', ['activacion'], 'Activación de ruta interinstitucional (Generar oficio)'],
  ['nueva «alerta»', ['alerta'], 'Alerta por barreras'],
  ['dos gestiones nuevas', ['activacion', 'alerta'], 'Activación de ruta interinstitucional (Generar oficio), Alerta por barreras'],
  ['histórico 01 (orientacion_llamada)', ['01'], 'Orientación y enrutamiento - Llamada'],
  ['histórico 02 (gestion_llamada)', ['02'], 'Gestión administrativa - Llamada'],
  ['histórico 03 (activacion_ruta_interinstitucional)', ['03'], 'Activación de ruta interinstitucional'],
  ['histórico 04 (articulacion_institucional)', ['04'], 'Articulación institucional'],
  ['histórico 05 (escalamiento_organismo_control)', ['05'], 'Escalamiento a organismo de control'],
  ['histórico 06 (alerta_barreras)', ['06'], 'Alerta por barreras'],
  ['valor desconocido se muestra tal cual', ['xyz_desconocido'], 'xyz_desconocido'],
  ['sin gestión', [], '—'],
];

(async () => {
  const t = await abrir({
    hook: async ({ api, icode }) => {
      await api('POST', `${icode}/entidades`, { kind: 'activacion', entityBranchId: 101, otherInfo: 'obj' });
      await api('POST', `${icode}/seccion/GS`, { answers: { GS_17: 'true' } });
      for (const [desc, management] of casos) {
        await api('POST', `${icode}/barreras`, { sectorCode: 'justicia', entityBranchId: 101, description: desc, management, date: '2026-10-21' });
      }
    },
  });
  await t.page.getByRole('link', { name: /Plan de acción Salvia/ }).click();
  const filas = await t.page.locator('table.table-sm tbody tr').evaluateAll((trs) => trs.map((tr) => Array.from(tr.querySelectorAll('td')).slice(0, 2).map((td) => td.innerText.trim())));
  check('hay una fila por barrera', filas.length, casos.length);
  for (const [desc, , esperado] of casos) {
    const fila = filas.find((f) => f[0] === desc);
    check(desc, fila && fila[1], esperado);
  }
  const celdas = filas.map((f) => f[1]).join(' | ');
  check('ninguna celda deja un slug con guion bajo (salvo el desconocido a propósito)', /activacion_ruta|orientacion_|alerta_barreras|escalamiento_|articulacion_institucional/.test(celdas), false);
  check('sin errores de JS', t.errores, []);
  await t.page.screenshot({ path: '/tmp/wp4-shots/t4-gestion.png', fullPage: false });
  await t.browser.close();
  console.log(`\nRESULTADO: ${ok} ok / ${fail} fallos`);
  process.exit(fail ? 1 : 0);
})().catch((e) => { console.error(e); process.exit(2); });
