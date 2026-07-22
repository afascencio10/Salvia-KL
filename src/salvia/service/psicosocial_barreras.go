package service

import (
	"bitsflow/internal/constants"
	"bitsflow/internal/models"
	"context"
	"fmt"
	"log"
	"strings"
)

// ─── Formularios Psicosociales — Barreras (sección "Identificación de Barreras") ─────────────
//
// Ver DocsMD/Screens/psicosocial-sesion/Flujos/flow-E02-cuando-se-guarda-formulario.md §9.
//
// Cada uno de los 4 formularios psicosociales (Primer Contacto, Primera Atención, Atención
// Psicosocial, Cierre) tiene su propia copia de la sección "Identificación de Barreras",
// clonada 1:1 desde el formulario de Seguimiento (hacer_seguimiento). Por eso el repeater
// group y las 22 preguntas tienen un ID distinto en cada formulario — este archivo mapea
// formID → esos IDs para poder reutilizar exactamente la misma lógica de creación de
// barrier_v2 / case_task / entity_letter que usa processFollowUpSubmission.

// psicosocialBarrierQuestions agrupa los IDs del repeater "Identificación de Barreras" (22
// preguntas + su repeater group) para un formulario psicosocial específico.
type psicosocialBarrierQuestions struct {
	RepeaterGroupID string

	QSector             string // Q1  dropdown
	QSalud              string // Q2  multiple
	QInstitucionSalud   string // Q3  multiple
	QOtraBarreraSalud   string // Q4  text
	QJusticia           string // Q5  multiple
	QInstitucionJusticia string // Q6  multiple
	QOtraBarreraJusticia string // Q7  text
	QProteccion         string // Q8  multiple
	QInstitucionProteccion string // Q9  multiple
	QOtraBarreraProteccion string // Q10 text
	QOtraInstitucion    string // Q11 text

	QDepartamento string // Q12 dropdown
	QCiudad       string // Q13 dropdown
	QMunicipio    string // Q14 dropdown

	QEstructuralInstitucional string // Q15 multiple
	QEstructuralEconomico     string // Q16 multiple
	QEstructuralTerritorial   string // Q17 multiple
	QEstructuralDiferencial   string // Q18 multiple

	QFecha       string // Q19 date
	QFuncionario string // Q20 text
	QDescripcion string // Q21 text
	QGestion     string // Q22 multiple
}

