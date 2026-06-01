# Changelog Mayo 2026 — `hacer-seguimiento`

---

## 2026-05-28 — Reasignación automática de caso según nivel de riesgo

Se implementó la lógica de reasignación de caso en el formulario de seguimiento. Al completar la Sección 1, el sistema evalúa factores protectores, de riesgo y extremos para determinar si el nivel de riesgo del caso debe cambiar. Para casos en riesgo bajo: si hay un factor extremo se reasigna automáticamente a nivel 4; si hay ≥4 factores de riesgo se solicita confirmación al profesional para subir a nivel 3. Para casos en riesgo alto: si hay ≥3 protectores sin extremo se solicita confirmación para bajar a nivel 2. El backend actualiza `victim_case_form2_risk_level`, regenera el calendario para niveles 3-4 (borra PENDIENTE + genera nuevos con `calcularAgente`) o solo reasigna el agente para nivel 2. Se agregaron 3 preguntas en BD (banner info + 2 boolean de confirmación) con sus `visibility_condition` y `render_modification`.

### Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/frontend/html/salvia/follow_up_v2/hacer_seguimiento.html` | `caseRiskLevel` en `data()`, lógica de reasignación en `onAnswersUpdated`, 4 nuevos campos en `formState` |
| `src/internal/repository/follow_up_v2_repository.go` | `RiskLevel` en `VictimCaseInfo`, campo `risk_level` en SQL de `LoadVictimInfoByCaseID` |
| `src/internal/repository/victim_case_light_repository.go` | Nuevos métodos `FindRiskLevelByICode` y `UpdateRiskLevelByICode` |
| `src/internal/repository/followup_repository.go` | Nuevos métodos `DeletePendingByCaseID` y `UpdateAgentForPendingByCaseID` |
| `src/salvia/service/followup_v2_service.go` | Nuevo método `ReasignarCalendario` en interfaz e implementación |
| `src/salvia/service/form_service.go` | Función `reasignarCaso`, PASO 9 en `processFollowUpSubmission`, dep `FollowUpV2Svc` |
| `src/main.go` | `followUpV2Svc` instanciado antes de `formSvc`, inyectado como `FollowUpV2Svc` |
| BD `salvia.question` | 3 preguntas nuevas en Sección 1 (order 7 banner, 8 confirm_high, 9 confirm_low); Q7 y Q8 anteriores movidos a order 10 y 11 |
| BD `salvia.visibility_condition` | 3 condiciones nuevas por `formState` (shouldReassignCase, canReassignHigh, canReassignLow) |
| BD `salvia.render_modification` | 1 fila REPLACE para el banner de reasignación (`search_string='Reasignación:'`, `state_path='reassingText'`) |

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
