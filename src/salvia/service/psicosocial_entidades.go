package service

import (
	"bitsflow/internal/constants"
	"bitsflow/internal/models"
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

// ─── Formularios Psicosociales — Entidades ("Seguimiento a Entidades" /
// "Identificación de Entidades") ────────────────────────────────────────────
//
// Ver DocsMD/Otros/temp/agregar-secciones-entidades-psicosocial.md.
//
// Cada uno de los 3 formularios en alcance (Primera Atención, Atención
// Psicosocial, Cierre — Primer Contacto queda fuera) tiene su propia copia de
// las 2 secciones "Seguimiento a Entidades" e "Identificación de Entidades",
// clonadas 1:1 desde el formulario de Hacer Seguimiento. Por eso el repeater
// group y las preguntas tienen un ID distinto en cada formulario — este
// archivo mapea formID → esos IDs, replicando exactamente el patrón de
// psicosocial_barreras.go.

// psicosocialEntidadQuestions agrupa los IDs de las secciones "Seguimiento a
// Entidades" e "Identificación de Entidades" para un formulario psicosocial
// específico.
type psicosocialEntidadQuestions struct {
	// Sección "Seguimiento a Entidades" (repeater state_items="currentEntities")
	RgSeguimientoEntidades string
	QSegRutaActualizada    string
	QSegMotivo             string
	QSegRequiereActivacion string
	QSegCanalActivacion    string

	// Sección "Identificación de Entidades"
	QReportoEntidadPrevio string

	RgEntidadesAcudidas string
	QAcSector           string
	QAcDepartamento     string
	QAcCiudad           string
	QAcMunicipio        string
	QAcEntidad          string
	QAcCompletadas      string
	QAcCompletadasOtra  string
	QAcPendientes       string
	QAcPendientesOtra   string
	QAcInfoRuta         string

	QAcuerdoActivacion string

	RgEntidadesActivacion string
	QAvSector             string
	QAvDepartamento       string
	QAvCiudad             string
	QAvMunicipio          string
	QAvEntidad            string
	QAvPendientes         string
	QAvPendientesOtra     string
	QAvCanalActivacion    string
}

// psicosocialEntidadQuestionsByForm — IDs reales capturados de Supabase (jul
// 2026, seed_entidades_psicosocial_template.sql ejecutado una vez por form).
// FormIDPrimerContacto NO se incluye — fuera de alcance.
var psicosocialEntidadQuestionsByForm = map[string]psicosocialEntidadQuestions{
	constants.FormIDPrimeraAtencion: {
		RgSeguimientoEntidades: "85080252-40ea-4e2a-9336-406fd84bcd18",
		QSegRutaActualizada:    "e21cc283-2a1d-4be2-8376-6dcf8effb11a",
		QSegMotivo:             "27030daf-9806-476e-9cc1-5ce442666df7",
		QSegRequiereActivacion: "8aeb7084-279d-44f3-afa7-d3a43f04dd24",
		QSegCanalActivacion:    "7b7729e7-14b6-407d-ab5f-3a90921b7389",

		QReportoEntidadPrevio: "c7a54a4c-38a9-4ad6-a0e7-4a4e2f38efe1",

		RgEntidadesAcudidas: "508260c4-b7de-4d95-ad2d-3f94cb54353a",
		QAcSector:           "15629935-3222-40ef-be18-cab414dacdff",
		QAcDepartamento:     "b9289ec8-0c68-43ef-9ab3-d9d39d316fd9",
		QAcCiudad:           "6f6192e1-60dd-466f-a0c7-003644780903",
		QAcMunicipio:        "0ac495ff-be8f-4443-805e-527d9f459dbd",
		QAcEntidad:          "98483b85-3e42-4e2a-9765-4e5b31851cf2",
		QAcCompletadas:      "4deb7fa3-caa5-417a-8cf7-2c0b38b0525f",
		QAcCompletadasOtra:  "14c8d4d9-770e-499b-8bc4-1dd30e2859be",
		QAcPendientes:       "8bd12eb5-c1b6-4bf1-8360-1fb96e0aeec1",
		QAcPendientesOtra:   "74d18d17-3656-494c-af78-e8f75078e0f8",
		QAcInfoRuta:         "03afdbfd-a0ef-4a1d-bf53-1049480670f0",

		QAcuerdoActivacion: "d5b50c54-01a8-4533-a9c3-304ad1edecb6",

		RgEntidadesActivacion: "0a59e505-f0c4-4994-a2a5-3bcf9e7fa6da",
		QAvSector:             "e30ff04c-9c9d-4bb5-bc2a-2d5de3c6c9c2",
		QAvDepartamento:       "9ee70bf9-d518-437d-b3fb-ca5a1d5dd91d",
		QAvCiudad:             "403f5b34-59ca-4c9b-bbf2-f3e95a271b9f",
		QAvMunicipio:          "c190136d-4407-4f84-9e37-ab05d67eb93d",
		QAvEntidad:            "3941db10-f4f5-4259-aa6c-39bc49dd48ee",
		QAvPendientes:         "3d95cfb9-11ee-4035-8913-9e7c5e794943",
		QAvPendientesOtra:     "67cc5310-64fd-40d7-9ecd-a790ffb00f1e",
		QAvCanalActivacion:    "bd5fd732-7129-413f-bebb-9ce7a8f34f01",
	},
	constants.FormIDAtencionPsicosocial: {
		RgSeguimientoEntidades: "367cb17f-c512-4127-98d3-a43822e04d39",
		QSegRutaActualizada:    "a81f9b95-e4ba-4b18-b3eb-3459d704ecfb",
		QSegMotivo:             "9a6d7a52-3b97-4130-82f0-6c5f25d6618b",
		QSegRequiereActivacion: "2a442945-5f46-4988-a2d5-8486c94f2d0a",
		QSegCanalActivacion:    "b071373b-a019-4e89-9542-56ee560c081a",

		QReportoEntidadPrevio: "0b3bd57d-7b07-478e-a0ca-2afcd70d9973",

		RgEntidadesAcudidas: "d67e72d2-fd1b-4199-a261-398bddd06dbc",
		QAcSector:           "614cd46f-dc80-4166-bac8-771e1ec548e6",
		QAcDepartamento:     "a1be03a8-6c0c-440c-acaf-3e103a838837",
		QAcCiudad:           "98f5871b-0ba6-4095-9456-834e97acfbb9",
		QAcMunicipio:        "542eed7e-69dc-45f8-8b8c-161dbdf25371",
		QAcEntidad:          "9658da08-101f-4a45-9204-f8f154cf36f0",
		QAcCompletadas:      "82948304-99ff-4e5b-a01e-1f8ebc4b2d61",
		QAcCompletadasOtra:  "473110b2-f2c0-417a-b1a3-07feb0400e3f",
		QAcPendientes:       "d6607084-e8fa-4e13-899b-d38053736cfd",
		QAcPendientesOtra:   "d64051b1-3e7e-47c8-8413-23cefd4c98ea",
		QAcInfoRuta:         "d3525a8c-80d4-4879-94eb-07692dbd6c66",

		QAcuerdoActivacion: "934ab7a1-c6d5-4e6a-b5ae-23ef79e62f2a",

		RgEntidadesActivacion: "10461694-7839-442e-91ea-139af9ad7f16",
		QAvSector:             "8175e98b-a759-4403-b426-c27e790c3a8b",
		QAvDepartamento:       "7d5898aa-447d-47d2-ac6e-fdcfe061397d",
		QAvCiudad:             "c2f51cef-ea36-42e9-8505-a0be8d2b401d",
		QAvMunicipio:          "784da3e0-83f2-472e-99e7-61ff083df353",
		QAvEntidad:            "45f1a39a-d212-4c2e-82a2-dbb66f6b471e",
		QAvPendientes:         "38eae94e-e7d5-45d7-9a67-2f95fd64a122",
		QAvPendientesOtra:     "8ef6a55d-96fc-450d-a509-f7a207737b09",
		QAvCanalActivacion:    "efb2bb18-df3c-4beb-a8c4-1ebc63fcf7e4",
	},
	constants.FormIDCierre: {
		RgSeguimientoEntidades: "c65f84c7-9751-431c-9f70-7fce9fe828a6",
		QSegRutaActualizada:    "a8568314-65f7-4f95-ab47-afcd6e0364bb",
		QSegMotivo:             "f8ed0ddf-8a86-490f-8a70-c83f4b9461da",
		QSegRequiereActivacion: "df718225-6e45-42b2-b7c5-8830ead5a6a5",
		QSegCanalActivacion:    "67ad3a7d-d5dc-46b0-9e6b-b9a1b87fcc9b",

		QReportoEntidadPrevio: "85a5373d-fd7b-4abd-a621-06539233cdb6",

		RgEntidadesAcudidas: "b4d5ed2c-19a9-43c7-b36e-93a44d9c6c71",
		QAcSector:           "74ed8a8a-3330-4e62-8788-9b84d2706d27",
		QAcDepartamento:     "6bc68182-f2d0-42f3-bb3d-0236c48cf535",
		QAcCiudad:           "4c08909d-9b5c-4ca0-bc8e-8f572cb24709",
		QAcMunicipio:        "5209d922-157a-411c-a47e-4f8d6063688b",
		QAcEntidad:          "db27821c-6276-4b45-b4cd-a8ebc24896cb",
		QAcCompletadas:      "1d3bd52d-a388-4720-9730-38a125431f0b",
		QAcCompletadasOtra:  "c6f91bd3-3e38-43ee-9840-4b9788f94f7b",
		QAcPendientes:       "e506c1cd-5b76-40b4-b7f2-7f09aeb4190c",
		QAcPendientesOtra:   "701acfb2-ee1f-4f4d-b354-066f24af4954",
		QAcInfoRuta:         "dc422813-3501-40af-9749-c1a985ae4d82",

		QAcuerdoActivacion: "be02e7f8-cf6c-4cb7-a198-804856a37f6d",

		RgEntidadesActivacion: "ffa24f17-9f27-4ebd-ba8b-bf4a5345ee46",
		QAvSector:             "cad27b82-1408-432e-9c2e-dde434f9d47a",
		QAvDepartamento:       "c4806bbf-fba7-4884-b0f8-17b71f28cf7a",
		QAvCiudad:             "2696dbe5-f50b-4d41-9364-63669ffa73ce",
		QAvMunicipio:          "3accf02b-50de-4a87-9c07-a6707559e61a",
		QAvEntidad:            "010cfd77-ea95-44cc-bb3d-8b0e3b484936",
		QAvPendientes:         "39682f44-1e74-43fc-8003-4a93e04239a1",
		QAvPendientesOtra:     "45d4e073-e90f-4943-b7dc-d1d25044ab78",
		QAvCanalActivacion:    "d31e400c-a628-4966-944e-65037970bb8e",
	},
}

// processPsicosocialEntidadEntries replica exactamente PASO 3c ("Seguimiento a
// Entidades") y PASO 3d ("Identificación de Entidades") de
// processFollowUpSubmission, relacionando los registros creados con la
// remisión psicosocial (ps) y el team_contact que se está completando, en vez
// del follow_up directamente. No es un error si el formulario no tiene estas
// secciones configuradas (Primer Contacto) — retorna nil sin hacer nada.
func (s *formService) processPsicosocialEntidadEntries(ctx context.Context, formID, submissionID, actorID, teamContactID string, ps *models.PsychosocialSupport) error {
	cfg, ok := psicosocialEntidadQuestionsByForm[formID]
	if !ok {
		return nil
	}
	if s.entityCaseSvc == nil || s.repeaterEntryRepo == nil || s.answerRepo == nil {
		log.Printf("⚠️  [processPsicosocialEntidadEntries] entityCaseSvc/repeaterEntryRepo/answerRepo es nil — entidades NO procesadas")
		return nil
	}

	// ── Seguimiento a Entidades (idéntico a PASO 3c de processFollowUpSubmission) ──
	entidadesDelCaso, err := s.entityCaseSvc.ListByCase(ctx, ps.CaseID)
	if err != nil {
		log.Printf("[processPsicosocialEntidadEntries] WARN: no se pudo listar entity_case del caso %s: %v", ps.CaseID, err)
	} else if len(entidadesDelCaso) > 0 {
		segEntries, err := s.repeaterEntryRepo.FindBySubmissionIDAndGroupIDs(ctx, submissionID, []string{cfg.RgSeguimientoEntidades})
		if err != nil {
			return fmt.Errorf("processPsicosocialEntidadEntries: leer entradas seguimiento a entidades: %w", err)
		}
		log.Printf("[processPsicosocialEntidadEntries] formID=%s seguimiento a entidades: %d entries (entidades del caso: %d)", formID, len(segEntries), len(entidadesDelCaso))

		for idx, entry := range segEntries {
			if idx >= len(entidadesDelCaso) {
				log.Printf("[processPsicosocialEntidadEntries] WARN: entry de seguimiento a entidades [%d] sin entity_case correspondiente — se omite", idx)
				continue
			}
			entityCaseID := entidadesDelCaso[idx].RelID
			entityBranchID := entidadesDelCaso[idx].EntityBranchID

			segAnswers, err := s.answerRepo.FindByRepeaterEntryID(ctx, entry.ID)
			if err != nil {
				return fmt.Errorf("processPsicosocialEntidadEntries: leer respuestas seguimiento a entidad [%s]: %w", entry.ID, err)
			}
			segMap := make(map[string]string, len(segAnswers))
			for _, a := range segAnswers {
				segMap[a.QuestionID] = a.Value
			}

			rutaActualizada := segMap[cfg.QSegRutaActualizada] == "true"
			requiereActivacion := segMap[cfg.QSegRequiereActivacion] == "true"
			motivo := strings.TrimSpace(segMap[cfg.QSegMotivo])
			canal := strings.TrimSpace(segMap[cfg.QSegCanalActivacion])

			if s.entityCaseFollowUpRepo != nil {
				tcID := teamContactID
				ecfu := &models.EntityCaseFollowUp{
					EntityCaseID:            entityCaseID,
					FollowUpID:              ps.FollowUpID,
					TeamContactID:           &tcID,
					RutaActualizada:         rutaActualizada,
					RequiereNuevaActivacion: requiereActivacion,
					CreatedByID:             actorID,
				}
				if motivo != "" {
					ecfu.MotivoActualizacion = &motivo
				}
				if canal != "" {
					ecfu.CanalActivacion = &canal
				}
				if err := s.entityCaseFollowUpRepo.Create(ctx, ecfu); err != nil {
					log.Printf("[processPsicosocialEntidadEntries] WARN: no se pudo crear entity_case_follow_up para entity_case %s: %v", entityCaseID, err)
				}
			}

			resumen := "Seguimiento registrado sin cambios en la ruta"
			if rutaActualizada && motivo != "" {
				resumen = fmt.Sprintf("Ruta actualizada: %s", motivo)
			}
			if err := s.entityCaseSvc.UpdateLastAction(ctx, entityCaseID, resumen); err != nil {
				log.Printf("[processPsicosocialEntidadEntries] WARN: no se pudo actualizar last_action de entity_case %s: %v", entityCaseID, err)
			}

			if requiereActivacion {
				psID := ps.ID
				s.procesarCanalActivacion(ctx, ps.CaseID, &ps.FollowUpID, actorID, entityCaseID, entityBranchID, canal, &psID)
			}
		}
	}

	// ── Identificación de Entidades (idéntico a PASO 3d) ──
	if s.entityCaseObligationRepo == nil {
		log.Printf("⚠️  [processPsicosocialEntidadEntries] entityCaseObligationRepo es nil — Identificación de Entidades NO procesada")
		return nil
	}

	type entidadRepeaterCfg struct {
		groupID          string
		qEntidad         string
		qCompletadas     string
		qCompletadasOtra string
		qPendientes      string
		qPendientesOtra  string
		qObjetivo        string // "" si no aplica
		qCanalActivacion string // "" si no aplica
	}
	repeaterConfigs := []entidadRepeaterCfg{
		{
			groupID: cfg.RgEntidadesAcudidas, qEntidad: cfg.QAcEntidad,
			qCompletadas: cfg.QAcCompletadas, qCompletadasOtra: cfg.QAcCompletadasOtra,
			qPendientes: cfg.QAcPendientes, qPendientesOtra: cfg.QAcPendientesOtra,
			qObjetivo: cfg.QAcInfoRuta,
		},
		{
			groupID: cfg.RgEntidadesActivacion, qEntidad: cfg.QAvEntidad,
			qPendientes: cfg.QAvPendientes, qPendientesOtra: cfg.QAvPendientesOtra,
			qCanalActivacion: cfg.QAvCanalActivacion,
		},
	}

	for _, rcfg := range repeaterConfigs {
		idEntries, err := s.repeaterEntryRepo.FindBySubmissionIDAndGroupIDs(ctx, submissionID, []string{rcfg.groupID})
		if err != nil {
			return fmt.Errorf("processPsicosocialEntidadEntries: leer entradas identificación de entidades [%s]: %w", rcfg.groupID, err)
		}
		log.Printf("[processPsicosocialEntidadEntries] identificación de entidades [%s]: %d entries", rcfg.groupID, len(idEntries))

		for _, entry := range idEntries {
			idAnswers, err := s.answerRepo.FindByRepeaterEntryID(ctx, entry.ID)
			if err != nil {
				return fmt.Errorf("processPsicosocialEntidadEntries: leer respuestas identificación de entidad [%s]: %w", entry.ID, err)
			}
			idMap := make(map[string]string, len(idAnswers))
			for _, a := range idAnswers {
				idMap[a.QuestionID] = a.Value
			}

			entityBranchIDStr := strings.TrimSpace(idMap[rcfg.qEntidad])
			if entityBranchIDStr == "" {
				continue
			}
			entityBranchID, err := strconv.ParseInt(entityBranchIDStr, 10, 64)
			if err != nil {
				log.Printf("[processPsicosocialEntidadEntries] WARN: entity_branch_id inválido %q en entry [%s]: %v", entityBranchIDStr, entry.ID, err)
				continue
			}

			var objetivo *string
			if rcfg.qObjetivo != "" {
				if v := strings.TrimSpace(idMap[rcfg.qObjetivo]); v != "" {
					objetivo = &v
				}
			}

			var entityCaseID string
			created, err := s.entityCaseSvc.Create(ctx, CreateEntityCaseInput{
				CaseID: ps.CaseID, EntityBranchID: entityBranchID,
				Objetivo: objetivo, CreatedByID: actorID,
			})
			if err != nil {
				if errors.Is(err, ErrEntityCaseDuplicate) {
					existentes, findErr := s.entityCaseSvc.ListByCase(ctx, ps.CaseID)
					if findErr != nil {
						log.Printf("[processPsicosocialEntidadEntries] WARN: no se pudo resolver entity_case existente tras duplicado (caso=%s, branch=%d): %v", ps.CaseID, entityBranchID, findErr)
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
						log.Printf("[processPsicosocialEntidadEntries] WARN: entity_case duplicado no encontrado en re-consulta (caso=%s, branch=%d)", ps.CaseID, entityBranchID)
						continue
					}
				} else {
					return fmt.Errorf("processPsicosocialEntidadEntries: crear entity_case (branch=%d): %w", entityBranchID, err)
				}
			} else {
				entityCaseID = created.ID
			}

			// Registrar obligaciones seleccionadas (completadas + pendientes)
			type obligGroup struct {
				raw     string
				otraTxt string
				status  string
			}
			var groups []obligGroup
			if rcfg.qCompletadas != "" {
				groups = append(groups, obligGroup{
					raw: idMap[rcfg.qCompletadas], otraTxt: idMap[rcfg.qCompletadasOtra],
					status: models.EntityCaseObligationStatusCompletada,
				})
			}
			groups = append(groups, obligGroup{
				raw: idMap[rcfg.qPendientes], otraTxt: idMap[rcfg.qPendientesOtra],
				status: models.EntityCaseObligationStatusPendiente,
			})

			for _, g := range groups {
				for _, val := range splitCSVValues(g.raw) {
					tcID := teamContactID
					eco := &models.EntityCaseObligation{
						EntityCaseID:  entityCaseID,
						Status:        g.status,
						FollowUpID:    &ps.FollowUpID,
						TeamContactID: &tcID,
						CreatedByID:   actorID,
					}
					if val == "otra" {
						label := strings.TrimSpace(g.otraTxt)
						eco.CustomLabel = &label
					} else {
						obligID := val
						eco.EntityObligationID = &obligID
					}
					if err := s.entityCaseObligationRepo.Create(ctx, eco); err != nil {
						log.Printf("[processPsicosocialEntidadEntries] WARN: no se pudo crear entity_case_obligation (%s) para entity_case %s: %v", g.status, entityCaseID, err)
					}
				}
			}

			// Canal de activación (solo aplica en el repeater de activación)
			if rcfg.qCanalActivacion != "" {
				psID := ps.ID
				s.procesarCanalActivacion(ctx, ps.CaseID, &ps.FollowUpID, actorID, entityCaseID, entityBranchID, idMap[rcfg.qCanalActivacion], &psID)
			}
		}
	}

	return nil
}
