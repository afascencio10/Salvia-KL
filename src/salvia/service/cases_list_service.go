package service

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"strings"
)

// CasesListInput parámetros de negocio para el listado de casos.
type CasesListInput struct {
	FilterKey           string
	FilterValue         string
	ChipFilter          string
	DropdownFilterKey   string
	DropdownFilterValue string
	Search      string
	Sort        string
	Order       string
	Page        int
	PageSize    int
}

// CasesListResponse respuesta JSON de GET /api/v1/cases/list.
type CasesListResponse struct {
	Cases    []models.CaseListItem `json:"cases"`
	Total    int64                 `json:"total"`
	Page     int                   `json:"page"`
	PageSize int                   `json:"pageSize"`
}

// CasesListService contrato del listado de casos para casos-component.
type CasesListService interface {
	List(ctx context.Context, input CasesListInput) (CasesListResponse, error)
}

type casesListService struct {
	repo repository.CasesListRepository
}

func NewCasesListService(repo repository.CasesListRepository) CasesListService {
	return &casesListService{repo: repo}
}

func (s *casesListService) List(ctx context.Context, input CasesListInput) (CasesListResponse, error) {
	sort := input.Sort
	if sort != "next_follow_up" {
		sort = "registration_date"
	}

	order := input.Order
	if order != "asc" {
		order = "desc"
	}

	page := input.Page
	if page < 1 {
		page = 1
	}

	pageSize := input.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	result, err := s.repo.List(ctx, repository.CasesListFilters{
		FilterKey:           input.FilterKey,
		FilterValue:         input.FilterValue,
		ChipFilter:          input.ChipFilter,
		DropdownFilterKey:   input.DropdownFilterKey,
		DropdownFilterValue: input.DropdownFilterValue,
		Search:              strings.TrimSpace(input.Search),
		Sort:        sort,
		Order:       order,
		Page:        page,
		PageSize:    pageSize,
	})
	if err != nil {
		return CasesListResponse{}, err
	}

	cases := result.Cases
	if cases == nil {
		cases = []models.CaseListItem{}
	}

	return CasesListResponse{
		Cases:    cases,
		Total:    result.Total,
		Page:     result.Page,
		PageSize: result.PageSize,
	}, nil
}
