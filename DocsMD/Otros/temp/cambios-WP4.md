# cambios-WP4 · Registro V3: Plan de acción, entidades y barreras

Sesión 4 de 4 (WP-4 «V3: Plan de acción, entidades y barreras»). Rama `claude/loving-wright-wyvvb3`, rebasada sobre `develop` @ `a7733d5` (ya trae WP-1, WP-2 y WP-3).
Planeación completa y aprobada por Felipe (3 rondas: 2, 3 y 4/11) en Notion: <https://app.notion.com/p/40cc8a4ba7d981729722dff3438337ab>.

## 1. Tareas

| Tarea | Requerimiento | Qué quedó | Estado |
|---|---|---|---|
| T1 | REG-14 | Se quita «gestiones interinstitucionales» (`GS_2`) del formulario y de las obligatorias de Finalizar. La columna y el catálogo se conservan para los datos históricos. | ✅ |
| T9 | REG-38* | Entidades, barreras y los datos de Seguimiento v4 que no tienen columna viven en un **borrador JSON** (`victim_case.victim_case_v3_draft_data`) mientras el Caso es Borrador. Al pasar de `bo` a `in` se materializan (entity_case, entity_case_obligation, barrier_v2) en una transacción y, después del COMMIT, se crean las tareas, los oficios y las remisiones a otros equipos. Después de materializar la sección es de solo lectura (409). | ✅ |
| T2 | REG-18 | La barrera es el formulario de Seguimiento v4 (sector, tipos e institución por sector, 4 listas transversales, fecha, funcionario, departamento, municipio, vereda o zona, descripción, gestión, perspectiva profesional con alerta, descripción de la gestión). Sin competencias. El catálogo `BARRERA_V4_V3` se **genera del seed** (`scripts/gen_registro_v3_barrera_catalog.py`). | ✅ |
| T3 | REG-20 | Canal de activación como Seguimiento v4: selección múltiple (Notificación SALVIA, Articulación institucional…, Denuncia fácil, Mi Salvia, Ninguna) + «¿Qué gestiones se realizaron…?», con las mismas tareas y oficios que Seguimiento. | ✅ |
| T4 | REG-04 (c) | La celda «Gestión» de la tabla de barreras muestra etiqueta, no código (2 de Seguimiento v4 + 6 históricas; vacío «—»; desconocido tal cual). | ✅ |
| T6 | REG-19 | Agregar entidad como Gestión institucional: Departamento\*, Ciudad\*, Municipio\* y Sector\* en cascada, «Buscar entidad» (lista pedida una vez, filtro local sin tildes, debounce 250 ms, descarte de respuestas atrasadas) y acciones institucionales como chips. Sin cambios de backend. | ✅ |
| T5 | REG-16 | Pregunta `GS_20` «elementos psicosociales» (14 opciones del cliente + «Ninguno» excluyente), migración propia, deseable (no obligatoria). | ✅ |
| T7 | REG-15 | «Plan de acción Salvia» reordenado a–i (síntesis, `GS_1`, `GS_15`, `GS_16` + activaciones, `GS_3`, remisión + justificación, `GS_20`, `GS_4`, alerta inminente), sin títulos A–D; formulario de activación de tipo fijo; tarjeta con acciones y su estado. | ✅ |
| T8 | REG-17 | Módulo nuevo «Ruta de atención integral» (entre Tamizaje y Plan de acción, sin candado): réplica de «Identificación de entidades» de Seguimiento v4 sin el componente dinámico, con el bloque de activación para la acudida con acciones pendientes (D17). | ✅ |

No quedó ninguna tarea sin hacer.

## 2. Decisiones aprobadas por Felipe

D3 (14 opciones psicosociales + «Ninguno») · D4 (deseable, «Ninguno» excluyente) · **D11** (sin «Canal institucional»; canales como Seguimiento) · **D13** (la cascada de borrado ocurre solo en el JSON) · **D15** (se materializa al pasar de `bo` a `in`) · **D16 (a)** (al materializar se crean las remisiones a otros equipos) · **D17** (la acudida con acciones pendientes pide también acuerdo, debida diligencia, canales y gestiones) · D19 (sin objeto) · **D20** (después de materializar, entidades, barreras y remisiones son solo lectura).

