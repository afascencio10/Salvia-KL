/**
 * migrate_users.js
 * ─────────────────────────────────────────────────────────────────────────────
 * Script de migración masiva de usuarios de Salvia.
 *
 * Lee el archivo Excel (hoja "Usuarios en BD") y llama a los endpoints:
 *   POST /api/v1/admin/migrate/users       → crea o actualiza usuarios reales
 *   POST /api/v1/admin/migrate/test-users  → crea credenciales temporales de prueba
 *
 * Uso:
 *   node migrate_users.js [--dry-run] [--only-test] [--url http://localhost:8080]
 *
 * Flags:
 *   --dry-run    Muestra qué haría pero NO llama a la API
 *   --only-test  Solo ejecuta el endpoint de test-users (omite migrate/users)
 *   --only-prod  Solo ejecuta el endpoint de migrate/users (omite test-users)
 *   --url        URL base del servidor (default: https://localhost)
 *
 * Dependencias: npm install xlsx node-fetch (o usa Node 18+ fetch nativo)
 * ─────────────────────────────────────────────────────────────────────────────
 */

'use strict';

const path = require('path');
const fs   = require('fs');
const XLSX = require('xlsx');

// ─── CONFIGURACIÓN ────────────────────────────────────────────────────────────

const CONFIG = {
  // Ruta al Excel — ajusta esta ruta a donde tengas el archivo
  excelPath: process.env.EXCEL_PATH ||
    path.join('C:\\Users\\braya\\Downloads\\Caricatura assets\\Videos', 'Reporte_Usuarios_Salvia.xlsx'),

  // Nombre de la hoja con los usuarios finales
  sheetName: 'Usuarios en BD',

  // URL base del servidor Salvia
  baseUrl: process.env.SALVIA_URL || 'https://localhost',

  // API Key (X-Security-Key)
  securityKey: process.env.SETUP_API_KEY || 's4lv1a_kr31vo',

  // Contraseña por defecto para usuarios nuevos en producción
  // (se usará SOLO si el usuario no existe — en update se ignora)
  defaultPassword: process.env.DEFAULT_PASSWORD || 'Salvia@Prod2026!',

  // Contraseña estándar para credenciales de prueba
  testPassword: process.env.TEST_PASSWORD || 'Salvia@Test2026!',

  // Tamaño del lote enviado por petición (para no sobrecargar la API)
  batchSize: 20,

  // Pausa entre lotes (ms)
  batchDelay: 500,
};

// ─── LECTURA DE FLAGS ─────────────────────────────────────────────────────────

const args = process.argv.slice(2);
const DRY_RUN   = args.includes('--dry-run');
const ONLY_TEST = args.includes('--only-test');
const ONLY_PROD = args.includes('--only-prod');

const urlIdx = args.indexOf('--url');
if (urlIdx !== -1 && args[urlIdx + 1]) CONFIG.baseUrl = args[urlIdx + 1];

// ─── MAPEO DE COLUMNAS DEL EXCEL ─────────────────────────────────────────────

/**
 * Convierte una fila del Excel en el formato esperado por la API.
 * Aplica la regla de roles: sv → ["sv"], cualquier otro → ["ro"]
 *
 * Columnas del Excel usadas:
 *  - "Login BD"       → login
 *  - "Nombre en BD"   → names + lastNames (split por primer espacio)
 *  - "RolKreivo"      → rol del sistema (sv o cualquier otro)
 *  - "Equipo"         → team
 *  - "Email(s) en BD" → email principal
 *  - "Correo (Excel)" → email alternativo (si no hay email en BD)
 *  - "Doc. Número"    → docNumber
 *  - "Doc. Tipo"      → docType (normalizado a mayúsculas)
 *  - "Estado BD"      → solo se procesan los activos ("e") y los nuevos
 */
function rowToUser(row) {
  const login = (row['Login BD'] || '').toString().trim();

  // Ignorar filas sin login o con login "Crear" (sin asignar)
  if (!login || login.toLowerCase() === 'crear') return null;

  // Nombre completo → names + lastNames
  const fullName  = (row['Nombre en BD'] || row['Nombre (Excel)'] || '').toString().trim();
  const nameParts = fullName.split(' ');
  const names     = nameParts.slice(0, 2).join(' ');
  const lastNames = nameParts.slice(2).join(' ') || names;

  // Email: preferir el de la BD, si no el del Excel
  const email = (
    (row['Email(s) en BD'] || '').toString().trim() ||
    (row['Correo (Excel)'] || '').toString().trim()
  );

  // Rol: regla sv → sv, cualquier otro → ro
  const rolKreivo = (row['RolKreivo'] || '').toString().trim();
  const roleCodes = rolKreivo === 'sv' ? ['sv'] : ['ro'];

  // Team: normalizar mayúsculas/minúsculas
  const teamRaw = (row['Equipo'] || '').toString().trim();
  const team    = normalizeTeam(teamRaw);

  // Documento
  const docNumber = (row['Doc. Número'] || '').toString().trim();
  const docType   = (row['Doc. Tipo'] || 'CC').toString().trim().toUpperCase();

  return {
    login,
    names:     names     || login,
    lastNames: lastNames || login,
    email:     email     || `${login}@salvia-temp.com`,
    phone:     '0000000000',
    docType,
    docNumber: docNumber || `AUTO_${login}`,
    gender:    'ma',       // género por defecto — el Excel no tiene este campo
    lang:      'sp',
    townCode:  '11001000', // Bogotá por defecto
    team,
    roleCodes,
    pass:      CONFIG.defaultPassword,
  };
}

