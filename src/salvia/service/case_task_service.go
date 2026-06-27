package service

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

var ErrCaseTaskNotFound = errors.New("case_task: registro no encontrado")

// CaseTaskService define las operaciones de negocio sobre CaseTask.
type CaseTaskService interface {
	GetByID(ctx context.Context, id string) (*models.CaseTask, error)
	List(ctx context.Context, page, limit int) (repository.PageResult[models.CaseTask], error)
	Create(ctx context.Context, ct *models.CaseTask) error
	Update(ctx context.Context, ct *models.CaseTask) error
	Delete(ctx context.Context, id string) error

	// ListByCaseID devuelve todas las tareas de un caso.
	ListByCaseID(ctx context.Context, caseID string) ([]models.CaseTask, error)

	// ListByCaseIDAndBarrier devuelve las tareas de un caso filtradas por barrera.
	ListByCaseIDAndBarrier(ctx context.Context, caseID, barrierID string) ([]models.CaseTask, error)

	// UpdateEntityLetter actualiza campos de un EntityLetter.
	UpdateEntityLetter(ctx context.Context, entityLetterID string, fields map[string]interface{}) error

	// GetDB retorna la instancia de GORM para queries directos.
	GetDB() *gorm.DB

	// Reassign reasigna una tarea a otro usuario.
	Reassign(ctx context.Context, taskID, newAssignedUserID string) error

	// ListByAssignedUserIDWithRelations es la consulta principal para la pantalla
	// "Mis Barreras": devuelve las tareas con barrierId asignadas al usuario,
	// enriquecidas con datos de barrier_v2 y victim_case.
	ListByAssignedUserIDWithRelations(ctx context.Context, assignedUserID string) ([]models.CaseTaskWithRelations, error)

	// Complete ejecuta el flujo "Cuando Guardar Gestionar":
	//   1. Actualiza CaseTask → status=Done, result, completedAt
	//   2. Actualiza BarrierV2 → status=Articulada  (si la tarea tiene barrierId)
	//   3. Crea un CaseTimelineEvent de tipo BARRERA_ARTICULADA
	Complete(ctx context.Context, id string, result string) error
}

// CaseTaskServiceDeps agrupa las dependencias necesarias para CaseTaskService.
type CaseTaskServiceDeps struct {
	CaseTaskRepo     repository.CaseTaskRepository
	BarrierV2Repo    repository.BarrierV2Repository
	CaseTimelineRepo repository.CaseTimelineEventRepository
	EntityLetterRepo repository.EntityLetterRepository
	DB               *gorm.DB
}

type caseTaskService struct {
	repo             repository.CaseTaskRepository
	barrierRepo      repository.BarrierV2Repository
	timelineRepo     repository.CaseTimelineEventRepository
	entityLetterRepo repository.EntityLetterRepository
	db               *gorm.DB
}

func NewCaseTaskService(deps CaseTaskServiceDeps) CaseTaskService {
	return &caseTaskService{
		repo:             deps.CaseTaskRepo,
		barrierRepo:      deps.BarrierV2Repo,
		timelineRepo:     deps.CaseTimelineRepo,
		entityLetterRepo: deps.EntityLetterRepo,
		db:               deps.DB,
	}
}

func (s *caseTaskService) GetDB() *gorm.DB {
	return s.db
}

func (s *caseTaskService) GetByID(ctx context.Context, id string) (*models.CaseTask, error) {
	ct, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCaseTaskNotFound
		}
		return nil, err
	}
	return ct, nil
}

func (s *caseTaskService) List(ctx context.Context, page, limit int) (repository.PageResult[models.CaseTask], error) {
	return s.repo.FindWithPagination(ctx, page, limit)
}

func (s *caseTaskService) ListByCaseID(ctx context.Context, caseID string) ([]models.CaseTask, error) {
	return s.repo.FindByCaseID(ctx, caseID)
}

func (s *caseTaskService) ListByCaseIDAndBarrier(ctx context.Context, caseID, barrierID string) ([]models.CaseTask, error) {
	return s.repo.FindByCaseIDAndBarrier(ctx, caseID, barrierID)
}

