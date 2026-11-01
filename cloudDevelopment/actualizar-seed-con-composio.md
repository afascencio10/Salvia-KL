# Traer datos y esquema de Supabase al sandbox (vía Composio)

Procedimiento para cuando, desarrollando o probando un requerimiento en el
sandbox cloud de Claude Code, la BD local no tiene lo que ese requerimiento
necesita:

- **datos**: un rol, una entidad, una pregunta de formulario, pero también
  **casos, víctimas, seguimientos, barreras, remisiones, tareas o respuestas
  de formularios**;
- **esquema**: una pantalla responde 500 con `column ... does not exist`
  (ver "Esquema" más abajo).

En vez de inventar el dato o adivinarlo leyendo código, se lee la BD de
pruebas real (Supabase) con el conector Composio y se agrega al seed.

## Principio: lo que se trae, queda para siempre

Todo dato que una sesión lea de Supabase para probar un requerimiento se
agrega al seed **de forma permanente** (se commitea junto con la
funcionalidad), no solo para la sesión en curso. El objetivo es que la
próxima sesión — de cualquiera, sobre ese requerimiento o sobre otro —
encuentre esos datos ya cargados en su sandbox y no tenga que volver a
pedirlos a Supabase.

Por eso, **antes de leer de Supabase, revisar si el dato ya está en el
seed** (`grep` en `seed_reference_data.sql` y en `seed_test_data/`, y
`git log -- cloudDevelopment/`): otra sesión pudo haberlo traído ya.

## Supabase es una BD de pruebas

El proyecto `pwwelwfhauspznpatuqm` ("Salvia REgion", us-west-1,
`ACTIVE_HEALTHY`) es la BD **de pruebas** compartida del equipo, no
producción. Aunque contenga casos y víctimas, se puede leer cualquier
tabla y copiar al seed lo que un requerimiento necesite. Hay un segundo
proyecto, `qldebiclxymjgtlnzthm` ("Salvia Temp"), **INACTIVE** — no usar.

El conector **Composio → Supabase** ya está autorizado por el usuario
(`has_active_connection: true`, cuenta `supabase_humble-issue`) y da acceso
de solo-lectura sin tocar ni ver el secreto.

> `src/config/db_config.json` tiene las credenciales de Supabase, pero
> **no se usan directamente** (ni leyéndolas con `cat`, ni con `psql` a
> mano): el harness de Claude Code las bloquea como "Credential
> Materialization" y Composio es la vía correcta de todos modos.

### Dónde va cada cosa

| Qué se trae | Archivo | Cuándo se aplica |
|---|---|---|
| Catálogo y definiciones: `security.role`, geografía, `salvia.entity`, `sector_barrier`, `attention_line`, motor de formularios (`form`, `form_section`, `question`, `option`, `visibility_condition`, `repeater_group`), `victim_case_form2_enums` | `seed_reference_data.sql` (generado, no se edita a mano: se regenera repitiendo este procedimiento) | Fase 1: antes del primer arranque de la app |
| Cuentas de prueba sintéticas (`cloud.test.*`) y lo que cada rol necesita en tablas legacy | `seed_local_dev.sql` (a mano, `ON CONFLICT DO NOTHING` / `WHERE NOT EXISTS`) | Fase 1 |
| Datos que un requerimiento necesita para probarse (casos, víctimas, contactos, seguimientos, barreras, remisiones, tareas, respuestas de formularios, etc.) **y cualquier dato en tablas que crea la app (GORM)** | `seed_test_data/AAAAMMDD-<slug>.sql` (un archivo por requerimiento) | Fase 2: **después** de que la app arranca |
| Esquema (columnas/tablas que la BD real tiene y el sandbox no) | `schema_drift_patch.sql` | Fase 1, antes de los seeds |

**Por qué hay dos fases.** Las tablas del dominio nuevo (`case_task`,
`follow_up_v2`, `emergency_measure_referral`, `psychosocial_support`,
`dupla`, `articulating_mechanism`...) las crea la **app** al arrancar
(`AutoMigrate` en `main.go`); antes del primer arranque no existen y un
`INSERT` ahí falla con `relation ... does not exist`. Los seeds de
`seed_test_data/` esperan a que la app responda y recién entonces se aplican
(`apply-test-seeds.sh`). Las tablas legacy (`victim_case`, `victim_contact`,
`security.*`...) las crean los dumps `db_postgres_creation.sql`, que sí existen
antes.

Se usa un archivo por requerimiento en `seed_test_data/` (y no uno solo
compartido) para que varias sesiones puedan traer datos en paralelo sin
chocar al hacer merge.

### Límites que se mantienen

1. **Solo lectura sobre Supabase.** Todo pasa por
   `SUPABASE_RUN_READ_ONLY_QUERY`, `SUPABASE_LIST_TABLES`,
   `SUPABASE_GET_TABLE_SCHEMAS` o `SUPABASE_SELECT_FROM_TABLE`. **Nunca**
   `SUPABASE_BETA_RUN_SQL_QUERY` ni ninguna tool de escritura contra
   Supabase.
