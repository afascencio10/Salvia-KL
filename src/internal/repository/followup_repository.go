package repository

import (
	"bitsflow/internal/models"
	"context"

	"gorm.io/gorm"
)

// FollowUpRepository extiende Repository[FollowUpV2] con métodos de negocio propios.
type FollowUpRepository interface {
	Repository[models.FollowUpV2]
	// FindByCaseID retorna todos los seguimientos activos asociados a un caso.
	FindByCaseID(ctx context.Context, caseID string) ([]models.FollowUpV2, error)
}

// followUpRepository implementa FollowUpRepository.
type followUpRepository struct {
	repository[models.FollowUpV2] // embebe el repositorio genérico
	db                            *gorm.DB
}

// NewFollowUpRepository construye un FollowUpRepository listo para usar.
func NewFollowUpRepository(db *gorm.DB) FollowUpRepository {
	return &followUpRepository{
		repository: repository[models.FollowUpV2]{db: db},
		db:         db,
	}
}

// FindByCaseID retorna los registros con case_id = caseID que no hayan sido soft-deleted.
func (r *followUpRepository) FindByCaseID(ctx context.Context, caseID string) ([]models.FollowUpV2, error) {
	var items []models.FollowUpV2
	result := r.db.WithContext(ctx).
		Where("case_id = ?", caseID).
		Find(&items)
	return items, result.Error
}