## 3. Supuestos (tomados con la opción recomendada; revisar)

1. **Canales: 5 opciones, no 4.** El seed de Seguimiento v4 incluye «Ninguna» (no crea tarea); la página de planeación decía 4. Se muestran las 5, y «Ninguna» es excluyente con las demás.
2. **Acudida: acciones.** Se exige marcar al menos una acción (completada o pendiente) cuando la sede tiene acciones en el catálogo. Seguimiento marca ambas listas como obligatorias, pero una entidad puede no tener ninguna completada.
3. **Una acción no puede estar completada y pendiente a la vez** (marcarla en una lista la quita de la otra).
4. **Solo lectura y Finalizar.** Tras materializar se bloquean las remisiones (`GS_6`–`GS_13`), pero **`GS_5` solo se bloquea si ya está respondida** y `GS_14`, `GS_15`, `GS_16` y `GS_17` siguen respondibles: Finalizar las exige y un Caso que pasó a Iniciado sin contestarlas no tendría salida.
5. **Remisión psicosocial:** igual que Seguimiento, se omite si también se pidió «Medidas de emergencia»; se exigen los mismos criterios mínimos (psicosocial: obligatorio + 3 puntos; hombres: obligatorio; estabilización: al menos uno). «Salvia Dignidad» no crea registro (Seguimiento tampoco).
6. **Tipo de entidad** (acudida o activación) en `entityMeta` del JSON (no en `solicitud`); los datos anteriores se reconocen porque guardaban «Canal de activación: …» en `solicitud`.
7. **Campos de Seguimiento sin columna** (vereda o zona, perspectiva, descripción de la gestión, fecha de reporte, acuerdo, debida diligencia, canales, gestiones, ids de las acciones) van a `entityMeta` / `barrierMeta`, que se conserva después de materializar.
8. **«Agregar garantía»** (comentario de Felipe) se leyó como «Agregar entidad» de Gestión institucional: en el código no existe un botón con esa palabra.
9. **REG-29 y REG-31 siguen pendientes del cliente:** se deja el flujo de barreras (botón «+ Agregar barrera a esta entidad» en las tarjetas de **ambas** listas, D8) y los textos actuales (D6).
10. **Reintento al leer:** si un Caso ya salió de borrador y quedó JSON pendiente (el vaciado anterior falló), `GetCase`/`ListEntities`/`ListBarriers` lo materializan, pero ese camino **no crea remisiones** (solo la transición `bo`→`in` lo hace, una vez).
11. Las preguntas obligatorias de Seguimiento lo son **al agregar la entidad** (validación del formulario), no en Finalizar.

## 4. Archivos

**Backend**
- `src/salvia/service/victim_case_v3_draft.go` (nuevo): tipos del JSON, lectura/escritura con `SELECT … FOR UPDATE`, ids `d~<icode>~<uuid>`, altas y bajas del borrador, unión con las tablas, vaciado (`flushDraft`) y decorador de `UpdateStatusIfCurrent`.
- `src/salvia/service/victim_case_v3_tasks.go` (nuevo): tareas y oficios por canal y por gestión de barrera y remisiones a otros equipos (replican `form_service_v4.go` y `processEquiposYMedidas`; no se editó código de Seguimiento).
- `src/salvia/service/victim_case_v3_service.go`: opciones variádicas del constructor, operaciones por modo, `GetCase` con `entityMeta`, `barrierMeta` y `entitiesReadOnly`, `GS_2` fuera de `gsObligatoriasFinalizarV3`.
- `src/salvia/controller/victim_case_v3_controller.go`: 409 `{"error":…,"readOnly":true}` y 404 al eliminar algo que no existe.
- `src/internal/repository/victim_case_v3_repository.go`: campos nuevos en `VictimCaseEntityInput` y `VictimCaseBarrierInput`.
- `src/internal/repository/victim_case_v3_fields.go`: `GS_20` (3 líneas al final de cada mapa).
- `src/main.go`: se pasan las dependencias del borrador al servicio.

