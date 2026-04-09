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

// ── Matriz de riesgo (HU-027) ─────────────────────────────────────────────────
// Días ACUMULADOS desde HOY para cada nivel: [S1, S2, S3, S4]
var riskMatrix = map[int][4]int{
	4: {1, 2, 3, 15},   // Extremo
	3: {1, 3, 15, 30},  // Alto
	2: {2, 15, 30, 45}, // Moderado
	1: {5, 15, 30, 60}, // Bajo
}

// ── Input ─────────────────────────────────────────────────────────────────────

// GenerateCalendarInput es el body de entrada para generar/recalcular el calendario.
type GenerateCalendarInput struct {
	RiskLevel int    `json:"risk_level" binding:"required,min=1,max=4"`
	AgentID   string `json:"agent_id"   binding:"required"`
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
}

// ── Implementación ────────────────────────────────────────────────────────────

type followUpV2Service struct {
	repo repository.FollowUpRepository
}

// NewFollowUpV2Service construye el servicio inyectando el repositorio.
func NewFollowUpV2Service(repo repository.FollowUpRepository) FollowUpV2Service {
	return &followUpV2Service{repo: repo}
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

	completed, err := s.repo.FindCompletedByCaseID(ctx, caseID)
	if err != nil {
		return nil, err
	}

	pending, err := s.repo.FindPendingByCaseID(ctx, caseID)
	if err != nil {
		return nil, err
	}

	riskLevelStr := riskLevelToString(input.RiskLevel)
	today := time.Now().Truncate(24 * time.Hour)

	// ── Caso 1: sin ningún seguimiento → generar S1..S4 completos ────────────
	if len(completed) == 0 && len(pending) == 0 {
		newFollowUps := buildFollowUps(caseID, input, riskLevelStr, offsets[:], today, 1)
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
	numFaltantes := 4 - numCompleted
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
		newFollowUps = buildFollowUps(caseID, input, riskLevelStr, faltantesOffsets, today, startSeq)
		return s.repo.BulkCreate(ctx, tx, newFollowUps)
	}); err != nil {
		return nil, err
	}

	return newFollowUps, nil
}

// ── helpers privados ──────────────────────────────────────────────────────────

func buildFollowUps(caseID string, input GenerateCalendarInput, riskLevelStr string, offsets []int, today time.Time, startSeq int) []models.FollowUpV2 {
	result := make([]models.FollowUpV2, len(offsets))
	for i, days := range offsets {
		result[i] = models.FollowUpV2{
			CaseID:         caseID,
			AgentID:        input.AgentID,
			Team:           input.Team,
			RiskStatus:     riskLevelStr,
			ScheduledDate:  today.AddDate(0, 0, days),
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
