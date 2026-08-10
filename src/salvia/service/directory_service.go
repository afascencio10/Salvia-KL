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
	ErrDirectoryNotFound     = errors.New("directory: registro no encontrado")
	ErrDirectoryInvalidType  = errors.New("directory: tipo inválido")
	ErrDirectorySectorType   = errors.New("directory: sector inválido")
	ErrDirectoryInvalidInput = errors.New("directory: datos inválidos")
)

type CreateDirectoryInput struct {
	CityName     string   `json:"cityName"`
	CityICode    string   `json:"cityICode,omitempty"`
	DepartmentID *uint64  `json:"departmentId,omitempty"`
	TownICode    *string  `json:"townICode,omitempty"`
	Name         string   `json:"name"`
	Address      string   `json:"address"`
	Phone        *string  `json:"phone,omitempty"`
	Email        *string  `json:"email,omitempty"`
	OpeningHours *string  `json:"openingHours,omitempty"`
	Type         string   `json:"type"`
	Sector       *string  `json:"sector,omitempty"`
	Latitude     *float64 `json:"latitude,omitempty"`
	Longitude    *float64 `json:"longitude,omitempty"`
}

type UpdateDirectoryInput struct {
	CityName     *string  `json:"cityName,omitempty"`
	CityICode    *string  `json:"cityICode,omitempty"`
	DepartmentID *uint64  `json:"departmentId,omitempty"`
	TownICode    *string  `json:"townICode,omitempty"`
	Name         *string  `json:"name,omitempty"`
	Address      *string  `json:"address,omitempty"`
	Phone        *string  `json:"phone,omitempty"`
	Email        *string  `json:"email,omitempty"`
	OpeningHours *string  `json:"openingHours,omitempty"`
	Type         *string  `json:"type,omitempty"`
	Sector       *string  `json:"sector,omitempty"`
	Latitude     *float64 `json:"latitude,omitempty"`
	Longitude    *float64 `json:"longitude,omitempty"`
}

type DirectoryResponse struct {
	ID           string   `json:"id"`
	CityID       string   `json:"cityId"`
	CityName     string   `json:"cityName"`
	DepartmentID *uint64  `json:"departmentId,omitempty"`
	TownICode    *string  `json:"townICode,omitempty"`
	Name         string   `json:"name"`
	Address      string   `json:"address"`
	Phone        *string  `json:"phone,omitempty"`
	Email        *string  `json:"email,omitempty"`
	OpeningHours *string  `json:"openingHours,omitempty"`
	Type         string   `json:"type"`
	TypeLabel    string   `json:"typeLabel"`
	Sector       *string  `json:"sector,omitempty"`
	SectorLabel  string   `json:"sectorLabel,omitempty"`
	Latitude     *float64 `json:"latitude,omitempty"`
	Longitude    *float64 `json:"longitude,omitempty"`
	Active       bool     `json:"active"`
	CreationDate string   `json:"creationDate,omitempty"`
	UpdateDate   string   `json:"updateDate,omitempty"`
}

type DirectoryListFilters struct {
	CityName     string
	CityICode    string
	Type         string
	Sector       string
	DepartmentID *uint64
	NameQuery    string
	IncludeInactive bool
}

type DirectoryListResponse struct {
	CityName  string              `json:"cityName,omitempty"`
	Type      string              `json:"type,omitempty"`
	TypeLabel string              `json:"typeLabel,omitempty"`
	Sector    string              `json:"sector,omitempty"`
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
	Update(ctx context.Context, id string, input UpdateDirectoryInput) (*DirectoryResponse, error)
	Disable(ctx context.Context, id string) error
	Enable(ctx context.Context, id string) error
	List(ctx context.Context, filters DirectoryListFilters, page, limit int) (*DirectoryListResponse, error)
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
	if err := validateCreateDirectoryInput(&input); err != nil {
		return nil, err
	}

	cityICode := strings.TrimSpace(input.CityICode)
	cityName := ""
	if cityICode == "" {
		city, err := s.resolveCityByName(ctx, input.CityName, cityCache)
		if err != nil {
			return nil, err
		}
		cityICode = city.CityICode
		cityName = city.CityName
	} else {
		if city, err := s.locRepo.FindCityByICode(ctx, cityICode); err == nil {
			cityName = city.CityName
		}
	}

	dir := models.Directory{
		CityID:       cityICode,
		DepartmentID: input.DepartmentID,
		TownICode:    trimOptionalString(input.TownICode),
		Name:         strings.TrimSpace(input.Name),
		Address:      strings.TrimSpace(input.Address),
		Phone:        trimOptionalString(input.Phone),
		Email:        trimOptionalString(input.Email),
		OpeningHours: trimOptionalString(input.OpeningHours),
		Type:         strings.ToUpper(strings.TrimSpace(input.Type)),
		Sector:       upperOptionalString(input.Sector),
		Latitude:     input.Latitude,
		Longitude:    input.Longitude,
		Active:       true,
	}

	if err := s.dirRepo.Create(ctx, &dir); err != nil {
		return nil, err
	}

	resp := toDirectoryResponse(dir, cityName)
	return &resp, nil
}

