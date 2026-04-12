// Package service contiene la lógica de negocio de la capa salvia (Fase 2).
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

// ── Errores de dominio ────────────────────────────────────────────────────────

// ErrFollowUpNotFound se retorna cuando el registro no existe o fue eliminado.
var ErrFollowUpNotFound = errors.New("followup: registro no encontrado")

// ErrFollowUpCaseEmpty se retorna cuando no hay seguimientos para el caso.
var ErrFollowUpCaseEmpty = errors.New("followup: no se encontraron seguimientos para el caso")

// ── Matriz de riesgo (HU-027) ─────────────────────────────────────────────────
// Días desde HOY para cada nivel. Extremo tiene 5 seguimientos (S1 = mismo día a las 4h).
// Los demás niveles tienen 4 seguimientos.
var riskMatrix = map[int][]int{
	4: {0, 1, 2, 3, 15}, // Extremo — 5 seguimientos (S1=+4h/hoy, S2=+1d, S3=+2d, S4=+3d, S5=+15d)
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
	RegisterFailedAttempt(ctx context.Context, followUpID string, reason string) (*models.FollowUpV2, error)

	// Seguimientos Área
	GetByTeamPaginated(ctx context.Context, team string, filters repository.FollowUpFilters, page, limit int) ([]models.FollowUpV2, int64, error)
	GetAgentWorkload(ctx context.Context, team string, fecha string) ([]repository.AgentWorkload, error)
	GetFilterOptions(ctx context.Context, team string) (FilterOptions, error)
	RescheduleFollowUp(ctx context.Context, id string, input RescheduleInput) error
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
	repo        repository.FollowUpRepository
	barrierRepo repository.BarrierV2Repository
	caseRepo    repository.VictimCaseLightRepository
	townRepo    repository.TownLightRepository
	attemptRepo repository.FollowUpAttemptRepository
	emRepo      repository.EmergencyMeasureRepository
	psRepo      repository.PsychosocialSupportRepository
	esRepo      repository.EconomicStabilizationRepository
	agentRepo   repository.AgentLightRepository
}

// NewFollowUpV2Service construye el servicio inyectando los repositorios.
func NewFollowUpV2Service(
	repo repository.FollowUpRepository,
	barrierRepo repository.BarrierV2Repository,
	caseRepo repository.VictimCaseLightRepository,
	townRepo repository.TownLightRepository,
	attemptRepo repository.FollowUpAttemptRepository,
	emRepo repository.EmergencyMeasureRepository,
	psRepo repository.PsychosocialSupportRepository,
	esRepo repository.EconomicStabilizationRepository,
	agentRepo repository.AgentLightRepository,
) FollowUpV2Service {
	return &followUpV2Service{
		repo:        repo,
		barrierRepo: barrierRepo,
		caseRepo:    caseRepo,
		townRepo:    townRepo,
		attemptRepo: attemptRepo,
		emRepo:      emRepo,
		psRepo:      psRepo,
		esRepo:      esRepo,
		agentRepo:   agentRepo,
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

	// Defaults para asignación diferida
	if input.AgentID == "" {
		input.AgentID = "SIN_ASIGNAR"
	}
	if input.Team == "" {
		input.Team = "SIN_EQUIPO"
	}

	totalExpected := maxFollowUps(input.RiskLevel)

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
		newFollowUps := buildFollowUps(caseID, input, riskLevelStr, offsets, now, today, 1)
		if err := s.repo.RunInTransaction(ctx, func(tx *gorm.DB) error {
			return s.repo.BulkCreate(ctx, tx, newFollowUps)
		}); err != nil {
			return nil, err
		}
		return newFollowUps, nil
	}

	// ── Caso 2: mismo risk_level → sin cambios ────────────────────────────────
	if len(pending) > 0 && pending[0].RiskStatus == riskLevelStr {
		return s.repo.FindByCaseIDOrdered(ctx, caseID)
	}

	// ── Caso 3: risk_level cambió → reprogramar pendientes + generar faltantes ─
	numCompleted := len(completed)
	numFaltantes := totalExpected - numCompleted
	if numFaltantes <= 0 {
		return completed, nil
	}

	faltantesOffsets := offsets[numCompleted : numCompleted+numFaltantes]
	startSeq := numCompleted + 1

	var newFollowUps []models.FollowUpV2
	if err := s.repo.RunInTransaction(ctx, func(tx *gorm.DB) error {
		if err := s.repo.SoftDeleteAndReprogramPending(ctx, tx, caseID); err != nil {
			return err
		}
		newFollowUps = buildFollowUps(caseID, input, riskLevelStr, faltantesOffsets, now, today, startSeq)
		return s.repo.BulkCreate(ctx, tx, newFollowUps)
	}); err != nil {
		return nil, err
	}

	return newFollowUps, nil
}

