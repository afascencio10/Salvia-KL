package repository

import (
	"bitsflow/internal/models"
	"time"
)

// ReportMilestone describe un hito del reporte en el timeline (categoría
// «Reportes»): código legacy, tipo visible, ícono y color. Ver la tabla de
// hitos de WP-1 T6.
type ReportMilestone struct {
	EventType string
	Type      string
	Icon      string
	Color     string
}

var (
	ReportMilestoneReceived = ReportMilestone{
		EventType: models.TimelineEventReporteRecibido, Type: models.TimelineTypeReporteRecibido,
		Icon: models.TimelineIconReporteRecibido, Color: models.TimelineColorBlue,
	}
	ReportMilestoneAssigned = ReportMilestone{
		EventType: models.TimelineEventReporteAsignado, Type: models.TimelineTypeReporteAsignado,
		Icon: models.TimelineIconReporteAsignado, Color: models.TimelineColorTeal,
	}
	ReportMilestoneReassigned = ReportMilestone{
		EventType: models.TimelineEventReporteReasignado, Type: models.TimelineTypeReporteReasignado,
		Icon: models.TimelineIconReasignacion, Color: models.TimelineColorPurple,
	}
	ReportMilestoneAttemptFailed = ReportMilestone{
		EventType: models.TimelineEventIntentoFallido, Type: models.TimelineTypeIntentoFallido,
		Icon: models.TimelineIconIntentoFallido, Color: models.TimelineColorOrange,
	}
	// Todo intento exitoso genera este mismo hito (no hay uno aparte para un segundo exitoso).
	ReportMilestoneContactMade = ReportMilestone{
		EventType: models.TimelineEventContactoLogrado, Type: models.TimelineTypeContactoLogrado,
		Icon: models.TimelineIconContactoLogrado, Color: models.TimelineColorGreen,
	}
	ReportMilestoneClosedInvalid = ReportMilestone{
		EventType: models.TimelineEventReporteCerradoInvalido, Type: models.TimelineTypeReporteCerradoInvalido,
		Icon: models.TimelineIconCierre, Color: models.TimelineColorRed,
	}
	ReportMilestoneDraftStarted = ReportMilestone{
		EventType: models.TimelineEventBorradorCasoIniciado, Type: models.TimelineTypeBorradorCasoIniciado,
		Icon: models.TimelineIconBorradorCaso, Color: models.TimelineColorPurple,
	}
)

// NewReportMilestoneEvent arma el evento de un hito del reporte (sub-flujo
// «Registrar un hito del reporte en el timeline»). El llamador lo guarda con
// repository.NewCaseTimelineEventRepository(tx).Create dentro de su transacción
// (eventos internos: asignación, intento, cierre, conversión) o con
// CaseTimelineRepo.Create después del commit (reporte recibido, best-effort).
//
// actorICode es quien origina el hito (el nombre se resuelve al leer, por JOIN).
// caseICode va vacío salvo en «Borrador de caso iniciado», que cuelga también del
// caso (vacío = NULL). at se guarda en UTC.
func NewReportMilestoneEvent(reportICode string, m ReportMilestone, description, actorICode, caseICode string, at time.Time) *models.CaseTimelineEvent {
	return &models.CaseTimelineEvent{
		VictimContactID: reportICode,
		CaseID:          caseICode,
		EventType:       m.EventType,
		Category:        models.TimelineCategoryReportes,
		Type:            m.Type,
		Icon:            m.Icon,
		Color:           m.Color,
		Description:     description,
		EventUserID:     actorICode,
		Date:            at.UTC(),
	}
}
