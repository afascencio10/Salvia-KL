package repository

import (
	"bitsflow/internal/models"
	"context"

	"gorm.io/gorm"
)

type PsychosocialSupportRepository interface {
	Repository[models.PsychosocialSupport]
	FindByCaseID(ctx context.Context, caseID string) ([]models.PsychosocialSupport, error)
	FindByFollowUpID(ctx context.Context, followUpID string) ([]models.PsychosocialSupport, error)
}

type psychosocialSupportRepository struct {
	repository[models.PsychosocialSupport]
	db *gorm.DB
}

func NewPsychosocialSupportRepository(db *gorm.DB) PsychosocialSupportRepository {
	return &psychosocialSupportRepository{
		repository: repository[models.PsychosocialSupport]{db: db},
		db:         db,
	}
}

func (r *psychosocialSupportRepository) FindByCaseID(ctx context.Context, caseID string) ([]models.PsychosocialSupport, error) {
	var items []models.PsychosocialSupport
	return items, r.db.WithContext(ctx).Where("case_id = ?", caseID).Find(&items).Error
}

func (r *psychosocialSupportRepository) FindByFollowUpID(ctx context.Context, followUpID string) ([]models.PsychosocialSupport, error) {
	var items []models.PsychosocialSupport
	return items, r.db.WithContext(ctx).Where("follow_up_id = ?", followUpID).Find(&items).Error
}
