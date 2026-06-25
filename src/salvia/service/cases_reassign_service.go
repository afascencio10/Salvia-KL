package service

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// CasesReassignInput parámetros de reasignación masiva (M-05).
type CasesReassignInput struct {
	CaseICodes      []string
	AgentICode      string
	SupervisorICode string
	SupervisorName  string
}

// CasesReassignResult resultado de POST /api/v1/cases/reasignar-bulk.
type CasesReassignResult struct {
	OK               bool  `json:"ok"`
	Updated          int   `json:"updated"`
	Skipped          int   `json:"skipped"`
	FollowUpsUpdated int64 `json:"follow_ups_updated"`
}

// CasesReassignService reasignación masiva de casos desde el modal.
type CasesReassignService interface {
	ReassignBulk(ctx context.Context, input CasesReassignInput) (CasesReassignResult, error)
}

type casesReassignService struct {
	reassignRepo repository.CasesReassignRepository
	timelineRepo repository.CaseTimelineEventRepository
	db           *gorm.DB
}

func NewCasesReassignService(
	reassignRepo repository.CasesReassignRepository,
	timelineRepo repository.CaseTimelineEventRepository,
	db *gorm.DB,
) CasesReassignService {
	return &casesReassignService{
		reassignRepo: reassignRepo,
		timelineRepo: timelineRepo,
		db:           db,
	}
}

func (s *casesReassignService) ReassignBulk(ctx context.Context, input CasesReassignInput) (CasesReassignResult, error) {
	caseICodes := uniqueNonEmptyStrings(input.CaseICodes)
	if len(caseICodes) == 0 {
		return CasesReassignResult{}, errors.New("debe indicar al menos un caso")
	}
	agentICode := strings.TrimSpace(input.AgentICode)
	if agentICode == "" {
		return CasesReassignResult{}, errors.New("debe indicar un agente")
	}

	agent, err := s.reassignRepo.FindAgentForReassign(ctx, agentICode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return CasesReassignResult{}, errors.New("agente no encontrado")
		}
		return CasesReassignResult{}, err
	}
	if agent.Status != "e" {
		return CasesReassignResult{}, errors.New("el agente seleccionado no está activo")
	}

	agentTeam := strings.TrimSpace(agent.Team)
	agentName := strings.TrimSpace(agent.FullName)
	if agentName == "" {
		agentName = agentICode
	}

	supervisorName := strings.TrimSpace(input.SupervisorName)
	if supervisorName == "" {
		supervisorName = "Supervisor"
	}
	supervisorICode := strings.TrimSpace(input.SupervisorICode)

	var result CasesReassignResult

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := repository.NewCasesReassignRepository(tx)
		txTimeline := repository.NewCaseTimelineEventRepository(tx)

		for _, caseICode := range caseICodes {
			_, findErr := txRepo.FindVictimCaseForReassign(ctx, caseICode)
			if findErr != nil {
				if errors.Is(findErr, gorm.ErrRecordNotFound) {
					result.Skipped++
					continue
				}
				return findErr
			}

			if err := txRepo.UpdateVictimCaseAgent(ctx, caseICode, agentICode, agentTeam); err != nil {
				return err
			}

			fuCount, err := txRepo.UpdateFollowUpsForReassign(ctx, caseICode, agentICode, agentTeam)
			if err != nil {
				return err
			}
			result.FollowUpsUpdated += fuCount

			now := time.Now()
			description := fmt.Sprintf(
				"Caso reasignado a %s por %s",
				agentName,
				supervisorName,
			)
			if err := txTimeline.Create(ctx, &models.CaseTimelineEvent{
				CaseID:      caseICode,
				EventType:   models.TimelineEventReasignacion,
				Category:    models.TimelineCategorySeguimientos,
				Type:        models.TimelineTypeReasignacionCaso,
				Icon:        models.TimelineIconReasignacion,
				Color:       models.TimelineColorBlue,
				Date:        now,
				Description: description,
				ActorID:     supervisorICode,
				ActorName:   supervisorName,
				EventUserID: supervisorICode,
				CreatedAt:   now,
			}); err != nil {
				return err
			}

			result.Updated++
		}

		return nil
	})
	if err != nil {
		return CasesReassignResult{}, err
	}

	result.OK = true
	return result, nil
}

func uniqueNonEmptyStrings(items []string) []string {
	seen := make(map[string]struct{}, len(items))
	out := make([]string, 0, len(items))
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		out = append(out, item)
	}
	return out
}
