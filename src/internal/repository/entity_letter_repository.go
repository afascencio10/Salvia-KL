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

	// FindByAgentIDWithRelations devuelve los oficios del agente enriquecidos con datos de
	// barrier_v2 (sector, description) y victim_case (nombres, doc, i_code).
	FindByAgentIDWithRelations(ctx context.Context, agentID string) ([]models.EntityLetterWithRelations, error)

	// FindByNotificationUserIDWithRelations devuelve los oficios del agente de notificaciones
	// enriquecidos con datos de barrier_v2 y victim_case.
	FindByNotificationUserIDWithRelations(ctx context.Context, notifUserID string) ([]models.EntityLetterWithRelations, error)

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

// withRelationsSQL es la query base que enriquece entity_letter con datos de
// barrier_v2 y victim_case mediante LEFT JOIN.
//
//	entity_letter.barrier_id  → barrier_v2.id (UUID almacenado como varchar)
//	entity_letter.case_id     → victim_case.victim_case_i_code
const withRelationsSQL = `
SELECT
    el.id,
    el.barrier_id,
    el.case_id,
    el.state,
    el.agent_id,
    el.notification_user_id,
    el.review_by,
    el.radicado_by,
    el.register_by,
    el.entidad,
    el.nivel,
    el.url_kofax,
    el.asunto_radicado,
    el.correo_entidad,
    el.numero_radicado,
    el.response_date,
    el.correo_remitente,
    el.asunto_respuesta,
    el.response_review_by,
    el.reason_correction,
    el.created_at,
    el.updated_at,
    COALESCE(b.sector, '')                         AS barrier_sector,
    COALESCE(b.description, '')                    AS barrier_description,
    COALESCE(vc.victim_case_victim_names, '')       AS victim_name,
    COALESCE(vc.victim_case_victim_last_names, '')  AS victim_last_name,
    COALESCE(vc.victim_case_victim_doc_number, '')  AS victim_doc_number,
    COALESCE(vc.victim_case_i_code, '')             AS case_code
FROM salvia.entity_letter el
LEFT JOIN salvia.barrier_v2 b
    ON b.id::text = el.barrier_id
LEFT JOIN salvia.victim_case vc
    ON vc.victim_case_i_code = el.case_id
WHERE el.deleted_at IS NULL
`

func (r *entityLetterRepository) FindByAgentIDWithRelations(ctx context.Context, agentID string) ([]models.EntityLetterWithRelations, error) {
	var items []models.EntityLetterWithRelations
	err := r.db.WithContext(ctx).
		Raw(withRelationsSQL+" AND el.agent_id = ? ORDER BY el.created_at DESC", agentID).
		Scan(&items).Error
	return items, err
}

func (r *entityLetterRepository) FindByNotificationUserIDWithRelations(ctx context.Context, notifUserID string) ([]models.EntityLetterWithRelations, error) {
	var items []models.EntityLetterWithRelations
	err := r.db.WithContext(ctx).
		Raw(withRelationsSQL+" AND el.notification_user_id = ? ORDER BY el.created_at DESC", notifUserID).
		Scan(&items).Error
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
