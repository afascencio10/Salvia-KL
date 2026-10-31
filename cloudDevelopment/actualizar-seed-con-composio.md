# Actualizar el seed con Composio-Supabase

Procedimiento para cuando, desarrollando o probando un feature en el sandbox
cloud de Claude Code, el seed local (`seed_reference_data.sql` +
`seed_local_dev.sql`) no tiene algún dato de catálogo que ese feature
necesita — un rol nuevo, un tipo de entidad, una pregunta de formulario, un
código de geografía. En vez de inventar el dato o adivinarlo leyendo código,
se lee la BD de pruebas real (Supabase) con el conector Composio y se agrega
al seed, para que la próxima sesión (de cualquiera) ya lo tenga.

## Por qué Composio y no credenciales directas

`src/config/db_config.json` tiene las credenciales reales de Supabase, pero
**no se deben usar directamente** (ni por `cat`/leerlas, ni por `psql` a
mano) — el harness de Claude Code las bloquea como "Credential
Materialization" y, aparte del bloqueo, es la vía correcta de todos modos:
el conector **Composio → Supabase** ya está autorizado por el usuario
(`has_active_connection: true`, cuenta `supabase_humble-issue`) y da acceso
de solo-lectura sin tocar ni ver el secreto.

- **Proyecto real**: `pwwelwfhauspznpatuqm` ("Salvia REgion", us-west-1,
  `ACTIVE_HEALTHY`). Hay un segundo proyecto, `qldebiclxymjgtlnzthm`
  ("Salvia Temp"), **INACTIVE** — no usar.
- Tamaño real de esa BD (2026-10-29): **122 MB** (`salvia` 53.7 MB +
  `security` 10.9 MB) — muy por encima del umbral de 30 MB del skill
  `setup-clouddev.md`, así que nunca se replica completa, solo catálogo.

## Regla innegociable: qué se puede copiar y qué no

Esta BD tiene datos reales de un programa de atención a víctimas de
violencia basada en género (incluye riesgo de feminicidio). Población
vulnerable. La regla no es de tamaño, es de **naturaleza del dato**:

| Seguro de copiar (estructura/catálogo) | Prohibido (dato de un caso o persona real) |
|---|---|
| `security.role`, `security.country/department/city/town` | `salvia.victim_case*`, `victim_contact*` |
| `salvia.entity`, `salvia.sector_barrier`, `salvia.attention_line` | `salvia.follow_up*`, `barrier*`, `feminicide*` |
| `salvia.form`, `form_section`, `question`, `option`, `visibility_condition`, `repeater_group` (el TEXTO de preguntas/opciones del motor de formularios) | `salvia.form_submission`, `answer`, `repeater_entry` (las RESPUESTAS de un caso) |
| `salvia.victim_case_form2_enums` (catálogo de opciones, ej. "pueblo indígena") | `security.general_user*` reales (usuarios reales del sistema) |
| | `salvia.entity_case*`, `psychosocial_support`, `team_contact`, `contact_attempts`, `moment`, `alert`, `case_log`, `case_owner`, `directories` |

Si una tabla no está en ninguna de las dos listas: por defecto se trata como
prohibida y se pregunta al usuario antes de copiar nada de ahí. El criterio
no es "¿cabe?", es "¿es información sobre una persona o caso específico?" —
si sí, no sale de Supabase bajo ninguna circunstancia, sin importar cuántas
filas sean.

**Nunca** usar una tool de escritura de Composio contra Supabase desde este
flujo — todo pasa por `SUPABASE_RUN_READ_ONLY_QUERY` / `SUPABASE_LIST_TABLES`
/ `SUPABASE_GET_TABLE_SCHEMAS` / `SUPABASE_SELECT_FROM_TABLE` (todas de solo
lectura).

## Procedimiento

1. **Detectar el gap**: un feature nuevo usa un rol/entidad/pregunta que no
   existe en el sandbox (error 404/vacío, un select sin opciones, un rol que
   el código referencia pero el seed no tiene).

