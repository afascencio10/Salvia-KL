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

// AgentsSearchService contrato de búsqueda de agentes para casos-component.
type AgentsSearchService interface {
	Search(ctx context.Context, query string, limit int) (AgentsSearchResponse, error)
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
