// Package salvia_facades contiene las funciones de fachada para la gestión de casos de víctimas,
// integrando la lógica de control, seguridad, y renderizado de plantillas.
package salvia_facades

import (
	"bytes"
	"math"
	"net/http"
	"strconv"
	"strings"

	common_config "bitsflow/common/config"
	"bitsflow/common/db"
	common_facades "bitsflow/common/facades"
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	security_config "bitsflow/security/config"
	security_ctrl "bitsflow/security/controllers"

	salvia_ctrl "bitsflow/salvia/controllers"
	salvia_daos "bitsflow/salvia/dao"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// FeminicideRiskPOST maneja la solicitud POST para crear o actualizar un caso de víctima.
// Obtiene la sesión actual, verifica los permisos necesarios, lee el cuerpo de la solicitud,
// y llama al controlador correspondiente para procesar el caso de víctima.
func FeminicideRiskPOST(c *gin.Context) {
	// Se obtiene la sesión del usuario desde el contexto de Gin.
	session := sessions.Default(c)
	// Se extrae el ID de sesión almacenado en "userData".
	var sessionID string = session.Get("userData").(string)
	// Se obtiene la sesión común asociada al ID.
	s, _ := utils.GetCommonSession(sessionID)

	// Verifica que el usuario tenga el permiso "set_feminicide_risk".
	if !utils.CheckPermission(salvia_config.PermissionsByRole, "set_feminicide_risk", s.CurrentRole, c) {
		return
	}

	// Lee el cuerpo de la solicitud en un buffer.
	buf := new(bytes.Buffer)
	buf.ReadFrom(c.Request.Body)
	// Llama al controlador para establecer el caso de víctima, pasando el contenido del buffer,
	// el origen "salvia", la sesión y la configuración de la base de datos.
	code, res := salvia_ctrl.SetFeminicideRisk(buf.String(), *s, dbClientConfig, dbServerConfig)
	// Envía la respuesta al cliente en formato JSON.
	c.DataFromReader(code, int64(len(res)), gin.MIMEJSON, strings.NewReader(res), nil)
}

// FeminicideRiskGET maneja la solicitud GET para obtener casos de víctima.
// Dependiendo de los parámetros de la URL y el rol del usuario, recupera información
// detallada de un caso específico o una lista de casos, y renderiza la respuesta en JSON o HTML.
func FeminicideRiskGET(c *gin.Context) {
	// Se obtiene la sesión del usuario y el ID de sesión.
	session := sessions.Default(c)
	var sessionID string = session.Get("userData").(string)
	s, err := utils.GetCommonSession(sessionID)

	// Variables para almacenar las respuestas y otros datos relevantes.
	var feminicideRiskRes string = "[]"
	var code int = http.StatusOK
	var tplName string
	var departments string = "{}"

	var count int

	// Establece cabeceras para evitar cache en la respuesta.
	common_facades.SetHeaderNoCache(c)

	_, departments = security_ctrl.GetDepartmentByAll(&db.ConnData{}, dbClientConfig, dbServerConfig)

	// Se obtiene el parámetro "id" de la URL.
	id := c.Param("id")
	page, _ := strconv.Atoi(c.Param("p"))

	//Si se consulta un caso particular
	if id != "" {
		// Si no se especifica "docType", se verifica el permiso "get_feminicide_risk".
		if !utils.CheckPermission(salvia_config.PermissionsByRole, "get_feminicide_risk", s.CurrentRole, c) {
			return
		}
		// Se obtiene el caso por su código interno.
		code, feminicideRiskRes, _ = salvia_ctrl.GetFeminicideRiskByICode(id, &db.ConnData{}, dbClientConfig, dbServerConfig)

	} else {
		//Si se consultan todos los casos
		// Si no se proporciona un ID, se verifica el permiso "get_feminicide_risks" para obtener todos los casos.
		if !utils.CheckPermission(salvia_config.PermissionsByRole, "get_feminicide_risks", s.CurrentRole, c) {
			return
		}
		code, feminicideRiskRes, count = salvia_ctrl.GetFeminicideRiskByAll(page, &db.ConnData{}, dbClientConfig, dbServerConfig)
	}

	// Se decide el formato de respuesta según el header "Accept" de la solicitud.
	if v, found := c.Request.Header["Accept"]; found && v[0] == "application/json" {
		// Respuesta en formato JSON.
		c.DataFromReader(code, int64(len(feminicideRiskRes)), gin.MIMEJSON, strings.NewReader(feminicideRiskRes), nil)
	} else {
		// Se prepara la información necesaria para renderizar la plantilla HTML.
		var menu map[string][]map[string]string
		var feminicideRiskMenuTools []map[string]string = []map[string]string{}

		if err == nil {
			menu = s.CurrentMenu
		}
		if v, found := salvia_config.MenuTools[s.Lang][s.CurrentRole]["menu_tool_get_fermincides"]; found {
			feminicideRiskMenuTools = v
		}

		// Renderiza la plantilla HTML con todos los parámetros necesarios para mostrar el caso de víctima.
		common_facades.RenderTemplate(c, salvia_daos.VictimContactEntityName, "salvia", "feminicideRisk/", salvia_config.HTML_Templates, tplName, utils.GetFullHtmlTemplates(), utils.DEFAULT_VIEW, utils.DEFAULT_PANIC_TEMPLATE,
			map[string]interface{}{
				"windowTitle":     salvia_config.Locale["sp"]["get_feminicide_risk_window_title"],
				"currentUser":     s.Names + " " + s.LastNames,
				"nav_rules":       salvia_config.TranslateNavigationRule(s.Lang, salvia_config.NAVIGATION_RULES["get_feminicide_risk"]),
				"locale":          salvia_config.Locale,
				"lang":            s.Lang,
				"feminicideRisks": feminicideRiskRes,
				"numPages":        int(math.Ceil(float64(count) / float64(salvia_config.NUM_ITEMS_PER_PAGE))),

				"departments":                  departments,
				"menu":                         menu,
				"menuToolsFeminicideRiskTable": feminicideRiskMenuTools,
				"momentCodes":                  salvia_config.MOMENT[s.Lang],
				"sectors":                      salvia_config.SECTOR[s.Lang],
				"salviaFeminicideRiskFormPath": salvia_config.FormPaths[s.Lang]["FeminicideRiskGET"],
			}, utils.GetFullHtmlFuncMap())
	}
}

// FeminicideRiskReportGET maneja la solicitud GET para la generación de reportes de casos de víctimas.
// Recupera información necesaria, como los departamentos, y renderiza la plantilla correspondiente.
func FeminicideRiskReportGET(c *gin.Context) {
	// Se obtiene la sesión y se extrae el ID de sesión.
	session := sessions.Default(c)
	var sessionID string = session.Get("userData").(string)
	s, err := utils.GetCommonSession(sessionID)
	// Establece cabeceras para evitar cache.
	common_facades.SetHeaderNoCache(c)
	var menu map[string][]map[string]string

	// Se obtiene el menú actual de la sesión si no hay error.
	if err == nil {
		menu = s.CurrentMenu
	}
	var departments string
	// Se obtienen todos los departamentos desde el controlador de seguridad.
	_, departments = security_ctrl.GetDepartmentByAll(&db.ConnData{}, dbClientConfig, dbServerConfig)

	// Renderiza la plantilla HTML para el reporte de casos de víctimas.
	common_facades.RenderTemplate(c, salvia_daos.VictimContactEntityName, "salvia", "feminicideRisk/", salvia_config.HTML_Templates, "report_feminicideRisks", utils.GetFullHtmlTemplates(), utils.REPORT_VIEW, utils.DEFAULT_PANIC_TEMPLATE,
		map[string]interface{}{
			"windowTitle":          salvia_config.Locale["sp"]["report_feminicideRisk_window_title"],
			"currentUser":          s.Names + " " + s.LastNames,
			"nav_rules":            salvia_config.TranslateNavigationRule(s.Lang, salvia_config.NAVIGATION_RULES["report_feminicideRisks"]),
			"locale":               salvia_config.Locale,
			"lang":                 s.Lang,
			"menu":                 menu,
			"departments":          departments,
			"docType":              common_config.DOCUMENT_TYPE,
			"violenceExperienced":  salvia_config.VIOLENCE_EXPERIENCED,
			"caseStatus":           salvia_config.VICTIM_CASE_STATUS[s.Lang],
			"salviaAlertFormPath":  salvia_config.FormPaths[s.Lang]["AlertGET"],
			"securityCityFormPath": security_config.FormPaths[s.Lang]["CityGET"],
			"securityTownFormPath": security_config.FormPaths[s.Lang]["TownGET"],
			"salviaReportFormPath": salvia_config.FormPaths[s.Lang]["FeminicideRiskReportGET"],
		}, utils.GetFullHtmlFuncMap())
}

// FeminicideRiskPUT_GET maneja la solicitud GET para obtener la información de un caso de víctima a actualizar.
// Recupera datos relacionados con el caso, contacto, ubicación y entidades, y renderiza la plantilla de actualización.
/*func FeminicideRiskPUT_GET(c *gin.Context) {
	// Se obtiene la sesión y se extrae el ID de sesión.
	session := sessions.Default(c)
	var sessionID string = session.Get("userData").(string)
	s, _ := utils.GetCommonSession(sessionID)

	// Verifica que el usuario tenga el permiso "update_feminicideRisk".
	if !utils.CheckPermission(salvia_config.PermissionsByRole, "update_feminicideRisk", s.CurrentRole, c) {
		return
	}

	// Declaración de variables para almacenar información del caso y contacto.
	var victimContact salvia_daos.VictimContactDTO = salvia_daos.VictimContactDTO{}
	var feminicideRisk salvia_daos.FeminicideRiskDTO = salvia_daos.FeminicideRiskDTO{}

	// Variables para almacenar datos geográficos.
	var departments string = "{}"
	var cities string = "{}"
	var towns string = "{}"

	var cities2 string = "{}"
	var towns2 string = "{}"

	var violenceCities string = "{}"
	var violenceTowns string = "{}"

	// Lista de entidades disponibles.
	var entities []salvia_daos.EntityDTO = []salvia_daos.EntityDTO{}
	// Mapa para agrupar ramas de entidad por momento y sector.
	var entitiesByMomentAndSector map[string]map[string]map[string][]salvia_daos.EntityBranchDTO = make(map[string]map[string]map[string][]salvia_daos.EntityBranchDTO)

	// Establece cabeceras para evitar cache.
	common_facades.SetHeaderNoCache(c)

	// Se obtienen los departamentos y las entidades.
	_, departments = security_ctrl.GetDepartmentByAll(&db.ConnData{}, dbClientConfig, dbServerConfig)
	_, _, entities = salvia_ctrl.GetEntityByAll(&db.ConnData{}, dbClientConfig, dbServerConfig)

	// Se obtiene el ID del caso de la URL.
	id := c.Param("id")
	if id != "" {
		// Se recupera la información del caso y del contacto asociado.
		_, _, feminicideRisk = salvia_ctrl.GetFeminicideRiskByICode(id, &db.ConnData{}, dbClientConfig, dbServerConfig)
		_, _, victimContact = salvia_ctrl.GetVictimContactByICode(strconv.FormatUint(feminicideRisk.FeminicideRiskVictimContact.VictimContactId, 10), &db.ConnData{}, dbClientConfig, dbServerConfig)
	}

	// Carga datos adicionales en feminicideRisk, excluyendo información redundante del contacto.
	feminicideRisk.LoadFromWithoutVictimContact(salvia_config.MOMENT[s.Lang], salvia_config.SECTOR[s.Lang], entities)

	// Para el lugar de los hechos, se obtienen las ciudades y pueblos basados en el departamento y ciudad seleccionados.
	if feminicideRisk.FeminicideRiskForm2.FeminicideRiskForm2FactsDepartment.DepartmentId > 0 {
		_, cities2, _ = security_ctrl.GetCitiesByDeparment(feminicideRisk.FeminicideRiskDepartment.DepartmentId, &db.ConnData{}, dbClientConfig, dbServerConfig)
	}
	if feminicideRisk.FeminicideRiskTown.TownId > 0 {
		_, towns2 = security_ctrl.GetTownsByCity(feminicideRisk.FeminicideRiskCity.CityId, &db.ConnData{}, dbClientConfig, dbServerConfig)
	}

	// Para el lugar de atención, se obtienen las ciudades y pueblos correspondientes.
	if feminicideRisk.FeminicideRiskDepartment.DepartmentId > 0 {
		_, cities, _ = security_ctrl.GetCitiesByDeparment(feminicideRisk.FeminicideRiskDepartment.DepartmentId, &db.ConnData{}, dbClientConfig, dbServerConfig)
	}
	if feminicideRisk.FeminicideRiskTown.TownId > 0 {
		_, towns = security_ctrl.GetTownsByCity(feminicideRisk.FeminicideRiskCity.CityId, &db.ConnData{}, dbClientConfig, dbServerConfig)

		// Se obtienen las ramas de entidad asociadas al pueblo para asignar la ruta.
		var branches []salvia_daos.EntityBranchDTO
		_, _, branches = salvia_ctrl.GetEntityBranchesByTownCodeWithMoments(feminicideRisk.FeminicideRiskTown.TownCode, &db.ConnData{}, dbClientConfig, dbServerConfig)

		for _, e := range branches {
			//if _, found := momentsWithSector[e.EntityBranchMoment][e.EntityBranchId]; found {
			if entitiesByMomentAndSector[e.EntityBranchMoment] == nil {
				entitiesByMomentAndSector[e.EntityBranchMoment] = map[string]map[string][]salvia_daos.EntityBranchDTO{}
			}
			if entitiesByMomentAndSector[e.EntityBranchMoment][e.EntityBranchSector] == nil {
				entitiesByMomentAndSector[e.EntityBranchMoment][e.EntityBranchSector] = map[string][]salvia_daos.EntityBranchDTO{}
				entitiesByMomentAndSector[e.EntityBranchMoment][e.EntityBranchSector][e.EntityBranchEntityICode] = []salvia_daos.EntityBranchDTO{}
			}
			entitiesByMomentAndSector[e.EntityBranchMoment][e.EntityBranchSector][e.EntityBranchEntityICode] =
				append(entitiesByMomentAndSector[e.EntityBranchMoment][e.EntityBranchSector][e.EntityBranchEntityICode], e)
		}
	}

	// Renderiza la plantilla HTML para actualizar el caso de víctima, pasando todos los datos recopilados.
	common_facades.RenderTemplate(c, salvia_daos.FeminicideRiskEntityName, "salvia", "feminicideRisk/", salvia_config.HTML_Templates, "update_feminicideRisk", utils.GetFullHtmlTemplates(), utils.DEFAULT_VIEW, utils.DEFAULT_PANIC_TEMPLATE,
		map[string]interface{}{
			"windowTitle":               salvia_config.Locale["sp"]["update_feminicideRisk_window_title"],
			"currentUser":               s.Names + " " + s.LastNames,
			"feminicideRisk":                feminicideRisk,
			"victimContact":             victimContact,
			"gender":                    common_config.GENDER_IDENTITY,
			"docType":                   common_config.DOCUMENT_TYPE,
			"language":                  common_config.LANGUAGE,
			"locale":                    salvia_config.Locale,
			"lang":                      s.Lang,
			"livingZone":                salvia_config.LIVING_ZONES,
			"genderIdentity":            common_config.GENDER_IDENTITY,
			"sexualOrientation":         salvia_config.SEXUAL_ORIENTATION,
			"origin":                    salvia_config.ORIGIN_PLACE,
			"occupation":                salvia_config.OCCUPATION,
			"ethnicGroup":               salvia_config.ETHNIC_GROUP,
			"maritalStatus":             common_config.MARITAL_STATUS,
			"disability":                salvia_config.DISABILITY,
			"occurrence":                salvia_config.OCCURRENCE,
			"violenceScope":             salvia_config.VIOLENCE_SCOPE,
			"aggressor":                 salvia_config.AGGRESSOR,
			"femicideRisk":              common_config.YES_NO,
			"violenceExperienced":       salvia_config.VIOLENCE_EXPERIENCED,
			"relationshipWithAggressor": salvia_config.RELATIONSHIP_WITH_AGGRESSOR,
			"weekDay":                   salvia_config.WEEK_DAY,
			"yes_no_na":                 common_config.YES_NO_NA,

			"feminicideRiskVictimNationality":                         salvia_config.VICTIM_CASE_VICTIM_NATIONALITY,
			"feminicideRiskVictimForeignerImmigrationStatus":          salvia_config.VICTIM_CASE_VICTIM_FOREIGNER_IMMIGRATION_STATUS,
			"feminicideRiskVictimGender":                              salvia_config.VICTIM_CASE_VICTIM_GENDER,
			"feminicideRiskVictimDependents":                          salvia_config.VICTIM_CASE_VICTIM_DEPENDENTS,
			"feminicideRiskVictimDeathThreats":                        common_config.YES_NO,
			"feminicideRiskVictimAggressorHasWeapons":                 common_config.YES_NO,
			"feminicideRiskExperiencedPhysicalOrSexualViolenceBefore": common_config.YES_NO,
			"feminicideRiskVictimImminentRisk":                        common_config.YES_NO,
			"feminicideRiskVictimPreviouslyReportedSituation":         common_config.YES_NO,
			"feminicideRiskVictimIfPreviouslyReported":                salvia_config.VICTIM_CASE_VICTIM_IF_PREVIOUSLY_REPORTED,

			"ifAfro":             salvia_config.AFRO_COMMUNTITIES[s.Lang],
			"ifIndigenous":       salvia_config.COLOMBIAN_INDIGENOUS[s.Lang],
			"ifIndigenousTongue": salvia_config.INDIGENOUS_TONGUES[s.Lang],
			"ifPeasant":          common_config.YES_NO,
			"ifArmedConflict":    common_config.YES_NO,
			"violenceScene":      salvia_config.VIOLENCE_SCENES[s.Lang],

			"moments":                    salvia_config.MOMENT[s.Lang],
			"sectors":                    salvia_config.SECTOR[s.Lang],
			"entities":                   entities,
			"entitiesByMomentAndSector":  entitiesByMomentAndSector,
			"salviaFormPath":             salvia_config.FormPaths[s.Lang]["FeminicideRiskPUT"],
			"securityCityFormPath":       security_config.FormPaths[s.Lang]["CityGET"],
			"salviaEntityBranchFormPath": salvia_config.FormPaths[s.Lang]["EntityBranchGET"],
			"securityTownFormPath":       security_config.FormPaths[s.Lang]["TownGET"],
			"departments":                departments,
			"cities":                     cities,
			"towns":                      towns,
			"cities2":                    cities2,
			"towns2":                     towns2,
			"violenceCities":             violenceCities,
			"violenceTowns":              violenceTowns,
			"nav_rules":                  salvia_config.TranslateNavigationRule(s.Lang, salvia_config.NAVIGATION_RULES["update_feminicideRisk"]),
		}, utils.GetFullHtmlFuncMap())
}*/

// FeminicideRiskPOST_GET maneja la solicitud GET para inicializar la creación de un caso de víctima.
// Obtiene la información de contacto de la víctima, carga datos por defecto para el caso,
// y renderiza la plantilla correspondiente para la creación del caso.
func FeminicideRiskPOST_GET(c *gin.Context) {
	// Se obtiene la sesión y el ID de sesión.
	session := sessions.Default(c)
	var sessionID string = session.Get("userData").(string)
	s, _ := utils.GetCommonSession(sessionID)

	// Verifica que el usuario tenga el permiso "set_feminicide_risk".
	if !utils.CheckPermission(salvia_config.PermissionsByRole, "set_feminicide_risk", s.CurrentRole, c) {
		return
	}

	id := c.Param("id")

	// Variables para almacenar el caso, contacto y datos geográficos.
	var departments string = "{}"
	var feminicideRiskEmpty salvia_daos.FeminicideRiskDTO = salvia_daos.FeminicideRiskDTO{}

	//Debe venir el id de un caso
	if id != "" {
		_, departments = security_ctrl.GetDepartmentByAll(&db.ConnData{}, dbClientConfig, dbServerConfig)
		feminicideRiskEmpty.FeminicideRiskVictimCase.VictimCaseICode = id
		// Establece cabeceras para evitar cache.
		common_facades.SetHeaderNoCache(c)
	} else {
		c.DataFromReader(http.StatusBadRequest, 0, gin.MIMEJSON, strings.NewReader(""), nil)
		return
	}

	// Se obtienen los departamentos y las entidades.

	// Renderiza la plantilla HTML para la creación de un nuevo caso de víctima.
	common_facades.RenderTemplate(c, salvia_daos.FeminicideRiskEntityName, "salvia", "feminicide/", salvia_config.HTML_Templates, "set_feminicide_risk", utils.GetFullHtmlTemplates(), utils.DEFAULT_VIEW, utils.DEFAULT_PANIC_TEMPLATE,
		map[string]interface{}{
			"windowTitle": salvia_config.Locale["sp"]["set_feminicide_risk_window_title"],
			"currentUser": s.Names + " " + s.LastNames,
			"gender":      common_config.GENDER_IDENTITY,
			"language":    common_config.LANGUAGE,
			"locale":      salvia_config.Locale,
			"lang":        s.Lang,

			"origin":      salvia_config.ORIGIN_PLACE,
			"ethnicGroup": salvia_config.ETHNIC_GROUP,

			"disability":                salvia_config.DISABILITY,
			"occurrence":                salvia_config.OCCURRENCE,
			"violenceScope":             salvia_config.VIOLENCE_SCOPE,
			"aggressor":                 salvia_config.AGGRESSOR,
			"femicideRisk":              common_config.YES_NO,
			"violenceExperienced":       salvia_config.VIOLENCE_EXPERIENCED,
			"relationshipWithAggressor": salvia_config.RELATIONSHIP_WITH_AGGRESSOR,
			"weekDay":                   salvia_config.WEEK_DAY,
			"feminicideRisk":            feminicideRiskEmpty,

			//New form ------
			//Régimen afiliación al SGSSS
			"livingZone":                   salvia_daos.VictimCaseForm2Enums["victim_case_form2_facts_zone"],
			"maritalStatus":                salvia_daos.VictimCaseForm2Enums["victim_case_form2_marital_status"],
			"assignedSexAtBirth":           salvia_daos.VictimCaseForm2Enums["victim_case_form2_assigned_sex_at_birth"],
			"genderIdentity":               salvia_daos.VictimCaseForm2Enums["victim_case_form2_gender_identity"],
			"sexualOrientation":            salvia_daos.VictimCaseForm2Enums["victim_case_form2_sexual_orientation"],
			"ethnicAffiliation":            salvia_daos.VictimCaseForm2Enums["victim_case_form2_ethnic_affiliation"],
			"indigenousPeople":             salvia_daos.VictimCaseForm2Enums["victim_case_form2_indigenous_people"],
			"migrationCondition":           salvia_daos.VictimCaseForm2Enums["victim_case_form2_migration_condition"],
			"lastEducationLevel":           salvia_daos.VictimCaseForm2Enums["victim_case_form2_last_education_level"],
			"speciallyProtectedPopulation": salvia_daos.VictimCaseForm2Enums["victim_case_form2_specially_protected_population"],

			"yes_no":    salvia_daos.VictimCaseForm2Enums["yes_no"],
			"docType":   salvia_daos.VictimCaseForm2Enums["feminicide_risk_form2_victim_doc_type"],
			"factsZone": salvia_daos.VictimCaseForm2Enums["feminicide_risk_form2_facts_zone"],

			//Campos múltiples

			"moments":                    salvia_config.MOMENT[s.Lang],
			"sectors":                    salvia_config.SECTOR[s.Lang],
			"salviaFormPath":             salvia_config.FormPaths[s.Lang]["FeminicideRiskPOST_GET"],
			"securityCityFormPath":       security_config.FormPaths[s.Lang]["CityGET"],
			"salviaEntityBranchFormPath": salvia_config.FormPaths[s.Lang]["EntityBranchGET"],
			"securityTownFormPath":       security_config.FormPaths[s.Lang]["TownGET"],
			"departments":                departments,
			"nav_rules":                  salvia_config.TranslateNavigationRule(s.Lang, salvia_config.NAVIGATION_RULES["set_feminicide_risk"]),
		}, utils.GetFullHtmlFuncMap())
}

// FeminicideRiskPUT maneja la solicitud PUT para actualizar o aprobar un caso de víctima.
// Según el parámetro "by" de la URL, decide si se aprueba el caso o se actualiza la información.
/*
func FeminicideRiskPUT(c *gin.Context) {
	var res string
	var code int
	// Se obtiene la sesión y se extrae el ID de sesión.
	session := sessions.Default(c)
	var sessionID string = session.Get("userData").(string)
	s, _ := utils.GetCommonSession(sessionID)

	// Se obtienen los parámetros "id" y "by" de la URL.
	id := c.Param("id")
	by := c.Param("by")

	// Si "by" indica "approve_case" tras la traducción, se procede a aprobar el caso.
	if by != "" && salvia_config.TranslateLocale(by, s.Lang) == "approve_case" {
		// Verifica el permiso "approve_feminicideRisk".
		if !utils.CheckPermission(salvia_config.PermissionsByRole, "approve_feminicideRisk", s.CurrentRole, c) {
			return
		}
		code, res = salvia_ctrl.ApproveFeminicideRisk(id, &db.ConnData{}, dbClientConfig, dbServerConfig)

	} else if by != "" && salvia_config.TranslateLocale(by, s.Lang) == "update" {
		// Si "by" indica "update", se verifica el permiso "update_feminicideRisk".
		if !utils.CheckPermission(salvia_config.PermissionsByRole, "update_feminicideRisk", s.CurrentRole, c) {
			return
		}
		// Se lee el cuerpo de la solicitud y se procesa la actualización.
		buf := new(bytes.Buffer)
		buf.ReadFrom(c.Request.Body)
		//code, res := salvia_ctrl.CtrlSetFeminicideRisk(buf.String())
		code, res = salvia_ctrl.UpdateFeminicideRisk(buf.String(), id, s, dbClientConfig, dbServerConfig)

	} else {
		// Si no se reconoce la acción, se responde con un error de solicitud incorrecta.
		code = http.StatusBadRequest
	}

	// Envía la respuesta al cliente en formato JSON.
	c.DataFromReader(code, int64(len(res)), gin.MIMEJSON, strings.NewReader(res), nil)
}
*/
