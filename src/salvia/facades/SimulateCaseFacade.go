package salvia_facades

// SimulateCaseFacade.go
//
// Endpoint de simulación para la creación de casos.
// Permite probar el flujo completo de creación sin necesidad del front-end.
//
// POST /salvia/casos/simular
//
// Body (JSON):
//
//	{
//	  "nombres":           "Ana María",
//	  "apellidos":         "Rodríguez López",
//	  "tipoDoc":           "cc",          // código del tipo de documento
//	  "numDoc":            "1012345678",
//	  "esPareja":          true,          // true = agresor es pareja (pi/ex) → usa tamizaje de pareja
//	  "nivelRiesgo":       3,             // 1=Bajo 2=Moderado 3=Alto 4=Extremo
//	  "municipioAtencion": "11001000",    // opcional, default = Bogotá
//	  "municipioResidencia": "11001000",  // opcional
//	  "municipioHechos":   "11001000"     // opcional
//	}
//
// Responde igual que POST /salvia/casos (201 + credenciales del usuario creado).

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	salvia_ctrl "bitsflow/salvia/controllers"
	salvia_daos "bitsflow/salvia/dao"
	"bitsflow/common/utils"

	"github.com/gin-gonic/gin"
)

// operadorICodDefault es el icode de un operador activo con case_owner en la DB.
// Se usa cuando el request no especifica "operadorICode".
const operadorICodeDefault = "019e07ac-5ff6-7956-b3a7-f22449c9a38f"

// ── Tipos del request simplificado ───────────────────────────────────────────

type SimulateCaseRequest struct {
	Nombres             string `json:"nombres"`
	Apellidos           string `json:"apellidos"`
	TipoDoc             string `json:"tipoDoc"`        // cc, ce, cd …
	NumDoc              string `json:"numDoc"`
	EsPareja            bool   `json:"esPareja"`       // true → tamizaje pareja íntima
	NivelRiesgo         int    `json:"nivelRiesgo"`    // 1=Bajo 2=Moderado 3=Alto 4=Extremo
	MunicipioAtencion   string `json:"municipioAtencion"`
	MunicipioResidencia string `json:"municipioResidencia"`
	MunicipioHechos     string `json:"municipioHechos"`
	OperadorICode       string `json:"operadorICode"`  // opcional — icode del operador
}

// ── Handler (sin autenticación) ───────────────────────────────────────────────

func SimulateCasePOST(c *gin.Context) {

	// 1. Parsear body simplificado
	buf := new(bytes.Buffer)
	buf.ReadFrom(c.Request.Body)
	var req SimulateCaseRequest
	if err := json.Unmarshal(buf.Bytes(), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Body inválido: " + err.Error()})
		return
	}

	// 2. Defaults
	if req.Nombres == "" { req.Nombres = "María Prueba" }
	if req.Apellidos == "" { req.Apellidos = "Simulación Test" }
	if req.TipoDoc == "" { req.TipoDoc = "cc" }
	if req.NumDoc == "" {
		req.NumDoc = fmt.Sprintf("%d", time.Now().UnixMilli()%9000000000+1000000000)
	}
	if req.NivelRiesgo < 1 || req.NivelRiesgo > 4 { req.NivelRiesgo = 3 }
	if req.OperadorICode == "" { req.OperadorICode = operadorICodeDefault }
	defaultMpio := "11001000" // Bogotá
	if req.MunicipioAtencion == "" { req.MunicipioAtencion = defaultMpio }
	if req.MunicipioResidencia == "" { req.MunicipioResidencia = defaultMpio }
	if req.MunicipioHechos == "" { req.MunicipioHechos = defaultMpio }

	// 3. Sesión mock — simula un operador activo sin pasar por login
	mockSession := utils.CommonSession{
		UserICode:   req.OperadorICode,
		CurrentRole: "op",
		Lang:        "sp",
		Names:       "Simulador",
		LastNames:   "Test",
	}

	// 4. Construir JSON completo del caso
	caseJSON, err := buildSimulatedCaseJSON(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error construyendo payload: " + err.Error()})
		return
	}

	// 5. Llamar al controlador existente (mismo flujo que el front)
	code, res := salvia_ctrl.SetVictimCase(caseJSON, mockSession, dbClientConfig, dbServerConfig)
	c.DataFromReader(code, int64(len(res)), gin.MIMEJSON, strings.NewReader(res), nil)
}

// ── Construcción del payload completo ────────────────────────────────────────

