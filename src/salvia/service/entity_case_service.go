package service

import (
	commondb "bitsflow/common/db"
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	salvia_config "bitsflow/salvia/config"
	"context"
	"errors"
	"log"
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
	Solicitud      *string
	Radicado       *string
	CreatedByID    string
	// ObligationIds — IDs del catálogo entity_obligation seleccionados en el
	// checklist del formulario. Se registran como entity_case_obligation en
	// estado "pendiente" al crear la entidad.
	ObligationIds []string
}

// EntityCaseListFilter filtros del listado Casos Entidad (por sede de sesión).
type EntityCaseListFilter struct {
	EntityBranchID int64
	Document       string
	Page           int
	PageSize       int
}

// EntityCaseListResult respuesta paginada del listado.
type EntityCaseListResult struct {
	Items          []models.EntityCaseListItem `json:"items"`
	Total          int64                       `json:"total"`
	Page           int                         `json:"page"`
	PageSize       int                         `json:"pageSize"`
	EntityBranchID int64                       `json:"entityBranchId,omitempty"`
	EntityName     string                      `json:"entityName,omitempty"`
	Sector         string                      `json:"sector,omitempty"`
	SectorName     string                      `json:"sectorName,omitempty"`
}

// EntityCaseService define las operaciones de negocio sobre EntityCase
// (componente case-entities y pantalla Casos Entidad).
type EntityCaseService interface {
	ListByCase(ctx context.Context, caseID string) ([]models.EntityCaseWithRelations, error)
	Create(ctx context.Context, input CreateEntityCaseInput) (*models.EntityCase, error)
	ListByEntity(ctx context.Context, filter EntityCaseListFilter) (*EntityCaseListResult, error)
	ListEntities(ctx context.Context) ([]models.EntityCatalogItem, error)
	ListCitiesByEntity(ctx context.Context, entityID int64) ([]models.EntityCityOption, error)

	// ListObligations lista las obligaciones de un entity_case, con labels
	// resueltos y el nombre de la organización. status vacío = todas.
	ListObligations(ctx context.Context, entityCaseID, status string) (*models.EntityCaseObligationsResponse, error)

	// UpdateLastAction sobreescribe entity_case.last_action con un resumen de
	// texto libre. Usado por processFollowUpSubmission al registrar un
	// entity_case_follow_up (Sección "Seguimiento a Entidades").
	UpdateLastAction(ctx context.Context, entityCaseID, lastAction string) error

	// GetObligacionesResumen agrega el checklist de obligaciones de todas las
	// entidades de un caso — usado por case-gestion-institucional.
	GetObligacionesResumen(ctx context.Context, caseID string) (*models.ObligacionesResumenResponse, error)

	// GetEntityCaseDetail devuelve el detalle completo de un entity_case por ID.
	GetEntityCaseDetail(ctx context.Context, entityCaseID string) (*models.EntityCaseDetailRow, error)
}

type entityCaseService struct {
	repo           repository.EntityCaseRepository
	obligationRepo repository.EntityCaseObligationRepository
	responseRepo   repository.EntityCaseObligationResponseRepository
}

func NewEntityCaseService(
	repo repository.EntityCaseRepository,
	obligationRepo repository.EntityCaseObligationRepository,
	responseRepo repository.EntityCaseObligationResponseRepository,
) EntityCaseService {
	return &entityCaseService{repo: repo, obligationRepo: obligationRepo, responseRepo: responseRepo}
}

func (s *entityCaseService) ListByCase(ctx context.Context, caseID string) ([]models.EntityCaseWithRelations, error) {
	items, err := s.repo.FindByCaseIDWithRelations(ctx, caseID)
	if err != nil || len(items) == 0 {
		return items, err
	}
	// Si falla lo de las respuestas, las tarjetas se devuelven igual sin ese bloque:
	// es información añadida, y tumbar toda la pestaña por ella sería peor.
	_ = s.attachEntityResponses(ctx, caseID, items)
	return items, nil
}

