package salvia_ctrl

import (
	"bitsflow/internal/repository"
	"context"
	"fmt"
	"log"
	"time"
)

// reportTypeV2Labels — etiqueta corta de cada tipo de reporte V2 para las
// descripciones del timeline.
var reportTypeV2Labels = map[string]string{
	reportTypeV2Victima:  "Víctima",
	reportTypeV2Tercero:  "Tercero",
	reportTypeV2Familiar: "Familiar",
	reportTypeV2Hombre:   "Hombre",
}

// recordReportReceived escribe el hito «Reporte recibido» (WP-1 T6) en el timeline
// del reporte. Se llama DESPUÉS del commit y corre en segundo plano, como el
// código legado: es un endpoint público, así que si el evento falla solo se deja
// un WARN y la respuesta a la persona no cambia. actorICode es quien reporta con
// sesión (vacío si es anónimo).
func recordReportReceived(reportICode, reportTypeCode, radicado, actorICode string) {
	repo := CaseTimelineRepo
	if repo == nil {
		return
	}
	label := reportTypeV2Labels[reportTypeCode]
	if label == "" {
		label = reportTypeCode
	}
	description := fmt.Sprintf("Reporte de %s recibido. Radicado %s.", label, radicado)
	go func() {
		event := repository.NewReportMilestoneEvent(reportICode, repository.ReportMilestoneReceived, description, actorICode, "", time.Now())
		if err := repo.Create(context.Background(), event); err != nil {
			log.Printf("[WARN] reporte recibido: no se pudo escribir el hito del timeline (reporte %s): %v", reportICode, err)
		}
	}()
}
