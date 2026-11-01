# Registro de Contacto Víctima V2 — Base de datos

> ✅ **Implementado y verificado contra la BD real — 2026-09-27; esquema re-verificado el 2026-10-31** (columnas de `victim_contact`, `victim_contact_form2` y `victim_contact_form2_v2`, y filas de enum).

Mapeo de cada campo del wizard V2 contra el esquema actual (`salvia.victim_contact`, `salvia.victim_contact_form2`, `salvia.victim_case_form2_enums`), verificado leyendo `VictimContactDAO.go`, `VictimContactForm2DAO.go` y `VictimContactController.go` en `SOG_SALVIA-develop/src`.

## Decisiones confirmadas (2026-09-25)

| # | Decisión | Elegida |
|---|---|---|
| 1 | Arquitectura del frontend | **A** — Vue wizard estático a la medida (no `dinamic-form`) |
| 2 | Aislamiento del código backend | **A** — Controller y DAO nuevos, sin compartir código con V1 |
| 3 | Aislar los 2 tipos de reporte nuevos (Familiar/Hombre) | **A** — Nueva clave de enum sobre la misma columna `ReportType` |
| 4 | Captcha en V2 | **A** — Se mantiene, igual que V1 |
| 5 | Campo "anonimato" (Tercero) | **B** — Campo nuevo y dedicado, no reutiliza `report_type_details` |
| 6 | Columnas de "disposición para llamada + hora" | **B** — Se separan: par propio para el reportante, par existente para la víctima |
| 7 | Nombre/apellido raíz sin "persona afectada" (tipoReporte = Hombre) | **A** — Se dejan vacíos, a completar por el agente |

La decisión 2A, llevada a su consecuencia lógica, resuelve también un punto que antes era ambiguo: **ninguna tabla ni archivo existente se modifica**, ni siquiera para agregar columnas nullable. Todo campo genuinamente nuevo (anonimato, parentesco, correo y apellido del reportante, disposición/hora del reportante) vive en una **tabla de extensión nueva**, 1:1 con `victim_contact_form2`. Esto es más estricto que el planteamiento original ("columnas aditivas en las tablas existentes") y evita por completo cualquier cambio, por mínimo que sea, en el código o esquema que usa V1 hoy.

## Regla general de diseño (actualizada)

> **Cero cambios en `victim_contact`, `victim_contact_form2`, `VictimContactDAO.go`, `VictimContactForm2DAO.go` y `VictimContactController.go`.** V2 solo agrega: (a) una tabla nueva de extensión, (b) filas nuevas bajo claves de enum nuevas en la tabla de catálogo ya genérica `victim_case_form2_enums` (sin tocar las claves que usa V1), y (c) archivos Go nuevos (facade, controller, DAO) que leen/escriben esas piezas nuevas más las columnas existentes que sí reutiliza en modo solo-lectura de estructura (no se altera su definición).

## Mapeo campo por campo

| Campo del wizard V2 | Tipo reportante | Destino confirmado |
|---|---|---|
| `form.autorizacion` | Todos | 🆕 `victim_contact_form2_v2.authorization_answer` (tabla nueva) — **corrección sobre el plan original**: `victim_contact_form2` no tiene esta columna; V1 valida el consentimiento pero nunca lo persiste (hallazgo durante la implementación, no un cambio introducido por V2). Se agregó como columna nueva en la extensión, cerrando ese vacío solo para V2 |
| `form.tipoReporte` | Todos | ✅ `victim_contact_form2.report_type` (columna existente) con FK a la clave de enum nueva `victim_case_form2_report_type_v2` (4 valores: `v`, `r`, `fm`, `ho` — nota: el código de Familiar es `fm`, no `fa`; `fa` quedó reservado para "Familia" en la clave de parentesco) — decisión 3A |
| `form.ajusteRazonable` | Todos | ✅ Relación multi-valor `rel_victim_case_form2_enums_victim_contact_form2` (existente; no hay columna en `victim_contact_form2`). UI single-select → se envía como arreglo de 1 y lo escribe `setVictimContactEnumsMultiple` |
| `form.anonimato` | Tercero | 🆕 `victim_contact_form2_v2.anonimato` (tabla nueva) + clave de enum nueva `victim_case_form2_anonimato` — decisión 5B |
| `form.nombreReportante` | Tercero(propio)/Familiar/Hombre | ✅ `victim_contact_form2.reporter_names` (columna **existente**, reutilizada — decisión tomada durante la implementación: la V1 ya tiene este campo de texto libre para "nombre del reportante", así que solo el apellido necesitaba columna nueva) |
| `form.apellidoReportante` | ídem | 🆕 `victim_contact_form2_v2.reporter_last_names` (tabla nueva) |
| `form.correoReportante` | ídem | 🆕 `victim_contact_form2_v2.reporter_email` (tabla nueva) |
| `form.parentesco` | Familiar | 🆕 `victim_contact_form2_v2.parentesco` (tabla nueva) + clave de enum nueva `victim_case_form2_parentesco` (3 valores) |
| `form.telefonoReportante` | Familiar/Hombre | ✅ `victim_contact_form2.reporter_phone` (existente, reutilizada) |
| `form.rolCuidado` | Tercero/Familiar | ✅ `victim_contact_form2.has_care_role` (existente, reutilizada) |
| `form.victimaConoce` | Tercero | ✅ `victim_contact_form2.victim_aware_of_report` (existente, reutilizada) |
| `form.disposicionReportante` + `form.horaReportante` | Familiar/Hombre | 🆕 `victim_contact_form2_v2.reporter_will_receive_call` + `reporter_best_contact_time` (tabla nueva) — decisión 6B |
| `form.nombreVictima` / `form.apellidoVictima` | Víctima/Tercero/Familiar | ✅ `victim_contact.names` / `victim_contact.last_names` (existentes, reutilizadas). Si tipoReporte = 'Hombre': ambas quedan como **cadena vacía `''`** (las columnas aceptan NULL, pero el DTO manda `""`; verificado en un reporte real tipo `ho`) — decisión 7A |
| `form.telefonoVictima` | Víctima/Tercero | ✅ `victim_contact_form2.victim_col_phone` (existente, reutilizada) |
| `form.disposicionLlamada` + `form.horaContacto` | Víctima/Tercero | ✅ `victim_contact_form2.will_receive_call` + `best_contact_time` (existentes, reutilizadas — ahora exclusivas de la víctima, sin doble uso) — decisión 6B |
| `form.descripcionHechos` | Todos | ✅ `victim_contact_form2.facts_description` (existente, reutilizada) |
| Captcha | Todos | ✅ `victim_contact.captcha_id` / `captcha_solution` (existentes, reutilizadas) — decisión 4A |

