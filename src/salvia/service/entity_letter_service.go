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

// ─── Errores de dominio ───────────────────────────────────────────────────────

var (
	ErrEntityLetterNotFound    = errors.New("entity_letter: registro no encontrado")
	ErrEntityLetterInvalidState = errors.New("entity_letter: transición de estado no permitida")
)

// ─── Inputs ───────────────────────────────────────────────────────────────────

// CreateEntityLetterInput contiene los campos requeridos para crear un oficio.
type CreateEntityLetterInput struct {
	BarrierID          string
	CaseID             string
	AgentID            *string
	NotificationUserID *string
}

// UpdateEntityLetterInput permite actualizar campos opcionales del oficio.
// Solo se aplica el campo si el puntero no es nil.
type UpdateEntityLetterInput struct {
	AgentID            *string
	NotificationUserID *string
	ReviewBy           *string
	RadicadoBy         *string
	RegisterBy         *string
	Entidad            *string
	Nivel              *string
	UrlKofax           *string
}

// ActionInput contiene los datos enviados desde cualquier modal de gestión.
// Cada acción solo utiliza los campos que le corresponden; los demás se ignoran.
//
// Acciones soportadas:
//
//	"proyectar"   — modal Proyectar oficio  (por_proyectar → para_revisar)
//	"revisar"     — modal Revisar oficio    (para_revisar  → aprobacion_juridica)
//	"por_corregir"— modal Revisar oficio    (para_revisar  → en_correccion)
//	"radicar"     — modal Radicar oficio    (aprobacion_juridica → radicado)
type ActionInput struct {
	Action string  // nombre de la acción
	UserID string  // ID del usuario que ejecuta la acción

	// Campos del modal "Proyectar oficio"
	Nivel             *string
	Entidad           *string
	UrlKofax          *string
	Priority          *string // "normal" | "alta"
	EntityBranchID    *int64
	EntityName        *string // nombre legible: entity_branch.name o texto libre "otra entidad"
	DepartmentID      *string
	CityID            *string
	TownID            *string
	OfficialDependency *string
	Subject           *string

	// Campos del modal "Radicar oficio"
	AsuntoRadicado *string
	CorreoEntidad  *string
	NumeroRadicado *string

	// Campos del modal "Revisar oficio" — acción por_corregir
	ReasonCorrection *string

	// Campos del modal "Registrar respuesta"
	ResponseDate     *string // fecha en formato "YYYY-MM-DD" — se parsea a time.Time
	CorreoRemitente  *string
	AsuntoRespuesta  *string
	ResponseReviewBy *string

	// TaskDescription — solo se llena cuando la acción viene de completar una
	// CaseTask vía case-task-modal (acciones "proyectar" y "corregir" disparadas
	// desde CompleteWithFormData). Si está presente, sobreescribe el texto
	// genérico por-estado del CaseTimelineEvent con "Tarea completada: {desc}".
	// Para el resto de acciones (revisar, por_corregir, radicar,
	// registrar_respuesta) queda nil y el evento conserva su texto genérico.
	TaskDescription *string
}

// UpdateStateInput contiene el nuevo estado y quién realiza la transición.
type UpdateStateInput struct {
	State  string
	UserID string // ID del usuario que ejecuta la acción (para auditoría)
}

// ─── Transiciones de estado válidas ──────────────────────────────────────────

// validTransitions define los estados destino permitidos desde cada estado origen.
var validTransitions = map[string][]string{
	models.EntityLetterStatePorProyectar:       {models.EntityLetterStateParaRevisar},
	models.EntityLetterStateParaRevisar:        {models.EntityLetterStateAprobacionJuridica, models.EntityLetterStateEnCorreccion},
	models.EntityLetterStateEnCorreccion:       {models.EntityLetterStateParaRevisar},
	models.EntityLetterStateAprobacionJuridica: {models.EntityLetterStateParaRadicar, models.EntityLetterStateEnCorreccion, models.EntityLetterStateRadicado},
	models.EntityLetterStateParaRadicar:        {models.EntityLetterStateRadicado},
	models.EntityLetterStateRadicado:           {models.EntityLetterStateRespondido},
}

func isValidTransition(from, to string) bool {
	allowed, ok := validTransitions[from]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}

