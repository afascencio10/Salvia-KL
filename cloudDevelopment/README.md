# cloudDevelopment

Permite correr SOG_SALVIA (backend Go) dentro de una sesión **cloud** de Claude
Code, donde no hay salida de red hacia la BD compartida de pruebas (Supabase,
`src/config/db_config.json`). La solución: una BD Postgres local dentro del
mismo sandbox, con esquema + datos de prueba sintéticos.

En local no cambia nada — `./run_dev.sh` sigue usando `db_config.json` igual
que siempre. La detección de ambiente es automática vía la variable
`CLAUDE_CODE_REMOTE` que exporta el propio sandbox de Claude Code.

## Cómo funciona

1. `run_dev.sh` fuentea `setup-local-postgres.sh` antes de `go run`.
2. Ese script es un no-op si `CLAUDE_CODE_REMOTE != "true"`.
3. En cloud: instala Postgres (una vez, vía apt — se cachea en la imagen del
   sandbox), inicializa un cluster local en `/var/lib/salvia-local-pg/data`
   (configurable con `SALVIA_LOCAL_PG_DIR`), crea la BD `salvia_local`, aplica
   el esquema legacy (`src/security/scripts/db_postgres_creation.sql` +
   `src/salvia/scripts/db_postgres_creation.sql`) y los seeds de este
   folder (`seed_reference_data.sql`, luego `seed_local_dev.sql`, luego cada
   archivo de `seed_test_data/` en orden alfabético), y exporta
   `DB_HOST/DB_PORT/DB_USER/DB_PASSWORD/DB_NAME`.
4. `common/utils/EntityOperator.go` → `LoadDBCLientConfig()`: si
   `CLAUDE_CODE_REMOTE=true`, arma la config de conexión a partir de esas
   variables en vez de leer el `db_config.json` embebido (Supabase). Es el
   único punto de la app que cambia — tanto el pool GORM como el pool legacy
   pgx pasan por esta misma función.
5. A partir de ahí arranca la app normal: su propio `main.go` hace
   `AutoMigrate` de los ~34 modelos GORM + columnas legacy
   (`migrateLegacyTables`) + seeds embebidos (`internal/seed`). Eso ya
   funcionaba igual en local — no se tocó.
6. Marcador `.schema_applied` en `SALVIA_LOCAL_PG_DIR`: el esquema/seed solo
   se aplican una vez por sandbox. Reinicios de `run_dev.sh` dentro de la
   misma sesión cloud solo verifican que Postgres esté arriba y arrancan la
   app — rápido. Borrar ese archivo fuerza recrear esquema + seed.

## Política de datos: Supabase es una BD de pruebas

La BD de Supabase (proyecto `pwwelwfhauspznpatuqm`, 122 MB al 2026-10-29)
es la BD **de pruebas** compartida del equipo, no producción. Tiene datos
de casos (víctimas, riesgo de feminicidio, seguimientos), y eso **no
impide leerlos ni copiarlos** al seed cuando un requerimiento los necesita
para probarse.

**Regla de trabajo:** cuando una sesión necesite datos para probar un
requerimiento, los lee de Supabase con el MCP de Composio (solo lectura) y
los agrega al seed **de forma permanente**: se commitean junto con la
funcionalidad, no solo para la sesión en curso. Así la próxima sesión —
sobre ese requerimiento o sobre otro — ya los encuentra en su sandbox sin
volver a pedirlos. Antes de leer de Supabase, revisar si el dato ya está
en el seed.

Límites que se mantienen:

- **Solo lectura** sobre Supabase, nunca escritura.
- **Sin credenciales**: no se copian hashes de contraseña, tokens ni
  secretos; las FK hacia usuarios se remapean a los usuarios `cloud.test.*`.
- **Lo mínimo que el requerimiento necesita**, no la BD completa (122 MB):
  cada sandbox nuevo aplica todo el seed al arrancar.

Qué va en qué archivo y el procedimiento paso a paso están en
`actualizar-seed-con-composio.md`.

## Qué carga cada seed

Tres fuentes, aplicadas en este orden (lo hace `setup-local-postgres.sh`):

**`seed_reference_data.sql`** — generado, no editar a mano. Catálogo/
referencia **real**, leído de Supabase vía Composio (ver
`actualizar-seed-con-composio.md`):
- Geografía DIVIPOLA **completa**: 1 país, 33 departamentos, 1121 ciudades,
  7838 municipios/corregimientos — no una muestra.
- Catálogo completo de roles: **19** `role_code` (la lectura de código de la
  primera versión de este setup se había quedado corta en 4: `sm`, `em`,
  `sc`, `cr`).
- Directorio de entidades (41) y sectores de barrera (13) — catálogo de
  instituciones, no datos de un caso.
