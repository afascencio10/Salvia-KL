package service

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"errors"
	"strings"
	"time"
)

var (
	ErrReferralNotFound      = errors.New("remision no encontrada")
	ErrReferralConflict      = errors.New("la remisión no está pendiente de revisión")
	ErrReferralInvalidFilter = errors.New("filtro inválido")
	ErrReferralInvalidAgent  = errors.New("agente inválido")
	ErrReferralInvalidReview = errors.New("revisión inválida")
)

// EmergencyReferralListInput es lo que el controller ya resolvió desde la sesión.
type EmergencyReferralListInput struct {
	Vista          string
	SessionAgentID string
	Q              string
	Status         string
	AgentID        string
	Urgency        string
}

// EmergencyReferralCard es una métrica del recorte, antes de los filtros de UI.
type EmergencyReferralCard struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Value int    `json:"value"`
}

// EmergencyReferralBanner es el aviso de equipo técnico, convenio o Cruz Roja.
type EmergencyReferralBanner struct {
	Tono  string `json:"tono"`
	Texto string `json:"texto"`
}

// EmergencyReferralAgent es el agente asignado de una fila.
type EmergencyReferralAgent struct {
	ID     string `json:"id"`
	Nombre string `json:"nombre"`
}

// EmergencyReferralReview es la revisión técnica que pinta la fila.
type EmergencyReferralReview struct {
	ProfessionalName string `json:"professionalName"`
	Concept          string `json:"concept"`
	Result           string `json:"result"`
	ResultLabel      string `json:"resultLabel"`
	Date             string `json:"date"`
}

// EmergencyReferralDetail es la lectura de una remisión para la pantalla de detalle.
type EmergencyReferralDetail struct {
	EmergencyReferralItem
	Reason              string `json:"reason"`
	RejectionReason     string `json:"rejectionReason"`
	SubmittedByTeam     string `json:"submittedByTeam"`
	ChildrenNumber      *int   `json:"childrenNumber"`
	ChildrenAges        string `json:"childrenAges"`
	RequiresInterpreter bool   `json:"requiresInterpreter"`
}

// EmergencyReferralItem es una fila del listado.
type EmergencyReferralItem struct {
	ID              string                   `json:"id"`
	CaseICode       string                   `json:"caseICode"`
	CaseName        string                   `json:"caseName"`
	DocNumber       string                   `json:"docNumber"`
	RiskLabel       string                   `json:"riskLabel"`
	Municipality    string                   `json:"municipality"`
	Measures        []string                 `json:"measures"`
	Urgency         string                   `json:"urgency"`
	UrgencyLabel    string                   `json:"urgencyLabel"`
	SubmittedByName string                   `json:"submittedByName"`
	ReferredAt      string                   `json:"referredAt"`
	Status          string                   `json:"status"`
	StatusLabel     string                   `json:"statusLabel"`
	AssignedAgent   *EmergencyReferralAgent  `json:"assignedAgent"`
	TechnicalReview *EmergencyReferralReview `json:"technicalReview"`
}

// EmergencyReferralListResponse es la respuesta de E-01.
type EmergencyReferralListResponse struct {
	Alcance int                      `json:"alcance"`
	Cards   []EmergencyReferralCard  `json:"cards"`
	Banner  *EmergencyReferralBanner `json:"banner"`
	Items   []EmergencyReferralItem  `json:"items"`
}

// EmergencyReferralCatalogItem es una opción de agente o de profesional.
type EmergencyReferralCatalogItem struct {
	ID     string `json:"id"`
	Nombre string `json:"nombre"`
}

// EmergencyReferralCatalogResponse envuelve el catálogo.
type EmergencyReferralCatalogResponse struct {
	Items []EmergencyReferralCatalogItem `json:"items"`
}

// EmergencyReferralReassignResult es la respuesta de E-10.
type EmergencyReferralReassignResult struct {
	Updated []string `json:"updated"`
}