**Base de datos (aplicar a mano; el sandbox las aplica solas al arrancar)**
- `src/migrations/20261103_wp4_borrador_json.sql`: columna `salvia.victim_case.victim_case_v3_draft_data jsonb`.
- `src/migrations/20261102_wp4_elementos_psicosociales.sql`: 15 filas en la categoría `victim_case_form2_psychosocial_elements`.

**Frontend**
- `src/frontend/html/salvia/registro_v3/registro_v3.html`: pestaña «Plan de acción Salvia», bloque «Ruta de atención integral» y fragmentos compartidos de la plantilla (`rv3_form_barrera`, `rv3_entidad_buscador`, `rv3_gs_campos`, `rv3_tarjeta_entidad`, `rv3_aviso_borrador`).
- `src/frontend/js/components/registro-v3-app.js` y `registro-v3-campos.js` (catálogos `BARRERA_V4_V3`, `ENTIDAD_V4_V3`, `CANAL_ACTIVACION_V3` generados; `CAMPOS_RA_V3`; `GS_20`).
- `scripts/gen_registro_v3_barrera_catalog.py` (nuevo): genera esos catálogos desde `src/internal/seed/seguimiento_v4_form.sql` (`--check` verifica que estén al día).

**Pruebas (material de apoyo)**: `DocsMD/Otros/temp/tests-wp4/` (ver sección 6).

## 5. Para integración

1. **Seguimiento (`form_service_v4.go`):** las tareas, oficios y remisiones se **replicaron** en `victim_case_v3_tasks.go`. Conviene extraerlas a un helper compartido. Además, los valores de gestión de barrera `activacion` y `alerta` **no tienen etiqueta** en `barrera_detalle.html`, `barrier-follow-up-timeline.js`, `psicosocial_barreras.go`, `followup_v2_service.go` ni `form_service.go`: no se editaron.
2. **Migraciones:** `20261102_wp4_elementos_psicosociales.sql` y `20261103_wp4_borrador_json.sql` están aplicadas solo en el sandbox. Falta aplicarlas en Supabase (develop) y copiarlas **idénticas** a `migrationsProd/` con su fila en el README (regla de `cloudDevelopment/README.md`). No lo hice porque implica escribir en la base compartida.
3. **WP-3 (`GetCase`, `Finalize`, `guardarYSalir`, `beforeunload`):** (a) `GetCase` devuelve tres campos opcionales nuevos (`entityMeta`, `barrierMeta`, `entitiesReadOnly`); (b) `Finalize` informa el módulo `GS` cuando falta `GS_14` o `GS_17`, que ahora se contestan en «Ruta de atención integral»; (c) `guardarYSalir` y `beforeunload` solo guardan CI, DP, CA, HV, AG y GS: **no guardan lo contestado en `RA`** si la persona sale sin pulsar «Guardar sección» (sugerencia: agregar `RA` y enviar `GS_14` y `GS_17`); (d) el constructor `NewVictimCaseV3Service` recibe opciones variádicas (`WithVictimCaseV3Draft`) y envuelve el repositorio para el vaciado, sin editar `SaveSection` ni `Finalize`.
4. **Tests Go de `salvia/service`:** al arrancar la sesión no compilaban (`NewFollowUpV2Service` esperaba 14 argumentos); se compiló con `go build . ./salvia/... ./internal/...` (`go build ./...` falla por `cmd/backup`, que pide un módulo que no está). No se tocaron.
5. **Remisiones y borrador:** las respuestas de remisión deben quedar guardadas **antes** de completar la tríada (D20). El backend rechaza altas y bajas de entidades y barreras (409); la remisión solo se bloquea en pantalla.
6. **Carpeta de MDs del módulo nuevo:** `DocsMD/Screens/Ruta de atención integral/` (se puede cambiar el nombre).
7. **Material de pruebas:** `DocsMD/Otros/temp/tests-wp4/` se puede borrar al integrar; no lo usa ninguna parte de la aplicación.

