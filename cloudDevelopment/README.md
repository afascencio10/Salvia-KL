# cloudDevelopment

Permite correr SOG_SALVIA (backend Go) dentro de una sesión **cloud** de Claude
Code, donde no hay salida de red hacia la BD compartida de pruebas (Supabase,
`src/config/db_config.json`). La solución: una BD Postgres local dentro del
mismo sandbox, con esquema + datos de prueba.

En local no cambia nada — `./run_dev.sh` sigue usando `db_config.json` igual
que siempre. La detección de ambiente es automática vía la variable
`CLAUDE_CODE_REMOTE` que exporta el propio sandbox de Claude Code.

## Inicio rápido (sandbox cloud)

```bash
./run_dev.sh
```

Levanta Postgres local, aplica esquema y seeds y arranca la app en
`http://127.0.0.1:9090` (HTTP, sin TLS). La primera vez tarda 1–3 minutos
(compila Go). La app corre en **primer plano**: ejecútalo en segundo plano y
espera a que `GET /seguridad/login` devuelva 200. Los datasets de la fase 2
(`seed_test_data/`) se aplican solos cuando la app responde; verás líneas
`[seed_test_data]` en la misma consola.

| Qué quieres | Comando |
|---|---|
| Comprobar que el ambiente está sano (login con los 19 roles, ~35 s) | `node cloudDevelopment/smoke-login.js` |
| Recorrer **Reporte V2** de inicio a fin (asistente público + ciclo del reporte, ~15 s) | `node cloudDevelopment/e2e/reporte-v2.js` |
| Recorrer **Registro V3** de inicio a fin (Borrador → Terminado, ~15 s) | `node cloudDevelopment/e2e/registro-v3.js` |
| Entrar a la app con Playwright CLI como un rol | `cloudDevelopment/playwright-cli.sh login sv` |
| Seguir con el CLI (estructura, clics, capturas) | `cloudDevelopment/playwright-cli.sh snapshot` / `click e12` / `goto <url>` / `screenshot --filename=/tmp/x.png` |
| Detener la app (para reiniciar tras un cambio de código) | `cloudDevelopment/stop-app.sh` y luego `./run_dev.sh` |
| Aplicar a mano un archivo nuevo de `seed_test_data/` | `cloudDevelopment/apply-test-seeds.sh` |

> **Reiniciar:** el backend es un binario embebido (`go:embed`), así que
> cualquier cambio de Go **o de frontend** exige reiniciar. Usa
> `stop-app.sh`: si se mata solo `go run`, el binario hijo (`bitsflow`) queda
> huérfano escuchando en el puerto y sigue sirviendo código viejo.
> `run_dev.sh` ya lo detiene solo al arrancar.

## Usuarios de prueba

Un usuario por cada uno de los 19 roles: `cloud.test.<rol>`. Clave para
todos: `CloudDev!Test2026`. No son personas reales. (`node
cloudDevelopment/smoke-login.js` los prueba todos.)

| Rol | Nombre | Aterriza en | Datos que lo hacen funcionar |
|---|---|---|---|
| `ad` | Administrador | `/seguridad/usuarios` | |
| `op` | Operador | `/salvia/casos` | |
| `sv` | Supervisor | `/salvia/casos` | |
| `ro` | Operador de seguimiento | `/salvia/casos` | equipo "Riesgo alto" (obligatorio) |
| `do` | Op. territorial departamental | `/salvia/casos` | |
| `no` | Op. territorial nacional | `/salvia/casos` | |
| `et` | Entidad | `/salvia/casos-entidad` | sede de prueba (`entity_branch`) |
| `en` | Enlace territorial | `/salvia/mis-barreras` | departamento `3` (Bogotá, D.C.) |
| `an` | Agente de notificaciones | `/salvia/notificaciones` | |
| `us` | Usuario externo | `/salvia/casos` | |
| `ps` | Psicólogo | `/salvia/mis-remisiones-psicosocial` | dupla de prueba con `ts` |
| `ts` | Trabajador social | `/salvia/mis-remisiones-psicosocial` | dupla de prueba con `ps` |
| `te` | Equipo Técnico | `/salvia/mis-tareas` | |
| `mc` | Mecanismo | `/salvia/casos-mecanismo` | mecanismo municipal de Bogotá |
| `sm` | Supervisor de Medidas | `/salvia/remisiones-medidas` | |
| `em` | Equipo Técnico de Medidas | `/salvia/remisiones-medidas` | |
| `sc` | Supervisión del Convenio | `/salvia/remisiones-medidas` | |
| `cr` | Cruz Roja | `/salvia/remisiones-medidas` | |
| `am` | Agente de Medidas | `/salvia/mis-remisiones-medidas` | |

