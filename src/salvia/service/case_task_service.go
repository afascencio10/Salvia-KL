package service

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var ErrCaseTaskNotFound = errors.New("case_task: registro no encontrado")
var ErrEntityUserNotFound = errors.New("no hay usuario de entidad asociado a esta sede")
var ErrEntityUserNotUnique = errors.New("hay más de un usuario de entidad para esta sede")

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

	// ListByReferralID devuelve las tareas de una remisión de medidas de emergencia.
	ListByReferralID(ctx context.Context, referralID string) ([]models.CaseTask, error)

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

	// ListByAssignedUserIDPaged es la consulta de la pantalla "Mis Tareas":
	// todas las tareas del agente (no solo las de barrera), paginadas y con
	// filtro opcional por estado.
	ListByAssignedUserIDPaged(ctx context.Context, assignedUserID, status, role string, page, pageSize int) (repository.PageResult[models.CaseTaskWithRelations], error)

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

	// SaveArticulacionTerritorial completa o crea+completa una tarea
	// articulacion_territorial (Mis Barreras y gestión propia en Detalle).
	// Si requiere seguimiento, crea un ToDo compromiso_con_enlace.
	// Completar una ToDo existente exige user_territory_link_id.
	// Crear sin taskId admite asignación al usuario o mismo departamento del enlace.
	SaveArticulacionTerritorial(ctx context.Context, userID, enlaceDepartmentID, barrierID, taskID string, formData datatypes.JSON) (*models.CaseTask, error)

	// SaveGestionMecanismo crea siempre una case_task Done para el formulario
	// de Casos Mecanismo. Nunca actualiza una tarea previa.
	SaveGestionMecanismo(ctx context.Context, userID, mechanismID string, formData datatypes.JSON) (*models.CaseTask, error)

	// SaveComentarioEquipoTecnico crea una case_task Done (tipo
	// comentario_equipo_tecnico) con el texto del comentario en Result, y
	// registra el evento correspondiente en el timeline del caso. Solo el rol
	// 'te' puede invocarla (verificado en el controller).
	SaveComentarioEquipoTecnico(ctx context.Context, userID, caseID, texto string) (*models.CaseTask, error)

	// AbrirMesaTecnica crea una case_task ToDo tipo comite_caso para analizar el
	// caso en mesa técnica, sin vincularla a ninguna barrera. justificacion es
	// el motivo de la remisión y queda guardado en Description. Queda en el
	// pool del Equipo Técnico (assigned_user_id vacío) — mismo efecto que
	// marcar la gestión "Alerta por barreras" al guardar un seguimiento. Solo
	// los roles 'sv' y 'te' pueden invocarla (verificado en el controller).
	AbrirMesaTecnica(ctx context.Context, userID, caseID, justificacion string) (*models.CaseTask, error)
}

// ErrBarrierDepartmentMismatch indica que la barrera no pertenece al
// departamento asignado del Enlace que intenta registrar la gestión propia.
var ErrBarrierDepartmentMismatch = errors.New("case_task: la barrera no pertenece al departamento asignado del enlace")

// ErrBarrierNotFound indica que la barrera referenciada no existe.
var ErrBarrierNotFound = errors.New("case_task: barrera no encontrada")

// ErrBarrierLinkMismatch indica que la barrera no está asignada al Enlace
// autenticado (user_territory_link_id).
var ErrBarrierLinkMismatch = errors.New("case_task: la barrera no está asignada a este enlace")

// ErrArticulacionInvalid indica formData incompleto o inválido.
var ErrArticulacionInvalid = errors.New("case_task: datos de articulación inválidos")

// ErrArticulacionTaskMismatch indica que taskId no es ToDo de articulación
// de esa barrera y usuario.
var ErrArticulacionTaskMismatch = errors.New("case_task: la tarea no corresponde a esta articulación")

// ErrComentarioEquipoTecnicoInvalid indica caseId o texto vacío.
var ErrComentarioEquipoTecnicoInvalid = errors.New("case_task: comentario del equipo técnico inválido")

// ErrAbrirMesaTecnicaInvalid indica caseId vacío.
var ErrAbrirMesaTecnicaInvalid = errors.New("case_task: caso inválido para abrir mesa técnica")

// CaseTaskServiceDeps agrupa las dependencias necesarias para CaseTaskService.
type CaseTaskServiceDeps struct {
	CaseTaskRepo     repository.CaseTaskRepository
	BarrierV2Repo    repository.BarrierV2Repository
	CaseTimelineRepo repository.CaseTimelineEventRepository
	EntityLetterSvc  EntityLetterService
	EntityLetterRepo repository.EntityLetterRepository
	EntityCaseSvc    EntityCaseService
	MechanismRepo    repository.CaseArticulatingMechanismRepository
	DB               *gorm.DB
}

type caseTaskService struct {
	repo             repository.CaseTaskRepository
	barrierRepo      repository.BarrierV2Repository
	timelineRepo     repository.CaseTimelineEventRepository
	entityLetterSvc  EntityLetterService
	entityLetterRepo repository.EntityLetterRepository
	entityCaseSvc    EntityCaseService
	mechanismRepo    repository.CaseArticulatingMechanismRepository
	db               *gorm.DB
}

