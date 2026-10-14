package service

import (
	"bitsflow/internal/repository"
	"context"
)

// CaseProfessionalService expone el listado consolidado de "Profesionales de
// atención" de un caso (pestaña Derivaciones del detalle del caso).
type CaseProfessionalService interface {
	ListByCase(ctx context.Context, caseID string) ([]repository.CaseProfessionalRow, error)
}

type caseProfessionalService struct {
	repo repository.CaseProfessionalRepository
}

func NewCaseProfessionalService(repo repository.CaseProfessionalRepository) CaseProfessionalService {
	return &caseProfessionalService{repo: repo}
}

func (s *caseProfessionalService) ListByCase(ctx context.Context, caseID string) ([]repository.CaseProfessionalRow, error) {
	return s.repo.FindByCaseID(ctx, caseID)
}
