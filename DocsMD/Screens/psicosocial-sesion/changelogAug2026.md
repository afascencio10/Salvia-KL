# Changelog Agosto 2026 — `psicosocial-sesion`

---

## 2026-08-12 — Diagrama: barreras §9.6, consentimiento, hechos, agenda condicional

Cambios alineados al diagrama de lógica de negocio (barreras = form Seguimiento, consentimiento PC/PA, hechos en timeline, agenda con validación 2h).

### Consentimiento (PC S4 / PA S4)

- El texto largo del consentimiento pasa a pregunta tipo `info` (banner `df-info-banner`).
- La aceptación queda en `single` Sí/No (“¿Acepta el consentimiento informado?”).
- Si responde **No**, se ocultan las preguntas posteriores de la sección (visibility_condition).
- `session_type` sin cambios: `PRIMER_CONTACTO_SIN_CONSENTIMIENTO` / `CIERRE_NO_CONSENTIMIENTO`.

### Agenda

- Nueva pregunta **¿Agendar nueva sesión?** (Sí/No) antes de fecha/hora, en todos los forms (S1 “Fecha nueva” y S4 “Fecha próxima atención”).
- Si Sí → muestra **Fecha próxima atención** + **Hora próxima atención**.
- Al responder fecha+hora → el frontend consulta `GET /api/v1/psychosocial-support/:id/availability` y muestra un mensaje informativo solo si el horario no está disponible (ventana de 2h; individual o dupla).
- Al guardar: solo se crea `team_contact` agendado si Agendar=Sí **y** el horario está libre. Si no hay cupo, el form se completa igual pero **no** se agenda.

### Barreras

- **Identificación:** sin cambio de lógica (ya crea `barrier_v2` + tareas/oficios según gestión).
- **Seguimiento a Barreras (§9.6):** al completar, crea `barrier_follow_up`, timeline “Seguimiento a Barrera” y cierra a `MANAGED` si aplica. No genera `case_task` (paridad con hacer-seguimiento).

### Hechos de violencia

- Si “Hay nuevos hechos de violencia” = true → evento de timeline `Hechos del caso` (paridad PASO 7b de hacer-seguimiento).

### Archivos principales

| Área | Archivos |
|---|---|
| Docs | `form-psicosocial-data.md`, `flow-E02`, `flujo-psicosocial.md`, `related-tables.md`, índice |
| Seed | `src/cmd/seed/seed_psicosocial.sql`, `src/cmd/seed/migrate_psicosocial_aug2026.sql` |
| Backend | `psicosocial_barreras.go`, `form_service.go`, `psychosocial_detail_*`, repos |
| Frontend | `dinamic-form.js` (tipo `time`), `registrar_sesion.html` |