func NewCaseTaskService(deps CaseTaskServiceDeps) CaseTaskService {
	return &caseTaskService{
		repo:             deps.CaseTaskRepo,
		barrierRepo:      deps.BarrierV2Repo,
		timelineRepo:     deps.CaseTimelineRepo,
		entityLetterSvc:  deps.EntityLetterSvc,
		entityLetterRepo: deps.EntityLetterRepo,
		entityCaseSvc:    deps.EntityCaseSvc,
		mechanismRepo:    deps.MechanismRepo,
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

func (s *caseTaskService) ListByReferralID(ctx context.Context, referralID string) ([]models.CaseTask, error) {
	return s.repo.FindByReferralID(ctx, referralID)
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

func (s *caseTaskService) ListByAssignedUserIDPaged(ctx context.Context, assignedUserID, status, role string, page, pageSize int) (repository.PageResult[models.CaseTaskWithRelations], error) {
	return s.repo.FindByAssignedUserIDPaged(ctx, assignedUserID, status, role, page, pageSize)
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

// completeDecidirVia asigna el agente, mueve el estado de la remisión,
// crea la tarea siguiente y cierra la de decidir vía. Si algo falla, no queda nada escrito.
func (s *caseTaskService) completeDecidirVia(ctx context.Context, task *models.CaseTask, formData datatypes.JSON, fd map[string]interface{}, now time.Time) error {
	if task.Status != models.CaseTaskStatusToDo {
		return fmt.Errorf("CompleteWithFormData: la tarea ya no está pendiente")
	}
	if task.EmergencyMeasureReferralID == nil || strings.TrimSpace(*task.EmergencyMeasureReferralID) == "" {
		return fmt.Errorf("CompleteWithFormData: la tarea no tiene remisión")
	}
	via, _ := fd["via"].(string)
	via = strings.TrimSpace(via)
	agentID, _ := fd["agentId"].(string)
	agentID = strings.TrimSpace(agentID)
	if (via != "A" && via != "B") || agentID == "" {
		return fmt.Errorf("CompleteWithFormData: vía o agente inválido")
	}
	if s.db == nil {
		return fmt.Errorf("CompleteWithFormData: sin base de datos")
	}

	var activos int
	if err := s.db.WithContext(ctx).Raw(`
		SELECT COUNT(*)::int
		FROM security.general_user gu
		JOIN security.rel_role_general_user rr ON rr.general_user_id = gu.general_user_id
		JOIN security.role role ON role.role_id = rr.role_id
		WHERE gu.general_user_i_code = ?
		  AND role.role_code = 'am'
		  AND gu.general_user_status = 'e'
	`, agentID).Scan(&activos).Error; err != nil {
		return fmt.Errorf("CompleteWithFormData: validar agente: %w", err)
	}
	if activos == 0 {
		return fmt.Errorf("CompleteWithFormData: el agente no es un Agente de Medidas activo")
	}

	refID := strings.TrimSpace(*task.EmergencyMeasureReferralID)
	status := models.EmergencyReferralStatusAsignada
	next := &models.CaseTask{
		Category:                   models.CaseTaskCategoryMedidas,
		Status:                     models.CaseTaskStatusToDo,
		CaseID:                     task.CaseID,
		FollowUpID:                 task.FollowUpID,
		EmergencyMeasureReferralID: &refID,
	}
	if via == "A" {
		next.Type = models.CaseTaskTypeIntentarContacto3x3
		next.Description = models.CaseTaskDescIntentarContacto3x3
		next.AssignedUserID = agentID
	} else {
		status = models.EmergencyReferralStatusRiesgoPorEscalarCruzRoja
		role := models.CaseTaskRoleSupervisionConvenio
		next.Type = models.CaseTaskTypeEscalarCruzRoja
		next.Description = models.CaseTaskDescEscalarCruzRoja
		next.AssignedRole = &role
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Exec(`
			UPDATE salvia.emergency_measure_referral
			   SET assigned_agent_id = ?,
			       status = ?,
			       updated_at = NOW()
			 WHERE id::text = ?
			   AND status = ?
			   AND deleted_at IS NULL
		`, agentID, status, refID, models.EmergencyReferralStatusRiesgoPorDecidirVia)
		if res.Error != nil {
			return fmt.Errorf("CompleteWithFormData: actualizar remisión: %w", res.Error)
		}
		if res.RowsAffected != 1 {
			return fmt.Errorf("CompleteWithFormData: la remisión ya no está por decidir vía")
		}
		if err := tx.Create(next).Error; err != nil {
			return fmt.Errorf("CompleteWithFormData: crear tarea siguiente: %w", err)
		}
		res = tx.Model(&models.CaseTask{}).
			Where("id = ? AND status = ?", task.ID, models.CaseTaskStatusToDo).
			Updates(map[string]interface{}{
				"form_data":    formData,
				"status":       models.CaseTaskStatusDone,
				"completed_at": now,
			})
		if res.Error != nil {
			return fmt.Errorf("CompleteWithFormData: cerrar la decisión: %w", res.Error)
		}
		if res.RowsAffected != 1 {
			return fmt.Errorf("CompleteWithFormData: la tarea ya no está pendiente")
		}
		return nil
	})
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
		// La tarea comite_caso nace en el pool del Equipo Técnico (assigned_user_id
		// vacío). Al completarla, el integrante 'te' que la gestionó la "reclama":
		// se registra su icode en assigned_user_id para dejar constancia de quién la
		// resolvió (se muestra luego como "Completado por").
		comiteUpdates := map[string]interface{}{
			"form_data":    formData,
			"status":       models.CaseTaskStatusDone,
			"completed_at": now,
		}
		if userId != "" {
			comiteUpdates["assigned_user_id"] = userId
		}
		if err := s.repo.UpdateFields(ctx, id, comiteUpdates); err != nil {
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

	case "solicitud_entidad":
		if err := s.completeSolicitudEntidad(ctx, task, userId, formData, fd, now); err != nil {
			return nil, err
		}

	case models.CaseTaskTypeDecidirViaRiesgoInminente:
		if err := s.completeDecidirVia(ctx, task, formData, fd, now); err != nil {
			return nil, err
		}

	case "respuesta_entidad":
		estado, _ := fd["estadoInstitucional"].(string)
		estado = strings.TrimSpace(estado)
		if estado != "iniciado" && estado != "en_proceso" && estado != "concluido" {
			return nil, fmt.Errorf("CompleteWithFormData: estado institucional inválido")
		}
		if err := s.repo.UpdateFields(ctx, id, map[string]interface{}{
			"form_data":    formData,
			"status":       models.CaseTaskStatusDone,
			"completed_at": now,
		}); err != nil {
			return nil, fmt.Errorf("CompleteWithFormData: marcar Done: %w", err)
		}
		label, _ := fd["label"].(string)
		label = strings.TrimSpace(label)
		if label == "" {
			label = "obligación institucional"
		}
		estadoTxt := map[string]string{
			"iniciado":   "Iniciado",
			"en_proceso": "En proceso",
			"concluido":  "Concluido",
		}[estado]
		s.createEntidadTimelineEvent(ctx, task, userId, now,
			models.TimelineEventRespuestaEntidad,
			models.TimelineTypeRespuestaEntidad,
			models.TimelineIconRespuestaEntidad,
			models.TimelineColorGreen,
			fmt.Sprintf("Respuesta sobre: %s (%s)", label, estadoTxt),
		)

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

type solicitudItem struct {
	ObligationID string
	Label        string
	Descripcion  string
}

func parseSolicitudes(fd map[string]interface{}) ([]solicitudItem, error) {
	raw, ok := fd["solicitudes"]
	if !ok {
		return nil, fmt.Errorf("CompleteWithFormData: solicitudes requeridas")
	}
	arr, ok := raw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("CompleteWithFormData: solicitudes debe ser un arreglo")
	}
	if len(arr) < 1 {
		return nil, fmt.Errorf("CompleteWithFormData: selecciona al menos un ítem del checklist")
	}
	out := make([]solicitudItem, 0, len(arr))
	for i, item := range arr {
		m, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("CompleteWithFormData: solicitud #%d inválida", i+1)
		}
		id, _ := m["obligationId"].(string)
		if strings.TrimSpace(id) == "" {
			return nil, fmt.Errorf("CompleteWithFormData: solicitud #%d sin obligationId", i+1)
		}
		desc, _ := m["descripcion"].(string)
		if strings.TrimSpace(desc) == "" {
			return nil, fmt.Errorf("CompleteWithFormData: la descripción es obligatoria en cada ítem seleccionado")
		}
		label, _ := m["label"].(string)
		out = append(out, solicitudItem{
			ObligationID: strings.TrimSpace(id),
			Label:        strings.TrimSpace(label),
			Descripcion:  strings.TrimSpace(desc),
		})
	}
	return out, nil
}

func (s *caseTaskService) findUsersByEntityBranchID(ctx context.Context, branchID int64) ([]string, error) {
	var ids []string
	err := s.db.WithContext(ctx).Raw(`
		SELECT general_user_i_code
		FROM security.general_user
		WHERE entity_branch_id = ?
		ORDER BY general_user_id ASC
	`, branchID).Scan(&ids).Error
	return ids, err
}

// completeSolicitudEntidad valida el checklist, resuelve el único usuario de
// la sede y, en la misma transacción, marca la tarea Done y crea una
// respuesta_entidad por ítem seleccionado.
func (s *caseTaskService) completeSolicitudEntidad(ctx context.Context, task *models.CaseTask, userId string, formData datatypes.JSON, fd map[string]interface{}, now time.Time) error {
	if task.EntityCaseID == nil || *task.EntityCaseID == "" {
		return fmt.Errorf("CompleteWithFormData: solicitud_entidad sin entity_case_id")
	}
	solicitudes, err := parseSolicitudes(fd)
	if err != nil {
		return err
	}

	var branchID int64
	if err := s.db.WithContext(ctx).Raw(
		`SELECT entity_branch_id FROM salvia.entity_case WHERE id = ? AND deleted_at IS NULL`,
		*task.EntityCaseID,
	).Scan(&branchID).Error; err != nil {
		return fmt.Errorf("CompleteWithFormData: cargar entity_case: %w", err)
	}
	if branchID == 0 {
		return fmt.Errorf("CompleteWithFormData: entity_case no encontrado")
	}

	users, err := s.findUsersByEntityBranchID(ctx, branchID)
	if err != nil {
		return fmt.Errorf("CompleteWithFormData: buscar usuario de entidad: %w", err)
	}
	if len(users) == 0 {
		return ErrEntityUserNotFound
	}
	if len(users) > 1 {
		return ErrEntityUserNotUnique
	}
	assignedTo := users[0]

	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.CaseTask{}).Where("id = ?", task.ID).Updates(map[string]interface{}{
			"form_data":    formData,
			"status":       models.CaseTaskStatusDone,
			"completed_at": now,
		}).Error; err != nil {
			return fmt.Errorf("CompleteWithFormData: marcar Done: %w", err)
		}

		for _, sol := range solicitudes {
			childFD, err := json.Marshal(map[string]interface{}{
				"solicitudTaskId":      task.ID,
				"obligationId":         sol.ObligationID,
				"label":                sol.Label,
				"descripcionSolicitud": sol.Descripcion,
			})
			if err != nil {
				return fmt.Errorf("CompleteWithFormData: serializar form_data respuesta: %w", err)
			}
			desc := sol.Label
			if desc == "" {
				desc = "Responder solicitud de la entidad"
			} else {
				desc = "Responder solicitud: " + desc
			}
			child := &models.CaseTask{
				Category:       "Entidades",
				Type:           "respuesta_entidad",
				Description:    desc,
				AssignedUserID: assignedTo,
				Status:         models.CaseTaskStatusToDo,
				FormData:       datatypes.JSON(childFD),
				CaseID:         task.CaseID,
				FollowUpID:     task.FollowUpID,
				EntityCaseID:   task.EntityCaseID,
			}
			if err := tx.Create(child).Error; err != nil {
				return fmt.Errorf("CompleteWithFormData: crear respuesta_entidad: %w", err)
			}
		}
		return nil
	}); err != nil {
		return err
	}

	entityName, _ := fd["entityName"].(string)
	entityName = strings.TrimSpace(entityName)
	desc := "Solicitudes enviadas a la entidad"
	if entityName != "" {
		desc = "Solicitudes enviadas a " + entityName
	}
	n := len(solicitudes)
	if n == 1 {
		desc += " (1 ítem)"
	} else {
		desc += fmt.Sprintf(" (%d ítems)", n)
	}
	s.createEntidadTimelineEvent(ctx, task, userId, now,
		models.TimelineEventSolicitudEntidad,
		models.TimelineTypeSolicitudEntidad,
		models.TimelineIconSolicitudEntidad,
		models.TimelineColorPurple,
		desc,
	)
	return nil
}

