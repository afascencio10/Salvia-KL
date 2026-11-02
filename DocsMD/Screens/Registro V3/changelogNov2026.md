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
| `cloudDevelopment/e2e/wp2-*.js` | Nuevos — un script por requerimiento: backend, BD y pantalla (93 comprobaciones en total). |