## 6. Pruebas

Contra la app local (`./run_dev.sh`, Postgres local), con la batería `bash DocsMD/Otros/temp/tests-wp4/run-all.sh` (API primero; Playwright solo para lo que la API no cubre), sobre el código **ya rebasado sobre `develop`**:

| Prueba | Qué cubre | Resultado |
|---|---|---|
| `t9-api-test.sh` | Borrador no toca las tablas; ids `d~`; lectura; duplicado → 409; cascada al borrar la entidad; materialización al completar la tríada (filas, 6 tareas, 2 oficios, JSON vaciado con la meta); idempotencia; solo lectura (409) | 63 ✅ |
| `t9b-api-test.sh` | Remisiones (psicosocial, hombres, discapacidad, estabilización, medidas) y sus criterios; acudida con y sin acuerdo; fallo a mitad (el Caso sigue en `bo` y el JSON queda intacto); sede ya existente; clientes anteriores; «Ninguna»; reintento al leer | 57 ✅ |
| `t5-api-test.sh` | Catálogo de `GS_20`, guardar y recargar, «Ninguno», no obligatoria | 13 ✅ |
| `t1-finalize-test.sh` | Finalizar sin `GS_2`; `GS_4`, `GS_5`, `GS_14`, `GS_16` y `GS_17` siguen exigidas; `GS_20` no | 9 ✅ |
| `ui-t2.js` | Formulario de barrera: condiciones por sector y por gestión, limpieza, departamento → municipio, obligatorias, payload | 25 ✅ |
| `ui-t3.js` | Canales: 5 opciones, «Ninguna» excluyente, gestiones, copias de `GS_15`/`GS_16`, dato histórico | 20 ✅ |
| `ui-t4.js` | Etiquetas de gestión (2 nuevas, 6 históricas, desconocido, vacío) | 14 ✅ |
| `ui-t5.js` | `GS_20`: 15 textos exactos, «Ninguno», guardar y reabrir | 10 ✅ |
| `ui-t6.js` | Cascada obligatoria, una sola petición de sedes, debounce 250 ms, sin tildes, respuesta atrasada, chips | 23 ✅ |
| `ui-t7.js` | Nombre de la pestaña, orden a–i, sin títulos A–D, formulario de activación, tarjeta con acciones | 24 ✅ |
| `ui-t8.js` | Módulo nuevo: pestaña sin candado, preguntas y obligatoriedad, bloque de activación (D17), barreras, guardar `GS_14`/`GS_17` | 34 ✅ |
| `ui-solo-lectura.js` | Solo lectura tras materializar (sin formulario ni botones; remisión deshabilitada); 409 con la pantalla atrasada; recarga al completar la tríada | 21 ✅ |
| `ui-cierre.js` | Cierre con canales y barreras tras el reintento al leer | 5 ✅ |
| **Total** | | **318 ✅ / 0 ❌** |

Durante el rebase sobre `develop` solo hubo que adaptar las pruebas de pantalla al menú vertical de WP-3 (y que ahora solo el Tamizaje tiene candado); el código de WP-4 no cambió. La integración con WP-3 se resolvió en `enviarSeccion` (se conservan sus errores de guardado y se agrega la recarga del borrador al completar la tríada) y en `guardarRutaAtencion` (usa el mismo patrón de mensajes de error).

## 7. MDs actualizados (en el último commit)

`DocsMD/Screens/Registro V3/`: `registro-v3-interface.md`, `registro-v3-index.md`, `registro-v3-secciones.md`, `registro-v3-data.md`, `registro-v3-pendientes.md`, `changelogNov2026.md`, `Flujos/` (E10, E11, E12, E15, E17, E18 y el flujo nuevo `flow-subflujo-materializar-borrador.md`) y `QA/` (casos nuevos y fila en el plan). Carpeta nueva `DocsMD/Screens/Ruta de atención integral/`. `DocsMD/General/proyecto-description.md` y `DocsMD/General/tables-database.md`.
