package service

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"strings"
)

// PsychosocialListInput parámetros del listado de remisiones psicosociales.
type PsychosocialListInput struct {
	FilterProfessionalID      string
	FilterDuplaID             string
	FilterEstadoRemision      string
	FilterSesionesCompletadas string
	FilterEquipoRemitente     string
	FilterNivelRiesgo         string
	FilterNumeroIdentidad     string
	FilterTelefono            string
	Sort                      string
	Order                     string
	Page                      int
	PageSize                  int
}

// PsychosocialListResponse respuesta JSON GET /api/v1/psychosocial-support/list.
type PsychosocialListResponse struct {
	Remisiones []models.PsychosocialListItem `json:"remisiones"`
	Total      int64                         `json:"total"`
	Page       int                           `json:"page"`
	PageSize   int                           `json:"pageSize"`
}

// DuplasListResponse respuesta JSON GET /api/v1/duplas.
type DuplasListResponse struct {
	Duplas []models.DuplaOption `json:"duplas"`
}

// EquiposRemitentesResponse respuesta JSON GET /api/v1/psychosocial-support/equipos-remitentes.
type EquiposRemitentesResponse struct {
	Teams []string `json:"teams"`
}

// PsychosocialListService contrato E-01.
type PsychosocialListService interface {
	List(ctx context.Context, input PsychosocialListInput) (PsychosocialListResponse, error)
	Stats(ctx context.Context, input PsychosocialListInput) (models.PsychosocialListStats, error)
	ListDuplas(ctx context.Context) (DuplasListResponse, error)
	ListEquiposRemitentes(ctx context.Context) (EquiposRemitentesResponse, error)
}

type psychosocialListService struct {
	listRepo  repository.PsychosocialListRepository
	duplaRepo repository.DuplaRepository
}

func NewPsychosocialListService(
	listRepo repository.PsychosocialListRepository,
	duplaRepo repository.DuplaRepository,
) PsychosocialListService {
	return &psychosocialListService{
		listRepo:  listRepo,
		duplaRepo: duplaRepo,
	}
}

func (s *psychosocialListService) List(ctx context.Context, input PsychosocialListInput) (PsychosocialListResponse, error) {
	filters := s.toFilters(input)
	result, err := s.listRepo.List(ctx, filters)
	if err != nil {
		return PsychosocialListResponse{}, err
	}
	remisiones := result.Remisiones
	if remisiones == nil {
		remisiones = []models.PsychosocialListItem{}
	}
	return PsychosocialListResponse{
		Remisiones: remisiones,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
	}, nil
}

func (s *psychosocialListService) Stats(ctx context.Context, input PsychosocialListInput) (models.PsychosocialListStats, error) {
	return s.listRepo.Stats(ctx, s.toFilters(input))
}

func (s *psychosocialListService) ListDuplas(ctx context.Context) (DuplasListResponse, error) {
	duplas, err := s.duplaRepo.ListActive(ctx)
	if err != nil {
		return DuplasListResponse{}, err
	}
	if duplas == nil {
		duplas = []models.DuplaOption{}
	}
	return DuplasListResponse{Duplas: duplas}, nil
}

func (s *psychosocialListService) ListEquiposRemitentes(ctx context.Context) (EquiposRemitentesResponse, error) {
	teams, err := s.listRepo.ListEquiposRemitentes(ctx)
	if err != nil {
		return EquiposRemitentesResponse{}, err
	}
	if teams == nil {
		teams = []string{}
	}
	return EquiposRemitentesResponse{Teams: teams}, nil
}

func (s *psychosocialListService) toFilters(input PsychosocialListInput) repository.PsychosocialListFilters {
	page := input.Page
	if page < 1 {
		page = 1
	}
	pageSize := input.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	sort := input.Sort
	if sort == "" {
		sort = "created_at"
	}
	order := input.Order
	if order != "asc" {
		order = "desc"
	}
	return repository.PsychosocialListFilters{
		FilterProfessionalID:      strings.TrimSpace(input.FilterProfessionalID),
		FilterDuplaID:             strings.TrimSpace(input.FilterDuplaID),
		FilterEstadoRemision:      strings.TrimSpace(input.FilterEstadoRemision),
		FilterSesionesCompletadas: strings.TrimSpace(input.FilterSesionesCompletadas),
		FilterEquipoRemitente:     strings.TrimSpace(input.FilterEquipoRemitente),
		FilterNivelRiesgo:         strings.TrimSpace(input.FilterNivelRiesgo),
		FilterNumeroIdentidad:     strings.TrimSpace(input.FilterNumeroIdentidad),
		FilterTelefono:            strings.TrimSpace(input.FilterTelefono),
		Sort:                      sort,
		Order:                     order,
		Page:                      page,
		PageSize:                  pageSize,
	}
}
