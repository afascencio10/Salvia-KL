package service

import "testing"

func TestReportAssignmentDescription(t *testing.T) {
	cases := []struct {
		agent, team, supName, supICode, want string
	}{
		{"Ana Pérez", "Riesgo alto", "Luis Gómez", "sup-1", "Asignado a Ana Pérez (equipo Riesgo alto) por Luis Gómez."},
		{"Ana Pérez", "", "Luis Gómez", "sup-1", "Asignado a Ana Pérez por Luis Gómez."},
		{"Ana Pérez", "Riesgo alto", "", "sup-1", "Asignado a Ana Pérez (equipo Riesgo alto) por sup-1."},
		{"Ana Pérez", "", "", "", "Asignado a Ana Pérez por el supervisor."},
	}
	for _, c := range cases {
		if got := reportAssignmentDescription(c.agent, c.team, c.supName, c.supICode); got != c.want {
			t.Errorf("got %q, want %q", got, c.want)
		}
	}
}
