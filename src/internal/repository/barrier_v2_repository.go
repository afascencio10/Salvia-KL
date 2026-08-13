package repository

import (
	"bitsflow/internal/models"
	"context"

	"gorm.io/gorm"
)

type BarrierV2Repository interface {
	Repository[models.BarrierV2]
	FindByCaseID(ctx context.Context, caseID string) ([]models.BarrierV2, error)
	FindByFollowUpID(ctx context.Context, followUpID string) ([]models.BarrierV2, error)
	FindActiveByCaseID(ctx context.Context, caseID string) ([]models.BarrierV2, error)
	FindByIDs(ctx context.Context, ids []string) ([]models.BarrierV2, error)
	// FindActiveByTeamContactIDs retorna las barreras (status != MANAGED) cuyo team_contact_id
	// esté entre los IDs dados — usado para resolver las "barreras activas" de una remisión
	// psicosocial (ver salvia/service/psychosocial_detail_service.go → LoadSession).
	FindActiveByTeamContactIDs(ctx context.Context, teamContactIDs []string) ([]models.BarrierV2, error)
	UpdateStatus(ctx context.Context, id string, status string) error
	// FindByCreatedByIDWithRelations devuelve barreras del agente enriquecidas
	// con datos de victim_case (nombres, documento, i_code).
	FindByCreatedByIDWithRelations(ctx context.Context, createdByID string) ([]models.BarrierV2WithRelations, error)
	// FindActiveByDepartmentWithRelations devuelve todas las barreras activas
	// (status != MANAGED) de los casos del departamento dado, enriquecidas con
	// datos de victim_case, entity_branch y ciudad. Es la consulta principal
	// para la pantalla "Barreras Departamento" — no filtra por asignación de
	// tareas, a diferencia de FindByAssignedUserIDWithRelations (Mis Barreras).
	FindActiveByDepartmentWithRelations(ctx context.Context, departmentID string, filters DepartmentBarrierFilters) ([]models.BarrierV2DepartmentItem, error)
	GetDB() *gorm.DB
}

// DepartmentBarrierFilters agrupa los filtros opcionales de la pantalla
// "Barreras Departamento": documento de identidad, entidad y ciudad.
type DepartmentBarrierFilters struct {
	DocNumber string
	Entidad   string
	CityID    string
}

type barrierV2Repository struct {
	repository[models.BarrierV2]
	db *gorm.DB
}

func NewBarrierV2Repository(db *gorm.DB) BarrierV2Repository {
	return &barrierV2Repository{
		repository: repository[models.BarrierV2]{db: db},
		db:         db,
	}
}

func (r *barrierV2Repository) GetDB() *gorm.DB {
	return r.db
}

func (r *barrierV2Repository) FindByCaseID(ctx context.Context, caseID string) ([]models.BarrierV2, error) {
	var items []models.BarrierV2
	return items, r.db.WithContext(ctx).Where("case_id = ?", caseID).Find(&items).Error
}

func (r *barrierV2Repository) FindByFollowUpID(ctx context.Context, followUpID string) ([]models.BarrierV2, error) {
	var items []models.BarrierV2
	return items, r.db.WithContext(ctx).Where("follow_up_id = ?", followUpID).Find(&items).Error
}

func (r *barrierV2Repository) FindActiveByCaseID(ctx context.Context, caseID string) ([]models.BarrierV2, error) {
	var items []models.BarrierV2
	return items, r.db.WithContext(ctx).
		Where("case_id = ? AND status != ? AND deleted_at IS NULL", caseID, models.BarrierV2StatusManaged).
		Find(&items).Error
}

func (r *barrierV2Repository) FindByIDs(ctx context.Context, ids []string) ([]models.BarrierV2, error) {
	var items []models.BarrierV2
	if len(ids) == 0 {
		return items, nil
	}
	return items, r.db.WithContext(ctx).Where("id IN ?", ids).Find(&items).Error
}

func (r *barrierV2Repository) FindActiveByTeamContactIDs(ctx context.Context, teamContactIDs []string) ([]models.BarrierV2, error) {
	var items []models.BarrierV2
	if len(teamContactIDs) == 0 {
		return items, nil
	}
	return items, r.db.WithContext(ctx).
		Where("team_contact_id IN ? AND status != ? AND deleted_at IS NULL", teamContactIDs, models.BarrierV2StatusManaged).
		Order("created_at ASC").
		Find(&items).Error
}