func (s *caseTaskService) UpdateEntityLetter(ctx context.Context, entityLetterID string, fields map[string]interface{}) error {
	if s.entityLetterRepo == nil {
		return nil
	}
	return s.entityLetterRepo.UpdateFields(ctx, entityLetterID, fields)
}

func (s *caseTaskService) Reassign(ctx context.Context, taskID, newAssignedUserID string) error {
	_, err := s.repo.FindByID(ctx, taskID)
	if err != nil {
		return ErrCaseTaskNotFound
	}
	return s.repo.UpdateFields(ctx, taskID, map[string]interface{}{
		"assigned_user_id": newAssignedUserID,
	})
}

func (s *caseTaskService) Create(ctx context.Context, ct *models.CaseTask) error {
	if ct.Status == "" {
		ct.Status = models.CaseTaskStatusToDo
	}
	return s.repo.Create(ctx, ct)
}

func (s *caseTaskService) Update(ctx context.Context, ct *models.CaseTask) error {
	return s.repo.Update(ctx, ct)
}

func (s *caseTaskService) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrCaseTaskNotFound
	}
	return err
}

func (s *caseTaskService) ListByAssignedUserIDWithRelations(ctx context.Context, assignedUserID string) ([]models.CaseTaskWithRelations, error) {
	return s.repo.FindByAssignedUserIDWithRelations(ctx, assignedUserID)
}

// Complete ejecuta el flujo "Cuando Guardar Gestionar":
//
//  1. Lee la CaseTask para obtener barrierID, caseID y assignedUserID.
//  2. Marca la CaseTask como Done (status, result, completedAt).
//  3. Actualiza la BarrierV2 relacionada a status "Articulada".
//  4. Crea un CaseTimelineEvent de tipo BARRERA_ARTICULADA.
//
// Los pasos 3 y 4 son best-effort: si fallan se loguea el error pero NO se
// revierte la actualización del CaseTask (que ya se completó exitosamente).
func (s *caseTaskService) Complete(ctx context.Context, id string, result string) error {
	// ── 1. Leer la tarea ──────────────────────────────────────────────────────
	task, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCaseTaskNotFound
		}
		return fmt.Errorf("complete: leer task: %w", err)
	}

	// ── 2. Actualizar CaseTask ────────────────────────────────────────────────
	now := time.Now()
	if err := s.repo.UpdateFields(ctx, id, map[string]interface{}{
		"status":       models.CaseTaskStatusDone,
		"result":       result,
		"completed_at": now,
	}); err != nil {
		return fmt.Errorf("complete: actualizar case_task: %w", err)
	}

	// ── 3. Actualizar BarrierV2 → Articulada ─────────────────────────────────
	if task.BarrierID != nil && *task.BarrierID != "" {
		barrierIDStr := (*task.BarrierID)
		if err := s.barrierRepo.UpdateFields(ctx, barrierIDStr, map[string]interface{}{
			"status": models.BarrierV2StatusArticulada,
		}); err != nil {
			// No bloqueante: la tarea ya quedó en Done
			log.Printf("[CaseTaskService.Complete] WARN: no se pudo actualizar barrier_v2 %s: %v", barrierIDStr, err)
		}
	}

	// ── 4. Crear CaseTimelineEvent ────────────────────────────────────────────
	barrierIDForEvent := ""
	if task.BarrierID != nil {
		barrierIDForEvent = *task.BarrierID
	}

	event := &models.CaseTimelineEvent{
		CaseID:      task.CaseID,
		Category:    models.TimelineCategoryBarreras,
		Type:        models.TimelineTypeBarreraArticulada,
		EventType:   models.TimelineEventBarreraArticulada,
		Icon:        models.TimelineIconBarreraArticulada,
		Color:       models.TimelineColorGreen,
		Date:        now,
		EventUserID: task.AssignedUserID,
		BarrierID:   barrierIDForEvent,
		TaskID:      task.ID,
		Description: result,
	}

	if err := s.timelineRepo.Create(ctx, event); err != nil {
		// No bloqueante
		log.Printf("[CaseTaskService.Complete] WARN: no se pudo crear timeline event para case %s: %v", task.CaseID, err)
	}

	return nil
}
