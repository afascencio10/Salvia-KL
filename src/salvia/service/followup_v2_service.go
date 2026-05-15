// Package service contiene la lógica de negocio de la capa salvia (Fase 2).
package service

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	salvia_config "bitsflow/salvia/config"
	"context"
	"errors"
	"fmt"
	"log"
	"time"
	"math"
	"strings"
    "sort"
	"gorm.io/gorm"
)

// ── Errores de dominio ────────────────────────────────────────────────────────

// ErrFollowUpNotFound se retorna cuando el registro no existe o fue eliminado.
var ErrFollowUpNotFound = errors.New("followup: registro no encontrado")

// ErrFollowUpCaseEmpty se retorna cuando no hay seguimientos para el caso.
var ErrFollowUpCaseEmpty = errors.New("followup: no se encontraron seguimientos para el caso")

// ErrFollowUpNotAssigned se retorna cuando el seguimiento no pertenece al agente.
var ErrFollowUpNotAssigned = errors.New("followup: este seguimiento no está asignado a ti")

// ErrFollowUpNotYetDue se retorna cuando la fecha programada aún no ha llegado.
var ErrFollowUpNotYetDue = errors.New("followup: la fecha programada aún no ha llegado")

// ── Matriz de riesgo (HU-027) ─────────────────────────────────────────────────
// Días desde HOY para cada nivel. Extremo tiene 5 seguimientos (S1 = mismo día a las 4h).
// Los demás niveles tienen 4 seguimientos.
var riskMatrix = map[int][]int{
	4: {0, 1, 2, 3, 15, 30}, // Extremo — 6 seguimientos (S1=+4h/hoy, S2=+1d, S3=+2d, S4=+3d, S5=+15d, S6=+30d)
	3: {1, 3, 15, 30},   // Alto    — 4 seguimientos
	2: {2, 15, 30, 45},  // Moderado — 4 seguimientos
	1: {5, 15, 30, 60},  // Bajo    — 4 seguimientos
}

// maxFollowUps retorna la cantidad máxima de seguimientos para un nivel de riesgo.
func maxFollowUps(riskLevel int) int {
	if offsets, ok := riskMatrix[riskLevel]; ok {
		return len(offsets)
	}
	return 4
}

// ── Input ─────────────────────────────────────────────────────────────────────

// GenerateCalendarInput es el body de entrada para generar/recalcular el calendario.
type GenerateCalendarInput struct {
	RiskLevel int    `json:"risk_level" binding:"required,min=1,max=4"`
	AgentID   string `json:"agent_id"`
	Team      string `json:"team"`
}

// ── Interfaz ──────────────────────────────────────────────────────────────────

// LoadFollowUpResult es la respuesta del endpoint hacer-seguimiento al cargar la página.
type LoadFollowUpResult struct {
	FollowUp   *models.FollowUpV2       `json:"followUp"`
	VictimInfo *repository.VictimCaseInfo `json:"victimInfo"`
}

// FollowUpV2Service define el contrato de negocio para FollowUpV2.
type FollowUpV2Service interface {
	// Existentes — NO modificar
	GetFollowUpByID(ctx context.Context, id string) (*models.FollowUpV2, error)
	GetPaginatedFollowUps(ctx context.Context, page, limit int) (repository.PageResult[models.FollowUpV2], error)

	// Nuevos HU-027
	GetByCaseID(ctx context.Context, caseID string) ([]models.FollowUpV2, error)
	GenerateOrRecalculate(ctx context.Context, caseID string, input GenerateCalendarInput) ([]models.FollowUpV2, error)
	GetFollowUpDetail(ctx context.Context, id string, isSupervisor bool) (*models.FollowUpDetailResponse, error)

	// Nuevos para "Mis Seguimientos" - Retornan entidades del dominio
	GetAgentDayFollowUps(ctx context.Context, agentID string, date time.Time) (pending []models.FollowUpV2, priority []models.FollowUpV2, completed []models.FollowUpV2, err error)
	GetMyDayFollowUpsEnriched(ctx context.Context, agentID string, date time.Time) (*models.MyDayResponse, error)
	RegisterContactAttempt(ctx context.Context, followUpID string, reason string, wasAnswered bool) (*models.FollowUpV2, error)

	// Hacer seguimiento
	LoadFollowUp(ctx context.Context, id string, agentID string, formID string) (*LoadFollowUpResult, error)

	// Seguimientos Área
	GetByTeamPaginated(ctx context.Context, team string, filters repository.FollowUpFilters, page, limit int) ([]models.FollowUpV2, int64, error)
	GetAgentWorkload(ctx context.Context, team string, fecha string) ([]repository.AgentWorkload, error)
	GetFilterOptions(ctx context.Context, team string) (FilterOptions, error)
	RescheduleFollowUp(ctx context.Context, id string, input RescheduleInput) error
	CloseCaseFollowUps(ctx context.Context, followUpID string, closureReason string) error
}

// RescheduleInput es el body para reagendar un seguimiento.
type RescheduleInput struct {
	NuevaFecha string `json:"nueva_fecha" binding:"required"`
	NuevaHora  string `json:"nueva_hora"`
	Prioridad  string `json:"prioridad"`
	Motivo     string `json:"motivo" binding:"required"`
}

