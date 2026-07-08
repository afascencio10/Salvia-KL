package repository

import (
	"bitsflow/internal/models"
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// EntityLetterListFilter agrupa los criterios de consulta paginada con relaciones.
type EntityLetterListFilter struct {
	AgentID            string
	NotificationUserID string // legacy v1
	ListAll            bool   // rol an — tab Todos / gestionar global
	MineOnly           bool   // rol an — tab Mis Oficios
	NotificationAgentID string
	State              string
	Identidad          string
	Entidad            string
	NumeroRadicado     string
	ManageableOnly     bool
	ManageableStates   []string
	NotificationUserIDReview   string
	NotificationUserIDRadicado string
	NotificationUserIDResponse string
}

// EntityLetterListResult extiende PageResult con el conteo de oficios pendientes de gestionar.
type EntityLetterListResult struct {
	PageResult[models.EntityLetterWithRelations]
	PendingCount int64
}

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

	// FindWithRelationsFilteredPaginated devuelve oficios enriquecidos con paginación y filtros en BD.
	FindWithRelationsFilteredPaginated(ctx context.Context, filter EntityLetterListFilter, page, pageSize int) (EntityLetterListResult, error)

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
    el.priority,
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
    el.notification_user_id_review,
    el.notification_user_id_radicado,
    el.notification_user_id_response,
    el.created_at,
    el.updated_at,
    el.entity_branch_id,
    el.department_id,
    el.city_id,
    el.town_id,
    el.official_dependency,
    el.subject,
    COALESCE(b.sector, '')                         AS barrier_sector,
    COALESCE(b.description, '')                    AS barrier_description,
    COALESCE(vc.victim_case_victim_names, '')       AS victim_name,
    COALESCE(vc.victim_case_victim_last_names, '')  AS victim_last_name,
    COALESCE(vc.victim_case_victim_doc_number, '')  AS victim_doc_number,
    COALESCE(vc.victim_case_i_code, '')             AS case_code,
    COALESCE(t.town_name, '')                       AS town_name
FROM salvia.entity_letter el
LEFT JOIN salvia.barrier_v2 b
    ON b.id::text = el.barrier_id
LEFT JOIN salvia.victim_case vc
    ON vc.victim_case_i_code = el.case_id
LEFT JOIN security.town t
    ON t.town_id::varchar = el.town_id
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

// buildRelationsWhere construye la cláusula WHERE adicional y sus argumentos para consultas enriquecidas.
func buildRelationsWhere(filter EntityLetterListFilter) (string, []interface{}) {
	var conds []string
	var args []interface{}

	if filter.AgentID != "" {
		conds = append(conds, "el.agent_id = ?")
		args = append(args, filter.AgentID)
	} else if filter.MineOnly {
		if filter.NotificationAgentID != "" {
			conds = append(conds, `(
				el.notification_user_id_review = ? OR
				el.notification_user_id_radicado = ? OR
				el.notification_user_id_response = ?
			)`)
			args = append(args, filter.NotificationAgentID, filter.NotificationAgentID, filter.NotificationAgentID)
		}
	} else if filter.NotificationUserID != "" {
		conds = append(conds, "el.notification_user_id = ?")
		args = append(args, filter.NotificationUserID)
	} else if filter.ListAll {
		// sin filtro por usuario — todos los oficios
	}

	if filter.ManageableOnly && len(filter.ManageableStates) > 0 {
		placeholders := strings.Repeat("?,", len(filter.ManageableStates))
		placeholders = placeholders[:len(placeholders)-1]
		conds = append(conds, fmt.Sprintf("el.state IN (%s)", placeholders))
		for _, s := range filter.ManageableStates {
			args = append(args, s)
		}
	} else if filter.State != "" {
		conds = append(conds, "el.state = ?")
		args = append(args, filter.State)
	}

	if filter.Identidad != "" {
		conds = append(conds, "vc.victim_case_victim_doc_number ILIKE ?")
		args = append(args, "%"+filter.Identidad+"%")
	}
	if filter.Entidad != "" {
		conds = append(conds, "b.sector ILIKE ?")
		args = append(args, "%"+filter.Entidad+"%")
	}
	if filter.NumeroRadicado != "" {
		conds = append(conds, "el.numero_radicado ILIKE ?")
		args = append(args, "%"+filter.NumeroRadicado+"%")
	}
	if filter.NotificationUserIDReview != "" {
		conds = append(conds, "el.notification_user_id_review = ?")
		args = append(args, filter.NotificationUserIDReview)
	}
	if filter.NotificationUserIDRadicado != "" {
		conds = append(conds, "el.notification_user_id_radicado = ?")
		args = append(args, filter.NotificationUserIDRadicado)
	}
	if filter.NotificationUserIDResponse != "" {
		conds = append(conds, "el.notification_user_id_response = ?")
		args = append(args, filter.NotificationUserIDResponse)
	}

	if len(conds) == 0 {
		return "", args
	}
	return " AND " + strings.Join(conds, " AND "), args
}

const withRelationsCountSQL = `
SELECT COUNT(*)
FROM salvia.entity_letter el
LEFT JOIN salvia.barrier_v2 b
    ON b.id::text = el.barrier_id
LEFT JOIN salvia.victim_case vc
    ON vc.victim_case_i_code = el.case_id
WHERE el.deleted_at IS NULL
`

func (r *entityLetterRepository) FindWithRelationsFilteredPaginated(
	ctx context.Context,
	filter EntityLetterListFilter,
	page, pageSize int,
) (EntityLetterListResult, error) {
	if pageSize <= 0 {
		pageSize = 20
	}
	offset := page * pageSize

	whereExtra, args := buildRelationsWhere(filter)
	if whereExtra == "" && !filter.ListAll {
		return EntityLetterListResult{}, fmt.Errorf("entity_letter: se requiere agentId, listAll o mineOnly")
	}

	var total int64
	if err := r.db.WithContext(ctx).
		Raw(withRelationsCountSQL+whereExtra, args...).
		Scan(&total).Error; err != nil {
		return EntityLetterListResult{}, err
	}

	dataSQL := withRelationsSQL + whereExtra + " ORDER BY el.created_at DESC LIMIT ? OFFSET ?"
	dataArgs := append(append([]interface{}{}, args...), pageSize, offset)

	var items []models.EntityLetterWithRelations
	if err := r.db.WithContext(ctx).Raw(dataSQL, dataArgs...).Scan(&items).Error; err != nil {
		return EntityLetterListResult{}, err
	}

	pendingCount, err := r.countManageableWithRelations(ctx, filter)
	if err != nil {
		return EntityLetterListResult{}, err
	}

	return EntityLetterListResult{
		PageResult: PageResult[models.EntityLetterWithRelations]{
			Items:    items,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
		PendingCount: pendingCount,
	}, nil
}

func (r *entityLetterRepository) countManageableWithRelations(ctx context.Context, filter EntityLetterListFilter) (int64, error) {
	if len(filter.ManageableStates) == 0 {
		return 0, nil
	}

	pendingFilter := EntityLetterListFilter{
		ManageableOnly:   true,
		ManageableStates: filter.ManageableStates,
	}
	if filter.AgentID != "" {
		pendingFilter.AgentID = filter.AgentID
	} else if filter.MineOnly && filter.NotificationAgentID != "" {
		pendingFilter.MineOnly = true
		pendingFilter.NotificationAgentID = filter.NotificationAgentID
	} else if filter.NotificationUserID != "" {
		pendingFilter.NotificationUserID = filter.NotificationUserID
	} else {
		pendingFilter.ListAll = true
	}
	whereExtra, args := buildRelationsWhere(pendingFilter)
	if whereExtra == "" {
		return 0, nil
	}

	var count int64
	err := r.db.WithContext(ctx).
		Raw(withRelationsCountSQL+whereExtra, args...).
		Scan(&count).Error
	return count, err
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