// psicosocialBarrierQuestionsByForm — IDs reales capturados de Supabase (Jul 2026), ver flow-E02 §9.
var psicosocialBarrierQuestionsByForm = map[string]psicosocialBarrierQuestions{
	constants.FormIDPrimerContacto: {
		RepeaterGroupID:          "4e9a169b-76d7-4f4d-ab9d-c807cff457d5",
		QSector:                  "8d2977dd-1b37-4848-b6b3-03fbc4074d59",
		QSalud:                   "f92f8730-f5f5-4f2b-8612-b752be222d61",
		QInstitucionSalud:        "6e08bcc5-5ff5-4091-ad45-7020fff9d59d",
		QOtraBarreraSalud:        "9b13fa75-5a0c-4a41-ae3c-1090b3fda297",
		QJusticia:                "a5640619-c3c0-48f6-8988-3682d026039e",
		QInstitucionJusticia:     "0c7a385b-97c3-432f-85cd-1420ab8cce78",
		QOtraBarreraJusticia:     "7ac4e83c-c9f2-4f61-9e13-551475493299",
		QProteccion:              "cf58e6c6-e597-4d15-bbd8-a736c55230dc",
		QInstitucionProteccion:   "660c8c0d-c724-4975-82dd-c958a741c636",
		QOtraBarreraProteccion:   "f5fd2e06-8f79-4071-b807-383c055cb7bf",
		QOtraInstitucion:         "addf75d0-9208-4cf0-bc1f-3d154adf7885",
		QDepartamento:            "e95358fa-f413-4521-9b20-2b1a8eeb9ea5",
		QCiudad:                  "3583a842-1c51-499c-a1ee-9b94d95aa7db",
		QMunicipio:               "6cc9fd56-93f3-4c25-bc8f-64ffb42f51af",
		QEstructuralInstitucional: "0806983e-aa77-4ced-86ef-ab28caccb72c",
		QEstructuralEconomico:    "8eddca56-3277-4832-9b7f-ea6afb00c7a0",
		QEstructuralTerritorial:  "a981a38a-adfe-4b3c-a6e3-81e2a5109dcc",
		QEstructuralDiferencial:  "c3ef5ce4-e83e-4a6c-b4ab-3282c9cc7f2e",
		QFecha:                   "9bb2b7b9-cfc8-4424-81af-28601d983653",
		QFuncionario:             "6514e23d-f356-4870-8f09-72aa7081dabb",
		QDescripcion:             "5fe94bb1-5e86-45d9-84d5-1dc5bff44a2c",
		QGestion:                 "6ffcf120-436d-42a7-a696-4eff179e439e",
	},
	constants.FormIDPrimeraAtencion: {
		RepeaterGroupID:          "b7fe655d-7e4e-4f67-b997-a1bf7512210e",
		QSector:                  "e3d1c715-5dcc-4337-8dac-891b9f7fd923",
		QSalud:                   "ab86ee00-29d5-4140-8bf1-a1abc497e8a5",
		QInstitucionSalud:        "db20d59c-3287-4a93-b024-665aa7208e73",
		QOtraBarreraSalud:        "ab0372d8-6e63-4f63-b4c9-87ffc2f3166d",
		QJusticia:                "86dbb377-6c08-4f22-914c-90c048e69c18",
		QInstitucionJusticia:     "902f2c6e-7565-4d88-9a2e-14571414dcb3",
		QOtraBarreraJusticia:     "64727dc3-bc1a-4b4f-b990-0dd747d7eeb7",
		QProteccion:              "1caab793-20b9-4e6f-8310-7e87981ef8e5",
		QInstitucionProteccion:   "8601cc22-94c4-4ff1-9662-12481b5a497f",
		QOtraBarreraProteccion:   "5a70ddab-bb5c-4ee4-b554-22bfe84836f4",
		QOtraInstitucion:         "7d8e9fc4-98d9-433b-960b-e286a8861d76",
		QDepartamento:            "87ed2814-8590-419b-a463-0109b1f4cfe2",
		QCiudad:                  "0bd19847-748c-4c7e-9b78-92de08f0a2c4",
		QMunicipio:               "1c3f68cc-ae2a-4130-9ce1-a11f8f51083f",
		QEstructuralInstitucional: "909f35db-05f5-4d80-85b4-5556c763ceac",
		QEstructuralEconomico:    "b08d9085-af25-4656-b2e7-0ed9757c8a17",
		QEstructuralTerritorial:  "d0f7d935-7fb8-4022-8f47-a2bbce812328",
		QEstructuralDiferencial:  "671718b4-b176-4017-a2e0-102c6c6edaac",
		QFecha:                   "17fc6777-5516-4d84-a68c-98cb74c5e506",
		QFuncionario:             "acc48434-4cdd-49dd-8bbc-5c90d3e183bc",
		QDescripcion:             "488560ad-83a6-4d2e-a476-1bb77cc1b833",
		QGestion:                 "7f03f169-c311-4a05-ac8f-9d9284e28d8d",
	},
	constants.FormIDAtencionPsicosocial: {
		RepeaterGroupID:          "0fbc7da0-8349-4bcd-acd2-da3c5240b4cd",
		QSector:                  "ca4d6f63-f317-4d83-8f0b-d088adcb64cb",
		QSalud:                   "28589ba4-e118-4e6f-97b4-5e50b5cb0d8b",
		QInstitucionSalud:        "b6dabd16-2cb4-4c3c-96d3-42fa8c37a793",
		QOtraBarreraSalud:        "17a49a9a-cdc2-4f9a-83c3-77be1c752d4c",
		QJusticia:                "d0540ae0-02f1-49d6-ab12-dccec6ee0f60",
		QInstitucionJusticia:     "6f7c4c77-07e3-4424-908d-53033188c266",
		QOtraBarreraJusticia:     "0cf0f198-d8f3-42d9-b15f-4444a625023f",
		QProteccion:              "2b4d5cc2-6e84-4118-ae5c-ed0062516a7e",
		QInstitucionProteccion:   "d087ac37-e13b-45a2-99b6-06dc7055093c",
		QOtraBarreraProteccion:   "263aae61-2ffe-4ada-be33-fcc1d6468932",
		QOtraInstitucion:         "27f9e894-e236-4bc1-b9eb-ebd6669c1cb8",
		QDepartamento:            "ebea5f98-11cc-484b-84c9-adde4d88e3a5",
		QCiudad:                  "672848b1-507b-47d5-9b09-d2dabf8d4617",
		QMunicipio:               "94fcd8b2-31a7-49ea-94e3-1aec6f0e799f",
		QEstructuralInstitucional: "8ed9aef5-992d-4ca1-99c9-116c6380cb22",
		QEstructuralEconomico:    "170d1e1a-01f1-47b6-85f2-9bc68dd66233",
		QEstructuralTerritorial:  "2d578a08-f775-4f25-846f-1c327e262764",
		QEstructuralDiferencial:  "833e2d18-7dca-4fa7-aa8c-7a159b5b59ef",
		QFecha:                   "02209c0a-5eba-4ac8-87fe-60d5b6318a03",
		QFuncionario:             "bb24fa2a-5b70-4387-990e-5065535f7c61",
		QDescripcion:             "91ee401c-0f24-4b14-bd9c-19a951d984e3",
		QGestion:                 "b649d994-bcdc-46d5-8171-974fd6c94d38",
	},
	constants.FormIDCierre: {
		RepeaterGroupID:          "cb9b49e5-91f0-430e-9c0f-dc0235ef396f",
		QSector:                  "c7895557-ca09-4eb1-9092-2396de7f04f7",
		QSalud:                   "e2ba58bb-c312-4d83-a9e0-b0dd24f47203",
		QInstitucionSalud:        "c930a521-5cc4-49fd-93d9-0650b1d854b7",
		QOtraBarreraSalud:        "145f6f13-2aea-4a5a-98a0-97cba20268bd",
		QJusticia:                "a6193bb3-42a8-480b-a7d2-5e80adf889a6",
		QInstitucionJusticia:     "361e40ea-2210-4073-8aca-44804c4ff8f1",
		QOtraBarreraJusticia:     "cb8de13d-7201-4bef-9092-2a553a0185c3",
		QProteccion:              "634418e7-a241-4932-a406-55065df5c3f7",
		QInstitucionProteccion:   "40ee6e8a-ea27-489b-a0fd-59a78f137e05",
		QOtraBarreraProteccion:   "772349a3-92e0-48d2-8c41-e44a3a7f56e6",
		QOtraInstitucion:         "94136345-8de1-4d5e-be79-df85c632358a",
		QDepartamento:            "99adb055-39dd-4ee2-ae7f-16633d2c123d",
		QCiudad:                  "2cf2d66b-d435-4b77-8e5b-7a09b5535a05",
		QMunicipio:               "dc56cb40-d432-4d5e-a913-3e0c0e1b4963",
		QEstructuralInstitucional: "0c1eb2a1-fefd-427f-9305-bc0d819190dc",
		QEstructuralEconomico:    "af8ed0ae-a576-4925-96a8-c928b768ceec",
		QEstructuralTerritorial:  "77d577d3-5cb0-46c1-aedc-14c61db03161",
		QEstructuralDiferencial:  "030dae03-0b0a-4b5a-bcb4-8b938f9722f8",
		QFecha:                   "2a4e9741-c112-42ad-a0ec-8b4a6902605c",
		QFuncionario:             "3430820d-f526-4517-8d89-5303a576a3c4",
		QDescripcion:             "a2e1d338-6234-43ed-9d41-f9f940f47a1a",
		QGestion:                 "ed67fc17-d6e6-4d51-aeba-6e0e438ed589",
	},
}

