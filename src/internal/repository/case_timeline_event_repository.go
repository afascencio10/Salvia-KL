package repository

import (
	"bitsflow/internal/models"
	"context"
	"log"
	"time"

	"gorm.io/gorm"
)

// TimelineEventRow extiende CaseTimelineEvent con el nombre del actor resuelto por JOIN.
type TimelineEventRow struct {
	ID                      string    `json:"id"`
	CaseID                  string    `json:"case_id"`
	Category                string    `json:"category"`
	Type                    string    `json:"type"`
	Icon                    string    `json:"icon"`
	Date                    time.Time `json:"date"`
	Description             string    `json:"description"`
	EventUserID             string    `json:"event_user_id"`
	ActorName               string    `json:"actor_name"`
	ActorFullName           string    `json:"actor_full_name"`
	Color                   string    `json:"color"`
	FollowUpID              string    `json:"follow_up_id"`
	BarrierID               string    `json:"barrier_id"`
	EmergencyMeasureID      string    `json:"emergency_measure_id"`
	PsychosocialSupportID   string    `json:"psychosocial_support_id"`
	EconomicStabilizationID string    `json:"economic_stabilization_id"`
	CreatedAt               time.Time `json:"created_at"`
}

// CaseTimelineEventRepository gestiona la persistencia de eventos del timeline de un caso.
// Es append-only: solo se crean registros, nunca se modifican ni eliminan.
type CaseTimelineEventRepository interface {
	Create(ctx context.Context, event *models.CaseTimelineEvent) error
	// GetByCaseID retorna eventos del caso en orden cronológico descendente.
	// Si barrierID no es vacío filtra además por ese barrier_id.
	// Si psychosocialID no es vacío filtra además por ese psychosocial_support_id.
	GetByCaseID(ctx context.Context, caseID string, barrierID string, psychosocialID ...string) ([]TimelineEventRow, error)
	// GetByVictimContactID retorna los eventos de un reporte (icode de
	// victim_contact), el más reciente primero. WP-1 T5.
	GetByVictimContactID(ctx context.Context, reportICode string) ([]TimelineEventRow, error)
}

type caseTimelineEventRepository struct {
	db *gorm.DB
}

func NewCaseTimelineEventRepository(db *gorm.DB) CaseTimelineEventRepository {
	return &caseTimelineEventRepository{db: db}
}

func (r *caseTimelineEventRepository) Create(ctx context.Context, event *models.CaseTimelineEvent) error {
	return r.db.WithContext(ctx).Create(event).Error
}

// timelineEventSelectFrom es el SELECT común del timeline (caso o reporte): resuelve
// el nombre del actor por JOIN con general_user y general_user_profile. case_id
// puede ser NULL en los eventos de reporte, por eso el COALESCE.
const timelineEventSelectFrom = `
		SELECT
			cte.id,
			COALESCE(cte.case_id, '')                 AS case_id,
			COALESCE(cte.category, '')                AS category,
			COALESCE(cte.type, cte.event_type, '')    AS type,
			COALESCE(cte.icon, '')                    AS icon,
			COALESCE(cte.date, cte.created_at)        AS date,
			COALESCE(cte.description, '')             AS description,
			COALESCE(cte.event_user_id, '')           AS event_user_id,
			COALESCE(cte.actor_name, '')              AS actor_name,
			COALESCE(
				gup.general_user_profile_names || ' ' || gup.general_user_profile_last_names,
				cte.actor_name,
				''
			)                                         AS actor_full_name,
			COALESCE(cte.color, '')                   AS color,
			COALESCE(cte.follow_up_id, '')            AS follow_up_id,
			COALESCE(cte.barrier_id, '')              AS barrier_id,
			COALESCE(cte.emergency_measure_id, '')    AS emergency_measure_id,
			COALESCE(cte.psychosocial_support_id, '') AS psychosocial_support_id,
			COALESCE(cte.economic_stabilization_id,'')AS economic_stabilization_id,
			cte.created_at
		FROM salvia.case_timeline_event cte
		LEFT JOIN security.general_user gu
			ON gu.general_user_i_code = cte.event_user_id
		LEFT JOIN security.general_user_profile gup
			ON gup.general_user_profile_id = gu.general_user_general_user_profile
`

func (r *caseTimelineEventRepository) GetByCaseID(ctx context.Context, caseID string, barrierID string, psychosocialID ...string) ([]TimelineEventRow, error) {
	q := timelineEventSelectFrom + `
		WHERE cte.case_id = ?
		  AND cte.deleted_at IS NULL
	`
	args := []interface{}{caseID}
	if barrierID != "" {
		q += " AND cte.barrier_id = ?"
		args = append(args, barrierID)
	}
	if len(psychosocialID) > 0 && psychosocialID[0] != "" {
		q += " AND cte.psychosocial_support_id = ?"
		args = append(args, psychosocialID[0])
	}
	q += " ORDER BY COALESCE(cte.date, cte.created_at) DESC"

	var rows []TimelineEventRow
	if err := r.db.WithContext(ctx).Raw(q, args...).Scan(&rows).Error; err != nil {
		return nil, err
	}
	if rows == nil {
		rows = []TimelineEventRow{}
	}
	for _, r := range rows {
		log.Printf("[timeline] id=%s type=%q event_user_id=%q actor_name=%q actor_full_name=%q",
			r.ID, r.Type, r.EventUserID, r.ActorName, r.ActorFullName)
	}
	return rows, nil
}

// GetByVictimContactID retorna el timeline de un reporte, el más reciente primero.
func (r *caseTimelineEventRepository) GetByVictimContactID(ctx context.Context, reportICode string) ([]TimelineEventRow, error) {
	q := timelineEventSelectFrom + `
		WHERE cte.victim_contact_id = ?
		  AND cte.deleted_at IS NULL
		ORDER BY COALESCE(cte.date, cte.created_at) DESC
	`
	var rows []TimelineEventRow
	if err := r.db.WithContext(ctx).Raw(q, reportICode).Scan(&rows).Error; err != nil {
		return nil, err
	}
	if rows == nil {
		rows = []TimelineEventRow{}
	}
	return rows, nil
}
