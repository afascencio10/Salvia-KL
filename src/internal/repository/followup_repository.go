package repository

import (
	"bitsflow/internal/models"
	"context"
	"fmt"

	"gorm.io/gorm"
)

// FollowUpRepository extiende el CRUD genérico con consultas específicas de follow_up_v2.
// Create, FindByID, Delete, Update, UpdateFields y FindWithPagination
// son heredados de base_repository — NO se reimplementan aquí.
type FollowUpRepository interface {
	Repository[models.FollowUpV2]

	// Existente — conservado
	FindByCaseID(ctx context.Context, caseID string) ([]models.FollowUpV2, error)

	// Nuevos HU-027
	FindByCaseIDOrdered(ctx context.Context, caseID string) ([]models.FollowUpV2, error)
	FindPendingByCaseID(ctx context.Context, caseID string) ([]models.FollowUpV2, error)
	FindCompletedByCaseID(ctx context.Context, caseID string) ([]models.FollowUpV2, error)
	BulkCreate(ctx context.Context, tx *gorm.DB, followUps []models.FollowUpV2) error
	SoftDeleteAndReprogramPending(ctx context.Context, tx *gorm.DB, caseID string) error
	RunInTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error
}

type followUpRepository struct {
	repository[models.FollowUpV2]
	db *gorm.DB
}

func NewFollowUpRepository(db *gorm.DB) FollowUpRepository {
	return &followUpRepository{
		repository: repository[models.FollowUpV2]{db: db},
		db:         db,
	}
}

// FindByCaseID retorna los seguimientos de un caso (sin orden garantizado).
// Conservado del patrón original.
func (r *followUpRepository) FindByCaseID(ctx context.Context, caseID string) ([]models.FollowUpV2, error) {
	var items []models.FollowUpV2
	err := r.db.WithContext(ctx).
		Where("case_id = ?", caseID).
		Find(&items).Error
	return items, err
}

// FindByCaseIDOrdered retorna todos los seguimientos del caso ordenados por scheduled_date ASC.
func (r *followUpRepository) FindByCaseIDOrdered(ctx context.Context, caseID string) ([]models.FollowUpV2, error) {
	var items []models.FollowUpV2
	err := r.db.WithContext(ctx).
		Where("case_id = ?", caseID).
		Order("scheduled_date ASC").
		Find(&items).Error
	return items, err
}

// FindPendingByCaseID retorna los seguimientos con status = 'PENDIENTE'.
func (r *followUpRepository) FindPendingByCaseID(ctx context.Context, caseID string) ([]models.FollowUpV2, error) {
	var items []models.FollowUpV2
	err := r.db.WithContext(ctx).
		Where("case_id = ? AND status = ?", caseID, models.FollowUpStatusPendiente).
		Order("scheduled_date ASC").
		Find(&items).Error
	return items, err
}

// FindCompletedByCaseID retorna los seguimientos con status REALIZADO o VENCIDO.
func (r *followUpRepository) FindCompletedByCaseID(ctx context.Context, caseID string) ([]models.FollowUpV2, error) {
	var items []models.FollowUpV2
	err := r.db.WithContext(ctx).
		Where("case_id = ? AND status IN ?", caseID,
			[]string{models.FollowUpStatusRealizado, models.FollowUpStatusVencido}).
		Order("scheduled_date ASC").
		Find(&items).Error
	return items, err
}

// BulkCreate inserta múltiples seguimientos en una sola operación usando la tx provista.
func (r *followUpRepository) BulkCreate(ctx context.Context, tx *gorm.DB, followUps []models.FollowUpV2) error {
	return tx.WithContext(ctx).Create(&followUps).Error
}

// SoftDeleteAndReprogramPending marca los PENDIENTES como REPROGRAMADO y aplica soft-delete.
// Sigue el patrón de submission_repository.go: primero Updates, luego Delete sobre el mismo Where.
func (r *followUpRepository) SoftDeleteAndReprogramPending(ctx context.Context, tx *gorm.DB, caseID string) error {
	return tx.WithContext(ctx).
		Where("case_id = ? AND status = ?", caseID, models.FollowUpStatusPendiente).
		Updates(map[string]interface{}{"status": models.FollowUpStatusReprogramado}).
		Delete(&models.FollowUpV2{}).Error
}

// RunInTransaction ejecuta fn dentro de una transacción GORM.
// Sigue exactamente el patrón de submission_repository.go.
func (r *followUpRepository) RunInTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("followup: no se pudo iniciar la transacción: %w", tx.Error)
	}

	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
