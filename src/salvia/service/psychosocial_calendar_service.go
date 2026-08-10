package service

import (
	"bitsflow/internal/repository"
	"context"
	"errors"
	"time"
)

// PsychosocialCalendarService reglas del calendario psicosocial.
type PsychosocialCalendarService interface {
	MineEvents(ctx context.Context, agentID string, from, to time.Time) ([]repository.CalendarEvent, error)
	TeamEvents(ctx context.Context, filterAgentID string, from, to time.Time) ([]repository.CalendarEvent, error)
	ListAgents(ctx context.Context) ([]repository.CalendarAgentOption, error)
	MyVictims(ctx context.Context, agentID string) ([]repository.CalendarVictimOption, error)
	CreateSession(ctx context.Context, in repository.CreateSessionInput) (string, error)
	UpdateSession(ctx context.Context, id string, in repository.CreateSessionInput) error
	DeleteSession(ctx context.Context, id string) error
	CreateMeeting(ctx context.Context, in repository.CreateMeetingInput) (string, error)
	GetMeeting(ctx context.Context, id string) (repository.MeetingDetail, error)
	UpdateMeeting(ctx context.Context, id string, in repository.CreateMeetingInput) error
	DeleteMeeting(ctx context.Context, id string) error
}

type psychosocialCalendarService struct {
	repo repository.PsychosocialCalendarRepository
}

func NewPsychosocialCalendarService(repo repository.PsychosocialCalendarRepository) PsychosocialCalendarService {
	return &psychosocialCalendarService{repo: repo}
}

func (s *psychosocialCalendarService) MineEvents(ctx context.Context, agentID string, from, to time.Time) ([]repository.CalendarEvent, error) {
	if agentID == "" {
		return nil, errors.New("agentID requerido")
	}
	return s.repo.MineEvents(ctx, agentID, from, to)
}

func (s *psychosocialCalendarService) TeamEvents(ctx context.Context, filterAgentID string, from, to time.Time) ([]repository.CalendarEvent, error) {
	return s.repo.TeamEvents(ctx, filterAgentID, from, to)
}

func (s *psychosocialCalendarService) ListAgents(ctx context.Context) ([]repository.CalendarAgentOption, error) {
	return s.repo.ListAgents(ctx)
}

func (s *psychosocialCalendarService) MyVictims(ctx context.Context, agentID string) ([]repository.CalendarVictimOption, error) {
	if agentID == "" {
		return nil, errors.New("agentID requerido")
	}
	return s.repo.MyVictims(ctx, agentID)
}

func (s *psychosocialCalendarService) CreateSession(ctx context.Context, in repository.CreateSessionInput) (string, error) {
	if in.AgentID == "" || in.CaseICode == "" || in.ScheduledDate.IsZero() || in.ScheduledTime == "" {
		return "", errors.New("agentId, caseICode, scheduledDate, scheduledTime son requeridos")
	}
	return s.repo.CreateSession(ctx, in)
}

func (s *psychosocialCalendarService) UpdateSession(ctx context.Context, id string, in repository.CreateSessionInput) error {
	if id == "" {
		return errors.New("id requerido")
	}
	if in.ScheduledDate.IsZero() || in.ScheduledTime == "" {
		return errors.New("scheduledDate y scheduledTime requeridos")
	}
	return s.repo.UpdateSession(ctx, id, in)
}

func (s *psychosocialCalendarService) DeleteSession(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id requerido")
	}
	return s.repo.DeleteSession(ctx, id)
}

func (s *psychosocialCalendarService) CreateMeeting(ctx context.Context, in repository.CreateMeetingInput) (string, error) {
	if in.Title == "" || in.ScheduledDate.IsZero() || in.ScheduledTime == "" || len(in.AgentIDs) == 0 {
		return "", errors.New("title, scheduledDate, scheduledTime y al menos un agente son requeridos")
	}
	if in.CreatedBy == "" {
		return "", errors.New("createdBy requerido")
	}
	return s.repo.CreateMeeting(ctx, in)
}

func (s *psychosocialCalendarService) GetMeeting(ctx context.Context, id string) (repository.MeetingDetail, error) {
	if id == "" {
		return repository.MeetingDetail{}, errors.New("id requerido")
	}
	return s.repo.GetMeeting(ctx, id)
}

func (s *psychosocialCalendarService) UpdateMeeting(ctx context.Context, id string, in repository.CreateMeetingInput) error {
	if id == "" {
		return errors.New("id requerido")
	}
	if in.Title == "" || in.ScheduledDate.IsZero() || in.ScheduledTime == "" || len(in.AgentIDs) == 0 {
		return errors.New("title, scheduledDate, scheduledTime y al menos un agente son requeridos")
	}
	return s.repo.UpdateMeeting(ctx, id, in)
}

func (s *psychosocialCalendarService) DeleteMeeting(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id requerido")
	}
	return s.repo.DeleteMeeting(ctx, id)
}
