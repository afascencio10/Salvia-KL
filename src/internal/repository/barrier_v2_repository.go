package repository

import (
	"bitsflow/internal/models"
	"context"

	"gorm.io/gorm"
)

type BarrierV2Repository interface {
	Repository[models.BarrierV2]
	FindByCaseID(ctx context.Context, caseID string) ([]models.BarrierV2, error)
	FindByFollowUpID(ctx context.Context, followUpID string) ([]models.BarrierV2, error)
	FindActiveByCaseID(ctx context.Context, caseID string) ([]models.BarrierV2, error)
	FindByIDs(ctx context.Context, ids []string) ([]models.BarrierV2, error)
	UpdateStatus(ctx context.Context, id string, status string) error
	// FindByCreatedByIDWithRelations devuelve barreras del agente enriquecidas
	// con datos de victim_case (nombres, documento, i_code).
	FindByCreatedByIDWithRelations(ctx context.Context, createdByID string) ([]models.BarrierV2WithRelations, error)
}

type barrierV2Repository struct {
	repository[models.BarrierV2]
	db *gorm.DB
}

func NewBarrierV2Repository(db *gorm.DB) BarrierV2Repository {
	return &barrierV2Repository{
		repository: repository[models.BarrierV2]{db: db},
		db:         db,
	}
}

func (r *barrierV2Repository) FindByCaseID(ctx context.Context, caseID string) ([]models.BarrierV2, error) {
	var items []models.BarrierV2
	return items, r.db.WithContext(ctx).Where("case_id = ?", caseID).Find(&items).Error
}

func (r *barrierV2Repository) FindByFollowUpID(ctx context.Context, followUpID string) ([]models.BarrierV2, error) {
	var items []models.BarrierV2
	return items, r.db.WithContext(ctx).Where("follow_up_id = ?", followUpID).Find(&items).Error
}

func (r *barrierV2Repository) FindActiveByCaseID(ctx context.Context, caseID string) ([]models.BarrierV2, error) {
	var items []models.BarrierV2
	return items, r.db.WithContext(ctx).
		Where("case_id = ? AND status != ? AND deleted_at IS NULL", caseID, models.BarrierV2StatusManaged).
		Find(&items).Error
}

func (r *barrierV2Repository) FindByIDs(ctx context.Context, ids []string) ([]models.BarrierV2, error) {
	var items []models.BarrierV2
	if len(ids) == 0 {
		return items, nil
	}
	return items, r.db.WithContext(ctx).Where("id IN ?", ids).Find(&items).Error
}

func (r *barrierV2Repository) UpdateStatus(ctx context.Context, id string, status string) error {
	return r.db.WithContext(ctx).
		Model(&models.BarrierV2{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// barrierWithRelationsSQL es la query base que enriquece barrier_v2 con datos
// de victim_case mediante LEFT JOIN.
//
//	barrier_v2.case_id → victim_case.victim_case_i_code
const barrierWithRelationsSQL = `
SELECT
    b.id,
    b.case_id,
    b.follow_up_id,
    b.sector,
    b.description,
    b.status,
    b.created_by_id,
    b.created_at,
    COALESCE(vc.victim_case_victim_names, '')       AS victim_name,
    COALESCE(vc.victim_case_victim_last_names, '')  AS victim_last_name,
    COALESCE(vc.victim_case_victim_doc_number, '')  AS victim_doc_number,
    COALESCE(vc.victim_case_i_code, '')             AS case_code
FROM salvia.barrier_v2 b
LEFT JOIN salvia.victim_case vc
    ON vc.victim_case_i_code = b.case_id
WHERE b.deleted_at IS NULL
`

func (r *barrierV2Repository) FindByCreatedByIDWithRelations(ctx context.Context, createdByID string) ([]models.BarrierV2WithRelations, error) {
	var items []models.BarrierV2WithRelations
	err := r.db.WithContext(ctx).
		Raw(barrierWithRelationsSQL+" AND b.created_by_id = ? ORDER BY b.created_at DESC", createdByID).
		Scan(&items).Error
	return items, err
}