2. **Inspeccionar la BD real** (Composio, toolkit `supabase`):
   - `SUPABASE_LIST_TABLES` (`project_ref: pwwelwfhauspznpatuqm`,
     `include_metadata: true`) para confirmar la tabla, su tamaño y conteo
     de filas antes de traer nada.
   - `SUPABASE_GET_TABLE_SCHEMAS` para columnas, PK y FKs exactas — no
     asumir nombres de columna por los modelos Go, confirmar contra la BD
     real.
   - Clasificar la tabla contra la lista de arriba. Si es ambigua, parar y
     preguntar al usuario — no copiar "por si acaso".

3. **Traer los datos** con `SUPABASE_RUN_READ_ONLY_QUERY`, generando
   directamente el `INSERT` listo para pegar (que Postgres escape los
   literales, no hacerlo a mano):
   ```sql
   SELECT string_agg(
     format('(%L,%L,%L)', col1, col2, col3),
     ',' ORDER BY id
   ) FROM esquema.tabla;
   ```
   Preservar los IDs reales tal cual (no dejar que un `SERIAL` nuevo
   reasigne — eso rompe las FK entre filas copiadas). Para tablas con PK
   `SERIAL`, después del insert correr
   `SELECT setval('esquema.tabla_id_seq', (SELECT MAX(id) FROM esquema.tabla));`
   para que la app pueda seguir creando filas nuevas sin colisionar.

4. **Si la respuesta es grande** (Composio la guarda en un archivo del
   sandbox remoto en vez de devolverla inline): procesarla ahí mismo con
   `COMPOSIO_REMOTE_WORKBENCH`/`COMPOSIO_REMOTE_BASH_TOOL` — nunca intentar
   traer el JSON crudo al chat. Armar el `.sql` final dentro del sandbox
   remoto y subirlo con `upload_local_file(...)` (da una `s3_url`); traerlo
   al repo real con `curl -L` desde una Bash normal. Así un dump de cientos
   de miles de filas no pasa entero por el contexto de la conversación.

5. **Agregarlo al seed**:
   - Catálogo/referencia real pulled de Supabase → `seed_reference_data.sql`
     (archivo generado, no se edita a mano — se regenera repitiendo este
     procedimiento).
   - Dato sintético nuevo (ej. un usuario de prueba más) → `seed_local_dev.sql`
     (archivo hand-maintained, sigue el patrón `ON CONFLICT DO NOTHING` /
     `WHERE NOT EXISTS` ya usado ahí).
   - Mantener el orden de aplicación: DDL (`db_postgres_creation.sql` x2) →
     `seed_reference_data.sql` → `seed_local_dev.sql` (ver
     `setup-local-postgres.sh`).

6. **Validar antes de commitear**: levantar una Postgres descartable local
   (`initdb`/`pg_ctl` en un directorio temporal, no el sandbox real) y
   correr el flujo completo — mismo método que se usó para armar este setup
   la primera vez. Confirmar que no hay errores nuevos y, si aplica, que la
   app (`go run .`) arranca contra esa BD.

7. **Commitear** el seed actualizado junto con el cambio de funcionalidad
   que lo motivó, para que quede disponible en la próxima sesión de
   sandbox — de cualquier persona del equipo, no solo quien lo generó.

## Gaps conocidos (no se resolvieron en la primera pasada, a propósito)

- `salvia.victim_case_form2_enums`: el dump de `db_postgres_creation.sql`
  trae 790 filas; la BD real tiene 912. Son solo catálogo (opciones de
  formulario), seguro de completar — no se hizo porque requiere resolver
  conflicto de ID contra las 790 ya cargadas por el DDL (deduplicar por
  `_name`, no por ID). Si un feature nuevo necesita una opción de las ~122
  faltantes, es el próximo candidato para este procedimiento.
- Geografía: se trajo `country`/`department`/`city`/`town` completos (1+33+
  1121+7838 filas) — el catálogo DIVIPOLA entero, no una muestra.