// EmergencyMeasureReferralService implementa la bandeja documentada.
type EmergencyMeasureReferralService interface {
	List(ctx context.Context, in EmergencyReferralListInput) (*EmergencyReferralListResponse, error)
	Detail(ctx context.Context, vista, sessionAgentID, id string) (*EmergencyReferralDetail, error)
	Agents(ctx context.Context) (*EmergencyReferralCatalogResponse, error)
	Reviewers(ctx context.Context) (*EmergencyReferralCatalogResponse, error)
	Reassign(ctx context.Context, ids []string, agentID, actorID string) (*EmergencyReferralReassignResult, error)
	Review(ctx context.Context, id, professionalID, concept, result string) (*EmergencyReferralItem, error)
}

type emergencyMeasureReferralService struct {
	repo   repository.EmergencyMeasureReferralRepository
	agents repository.AgentLightRepository
}

func NewEmergencyMeasureReferralService(repo repository.EmergencyMeasureReferralRepository, agents repository.AgentLightRepository) EmergencyMeasureReferralService {
	return &emergencyMeasureReferralService{repo: repo, agents: agents}
}

func (s *emergencyMeasureReferralService) List(ctx context.Context, in EmergencyReferralListInput) (*EmergencyReferralListResponse, error) {
	if err := validateListFilters(in); err != nil {
		return nil, err
	}
	counts, err := s.repo.CountByStatus(ctx, in.Vista, in.SessionAgentID)
	if err != nil {
		return nil, err
	}
	rows, err := s.repo.List(ctx, repository.EmergencyReferralListFilter{
		Vista:      in.Vista,
		AgentICode: in.SessionAgentID,
		Q:          strings.TrimSpace(in.Q),
		Status:     in.Status,
		AgentID:    in.AgentID,
		Urgency:    in.Urgency,
	})
	if err != nil {
		return nil, err
	}
	items := make([]EmergencyReferralItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, toItem(row))
	}
	resp := &EmergencyReferralListResponse{
		Alcance: sumCounts(counts),
		Items:   items,
	}
	if in.Vista == "supervisor" {
		resp.Cards = supervisorCards(counts, resp.Alcance)
	} else if in.Vista == "agente" {
		resp.Cards = agenteCards(counts, resp.Alcance)
	} else {
		resp.Banner = bannerDe(in.Vista)
	}
	return resp, nil
}

func (s *emergencyMeasureReferralService) Detail(ctx context.Context, vista, sessionAgentID, id string) (*EmergencyReferralDetail, error) {
	row, err := s.repo.FindDetail(ctx, vista, sessionAgentID, id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrReferralNotFound
	}
	detail := &EmergencyReferralDetail{
		EmergencyReferralItem: toItem(row.EmergencyReferralRow),
		Reason:                row.Reason,
		RejectionReason:       row.RejectionReason,
		SubmittedByTeam:       row.SubmittedByTeam,
		ChildrenNumber:        row.ChildrenNumber,
		ChildrenAges:          row.ChildrenAges,
		RequiresInterpreter:   row.RequiresInterpreter,
	}
	return detail, nil
}

func (s *emergencyMeasureReferralService) Agents(ctx context.Context) (*EmergencyReferralCatalogResponse, error) {
	return s.catalog(ctx, "am")
}

func (s *emergencyMeasureReferralService) Reviewers(ctx context.Context) (*EmergencyReferralCatalogResponse, error) {
	return s.catalog(ctx, "em")
}

func (s *emergencyMeasureReferralService) catalog(ctx context.Context, role string) (*EmergencyReferralCatalogResponse, error) {
	people, err := s.agents.FindAllByRole(ctx, role)
	if err != nil {
		return nil, err
	}
	items := make([]EmergencyReferralCatalogItem, 0, len(people))
	for _, person := range people {
		items = append(items, EmergencyReferralCatalogItem{
			ID:     person.ICode,
			Nombre: strings.TrimSpace(person.Names + " " + person.LastNames),
		})
	}
	return &EmergencyReferralCatalogResponse{Items: items}, nil
}