// ─── Interface ────────────────────────────────────────────────────────────────

type EntityLetterService interface {
	// CRUD
	Create(ctx context.Context, input CreateEntityLetterInput) (*models.EntityLetter, error)
	GetByID(ctx context.Context, id string) (*models.EntityLetter, error)
	List(ctx context.Context, page, limit int) (repository.PageResult[models.EntityLetter], error)
	Update(ctx context.Context, id string, input UpdateEntityLetterInput) (*models.EntityLetter, error)
	Delete(ctx context.Context, id string) error

	// Consultas específicas
	ListByCase(ctx context.Context, caseID string) ([]models.EntityLetter, error)
	ListByBarrier(ctx context.Context, barrierID string) ([]models.EntityLetter, error)
	ListByState(ctx context.Context, state string, page, limit int) (repository.PageResult[models.EntityLetter], error)
	ListByAgent(ctx context.Context, agentID string) ([]models.EntityLetter, error)
	ListByNotificationUser(ctx context.Context, notificationUserID string) ([]models.EntityLetter, error)

	// Consultas enriquecidas con relaciones (barrier_v2 + victim_case)
	ListByAgentWithRelations(ctx context.Context, agentID string) ([]models.EntityLetterWithRelations, error)
	ListByNotificationUserWithRelations(ctx context.Context, notifUserID string) ([]models.EntityLetterWithRelations, error)

	// ListWithRelationsFiltered devuelve oficios paginados con filtros en BD y conteo de pendientes.
	ListWithRelationsFiltered(ctx context.Context, filter repository.EntityLetterListFilter, page, limit int) (repository.EntityLetterListResult, error)

	// Transición de estado con validación
	UpdateState(ctx context.Context, id string, input UpdateStateInput) (*models.EntityLetter, error)

	// PerformAction ejecuta la acción de un modal: actualiza los campos relevantes
	// y realiza la transición de estado en una sola operación atómica.
	PerformAction(ctx context.Context, id string, input ActionInput) (*models.EntityLetter, error)
}

// ─── Implementación ───────────────────────────────────────────────────────────

type entityLetterService struct {
	repo         repository.EntityLetterRepository
	timelineRepo repository.CaseTimelineEventRepository
	taskRepo     repository.CaseTaskRepository
	barrierRepo  repository.BarrierV2Repository
}

func NewEntityLetterService(
	repo repository.EntityLetterRepository,
	timelineRepo repository.CaseTimelineEventRepository,
	taskRepo repository.CaseTaskRepository,
	barrierRepo repository.BarrierV2Repository,
) EntityLetterService {
	return &entityLetterService{
		repo:         repo,
		timelineRepo: timelineRepo,
		taskRepo:     taskRepo,
		barrierRepo:  barrierRepo,
	}
}

func (s *entityLetterService) Create(ctx context.Context, input CreateEntityLetterInput) (*models.EntityLetter, error) {
	letter := &models.EntityLetter{
		BarrierID:          input.BarrierID,
		CaseID:             input.CaseID,
		State:              models.EntityLetterStatePorProyectar,
		AgentID:            input.AgentID,
		NotificationUserID: input.NotificationUserID,
	}
	if err := s.repo.Create(ctx, letter); err != nil {
		return nil, fmt.Errorf("entity_letter: crear: %w", err)
	}
	return letter, nil
}

func (s *entityLetterService) GetByID(ctx context.Context, id string) (*models.EntityLetter, error) {
	letter, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEntityLetterNotFound
		}
		return nil, err
	}
	return letter, nil
}

func (s *entityLetterService) List(ctx context.Context, page, limit int) (repository.PageResult[models.EntityLetter], error) {
	return s.repo.FindWithPagination(ctx, page, limit)
}

func (s *entityLetterService) Update(ctx context.Context, id string, input UpdateEntityLetterInput) (*models.EntityLetter, error) {
	fields := map[string]interface{}{}
	if input.AgentID != nil            { fields["agent_id"] = *input.AgentID }
	if input.NotificationUserID != nil { fields["notification_user_id"] = *input.NotificationUserID }
	if input.ReviewBy != nil           { fields["review_by"] = *input.ReviewBy }
	if input.RadicadoBy != nil         { fields["radicado_by"] = *input.RadicadoBy }
	if input.RegisterBy != nil         { fields["register_by"] = *input.RegisterBy }
	if input.Entidad != nil            { fields["entidad"] = *input.Entidad }
	if input.Nivel != nil              { fields["nivel"] = *input.Nivel }
	if input.UrlKofax != nil           { fields["url_kofax"] = *input.UrlKofax }

	if err := s.repo.UpdateFields(ctx, id, fields); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEntityLetterNotFound
		}
		return nil, err
	}
	return s.repo.FindByID(ctx, id)
}

