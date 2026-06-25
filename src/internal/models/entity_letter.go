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
//	    priority              VARCHAR(30) NOT NULL DEFAULT 'normal',
//	    agent_id              VARCHAR(36),
//	    notification_user_id  VARCHAR(36),
//	    review_by             VARCHAR(36),
//	    radicado_by           VARCHAR(36),
//	    register_by           VARCHAR(36),
//	    entidad               VARCHAR(255),
//	    nivel                 VARCHAR(50),
//	    url_kofax             TEXT,
//	    created_at            TIMESTAMPTZ,
//	    updated_at            TIMESTAMPTZ,
//	    reason_correction     VARCHAR(255),
//	    deleted_at            TIMESTAMPTZ
//	);
type EntityLetter struct {
	ID                 string         `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()" json:"id"`
	BarrierID          string         `gorm:"type:varchar(36);not null;column:barrier_id"          json:"barrierId"`
	CaseID             string         `gorm:"type:varchar(36);not null;column:case_id"             json:"caseId"`
	State              string         `gorm:"type:varchar(30);not null;default:'por_proyectar'"    json:"state"`
	Priority           string         `gorm:"type:varchar(30);not null;default:'normal';column:priority" json:"priority"`
	AgentID            *string        `gorm:"type:varchar(36);column:agent_id"                     json:"agentId,omitempty"`
	NotificationUserID *string        `gorm:"type:varchar(36);column:notification_user_id"         json:"notificationUserId,omitempty"`
	ReviewBy           *string        `gorm:"type:varchar(36);column:review_by"                    json:"reviewBy,omitempty"`
	RadicadoBy         *string        `gorm:"type:varchar(36);column:radicado_by"                  json:"radicadoBy,omitempty"`
	RegisterBy         *string        `gorm:"type:varchar(36);column:register_by"                  json:"registerBy,omitempty"`
	Entidad            *string        `gorm:"type:varchar(255);column:entidad"                     json:"entidad,omitempty"`
	Nivel              *string        `gorm:"type:varchar(50);column:nivel"                        json:"nivel,omitempty"`
	UrlKofax           *string        `gorm:"type:text;column:url_kofax"                           json:"urlKofax,omitempty"`
	EntityBranchID     *int64         `gorm:"type:integer;column:entity_branch_id"                 json:"entityBranchId,omitempty"`
	DepartmentID       *string        `gorm:"type:varchar(20);column:department_id"                json:"departmentId,omitempty"`
	CityID             *string        `gorm:"type:varchar(20);column:city_id"                      json:"cityId,omitempty"`
	TownID             *string        `gorm:"type:varchar(20);column:town_id"                      json:"townId,omitempty"`
	OfficialDependency *string        `gorm:"type:text;column:official_dependency"                 json:"officialDependency,omitempty"`
	Subject            *string        `gorm:"type:varchar(255);column:subject"                     json:"subject,omitempty"`
	AsuntoRadicado     *string        `gorm:"type:varchar(255);column:asunto_radicado"             json:"asuntoRadicado,omitempty"`
	CorreoEntidad      *string        `gorm:"type:varchar(255);column:correo_entidad"              json:"correoEntidad,omitempty"`
	NumeroRadicado     *string        `gorm:"type:varchar(100);column:numero_radicado"             json:"numeroRadicado,omitempty"`
	ResponseDate       *time.Time     `gorm:"column:response_date"                                 json:"responseDate,omitempty"`
	CorreoRemitente    *string        `gorm:"type:varchar(255);column:correo_remitente"            json:"correoRemitente,omitempty"`
	AsuntoRespuesta    *string        `gorm:"type:varchar(255);column:asunto_respuesta"            json:"asuntoRespuesta,omitempty"`
	ResponseReviewBy   *string        `gorm:"type:varchar(36);column:response_review_by"           json:"responseReviewBy,omitempty"`
	ReasonCorrection   *string        `gorm:"type:varchar(255);column:reason_correction"           json:"reasonCorrection,omitempty"`
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
