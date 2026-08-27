# Tablas de Base de Datos — Pantalla Psicosocial Sesión

Describe las tablas involucradas en el registro de sesiones psicosociales. Se reutilizan tablas existentes y se agregan campos nuevos donde corresponde.

---

## Tablas existentes (reutilizadas)

### `salvia.psychosocial_support` — Estado acumulado del proceso psicosocial

Tabla ya existente y en uso por `remisiones-psicosocial-component`. Actúa como el registro maestro de una remisión de atención psicosocial. Para este flujo se agregan dos campos de estado.

**Campos actuales relevantes:**

| Campo | Tipo | Descripción |
|---|---|---|
| `id` | uuid PK | Identificador único de la remisión psicosocial |
| `case_id` | varchar(36) | Caso de seguimiento relacionado |
| `follow_up_id` | uuid | Seguimiento relacionado |
| `session_count` | int DEFAULT 0 | Contador de sesiones completadas (Primera Atención + Seguimientos + Cierre) |
| `status` | varchar(30) DEFAULT 'abierto' | `abierto` / `en_gestion` / `en_devolucion` / `cerrado` |
| `dupla_id` | varchar(36) nullable | Dupla asignada por supervisor |
| `professional_id` | varchar(36) nullable | Profesional asignado directamente |
| `notes` | text nullable | Notas generales |
| `submitted_by` | varchar(36) nullable | Agent que remitió el caso |
| `submitted_by_team` | varchar(50) nullable | Equipo del remitente |
| `scheduled_at` | timestamptz nullable | Fecha de la próxima sesión agendada |
| `created_at`, `updated_at`, `deleted_at` | timestamps | Auditoría |

**Campos nuevos a agregar (migración):**

| Campo | Tipo | Default | Descripción |
|---|---|---|---|
| `ya_hizo_primer_contacto` | boolean | `false` | Se activa al completar el formulario Primer Contacto |
| `ya_hizo_primera_atencion` | boolean | `false` | Se activa al completar el formulario Primera Atención con consentimiento |

**Transiciones de estado:**

| Evento | `status` resultante |
|---|---|
| Remisión creada | `abierto` |
| Primer Contacto completado | `en_gestion` |
| Primera Atención: consentimiento = No | `en_devolucion` |
| Sesión de Cierre completada | `cerrado` |

### Aug 2026 — tablas adicionales tocadas al guardar

| Tabla | Uso |
|---|---|
| `salvia.barrier_follow_up` | Seguimiento a Barreras (§9.6) |
| `salvia.barrier_v2` | Identificación + cierre `MANAGED`; columna `team_contact_id` |
| `salvia.case_task` / `entity_letter` | Solo Identificación (según gestión) |
| `salvia.case_timeline_event` | Sesión, Seguimiento a Barrera, Hechos del caso |
| `salvia.team_contact` | Agenda: `scheduled_date` + `scheduled_time` (ventana 2h) |

Migración datos form: `src/cmd/seed/migrate_psicosocial_aug2026.sql`.

**Cambio en Go (`psychosocial_support.go`):**
```go
// Agregar al struct PsychosocialSupport:
YaHizoPrimerContacto  bool `gorm:"column:ya_hizo_primer_contacto;default:false" json:"yaHizoPrimerContacto"`
YaHizoPrimeraAtencion bool `gorm:"column:ya_hizo_primera_atencion;default:false" json:"yaHizoPrimeraAtencion"`
```
AutoMigrate maneja la migración en Supabase.

---

### `salvia.team_contact` — Registro de cada sesión individual

Tabla ya existente. Cada llamada/sesión del profesional genera un `team_contact`. Se agrega el campo `session_type` para distinguir el tipo de sesión.

**Campos actuales relevantes:**

| Campo | Tipo | Descripción |
|---|---|---|
| `id` | varchar(36) PK | — |
| `case_id` | varchar(36) | Caso relacionado |
| `psicosocial_id` | varchar(36) | FK a `psychosocial_support.id` |
| `form_submission_id` | varchar(36) | FK al formulario completado |
| `professional_id` | varchar(36) | Profesional que atendió |
| `dupla_id` | varchar(36) | Dupla del profesional |
| `scheduled_date` | timestamptz | Fecha agendada de la sesión |
| `scheduled_time` | varchar(8) | Hora agendada |
| `is_completed` | boolean DEFAULT false | `true` cuando la sesión se completó |
| `is_psico_session` | boolean DEFAULT false | `true` para contar como sesión psicosocial |
| `completed_at` | timestamp | Timestamp de completado |
| `status` | varchar(20) | Estado del contacto |
| `summary` | text | Resumen de la sesión |
| `created_at`, `updated_at`, `deleted_at` | timestamps | Auditoría |

