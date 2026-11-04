## 2026-11-02 — WP-2: datos, catálogos y Tamizaje de Registro V3 (REG-07/08/09/10/11/12/34)

Siete correcciones de Registro V3. La zona de residencia ofrece Rural y Urbano (REG-07); nueva pregunta «Tipo de discapacidad» con 7 opciones (REG-08); la hora del hecho es un selector (REG-09); los subtipos de violencia muestran su texto (REG-10). Guardar una sección es atómico y valida lo escrito, y si algo falla la persona lo ve (REG-11); los departamentos de hechos y residencia ya se guardan. El Tamizaje se presenta en tablas Sí/No y guarda sus preguntas generales (REG-12), y se puede guardar con las 5 iniciales en «No» (REG-34).

### Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/internal/repository/victim_case_v3_repository.go` | SaveSection y SaveTamizajeBranch en transacción con validación previa; GetAnswers con fechas y horas de los controles; resto del Tamizaje «No». |
| `src/internal/repository/victim_case_v3_validation.go` | Nuevo — `ErrInvalidAnswerV3` y validadores de fecha, hora, número y largo de texto contra la columna. |
| `src/internal/repository/victim_case_v3_fields.go` | Mapeo de HV_4, DP_27 y DP_51; columnas nullable; columnas generales y del resto del Tamizaje. |
| `src/salvia/controller/victim_case_v3_controller.go` | SaveSection responde 400 con `campo` y `detalle` si el dato es inválido; registra rechazos y fallas en el log. |
| `src/frontend/js/components/registro-v3-app.js` | Avisos de error, `etiqueta()` solo en las categorías del campo, limpieza de DP_51 y Tamizaje por etapas. |
| `src/frontend/js/components/registro-v3-campos.js` | HV_3 pasa a tipo `time` y se agrega la pregunta DP_51. |
| `src/frontend/html/salvia/registro_v3/registro_v3.html` | Aviso rojo de guardado, input de hora, Tamizaje en tablas Sí/No y botón visible con las 5 iniciales respondidas. |
| `src/salvia/config/Locale.go` | Textos «Rural», «Urbano» y los 7 tipos de discapacidad. |
| `src/migrations/20261102_wp2_registro_v3_catalogos_y_columnas.sql` | Nuevo — opciones de zona y discapacidad, columnas nuevas y municipios a `varchar(100)`. Aplicada en Supabase el 2026-11-02. |
| `src/internal/repository/victim_case_v3_validation_test.go · victim_case_v3_fields_test.go` | Nuevos — pruebas unitarias de los validadores y de las columnas del Tamizaje. |
| `scripts/e2e/wp2-*.js` | Nuevos — un script por requerimiento: backend, BD y pantalla (93 comprobaciones en total). |

## 2026-11-02 — «Registro de caso»: borrador de Registro V3 desde un reporte (WP-1 T7, REP-03)

Desde el modal de gestión de un reporte (Víctima o Tercero con contacto exitoso) se crea un borrador de Registro V3 prellenado y se abre en el formulario. Corte vertical nuevo: no se modifica `CreateDraft` ni los archivos compartidos de V3 (solo +1 línea y +1 método en `registro-v3-app.js`).

| Archivo | Cambio |
|---|---|
| `src/salvia/controller/victim_case_v3_from_report_controller.go` | Nuevo — `POST /api/v1/victim-cases-v3/desde-reporte/:icode` |
| `src/salvia/service/victim_case_v3_from_report_service.go` | Nuevo — validación previa, credenciales y orquestación |
| `src/internal/repository/victim_case_v3_from_report_repository.go` | Nuevo — la transacción de conversión (reporte bloqueado `FOR UPDATE`, `CreateDraft` en un savepoint, vínculo, prellenado y hito) |
| `src/frontend/js/components/registro-v3-app.js` | +1 llamada en `mounted()` y el método `abrirDesdeReporte()` |
| `src/main.go` | cableado del controller nuevo |

Los borradores creados desde un reporte quedan **sin documento** (tipo y número vacíos) hasta que el agente complete CI_5 y CI_6: el documento temporal `cc/0000000000` de `CreateDraft` hacía que `fanOutPerfilV3` replicara el Perfil entre borradores sin documento (ver el GAP de `Flujos/flow-E19-cuando-carga-con-registro-de-origen.md`).

El reporte NO se cierra: queda «Borrador iniciado» (estado derivado) y se lee «Cerrado — válido» cuando el caso queda en `te` o `cd`; ningún código escribe `ap`.

## 2026-11-02 — WP-3: navegación, guardado y cierre de Registro V3 (REG-01/02/03/04/05/06/13/21/35/36/37)

La pantalla pasa a llamarse «Registro de víctima». Los módulos se eligen en un menú vertical (armado con `modulos[]`) con porcentaje de avance, hay flechas «Anterior/Siguiente» y guardar con éxito pasa al siguiente módulo no bloqueado (REG-01/02). El encabezado muestra nombre y documento reales y el estado con nombre (REG-03/04/06). Solo el Tamizaje se bloquea (REG-13) y guardar GS o TZ ya no repone los candados (REG-35: el servidor devuelve la tríada real). «Finalizar» guarda GS antes de validar y, si falta algo, lista los campos (REG-37); un 401 avisa con claridad sin perder lo escrito (REG-21; era la cookie sin sesión, no un defecto de Finalizar). Los faltantes salen en rojo (REG-05). REG-36 quedó cubierto por REG-11 (WP-2).

### Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/salvia/service/victim_case_v3_service.go` | `SaveSection` evalúa la tríada con cualquier módulo; `ErrModulosIncompletos`/`Finalize` informan `Campos` ordenados. |
| `src/salvia/controller/victim_case_v3_controller.go` | El 400 de `Finalize` incluye `camposFaltantes`. |
| `src/salvia/facades/RegistroV3Facade.go` | `windowTitle` «Registro de víctima». |
| `src/frontend/js/components/registro-v3-app.js` | Menú, navegación y avance; porcentaje; encabezado y estados; `finalizarRegistro` guarda GS y maneja 401; faltantes en rojo; estilos inyectados. |
| `src/frontend/html/salvia/registro_v3/registro_v3.html` | Título y avisos; columna Estado; encabezado con avance; menú vertical y barra de navegación; clase del mensaje de guardado. |
| `scripts/e2e/wp3-*.js` | Pruebas E2E de WP-3 (guardado y cierre, navegación y encabezado, sesión). |

## 2026-11-04 — WP-4: Plan de acción, Ruta de atención integral, entidades y barreras de Registro V3 (REG-14/15/16/17/18/19/20/38\*, REG-04 c)

«Gestión Salvia» pasa a **«Plan de acción Salvia»** con el orden a–i del cliente y sin los títulos A–D (REG-15); nace el módulo **«Ruta de atención integral»** (REG-17); la barrera es el formulario de Seguimiento v4, sin competencias (REG-18); el canal de activación es de selección múltiple y crea las mismas tareas que Seguimiento (REG-20); agregar entidad funciona como en Gestión institucional, con cascada obligatoria y búsqueda con debounce (REG-19); nueva pregunta `GS_20` de elementos psicosociales (REG-16); se retira `GS_2` (REG-14); la «Gestión» de la barrera muestra etiqueta (REG-04 c). Entidades y barreras se guardan en un **borrador JSON** mientras el Caso es Borrador y se materializan al iniciarse, con sus tareas, oficios y remisiones; después son de solo lectura (REG-38\*). Cierra G-V3-22. Plan completo en Notion («Ajustes REgistro»).

### Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/salvia/service/victim_case_v3_draft.go` | Nuevo — borrador JSON, modos por estado, unión con las tablas, vaciado (`flushDraft`) y decorador de `UpdateStatusIfCurrent`. |
| `src/salvia/service/victim_case_v3_tasks.go` | Nuevo — tareas y oficios por canal y por gestión de barrera, y remisiones a otros equipos (replican Seguimiento v4). |
| `src/salvia/service/victim_case_v3_service.go` | Opciones variádicas del constructor, `AddEntity`/`RemoveEntity`/`AddBarrier`/`RemoveBarrier`/`List*` por modo, `GetCase` con `entityMeta`/`barrierMeta`/`entitiesReadOnly`, `GS_2` fuera de las obligatorias de Finalizar. |
| `src/salvia/controller/victim_case_v3_controller.go` | 409 `{error, readOnly: true}` en solo lectura y 404 al eliminar algo inexistente. |
| `src/internal/repository/victim_case_v3_repository.go` · `victim_case_v3_fields.go` | Campos nuevos de `VictimCaseEntityInput`/`VictimCaseBarrierInput`; `GS_20`. |
| `src/main.go` | Dependencias del borrador. |
| `src/migrations/20261103_wp4_borrador_json.sql` | Nuevo — columna `victim_case.victim_case_v3_draft_data jsonb`. |
| `src/migrations/20261102_wp4_elementos_psicosociales.sql` | Nuevo — 15 filas de `victim_case_form2_psychosocial_elements`. |
| `src/frontend/html/salvia/registro_v3/registro_v3.html` | Plan de acción, módulo «Ruta de atención integral» y fragmentos compartidos de la plantilla. |
| `src/frontend/js/components/registro-v3-app.js` · `registro-v3-campos.js` | Módulo RA, formularios de entidad y barrera, canales, solo lectura, `GS_20` y catálogos generados. |
| `scripts/gen_registro_v3_barrera_catalog.py` | Nuevo — genera `BARRERA_V4_V3`, `ENTIDAD_V4_V3` y `CANAL_ACTIVACION_V3` desde el seed de Seguimiento v4. |

**Migraciones:** las dos están aplicadas solo en el sandbox; falta aplicarlas en Supabase (develop) y copiarlas a `migrationsProd/`.
**Pruebas:** 318 comprobaciones (API + BD y Playwright) en verde sobre el código rebasado en `develop`; casos y resultados en [`QA/wp4-plan-de-accion-y-ruta-qa-cases.md`](QA/wp4-plan-de-accion-y-ruta-qa-cases.md).