func (s *directoryService) Update(ctx context.Context, id string, input UpdateDirectoryInput) (*DirectoryResponse, error) {
	dir, err := s.dirRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDirectoryNotFound
		}
		return nil, err
	}

	cityCache := make(map[string]*models.CityICodeLight)
	cityName := ""

	if input.CityICode != nil && strings.TrimSpace(*input.CityICode) != "" {
		dir.CityID = strings.TrimSpace(*input.CityICode)
		if city, err := s.locRepo.FindCityByICode(ctx, dir.CityID); err == nil {
			cityName = city.CityName
		}
	} else if input.CityName != nil && strings.TrimSpace(*input.CityName) != "" {
		city, err := s.resolveCityByName(ctx, *input.CityName, cityCache)
		if err != nil {
			return nil, err
		}
		dir.CityID = city.CityICode
		cityName = city.CityName
	}

	if input.DepartmentID != nil {
		dir.DepartmentID = input.DepartmentID
	}
	if input.TownICode != nil {
		dir.TownICode = trimOptionalString(input.TownICode)
	}
	if input.Name != nil {
		trimmed := strings.TrimSpace(*input.Name)
		if trimmed == "" {
			return nil, fmt.Errorf("%w: name no puede estar vacío", ErrDirectoryInvalidInput)
		}
		dir.Name = trimmed
	}
	if input.Address != nil {
		trimmed := strings.TrimSpace(*input.Address)
		if trimmed == "" {
			return nil, fmt.Errorf("%w: address no puede estar vacío", ErrDirectoryInvalidInput)
		}
		dir.Address = trimmed
	}
	if input.Phone != nil {
		dir.Phone = trimOptionalString(input.Phone)
	}
	if input.Email != nil {
		dir.Email = trimOptionalString(input.Email)
	}
	if input.OpeningHours != nil {
		dir.OpeningHours = trimOptionalString(input.OpeningHours)
	}
	if input.Type != nil {
		code := strings.ToUpper(strings.TrimSpace(*input.Type))
		if !models.IsValidDirectoryType(code) {
			return nil, ErrDirectoryInvalidType
		}
		dir.Type = code
	}
	if input.Sector != nil {
		code := strings.ToUpper(strings.TrimSpace(*input.Sector))
		if code == "" {
			dir.Sector = nil
		} else {
			if !models.IsValidDirectorySector(code) {
				return nil, ErrDirectorySectorType
			}
			dir.Sector = &code
		}
	}
	if input.Latitude != nil {
		dir.Latitude = input.Latitude
	}
	if input.Longitude != nil {
		dir.Longitude = input.Longitude
	}

	if err := s.dirRepo.Update(ctx, dir); err != nil {
		return nil, err
	}

	if cityName == "" {
		if city, err := s.locRepo.FindCityByICode(ctx, dir.CityID); err == nil {
			cityName = city.CityName
		}
	}

	resp := toDirectoryResponse(*dir, cityName)
	return &resp, nil
}

func (s *directoryService) Disable(ctx context.Context, id string) error {
	if _, err := s.dirRepo.FindByID(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrDirectoryNotFound
		}
		return err
	}
	return s.dirRepo.SetActive(ctx, id, false)
}

func (s *directoryService) Enable(ctx context.Context, id string) error {
	if _, err := s.dirRepo.FindByID(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrDirectoryNotFound
		}
		return err
	}
	return s.dirRepo.SetActive(ctx, id, true)
}