**Campos nuevos agregados (migración — implementado Jul 2026):**

| Campo | Tipo | Default | Valores | Descripción |
|---|---|---|---|---|
| `session_type` | varchar(30) | `null` | `PRIMER_CONTACTO` / `PRIMERA_ATENCION` / `ATENCION_PSICOSOCIAL` / `CIERRE` | Tipo de sesión registrada |
| `form_id` | varchar(36) | `null` | UUID de `salvia.form` | Formulario con el que se inició esta sesión (evento E-01). Se fija una sola vez y se reutiliza en cargas posteriores para que la sesión no cambie de formulario si el estado del proceso avanza mientras el contacto sigue pendiente. |

**Cambio en Go (`team_contact.go`):**
```go
// Agregado al struct TeamContact:
FormID      *string `gorm:"type:varchar(36);column:form_id" json:"formId,omitempty"`
SessionType *string `gorm:"type:varchar(30);column:session_type" json:"sessionType,omitempty"`
```

**Regla de asignación profesional/dupla (GAP resuelto):** un `team_contact` nunca
guarda `professional_id` y `dupla_id` a la vez. Al crear uno nuevo, se copia
exactamente el modo de asignación de `psychosocial_support`: si tiene `dupla_id`,
el contacto usa `dupla_id` (cualquier miembro de la dupla puede completarlo); si
tiene `professional_id`, el contacto usa ese `professional_id`.

**Regla de conteo de sesiones** (sin cambios):
```sql
COUNT(*) FROM salvia.team_contact
WHERE psicosocial_id = {id}
  AND is_psico_session = true
  AND is_completed = true
  AND deleted_at IS NULL
```

---

## Tablas de formulario (sistema DinamicForm)

Estas tablas son las que se popula con el seed (`seed_psicosocial.sql`):

### `salvia.form`

| Campo | Tipo | Descripción |
|---|---|---|
| `id` | varchar(36) PK | UUID generado por seed |
| `name` | varchar(255) | Nombre del formulario |
| `description` | text nullable | Descripción |
| `created_at`, `updated_at` | timestamps | — |

**4 registros a insertar** (uno por formulario):
- `Primer Contacto Psicosocial`
- `Primera Atención Psicosocial`
- `Atención Psicosocial` (antes "Seguimiento Psicosocial" — rebranding Jul 2026)
- `Cierre Psicosocial`

---

### `salvia.form_section`

| Campo | Tipo | Descripción |
|---|---|---|
| `id` | uuid PK | — |
| `form_id` | uuid FK → `salvia.form.id` | — |
| `name` | varchar(255) | Nombre visible de la sección |
| `order_index` | int | Orden de la sección dentro del form |
| `is_visible_by_default` | boolean | Si `false`, se oculta hasta que una visibility_condition la habilite |

**Secciones planificadas:**

| Formulario | Sección | `is_visible_by_default` |
|---|---|---|
| Primer Contacto | Primer contacto | true |
| Primera Atención | Contacto Primera Atención | true (false cuando viene de Escenario A continuar=Sí) |
| Primera Atención | Primera Atención | true |
| Seguimiento | Contacto Seguimiento | true |
| Seguimiento | Seguimientos | true |
| Cierre | Contacto Cierre | true |
| Cierre | Seguimiento (Cierre) | true |
| Cierre | Cierre | true |

> La sección "Contacto Primera Atención" se oculta vía `formState.skip_contact = true` cuando el usuario llega desde el Escenario A (Continuar = Sí). El backend envía este formState al cargar el form.

---

### `salvia.question`