func (s *emergencyMeasureReferralService) Reassign(ctx context.Context, ids []string, agentID, actorID string) (*EmergencyReferralReassignResult, error) {
	clean := uniqueIDs(ids)
	if len(clean) == 0 || strings.TrimSpace(agentID) == "" {
		return nil, ErrReferralInvalidFilter
	}
	ok, err := s.repo.UserHasActiveRole(ctx, agentID, "am")
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrReferralInvalidAgent
	}
	if err := s.repo.Reassign(ctx, clean, agentID, actorID); err != nil {
		if errors.Is(err, repository.ErrReferralMissing) {
			return nil, ErrReferralNotFound
		}
		return nil, err
	}
	return &EmergencyReferralReassignResult{Updated: clean}, nil
}

func (s *emergencyMeasureReferralService) Review(ctx context.Context, id, professionalID, concept, result string) (*EmergencyReferralItem, error) {
	concept = strings.TrimSpace(concept)
	newStatus, ok := models.EmergencyReferralStatusPorRevision(result)
	if strings.TrimSpace(id) == "" || strings.TrimSpace(professionalID) == "" || concept == "" || !ok {
		return nil, ErrReferralInvalidReview
	}
	active, err := s.repo.UserHasActiveRole(ctx, professionalID, "em")
	if err != nil {
		return nil, err
	}
	if !active {
		return nil, ErrReferralInvalidReview
	}
	today := time.Now().Format("2006-01-02")
	found, pending, err := s.repo.Review(ctx, id, professionalID, concept, result, newStatus, today)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, ErrReferralNotFound
	}
	if !pending {
		return nil, ErrReferralConflict
	}
	rows, err := s.repo.List(ctx, repository.EmergencyReferralListFilter{Vista: "supervisor", ID: id})
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		if row.ID == id {
			item := toItem(row)
			return &item, nil
		}
	}
	return nil, ErrReferralNotFound
}

func validateListFilters(in EmergencyReferralListInput) error {
	if in.Urgency != "" && in.Urgency != models.EmergencyReferralUrgencyRiesgoInminente && in.Urgency != models.EmergencyReferralUrgencyNormal {
		return ErrReferralInvalidFilter
	}
	if in.Status == "" {
		return nil
	}
	if !statusKnown(in.Status) {
		return ErrReferralInvalidFilter
	}
	if in.Vista == "agente" && !contains(models.EmergencyReferralStatusesAgente(), in.Status) {
		return ErrReferralInvalidFilter
	}
	return nil
}

func statusKnown(status string) bool {
	return models.EmergencyReferralStatusLabel(status) != status
}

func contains(list []string, value string) bool {
	for _, item := range list {
		if item == value {
			return true
		}
	}
	return false
}

func uniqueIDs(ids []string) []string {
	seen := map[string]bool{}
	out := make([]string, 0, len(ids))
	for _, id := range ids {
		id = strings.TrimSpace(id)
		if id == "" || seen[id] {
			continue
		}
		seen[id] = true
		out = append(out, id)
	}
	return out
}

func sumCounts(counts map[string]int) int {
	n := 0
	for _, v := range counts {
		n += v
	}
	return n
}

func countOf(counts map[string]int, statuses ...string) int {
	n := 0
	for _, status := range statuses {
		n += counts[status]
	}
	return n
}

