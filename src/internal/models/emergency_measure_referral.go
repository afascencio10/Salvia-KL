package models

import (
	"time"

	"gorm.io/gorm"
)

// EmergencyMeasureReferral es la remisión que pinta una fila de la bandeja.
// Las medidas solicitadas viven en salvia.emergency_measure.referral_id.
// El status de la medida no alimenta esta bandeja.
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

const (
	EmergencyReferralStatusRiesgoPorDecidirVia         = "riesgo_inminente_por_decidir_via"
	EmergencyReferralStatusRiesgoPorEscalarCruzRoja    = "riesgo_inminente_por_escalar_cruz_roja"
	EmergencyReferralStatusRiesgoPorEntregarMedida     = "riesgo_inminente_por_entregar_medida"
	EmergencyReferralStatusRiesgoPorConfirmarSeguridad = "riesgo_inminente_por_confirmar_seguridad"
	EmergencyReferralStatusPendienteDeRevision         = "pendiente_de_revision"
	EmergencyReferralStatusPorAsignar                  = "por_asignar"
	EmergencyReferralStatusAsignada                    = "asignada"
	EmergencyReferralStatusEnGestionDocumental         = "en_gestion_documental"
	EmergencyReferralStatusCaracterizado               = "caracterizado"
	EmergencyReferralStatusPorRevisarTrazabilidad      = "por_revisar_trazabilidad"
	EmergencyReferralStatusPorCotizar                  = "por_cotizar"
	EmergencyReferralStatusPorAprobarCotizaciones      = "por_aprobar_cotizaciones"
	EmergencyReferralStatusAprobadas                   = "aprobadas"
	EmergencyReferralStatusEnSeguimiento               = "en_seguimiento"
	EmergencyReferralStatusDevueltaParaJustificacion   = "devuelta_para_justificacion"
	EmergencyReferralStatusDevueltaASeguimiento        = "devuelta_a_seguimiento"
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

// EmergencyReferralStatusLabel es el texto del badge, la card y el select.
func EmergencyReferralStatusLabel(status string) string {
	switch status {
	case EmergencyReferralStatusRiesgoPorDecidirVia:
		return "Riesgo inminente: por decidir vía"
	case EmergencyReferralStatusRiesgoPorEscalarCruzRoja:
		return "Riesgo inminente: por escalar a Cruz Roja"
	case EmergencyReferralStatusRiesgoPorEntregarMedida:
		return "Riesgo inminente: por entregar medida"
	case EmergencyReferralStatusRiesgoPorConfirmarSeguridad:
		return "Riesgo inminente: por confirmar seguridad"
	case EmergencyReferralStatusPendienteDeRevision:
		return "Pendiente de Revisión"
	case EmergencyReferralStatusPorAsignar:
		return "Por asignar"
	case EmergencyReferralStatusAsignada:
		return "Asignada"
	case EmergencyReferralStatusEnGestionDocumental:
		return "En gestión documental"
	case EmergencyReferralStatusCaracterizado:
		return "Caracterizado"
	case EmergencyReferralStatusPorRevisarTrazabilidad:
		return "Por revisar trazabilidad"
	case EmergencyReferralStatusPorCotizar:
		return "Por cotizar"
	case EmergencyReferralStatusPorAprobarCotizaciones:
		return "Por aprobar cotizaciones"
	case EmergencyReferralStatusAprobadas:
		return "Aprobadas"
	case EmergencyReferralStatusEnSeguimiento:
		return "En Seguimiento"
	case EmergencyReferralStatusDevueltaParaJustificacion:
		return "Devuelta para justificación"
	case EmergencyReferralStatusDevueltaASeguimiento:
		return "Devuelta a seguimiento"
	default:
		return status
	}
}

// EmergencyReferralStatusKnown dice si el código pertenece al catálogo de la bandeja.
func EmergencyReferralStatusKnown(status string) bool {
	return EmergencyReferralStatusLabel(status) != status || status == EmergencyReferralStatusPendienteDeRevision
}

// EmergencyReferralUrgencyLabel es el texto del badge de urgencia.
func EmergencyReferralUrgencyLabel(urgency string) string {
	switch urgency {
	case EmergencyReferralUrgencyRiesgoInminente:
		return "Riesgo inminente"
	case EmergencyReferralUrgencyNormal:
		return "Normal"
	default:
		return urgency
	}
}

// EmergencyReferralReviewLabel es el texto del concepto técnico.
func EmergencyReferralReviewLabel(result string) string {
	switch result {
	case EmergencyReferralReviewAprobadoCaracterizacion:
		return "Aprobado para caracterización ME"
	case EmergencyReferralReviewDevueltoAjuste:
		return "Devuelto para ajuste"
	case EmergencyReferralReviewDevueltoAcompanamiento:
		return "Devuelto para acompañamiento"
	default:
		return result
	}
}

// EmergencyReferralStatusPorRevision mapea el resultado de E-13 al estado nuevo.
func EmergencyReferralStatusPorRevision(result string) (string, bool) {
	switch result {
	case EmergencyReferralReviewAprobadoCaracterizacion:
		return EmergencyReferralStatusPorAsignar, true
	case EmergencyReferralReviewDevueltoAjuste:
		return EmergencyReferralStatusDevueltaParaJustificacion, true
	case EmergencyReferralReviewDevueltoAcompanamiento:
		return EmergencyReferralStatusDevueltaASeguimiento, true
	default:
		return "", false
	}
}

// EmergencyReferralStatusesRiesgo son los cuatro estados de la card del supervisor.
func EmergencyReferralStatusesRiesgo() []string {
	return []string{
		EmergencyReferralStatusRiesgoPorDecidirVia,
		EmergencyReferralStatusRiesgoPorEscalarCruzRoja,
		EmergencyReferralStatusRiesgoPorEntregarMedida,
		EmergencyReferralStatusRiesgoPorConfirmarSeguridad,
	}
}

// EmergencyReferralStatusesAgente es el recorte y el filtro de la bandeja del agente.
func EmergencyReferralStatusesAgente() []string {
	return []string{
		EmergencyReferralStatusRiesgoPorConfirmarSeguridad,
		EmergencyReferralStatusAsignada,
		EmergencyReferralStatusEnGestionDocumental,
		EmergencyReferralStatusAprobadas,
		EmergencyReferralStatusEnSeguimiento,
	}
}

// EmergencyReferralStatusesConvenio es el recorte de Supervisión del Convenio.
func EmergencyReferralStatusesConvenio() []string {
	return []string{
		EmergencyReferralStatusPorRevisarTrazabilidad,
		EmergencyReferralStatusPorAprobarCotizaciones,
		EmergencyReferralStatusRiesgoPorEscalarCruzRoja,
	}
}

// EmergencyReferralStatusesCruzRoja es el recorte de Cruz Roja.
func EmergencyReferralStatusesCruzRoja() []string {
	return []string{
		EmergencyReferralStatusAprobadas,
		EmergencyReferralStatusPorCotizar,
		EmergencyReferralStatusRiesgoPorEntregarMedida,
	}
}
