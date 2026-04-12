package repository

import (
	"bitsflow/internal/models"
	"context"

	"gorm.io/gorm"
)

type OptionRepository interface {
	Repository[models.Option]
	FindByQuestionID(ctx context.Context, questionID string) ([]models.Option, error)
	FindByQuestionIDs(ctx context.Context, questionIDs []string) ([]models.Option, error)
}

type optionRepository struct {
	repository[models.Option]
	db *gorm.DB
}

func NewOptionRepository(db *gorm.DB) OptionRepository {
	return &optionRepository{
		repository: repository[models.Option]{db: db},
		db:         db,
	}
}

func (r *optionRepository) FindByQuestionID(ctx context.Context, questionID string) ([]models.Option, error) {
	var items []models.Option
	return items, r.db.WithContext(ctx).
		Where("question_id = ?", questionID).
		Order(`"order" ASC`).
		Find(&items).Error
}

func (r *optionRepository) FindByQuestionIDs(ctx context.Context, questionIDs []string) ([]models.Option, error) {
	var items []models.Option
	if len(questionIDs) == 0 {
		return items, nil
	}
	return items, r.db.WithContext(ctx).
		Where("question_id IN ?", questionIDs).
		Order(`"order" ASC`).
		Find(&items).Error
}
