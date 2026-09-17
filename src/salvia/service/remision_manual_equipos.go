package service

import (
	"bitsflow/internal/models"
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

// ErrRemisionManualInvalid indica que faltan campos requeridos: caseId,
// actorId, o ningún equipo seleccionado.
var ErrRemisionManualInvalid = errors.New("form_service: remisión manual inválida")

// qEquiposRemisionManual es la misma pregunta "¿Cuáles equipos?" (Sección 5,
// Seguimiento de Caso) que usa processEquiposYMedidas — repetida acá solo para
// validar que el modal mandó al menos un equipo antes de procesar nada.
const qEquiposRemisionManual = "e0d38cf5-fe3f-45cb-9fd3-f5b8f7b2f7dc"

// ProcessManualRemision crea derivaciones a equipos SALVIA (y medidas de
// emergencia) directamente para un caso, sin pasar por el formulario de
// seguimiento — usado por el modal "Remisión a equipos" del Equipo Técnico en
// Detalle del Caso → Derivaciones (preguntas de orden 5 a 13 de la sección
// Seguimiento de Caso, replicadas fuera de dinamic-form). Reutiliza
// processEquiposYMedidas (misma lógica y reglas de PASO 4-5 de
// processFollowUpSubmission) con followUpID nil, porque no existe un
// seguimiento de origen.
func (s *formService) ProcessManualRemision(ctx context.Context, caseID, actorID string, answerMap map[string]string) (int, error) {
	caseID = strings.TrimSpace(caseID)
	actorID = strings.TrimSpace(actorID)
	if caseID == "" || actorID == "" || strings.TrimSpace(answerMap[qEquiposRemisionManual]) == "" {
		return 0, ErrRemisionManualInvalid
	}

	if s.caseRepo != nil {
		if _, err := s.caseRepo.FindByICode(ctx, caseID); err != nil {
			return 0, ErrCaseNotFound
		}
	}

	remisionCount, err := s.processEquiposYMedidas(ctx, caseID, nil, actorID, answerMap)
	if err != nil {
		return remisionCount, fmt.Errorf("ProcessManualRemision: %w", err)
	}

	if s.caseTimelineRepo != nil {
		actorName := actorID
		if s.agentLightRepo != nil {
			if agent, aErr := s.agentLightRepo.FindByICode(ctx, actorID); aErr == nil && agent != nil {
				actorName = strings.TrimSpace(agent.Names + " " + agent.LastNames)
			}
		}
		now := time.Now()
		event := &models.CaseTimelineEvent{
			CaseID:      caseID,
			Category:    models.TimelineCategoryEquipoTecnico,
			Type:        models.TimelineTypeRemisionEquipos,
			Icon:        models.TimelineIconRegistro,
			Date:        now,
			Description: fmt.Sprintf("Remisión a equipos registrada por el Equipo Técnico. Se realizaron %d %s a Equipos Salvia", remisionCount, pluralize(remisionCount, "remisión", "remisiones")),
			EventUserID: actorID,
			ActorName:   actorName,
			Color:       models.TimelineColorPurple,
			CreatedAt:   now,
		}
		if err := s.caseTimelineRepo.Create(ctx, event); err != nil {
			log.Printf("[ProcessManualRemision] WARN: no se pudo crear evento timeline caso %s: %v", caseID, err)
		}
	}

	return remisionCount, nil
}