// FilterOptions contiene los catálogos para los filtros del frontend.
type FilterOptions struct {
	Agentes []repository.AgentOption `json:"agentes"`
	Estados []string                 `json:"estados"`
}

// ── Implementación ────────────────────────────────────────────────────────────

type followUpV2Service struct {
	repo         repository.FollowUpRepository
	fsRepo       repository.FormSubmissionRepository
	barrierRepo  repository.BarrierV2Repository
	caseRepo     repository.VictimCaseLightRepository
	townRepo     repository.TownLightRepository
	attemptRepo  repository.FollowUpAttemptRepository
	emRepo       repository.EmergencyMeasureRepository
	psRepo       repository.PsychosocialSupportRepository
	esRepo       repository.EconomicStabilizationRepository
	agentRepo    repository.AgentLightRepository
	timelineRepo repository.CaseTimelineEventRepository
}

// NewFollowUpV2Service construye el servicio inyectando los repositorios.
func NewFollowUpV2Service(
	repo repository.FollowUpRepository,
	fsRepo repository.FormSubmissionRepository,
	barrierRepo repository.BarrierV2Repository,
	caseRepo repository.VictimCaseLightRepository,
	townRepo repository.TownLightRepository,
	attemptRepo repository.FollowUpAttemptRepository,
	emRepo repository.EmergencyMeasureRepository,
	psRepo repository.PsychosocialSupportRepository,
	esRepo repository.EconomicStabilizationRepository,
	agentRepo repository.AgentLightRepository,
	timelineRepo repository.CaseTimelineEventRepository,
) FollowUpV2Service {
	return &followUpV2Service{
		repo:         repo,
		fsRepo:       fsRepo,
		barrierRepo:  barrierRepo,
		caseRepo:     caseRepo,
		townRepo:     townRepo,
		attemptRepo:  attemptRepo,
		emRepo:       emRepo,
		psRepo:       psRepo,
		esRepo:       esRepo,
		agentRepo:    agentRepo,
		timelineRepo: timelineRepo,
	}
}

// GetFollowUpByID busca un FollowUpV2 por su ID.
// Traduce gorm.ErrRecordNotFound al error de dominio ErrFollowUpNotFound.
func (s *followUpV2Service) GetFollowUpByID(ctx context.Context, id string) (*models.FollowUpV2, error) {
	fu, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrFollowUpNotFound
		}
		return nil, err
	}
	return fu, nil
}

// GetPaginatedFollowUps retorna una página de FollowUpV2.
func (s *followUpV2Service) GetPaginatedFollowUps(ctx context.Context, page, limit int) (repository.PageResult[models.FollowUpV2], error) {
	return s.repo.FindWithPagination(ctx, page, limit)
}

// LoadFollowUp carga toda la información necesaria para la pantalla hacer-seguimiento.
// Valida que el seguimiento pertenezca al agente y que la fecha programada ya llegó.
// Si no tiene formSubmissionId, crea uno y lo asigna.
func (s *followUpV2Service) LoadFollowUp(ctx context.Context, id, agentID, formID string) (*LoadFollowUpResult, error) {
	log.Printf("[SVC] LoadFollowUp → id=%s agentID=%s formID=%s", id, agentID, formID)

	fu, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrFollowUpNotFound
		}
		return nil, err
	}
	log.Printf("[SVC] LoadFollowUp → followUp encontrado: id=%s caseID=%s status=%s agentID=%v", fu.ID, fu.CaseID, fu.Status, fu.AgentID)

	if fu.AgentID == nil || *fu.AgentID != agentID {
		return nil, ErrFollowUpNotAssigned
	}

	today := time.Now().Truncate(24 * time.Hour)
	if fu.ScheduledDate.Truncate(24 * time.Hour).After(today) {
		return nil, ErrFollowUpNotYetDue
	}

	log.Printf("[SVC] LoadFollowUp → llamando LoadVictimInfoByCaseID con caseID=%s", fu.CaseID)
	victimInfo, err := s.repo.LoadVictimInfoByCaseID(ctx, fu.CaseID)
	if err != nil {
		log.Printf("[SVC] LoadFollowUp → ERROR en LoadVictimInfoByCaseID: %v", err)
		return nil, fmt.Errorf("loadFollowUp: leer info víctima: %w", err)
	}
	log.Printf("[SVC] LoadFollowUp → victimInfo raw: %+v", victimInfo)

	if victimInfo != nil {
		locale := salvia_config.Locale["sp"]
		log.Printf("[SVC] LoadFollowUp → resolviendo locale: genderKey=%q → %q | orientationKey=%q → %q",
			victimInfo.GenderIdentity, locale[victimInfo.GenderIdentity],
			victimInfo.SexualOrientation, locale[victimInfo.SexualOrientation])
		victimInfo.GenderIdentity    = locale[victimInfo.GenderIdentity]
		victimInfo.SexualOrientation = locale[victimInfo.SexualOrientation]
	}
	log.Printf("[SVC] LoadFollowUp → victimInfo final: %+v", victimInfo)

	if fu.FormSubmissionID == nil || *fu.FormSubmissionID == "" {
		fs := &models.FormSubmission{FormID: formID}
		if err := s.fsRepo.Create(ctx, fs); err != nil {
			return nil, fmt.Errorf("loadFollowUp: crear form submission: %w", err)
		}
		if err := s.repo.UpdateFormSubmissionID(ctx, fu.ID, fs.ID); err != nil {
			return nil, fmt.Errorf("loadFollowUp: actualizar formSubmissionId: %w", err)
		}
		fu.FormSubmissionID = &fs.ID
	}

	return &LoadFollowUpResult{FollowUp: fu, VictimInfo: victimInfo}, nil
}

