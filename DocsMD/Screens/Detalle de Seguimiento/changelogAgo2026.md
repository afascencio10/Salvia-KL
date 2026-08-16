## 2026-08-14 — Respuestas del formulario en el tab Resumen

Se agregó una nueva sub-sección "Respuestas del Formulario" al tab Resumen que muestra, en texto legible, 6 respuestas puntuales del formulario dinámico de seguimiento asociado (leídas directamente por `formSubmissionId`, sin pasar por `dinamic-form`): el análisis de la situación de riesgo, la gestión y actuaciones de cada barrera en seguimiento/identificada (2 repeaters distintos), la gestión realizada en el seguimiento y los elementos que evidencian la remisión. Las opciones de "Gestión de la barrera" (multi-select) se resuelven a labels legibles en el backend.

### Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/internal/models/follow_up_detail.go` | Nuevo `FormAnswers` en `FollowUpDetailResponse` + tipos `FollowUpFormAnswers`/`BarrierAnswerSummary` |
| `src/salvia/service/followup_v2_service.go` | Inyecta `answerRepo`/`repeaterEntryRepo`; nuevo `loadFormAnswers`/`loadBarrierAnswers`/`resolveGestionLabels` llamado desde `GetFollowUpDetail` |
| `src/main.go` | `NewFollowUpV2Service` recibe los 2 repos nuevos |
| `src/frontend/html/salvia/follow_up_detail/get_follow_up_detail.html` | Nueva sección "Respuestas del Formulario" en el tab Resumen + wiring de `formAnswers` |