**El sandbox arranca sin casos** (0 casos, 0 seguimientos, 0 remisiones): hay
usuarios, catálogos, formularios y configuración, pero no datos de casos. Los
datos para probar un requerimiento se traen de Supabase (ver "Política de
datos") y quedan en `seed_test_data/` para las sesiones siguientes.

### Captcha y Playwright

El login tiene un captcha de imagen (`github.com/dchest/captcha`, validado en
el servidor). El backend lo omite cuando el `User-Agent` es `flutter-client`
(es el mecanismo de la app móvil, `GeneralUserFacade.go`). `smoke-login.js` y
`playwright-cli.sh login` mandan ese `User-Agent` **solo en el POST de
`/seguridad/login`**; el resto de la navegación va con el del navegador.

`playwright-cli.sh` instala `@playwright/cli` (una vez, ~2 s, en
`/opt/playwright-cli`) y lo apunta al Chromium que ya trae el sandbox
(`/opt/pw-browsers`). El CLI necesita un Playwright alfa distinto del
preinstalado, por eso no se instala encima del existente. Sus capturas y
snapshots van a `/tmp/playwright-cli`, fuera del repo.

Para llamar a la API con la sesión del navegador (tras `login <rol>`):

```bash
cloudDevelopment/playwright-cli.sh run-code "async page => {
  const r = await page.evaluate(() => fetch('/api/v1/cases/list').then(async r => r.status + ' ' + (await r.text()).slice(0, 200)));
  return r;
}"
```

## Cómo funciona

1. `run_dev.sh` fuentea `setup-local-postgres.sh` antes de `go run`.
2. Ese script es un no-op si `CLAUDE_CODE_REMOTE != "true"`.
3. **Postgres**: busca los binarios (`pg_ctl` no está en el PATH en
   Debian/Ubuntu: viven en `/usr/lib/postgresql/<v>/bin`; si no hay Postgres,
   lo instala por apt), inicializa un cluster en
   `/var/lib/salvia-local-pg/data` (configurable con `SALVIA_LOCAL_PG_DIR`) y
   lo arranca en `127.0.0.1:5432`. El sandbox corre como **root** y Postgres se
   niega a correr como root, así que `initdb`/`pg_ctl` se ejecutan como el
   usuario del SO `postgres`; `psql` sí corre como root (TCP + auth trust).
4. **Fase 1 del seed** (solo la primera vez por sandbox, marcador
   `.schema_applied`): crea la BD `salvia_local`; aplica el esquema legacy
   (`src/security/scripts/db_postgres_creation.sql` +
   `src/salvia/scripts/db_postgres_creation.sql`); aplica
   `schema_drift_patch.sql`; pre-crea las columnas de `security.general_user`
   que la app agregaría después; y carga `seed_reference_data.sql`,
   `seed_reference_enums.sql` y `seed_local_dev.sql`. Exporta
   `DB_HOST/DB_PORT/DB_USER/DB_PASSWORD/DB_NAME`.
5. `common/utils/EntityOperator.go` → `LoadDBCLientConfig()`: si
   `CLAUDE_CODE_REMOTE=true`, arma la config de conexión a partir de esas
   variables en vez de leer el `db_config.json` embebido (Supabase). Es el
   único punto de la app que cambia — tanto el pool GORM como el pool legacy
   pgx pasan por esta misma función.
6. Si había una instancia anterior de la app escuchando, `stop-app.sh` la
   detiene, y el setup lanza en segundo plano `apply-test-seeds.sh --wait`.
7. La app arranca normal: su propio `main.go` hace `AutoMigrate` de los ~34
   modelos GORM + columnas legacy (`migrateLegacyTables`) + seeds embebidos
   (`internal/seed`). Eso ya funcionaba igual en local — no se tocó.
8. **Fase 2 del seed**: apenas la app responde (`/seguridad/login` = 200, o sea
   que ya migró), `apply-test-seeds.sh` aplica `seed_test_data/*.sql`. Es
   aparte porque las tablas del dominio nuevo (GORM: `case_task`,
   `follow_up_v2`, `dupla`, `articulating_mechanism`...) las crea la app, y
   antes del primer arranque no existen. Cada archivo se registra en
   `public.cloud_seed_applied` (nombre + sha256): se aplican solo los nuevos o
   modificados, en cada arranque.

Reinicios de `run_dev.sh` dentro de la misma sesión solo verifican que
Postgres esté arriba y arrancan la app — rápido. Borrar `.schema_applied`
fuerza reaplicar la fase 1; para empezar de cero de verdad, detener Postgres
(`runuser -u postgres -- pg_ctl -D /var/lib/salvia-local-pg/data stop`) y borrar
`/var/lib/salvia-local-pg`.

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

## Qué carga cada fuente

En este orden:

**`schema_drift_patch.sql`** — esquema, no datos. Los dumps legacy
(`db_postgres_creation.sql`) quedaron atrás de la BD real: allá se agregaron a
mano tablas y columnas. Sin este parche varias pantallas responden 500
(`column victim_case_victim_town_code does not exist`, `relation
victim_contact_form1 does not exist`...). Agrega 4 tablas, 82 columnas, 5
anchos de columna, 1 `UNIQUE` y 1 índice; hace nullable 4 columnas que en
Supabase aceptan NULL (sin eso crear un caso nuevo falla: `victim_case_approved_by`
es `NOT NULL` en los dumps) y agrega 4 defaults; y renombra
`security.email.email_id` → `e_mail_id`. Generado el 2026-11-01 comparando el
catálogo de Supabase con el local. Ver en el propio archivo qué **no** se
replica y por qué (FKs, tablas sin uso). Verificado con una escritura real:
`POST /api/v1/victim-cases-v3/draft` como `op` crea el caso (inserta en
`general_user_profile`, `general_user`, `victim_case` y `victim_case_form2`) y
aparece en `/api/v1/cases/list`.

**`seed_reference_data.sql`** — generado, no editar a mano. Catálogo/
referencia **real**, leído de Supabase vía Composio (ver
`actualizar-seed-con-composio.md`):
- Geografía DIVIPOLA **completa**: 1 país, 33 departamentos, 1121 ciudades,
  7838 municipios/corregimientos — no una muestra.
- Catálogo completo de roles: **19** `role_code`.
- Directorio de entidades (41) y sectores de barrera (13) — catálogo de
  instituciones, no datos de un caso.
- Motor de formularios **completo**: `form` (27), `form_section` (89),
  `question`, `option`, `visibility_condition`, `repeater_group` — el texto
  de las preguntas/opciones de todos los formularios reales de la app.
  Son definiciones, no respuestas de ningún caso.

**`seed_reference_enums.sql`** — generado. El catálogo de opciones real
(`salvia.victim_case_form2_enums`, 912 filas y 74 categorías: sí/no, tipos de
documento, tipo de reporte, ajustes razonables, sectores, barreras, criterios de
remisión...). **Reemplaza** el del dump legacy, que tenía otra numeración y
nombres viejos (788 de sus 790 ids significaban otra cosa) y le faltaban 27
categorías: con él, los selects de las pantallas llegaban vacíos (el tipo de
reporte del Reporte V2, p. ej.). Se conservan los ids reales. **La app carga este
catálogo a memoria una sola vez al arrancar**: si lo cambias con la app corriendo,
reiníciala.

**`seed_local_dev.sql`** — hand-maintained, sigue el patrón `ON CONFLICT
DO NOTHING`. Solo lo genuinamente sintético: los **19 usuarios de prueba** y
lo que cada rol necesita antes de que arranque la app (equipo de `ro`,
departamento de `en`, sede de prueba de `et`).

**`seed_test_data/*.sql`** (fase 2, después de que la app arranca) —
datos que dependen de tablas que crea la app y datos leídos de Supabase para
probar requerimientos concretos. Un archivo por requerimiento,
`AAAAMMDD-<slug>.sql`, idempotente y autocontenido. Hoy trae
`20261101-base-mecanismos-y-duplas.sql` (catálogo real de mecanismos, el mecanismo del
usuario `mc` y la dupla `ps` + `ts`) y `20261101-registro-v3-entidades-bogota.sql` (7
sedes reales de Bogotá y las obligaciones de sus entidades, para el paso de
entidades de Registro V3). Es permanente: lo que una sesión agrega aquí lo encuentran todas las
siguientes. Ver `seed_test_data/README.md`.

## Recorridos E2E de las pantallas

Dos recorridos por la UI (Playwright, sin test runner) que usan los datos del seed y
verifican en la BD lo que queda guardado. Se pueden repetir y correr uno tras otro: cada
corrida crea sus datos con un sufijo único (`RUN_ID`). `SCREENSHOTS=/ruta` guarda capturas;
`E2E_VERBOSE=1` muestra la traza completa cuando un escenario falla.

**`e2e/reporte-v2.js` — Registro de Contacto Víctima V2** (el reporte público,
`/salvia/public-v2/nuevo`; 9 escenarios):
- E-02 "No autorizo" corta el flujo · E-06/E-09 "Atrás" conserva lo diligenciado · E-11 "Volver al Inicio".
- E-10 radicar un reporte de cada tipo de reportante (Víctima con hora de contacto, Tercero
  anónimo, Tercero a nombre propio, Familiar y Hombre). Las listas del asistente salen del
  seed y se verifica `victim_contact`, `victim_contact_form2` (tipo de reporte = enum del
  seed) y `victim_contact_form2_v2`.
- Ciclo de vida: radicación pública → el supervisor (`sv`) lo ve "Sin asignar" en Listado
  Reportes y lo reasigna a `op` → `op` lo gestiona en Mis Reportes (esquema 3x3: intento
  fallido, intento exitoso, "Cierre por decisión de la persona") → `sv` lo ve "Cerrado —
  inválido" en solo lectura. Verifica el motivo de cierre (enum del seed) y los 2 intentos.

**`e2e/registro-v3.js` — Registro V3** (`/salvia/registro-v3`, como `op`; 9 escenarios). El
camino feliz E2E-1 de `DocsMD/Screens/Registro V3/QA/`: habeas data → tipo de contacto →
Cabezote (búsqueda sin coincidencias + "Nuevo registro") → CI → DP → CA → HV → AG (la tríada
completa pasa el Caso a Iniciado y desbloquea Tamizaje y Gestión Salvia) → Tamizaje (rama
pareja) → Gestión Salvia (denuncia, remisión a Atención Psicosocial con criterios, entidad
real de Bogotá, barrera, alerta inminente) → Finalizar → Cierre "Terminado". Los campos se
llenan leyendo el catálogo de la propia pantalla (`window.CAMPOS_*_V3`) y los condicionales
en pasadas sucesivas. Verifica en la BD `bo` → `in` → `te`, la entidad y la barrera.

No cubren: E2E-2 (abandonar un Borrador y retomarlo desde el Cabezote), E2E-3 (rama
no-pareja con OSIGD y riesgo extremo) ni la limitación L1 del Reporte V2 (Familiar/Hombre con
disposición "Sí" y hora), ya documentada en `DocsMD`.

## Problemas conocidos del código (no son del sandbox)

Aparecen en la consola al arrancar, también ocurren contra Supabase y **no
impiden usar la app**. No los investigues como si fueran del setup:

- `[WARN] pregunta de riesgo inminente del seguimiento V.4 ... operator does
  not exist: uuid = character varying` (`internal/seed/seed.go`,
  `seguimiento_v4_riesgo_inminente.sql`): compara `form_section.id` (uuid)
  con una variable `varchar` sin cast. La pregunta de riesgo inminente de la V.4
  nunca se agrega, ni aquí ni en Supabase (919 preguntas en ambos lados).
  Arreglo: `WHERE id = sec_id::uuid`.
- `[WARN] reparación de tareas de validación de remisiones ... operator does
  not exist: character varying = uuid` (`validation_tasks.sql`): compara
  `case_timeline_event.psychosocial_support_id` (varchar) con `ps.id` (uuid).
  Arreglo: `e.psychosocial_support_id = ps.id::text`. Ojo: al corregirlo, esa
  reparación **empezaría a escribir** tareas `validar_remision` en cada
  arranque, también en producción.
- `go test` falla en `salvia/service` (el test `followup_v2_service_test.go`
  llama a `NewFollowUpV2Service` con 13 argumentos y pide 14, y eso impide
  compilar todos los tests del paquete), `salvia/dao` (error de `vet` en
  `VictimCaseDAO.go:1569`) e `internal/seed` (un test de la V.4 quedó viejo).
  Corre los tests por paquete.
- `go build ./...` falla en `cmd/backup` (falta `godotenv`); la app (`go build .`)
  compila bien.
- Registro V3 · Cierre: "Se confirma remisión a: victim_case_form2_referral_team_ap" muestra
  el nombre interno del enum en vez de "Atención Psicosocial" (el formulario sí usa la
  etiqueta; el mensaje del cierre, no). Cosmético; no depende del sandbox.
- Registro V3 · DP_30 ("¿En qué zona se encuentra su lugar de residencia?") queda sin
  opciones: la categoría `victim_case_form2_residence_zone` tampoco existe en Supabase. Es un
  campo "deseable", no bloquea nada.

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

Lo que esos dumps **no** traen (porque la BD real cambió a mano después) lo
cierra `schema_drift_patch.sql`, que sí va estricto (`ON_ERROR_STOP=1`).

## Actualizar el seed

- Falta un dato de **catálogo/referencia real** (un rol, una entidad, una
  pregunta de formulario, más geografía): seguir
  [`actualizar-seed-con-composio.md`](actualizar-seed-con-composio.md) —
  se lee de Supabase vía Composio y se agrega a `seed_reference_data.sql`.
- Faltan **datos para probar un requerimiento** (un caso en cierto estado,
  seguimientos, una remisión, respuestas de un formulario): se leen de
  Supabase vía Composio y se guardan en un archivo nuevo
  `seed_test_data/AAAAMMDD-<slug>.sql`, con el mismo procedimiento.
- Una pantalla responde 500 con `column ... does not exist` o `relation ...
  does not exist`: es drift de esquema. Ver "Esquema: cuando aparece 'column ...
  does not exist'" en `actualizar-seed-con-composio.md`.
- Falta un dato **sintético** (otro usuario de prueba, por ejemplo): editar
  `seed_local_dev.sql` directamente (mantener `ON CONFLICT DO NOTHING` /
  `WHERE NOT EXISTS`).

En todos los casos: commitear junto con la funcionalidad que lo motivó. Un
sandbox nuevo aplica todo solo. En uno que ya está corriendo, la fase 1
(esquema, referencia, usuarios) no se reaplica; los archivos nuevos de
`seed_test_data/` sí se aplican en el próximo `./run_dev.sh`, o al momento con
`cloudDevelopment/apply-test-seeds.sh`.

## Archivos de esta carpeta

| Archivo | Para qué |
|---|---|
| `setup-local-postgres.sh` | Lo fuentea `run_dev.sh`. Postgres + fase 1 del seed. |
| `schema_drift_patch.sql` | Esquema: cierra el drift con la BD real. |
| `seed_reference_data.sql` | Catálogo real (generado). |
| `seed_reference_enums.sql` | Catálogo de opciones real (912 filas); reemplaza el del dump. |
| `seed_local_dev.sql` | Los 19 usuarios de prueba. |
| `seed_test_data/` | Fase 2: datos de requerimientos y base de mecanismos/duplas. |
| `apply-test-seeds.sh` | Aplica `seed_test_data/` (lo lanza el setup; también a mano). |
| `stop-app.sh` | Detiene la app que ocupa el puerto (sin matar por nombre). |
| `smoke-login.js` | Login por la UI con los 19 roles (Playwright). |
| `e2e/` | Recorridos de pantallas de inicio a fin (`reporte-v2.js`, `registro-v3.js`) y su `lib.js`. |
| `playwright-cli.sh` | Playwright CLI listo para usar, con `login <rol>`. |
| `actualizar-seed-con-composio.md` | Cómo traer datos y esquema de Supabase. |

## Utilidades que NO pasan por este setup

`src/cmd/seed/main.go`, `src/cmd/migrate/main.go` y `src/cmd/seed/verify/`
son binarios aparte con sus propias credenciales hardcodeadas apuntando a
otra BD remota (Render, `salvia_pruebas_q9xt`) — no pasan por
`LoadDBCLientConfig()`. No funcionan en cloud (sin salida de red) y están
fuera del alcance de este setup, que solo cubre el flujo normal de
`run_dev.sh` / `main.go`.
