package service

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"fmt"
	"log"
	"time"
)

// AssignCaseInput contiene los datos recibidos por el endpoint.
type AssignCaseInput struct {
	DocNumber     string    `json:"docNumber"`
	Team          string    `json:"team"`
	AgentLogin    string    `json:"agentLogin"`
	FollowUpDates []string  `json:"followUpDates"`
}

// AssignCaseResult resume el resultado de la operación.
type AssignCaseResult struct {
	CaseICode       string   `json:"caseICode"`
	AgentICode      string   `json:"agentICode"`
	FollowUpsCreated []string `json:"followUpsCreated"`
	FollowUpsSkipped []string `json:"followUpsSkipped"`
}

type AssignCaseService interface {
	AssignCase(ctx context.Context, input AssignCaseInput) (*AssignCaseResult, error)
}

type assignCaseService struct {
	caseRepo   repository.VictimCaseLightRepository
	agentRepo  repository.AgentLightRepository
	followRepo repository.FollowUpV2Repository
}

func NewAssignCaseService(
	caseRepo repository.VictimCaseLightRepository,
	agentRepo repository.AgentLightRepository,
	followRepo repository.FollowUpV2Repository,
) AssignCaseService {
	return &assignCaseService{
		caseRepo:   caseRepo,
		agentRepo:  agentRepo,
		followRepo: followRepo,
	}
}

func (s *assignCaseService) AssignCase(ctx context.Context, input AssignCaseInput) (*AssignCaseResult, error) {
	log.Printf("[AssignCase] START docNumber=%s team=%s agentLogin=%s dates=%v", input.DocNumber, input.Team, input.AgentLogin, input.FollowUpDates)

	// 1. Buscar el caso abierto más reciente por número de documento
	vcase, err := s.caseRepo.FindOpenByDocNumber(ctx, input.DocNumber)
	if err != nil {
		log.Printf("[AssignCase] ERROR FindOpenByDocNumber docNumber=%s: %v", input.DocNumber, err)
		return nil, fmt.Errorf("caso no encontrado para el documento %q: %w", input.DocNumber, err)
	}
	log.Printf("[AssignCase] caso encontrado icode=%s status=%s", vcase.VictimICode, vcase.Status)

	// 2. Buscar el agente por login
	agent, err := s.agentRepo.FindByLogin(ctx, input.AgentLogin)
	if err != nil {
		log.Printf("[AssignCase] ERROR FindByLogin login=%s: %v", input.AgentLogin, err)
		return nil, fmt.Errorf("agente no encontrado con login %q: %w", input.AgentLogin, err)
	}
	log.Printf("[AssignCase] agente encontrado icode=%s nombre=%s %s team=%s", agent.ICode, agent.Names, agent.LastNames, agent.Team)

	// 3. Actualizar team y agent_id en victim_case
	if err := s.caseRepo.UpdateTeamAndAgent(ctx, vcase.VictimICode, input.Team, agent.ICode); err != nil {
		log.Printf("[AssignCase] ERROR UpdateTeamAndAgent: %v", err)
		return nil, fmt.Errorf("error actualizando team/agent en victim_case: %w", err)
	}
	log.Printf("[AssignCase] victim_case actualizado team=%s agentID=%s", input.Team, agent.ICode)

	result := &AssignCaseResult{
		CaseICode:  vcase.VictimICode,
		AgentICode: agent.ICode,
	}

	// 4. Para cada fecha de seguimiento: validar y crear si no existe
	minDate := time.Now().UTC().Truncate(24 * time.Hour).AddDate(0, 0, -10)
	log.Printf("[AssignCase] fecha mínima permitida: %s", minDate.Format("2006-01-02"))

	for _, rawDate := range input.FollowUpDates {
		date, parseErr := time.Parse(time.RFC3339, rawDate)
		if parseErr != nil {
			date, parseErr = time.Parse("2006-01-02", rawDate)
			if parseErr != nil {
				log.Printf("[AssignCase] ERROR parse fecha %q: %v", rawDate, parseErr)
				return nil, fmt.Errorf("fecha inválida %q (usar RFC3339 o YYYY-MM-DD): %w", rawDate, parseErr)
			}
		}

		// Descartar fechas anteriores a (hoy - 10 días)
		if date.UTC().Before(minDate) {
			log.Printf("[AssignCase] SKIP fecha=%s es anterior a la fecha mínima %s", rawDate, minDate.Format("2006-01-02"))
			result.FollowUpsSkipped = append(result.FollowUpsSkipped, rawDate)
			continue
		}

		// Verificar que no exista un seguimiento para ese caso en ese día calendario
		exists, checkErr := s.followRepo.ExistsForCaseOnDate(ctx, vcase.VictimICode, date)
		if checkErr != nil {
			log.Printf("[AssignCase] ERROR ExistsForCaseOnDate fecha=%s: %v", rawDate, checkErr)
			return nil, fmt.Errorf("error verificando seguimientos para fecha %s: %w", rawDate, checkErr)
		}
		if exists {
			log.Printf("[AssignCase] SKIP fecha=%s ya existe seguimiento", rawDate)
			result.FollowUpsSkipped = append(result.FollowUpsSkipped, rawDate)
			continue
		}

		// Crear el seguimiento
		agentICode := agent.ICode
		fu := &models.FollowUpV2{
			CaseID:        vcase.VictimICode,
			AgentID:       &agentICode,
			Team:          input.Team,
			Status:        models.FollowUpStatusPendiente,
			ScheduledDate: date,
			ScheduledTime: "08:00",
		}
		if createErr := s.followRepo.Create(ctx, fu); createErr != nil {
			log.Printf("[AssignCase] ERROR Create followUp fecha=%s: %v", rawDate, createErr)
			return nil, fmt.Errorf("error creando seguimiento para fecha %s: %w", rawDate, createErr)
		}
		log.Printf("[AssignCase] followUp creado fecha=%s id=%s", rawDate, fu.ID)
		result.FollowUpsCreated = append(result.FollowUpsCreated, rawDate)
	}

	if result.FollowUpsCreated == nil {
		result.FollowUpsCreated = []string{}
	}
	if result.FollowUpsSkipped == nil {
		result.FollowUpsSkipped = []string{}
	}

	log.Printf("[AssignCase] DONE created=%d skipped=%d", len(result.FollowUpsCreated), len(result.FollowUpsSkipped))
	return result, nil
}
