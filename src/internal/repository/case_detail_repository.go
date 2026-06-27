// Package repository — case_detail_repository.go
// Repositorio GORM para la pantalla de detalle de caso (rol sv).
package repository

import (
	"bitsflow/internal/models"
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"
)

// CaseDetailData agrupa los datos que necesita la pantalla de detalle.
type CaseDetailData struct {
	Case       models.VictimCase
	Form1      *models.VictimCaseForm1
	Form2      *models.VictimCaseForm2
	FollowUp   *models.FollowUp
	Entries    []models.FollowUpEntry
	FollowUpsV2 []models.FollowUpV2
	TownName    string `json:"townName"`
	DeptName    string `json:"deptName"`
	DenunciasAnteriores int `json:"denunciasAnteriores"`
	AgentName           string `json:"agentName"`
	TimelineEvents []models.CaseTimelineEvent `json:"timelineEvents"`
	EmergencyMeasures      []models.EmergencyMeasure      `json:"emergencyMeasures"`
	PsychosocialSupports   []models.PsychosocialSupport   `json:"psychosocialSupports"`
	EconomicStabilizations []models.EconomicStabilization `json:"economicStabilizations"`
	Barriers               []models.BarrierV2              `json:"barriers"`
	// Campos resumen del caso (form2 enums resueltos)
	TipoViolencia     []string `json:"tipoViolencia"`
	SubtipoViolencia  []string `json:"subtipoViolencia"`
	AmbitoViolencia   []string `json:"ambitoViolencia"`
	Nacionalidad      string   `json:"nacionalidadResumen"`
	Genero            string   `json:"generoResumen"`
	Diversidad        string   `json:"diversidadResumen"`
	LugarHechos       string   `json:"lugarHechos"`
	TerritorioOcurrencia string `json:"territorioOcurrencia"`
	EdadCalculada        *int64 `json:"edadCalculada"`
	TipoAgresorResumen   string `json:"tipoAgresorResumen"`
	NombreIdentitario    string `json:"nombreIdentitario"`
	PlanAtencion         []string `json:"planAtencion"`
	AjusteRazonable      []string `json:"ajusteRazonable"`
}

type CaseDetailRepository interface {
	GetByICode(ctx context.Context, caseICode string) (*CaseDetailData, error)
	CreateFollowUpV2(ctx context.Context, followUp *models.FollowUpV2) error
	CountFollowUpsByCaseID(ctx context.Context, caseID string) (int, error)
	CreateTimelineEvent(ctx context.Context, event *models.CaseTimelineEvent) error
	UpdateFollowUpAgent(ctx context.Context, followUpID string, newAgentID string) error
}

type caseDetailRepository struct {
	db *gorm.DB
}

func NewCaseDetailRepository(db *gorm.DB) CaseDetailRepository {
	return &caseDetailRepository{db: db}
}