/**
 * Normaliza el nombre del equipo del Excel al valor exacto del sistema.
 * "Riesgo Alto" → "Riesgo alto"
 * "Riesgo Bajo" → "Riesgo bajo"
 * "Hombres"     → "Hombres"
 * vacío / otro  → ""
 */
function normalizeTeam(raw) {
  const lower = raw.toLowerCase().trim();
  if (lower === 'riesgo alto') return 'Riesgo alto';
  if (lower === 'riesgo bajo') return 'Riesgo bajo';
  if (lower === 'hombres')     return 'Hombres';
  return '';
}

// ─── LECTURA DEL EXCEL ────────────────────────────────────────────────────────

function readExcel() {
  if (!fs.existsSync(CONFIG.excelPath)) {
    console.error(`❌ No se encontró el Excel en: ${CONFIG.excelPath}`);
    console.error(`   Ajusta CONFIG.excelPath o coloca el Excel en la ruta correcta.`);
    process.exit(1);
  }

  const wb   = XLSX.readFile(CONFIG.excelPath);
  const ws   = wb.Sheets[CONFIG.sheetName];
  if (!ws) {
    console.error(`❌ No se encontró la hoja "${CONFIG.sheetName}" en el Excel.`);
    console.error(`   Hojas disponibles: ${wb.SheetNames.join(', ')}`);
    process.exit(1);
  }

  const rows = XLSX.utils.sheet_to_json(ws);
  const users = [];
  const skipped = [];

  for (const row of rows) {
    const user = rowToUser(row);
    if (user) {
      users.push(user);
    } else {
      skipped.push(row['Login BD'] || row['Nombre (Excel)'] || '(sin login)');
    }
  }

  return { users, skipped };
}

// ─── LLAMADAS A LA API ────────────────────────────────────────────────────────

async function callAPI(endpoint, body) {
  const url  = `${CONFIG.baseUrl}${endpoint}`;
  const opts = {
    method:  'POST',
    headers: {
      'Content-Type':   'application/json',
      'X-Security-Key': CONFIG.securityKey,
    },
    body: JSON.stringify(body),
  };

  // Node 18+ tiene fetch nativo; para versiones anteriores instala node-fetch
  let fetchFn;
  try {
    fetchFn = fetch; // nativo
  } catch {
    fetchFn = require('node-fetch');
  }

  // Ignorar certificados autofirmados en local (no usar en producción real)
  if (CONFIG.baseUrl.startsWith('https://localhost') || CONFIG.baseUrl.includes('127.0.0.1')) {
    process.env.NODE_TLS_REJECT_UNAUTHORIZED = '0';
  }

  const res  = await fetchFn(url, opts);
  const text = await res.text();

  try {
    return { status: res.status, data: JSON.parse(text) };
  } catch {
    return { status: res.status, data: { raw: text } };
  }
}

