package service

import (
	"bitsflow/internal/repository"
	"context"
	"strings"
)

const agentsSearchMinQueryLen = 3

// AgentSearchItem agente en la respuesta de GET /api/v1/agents/search (E-07).
type AgentSearchItem struct {
	ICode     string `json:"icode"`
	Names     string `json:"names"`
	LastNames string `json:"lastNames"`
	Team      string `json:"team"`
}

// AgentsSearchResponse respuesta JSON del autocomplete de persona asignada.
type AgentsSearchResponse struct {
	Agents []AgentSearchItem `json:"agents"`
}

// AgentPsicosocialSearchItem agente en GET /api/v1/agents/search-psicosocial (E-07).
type AgentPsicosocialSearchItem struct {
	ICode     string `json:"icode"`
	Names     string `json:"names"`
	LastNames string `json:"lastNames"`
	Team      string `json:"team"`
	Specialty string `json:"specialty"`
	RoleLabel string `json:"roleLabel"`
}

// AgentsPsicosocialSearchResponse respuesta JSON del autocomplete de profesional asignada.
type AgentsPsicosocialSearchResponse struct {
	Agents []AgentPsicosocialSearchItem `json:"agents"`
}

// AgentsSearchService contrato de búsqueda de agentes para casos-component.
type AgentsSearchService interface {
	Search(ctx context.Context, query string, limit int) (AgentsSearchResponse, error)
	SearchPsicosocial(ctx context.Context, query string, limit int) (AgentsPsicosocialSearchResponse, error)
}

type agentsSearchService struct {
	repo repository.AgentLightRepository
}

func NewAgentsSearchService(repo repository.AgentLightRepository) AgentsSearchService {
	return &agentsSearchService{repo: repo}
}

func (s *agentsSearchService) Search(ctx context.Context, query string, limit int) (AgentsSearchResponse, error) {
	query = strings.TrimSpace(query)
	if len(query) < agentsSearchMinQueryLen {
		return AgentsSearchResponse{Agents: []AgentSearchItem{}}, nil
	}

	rows, err := s.repo.SearchByName(ctx, query, limit)
	if err != nil {
		return AgentsSearchResponse{}, err
	}

	items := make([]AgentSearchItem, len(rows))
	for i, row := range rows {
		items[i] = AgentSearchItem{
			ICode:     row.ICode,
			Names:     row.Names,
			LastNames: row.LastNames,
			Team:      row.Team,
		}
	}

	return AgentsSearchResponse{Agents: items}, nil
}

func (s *agentsSearchService) SearchPsicosocial(ctx context.Context, query string, limit int) (AgentsPsicosocialSearchResponse, error) {
	query = strings.TrimSpace(query)
	if len(query) < agentsSearchMinQueryLen {
		return AgentsPsicosocialSearchResponse{Agents: []AgentPsicosocialSearchItem{}}, nil
	}

	rows, err := s.repo.SearchPsicosocialByName(ctx, query, limit)
	if err != nil {
		return AgentsPsicosocialSearchResponse{}, err
	}

	items := make([]AgentPsicosocialSearchItem, len(rows))
	for i, row := range rows {
		items[i] = AgentPsicosocialSearchItem{
			ICode:     row.ICode,
			Names:     row.Names,
			LastNames: row.LastNames,
			Team:      row.Team,
			Specialty: row.Specialty,
			RoleLabel: psicosocialRoleLabel(row.Specialty),
		}
	}

	return AgentsPsicosocialSearchResponse{Agents: items}, nil
}

func psicosocialRoleLabel(specialty string) string {
	switch strings.ToLower(strings.TrimSpace(specialty)) {
	case "psicologia":
		return "Psicóloga"
	case "trab. social":
		return "Trab. Social"
	default:
		return "Profesional"
	}
}
