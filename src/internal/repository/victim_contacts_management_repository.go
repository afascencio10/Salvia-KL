package repository

import (
	"bitsflow/internal/models"
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

var ErrVictimContactNotFound = errors.New("reporte no encontrado")
var ErrVictimContactAlreadyClosed = errors.New("el reporte ya está cerrado")
var ErrSchemeNotComplete = errors.New("el esquema 3x3 no está completo todavía")
var ErrCloseReasonNotFound = errors.New("motivo de cierre no encontrado")

// CloseReasonLegacy — "código" sintético que se devuelve para los reportes invalidados
// por el flujo legado (no existe como fila de enum): hace que el modal abra en 'cerrado'.
const CloseReasonLegacy = "legado"

// closeReasonLabels — texto en UI de cada motivo de cierre (related-tables.md). La
// observación (co/dp) se concatena en el cliente, no aquí.
var closeReasonLabels = map[string]string{
	"ni": "Inválido por No información Mínima",
	"nc": "Inválido por Cierre por no contacto",
	"co": "Inválido por Cierre por no competencia",
	"dp": "Inválido por Cierre por decisión de la persona",
	"ap": "Válido - Apertura de Caso Salvia",
}

// VictimContactsManagementRepository acceso a datos para el detalle,
// intentos y cierre de un reporte (gestion-reporte-modal). Ver
// DocsMD/Componentes/Gestion Reporte Modal/related-tables.md.
type VictimContactsManagementRepository interface {
	GetDetail(ctx context.Context, icode string) (*models.ReportDetail, error)
	InsertAttempt(ctx context.Context, icode, result, agentICode string) (intentosHoy int, diasCompletos int, err error)
	Close(ctx context.Context, icode, closeReasonCode, observation string) (categoryLabel string, err error)
}

type victimContactsManagementRepository struct {
	db *gorm.DB
}

func NewVictimContactsManagementRepository(db *gorm.DB) VictimContactsManagementRepository {
	return &victimContactsManagementRepository{db: db}
}

type victimContactRow struct {
	ID              int64   `gorm:"column:victim_contact_id"`
	ICode           string  `gorm:"column:victim_contact_i_code"`
	Names           string  `gorm:"column:victim_contact_names"`
	LastNames       string  `gorm:"column:victim_contact_last_names"`
	AgentID         *string `gorm:"column:victim_contact_agent_id"`
	CloseReasonID   *int64  `gorm:"column:victim_contact_close_reason"`
	CloseReasonCode *string `gorm:"column:close_reason_code"`
	CloseReasonName *string `gorm:"column:close_reason_name"`
	CloseObservation *string `gorm:"column:victim_contact_close_observation"`
	LegacyStatus    string  `gorm:"column:victim_contact_status"`
	LegacyStatusDescription *string `gorm:"column:victim_contact_status_description"`
	ReportType      *string `gorm:"column:report_type"`
	ReporterNames   *string `gorm:"column:reporter_names"`
	ReporterLast    *string `gorm:"column:reporter_last_names"`
	Parentesco      *string `gorm:"column:parentesco"`
	ReporterEmail   *string `gorm:"column:reporter_email"`
	Telefono        *string `gorm:"column:telefono"`
	AjusteRazonable *string `gorm:"column:ajuste_razonable"`
	Hechos          *string `gorm:"column:descripcion_hechos"`
}

// GetDetail carga el reporte + su origen (form1/form2/form2_v2, GAP-1) + historial de intentos +
// contadores del esquema 3x3. Parentesco y ajustes razonables se devuelven como clave de traducción
// `categoría_código` (la misma que usa el wizard); el servicio las convierte a texto con Locale.
func (r *victimContactsManagementRepository) GetDetail(ctx context.Context, icode string) (*models.ReportDetail, error) {
	const query = `
		SELECT
			vc.victim_contact_id,
			vc.victim_contact_i_code,
			vc.victim_contact_names,
			vc.victim_contact_last_names,
			vc.victim_contact_agent_id,
			vc.victim_contact_close_reason,
			closeEnum.victim_case_form2_enums_code AS close_reason_code,
			closeEnum.victim_case_form2_enums_name AS close_reason_name,
			vc.victim_contact_close_observation,
			vc.victim_contact_status,
			vc.victim_contact_status_description,
			(` + reportTypeSelect + `) AS report_type,
			(SELECT vf.victim_contact_form2_reporter_names FROM salvia.victim_contact_form2_v2 v2
			   JOIN salvia.victim_contact_form2 vf ON vf.victim_contact_form2_id = v2.victim_contact_form2_v2_victim_contact_form2
			  WHERE vf.victim_contact_form2_victim_contact = vc.victim_contact_id LIMIT 1) AS reporter_names,
			(SELECT v2.victim_contact_form2_v2_reporter_last_names FROM salvia.victim_contact_form2_v2 v2
			   JOIN salvia.victim_contact_form2 vf ON vf.victim_contact_form2_id = v2.victim_contact_form2_v2_victim_contact_form2
			  WHERE vf.victim_contact_form2_victim_contact = vc.victim_contact_id LIMIT 1) AS reporter_last_names,
			(SELECT pe.victim_case_form2_enums_category || '_' || pe.victim_case_form2_enums_code FROM salvia.victim_contact_form2_v2 v2
			   JOIN salvia.victim_contact_form2 vf ON vf.victim_contact_form2_id = v2.victim_contact_form2_v2_victim_contact_form2
			   LEFT JOIN salvia.victim_case_form2_enums pe ON pe.victim_case_form2_enums_id = v2.victim_contact_form2_v2_parentesco
			  WHERE vf.victim_contact_form2_victim_contact = vc.victim_contact_id LIMIT 1) AS parentesco,
			(SELECT v2.victim_contact_form2_v2_reporter_email FROM salvia.victim_contact_form2_v2 v2
			   JOIN salvia.victim_contact_form2 vf ON vf.victim_contact_form2_id = v2.victim_contact_form2_v2_victim_contact_form2
			  WHERE vf.victim_contact_form2_victim_contact = vc.victim_contact_id LIMIT 1) AS reporter_email,
			COALESCE(
				(SELECT NULLIF(vf.victim_contact_form2_victim_col_phone, 0)::text FROM salvia.victim_contact_form2 vf
				  WHERE vf.victim_contact_form2_victim_contact = vc.victim_contact_id
				  ORDER BY vf.victim_contact_form2_id DESC LIMIT 1),
				(SELECT NULLIF(vf.victim_contact_form2_reporter_phone, 0)::text FROM salvia.victim_contact_form2 vf
				  WHERE vf.victim_contact_form2_victim_contact = vc.victim_contact_id
				  ORDER BY vf.victim_contact_form2_id DESC LIMIT 1),
				(SELECT vf1.victim_contact_form1_phone FROM salvia.victim_contact_form1 vf1
				  WHERE vf1.victim_contact_form1_victim_contact = vc.victim_contact_id
				  ORDER BY vf1.victim_contact_form1_id DESC LIMIT 1)
			) AS telefono,
			(SELECT string_agg(e.victim_case_form2_enums_category || '_' || e.victim_case_form2_enums_code, ', ')
			   FROM salvia.victim_contact_form2 vf
			   JOIN salvia.rel_victim_case_form2_enums_victim_contact_form2 rel
			     ON rel.victim_contact_form2_id = vf.victim_contact_form2_id
			   JOIN salvia.victim_case_form2_enums e ON e.victim_case_form2_enums_id = rel.victim_case_form2_enums_id
			  WHERE vf.victim_contact_form2_victim_contact = vc.victim_contact_id) AS ajuste_razonable,
			COALESCE(
				(SELECT vf.victim_contact_form2_facts_description FROM salvia.victim_contact_form2 vf
				  WHERE vf.victim_contact_form2_victim_contact = vc.victim_contact_id
				  ORDER BY vf.victim_contact_form2_id DESC LIMIT 1),
				(SELECT vf1.victim_contact_form1_facts_description FROM salvia.victim_contact_form1 vf1
				  WHERE vf1.victim_contact_form1_victim_contact = vc.victim_contact_id
				  ORDER BY vf1.victim_contact_form1_id DESC LIMIT 1)
			) AS descripcion_hechos
		FROM salvia.victim_contact vc
		LEFT JOIN salvia.victim_case_form2_enums closeEnum
		       ON closeEnum.victim_case_form2_enums_id = vc.victim_contact_close_reason
		WHERE vc.victim_contact_i_code = ?
	`

	var row victimContactRow
	if err := r.db.WithContext(ctx).Raw(query, icode).Scan(&row).Error; err != nil {
		return nil, err
	}
	if row.ICode == "" {
		return nil, ErrVictimContactNotFound
	}

	var attempts []struct {
		Date      time.Time `gorm:"column:victim_contact_attempt_date"`
		Result    string    `gorm:"column:victim_contact_attempt_result"`
		AgentName string    `gorm:"column:agent_name"`
	}
	const attemptsQuery = `
		SELECT vca.victim_contact_attempt_date, vca.victim_contact_attempt_result,
		       COALESCE(gup.general_user_profile_names || ' ' || gup.general_user_profile_last_names, '') AS agent_name
		FROM salvia.victim_contact_attempt vca
		LEFT JOIN security.general_user gu ON gu.general_user_i_code = vca.victim_contact_attempt_agent
		LEFT JOIN security.general_user_profile gup ON gup.general_user_profile_id = gu.general_user_general_user_profile
		WHERE vca.victim_contact_attempt_victim_contact = ?
		  AND vca.victim_contact_attempt_deleted_at IS NULL
		ORDER BY vca.victim_contact_attempt_date ASC
	`
	if err := r.db.WithContext(ctx).Raw(attemptsQuery, icode).Scan(&attempts).Error; err != nil {
		return nil, err
	}

	items := make([]models.ReportAttemptItem, len(attempts))
	intentosHoy := 0
	fallidosPorDia := map[string]int{}
	today := time.Now().In(bogotaLocation()).Format("2006-01-02")
	for i, a := range attempts {
		items[i] = models.ReportAttemptItem{Date: a.Date, Result: a.Result, AgentName: a.AgentName}
		if a.Result == "fallido" {
			dia := a.Date.In(bogotaLocation()).Format("2006-01-02")
			fallidosPorDia[dia]++
			if dia == today {
				intentosHoy++
			}
		}
	}
	diasCompletos := 0
	for _, n := range fallidosPorDia {
		if n >= 3 {
			diasCompletos++
		}
	}

	nombreMostrado := strings.TrimSpace(row.Names + " " + row.LastNames)
	closeCategory := ""
	closeCode := derefStr(row.CloseReasonCode)
	closeLabel := closeReasonLabels[closeCode]
	if closeCode != "" {
		if closeCode == "ap" {
			closeCategory = "Válido"
		} else {
			closeCategory = "Inválido"
		}
	} else if row.LegacyStatus == victimContactLegacyInvalidated {
		// Reporte invalidado por el flujo legado de /salvia/primer_contacto: se lee como
		// "Cerrado - inválido" de solo lectura, con su descripción como motivo (related-tables.md).
		closeCode = CloseReasonLegacy
		closeCategory = "Inválido"
		closeLabel = strings.TrimSpace(derefStr(row.LegacyStatusDescription))
		if closeLabel == "" {
			closeLabel = "Invalidado en el flujo anterior"
		}
	}

	canManage := closeCode == ""

	return &models.ReportDetail{
		ICode:             row.ICode,
		NombreMostrado:    nombreMostrado,
		ReportType:        derefStr(row.ReportType),
		ReporterNames:     derefStr(row.ReporterNames),
		ReporterLastNames: derefStr(row.ReporterLast),
		Parentesco:        derefStr(row.Parentesco),
		Telefono:          derefStr(row.Telefono),
		ReporterEmail:     derefStr(row.ReporterEmail),
		AjusteRazonable:   derefStr(row.AjusteRazonable),
		DescripcionHechos: derefStr(row.Hechos),
		Attempts:          items,
		IntentosHoy:       intentosHoy,
		DiasCompletos:     diasCompletos,
		CloseReasonCode:   closeCode,
		CloseReasonLabel:  closeLabel,
		CloseCategory:     closeCategory,
		CloseObservation:  derefStr(row.CloseObservation),
		CanManage:         canManage,
	}, nil
}

func bogotaLocation() *time.Location {
	loc, err := time.LoadLocation("America/Bogota")
	if err != nil {
		return time.UTC
	}
	return loc
}

// InsertAttempt registra un intento y recalcula los contadores del esquema
// 3x3. "Día completo" = día con >= 3 intentos fallidos — ver corrección
// 2026-10-12 en related-tables.md.
func (r *victimContactsManagementRepository) InsertAttempt(ctx context.Context, icode, result, agentICode string) (int, int, error) {
	var vc struct {
		ID          int64  `gorm:"column:victim_contact_id"`
		ICode       string `gorm:"column:victim_contact_i_code"`
		CloseReason *int64 `gorm:"column:victim_contact_close_reason"`
		Status      string `gorm:"column:victim_contact_status"`
	}
	if err := r.db.WithContext(ctx).Raw(`SELECT victim_contact_id, victim_contact_i_code, victim_contact_close_reason, victim_contact_status FROM salvia.victim_contact WHERE victim_contact_i_code = ?`, icode).Scan(&vc).Error; err != nil {
		return 0, 0, err
	}
	if vc.ICode == "" {
		return 0, 0, ErrVictimContactNotFound
	}
	if vc.CloseReason != nil || vc.Status == victimContactLegacyInvalidated {
		return 0, 0, ErrVictimContactAlreadyClosed
	}

	// Se guarda en UTC explícito: la columna es `timestamp without time zone` y todo el pipeline de
	// lectura asume UTC; con time.Now() quedaba la hora local del servidor (distinta en dev y en Docker).
	now := time.Now().UTC()
	if err := r.db.WithContext(ctx).Exec(`
		INSERT INTO salvia.victim_contact_attempt (
			victim_contact_attempt_i_code, victim_contact_attempt_victim_contact,
			victim_contact_attempt_result, victim_contact_attempt_date,
			victim_contact_attempt_agent, victim_contact_attempt_creation_date
		) VALUES (gen_random_uuid(), ?, ?, ?, ?, ?)
	`, icode, result, now, agentICode, now).Error; err != nil {
		return 0, 0, err
	}

	intentosHoy, diasCompletos, err := r.computeCounters(ctx, icode)
	return intentosHoy, diasCompletos, err
}

func (r *victimContactsManagementRepository) computeCounters(ctx context.Context, icode string) (int, int, error) {
	var intentosHoy int
	if err := r.db.WithContext(ctx).Raw(`
		SELECT COUNT(*) FROM salvia.victim_contact_attempt
		WHERE victim_contact_attempt_victim_contact = ?
		  AND victim_contact_attempt_result = 'fallido'
		  AND victim_contact_attempt_deleted_at IS NULL
		  AND (victim_contact_attempt_date AT TIME ZONE 'UTC' AT TIME ZONE 'America/Bogota')::date = (NOW() AT TIME ZONE 'America/Bogota')::date
	`, icode).Scan(&intentosHoy).Error; err != nil {
		return 0, 0, err
	}

	var diasCompletos int
	if err := r.db.WithContext(ctx).Raw(`
		SELECT COUNT(*) FROM (
			SELECT (victim_contact_attempt_date AT TIME ZONE 'UTC' AT TIME ZONE 'America/Bogota')::date AS dia, COUNT(*) AS n
			FROM salvia.victim_contact_attempt
			WHERE victim_contact_attempt_victim_contact = ?
			  AND victim_contact_attempt_result = 'fallido'
			  AND victim_contact_attempt_deleted_at IS NULL
			GROUP BY dia
		) dias WHERE dias.n >= 3
	`, icode).Scan(&diasCompletos).Error; err != nil {
		return 0, 0, err
	}
	return intentosHoy, diasCompletos, nil
}

// Close resuelve el enum de close_reason y actualiza victim_contact. Para
// code='nc' (cierre por no contacto) valida server-side que el esquema 3x3
// esté realmente completo — ver flow-M07.
func (r *victimContactsManagementRepository) Close(ctx context.Context, icode, closeReasonCode, observation string) (string, error) {
	var vc struct {
		ID          int64  `gorm:"column:victim_contact_id"`
		ICode       string `gorm:"column:victim_contact_i_code"`
		CloseReason *int64 `gorm:"column:victim_contact_close_reason"`
		Status      string `gorm:"column:victim_contact_status"`
	}
	if err := r.db.WithContext(ctx).Raw(`SELECT victim_contact_id, victim_contact_i_code, victim_contact_close_reason, victim_contact_status FROM salvia.victim_contact WHERE victim_contact_i_code = ?`, icode).Scan(&vc).Error; err != nil {
		return "", err
	}
	if vc.ICode == "" {
		return "", ErrVictimContactNotFound
	}
	if vc.CloseReason != nil || vc.Status == victimContactLegacyInvalidated {
		return "", ErrVictimContactAlreadyClosed
	}

	if closeReasonCode == "nc" {
		_, diasCompletos, err := r.computeCounters(ctx, icode)
		if err != nil {
			return "", err
		}
		if diasCompletos < 3 {
			return "", ErrSchemeNotComplete
		}
	}

	var enumRow struct {
		ID       int64  `gorm:"column:victim_case_form2_enums_id"`
		Code     string `gorm:"column:victim_case_form2_enums_code"`
	}
	if err := r.db.WithContext(ctx).Raw(`
		SELECT victim_case_form2_enums_id, victim_case_form2_enums_code
		FROM salvia.victim_case_form2_enums
		WHERE victim_case_form2_enums_category = 'victim_contact_close_reason' AND victim_case_form2_enums_code = ?
	`, closeReasonCode).Scan(&enumRow).Error; err != nil {
		return "", err
	}
	if enumRow.ID == 0 {
		return "", ErrCloseReasonNotFound
	}

	var observationArg interface{}
	if strings.TrimSpace(observation) != "" {
		observationArg = strings.TrimSpace(observation)
	}

	if err := r.db.WithContext(ctx).Exec(`
		UPDATE salvia.victim_contact SET
			victim_contact_close_reason = ?,
			victim_contact_close_observation = ?,
			victim_contact_closed_date = ?
		WHERE victim_contact_i_code = ?
	`, enumRow.ID, observationArg, time.Now().UTC(), icode).Error; err != nil {
		return "", err
	}

	if enumRow.Code == "ap" {
		return "Válido", nil
	}
	return "Inválido", nil
}