// GetByCaseID retorna todos los seguimientos del caso ordenados por fecha ASC.
func (s *followUpV2Service) GetByCaseID(ctx context.Context, caseID string) ([]models.FollowUpV2, error) {
	items, err := s.repo.FindByCaseIDOrdered(ctx, caseID)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, ErrFollowUpCaseEmpty
	}
	return items, nil
}

// GenerateOrRecalculate implementa la lógica central de la HU-027.
//
//  1. Sin seguimientos → genera S1..S4 desde HOY.
//  2. Con seguimientos y mismo risk_level → retorna los existentes sin cambios.
//  3. Con seguimientos y risk_level distinto → reprograma pendientes y genera los faltantes.
func (s *followUpV2Service) GenerateOrRecalculate(ctx context.Context, caseID string, input GenerateCalendarInput) ([]models.FollowUpV2, error) {
	offsets, ok := riskMatrix[input.RiskLevel]
	if !ok {
		return nil, fmt.Errorf("risk_level inválido: %d (debe ser 1-4)", input.RiskLevel)
	}

	// AgentID vacío se persiste como NULL (asignación diferida por supervisor)

	// El equipo se define por el nivel de riesgo, ignorando lo que venga en el input
	if input.RiskLevel >= 3 {
		input.Team = "Riesgo alto"
	} else {
		input.Team = "Riesgo bajo"
	}

	completed, err := s.repo.FindCompletedByCaseID(ctx, caseID)
	if err != nil {
		return nil, err
	}

	pending, err := s.repo.FindPendingByCaseID(ctx, caseID)
	if err != nil {
		return nil, err
	}

	riskLevelStr := riskLevelToString(input.RiskLevel)
	now := time.Now()
	today := now.Truncate(24 * time.Hour)

	// ── Caso 1: sin ningún seguimiento → generar todos ───────────────────────
	if len(completed) == 0 && len(pending) == 0 {

		// [Auto-asignación] Elige el agente con rol "ro" del mismo equipo que tenga
        // el menor promedio de posición de carga en las fechas a generar.
        // Descomentar cuando el evento esté listo para activarse.
        //
        scheduledDates := computeScheduledDates(offsets, now, today, input.RiskLevel)
        assignedAgentID, autoErr := s.calcularAgente(ctx, scheduledDates, input.Team)
        if autoErr != nil {
            log.Printf("[WARN] AutoAsignacion: %v — se usará el agentID del input", autoErr)
        } else {
            input.AgentID = assignedAgentID
        }
		// ───────────────────────────────────────────────────────────────────────

		newFollowUps := buildFollowUps(caseID, input, riskLevelStr, offsets, now, today, 1)
		if err := s.repo.RunInTransaction(ctx, func(tx *gorm.DB) error {
			return s.repo.BulkCreate(ctx, tx, newFollowUps)
		}); err != nil {
			return nil, err
		}

		// Registrar eventos del timeline (no bloquea si falla)
		s.registrarEventosCreacion(ctx, caseID, input.AgentID, riskLevelStr, newFollowUps, now)

		return newFollowUps, nil
	}

	// ── Caso 2: mismo risk_level o distinto → sin cambios ───────────────────
	return s.repo.FindByCaseIDOrdered(ctx, caseID)
}