## Tablas involucradas

| Tabla | Cambio |
|---|---|
| `salvia.victim_contact` | **Ninguno por parte de V2.** V2 escribe `names`, `last_names` y `status = 'v'` (`SetVictimContactDefaults`). Las columnas `captcha_id`/`captcha_solution` solo viajan en el DTO (no existen en la tabla). *Después* de V2, la feature Gestión de reportes agregó a esta tabla `agent_id`, `team`, `close_reason`, `close_observation` y `closed_date` (más la tabla `victim_contact_attempt`); ver [`../Listado Reportes/related-tables.md`](../Listado%20Reportes/related-tables.md) |
| `salvia.victim_contact_form2` | **Ninguno.** Solo lectura de estructura (columnas `report_type`, `reporter_names`, `reporter_phone`, `has_care_role`, `victim_aware_of_report`, `victim_col_phone`, `will_receive_call`, `best_contact_time`, `facts_description`; `adjustments_gbv` va por la tabla de relación `rel_victim_case_form2_enums_victim_contact_form2`). **No** tiene `authorization_answer`: esa respuesta vive en la extensión |
| `salvia.victim_contact_form2_v2` *(NUEVA — creada 2026-09-27)* | Tabla de extensión 1:1 (FK `victim_contact_form2_v2_victim_contact_form2`, NOT NULL; columnas con el prefijo `victim_contact_form2_v2_`): `authorization_answer` *(NOT NULL — ver fila `form.autorizacion` arriba)*, `anonimato`, `parentesco`, `reporter_last_names`, `reporter_email`, `reporter_will_receive_call`, `reporter_best_contact_time` |
| `salvia.victim_case_form2_enums` | Sin tocar filas existentes. + filas nuevas bajo `victim_case_form2_report_type_v2` (4), `victim_case_form2_anonimato` (2), `victim_case_form2_parentesco` (3) |

## Archivos Go nuevos (ninguno existente se modifica)

| Archivo | Rol |
|---|---|
| `salvia/facades/VictimContactFacadeV2.go` | `VictimContactPOST_Public_V2` / `VictimContactPOST_GET_Public_V2` |
| `salvia/controllers/VictimContactControllerV2.go` | `SetVictimContactV2` — validación propia (incluye captcha), completamente independiente de `SetVictimContact` |
| `salvia/dao/VictimContactForm2V2DAO.go` | DTO + funciones SQL para `victim_contact_form2_v2`; internamente llama a las funciones DAO existentes de `VictimContact`/`VictimContactForm2` **sin modificarlas**, solo invocándolas, y adicionalmente inserta la fila de extensión en la misma transacción |
| `salvia/config/Locale.go` | Se **añaden** entradas nuevas (labels de los enums/campos nuevos); no se modifica ninguna entrada existente |

## Migración ejecutada — 2026-09-27

[`migrations/add_victim_contact_v2.sql`](../../../src/migrations/add_victim_contact_v2.sql) (no está en git — `*.sql` está en `.gitignore` del repo, igual que las demás migraciones locales). Corrida contra la BD real con `salvia_gorm` (rol dueño, necesario para el `CREATE TABLE`/`ALTER TABLE`, puerto 6543 — modo sesión de transacción, requerido para que el DDL persista).

**GAP-6 — RESUELTO.** Verificado antes de correr la migración:
```sql
SELECT * FROM salvia.victim_case_form2_enums
WHERE victim_case_form2_enums_name ILIKE '%familiar%' OR victim_case_form2_enums_name ILIKE '%hombre%';
-- 0 filas: no había significados previos que chocaran.
```