func supervisorCards(counts map[string]int, total int) []EmergencyReferralCard {
	return []EmergencyReferralCard{
		{Key: "total", Label: "Total", Value: total},
		{Key: "riesgo", Label: "🔥 Riesgo inminente", Value: countOf(counts, models.EmergencyReferralStatusesRiesgo()...)},
		{Key: "revision", Label: "Pendiente de Revisión", Value: counts[models.EmergencyReferralStatusPendienteDeRevision]},
		{Key: "porAsignar", Label: "Por asignar", Value: counts[models.EmergencyReferralStatusPorAsignar]},
		{Key: "asignadas", Label: "Asignadas", Value: counts[models.EmergencyReferralStatusAsignada]},
		{Key: "gestion", Label: "En gestión documental", Value: counts[models.EmergencyReferralStatusEnGestionDocumental]},
		{Key: "caracterizado", Label: "Caracterizado", Value: counts[models.EmergencyReferralStatusCaracterizado]},
		{Key: "aprobadas", Label: "Aprobadas", Value: counts[models.EmergencyReferralStatusAprobadas]},
	}
}

func agenteCards(counts map[string]int, total int) []EmergencyReferralCard {
	return []EmergencyReferralCard{
		{Key: "total", Label: "Mis remisiones", Value: total},
		{Key: "riesgo", Label: "🔥 Riesgo inminente", Value: counts[models.EmergencyReferralStatusRiesgoPorConfirmarSeguridad]},
		{Key: "asignadas", Label: "Asignadas", Value: counts[models.EmergencyReferralStatusAsignada]},
		{Key: "gestion", Label: "En gestión documental", Value: counts[models.EmergencyReferralStatusEnGestionDocumental]},
		{Key: "aprobadas", Label: "Aprobadas", Value: counts[models.EmergencyReferralStatusAprobadas]},
		{Key: "seguimiento", Label: "En Seguimiento", Value: counts[models.EmergencyReferralStatusEnSeguimiento]},
	}
}

func bannerDe(vista string) *EmergencyReferralBanner {
	switch vista {
	case "equipo-tecnico":
		return &EmergencyReferralBanner{
			Tono:  "equipo",
			Texto: "Bandeja \"Pendiente de Revisión\": casos remitidos a Medidas de Emergencia a la espera del concepto cualificado del Equipo Técnico.",
		}
	case "supervision-convenio":
		return &EmergencyReferralBanner{
			Tono:  "convenio",
			Texto: "Mostrando casos \"Por revisar trazabilidad\", \"Por aprobar cotizaciones\" y \"Riesgo inminente: por escalar a Cruz Roja\" (prioritario) a cargo de Supervisión del Convenio.",
		}
	case "cruz-roja":
		return &EmergencyReferralBanner{
			Tono:  "cruz-roja",
			Texto: "Mostrando casos \"Por cotizar\", \"Riesgo inminente: por entregar medida\" (prioritario) y \"Aprobadas\" pendientes de otorgamiento.",
		}
	default:
		return nil
	}
}

func toItem(row repository.EmergencyReferralRow) EmergencyReferralItem {
	item := EmergencyReferralItem{
		ID:              row.ID,
		CaseICode:       row.CaseICode,
		CaseName:        row.CaseName,
		DocNumber:       row.DocNumber,
		RiskLabel:       row.RiskLabel,
		Municipality:    row.Municipality,
		Measures:        row.Measures,
		Urgency:         row.Urgency,
		UrgencyLabel:    models.EmergencyReferralUrgencyLabel(row.Urgency),
		SubmittedByName: row.SubmittedByName,
		ReferredAt:      row.ReferredAt,
		Status:          row.Status,
		StatusLabel:     models.EmergencyReferralStatusLabel(row.Status),
	}
	if item.Measures == nil {
		item.Measures = []string{}
	}
	if row.AssignedAgentID != "" {
		item.AssignedAgent = &EmergencyReferralAgent{ID: row.AssignedAgentID, Nombre: row.AssignedAgentName}
	}
	if row.ReviewResult != "" || row.ReviewProfessionalName != "" {
		item.TechnicalReview = &EmergencyReferralReview{
			ProfessionalName: row.ReviewProfessionalName,
			Concept:          row.ReviewConcept,
			Result:           row.ReviewResult,
			ResultLabel:      models.EmergencyReferralReviewLabel(row.ReviewResult),
			Date:             row.ReviewDate,
		}
	}
	return item
}