// GetFollowUpDetail ensambla el modelo de detalle de seguimiento (CSR para carga de pantalla)
func (s *followUpV2Service) GetFollowUpDetail(ctx context.Context, id string, isSupervisor bool) (*models.FollowUpDetailResponse, error) {
	// 1. Obtener los datos base del Seguimiento
	fu, err := s.GetFollowUpByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 2. Obtener información del Caso
	vcase, err := s.caseRepo.FindByICode(ctx, fu.CaseID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener el caso: %v", err)
	}

	// 3. Consultar Barreras activas
	barriers, err := s.barrierRepo.FindByFollowUpID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error al obtener las barreras: %v", err)
	}

	// 4. Obtener información del Agente asignado
	if fu.AgentID != nil && *fu.AgentID != "" && *fu.AgentID != "SIN_ASIGNAR" {
		agent, err := s.agentRepo.FindByICode(ctx, *fu.AgentID)
		if err == nil && agent != nil {
			fu.AgentNames = agent.Names
			fu.AgentLastNames = agent.LastNames
		}
	}

	// 5. Consultar Remisiones (EmergencyMeasure, PsychosocialSupport, EconomicStabilization)
	emergencyMeasures, err := s.emRepo.FindByFollowUpID(ctx, id)
	if err != nil {
		log.Printf("[WARN] Error al obtener medidas de emergencia: %v", err)
	}

	psychosocialSupports, err := s.psRepo.FindByFollowUpID(ctx, id)
	if err != nil {
		log.Printf("[WARN] Error al obtener apoyos psicosociales: %v", err)
	}

	economicStabilizations, err := s.esRepo.FindByFollowUpID(ctx, id)
	if err != nil {
		log.Printf("[WARN] Error al obtener estabilizaciones económicas: %v", err)
	}

	// 6. Cargar info completa de la víctima (teléfono, género, edad, etc.)
	var victimInfo *models.FollowUpVictimInfo
	if vi, err := s.repo.LoadVictimInfoByCaseID(ctx, fu.CaseID); err != nil {
		log.Printf("[WARN] GetFollowUpDetail: no se pudo cargar info víctima: %v", err)
	} else {
		locale := salvia_config.Locale["sp"]
		victimInfo = &models.FollowUpVictimInfo{
			Names:             vi.Names,
			LastNames:         vi.LastNames,
			TownName:          vi.TownName,
			Phone:             vi.Phone,
			GenderIdentity:    locale[vi.GenderIdentity],
			SexualOrientation: locale[vi.SexualOrientation],
			ContactPhone:      vi.ContactPhone,
			Age:               vi.Age,
		}
	}

	// 7. Determinar permisos
	perms := models.Permissions{
		CanEdit: isSupervisor,
	}

	return &models.FollowUpDetailResponse{
		FollowUp:               *fu,
		CaseInfo:               *vcase,
		VictimInfo:             victimInfo,
		Barriers:               barriers,
		Permissions:            perms,
		EmergencyMeasures:      emergencyMeasures,
		PsychosocialSupports:   psychosocialSupports,
		EconomicStabilizations: economicStabilizations,
	}, nil
}

// ── helpers privados ──────────────────────────────────────────────────────────

func buildFollowUps(caseID string, input GenerateCalendarInput, riskLevelStr string, offsets []int, now time.Time, today time.Time, startSeq int) []models.FollowUpV2 {
	result := make([]models.FollowUpV2, len(offsets))
	for i, days := range offsets {
		var scheduledDate time.Time
		var scheduledTime string
		if days == 0 && input.RiskLevel == 4 {
			// Extremo S1: programar a las 4 horas desde ahora
			scheduledDate = now.Add(4 * time.Hour)
			scheduledTime = scheduledDate.Format("15:04:05")
		} else {
			scheduledDate = today.AddDate(0, 0, days)
		}
		riskStr := riskLevelStr
		var agentID *string
		if input.AgentID != "" {
			agentID = &input.AgentID
		}
		result[i] = models.FollowUpV2{
			CaseID:        caseID,
			AgentID:       agentID,
			Team:          input.Team,
			RiskStatus:    &riskStr,
			ScheduledDate: scheduledDate,
			ScheduledTime: scheduledTime,
			Status:        models.FollowUpStatusPendiente,
		}
	}
	return result
}

func riskLevelToString(level int) string {
	switch level {
	case 4:
		return "EXTREMO"
	case 3:
		return "ALTO"
	case 2:
		return "MODERADO"
	default:
		return "BAJO"
	}
}

// GetAgentDayFollowUps obtiene los seguimientos del agente para una fecha específica
// y los clasifica en pendientes, priorizados y realizados.
// Retorna entidades del dominio (NO DTOs).
func (s *followUpV2Service) GetAgentDayFollowUps(ctx context.Context, agentID string, date time.Time) (pending []models.FollowUpV2, priority []models.FollowUpV2, completed []models.FollowUpV2, err error) {
	// 1. Obtener todos los seguimientos del agente para esa fecha (Pendientes y Reprogramados)
	followUps, err := s.repo.FindByAgentAndDate(ctx, agentID, date)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("followup: error obteniendo seguimientos del agente: %w", err)
	}

	// 2. Separar en pendientes (sin hora) y priorizados (con hora)
	for _, fu := range followUps {
		// Un seguimiento es prioritario si tiene una hora programada específica
		if fu.ScheduledTime != "" && fu.ScheduledTime != "00:00:00" {
			priority = append(priority, fu)
		} else {
			pending = append(pending, fu)
		}
	}

	// 3. Obtener completados hoy (Cualquier seguimiento con intento hoy y estado REALIZADO)
	completed, err = s.getCompletedByAgentAndDate(ctx, agentID, date)
	if err != nil {
		log.Printf("Warning: error obteniendo completados: %v", err)
		completed = []models.FollowUpV2{}
	}

	return pending, priority, completed, nil
}

