package repository

import (
	internaldb "bitsflow/internal/db"
	"bitsflow/internal/models"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

const psychosocialSessionCountSubquery = `(
	SELECT COUNT(*)::int
	FROM salvia.team_contact tc
	WHERE tc.psicosocial_id = ps.id::text
	  AND tc.is_psico_session = true
	  AND tc.is_completed = true
	  AND tc.deleted_at IS NULL
)`

const psychosocialVictimPhoneSelect = `
	COALESCE(
		NULLIF(vf2.victim_case_form2_victim_phone::text, ''),
		(
			SELECT vcf1.victim_contact_form1_phone::text
			FROM salvia.victim_contact_form1 vcf1
			WHERE vcf1.victim_contact_form1_victim_contact = vc.victim_case_victim_contact
			ORDER BY vcf1.victim_contact_form1_id DESC
			LIMIT 1
		),
		''
	)
`

const psychosocialListBaseFrom = `
FROM salvia.psychosocial_support ps
INNER JOIN salvia.victim_case vc
        ON NULLIF(BTRIM(ps.case_id::text), '') IS NOT NULL
       AND vc.victim_case_i_code = BTRIM(ps.case_id::text)
LEFT JOIN salvia.victim_case_form2 vf2
       ON vf2.victim_case_form2_victim_case = vc.victim_case_id
LEFT JOIN security.town t
       ON t.town_code = vc.victim_case_victim_town_code
LEFT JOIN security.general_user submitter_gu
       ON NULLIF(BTRIM(ps.submitted_by::text), '') IS NOT NULL
      AND submitter_gu.general_user_i_code::text = BTRIM(ps.submitted_by::text)
LEFT JOIN security.general_user_profile submitter_gup
       ON submitter_gu.general_user_i_code IS NOT NULL
      AND submitter_gup.general_user_profile_id = submitter_gu.general_user_general_user_profile
LEFT JOIN salvia.dupla d
       ON NULLIF(BTRIM(ps.dupla_id::text), '') IS NOT NULL
      AND d.id = BTRIM(ps.dupla_id::text)
      AND d.deleted_at IS NULL
LEFT JOIN security.general_user psych_gu
       ON d.id IS NOT NULL
      AND NULLIF(BTRIM(d.psychologist_id::text), '') IS NOT NULL
      AND psych_gu.general_user_i_code::text = BTRIM(d.psychologist_id::text)
LEFT JOIN security.general_user_profile psych_gup
       ON psych_gu.general_user_i_code IS NOT NULL
      AND psych_gup.general_user_profile_id = psych_gu.general_user_general_user_profile
LEFT JOIN security.general_user sw_gu
       ON d.id IS NOT NULL
      AND NULLIF(BTRIM(d.social_worker_id::text), '') IS NOT NULL
      AND sw_gu.general_user_i_code::text = BTRIM(d.social_worker_id::text)
LEFT JOIN security.general_user_profile sw_gup
       ON sw_gu.general_user_i_code IS NOT NULL
      AND sw_gup.general_user_profile_id = sw_gu.general_user_general_user_profile
LEFT JOIN security.general_user prof_gu
       ON NULLIF(BTRIM(ps.professional_id::text), '') IS NOT NULL
      AND prof_gu.general_user_i_code::text = BTRIM(ps.professional_id::text)
LEFT JOIN security.general_user_profile prof_gup
       ON prof_gu.general_user_i_code IS NOT NULL
      AND prof_gup.general_user_profile_id = prof_gu.general_user_general_user_profile
`

const psychosocialProfSpecialtySelect = `(
	SELECT CASE
		WHEN COUNT(*) FILTER (
			WHERE LOWER(r.role_name) LIKE '%psicolog%'
			   OR LOWER(r.role_code) LIKE '%psi%'
			   OR LOWER(r.role_code) = 'psicologia'
		) > 0 THEN 'psicologia'
		WHEN COUNT(*) FILTER (
			WHERE LOWER(r.role_name) LIKE '%trab%social%'
			   OR LOWER(r.role_code) LIKE '%ts%'
			   OR LOWER(r.role_code) IN ('trab. social', 'trab_social', 'tss')
		) > 0 THEN 'trab. social'
		ELSE ''
	END
	FROM security.rel_role_general_user rr
	JOIN security.role r ON r.role_id = rr.role_id
	WHERE rr.general_user_id = prof_gu.general_user_id
)`

const psychosocialListSelectCols = `
	ps.id,
	ps.case_id,
	ps.follow_up_id,
	ps.status,
	ps.created_at,
	ps.submitted_by_team,
	ps.dupla_id,
	ps.professional_id,
	COALESCE(submitter_gup.general_user_profile_names, '') AS submitted_by_names,
	COALESCE(submitter_gup.general_user_profile_last_names, '') AS submitted_by_last_names,
	COALESCE(d.name, '') AS dupla_name,
	COALESCE(psych_gup.general_user_profile_names || ' ' || psych_gup.general_user_profile_last_names, '') AS psychologist_name,
	COALESCE(sw_gup.general_user_profile_names || ' ' || sw_gup.general_user_profile_last_names, '') AS social_worker_name,
	COALESCE(prof_gup.general_user_profile_names || ' ' || prof_gup.general_user_profile_last_names, '') AS professional_name,
	COALESCE(prof_gu.general_user_team, '') AS professional_team,
	` + psychosocialProfSpecialtySelect + ` AS professional_specialty,
	COALESCE(vc.victim_case_victim_names, '') AS victim_names,
	COALESCE(vc.victim_case_victim_last_names, '') AS victim_last_names,
	COALESCE(vc.victim_case_victim_doc_number, '') AS doc_number,
	(` + psychosocialVictimPhoneSelect + `) AS victim_phone,
	COALESCE(t.town_name, '') AS municipality,
	COALESCE(vf2.victim_case_form2_risk_level, 0) AS risk_level,
	` + psychosocialSessionCountSubquery + ` AS session_count
`

// PsychosocialListFilters parámetros de consulta E-01 … E-13.
type PsychosocialListFilters struct {
	FilterProfessionalID     string
	FilterDuplaID            string
	FilterEstadoRemision     string
	FilterSesionesCompletadas string
	FilterEquipoRemitente    string
	FilterNivelRiesgo        string
	FilterNumeroIdentidad    string
	FilterTelefono           string
	Sort                     string
	Order                    string
	Page                     int
	PageSize                 int
}

// PsychosocialListResult respuesta paginada del listado.
type PsychosocialListResult struct {
	Remisiones []models.PsychosocialListItem
	Total      int64
	Page       int
	PageSize   int
}

// PsychosocialListRepository acceso a datos del componente remisiones-psicosocial.
type PsychosocialListRepository interface {
	List(ctx context.Context, filters PsychosocialListFilters) (PsychosocialListResult, error)
	Stats(ctx context.Context, filters PsychosocialListFilters) (models.PsychosocialListStats, error)
	ListEquiposRemitentes(ctx context.Context) ([]string, error)
}

type psychosocialListRepository struct {
	db *gorm.DB
}

func NewPsychosocialListRepository(db *gorm.DB) PsychosocialListRepository {
	return &psychosocialListRepository{db: db}
}

type psychosocialListRow struct {
	ID                 string  `gorm:"column:id"`
	CaseID             string  `gorm:"column:case_id"`
	FollowUpID         string  `gorm:"column:follow_up_id"`
	Status             string  `gorm:"column:status"`
	CreatedAt            time.Time `gorm:"column:created_at"`
	SubmittedByTeam    *string `gorm:"column:submitted_by_team"`
	DuplaID            *string `gorm:"column:dupla_id"`
	ProfessionalID     *string `gorm:"column:professional_id"`
	SubmittedByNames   string  `gorm:"column:submitted_by_names"`
	SubmittedByLastNames string `gorm:"column:submitted_by_last_names"`
	DuplaName          string  `gorm:"column:dupla_name"`
	PsychologistName   string  `gorm:"column:psychologist_name"`
	SocialWorkerName   string  `gorm:"column:social_worker_name"`
	ProfessionalName     string  `gorm:"column:professional_name"`
	ProfessionalTeam     string  `gorm:"column:professional_team"`
	ProfessionalSpecialty  string  `gorm:"column:professional_specialty"`
	VictimNames        string  `gorm:"column:victim_names"`
	VictimLastNames    string  `gorm:"column:victim_last_names"`
	DocNumber          string  `gorm:"column:doc_number"`
	VictimPhone        string  `gorm:"column:victim_phone"`
	Municipality       string  `gorm:"column:municipality"`
	RiskLevel          int     `gorm:"column:risk_level"`
	SessionCount       int     `gorm:"column:session_count"`
}

func (r *psychosocialListRepository) List(ctx context.Context, filters PsychosocialListFilters) (PsychosocialListResult, error) {
	page := filters.Page
	if page < 1 {
		page = 1
	}
	pageSize := filters.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	whereSQL, args := buildPsychosocialListWhere(filters)
	orderSQL := buildPsychosocialListOrder(filters)

	countSQL := "SELECT COUNT(*) " + psychosocialListBaseFrom + " " + whereSQL
	var total int64
	if err := internaldb.WithRetry(func() error {
		return r.db.WithContext(ctx).Raw(countSQL, args...).Scan(&total).Error
	}); err != nil {
		return PsychosocialListResult{}, err
	}

	offset := (page - 1) * pageSize
	dataSQL := fmt.Sprintf(
		"SELECT %s %s %s %s LIMIT %d OFFSET %d",
		psychosocialListSelectCols, psychosocialListBaseFrom, whereSQL, orderSQL, pageSize, offset,
	)

	var rows []psychosocialListRow
	if err := internaldb.WithRetry(func() error {
		return r.db.WithContext(ctx).Raw(dataSQL, args...).Scan(&rows).Error
	}); err != nil {
		return PsychosocialListResult{}, err
	}

	items := make([]models.PsychosocialListItem, len(rows))
	for i, row := range rows {
		submittedByTeam := ""
		if row.SubmittedByTeam != nil {
			submittedByTeam = *row.SubmittedByTeam
		}
		submittedName := strings.TrimSpace(row.SubmittedByNames + " " + row.SubmittedByLastNames)
		items[i] = models.PsychosocialListItem{
			ID:               row.ID,
			CaseICode:        row.CaseID,
			FollowUpID:       row.FollowUpID,
			Status:           row.Status,
			SessionCount:     row.SessionCount,
			CreatedAt:        row.CreatedAt,
			SubmittedByName:  submittedName,
			SubmittedByTeam:  submittedByTeam,
			DuplaID:          row.DuplaID,
			DuplaName:        row.DuplaName,
			PsychologistName: strings.TrimSpace(row.PsychologistName),
			SocialWorkerName: strings.TrimSpace(row.SocialWorkerName),
			ProfessionalID:   row.ProfessionalID,
			ProfessionalName: strings.TrimSpace(row.ProfessionalName),
			ProfessionalTeam: row.ProfessionalTeam,
			ProfessionalRole: psychosocialRoleLabelFromSpecialty(row.ProfessionalSpecialty),
			VictimNames:      row.VictimNames,
			VictimLastNames:  row.VictimLastNames,
			DocNumber:        row.DocNumber,
			VictimPhone:      row.VictimPhone,
			Municipality:     row.Municipality,
			RiskLevel:        row.RiskLevel,
		}
	}

	return PsychosocialListResult{
		Remisiones: items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (r *psychosocialListRepository) Stats(ctx context.Context, filters PsychosocialListFilters) (models.PsychosocialListStats, error) {
	whereSQL, args := buildPsychosocialListWhere(filters)
	statsSQL := `
SELECT
	COUNT(*)::bigint AS total,
	COUNT(*) FILTER (WHERE ps.status = 'abierto')::bigint AS abierto,
	COUNT(*) FILTER (WHERE ps.status = 'en_gestion')::bigint AS en_gestion,
	COUNT(*) FILTER (WHERE ps.status = 'en_devolucion')::bigint AS en_devolucion,
	COUNT(*) FILTER (WHERE ps.status = 'cerrado')::bigint AS cerrado
` + psychosocialListBaseFrom + " " + whereSQL

	var stats models.PsychosocialListStats
	if err := internaldb.WithRetry(func() error {
		row := r.db.WithContext(ctx).Raw(statsSQL, args...).Row()
		return row.Scan(&stats.Total, &stats.Abierto, &stats.EnGestion, &stats.EnDevolucion, &stats.Cerrado)
	}); err != nil {
		return models.PsychosocialListStats{}, err
	}
	return stats, nil
}

func (r *psychosocialListRepository) ListEquiposRemitentes(ctx context.Context) ([]string, error) {
	var teams []string
	err := internaldb.WithRetry(func() error {
		return r.db.WithContext(ctx).Raw(`
SELECT DISTINCT ps.submitted_by_team AS team
FROM salvia.psychosocial_support ps
WHERE ps.deleted_at IS NULL
  AND ps.submitted_by_team IS NOT NULL
  AND TRIM(ps.submitted_by_team) <> ''
ORDER BY team
`).Scan(&teams).Error
	})
	if err != nil {
		return nil, err
	}
	if teams == nil {
		teams = []string{}
	}
	return teams, nil
}

const psychosocialProfessionalIDFilterClause = `AND (
	BTRIM(ps.professional_id::text) = BTRIM(?)
	OR ps.dupla_id IN (
		SELECT d.id
		FROM salvia.dupla d
		WHERE d.deleted_at IS NULL
		  AND (
			BTRIM(d.psychologist_id::text) = BTRIM(?)
			OR BTRIM(d.social_worker_id::text) = BTRIM(?)
		  )
	)
)`

func buildPsychosocialListWhere(filters PsychosocialListFilters) (string, []interface{}) {
	clauses := []string{"WHERE ps.deleted_at IS NULL"}
	args := []interface{}{}

	if filters.FilterProfessionalID != "" {
		clauses = append(clauses, psychosocialProfessionalIDFilterClause)
		pid := strings.TrimSpace(filters.FilterProfessionalID)
		args = append(args, pid, pid, pid)
	}
	if filters.FilterDuplaID != "" {
		clauses = append(clauses, "AND BTRIM(ps.dupla_id::text) = BTRIM(?)")
		args = append(args, filters.FilterDuplaID)
	}
	if status := psychosocialValidStatus(filters.FilterEstadoRemision); status != "" {
		clauses = append(clauses, "AND ps.status = ?")
		args = append(args, status)
	}
	if filters.FilterEquipoRemitente != "" {
		clauses = append(clauses, "AND ps.submitted_by_team = ?")
		args = append(args, filters.FilterEquipoRemitente)
	}
	if filters.FilterNumeroIdentidad != "" {
		clauses = append(clauses, "AND vc.victim_case_victim_doc_number ILIKE ?")
		args = append(args, "%"+filters.FilterNumeroIdentidad+"%")
	}
	if filters.FilterTelefono != "" {
		clauses = append(clauses, "AND ("+psychosocialVictimPhoneSelect+") ILIKE ?")
		args = append(args, "%"+filters.FilterTelefono+"%")
	}
	if filters.FilterNivelRiesgo != "" {
		if level := psychosocialRiskLevelFromSlug(filters.FilterNivelRiesgo); level > 0 {
			clauses = append(clauses, "AND vf2.victim_case_form2_risk_level = ?")
			args = append(args, level)
		}
	}
	if filters.FilterSesionesCompletadas != "" {
		if n, err := strconv.Atoi(filters.FilterSesionesCompletadas); err == nil {
			clauses = append(clauses, "AND "+psychosocialSessionCountSubquery+" = ?")
			args = append(args, n)
		}
	}

	return strings.Join(clauses, " "), args
}

func buildPsychosocialListOrder(filters PsychosocialListFilters) string {
	sortCol := "ps.created_at"
	if filters.Sort == "created_at" {
		sortCol = "ps.created_at"
	}
	order := "DESC"
	if strings.EqualFold(filters.Order, "asc") {
		order = "ASC"
	}
	return fmt.Sprintf("ORDER BY %s %s", sortCol, order)
}

func psychosocialRoleLabelFromSpecialty(specialty string) string {
	switch strings.ToLower(strings.TrimSpace(specialty)) {
	case "psicologia":
		return "Psicóloga"
	case "trab. social":
		return "Trab. Social"
	default:
		return "Profesional"
	}
}

func psychosocialValidStatus(raw string) string {
	switch strings.TrimSpace(raw) {
	case models.PsychosocialSupportStatusAbierto,
		models.PsychosocialSupportStatusEnGestion,
		models.PsychosocialSupportStatusEnDevolucion,
		models.PsychosocialSupportStatusCerrado:
		return strings.TrimSpace(raw)
	default:
		return ""
	}
}

func psychosocialRiskLevelFromSlug(slug string) int {
	switch strings.ToLower(strings.TrimSpace(slug)) {
	case "bajo":
		return 1
	case "moderado":
		return 2
	case "alto":
		return 3
	case "extremo":
		return 4
	default:
		return 0
	}
}