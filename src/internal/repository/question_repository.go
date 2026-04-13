package repository

import (
	"bitsflow/internal/models"
	"context"

	"gorm.io/gorm"
)

type QuestionRepository interface {
	Repository[models.Question]
	FindByFormID(ctx context.Context, formID string) ([]models.Question, error)
	FindBySectionID(ctx context.Context, sectionID string) ([]models.Question, error)
}

type questionRepository struct {
	repository[models.Question]
	db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) QuestionRepository {
	return &questionRepository{
		repository: repository[models.Question]{db: db},
		db:         db,
	}
}

func (r *questionRepository) FindByFormID(ctx context.Context, formID string) ([]models.Question, error) {
	var items []models.Question
	return items, r.db.WithContext(ctx).Where("form_id = ?", formID).Order(`"order" ASC`).Find(&items).Error
}

func (r *questionRepository) FindBySectionID(ctx context.Context, sectionID string) ([]models.Question, error) {
	var items []models.Question
	return items, r.db.WithContext(ctx).Where("form_section_id = ?", sectionID).Order(`"order" ASC`).Find(&items).Error
}