// GetMyDayFollowUpsEnriched obtiene los seguimientos del día enriquecidos con datos del caso
// y retorna la estructura esperada por el frontend
func (s *followUpV2Service) GetMyDayFollowUpsEnriched(ctx context.Context, agentID string, date time.Time) (*models.MyDayResponse, error) {
	// 1. Obtener seguimientos clasificados
	pending, priority, completed, err := s.GetAgentDayFollowUps(ctx, agentID, date)
	if err != nil {
		return nil, fmt.Errorf("followup: error obteniendo seguimientos del día: %w", err)
	}

	// 2. Enriquecer con datos del caso
	enrichFollowUps := func(followUps []models.FollowUpV2) ([]models.MyDayFollowUpResponse, error) {
		var enriched []models.MyDayFollowUpResponse
		for _, fu := range followUps {
			riskStatus := ""
			if fu.RiskStatus != nil {
				riskStatus = *fu.RiskStatus
			}
			resp := models.MyDayFollowUpResponse{
				ID:             fu.ID,
				CaseID:         fu.CaseID,
				RiskStatus:     riskStatus,
				ScheduledTime:  fu.ScheduledTime,
				Attempts:       fu.Attempts,
				IsPriority:     fu.ScheduledTime != "" && fu.ScheduledTime != "00:00:00", // Nueva regla
				Status:             fu.Status,
				SequenceNumber:     fu.SequenceNumber,
				LastAttemptAt:      fu.LastAttemptAt,
				ReassignmentReason: fu.ReassignmentReason,
				Team:               fu.Team,
			}

			// Obtener datos del caso
			vc, err := s.caseRepo.FindByICode(ctx, fu.CaseID)
			if err != nil {
				log.Printf("Warning: no se pudo obtener caso %s: %v", fu.CaseID, err)
				resp.Case = nil
			} else {
				// DEBUG: Verificar town_code del caso
				log.Printf("[Municipio] CaseID=%s TownCode=%s", fu.CaseID, vc.TownCode)

				// Si el caso tiene town_code, obtener el nombre del municipio
				if vc.TownCode != "" && s.townRepo != nil {
					town, err := s.townRepo.FindByCode(ctx, vc.TownCode)
					if err != nil {
						log.Printf("Warning: no se pudo obtener municipio %s: %v", vc.TownCode, err)
						vc.Municipality = vc.TownCode // Fallback: mostrar el código
					} else {
						log.Printf("[Municipio] TownCode=%s -> TownName=%s", vc.TownCode, town.TownName)
						vc.Municipality = town.TownName
					}
				} else if vc.TownCode == "" {
					log.Printf("[Municipio] CaseID=%s no tiene town_code", fu.CaseID)
					vc.Municipality = "N/A"
				}
				resp.Case = vc
			}

			// Obtener historial de intentos
			attempts, err := s.attemptRepo.GetByFollowUpID(ctx, fu.ID)
			if err != nil {
				log.Printf("Warning: no se pudieron obtener intentos para %s: %v", fu.ID, err)
				resp.FollowUpAttempts = []models.FollowUpAttempt{}
			} else {
				resp.FollowUpAttempts = attempts
			}

			enriched = append(enriched, resp)
		}
		return enriched, nil
	}

	enrichedPending, err := enrichFollowUps(pending)
	if err != nil {
		return nil, err
	}

	enrichedPriority, err := enrichFollowUps(priority)
	if err != nil {
		return nil, err
	}

	enrichedCompleted, err := enrichFollowUps(completed)
	if err != nil {
		return nil, err
	}

	// 3. Construir respuesta
	return &models.MyDayResponse{
		Success: "true",
		Data: models.MyDayDataResponse{
			FollowUpsPendingCount:         len(enrichedPending),
			FollowUpsPendingPriorityCount: len(enrichedPriority),
			FollowUpsCompletedCount:       len(enrichedCompleted),
			FollowUpsPending:              enrichedPending,
			FollowUpsPendingPriority:      enrichedPriority,
			FollowUpsCompleted:            enrichedCompleted,
		},
	}, nil
}

// getCompletedByAgentAndDate obtiene seguimientos completados del agente para una fecha
// Un seguimiento se considera "realizado" si tiene al menos un intento registrado hoy
func (s *followUpV2Service) getCompletedByAgentAndDate(ctx context.Context, agentID string, date time.Time) ([]models.FollowUpV2, error) {
	return s.repo.FindRealizedTodayByAgent(ctx, agentID, date)
}