- Motor de formularios **completo**: `form` (27), `form_section` (89),
  `question`, `option`, `visibility_condition`, `repeater_group` — el texto
  de las preguntas/opciones de todos los formularios reales de la app
  (antes solo existían los 2 que la app crea sola al arrancar). Son
  definiciones, no respuestas de ningún caso.

**`seed_local_dev.sql`** — hand-maintained, sigue el patrón `ON CONFLICT
DO NOTHING` existente. Solo lo genuinamente sintético: **3 usuarios de
prueba** — `cloud.test.ad` (Administrador), `cloud.test.te` (Equipo
Técnico), `cloud.test.ro` (Operador de riesgo, equipo "Riesgo alto"). Clave
para los tres: `CloudDev!Test2026`. No son personas reales.

**`seed_test_data/*.sql`** — datos leídos de Supabase para probar
requerimientos concretos: casos, víctimas, seguimientos, barreras,
remisiones, tareas, respuestas de formularios, etc. Un archivo por
requerimiento, `AAAAMMDD-<slug>.sql`, idempotente y autocontenido; se
aplican en orden alfabético, cada uno en una sola transacción. Un archivo
que falle se omite con un aviso y no detiene el arranque. Es permanente:
lo que una sesión agrega aquí lo encuentran todas las siguientes. Ver
`seed_test_data/README.md` y `actualizar-seed-con-composio.md`.

## Errores esperados al aplicar el esquema

`setup-local-postgres.sh` aplica los dos `db_postgres_creation.sql` con
`ON_ERROR_STOP=0` porque esos dumps (preexistentes, no generados por este
setup) traen algunos defectos ya conocidos que no afectan la app:

- `security`: una FK (`general_user_profile_city`) referencia una columna que
  ya no existe (el código real usa `general_user_profile_town`) — vestigio
  de un rename viejo. Mismo problema que ya tenía la BD real en Supabase.
- `salvia`: dos FK con nombres de columna desactualizados
  (`feminicide_risk_form1_victim_zone` → hoy `..._victim_living_zone`, y una
  referencia similar sobre "subsidy/help"); y los `GRANT`/`ALTER ROLE` a
  `salvia_admin` fallan porque ese rol no existe en una BD nueva (es
  específico del Postgres administrado original).
- `salvia/scripts/db_postgres_creation.sql` se corre **dos veces**: varias
  `ALTER TABLE ... ADD CONSTRAINT` (FK) aparecen antes de que la tabla que
  referencian se cree más abajo en el mismo archivo. La 2da pasada resuelve
  esas FK (y solo genera ruido de "already exists" en lo demás, inocuo).

Se corrigieron en el archivo original (no en una copia) 3 bugs que sí
bloqueaban la creación de tablas completas — no son específicos del
sandbox, habrían fallado igual en cualquier restauración desde cero:
`victim_case_form1` sin `;` de cierre (tumbaba también `case_owner`), una FK
que apuntaba a una tabla `case` inexistente (typo de `victim_case`), y dos
tablas `rel_case_owner_feminicide*` con un `UNIQUE` sobre una columna
`_i_code` que nunca se definió.

## Actualizar el seed

- Falta un dato de **catálogo/referencia real** (un rol, una entidad, una
  pregunta de formulario, más geografía): seguir
  [`actualizar-seed-con-composio.md`](actualizar-seed-con-composio.md) —
  se lee de Supabase vía Composio y se agrega a `seed_reference_data.sql`.
- Faltan **datos para probar un requerimiento** (un caso en cierto estado,
  seguimientos, una remisión, respuestas de un formulario): se leen de
  Supabase vía Composio y se guardan en un archivo nuevo
  `seed_test_data/AAAAMMDD-<slug>.sql`, con el mismo procedimiento.
- Falta un dato **sintético** (otro usuario de prueba, por ejemplo): editar
  `seed_local_dev.sql` directamente (mantener `ON CONFLICT DO NOTHING` /
  `WHERE NOT EXISTS`).

En todos los casos: commitear junto con la funcionalidad que lo motivó. La
próxima vez que `.schema_applied` no exista (sandbox nuevo), se aplica solo.
En un sandbox que ya está corriendo, un archivo nuevo de `seed_test_data/`
no se carga solo: ver "Aplicar un archivo nuevo en un sandbox que ya está
corriendo" en `actualizar-seed-con-composio.md`.

## Utilidades que NO pasan por este setup

`src/cmd/seed/main.go`, `src/cmd/migrate/main.go` y `src/cmd/seed/verify/`
son binarios aparte con sus propias credenciales hardcodeadas apuntando a
otra BD remota (Render, `salvia_pruebas_q9xt`) — no pasan por
`LoadDBCLientConfig()`. No funcionan en cloud (sin salida de red) y están
fuera del alcance de este setup, que solo cubre el flujo normal de
`run_dev.sh` / `main.go`.
