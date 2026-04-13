// Package repository — case_detail_repository.go
// Repositorio GORM para la pantalla de detalle de caso (rol sv).
package repository

import (
	"bitsflow/internal/models"
	"context"
	"errors"

	"gorm.io/gorm"
)

// CaseDetailData agrupa los datos que necesita la pantalla de detalle.
type CaseDetailData struct {
	Case       models.VictimCase
	Form1      *models.VictimCaseForm1
	FollowUp   *models.FollowUp
	Entries    []models.FollowUpEntry
	FollowUpsV2 []models.FollowUpV2
	TownName    string `json:"townName"`
	DeptName    string `json:"deptName"`
	DenunciasAnteriores int `json:"denunciasAnteriores"`
}

type CaseDetailRepository interface {
	GetByICode(ctx context.Context, caseICode string) (*CaseDetailData, error)
	CreateFollowUpV2(ctx context.Context, followUp *models.FollowUpV2) error
	CountFollowUpsByCaseID(ctx context.Context, caseID string) (int, error)
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

	// Query 4: seguimientos v2 del caso (por icode) — siempre se ejecuta
	var followUpsV2 []models.FollowUpV2
	r.db.WithContext(ctx).
		Where("case_id = ?", vc.VictimCaseICode).
		Order("sequence_number ASC").
		Find(&followUpsV2)
	result.FollowUpsV2 = followUpsV2

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