// RegisterContactAttempt registra un intento de contacto (exitoso o fallido)
func (s *followUpV2Service) RegisterContactAttempt(ctx context.Context, followUpID string, reason string, wasAnswered bool) (*models.FollowUpV2, error) {
	// 1. Validar que el seguimiento existe
	fu, err := s.repo.FindByID(ctx, followUpID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrFollowUpNotFound
		}
		return nil, fmt.Errorf("followup: error buscando seguimiento: %w", err)
	}

	// 2. Validar que no haya excedido el máximo de intentos (9) si no contestó
	if !wasAnswered && fu.Attempts >= 9 {
		return nil, fmt.Errorf("followup: se alcanzó el máximo de intentos permitidos")
	}

	// 3. Guardar el intento en la tabla follow_up_attempts
	attempt := &models.FollowUpAttempt{
		FollowUpID:  followUpID,
		Reason:      reason,
		WasAnswered: wasAnswered,
	}
	if err := s.attemptRepo.CreateAttempt(ctx, attempt); err != nil {
		return nil, fmt.Errorf("followup: error guardando intento: %w", err)
	}

	// 4. Contar el total de intentos reales en la base de datos para ser precisos
	totalAttempts, err := s.attemptRepo.CountByFollowUpID(ctx, followUpID)
	if err != nil {
		return nil, fmt.Errorf("followup: error contando intentos: %w", err)
	}

	// 5. Actualizar únicamente el campo attempts y last_attempt_at en follow_up_v2
	// Usamos UpdateFields para evitar problemas de tipos con strings vacíos en PostgreSQL (ej. scheduled_time)
	now := time.Now().UTC()
	err = s.repo.UpdateFields(ctx, followUpID, map[string]interface{}{
		"attempts":        int(totalAttempts),
		"last_attempt_at": now,
	})
	if err != nil {
		return nil, fmt.Errorf("followup: error actualizando intento en seguimiento: %w", err)
	}

	// Sincronizar el objeto local para el retorno
	fu.Attempts = int(totalAttempts)
	fu.LastAttemptAt = &now

	log.Printf("Intento #%d registrado para seguimiento %s: %s", fu.Attempts, followUpID, reason)

	// 6. Si el intento fue fallido (no contestó), registrar en el timeline (HU-027 / Requerimiento adicional)
	if !wasAnswered {
		// Obtener info de la víctima para la descripción
		// fu.CaseID almacena el victim_case_i_code (UUID), no el victim_case_id (entero)
		victimName := "la víctima"
		if vc, vcErr := s.caseRepo.FindByICode(ctx, fu.CaseID); vcErr == nil && vc != nil {
			name := strings.TrimSpace(vc.VictimNames + " " + vc.VictimLastNames)
			if name != "" {
				victimName = name
			}
		}

		agentID := ""
		if fu.AgentID != nil {
			agentID = *fu.AgentID
		}

		timelineEvent := &models.CaseTimelineEvent{
			CaseID:      fu.CaseID,
			Category:    "Seguimientos",
			Type:        "Intento de Seguimiento",
			Icon:        "fa fa-calendar",
			Date:        time.Now(),
			Description: fmt.Sprintf("Llamada realizada sin éxito. Se intentó contactar a %s. Motivo: %s", victimName, reason),
			EventUserID: agentID,
			Color:       "#f8a625",
			FollowUpID:  fu.ID,
		}

		if err := s.repo.CreateTimelineEvent(ctx, timelineEvent); err != nil {
			log.Printf("[WARN] No se pudo crear evento en timeline para intento fallido: %v", err)
			// No retornamos error para no bloquear el registro del intento principal
		}
	}

	return fu, nil
}

// ── Seguimientos Área ─────────────────────────────────────────────────────────

func (s *followUpV2Service) GetByTeamPaginated(ctx context.Context, team string, filters repository.FollowUpFilters, page, limit int) ([]models.FollowUpV2, int64, error) {
	return s.repo.FindByTeamPaginated(ctx, team, filters, page, limit)
}

func (s *followUpV2Service) GetAgentWorkload(ctx context.Context, team string, fecha string) ([]repository.AgentWorkload, error) {
	return s.repo.FindPendingByTeamGroupedByAgent(ctx, team, fecha)
}

func (s *followUpV2Service) GetFilterOptions(ctx context.Context, team string) (FilterOptions, error) {
	agentes, err := s.repo.FindAgentsByTeam(ctx, team)
	if err != nil {
		return FilterOptions{}, err
	}
	return FilterOptions{
		Agentes: agentes,
		Estados: []string{
			models.FollowUpStatusPendiente,
			models.FollowUpStatusRealizado,
			models.FollowUpStatusVencido,
			models.FollowUpStatusReprogramado,
		},
	}, nil
}

func (s *followUpV2Service) RescheduleFollowUp(ctx context.Context, id string, input RescheduleInput) error {
	// 1. Obtener el seguimiento para tener el case_id y agent_id
	fu, err := s.GetFollowUpByID(ctx, id)
	if err != nil {
		return err
	}

	fields := map[string]interface{}{
		"scheduled_date": input.NuevaFecha,
		"status":         models.FollowUpStatusReprogramado,
	}
	if input.NuevaHora != "" {
		fields["scheduled_time"] = input.NuevaHora
	}
	if input.Prioridad == "ALTA" {
		fields["is_priority"] = true
	} else if input.Prioridad == "NORMAL" {
		fields["is_priority"] = false
	}

	err = s.repo.Reschedule(ctx, id, fields)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrFollowUpNotFound
		}
		return err
	}

	// 2. Registrar evento en el timeline (HU-027 / Requerimiento adicional)
	agentID := ""
	if fu.AgentID != nil {
		agentID = *fu.AgentID
	}

	timelineEvent := &models.CaseTimelineEvent{
		CaseID:      fu.CaseID,
		Category:    "Seguimientos",
		Type:        "Seguimiento Pospuesto",
		Icon:        "fa fa-calendar",
		Date:        time.Now(),
		Description: fmt.Sprintf("Se reprogamo el seguimiento para el %s", input.NuevaFecha),
		EventUserID: agentID,
		Color:       "#f8a625",
		FollowUpID:  fu.ID,
	}

	if err := s.repo.CreateTimelineEvent(ctx, timelineEvent); err != nil {
		log.Printf("[WARN] No se pudo crear evento en timeline para seguimiento pospuesto: %v", err)
		// No retornamos error para no bloquear la operación principal
	}

	return nil
}

