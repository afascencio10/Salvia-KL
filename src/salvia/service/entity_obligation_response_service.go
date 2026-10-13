package service

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"errors"
	"strings"
)

// Errores del servicio de acciones institucionales de la entidad.
var (
	// ErrEntityObligationEstadoInvalido — el estado general de la respuesta no
	// está en el catálogo.
	ErrEntityObligationEstadoInvalido = errors.New("estado de respuesta inválido")

	// ErrEntityObligationNoAsignada — la acción no está asignada a la sede del
	// usuario. Se devuelve igual cuando la acción no existe y cuando existe pero es
	// de otra entidad, para no revelar la existencia de datos ajenos.
	ErrEntityObligationNoAsignada = errors.New("la acción institucional no está asignada a esta entidad")
)

// EntityObligationListResult es la respuesta de "Mis Acciones Institucionales":
// las filas de la pestaña pedida, los contadores de las dos pestañas y las
// opciones del filtro de ciudad.
type EntityObligationListResult struct {
	EntityName  string                            `json:"entityName"`
	Status      string                            `json:"status"`
	Counts      repository.EntityObligationCounts `json:"counts"`
	Items       []repository.EntityObligationRow  `json:"items"`
	Total       int64                             `json:"total"`
	Page        int                               `json:"page"`
	PageSize    int                               `json:"pageSize"`
	CityOptions []models.EntityCityOption         `json:"cityOptions"`
	// EstadoOptions es el catálogo del select "estado general de la respuesta" del
	// modal. Va en esta misma respuesta y no en un endpoint aparte: el modal se abre
	// desde esta pantalla, que ya hizo esta consulta.
	EstadoOptions []EntityObligationEstadoOption `json:"estadoOptions"`

	// Checklists alimenta el panel "Relación con entidad", indexado por
	// entityCaseId. Va como mapa y no repetido en cada fila porque varias acciones
	// del mismo caso comparten checklist, y repetirlo multiplicaría el payload.
	Checklists map[string]EntityObligationChecklist `json:"checklists"`
}

// EntityObligationChecklist es el checklist de un caso con su conteo ya resuelto,
// para que el frontend no tenga que contarlo.
type EntityObligationChecklist struct {
	Total       int                                        `json:"total"`
	Completadas int                                        `json:"completadas"`
	Items       []repository.EntityObligationChecklistItem `json:"items"`
}

// EntityObligationEstadoOption es una opción del select de estado, en el formato
// { label, value } que ya usan los demás catálogos del frontend.
type EntityObligationEstadoOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// entityObligationEstadoOptions arma el catálogo en un orden fijo. No se recorre
// el mapa directamente porque el orden de un mapa en Go es aleatorio y el select
// cambiaría de orden en cada carga.
func entityObligationEstadoOptions() []EntityObligationEstadoOption {
	orden := []string{
		models.EntityCaseObligationEstadoCompletada,
		models.EntityCaseObligationEstadoEnTramite,
		models.EntityCaseObligationEstadoFueraCompetencia,
		models.EntityCaseObligationEstadoRequiereInfo,
	}
	options := make([]EntityObligationEstadoOption, 0, len(orden))
	for _, code := range orden {
		options = append(options, EntityObligationEstadoOption{
			Label: models.EntityCaseObligationEstadoLabels[code],
			Value: code,
		})
	}
	return options
}

// SaveEntityObligationResponseInput es lo que diligencia la entidad en el modal
// "Completar acción institucional".
type SaveEntityObligationResponseInput struct {
	EntityBranchID     int64
	EntityCaseID       string
	EntityObligationID string
	Radicado           string
	Ampliacion         string
	Estado             string
	RespondedByID      string
}

// EntityObligationResponseService expone la pantalla de acciones institucionales
// de la entidad y el registro de su gestión.
type EntityObligationResponseService interface {
	ListForEntityBranch(ctx context.Context, filter repository.EntityObligationListFilter) (*EntityObligationListResult, error)
	SaveResponse(ctx context.Context, input SaveEntityObligationResponseInput) error
}

type entityObligationResponseService struct {
	repo           repository.EntityCaseObligationResponseRepository
	entityCaseRepo repository.EntityCaseRepository
}

func NewEntityObligationResponseService(
	repo repository.EntityCaseObligationResponseRepository,
	entityCaseRepo repository.EntityCaseRepository,
) EntityObligationResponseService {
	return &entityObligationResponseService{repo: repo, entityCaseRepo: entityCaseRepo}
}

