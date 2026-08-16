package service

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"gorm.io/datatypes"
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

	// ListByCaseIDAndPsychosocial devuelve las tareas de un caso filtradas por remisión psicosocial.
	ListByCaseIDAndPsychosocial(ctx context.Context, caseID, psychosocialID string) ([]models.CaseTask, error)

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

	// CompleteWithFormData completa una CaseTask del tipo gestion_llamada,
	// proyectar_oficio, comite_caso o Corregir oficio guardando el JSON del
	// formulario y ejecutando los efectos de lado correspondientes a cada tipo.
	CompleteWithFormData(ctx context.Context, id string, userId string, formData datatypes.JSON) (*models.CaseTask, error)

	// CreateGestionPropia crea una CaseTask ya completada (sin pasar por ToDo)
	// para una gestión que el Enlace Territorial realizó por iniciativa propia
	// sobre una barrera. Valida que la barrera pertenezca a enlaceDepartmentID
	// antes de crear nada. Ver DocsMD/Otros/temp/req-registrar-gestion-propia-enlace.md.
	CreateGestionPropia(ctx context.Context, caseID, barrierID, tipo, descripcion, userID, enlaceDepartmentID string) (*models.CaseTask, error)
}

// ErrBarrierDepartmentMismatch indica que la barrera no pertenece al
// departamento asignado del Enlace que intenta registrar la gestión propia.
var ErrBarrierDepartmentMismatch = errors.New("case_task: la barrera no pertenece al departamento asignado del enlace")

// ErrBarrierNotFound indica que la barrera referenciada no existe.
var ErrBarrierNotFound = errors.New("case_task: barrera no encontrada")

// CaseTaskServiceDeps agrupa las dependencias necesarias para CaseTaskService.
type CaseTaskServiceDeps struct {
	CaseTaskRepo     repository.CaseTaskRepository
	BarrierV2Repo    repository.BarrierV2Repository
	CaseTimelineRepo repository.CaseTimelineEventRepository
	EntityLetterSvc  EntityLetterService
	EntityLetterRepo repository.EntityLetterRepository
	DB               *gorm.DB
}

type caseTaskService struct {
	repo             repository.CaseTaskRepository
	barrierRepo      repository.BarrierV2Repository
	timelineRepo     repository.CaseTimelineEventRepository
	entityLetterSvc  EntityLetterService
	entityLetterRepo repository.EntityLetterRepository
	db               *gorm.DB
}

