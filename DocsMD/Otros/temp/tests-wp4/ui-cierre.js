// Pantalla · Cierre: la tabla de actuaciones muestra los canales de activación de cada entidad (REG-20) y la gestión de las barreras.
const { abrir } = require('./ui-lib');
let ok = 0, fail = 0;
const check = (d, a, e) => { if (JSON.stringify(a) === JSON.stringify(e)) { ok++; console.log('  ✅', d); } else { fail++; console.log('  ❌', d, '\n      esperado:', JSON.stringify(e), '\n      obtenido:', JSON.stringify(a)); } };
const TODAS = ['CI:20', 'DP:60', 'CA:20', 'HV:30', 'AG:20', 'TZ:10', 'GS:20'].flatMap((pm) => { const [p, m] = pm.split(':'); return Array.from({ length: +m }, (_, i) => `${p}_${i + 1}`); }).filter((k) => k !== 'GS_2').join(',');

(async () => {
  const t = await abrir({
    triada: false,
    hook: async ({ api, icode, sql }) => {
      await api('POST', `${icode}/entidades`, { kind: 'activacion', entityBranchId: 102, agreement: true, dueDiligence: true, activationChannels: ['notificacion', 'denuncia'], managementNotes: 'x', otherInfo: 'o' });
      await api('POST', `${icode}/barreras`, { sectorCode: 'justicia', entityBranchId: 102, description: 'Barrera del cierre', management: ['activacion'], date: '2026-10-21' });
      // El Caso sale de borrador sin pasar por SaveSection: la lectura materializa lo pendiente (reintento).
      sql(`update salvia.victim_case set victim_case_status='in', victim_case_v3_answered_keys='${TODAS}' where victim_case_i_code='${icode}'`);
    },
  });
  const { page } = t;
  await page.getByRole('link', { name: /Plan de acción Salvia/ }).click();
  await page.getByRole('button', { name: /Finalizar registro y enrutar a Seguimiento/ }).click();
  await page.waitForSelector('text=Registro finalizado.');
  const tabla = await page.locator('table.table-sm tbody tr').allInnerTexts();
  check('«Registro finalizado.»', await page.locator('.alert-success').first().innerText(), 'Registro finalizado.');
  check('la fila de la entidad muestra los canales de activación', tabla.some((f) => f.includes('Entidad') && f.includes('Canales: Notificación SALVIA, Denuncia fácil')), true);
  check('la fila de la barrera muestra su descripción', tabla.some((f) => f.includes('Barrera') && f.includes('Barrera del cierre')), true);
  check('el Caso quedó Terminado', t.sql(`select victim_case_status from salvia.victim_case where victim_case_i_code='${t.icode}'`), 'te');
  check('sin errores de JS', t.errores, []);
  await t.browser.close();
  console.log(`\nRESULTADO: ${ok} ok / ${fail} fallos`);
  process.exit(fail ? 1 : 0);
})().catch((e) => { console.error(e); process.exit(2); });