func (s *caseTaskService) createEntidadTimelineEvent(ctx context.Context, task *models.CaseTask, userId string, now time.Time, eventType, typeLabel, icon, color, description string) {
	followUpID := ""
	if task.FollowUpID != nil {
		followUpID = *task.FollowUpID
	}
	event := &models.CaseTimelineEvent{
		CaseID:      task.CaseID,
		Category:    models.TimelineCategoryEntidades,
		Type:        typeLabel,
		EventType:   eventType,
		Icon:        icon,
		Color:       color,
		Date:        now,
		EventUserID: userId,
		TaskID:      task.ID,
		FollowUpID:  followUpID,
		Description: description,
	}
	if err := s.timelineRepo.Create(ctx, event); err != nil {
		log.Printf("[CaseTaskService] WARN: no se pudo crear timeline event entidad para case %s task %s: %v", task.CaseID, task.ID, err)
	}
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

type articulacionFormPayload struct {
	FechaContacto       string `json:"fechaContacto"`
	Compromiso          string `json:"compromiso"`
	RequiereSeguimiento *bool  `json:"requiereSeguimiento"`
	FechaSeguimiento    string `json:"fechaSeguimiento"`
	SugerenciasRemision string `json:"sugerenciasRemision"`
	EntityBranchID      *int64 `json:"entityBranchId"`
}

func parseArticulacionForm(formData datatypes.JSON) (*articulacionFormPayload, error) {
	var fd articulacionFormPayload
	if err := json.Unmarshal(formData, &fd); err != nil {
		return nil, ErrArticulacionInvalid
	}
	fd.FechaContacto = strings.TrimSpace(fd.FechaContacto)
	fd.Compromiso = strings.TrimSpace(fd.Compromiso)
	fd.FechaSeguimiento = strings.TrimSpace(fd.FechaSeguimiento)
	fd.SugerenciasRemision = strings.TrimSpace(fd.SugerenciasRemision)
	if fd.FechaContacto == "" || fd.Compromiso == "" || fd.RequiereSeguimiento == nil {
		return nil, ErrArticulacionInvalid
	}
	if *fd.RequiereSeguimiento && fd.FechaSeguimiento == "" {
		return nil, ErrArticulacionInvalid
	}
	return &fd, nil
}

func barrierAssignedToUser(barrier *models.BarrierV2, userID string) bool {
	return barrier.UserTerritoryLinkID != nil && *barrier.UserTerritoryLinkID == userID
}

// SaveArticulacionTerritorial completa una ToDo existente o crea+completa
// articulacion_territorial. Si hay seguimiento, crea ToDo compromiso_con_enlace.
func (s *caseTaskService) SaveArticulacionTerritorial(ctx context.Context, userID, enlaceDepartmentID, barrierID, taskID string, formData datatypes.JSON) (*models.CaseTask, error) {
	fd, err := parseArticulacionForm(formData)
	if err != nil {
		return nil, err
	}

	barrier, err := s.barrierRepo.FindByID(ctx, barrierID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBarrierNotFound
		}
		return nil, fmt.Errorf("SaveArticulacionTerritorial: leer barrera: %w", err)
	}
	// Autorización: la barrera debe estar asignada al enlace (user_territory_link_id)
	// o pertenecer a su departamento asignado.
	if !barrierAssignedToUser(barrier, userID) && barrier.DepartmentID != enlaceDepartmentID {
		return nil, ErrBarrierDepartmentMismatch
	}

	now := time.Now()

	// El campo entidad no se toma del cliente: siempre se fija con la entidad real
	// de la barrera (solo lectura, no editable ni por API).
	fd.EntityBranchID = barrier.EntityBranchID

	// Múltiples gestiones sin sobreescribir: CADA gestión crea un nuevo registro
	// Done (un compromiso adquirido). No se completa una tarea existente — así el
	// enlace puede registrar cuantas gestiones necesite sobre la misma barrera.
	// Si venía un taskId de una ToDo abierta (la que creó el comité), esa ToDo se
	// cierra la primera vez para que no quede pendiente; las siguientes gestiones
	// simplemente agregan nuevos registros.
	if tid := strings.TrimSpace(taskID); tid != "" {
		if existing, ferr := s.repo.FindByID(ctx, tid); ferr == nil &&
			existing.Type == models.CaseTaskTypeArticulacionTerritorial &&
			existing.Status == models.CaseTaskStatusToDo &&
			existing.BarrierID != nil && *existing.BarrierID == barrierID {
			// Cierra la ToDo inicial sin abortar si falla.
			if err := s.repo.UpdateFields(ctx, existing.ID, map[string]interface{}{
				"status":       models.CaseTaskStatusDone,
				"completed_at": now,
			}); err != nil {
				log.Printf("[CaseTaskService.SaveArticulacionTerritorial] WARN: no se pudo cerrar ToDo inicial %s: %v", existing.ID, err)
			}
		}
	}

	// Reserializar formData con la entidad forzada de la barrera.
	formDataFixed, mErr := json.Marshal(fd)
	if mErr != nil {
		formDataFixed = formData // fallback: usa el original
	}

	newTask := &models.CaseTask{
		Category:       "Barreras",
		Type:           models.CaseTaskTypeArticulacionTerritorial,
		Description:    "Gestión propia",
		Result:         &fd.Compromiso,
		Status:         models.CaseTaskStatusDone,
		CompletedAt:    &now,
		AssignedUserID: userID,
		CaseID:         barrier.CaseID,
		BarrierID:      &barrierID,
		FormData:       datatypes.JSON(formDataFixed),
	}
	if err := s.repo.Create(ctx, newTask); err != nil {
		return nil, fmt.Errorf("SaveArticulacionTerritorial: crear tarea: %w", err)
	}
	task := newTask

	if *fd.RequiereSeguimiento {
		seguimientoData, err := json.Marshal(map[string]interface{}{
			"fechaSeguimiento":   fd.FechaSeguimiento,
			"articulacionTaskId": task.ID,
		})
		if err != nil {
			return nil, fmt.Errorf("SaveArticulacionTerritorial: serializar seguimiento: %w", err)
		}
		follow := &models.CaseTask{
			Category:       "Barreras",
			Type:           models.CaseTaskTypeCompromisoConEnlace,
			Description:    "Seguimiento al compromiso adquirido",
			Status:         models.CaseTaskStatusToDo,
			AssignedUserID: userID,
			CaseID:         barrier.CaseID,
			BarrierID:      &barrierID,
			FormData:       datatypes.JSON(seguimientoData),
		}
		if err := s.repo.Create(ctx, follow); err != nil {
			log.Printf("[CaseTaskService.SaveArticulacionTerritorial] WARN: no se pudo crear compromiso_con_enlace: %v", err)
		}
	}

	event := &models.CaseTimelineEvent{
		CaseID:      barrier.CaseID,
		Category:    models.TimelineCategoryBarreras,
		Type:        models.TimelineTypeArticulacionTerritorial,
		EventType:   models.TimelineEventArticulacionTerritorial,
		Icon:        models.TimelineIconGestionPropia,
		Color:       models.TimelineColorPurple,
		Date:        now,
		EventUserID: userID,
		BarrierID:   barrierID,
		TaskID:      task.ID,
		Description: fd.Compromiso,
	}
	if err := s.timelineRepo.Create(ctx, event); err != nil {
		log.Printf("[CaseTaskService.SaveArticulacionTerritorial] WARN: timeline: %v", err)
	}

	go s.actualizarBarreraEnGestion(context.Background(), barrierID, task.ID)

	return task, nil
}