func (s *followUpV2Service) CloseCaseFollowUps(ctx context.Context, followUpID string, closureReason string) error {
	// 1. Obtener el follow-up para extraer case_id y agent_id
	fu, err := s.repo.FindByID(ctx, followUpID)
	if err != nil {
		return s.repo.CloseCaseFollowUps(ctx, followUpID) // fallback: cerrar sin timeline
	}

	// 2. Obtener nombre de la víctima
	victimName := "la víctima"
	if vc, vcErr := s.caseRepo.FindByICode(ctx, fu.CaseID); vcErr == nil && vc != nil {
		name := strings.TrimSpace(vc.VictimNames + " " + vc.VictimLastNames)
		if name != "" {
			victimName = name
		}
	}

	// 3. Obtener agent_id
	agentID := ""
	if fu.AgentID != nil {
		agentID = *fu.AgentID
	}

	// 4. Determinar motivo de cierre
	motivo := "No se logró contactar a la víctima"
	if strings.TrimSpace(closureReason) != "" {
		motivo = strings.TrimSpace(closureReason)
	}

	// 5. Crear evento de cierre en el timeline
	timelineEvent := &models.CaseTimelineEvent{
		CaseID:      fu.CaseID,
		Category:    "Seguimientos",
		Type:        "Cierre de Caso",
		Icon:        "fa fa-calendar-xmark",
		Date:        time.Now(),
		Description: fmt.Sprintf("Se procede con cierre de caso de %s. Motivo: %s", victimName, motivo),
		EventUserID: agentID,
		Color:       "#d62d20",
		FollowUpID:  fu.ID,
	}

	if err := s.repo.CreateTimelineEvent(ctx, timelineEvent); err != nil {
		log.Printf("[WARN] No se pudo crear evento en timeline para cierre de caso: %v", err)
	}

	// 5. Cerrar los seguimientos del caso
	return s.repo.CloseCaseFollowUps(ctx, followUpID)
}

// ── Timeline de eventos ───────────────────────────────────────────────────────

// registrarEventosCreacion persiste en el timeline:
//   - 1 evento de "Creación de Caso" con el resumen del calendario generado.
//   - 1 evento de "Seguimiento Programado" por cada followup creado.
//
// Los errores se loguean como warnings y no interrumpen el flujo principal.
func (s *followUpV2Service) registrarEventosCreacion(
	ctx context.Context,
	caseID, actorID, riskLevelStr string,
	followUps []models.FollowUpV2,
	now time.Time,
) {
	if s.timelineRepo == nil {
		return
	}

	// Evento del caso
	caseEvent := &models.CaseTimelineEvent{
		CaseID:      caseID,
		EventType:   models.TimelineEventRegistro,
		Category:    "Seguimientos",
		Type:        "Creación de Caso",
		Description: fmt.Sprintf("Caso registrado con %d seguimientos programados (riesgo %s)", len(followUps), riskLevelStr),
		ActorID:     actorID,
		EventUserID: actorID,
		Date:        now,
	}
	if err := s.timelineRepo.Create(ctx, caseEvent); err != nil {
		log.Printf("[WARN] timeline: evento caso %s: %v", caseID, err)
	}

	// Un evento por cada seguimiento creado
	for i, fu := range followUps {
		seq := i + 1
		fecha := fu.ScheduledDate.Format("02/01/2006")
		fuEvent := &models.CaseTimelineEvent{
			CaseID:      caseID,
			FollowUpID:  fu.ID,
			EventType:   models.TimelineEventSeguimiento,
			Category:    "Seguimientos",
			Type:        "Seguimiento Programado",
			Description: fmt.Sprintf("Seguimiento #%d programado para el %s", seq, fecha),
			ActorID:     actorID,
			EventUserID: actorID,
			Date:        now,
		}
		if err := s.timelineRepo.Create(ctx, fuEvent); err != nil {
			log.Printf("[WARN] timeline: evento seguimiento #%d (caso %s): %v", seq, caseID, err)
		}
	}
}

// ── Auto-asignación de agente ─────────────────────────────────────────────────

// computeScheduledDates calcula las fechas absolutas que tendrán los seguimientos
// dado el slice de offsets de la riskMatrix, la referencia de tiempo y el nivel de riesgo.
func computeScheduledDates(offsets []int, now, today time.Time, riskLevel int) []time.Time {
    dates := make([]time.Time, len(offsets))
    for i, days := range offsets {
        if days == 0 && riskLevel == 4 {
            // Extremo S1: +4 h desde ahora
            dates[i] = now.Add(4 * time.Hour)
        } else {
            dates[i] = today.AddDate(0, 0, days)
        }
    }
    return dates
}

