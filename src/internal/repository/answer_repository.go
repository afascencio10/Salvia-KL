package repository

import (
	"bitsflow/internal/models"
	"context"

	"gorm.io/gorm"
)

type AnswerRepository interface {
	Repository[models.Answer]
	FindBySubmissionID(ctx context.Context, submissionID string) ([]models.Answer, error)
	FindDirectBySubmissionID(ctx context.Context, submissionID string) ([]models.Answer, error)
	FindByRepeaterEntryID(ctx context.Context, entryID string) ([]models.Answer, error)
	FindByRepeaterEntryIDs(ctx context.Context, entryIDs []string) ([]models.Answer, error)
	FindByQuestionID(ctx context.Context, questionID string) ([]models.Answer, error)
}

type answerRepository struct {
	repository[models.Answer]
	db *gorm.DB
}

func NewAnswerRepository(db *gorm.DB) AnswerRepository {
	return &answerRepository{
		repository: repository[models.Answer]{db: db},
		db:         db,
	}
}

func (r *answerRepository) FindBySubmissionID(ctx context.Context, submissionID string) ([]models.Answer, error) {
	var items []models.Answer
	return items, r.db.WithContext(ctx).Where("form_submission_id = ?", submissionID).Find(&items).Error
}

func (r *answerRepository) FindDirectBySubmissionID(ctx context.Context, submissionID string) ([]models.Answer, error) {
	var items []models.Answer
	return items, r.db.WithContext(ctx).
		Where("form_submission_id = ? AND repeater_entry_id IS NULL", submissionID).
		Find(&items).Error
}

func (r *answerRepository) FindByRepeaterEntryID(ctx context.Context, entryID string) ([]models.Answer, error) {
	var items []models.Answer
	return items, r.db.WithContext(ctx).
		Where("repeater_entry_id = ?", entryID).
		Find(&items).Error
}

func (r *answerRepository) FindByRepeaterEntryIDs(ctx context.Context, entryIDs []string) ([]models.Answer, error) {
	var items []models.Answer
	if len(entryIDs) == 0 {
		return items, nil
	}
	return items, r.db.WithContext(ctx).
		Where("repeater_entry_id IN ?", entryIDs).
		Find(&items).Error
}

func (r *answerRepository) FindByQuestionID(ctx context.Context, questionID string) ([]models.Answer, error) {
	var items []models.Answer
	return items, r.db.WithContext(ctx).Where("question_id = ?", questionID).Find(&items).Error
}