// processPsicosocialBarrierEntries crea un BarrierV2 (y sus tareas/oficios derivados) por cada
// entrada del repeater "Identificación de Barreras" del formulario recién completado. Replica
// exactamente la lógica de la sección "3. Crear un BarrierV2 por cada entrada" de
// processFollowUpSubmission, pero relacionando el registro con la remisión psicosocial (ps) en
// lugar del follow_up directamente — usa ps.FollowUpID (el follow_up de origen de la remisión)
// para mantener consistencia con el resto de barreras del caso, y además fija
// CaseTask.PsychosocialSupportID para poder filtrar tareas por remisión psicosocial.
//
// No es un error si el formulario no tiene barreras configuradas o si no se agregó ninguna
// entrada — en ambos casos retorna (0, nil).
//
// teamContactID es el team_contact que se está completando en este guardado (el mismo `tc` de
// processPsicosocialSessionSubmission) — se fija en cada BarrierV2 creado para poder resolverlas
// después como "barreras activas" de la remisión (ver LoadSession en
// psychosocial_detail_service.go y flow-E02 §9).
func (s *formService) processPsicosocialBarrierEntries(ctx context.Context, formID, submissionID, actorID, teamContactID string, ps *models.PsychosocialSupport) (int, error) {
	cfg, ok := psicosocialBarrierQuestionsByForm[formID]
	if !ok {
		return 0, nil
	}
	if s.repeaterEntryRepo == nil || s.answerRepo == nil || s.barrierV2Repo == nil {
		log.Printf("⚠️  [processPsicosocialBarrierEntries] repeaterEntryRepo/answerRepo/barrierV2Repo es nil — barreras NO procesadas")
		return 0, nil
	}

	sectorBarrierQ := map[string]string{
		"salud":      cfg.QSalud,
		"justicia":   cfg.QJusticia,
		"proteccion": cfg.QProteccion,
	}
	sectorInstitutionQ := map[string]string{
		"salud":      cfg.QInstitucionSalud,
		"justicia":   cfg.QInstitucionJusticia,
		"proteccion": cfg.QInstitucionProteccion,
	}
	sectorOtherBarrierQ := map[string]string{
		"salud":      cfg.QOtraBarreraSalud,
		"justicia":   cfg.QOtraBarreraJusticia,
		"proteccion": cfg.QOtraBarreraProteccion,
	}

	entries, err := s.repeaterEntryRepo.FindBySubmissionIDAndGroupIDs(ctx, submissionID, []string{cfg.RepeaterGroupID})
	if err != nil {
		return 0, fmt.Errorf("processPsicosocialBarrierEntries: leer entradas de barreras: %w", err)
	}
	log.Printf("[processPsicosocialBarrierEntries] formID=%s psicosocialId=%s barreras: %d entradas encontradas", formID, ps.ID, len(entries))

	created := 0
	for _, entry := range entries {
		entryAnswers, err := s.answerRepo.FindByRepeaterEntryID(ctx, entry.ID)
		if err != nil {
			return created, fmt.Errorf("processPsicosocialBarrierEntries: leer respuestas entrada [%s]: %w", entry.ID, err)
		}
		entryMap := make(map[string]string, len(entryAnswers))
		for _, a := range entryAnswers {
			entryMap[a.QuestionID] = a.Value
		}

		sector := strings.TrimSpace(entryMap[cfg.QSector])

		b := &models.BarrierV2{
			CaseID:        ps.CaseID,
			FollowUpID:    ps.FollowUpID,
			TeamContactID: &teamContactID,
			CreatedByID:   actorID,
			Status:        models.BarrierV2StatusOpen,

			Sector:               sector,
			SpecificBarriers:     entryMap[sectorBarrierQ[sector]],
			SpecificInstitutions: entryMap[sectorInstitutionQ[sector]],
			OtherBarrierDesc:     entryMap[sectorOtherBarrierQ[sector]],
			InstitutionName:      entryMap[cfg.QOtraInstitucion],

			DepartmentID: entryMap[cfg.QDepartamento],
			CityID:       entryMap[cfg.QCiudad],
			TownID:       entryMap[cfg.QMunicipio],

			StructuralInstitutional: entryMap[cfg.QEstructuralInstitucional],
			StructuralEconomic:      entryMap[cfg.QEstructuralEconomico],
			StructuralTerritorial:   entryMap[cfg.QEstructuralTerritorial],
			StructuralDifferential:  entryMap[cfg.QEstructuralDiferencial],

			BarrierDate:        entryMap[cfg.QFecha],
			OfficialDependency: entryMap[cfg.QFuncionario],
			Description:        entryMap[cfg.QDescripcion],
			ManagementActions:  entryMap[cfg.QGestion],
		}
		log.Printf("[processPsicosocialBarrierEntries] creando barrera sector=%s entry=%s", sector, entry.ID)
		if err := s.barrierV2Repo.Create(ctx, b); err != nil {
			return created, fmt.Errorf("processPsicosocialBarrierEntries: crear barrera entry [%s]: %w", entry.ID, err)
		}
		created++

		// Crear tareas (y oficios si aplica) por cada opción seleccionada en gestión — idéntico
		// a la lógica de processFollowUpSubmission, agregando además PsychosocialSupportID.
		gestionRaw := strings.TrimSpace(entryMap[cfg.QGestion])
		if gestionRaw == "" {
			continue
		}
		gestionSinTarea := map[string]bool{"orientacion_llamada": true}
		gestionConOficio := map[string]bool{
			"activacion_ruta_interinstitucional": true,
			"articulacion_institucional":         true,
			"escalamiento_organismo_control":     true,
		}
		gestionLabels := map[string]string{
			"gestion_llamada":                    "Gestión administrativa - Llamada",
			"activacion_ruta_interinstitucional": "Activación de ruta interinstitucional",
			"articulacion_institucional":         "Articulación institucional",
			"escalamiento_organismo_control":     "Escalamiento a organismo de control",
			"alerta_barreras":                    "Alerta por barreras",
		}

		for _, gVal := range strings.Split(gestionRaw, ",") {
			gVal = strings.TrimSpace(gVal)
			if gVal == "" || gestionSinTarea[gVal] {
				continue
			}

			label, ok := gestionLabels[gVal]
			if !ok {
				label = gVal
			}

			institution := b.InstitutionName
			if institution == "" {
				institution = b.SpecificInstitutions
			}
			label = label + " (Barreras Psicosocial)"
			if sector != "" || institution != "" {
				label = fmt.Sprintf("%s | Sector: %s | Institución: %s", label, sector, institution)
			}

			barrierID := b.ID
			followUpID := ps.FollowUpID
			psychosocialSupportID := ps.ID
			var entityLetterID *string

			if gestionConOficio[gVal] {
				letter := &models.EntityLetter{
					BarrierID: barrierID,
					CaseID:    ps.CaseID,
					State:     models.EntityLetterStatePorProyectar,
					AgentID:   &actorID,
				}
				if s.entityLetterRepo != nil {
					if err := s.entityLetterRepo.Create(ctx, letter); err != nil {
						log.Printf("[processPsicosocialBarrierEntries] WARN: no se pudo crear entity_letter para barrera %s gestion=%s: %v", barrierID, gVal, err)
					} else {
						entityLetterID = &letter.ID
					}
				}
			}

			taskType := gVal
			if gestionConOficio[gVal] {
				taskType = "proyectar_oficio"
			} else if gVal == "alerta_barreras" {
				taskType = "comite_caso"
			}

			task := &models.CaseTask{
				Category:              "Barreras",
				Type:                  taskType,
				Description:           label,
				AssignedUserID:        actorID,
				Status:                models.CaseTaskStatusToDo,
				CaseID:                ps.CaseID,
				FollowUpID:            &followUpID,
				BarrierID:             &barrierID,
				EntityLetterID:        entityLetterID,
				PsychosocialSupportID: &psychosocialSupportID,
			}
			if s.caseTaskRepo != nil {
				if err := s.caseTaskRepo.Create(ctx, task); err != nil {
					log.Printf("[processPsicosocialBarrierEntries] WARN: no se pudo crear case_task para barrera %s gestion=%s: %v", barrierID, gVal, err)
				}
			}
		}
	}

	return created, nil
}