// GetFollowUpDetail ensambla el modelo de detalle de seguimiento (CSR para carga de pantalla)
func (s *followUpV2Service) GetFollowUpDetail(ctx context.Context, id string, isSupervisor bool) (*models.FollowUpDetailResponse, error) {
	// 1. Obtener los datos base del Seguimiento
	fu, err := s.GetFollowUpByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 2. Obtener información del Caso
	vcase, err := s.caseRepo.FindByID(ctx, fu.CaseID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener el caso: %v", err)
	}

	// 3. Consultar Barreras activas
	barriers, err := s.barrierRepo.FindByFollowUpID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error al obtener las barreras: %v", err)
	}

	// 4. Obtener información del Agente asignado
	if fu.AgentID != "" && fu.AgentID != "SIN_ASIGNAR" {
		agent, err := s.agentRepo.FindByICode(ctx, fu.AgentID)
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

	// 5. Determinar permisos
	perms := models.Permissions{
		CanEdit: isSupervisor,
	}

	return &models.FollowUpDetailResponse{
		FollowUp:               *fu,
		CaseInfo:               *vcase,
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
		if days == 0 && input.RiskLevel == 4 {
			// Extremo S1: programar a las 4 horas desde ahora
			scheduledDate = now.Add(4 * time.Hour)
		} else {
			scheduledDate = today.AddDate(0, 0, days)
		}
		result[i] = models.FollowUpV2{
			CaseID:         caseID,
			AgentID:        input.AgentID,
			Team:           input.Team,
			RiskStatus:     riskLevelStr,
			ScheduledDate:  scheduledDate,
			IsCompleted:    false,
			Status:         models.FollowUpStatusPendiente,
			SequenceNumber: startSeq + i,
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
	// 1. Obtener todos los seguimientos del agente para esa fecha
	followUps, err := s.repo.FindByAgentAndDate(ctx, agentID, date)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("followup: error obteniendo seguimientos del agente: %w", err)
	}

	// 2. Separar en pendientes y priorizados
	for _, fu := range followUps {
		if fu.IsPriority {
			priority = append(priority, fu)
		} else {
			pending = append(pending, fu)
		}
	}

	// 3. Obtener completados del día (REALIZADO o VENCIDO)
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
			resp := models.MyDayFollowUpResponse{
				ID:             fu.ID,
				CaseID:         fu.CaseID,
				RiskStatus:     fu.RiskStatus,
				ScheduledTime:  fu.ScheduledTime,
				Attempts:       fu.Attempts,
				IsPriority:     fu.IsPriority,
				Status:         fu.Status,
				SequenceNumber: fu.SequenceNumber,
				LastAttemptAt:  fu.LastAttemptAt,
			}

			// Obtener datos del caso
			vc, err := s.caseRepo.FindByID(ctx, fu.CaseID)
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
	// Obtener todos los seguimientos del agente (page 0, limit 1000)
	allFollowUps, err := s.repo.FindWithPagination(ctx, 0, 1000)
	if err != nil {
		return nil, err
	}

	// Filtrar por agente, fecha y al menos un intento
	dateOnly := date.Format("2006-01-02")
	var completed []models.FollowUpV2

	for _, fu := range allFollowUps.Items {
		if fu.AgentID == agentID &&
			fu.ScheduledDate.Format("2006-01-02") == dateOnly &&
			fu.Attempts > 0 {
			completed = append(completed, fu)
		}
	}

	return completed, nil
}

// RegisterFailedAttempt registra un intento fallido de contacto
func (s *followUpV2Service) RegisterFailedAttempt(ctx context.Context, followUpID string, reason string) (*models.FollowUpV2, error) {
	// 1. Validar que el seguimiento existe
	fu, err := s.repo.FindByID(ctx, followUpID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrFollowUpNotFound
		}
		return nil, fmt.Errorf("followup: error buscando seguimiento: %w", err)
	}

	// 2. Validar que no haya excedido el máximo de intentos (9)
	if fu.Attempts >= 9 {
		return nil, fmt.Errorf("followup: se alcanzó el máximo de intentos permitidos")
	}

	// 3. Guardar el intento en la tabla follow_up_attempts
	attempt := &models.FollowUpAttempt{
		FollowUpID:  followUpID,
		Reason:      reason,
		WasAnswered: false, // Este método es para intentos fallidos
	}
	if err := s.attemptRepo.CreateAttempt(ctx, attempt); err != nil {
		return nil, fmt.Errorf("followup: error guardando intento: %w", err)
	}

	// 4. Contar el total de intentos reales en la base de datos para ser precisos
	totalAttempts, err := s.attemptRepo.CountByFollowUpID(ctx, followUpID)
	if err != nil {
		return nil, fmt.Errorf("followup: error contando intentos: %w", err)
	}

	// 5. Actualizar el campo attempts y last_attempt_at en follow_up_v2
	// Usamos UTC para almacenamiento consistente, el frontend se encarga de la conversión
	now := time.Now().UTC()
	fu.Attempts = int(totalAttempts)
	fu.LastAttemptAt = &now
	err = s.repo.Update(ctx, fu)
	if err != nil {
		return nil, fmt.Errorf("followup: error actualizando intento en seguimiento: %w", err)
	}

	log.Printf("Intento #%d registrado para seguimiento %s: %s", fu.Attempts, followUpID, reason)

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
	fields := map[string]interface{}{
		"scheduled_date": input.NuevaFecha,
		"status":         models.FollowUpStatusReprogramado,
	}
	if input.NuevaHora != "" {
		fields["scheduled_time"] = input.NuevaHora
	}

	err := s.repo.Reschedule(ctx, id, fields)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrFollowUpNotFound
	}
	return err
}
