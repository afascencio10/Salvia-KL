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
//	"ToDo" → tarea pendiente
//	"Done" → tarea completada
type CaseTask struct {
	ID string `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	// Clasificación
	Category string `gorm:"type:varchar(100)" json:"category"`
	Type     string `gorm:"type:varchar(100)" json:"type"`

	// Contenido
	Description string     `gorm:"type:text"                             json:"description"`
	Result      *string    `gorm:"type:text"                             json:"result,omitempty"`
	CompletedAt *time.Time `gorm:"column:completed_at"                   json:"completedAt,omitempty"`

	// Asignación y estado
	AssignedUserID string `gorm:"type:varchar(36);not null;index;column:assigned_user_id" json:"assignedUserId"`
	// AssignedRole es el role_code que puede ver y gestionar la tarea.
	// Vacío cuando la tarea sigue asignada a un usuario concreto.
	AssignedRole *string        `gorm:"type:varchar(10);index;column:assigned_role"             json:"assignedRole,omitempty"`
	Status       string         `gorm:"type:varchar(10);not null;default:'ToDo'"                 json:"status"`
	FormData     datatypes.JSON `gorm:"type:jsonb;column:form_data"                              json:"formData,omitempty"`

	// Relaciones (todas opcionales)
	CaseID                      string  `gorm:"type:varchar(36);not null;index"                        json:"caseId"`
	FollowUpID                  *string `gorm:"type:uuid;column:follow_up_id"                          json:"followUpId,omitempty"`
	EntityLetterID              *string `gorm:"type:uuid;column:entity_letter_id"                      json:"entityLetterId,omitempty"`
	BarrierID                   *string `gorm:"type:uuid;column:barrier_id;index"                      json:"barrierId,omitempty"`
	EmergencyMeasureID          *string `gorm:"type:uuid;column:emergency_measure_id"                       json:"emergencyMeasureId,omitempty"`
	EmergencyMeasureReferralID  *string `gorm:"type:uuid;column:emergency_measure_referral_id;index"        json:"emergencyMeasureReferralId,omitempty"`
	PsychosocialSupportID       *string `gorm:"type:uuid;column:psychosocial_support_id"               json:"psychosocialSupportId,omitempty"`
	EconomicStabilizationID     *string `gorm:"type:uuid;column:economic_stabilization_id"             json:"economicStabilizationId,omitempty"`
	EntityCaseID                *string `gorm:"type:varchar(36);column:entity_case_id;index"           json:"entityCaseId,omitempty"`
	CaseArticulatingMechanismID *string `gorm:"type:varchar(36);column:case_articulating_mechanism_id;index" json:"caseArticulatingMechanismId,omitempty"`

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

// Tipos de tarea usados por Mis Barreras (Enlace Territorial).
const (
	CaseTaskTypeArticulacionTerritorial = "articulacion_territorial"
	CaseTaskTypeCompromisoConEnlace     = "compromiso_con_enlace"
	CaseTaskTypeGestionMecanismo        = "gestion_mecanismo"
	CaseTaskCategoryMecanismos          = "Mecanismos"
)

// Comentarios del Equipo Técnico en el tab Derivaciones de Detalle del Caso.
// Cada comentario es una case_task Done: el texto va en Result, sin pasar por ToDo.
const (
	CaseTaskTypeComentarioEquipoTecnico = "comentario_equipo_tecnico"
	CaseTaskCategoryEquipoTecnico       = "Equipo Técnico"
)

// Tarea de riesgo inminente: la ve y la gestiona cualquier usuario del rol.
const (
	CaseTaskTypeDecidirViaRiesgoInminente = "decidir_via_riesgo_inminente"
	CaseTaskTypeIntentarContacto3x3       = "intentar_contacto_3x3"
	CaseTaskTypeEscalarCruzRoja           = "escalar_cruz_roja"
	CaseTaskTypeEntregarMedidaInmediata   = "entregar_medida_inmediata"
	CaseTaskTypeConfirmarLugarSeguro      = "confirmar_lugar_seguro"
	CaseTaskTypeRevisionEquipoTecnico     = "revision_equipo_tecnico"
	CaseTaskTypeAsignarAgenteMedidas      = "asignar_agente_medidas"
	CaseTaskTypeAjustarRemisionDevuelta   = "ajustar_remision_devuelta"
	CaseTaskTypeEntrevistaCaracterizacion = "entrevista_caracterizacion"
	CaseTaskCategoryMedidas               = "Medidas de Emergencia"
	CaseTaskRoleEquipoTecnicoMedidas      = "em"
	CaseTaskRoleSupervisorMedidas         = "sm"
	CaseTaskRoleSupervisionConvenio       = "sc"
	CaseTaskRoleCruzRoja                  = "cr"
	CaseTaskDescDecidirViaRiesgoInminente = "Decidir vía de atención — Riesgo inminente"
	CaseTaskDescIntentarContacto3x3       = "Intentar contacto (3×3)"
	CaseTaskDescEscalarCruzRoja           = "Escalar a Cruz Roja — Riesgo inminente"
	CaseTaskDescEntregarMedidaInmediata   = "Entregar medida inmediata (sin caracterización previa)"
	CaseTaskDescConfirmarLugarSeguro      = "Confirmar que la persona está en lugar seguro"
	CaseTaskDescRevisionEquipoTecnico     = "Revisión Equipo Técnico"
	CaseTaskDescAsignarAgenteMedidas      = "Asignar agente de medidas"
	CaseTaskDescAjustarRemisionDevuelta   = "Ajustar remisión devuelta por el Equipo Técnico"
	CaseTaskDescEntrevistaCaracterizacion = "Contactar víctima y realizar entrevista de caracterización"
)

// CaseTaskCommentItem es un comentario del Equipo Técnico — una case_task Done
// del tipo comentario_equipo_tecnico — con el nombre de quien lo creó resuelto.
type CaseTaskCommentItem struct {
	ID               string     `gorm:"column:id" json:"id"`
	Result           *string    `gorm:"column:result" json:"result"`
	CompletedAt      *time.Time `gorm:"column:completed_at" json:"completedAt"`
	AssignedUserName string     `gorm:"column:assigned_user_name" json:"assignedUserName"`
}
