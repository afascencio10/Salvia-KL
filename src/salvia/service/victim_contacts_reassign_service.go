package service

import (
	"bitsflow/internal/repository"
	"context"
	"errors"
	"log"
	"strings"
)

// VictimContactsReassignInput parámetros de reasignación masiva de reportes.
type VictimContactsReassignInput struct {
	ReportICodes    []string
	AgentICode      string
	SupervisorICode string
	SupervisorName  string
}

// VictimContactsReassignResult resultado de POST /api/v1/victim-contacts/reasignar-bulk.
type VictimContactsReassignResult struct {
	OK      bool `json:"ok"`
	Updated int  `json:"updated"`
	Skipped int  `json:"skipped"`
}

// VictimContactsReassignService reasignación masiva de reportes desde
// reasignar-casos-modal (entityKind='reportes'). Mismo patrón que
// CasesReassignService, sin equipo previo que proteger (GAP-3) y sin
// follow_up_v2/case_timeline_event (no aplican a reportes).
type VictimContactsReassignService interface {
	ReassignBulk(ctx context.Context, input VictimContactsReassignInput) (VictimContactsReassignResult, error)
}

type victimContactsReassignService struct {
	repo    repository.VictimContactsReassignRepository
	logRepo repository.AppLogRepository
}

func NewVictimContactsReassignService(repo repository.VictimContactsReassignRepository, logRepo repository.AppLogRepository) VictimContactsReassignService {
	return &victimContactsReassignService{repo: repo, logRepo: logRepo}
}

func (s *victimContactsReassignService) ReassignBulk(ctx context.Context, input VictimContactsReassignInput) (VictimContactsReassignResult, error) {
	icodes := uniqueNonEmptyStrings(input.ReportICodes)
	if len(icodes) == 0 {
		return VictimContactsReassignResult{}, errors.New("debe indicar al menos un reporte")
	}
	agentICode := strings.TrimSpace(input.AgentICode)
	if agentICode == "" {
		return VictimContactsReassignResult{}, errors.New("debe indicar un agente")
	}

	agent, err := s.repo.FindAgentForReassign(ctx, agentICode)
	if err != nil {
		return VictimContactsReassignResult{}, errors.New("agente no encontrado")
	}
	if agent.Status != "e" {
		return VictimContactsReassignResult{}, errors.New("el agente seleccionado no está activo")
	}

	agentTeam := strings.TrimSpace(agent.Team)
	agentName := strings.TrimSpace(agent.FullName)
	if agentName == "" {
		agentName = agentICode
	}
	supervisorICode := strings.TrimSpace(input.SupervisorICode)
	supervisorName := strings.TrimSpace(input.SupervisorName)

	var result VictimContactsReassignResult
	for _, icode := range icodes {
		current, err := s.repo.FindVictimContactForReassign(ctx, icode)
		if err != nil {
			result.Skipped++
			continue
		}
		// Hito del timeline (WP-1 T6): «asignado» si el reporte no tenía agente, «reasignado» si ya tenía
		// (aunque sea el mismo: lo más simple). Va en la misma transacción que el cambio de agente.
		if err := s.repo.AssignWithTimeline(ctx, icode, agentICode, agentTeam,
			reportAssignmentDescription(agentName, agentTeam, supervisorName, supervisorICode),
			supervisorICode, strings.TrimSpace(current.AgentICode) != ""); err != nil {
			return VictimContactsReassignResult{}, err
		}
		if s.logRepo != nil {
			if logErr := s.logRepo.Log(ctx, icode, "Reasignación manual de reporte (bulk): agente actualizado", map[string]interface{}{
				"agentICode": agentICode,
				"agentName":  agentName,
				"team":       agentTeam,
				"actorICode": supervisorICode,
				"actorName":  supervisorName,
			}); logErr != nil {
				log.Printf("[WARN] VictimContactsReassignBulk: no se pudo guardar log icode=%s: %v", icode, logErr)
			}
		}
		result.Updated++
	}

	result.OK = true
	return result, nil
}

// reportAssignmentDescription arma el texto del hito de asignación:
// «Asignado a {agente} (equipo {equipo}) por {supervisor}.» Si el agente no tiene
// equipo se omite el paréntesis.
func reportAssignmentDescription(agentName, team, supervisorName, supervisorICode string) string {
	text := "Asignado a " + agentName
	if team != "" {
		text += " (equipo " + team + ")"
	}
	by := supervisorName
	if by == "" {
		by = supervisorICode
	}
	if by == "" {
		by = "el supervisor"
	}
	return text + " por " + by + "."
}