type gestionMecanismoForm struct {
	AbordadoEnSesion      string `json:"abordadoEnSesion"`
	AbordadoEnSesionTexto string `json:"abordadoEnSesionTexto"`
	FechaSesion           string `json:"fechaSesion"`
	Compromisos           string `json:"compromisos"`
	FechaProxima          string `json:"fechaProxima"`
	DetalleProxima        string `json:"detalleProxima"`
	ReporteGestion        string `json:"reporteGestion"`
	RemisionResponsables  string `json:"remisionResponsables"`
	AsistenciasFecha      string `json:"asistenciasFecha"`
	AsistenciasTexto      string `json:"asistenciasTexto"`
}

func parseGestionMecanismoForm(scope string, formData datatypes.JSON) (*gestionMecanismoForm, error) {
	var fd gestionMecanismoForm
	if len(formData) > 0 {
		if err := json.Unmarshal(formData, &fd); err != nil {
			return nil, ErrGestionMecanismoInvalid
		}
	}
	fd.AbordadoEnSesion = strings.TrimSpace(fd.AbordadoEnSesion)
	fd.AbordadoEnSesionTexto = strings.TrimSpace(fd.AbordadoEnSesionTexto)
	fd.FechaSesion = strings.TrimSpace(fd.FechaSesion)
	fd.Compromisos = strings.TrimSpace(fd.Compromisos)
	fd.FechaProxima = strings.TrimSpace(fd.FechaProxima)
	fd.DetalleProxima = strings.TrimSpace(fd.DetalleProxima)
	fd.ReporteGestion = strings.TrimSpace(fd.ReporteGestion)
	fd.RemisionResponsables = strings.TrimSpace(fd.RemisionResponsables)
	fd.AsistenciasFecha = strings.TrimSpace(fd.AsistenciasFecha)
	fd.AsistenciasTexto = strings.TrimSpace(fd.AsistenciasTexto)

	switch scope {
	case models.ArticulatingMechanismScopeMunicipal:
		switch fd.AbordadoEnSesion {
		case "si_ordinaria", "si_extraordinaria", "no":
		default:
			return nil, ErrGestionMecanismoInvalid
		}
		if fd.ReporteGestion == "" {
			return nil, ErrGestionMecanismoInvalid
		}
	case models.ArticulatingMechanismScopeDepartamental:
		if fd.FechaProxima == "" || fd.DetalleProxima == "" || fd.ReporteGestion == "" {
			return nil, ErrGestionMecanismoInvalid
		}
	case models.ArticulatingMechanismScopeNacional:
		if fd.ReporteGestion == "" || fd.AsistenciasFecha == "" || fd.AsistenciasTexto == "" {
			return nil, ErrGestionMecanismoInvalid
		}
	default:
		return nil, ErrMechanismInvalidScope
	}
	return &fd, nil
}

