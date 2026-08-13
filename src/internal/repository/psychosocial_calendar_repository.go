package repository

import (
	"bitsflow/internal/models"
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// CalendarEvent evento normalizado para el frontend (sesión con víctima o reunión de equipo).
type CalendarEvent struct {
	ID            string    `json:"id"                       gorm:"column:id"`
	Kind          string    `json:"kind"                     gorm:"column:kind"` // "session" | "meeting"
	Title         string    `json:"title"                    gorm:"column:title"`
	Start         time.Time `json:"start"                    gorm:"column:start"`
	End           time.Time `json:"end"                      gorm:"column:end"`
	ScheduledTime string    `json:"scheduledTime"            gorm:"column:scheduled_time"`
	ScheduledDate string    `json:"scheduledDate"            gorm:"column:scheduled_date_str"`
	CaseICode     string    `json:"caseICode,omitempty"      gorm:"column:case_icode"`
	VictimName    string    `json:"victimName,omitempty"     gorm:"column:victim_name"`
	AgentID       string    `json:"agentId,omitempty"        gorm:"column:agent_id"`
	AgentName     string    `json:"agentName,omitempty"      gorm:"column:agent_name"`
	RiskLevel     int       `json:"riskLevel"                gorm:"column:risk_level"`
	Status        string    `json:"status,omitempty"         gorm:"column:status"`
	Description   string    `json:"description,omitempty"    gorm:"column:description"`
}

// CalendarAgentOption ítem para dropdown filtro por agente (supervisor).
type CalendarAgentOption struct {
	ID   string `json:"id"   gorm:"column:id"`
	Name string `json:"name" gorm:"column:name"`
}

// CalendarVictimOption víctima con remisión activa del agente (para programar sesión).
type CalendarVictimOption struct {
	CaseICode    string `json:"caseICode"    gorm:"column:case_icode"`
	FollowUpID   string `json:"followUpId"   gorm:"column:follow_up_id"`
	VictimName   string `json:"victimName"   gorm:"column:victim_name"`
	RiskLevel    int    `json:"riskLevel"    gorm:"column:risk_level"`
	RemissionID  string `json:"remissionId"  gorm:"column:remission_id"`
}

// CreateSessionInput registro de sesión con víctima.
type CreateSessionInput struct {
	AgentID       string
	AgentName     string
	CaseICode     string
	FollowUpID    string
	RemissionID   string
	ScheduledDate time.Time
	ScheduledTime string
	Summary       *string
}

// CreateMeetingInput reunión de equipo con N agentes.
type CreateMeetingInput struct {
	Title         string
	Description   *string
	ScheduledDate time.Time
	ScheduledTime string
	DurationMin   int
	CreatedBy     string
	CreatedByName string
	AgentIDs      []string
}

// MeetingDetail info completa de reunión + agentes.
type MeetingDetail struct {
	ID            string             `json:"id"`
	Title         string             `json:"title"`
	Description   string             `json:"description"`
	ScheduledDate time.Time          `json:"scheduledDate"`
	ScheduledTime string             `json:"scheduledTime"`
	DurationMin   int                `json:"durationMin"`
	AgentIDs      []string           `json:"agentIds"`
	Agents        []MeetingAgentInfo `json:"agents"`
}

// MeetingAgentInfo — id + nombre resuelto.
type MeetingAgentInfo struct {
	ID   string `json:"id"   gorm:"column:id"`
	Name string `json:"name" gorm:"column:name"`
}

// PsychosocialCalendarRepository acceso a datos del calendario psicosocial.
type PsychosocialCalendarRepository interface {
	MineEvents(ctx context.Context, agentID string, from, to time.Time) ([]CalendarEvent, error)
	TeamEvents(ctx context.Context, filterAgentID string, from, to time.Time) ([]CalendarEvent, error)
	ListAgents(ctx context.Context) ([]CalendarAgentOption, error)
	MyVictims(ctx context.Context, agentID string) ([]CalendarVictimOption, error)
	CreateSession(ctx context.Context, in CreateSessionInput) (string, error)
	UpdateSession(ctx context.Context, id string, in CreateSessionInput) error
	DeleteSession(ctx context.Context, id string) error
	CreateMeeting(ctx context.Context, in CreateMeetingInput) (string, error)
	GetMeeting(ctx context.Context, id string) (MeetingDetail, error)
	UpdateMeeting(ctx context.Context, id string, in CreateMeetingInput) error
	DeleteMeeting(ctx context.Context, id string) error
}

type psychosocialCalendarRepository struct {
	db *gorm.DB
}

func NewPsychosocialCalendarRepository(db *gorm.DB) PsychosocialCalendarRepository {
	return &psychosocialCalendarRepository{db: db}
}

// ─── Query base sesiones ──────────────────────────────────────────────────────
// Fuente de verdad: team_contact con is_psico_session=true (sesiones dentro de remisión psicosocial).
const calSessionsFrom = `
FROM salvia.team_contact tc
LEFT JOIN salvia.psychosocial_support ps
       ON ps.id::text = tc.psicosocial_id
LEFT JOIN salvia.victim_case vc
       ON vc.victim_case_i_code = BTRIM(COALESCE(ps.case_id, tc.case_id)::text)
LEFT JOIN salvia.victim_case_form2 vf2
       ON vf2.victim_case_form2_victim_case = vc.victim_case_id
LEFT JOIN security.general_user prof_gu
       ON prof_gu.general_user_id::text = BTRIM(COALESCE(tc.professional_id, ''))
LEFT JOIN security.general_user_profile prof_gup
       ON prof_gup.general_user_profile_id = prof_gu.general_user_general_user_profile
`

const calSessionsWhereBase = `
	tc.deleted_at IS NULL
	AND tc.is_psico_session = true
	AND tc.scheduled_date IS NOT NULL
`

const calSessionsSelect = `
	tc.id                              AS id,
	'session'                          AS kind,
	COALESCE(NULLIF(TRIM(CONCAT(vc.victim_case_victim_names, ' ', vc.victim_case_victim_last_names)), ''), 'Víctima') AS title,
	tc.scheduled_date                  AS start,
	(tc.scheduled_date + interval '60 minutes') AS end,
	COALESCE(tc.scheduled_time, '')    AS scheduled_time,
	TO_CHAR(tc.scheduled_date, 'YYYY-MM-DD') AS scheduled_date_str,
	COALESCE(vc.victim_case_i_code, tc.case_id, '') AS case_icode,
	COALESCE(NULLIF(TRIM(CONCAT(vc.victim_case_victim_names, ' ', vc.victim_case_victim_last_names)), ''), '') AS victim_name,
	COALESCE(prof_gu.general_user_i_code::text, tc.professional_id, '') AS agent_id,
	COALESCE(TRIM(CONCAT(prof_gup.general_user_profile_names, ' ', prof_gup.general_user_profile_last_names)), '') AS agent_name,
	COALESCE(vf2.victim_case_form2_risk_level, 0) AS risk_level,
	CASE WHEN tc.is_completed THEN 'realizada' ELSE COALESCE(tc.status, 'programada') END AS status,
	COALESCE(tc.summary, '')           AS description
`

// ─── Query base reuniones ─────────────────────────────────────────────────────
// LEFT JOIN opcional para filtrar por agente sin duplicar cuando no hay filtro
const calMeetingsFrom = `
FROM salvia.team_meeting tm
LEFT JOIN salvia.team_meeting_agent tma
       ON tma.meeting_id = tm.id
      AND tma.deleted_at IS NULL
LEFT JOIN security.general_user tma_gu
       ON tma_gu.general_user_id::text = tma.agent_id
`

const calMeetingsSelect = `
	tm.id                              AS id,
	'meeting'                          AS kind,
	tm.title                           AS title,
	(tm.scheduled_date + COALESCE(NULLIF(tm.scheduled_time,'')::time, '00:00'::time)) AS start,
	(tm.scheduled_date + COALESCE(NULLIF(tm.scheduled_time,'')::time, '00:00'::time) + (COALESCE(tm.duration_min, 60) * interval '1 minute')) AS end,
	COALESCE(tm.scheduled_time, '')    AS scheduled_time,
	TO_CHAR(tm.scheduled_date, 'YYYY-MM-DD') AS scheduled_date_str,
	''                                 AS case_icode,
	''                                 AS victim_name,
	COALESCE(tma_gu.general_user_i_code::text, tma.agent_id, '') AS agent_id,
	''                                 AS agent_name,
	0                                  AS risk_level,
	'programada'                       AS status,
	COALESCE(tm.description, '')       AS description
`

// Igual que calMeetingsSelect pero sin JOIN a tma (agent_id vacío) — usado en TeamEvents sin filtro.
const calMeetingsSelectNoAgent = `
	tm.id                              AS id,
	'meeting'                          AS kind,
	tm.title                           AS title,
	(tm.scheduled_date + COALESCE(NULLIF(tm.scheduled_time,'')::time, '00:00'::time)) AS start,
	(tm.scheduled_date + COALESCE(NULLIF(tm.scheduled_time,'')::time, '00:00'::time) + (COALESCE(tm.duration_min, 60) * interval '1 minute')) AS end,
	COALESCE(tm.scheduled_time, '')    AS scheduled_time,
	TO_CHAR(tm.scheduled_date, 'YYYY-MM-DD') AS scheduled_date_str,
	''                                 AS case_icode,
	''                                 AS victim_name,
	''                                 AS agent_id,
	''                                 AS agent_name,
	0                                  AS risk_level,
	'programada'                       AS status,
	COALESCE(tm.description, '')       AS description
`

// tzBogota — todas las fechas se interpretan en zona horaria Colombia para el cast a ::date.
const tzBogota = "America/Bogota"

// userUUIDLookup — resuelve el UUID del general_user desde el i_code. Usado en filtros de agente.
const userUUIDLookup = `(SELECT gu.general_user_id::text FROM security.general_user gu WHERE gu.general_user_i_code::text = ? LIMIT 1)`

func (r *psychosocialCalendarRepository) MineEvents(ctx context.Context, agentICode string, from, to time.Time) ([]CalendarEvent, error) {
	// agentICode = i_code (viene de la sesión). Resolvemos a UUID via subquery.
	// Aceptamos coincidencia contra UUID (nuevo estándar) o contra i_code (rows legacy).
	sql := `
	SELECT ` + calSessionsSelect + `
	` + calSessionsFrom + `
	WHERE ` + calSessionsWhereBase + `
	  AND (tc.scheduled_date AT TIME ZONE '` + tzBogota + `')::date BETWEEN ? AND ?
	  AND tc.professional_id IN (` + userUUIDLookup + `, ?)
	UNION ALL
	SELECT ` + calMeetingsSelect + `
	` + calMeetingsFrom + `
	WHERE tm.deleted_at IS NULL
	  AND tm.scheduled_date BETWEEN ? AND ?
	  AND tma.agent_id IN (` + userUUIDLookup + `, ?)
	ORDER BY start ASC
	`
	var rows []CalendarEvent
	if err := r.db.WithContext(ctx).Raw(sql, from, to, agentICode, agentICode, from, to, agentICode, agentICode).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *psychosocialCalendarRepository) TeamEvents(ctx context.Context, filterAgentICode string, from, to time.Time) ([]CalendarEvent, error) {
	sessionWhere := ""
	args := []interface{}{from, to}
	if filterAgentICode != "" {
		sessionWhere = " AND tc.professional_id IN (" + userUUIDLookup + ", ?)"
		args = append(args, filterAgentICode, filterAgentICode)
	}
	args = append(args, from, to)

	// Sin filtro por agente: mostrar reuniones una vez (sin JOIN a tma).
	// Con filtro: reuniones donde el agente participa.
	meetingsQuery := ""
	if filterAgentICode == "" {
		meetingsQuery = `
		SELECT ` + calMeetingsSelectNoAgent + `
		FROM salvia.team_meeting tm
		WHERE tm.deleted_at IS NULL
		  AND tm.scheduled_date BETWEEN ? AND ?
		`
	} else {
		meetingsQuery = `
		SELECT ` + calMeetingsSelect + `
		` + calMeetingsFrom + `
		WHERE tm.deleted_at IS NULL
		  AND tm.scheduled_date BETWEEN ? AND ?
		  AND tma.agent_id IN (` + userUUIDLookup + `, ?)
		`
		args = append(args, filterAgentICode, filterAgentICode)
	}

	sql := `
	SELECT * FROM (
		SELECT ` + calSessionsSelect + `
		` + calSessionsFrom + `
		WHERE ` + calSessionsWhereBase + `
		  AND (tc.scheduled_date AT TIME ZONE '` + tzBogota + `')::date BETWEEN ? AND ?
		  ` + sessionWhere + `
		UNION ALL
		` + meetingsQuery + `
	) e
	ORDER BY start ASC
	`
	var rows []CalendarEvent
	if err := r.db.WithContext(ctx).Raw(sql, args...).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *psychosocialCalendarRepository) ListAgents(ctx context.Context) ([]CalendarAgentOption, error) {
	sql := `
	SELECT DISTINCT
		gu.general_user_i_code AS id,
		TRIM(CONCAT(gup.general_user_profile_names, ' ', gup.general_user_profile_last_names)) AS name
	FROM security.general_user gu
	LEFT JOIN security.general_user_profile gup ON gup.general_user_profile_id = gu.general_user_general_user_profile
	LEFT JOIN security.rel_role_general_user rr ON rr.general_user_id = gu.general_user_id
	LEFT JOIN security.role r ON r.role_id = rr.role_id
	WHERE (LOWER(r.role_name) LIKE '%psicolog%'
	    OR LOWER(r.role_code) IN ('psicologia','ps'))
	ORDER BY name ASC
	`
	var rows []CalendarAgentOption
	if err := r.db.WithContext(ctx).Raw(sql).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *psychosocialCalendarRepository) MyVictims(ctx context.Context, agentICode string) ([]CalendarVictimOption, error) {
	sql := `
	SELECT
		ps.case_id                                                                     AS case_icode,
		ps.follow_up_id                                                                AS follow_up_id,
		ps.id                                                                          AS remission_id,
		TRIM(CONCAT(vc.victim_case_victim_names, ' ', vc.victim_case_victim_last_names)) AS victim_name,
		COALESCE(vf2.victim_case_form2_risk_level, 0)                                  AS risk_level
	FROM salvia.psychosocial_support ps
	LEFT JOIN salvia.victim_case vc          ON vc.victim_case_i_code = BTRIM(ps.case_id::text)
	LEFT JOIN salvia.victim_case_form2 vf2   ON vf2.victim_case_form2_victim_case = vc.victim_case_id
	WHERE ps.deleted_at IS NULL
	  AND ps.status <> 'cerrado'
	  AND ps.professional_id IN (` + userUUIDLookup + `, ?)
	ORDER BY vf2.victim_case_form2_risk_level DESC NULLS LAST, victim_name ASC
	`
	var rows []CalendarVictimOption
	if err := r.db.WithContext(ctx).Raw(sql, agentICode, agentICode).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// resolveUserUUID convierte un i_code a general_user_id (UUID). Si ya es UUID válido lo devuelve tal cual.
func (r *psychosocialCalendarRepository) resolveUserUUID(ctx context.Context, icodeOrUUID string) string {
	if icodeOrUUID == "" {
		return ""
	}
	var uuid string
	err := r.db.WithContext(ctx).Raw(
		`SELECT general_user_id::text FROM security.general_user WHERE general_user_i_code::text = ? OR general_user_id::text = ? LIMIT 1`,
		icodeOrUUID, icodeOrUUID,
	).Scan(&uuid).Error
	if err != nil || uuid == "" {
		return icodeOrUUID
	}
	return uuid
}

func (r *psychosocialCalendarRepository) CreateSession(ctx context.Context, in CreateSessionInput) (string, error) {
	scheduledAt := combineDateTime(in.ScheduledDate, in.ScheduledTime)
	agentUUID := r.resolveUserUUID(ctx, in.AgentID)
	team := "psicosocial"
	status := "programada"
	psicosocialID := in.RemissionID
	rec := models.TeamContact{
		CaseID:         in.CaseICode,
		PsicosocialID:  &psicosocialID,
		ProfessionalID: &agentUUID,
		Team:           &team,
		ScheduledDate:  &scheduledAt,
		ScheduledTime:  &in.ScheduledTime,
		IsPsicoSession: true,
		IsCompleted:    false,
		Status:         &status,
		Summary:        in.Summary,
	}
	if err := r.db.WithContext(ctx).Create(&rec).Error; err != nil {
		return "", err
	}
	return rec.ID, nil
}

func combineDateTime(date time.Time, hhmm string) time.Time {
	h, m := 0, 0
	if len(hhmm) >= 5 {
		_, _ = fmt.Sscanf(hhmm, "%d:%d", &h, &m)
	}
	return time.Date(date.Year(), date.Month(), date.Day(), h, m, 0, 0, date.Location())
}

func (r *psychosocialCalendarRepository) UpdateSession(ctx context.Context, id string, in CreateSessionInput) error {
	scheduledAt := combineDateTime(in.ScheduledDate, in.ScheduledTime)
	updates := map[string]interface{}{
		"scheduled_date": scheduledAt,
		"scheduled_time": in.ScheduledTime,
	}
	if in.AgentID != "" {
		updates["professional_id"] = r.resolveUserUUID(ctx, in.AgentID)
	}
	if in.Summary != nil {
		updates["summary"] = in.Summary
	}
	return r.db.WithContext(ctx).Model(&models.TeamContact{}).
		Where("id = ? AND is_completed = false", id).
		Updates(updates).Error
}

func (r *psychosocialCalendarRepository) DeleteSession(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND is_completed = false", id).
		Delete(&models.TeamContact{}).Error
}

func (r *psychosocialCalendarRepository) GetMeeting(ctx context.Context, id string) (MeetingDetail, error) {
	var meeting models.TeamMeeting
	if err := r.db.WithContext(ctx).First(&meeting, "id = ?", id).Error; err != nil {
		return MeetingDetail{}, err
	}
	// Devolvemos id (i_code) + nombre resuelto para render directo en frontend.
	var agentsInfo []MeetingAgentInfo
	if err := r.db.WithContext(ctx).Raw(`
		SELECT
			COALESCE(gu.general_user_i_code::text, tma.agent_id) AS id,
			COALESCE(NULLIF(TRIM(CONCAT(gup.general_user_profile_names, ' ', gup.general_user_profile_last_names)), ''), tma.agent_id) AS name
		FROM salvia.team_meeting_agent tma
		LEFT JOIN security.general_user gu
		       ON gu.general_user_id::text = BTRIM(tma.agent_id)
		       OR gu.general_user_i_code::text = BTRIM(tma.agent_id)
		LEFT JOIN security.general_user_profile gup
		       ON gup.general_user_profile_id = gu.general_user_general_user_profile
		WHERE tma.meeting_id = ?
		  AND tma.deleted_at IS NULL
	`, id).Scan(&agentsInfo).Error; err != nil {
		return MeetingDetail{}, err
	}
	ids := make([]string, 0, len(agentsInfo))
	for _, a := range agentsInfo {
		ids = append(ids, a.ID)
	}
	desc := ""
	if meeting.Description != nil {
		desc = *meeting.Description
	}
	return MeetingDetail{
		ID:            meeting.ID,
		Title:         meeting.Title,
		Description:   desc,
		ScheduledDate: meeting.ScheduledDate,
		ScheduledTime: meeting.ScheduledTime,
		DurationMin:   meeting.DurationMin,
		AgentIDs:      ids,
		Agents:        agentsInfo,
	}, nil
}

func (r *psychosocialCalendarRepository) UpdateMeeting(ctx context.Context, id string, in CreateMeetingInput) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		duration := in.DurationMin
		if duration <= 0 {
			duration = 60
		}
		updates := map[string]interface{}{
			"title":          in.Title,
			"description":    in.Description,
			"scheduled_date": in.ScheduledDate,
			"scheduled_time": in.ScheduledTime,
			"duration_min":   duration,
		}
		if err := tx.Model(&models.TeamMeeting{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Where("meeting_id = ?", id).Delete(&models.TeamMeetingAgent{}).Error; err != nil {
			return err
		}
		for _, aid := range in.AgentIDs {
			if aid == "" {
				continue
			}
			if err := tx.Create(&models.TeamMeetingAgent{MeetingID: id, AgentID: r.resolveUserUUID(ctx, aid)}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *psychosocialCalendarRepository) DeleteMeeting(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("meeting_id = ?", id).Delete(&models.TeamMeetingAgent{}).Error; err != nil {
			return err
		}
		return tx.Delete(&models.TeamMeeting{}, "id = ?", id).Error
	})
}

func (r *psychosocialCalendarRepository) CreateMeeting(ctx context.Context, in CreateMeetingInput) (string, error) {
	var meetingID string
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		duration := in.DurationMin
		if duration <= 0 {
			duration = 60
		}
		createdByName := in.CreatedByName
		meeting := models.TeamMeeting{
			Title:         in.Title,
			Description:   in.Description,
			ScheduledDate: in.ScheduledDate,
			ScheduledTime: in.ScheduledTime,
			DurationMin:   duration,
			CreatedBy:     in.CreatedBy,
			CreatedByName: &createdByName,
		}
		if err := tx.Create(&meeting).Error; err != nil {
			return err
		}
		meetingID = meeting.ID
		for _, aid := range in.AgentIDs {
			if aid == "" {
				continue
			}
			link := models.TeamMeetingAgent{MeetingID: meeting.ID, AgentID: r.resolveUserUUID(ctx, aid)}
			if err := tx.Create(&link).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return meetingID, nil
}
