package repository

import (
	"bitsflow/internal/models"
	"context"

	"gorm.io/gorm"
)

// CaseTimelineEventRepository gestiona la persistencia de eventos del timeline de un caso.
// Es append-only: solo se crean registros, nunca se modifican ni eliminan.
type CaseTimelineEventRepository interface {
	Create(ctx context.Context, event *models.CaseTimelineEvent) error
}

type caseTimelineEventRepository struct {
	db *gorm.DB
}

func NewCaseTimelineEventRepository(db *gorm.DB) CaseTimelineEventRepository {
	return &caseTimelineEventRepository{db: db}
}

func (r *caseTimelineEventRepository) Create(ctx context.Context, event *models.CaseTimelineEvent) error {
	return r.db.WithContext(ctx).Create(event).Error
}
