package service

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

var (
	ErrDirectoryNotFound      = errors.New("directory: registro no encontrado")
	ErrDirectoryInvalidType   = errors.New("directory: tipo inválido")
	ErrDirectoryInvalidInput  = errors.New("directory: datos inválidos")
)

type CreateDirectoryInput struct {
	CityName     string  `json:"cityName"`
	Name         string  `json:"name"`
	Address      string  `json:"address"`
	Phone        *string `json:"phone,omitempty"`
	Email        *string `json:"email,omitempty"`
	OpeningHours *string `json:"openingHours,omitempty"`
	Type         string  `json:"type"`
}

type DirectoryResponse struct {
	ID           string  `json:"id"`
	CityID       string  `json:"cityId"`
	CityName     string  `json:"cityName"`
	Name         string  `json:"name"`
	Address      string  `json:"address"`
	Phone        *string `json:"phone,omitempty"`
	Email        *string `json:"email,omitempty"`
	OpeningHours *string `json:"openingHours,omitempty"`
	Type         string  `json:"type"`
	TypeLabel    string  `json:"typeLabel"`
	CreationDate string  `json:"creationDate,omitempty"`
	UpdateDate   string  `json:"updateDate,omitempty"`
}

type DirectoryListResponse struct {
	CityName  string              `json:"cityName,omitempty"`
	Type      string              `json:"type"`
	TypeLabel string              `json:"typeLabel"`
	Total     int64               `json:"total"`
	Items     []DirectoryResponse `json:"items"`
}

type DirectoryBulkCreateError struct {
	Index    int    `json:"index"`
	CityName string `json:"cityName,omitempty"`
	Name     string `json:"name,omitempty"`
	Error    string `json:"error"`
}

type DirectoryBulkCreateResult struct {
	Total   int                        `json:"total"`
	Created int                        `json:"created"`
	Failed  int                        `json:"failed"`
	Items   []DirectoryResponse        `json:"items"`
	Errors  []DirectoryBulkCreateError `json:"errors,omitempty"`
}

const maxDirectoryBulkSize = 500

type DirectoryService interface {
	Create(ctx context.Context, input CreateDirectoryInput) (*DirectoryResponse, error)
	CreateBulk(ctx context.Context, inputs []CreateDirectoryInput) (*DirectoryBulkCreateResult, error)
	List(ctx context.Context, cityName, dirType string, page, limit int) (*DirectoryListResponse, error)
	GetByID(ctx context.Context, id string) (*DirectoryResponse, error)
}

type directoryService struct {
	dirRepo repository.DirectoryRepository
	locRepo repository.LocationRepository
}

func NewDirectoryService(dirRepo repository.DirectoryRepository, locRepo repository.LocationRepository) DirectoryService {
	return &directoryService{
		dirRepo: dirRepo,
		locRepo: locRepo,
	}
}

func (s *directoryService) Create(ctx context.Context, input CreateDirectoryInput) (*DirectoryResponse, error) {
	cityCache := make(map[string]*models.CityICodeLight)
	return s.createOne(ctx, input, cityCache)
}

func (s *directoryService) CreateBulk(ctx context.Context, inputs []CreateDirectoryInput) (*DirectoryBulkCreateResult, error) {
	if len(inputs) == 0 {
		return nil, fmt.Errorf("%w: el array no puede estar vacío", ErrDirectoryInvalidInput)
	}
	if len(inputs) > maxDirectoryBulkSize {
		return nil, fmt.Errorf("%w: máximo %d directorios por solicitud", ErrDirectoryInvalidInput, maxDirectoryBulkSize)
	}

	result := &DirectoryBulkCreateResult{
		Total: len(inputs),
		Items: make([]DirectoryResponse, 0, len(inputs)),
	}
	cityCache := make(map[string]*models.CityICodeLight)

	for i, input := range inputs {
		created, err := s.createOne(ctx, input, cityCache)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, DirectoryBulkCreateError{
				Index:    i,
				CityName: strings.TrimSpace(input.CityName),
				Name:     strings.TrimSpace(input.Name),
				Error:    bulkDirectoryErrorMessage(err),
			})
			continue
		}
		result.Created++
		result.Items = append(result.Items, *created)
	}

	return result, nil
}

func (s *directoryService) createOne(ctx context.Context, input CreateDirectoryInput, cityCache map[string]*models.CityICodeLight) (*DirectoryResponse, error) {
	if err := validateCreateDirectoryInput(input); err != nil {
		return nil, err
	}

	city, err := s.resolveCity(ctx, input.CityName, cityCache)
	if err != nil {
		return nil, err
	}

	dir := models.Directory{
		CityID:       city.CityICode,
		Name:         strings.TrimSpace(input.Name),
		Address:      strings.TrimSpace(input.Address),
		Phone:        trimOptionalString(input.Phone),
		Email:        trimOptionalString(input.Email),
		OpeningHours: trimOptionalString(input.OpeningHours),
		Type:         strings.ToUpper(strings.TrimSpace(input.Type)),
	}

	if err := s.dirRepo.Create(ctx, &dir); err != nil {
		return nil, err
	}

	resp := toDirectoryResponse(dir, city.CityName)
	return &resp, nil
}