// attachEntityResponses llena RespuestaEntidad leyendo la misma tabla que "Mis
// Acciones Institucionales". Son dos consultas para toda la grilla, no una por
// tarjeta.
func (s *entityCaseService) attachEntityResponses(ctx context.Context, caseID string, items []models.EntityCaseWithRelations) error {
	if s.responseRepo == nil {
		return nil
	}

	ids := make([]string, 0, len(items))
	for _, item := range items {
		if item.RelID != "" {
			ids = append(ids, item.RelID)
		}
	}

	// Los conteos salen del checklist y no de contar respuestas, para que la tarjeta
	// y la pantalla de la entidad usen el mismo criterio de "completada".
	checklists, err := s.responseRepo.FindChecklistsByEntityCaseIDs(ctx, ids)
	if err != nil {
		return err
	}
	respuestas, err := s.responseRepo.FindByCaseID(ctx, caseID)
	if err != nil {
		return err
	}

	totales := make(map[string]int, len(ids))
	respondidas := make(map[string]int, len(ids))
	labels := make(map[string]string, len(checklists))
	for _, accion := range checklists {
		totales[accion.EntityCaseID]++
		if accion.Completada {
			respondidas[accion.EntityCaseID]++
		}
		labels[accion.EntityCaseID+"|"+accion.EntityObligationID] = accion.Label
	}

	// La más reciente por entidad se decide comparando fechas, no confiando en el
	// orden en que vengan las filas.
	ultimas := make(map[string]models.EntityCaseObligationResponse, len(ids))
	for _, respuesta := range respuestas {
		anterior, ok := ultimas[respuesta.EntityCaseID]
		if !ok || respuesta.UpdatedAt.After(anterior.UpdatedAt) {
			ultimas[respuesta.EntityCaseID] = respuesta
		}
	}

	for i := range items {
		relID := items[i].RelID
		ultima, respondio := ultimas[relID]
		if totales[relID] == 0 && !respondio {
			continue // sin acciones institucionales: la tarjeta no muestra el bloque
		}

		resumen := &models.EntityCaseObligationResponseSummary{
			Total:       totales[relID],
			Respondidas: respondidas[relID],
		}
		if respondio {
			respondedAt := ultima.UpdatedAt
			resumen.AccionLabel = labels[relID+"|"+ultima.EntityObligationID]
			resumen.Estado = ultima.Estado
			resumen.EstadoLabel = models.EntityCaseObligationEstadoLabels[ultima.Estado]
			resumen.RespondedAt = &respondedAt
			if ultima.Radicado != nil {
				resumen.Radicado = *ultima.Radicado
			}
			if ultima.Ampliacion != nil {
				resumen.Ampliacion = *ultima.Ampliacion
			}
		}
		items[i].RespuestaEntidad = resumen
	}

	return nil
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
		Solicitud:      input.Solicitud,
		Radicado:       input.Radicado,
		CreatedByID:    input.CreatedByID,
	}
	if err := s.repo.Create(ctx, ec); err != nil {
		if commondb.IsDuplicateKeyError(err) {
			return nil, ErrEntityCaseDuplicate
		}
		return nil, err
	}

	// Registrar las tareas seleccionadas del checklist como obligaciones
	// pendientes de este entity_case. Si alguna falla, se registra en log pero
	// no se aborta (la entidad ya quedó creada).
	if s.obligationRepo != nil {
		for _, obligID := range input.ObligationIds {
			obligID = strings.TrimSpace(obligID)
			if obligID == "" || obligID == "otra" {
				continue
			}
			oid := obligID
			eco := &models.EntityCaseObligation{
				EntityCaseID:       ec.ID,
				EntityObligationID: &oid,
				Status:             models.EntityCaseObligationStatusPendiente,
				CreatedByID:        input.CreatedByID,
			}
			if err := s.obligationRepo.Create(ctx, eco); err != nil {
				log.Printf("[EntityCaseService.Create] WARN: no se pudo crear entity_case_obligation %s para entity_case %s: %v", obligID, ec.ID, err)
			}
		}
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

	meta, err := s.repo.GetEntityMetaByBranchID(ctx, filter.EntityBranchID)
	if err != nil {
		return nil, err
	}

	items, total, err := s.repo.FindByEntityBranchIDPaginated(ctx, filter.EntityBranchID, filter.Document, page, pageSize)
	if err != nil {
		return nil, err
	}

	for i := range items {
		items[i].CaseStatusLabel = victimCaseStatusLabel(items[i].CaseStatus)
	}

	return &EntityCaseListResult{
		Items:          items,
		Total:          total,
		Page:           page,
		PageSize:       pageSize,
		EntityBranchID: filter.EntityBranchID,
		EntityName:     meta.EntityName,
		Sector:         meta.Sector,
		SectorName:     entitySectorLabel(meta.Sector),
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

func (s *entityCaseService) ListObligations(ctx context.Context, entityCaseID, status string) (*models.EntityCaseObligationsResponse, error) {
	if _, err := s.repo.FindByID(ctx, entityCaseID); err != nil {
		return nil, err
	}
	entityName, err := s.repo.GetEntityNameByEntityCaseID(ctx, entityCaseID)
	if err != nil {
		return nil, err
	}
	items, err := s.obligationRepo.FindByEntityCaseIDWithLabels(ctx, entityCaseID, status)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []models.EntityCaseObligationItem{}
	}
	return &models.EntityCaseObligationsResponse{
		EntityName: entityName,
		Items:      items,
	}, nil
}

func (s *entityCaseService) UpdateLastAction(ctx context.Context, entityCaseID, lastAction string) error {
	return s.repo.UpdateFields(ctx, entityCaseID, map[string]interface{}{"last_action": lastAction})
}

func (s *entityCaseService) GetObligacionesResumen(ctx context.Context, caseID string) (*models.ObligacionesResumenResponse, error) {
	return s.repo.GetObligacionesResumen(ctx, caseID)
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

func (s *entityCaseService) GetEntityCaseDetail(ctx context.Context, entityCaseID string) (*models.EntityCaseDetailRow, error) {
	row, err := s.repo.GetDetailByID(ctx, entityCaseID)
	if err != nil {
		return nil, err
	}
	row.CaseStatusLabel = victimCaseStatusLabel(row.CaseStatus)
	row.SectorLabel = entitySectorLabel(row.Sector)
	return row, nil
}