**Dos ajustes que el script original no tenía y que hicieron falta al correrlo contra la BD real** (ya corregidos en el archivo):
- `public.uuid_generate_v4()` no existe en este ambiente — la extensión `uuid-ossp` vive en el schema `extensions`, no en `public`. Se cambió a `gen_random_uuid()` (built-in de Postgres, sin dependencia de schema).
- Los `GRANT` originales apuntaban a un rol `salvia_admin` que no existe en esta BD. Se corrigieron a `salvia_legacy` (el rol real, sin el sufijo de enrutamiento del pooler de Supabase) — es el rol con el que la app hace sus queries en runtime (`src/config/db_config.json`, campo `"user"`), así que sin este GRANT la app fallaba con `permission denied for table victim_contact_form2_v2` a pesar de que la migración se veía exitosa.
- **Encontrado después, al escribir el test E2E en `qa-salvia`**: el mismo `permission denied` le pasaba a `helpers/db.ts` (`resetVictimContactV2Fixture`), que limpia los datos de prueba usando el rol `salvia_gorm2` — un cuarto rol más, distinto de `salvia_gorm`/`salvia_legacy`. Se agregó el mismo `GRANT` para `salvia_gorm2`.

**Contrato de error 400** — confirmado durante implementación, misma forma que V1: `{"default":{"captchaID":"...", "global_msg":"..."}, "form2":{"<campo>":"<mensaje>"}, "form2V2":{"<campo>":"<mensaje>"}}`. El frontend lee `errors.default.captchaID` para regenerar el captcha y mapea el resto por `data-entity`/`data-field`.

## Bugs encontrados durante implementación y QA (2026-09-27)

Ninguno estaba en la planeación — surgieron al escribir el código y al probarlo con Playwright + `curl` contra el backend real. Casos de prueba completos en [`QA/E10-cuando-envia-reporte-qa-cases.md`](QA/E10-cuando-envia-reporte-qa-cases.md).

| # | Bug | Causa raíz | Fix |
|---|---|---|---|
| 1 | Pantalla V2 completamente en blanco | `<body>` del layout compartido (`header.html:30`) ya usa `id="app" style="display:none"`; el `<div>` propio de `set_victim_contact_v2.html` duplicaba ese id. `appV2.mount('#app')` montaba sobre el `<body>` (primera coincidencia), pero el `mounted()` de la pantalla intentaba revelar SU PROPIO div duplicado, que `getElementById` nunca encontraba | Se quitó `id="app" style="display:none"` del `<div>` propio — igual que V1, que nunca lo tuvo |
| 2 | `form2V2` llegaba siempre vacío al backend, aunque el JSON de entrada lo traía completo | `VictimContactDTO` implementa `UnmarshalJSON` propio; al embeberlo anónimamente dentro de un struct wrapper, `encoding/json` deja de "aplanar" sus campos promovidos, y esto interfería con la decodificación del campo hermano `form2V2` | Se reemplazó el embedding por dos decodificaciones independientes de los mismos bytes crudos (`json.RawMessage`) — una para `VictimContactDTO`, otra para el wrapper de `form2V2` |
| 3 | Un captcha incorrecto no bloqueaba el envío | `utils.SetError(..., "security_general_user_captcha_fail", salvia_config.Locale)` — ese código de error vive en `security_config.Locale`, no en `salvia_config.Locale`; al no encontrarlo, `SetError` no hace nada (falla en silencio) | Cambiado a `security_config.Locale`, igual que hace V1 para el mismo código |
| 4 | ~8 validaciones de "campo requerido" nunca bloqueaban (nombres vacíos, teléfono en 0, descripción muy corta, "No autorizo") | Mismo patrón que el bug 3: `"common_validation_field_empty_error"` y `"common_validation_field_number_error"` viven en `common_config.Locale`, no en `salvia_config.Locale` | Cambiadas todas a `common_config.Locale` |
| 5 | Primer envío exitoso fallaba con `permission denied for table victim_contact_form2_v2` | Ver "Migración ejecutada" arriba — faltaba el `GRANT` a `salvia_legacy` | `GRANT SELECT, INSERT, UPDATE, DELETE ... TO salvia_legacy` + `GRANT USAGE, SELECT ON SEQUENCE ...` |
| 6 | El botón "Volver al Inicio" en la confirmación quedaba bloqueado (no clickeable) | El popup de éxito compartido (`overlay.html`, usado por `dao.js` para TODO envío exitoso en la app) no tiene lógica de auto-ocultado — a diferencia del popup de error, que sí se oculta solo a los 2s | Se oculta explícitamente `#successOverlay` al llegar a `step = 'confirmacion'` |
| 7 | Validación de longitud mínima de `factsDescription` nunca se disparaba | Comparaba `len(...)` contra `VictimContactForm2FieldDefinitions[...].MinSize`, que es `0` para ese campo — `len(x) < 0` nunca es verdadero | Cambiado a una constante literal (`< 10`), igual que la regla ya usada en el frontend |
