package repository

import (
	"bitsflow/internal/models"
	"context"

	"gorm.io/gorm"
)

// EntityCaseRepository extiende el contrato CRUD genérico con las consultas
// específicas del componente case-entities.
type EntityCaseRepository interface {
	Repository[models.EntityCase]

	// FindByCaseIDWithRelations devuelve las entidades relacionadas con un caso,
	// enriquecidas con datos de entity_branch/entity, la cadena de ubicación
	// town/city/department, y los conteos de oficios y barreras activas.
	FindByCaseIDWithRelations(ctx context.Context, caseID string) ([]models.EntityCaseWithRelations, error)

	// ExistsActive indica si ya existe una relación activa (no eliminada) entre
	// el caso y la sede dados.
	ExistsActive(ctx context.Context, caseID string, entityBranchID int64) (bool, error)
}

type entityCaseRepository struct {
	repository[models.EntityCase]
	db *gorm.DB
}

func NewEntityCaseRepository(db *gorm.DB) EntityCaseRepository {
	return &entityCaseRepository{
		repository: repository[models.EntityCase]{db: db},
		db:         db,
	}
}

func (r *entityCaseRepository) ExistsActive(ctx context.Context, caseID string, entityBranchID int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.EntityCase{}).
		Where("case_id = ? AND entity_branch_id = ?", caseID, entityBranchID).
		Count(&count).Error
	return count > 0, err
}

// entityCaseWithRelationsSQL enriquece entity_case con:
//   - entity_branch + entity (nombre, sector, dirección)
//   - town → city → department (ubicación completa)
//   - COUNT(entity_letter) por case_id + entity_branch_id → oficiosCount
//   - COUNT(barrier_v2 activas) por case_id + entity_branch_id → barrerasActivasCount
const entityCaseWithRelationsSQL = `
SELECT
    ec.id                       AS rel_id,
    ec.entity_branch_id         AS entity_branch_id,
    eb.entity_branch_i_code     AS entity_branch_icode,
    eb.entity_branch_name       AS entity_branch_name,
    e.entity_sector             AS sector,
    eb.entity_branch_address    AS address,
    COALESCE(d.department_name, '') AS department_name,
    COALESCE(c.city_name, '')       AS city_name,
    COALESCE(t.town_name, '')       AS town_name,
    ec.objetivo                 AS objetivo,
    ec.last_action              AS last_action,
    ec.created_by_id            AS created_by_id,
    ec.created_at               AS created_at,
    (SELECT COUNT(*) FROM salvia.entity_letter el
       WHERE el.case_id = ec.case_id AND el.entity_branch_id = ec.entity_branch_id AND el.deleted_at IS NULL
    ) AS oficios_count,
    (SELECT COUNT(*) FROM salvia.barrier_v2 b
       WHERE b.case_id = ec.case_id AND b.entity_branch_id = ec.entity_branch_id
         AND b.status IN ('OPEN', 'En Gestion') AND b.deleted_at IS NULL
    ) AS barreras_activas_count
FROM salvia.entity_case ec
JOIN salvia.entity_branch eb ON eb.entity_branch_id = ec.entity_branch_id
JOIN salvia.entity e ON e.entity_id = eb.entity_id
LEFT JOIN security.town t ON t.town_code = eb.entity_branch_town_code
LEFT JOIN security.city c ON c.city_id = t.city_id
LEFT JOIN security.department d ON d.department_id = c.department_id
WHERE ec.case_id = ? AND ec.deleted_at IS NULL
ORDER BY ec.created_at ASC
`

func (r *entityCaseRepository) FindByCaseIDWithRelations(ctx context.Context, caseID string) ([]models.EntityCaseWithRelations, error) {
	items := make([]models.EntityCaseWithRelations, 0)
	err := r.db.WithContext(ctx).Raw(entityCaseWithRelationsSQL, caseID).Scan(&items).Error
	return items, err
}
