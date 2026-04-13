package repository

import (
	"bitsflow/internal/models"
	"context"

	"gorm.io/gorm"
)

// FollowUpAttemptRepository define el acceso a datos para los intentos de seguimiento
type FollowUpAttemptRepository interface {
	CreateAttempt(ctx context.Context, attempt *models.FollowUpAttempt) error
	GetByFollowUpID(ctx context.Context, followUpID string) ([]models.FollowUpAttempt, error)
	CountByFollowUpID(ctx context.Context, followUpID string) (int64, error)
}

type followUpAttemptRepository struct {
	db *gorm.DB
}

// NewFollowUpAttemptRepository crea la instancia del repositorio
func NewFollowUpAttemptRepository(db *gorm.DB) FollowUpAttemptRepository {
	return &followUpAttemptRepository{db: db}
}

// CreateAttempt crea un nuevo registro de intento
func (r *followUpAttemptRepository) CreateAttempt(ctx context.Context, attempt *models.FollowUpAttempt) error {
	return r.db.WithContext(ctx).Create(attempt).Error
}

// GetByFollowUpID obtiene todos los intentos de un seguimiento específico ordenados por fecha
func (r *followUpAttemptRepository) GetByFollowUpID(ctx context.Context, followUpID string) ([]models.FollowUpAttempt, error) {
	var attempts []models.FollowUpAttempt
	err := r.db.WithContext(ctx).
		Where("follow_up_id = ?", followUpID).
		Order("created_at DESC").
		Find(&attempts).Error
	return attempts, err
}

// CountByFollowUpID cuenta el número total de intentos para un seguimiento específico
func (r *followUpAttemptRepository) CountByFollowUpID(ctx context.Context, followUpID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.FollowUpAttempt{}).
		Where("follow_up_id = ?", followUpID).
		Count(&count).Error
	return count, err
}