func buildGestionMecanismoResult(scope string, fd *gestionMecanismoForm) string {
	parts := make([]string, 0, 6)
	appendPart := func(label, value string) {
		if strings.TrimSpace(value) == "" {
			return
		}
		parts = append(parts, label+": "+strings.TrimSpace(value))
	}
	switch scope {
	case models.ArticulatingMechanismScopeMunicipal:
		label := fd.AbordadoEnSesionTexto
		if label == "" {
			switch fd.AbordadoEnSesion {
			case "si_ordinaria":
				label = "Sí, sesión ordinaria"
			case "si_extraordinaria":
				label = "Sí, sesión extraordinaria"
			case "no":
				label = "No"
			}
		}
		appendPart("Abordado en sesión", label)
		appendPart("Fecha", fd.FechaSesion)
		appendPart("Compromisos", fd.Compromisos)
		appendPart("Próxima sesión", strings.TrimSpace(fd.FechaProxima+" "+fd.DetalleProxima))
		appendPart("Reporte de gestión", fd.ReporteGestion)
		appendPart("Remisión a responsables", fd.RemisionResponsables)
	case models.ArticulatingMechanismScopeDepartamental:
		appendPart("Próxima sesión", strings.TrimSpace(fd.FechaProxima+" "+fd.DetalleProxima))
		appendPart("Reporte de gestión", fd.ReporteGestion)
		appendPart("Remisión a responsables", fd.RemisionResponsables)
	case models.ArticulatingMechanismScopeNacional:
		appendPart("Reporte de gestión", fd.ReporteGestion)
		appendPart("Asistencias técnicas", strings.TrimSpace(fd.AsistenciasFecha+" "+fd.AsistenciasTexto))
	}
	return strings.Join(parts, ". ")
}

