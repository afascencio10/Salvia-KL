package repository

import (
	"bitsflow/internal/models"
	"context"

	"gorm.io/gorm"
)

// ContactAttemptCounters agrupa los contadores que gobiernan el flujo 3x3.
type ContactAttemptCounters struct {
	DailyFailedCount  int64
	TotalCount        int64
	DistinctDaysCount int64
}

// ContactAttemptRepository define el acceso a datos de contact_attempts.
type ContactAttemptRepository interface {
	Create(ctx context.Context, attempt *models.ContactAttempt) error
	GetByPsicosocialID(ctx context.Context, psicosocialID string) ([]models.ContactAttempt, error)
	GetByID(ctx context.Context, id string) (*models.ContactAttempt, error)
	UpdateConsent(ctx context.Context, id string, consentGiven bool) error
	Counters(ctx context.Context, psicosocialID string) (ContactAttemptCounters, error)
}

type contactAttemptRepository struct {
	db *gorm.DB
}

// NewContactAttemptRepository crea la instancia del repositorio.
func NewContactAttemptRepository(db *gorm.DB) ContactAttemptRepository {
	return &contactAttemptRepository{db: db}
}

func (r *contactAttemptRepository) Create(ctx context.Context, attempt *models.ContactAttempt) error {
	return r.db.WithContext(ctx).Create(attempt).Error
}

// GetByPsicosocialID devuelve los intentos ordenados ascendente por attempt_at.
func (r *contactAttemptRepository) GetByPsicosocialID(ctx context.Context, psicosocialID string) ([]models.ContactAttempt, error) {
	var attempts []models.ContactAttempt
	err := r.db.WithContext(ctx).
		Where("psicosocial_id = ?", psicosocialID).
		Order("attempt_at ASC").
		Find(&attempts).Error
	return attempts, err
}

func (r *contactAttemptRepository) GetByID(ctx context.Context, id string) (*models.ContactAttempt, error) {
	var attempt models.ContactAttempt
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&attempt).Error
	if err != nil {
		return nil, err
	}
	return &attempt, nil
}

func (r *contactAttemptRepository) UpdateConsent(ctx context.Context, id string, consentGiven bool) error {
	return r.db.WithContext(ctx).
		Model(&models.ContactAttempt{}).
		Where("id = ?", id).
		Update("consent_given", consentGiven).Error
}

// Counters calcula intentos fallidos del día, total y días distintos en una sola pasada.
func (r *contactAttemptRepository) Counters(ctx context.Context, psicosocialID string) (ContactAttemptCounters, error) {
	var c ContactAttemptCounters
	row := struct {
		DailyFailedCount  int64
		TotalCount        int64
		DistinctDaysCount int64
	}{}
	err := r.db.WithContext(ctx).
		Model(&models.ContactAttempt{}).
		Select(`
			COUNT(*) FILTER (WHERE was_answered = false AND attempt_at::date = CURRENT_DATE) AS daily_failed_count,
			COUNT(*) AS total_count,
			COUNT(DISTINCT attempt_at::date) AS distinct_days_count`).
		Where("psicosocial_id = ?", psicosocialID).
		Scan(&row).Error
	if err != nil {
		return c, err
	}
	c.DailyFailedCount = row.DailyFailedCount
	c.TotalCount = row.TotalCount
	c.DistinctDaysCount = row.DistinctDaysCount
	return c, nil
}
