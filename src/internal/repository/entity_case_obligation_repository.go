package repository

import (
	"bitsflow/internal/models"
	"context"

	"gorm.io/gorm"
)

type EntityCaseObligationRepository interface {
	Repository[models.EntityCaseObligation]
	FindByEntityCaseID(ctx context.Context, entityCaseID string) ([]models.EntityCaseObligation, error)
	FindByEntityCaseIDWithLabels(ctx context.Context, entityCaseID, status string) ([]models.EntityCaseObligationItem, error)
}

type entityCaseObligationRepository struct {
	repository[models.EntityCaseObligation]
	db *gorm.DB
}

func NewEntityCaseObligationRepository(db *gorm.DB) EntityCaseObligationRepository {
	return &entityCaseObligationRepository{
		repository: repository[models.EntityCaseObligation]{db: db},
		db:         db,
	}
}

func (r *entityCaseObligationRepository) FindByEntityCaseID(ctx context.Context, entityCaseID string) ([]models.EntityCaseObligation, error) {
	var items []models.EntityCaseObligation
	return items, r.db.WithContext(ctx).Where("entity_case_id = ?", entityCaseID).Order("created_at ASC").Find(&items).Error
}

func (r *entityCaseObligationRepository) FindByEntityCaseIDWithLabels(ctx context.Context, entityCaseID, status string) ([]models.EntityCaseObligationItem, error) {
	items := make([]models.EntityCaseObligationItem, 0)
	q := r.db.WithContext(ctx).Raw(`
SELECT
    eco.id,
    eco.entity_obligation_id,
    eco.status,
    COALESCE(NULLIF(TRIM(eco.custom_label), ''), eo.label, '') AS label
FROM salvia.entity_case_obligation eco
LEFT JOIN salvia.entity_obligation eo ON eo.id = eco.entity_obligation_id
WHERE eco.entity_case_id = ?
  AND (? = '' OR eco.status = ?)
ORDER BY eco.created_at ASC
`, entityCaseID, status, status)
	return items, q.Scan(&items).Error
}
