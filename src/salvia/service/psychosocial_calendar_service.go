package service

import (
	"bitsflow/internal/repository"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

// PsychosocialCalendarService reglas del calendario psicosocial.
type PsychosocialCalendarService interface {
	MineEvents(ctx context.Context, agentID string, from, to time.Time) ([]repository.CalendarEvent, error)
	TeamEvents(ctx context.Context, filterAgentID, filterDuplaID string, from, to time.Time) ([]repository.CalendarEvent, error)
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

// TeamEvents — el filtro del supervisor es uno solo: o una dupla o un profesional.
// La pantalla manda uno u otro, pero el endpoint acepta los dos; si llegan juntos
// gana la dupla, en vez de cruzarlos y devolver un calendario vacío que parecería
// "esta dupla no tiene nada agendado".
func (s *psychosocialCalendarService) TeamEvents(ctx context.Context, filterAgentID, filterDuplaID string, from, to time.Time) ([]repository.CalendarEvent, error) {
	filterDuplaID = strings.TrimSpace(filterDuplaID)
	if filterDuplaID != "" {
		filterAgentID = ""
	}
	return s.repo.TeamEvents(ctx, strings.TrimSpace(filterAgentID), filterDuplaID, from, to)
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

// errPastDate — la sesión no puede agendarse en una fecha ya transcurrida.
// Compara por día calendario, no por instante: agendar "hoy" a cualquier hora
// del día sigue siendo válido aunque esa hora ya haya pasado.
var errPastDate = errors.New("la fecha de la sesión no puede ser anterior a hoy")

func isPastDate(d time.Time) bool {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return d.Before(today)
}

func (s *psychosocialCalendarService) CreateSession(ctx context.Context, in repository.CreateSessionInput) (string, error) {
	if in.AgentID == "" || in.CaseICode == "" || in.ScheduledDate.IsZero() || in.ScheduledTime == "" {
		return "", errors.New("agentId, caseICode, scheduledDate, scheduledTime son requeridos")
	}
	if isPastDate(in.ScheduledDate) {
		return "", errPastDate
	}
	if err := s.checkSessionConflict(ctx, in.AgentID, in.ScheduledDate, in.ScheduledTime, ""); err != nil {
		return "", err
	}
	return s.repo.CreateSession(ctx, in)
}

// checkSessionConflict aplica el mismo criterio de choque de horario que
// evaluatePsicosocialAvailability (psicosocial_availability.go): ventana de 2h,
// por agente. excludeID evita que una sesión choque consigo misma al editarla.
func (s *psychosocialCalendarService) checkSessionConflict(ctx context.Context, agentID string, date time.Time, timeStr, excludeID string) error {
	reqStart, ok := parseTimeToMinutes(timeStr)
	if !ok {
		return nil // formato de hora ya se validó antes de llegar aquí
	}
	times, err := s.repo.AgentSessionTimesOnDate(ctx, agentID, date, excludeID)
	if err != nil {
		return err
	}
	for _, t := range times {
		existStart, ok := parseTimeToMinutes(t)
		if !ok {
			continue
		}
		if sessionWindowsOverlap(reqStart, existStart, psicosocialSessionDurationMin) {
			existEnd := existStart + psicosocialSessionDurationMin
			return fmt.Errorf(
				"el agente ya tiene una sesión de %02d:%02d a %02d:%02d (ventana de 2 horas)",
				existStart/60, existStart%60, existEnd/60, existEnd%60,
			)
		}
	}
	return nil
}

func (s *psychosocialCalendarService) UpdateSession(ctx context.Context, id string, in repository.CreateSessionInput) error {
	if id == "" {
		return errors.New("id requerido")
	}
	if in.ScheduledDate.IsZero() || in.ScheduledTime == "" {
		return errors.New("scheduledDate y scheduledTime requeridos")
	}
	// No se valida fecha pasada aquí a propósito: una sesión PENDIENTE cuya fecha ya
	// venció sigue siendo editable (ver canEditSession en psicosocial_calendar.html),
	// y el agente puede necesitar corregir solo el resumen sin reprogramarla.
	// El choque de horario sí se valida: mover una sesión a un horario que ya
	// ocupa otra es el mismo bug que crearla así, sin el caso borde de la fecha.
	if err := s.checkSessionConflict(ctx, in.AgentID, in.ScheduledDate, in.ScheduledTime, id); err != nil {
		return err
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
