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
	FindGestionesByMechanismIDs(ctx context.Context, ids []string) ([]models.MechanismGestionItem, error)
	FindByIDAndCaseID(ctx context.Context, id, caseID string) (*models.CaseArticulatingMechanism, error)
	FindByID(ctx context.Context, id string) (*models.CaseArticulatingMechanism, error)
	FindMinePaginated(ctx context.Context, assignedUserID, document, city string, page, pageSize int) ([]models.CaseMechanismMineItem, int64, error)
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
		Order("created_at DESC").
		Find(&rows).Error
	return rows, err
}

func (r *caseArticulatingMechanismRepository) FindGestionesByMechanismIDs(ctx context.Context, ids []string) ([]models.MechanismGestionItem, error) {
	if len(ids) == 0 {
		return []models.MechanismGestionItem{}, nil
	}
	const q = `
SELECT
    ct.id AS id,
    ct.case_articulating_mechanism_id AS mechanism_id,
    ct.result AS result,
    ct.form_data AS form_data,
    ct.completed_at AS completed_at,
    TRIM(COALESCE(gup.general_user_profile_names, '') || ' ' || COALESCE(gup.general_user_profile_last_names, '')) AS assigned_user_name
FROM salvia.case_task ct
LEFT JOIN security.general_user gu ON gu.general_user_i_code = ct.assigned_user_id
LEFT JOIN security.general_user_profile gup ON gup.general_user_profile_id = gu.general_user_general_user_profile
WHERE ct.deleted_at IS NULL
  AND ct.status = ?
  AND ct.type = ?
  AND ct.case_articulating_mechanism_id IN ?
ORDER BY ct.completed_at ASC NULLS LAST, ct.created_at ASC
`
	items := make([]models.MechanismGestionItem, 0)
	err := r.db.WithContext(ctx).Raw(q, models.CaseTaskStatusDone, models.CaseTaskTypeGestionMecanismo, ids).Scan(&items).Error
	return items, err
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

func (r *caseArticulatingMechanismRepository) FindByID(ctx context.Context, id string) (*models.CaseArticulatingMechanism, error) {
	var row models.CaseArticulatingMechanism
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

const mechanismMineSelectSQL = `
SELECT
    cam.id AS id,
    cam.case_id AS case_id,
    vc.victim_case_i_code AS case_code,
    TRIM(CONCAT(COALESCE(vc.victim_case_victim_names, ''), ' ', COALESCE(vc.victim_case_victim_last_names, ''))) AS victim_full_name,
    COALESCE(vc.victim_case_victim_doc_number, '') AS document,
    COALESCE(NULLIF(vf2.victim_case_form2_victim_phone::text, ''), '') AS phone,
    COALESCE(case_city.city_name, '') AS city,
    (` + riskLevelSelect + `) AS risk_level,
    COALESCE(vc.victim_case_status, '') AS case_status,
    cam.scope AS scope,
    cam.salvia_commitment AS salvia_commitment,
    last_task.result AS last_comment,
    last_task.completed_at AS last_comment_at
FROM salvia.case_articulating_mechanism cam
JOIN salvia.victim_case vc ON vc.victim_case_i_code = cam.case_id
LEFT JOIN salvia.victim_case_form2 vf2 ON vf2.victim_case_form2_victim_case = vc.victim_case_id
LEFT JOIN security.town case_town ON case_town.town_code = vc.victim_case_victim_town_code
LEFT JOIN security.city case_city ON case_city.city_id = case_town.city_id
LEFT JOIN LATERAL (
    SELECT ct.result, ct.completed_at
    FROM salvia.case_task ct
    WHERE ct.case_articulating_mechanism_id = cam.id
      AND ct.status = 'Done'
      AND ct.deleted_at IS NULL
    ORDER BY ct.completed_at DESC NULLS LAST, ct.created_at DESC
    LIMIT 1
) last_task ON true
WHERE cam.deleted_at IS NULL
  AND cam.status = ?
  AND cam.assigned_user_id = ?
`

func (r *caseArticulatingMechanismRepository) FindMinePaginated(ctx context.Context, assignedUserID, document, city string, page, pageSize int) ([]models.CaseMechanismMineItem, int64, error) {
	if page < 0 {
		page = 0
	}
	if pageSize <= 0 {
		pageSize = 5
	}

	args := []interface{}{models.CaseArticulatingMechanismStatusActive, assignedUserID}
	whereExtra := ""

	doc := strings.TrimSpace(document)
	if doc != "" {
		whereExtra += " AND vc.victim_case_victim_doc_number ILIKE ?"
		args = append(args, "%"+doc+"%")
	}
	cityQ := strings.TrimSpace(city)
	if cityQ != "" {
		whereExtra += " AND COALESCE(case_city.city_name, '') ILIKE ?"
		args = append(args, "%"+cityQ+"%")
	}

	countSQL := `
SELECT COUNT(*)
FROM salvia.case_articulating_mechanism cam
JOIN salvia.victim_case vc ON vc.victim_case_i_code = cam.case_id
LEFT JOIN security.town case_town ON case_town.town_code = vc.victim_case_victim_town_code
LEFT JOIN security.city case_city ON case_city.city_id = case_town.city_id
WHERE cam.deleted_at IS NULL
  AND cam.status = ?
  AND cam.assigned_user_id = ?` + whereExtra

	var total int64
	if err := r.db.WithContext(ctx).Raw(countSQL, args...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	listSQL := mechanismMineSelectSQL + whereExtra + `
ORDER BY cam.created_at DESC
LIMIT ? OFFSET ?`
	listArgs := append(append([]interface{}{}, args...), pageSize, page*pageSize)

	items := make([]models.CaseMechanismMineItem, 0)
	if err := r.db.WithContext(ctx).Raw(listSQL, listArgs...).Scan(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
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
