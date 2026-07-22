package repository

import (
	"bitsflow/internal/models"
	"context"
	"errors"

	"gorm.io/gorm"
)

// CaseTaskRepository extiende el contrato CRUD genérico con consultas
// específicas del dominio de tareas de caso (CaseTask).
type CaseTaskRepository interface {
	Repository[models.CaseTask]

	// FindByAssignedUserID devuelve todas las tareas asignadas a un usuario.
	FindByAssignedUserID(ctx context.Context, assignedUserID string) ([]models.CaseTask, error)

	// FindByAssignedUserIDWithRelations devuelve tareas del usuario que tienen
	// barrierId asignado, enriquecidas con datos de barrier_v2 y victim_case.
	// Es la consulta principal para la pantalla "Mis Barreras".
	FindByAssignedUserIDWithRelations(ctx context.Context, assignedUserID string) ([]models.CaseTaskWithRelations, error)

	// FindByCaseID devuelve todas las tareas de un caso.
	FindByCaseID(ctx context.Context, caseID string) ([]models.CaseTask, error)

	// FindByCaseIDAndBarrier devuelve las tareas de un caso filtradas por barrera.
	FindByCaseIDAndBarrier(ctx context.Context, caseID, barrierID string) ([]models.CaseTask, error)

	// FindByCaseIDAndPsychosocial devuelve las tareas de un caso filtradas por remisión psicosocial.
	FindByCaseIDAndPsychosocial(ctx context.Context, caseID, psychosocialID string) ([]models.CaseTask, error)

	// FindByBarrierID devuelve todas las tareas relacionadas a una barrera.
	FindByBarrierID(ctx context.Context, barrierID string) ([]models.CaseTask, error)

	// FindTodoByEntityLetterID devuelve la primera tarea en estado "ToDo"
	// asociada a un EntityLetter dado. Devuelve nil, nil si no existe.
	FindTodoByEntityLetterID(ctx context.Context, entityLetterID string) (*models.CaseTask, error)
}

type caseTaskRepository struct {
	repository[models.CaseTask]
	db *gorm.DB
}

func NewCaseTaskRepository(db *gorm.DB) CaseTaskRepository {
	return &caseTaskRepository{
		repository: repository[models.CaseTask]{db: db},
		db:         db,
	}
}

func (r *caseTaskRepository) FindByAssignedUserID(ctx context.Context, assignedUserID string) ([]models.CaseTask, error) {
	var items []models.CaseTask
	err := r.db.WithContext(ctx).
		Where("assigned_user_id = ?", assignedUserID).
		Order("created_at DESC").
		Find(&items).Error
	return items, err
}

func (r *caseTaskRepository) FindByCaseID(ctx context.Context, caseID string) ([]models.CaseTask, error) {
	var items []models.CaseTask
	err := r.db.WithContext(ctx).
		Where("case_id = ?", caseID).
		Order("created_at DESC").
		Find(&items).Error
	return items, err
}

func (r *caseTaskRepository) FindByCaseIDAndBarrier(ctx context.Context, caseID, barrierID string) ([]models.CaseTask, error) {
	var items []models.CaseTask
	err := r.db.WithContext(ctx).
		Where("case_id = ? AND barrier_id = ?", caseID, barrierID).
		Order("created_at DESC").
		Find(&items).Error
	return items, err
}

func (r *caseTaskRepository) FindByCaseIDAndPsychosocial(ctx context.Context, caseID, psychosocialID string) ([]models.CaseTask, error) {
	var items []models.CaseTask
	err := r.db.WithContext(ctx).
		Where("case_id = ? AND psychosocial_support_id = ?", caseID, psychosocialID).
		Order("created_at DESC").
		Find(&items).Error
	return items, err
}

func (r *caseTaskRepository) FindByBarrierID(ctx context.Context, barrierID string) ([]models.CaseTask, error) {
	var items []models.CaseTask
	err := r.db.WithContext(ctx).
		Where("barrier_id = ?", barrierID).
		Order("created_at DESC").
		Find(&items).Error
	return items, err
}

func (r *caseTaskRepository) FindTodoByEntityLetterID(ctx context.Context, entityLetterID string) (*models.CaseTask, error) {
	var task models.CaseTask
	err := r.db.WithContext(ctx).
		Where("entity_letter_id = ? AND status = ?", entityLetterID, models.CaseTaskStatusToDo).
		First(&task).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &task, nil
}

// caseTaskWithRelationsSQL es la query base que enriquece case_task con datos
// de barrier_v2 y victim_case mediante LEFT JOIN.
//
//	case_task.barrier_id (UUID) → barrier_v2.id (UUID)  — comparación directa sin cast
//	case_task.case_id  (varchar) → victim_case.victim_case_i_code (varchar)
const caseTaskWithRelationsSQL = `
SELECT
    ct.id,
    ct.category,
    ct.type,
    ct.description,
    ct.result,
    ct.completed_at,
    ct.assigned_user_id,
    ct.status,
    ct.case_id,
    ct.barrier_id,
    ct.created_at,
    COALESCE(b.sector, '')                         AS barrier_sector,
    COALESCE(b.description, '')                    AS barrier_description,
    COALESCE(b.status, '')                         AS barrier_status,
    COALESCE(vc.victim_case_victim_names, '')       AS victim_name,
    COALESCE(vc.victim_case_victim_last_names, '')  AS victim_last_name,
    COALESCE(vc.victim_case_victim_doc_number, '')  AS victim_doc_number,
    COALESCE(vc.victim_case_i_code, '')             AS case_code
FROM salvia.case_task ct
LEFT JOIN salvia.barrier_v2 b
    ON b.id = ct.barrier_id
LEFT JOIN salvia.victim_case vc
    ON vc.victim_case_i_code = ct.case_id
WHERE ct.deleted_at IS NULL
  AND ct.barrier_id IS NOT NULL
`

func (r *caseTaskRepository) FindByAssignedUserIDWithRelations(ctx context.Context, assignedUserID string) ([]models.CaseTaskWithRelations, error) {
	var items []models.CaseTaskWithRelations
	err := r.db.WithContext(ctx).
		Raw(caseTaskWithRelationsSQL+" AND ct.assigned_user_id = ? ORDER BY ct.created_at DESC", assignedUserID).
		Scan(&items).Error
	return items, err
}