func (s *directoryService) resolveCityByName(ctx context.Context, cityName string, cache map[string]*models.CityICodeLight) (*models.CityICodeLight, error) {
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

func (s *directoryService) List(ctx context.Context, filters DirectoryListFilters, page, limit int) (*DirectoryListResponse, error) {
	dirType := strings.ToUpper(strings.TrimSpace(filters.Type))
	if dirType != "" && !models.IsValidDirectoryType(dirType) {
		return nil, ErrDirectoryInvalidType
	}

	sector := strings.ToUpper(strings.TrimSpace(filters.Sector))
	if sector != "" && !models.IsValidDirectorySector(sector) {
		return nil, ErrDirectorySectorType
	}

	cityICode := strings.TrimSpace(filters.CityICode)
	cityName := strings.TrimSpace(filters.CityName)
	resolvedCityName := ""

	if cityICode == "" && cityName != "" {
		city, err := s.locRepo.FindBestCityMatchByName(ctx, cityName)
		if err != nil {
			if errors.Is(err, repository.ErrCityNotFound) {
				return nil, fmt.Errorf("%w: %s", repository.ErrCityNotFound, cityName)
			}
			return nil, err
		}
		cityICode = city.CityICode
		resolvedCityName = city.CityName
	} else if cityICode != "" {
		if city, err := s.locRepo.FindCityByICode(ctx, cityICode); err == nil {
			resolvedCityName = city.CityName
		}
	}

	repoFilters := repository.DirectoryFilters{
		CityICode:       cityICode,
		Type:            dirType,
		Sector:          sector,
		DepartmentID:    filters.DepartmentID,
		NameQuery:       filters.NameQuery,
		IncludeInactive: filters.IncludeInactive,
	}

	items, total, err := s.dirRepo.FindByFilters(ctx, repoFilters, page, limit)
	if err != nil {
		return nil, err
	}

	cityInfos, err := s.resolveCityInfosForItems(ctx, items)
	if err != nil {
		return nil, err
	}

	out := make([]DirectoryResponse, len(items))
	for i, item := range items {
		info := cityInfos[item.CityID]
		out[i] = toDirectoryResponse(item, info.CityName)
		if out[i].DepartmentID == nil && info.DepartmentId != 0 {
			depID := info.DepartmentId
			out[i].DepartmentID = &depID
		}
	}

	return &DirectoryListResponse{
		CityName:  resolvedCityName,
		Type:      dirType,
		TypeLabel: models.DirectoryTypeLabels[dirType],
		Sector:    sector,
		Total:     total,
		Items:     out,
	}, nil
}

func (s *directoryService) resolveCityInfosForItems(ctx context.Context, items []models.Directory) (map[string]models.CityICodeLight, error) {
	unique := make(map[string]struct{}, len(items))
	icodes := make([]string, 0, len(items))
	for _, item := range items {
		if _, ok := unique[item.CityID]; ok {
			continue
		}
		unique[item.CityID] = struct{}{}
		icodes = append(icodes, item.CityID)
	}
	return s.locRepo.FindCitiesByICodes(ctx, icodes)
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
	var depID *uint64
	if city, err := s.locRepo.FindCityByICode(ctx, dir.CityID); err == nil {
		cityName = city.CityName
		if city.DepartmentId != 0 {
			d := city.DepartmentId
			depID = &d
		}
	}

	resp := toDirectoryResponse(*dir, cityName)
	if resp.DepartmentID == nil {
		resp.DepartmentID = depID
	}
	return &resp, nil
}

func validateCreateDirectoryInput(input *CreateDirectoryInput) error {
	input.CityName = strings.TrimSpace(input.CityName)
	input.CityICode = strings.TrimSpace(input.CityICode)
	input.Name = strings.TrimSpace(input.Name)
	input.Address = strings.TrimSpace(input.Address)
	input.Type = strings.ToUpper(strings.TrimSpace(input.Type))

	if (input.CityName == "" && input.CityICode == "") || input.Name == "" || input.Address == "" || input.Type == "" {
		return fmt.Errorf("%w: campos requeridos (cityName o cityICode), name, address, type", ErrDirectoryInvalidInput)
	}
	if !models.IsValidDirectoryType(input.Type) {
		return ErrDirectoryInvalidType
	}
	if input.Sector != nil {
		code := strings.ToUpper(strings.TrimSpace(*input.Sector))
		if code != "" && !models.IsValidDirectorySector(code) {
			return ErrDirectorySectorType
		}
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

func upperOptionalString(value *string) *string {
	trimmed := trimOptionalString(value)
	if trimmed == nil {
		return nil
	}
	upper := strings.ToUpper(*trimmed)
	return &upper
}

func toDirectoryResponse(dir models.Directory, cityName string) DirectoryResponse {
	resp := DirectoryResponse{
		ID:           dir.ID,
		CityID:       dir.CityID,
		CityName:     cityName,
		DepartmentID: dir.DepartmentID,
		TownICode:    dir.TownICode,
		Name:         dir.Name,
		Address:      dir.Address,
		Phone:        dir.Phone,
		Email:        dir.Email,
		OpeningHours: dir.OpeningHours,
		Type:         dir.Type,
		TypeLabel:    models.DirectoryTypeLabels[dir.Type],
		Sector:       dir.Sector,
		Latitude:     dir.Latitude,
		Longitude:    dir.Longitude,
		Active:       dir.Active,
		CreationDate: dir.CreationDate.UTC().Format("2006-01-02T15:04:05Z"),
		UpdateDate:   dir.UpdateDate.UTC().Format("2006-01-02T15:04:05Z"),
	}
	if dir.Sector != nil {
		resp.SectorLabel = models.DirectorySectorLabels[*dir.Sector]
	}
	return resp
}

func bulkDirectoryErrorMessage(err error) string {
	switch {
	case errors.Is(err, repository.ErrCityNotFound):
		return err.Error()
	case errors.Is(err, ErrDirectoryInvalidType):
		return "tipo de directorio inválido"
	case errors.Is(err, ErrDirectorySectorType):
		return "sector inválido"
	case errors.Is(err, ErrDirectoryInvalidInput):
		return err.Error()
	default:
		return "error al guardar el directorio"
	}
}
