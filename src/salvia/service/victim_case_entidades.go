package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"bitsflow/internal/models"
)

// VictimCaseEntidadEntryInput es la entrada de una fila del repeater
// "Identificación de Entidades" en Registro de Caso — ya parseada por el
// controller (no hay repeaterEntryRepo/answerRepo, esto no vino de
// dinamic-form). Ver DocsMD/Otros/temp/agregar-identificacion-entidades-registro-caso-legacy.md.
type VictimCaseEntidadEntryInput struct {
	EntityBranchID  string
	Completadas     []string // vacío en el repeater de activación (no tiene esa pregunta)
	CompletadasOtra string
	Pendientes      []string
	PendientesOtra  string
	InfoRuta        string   // solo "acudidas" → entity_case.objetivo
	CanalActivacion []string // solo "activacion"
}

// ProcessVictimCaseEntidadEntries replica PASO 3d de processFollowUpSubmission
// (Identificación de Entidades) para el flujo de Registro de Caso — mismo
// código de negocio (entityCaseSvc, entityCaseObligationRepo,
// procesarCanalActivacion), pero:
//   - la entrada ya viene parseada (sin submission/answers en BD)
//   - FollowUpID queda nil en EntityCaseObligation y en procesarCanalActivacion
//     (el caso recién se crea, no existe ningún follow_up_v2 todavía)
//
// El frontend garantiza (validación bloqueante, ver F-5 del plan) que toda
// fila agregada por el agente esté completa — por eso acá no se rechaza nada,
// solo se ignora defensivamente cualquier fila sin EntityBranchID (mismo
// criterio que ya usa "Asignación de Ruta" con branchICode == "").
func (s *formService) ProcessVictimCaseEntidadEntries(ctx context.Context, caseID, actorID string, acudidas, activacion []VictimCaseEntidadEntryInput) error {
	if s.entityCaseSvc == nil {
		log.Printf("⚠️  [ProcessVictimCaseEntidadEntries] entityCaseSvc es nil — entidades NO procesadas para caso %s", caseID)
		return nil
	}

	type entidadGroup struct {
		entries      []VictimCaseEntidadEntryInput
		conCanal     bool // solo el repeater de activación tiene canal de activación
		conObjetivo  bool // solo el repeater de acudidas alimenta entity_case.objetivo
		conCompletad bool // solo el repeater de acudidas tiene checklist "completadas"
	}
	groups := []entidadGroup{
		{entries: acudidas, conObjetivo: true, conCompletad: true},
		{entries: activacion, conCanal: true},
	}

	for _, grp := range groups {
		for _, entry := range grp.entries {
			entityBranchIDStr := strings.TrimSpace(entry.EntityBranchID)
			if entityBranchIDStr == "" {
				continue
			}
			entityBranchID, err := strconv.ParseInt(entityBranchIDStr, 10, 64)
			if err != nil {
				log.Printf("[ProcessVictimCaseEntidadEntries] WARN: entity_branch_id inválido %q para caso %s: %v", entityBranchIDStr, caseID, err)
				continue
			}

			var objetivo *string
			if grp.conObjetivo {
				if v := strings.TrimSpace(entry.InfoRuta); v != "" {
					objetivo = &v
				}
			}

			var entityCaseID string
			created, err := s.entityCaseSvc.Create(ctx, CreateEntityCaseInput{
				CaseID: caseID, EntityBranchID: entityBranchID,
				Objetivo: objetivo, CreatedByID: actorID,
			})
			if err != nil {
				if errors.Is(err, ErrEntityCaseDuplicate) {
					existentes, findErr := s.entityCaseSvc.ListByCase(ctx, caseID)
					if findErr != nil {
						log.Printf("[ProcessVictimCaseEntidadEntries] WARN: no se pudo resolver entity_case existente tras duplicado (caso=%s, branch=%d): %v", caseID, entityBranchID, findErr)
						continue
					}
					found := false
					for _, e := range existentes {
						if e.EntityBranchID == entityBranchID {
							entityCaseID = e.RelID
							found = true
							break
						}
					}
					if !found {
						log.Printf("[ProcessVictimCaseEntidadEntries] WARN: entity_case duplicado no encontrado en re-consulta (caso=%s, branch=%d)", caseID, entityBranchID)
						continue
					}
				} else {
					return fmt.Errorf("ProcessVictimCaseEntidadEntries: crear entity_case (branch=%d): %w", entityBranchID, err)
				}
			} else {
				entityCaseID = created.ID
			}

			// Registrar obligaciones seleccionadas (completadas + pendientes)
			type obligGroup struct {
				values  []string
				otraTxt string
				status  string
			}
			var oGroups []obligGroup
			if grp.conCompletad {
				oGroups = append(oGroups, obligGroup{
					values: entry.Completadas, otraTxt: entry.CompletadasOtra,
					status: models.EntityCaseObligationStatusCompletada,
				})
			}
			oGroups = append(oGroups, obligGroup{
				values: entry.Pendientes, otraTxt: entry.PendientesOtra,
				status: models.EntityCaseObligationStatusPendiente,
			})

			if s.entityCaseObligationRepo != nil {
				for _, g := range oGroups {
					for _, val := range g.values {
						val = strings.TrimSpace(val)
						if val == "" {
							continue
						}
						eco := &models.EntityCaseObligation{
							EntityCaseID: entityCaseID,
							Status:       g.status,
							CreatedByID:  actorID,
						}
						if val == "otra" {
							label := strings.TrimSpace(g.otraTxt)
							eco.CustomLabel = &label
						} else {
							obligID := val
							eco.EntityObligationID = &obligID
						}
						if err := s.entityCaseObligationRepo.Create(ctx, eco); err != nil {
							log.Printf("[ProcessVictimCaseEntidadEntries] WARN: no se pudo crear entity_case_obligation (%s) para entity_case %s: %v", g.status, entityCaseID, err)
						}
					}
				}
			}

			// Canal de activación (solo aplica en el repeater de activación)
			if grp.conCanal {
				s.procesarCanalActivacion(ctx, caseID, nil, actorID, entityCaseID, entityBranchID, strings.Join(entry.CanalActivacion, ","), nil)
			}
		}
	}

	return nil
}
