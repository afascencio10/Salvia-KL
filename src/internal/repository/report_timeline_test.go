package repository

import (
	"bitsflow/internal/models"
	"testing"
	"time"
)

func TestNewReportMilestoneEvent(t *testing.T) {
	at := time.Date(2026, 11, 2, 15, 30, 0, 0, time.FixedZone("COT", -5*3600))
	ev := NewReportMilestoneEvent("report-1", ReportMilestoneContactMade, "Contacto logrado el 02/11/2026 15:30.", "agent-1", "", at)

	if ev.VictimContactID != "report-1" || ev.CaseID != "" {
		t.Fatalf("ids: reporte=%q caso=%q", ev.VictimContactID, ev.CaseID)
	}
	if ev.Category != models.TimelineCategoryReportes || ev.Type != "Contacto logrado" || ev.Icon != "handshake" || ev.EventType != "CONTACTO_LOGRADO" {
		t.Fatalf("hito mal armado: %+v", ev)
	}
	if ev.EventUserID != "agent-1" || ev.Color == "" || ev.Description == "" {
		t.Fatalf("actor, color o descripción faltan: %+v", ev)
	}
	if ev.Date.Location() != time.UTC || !ev.Date.Equal(at) {
		t.Fatalf("la fecha debe guardarse en UTC y ser el mismo instante: %v", ev.Date)
	}
}

func TestReportMilestonesAreDistinctAndFit(t *testing.T) {
	all := []ReportMilestone{
		ReportMilestoneReceived, ReportMilestoneAssigned, ReportMilestoneReassigned,
		ReportMilestoneAttemptFailed, ReportMilestoneContactMade, ReportMilestoneClosedInvalid, ReportMilestoneDraftStarted,
	}
	types := map[string]bool{}
	for _, m := range all {
		if m.EventType == "" || m.Type == "" || m.Icon == "" || m.Color == "" {
			t.Fatalf("hito incompleto: %+v", m)
		}
		if len(m.EventType) > 30 { // case_timeline_event.event_type es varchar(30)
			t.Fatalf("event_type excede varchar(30): %q", m.EventType)
		}
		if len(m.Type) > 50 { // type es varchar(50)
			t.Fatalf("type excede varchar(50): %q", m.Type)
		}
		if types[m.Type] {
			t.Fatalf("tipo repetido: %q", m.Type)
		}
		types[m.Type] = true
	}
}

func TestReportMilestoneEventForDraftCarriesBothIDs(t *testing.T) {
	ev := NewReportMilestoneEvent("report-1", ReportMilestoneDraftStarted, "x", "agent-1", "case-1", time.Now())
	if ev.VictimContactID != "report-1" || ev.CaseID != "case-1" {
		t.Fatalf("el borrador debe llevar reporte y caso: %+v", ev)
	}
}
