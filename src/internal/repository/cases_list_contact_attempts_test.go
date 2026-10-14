package repository

import (
	"database/sql"
	"strings"
	"testing"
)

func TestContactAttemptsPtr(t *testing.T) {
	if got := contactAttemptsPtr(sql.NullInt64{}); got != nil {
		t.Fatalf("NULL debe quedar nil, obtuvo %v", *got)
	}

	got := contactAttemptsPtr(sql.NullInt64{Int64: 0, Valid: true})
	if got == nil || *got != 0 {
		t.Fatalf("0 intentos debe ser puntero a 0, obtuvo %v", got)
	}

	got = contactAttemptsPtr(sql.NullInt64{Int64: 9, Valid: true})
	if got == nil || *got != 9 {
		t.Fatalf("9 intentos debe ser puntero a 9, obtuvo %v", got)
	}
}

func TestCurrentContactAttemptsSubquery(t *testing.T) {
	sqlText := currentContactAttemptsSubquery
	for _, fragment := range []string{
		"fu.attempts",
		"fu.status IN ('PENDIENTE', 'REPROGRAMADO')",
		"fu.deleted_at IS NULL",
		"ORDER BY fu.scheduled_date ASC, fu.created_at ASC",
	} {
		if !strings.Contains(sqlText, fragment) {
			t.Fatalf("la subconsulta no contiene %q:\n%s", fragment, sqlText)
		}
	}
	if strings.Contains(sqlText, "REALIZADO") {
		t.Fatal("los intentos 3x3 no deben contar seguimientos REALIZADO")
	}
	if strings.Contains(sqlText, "scheduled_date AT TIME ZONE") {
		t.Fatal("el seguimiento en curso incluye fechas vencidas; no debe filtrar por fecha futura")
	}
	if !strings.Contains(completedFollowUpsCountSubquery, "fu.status = 'REALIZADO'") {
		t.Fatal("seguimientos ejecutados debe seguir contando solo REALIZADO")
	}
}