func buildSimulatedCaseJSON(req SimulateCaseRequest) (string, error) {

	// Helpers: obtiene el primer enum de una categoría con el código deseado.
	// Si no lo encuentra, usa el primero de la categoría.
	enumByCode := func(category, code string) salvia_daos.VictimCaseForm2EnumsDTO {
		candidates := salvia_daos.GetLocalVictimCaseForm2EnumsByCategory(
			salvia_daos.VictimCaseForm2EnumsDTO{VictimCaseForm2EnumsCategory: category},
		)
		for _, e := range candidates {
			if e.VictimCaseForm2EnumsCode == code {
				return e
			}
		}
		if len(candidates) > 0 {
			return candidates[0]
		}
		return salvia_daos.VictimCaseForm2EnumsDTO{}
	}

	// yes/no shortcuts
	yes := enumByCode("yes_no", "y")
	no  := enumByCode("yes_no", "n")

	// ── Tamizaje (answers según nivel de riesgo) ──────────────────────────────
	// Distribuye respuestas "sí" para alcanzar el nivel de riesgo pedido.
	//
	// Pareja (pi/ex):  0-4=Nivel1  5-8=Nivel2  9-15=Nivel3  16-24=Nivel4
	// No pareja:       0-2=Nivel1  3-5=Nivel2  6-8=Nivel3   9-20=Nivel4
	//
	// Estrategia: definimos cuántos "sí" poner; los primeros N de la lista de
	// preguntas quedan en "sí" y el resto en "no".

	targetScore := 0
	switch req.NivelRiesgo {
	case 1:
		targetScore = 0
	case 2:
		if req.EsPareja { targetScore = 5 } else { targetScore = 3 }
	case 3:
		if req.EsPareja { targetScore = 9 } else { targetScore = 6 }
	case 4:
		if req.EsPareja { targetScore = 16 } else { targetScore = 9 }
	}

	// yn devuelve yes/no según si idx < targetScore
	idx := 0
	yn := func() salvia_daos.VictimCaseForm2EnumsDTO {
		ans := no
		if idx < targetScore {
			ans = yes
		}
		idx++
		return ans
	}

	// Preguntas comunes (6)
	aggressorViolencePhysicalIncrease := yn()
	aggressorWeaponUsed               := yn()
	aggressorThreatKill               := yn()
	aggressorPursuesSpiesDestroys     := yn()
	aggressorCapableOfKilling         := yn()
	aggressorHasAccessToWeapons       := yn()

	// Preguntas específicas de pareja (18 preguntas)
	var (
		partnerUnemployed                     salvia_daos.VictimCaseForm2EnumsDTO
		partnerOtherDenunciations             salvia_daos.VictimCaseForm2EnumsDTO
		aggressorHasPenalBackground           salvia_daos.VictimCaseForm2EnumsDTO
		stoppedSeekingHelp                    salvia_daos.VictimCaseForm2EnumsDTO
		aggressorForcedSex                    salvia_daos.VictimCaseForm2EnumsDTO
		aggressorAttemptedStrangulation       salvia_daos.VictimCaseForm2EnumsDTO
		aggressorConsumesDrugs                salvia_daos.VictimCaseForm2EnumsDTO
		aggressorIsAlcoholic                  salvia_daos.VictimCaseForm2EnumsDTO
		partnerControls                       salvia_daos.VictimCaseForm2EnumsDTO
		aggressorHadHitInVulnerability        salvia_daos.VictimCaseForm2EnumsDTO
		victimHealthToBlackmail               salvia_daos.VictimCaseForm2EnumsDTO
		partnerThreatenedSuicide              salvia_daos.VictimCaseForm2EnumsDTO
		partnerThreatenedDamageMembers        salvia_daos.VictimCaseForm2EnumsDTO
		thoughtsOfSelfHarm                    salvia_daos.VictimCaseForm2EnumsDTO
		aggressorLimitsContactSupportNetworks salvia_daos.VictimCaseForm2EnumsDTO
		threatenedRevealSexualOrientation     salvia_daos.VictimCaseForm2EnumsDTO
		stillLivesWithAggressor               salvia_daos.VictimCaseForm2EnumsDTO
		aggressorViolentlyJealous             salvia_daos.VictimCaseForm2EnumsDTO

		// No-pareja (14 preguntas)
		aggressorTakenAdvantagePhysicalVulnerability salvia_daos.VictimCaseForm2EnumsDTO
		aggressorUnemployed                          salvia_daos.VictimCaseForm2EnumsDTO
		aggressorHasPenalBackground2                 salvia_daos.VictimCaseForm2EnumsDTO
		aggressorSexuallyHarassment                  salvia_daos.VictimCaseForm2EnumsDTO
		aggressorSexuallyHarassment2                 salvia_daos.VictimCaseForm2EnumsDTO
		violenceMotivatedByGender2                   salvia_daos.VictimCaseForm2EnumsDTO
		aggressorUseDrugs                            salvia_daos.VictimCaseForm2EnumsDTO
		aggressorIsAlcoholic2                        salvia_daos.VictimCaseForm2EnumsDTO
		aggressorControls                            salvia_daos.VictimCaseForm2EnumsDTO
		aggressorThreatenedDamageMembers             salvia_daos.VictimCaseForm2EnumsDTO
		thoughtsOfSelfHarm2                          salvia_daos.VictimCaseForm2EnumsDTO
		aggressorCommonSpaces                        salvia_daos.VictimCaseForm2EnumsDTO
		aggressorHierarchy                           salvia_daos.VictimCaseForm2EnumsDTO
		aggressorUsedPositionAuthority               salvia_daos.VictimCaseForm2EnumsDTO
	)

	// Código de relación con agresor
	relationshipCode := "am" // "am" = otro conocido (no pareja)
	if req.EsPareja {
		relationshipCode = "pi" // pareja íntima → activa tamizaje de pareja
		partnerUnemployed                     = yn()
		partnerOtherDenunciations             = yn()
		aggressorHasPenalBackground           = yn()
		stoppedSeekingHelp                    = yn()
		aggressorForcedSex                    = yn()
		aggressorAttemptedStrangulation       = yn()
		aggressorConsumesDrugs                = yn()
		aggressorIsAlcoholic                  = yn()
		partnerControls                       = yn()
		aggressorHadHitInVulnerability        = yn()
		victimHealthToBlackmail               = yn()
		partnerThreatenedSuicide              = yn()
		partnerThreatenedDamageMembers        = yn()
		thoughtsOfSelfHarm                    = yn()
		aggressorLimitsContactSupportNetworks = yn()
		threatenedRevealSexualOrientation     = yn()
		stillLivesWithAggressor               = yn()
		aggressorViolentlyJealous             = yn()
		// No-pareja quedan en "no"
		aggressorTakenAdvantagePhysicalVulnerability = no
		aggressorUnemployed                          = no
		aggressorHasPenalBackground2                 = no
		aggressorSexuallyHarassment                  = no
		aggressorSexuallyHarassment2                 = no
		violenceMotivatedByGender2                   = no
		aggressorUseDrugs                            = no
		aggressorIsAlcoholic2                        = no
		aggressorControls                            = no
		aggressorThreatenedDamageMembers             = no
		thoughtsOfSelfHarm2                          = no
		aggressorCommonSpaces                        = no
		aggressorHierarchy                           = no
		aggressorUsedPositionAuthority               = no
	} else {
		// Pareja quedan en "no"
		partnerUnemployed                     = no
		partnerOtherDenunciations             = no
		aggressorHasPenalBackground           = no
		stoppedSeekingHelp                    = no
		aggressorForcedSex                    = no
		aggressorAttemptedStrangulation       = no
		aggressorConsumesDrugs                = no
		aggressorIsAlcoholic                  = no
		partnerControls                       = no
		aggressorHadHitInVulnerability        = no
		victimHealthToBlackmail               = no
		partnerThreatenedSuicide              = no
		partnerThreatenedDamageMembers        = no
		thoughtsOfSelfHarm                    = no
		aggressorLimitsContactSupportNetworks = no
		threatenedRevealSexualOrientation     = no
		stillLivesWithAggressor               = no
		aggressorViolentlyJealous             = no
		// No-pareja activas
		aggressorTakenAdvantagePhysicalVulnerability = yn()
		aggressorUnemployed                          = yn()
		aggressorHasPenalBackground2                 = yn()
		aggressorSexuallyHarassment                  = yn()
		aggressorSexuallyHarassment2                 = yn()
		violenceMotivatedByGender2                   = yn()
		aggressorUseDrugs                            = yn()
		aggressorIsAlcoholic2                        = yn()
		aggressorControls                            = yn()
		aggressorThreatenedDamageMembers             = yn()
		thoughtsOfSelfHarm2                          = yn()
		aggressorCommonSpaces                        = yn()
		aggressorHierarchy                           = yn()
		aggressorUsedPositionAuthority               = yn()
	}

	// ── Enums fijos ───────────────────────────────────────────────────────────
	docTypeEnum        := enumByCode("victim_case_form2_victim_doc_type", req.TipoDoc)
	factsZone          := enumByCode("victim_case_form2_facts_zone", "cm")        // cabecera municipal
	residenceZone      := enumByCode("victim_case_form2_facts_zone", "cm")
	scenarioViolence   := enumByCode("victim_case_form2_scenario_violence", "ho") // hogar
	recurrenceAgg      := enumByCode("victim_case_form2_recurrence_aggression", "re") // reiterativa
	numAgressors       := enumByCode("victim_case_form2_num_agressors", "01")
	// Proximidad: pareja conviviente si esPareja, persona diferente si no
	proximityCode := "pd"
	if req.EsPareja { proximityCode = "pc" }
	proximity          := enumByCode("victim_case_form2_proximity_principal_aggressor", proximityCode)
	relationship       := enumByCode("victim_case_form2_relationship_with_presumed_aggressor_01", relationshipCode)
	aggressorGender    := enumByCode("victim_case_form2_aggressor_gender_identity", "ho") // hombre
	aggressorOccup     := enumByCode("victim_case_form2_occupation", "ni")               // ninguna
	economicallyDep    := no
	if req.EsPareja { economicallyDep = no }

	nationality        := enumByCode("victim_case_form2_nationality", "co")          // colombiana
	genderIdentity     := enumByCode("victim_case_form2_gender_identity", "mu")      // mujer
	sexualOrientation  := enumByCode("victim_case_form2_sexual_orientation", "he")   // heterosexual
	assignedSex        := enumByCode("victim_case_form2_assigned_sex_at_birth", "mu")
	ethnicAffiliation  := enumByCode("victim_case_form2_ethnic_affiliation", "ne")   // ninguna
	maritalStatus      := enumByCode("victim_case_form2_marital_status", "so")       // soltera
	educationLevel     := enumByCode("victim_case_form2_last_education_level", "bm") // bachillerato
	occupation         := enumByCode("victim_case_form2_occupation", "de")           // desempleada (evita requerir employmentRelationship)
	incomeMethod       := enumByCode("victim_case_form2_income_generation_method", "de") // desempleada
	employmentRel      := enumByCode("victim_case_form2_employment_relationship", "") // primer valor disponible (por si income="em")
	housingTenancy     := enumByCode("victim_case_form2_housing_tenancy_form", "ar") // arrendada
	housingStratum     := enumByCode("victim_case_form2_housing_stratum", "02")      // estrato 2

	// Campos multi-selección obligatorios (siempre ≥ 1 item)
	adjustmentsGBVEnum          := enumByCode("victim_case_form2_adjustments_gbv", "na")          // ninguna
	speciallyProtectedPopEnum   := enumByCode("victim_case_form2_specially_protected_population", "na") // ninguna

	typeViolence       := enumByCode("victim_case_form2_type_violence_experienced", "fi") // física
	subtypeViolence    := enumByCode("victim_case_form2_subtype_violence_experienced_fi", "ac")
	scopeViolence      := enumByCode("victim_case_form2_scope_of_violence", "ac")
	actionPlan         := enumByCode("victim_case_form2_action_plan", "ap")              // atención psicosocial
	hasDependents      := enumByCode("victim_case_form2_has_dependents", "nn")           // ninguno

	// Fechas
	now      := time.Now()
	today    := now.Format("02/01/2006")
	timeNow  := now.Format("15:04")
	birthDate := "15/03/1990"

	// ── Estructura del payload ────────────────────────────────────────────────
	type enumField struct {
		ICode string `json:"icode"`
		Code  string `json:"code"`
	}
	toEF := func(e salvia_daos.VictimCaseForm2EnumsDTO) enumField {
		return enumField{ICode: e.VictimCaseForm2EnumsICode, Code: e.VictimCaseForm2EnumsCode}
	}

	payload := map[string]interface{}{
		"victimCase": map[string]interface{}{
			"names":     req.Nombres,
			"lastNames": req.Apellidos,
			"docType":   req.TipoDoc,
			"docNumber": req.NumDoc,
			"townCode":  req.MunicipioAtencion,
			"entityBranchesByMomentIdx": map[string]interface{}{},
			"form2": map[string]interface{}{
				// ── Identidad / contacto ──────────────────────────────────────
				"identityName":   req.Nombres,
				"phone":          3001234567,
				"factsDate":      today,
				"factsStartTime": timeNow,
				"factsDescription": fmt.Sprintf(
					"Caso simulado generado automáticamente. Nivel de riesgo %d. Agresor %s.",
					req.NivelRiesgo,
					map[bool]string{true: "pareja íntima", false: "no pareja"}[req.EsPareja],
				),
				"factsZone":    toEF(factsZone),
				"factsTownCode": req.MunicipioHechos,
				"factsAddress": "Calle 1 # 2-3",

				// ── Residencia ────────────────────────────────────────────────
				"residenceZone":    toEF(residenceZone),
				"residenceTownCode": req.MunicipioResidencia,
				"residenceAddress": "Calle 4 # 5-6",

				// ── Accesibilidad ─────────────────────────────────────────────
				"personWithDisability":       toEF(no),
				"requireLanguageInterpreter": toEF(no),
				"adjustmentsGBV":             []interface{}{toEF(adjustmentsGBVEnum)}, // ≥1 item obligatorio

				// ── Violencia ─────────────────────────────────────────────────
				"scenarioViolence":       toEF(scenarioViolence),
				"typeViolenceExperienced": []interface{}{toEF(typeViolence)},
				"subtypeViolenceExperienced": []interface{}{toEF(subtypeViolence)},
				"scopeOfViolence":         []interface{}{toEF(scopeViolence)},
				"violenceMotivatedByGender": toEF(yes),
				"reportedPreviously":      toEF(no),
				"whoReportTo":             []interface{}{},
				"attentionWasAppropriate": toEF(yes),
				"recurrenceAggression":    toEF(recurrenceAgg),

				// ── Agresor ───────────────────────────────────────────────────
				"numAgressors":                       toEF(numAgressors),
				"proximityPrincipalAggressor":        toEF(proximity),
				"relationshipWithPresumedAggressor":  toEF(relationship),
				"economicallyDependent":              toEF(economicallyDep),
				"aggressorGenderIdentity":            toEF(aggressorGender),
				"aggressorOccupation":                toEF(aggressorOccup),
				"aggressorNames":    "Agresor Simulado",
				"aggressorDocType":  toEF(docTypeEnum),
				"aggressorDocNumber": "9999999999",
				"aggressorAddress":  "",
				"aggressorPhone":    0,

				// ── Tamizaje — preguntas comunes ──────────────────────────────
				"aggressorViolencePhysicalIncrease": toEF(aggressorViolencePhysicalIncrease),
				"aggressorWeaponUsed":               toEF(aggressorWeaponUsed),
				"aggressorThreatKill":               toEF(aggressorThreatKill),
				"aggressorPursuesSpiesDestroys":     toEF(aggressorPursuesSpiesDestroys),
				"aggressorCapableOfKilling":         toEF(aggressorCapableOfKilling),
				"aggressorHasAccessToWeapons":       toEF(aggressorHasAccessToWeapons),

				// ── Tamizaje — pareja ─────────────────────────────────────────
				"partnerUnemployed":                     toEF(partnerUnemployed),
				"partnerOtherDenunciations":             toEF(partnerOtherDenunciations),
				"aggressorHasPenalBackground":           toEF(aggressorHasPenalBackground),
				"stoppedSeekingHelp":                    toEF(stoppedSeekingHelp),
				"aggressorForcedSex":                    toEF(aggressorForcedSex),
				"aggressorAttemptedStrangulation":       toEF(aggressorAttemptedStrangulation),
				"aggressorConsumesDrugs":                toEF(aggressorConsumesDrugs),
				"aggressorIsAlcoholic":                  toEF(aggressorIsAlcoholic),
				"partnerControls":                       toEF(partnerControls),
				"aggressorHadHitInVulnerability":        toEF(aggressorHadHitInVulnerability),
				"victimHealthToBlackmail":               toEF(victimHealthToBlackmail),
				"partnerThreatenedSuicide":              toEF(partnerThreatenedSuicide),
				"partnerThreatenedDamageMembers":        toEF(partnerThreatenedDamageMembers),
				"thoughtsOfSelfHarm":                    toEF(thoughtsOfSelfHarm),
				"aggressorLimitsContactSupportNetworks": toEF(aggressorLimitsContactSupportNetworks),
				"threatenedRevealSexualOrientation":     toEF(threatenedRevealSexualOrientation),
				"stillLivesWithAggressor":               toEF(stillLivesWithAggressor),
				"aggressorViolentlyJealous":             toEF(aggressorViolentlyJealous),

				// ── Tamizaje — no pareja ──────────────────────────────────────
				"aggressorTakenAdvantagePhysicalVulnerability": toEF(aggressorTakenAdvantagePhysicalVulnerability),
				"aggressorUnemployed":                          toEF(aggressorUnemployed),
				"aggressorHasPenalBackground2":                 toEF(aggressorHasPenalBackground2),
				"aggressorSexuallyHarassment":                  toEF(aggressorSexuallyHarassment),
				"aggressorSexuallyHarassment2":                 toEF(aggressorSexuallyHarassment2),
				"violenceMotivatedByGender2":                   toEF(violenceMotivatedByGender2),
				"aggressorUseDrugs":                            toEF(aggressorUseDrugs),
				"aggressorIsAlcoholic2":                        toEF(aggressorIsAlcoholic2),
				"aggressorControls":                            toEF(aggressorControls),
				"aggressorThreatenedDamageMembers":             toEF(aggressorThreatenedDamageMembers),
				"thoughtsOfSelfHarm2":                          toEF(thoughtsOfSelfHarm2),
				"aggressorCommonSpaces":                        toEF(aggressorCommonSpaces),
				"aggressorHierarchy":                           toEF(aggressorHierarchy),
				"aggressorUsedPositionAuthority":               toEF(aggressorUsedPositionAuthority),

				// ── Score (el backend recalcula y verifica) ───────────────────
				"riskScore": targetScore,
				"riskLevel": req.NivelRiesgo,

				// ── Datos personales ──────────────────────────────────────────
				"birthDate":                       birthDate,
				"physicalMentalSensoryDifficulties": toEF(no),
				"activitiesUnableToHear":          0,
				"activitiesUnableToTalk":          0,
				"activitiesUnableToSee":           0,
				"activitiesUnableToMove":          0,
				"activitiesUnableToTake":          0,
				"activitiesUnableToUnderstand":    0,
				"activitiesUnableToEat":           0,
				"activitiesUnableToInteract":      0,
				"activitiesUnableToDoEveryday":    0,
				"activitiesUnableToPerform":       []interface{}{},
				"law1996":                         []interface{}{},
				"nationality":                     toEF(nationality),
				"specifiedNationality":            toEF(salvia_daos.VictimCaseForm2EnumsDTO{}),
				"migrationCondition":              toEF(salvia_daos.VictimCaseForm2EnumsDTO{}),
				"genderIdentity":                  toEF(genderIdentity),
				"sexualOrientation":               toEF(sexualOrientation),
				"assignedSexAtBirth":              toEF(assignedSex),
				"speciallyProtectedPopulation":    []interface{}{toEF(speciallyProtectedPopEnum)}, // ≥1 item obligatorio
				"ethnicAffiliation":               toEF(ethnicAffiliation),
				"indigenousPeople":                toEF(salvia_daos.VictimCaseForm2EnumsDTO{}),
				"campesinoRecognition":             toEF(no),
				"maritalStatus":                   toEF(maritalStatus),
				"lastEducationLevel":              toEF(educationLevel),
				"occupation":                      toEF(occupation),
				"incomeGenerationMethod":          toEF(incomeMethod),
				"employmentRelationship":          toEF(employmentRel),             // primer valor disponible
				"aspMode":                         []interface{}{},
				"approxStartAsp":                  today,                           // fecha válida por si income="pr"
				"reasonASP":                       []interface{}{},
				"housingTenancyForm":              toEF(housingTenancy),
				"housingStratum":                  toEF(housingStratum),
				"hasDependents":                   []interface{}{toEF(hasDependents)},
				"currentlyPregnant":               toEF(no),

				// ── Plan de acción ────────────────────────────────────────────
				"actionPlan":                 []interface{}{toEF(actionPlan)},
				"salivaManagementExplanation": "Caso simulado automáticamente para pruebas.",
				"allowsEasyReport":           toEF(yes),

				// ── Contacto de apoyo (opcional) ──────────────────────────────
				"supportContactNames":   "",
				"supportContactPhone":   0,
				"supportContactEmail":   "",
				"supportContactKinship": toEF(salvia_daos.VictimCaseForm2EnumsDTO{}),
			},
		},
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}
