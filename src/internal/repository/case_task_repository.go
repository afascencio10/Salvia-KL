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

	// FindByAssignedUserIDPaged devuelve, paginadas, todas las tareas del usuario
	// —no solo las de barrera— con el caso de origen resuelto. Es la consulta de la
	// pantalla "Mis Tareas". status vacío trae todos los estados.
	//
	// role es el rol del usuario en sesión: cuando es 'te' (Equipo Técnico) se
	// incluyen además las tareas comite_caso del pool (assigned_user_id vacío),
	// disponibles para que cualquier integrante del equipo técnico las gestione.
	FindByAssignedUserIDPaged(ctx context.Context, assignedUserID, status, role string, page, pageSize int) (PageResult[models.CaseTaskWithRelations], error)

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

	// FindCommentsByCaseID devuelve los comentarios del Equipo Técnico (case_task
	// Done de tipo comentario_equipo_tecnico) de un caso, con el nombre de quien
	// los creó resuelto, del más antiguo al más reciente.
	FindCommentsByCaseID(ctx context.Context, caseID string) ([]models.CaseTaskCommentItem, error)
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

func (r *caseTaskRepository) FindCommentsByCaseID(ctx context.Context, caseID string) ([]models.CaseTaskCommentItem, error) {
	const q = `
SELECT
    ct.id AS id,
    ct.result AS result,
    ct.completed_at AS completed_at,
    TRIM(COALESCE(gup.general_user_profile_names, '') || ' ' || COALESCE(gup.general_user_profile_last_names, '')) AS assigned_user_name
FROM salvia.case_task ct
LEFT JOIN security.general_user gu ON gu.general_user_i_code = ct.assigned_user_id
LEFT JOIN security.general_user_profile gup ON gup.general_user_profile_id = gu.general_user_general_user_profile
WHERE ct.deleted_at IS NULL
  AND ct.status = ?
  AND ct.type = ?
  AND ct.case_id = ?
ORDER BY ct.completed_at ASC NULLS LAST, ct.created_at ASC
`
	items := make([]models.CaseTaskCommentItem, 0)
	err := r.db.WithContext(ctx).Raw(q, models.CaseTaskStatusDone, models.CaseTaskTypeComentarioEquipoTecnico, caseID).Scan(&items).Error
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

// caseTaskWithRelationsBaseSQL enriquece case_task con datos de barrier_v2 y
// victim_case mediante LEFT JOIN, sin filtrar por tipo de tarea.
//
//	case_task.barrier_id (UUID) → barrier_v2.id (UUID)  — comparación directa sin cast
//	case_task.case_id  (varchar) → victim_case.victim_case_i_code (varchar)
//
// Los campos de barrera quedan vacíos en las tareas que no nacen de una barrera:
// el JOIN es LEFT y las columnas van con COALESCE.
const caseTaskWithRelationsBaseSQL = `
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
`

// caseTaskWithRelationsSQL restringe la base a las tareas nacidas de una barrera:
// es la consulta de la pantalla "Mis Barreras", no la bandeja general de tareas.
const caseTaskWithRelationsSQL = caseTaskWithRelationsBaseSQL + `
  AND ct.barrier_id IS NOT NULL
`

func (r *caseTaskRepository) FindByAssignedUserIDWithRelations(ctx context.Context, assignedUserID string) ([]models.CaseTaskWithRelations, error) {
	var items []models.CaseTaskWithRelations
	err := r.db.WithContext(ctx).
		Raw(caseTaskWithRelationsSQL+" AND ct.assigned_user_id = ? ORDER BY ct.created_at DESC", assignedUserID).
		Scan(&items).Error
	return items, err
}

func (r *caseTaskRepository) FindByAssignedUserIDPaged(ctx context.Context, assignedUserID, status, role string, page, pageSize int) (PageResult[models.CaseTaskWithRelations], error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	// Base: tareas asignadas al usuario. Para el Equipo Técnico ('te') se suman
	// las tareas comite_caso del pool (sin dueño), visibles para todo el equipo.
	var where string
	var args []interface{}
	if role == "te" {
		where = " AND (ct.assigned_user_id = ? OR (ct.type = 'comite_caso' AND ct.assigned_user_id = ''))"
		args = []interface{}{assignedUserID}
	} else {
		where = " AND ct.assigned_user_id = ?"
		args = []interface{}{assignedUserID}
	}
	if status != "" {
		where += " AND ct.status = ?"
		args = append(args, status)
	}

	result := PageResult[models.CaseTaskWithRelations]{
		Items:    []models.CaseTaskWithRelations{},
		Page:     page,
		PageSize: pageSize,
	}

	// El COUNT va sobre case_task sola: los LEFT JOIN no cambian el número de filas
	// y traerlos aquí solo encarece la consulta.
	countSQL := "SELECT COUNT(*) FROM salvia.case_task ct WHERE ct.deleted_at IS NULL" + where
	if err := r.db.WithContext(ctx).Raw(countSQL, args...).Scan(&result.Total).Error; err != nil {
		return result, err
	}

	dataSQL := caseTaskWithRelationsBaseSQL + where + " ORDER BY ct.created_at DESC LIMIT ? OFFSET ?"
	dataArgs := append(append([]interface{}{}, args...), pageSize, (page-1)*pageSize)

	err := r.db.WithContext(ctx).Raw(dataSQL, dataArgs...).Scan(&result.Items).Error
	return result, err
}
