---
okf_version: "1.0"
type: Service_Logic
title: "Servicio: Psychosocial3x3Service"
description: "Lógica de dominio del flujo 3x3 de Atención Psicosocial: registro de intentos, umbral diario, tope de 50, elegibilidad de cierre (9 en 3 días), transiciones de estado del proceso, consentimiento, agendamiento de sesión y próximo intento."
owner: "@backend-squad"
status: active
tags: [backend, service, psicosocial, 3x3]
dependencies:
  - Model: 2-data-dictionary/contact-attempts-model.md
  - Model: 2-data-dictionary/team-contact-model.md
code_refs:
  - "src/salvia/service/psychosocial_3x3_service.go"
  - "src/internal/repository/contact_attempt_repository.go"
  - "src/main.go"
  - "src/salvia/config/Menu.go"
  - "src/tests/psychosocial3x3/routes_test.go"
last_updated: "2026-07-09"
---

# Servicio: Psychosocial3x3Service

Concentra las reglas del flujo 3x3 psicosocial. **No** depende de `follow_up_v2` ni de `follow_up_attempts`.

## Responsabilidades

1. **Registrar intento** (exitoso/fallido) en `contact_attempts`, respetando el tope de 50.
2. **Calcular contadores** que gobiernan la UI: intentos fallidos del día, total, días distintos, y los flags derivados (`dailyThresholdReached`, `maxAttemptsReached`, `closureEligible`).
3. **Transicionar el proceso**: a `en_gestion` al primer contacto exitoso; los intentos fallidos lo mantienen en `abierto`.
4. **Registrar consentimiento** (`consent_given`) sobre el intento exitoso; señalar cierre cuando no acepta.
5. **Agendar sesión** en `team_contact` (inmediata o flexible), solo si el consentimiento fue aceptado.
6. **Fijar próximo intento** en `psychosocial_support.next_contact_attempt_at` (acción 3a).
7. **Iniciar el formulario de cierre** (crear `form_submission` del form `fcc7dc8d…`, preseleccionar (editable) el Motivo según el disparador) y **completar el cierre** (mapear Motivo → `cerrado`/`en_devolucion`).
8. **Emitir eventos de timeline** por cada intento y transición.

## Interfaz Pública (Contract)

```go
type Psychosocial3x3Service interface {
    RegisterContactAttempt(ctx context.Context, psicosocialID string, in RegisterAttemptInput) (*AttemptResult, error)
    SetConsent(ctx context.Context, attemptID string, consentGiven bool) (*ConsentResult, error)
    ScheduleSession(ctx context.Context, psicosocialID string, in ScheduleSessionInput) (*SessionResult, error)
    SetNextContactAttempt(ctx context.Context, psicosocialID string, at time.Time) (*Psychosocial, error)
    GetCounters(ctx context.Context, psicosocialID string) (*Counters, error)
    InitClosureForm(ctx context.Context, psicosocialID, reason string) (*ClosureFormInit, error)
    CloseProcess(ctx context.Context, psicosocialID, submissionID string) (status, motivo string, err error)
}
```

### Cierre — mapeo Motivo → estado (RN-12)
```
motivo == 'imposibilidad_contacto_3x3' | 'cumplimiento_objetivos' | 'cumplimiento_esquema'
    -> status = 'cerrado'
motivo == 'no_consentimiento' | 'desistimiento_proceso'
    -> status = 'en_devolucion'
```
`InitClosureForm` crea la `form_submission` del form de cierre (`fcc7dc8d…`) y pre-responde el Motivo (`1c6d97ac…`) con la opción del `reason`, quedando preseleccionado pero editable (dinamic-form no soporta lock por pregunta). `CloseProcess` lee el Motivo de la submission y aplica la transición. El cierre afecta solo al proceso psicosocial (no al caso de víctima).

## Reglas de Negocio (algoritmos)

### Umbral diario (RN-04)
```
dailyFailedCount = COUNT(contact_attempts)
    WHERE psicosocial_id = X AND was_answered = false
      AND attempt_at::date = CURRENT_DATE AND deleted_at IS NULL
dailyThresholdReached = dailyFailedCount >= 3
```
El conteo usa `attempt_at::date` (editable), no `created_at`. Se reinicia cada día natural.

### Tope de intentos (RN-05)
```
totalCount = COUNT(contact_attempts) WHERE psicosocial_id = X AND deleted_at IS NULL
maxAttemptsReached = totalCount >= 50
```
`RegisterContactAttempt` rechaza con `RESOURCE_CONFLICT` si ya hay 50.

### Elegibilidad de cierre por imposibilidad (RN-06)
```
distinctDaysCount = COUNT(DISTINCT attempt_at::date) WHERE psicosocial_id = X AND deleted_at IS NULL
closureEligible = (totalCount >= 9) AND (distinctDaysCount >= 3)
```

### Transición de estado (RN-03)
```
if was_answered == true and process.status == 'abierto':
    process.status = 'en_gestion'
```
El paso a `cerrado` (no acepta / imposibilidad) depende del **formulario de cierre** (🔴 futuro) y no lo ejecuta este servicio todavía.

### Limpieza del próximo intento al registrar (RN-10)
```
al registrar un intento (was_answered true o false):
  if process.next_contact_attempt_at IS NOT NULL
     and process.next_contact_attempt_at::date <= attempt_at::date:
        process.next_contact_attempt_at = NULL
```
El intento programado ya se ejecutó (mismo día o el agente lo hizo días después). Solo se conserva un `next_contact_attempt_at` **a futuro** respecto al intento registrado.

### Orden del historial (GetHistory)
Los intentos se devuelven **descendente por `attempt_at`** (más nuevo → más viejo). El `sequenceNumber` se calcula cronológicamente (1 = más antiguo) antes de invertir el arreglo, así que no cambia con el orden de presentación.

## Errores que Lanza

| Error | Código Catálogo | Cuándo |
| :--- | :--- | :--- |
| `ErrValidation` | `VALIDATION_FAILED` | `was_answered` ausente; fecha/hora inválida; consentimiento sobre intento no exitoso |
| `ErrNotFound` | `RESOURCE_NOT_FOUND` | `psicosocialId` o `attemptId` inexistente |
| `ErrMaxAttempts` | `RESOURCE_CONFLICT` | Tope de 50 intentos alcanzado |
| `ErrConsentRequired` | `RESOURCE_CONFLICT` | Agendar sesión sin consentimiento aceptado |
| `ErrInvalidReason` | `VALIDATION_FAILED` | `reason` inválido al iniciar cierre; submission sin Motivo respondido |
| `ErrInternal` | `SYSTEM_INTERNAL_ERROR` | Fallo inesperado de persistencia |

## Dependencias Externas

- `ContactAttemptRepository` (nuevo) — persistencia de `contact_attempts`.
- Acceso GORM a `psychosocial_support` (estado y `next_contact_attempt_at`), `team_contact` (sesión) y `form_submission`/`answer` (formulario de cierre `fcc7dc8d…`, ver [closure-form-model](/.knowledge/2-data-dictionary/closure-form-model.md)).
- Servicio de timeline del caso — eventos por intento y transición.

## Dependencias 🔴 PENDIENTE

- Destino del redirect de sesión inmediata → **formulario de atención** (spec futuro).