function sleep(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

// ─── PROCESAMIENTO EN LOTES ───────────────────────────────────────────────────

async function runMigrate(users) {
  console.log(`\n📤 Enviando ${users.length} usuarios a POST /api/v1/admin/migrate/users`);
  console.log(`   Lotes de ${CONFIG.batchSize} usuarios con ${CONFIG.batchDelay}ms de pausa\n`);

  const totals = { created: 0, updated: 0, skipped: 0, errors: [] };

  for (let i = 0; i < users.length; i += CONFIG.batchSize) {
    const batch = users.slice(i, i + CONFIG.batchSize);
    const batchNum = Math.floor(i / CONFIG.batchSize) + 1;
    const totalBatches = Math.ceil(users.length / CONFIG.batchSize);

    process.stdout.write(`   Lote ${batchNum}/${totalBatches} (${batch.length} usuarios)... `);

    if (DRY_RUN) {
      console.log('⏭ [dry-run]');
      batch.forEach(u => console.log(`     → ${u.login} [${u.roleCodes}] team:${u.team || 'sin'}`));
      continue;
    }

    try {
      const { status, data } = await callAPI('/api/v1/admin/migrate/users', { users: batch });

      if (status === 200) {
        totals.created += data.created_count || 0;
        totals.updated += data.updated_count || 0;
        totals.skipped += data.skipped_count || 0;
        if (data.skipped && data.skipped.length > 0) {
          totals.errors.push(...data.skipped.map(s => `${s.login}: ${s.reason}`));
        }
        console.log(`✅ creados:${data.created_count} actualizados:${data.updated_count} omitidos:${data.skipped_count}`);
      } else {
        console.log(`❌ HTTP ${status}:`, JSON.stringify(data));
      }
    } catch (err) {
      console.log(`❌ Error de red:`, err.message);
    }

    if (i + CONFIG.batchSize < users.length) {
      await sleep(CONFIG.batchDelay);
    }
  }

  return totals;
}

async function runTestMigrate(users) {
  console.log(`\n🧪 Creando credenciales de prueba para ${users.length} usuarios`);
  console.log(`   Endpoint: POST /api/v1/admin/migrate/test-users`);
  console.log(`   Contraseña de prueba: ${CONFIG.testPassword}\n`);

  const totals = { created: 0, skipped: 0, errors: [] };

  for (let i = 0; i < users.length; i += CONFIG.batchSize) {
    const batch = users.slice(i, i + CONFIG.batchSize);
    const batchNum = Math.floor(i / CONFIG.batchSize) + 1;
    const totalBatches = Math.ceil(users.length / CONFIG.batchSize);

    process.stdout.write(`   Lote ${batchNum}/${totalBatches} (${batch.length} usuarios)... `);

    if (DRY_RUN) {
      console.log('⏭ [dry-run]');
      batch.forEach(u => console.log(`     → test.${u.login}`));
      continue;
    }

    try {
      const { status, data } = await callAPI('/api/v1/admin/migrate/test-users', {
        testPassword: CONFIG.testPassword,
        users: batch,
      });

      if (status === 200) {
        totals.created += data.created_count || 0;
        totals.skipped += data.skipped_count || 0;
        if (data.skipped && data.skipped.length > 0) {
          totals.errors.push(...data.skipped.map(s => `${s.login}: ${s.reason}`));
        }
        console.log(`✅ creados:${data.created_count} omitidos:${data.skipped_count}`);
      } else {
        console.log(`❌ HTTP ${status}:`, JSON.stringify(data));
      }
    } catch (err) {
      console.log(`❌ Error de red:`, err.message);
    }

    if (i + CONFIG.batchSize < users.length) {
      await sleep(CONFIG.batchDelay);
    }
  }

  return totals;
}

// ─── MAIN ─────────────────────────────────────────────────────────────────────

async function main() {
  console.log('════════════════════════════════════════════════════════');
  console.log('  Salvia — Script de Migración de Usuarios');
  console.log('════════════════════════════════════════════════════════');
  console.log(`  Servidor  : ${CONFIG.baseUrl}`);
  console.log(`  Excel     : ${CONFIG.excelPath}`);
  console.log(`  Hoja      : ${CONFIG.sheetName}`);
  console.log(`  Dry-run   : ${DRY_RUN ? 'SÍ — no se hará ningún cambio' : 'NO'}`);
  console.log('────────────────────────────────────────────────────────\n');

  // 1. Leer Excel
  const { users, skipped } = readExcel();
  console.log(`📊 Excel leído: ${users.length} usuarios válidos, ${skipped.length} omitidos`);
  if (skipped.length > 0) {
    console.log('   Omitidos (sin login o con "Crear"):');
    skipped.forEach(s => console.log(`     - ${s}`));
  }

  // 2. Mostrar resumen de lo que se enviará
  const svUsers = users.filter(u => u.roleCodes[0] === 'sv');
  const roUsers = users.filter(u => u.roleCodes[0] === 'ro');
  const conTeam = users.filter(u => u.team !== '');
  console.log(`\n📋 Distribución:`);
  console.log(`   Supervisores (sv): ${svUsers.length}`);
  console.log(`   Operadores   (ro): ${roUsers.length}`);
  console.log(`   Con equipo       : ${conTeam.length}`);

  if (users.length === 0) {
    console.log('\n⚠️  No hay usuarios para procesar. Revisa el Excel.');
    return;
  }

  // 3. Migración de usuarios reales
  if (!ONLY_TEST) {
    const prodResult = await runMigrate(users);
    if (!DRY_RUN) {
      console.log('\n📊 Resultado migrate/users:');
      console.log(`   Creados   : ${prodResult.created}`);
      console.log(`   Actualizados: ${prodResult.updated}`);
      console.log(`   Omitidos  : ${prodResult.skipped}`);
      if (prodResult.errors.length > 0) {
        console.log('\n   ⚠️ Errores/Advertencias:');
        prodResult.errors.forEach(e => console.log(`     - ${e}`));
      }
    }
  }

  // 4. Credenciales de prueba
  if (!ONLY_PROD) {
    const testResult = await runTestMigrate(users);
    if (!DRY_RUN) {
      console.log('\n📊 Resultado migrate/test-users:');
      console.log(`   Creados  : ${testResult.created}`);
      console.log(`   Omitidos : ${testResult.skipped}`);
      if (testResult.errors.length > 0) {
        console.log('\n   ⚠️ Errores/Advertencias:');
        testResult.errors.forEach(e => console.log(`     - ${e}`));
      }
    }
  }

  console.log('\n════════════════════════════════════════════════════════');
  console.log('  ✅ Migración completada');
  console.log('════════════════════════════════════════════════════════\n');
}

main().catch(err => {
  console.error('\n❌ Error fatal:', err.message);
  process.exit(1);
});
