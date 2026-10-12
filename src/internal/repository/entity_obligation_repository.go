package repository

import (
	"bitsflow/internal/models"
	"context"

	"gorm.io/gorm"
)

// EntityObligationRepository extiende el contrato CRUD genérico con la consulta
// de catálogo por sede (entity_branch → entity → entity_obligation).
type EntityObligationRepository interface {
	Repository[models.EntityObligation]

	// FindByEntityBranchID resuelve la organización dueña de la sede y devuelve
	// su catálogo de obligaciones, ordenado por "order".
	FindByEntityBranchID(ctx context.Context, entityBranchID int64) ([]models.EntityObligation, error)
}

type entityObligationRepository struct {
	repository[models.EntityObligation]
	db *gorm.DB
}

func NewEntityObligationRepository(db *gorm.DB) EntityObligationRepository {
	return &entityObligationRepository{
		repository: repository[models.EntityObligation]{db: db},
		db:         db,
	}
}

func (r *entityObligationRepository) FindByEntityBranchID(ctx context.Context, entityBranchID int64) ([]models.EntityObligation, error) {
	items := make([]models.EntityObligation, 0)
	err := r.db.WithContext(ctx).Raw(`
SELECT eo.*
FROM salvia.entity_obligation eo
JOIN salvia.entity_branch eb ON eb.entity_id = eo.entity_id
WHERE eb.entity_branch_id = ? AND eo.deleted_at IS NULL
ORDER BY eo."order" ASC
`, entityBranchID).Scan(&items).Error
	return items, err
}