// PerformAction ejecuta la acción del modal correspondiente:
// valida el estado actual, actualiza los campos propios de la acción
// y realiza la transición de estado en una sola operación.
func (s *entityLetterService) PerformAction(ctx context.Context, id string, input ActionInput) (*models.EntityLetter, error) {
	letter, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	fields := map[string]interface{}{}

	switch input.Action {
	case "proyectar":
		// por_proyectar → para_revisar
		if letter.State != models.EntityLetterStatePorProyectar {
			return nil, fmt.Errorf("%w: acción 'proyectar' requiere estado '%s', estado actual: '%s'",
				ErrEntityLetterInvalidState, models.EntityLetterStatePorProyectar, letter.State)
		}
		if input.DepartmentID == nil || *input.DepartmentID == "" {
			return nil, fmt.Errorf("entity_letter: el campo 'departmentId' es requerido para proyectar")
		}
		if input.CityID == nil || *input.CityID == "" {
			return nil, fmt.Errorf("entity_letter: el campo 'cityId' es requerido para proyectar")
		}
		if input.TownID == nil || *input.TownID == "" {
			return nil, fmt.Errorf("entity_letter: el campo 'townId' es requerido para proyectar")
		}
		if input.EntityName == nil || *input.EntityName == "" {
			return nil, fmt.Errorf("entity_letter: el campo 'entityName' es requerido para proyectar")
		}
		if input.OfficialDependency == nil || *input.OfficialDependency == "" {
			return nil, fmt.Errorf("entity_letter: el campo 'officialDependency' es requerido para proyectar")
		}
		if input.Subject == nil || *input.Subject == "" {
			return nil, fmt.Errorf("entity_letter: el campo 'subject' es requerido para proyectar")
		}
		if input.UrlKofax == nil || *input.UrlKofax == "" {
			return nil, fmt.Errorf("entity_letter: el campo 'urlKofax' es requerido para proyectar")
		}
		fields["department_id"]       = *input.DepartmentID
		fields["city_id"]             = *input.CityID
		fields["town_id"]             = *input.TownID
		fields["entidad"]             = *input.EntityName
		fields["official_dependency"] = *input.OfficialDependency
		fields["subject"]             = *input.Subject
		fields["url_kofax"]           = *input.UrlKofax
		if input.EntityBranchID != nil {
			fields["entity_branch_id"] = *input.EntityBranchID
		}
		if input.Priority != nil && *input.Priority != "" {
			fields["priority"] = *input.Priority
		}
		if input.UserID != "" {
			fields["register_by"] = input.UserID
		}
		fields["state"] = models.EntityLetterStateParaRevisar
		log.Printf("[DEBUG] proyectar: id=%s dept=%s city=%s town=%s entity=%q dep=%q subj=%q",
			id, *input.DepartmentID, *input.CityID, *input.TownID,
			*input.EntityName, *input.OfficialDependency, *input.Subject)

	case "revisar":
		// para_revisar → aprobacion_juridica
		if letter.State != models.EntityLetterStateParaRevisar {
			return nil, fmt.Errorf("%w: acción 'revisar' requiere estado '%s', estado actual: '%s'",
				ErrEntityLetterInvalidState, models.EntityLetterStateParaRevisar, letter.State)
		}
		if input.UserID != "" {
			fields["review_by"] = input.UserID
		}
		fields["state"] = models.EntityLetterStateAprobacionJuridica

	case "por_corregir":
		// para_revisar → en_correccion (agente de notificaciones)
		// aprobacion_juridica → en_correccion (abogado devuelve por inconsistencia)
		if letter.State != models.EntityLetterStateParaRevisar &&
			letter.State != models.EntityLetterStateAprobacionJuridica {
			return nil, fmt.Errorf("%w: acción 'por_corregir' requiere estado '%s' o '%s', estado actual: '%s'",
				ErrEntityLetterInvalidState, models.EntityLetterStateParaRevisar,
				models.EntityLetterStateAprobacionJuridica, letter.State)
		}
		if input.ReasonCorrection == nil || *input.ReasonCorrection == "" {
			return nil, fmt.Errorf("entity_letter: el campo 'reasonCorrection' es requerido para marcar por corregir")
		}
		fields["reason_correction"] = *input.ReasonCorrection
		fields["state"] = models.EntityLetterStateEnCorreccion

	case "radicar":
		// aprobacion_juridica → radicado (el paso "para_radicar" lo omitimos mientras solo existen roles op/an)
		if letter.State != models.EntityLetterStateAprobacionJuridica {
			return nil, fmt.Errorf("%w: acción 'radicar' requiere estado '%s', estado actual: '%s'",
				ErrEntityLetterInvalidState, models.EntityLetterStateAprobacionJuridica, letter.State)
		}
		if input.AsuntoRadicado == nil || *input.AsuntoRadicado == "" {
			return nil, fmt.Errorf("entity_letter: el campo 'asuntoRadicado' es requerido para radicar")
		}
		if input.CorreoEntidad == nil || *input.CorreoEntidad == "" {
			return nil, fmt.Errorf("entity_letter: el campo 'correoEntidad' es requerido para radicar")
		}
		if input.NumeroRadicado == nil || *input.NumeroRadicado == "" {
			return nil, fmt.Errorf("entity_letter: el campo 'numeroRadicado' es requerido para radicar")
		}
		fields["asunto_radicado"] = *input.AsuntoRadicado
		fields["correo_entidad"]  = *input.CorreoEntidad
		fields["numero_radicado"] = *input.NumeroRadicado
		if input.UserID != "" {
			fields["radicado_by"] = input.UserID
		}
		fields["state"] = models.EntityLetterStateRadicado

	case "corregir":
		// en_correccion → para_revisar
		if letter.State != models.EntityLetterStateEnCorreccion {
			return nil, fmt.Errorf("%w: acción 'corregir' requiere estado '%s', estado actual: '%s'",
				ErrEntityLetterInvalidState, models.EntityLetterStateEnCorreccion, letter.State)
		}
		fields["state"] = models.EntityLetterStateParaRevisar

	case "registrar_respuesta":
		// radicado → respondido
		if letter.State != models.EntityLetterStateRadicado {
			return nil, fmt.Errorf("%w: acción 'registrar_respuesta' requiere estado '%s', estado actual: '%s'",
				ErrEntityLetterInvalidState, models.EntityLetterStateRadicado, letter.State)
		}
		if input.ResponseDate == nil || *input.ResponseDate == "" {
			return nil, fmt.Errorf("entity_letter: el campo 'responseDate' es requerido para registrar respuesta")
		}
		if input.CorreoRemitente == nil || *input.CorreoRemitente == "" {
			return nil, fmt.Errorf("entity_letter: el campo 'correoRemitente' es requerido para registrar respuesta")
		}
		if input.AsuntoRespuesta == nil || *input.AsuntoRespuesta == "" {
			return nil, fmt.Errorf("entity_letter: el campo 'asuntoRespuesta' es requerido para registrar respuesta")
		}
		if input.ResponseReviewBy == nil || *input.ResponseReviewBy == "" {
			return nil, fmt.Errorf("entity_letter: el campo 'responseReviewBy' es requerido para registrar respuesta")
		}
		parsedDate, err := time.Parse("2006-01-02", *input.ResponseDate)
		if err != nil {
			return nil, fmt.Errorf("entity_letter: formato de fecha inválido para 'responseDate' (esperado YYYY-MM-DD): %w", err)
		}
		fields["correo_remitente"]   = *input.CorreoRemitente
		fields["asunto_respuesta"]   = *input.AsuntoRespuesta
		fields["response_review_by"] = *input.ResponseReviewBy
		fields["response_date"]      = parsedDate
		fields["state"]              = models.EntityLetterStateRespondido

	default:
		return nil, fmt.Errorf("entity_letter: acción desconocida: %s", input.Action)
	}

	if err := s.repo.UpdateFields(ctx, id, fields); err != nil {
		return nil, err
	}

	updated, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Al proyectar o corregir: completar la CaseTask pendiente asociada al oficio (fire-and-forget)
	if (input.Action == "proyectar" || input.Action == "corregir") && s.taskRepo != nil {
		s.completarCaseTask(ctx, id)
	}

	// Al proyectar: si la barrera vinculada está OPEN, pasarla a "En Gestion" (fire-and-forget)
	if input.Action == "proyectar" && s.barrierRepo != nil {
		s.actualizarBarreraEnGestion(ctx, updated)
	}

	// Al marcar por corregir: crear nueva CaseTask para el agente de seguimiento (fire-and-forget)
	if input.Action == "por_corregir" && s.taskRepo != nil {
		s.crearCaseTaskCorreccion(ctx, updated, input.ReasonCorrection)
	}

	// Registrar evento en el timeline (fire-and-forget: no bloquea si falla)
	s.registrarEventoOficio(ctx, updated, input.UserID, input.ReasonCorrection, input.TaskDescription)

	return updated, nil
}