// getDateMatrix hace UNA sola query a la BD para obtener la carga de todos los agentes
// en todas las fechas dadas. Retorna map[dateStr(UTC)]map[agentID]count.
// Usar AT TIME ZONE 'UTC' en PostgreSQL garantiza que la fecha se compare siempre en UTC,
// independientemente de la configuración de timezone del servidor.
func (s *followUpV2Service) getDateMatrix(ctx context.Context, dates []time.Time, team string) (map[string]map[string]int64, error) {
	workload, err := s.repo.FindWorkloadByDates(ctx, team, dates)
	if err != nil {
		return nil, fmt.Errorf("getDateMatrix: %w", err)
	}
	matrix := make(map[string]map[string]int64)
	for _, w := range workload {
		if w.AgentID == "" || w.DateStr == "" {
			continue
		}
		if _, ok := matrix[w.DateStr]; !ok {
			matrix[w.DateStr] = make(map[string]int64)
		}
		matrix[w.DateStr][w.AgentID] = w.Total
	}
	return matrix, nil
}

// calcularAgente implementa el algoritmo de auto-asignación balanceada:
//
//  1. Obtiene todos los agentes con rol "ro" del mismo equipo.
//  2. Consulta la carga de todos los agentes en todas las fechas en UNA sola query.
//  3. Por cada fecha asigna posición dense-rank (menor carga = posición 1).
//  4. Calcula avg_pos(Aj) = suma_posiciones / n_fechas.
//  5. Elige el agente con menor avg_pos.
//     Desempate 1: menor carga en las fechas específicas.
//     Desempate 2: menor carga global en el equipo (todos sus seguimientos pendientes).
func (s *followUpV2Service) calcularAgente(ctx context.Context, dates []time.Time, team string) (string, error) {
	if len(dates) == 0 {
		return "", fmt.Errorf("calcularAgente: no se recibieron fechas")
	}

	// 1. Agentes del equipo con rol "ro"
	agents, err := s.agentRepo.FindAllByRoleAndTeam(ctx, "ro", team)
	if err != nil {
		return "", fmt.Errorf("calcularAgente: obtener agentes ro del equipo %q: %w", team, err)
	}
	if len(agents) == 0 {
		return "", fmt.Errorf("calcularAgente: no hay agentes con rol ro en el equipo %q", team)
	}

	// 2. Matriz de carga por fecha en una sola query
	fullMatrix, err := s.getDateMatrix(ctx, dates, team)
	if err != nil {
		return "", err
	}

	// 3. Carga global por agente (tiebreaker final)
	globalRows, err := s.repo.FindGlobalWorkloadByTeam(ctx, team)
	if err != nil {
		log.Printf("[WARN] calcularAgente: carga global no disponible: %v", err)
	}
	globalLoad := make(map[string]int64, len(agents))
	for _, w := range globalRows {
		globalLoad[w.AgentID] = w.Total
	}

	n := len(dates)
	agentPosSum := make(map[string]float64, len(agents))
	agentDateLoad := make(map[string]int64, len(agents))

	// 4. Dense-rank por fecha
	for _, date := range dates {
		dateStr := date.UTC().Format("2006-01-02")
		dayMatrix := fullMatrix[dateStr] // nil si nadie tiene carga ese día

		loadByAgent := make(map[string]int64, len(agents))
		for _, a := range agents {
			loadByAgent[a.ICode] = dayMatrix[a.ICode]
			agentDateLoad[a.ICode] += dayMatrix[a.ICode]
		}

		loadSet := make(map[int64]struct{}, len(agents))
		for _, v := range loadByAgent {
			loadSet[v] = struct{}{}
		}
		sortedLoads := make([]int64, 0, len(loadSet))
		for v := range loadSet {
			sortedLoads = append(sortedLoads, v)
		}
		sort.Slice(sortedLoads, func(i, j int) bool { return sortedLoads[i] < sortedLoads[j] })

		rankOf := make(map[int64]int, len(sortedLoads))
		for i, v := range sortedLoads {
			rankOf[v] = i + 1
		}

		for _, a := range agents {
			agentPosSum[a.ICode] += float64(rankOf[loadByAgent[a.ICode]])
		}

		log.Printf("[AutoAsignacion] Fecha %s | cargas: %v | ranks: %v", dateStr, loadByAgent, rankOf)
	}

	// 5. Seleccionar agente con criterios en cascada
	bestID := ""
	bestAvg := math.MaxFloat64
	bestDateLoad := int64(math.MaxInt64)
	bestGlobalLoad := int64(math.MaxInt64)

	for _, a := range agents {
		avg := agentPosSum[a.ICode] / float64(n)
		dLoad := agentDateLoad[a.ICode]
		gLoad := globalLoad[a.ICode]

		isBetter := avg < bestAvg-1e-9
		isTie1 := math.Abs(avg-bestAvg) < 1e-9 && dLoad < bestDateLoad
		isTie2 := math.Abs(avg-bestAvg) < 1e-9 && dLoad == bestDateLoad && gLoad < bestGlobalLoad
		if isBetter || isTie1 || isTie2 {
			bestID = a.ICode
			bestAvg = avg
			bestDateLoad = dLoad
			bestGlobalLoad = gLoad
		}

		log.Printf("[AutoAsignacion] Agente %s | avg_pos=%.3f | carga_fechas=%d | carga_global=%d",
			a.ICode, avg, dLoad, gLoad)
	}

	log.Printf("[AutoAsignacion] ✓ Seleccionado: %s (equipo=%q, avg_pos=%.3f, carga_fechas=%d, carga_global=%d)",
		bestID, team, bestAvg, bestDateLoad, bestGlobalLoad)

	return bestID, nil
}