package repository

import (
	"bitsflow/internal/models"
	"context"

	"gorm.io/gorm"
)

type VisibilityConditionRepository interface {
	Repository[models.VisibilityCondition]
	FindByTriggerQuestion(ctx context.Context, questionID string) ([]models.VisibilityCondition, error)
}

type visibilityConditionRepository struct {
	repository[models.VisibilityCondition]
	db *gorm.DB
}

func NewVisibilityConditionRepository(db *gorm.DB) VisibilityConditionRepository {
	return &visibilityConditionRepository{
		repository: repository[models.VisibilityCondition]{db: db},
		db:         db,
	}
}

func (r *visibilityConditionRepository) FindByTriggerQuestion(ctx context.Context, questionID string) ([]models.VisibilityCondition, error) {
	var items []models.VisibilityCondition
	return items, r.db.WithContext(ctx).Where("trigger_question_id = ?", questionID).Find(&items).Error
}
