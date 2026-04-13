// Package service contiene la lógica de negocio de la capa salvia (Fase 2).
package service

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"errors"
	"fmt"
	"time"

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
	4: {0, 1, 2, 3, 15},  // Extremo — 5 seguimientos (S1=+4h/hoy, S2=+1d, S3=+2d, S4=+3d, S5=+15d)
	3: {1, 3, 15, 30},    // Alto    — 4 seguimientos
	2: {2, 15, 30, 45},   // Moderado — 4 seguimientos
	1: {5, 15, 30, 60},   // Bajo    — 4 seguimientos
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

	// Hacer seguimiento
	LoadFollowUp(ctx context.Context, id string, agentID string, formID string) (*LoadFollowUpResult, error)

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
	repo   repository.FollowUpRepository
	fsRepo repository.FormSubmissionRepository
}

// NewFollowUpV2Service construye el servicio inyectando los repositorios.
func NewFollowUpV2Service(repo repository.FollowUpRepository, fsRepo repository.FormSubmissionRepository) FollowUpV2Service {
	return &followUpV2Service{repo: repo, fsRepo: fsRepo}
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
	fu, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrFollowUpNotFound
		}
		return nil, err
	}

	if fu.AgentID != agentID {
		return nil, ErrFollowUpNotAssigned
	}

	today := time.Now().Truncate(24 * time.Hour)
	if fu.ScheduledDate.Truncate(24 * time.Hour).After(today) {
		return nil, ErrFollowUpNotYetDue
	}

	victimInfo, err := s.repo.LoadVictimInfoByCaseID(ctx, fu.CaseID)
	if err != nil {
		return nil, fmt.Errorf("loadFollowUp: leer info víctima: %w", err)
	}

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
	if len(pending) > 0 && pending[0].RiskStatus != nil && *pending[0].RiskStatus == riskLevelStr {
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
		riskStr := riskLevelStr
		result[i] = models.FollowUpV2{
			CaseID:        caseID,
			AgentID:       input.AgentID,
			Team:          input.Team,
			RiskStatus:    &riskStr,
			ScheduledDate: scheduledDate,
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
		fields["scheduled_date"] = input.NuevaFecha + " " + input.NuevaHora
	}
	err := s.repo.Reschedule(ctx, id, fields)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrFollowUpNotFound
	}
	return err
}