func NewCaseTaskService(deps CaseTaskServiceDeps) CaseTaskService {
	return &caseTaskService{
		repo:             deps.CaseTaskRepo,
		barrierRepo:      deps.BarrierV2Repo,
		timelineRepo:     deps.CaseTimelineRepo,
		entityLetterSvc:  deps.EntityLetterSvc,
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

func (s *caseTaskService) ListByCaseIDAndPsychosocial(ctx context.Context, caseID, psychosocialID string) ([]models.CaseTask, error) {
	return s.repo.FindByCaseIDAndPsychosocial(ctx, caseID, psychosocialID)
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
	if err := s.repo.Create(ctx, ct); err != nil {
		return err
	}
	// Si la tarea está asociada a una barrera y es pendiente, marcar barrera como "En Gestion"
	if ct.BarrierID != nil && *ct.BarrierID != "" && ct.Status == models.CaseTaskStatusToDo {
		s.barrierRepo.UpdateFields(ctx, *ct.BarrierID, map[string]interface{}{
			"status": models.BarrierV2StatusEnGestion,
		})
	}
	return nil
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

// CompleteWithFormData completa una CaseTask de tipo gestion_llamada,
// proyectar_oficio, comite_caso o Corregir oficio. Persiste el JSON del
// formulario y ejecuta los efectos de lado propios de cada tipo.
//
// Para proyectar_oficio: guarda form_data en la tarea y delega en
// EntityLetterService.PerformAction("proyectar") que actualiza el oficio
// y marca la tarea Done a través de completarCaseTask().
//
// Para Corregir oficio: delega en PerformAction("corregir") que transiciona
// el entity_letter (en_correccion → para_revisar) y marca la tarea Done.
//
// Para gestion_llamada y comite_caso: marca la tarea Done + guarda form_data
// primero; los efectos de lado son best-effort (se loguean pero no abortan).
func (s *caseTaskService) CompleteWithFormData(ctx context.Context, id string, userId string, formData datatypes.JSON) (*models.CaseTask, error) {
	// ── 1. Cargar la tarea ────────────────────────────────────────────────────
	task, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCaseTaskNotFound
		}
		return nil, fmt.Errorf("CompleteWithFormData: cargar tarea: %w", err)
	}

	// ── 2. Parsear formData ───────────────────────────────────────────────────
	var fd map[string]interface{}
	if err := json.Unmarshal(formData, &fd); err != nil {
		return nil, fmt.Errorf("CompleteWithFormData: formData inválido: %w", err)
	}

	now := time.Now()

	// ── 3. Lógica por tipo ────────────────────────────────────────────────────
	switch task.Type {

	case "proyectar_oficio":
		// Guardar form_data antes de llamar a PerformAction (sin cambiar status aún
		// — PerformAction llamará completarCaseTask() que lo marcará Done).
		if err := s.repo.UpdateFields(ctx, id, map[string]interface{}{
			"form_data": formData,
		}); err != nil {
			return nil, fmt.Errorf("CompleteWithFormData: guardar form_data: %w", err)
		}

		if task.EntityLetterID == nil {
			return nil, fmt.Errorf("CompleteWithFormData: proyectar_oficio sin entity_letter_id")
		}

		input := s.buildProyectarInput(userId, fd)
		input.TaskDescription = &task.Description
		if _, err := s.entityLetterSvc.PerformAction(ctx, *task.EntityLetterID, input); err != nil {
			return nil, fmt.Errorf("CompleteWithFormData: proyectar oficio: %w", err)
		}

	case "Corregir oficio":
		// Sin form_data relevante: delega directamente en PerformAction("corregir")
		// que transiciona el entity_letter (en_correccion → para_revisar) y llama
		// completarCaseTask() [fire-and-forget] para marcar esta tarea Done.
		if task.EntityLetterID == nil {
			return nil, fmt.Errorf("CompleteWithFormData: Corregir oficio sin entity_letter_id")
		}

		input := ActionInput{Action: "corregir", UserID: userId, TaskDescription: &task.Description}
		if _, err := s.entityLetterSvc.PerformAction(ctx, *task.EntityLetterID, input); err != nil {
			return nil, fmt.Errorf("CompleteWithFormData: corregir oficio: %w", err)
		}

	case "gestion_llamada":
		if err := s.repo.UpdateFields(ctx, id, map[string]interface{}{
			"form_data":    formData,
			"status":       models.CaseTaskStatusDone,
			"completed_at": now,
		}); err != nil {
			return nil, fmt.Errorf("CompleteWithFormData: marcar Done: %w", err)
		}
		go s.sideEffectsGestionLlamada(context.Background(), task, fd, userId, now)

	case "comite_caso":
		if err := s.repo.UpdateFields(ctx, id, map[string]interface{}{
			"form_data":    formData,
			"status":       models.CaseTaskStatusDone,
			"completed_at": now,
		}); err != nil {
			return nil, fmt.Errorf("CompleteWithFormData: marcar Done: %w", err)
		}
		go s.sideEffectsComiteCaso(context.Background(), task, fd, userId, now)

	case "justificar_remision":
		if err := s.repo.UpdateFields(ctx, id, map[string]interface{}{
			"form_data":    formData,
			"status":       models.CaseTaskStatusDone,
			"completed_at": now,
		}); err != nil {
			return nil, fmt.Errorf("CompleteWithFormData: marcar Done: %w", err)
		}
		go s.sideEffectsJustificarRemision(context.Background(), task, userId, now)

	default:
		return nil, fmt.Errorf("CompleteWithFormData: tipo de tarea no soportado: %s", task.Type)
	}

	// ── 4. Barrier OPEN → En Gestion (fire-and-forget, aplica a cualquier tipo) ──
	if task.BarrierID != nil && *task.BarrierID != "" {
		go s.actualizarBarreraEnGestion(context.Background(), *task.BarrierID, id)
	}

	// ── 5. Retornar tarea actualizada ─────────────────────────────────────────
	return s.repo.FindByID(ctx, id)
}

// CreateGestionPropia crea una CaseTask ya completada para una gestión que el
// Enlace Territorial registra por iniciativa propia sobre una barrera — a
// diferencia de Complete/CompleteWithFormData, no existe una CaseTask "ToDo"
// previa: se crea y se completa en un solo paso.
//
//  1. Verifica que la barrera exista y pertenezca a enlaceDepartmentID.
//  2. Crea la CaseTask con status=Done directo (category=Barreras, type=gestion_propia).
//  3. Crea un CaseTimelineEvent (fire-and-forget).
//  4. Si la barrera estaba OPEN, la transiciona a "En Gestion" (fire-and-forget,
//     mismo mecanismo que actualizarBarreraEnGestion).
func (s *caseTaskService) CreateGestionPropia(ctx context.Context, caseID, barrierID, tipo, descripcion, userID, enlaceDepartmentID string) (*models.CaseTask, error) {
	// ── 1. Verificar barrera y departamento ──────────────────────────────────
	barrier, err := s.barrierRepo.FindByID(ctx, barrierID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBarrierNotFound
		}
		return nil, fmt.Errorf("CreateGestionPropia: leer barrera: %w", err)
	}
	if barrier.DepartmentID != enlaceDepartmentID {
		return nil, ErrBarrierDepartmentMismatch
	}

	// ── 2. Crear la CaseTask ya completada ───────────────────────────────────
	now := time.Now()
	formData, err := json.Marshal(map[string]string{"subtipo": tipo})
	if err != nil {
		return nil, fmt.Errorf("CreateGestionPropia: serializar formData: %w", err)
	}

	newTask := &models.CaseTask{
		Category:       "Barreras",
		Type:           "gestion_propia",
		Description:    descripcion,
		Status:         models.CaseTaskStatusDone,
		CompletedAt:    &now,
		AssignedUserID: userID,
		CaseID:         caseID,
		BarrierID:      &barrierID,
		FormData:       datatypes.JSON(formData),
	}
	if err := s.repo.Create(ctx, newTask); err != nil {
		return nil, fmt.Errorf("CreateGestionPropia: crear case_task: %w", err)
	}

	// ── 3. Timeline event (fire-and-forget) ──────────────────────────────────
	event := &models.CaseTimelineEvent{
		CaseID:      caseID,
		Category:    models.TimelineCategoryBarreras,
		Type:        models.TimelineTypeGestionPropia,
		EventType:   models.TimelineEventGestionPropia,
		Icon:        models.TimelineIconGestionPropia,
		Color:       models.TimelineColorOrange,
		Date:        now,
		EventUserID: userID,
		BarrierID:   barrierID,
		TaskID:      newTask.ID,
		Description: "Enlace territorial registró una gestión por iniciativa propia (" + tipo + "): " + descripcion,
	}
	if err := s.timelineRepo.Create(ctx, event); err != nil {
		log.Printf("[CaseTaskService.CreateGestionPropia] WARN: no se pudo crear timeline event para case %s: %v", caseID, err)
	}

	// ── 4. Barrier OPEN → En Gestion (fire-and-forget) ───────────────────────
	go s.actualizarBarreraEnGestion(context.Background(), barrierID, newTask.ID)

	return newTask, nil
}