func gestionMecanismoDescription(scope string) string {
	switch scope {
	case models.ArticulatingMechanismScopeMunicipal:
		return "Gestión del Comité Municipal"
	case models.ArticulatingMechanismScopeDepartamental:
		return "Gestión del Comité Departamental"
	default:
		return "Gestión de la Instancia Nacional"
	}
}

// SaveGestionMecanismo crea una case_task Done por cada envío del formulario del mecanismo.
func (s *caseTaskService) SaveGestionMecanismo(ctx context.Context, userID, mechanismID string, formData datatypes.JSON) (*models.CaseTask, error) {
	mechanismID = strings.TrimSpace(mechanismID)
	userID = strings.TrimSpace(userID)
	if mechanismID == "" || userID == "" {
		return nil, ErrGestionMecanismoInvalid
	}
	row, err := s.mechanismRepo.FindByID(ctx, mechanismID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMechanismNotFound
		}
		return nil, fmt.Errorf("SaveGestionMecanismo: leer mecanismo: %w", err)
	}
	if row.AssignedUserID != userID {
		return nil, ErrMechanismNotAssigned
	}
	if row.Status != models.CaseArticulatingMechanismStatusActive {
		return nil, ErrMechanismNotActive
	}

	fd, err := parseGestionMecanismoForm(row.Scope, formData)
	if err != nil {
		return nil, err
	}
	formFixed, mErr := json.Marshal(fd)
	if mErr != nil {
		formFixed = formData
	}
	result := buildGestionMecanismoResult(row.Scope, fd)
	now := time.Now()
	mechID := row.ID
	task := &models.CaseTask{
		Category:                    models.CaseTaskCategoryMecanismos,
		Type:                        models.CaseTaskTypeGestionMecanismo,
		Description:                 gestionMecanismoDescription(row.Scope),
		Result:                      &result,
		Status:                      models.CaseTaskStatusDone,
		CompletedAt:                 &now,
		AssignedUserID:              userID,
		CaseID:                      row.CaseID,
		CaseArticulatingMechanismID: &mechID,
		FormData:                    datatypes.JSON(formFixed),
	}
	if err := s.repo.Create(ctx, task); err != nil {
		return nil, fmt.Errorf("SaveGestionMecanismo: crear tarea: %w", err)
	}

	actorName := ""
	if name, nErr := s.mechanismRepo.UserFullName(ctx, userID); nErr == nil {
		actorName = strings.TrimSpace(name)
	}
	event := &models.CaseTimelineEvent{
		CaseID:                  row.CaseID,
		EventType:               models.TimelineEventGestionMecanismo,
		Description:             result,
		ActorID:                 userID,
		ActorName:               actorName,
		Category:                models.TimelineCategoryMecanismos,
		Type:                    models.TimelineTypeGestionMecanismo,
		Icon:                    models.TimelineIconMecanismo,
		Date:                    now,
		EventUserID:             userID,
		Color:                   models.TimelineColorPurple,
		TaskID:                  task.ID,
		ArticulatingMechanismID: row.ID,
		CreatedAt:               now,
	}
	if err := s.timelineRepo.Create(ctx, event); err != nil {
		log.Printf("[CaseTaskService.SaveGestionMecanismo] WARN: timeline caso %s tarea %s: %v", row.CaseID, task.ID, err)
	}

	return task, nil
}