// completarCaseTask busca la CaseTask en estado "ToDo" vinculada al EntityLetter
// y la marca como "Done" con la fecha actual. Los errores se loguean sin interrumpir el flujo.
func (s *entityLetterService) completarCaseTask(ctx context.Context, entityLetterID string) {
	task, err := s.taskRepo.FindTodoByEntityLetterID(ctx, entityLetterID)
	if err != nil {
		log.Printf("[WARN] entity_letter: no se pudo buscar CaseTask para oficio %s: %v", entityLetterID, err)
		return
	}
	if task == nil {
		return
	}

	now := time.Now()
	updateFields := map[string]interface{}{
		"status":       models.CaseTaskStatusDone,
		"completed_at": now,
	}
	if err := s.taskRepo.UpdateFields(ctx, task.ID, updateFields); err != nil {
		log.Printf("[WARN] entity_letter: no se pudo completar CaseTask %s (oficio %s): %v",
			task.ID, entityLetterID, err)
	}
}

// actualizarBarreraEnGestion cambia el status de barrier_v2 a "En Gestion" cuando el oficio
// se proyecta y la barrera asociada aún está en OPEN. Los errores se loguean sin interrumpir el flujo.
func (s *entityLetterService) actualizarBarreraEnGestion(ctx context.Context, letter *models.EntityLetter) {
	if letter == nil || letter.BarrierID == "" {
		return
	}

	barrier, err := s.barrierRepo.FindByID(ctx, letter.BarrierID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return
		}
		log.Printf("[WARN] entity_letter: no se pudo leer barrier_v2 %s (oficio %s): %v",
			letter.BarrierID, letter.ID, err)
		return
	}
	if barrier.Status != models.BarrierV2StatusOpen {
		return
	}

	if err := s.barrierRepo.UpdateStatus(ctx, letter.BarrierID, models.BarrierV2StatusEnGestion); err != nil {
		log.Printf("[WARN] entity_letter: no se pudo actualizar barrier_v2 %s a En Gestion (oficio %s): %v",
			letter.BarrierID, letter.ID, err)
	}
}

