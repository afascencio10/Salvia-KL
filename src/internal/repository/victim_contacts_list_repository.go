package repository

import (
	"bitsflow/internal/models"
	"context"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// victimContactsListTimezone — misma zona que casesListTimezone (cases_list_repository.go),
// usada para el cálculo de "día completo" del esquema 3x3 (America/Bogota).
const victimContactsListTimezone = "America/Bogota"

// victimContactLegacyInvalidated — valor de victim_contact_status que deja el botón
// "Invalidar" legado de /salvia/primer_contacto ('v' = vigente).
const victimContactLegacyInvalidated = "i"

// VictimContactsListFilters parámetros de consulta de GET /api/v1/victim-contacts/list.
// Mismo patrón que CasesListFilters — filter_key/filter_value genérico +
// filtros con nombre propio para los que necesitan combinarse (estado_reporte,
// asignación). Ver DocsMD/Screens/Listado Reportes/listado-reportes-interface.md.
type VictimContactsListFilters struct {
	FilterKey          string
	FilterValue        string
	FilterEstadoReporte string
	Search             string
	Sort               string
	Order              string
	Page               int
	PageSize           int
}

// VictimContactsListResult respuesta paginada del listado de reportes.
type VictimContactsListResult struct {
	Reports  []models.ReportListItem
	Total    int64
	Page     int
	PageSize int
}

// VictimContactsListRepository acceso a datos para GET /api/v1/victim-contacts/list.
type VictimContactsListRepository interface {
	List(ctx context.Context, filters VictimContactsListFilters) (VictimContactsListResult, error)
}

type victimContactsListRepository struct {
	db *gorm.DB
}

func NewVictimContactsListRepository(db *gorm.DB) VictimContactsListRepository {
	return &victimContactsListRepository{db: db}
}

type reportListRow struct {
	ID                 int64     `gorm:"column:victim_contact_id"`
	ICode              string    `gorm:"column:victim_contact_i_code"`
	Names              string    `gorm:"column:victim_contact_names"`
	LastNames          string    `gorm:"column:victim_contact_last_names"`
	CreationDate       time.Time `gorm:"column:victim_contact_creation_date"`
	ReportType         *string `gorm:"column:report_type"`
	TotalFallidos      int     `gorm:"column:total_fallidos"`
	DiasCompletos      int     `gorm:"column:dias_completos"`
	AgentID            *string `gorm:"column:agent_id"`
	CloseReasonCode    *string `gorm:"column:close_reason_code"`
	CloseCategory      *string `gorm:"column:close_category"`
	AssignedAgentNames *string `gorm:"column:assigned_agent_names"`
	AssignedAgentLast  *string `gorm:"column:assigned_agent_last_names"`
	Team               *string `gorm:"column:team"`
	LegacyStatus       string  `gorm:"column:legacy_status"`
}

// reportTypeSelect devuelve el NOMBRE del enum del tipo de reporte (V1:
// victim_case_form2_report_type_v|_r — V2: victim_case_form2_report_type_v2_victima|
// _tercero|_familiar|_hombre). Ambos viven en victim_contact_form2.victim_contact_form2_report_type
// (la columna ..._report_type_details es otra cosa: "Usted es una persona…").
// El servicio lo traduce a la etiqueta corta (Víctima/Tercero/Familiar/Hombre).
const reportTypeSelect = `
    (SELECT e.victim_case_form2_enums_name
       FROM salvia.victim_contact_form2 vf
       JOIN salvia.victim_case_form2_enums e
         ON e.victim_case_form2_enums_id = vf.victim_contact_form2_report_type
      WHERE vf.victim_contact_form2_victim_contact = vc.victim_contact_id
      ORDER BY vf.victim_contact_form2_id DESC LIMIT 1)
`

// attemptsByDaySubquery cuenta intentos fallidos agrupados por día calendario
// (America/Bogota), para derivar diasCompletos (días con >= 3 fallidos) y
// totalFallidos. "Día completo" exige 3 por día, no 1 — ver corrección
// 2026-10-12 en related-tables.md (comentario de Felipe en flow-M07).
const attemptsAggregateSelect = `
    (SELECT COUNT(*)::int FROM salvia.victim_contact_attempt vca
       WHERE vca.victim_contact_attempt_victim_contact = vc.victim_contact_i_code
         AND vca.victim_contact_attempt_result = 'fallido'
         AND vca.victim_contact_attempt_deleted_at IS NULL
    ) AS total_fallidos,
    (SELECT COUNT(*)::int FROM (
        SELECT (vca.victim_contact_attempt_date AT TIME ZONE 'UTC' AT TIME ZONE '` + victimContactsListTimezone + `')::date AS dia,
               COUNT(*) AS n
        FROM salvia.victim_contact_attempt vca
        WHERE vca.victim_contact_attempt_victim_contact = vc.victim_contact_i_code
          AND vca.victim_contact_attempt_result = 'fallido'
          AND vca.victim_contact_attempt_deleted_at IS NULL
        GROUP BY dia
     ) dias WHERE dias.n >= 3
    ) AS dias_completos
`

const victimContactsListBaseFrom = `
FROM salvia.victim_contact vc
LEFT JOIN security.general_user gu
       ON gu.general_user_i_code = vc.victim_contact_agent_id
LEFT JOIN security.general_user_profile gup
       ON gup.general_user_profile_id = gu.general_user_general_user_profile
LEFT JOIN salvia.victim_case_form2_enums closeEnum
       ON closeEnum.victim_case_form2_enums_id = vc.victim_contact_close_reason
`

const victimContactsListSelectCols = `
    vc.victim_contact_id,
    vc.victim_contact_i_code,
    vc.victim_contact_names,
    vc.victim_contact_last_names,
    vc.victim_contact_creation_date,
    (` + reportTypeSelect + `) AS report_type,
    ` + attemptsAggregateSelect + `,
    vc.victim_contact_agent_id AS agent_id,
    closeEnum.victim_case_form2_enums_code AS close_reason_code,
    CASE WHEN closeEnum.victim_case_form2_enums_code = 'ap' THEN 'valido'
         WHEN closeEnum.victim_case_form2_enums_code IS NOT NULL THEN 'invalido'
         ELSE NULL END AS close_category,
    gup.general_user_profile_names AS assigned_agent_names,
    gup.general_user_profile_last_names AS assigned_agent_last_names,
    vc.victim_contact_team AS team,
    vc.victim_contact_status AS legacy_status
`

func (r *victimContactsListRepository) List(ctx context.Context, filters VictimContactsListFilters) (VictimContactsListResult, error) {
	page := filters.Page
	if page < 1 {
		page = 1
	}
	pageSize := filters.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	whereSQL, args := buildVictimContactsListWhere(filters)
	orderSQL := buildVictimContactsListOrder(filters)

	countSQL := "SELECT COUNT(*) " + victimContactsListBaseFrom + " " + whereSQL
	var total int64
	if err := r.db.WithContext(ctx).Raw(countSQL, args...).Scan(&total).Error; err != nil {
		return VictimContactsListResult{}, err
	}

	offset := (page - 1) * pageSize
	dataSQL := fmt.Sprintf(
		"SELECT %s %s %s %s LIMIT %d OFFSET %d",
		victimContactsListSelectCols, victimContactsListBaseFrom, whereSQL, orderSQL, pageSize, offset,
	)

	var rows []reportListRow
	if err := r.db.WithContext(ctx).Raw(dataSQL, args...).Scan(&rows).Error; err != nil {
		return VictimContactsListResult{}, err
	}

	items := make([]models.ReportListItem, len(rows))
	for i, row := range rows {
		items[i] = models.ReportListItem{
			ID:                 row.ID,
			ICode:              row.ICode,
			Names:              row.Names,
			LastNames:          row.LastNames,
			CreationDate:       row.CreationDate,
			ReportType:         derefStr(row.ReportType),
			TotalFallidos:      row.TotalFallidos,
			DiasCompletos:      row.DiasCompletos,
			ReportStatus:       reportStatusFrom(row.AgentID, row.CloseReasonCode, row.CloseCategory, row.LegacyStatus),
			AssignedAgentICode: derefStr(row.AgentID),
			AssignedAgentName:  strings.TrimSpace(derefStr(row.AssignedAgentNames) + " " + derefStr(row.AssignedAgentLast)),
			Team:               derefStr(row.Team),
		}
	}

	return VictimContactsListResult{
		Reports:  items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// reportStatusFrom deriva el estado del reporte — 100% derivado, no columna
// propia (related-tables.md: "Por qué no hay una columna estado separada").
// Regla de convivencia con el flujo legado de /salvia/primer_contacto: un
// reporte con victim_contact_status='i' (invalidado) y sin motivo de cierre
// nuevo se lee como "Cerrado - inválido" (related-tables.md, "Nota — convivencia
// con el botón Invalidar legado"). Son ~65% de los reportes existentes.
func reportStatusFrom(agentID, closeReasonCode, closeCategory *string, legacyStatus string) string {
	if closeReasonCode != nil && *closeReasonCode != "" {
		if closeCategory != nil && *closeCategory == "valido" {
			return "cerrado_valido"
		}
		return "cerrado_invalido"
	}
	if legacyStatus == victimContactLegacyInvalidated {
		return "cerrado_invalido"
	}
	if agentID != nil && *agentID != "" {
		return "asignado"
	}
	return "sin_asignar"
}

func buildVictimContactsListWhere(filters VictimContactsListFilters) (string, []interface{}) {
	clauses := []string{"1=1"}
	args := []interface{}{}

	// estado_reporte (E-02)
	estado := filters.FilterEstadoReporte
	if estado == "" && filters.FilterKey == "estado_reporte" {
		estado = filters.FilterValue
	}
	// Un reporte legado invalidado (status 'i') cuenta como cerrado: no está
	// "asignado" ni "sin asignar" (ver reportStatusFrom).
	const abierto = "vc.victim_contact_close_reason IS NULL AND vc.victim_contact_status <> '" + victimContactLegacyInvalidated + "'"
	switch estado {
	case "asignado":
		clauses = append(clauses, "vc.victim_contact_agent_id IS NOT NULL AND "+abierto)
	case "sin_asignar":
		clauses = append(clauses, "vc.victim_contact_agent_id IS NULL AND "+abierto)
	case "cerrado_valido":
		clauses = append(clauses, `vc.victim_contact_close_reason IN (
			SELECT victim_case_form2_enums_id FROM salvia.victim_case_form2_enums
			WHERE victim_case_form2_enums_category = 'victim_contact_close_reason' AND victim_case_form2_enums_code = 'ap'
		)`)
	case "cerrado_invalido":
		clauses = append(clauses, `(vc.victim_contact_close_reason IN (
			SELECT victim_case_form2_enums_id FROM salvia.victim_case_form2_enums
			WHERE victim_case_form2_enums_category = 'victim_contact_close_reason' AND victim_case_form2_enums_code != 'ap'
		) OR (vc.victim_contact_close_reason IS NULL AND vc.victim_contact_status = '`+victimContactLegacyInvalidated+`'))`)
	}

	// asignacion (E-03): filter_key=asignacion&filter_value={icode|__unassigned__}
	if filters.FilterKey == "asignacion" && filters.FilterValue != "" {
		if filters.FilterValue == "__unassigned__" {
			clauses = append(clauses, "vc.victim_contact_agent_id IS NULL")
		} else {
			clauses = append(clauses, "vc.victim_contact_agent_id = ?")
			args = append(args, filters.FilterValue)
		}
	}

	// persona_asignada (usado por "Mis Reportes" — mismo nombre de filtro que
	// casos, comparado aquí contra victim_contact_agent_id)
	if filters.FilterKey == "persona_asignada" && filters.FilterValue != "" {
		clauses = append(clauses, "vc.victim_contact_agent_id = ?")
		args = append(args, filters.FilterValue)
	}

	// Búsqueda combinada (E-07) — adaptado de buildVictimContactsWithoutCaseWhere
	// (VictimContactDAO.go:316-372): nombres/apellidos directos, teléfono en
	// Form1 O Form2, documento solo en Form1. Aquí los 4 se combinan con OR
	// en un solo parámetro search (allá son 4 inputs AND-combinados).
	search := strings.TrimSpace(filters.Search)
	if search != "" {
		pattern := "%" + search + "%"
		clauses = append(clauses, `(
			vc.victim_contact_names ILIKE ?
			OR vc.victim_contact_last_names ILIKE ?
			OR (vc.victim_contact_names || ' ' || vc.victim_contact_last_names) ILIKE ?
			OR EXISTS (
				SELECT 1 FROM salvia.victim_contact_form2 vcf2
				WHERE vcf2.victim_contact_form2_victim_contact = vc.victim_contact_id
				  AND vcf2.victim_contact_form2_victim_col_phone::text ILIKE ?
			)
			OR EXISTS (
				SELECT 1 FROM salvia.victim_contact_form1 vcf1
				WHERE vcf1.victim_contact_form1_victim_contact = vc.victim_contact_id
				  AND vcf1.victim_contact_form1_phone::text ILIKE ?
			)
			OR EXISTS (
				SELECT 1 FROM salvia.victim_contact_form1 vcf1d
				WHERE vcf1d.victim_contact_form1_victim_contact = vc.victim_contact_id
				  AND vcf1d.victim_contact_form1_doc_number::text ILIKE ?
			)
		)`)
		args = append(args, pattern, pattern, pattern, pattern, pattern, pattern)
	}

	return "WHERE " + strings.Join(clauses, " AND "), args
}

func buildVictimContactsListOrder(filters VictimContactsListFilters) string {
	orderDir := "DESC"
	if strings.EqualFold(filters.Order, "asc") {
		orderDir = "ASC"
	}
	// Única opción soportada: fecha de reporte. "Próximo seguimiento" no
	// aplica a un reporte (no existe follow_up_v2 para victim_contact) —
	// ver flow-E08 de Listado Reportes.
	return fmt.Sprintf("ORDER BY vc.victim_contact_creation_date %s", orderDir)
}
