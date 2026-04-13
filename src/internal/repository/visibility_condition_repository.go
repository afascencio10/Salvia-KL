package repository

import (
	"bitsflow/internal/models"
	"context"

	"gorm.io/gorm"
)

type VisibilityConditionRepository interface {
	Repository[models.VisibilityCondition]
	FindByTargetID(ctx context.Context, targetID string) ([]models.VisibilityCondition, error)
	FindByTargetIDs(ctx context.Context, ids []string) ([]models.VisibilityCondition, error)
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

func (r *visibilityConditionRepository) FindByTargetID(ctx context.Context, targetID string) ([]models.VisibilityCondition, error) {
	var items []models.VisibilityCondition
	return items, r.db.WithContext(ctx).Where("target_id = ?", targetID).Find(&items).Error
}

func (r *visibilityConditionRepository) FindByTargetIDs(ctx context.Context, ids []string) ([]models.VisibilityCondition, error) {
	var items []models.VisibilityCondition
	if len(ids) == 0 {
		return items, nil
	}
	return items, r.db.WithContext(ctx).Where("target_id IN ?", ids).Find(&items).Error
}

func (r *visibilityConditionRepository) FindByTriggerQuestion(ctx context.Context, questionID string) ([]models.VisibilityCondition, error) {
	var items []models.VisibilityCondition
	return items, r.db.WithContext(ctx).Where("trigger_question_id = ?", questionID).Find(&items).Error
}