func (r *caseDetailRepository) GetByICode(ctx context.Context, caseICode string) (*CaseDetailData, error) {
	result := &CaseDetailData{}

	// Query 1: caso por icode
	var vc models.VictimCase
	if err := r.db.WithContext(ctx).
		Where("victim_case_i_code = ?", caseICode).
		First(&vc).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	result.Case = vc

	// Resolver nombre del agente asignado (agent_id en victim_case)
	if vc.AgentId != nil && *vc.AgentId != "" {
		var agentName string
		r.db.WithContext(ctx).Raw(`
			SELECT COALESCE(gup.general_user_profile_names, '') || ' ' || COALESCE(gup.general_user_profile_last_names, '')
			FROM security.general_user gu
			JOIN security.general_user_profile gup ON gup.general_user_profile_id = gu.general_user_general_user_profile
			WHERE gu.general_user_i_code = ?
		`, *vc.AgentId).Scan(&agentName)
		result.AgentName = strings.TrimSpace(agentName)
	}

	// Resolver nombre del municipio y departamento
	// Cadena: victim_case.town_code → town.city_id → city.city_name + city.department_id → department.department_name
	if vc.VictimCaseTownCode != "" {
		type townCityDept struct {
			CityName string `gorm:"column:city_name"`
			DeptName string `gorm:"column:department_name"`
		}
		var tcd townCityDept
		err := r.db.WithContext(ctx).Raw(`
			SELECT c.city_name, d.department_name
			FROM security.town t
			JOIN security.city c ON c.city_id = t.city_id
			JOIN security.department d ON d.department_id = c.department_id
			WHERE t.town_code = ?
			LIMIT 1
		`, vc.VictimCaseTownCode).Scan(&tcd).Error
		if err == nil {
			result.TownName = tcd.CityName
			result.DeptName = tcd.DeptName
		}
	}

	// Contar denuncias anteriores: cuántos victim_case tienen el mismo doc_number
	if vc.VictimCaseVictimDocNumber != "" {
		var count int64
		r.db.WithContext(ctx).
			Model(&models.VictimCase{}).
			Where("victim_case_victim_doc_number = ?", vc.VictimCaseVictimDocNumber).
			Count(&count)
		if count > 1 {
			result.DenunciasAnteriores = int(count - 1) // restar el caso actual
		}
	}

	// Query 2: form_1 del caso — usa el ID numérico como FK
	var form1 models.VictimCaseForm1
	if err := r.db.WithContext(ctx).
		Where("victim_case_form1_victim_case = ?", vc.VictimCaseId).
		First(&form1).Error; err == nil {
		result.Form1 = &form1
	}

	// Query: form2 del caso (nivel de riesgo)
	var form2 models.VictimCaseForm2
	if err := r.db.WithContext(ctx).
		Where("victim_case_form2_victim_case = ?", vc.VictimCaseId).
		First(&form2).Error; err == nil {
		result.Form2 = &form2

		// Cargar enums relacionados al form2 (tipo violencia, subtipo, ámbito)
		// La tabla relacional usa victim_case_form2_id (PK de form2)
		type enumRow struct {
			Name     string `gorm:"column:enum_name"`
			Category string `gorm:"column:enum_category"`
		}
		var enums []enumRow
		r.db.WithContext(ctx).Raw(`
			SELECT e.victim_case_form2_enums_name AS enum_name,
			       e.victim_case_form2_enums_category AS enum_category
			FROM salvia.rel_victim_case_form2_enums_victim_case_form2 rel
			JOIN salvia.victim_case_form2_enums e ON e.victim_case_form2_enums_id = rel.victim_case_form2_enums_id
			WHERE rel.victim_case_form2_id = (
				SELECT victim_case_form2_id FROM salvia.victim_case_form2
				WHERE victim_case_form2_victim_case = ? LIMIT 1
			)
		`, vc.VictimCaseId).Scan(&enums)

		for _, en := range enums {
			cat := en.Category
			switch {
			case strings.Contains(cat, "subtype_violence_experienced"):
				result.SubtipoViolencia = append(result.SubtipoViolencia, en.Name)
			case strings.Contains(cat, "type_violence_experienced"):
				result.TipoViolencia = append(result.TipoViolencia, en.Name)
			case strings.Contains(cat, "scope_of_violence"):
				result.AmbitoViolencia = append(result.AmbitoViolencia, en.Name)
			case strings.Contains(cat, "action_plan"):
				result.PlanAtencion = append(result.PlanAtencion, en.Name)
			case strings.Contains(cat, "adjustments_gbv"):
				result.AjusteRazonable = append(result.AjusteRazonable, en.Name)
			}
		}

		// Resolver género y nacionalidad (campos de selección única ya resueltos por el componente case-info)
		// Para el resumen usamos los mismos JOINs que case_info_repository
		type resumenRow struct {
			Genero       string `gorm:"column:genero"`
			Nacionalidad string `gorm:"column:nacionalidad"`
			Edad         *int64 `gorm:"column:edad"`
			RelAgresor   string `gorm:"column:rel_agresor"`
		}
		var resumen resumenRow
		r.db.WithContext(ctx).Raw(`
			SELECT
				COALESCE(gi.victim_case_form2_enums_name, '') AS genero,
				COALESCE(na.victim_case_form2_enums_name, '') AS nacionalidad,
				EXTRACT(YEAR FROM AGE(NOW(), f2.victim_case_form2_birth_date))::int AS edad,
				COALESCE(ra.victim_case_form2_enums_name, '') AS rel_agresor
			FROM salvia.victim_case_form2 f2
			LEFT JOIN salvia.victim_case_form2_enums gi ON gi.victim_case_form2_enums_id = f2.victim_case_form2_gender_identity
			LEFT JOIN salvia.victim_case_form2_enums na ON na.victim_case_form2_enums_id = f2.victim_case_form2_nationality
			LEFT JOIN salvia.victim_case_form2_enums ra ON ra.victim_case_form2_enums_id = f2.victim_case_form2_relationship_with_presumed_aggressor
			WHERE f2.victim_case_form2_victim_case = ? LIMIT 1
		`, vc.VictimCaseId).Scan(&resumen)
		result.Genero = resumen.Genero
		result.Nacionalidad = resumen.Nacionalidad
		result.EdadCalculada = resumen.Edad
		result.TipoAgresorResumen = resumen.RelAgresor

		// Nombre identitario (solo form2)
		var identityName string
		r.db.WithContext(ctx).Raw(`SELECT COALESCE(victim_case_form2_identity_name, '') FROM salvia.victim_case_form2 WHERE victim_case_form2_victim_case = ?`, vc.VictimCaseId).Scan(&identityName)
		if identityName != "" {
			result.NombreIdentitario = identityName
		} else {
			result.NombreIdentitario = "No registra"
		}

		// Territorio de ocurrencia (municipio de los hechos)
		if form2FactsTown := ""; true {
			var factsTown string
			r.db.WithContext(ctx).Raw(`
				SELECT COALESCE(c.city_name, '') FROM salvia.victim_case_form2 f2
				LEFT JOIN security.town t ON t.town_code = f2.victim_case_form2_facts_town_code
				LEFT JOIN security.city c ON c.city_id = t.city_id
				WHERE f2.victim_case_form2_victim_case = ? LIMIT 1
			`, vc.VictimCaseId).Scan(&factsTown)
			_ = form2FactsTown
			result.TerritorioOcurrencia = factsTown
		}
	}

	// Query 4: seguimientos v2 del caso (por icode) — ordenados por fecha programada
	var followUpsV2 []models.FollowUpV2
	r.db.WithContext(ctx).
		Where("case_id = ?", vc.VictimCaseICode).
		Order("sequence_number ASC").
		Find(&followUpsV2)

	// Cargar intentos de contacto para cada seguimiento (esquema 3×3)
	for i := range followUpsV2 {
		var attempts []models.FollowUpAttempt
		r.db.WithContext(ctx).
			Where("follow_up_id = ?", followUpsV2[i].ID).
			Order("created_at ASC").
			Find(&attempts)
		followUpsV2[i].FollowUpAttempts = attempts
	}

	result.FollowUpsV2 = followUpsV2

	// Query: eventos del timeline del caso
	var events []models.CaseTimelineEvent
	r.db.WithContext(ctx).
		Where("case_id = ?", vc.VictimCaseICode).
		Order("created_at DESC").
		Find(&events)
	result.TimelineEvents = events

	// Derivaciones: medidas de emergencia, psicosocial, estabilizacion economica
	var em []models.EmergencyMeasure
	r.db.WithContext(ctx).Where("case_id = ?", vc.VictimCaseICode).Find(&em)
	result.EmergencyMeasures = em

	var ps []models.PsychosocialSupport
	r.db.WithContext(ctx).Where("case_id = ?", vc.VictimCaseICode).Find(&ps)
	result.PsychosocialSupports = ps

	var es []models.EconomicStabilization
	r.db.WithContext(ctx).Where("case_id = ?", vc.VictimCaseICode).Find(&es)
	result.EconomicStabilizations = es

	var barriers []models.BarrierV2
	r.db.WithContext(ctx).Where("case_id = ?", vc.VictimCaseICode).Find(&barriers)
	result.Barriers = barriers

	// Si no tiene follow_up asignado, retornar solo el caso, form1 y followUpsV2
	if vc.VictimCaseFollowUpId == nil || *vc.VictimCaseFollowUpId == 0 {
		return result, nil
	}

	// Query 2: follow_up
	var fu models.FollowUp
	if err := r.db.WithContext(ctx).
		Where("follow_up_id = ?", *vc.VictimCaseFollowUpId).
		First(&fu).Error; err == nil {
		result.FollowUp = &fu
	}

	// Query 3: entries del follow_up
	var entries []models.FollowUpEntry
	r.db.WithContext(ctx).
		Where("follow_up_id = ?", *vc.VictimCaseFollowUpId).
		Order("follow_up_entry_completion_date ASC").
		Find(&entries)
	result.Entries = entries

	return result, nil
}

func (r *caseDetailRepository) CreateFollowUpV2(ctx context.Context, followUp *models.FollowUpV2) error {
	return r.db.WithContext(ctx).Create(followUp).Error
}

// CountFollowUpsByCaseID cuenta cuántos seguimientos tiene un caso para calcular el sequence_number.
func (r *caseDetailRepository) CountFollowUpsByCaseID(ctx context.Context, caseID string) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.FollowUpV2{}).
		Where("case_id = ?", caseID).
		Count(&count).Error
	return int(count), err
}

func (r *caseDetailRepository) CreateTimelineEvent(ctx context.Context, event *models.CaseTimelineEvent) error {
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *caseDetailRepository) UpdateFollowUpAgent(ctx context.Context, followUpID string, newAgentID string) error {
	return r.db.WithContext(ctx).
		Model(&models.FollowUpV2{}).
		Where("id = ?", followUpID).
		Update("agent_id", newAgentID).Error
}