// crearCaseTaskCorreccion crea una CaseTask de tipo "Corregir oficio" asignada al agente
// de seguimiento del EntityLetter cuando el oficio es devuelto para corrección.
// Los errores se loguean sin interrumpir el flujo principal.
func (s *entityLetterService) crearCaseTaskCorreccion(ctx context.Context, letter *models.EntityLetter, reason *string) {
	if s.taskRepo == nil {
		return
	}

	assignedUserID := ""
	if letter.AgentID != nil {
		assignedUserID = *letter.AgentID
	}
	if assignedUserID == "" {
		log.Printf("[WARN] entity_letter: no se puede crear CaseTask de corrección para oficio %s: agentId vacío", letter.ID)
		return
	}

	desc := "Corregir oficio devuelto para revisión"
	if reason != nil && *reason != "" {
		desc = fmt.Sprintf("Corregir oficio — Razón de devolución: %s", *reason)
	}

	barrierID := letter.BarrierID
	letterID  := letter.ID

	task := &models.CaseTask{
		Category:       models.TimelineCategoryBarreras,
		Type:           "Corregir oficio",
		Description:    desc,
		Status:         models.CaseTaskStatusToDo,
		AssignedUserID: assignedUserID,
		CaseID:         letter.CaseID,
		BarrierID:      &barrierID,
		EntityLetterID: &letterID,
	}

	if err := s.taskRepo.Create(ctx, task); err != nil {
		log.Printf("[WARN] entity_letter: no se pudo crear CaseTask de corrección para oficio %s: %v", letter.ID, err)
	}
}

