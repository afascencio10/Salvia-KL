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
| `src/migrations/20261102_wp2_registro_v3_catalogos_y_columnas.sql` | Nuevo — opciones de zona y discapacidad, columnas nuevas y municipios a `varchar(100)`. Se aplica a mano en la BD. |
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
