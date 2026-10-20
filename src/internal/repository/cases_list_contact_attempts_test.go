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

func TestEmptyCaseStatusBindsToActive(t *testing.T) {
	if !strings.Contains(caseStatusSelect, "THEN 'ra'") {
		t.Fatalf("el listado debe presentar el estado vacío como ra:\n%s", caseStatusSelect)
	}
	if !strings.Contains(caseStatusSelect, "btrim(vc.victim_case_status) = ''") {
		t.Fatal("el bind debe cubrir el estado en blanco")
	}

	where, args := buildCasesListWhere(CasesListFilters{FilterEstadoCaso: "ra"})
	if !strings.Contains(where, "btrim(vc.victim_case_status) = ''") {
		t.Fatalf("filtrar Activo debe incluir el estado vacío:\n%s", where)
	}
	if len(args) != 1 || args[0] != "ra" {
		t.Fatalf("el filtro Activo debe comparar contra ra, obtuvo %#v", args)
	}

	where, args = buildCasesListWhere(CasesListFilters{FilterEstadoCaso: "in"})
	if strings.Contains(where, "btrim(vc.victim_case_status)") {
		t.Fatalf("Iniciado no debe arrastrar los vacíos:\n%s", where)
	}
	if !strings.Contains(where, "vc.victim_case_status = ?") || len(args) != 1 || args[0] != "in" {
		t.Fatalf("Iniciado debe ser igualdad exacta, where=%s args=%#v", where, args)
	}

	where, args = buildCasesListWhere(CasesListFilters{FilterEstadoCaso: "cd"})
	if strings.Contains(where, "btrim(vc.victim_case_status)") {
		t.Fatalf("otro estado no debe arrastrar los vacíos:\n%s", where)
	}
	if !strings.Contains(where, "vc.victim_case_status = ?") || len(args) != 1 || args[0] != "cd" {
		t.Fatalf("Cerrado debe seguir siendo igualdad exacta, where=%s args=%#v", where, args)
	}
}