// SaveComentarioEquipoTecnico crea una case_task Done por cada comentario que el
// Equipo Técnico agrega en el tab Derivaciones de Detalle del Caso. El texto va
// en Result; nunca pasa por ToDo.
func (s *caseTaskService) SaveComentarioEquipoTecnico(ctx context.Context, userID, caseID, texto string) (*models.CaseTask, error) {
	userID = strings.TrimSpace(userID)
	caseID = strings.TrimSpace(caseID)
	texto = strings.TrimSpace(texto)
	if userID == "" || caseID == "" || texto == "" {
		return nil, ErrComentarioEquipoTecnicoInvalid
	}

	now := time.Now()
	task := &models.CaseTask{
		Category:       models.CaseTaskCategoryEquipoTecnico,
		Type:           models.CaseTaskTypeComentarioEquipoTecnico,
		Description:    "Comentario del Equipo Técnico",
		Result:         &texto,
		Status:         models.CaseTaskStatusDone,
		CompletedAt:    &now,
		AssignedUserID: userID,
		CaseID:         caseID,
	}
	if err := s.repo.Create(ctx, task); err != nil {
		return nil, fmt.Errorf("SaveComentarioEquipoTecnico: crear tarea: %w", err)
	}

	// UserFullName es una consulta genérica por i_code (vive en mechanismRepo por
	// historia, no por acoplamiento a mecanismos — ver su implementación).
	actorName := ""
	if name, nErr := s.mechanismRepo.UserFullName(ctx, userID); nErr == nil {
		actorName = strings.TrimSpace(name)
	}
	event := &models.CaseTimelineEvent{
		CaseID:      caseID,
		EventType:   models.TimelineEventComentarioEquipoTecnico,
		Description: texto,
		ActorID:     userID,
		ActorName:   actorName,
		Category:    models.TimelineCategoryEquipoTecnico,
		Type:        models.TimelineTypeComentarioEquipoTecnico,
		Icon:        models.TimelineIconComentario,
		Date:        now,
		EventUserID: userID,
		Color:       models.TimelineColorPurple,
		TaskID:      task.ID,
		CreatedAt:   now,
	}
	if err := s.timelineRepo.Create(ctx, event); err != nil {
		log.Printf("[CaseTaskService.SaveComentarioEquipoTecnico] WARN: timeline caso %s tarea %s: %v", caseID, task.ID, err)
	}

	return task, nil
}