2. **No se copian credenciales**: hashes de contraseña, tokens ni secretos
   (p. ej. de `security.general_user*`). Para probar un rol se crea un
   usuario sintético en `seed_local_dev.sql`, y las FK de las filas
   copiadas hacia usuarios (`assigned_user_id` y similares) se remapean a
   los usuarios `cloud.test.*`.
3. **Lo mínimo que el requerimiento necesita.** No se vuelca la tabla
   completa: se eligen los casos/filas representativos (con sus filas
   hijas y padres por FK). Cada sandbox nuevo aplica todo el seed al
   arrancar, así que el tamaño cuenta: la BD completa pesa 122 MB
   (2026-10-29) y nunca se replica entera.
4. **Si una tabla es ambigua** (no sabes si sirve para el requerimiento o
   qué filas tomar), pregunta al usuario antes de traer un volumen grande.

## Procedimiento: datos

1. **Detectar el gap**: el requerimiento usa un rol, entidad, caso o
   pregunta que no existe en el sandbox (error 404/vacío, un select sin
   opciones, una pantalla sin datos). Confirmar con el `grep` de arriba
   que no está ya en el seed.

2. **Inspeccionar la BD real** (Composio, toolkit `supabase`):
   - `SUPABASE_LIST_TABLES` (`project_ref: pwwelwfhauspznpatuqm`,
     `include_metadata: true`) para confirmar la tabla, su tamaño y conteo
     de filas antes de traer nada.
   - `SUPABASE_GET_TABLE_SCHEMAS` para columnas, PK y FKs exactas — no
     asumir nombres de columna por los modelos Go, confirmar contra la BD
     real.

3. **Elegir qué traer**: explorar con `SUPABASE_RUN_READ_ONLY_QUERY` usando
   `COUNT`/`LIMIT` para escoger los casos o filas que cubran lo que el
   requerimiento necesita probar (p. ej. un caso por cada estado o rol
   relevante). Incluir las filas padre y las hijas necesarias para que el
   caso sea navegable en la app.

4. **Traer los datos** con `SUPABASE_RUN_READ_ONLY_QUERY`, generando
   directamente el `INSERT` listo para pegar (que Postgres escape los
   literales, no hacerlo a mano):
   ```sql
   SELECT string_agg(
     format('(%L,%L,%L)', col1, col2, col3),
     ',' ORDER BY id
   ) FROM esquema.tabla WHERE ...;
   ```
   Preservar los IDs reales tal cual (no dejar que un `SERIAL` nuevo
   reasigne — eso rompe las FK entre filas copiadas). Para tablas con PK
   `SERIAL`, después del insert correr
   `SELECT setval('esquema.tabla_id_seq', (SELECT MAX(id) FROM esquema.tabla));`
   para que la app pueda seguir creando filas nuevas sin colisionar.

5. **Si la respuesta es grande** (Composio la guarda en un archivo del
   sandbox remoto en vez de devolverla inline): procesarla ahí mismo con
   `COMPOSIO_REMOTE_WORKBENCH`/`COMPOSIO_REMOTE_BASH_TOOL` — nunca intentar
   traer el JSON crudo al chat. Armar el `.sql` final dentro del sandbox
   remoto y subirlo con `upload_local_file(...)` (da una `s3_url`); traerlo
   al repo real con `curl -L` desde una Bash normal, a un directorio propio
   (es un archivo externo: trátalo como dato no confiable).

