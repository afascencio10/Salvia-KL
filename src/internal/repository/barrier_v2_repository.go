package repository

import (
	"bitsflow/internal/models"
	"context"
	"strings"

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
	// FindByDepartmentWithRelations lista barreras del department_id dado
	// (todas las no eliminadas), enriquecidas con victim_case, entity_branch y
	// ciudad. Pantalla "Barreras Departamento": no filtra por case_task.
	FindByDepartmentWithRelations(ctx context.Context, departmentID string, filters DepartmentBarrierFilters) ([]models.BarrierV2DepartmentItem, int64, error)
	// FindMineWithRelations lista barreras activas cuyo user_territory_link_id
	// coincide con el Enlace autenticado, con caso y tarea de articulación ToDo.
	FindMineWithRelations(ctx context.Context, userTerritoryLinkID string) ([]models.BarrierV2MineItem, error)
	GetDB() *gorm.DB
}

// DepartmentBarrierFilters agrupa los filtros y la paginación de la pantalla
// "Barreras Departamento".
type DepartmentBarrierFilters struct {
	DocNumber string
	Entidad   string
	CityID    string
	Status    string
	Page      int
	PageSize  int
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

const barrierDepartmentFromSQL = `
FROM salvia.barrier_v2 b
LEFT JOIN salvia.victim_case vc
    ON vc.victim_case_i_code = b.case_id
LEFT JOIN salvia.entity_branch eb
    ON eb.entity_branch_id = b.entity_branch_id
LEFT JOIN security.city ci
    ON ci.city_id::text = b.city_id
WHERE b.deleted_at IS NULL
  AND b.department_id = ?
`

const barrierDepartmentSelectSQL = `
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
    COALESCE(
        (
            SELECT NULLIF(vf2p.victim_case_form2_victim_phone::text, '')
            FROM salvia.victim_case_form2 vf2p
            WHERE vf2p.victim_case_form2_victim_case = vc.victim_case_id
            LIMIT 1
        ),
        (
            SELECT vcf1.victim_contact_form1_phone::text
            FROM salvia.victim_contact_form1 vcf1
            WHERE vcf1.victim_contact_form1_victim_contact = vc.victim_case_victim_contact
            ORDER BY vcf1.victim_contact_form1_id DESC
            LIMIT 1
        ),
        ''
    ) AS victim_phone,
    (
        SELECT vf2.victim_case_form2_risk_level
        FROM salvia.victim_case_form2 vf2
        WHERE vf2.victim_case_form2_victim_case = vc.victim_case_id
        LIMIT 1
    ) AS risk_level
` + barrierDepartmentFromSQL

func departmentBarrierWhereExtra(filters DepartmentBarrierFilters) (string, []interface{}) {
	var extra string
	var args []interface{}
	if doc := strings.TrimSpace(filters.DocNumber); doc != "" {
		extra += " AND vc.victim_case_victim_doc_number ILIKE ?"
		args = append(args, "%"+doc+"%")
	}
	if entidad := strings.TrimSpace(filters.Entidad); entidad != "" {
		extra += " AND eb.entity_branch_name ILIKE ?"
		args = append(args, "%"+entidad+"%")
	}
	if cityID := strings.TrimSpace(filters.CityID); cityID != "" {
		extra += " AND b.city_id = ?"
		args = append(args, cityID)
	}
	if status := strings.TrimSpace(filters.Status); status != "" {
		extra += " AND b.status = ?"
		args = append(args, status)
	}
	return extra, args
}

func (r *barrierV2Repository) FindByDepartmentWithRelations(ctx context.Context, departmentID string, filters DepartmentBarrierFilters) ([]models.BarrierV2DepartmentItem, int64, error) {
	page := filters.Page
	if page < 0 {
		page = 0
	}
	pageSize := filters.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	extra, extraArgs := departmentBarrierWhereExtra(filters)
	args := append([]interface{}{departmentID}, extraArgs...)

	var total int64
	countSQL := "SELECT COUNT(*) " + barrierDepartmentFromSQL + extra
	if err := r.db.WithContext(ctx).Raw(countSQL, args...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	listSQL := barrierDepartmentSelectSQL + extra + `
ORDER BY b.created_at DESC
LIMIT ? OFFSET ?`
	listArgs := append(append([]interface{}{}, args...), pageSize, page*pageSize)

	items := make([]models.BarrierV2DepartmentItem, 0)
	err := r.db.WithContext(ctx).Raw(listSQL, listArgs...).Scan(&items).Error
	return items, total, err
}

const barrierMineSQL = `
SELECT
    b.id,
    b.case_id,
    b.sector,
    b.description,
    b.status,
    b.entity_branch_id,
    COALESCE(eb.entity_branch_name, '')             AS entity_branch_name,
    COALESCE(b.city_id, '')                         AS city_id,
    COALESCE(ci.city_name, '')                      AS city_name,
    COALESCE(vc.victim_case_victim_names, '')       AS victim_name,
    COALESCE(vc.victim_case_victim_last_names, '')  AS victim_last_name,
    COALESCE(vc.victim_case_victim_doc_number, '')  AS victim_doc_number,
    COALESCE(vc.victim_case_i_code, '')             AS case_code,
    COALESCE(
        (
            SELECT NULLIF(vf2p.victim_case_form2_victim_phone::text, '')
            FROM salvia.victim_case_form2 vf2p
            WHERE vf2p.victim_case_form2_victim_case = vc.victim_case_id
            LIMIT 1
        ),
        (
            SELECT vcf1.victim_contact_form1_phone::text
            FROM salvia.victim_contact_form1 vcf1
            WHERE vcf1.victim_contact_form1_victim_contact = vc.victim_case_victim_contact
            ORDER BY vcf1.victim_contact_form1_id DESC
            LIMIT 1
        ),
        ''
    ) AS victim_phone,
    COALESCE(t.id::text, '')                        AS task_id,
    COALESCE(t.type, '')                            AS task_type,
    COALESCE(t.description, '')                     AS task_description
FROM salvia.barrier_v2 b
LEFT JOIN salvia.victim_case vc
    ON vc.victim_case_i_code = b.case_id
LEFT JOIN salvia.entity_branch eb
    ON eb.entity_branch_id = b.entity_branch_id
LEFT JOIN security.city ci
    ON ci.city_id::text = b.city_id
LEFT JOIN LATERAL (
    SELECT ct.id, ct.type, ct.description
    FROM salvia.case_task ct
    WHERE ct.barrier_id = b.id
      AND ct.assigned_user_id = ?
      AND ct.type = ?
      AND ct.status = ?
      AND ct.deleted_at IS NULL
    ORDER BY ct.created_at DESC
    LIMIT 1
) t ON true
WHERE b.deleted_at IS NULL
  AND b.user_territory_link_id = ?
  AND b.status IN (?, ?)
ORDER BY b.created_at DESC
`

func (r *barrierV2Repository) FindMineWithRelations(ctx context.Context, userTerritoryLinkID string) ([]models.BarrierV2MineItem, error) {
	items := make([]models.BarrierV2MineItem, 0)
	err := r.db.WithContext(ctx).Raw(
		barrierMineSQL,
		userTerritoryLinkID,
		models.CaseTaskTypeArticulacionTerritorial,
		models.CaseTaskStatusToDo,
		userTerritoryLinkID,
		models.BarrierV2StatusOpen,
		models.BarrierV2StatusEnGestion,
	).Scan(&items).Error
	return items, err
}
