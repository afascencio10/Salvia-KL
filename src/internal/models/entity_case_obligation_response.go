package models

import "time"

// EntityCaseObligationResponse es la respuesta de la ENTIDAD a una acción
// institucional concreta de un caso: el radicado con el que quedó su gestión, la
// ampliación de lo actuado y el estado general de su respuesta.
//
// Es la fuente única de la gestión de la entidad: la consumen tanto la pantalla
// "Mis Acciones Institucionales" (rol et) como la tarjeta de la entidad en la
// pestaña Gestión Institucional del detalle del caso, de modo que SALVIA y la
// entidad siempre vean lo mismo.
//
// No se confunde con EntityCaseObligation, que es append-only y guarda una foto
// del checklist por cada seguimiento: eso es lo que SALVIA observó en un momento
// dado. Esto es lo que la entidad reporta, y hay una sola fila por acción, que se
// actualiza si la entidad corrige su respuesta.
//
// Tampoco se confunde con EntityCase.Radicado, que es el radicado con el que la
// víctima ya había acudido a la entidad antes de que SALVIA interviniera.
//
// SQL equivalente:
//
//	CREATE TABLE salvia.entity_case_obligation_response (
//	    id                     VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid(),
//	    entity_case_id         VARCHAR(36) NOT NULL,
//	    entity_obligation_id   VARCHAR(36) NOT NULL,
//	    radicado               TEXT,
//	    ampliacion             TEXT,
//	    estado                 VARCHAR(30) NOT NULL,
//	    responded_by_id        VARCHAR(36) NOT NULL,
//	    created_at             TIMESTAMPTZ,
//	    updated_at             TIMESTAMPTZ,
//	    UNIQUE (entity_case_id, entity_obligation_id)
//	);
//
// Sin DeletedAt a propósito: una respuesta no se borra, se corrige. Así el índice
// único compuesto funciona con los tags de GORM y no hace falta crear a mano un
// índice único parcial como en entity_case.
type EntityCaseObligationResponse struct {
	ID string `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()" json:"id"`

	// EntityCaseID + EntityObligationID identifican la acción respondida. Solo se
	// admiten acciones del catálogo: las del checklist marcadas como "Otra" (sin
	// entity_obligation_id) no se exponen en la pantalla de la entidad.
	EntityCaseID       string `gorm:"type:varchar(36);not null;uniqueIndex:ux_ecor_caso_accion;column:entity_case_id"       json:"entityCaseId"`
	EntityObligationID string `gorm:"type:varchar(36);not null;uniqueIndex:ux_ecor_caso_accion;column:entity_obligation_id" json:"entityObligationId"`

	// Radicado es el número de radicado/NUNC/MP con el que quedó la gestión de la
	// entidad. Opcional: el modal lo pide "si aplica", porque no toda acción
	// institucional genera un radicado.
	Radicado *string `gorm:"type:text;column:radicado" json:"radicado,omitempty"`

	// Ampliacion es el detalle libre de la gestión realizada.
	Ampliacion *string `gorm:"type:text;column:ampliacion" json:"ampliacion,omitempty"`

	// Estado es el estado general de la respuesta institucional — ver
	// EntityCaseObligationEstadoLabels.
	Estado string `gorm:"type:varchar(30);not null;column:estado" json:"estado"`

	// RespondedByID es el usuario de la entidad que diligenció la respuesta.
	RespondedByID string `gorm:"type:varchar(36);not null;column:responded_by_id" json:"respondedById"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (EntityCaseObligationResponse) TableName() string {
	return "salvia.entity_case_obligation_response"
}

// Estado general de la respuesta institucional.
//
// ⚠️ Catálogo inicial: el ticket pide el campo pero no enumera los valores.
// Están acá, en un solo lugar, para que cambiarlos sea editar este mapa. No se
// reusaron los valores de la pregunta "¿Hubo respuesta institucional?" del
// seguimiento (respuesta_oficial / contacto_victima / sin_respuesta) porque esos
// los responde SALVIA sobre la entidad, y varios no tienen sentido dichos por la
// entidad sobre sí misma.
const (
	EntityCaseObligationEstadoCompletada       = "completada"
	EntityCaseObligationEstadoEnTramite        = "en_tramite"
	EntityCaseObligationEstadoFueraCompetencia = "fuera_competencia"
	EntityCaseObligationEstadoRequiereInfo     = "requiere_info"
)

// EntityCaseObligationEstadoLabels mapea código → etiqueta para la UI y es la
// fuente de verdad del catálogo (IsValidEntityCaseObligationEstado).
var EntityCaseObligationEstadoLabels = map[string]string{
	EntityCaseObligationEstadoCompletada:       "Gestión completada",
	EntityCaseObligationEstadoEnTramite:        "En trámite",
	EntityCaseObligationEstadoFueraCompetencia: "Fuera de la competencia de la entidad",
	EntityCaseObligationEstadoRequiereInfo:     "Requiere información adicional",
}

// IsValidEntityCaseObligationEstado indica si el estado es uno de los permitidos.
func IsValidEntityCaseObligationEstado(code string) bool {
	_, ok := EntityCaseObligationEstadoLabels[code]
	return ok
}

// EntityCaseObligationResponseSummary resume, para la tarjeta de la entidad en
// Gestión Institucional, lo que la propia entidad reportó desde "Mis Acciones
// Institucionales". Lo llena el servicio leyendo la misma tabla que esa pantalla:
// es lo que pide el DoD, una sola fuente para las dos vistas.
type EntityCaseObligationResponseSummary struct {
	// Total y Respondidas salen del mismo checklist que usa la pantalla de la
	// entidad, no de contar respuestas acá, para que los dos conteos no puedan
	// discrepar.
	Total       int `json:"total"`
	Respondidas int `json:"respondidas"`

	// El resto describe la respuesta más reciente de esa entidad en el caso.
	// ponytail: la tarjeta resume solo la última; si piden ver todas, van en la
	// pantalla de detalle, que ya existe.
	AccionLabel string     `json:"accionLabel,omitempty"`
	Estado      string     `json:"estado,omitempty"`
	EstadoLabel string     `json:"estadoLabel,omitempty"`
	Radicado    string     `json:"radicado,omitempty"`
	Ampliacion  string     `json:"ampliacion,omitempty"`
	RespondedAt *time.Time `json:"respondedAt,omitempty"`
}
