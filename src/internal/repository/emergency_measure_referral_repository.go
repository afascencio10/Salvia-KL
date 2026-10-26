package repository

import (
	"bitsflow/internal/models"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// EmergencyReferralListFilter es el recorte de la vista más los filtros de UI.
// AgentICode es el icode de sesión y solo aplica en la vista agente.
// AgentID es el filtro del supervisor: un icode o "sin-asignar".
type EmergencyReferralListFilter struct {
	Vista      string
	AgentICode string
	Q          string
	Status     string
	AgentID    string
	Urgency    string
	ID         string
}

// EmergencyReferralRow es una fila ya unida al caso, al agente y a la revisión.
type EmergencyReferralRow struct {
	ID                     string
	CaseICode              string
	CaseName               string
	DocNumber              string
	RiskLabel              string
	Municipality           string
	Measures               []string
	Urgency                string
	SubmittedByName        string
	ReferredAt             string
	Status                 string
	AssignedAgentID        string
	AssignedAgentName      string
	ReviewProfessionalName string
	ReviewConcept          string
	ReviewResult           string
	ReviewDate             string
}

// EmergencyReferralDetailRow es la fila del detalle: la del listado más motivo,
// devolución, equipo del remitente e hijos e intérprete del caso.
type EmergencyReferralDetailRow struct {
	EmergencyReferralRow
	Reason              string
	RejectionReason     string
	SubmittedByTeam     string
	ChildrenNumber      *int
	ChildrenAges        string
	RequiresInterpreter bool
}

type emergencyReferralScan struct {
	ID                     string `gorm:"column:id"`
	CaseICode              string `gorm:"column:case_i_code"`
	CaseName               string `gorm:"column:case_name"`
	DocNumber              string `gorm:"column:doc_number"`
	RiskLabel              string `gorm:"column:risk_label"`
	Municipality           string `gorm:"column:municipality"`
	MeasuresJSON           string `gorm:"column:measures"`
	Urgency                string `gorm:"column:urgency"`
	SubmittedByName        string `gorm:"column:submitted_by_name"`
	ReferredAt             string `gorm:"column:referred_at"`
	Status                 string `gorm:"column:status"`
	AssignedAgentID        string `gorm:"column:assigned_agent_id"`
	AssignedAgentName      string `gorm:"column:assigned_agent_name"`
	ReviewProfessionalName string `gorm:"column:review_professional_name"`
	ReviewConcept          string `gorm:"column:review_concept"`
	ReviewResult           string `gorm:"column:review_result"`
	ReviewDate             string `gorm:"column:review_date"`
}

type emergencyReferralDetailScan struct {
	ID                     string        `gorm:"column:id"`
	CaseICode              string        `gorm:"column:case_i_code"`
	CaseName               string        `gorm:"column:case_name"`
	DocNumber              string        `gorm:"column:doc_number"`
	RiskLabel              string        `gorm:"column:risk_label"`
	Municipality           string        `gorm:"column:municipality"`
	MeasuresJSON           string        `gorm:"column:measures"`
	Urgency                string        `gorm:"column:urgency"`
	SubmittedByName        string        `gorm:"column:submitted_by_name"`
	ReferredAt             string        `gorm:"column:referred_at"`
	Status                 string        `gorm:"column:status"`
	AssignedAgentID        string        `gorm:"column:assigned_agent_id"`
	AssignedAgentName      string        `gorm:"column:assigned_agent_name"`
	ReviewProfessionalName string        `gorm:"column:review_professional_name"`
	ReviewConcept          string        `gorm:"column:review_concept"`
	ReviewResult           string        `gorm:"column:review_result"`
	ReviewDate             string        `gorm:"column:review_date"`
	Reason                 string        `gorm:"column:reason"`
	RejectionReason        string        `gorm:"column:rejection_reason"`
	SubmittedByTeam        string        `gorm:"column:submitted_by_team"`
	ChildrenNumber         sql.NullInt64 `gorm:"column:children_number"`
	ChildrenAges           string        `gorm:"column:children_ages"`
	InterpreterCode        string        `gorm:"column:interpreter_code"`
}

type statusCountRow struct {
	Status string `gorm:"column:status"`
	N      int    `gorm:"column:n"`
}

// EmergencyMeasureReferralRepository lee y actualiza la bandeja de medidas.
type EmergencyMeasureReferralRepository interface {
	CountByStatus(ctx context.Context, vista, agentICode string) (map[string]int, error)
	List(ctx context.Context, f EmergencyReferralListFilter) ([]EmergencyReferralRow, error)
	// FindDetail lee una remisión dentro del recorte de la vista. Nil si no entra.
	FindDetail(ctx context.Context, vista, agentICode, id string) (*EmergencyReferralDetailRow, error)
	UserHasActiveRole(ctx context.Context, iCode, roleCode string) (bool, error)
	Reassign(ctx context.Context, ids []string, agentID string) error
	Review(ctx context.Context, id, professionalID, concept, result, newStatus, reviewDate string) (found bool, pending bool, err error)
	// CreateWithMeasures inserta la cabecera y sus medidas en una transacción.
	// El id de la cabecera queda en referral.ID.
	CreateWithMeasures(ctx context.Context, referral *models.EmergencyMeasureReferral, measureTypes []string) error
}

type emergencyMeasureReferralRepository struct {
	db *gorm.DB
}

func NewEmergencyMeasureReferralRepository(db *gorm.DB) EmergencyMeasureReferralRepository {
	return &emergencyMeasureReferralRepository{db: db}
}

// CreateWithMeasures crea una remisión y una medida ACTIVE por cada tipo.
// Si no queda ningún tipo, no inserta nada.
func (r *emergencyMeasureReferralRepository) CreateWithMeasures(ctx context.Context, referral *models.EmergencyMeasureReferral, measureTypes []string) error {
	types := make([]string, 0, len(measureTypes))
	for _, tipo := range measureTypes {
		if tipo = strings.TrimSpace(tipo); tipo != "" {
			types = append(types, tipo)
		}
	}
	if referral == nil || strings.TrimSpace(referral.CaseID) == "" || len(types) == 0 {
		return fmt.Errorf("remisión sin caso o sin medidas")
	}
	if referral.ID == "" {
		referral.ID = uuid.NewString()
	}
	if referral.Status == "" {
		referral.Status = models.EmergencyReferralStatusPendienteDeRevision
	}
	if referral.Urgency == "" {
		referral.Urgency = models.EmergencyReferralUrgencyNormal
	}
	if referral.ReferredAt == nil {
		now := time.Now()
		day := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		referral.ReferredAt = &day
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(referral).Error; err != nil {
			return err
		}
		referralID := referral.ID
		for _, tipo := range types {
			measure := &models.EmergencyMeasure{
				CaseID:     referral.CaseID,
				FollowUpID: referral.FollowUpID,
				Type:       tipo,
				Status:     "ACTIVE",
				ReferralID: &referralID,
			}
			if err := tx.Create(measure).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

const referralListSelect = `
	r.id::text AS id,
	COALESCE(vc.victim_case_i_code, r.case_id, '') AS case_i_code,
	TRIM(CONCAT(COALESCE(vc.victim_case_victim_names, ''), ' ', COALESCE(vc.victim_case_victim_last_names, ''))) AS case_name,
	COALESCE(vc.victim_case_victim_doc_number, '') AS doc_number,
	CASE
		WHEN risk_level.level IS NOT NULL THEN
			CASE risk_level.level
				WHEN 1 THEN 'Bajo'
				WHEN 2 THEN 'Moderado'
				WHEN 3 THEN 'Alto'
				WHEN 4 THEN 'Extremo'
				ELSE ''
			END
		WHEN COALESCE(f1.victim_case_form1_victim_femicide_risk, '') = '1' THEN 'Alto (feminicidio)'
		ELSE ''
	END AS risk_label,
	COALESCE(t.town_name, '') AS municipality,
	COALESCE((
		SELECT json_agg(em.type ORDER BY em.created_at)
		FROM salvia.emergency_measure em
		WHERE em.referral_id = r.id AND em.deleted_at IS NULL
	), '[]'::json)::text AS measures,
	COALESCE(r.urgency, '') AS urgency,
	TRIM(CONCAT(COALESCE(submitter_gup.general_user_profile_names, ''), ' ', COALESCE(submitter_gup.general_user_profile_last_names, ''))) AS submitted_by_name,
	COALESCE(TO_CHAR(r.referred_at, 'YYYY-MM-DD'), '') AS referred_at,
	COALESCE(r.status, '') AS status,
	COALESCE(r.assigned_agent_id, '') AS assigned_agent_id,
	TRIM(CONCAT(COALESCE(agent_gup.general_user_profile_names, ''), ' ', COALESCE(agent_gup.general_user_profile_last_names, ''))) AS assigned_agent_name,
	TRIM(CONCAT(COALESCE(review_gup.general_user_profile_names, ''), ' ', COALESCE(review_gup.general_user_profile_last_names, ''))) AS review_professional_name,
	COALESCE(r.review_concept, '') AS review_concept,
	COALESCE(r.review_result, '') AS review_result,
	COALESCE(TO_CHAR(r.review_date, 'YYYY-MM-DD'), '') AS review_date
`

const referralListFrom = `
FROM salvia.emergency_measure_referral r
LEFT JOIN salvia.victim_case vc
       ON NULLIF(BTRIM(r.case_id::text), '') IS NOT NULL
      AND vc.victim_case_i_code = BTRIM(r.case_id::text)
LEFT JOIN LATERAL (
       SELECT vf2.victim_case_form2_risk_level AS level
       FROM salvia.victim_case_form2 vf2
       WHERE vc.victim_case_id IS NOT NULL
         AND vf2.victim_case_form2_victim_case = vc.victim_case_id
       ORDER BY vf2.victim_case_form2_id DESC
       LIMIT 1
) risk_level ON true
LEFT JOIN LATERAL (
       SELECT f1.victim_case_form1_victim_femicide_risk
       FROM salvia.victim_case_form1 f1
       WHERE vc.victim_case_id IS NOT NULL
         AND f1.victim_case_form1_victim_case = vc.victim_case_id
       ORDER BY f1.victim_case_form1_id DESC
       LIMIT 1
) f1 ON true
LEFT JOIN security.town t
       ON t.town_code = vc.victim_case_victim_town_code
LEFT JOIN security.general_user submitter_gu
       ON NULLIF(BTRIM(r.submitted_by::text), '') IS NOT NULL
      AND submitter_gu.general_user_i_code::text = BTRIM(r.submitted_by::text)
LEFT JOIN security.general_user_profile submitter_gup
       ON submitter_gu.general_user_i_code IS NOT NULL
      AND submitter_gup.general_user_profile_id = submitter_gu.general_user_general_user_profile
LEFT JOIN security.general_user agent_gu
       ON NULLIF(BTRIM(r.assigned_agent_id::text), '') IS NOT NULL
      AND agent_gu.general_user_i_code::text = BTRIM(r.assigned_agent_id::text)
LEFT JOIN security.general_user_profile agent_gup
       ON agent_gu.general_user_i_code IS NOT NULL
      AND agent_gup.general_user_profile_id = agent_gu.general_user_general_user_profile
LEFT JOIN security.general_user review_gu
       ON NULLIF(BTRIM(r.review_professional_id::text), '') IS NOT NULL
      AND review_gu.general_user_i_code::text = BTRIM(r.review_professional_id::text)
LEFT JOIN security.general_user_profile review_gup
       ON review_gu.general_user_i_code IS NOT NULL
      AND review_gup.general_user_profile_id = review_gu.general_user_general_user_profile
`

const referralDetailExtraSelect = `
	COALESCE(r.reason, '') AS reason,
	COALESCE(r.rejection_reason, '') AS rejection_reason,
	COALESCE(r.submitted_by_team, '') AS submitted_by_team,
	kids.children_number AS children_number,
	COALESCE(kids.children_ages, '') AS children_ages,
	COALESCE(interp.interpreter_code, '') AS interpreter_code
`

const referralDetailExtraFrom = `
LEFT JOIN LATERAL (
       SELECT f1c.victim_case_form1_victim_children_number AS children_number,
              COALESCE(f1c.victim_case_form1_victim_children_age, '') AS children_ages
       FROM salvia.victim_case_form1 f1c
       WHERE vc.victim_case_id IS NOT NULL
         AND f1c.victim_case_form1_victim_case = vc.victim_case_id
       ORDER BY f1c.victim_case_form1_id DESC
       LIMIT 1
) kids ON true
LEFT JOIN LATERAL (
       SELECT COALESCE(ri.victim_case_form2_enums_code, '') AS interpreter_code
       FROM salvia.victim_case_form2 vf2
       LEFT JOIN salvia.victim_case_form2_enums ri
              ON ri.victim_case_form2_enums_id = vf2.victim_case_form2_require_language_interpreter
       WHERE vc.victim_case_id IS NOT NULL
         AND vf2.victim_case_form2_victim_case = vc.victim_case_id
       ORDER BY vf2.victim_case_form2_id DESC
       LIMIT 1
) interp ON true
`

func (r *emergencyMeasureReferralRepository) CountByStatus(ctx context.Context, vista, agentICode string) (map[string]int, error) {
	clause, args := scopeClause(vista, agentICode)
	var rows []statusCountRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT r.status AS status, COUNT(*)::int AS n
		FROM salvia.emergency_measure_referral r
		WHERE r.deleted_at IS NULL
		`+clause+`
		GROUP BY r.status
	`, args...).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := map[string]int{}
	for _, row := range rows {
		out[row.Status] = row.N
	}
	return out, nil
}

func (r *emergencyMeasureReferralRepository) List(ctx context.Context, f EmergencyReferralListFilter) ([]EmergencyReferralRow, error) {
	clause, args := scopeClause(f.Vista, f.AgentICode)
	clause, args = appendFilters(clause, args, f)
	var scanned []emergencyReferralScan
	err := r.db.WithContext(ctx).Raw(`
		SELECT `+referralListSelect+`
		`+referralListFrom+`
		WHERE r.deleted_at IS NULL
		`+clause+`
		ORDER BY (r.urgency = ?) DESC, r.referred_at DESC NULLS LAST, r.created_at DESC
	`, append(args, models.EmergencyReferralUrgencyRiesgoInminente)...).Scan(&scanned).Error
	if err != nil {
		return nil, err
	}
	out := make([]EmergencyReferralRow, 0, len(scanned))
	for _, row := range scanned {
		measures := []string{}
		if row.MeasuresJSON != "" && row.MeasuresJSON != "null" {
			_ = json.Unmarshal([]byte(row.MeasuresJSON), &measures)
		}
		if measures == nil {
			measures = []string{}
		}
		out = append(out, EmergencyReferralRow{
			ID:                     row.ID,
			CaseICode:              strings.TrimSpace(row.CaseICode),
			CaseName:               strings.TrimSpace(row.CaseName),
			DocNumber:              row.DocNumber,
			RiskLabel:              row.RiskLabel,
			Municipality:           row.Municipality,
			Measures:               measures,
			Urgency:                row.Urgency,
			SubmittedByName:        strings.TrimSpace(row.SubmittedByName),
			ReferredAt:             row.ReferredAt,
			Status:                 row.Status,
			AssignedAgentID:        row.AssignedAgentID,
			AssignedAgentName:      strings.TrimSpace(row.AssignedAgentName),
			ReviewProfessionalName: strings.TrimSpace(row.ReviewProfessionalName),
			ReviewConcept:          row.ReviewConcept,
			ReviewResult:           row.ReviewResult,
			ReviewDate:             row.ReviewDate,
		})
	}
	return out, nil
}

// FindDetail aplica el mismo recorte que el listado. Una remisión fuera de la
// bandeja del rol devuelve nil, igual que si no existiera.
func (r *emergencyMeasureReferralRepository) FindDetail(ctx context.Context, vista, agentICode, id string) (*EmergencyReferralDetailRow, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, nil
	}
	clause, args := scopeClause(vista, agentICode)
	args = append(args, id)
	var scanned emergencyReferralDetailScan
	tx := r.db.WithContext(ctx).Raw(`
		SELECT `+referralListSelect+`,
		`+referralDetailExtraSelect+`
		`+referralListFrom+`
		`+referralDetailExtraFrom+`
		WHERE r.deleted_at IS NULL
		`+clause+`
		  AND r.id::text = ?
		LIMIT 1
	`, args...).Scan(&scanned)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if strings.TrimSpace(scanned.ID) == "" {
		return nil, nil
	}
	measures := []string{}
	if scanned.MeasuresJSON != "" && scanned.MeasuresJSON != "null" {
		_ = json.Unmarshal([]byte(scanned.MeasuresJSON), &measures)
	}
	if measures == nil {
		measures = []string{}
	}
	row := &EmergencyReferralDetailRow{
		EmergencyReferralRow: EmergencyReferralRow{
			ID:                     scanned.ID,
			CaseICode:              strings.TrimSpace(scanned.CaseICode),
			CaseName:               strings.TrimSpace(scanned.CaseName),
			DocNumber:              scanned.DocNumber,
			RiskLabel:              scanned.RiskLabel,
			Municipality:           scanned.Municipality,
			Measures:               measures,
			Urgency:                scanned.Urgency,
			SubmittedByName:        strings.TrimSpace(scanned.SubmittedByName),
			ReferredAt:             scanned.ReferredAt,
			Status:                 scanned.Status,
			AssignedAgentID:        scanned.AssignedAgentID,
			AssignedAgentName:      strings.TrimSpace(scanned.AssignedAgentName),
			ReviewProfessionalName: strings.TrimSpace(scanned.ReviewProfessionalName),
			ReviewConcept:          scanned.ReviewConcept,
			ReviewResult:           scanned.ReviewResult,
			ReviewDate:             scanned.ReviewDate,
		},
		Reason:              strings.TrimSpace(scanned.Reason),
		RejectionReason:     strings.TrimSpace(scanned.RejectionReason),
		SubmittedByTeam:     strings.TrimSpace(scanned.SubmittedByTeam),
		ChildrenAges:        strings.TrimSpace(scanned.ChildrenAges),
		RequiresInterpreter: scanned.InterpreterCode == "y",
	}
	if scanned.ChildrenNumber.Valid {
		n := int(scanned.ChildrenNumber.Int64)
		row.ChildrenNumber = &n
	}
	return row, nil
}

func (r *emergencyMeasureReferralRepository) UserHasActiveRole(ctx context.Context, iCode, roleCode string) (bool, error) {
	var n int
	err := r.db.WithContext(ctx).Raw(`
		SELECT COUNT(*)::int
		FROM security.general_user gu
		JOIN security.rel_role_general_user rr ON rr.general_user_id = gu.general_user_id
		JOIN security.role role ON role.role_id = rr.role_id
		WHERE gu.general_user_i_code = ?
		  AND role.role_code = ?
		  AND gu.general_user_status = 'e'
	`, iCode, roleCode).Scan(&n).Error
	return n > 0, err
}

func (r *emergencyMeasureReferralRepository) Reassign(ctx context.Context, ids []string, agentID string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Exec(`
			UPDATE salvia.emergency_measure_referral
			   SET assigned_agent_id = ?,
			       status = ?,
			       updated_at = NOW()
			 WHERE id::text IN ?
			   AND deleted_at IS NULL
		`, agentID, models.EmergencyReferralStatusAsignada, ids)
		if res.Error != nil {
			return res.Error
		}
		if int(res.RowsAffected) != len(ids) {
			return errReferralMissing
		}
		return nil
	})
}

func (r *emergencyMeasureReferralRepository) Review(ctx context.Context, id, professionalID, concept, result, newStatus, reviewDate string) (bool, bool, error) {
	found := false
	pending := false
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var current struct {
			Status string `gorm:"column:status"`
		}
		if err := tx.Raw(`
			SELECT status FROM salvia.emergency_measure_referral
			 WHERE id::text = ? AND deleted_at IS NULL
		`, id).Scan(&current).Error; err != nil {
			return err
		}
		if current.Status == "" {
			return nil
		}
		found = true
		if current.Status != models.EmergencyReferralStatusPendienteDeRevision {
			return nil
		}
		pending = true
		return tx.Exec(`
			UPDATE salvia.emergency_measure_referral
			   SET review_professional_id = ?,
			       review_concept = ?,
			       review_result = ?,
			       review_date = ?::date,
			       status = ?,
			       updated_at = NOW()
			 WHERE id::text = ?
			   AND status = ?
			   AND deleted_at IS NULL
		`, professionalID, concept, result, reviewDate, newStatus, id, models.EmergencyReferralStatusPendienteDeRevision).Error
	})
	return found, pending, err
}

// errReferralMissing se traduce a 404 cuando algún id del lote no existe.
var errReferralMissing = errSentinel("remision no encontrada")

type errSentinel string

func (e errSentinel) Error() string { return string(e) }

// ErrReferralMissing es el error de un reasignar con un id inexistente.
var ErrReferralMissing = errReferralMissing

func scopeClause(vista, agentICode string) (string, []interface{}) {
	switch vista {
	case "equipo-tecnico":
		return "AND r.status = ?", []interface{}{models.EmergencyReferralStatusPendienteDeRevision}
	case "supervision-convenio":
		return "AND r.status IN ?", []interface{}{models.EmergencyReferralStatusesConvenio()}
	case "cruz-roja":
		return "AND r.status IN ?", []interface{}{models.EmergencyReferralStatusesCruzRoja()}
	case "agente":
		return "AND r.assigned_agent_id = ? AND r.status IN ?", []interface{}{agentICode, models.EmergencyReferralStatusesAgente()}
	default:
		return "", nil
	}
}

func appendFilters(clause string, args []interface{}, f EmergencyReferralListFilter) (string, []interface{}) {
	if q := strings.TrimSpace(f.Q); q != "" {
		like := "%" + q + "%"
		clause += ` AND (
			TRIM(CONCAT(COALESCE(vc.victim_case_victim_names, ''), ' ', COALESCE(vc.victim_case_victim_last_names, ''))) ILIKE ?
			OR COALESCE(vc.victim_case_victim_doc_number, '') ILIKE ?
		)`
		args = append(args, like, like)
	}
	if f.Status != "" {
		clause += " AND r.status = ?"
		args = append(args, f.Status)
	}
	if f.Urgency != "" {
		clause += " AND r.urgency = ?"
		args = append(args, f.Urgency)
	}
	if f.Vista == "supervisor" && f.AgentID != "" {
		if f.AgentID == "sin-asignar" {
			clause += " AND (r.assigned_agent_id IS NULL OR BTRIM(r.assigned_agent_id) = '')"
		} else {
			clause += " AND r.assigned_agent_id = ?"
			args = append(args, f.AgentID)
		}
	}
	if f.ID != "" {
		clause += " AND r.id::text = ?"
		args = append(args, f.ID)
	}
	return clause, args
}
