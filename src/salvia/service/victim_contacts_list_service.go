package service

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"strings"
)

// VictimContactsListInput parámetros de negocio para el listado de reportes.
type VictimContactsListInput struct {
	FilterKey           string
	FilterValue         string
	FilterEstadoReporte string
	Search              string
	Sort                string
	Order               string
	Page                int
	PageSize            int
}

// VictimContactsListResponse respuesta JSON de GET /api/v1/victim-contacts/list.
type VictimContactsListResponse struct {
	Reports  []models.ReportListItem `json:"reports"`
	Total    int64                   `json:"total"`
	Page     int                     `json:"page"`
	PageSize int                     `json:"pageSize"`
}

// VictimContactsListService contrato del listado de reportes para casos-component.
type VictimContactsListService interface {
	List(ctx context.Context, input VictimContactsListInput) (VictimContactsListResponse, error)
}

type victimContactsListService struct {
	repo repository.VictimContactsListRepository
}

func NewVictimContactsListService(repo repository.VictimContactsListRepository) VictimContactsListService {
	return &victimContactsListService{repo: repo}
}

func (s *victimContactsListService) List(ctx context.Context, input VictimContactsListInput) (VictimContactsListResponse, error) {
	page := input.Page
	if page < 1 {
		page = 1
	}
	pageSize := input.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	result, err := s.repo.List(ctx, repository.VictimContactsListFilters{
		FilterKey:           input.FilterKey,
		FilterValue:         input.FilterValue,
		FilterEstadoReporte: input.FilterEstadoReporte,
		Search:              strings.TrimSpace(input.Search),
		Sort:                input.Sort,
		Order:               input.Order,
		Page:                page,
		PageSize:            pageSize,
	})
	if err != nil {
		return VictimContactsListResponse{}, err
	}

	reports := result.Reports
	if reports == nil {
		reports = []models.ReportListItem{}
	}
	for i := range reports {
		reports[i].ReportType = reportTypeLabel(reports[i].ReportType)
	}

	return VictimContactsListResponse{
		Reports:  reports,
		Total:    result.Total,
		Page:     result.Page,
		PageSize: result.PageSize,
	}, nil
}
