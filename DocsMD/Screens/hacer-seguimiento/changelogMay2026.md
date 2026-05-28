# Changelog Mayo 2026 — `hacer-seguimiento`

---

## 2026-05-26 — Nuevos equipos de remisión: Hombres y Discapacidad

Se agregaron dos nuevos equipos al flujo de derivaciones del formulario de seguimiento. El equipo de **Atención Hombres** requiere que el profesional marque `criterio_hombres` para generar la remisión. El equipo de **Discapacidad** permite seleccionar uno o ambos servicios (`apoyo_lsc`, `enfoque_discapacidad`), creando un registro `DiscapacidadRemision` por cada servicio. Se crearon los modelos GORM, repositorios y se registraron en AutoMigrate.

### Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/internal/models/men_team_remision.go` | Nuevo modelo GORM para remisiones al equipo de hombres |
| `src/internal/models/discapacidad_remision.go` | Nuevo modelo GORM con campo `service` para remisiones a discapacidad |
| `src/internal/repository/men_team_remision_repository.go` | Nuevo repositorio con `FindByCaseID` y `FindByFollowUpID` |
| `src/internal/repository/discapacidad_remision_repository.go` | Nuevo repositorio con `FindByCaseID` y `FindByFollowUpID` |
| `src/salvia/service/form_service.go` | Casos `atencion_hombres` y `discapacidad` en `processFollowUpSubmission` |
| `src/main.go` | AutoMigrate + repos + inyección en `FormServiceDeps` |

---

## 2026-05-26 — Regla de exclusión: medidas de emergencia + psicosocial

Si el profesional selecciona tanto `medidas_emergencia` como `atencion_psico` en la pregunta de equipos, la remisión psicosocial se omite y solo se crean las medidas de emergencia. En el frontend se agregó una pregunta de tipo `info` con visibilidad AND que avisa al usuario de esta incompatibilidad. En el backend se agregó la verificación antes de procesar `atencion_psico`.

### Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/salvia/service/form_service.go` | Verificación de exclusión al inicio del caso `atencion_psico` |
| BD `salvia.question` | Nueva pregunta `info` de advertencia con 2 visibility_conditions AND |

---

## 2026-05-26 — Banner de estado de remisión psicosocial (E-05)

Se agregó un banner informativo en el formulario que muestra en tiempo real si los criterios de remisión al equipo psicosocial se cumplen. El frontend computa el estado en `onAnswersUpdated` y lo pasa como `formState.psysocialRemisionState`. El componente `dinamic-form` lo recibe como prop `:form-state` y lo inyecta en el `render_modification` de tipo REPLACE del banner.

### Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/frontend/html/salvia/follow_up_v2/hacer_seguimiento.html` | `formState`, `onAnswersUpdated`, `:form-state` y `@answers-updated` en dinamic-form |
| `src/frontend/html/salvia/follow_up_detail/get_follow_up_detail.html` | Mismo cambio aplicado para la ruta de detalle/edición |
| BD `salvia.question` | Nueva pregunta `info` de banner con `render_modification` tipo REPLACE |
| BD `salvia.render_modification` | Fila: `search_string='Remisión:'`, `state_path='psysocialRemisionState'` |

---

## 2026-05-26 — Validación de criterios por equipo en backend

`processFollowUpSubmission` ahora valida criterios antes de crear cada remisión: psicosocial requiere `criterio_obligatorio` + ≥ 3 puntos; hombres requiere `criterio_hombres`; estabilización requiere al menos 1 criterio. Si no se cumplen, la remisión se omite y se registra en logs.

### Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/salvia/service/form_service.go` | Validaciones de criterios en los casos `atencion_psico`, `atencion_hombres` y `estabilizacion` |

---

## 2026-05-26 — Logs de trazabilidad en processFollowUpSubmission

Se agregaron 4 puntos de log: equipos seleccionados, criterios cumplidos, criterios NO cumplidos y confirmación después de crear cada registro. Facilitan el diagnóstico durante pruebas y en producción.

### Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/salvia/service/form_service.go` | `log.Printf` en los 4 puntos de trazabilidad de cada equipo |
