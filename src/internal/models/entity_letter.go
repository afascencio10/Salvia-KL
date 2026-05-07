package models

import (
	"time"

	"gorm.io/gorm"
)

// EntityLetter representa un oficio emitido por el sistema Salvia hacia una entidad externa
// en el marco de la gestión de barreras de acceso de una víctima.
//
// Flujo de estados:
//
//	por_proyectar  → para_revisar      (el agente de seguimiento proyecta y sube el oficio)
//	para_revisar   → aprobacion_juridica (el agente de notificaciones revisa y lo envía a jurídica)
//	para_revisar   → en_correccion     (jurídica o notificaciones devuelve para corrección)
//	en_correccion  → para_revisar      (el agente de seguimiento corrige y re-envía)
//	aprobacion_juridica → para_radicar  (el abogado aprueba)
//	para_radicar   → radicado          (el agente de notificaciones radica ante la entidad)
//	radicado       → respondido        (se registra la respuesta de la entidad)
//
// SQL equivalente:
//
//	CREATE TABLE salvia.entity_letter (
//	    id                    VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid(),
//	    barrier_id            VARCHAR(36) NOT NULL,
//	    case_id               VARCHAR(36) NOT NULL,
//	    state                 VARCHAR(30) NOT NULL DEFAULT 'por_proyectar',
//	    agent_id              VARCHAR(36),
//	    notification_user_id  VARCHAR(36),
//	    review_by             VARCHAR(36),
//	    radicado_by           VARCHAR(36),
//	    register_by           VARCHAR(36),
//	    created_at            TIMESTAMPTZ,
//	    updated_at            TIMESTAMPTZ,
//	    deleted_at            TIMESTAMPTZ
//	);
type EntityLetter struct {
	ID                 string         `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()" json:"id"`
	BarrierID          string         `gorm:"type:varchar(36);not null;column:barrier_id"          json:"barrierId"`
	CaseID             string         `gorm:"type:varchar(36);not null;column:case_id"             json:"caseId"`
	State              string         `gorm:"type:varchar(30);not null;default:'por_proyectar'"    json:"state"`
	AgentID            *string        `gorm:"type:varchar(36);column:agent_id"                     json:"agentId,omitempty"`
	NotificationUserID *string        `gorm:"type:varchar(36);column:notification_user_id"         json:"notificationUserId,omitempty"`
	ReviewBy           *string        `gorm:"type:varchar(36);column:review_by"                    json:"reviewBy,omitempty"`
	RadicadoBy         *string        `gorm:"type:varchar(36);column:radicado_by"                  json:"radicadoBy,omitempty"`
	RegisterBy         *string        `gorm:"type:varchar(36);column:register_by"                  json:"registerBy,omitempty"`
	CreatedAt          time.Time      `                                                            json:"createdAt"`
	UpdatedAt          time.Time      `                                                            json:"updatedAt"`
	DeletedAt          gorm.DeletedAt `gorm:"index"                                                json:"deletedAt,omitempty"`
}

func (EntityLetter) TableName() string { return "salvia.entity_letter" }

// Estados válidos de un EntityLetter.
const (
	EntityLetterStatePorProyectar      = "por_proyectar"
	EntityLetterStateParaRevisar       = "para_revisar"
	EntityLetterStateEnCorreccion      = "en_correccion"
	EntityLetterStateAprobacionJuridica = "aprobacion_juridica"
	EntityLetterStateParaRadicar       = "para_radicar"
	EntityLetterStateRadicado          = "radicado"
	EntityLetterStateRespondido        = "respondido"
)
