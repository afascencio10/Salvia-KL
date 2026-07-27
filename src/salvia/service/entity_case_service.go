package service

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	salvia_config "bitsflow/salvia/config"
	"context"
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"
)

// ErrEntityCaseDuplicate se retorna cuando ya existe una relación activa entre
// el caso y la sede de entidad indicados.
var ErrEntityCaseDuplicate = errors.New("entity_case: la entidad ya está asociada a este caso")

// CreateEntityCaseInput agrupa los datos necesarios para asociar una entidad a un caso.
type CreateEntityCaseInput struct {
	CaseID         string
	EntityBranchID int64
	Objetivo       *string
	CreatedByID    string
}

// EntityCaseListFilter filtros del listado Casos Entidad.
type EntityCaseListFilter struct {
	EntityID int64
	Document string
	City     string
	Page     int
	PageSize int
}

// EntityCaseListResult respuesta paginada del listado.
type EntityCaseListResult struct {
	Items      []models.EntityCaseListItem `json:"items"`
	Total      int64                       `json:"total"`
	Page       int                         `json:"page"`
	PageSize   int                         `json:"pageSize"`
	Sector     string                      `json:"sector,omitempty"`
	SectorName string                      `json:"sectorName,omitempty"`
}

// EntityCaseService define las operaciones de negocio sobre EntityCase
// (componente case-entities y pantalla Casos Entidad).
type EntityCaseService interface {
	ListByCase(ctx context.Context, caseID string) ([]models.EntityCaseWithRelations, error)
	Create(ctx context.Context, input CreateEntityCaseInput) (*models.EntityCase, error)
	ListByEntity(ctx context.Context, filter EntityCaseListFilter) (*EntityCaseListResult, error)
	ListEntities(ctx context.Context) ([]models.EntityCatalogItem, error)
	ListCitiesByEntity(ctx context.Context, entityID int64) ([]models.EntityCityOption, error)
}

type entityCaseService struct {
	repo repository.EntityCaseRepository
}

func NewEntityCaseService(repo repository.EntityCaseRepository) EntityCaseService {
	return &entityCaseService{repo: repo}
}

func (s *entityCaseService) ListByCase(ctx context.Context, caseID string) ([]models.EntityCaseWithRelations, error) {
	return s.repo.FindByCaseIDWithRelations(ctx, caseID)
}

func (s *entityCaseService) Create(ctx context.Context, input CreateEntityCaseInput) (*models.EntityCase, error) {
	exists, err := s.repo.ExistsActive(ctx, input.CaseID, input.EntityBranchID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEntityCaseDuplicate
	}

	ec := &models.EntityCase{
		CaseID:         input.CaseID,
		EntityBranchID: input.EntityBranchID,
		Objetivo:       input.Objetivo,
		CreatedByID:    input.CreatedByID,
	}
	if err := s.repo.Create(ctx, ec); err != nil {
		return nil, err
	}
	return ec, nil
}

func (s *entityCaseService) ListByEntity(ctx context.Context, filter EntityCaseListFilter) (*EntityCaseListResult, error) {
	pageSize := filter.PageSize
	if pageSize <= 0 {
		pageSize = 5
	}
	page := filter.Page
	if page < 0 {
		page = 0
	}

	items, total, err := s.repo.FindByEntityIDPaginated(ctx, filter.EntityID, filter.Document, filter.City, page, pageSize)
	if err != nil {
		return nil, err
	}

	sector := ""
	for i := range items {
		items[i].CaseStatusLabel = victimCaseStatusLabel(items[i].CaseStatus)
		if sector == "" {
			sector = items[i].Sector
		}
	}

	return &EntityCaseListResult{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		Sector:     sector,
		SectorName: entitySectorLabel(sector),
	}, nil
}

func (s *entityCaseService) ListEntities(ctx context.Context) ([]models.EntityCatalogItem, error) {
	items, err := s.repo.ListEntities(ctx)
	if err != nil {
		return nil, err
	}
	for i := range items {
		items[i].SectorName = entitySectorLabel(items[i].Sector)
	}
	return items, nil
}

func (s *entityCaseService) ListCitiesByEntity(ctx context.Context, entityID int64) ([]models.EntityCityOption, error) {
	return s.repo.ListCitiesByEntityID(ctx, entityID)
}

func entitySectorLabel(code string) string {
	key := "sector_" + strings.TrimSpace(code)
	if label, ok := salvia_config.Locale["sp"][key]; ok {
		return label
	}
	return code
}

func victimCaseStatusLabel(code string) string {
	if label, ok := salvia_config.VICTIM_CASE_STATUS["sp"][code]; ok {
		return capitalizeFirst(label)
	}
	return code
}

func capitalizeFirst(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}
	r, size := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[size:]
}