func (s *entityObligationResponseService) ListForEntityBranch(ctx context.Context, filter repository.EntityObligationListFilter) (*EntityObligationListResult, error) {
	// Cualquier valor distinto de "completada" se trata como la pestaña de
	// pendientes, que es la que abre por defecto.
	if filter.Status != models.EntityCaseObligationStatusCompletada {
		filter.Status = models.EntityCaseObligationStatusPendiente
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 10
	}
	if filter.Page < 0 {
		filter.Page = 0
	}

	items, total, err := s.repo.ListForEntityBranch(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Los contadores van sin el filtro de estado: las dos pestañas se muestran a la
	// vez y deben reflejar los mismos filtros de documento y ciudad que el listado.
	counts, err := s.repo.CountForEntityBranch(ctx, filter)
	if err != nil {
		return nil, err
	}

	cities, err := s.repo.ListCaseCitiesForEntityBranch(ctx, filter.EntityBranchID)
	if err != nil {
		return nil, err
	}

	checklists, err := s.buildChecklists(ctx, items)
	if err != nil {
		return nil, err
	}

	result := &EntityObligationListResult{
		Status:        filter.Status,
		Counts:        counts,
		Items:         items,
		Total:         total,
		Page:          filter.Page,
		PageSize:      filter.PageSize,
		CityOptions:   cities,
		EstadoOptions: entityObligationEstadoOptions(),
		Checklists:    checklists,
	}

	// El nombre de la organización es para el encabezado; si falla no vale tumbar
	// la pantalla completa.
	if meta, metaErr := s.entityCaseRepo.GetEntityMetaByBranchID(ctx, filter.EntityBranchID); metaErr == nil && meta != nil {
		result.EntityName = meta.EntityName
	}

	return result, nil
}

// buildChecklists agrupa el checklist por entity_case para las filas de la página
// actual. Los ids se deduplican: varias acciones del mismo caso aparecen como filas
// distintas pero comparten un solo checklist.
func (s *entityObligationResponseService) buildChecklists(ctx context.Context, rows []repository.EntityObligationRow) (map[string]EntityObligationChecklist, error) {
	checklists := make(map[string]EntityObligationChecklist)
	if len(rows) == 0 {
		return checklists, nil
	}

	vistos := make(map[string]bool, len(rows))
	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		if row.EntityCaseID == "" || vistos[row.EntityCaseID] {
			continue
		}
		vistos[row.EntityCaseID] = true
		ids = append(ids, row.EntityCaseID)
	}

	items, err := s.repo.FindChecklistsByEntityCaseIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		checklist := checklists[item.EntityCaseID]
		checklist.Items = append(checklist.Items, item)
		checklist.Total++
		if item.Completada {
			checklist.Completadas++
		}
		checklists[item.EntityCaseID] = checklist
	}

	// Un caso cuyas acciones no devolvieran checklist igual debe tener entrada, para
	// que el frontend no tenga que distinguir entre "sin datos" y "no consultado".
	for _, id := range ids {
		if _, ok := checklists[id]; !ok {
			checklists[id] = EntityObligationChecklist{Items: []repository.EntityObligationChecklistItem{}}
		}
	}

	return checklists, nil
}

func (s *entityObligationResponseService) SaveResponse(ctx context.Context, input SaveEntityObligationResponseInput) error {
	if !models.IsValidEntityCaseObligationEstado(input.Estado) {
		return ErrEntityObligationEstadoInvalido
	}

	asignada, err := s.repo.BelongsToEntityBranch(ctx, input.EntityBranchID, input.EntityCaseID, input.EntityObligationID)
	if err != nil {
		return err
	}
	if !asignada {
		return ErrEntityObligationNoAsignada
	}

	return s.repo.Upsert(ctx, &models.EntityCaseObligationResponse{
		EntityCaseID:       input.EntityCaseID,
		EntityObligationID: input.EntityObligationID,
		Radicado:           optionalText(input.Radicado),
		Ampliacion:         optionalText(input.Ampliacion),
		Estado:             input.Estado,
		RespondedByID:      input.RespondedByID,
	})
}

// optionalText normaliza un campo de texto opcional: en blanco se guarda como
// nulo, no como cadena vacía, para que "sin radicado" se distinga de "radicado
// vacío" al leerlo.
func optionalText(value string) *string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}
