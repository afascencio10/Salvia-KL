package service

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"errors"
	"fmt"

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
	models.EntityLetterStateAprobacionJuridica: {models.EntityLetterStateParaRadicar, models.EntityLetterStateEnCorreccion},
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

	// Transición de estado con validación
	UpdateState(ctx context.Context, id string, input UpdateStateInput) (*models.EntityLetter, error)
}

// ─── Implementación ───────────────────────────────────────────────────────────

type entityLetterService struct {
	repo repository.EntityLetterRepository
}

func NewEntityLetterService(repo repository.EntityLetterRepository) EntityLetterService {
	return &entityLetterService{repo: repo}
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

	if err := s.repo.UpdateFields(ctx, id, fields); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEntityLetterNotFound
		}
		return nil, err
	}
	return s.repo.FindByID(ctx, id)
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
