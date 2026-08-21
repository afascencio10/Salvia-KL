## 2026-08-20 — Hidratar Seguimiento a Barreras en el tab Formulario

El repeater "Seguimiento a Barreras" (`stateItems = currentBarriers`) no recargaba las respuestas en Detalle porque `formState` no traía el snapshot de barreras. `GetFollowUpDetail` ahora devuelve `activeBarriers` vía `loadActiveBarriers` (IDs fijados en `follow_up_v2.active_barrier_ids`) y el frontend los copia a `formState.currentBarriers` antes de montar `dinamic-form`. Distinto de `barriers` (identificadas en este follow-up).

### Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/internal/models/follow_up_detail.go` | `ActiveBarrierInfo` + campo `activeBarriers` en la respuesta de detalle |
| `src/salvia/service/followup_v2_service.go` | `GetFollowUpDetail` llama `loadActiveBarriers`; DTO movido a models |
| `src/salvia/service/psychosocial_detail_service.go` | Usa `models.ActiveBarrierInfo` (mismo DTO) |
| `src/frontend/html/salvia/follow_up_detail/get_follow_up_detail.html` | `formState.currentBarriers` se asigna antes de `submissionId` |
| `DocsMD/Screens/Detalle de Seguimiento/detalle-seguimiento-interface.md` | Documenta `currentBarriers` y la diferencia vs `barriers` |
| `DocsMD/Screens/Detalle de Seguimiento/QA/resumen-respuestas-formulario-qa-cases.md` | C4/C5: tab Formulario hidrata el repeater; Resumen con barrierFollowUps |

## 2026-08-14 — Respuestas del formulario en el tab Resumen

Se agregó una nueva sub-sección "Respuestas del Formulario" al tab Resumen que muestra, en texto legible, 6 respuestas puntuales del formulario dinámico de seguimiento asociado (leídas directamente por `formSubmissionId`, sin pasar por `dinamic-form`): el análisis de la situación de riesgo, la gestión y actuaciones de cada barrera en seguimiento/identificada (2 repeaters distintos), la gestión realizada en el seguimiento y los elementos que evidencian la remisión. Las opciones de "Gestión de la barrera" (multi-select) se resuelven a labels legibles en el backend.

### Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/internal/models/follow_up_detail.go` | Nuevo `FormAnswers` en `FollowUpDetailResponse` + tipos `FollowUpFormAnswers`/`BarrierAnswerSummary` |
| `src/salvia/service/followup_v2_service.go` | Inyecta `answerRepo`/`repeaterEntryRepo`; nuevo `loadFormAnswers`/`loadBarrierAnswers`/`resolveGestionLabels` llamado desde `GetFollowUpDetail` |
| `src/main.go` | `NewFollowUpV2Service` recibe los 2 repos nuevos |
| `src/frontend/html/salvia/follow_up_detail/get_follow_up_detail.html` | Nueva sección "Respuestas del Formulario" en el tab Resumen + wiring de `formAnswers` |