// actualizarBarreraEnGestion transiciona barrier_v2 de OPEN a "En Gestion" si aplica.
// Es fire-and-forget: los errores se loguean pero no abortan el flujo principal.
func (s *caseTaskService) actualizarBarreraEnGestion(ctx context.Context, barrierID string, taskID string) {
	log.Printf("[CaseTaskService] DEBUG: actualizarBarreraEnGestion — barrier=%s task=%s", barrierID, taskID)

	barrier, err := s.barrierRepo.FindByID(ctx, barrierID)
	if err != nil {
		log.Printf("[CaseTaskService] WARN: actualizarBarreraEnGestion no pudo leer barrier %s (task %s): %v", barrierID, taskID, err)
		return
	}

	log.Printf("[CaseTaskService] DEBUG: barrier %s status actual=%s", barrierID, barrier.Status)

	if barrier.Status != models.BarrierV2StatusOpen {
		log.Printf("[CaseTaskService] DEBUG: barrier %s no está OPEN — sin cambios", barrierID)
		return
	}

	if err := s.barrierRepo.UpdateStatus(ctx, barrierID, models.BarrierV2StatusEnGestion); err != nil {
		log.Printf("[CaseTaskService] WARN: actualizarBarreraEnGestion no pudo actualizar barrier %s (task %s): %v", barrierID, taskID, err)
		return
	}

	log.Printf("[CaseTaskService] DEBUG: barrier %s actualizada OPEN → En Gestion (task %s)", barrierID, taskID)
}

