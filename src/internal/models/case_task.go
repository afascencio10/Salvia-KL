package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// CaseTask representa una tarea asociada a un caso dentro del sistema Salvia.
// Puede estar relacionada con una barrera, oficio, medida de emergencia,
// apoyo psicosocial o estabilización económica.
//
// Estados posibles:
//
//	"ToDo" → tarea pendiente   (tab "Por Articular" en Mis Barreras)
//	"Done" → tarea completada  (tab "Articuladas" en Mis Barreras)
type CaseTask struct {
	ID          string `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	// Clasificación
	Category string `gorm:"type:varchar(100)" json:"category"`
	Type     string `gorm:"type:varchar(100)" json:"type"`

	// Contenido
	Description string     `gorm:"type:text"                             json:"description"`
	Result      *string    `gorm:"type:text"                             json:"result,omitempty"`
	CompletedAt *time.Time `gorm:"column:completed_at"                   json:"completedAt,omitempty"`

	// Asignación y estado
	AssignedUserID string         `gorm:"type:varchar(36);not null;index;column:assigned_user_id" json:"assignedUserId"`
	Status         string         `gorm:"type:varchar(10);not null;default:'ToDo'"                 json:"status"`
	FormData       datatypes.JSON `gorm:"type:jsonb;column:form_data"                              json:"formData,omitempty"`

	// Relaciones (todas opcionales)
	CaseID                  string  `gorm:"type:varchar(36);not null;index"                        json:"caseId"`
	FollowUpID              *string `gorm:"type:uuid;column:follow_up_id"                          json:"followUpId,omitempty"`
	EntityLetterID          *string `gorm:"type:uuid;column:entity_letter_id"                      json:"entityLetterId,omitempty"`
	BarrierID               *string `gorm:"type:uuid;column:barrier_id;index"                      json:"barrierId,omitempty"`
	EmergencyMeasureID      *string `gorm:"type:uuid;column:emergency_measure_id"                  json:"emergencyMeasureId,omitempty"`
	PsychosocialSupportID   *string `gorm:"type:uuid;column:psychosocial_support_id"               json:"psychosocialSupportId,omitempty"`
	EconomicStabilizationID *string `gorm:"type:uuid;column:economic_stabilization_id"             json:"economicStabilizationId,omitempty"`
	EntityCaseID            *string `gorm:"type:varchar(36);column:entity_case_id;index"           json:"entityCaseId,omitempty"`

	// Timestamps
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (CaseTask) TableName() string {
	return "salvia.case_task"
}

// Estados de CaseTask.
const (
	CaseTaskStatusToDo = "ToDo"
	CaseTaskStatusDone = "Done"
)
