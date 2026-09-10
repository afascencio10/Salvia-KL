package repository

import (
	"bitsflow/internal/models"
	"context"
	"strings"

	"gorm.io/gorm"
)

// MechanismUserFilter filtra usuarios mc por el tipo y la ubicación de su articulating_mechanism.
type MechanismUserFilter struct {
	Scope        string
	DepartmentID string
	CityID       string
	TownCode     string
}

type CaseArticulatingMechanismRepository interface {
	Create(ctx context.Context, entity *models.CaseArticulatingMechanism) error
	FindByCaseID(ctx context.Context, caseID string) ([]models.CaseArticulatingMechanism, error)
	FindByIDAndCaseID(ctx context.Context, id, caseID string) (*models.CaseArticulatingMechanism, error)
	CreateComment(ctx context.Context, comment *models.CaseArticulatingMechanismComment) error
	ListMechanismUsers(ctx context.Context, filter MechanismUserFilter) ([]models.MechanismUserOption, error)
	FindMechanismUser(ctx context.Context, userICode string) (*models.MechanismUserOption, error)
	ResolveDepartmentName(ctx context.Context, departmentID string) (string, error)
	ResolveCityName(ctx context.Context, cityID string) (string, error)
	ResolveTownName(ctx context.Context, townCode string) (string, error)
	UserFullName(ctx context.Context, userICode string) (string, error)
}

type caseArticulatingMechanismRepository struct {
	db *gorm.DB
}

func NewCaseArticulatingMechanismRepository(db *gorm.DB) CaseArticulatingMechanismRepository {
	return &caseArticulatingMechanismRepository{db: db}
}

func (r *caseArticulatingMechanismRepository) Create(ctx context.Context, entity *models.CaseArticulatingMechanism) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *caseArticulatingMechanismRepository) FindByCaseID(ctx context.Context, caseID string) ([]models.CaseArticulatingMechanism, error) {
	var rows []models.CaseArticulatingMechanism
	err := r.db.WithContext(ctx).
		Where("case_id = ?", caseID).
		Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		Order("created_at DESC").
		Find(&rows).Error
	return rows, err
}

func (r *caseArticulatingMechanismRepository) FindByIDAndCaseID(ctx context.Context, id, caseID string) (*models.CaseArticulatingMechanism, error) {
	var row models.CaseArticulatingMechanism
	err := r.db.WithContext(ctx).
		Where("id = ? AND case_id = ?", id, caseID).
		First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *caseArticulatingMechanismRepository) CreateComment(ctx context.Context, comment *models.CaseArticulatingMechanismComment) error {
	return r.db.WithContext(ctx).Create(comment).Error
}

const mechanismUserSQL = `
SELECT
    gu.general_user_i_code AS icode,
    TRIM(COALESCE(gup.general_user_profile_names, '') || ' ' || COALESCE(gup.general_user_profile_last_names, '')) AS full_name,
    COALESCE(am.town_code, '') AS town_code,
    COALESCE(am.town_name, '') AS town_name,
    COALESCE(am.city_name, '') AS city_name,
    COALESCE(am.department_name, '') AS department_name
FROM security.general_user gu
JOIN security.rel_role_general_user rgu ON rgu.general_user_id = gu.general_user_id
JOIN security.role r ON r.role_id = rgu.role_id
JOIN security.general_user_profile gup ON gup.general_user_profile_id = gu.general_user_general_user_profile
JOIN salvia.articulating_mechanism am ON am.id = gu.articulating_mechanism_id AND am.deleted_at IS NULL
WHERE r.role_code = ?
  AND gu.general_user_status = 'e'
`

func (r *caseArticulatingMechanismRepository) ListMechanismUsers(ctx context.Context, filter MechanismUserFilter) ([]models.MechanismUserOption, error) {
	query := mechanismUserSQL
	args := []interface{}{models.ArticulatingMechanismRoleCode}

	query += " AND am.type = ?"
	args = append(args, strings.TrimSpace(filter.Scope))

	switch filter.Scope {
	case models.ArticulatingMechanismScopeMunicipal:
		query += " AND am.town_code = ?"
		args = append(args, strings.TrimSpace(filter.TownCode))
	case models.ArticulatingMechanismScopeDepartamental:
		query += " AND am.department_id = ?"
		args = append(args, strings.TrimSpace(filter.DepartmentID))
	}

	query += " ORDER BY full_name ASC"

	var rows []models.MechanismUserOption
	if err := r.db.WithContext(ctx).Raw(query, args...).Scan(&rows).Error; err != nil {
		return nil, err
	}
	if rows == nil {
		rows = []models.MechanismUserOption{}
	}
	for i := range rows {
		rows[i].ID = rows[i].ICode
	}
	return rows, nil
}

func (r *caseArticulatingMechanismRepository) FindMechanismUser(ctx context.Context, userICode string) (*models.MechanismUserOption, error) {
	query := mechanismUserSQL + " AND gu.general_user_i_code = ? LIMIT 1"
	var row models.MechanismUserOption
	err := r.db.WithContext(ctx).Raw(query, models.ArticulatingMechanismRoleCode, userICode).Scan(&row).Error
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(row.ICode) == "" {
		return nil, gorm.ErrRecordNotFound
	}
	row.ID = row.ICode
	return &row, nil
}

func (r *caseArticulatingMechanismRepository) ResolveDepartmentName(ctx context.Context, departmentID string) (string, error) {
	var name string
	err := r.db.WithContext(ctx).
		Raw("SELECT department_name FROM security.department WHERE department_id::text = ? LIMIT 1", departmentID).
		Scan(&name).Error
	return name, err
}

func (r *caseArticulatingMechanismRepository) ResolveCityName(ctx context.Context, cityID string) (string, error) {
	var name string
	err := r.db.WithContext(ctx).
		Raw("SELECT city_name FROM security.city WHERE city_id::text = ? LIMIT 1", cityID).
		Scan(&name).Error
	return name, err
}

func (r *caseArticulatingMechanismRepository) ResolveTownName(ctx context.Context, townCode string) (string, error) {
	var name string
	err := r.db.WithContext(ctx).
		Raw("SELECT town_name FROM security.town WHERE town_code = ? LIMIT 1", townCode).
		Scan(&name).Error
	return name, err
}

func (r *caseArticulatingMechanismRepository) UserFullName(ctx context.Context, userICode string) (string, error) {
	var name string
	err := r.db.WithContext(ctx).Raw(`
		SELECT TRIM(COALESCE(gup.general_user_profile_names, '') || ' ' || COALESCE(gup.general_user_profile_last_names, ''))
		FROM security.general_user gu
		JOIN security.general_user_profile gup ON gup.general_user_profile_id = gu.general_user_general_user_profile
		WHERE gu.general_user_i_code = ?
		LIMIT 1
	`, userICode).Scan(&name).Error
	return name, err
}