// registrarEventoOficio persiste un CaseTimelineEvent cada vez que un EntityLetter
// cambia de estado. Los errores se loguean como warnings sin interrumpir el flujo.
func (s *entityLetterService) registrarEventoOficio(
	ctx context.Context,
	letter *models.EntityLetter,
	userID string,
	reasonCorrection *string,
	taskDescription *string,
) {
	if s.timelineRepo == nil {
		return
	}

	type stateConfig struct {
		eventType   string
		description string
		icon        string
		color       string
	}

	cfgByState := map[string]stateConfig{
		models.EntityLetterStateParaRevisar: {
			eventType:   models.TimelineTypeOficioParaRevisar,
			description: "Oficio enviado a revisión",
			icon:        models.TimelineIconOficioRevisar,
			color:       models.TimelineColorBlue,
		},
		models.EntityLetterStateEnCorreccion: {
			eventType:   models.TimelineTypeOficioEnCorreccion,
			description: "Oficio devuelto para corrección",
			icon:        models.TimelineIconOficioCorreccion,
			color:       models.TimelineColorOrange,
		},
		models.EntityLetterStateAprobacionJuridica: {
			eventType:   models.TimelineTypeOficioAprobacionJuridica,
			description: "Oficio en aprobación jurídica",
			icon:        models.TimelineIconOficioJuridica,
			color:       models.TimelineColorPurple,
		},
		models.EntityLetterStateParaRadicar: {
			eventType:   models.TimelineTypeOficioParaRadicar,
			description: "Oficio listo para radicar",
			icon:        models.TimelineIconOficioRadicar,
			color:       models.TimelineColorTeal,
		},
		models.EntityLetterStateRadicado: {
			eventType:   models.TimelineTypeOficioRadicado,
			description: "Oficio radicado ante la entidad",
			icon:        models.TimelineIconOficioRadicar,
			color:       models.TimelineColorGreen,
		},
		models.EntityLetterStateRespondido: {
			eventType:   models.TimelineTypeOficioRespondido,
			description: "Oficio respondido por la entidad",
			icon:        models.TimelineIconOficioRespondido,
			color:       models.TimelineColorGreen,
		},
	}

	cfg, ok := cfgByState[letter.State]
	if !ok {
		return
	}

	// Si la acción fue por_corregir e incluye razón, enriquecer la descripción
	if letter.State == models.EntityLetterStateEnCorreccion &&
		reasonCorrection != nil && *reasonCorrection != "" {
		cfg.description = fmt.Sprintf("%s — Razón: %s", cfg.description, *reasonCorrection)
	}

	// Si esta transición vino de completar una CaseTask (proyectar/corregir vía
	// case-task-modal), el evento debe reflejar la tarea completada, no el
	// texto genérico del estado.
	if taskDescription != nil && *taskDescription != "" {
		cfg.description = "Tarea completada: " + *taskDescription
	}

	now := time.Now()
	event := &models.CaseTimelineEvent{
		CaseID:         letter.CaseID,
		EventType:      models.TimelineEventOficioActualizado,
		Category:       models.TimelineCategoryBarreras,
		Type:           cfg.eventType,
		Icon:           cfg.icon,
		Date:           now,
		Description:    cfg.description,
		EventUserID:    userID,
		Color:          cfg.color,
		EntityLetterID: letter.ID,
		BarrierID:      letter.BarrierID,
		CreatedAt:      now,
	}

	if err := s.timelineRepo.Create(ctx, event); err != nil {
		log.Printf("[WARN] entity_letter: no se pudo registrar evento en timeline (id=%s estado=%s): %v",
			letter.ID, letter.State, err)
	}
}

func (s *entityLetterService) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrEntityLetterNotFound
	}
	return err
}

