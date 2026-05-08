package repository

import (
	"bitsflow/internal/models"
	"context"

	"gorm.io/gorm"
)

// EntityLetterRepository extiende el contrato CRUD genérico con consultas
// específicas del dominio de oficios (EntityLetter).
type EntityLetterRepository interface {
	Repository[models.EntityLetter]

	// FindByCaseID devuelve todos los oficios asociados a un caso, ordenados por fecha de creación.
	FindByCaseID(ctx context.Context, caseID string) ([]models.EntityLetter, error)

	// FindByBarrierID devuelve todos los oficios asociados a una barrera.
	FindByBarrierID(ctx context.Context, barrierID string) ([]models.EntityLetter, error)

	// FindByState devuelve todos los oficios en un estado específico, paginados.
	FindByState(ctx context.Context, state string, page, pageSize int) (PageResult[models.EntityLetter], error)

	// FindByAgentID devuelve los oficios asignados a un agente de seguimiento.
	FindByAgentID(ctx context.Context, agentID string) ([]models.EntityLetter, error)

	// FindByNotificationUserID devuelve los oficios asignados a un agente de notificaciones.
	FindByNotificationUserID(ctx context.Context, notificationUserID string) ([]models.EntityLetter, error)

	// UpdateState actualiza únicamente el campo state del oficio.
	UpdateState(ctx context.Context, id, state string) error
}

type entityLetterRepository struct {
	repository[models.EntityLetter]
	db *gorm.DB
}

func NewEntityLetterRepository(db *gorm.DB) EntityLetterRepository {
	return &entityLetterRepository{
		repository: repository[models.EntityLetter]{db: db},
		db:         db,
	}
}

func (r *entityLetterRepository) FindByCaseID(ctx context.Context, caseID string) ([]models.EntityLetter, error) {
	var items []models.EntityLetter
	err := r.db.WithContext(ctx).
		Where("case_id = ?", caseID).
		Order("created_at ASC").
		Find(&items).Error
	return items, err
}

func (r *entityLetterRepository) FindByBarrierID(ctx context.Context, barrierID string) ([]models.EntityLetter, error) {
	var items []models.EntityLetter
	err := r.db.WithContext(ctx).
		Where("barrier_id = ?", barrierID).
		Order("created_at ASC").
		Find(&items).Error
	return items, err
}

func (r *entityLetterRepository) FindByState(ctx context.Context, state string, page, pageSize int) (PageResult[models.EntityLetter], error) {
	if pageSize <= 0 {
		pageSize = 20
	}
	offset := page * pageSize

	var total int64
	var items []models.EntityLetter

	base := r.db.WithContext(ctx).Model(&models.EntityLetter{}).Where("state = ?", state)

	if err := base.Count(&total).Error; err != nil {
		return PageResult[models.EntityLetter]{}, err
	}
	if err := base.Order("created_at ASC").Offset(offset).Limit(pageSize).Find(&items).Error; err != nil {
		return PageResult[models.EntityLetter]{}, err
	}

	return PageResult[models.EntityLetter]{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (r *entityLetterRepository) FindByAgentID(ctx context.Context, agentID string) ([]models.EntityLetter, error) {
	var items []models.EntityLetter
	err := r.db.WithContext(ctx).
		Where("agent_id = ?", agentID).
		Order("created_at ASC").
		Find(&items).Error
	return items, err
}

func (r *entityLetterRepository) FindByNotificationUserID(ctx context.Context, notificationUserID string) ([]models.EntityLetter, error) {
	var items []models.EntityLetter
	err := r.db.WithContext(ctx).
		Where("notification_user_id = ?", notificationUserID).
		Order("created_at ASC").
		Find(&items).Error
	return items, err
}

func (r *entityLetterRepository) UpdateState(ctx context.Context, id, state string) error {
	result := r.db.WithContext(ctx).
		Model(&models.EntityLetter{}).
		Where("id = ?", id).
		Update("state", state)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
