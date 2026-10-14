package controller

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"testing"
)

// TestApplyEntidadMecanismoRestriction reproduce la tarea "Filtrar vista de detalle del
// caso para entidades y mecanismos": et/mn no deben recibir del backend los datos de
// Derivaciones/Seguimientos, pero sí deben conservar Info General y Gestión Institucional.
func TestApplyEntidadMecanismoRestriction(t *testing.T) {
	build := func() *repository.CaseDetailData {
		return &repository.CaseDetailData{
			Case:                   models.VictimCase{VictimCaseICode: "case-1"},
			FollowUp:               &models.FollowUp{},
			Entries:                []models.FollowUpEntry{{}},
			FollowUpsV2:            []models.FollowUpV2{{}},
			EmergencyMeasures:      []models.EmergencyMeasure{{}},
			PsychosocialSupports:   []repository.PsychosocialSupportView{{}},
			EconomicStabilizations: []models.EconomicStabilization{{}},
			Barriers:               []models.BarrierV2{{}},
			TimelineEvents:         []models.CaseTimelineEvent{{}},
		}
	}

	for _, role := range []string{"et", "mn", models.ArticulatingMechanismRoleCode} {
		detail := build()
		applyEntidadMecanismoRestriction(detail, role)

		if detail.FollowUp != nil || detail.Entries != nil || detail.FollowUpsV2 != nil {
			t.Fatalf("role %q: Seguimientos debería quedar oculto, obtuve FollowUp=%v Entries=%v FollowUpsV2=%v", role, detail.FollowUp, detail.Entries, detail.FollowUpsV2)
		}
		if detail.EmergencyMeasures != nil || detail.PsychosocialSupports != nil || detail.EconomicStabilizations != nil {
			t.Fatalf("role %q: Derivaciones debería quedar oculto", role)
		}
		if len(detail.Barriers) == 0 {
			t.Fatalf("role %q: Gestión Institucional (Barriers) no debería tocarse", role)
		}
		if len(detail.TimelineEvents) == 0 {
			t.Fatalf("role %q: TimelineEvents no debería tocarse (lo usa Info General para 'Hechos del caso')", role)
		}
		if detail.Case.VictimCaseICode != "case-1" {
			t.Fatalf("role %q: Case (Info General) no debería tocarse", role)
		}
	}

	for _, role := range []string{"sv", "op", "ro", "ps"} {
		detail := build()
		applyEntidadMecanismoRestriction(detail, role)
		if detail.FollowUp == nil || detail.FollowUpsV2 == nil || detail.EmergencyMeasures == nil {
			t.Fatalf("role %q: no debería restringirse, los datos internos deberían seguir presentes", role)
		}
	}
}