func (s *entityLetterService) ListByCase(ctx context.Context, caseID string) ([]models.EntityLetter, error) {
	return s.repo.FindByCaseID(ctx, caseID)
}

func (s *entityLetterService) ListByBarrier(ctx context.Context, barrierID string) ([]models.EntityLetter, error) {
	return s.repo.FindByBarrierID(ctx, barrierID)
}

func (s *entityLetterService) ListByState(ctx context.Context, state string, page, limit int) (repository.PageResult[models.EntityLetter], error) {
	return s.repo.FindByState(ctx, state, page, limit)
}

func (s *entityLetterService) ListByAgent(ctx context.Context, agentID string) ([]models.EntityLetter, error) {
	return s.repo.FindByAgentID(ctx, agentID)
}

func (s *entityLetterService) ListByNotificationUser(ctx context.Context, notificationUserID string) ([]models.EntityLetter, error) {
	return s.repo.FindByNotificationUserID(ctx, notificationUserID)
}

func (s *entityLetterService) ListByAgentWithRelations(ctx context.Context, agentID string) ([]models.EntityLetterWithRelations, error) {
	return s.repo.FindByAgentIDWithRelations(ctx, agentID)
}

func (s *entityLetterService) ListByNotificationUserWithRelations(ctx context.Context, notifUserID string) ([]models.EntityLetterWithRelations, error) {
	return s.repo.FindByNotificationUserIDWithRelations(ctx, notifUserID)
}

// ManageableStatesForAgent devuelve los estados que el rol op/ro puede gestionar.
func ManageableStatesForAgent() []string {
	return []string{
		models.EntityLetterStatePorProyectar,
		models.EntityLetterStateEnCorreccion,
	}
}

// ManageableStatesForNotificationUser devuelve los estados que el rol an puede gestionar.
func ManageableStatesForNotificationUser() []string {
	return []string{
		models.EntityLetterStateParaRevisar,
		models.EntityLetterStateAprobacionJuridica,
		models.EntityLetterStateParaRadicar,
		models.EntityLetterStateRadicado,
	}
}

func (s *entityLetterService) ListWithRelationsFiltered(
	ctx context.Context,
	filter repository.EntityLetterListFilter,
	page, limit int,
) (repository.EntityLetterListResult, error) {
	if filter.AgentID == "" && filter.NotificationUserID == "" {
		return repository.EntityLetterListResult{}, fmt.Errorf("entity_letter: se requiere agentId o notificationUserId")
	}
	if filter.AgentID != "" && len(filter.ManageableStates) == 0 {
		filter.ManageableStates = ManageableStatesForAgent()
	}
	if filter.NotificationUserID != "" && len(filter.ManageableStates) == 0 {
		filter.ManageableStates = ManageableStatesForNotificationUser()
	}
	return s.repo.FindWithRelationsFilteredPaginated(ctx, filter, page, limit)
}

// UpdateState valida que la transición sea permitida y actualiza el estado.
// Además asigna automáticamente el campo de auditoría correspondiente según el estado destino.
func (s *entityLetterService) UpdateState(ctx context.Context, id string, input UpdateStateInput) (*models.EntityLetter, error) {
	letter, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !isValidTransition(letter.State, input.State) {
		return nil, fmt.Errorf("%w: %s → %s", ErrEntityLetterInvalidState, letter.State, input.State)
	}

	if err := s.repo.UpdateState(ctx, id, input.State); err != nil {
		return nil, err
	}

	// Asignar campo de auditoría según el estado destino
	auditFields := map[string]interface{}{}
	switch input.State {
	case models.EntityLetterStateParaRevisar:
		if input.UserID != "" {
			auditFields["agent_id"] = input.UserID
		}
	case models.EntityLetterStateAprobacionJuridica:
		if input.UserID != "" {
			auditFields["review_by"] = input.UserID
		}
	case models.EntityLetterStateRadicado:
		if input.UserID != "" {
			auditFields["radicado_by"] = input.UserID
		}
	case models.EntityLetterStateRespondido:
		if input.UserID != "" {
			auditFields["register_by"] = input.UserID
		}
	}

	if len(auditFields) > 0 {
		if err := s.repo.UpdateFields(ctx, id, auditFields); err != nil {
			return nil, err
		}
	}

	return s.repo.FindByID(ctx, id)
}