| Campo | Tipo | Descripción |
|---|---|---|
| `id` | varchar(36) PK | — |
| `form_section_id` | uuid FK → `salvia.form_section.id` | — |
| `form_id` | varchar(36) FK → `salvia.form.id` | Desnormalizado para queries rápidas |
| `label` | text | Texto de la pregunta |
| `type` | varchar(50) | `single`, `multiple`, `text`, `date`, `boolean` |
| `is_required` | boolean | — |
| `order_index` | int | Orden dentro de la sección |
| `placeholder` | text nullable | — |

---

### `salvia.option`

| Campo | Tipo | Descripción |
|---|---|---|
| `id` | uuid PK | — |
| `question_id` | varchar(36) FK → `salvia.question.id` | — |
| `label` | varchar(255) | Texto de la opción |
| `value` | varchar(255) | Valor guardado en `form_submission_answer` |
| `order_index` | int | — |

---

### `salvia.visibility_condition`

| Campo | Tipo | Descripción |
|---|---|---|
| `id` | uuid PK | — |
| `target_question_id` | varchar(36) nullable | Pregunta que se muestra/oculta |
| `target_section_id` | uuid nullable | Sección que se muestra/oculta |
| `trigger_question_id` | varchar(36) nullable | Pregunta que dispara la condición |
| `trigger_state_path` | varchar(255) nullable | Ruta en `formState` que dispara la condición (ej: `formState.skip_contact`) |
| `operator` | varchar(20) | `eq`, `neq`, `truthy`, `falsy` |
| `expected_value` | varchar(255) nullable | Valor esperado del trigger |

**Condiciones planificadas:**

| Objetivo | Trigger | Operador | Valor |
|---|---|---|---|
| Mostrar "Acciones riesgo inminente" | `¿Se encuentra en riesgo inminente?` | `eq` | `Sí` |
| Mostrar "Descripción hechos" | `Hay nuevos hechos de violencia` | `eq` | `true` |
| Mostrar "Fecha hechos" | `Hay nuevos hechos de violencia` | `eq` | `true` |
| Mostrar "Plan de orientación" | `¿Hay voluntariedad?` | `eq` | `Sí` |
| Mostrar "Compromisos" (PC) | `¿Hay voluntariedad?` | `eq` | `Sí` |
| Mostrar "Fecha próxima atención" (PC) | `¿Hay voluntariedad?` | `eq` | `Sí` |
| Mostrar "Tipo conducta suicida" | `Ingresa por conducta suicida?` | `eq` | `Sí` |
| Ocultar sección "Contacto PA" | `formState.skip_contact` | `truthy` | — |

---

## Diagrama de relaciones (psicosocial)

```
salvia.psychosocial_support (1)
    ├── session_count         [contador de sesiones]
    ├── ya_hizo_primer_contacto  [NUEVO]
    ├── ya_hizo_primera_atencion [NUEVO]
    └── status

    tiene muchos:

    salvia.team_contact (N)
        ├── psicosocial_id → psychosocial_support.id
        ├── form_submission_id → salvia.form_submission.id
        ├── session_type    [NUEVO: PRIMER_CONTACTO / PRIMERA_ATENCION / SEGUIMIENTO / CIERRE]
        └── is_psico_session = true

    cada team_contact referencia:

    salvia.form_submission (1)
        └── form_id → salvia.form.id (uno de los 4 forms psicosociales)
```

---

## Resumen de cambios de migración

| Tabla | Acción | Detalle |
|---|---|---|
| `salvia.psychosocial_support` | ADD COLUMN | `ya_hizo_primer_contacto boolean DEFAULT false` |
| `salvia.psychosocial_support` | ADD COLUMN | `ya_hizo_primera_atencion boolean DEFAULT false` |
| `salvia.team_contact` | ADD COLUMN | `session_type varchar(30) NULL` |
| `salvia.team_contact` | ADD COLUMN | `form_id varchar(36) NULL` |
| `salvia.form` | INSERT | 4 registros (Primer Contacto, Primera Atención, Seguimiento, Cierre) |
| `salvia.form_section` | INSERT | 8 secciones distribuidas en los 4 formularios |
| `salvia.question` | INSERT | ~70 preguntas |
| `salvia.option` | INSERT | ~80 opciones |
| `salvia.visibility_condition` | INSERT | ~15 condiciones |

Los ADD COLUMN son manejados automáticamente por **GORM AutoMigrate** al agregar los campos al struct de Go. Los INSERT se hacen con `seed_psicosocial.sql`.