// buildProyectarInput construye el ActionInput para PerformAction("proyectar")
// a partir del formData parseado.
func (s *caseTaskService) buildProyectarInput(userId string, fd map[string]interface{}) ActionInput {
	input := ActionInput{
		Action: "proyectar",
		UserID: userId,
	}
	if v, ok := fd["departamentoId"].(string); ok {
		input.DepartmentID = &v
	}
	if v, ok := fd["ciudadId"].(string); ok {
		input.CityID = &v
	}
	if v, ok := fd["municipioId"].(string); ok {
		input.TownID = &v
	}
	if v, ok := fd["entidadNombre"].(string); ok {
		input.EntityName = &v
	}
	if v, ok := fd["funcionario"].(string); ok {
		input.OfficialDependency = &v
	}
	if v, ok := fd["asunto"].(string); ok {
		input.Subject = &v
	}
	if v, ok := fd["rutaKofax"].(string); ok {
		input.UrlKofax = &v
	}
	if v, ok := fd["entidadId"].(float64); ok {
		id64 := int64(v)
		input.EntityBranchID = &id64
	}
	return input
}

// sideEffectsGestionLlamada ejecuta los efectos de lado de una tarea gestion_llamada
// de forma asíncrona (fire-and-forget). Los errores se loguean pero no abortan.
func (s *caseTaskService) sideEffectsGestionLlamada(ctx context.Context, task *models.CaseTask, fd map[string]interface{}, userId string, now time.Time) {
	generaOficio, _ := fd["generaOficio"].(bool)
	if generaOficio {
		letter := s.buildEntityLetterFromFormData(task, fd, userId, models.EntityLetterStateParaRevisar)
		if err := s.entityLetterRepo.Create(ctx, letter); err != nil {
			log.Printf("[CaseTaskService] WARN: gestion_llamada no pudo crear entity_letter para tarea %s: %v", task.ID, err)
		}
	}

	barrierID := ""
	if task.BarrierID != nil {
		barrierID = *task.BarrierID
	}
	event := &models.CaseTimelineEvent{
		CaseID:      task.CaseID,
		Category:    models.TimelineCategoryBarreras,
		Type:        "Gestión de Llamada",
		EventType:   "GESTION_LLAMADA",
		Icon:        "phone",
		Color:       models.TimelineColorBlue,
		Date:        now,
		Description: "Tarea completada: " + task.Description,
		EventUserID: userId,
		BarrierID:   barrierID,
		TaskID:      task.ID,
	}
	if err := s.timelineRepo.Create(ctx, event); err != nil {
		log.Printf("[CaseTaskService] WARN: gestion_llamada no pudo crear timeline event para case %s: %v", task.CaseID, err)
	}
}

