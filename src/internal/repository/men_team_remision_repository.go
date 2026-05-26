package repository

import (
	"bitsflow/internal/models"
	"context"

	"gorm.io/gorm"
)

type MenTeamRemisionRepository interface {
	Repository[models.MenTeamRemision]
	FindByCaseID(ctx context.Context, caseID string) ([]models.MenTeamRemision, error)
	FindByFollowUpID(ctx context.Context, followUpID string) ([]models.MenTeamRemision, error)
}

type menTeamRemisionRepository struct {
	repository[models.MenTeamRemision]
	db *gorm.DB
}

func NewMenTeamRemisionRepository(db *gorm.DB) MenTeamRemisionRepository {
	return &menTeamRemisionRepository{
		repository: repository[models.MenTeamRemision]{db: db},
		db:         db,
	}
}

func (r *menTeamRemisionRepository) FindByCaseID(ctx context.Context, caseID string) ([]models.MenTeamRemision, error) {
	var items []models.MenTeamRemision
	return items, r.db.WithContext(ctx).Where("case_id = ?", caseID).Find(&items).Error
}

func (r *menTeamRemisionRepository) FindByFollowUpID(ctx context.Context, followUpID string) ([]models.MenTeamRemision, error) {
	var items []models.MenTeamRemision
	return items, r.db.WithContext(ctx).Where("follow_up_id = ?", followUpID).Find(&items).Error
}
