//go:build availability

package service

import (
	"bitsflow/internal/models"
	"testing"
)

func strPtr(s string) *string { return &s }

func TestEvaluatePsicosocialAvailabilityWithScope_IndividualIncludesDuplaSession(t *testing.T) {
	prof := "PROF-1"
	duplaID := "DUPLA-1"
	timeStr := "10:00"

	// Remisión individual: solo personIDs=[PROF-1], pero duplaIDs expandido incluye DUPLA-1
	personIDs := []string{prof}
	duplaIDs := []string{duplaID}

	contacts := []models.TeamContact{
		{
			IsPsicoSession: true,
			DuplaID:        &duplaID,
			ScheduledTime:  strPtr("10:30"), // solapa con 10:00–12:00
		},
	}

	ok, msg := evaluatePsicosocialAvailabilityWithScope(contacts, timeStr, personIDs, duplaIDs)
	if ok {
		t.Fatalf("expected unavailable when overlapping dupla session exists; msg=%s", msg)
	}
}

func TestEvaluatePsicosocialAvailabilityWithScope_OwnProfessionalSession(t *testing.T) {
	prof := "PROF-1"
	contacts := []models.TeamContact{
		{
			IsPsicoSession: true,
			ProfessionalID: &prof,
			ScheduledTime:  strPtr("09:30"),
		},
	}
	ok, _ := evaluatePsicosocialAvailabilityWithScope(contacts, "10:00", []string{prof}, nil)
	if ok {
		t.Fatal("expected unavailable for own professional session overlap")
	}
}

func TestEvaluatePsicosocialAvailabilityWithScope_UnrelatedDuplaFree(t *testing.T) {
	prof := "PROF-1"
	otherDupla := "DUPLA-OTHER"
	contacts := []models.TeamContact{
		{
			IsPsicoSession: true,
			DuplaID:        &otherDupla,
			ScheduledTime:  strPtr("10:00"),
		},
	}
	ok, _ := evaluatePsicosocialAvailabilityWithScope(contacts, "10:00", []string{prof}, []string{"DUPLA-1"})
	if !ok {
		t.Fatal("expected available when overlapping session belongs to unrelated dupla")
	}
}

func TestSeedAvailabilityActors_RemisionDupla(t *testing.T) {
	duplaID := "DUPLA-1"
	ps := &models.PsychosocialSupport{DuplaID: &duplaID}
	persons, duplas := seedAvailabilityActors(ps, "PS-1", "TS-1")
	if len(persons) != 2 || len(duplas) != 1 || duplas[0] != duplaID {
		t.Fatalf("persons=%v duplas=%v", persons, duplas)
	}
}

func TestMergeDuplaIDs(t *testing.T) {
	got := mergeDuplaIDs([]string{"A", "B"}, []string{"B", "C"})
	if len(got) != 3 {
		t.Fatalf("got %v", got)
	}
}