// sideEffectsComiteCaso ejecuta los efectos de lado de una tarea comite_caso
// de forma asíncrona (fire-and-forget). Los errores se loguean pero no abortan.
func (s *caseTaskService) sideEffectsComiteCaso(ctx context.Context, task *models.CaseTask, fd map[string]interface{}, userId string, now time.Time) {
	decisiones, _ := fd["decisiones"].([]interface{})

	for _, d := range decisiones {
		decision, _ := d.(string)
		switch decision {

		case "activar_enlace":
			if task.BarrierID != nil && *task.BarrierID != "" {
				if err := s.barrierRepo.UpdateFields(ctx, *task.BarrierID, map[string]interface{}{
					"enlace_activado": true,
				}); err != nil {
					log.Printf("[CaseTaskService] WARN: comite_caso no pudo activar enlace en barrier %s: %v", *task.BarrierID, err)
				}
			}

		case "oficio":
			barrierID := ""
			if task.BarrierID != nil {
				barrierID = *task.BarrierID
			}
			letter := &models.EntityLetter{
				BarrierID: barrierID,
				CaseID:    task.CaseID,
				AgentID:   &userId,
				State:     models.EntityLetterStatePorProyectar,
			}
			if err := s.entityLetterRepo.Create(ctx, letter); err != nil {
				log.Printf("[CaseTaskService] WARN: comite_caso no pudo crear entity_letter para tarea %s: %v", task.ID, err)
				continue
			}
			newTask := &models.CaseTask{
				Category:       "Barreras",
				Type:           "proyectar_oficio",
				Description:    "Proyectar oficio — Decisión de comité",
				Status:         models.CaseTaskStatusToDo,
				AssignedUserID: userId,
				CaseID:         task.CaseID,
				BarrierID:      task.BarrierID,
				EntityLetterID: &letter.ID,
			}
			if err := s.repo.Create(ctx, newTask); err != nil {
				log.Printf("[CaseTaskService] WARN: comite_caso no pudo crear task proyectar_oficio para tarea %s: %v", task.ID, err)
			}
		}
	}

	barrierID := ""
	if task.BarrierID != nil {
		barrierID = *task.BarrierID
	}
	event := &models.CaseTimelineEvent{
		CaseID:      task.CaseID,
		Category:    models.TimelineCategoryGeneral,
		Type:        "Decisiones del Comité",
		EventType:   "COMITE_CASO",
		Icon:        "users",
		Color:       models.TimelineColorPurple,
		Date:        now,
		Description: "Tarea completada: " + task.Description,
		EventUserID: userId,
		BarrierID:   barrierID,
		TaskID:      task.ID,
	}
	if err := s.timelineRepo.Create(ctx, event); err != nil {
		log.Printf("[CaseTaskService] WARN: comite_caso no pudo crear timeline event para case %s: %v", task.CaseID, err)
	}
}