func (r *barrierV2Repository) UpdateStatus(ctx context.Context, id string, status string) error {
	return r.db.WithContext(ctx).
		Model(&models.BarrierV2{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// barrierWithRelationsSQL es la query base que enriquece barrier_v2 con datos
// de victim_case mediante LEFT JOIN.
//
//	barrier_v2.case_id → victim_case.victim_case_i_code
const barrierWithRelationsSQL = `
SELECT
    b.id,
    b.case_id,
    b.follow_up_id,
    b.sector,
    b.description,
    b.status,
    b.created_by_id,
    b.created_at,
    COALESCE(vc.victim_case_victim_names, '')       AS victim_name,
    COALESCE(vc.victim_case_victim_last_names, '')  AS victim_last_name,
    COALESCE(vc.victim_case_victim_doc_number, '')  AS victim_doc_number,
    COALESCE(vc.victim_case_i_code, '')             AS case_code
FROM salvia.barrier_v2 b
LEFT JOIN salvia.victim_case vc
    ON vc.victim_case_i_code = b.case_id
WHERE b.deleted_at IS NULL
`

func (r *barrierV2Repository) FindByCreatedByIDWithRelations(ctx context.Context, createdByID string) ([]models.BarrierV2WithRelations, error) {
	var items []models.BarrierV2WithRelations
	err := r.db.WithContext(ctx).
		Raw(barrierWithRelationsSQL+" AND b.created_by_id = ? ORDER BY b.created_at DESC", createdByID).
		Scan(&items).Error
	return items, err
}

// barrierDepartmentWithRelationsSQL enriquece barrier_v2 con datos de
// victim_case, victim_case_form2 (nivel de riesgo/prioridad), entity_branch
// (entidad involucrada) y security.city (ciudad) mediante LEFT JOIN.
const barrierDepartmentWithRelationsSQL = `
SELECT
    b.id,
    b.case_id,
    b.sector,
    b.description,
    b.status,
    b.barrier_date,
    b.created_at,
    COALESCE(eb.entity_branch_name, '')             AS entity_branch_name,
    COALESCE(ci.city_name, '')                      AS city_name,
    COALESCE(vc.victim_case_victim_names, '')       AS victim_name,
    COALESCE(vc.victim_case_victim_last_names, '')  AS victim_last_name,
    COALESCE(vc.victim_case_victim_doc_number, '')  AS victim_doc_number,
    COALESCE(vc.victim_case_i_code, '')             AS case_code,
    vf2.victim_case_form2_risk_level                AS risk_level
FROM salvia.barrier_v2 b
LEFT JOIN salvia.victim_case vc
    ON vc.victim_case_i_code = b.case_id
LEFT JOIN salvia.victim_case_form2 vf2
    ON vf2.victim_case_form2_victim_case = vc.victim_case_id
LEFT JOIN salvia.entity_branch eb
    ON eb.entity_branch_id = b.entity_branch_id
LEFT JOIN security.city ci
    ON ci.city_id = b.city_id
WHERE b.deleted_at IS NULL
  AND b.status != ?
  AND b.department_id = ?
`

func (r *barrierV2Repository) FindActiveByDepartmentWithRelations(ctx context.Context, departmentID string, filters DepartmentBarrierFilters) ([]models.BarrierV2DepartmentItem, error) {
	sql := barrierDepartmentWithRelationsSQL
	args := []interface{}{models.BarrierV2StatusManaged, departmentID}

	if filters.DocNumber != "" {
		sql += " AND vc.victim_case_victim_doc_number ILIKE ?"
		args = append(args, "%"+filters.DocNumber+"%")
	}
	if filters.Entidad != "" {
		sql += " AND eb.entity_branch_name ILIKE ?"
		args = append(args, "%"+filters.Entidad+"%")
	}
	if filters.CityID != "" {
		sql += " AND b.city_id = ?"
		args = append(args, filters.CityID)
	}

	sql += " ORDER BY b.created_at DESC"

	var items []models.BarrierV2DepartmentItem
	err := r.db.WithContext(ctx).Raw(sql, args...).Scan(&items).Error
	return items, err
}
