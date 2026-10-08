package repository

import (
	"strings"
	"testing"
	"time"
)

// Cada '?' del SQL necesita su argumento. Es lo único de esta consulta que se puede
// comprobar sin base de datos, y es justo donde duele: un parámetro de más o de
// menos no lo ve el compilador y revienta al abrir el calendario.
func TestTeamEventsQueryParametrosCuadran(t *testing.T) {
	from := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	to := from.AddDate(0, 1, 0)

	casos := []struct {
		nombre    string
		agente    string
		dupla     string
		esperados int
	}{
		{"sin filtro", "", "", 4},
		{"por profesional", "agente-7", "", 10},
		{"por dupla", "", "dupla-3", 7},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			sql, args := buildTeamEventsQuery(caso.agente, caso.dupla, from, to)

			marcadores := strings.Count(sql, "?")
			if marcadores != len(args) {
				t.Fatalf("%d marcadores '?' contra %d argumentos", marcadores, len(args))
			}
			if marcadores != caso.esperados {
				t.Fatalf("se esperaban %d parámetros y hay %d: ¿cambió la consulta?", caso.esperados, marcadores)
			}
		})
	}
}

// El filtro por dupla no puede apoyarse en el JOIN a team_meeting_agent: si los dos
// integrantes están invitados a la misma reunión, el calendario la pintaría dos
// veces. Por eso usa EXISTS y el SELECT sin agente.
func TestTeamEventsPorDuplaNoDuplicaReuniones(t *testing.T) {
	sql, _ := buildTeamEventsQuery("", "dupla-3", time.Now(), time.Now())

	if strings.Contains(sql, calMeetingsFrom) {
		t.Fatal("el filtro por dupla no debe unir con team_meeting_agent: duplicaría las reuniones")
	}
	if !strings.Contains(sql, calMeetingsDuplaMatch) {
		t.Fatal("falta la condición EXISTS que limita las reuniones a las de la dupla")
	}
}

// Sin filtro la consulta no debe mencionar duplas ni agentes: es la agenda completa
// del equipo, y cualquier condición colada acá esconde eventos al supervisor.
func TestTeamEventsSinFiltroNoAplicaCondiciones(t *testing.T) {
	sql, _ := buildTeamEventsQuery("", "", time.Now(), time.Now())

	if strings.Contains(sql, calSessionsDuplaMatch) || strings.Contains(sql, calSessionsAgentMatch) {
		t.Fatal("sin filtro no debería aplicarse ninguna condición de dupla o agente")
	}
}
