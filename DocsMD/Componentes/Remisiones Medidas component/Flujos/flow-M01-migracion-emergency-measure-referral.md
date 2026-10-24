━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🔧 EVENTO: Cabecera de remisión y referral_id
   Tipo: Infrastructure / Backend
   Código: M-01
   Prerequisito de: E-01, E-10, E-13 y todos los flujos que lean la bandeja
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

> **Alcance:** Propone el modelo. No está migrado. La fila de la bandeja es
> una remisión con varias medidas. `salvia.emergency_measure` sigue siendo
> una fila por tipo; su `status` (hoy `ACTIVE` al crear) no alimenta las cards
> ni los filtros de esta pantalla.
>
> Las filas de `emergency_measure` que no tengan `referral_id` no entran al
> listado. Crear la cabecera al remitir un caso queda fuera de estos eventos:
> este componente solo lee y actualiza remisiones que ya existen.


INPUT: {
  gormDB:   conexión GORM al arrancar la aplicación   → src/main.go
  models:   structs Go a migrar                      → src/internal/models/
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  PASO 1 — Crear emergency_measure_referral
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Archivo propuesto: `src/internal/models/emergency_measure_referral.go`

```go
// EmergencyMeasureReferral es la remisión que pinta una fila de la bandeja.
// Las medidas solicitadas viven en salvia.emergency_measure.referral_id.
//
// SQL equivalente:
//
//	CREATE TABLE salvia.emergency_measure_referral (
//	    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
//	    case_id                 VARCHAR(36) NOT NULL,
//	    follow_up_id            UUID,
//	    status                  VARCHAR(60) NOT NULL DEFAULT 'pendiente_de_revision',
//	    urgency                 VARCHAR(30) NOT NULL DEFAULT 'normal',
//	    assigned_agent_id       VARCHAR(36),
//	    submitted_by            VARCHAR(36),
//	    submitted_by_team       VARCHAR(50),
//	    reason                  TEXT,
//	    rejection_reason        TEXT,
//	    referred_at             DATE,
//	    review_professional_id  VARCHAR(36),
//	    review_concept          TEXT,
//	    review_result           VARCHAR(60),
//	    review_date             DATE,
//	    created_at              TIMESTAMPTZ,
//	    updated_at              TIMESTAMPTZ,
//	    deleted_at              TIMESTAMPTZ
//	);
type EmergencyMeasureReferral struct {
	ID                   string         `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CaseID               string         `gorm:"type:varchar(36);not null;column:case_id" json:"caseId"`
	FollowUpID           *string        `gorm:"type:uuid;column:follow_up_id" json:"followUpId,omitempty"`
	Status               string         `gorm:"type:varchar(60);not null;default:'pendiente_de_revision'" json:"status"`
	Urgency              string         `gorm:"type:varchar(30);not null;default:'normal'" json:"urgency"`
	AssignedAgentID      *string        `gorm:"type:varchar(36);column:assigned_agent_id" json:"assignedAgentId,omitempty"`
	SubmittedBy          *string        `gorm:"type:varchar(36);column:submitted_by" json:"submittedBy,omitempty"`
	SubmittedByTeam      *string        `gorm:"type:varchar(50);column:submitted_by_team" json:"submittedByTeam,omitempty"`
	Reason               *string        `gorm:"type:text" json:"reason,omitempty"`
	RejectionReason      *string        `gorm:"type:text;column:rejection_reason" json:"rejectionReason,omitempty"`
	ReferredAt           *time.Time     `gorm:"type:date;column:referred_at" json:"referredAt,omitempty"`
	ReviewProfessionalID *string        `gorm:"type:varchar(36);column:review_professional_id" json:"reviewProfessionalId,omitempty"`
	ReviewConcept        *string        `gorm:"type:text;column:review_concept" json:"reviewConcept,omitempty"`
	ReviewResult         *string        `gorm:"type:varchar(60);column:review_result" json:"reviewResult,omitempty"`
	ReviewDate           *time.Time     `gorm:"type:date;column:review_date" json:"reviewDate,omitempty"`
	CreatedAt            time.Time      `json:"createdAt"`
	UpdatedAt            time.Time      `json:"updatedAt"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"-"`
}

