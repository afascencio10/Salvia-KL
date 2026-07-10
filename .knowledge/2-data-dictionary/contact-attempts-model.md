---
okf_version: "1.0"
type: Data_Model
title: "Modelo: contact_attempts"
description: "Registro de intentos de contacto (exitosos y fallidos) del flujo 3x3 de Atención Psicosocial. Tabla nueva, independiente de follow_up_attempts (que es exclusiva de Seguimientos)."
owner: "@platform-team"
status: active
tags: [database, data, psicosocial, 3x3, contact-attempts, pii]
engine: PostgreSQL
schema: salvia
table: contact_attempts
code_refs:
  - "src/internal/models/contact_attempt.go"
  - "src/internal/models/psychosocial_support.go"
  - "src/internal/repository/contact_attempt_repository.go"
  - "src/migrations/create_contact_attempts.sql"
last_updated: "2026-07-09"
---

# Modelo de Datos: `salvia.contact_attempts`

Registra **cada** intento de contacto de un proceso de Atención Psicosocial (`psychosocial_support`), tanto los **exitosos** (`was_answered = true`) como los **fallidos** (`was_answered = false`). Reemplaza, para el módulo psicosocial, el uso de `follow_up_attempts` (que queda reservada a Seguimientos).

## Esquema

| Campo | Tipo | Null | Default | Clasificación | Descripción |
| :--- | :--- | :---: | :--- | :--- | :--- |
| `id` | `uuid` | No | `gen_random_uuid()` | — | PK. |
| `psicosocial_id` | `varchar(36)` | No | — | — | FK lógica a `psychosocial_support.id`. Proceso al que pertenece el intento. |
| `case_id` | `varchar(36)` | No | — | — | FK lógica al caso de víctima (denormalizado para consultas/timeline). |
| `was_answered` | `boolean` | No | `false` | — | `true` = la ciudadana contestó (contacto exitoso); `false` = intento fallido. |
| `note` | `text` | Sí | `NULL` | Nivel 2 (PII contextual) | Nota libre de la profesional, capturada en el modal "¿Contestó?". Aplica a intentos **exitosos y fallidos** (ej. "No contestó, celular apagado" / "Sí contestó, pide que lo llamen otro día"). Opcional. Si la profesional quiere indicar a quién intentó contactar (víctima, familiar, institución), lo escribe aquí. |
| `consent_given` | `boolean` | Sí | `NULL` | — | Solo en exitosos: `true`/`false` según Consentimiento Informado. `NULL` mientras no se registre o si no aplica. |
| `professional_id` | `varchar(36)` | Sí | `NULL` | — | Profesional (`ps`/`ts`) que registró el intento. |
| `team` | `varchar(50)` | Sí | `NULL` | — | Equipo de la profesional. |
| `attempt_at` | `timestamptz` | No | `now()` | — | Fecha y hora del intento. Default = ahora, **editable** por la profesional. Base del conteo diario. |
| `created_at` | `timestamptz` | No | `now()` | — | Auditoría de creación de la fila. |
| `updated_at` | `timestamptz` | No | `now()` | — | Auditoría de actualización. |
| `deleted_at` | `timestamptz` | Sí | `NULL` | — | Soft delete (GORM). |

## Reglas e invariantes

- **RN-01** — El conteo diario que dispara el modal de acciones usa `attempt_at::date` (no `created_at`), porque `attempt_at` es editable y representa el momento real del intento.
- **RN-02** — `note` es **opcional** en ambos casos (exitoso y fallido) y se captura en el modal "¿Contestó?". No hay campo estructurado de destinatario: el intento fallido solo requiere `was_answered = false` (y `attempt_at`).
- **RN-03** — `consent_given` se fija vía el endpoint de consentimiento sobre la fila exitosa más reciente.
- **RN-04** — Tope de 50 filas activas por `psicosocial_id` (regla de negocio, no constraint de BD): `COUNT(*) WHERE psicosocial_id = X AND deleted_at IS NULL`.

## Índices sugeridos

- `idx_contact_attempts_psicosocial_id` sobre (`psicosocial_id`).
- `idx_contact_attempts_daily` sobre (`psicosocial_id`, `attempt_at`) para el conteo diario y de días distintos.

## DDL de referencia

```sql
CREATE TABLE salvia.contact_attempts (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    psicosocial_id     VARCHAR(36) NOT NULL,
    case_id            VARCHAR(36) NOT NULL,
    was_answered       BOOLEAN NOT NULL DEFAULT false,
    note               TEXT,
    consent_given      BOOLEAN,
    professional_id    VARCHAR(36),
    team               VARCHAR(50),
    attempt_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at         TIMESTAMPTZ
);

CREATE INDEX idx_contact_attempts_psicosocial_id ON salvia.contact_attempts (psicosocial_id);
CREATE INDEX idx_contact_attempts_daily ON salvia.contact_attempts (psicosocial_id, attempt_at);
```

## Struct GORM de referencia

```go
// ContactAttempt representa un intento de contacto del flujo 3x3 psicosocial.
type ContactAttempt struct {
    ID              string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
    PsicosocialID   string         `gorm:"type:varchar(36);not null;index;column:psicosocial_id" json:"psicosocialId"`
    CaseID          string         `gorm:"type:varchar(36);not null;column:case_id" json:"caseId"`
    WasAnswered     bool           `gorm:"default:false;column:was_answered" json:"wasAnswered"`
    Note            *string        `gorm:"type:text" json:"note,omitempty"`
    ConsentGiven    *bool          `gorm:"column:consent_given" json:"consentGiven,omitempty"`
    ProfessionalID  *string        `gorm:"type:varchar(36);column:professional_id" json:"professionalId,omitempty"`
    Team            *string        `gorm:"type:varchar(50)" json:"team,omitempty"`
    AttemptAt       time.Time      `gorm:"column:attempt_at" json:"attemptAt"`
    CreatedAt       time.Time      `json:"createdAt"`
    UpdatedAt       time.Time      `json:"updatedAt"`
    DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ContactAttempt) TableName() string { return "salvia.contact_attempts" }
```

## Relaciones

- `psicosocial_id` → [`salvia.psychosocial_support`](#) (proceso psicosocial; ver `src/internal/models/psychosocial_support.go`).
- `case_id` → caso de víctima.

## Cambio asociado en `psychosocial_support`

Este flujo añade la columna **`next_contact_attempt_at TIMESTAMPTZ NULL`** a `salvia.psychosocial_support` (fecha/hora del próximo intento de contacto, fijada en la acción 3a). Migración:

```sql
ALTER TABLE salvia.psychosocial_support ADD COLUMN next_contact_attempt_at TIMESTAMPTZ;
```

> [!WARNING]
> **Ejecutar con rol privilegiado.** En el entorno actual `salvia.psychosocial_support` es propiedad del rol `postgres`; el rol de la app (`salvia_gorm`) **no** puede ALTERarla, por lo que `AutoMigrate` **no** agrega esta columna (falla con `must be owner of table psychosocial_support`, SQLSTATE 42501). Esta `ALTER` debe correrse manualmente con `postgres` (p. ej. el SQL editor de Supabase) antes de usar el endpoint `PUT /psychosocial/{id}/next-attempt`. La tabla nueva `contact_attempts` sí se crea vía AutoMigrate (la app es su owner).