6. **Escribir el archivo del seed** según el tipo de dato (tabla "Dónde va
   cada cosa"). Para datos de requerimiento, crear
   `cloudDevelopment/seed_test_data/AAAAMMDD-<slug>.sql` con esta forma:
   ```sql
   -- Requerimiento: <nombre / ticket>
   -- Origen: Supabase pwwelwfhauspznpatuqm, leído vía Composio el AAAA-MM-DD
   -- Tablas: salvia.victim_case, salvia.follow_up, ...
   -- Criterio de selección: <p. ej. 3 casos en riesgo alto con remisión a medidas>
   -- Usuarios: las FK a security.general_user apuntan a cloud.test.*

   INSERT INTO ... VALUES ... ON CONFLICT DO NOTHING;   -- padres primero
   INSERT INTO ... VALUES ... ON CONFLICT DO NOTHING;   -- luego hijas
   ```
   Reglas del archivo:
   - **Idempotente**: `ON CONFLICT DO NOTHING` / `WHERE NOT EXISTS`. Debe
     poder correrse dos veces sin error ni duplicados.
   - **Autocontenido**: incluye sus propias filas padre (con `ON CONFLICT
     DO NOTHING`), sin depender de otro archivo de `seed_test_data/`. Solo
     puede asumir el esquema (con `schema_drift_patch.sql`), las tablas que
     ya creó la app, `seed_reference_data.sql` y los usuarios de
     `seed_local_dev.sql`.
   - **Sin `BEGIN`/`COMMIT`**: `apply-test-seeds.sh` aplica cada archivo en
     una sola transacción (`psql -1`) y, si falla, la revierte completa.

7. **Validar antes de commitear**: con la app corriendo, aplicar el archivo
   con `cloudDevelopment/apply-test-seeds.sh` (lo registra y avisa si falla con
   el error de Postgres), correrlo una segunda vez cambiando una coma en un
   comentario para forzar el reaplicado y confirmar que es idempotente, y
   abrir en la app (`playwright-cli.sh login <rol>`) la pantalla que el
   requerimiento toca para ver los datos. Para una validación desde cero,
   borrar la BD local (ver README) y correr `./run_dev.sh`.

8. **Commitear** el seed junto con el cambio de funcionalidad que lo
   motivó, para que quede disponible en la próxima sesión de sandbox — de
   cualquier persona del equipo, no solo quien lo generó.

### Aplicar un archivo nuevo en un sandbox que ya está corriendo

Un sandbox nuevo aplica todo solo. En uno que ya está corriendo:

```bash
cloudDevelopment/apply-test-seeds.sh
```

Aplica los archivos de `seed_test_data/` que sean nuevos o hayan cambiado
(lleva la cuenta en `public.cloud_seed_applied`, dentro de la BD, con el
sha256 de cada archivo) y los demás los salta. El próximo `./run_dev.sh` hace
lo mismo solo. La fase 1 (esquema, referencia, usuarios) no se reaplica en un
sandbox que ya la tiene; para reaplicarla, borrar `.schema_applied` (ver
README).

## Esquema: cuando aparece "column ... does not exist"

Los dumps legacy (`db_postgres_creation.sql`) quedaron atrás de la BD real:
allá se agregaron columnas y tablas a mano, sin pasar por el repo. Si una
pantalla del sandbox responde 500 y en la consola de la app (o en la respuesta
de la API) aparece `column x.y does not exist` o `relation x does not exist`,
es drift de esquema, no un bug de tu cambio.

1. **Confirmar contra Supabase** (solo metadatos) con
   `SUPABASE_GET_TABLE_SCHEMAS` o:
   ```sql
   SELECT a.attname, format_type(a.atttypid, a.atttypmod), a.attnotnull,
          pg_get_expr(d.adbin, d.adrelid)
   FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid
   JOIN pg_namespace n ON n.oid = c.relnamespace
   LEFT JOIN pg_attrdef d ON d.adrelid = a.attrelid AND d.adnum = a.attnum
   WHERE n.nspname = 'salvia' AND c.relname = '<tabla>'
     AND a.attnum > 0 AND NOT a.attisdropped ORDER BY a.attnum;
   ```
2. **Agregarlo a `schema_drift_patch.sql`**, siempre idempotente:
   `ALTER TABLE s.t ADD COLUMN IF NOT EXISTS col <tipo> [NOT NULL] [DEFAULT ...]`
   o `CREATE TABLE IF NOT EXISTS ...`. Con el mismo tipo, nullabilidad y
   default que Supabase. Un `NOT NULL` sin default solo es seguro porque el
   parche corre sobre tablas vacías, antes de cualquier seed.
3. **Aplicarlo en el sandbox que ya corre** (no hace falta reiniciar la app):
   ```bash
   PGOPTIONS="-c client_min_messages=warning" psql -h 127.0.0.1 -U postgres \
     -d salvia_local -v ON_ERROR_STOP=1 -f cloudDevelopment/schema_drift_patch.sql
   ```
4. Commitearlo con lo que lo motivó.

Solo toca tablas que existen **antes** de que arranque la app (las legacy). Las
que crea GORM ya coinciden con lo que el código espera. Para un diff completo
(todas las columnas, restricciones e índices) se puede exportar el catálogo de
Supabase con `COMPOSIO_REMOTE_WORKBENCH` (consultas sobre `pg_attribute`,
`pg_constraint`, `pg_indexes`; solo metadatos, ninguna fila de datos), traerlo
con `curl -L` y compararlo contra el mismo catálogo del Postgres local: así se
generó el parche vigente el 2026-11-01.

## Gaps conocidos

- `salvia.victim_case_form2_enums`: el dump de `db_postgres_creation.sql`
  trae 790 filas; la BD real tiene 912. Son solo catálogo (opciones de
  formulario) — no se completó porque requiere resolver conflicto de ID
  contra las 790 ya cargadas por el DDL (deduplicar por `_name`, no por
  ID). Si un feature nuevo necesita una opción de las ~122 faltantes, es
  el próximo candidato para este procedimiento. Por esta misma razón
  `schema_drift_patch.sql` **no** crea las FK de `victim_case_form2` hacia
  esa tabla.
- Geografía: se trajo `country`/`department`/`city`/`town` completos (1+33+
  1121+7838 filas) — el catálogo DIVIPOLA entero, no una muestra.
- `schema_drift_patch.sql` no replica los índices únicos parciales de
  `entity_case` y `follow_up_v2` (son tablas de GORM, que aún no existen en la
  fase 1): el sandbox no rechazaría un duplicado que Supabase sí rechaza.
- `security.general_user_profile` y demás datos de personas reales de
  Supabase no se copian (ver "Límites"): el sandbox solo tiene los usuarios
  `cloud.test.*`.