func (EmergencyMeasureReferral) TableName() string { return "salvia.emergency_measure_referral" }
```

`assigned_agent_id`, `submitted_by` y `review_professional_id` guardan `security.general_user.general_user_i_code`. El nombre se resuelve con `general_user_profile_names` + `general_user_profile_last_names`, igual que el listado de Atención Psicosocial.

`case_id` guarda `victim_case.victim_case_i_code`.

La revisión técnica es una sola por remisión (la que pinta la fila). Guardar otra vez (E-13) reemplaza profesional, concepto, resultado y fecha. No hay tabla de historial.

`reason` y `rejection_reason` se guardan y el listado no los pinta. `rejection_reason` corresponde al `motivoRechazo` del seed.


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  PASO 2 — Constantes
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Mismo patrón que los estados de `psychosocial_support`: español snake_case.

```go
const (
	EmergencyReferralStatusRiesgoPorDecidirVia          = "riesgo_inminente_por_decidir_via"
	EmergencyReferralStatusRiesgoPorEscalarCruzRoja     = "riesgo_inminente_por_escalar_cruz_roja"
	EmergencyReferralStatusRiesgoPorEntregarMedida      = "riesgo_inminente_por_entregar_medida"
	EmergencyReferralStatusRiesgoPorConfirmarSeguridad  = "riesgo_inminente_por_confirmar_seguridad"
	EmergencyReferralStatusPendienteDeRevision          = "pendiente_de_revision"
	EmergencyReferralStatusPorAsignar                   = "por_asignar"
	EmergencyReferralStatusAsignada                     = "asignada"
	EmergencyReferralStatusEnGestionDocumental          = "en_gestion_documental"
	EmergencyReferralStatusCaracterizado                = "caracterizado"
	EmergencyReferralStatusPorRevisarTrazabilidad       = "por_revisar_trazabilidad"
	EmergencyReferralStatusPorCotizar                   = "por_cotizar"
	EmergencyReferralStatusPorAprobarCotizaciones       = "por_aprobar_cotizaciones"
	EmergencyReferralStatusAprobadas                    = "aprobadas"
	EmergencyReferralStatusEnSeguimiento                = "en_seguimiento"
	EmergencyReferralStatusDevueltaParaJustificacion    = "devuelta_para_justificacion"
	EmergencyReferralStatusDevueltaASeguimiento         = "devuelta_a_seguimiento"
)

const (
	EmergencyReferralUrgencyRiesgoInminente = "riesgo_inminente"
	EmergencyReferralUrgencyNormal          = "normal"
)

const (
	EmergencyReferralReviewAprobadoCaracterizacion = "aprobado_para_caracterizacion_me"
	EmergencyReferralReviewDevueltoAjuste          = "devuelto_para_ajuste"
	EmergencyReferralReviewDevueltoAcompanamiento  = "devuelto_para_acompanamiento"
)
```

| Valor BD | Label en badge, card y select |
|---|---|
| `riesgo_inminente_por_decidir_via` | Riesgo inminente: por decidir vía |
| `riesgo_inminente_por_escalar_cruz_roja` | Riesgo inminente: por escalar a Cruz Roja |
| `riesgo_inminente_por_entregar_medida` | Riesgo inminente: por entregar medida |
| `riesgo_inminente_por_confirmar_seguridad` | Riesgo inminente: por confirmar seguridad |
| `pendiente_de_revision` | Pendiente de Revisión |
| `por_asignar` | Por asignar |
| `asignada` | Asignada |
| `en_gestion_documental` | En gestión documental |
| `caracterizado` | Caracterizado |
| `por_revisar_trazabilidad` | Por revisar trazabilidad |
| `por_cotizar` | Por cotizar |
| `por_aprobar_cotizaciones` | Por aprobar cotizaciones |
| `aprobadas` | Aprobadas |
| `en_seguimiento` | En Seguimiento |
| `devuelta_para_justificacion` | Devuelta para justificación |
| `devuelta_a_seguimiento` | Devuelta a seguimiento |

| `urgency` | Label |
|---|---|
| `riesgo_inminente` | Riesgo inminente |
| `normal` | Normal |

| `review_result` | Label | Estado que deja E-13 |
|---|---|---|
| `aprobado_para_caracterizacion_me` | Aprobado para caracterización ME | `por_asignar` |
| `devuelto_para_ajuste` | Devuelto para ajuste | `devuelta_para_justificacion` |
| `devuelto_para_acompanamiento` | Devuelto para acompañamiento | `devuelta_a_seguimiento` |

Los cuatro estados `riesgo_inminente_*` son los que suma la card "Riesgo inminente" del supervisor. La card del agente solo suma `riesgo_inminente_por_confirmar_seguridad`.


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  PASO 3 — Agregar referral_id a emergency_measure
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Archivo: `src/internal/models/emergency_measure.go`

Campo nuevo:

```go
ReferralID *string `gorm:"type:uuid;column:referral_id" json:"referralId,omitempty"`
```

`type` sigue siendo la medida solicitada (el chip de la fila). `status` de la medida no cambia de significado: al crear una medida desde el formulario sigue en `ACTIVE`. La bandeja no lo lee.

Un `referral_id` agrupa las medidas de una misma remisión. El listado las devuelve ordenadas por `created_at`.


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  PASO 4 — Registrar AutoMigrate
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

En `src/main.go`, junto al AutoMigrate de `EmergencyMeasure`:

```go
&models.EmergencyMeasureReferral{},
&models.EmergencyMeasure{},
```

`EmergencyMeasure` ya se migra. Volver a incluirlo aplica la columna `referral_id`.


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  PASO 5 — Catálogos de usuario (sin rol nuevo)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

E-04 y E-08 leen agentes. E-11 lee profesionales de revisión. Los dos endpoints devuelven:

```json
{ "id": "general_user_i_code", "nombre": "Nombres Apellidos" }
```

El filtro por rol queda pendiente. Este documento no asigna un código (`sv`, `te`, u otro) a Agente de Medidas, Equipo Técnico de medidas, Supervisión del Convenio ni Cruz Roja. Hasta esa decisión, los endpoints responden el catálogo que corresponda a ese rol y, si el rol no existe, una lista vacía. El seed actual (Pedro Ramírez, Lucía Mendoza, Jorge Castillo, y los tres profesionales de revisión) no se persiste.

  → FIN DE LA MIGRACIÓN PROPUESTA ✓
