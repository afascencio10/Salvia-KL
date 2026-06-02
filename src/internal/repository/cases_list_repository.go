package repository

import (
	"bitsflow/internal/models"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

// casesListTimezone zona horaria para filtros de fecha calendario (E-09: casos nuevos = hoy).
const casesListTimezone = "America/Bogota"

// CasesListFilters parámetros de consulta para el listado de casos del componente casos-component.
// Los filtros chip y dropdown son aditivos (combinables entre sí).
type CasesListFilters struct {
	FilterKey           string
	FilterValue         string
	ChipFilter          string // Filtro chip aditivo (p. ej. casos_nuevos) combinable con FilterKey
	DropdownFilterKey   string // Legacy: un solo dropdown (filter_* preferido)
	DropdownFilterValue string
	FilterRiesgo                 string // filter_riesgo (E-10)
	FilterEquipo                 string // filter_equipo (E-11)
	FilterSeguimientosEjecutados string // filter_seguimientos_ejecutados (E-12)
	Search      string
	Sort        string
	Order       string
	Page        int
	PageSize    int
}

// CasesListResult respuesta paginada del listado de casos.
type CasesListResult struct {
	Cases    []models.CaseListItem
	Total    int64
	Page     int
	PageSize int
}

// CasesListRepository acceso a datos para GET /api/v1/cases/list.
type CasesListRepository interface {
	List(ctx context.Context, filters CasesListFilters) (CasesListResult, error)
}

type casesListRepository struct {
	db *gorm.DB
}

func NewCasesListRepository(db *gorm.DB) CasesListRepository {
	return &casesListRepository{db: db}
}

type caseListRow struct {
	ID               int64      `gorm:"column:victim_case_id"`
	ICode            string     `gorm:"column:victim_case_i_code"`
	Names            string     `gorm:"column:victim_case_victim_names"`
	LastNames        string     `gorm:"column:victim_case_victim_last_names"`
	DocNumber        string     `gorm:"column:victim_case_victim_doc_number"`
	VictimPhone      string     `gorm:"column:victim_phone"`
	CreationDate     time.Time  `gorm:"column:victim_case_creation_date"`
	Status           string     `gorm:"column:victim_case_status"`
	OwnerNames       string     `gorm:"column:owner_names"`
	OwnerLastNames   string     `gorm:"column:owner_last_names"`
	OwnerTeam        string     `gorm:"column:owner_team"`
	CaseTeam         string     `gorm:"column:case_team"`
	RiskStatus       *string    `gorm:"column:risk_status"`
	NextFollowUpDate          *time.Time `gorm:"column:next_follow_up_date"`
	CompletedFollowUpsCount int        `gorm:"column:completed_follow_ups_count"`
}

// riskStatusSelect mapea victim_case_form2_risk_level (1-4) al slug usado por el frontend.
// Sin form2 o nivel ausente → 'desconocido' (E-10).
const riskStatusSelect = `
    CASE vf2.victim_case_form2_risk_level
        WHEN 1 THEN 'bajo'
        WHEN 2 THEN 'moderado'
        WHEN 3 THEN 'alto'
        WHEN 4 THEN 'extremo'
        ELSE 'desconocido'
    END
`

// victimPhoneSelect teléfono de la víctima: form2 del caso (copiado al crear) o form1 del contacto vinculado.
const victimPhoneSelect = `
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

const casesListBaseFrom = `
FROM salvia.victim_case vc
LEFT JOIN security.general_user gu
       ON gu.general_user_i_code = vc.agent_id
LEFT JOIN security.general_user_profile gup
       ON gup.general_user_profile_id = gu.general_user_general_user_profile
LEFT JOIN salvia.victim_case_form2 vf2
       ON vf2.victim_case_form2_victim_case = vc.victim_case_id
`

const casesListSelectCols = `
    vc.victim_case_id,
    vc.victim_case_i_code,
    vc.victim_case_victim_names,
    vc.victim_case_victim_last_names,
    vc.victim_case_victim_doc_number,
    (` + victimPhoneSelect + `) AS victim_phone,
    vc.victim_case_creation_date,
    vc.victim_case_status,
    COALESCE(gup.general_user_profile_names, '')      AS owner_names,
    COALESCE(gup.general_user_profile_last_names, '') AS owner_last_names,
    COALESCE(gu.general_user_team, '')              AS owner_team,
    COALESCE(vc.victim_case_team, '')               AS case_team,
    (` + riskStatusSelect + `) AS risk_status,
    (
        SELECT scheduled_date
        FROM salvia.follow_up_v2
        WHERE case_id = vc.victim_case_i_code
          AND status = 'PENDIENTE'
          AND deleted_at IS NULL
        ORDER BY scheduled_date ASC
        LIMIT 1
    ) AS next_follow_up_date,
    (` + completedFollowUpsCountSubquery + `) AS completed_follow_ups_count
`

// completedFollowUpsCountSubquery conteo de follow_up_v2 REALIZADO por caso (E-01, E-12).
const completedFollowUpsCountSubquery = `(
    SELECT COUNT(*)::int
    FROM salvia.follow_up_v2 fu
    WHERE fu.case_id = vc.victim_case_i_code
      AND fu.status = 'REALIZADO'
      AND fu.deleted_at IS NULL
)`

func (r *casesListRepository) List(ctx context.Context, filters CasesListFilters) (CasesListResult, error) {
	page := filters.Page
	if page < 1 {
		page = 1
	}
	pageSize := filters.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	whereSQL, args := buildCasesListWhere(filters)
	orderSQL := buildCasesListOrder(filters)

	countSQL := "SELECT COUNT(*) " + casesListBaseFrom + " " + whereSQL
	var total int64
	if err := r.db.WithContext(ctx).Raw(countSQL, args...).Scan(&total).Error; err != nil {
		return CasesListResult{}, err
	}

	offset := (page - 1) * pageSize
	dataSQL := fmt.Sprintf(
		"SELECT %s %s %s %s LIMIT %d OFFSET %d",
		casesListSelectCols, casesListBaseFrom, whereSQL, orderSQL, pageSize, offset,
	)

	var rows []caseListRow
	if err := r.db.WithContext(ctx).Raw(dataSQL, args...).Scan(&rows).Error; err != nil {
		return CasesListResult{}, err
	}

	items := make([]models.CaseListItem, len(rows))
	for i, row := range rows {
		items[i] = models.CaseListItem{
			ID:               row.ID,
			ICode:            row.ICode,
			Names:            row.Names,
			LastNames:        row.LastNames,
			DocNumber:        row.DocNumber,
			VictimPhone:      row.VictimPhone,
			CreationDate:     row.CreationDate,
			Status:           row.Status,
			OwnerNames:       row.OwnerNames,
			OwnerLastNames:   row.OwnerLastNames,
			OwnerTeam:        row.OwnerTeam,
			CaseTeam:         row.CaseTeam,
			RiskStatus:       row.RiskStatus,
			NextFollowUpDate:        row.NextFollowUpDate,
			CompletedFollowUpsCount: row.CompletedFollowUpsCount,
		}
	}

	return CasesListResult{
		Cases:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// E-09: casos creados hoy y hasta 5 días calendario antes (zona Bogotá).
const casosNuevosRangeClause = `(vc.victim_case_creation_date AT TIME ZONE '` + casesListTimezone + `')::date BETWEEN (NOW() AT TIME ZONE '` + casesListTimezone + `')::date - INTERVAL '5 days' AND (NOW() AT TIME ZONE '` + casesListTimezone + `')::date`

func buildCasesListWhere(filters CasesListFilters) (string, []interface{}) {
	clauses := []string{"1=1"}
	args := []interface{}{}

	casosNuevosActive := filters.ChipFilter == "casos_nuevos" || filters.FilterKey == "casos_nuevos"
	if casosNuevosActive {
		clauses = append(clauses, casosNuevosRangeClause)
	}

	if riskValue := casesListRiskFilterValue(filters); riskValue != "" {
		if level, ok := riskLevelFromFilterValue(riskValue); ok {
			clauses = append(clauses, "vf2.victim_case_form2_risk_level = ?")
			args = append(args, level)
		}
	}

	// E-11: equipo del caso; no combinable con scope agentId (Mis casos)
	if teamValue := casesListTeamFilterValue(filters); teamValue != "" && !casesListHasAgentScope(filters) {
		clauses = append(clauses, "vc.victim_case_team = ?")
		args = append(args, teamValue)
	}

	// E-12: igualdad exacta con el conteo de seguimientos REALIZADO
	if segValue := casesListSeguimientosFilterValue(filters); segValue != "" {
		if n, ok := seguimientosCountFromFilterValue(segValue); ok {
			clauses = append(clauses, completedFollowUpsCountSubquery+" = ?")
			args = append(args, n)
		}
	}

	switch filters.FilterKey {
	case "casos_nuevos":
		// Ya aplicado arriba (chip o filter_key legacy)
	case "riesgo":
		// Ya aplicado arriba
	case "equipo":
		// Ya aplicado arriba
	case "seguimientos_ejecutados":
		// Ya aplicado arriba
	case "persona_asignada":
		if filters.FilterValue != "" {
			clauses = append(clauses, "vc.agent_id = ?")
			args = append(args, filters.FilterValue)
		}
	}

	search := strings.TrimSpace(filters.Search)
	if search != "" {
		pattern := "%" + search + "%"
		clauses = append(clauses, `(
            vc.victim_case_victim_doc_number ILIKE ?
            OR vf2.victim_case_form2_victim_phone::text ILIKE ?
            OR EXISTS (
                SELECT 1
                FROM salvia.victim_contact_form1 vcf1
                WHERE vcf1.victim_contact_form1_victim_contact = vc.victim_case_victim_contact
                  AND vcf1.victim_contact_form1_phone::text ILIKE ?
            )
        )`)
		args = append(args, pattern, pattern, pattern)
	}

	return "WHERE " + strings.Join(clauses, " AND "), args
}

func buildCasesListOrder(filters CasesListFilters) string {
	orderDir := "DESC"
	if strings.EqualFold(filters.Order, "asc") {
		orderDir = "ASC"
	}

	switch filters.Sort {
	case "next_follow_up":
		return fmt.Sprintf("ORDER BY next_follow_up_date %s NULLS LAST", orderDir)
	default:
		return fmt.Sprintf("ORDER BY vc.victim_case_creation_date %s", orderDir)
	}
}

// casesListTeamFilterValue obtiene el valor del filtro equipo (combinable salvo Mis casos).
func casesListTeamFilterValue(filters CasesListFilters) string {
	if filters.FilterEquipo != "" {
		return filters.FilterEquipo
	}
	if filters.DropdownFilterKey == "equipo" && filters.DropdownFilterValue != "" {
		return filters.DropdownFilterValue
	}
	if filters.FilterKey == "equipo" && filters.FilterValue != "" {
		return filters.FilterValue
	}
	return ""
}

// casesListHasAgentScope indica si la consulta está acotada a un agente (Mis casos).
func casesListHasAgentScope(filters CasesListFilters) bool {
	return filters.FilterKey == "persona_asignada" && filters.FilterValue != ""
}

// casesListRiskFilterValue obtiene el valor del filtro riesgo (combinable con otros dropdowns).
func casesListRiskFilterValue(filters CasesListFilters) string {
	if filters.FilterRiesgo != "" {
		return filters.FilterRiesgo
	}
	if filters.DropdownFilterKey == "riesgo" && filters.DropdownFilterValue != "" {
		return filters.DropdownFilterValue
	}
	if filters.FilterKey == "riesgo" && filters.FilterValue != "" {
		return filters.FilterValue
	}
	return ""
}

// casesListSeguimientosFilterValue obtiene el valor del filtro E-12 (0–10).
func casesListSeguimientosFilterValue(filters CasesListFilters) string {
	if filters.FilterSeguimientosEjecutados != "" {
		return filters.FilterSeguimientosEjecutados
	}
	if filters.DropdownFilterKey == "seguimientos_ejecutados" && filters.DropdownFilterValue != "" {
		return filters.DropdownFilterValue
	}
	if filters.FilterKey == "seguimientos_ejecutados" && filters.FilterValue != "" {
		return filters.FilterValue
	}
	return ""
}

// seguimientosCountFromFilterValue valida el entero 0–10 del dropdown E-12.
func seguimientosCountFromFilterValue(value string) (int, bool) {
	n, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil || n < 0 || n > 10 {
		return 0, false
	}
	return n, true
}

// riskLevelFromFilterValue traduce el valor del filtro UI al entero 1-4 de victim_case_form2_risk_level.
// Escala: 1=Bajo, 2=Moderado, 3=Alto, 4=Extremo (HU-027).
func riskLevelFromFilterValue(value string) (int, bool) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "bajo", "1":
		return 1, true
	case "medio", "moderado", "2":
		return 2, true
	case "alto", "3":
		return 3, true
	case "extremo", "4":
		return 4, true
	default:
		return 0, false
	}
}