// AbrirMesaTecnica crea una case_task ToDo tipo "comite_caso" para que el
// Equipo Técnico analice el caso en mesa técnica, igual que cuando se marca la
// gestión "Alerta por barreras" al guardar un seguimiento (ver
// processFollowUpSubmission en form_service.go). A diferencia de esa, no nace
// de una barrera: BarrierID queda nil. Nace sin dueño (assigned_user_id vacío)
// — pool del Equipo Técnico, la toma y gestiona cualquier 'te' — y se
// completa por el mismo flujo que CompleteWithFormData ya resuelve para
// "comite_caso". justificacion es el motivo por el que se remite el caso a
// mesa técnica; se guarda en Description (mismo patrón que "justificar_remision",
// donde el motivo vive en Description y se muestra como contexto de solo
// lectura al completar la tarea).
func (s *caseTaskService) AbrirMesaTecnica(ctx context.Context, userID, caseID, justificacion string) (*models.CaseTask, error) {
	userID = strings.TrimSpace(userID)
	caseID = strings.TrimSpace(caseID)
	justificacion = strings.TrimSpace(justificacion)
	if caseID == "" || justificacion == "" {
		return nil, ErrAbrirMesaTecnicaInvalid
	}

	task := &models.CaseTask{
		Category:       "Barreras",
		Type:           "comite_caso",
		Description:    justificacion,
		Status:         models.CaseTaskStatusToDo,
		AssignedUserID: "",
		CaseID:         caseID,
	}
	if err := s.repo.Create(ctx, task); err != nil {
		return nil, fmt.Errorf("AbrirMesaTecnica: crear tarea: %w", err)
	}

	actorName := ""
	if name, nErr := s.mechanismRepo.UserFullName(ctx, userID); nErr == nil {
		actorName = strings.TrimSpace(name)
	}
	now := time.Now()
	event := &models.CaseTimelineEvent{
		CaseID:      caseID,
		EventType:   models.TimelineEventMesaTecnicaAbierta,
		Description: justificacion,
		ActorID:     userID,
		ActorName:   actorName,
		Category:    models.TimelineCategoryGeneral,
		Type:        models.TimelineTypeMesaTecnicaAbierta,
		Icon:        "users",
		Date:        now,
		EventUserID: userID,
		Color:       models.TimelineColorPurple,
		TaskID:      task.ID,
		CreatedAt:   now,
	}
	if err := s.timelineRepo.Create(ctx, event); err != nil {
		log.Printf("[CaseTaskService.AbrirMesaTecnica] WARN: timeline caso %s tarea %s: %v", caseID, task.ID, err)
	}

	return task, nil
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
	if v, ok := fd["prioridad"].(string); ok {
		input.Priority = &v
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
			enlaceID, _ := fd["enlaceId"].(string)
			solicitud, _ := fd["solicitudEnlace"].(string)
			solicitud = strings.TrimSpace(solicitud)

			if task.BarrierID != nil && *task.BarrierID != "" {
				// Marca la barrera como enlace activado y le asigna el enlace destino
				// (user_territory_link_id) para que aparezca en su pantalla "Mis Barreras".
				fields := map[string]interface{}{"enlace_activado": true}
				if enlaceID != "" {
					fields["user_territory_link_id"] = enlaceID
				}
				if err := s.barrierRepo.UpdateFields(ctx, *task.BarrierID, fields); err != nil {
					log.Printf("[CaseTaskService] WARN: comite_caso no pudo activar enlace en barrier %s: %v", *task.BarrierID, err)
				}

				// Crea la tarea articulacion_territorial (ToDo) para el enlace destino.
				if enlaceID != "" {
					descripcion := "Articulación territorial"
					if solicitud != "" {
						descripcion = solicitud
					}
					barrierID := *task.BarrierID
					artTask := &models.CaseTask{
						Category:       "Barreras",
						Type:           models.CaseTaskTypeArticulacionTerritorial,
						Description:    descripcion,
						Status:         models.CaseTaskStatusToDo,
						AssignedUserID: enlaceID,
						CaseID:         task.CaseID,
						BarrierID:      &barrierID,
					}
					if err := s.repo.Create(ctx, artTask); err != nil {
						log.Printf("[CaseTaskService] WARN: comite_caso no pudo crear articulacion_territorial para enlace %s: %v", enlaceID, err)
					}
				}
			}

			// Si el comité eligió "otra entidad", asociarla al caso (entity_case).
			if destino, _ := fd["entidadDestino"].(string); destino == "otra" {
				if otraID, ok := fd["otraEntidadId"].(float64); ok && otraID > 0 && s.entityCaseSvc != nil {
					_, err := s.entityCaseSvc.Create(ctx, CreateEntityCaseInput{
						CaseID:         task.CaseID,
						EntityBranchID: int64(otraID),
						CreatedByID:    userId,
					})
					if err != nil && !errors.Is(err, ErrEntityCaseDuplicate) {
						log.Printf("[CaseTaskService] WARN: comite_caso no pudo asociar otra entidad %d al caso %s: %v", int64(otraID), task.CaseID, err)
					}
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
// nueva case_task "validar_remision" para que la profesional psicosocial vuelva a
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
