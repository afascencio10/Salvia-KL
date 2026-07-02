━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🔧 EVENTO: Migración de schema — submitted_by_team, professional_id y team_contact
   Tipo: Infrastructure / Backend
   Código: M-02
   Prerequisito de: E-01, E-05, E-08, E-10, E-13 (ajustes post-reunión)
   Depende de: M-01 (dupla + psychosocial_support base)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

> **Alcance:** Extiende M-01 con cambios acordados en reunión. M-01 puede estar
> parcialmente implementado con `agent_id`; esta migración renombra a
> `professional_id`, agrega `submitted_by_team` y crea la tabla `team_contact`
> para el conteo de sesiones psicosociales completadas.

INPUT: {
  gormDB:   conexión GORM al arrancar la aplicación   → src/main.go
  models:   structs Go a migrar                      → src/internal/models/
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  PASO 1 — Actualizar modelo PsychosocialSupport
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Archivo: `src/internal/models/psychosocial_support.go`

| Campo Go | Columna BD | Tipo | Cambio |
|---|---|---|---|
| `SubmittedByTeam` | `submitted_by_team` | varchar(50) | **Nuevo** — equipo del remitente denormalizado |
| `ProfessionalID` | `professional_id` | varchar(36) | **Renombrar** desde `agent_id` |
| `AgentID` | — | — | **Eliminar** del modelo (migración de columna) |

SQL equivalente:

```sql
ALTER TABLE salvia.psychosocial_support
  ADD COLUMN IF NOT EXISTS submitted_by_team VARCHAR(50);

ALTER TABLE salvia.psychosocial_support
  RENAME COLUMN agent_id TO professional_id;
```

> `submitted_by` se mantiene para resolver el **nombre** del remitente vía JOIN.
> `submitted_by_team` evita JOIN adicional solo para renderizar el equipo en UI
> y simplifica el filtro E-13.


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  PASO 2 — Crear modelo TeamContact
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Archivo: `src/internal/models/team_contact.go`

```go
// TeamContact registra contactos/sesiones del equipo de Atención Psicosocial.
//
// SQL equivalente:
//
//	CREATE TABLE salvia.team_contact (
//	    id                  VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid(),
//	    case_id             VARCHAR(36) NOT NULL,
//	    form_submission_id  VARCHAR(36),
//	    dupla_id            VARCHAR(36),
//	    psicosocial_id      VARCHAR(36),
//	    professional_id     VARCHAR(36),
//	    team                VARCHAR(50),
//	    scheduled_date      TIMESTAMPTZ,
//	    is_completed        BOOLEAN DEFAULT false,
//	    created_at          TIMESTAMPTZ,
//	    updated_at          TIMESTAMPTZ,
//	    deleted_at          TIMESTAMPTZ,
//	    status              VARCHAR(20),
//	    scheduled_time      VARCHAR(8),
//	    completed_at        TIMESTAMP,
//	    summary             TEXT,
//	    is_psico_session    BOOLEAN DEFAULT false
//	);
type TeamContact struct {
	ID               string         `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()" json:"id"`
	CaseID           string         `gorm:"type:varchar(36);not null;column:case_id"              json:"caseId"`
	FormSubmissionID *string        `gorm:"type:varchar(36);column:form_submission_id"            json:"formSubmissionId"`
	DuplaID          *string        `gorm:"type:varchar(36);column:dupla_id"                      json:"duplaId"`
	PsicosocialID    *string        `gorm:"type:varchar(36);column:psicosocial_id"                json:"psicosocialId"`
	ProfessionalID   *string        `gorm:"type:varchar(36);column:professional_id"               json:"professionalId"`
	Team             *string        `gorm:"type:varchar(50)"                                      json:"team"`
	ScheduledDate    *time.Time     `gorm:"column:scheduled_date"                                   json:"scheduledDate"`
	IsCompleted      bool           `gorm:"column:is_completed;default:false"                       json:"isCompleted"`
	Status           *string        `gorm:"type:varchar(20)"                                        json:"status"`
	ScheduledTime    *string        `gorm:"type:varchar(8);column:scheduled_time"                   json:"scheduledTime"`
	CompletedAt      *time.Time     `gorm:"column:completed_at"                                     json:"completedAt"`
	Summary          *string        `gorm:"type:text"                                               json:"summary"`
	IsPsicoSession   bool           `gorm:"column:is_psico_session;default:false"                  json:"isPsicoSession"`
	CreatedAt        time.Time      `                                                               json:"createdAt"`
	UpdatedAt        time.Time      `                                                               json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `gorm:"index"                                                   json:"deletedAt,omitempty"`
}

func (TeamContact) TableName() string { return "salvia.team_contact" }
```


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  PASO 3 — Registrar AutoMigrate en main.go
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

```go
&models.Dupla{},
&models.PsychosocialSupport{},
&models.TeamContact{},   // nuevo
```


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  PASO 4 — Migración de datos (si M-01 ya aplicó agent_id)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Si la columna aún se llama `agent_id`, ejecutar RENAME (PASO 1).
No hay backfill obligatorio de `submitted_by_team` en esta migración;
población al crear/editar remisión queda fuera de alcance del componente.


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  PASO 5 — Verificación post-migración
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Checklist:

- [ ] Columna `submitted_by_team` existe en `psychosocial_support`
- [ ] Columna `professional_id` existe (ya no `agent_id`)
- [ ] Tabla `salvia.team_contact` existe con `is_psico_session` e `is_completed`
- [ ] Modelo Go `TeamContact` registrado en AutoMigrate

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  IMPACTO EN EVENTOS DEL COMPONENTE
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Evento | Cambio |
|---|---|
| E-01 | Stats por status; sesiones vía `team_contact`; equipo remitente desde columna |
| E-05 | `defaultFilter.professional_id` en lugar de `agent_id` |
| E-08 | Filtro `filter_professional_id` |
| E-10 | COUNT `team_contact` con `is_psico_session AND is_completed` |
| E-13 | Filtro directo `submitted_by_team`; catálogo sin JOIN a `general_user` |
