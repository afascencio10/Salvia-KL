package repository

import (
	"bitsflow/internal/models"
	"context"

	"gorm.io/gorm"
)

// ─── RepeaterGroup ────────────────────────────────────────────────────────────

type RepeaterGroupRepository interface {
	Repository[models.RepeaterGroup]
	FindByFormSectionID(ctx context.Context, formSectionID string) ([]models.RepeaterGroup, error)
}

type repeaterGroupRepository struct {
	repository[models.RepeaterGroup]
	db *gorm.DB
}

func NewRepeaterGroupRepository(db *gorm.DB) RepeaterGroupRepository {
	return &repeaterGroupRepository{repository: repository[models.RepeaterGroup]{db: db}, db: db}
}

func (r *repeaterGroupRepository) FindByFormSectionID(ctx context.Context, formSectionID string) ([]models.RepeaterGroup, error) {
	var items []models.RepeaterGroup
	err := r.db.WithContext(ctx).Where("form_section_id = ?", formSectionID).Order(`"order" ASC`).Find(&items).Error
	return items, err
}

// ─── Question ─────────────────────────────────────────────────────────────────

type QuestionRepository interface {
	Repository[models.Question]
	FindByFormID(ctx context.Context, formID string) ([]models.Question, error)
	FindByFormSectionID(ctx context.Context, formSectionID string) ([]models.Question, error)
}

type questionRepository struct {
	repository[models.Question]
	db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) QuestionRepository {
	return &questionRepository{repository: repository[models.Question]{db: db}, db: db}
}

func (r *questionRepository) FindByFormID(ctx context.Context, formID string) ([]models.Question, error) {
	var items []models.Question
	err := r.db.WithContext(ctx).Where("form_id = ?", formID).Order(`"order" ASC`).Find(&items).Error
	return items, err
}

func (r *questionRepository) FindByFormSectionID(ctx context.Context, formSectionID string) ([]models.Question, error) {
	var items []models.Question
	err := r.db.WithContext(ctx).Where("form_section_id = ?", formSectionID).Order(`"order" ASC`).Find(&items).Error
	return items, err
}

// ─── VisibilityCondition ──────────────────────────────────────────────────────

type VisibilityConditionRepository interface {
	Repository[models.VisibilityCondition]
	FindByTargetID(ctx context.Context, targetID string) ([]models.VisibilityCondition, error)
}

type visibilityConditionRepository struct {
	repository[models.VisibilityCondition]
	db *gorm.DB
}

func NewVisibilityConditionRepository(db *gorm.DB) VisibilityConditionRepository {
	return &visibilityConditionRepository{repository: repository[models.VisibilityCondition]{db: db}, db: db}
}

func (r *visibilityConditionRepository) FindByTargetID(ctx context.Context, targetID string) ([]models.VisibilityCondition, error) {
	var items []models.VisibilityCondition
	err := r.db.WithContext(ctx).Where("target_id = ?", targetID).Find(&items).Error
	return items, err
}