// sideEffectsJustificarRemision se ejecuta cuando el agente responde al
// motivo de devolución de una remisión psicosocial (case_task
// "justificar_remision"): reabre la remisión (status → 'abierto') y crea una
// nueva case_task "validar_remision" para que el psicólogo/a vuelva a
// evaluarla. La nueva tarea se asigna al profesional de la dupla (o al
// profesional individual) que tiene la remisión — si no hay ninguno asignado,
// la tarea queda sin asignar, igual que la remisión.
func (s *caseTaskService) sideEffectsJustificarRemision(ctx context.Context, task *models.CaseTask, userId string, now time.Time) {
	if task.PsychosocialSupportID == nil || *task.PsychosocialSupportID == "" {
		log.Printf("[CaseTaskService] WARN: justificar_remision sin psychosocial_support_id (tarea %s)", task.ID)
		return
	}
	psID := *task.PsychosocialSupportID

	var duplaID, professionalID string
	s.db.WithContext(ctx).Raw(`
		SELECT COALESCE(dupla_id, ''), COALESCE(professional_id, '')
		FROM salvia.psychosocial_support WHERE id = ?
	`, psID).Row().Scan(&duplaID, &professionalID)

	assignedTo := professionalID
	if duplaID != "" {
		var psychologistID string
		s.db.WithContext(ctx).Raw(`SELECT COALESCE(psychologist_id, '') FROM salvia.dupla WHERE id = ?`, duplaID).Row().Scan(&psychologistID)
		assignedTo = psychologistID
	}

	if err := s.db.WithContext(ctx).Exec(`
		UPDATE salvia.psychosocial_support SET status = 'abierto', updated_at = NOW() WHERE id = ?
	`, psID).Error; err != nil {
		log.Printf("[CaseTaskService] WARN: justificar_remision no pudo reabrir remisión %s: %v", psID, err)
	}

	newTask := &models.CaseTask{
		Category:              "Psicosocial",
		Type:                  "validar_remision",
		Description:           "Validar que remision a psicosocial es valida",
		Status:                models.CaseTaskStatusToDo,
		AssignedUserID:        assignedTo,
		CaseID:                task.CaseID,
		FollowUpID:            task.FollowUpID,
		PsychosocialSupportID: task.PsychosocialSupportID,
	}
	if err := s.repo.Create(ctx, newTask); err != nil {
		log.Printf("[CaseTaskService] WARN: justificar_remision no pudo crear nueva task validar_remision para remisión %s: %v", psID, err)
	}

	event := &models.CaseTimelineEvent{
		CaseID:                task.CaseID,
		Category:              models.TimelineCategoryPsicosocial,
		Type:                  "Remisión Justificada",
		Icon:                  "reply",
		Color:                 models.TimelineColorBlue,
		Date:                  now,
		Description:           "Agente respondió al motivo de devolución — remisión vuelve a validación",
		EventUserID:           userId,
		TaskID:                task.ID,
		PsychosocialSupportID: psID,
	}
	if err := s.timelineRepo.Create(ctx, event); err != nil {
		log.Printf("[CaseTaskService] WARN: justificar_remision no pudo crear timeline event para remisión %s: %v", psID, err)
	}
}

// buildEntityLetterFromFormData construye un EntityLetter a partir del formData
// de gestion_llamada. El state se pasa explícitamente.
func (s *caseTaskService) buildEntityLetterFromFormData(task *models.CaseTask, fd map[string]interface{}, userId string, state string) *models.EntityLetter {
	barrierID := ""
	if task.BarrierID != nil {
		barrierID = *task.BarrierID
	}
	letter := &models.EntityLetter{
		BarrierID:  barrierID,
		CaseID:     task.CaseID,
		AgentID:    &userId,
		RegisterBy: &userId,
		State:      state,
	}
	if v, ok := fd["departamentoId"].(string); ok {
		letter.DepartmentID = &v
	}
	if v, ok := fd["ciudadId"].(string); ok {
		letter.CityID = &v
	}
	if v, ok := fd["municipioId"].(string); ok {
		letter.TownID = &v
	}
	if v, ok := fd["entidadNombre"].(string); ok {
		letter.Entidad = &v
	}
	if v, ok := fd["funcionario"].(string); ok {
		letter.OfficialDependency = &v
	}
	if v, ok := fd["asunto"].(string); ok {
		letter.Subject = &v
	}
	if v, ok := fd["rutaKofax"].(string); ok {
		letter.UrlKofax = &v
	}
	if v, ok := fd["entidadId"].(float64); ok {
		id64 := int64(v)
		letter.EntityBranchID = &id64
	}
	return letter
}
