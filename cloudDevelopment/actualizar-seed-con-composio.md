# Traer datos de Supabase al seed (vía Composio)

Procedimiento para cuando, desarrollando o probando un requerimiento en el
sandbox cloud de Claude Code, la BD local no tiene los datos que necesita:
un rol, una entidad, una pregunta de formulario, pero también **casos,
víctimas, seguimientos, barreras, remisiones, tareas o respuestas de
formularios**. En vez de inventar el dato o adivinarlo leyendo código, se
lee la BD de pruebas real (Supabase) con el conector Composio y se agrega
al seed.

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

| Qué se trae | Archivo | Cómo se mantiene |
|---|---|---|
| Catálogo y definiciones: `security.role`, geografía, `salvia.entity`, `sector_barrier`, `attention_line`, motor de formularios (`form`, `form_section`, `question`, `option`, `visibility_condition`, `repeater_group`), `victim_case_form2_enums` | `seed_reference_data.sql` | Generado — no se edita a mano, se regenera repitiendo este procedimiento |
| Datos que un requerimiento necesita para probarse: casos, víctimas, contactos, seguimientos, barreras, remisiones, tareas, respuestas de formularios, etc. | `seed_test_data/AAAAMMDD-<slug>.sql` (un archivo por requerimiento) | Archivo nuevo y autocontenido por requerimiento — ver abajo |
| Cuentas de prueba sintéticas (usuarios `cloud.test.*`) | `seed_local_dev.sql` | A mano, con `ON CONFLICT DO NOTHING` / `WHERE NOT EXISTS` |

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

## Procedimiento

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
   al repo real con `curl -L` desde una Bash normal.

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
     puede asumir el esquema, `seed_reference_data.sql` y los usuarios de
     `seed_local_dev.sql`.
   - **Sin `BEGIN`/`COMMIT`**: el setup aplica cada archivo en una sola
     transacción (`psql -1`).

7. **Validar antes de commitear**: levantar una Postgres descartable local
   (`initdb`/`pg_ctl` en un directorio temporal, no el sandbox real) y
   correr el flujo completo, o al menos aplicar el archivo dos veces
   (`psql -1 -v ON_ERROR_STOP=1 -f ...`) sobre una BD con el esquema y el
   seed base. Confirmar que no hay errores nuevos y, si aplica, que la app
   (`go run .`) arranca contra esa BD y el requerimiento se puede probar.

8. **Commitear** el seed junto con el cambio de funcionalidad que lo
   motivó, para que quede disponible en la próxima sesión de sandbox — de
   cualquier persona del equipo, no solo quien lo generó.

## Aplicar un archivo nuevo en un sandbox que ya está corriendo

`setup-local-postgres.sh` aplica el esquema y todo el seed **una sola vez
por sandbox** (marcador `.schema_applied`). Un archivo de `seed_test_data/`
que llegue después (porque lo agregaste tú en esta sesión, o porque hiciste
`git pull` y trajo uno nuevo) no se carga solo. Para aplicar uno sin
recrear la BD:

```bash
PGPASSWORD=postgres psql -h 127.0.0.1 -p 5432 -U postgres -d salvia_local \
  -1 -v ON_ERROR_STOP=1 -f cloudDevelopment/seed_test_data/<archivo>.sql
```

Un sandbox nuevo los aplica todos, en orden alfabético (la fecha del
nombre da el orden).

## Gaps conocidos

- `salvia.victim_case_form2_enums`: el dump de `db_postgres_creation.sql`
  trae 790 filas; la BD real tiene 912. Son solo catálogo (opciones de
  formulario) — no se completó porque requiere resolver conflicto de ID
  contra las 790 ya cargadas por el DDL (deduplicar por `_name`, no por
  ID). Si un feature nuevo necesita una opción de las ~122 faltantes, es
  el próximo candidato para este procedimiento.
- Geografía: se trajo `country`/`department`/`city`/`town` completos (1+33+
  1121+7838 filas) — el catálogo DIVIPOLA entero, no una muestra.