func (s *directoryService) resolveCity(ctx context.Context, cityName string, cache map[string]*models.CityICodeLight) (*models.CityICodeLight, error) {
	trimmed := strings.TrimSpace(cityName)
	cacheKey := strings.ToLower(trimmed)

	if city, ok := cache[cacheKey]; ok {
		return city, nil
	}

	city, err := s.locRepo.FindBestCityMatchByName(ctx, trimmed)
	if err != nil {
		if errors.Is(err, repository.ErrCityNotFound) {
			return nil, fmt.Errorf("%w: %s", repository.ErrCityNotFound, trimmed)
		}
		return nil, err
	}

	cache[cacheKey] = city
	return city, nil
}

func (s *directoryService) List(ctx context.Context, cityName, dirType string, page, limit int) (*DirectoryListResponse, error) {
	dirType = strings.ToUpper(strings.TrimSpace(dirType))
	if dirType == "" {
		return nil, fmt.Errorf("%w: parámetro type requerido", ErrDirectoryInvalidInput)
	}
	if !models.IsValidDirectoryType(dirType) {
		return nil, ErrDirectoryInvalidType
	}

	cityName = strings.TrimSpace(cityName)
	resolvedCityName := ""
	cityICode := ""

	if cityName != "" {
		city, err := s.locRepo.FindBestCityMatchByName(ctx, cityName)
		if err != nil {
			if errors.Is(err, repository.ErrCityNotFound) {
				return nil, fmt.Errorf("%w: %s", repository.ErrCityNotFound, cityName)
			}
			return nil, err
		}
		cityICode = city.CityICode
		resolvedCityName = city.CityName
	}

	items, total, err := s.dirRepo.FindByFilters(ctx, cityICode, dirType, page, limit)
	if err != nil {
		return nil, err
	}

	cityNames, err := s.resolveCityNamesForItems(ctx, items, cityName != "")
	if err != nil {
		return nil, err
	}

	out := make([]DirectoryResponse, len(items))
	for i, item := range items {
		name := cityNames[item.CityID]
		out[i] = toDirectoryResponse(item, name)
	}

	return &DirectoryListResponse{
		CityName:  resolvedCityName,
		Type:      dirType,
		TypeLabel: models.DirectoryTypeLabels[dirType],
		Total:     total,
		Items:     out,
	}, nil
}

func (s *directoryService) resolveCityNamesForItems(ctx context.Context, items []models.Directory, singleCity bool) (map[string]string, error) {
	if singleCity && len(items) > 0 {
		city, err := s.locRepo.FindCityByICode(ctx, items[0].CityID)
		if err != nil {
			return map[string]string{}, nil
		}
		names := make(map[string]string, len(items))
		for _, item := range items {
			names[item.CityID] = city.CityName
		}
		return names, nil
	}

	unique := make(map[string]struct{}, len(items))
	icodes := make([]string, 0, len(items))
	for _, item := range items {
		if _, ok := unique[item.CityID]; ok {
			continue
		}
		unique[item.CityID] = struct{}{}
		icodes = append(icodes, item.CityID)
	}
	return s.locRepo.FindCityNamesByICodes(ctx, icodes)
}

func (s *directoryService) GetByID(ctx context.Context, id string) (*DirectoryResponse, error) {
	dir, err := s.dirRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDirectoryNotFound
		}
		return nil, err
	}

	cityName := ""
	if city, err := s.locRepo.FindCityByICode(ctx, dir.CityID); err == nil {
		cityName = city.CityName
	}

	resp := toDirectoryResponse(*dir, cityName)
	return &resp, nil
}

func validateCreateDirectoryInput(input CreateDirectoryInput) error {
	input.CityName = strings.TrimSpace(input.CityName)
	input.Name = strings.TrimSpace(input.Name)
	input.Address = strings.TrimSpace(input.Address)
	input.Type = strings.ToUpper(strings.TrimSpace(input.Type))

	if input.CityName == "" || input.Name == "" || input.Address == "" || input.Type == "" {
		return fmt.Errorf("%w: campos requeridos cityName, name, address, type", ErrDirectoryInvalidInput)
	}
	if !models.IsValidDirectoryType(input.Type) {
		return ErrDirectoryInvalidType
	}
	return nil
}

func trimOptionalString(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func toDirectoryResponse(dir models.Directory, cityName string) DirectoryResponse {
	resp := DirectoryResponse{
		ID:           dir.ID,
		CityID:       dir.CityID,
		CityName:     cityName,
		Name:         dir.Name,
		Address:      dir.Address,
		Phone:        dir.Phone,
		Email:        dir.Email,
		OpeningHours: dir.OpeningHours,
		Type:         dir.Type,
		TypeLabel:    models.DirectoryTypeLabels[dir.Type],
		CreationDate: dir.CreationDate.UTC().Format("2006-01-02T15:04:05Z"),
		UpdateDate:   dir.UpdateDate.UTC().Format("2006-01-02T15:04:05Z"),
	}
	return resp
}

func bulkDirectoryErrorMessage(err error) string {
	switch {
	case errors.Is(err, repository.ErrCityNotFound):
		return err.Error()
	case errors.Is(err, ErrDirectoryInvalidType):
		return "tipo de directorio inválido"
	case errors.Is(err, ErrDirectoryInvalidInput):
		return err.Error()
	default:
		return "error al guardar el directorio"
	}
}
