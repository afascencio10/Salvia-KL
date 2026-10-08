package service

import (
	"bitsflow/internal/repository"
	"context"
	"testing"
	"time"
)

type fakeCalendarRepo struct {
	repository.PsychosocialCalendarRepository
	agenteRecibido string
	duplaRecibida  string
}

func (f *fakeCalendarRepo) TeamEvents(_ context.Context, agentID, duplaID string, _, _ time.Time) ([]repository.CalendarEvent, error) {
	f.agenteRecibido = agentID
	f.duplaRecibida = duplaID
	return []repository.CalendarEvent{}, nil
}

func teamEvents(t *testing.T, agentID, duplaID string) *fakeCalendarRepo {
	t.Helper()
	repo := &fakeCalendarRepo{}
	if _, err := NewPsychosocialCalendarService(repo).TeamEvents(context.Background(), agentID, duplaID, time.Now(), time.Now()); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	return repo
}

// Los dos filtros son excluyentes. Si se cruzaran, la consulta pediría las sesiones
// de un profesional que además sean de otra dupla: casi siempre ninguna, y el
// supervisor leería el calendario vacío como "esta dupla no tiene nada agendado".
func TestTeamEventsLaDuplaGanaSobreElAgente(t *testing.T) {
	repo := teamEvents(t, "agente-7", "dupla-3")

	if repo.duplaRecibida != "dupla-3" {
		t.Fatalf("la dupla debería llegar al repositorio, llegó %q", repo.duplaRecibida)
	}
	if repo.agenteRecibido != "" {
		t.Fatalf("con dupla seleccionada el agente debe ir vacío, llegó %q", repo.agenteRecibido)
	}
}

func TestTeamEventsSinDuplaConservaElAgente(t *testing.T) {
	repo := teamEvents(t, "agente-7", "")

	if repo.agenteRecibido != "agente-7" {
		t.Fatalf("el agente debería llegar al repositorio, llegó %q", repo.agenteRecibido)
	}
	if repo.duplaRecibida != "" {
		t.Fatalf("sin dupla seleccionada la dupla debe ir vacía, llegó %q", repo.duplaRecibida)
	}
}

// Un valor en blanco no es un filtro: sin recortarlo, el repositorio lo trataría
// como un id y devolvería un calendario vacío en vez de la agenda completa.
func TestTeamEventsIgnoraFiltrosEnBlanco(t *testing.T) {
	repo := teamEvents(t, "   ", "  ")

	if repo.agenteRecibido != "" || repo.duplaRecibida != "" {
		t.Fatalf("los filtros en blanco deben quedar vacíos, llegaron %q y %q", repo.agenteRecibido, repo.duplaRecibida)
	}
}
