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

// VictimCasePOST maneja la solicitud POST para crear o actualizar un caso de víctima.
// Obtiene la sesión actual, verifica los permisos necesarios, lee el cuerpo de la solicitud,
// y llama al controlador correspondiente para procesar el caso de víctima.
func VictimCasePOST(c *gin.Context) {
	// Se obtiene la sesión del usuario desde el contexto de Gin.
	session := sessions.Default(c)
	// Se extrae el ID de sesión almacenado en "userData".
	var sessionID string = session.Get("userData").(string)
	// Se obtiene la sesión común asociada al ID.
	s, _ := utils.GetCommonSession(sessionID)

	// Verifica que el usuario tenga el permiso "set_victim_case".
	if !utils.CheckPermission(salvia_config.PermissionsByRole, "set_victim_case", s.CurrentRole, c) {
		return
	}

	// Lee el cuerpo de la solicitud en un buffer.
	buf := new(bytes.Buffer)
	buf.ReadFrom(c.Request.Body)
	// Llama al controlador para establecer el caso de víctima, pasando el contenido del buffer,
	// el origen "salvia", la sesión y la configuración de la base de datos.
	code, res := salvia_ctrl.SetVictimCase(buf.String(), *s, dbClientConfig, dbServerConfig)
	// Envía la respuesta al cliente en formato JSON.
	c.DataFromReader(code, int64(len(res)), gin.MIMEJSON, strings.NewReader(res), nil)
}

// VictimCaseGET maneja la solicitud GET para obtener casos de víctima.
// Dependiendo de los parámetros de la URL y el rol del usuario, recupera información
// detallada de un caso específico o una lista de casos, y renderiza la respuesta en JSON o HTML.
func VictimCaseGET(c *gin.Context) {
	// Se obtiene la sesión del usuario y el ID de sesión.
	session := sessions.Default(c)
	var sessionID string = session.Get("userData").(string)
	s, err := utils.GetCommonSession(sessionID)

	// Variables para almacenar las respuestas y otros datos relevantes.
	var victimContactRes string = "[]"
	var victimCaseRes string = "[]"
	var code int = http.StatusOK
	var tplName string
	// Objeto que contendrá la información detallada del caso de víctima.
	var vCase salvia_daos.VictimCaseDTO = salvia_daos.VictimCaseDTO{}
	// Mapa para agrupar ramas de entidades según el momento y sector.
	var entitiesByMomentAndSector map[string]map[string]map[string][]salvia_daos.EntityBranchDTO = make(map[string]map[string]map[string][]salvia_daos.EntityBranchDTO)
	// Listas para almacenar entidades y ramas.
	var entities []salvia_daos.EntityDTO = []salvia_daos.EntityDTO{}
	var branches []salvia_daos.EntityBranchDTO = []salvia_daos.EntityBranchDTO{}
	var operators string = "[]"
	var sectorBarriers []salvia_daos.SectorBarrierDTO = []salvia_daos.SectorBarrierDTO{}

	var count int

	// Mapa para almacenar momentos filtrados por sector.
	var momentsWithSector map[string]map[uint64]salvia_daos.MomentDTO = make(map[string]map[uint64]salvia_daos.MomentDTO)

	// Establece cabeceras para evitar cache en la respuesta.
	common_facades.SetHeaderNoCache(c)

	// Se obtiene el parámetro "id" de la URL.
	id := c.Param("id")
	filter := c.Param("f")
	page, _ := strconv.Atoi(c.Param("p"))

	//Si se consulta un caso particular
	if id != "" {
		// Si se proporciona el parámetro "docType", se realiza la búsqueda por documento.
		docType := c.Param("docType")
		if docType != "" {
			//Revisamos si es la opción no aplica ("na") caso en el cual se busca sólo por el número de doc.
			if docType == "na" {
				docType = ""
			}
			// Verifica el permiso "get_victim_case_by_document" para el rol actual.
			if !utils.CheckPermission(salvia_config.PermissionsByRole, "get_victim_case_by_document", s.CurrentRole, c) {
				return
			}
			// Si el rol es "et", se utiliza una función específica que incluye información de atención.
			if s.CurrentRole == "et" {
				code, victimCaseRes, count = salvia_ctrl.GetVictimCasesByDocumentAndTownCodeWithAttend(docType, id, s.TownCode, filter, page, &db.ConnData{}, dbClientConfig, dbServerConfig)
			} else {
				// En otro caso, se obtiene el caso por documento.
				code, victimCaseRes, count = salvia_ctrl.GetVictimCasesByDocument(docType, id, page, &db.ConnData{}, dbClientConfig, dbServerConfig)
			}
		} else {
			// Si no se especifica "docType", se verifica el permiso "get_victim_case".
			if !utils.CheckPermission(salvia_config.PermissionsByRole, "get_victim_case", s.CurrentRole, c) {
				return
			}
			// Se obtiene el caso de víctima por su código interno.
			code, victimCaseRes, vCase = salvia_ctrl.GetVictimCaseByICode(id, &db.ConnData{}, dbClientConfig, dbServerConfig)
		}

		if code == http.StatusOK {
			// Se filtran los momentos asociados al caso y se agrupan por código de momento y entidad.
			for _, m := range vCase.VictimCaseMoments {
				if momentsWithSector[m.MomentCode] == nil {
					momentsWithSector[m.MomentCode] = make(map[uint64]salvia_daos.MomentDTO)
				}
				momentsWithSector[m.MomentCode][m.MomentEntityBranch.EntityBranchId] = m
			}
			// Dependiendo del rol, se obtienen las ramas de las entidades de acuerdo al código del pueblo.
			switch s.CurrentRole {
			case "op":
				//Por defecto la última versión de plantilla
				tplName = "get_victim_case_op"
				if vCase.VictimCaseForm1.VictimCaseForm1ICode != "" {
					//Si se diligenció el formulario 1 vamos a su plantilla
					tplName = "get_victim_case_op_v1"
				}
				_, _, branches = salvia_ctrl.GetEntityBranchesByTownCodeWithMoments(vCase.VictimCaseTownCode, &db.ConnData{}, dbClientConfig, dbServerConfig)
			case "us":
				//Por defecto la última versión de plantilla
				tplName = "get_victim_case_us"
				if vCase.VictimCaseForm1.VictimCaseForm1ICode != "" {
					//Si se diligenció el formulario 1 vamos a su plantilla
					tplName = "get_victim_case_us_v1"
				}
				_, _, branches = salvia_ctrl.GetEntityBranchesByTownCodeWithMoments(vCase.VictimCaseTownCode, &db.ConnData{}, dbClientConfig, dbServerConfig)
			case "do":
				tplName = "get_victim_case_do"
				_, _, branches = salvia_ctrl.GetEntityBranchesByTownCodeWithMoments(vCase.VictimCaseTownCode, &db.ConnData{}, dbClientConfig, dbServerConfig)
			case "sv":
				tplName = "get_victim_case_sv"
				if vCase.VictimCaseForm1.VictimCaseForm1ICode != "" {
					//Si se diligenció el formulario 1 vamos a su plantilla
					tplName = "get_victim_case_sv_v1"
				}
				_, _, branches = salvia_ctrl.GetEntityBranchesByTownCodeWithMoments(vCase.VictimCaseTownCode, &db.ConnData{}, dbClientConfig, dbServerConfig)
			case "no":
				tplName = "get_victim_case_no"
				if vCase.VictimCaseForm1.VictimCaseForm1ICode != "" {
					//Si se diligenció el formulario 1 vamos a su plantilla
					tplName = "get_victim_case_no_v1"
				}
				_, _, branches = salvia_ctrl.GetEntityBranchesByTownCodeWithMoments(vCase.VictimCaseTownCode, &db.ConnData{}, dbClientConfig, dbServerConfig)
			case "ro":
				//Por defecto la última versión de plantilla
				tplName = "get_victim_case_ro"
				if vCase.VictimCaseForm1.VictimCaseForm1ICode != "" {
					//Si se diligenció el formulario 1 vamos a su plantilla
					tplName = "get_victim_case_ro_v1"
				}

				_, _, branches = salvia_ctrl.GetEntityBranchesByTownCodeWithMoments(vCase.VictimCaseTownCode, &db.ConnData{}, dbClientConfig, dbServerConfig)

				_, _, sectorBarriers = salvia_ctrl.GetAllSectorBarriers(&db.ConnData{}, dbClientConfig, dbServerConfig)
			case "et":
				tplName = "get_victim_case_et"
				if vCase.VictimCaseForm1.VictimCaseForm1ICode != "" {
					//Si se diligenció el formulario 1 vamos a su plantilla
					tplName = "get_victim_case_et_v1"
				}
				_, _, branches = salvia_ctrl.GetEntityBranchesByTownCodeAndEntityBranchIcodeWithMoments(vCase.VictimCaseTownCode, s.EntityBrandICode, &db.ConnData{}, dbClientConfig, dbServerConfig)
			case "fo":
				tplName = "get_victim_case_fo"
				if vCase.VictimCaseForm1.VictimCaseForm1ICode != "" {
					//Si se diligenció el formulario 1 vamos a su plantilla
					tplName = "get_victim_case_fo_v1"
				}
				_, _, branches = salvia_ctrl.GetEntityBranchesByTownCodeWithMoments(vCase.VictimCaseTownCode, &db.ConnData{}, dbClientConfig, dbServerConfig)
			}

			// Se agrupan las ramas de entidad por momento y sector.
			for _, e := range branches {
				if _, found := momentsWithSector[e.EntityBranchMoment][e.EntityBranchId]; found {
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
		}

		// Se obtienen todas las entidades registradas.
		_, _, entities = salvia_ctrl.GetEntityByAll(&db.ConnData{}, dbClientConfig, dbServerConfig)

	} else {
		//Si se consultan todos los casos
		// Si no se proporciona un ID, se verifica el permiso "get_victim_cases" para obtener todos los casos.
		if !utils.CheckPermission(salvia_config.PermissionsByRole, "get_victim_cases", s.CurrentRole, c) {
			return
		}
		// Se ejecutan distintas consultas según el rol del usuario.
		switch s.CurrentRole {
		case "op":
			tplName = "get_victim_cases"

			if filter == "" || filter == "fcv" {
				//Es la primera vez y carga por defecto los recontactos, o explícitamente se solicita el primer contacto.
				code, victimContactRes, count = salvia_ctrl.GetVictimContactsWithoutVictimCase(page, "v", &db.ConnData{}, dbClientConfig, dbServerConfig)
			} else if filter == "fci" {
				//Es la primera vez y carga por defecto los recontactos, o explícitamente se solicita el primer contacto.
				code, victimContactRes, count = salvia_ctrl.GetVictimContactsWithoutVictimCase(page, "i", &db.ConnData{}, dbClientConfig, dbServerConfig)
			} else {
				code, victimCaseRes, count = salvia_ctrl.GetVictimCasesByOwnerUserICode(s.UserICode, filter, page, &db.ConnData{}, dbClientConfig, dbServerConfig)
			}
		case "ro":
			tplName = "get_victim_cases_ro"

			if filter == "" || filter == "fcv" {
				//Es la primera vez y carga por defecto los recontactos, o explícitamente se solicita el primer contacto.
				code, victimContactRes, count = salvia_ctrl.GetVictimContactsWithoutVictimCase(page, "v", &db.ConnData{}, dbClientConfig, dbServerConfig)
			} else if filter == "fci" {
				//Es la primera vez y carga por defecto los recontactos, o explícitamente se solicita el primer contacto.
				code, victimContactRes, count = salvia_ctrl.GetVictimContactsWithoutVictimCase(page, "i", &db.ConnData{}, dbClientConfig, dbServerConfig)
			} else {
				code, victimCaseRes, count = salvia_ctrl.GetVictimCaseByAll(filter, page, &db.ConnData{}, dbClientConfig, dbServerConfig)
			}
		case "do":
			tplName = "get_victim_cases_do"
			if filter == "" {
				//traemos todos menos el filtro especificado. Por lo tanto "d" es cuando está realizado. Entonces traemos todos menos los realizados
				filter = "d"
			}
			code, victimCaseRes, count = salvia_ctrl.GetVictimCasesByTownCodeAndFollowUpStatusExcluded(s.TownCode, filter, page, &db.ConnData{}, dbClientConfig, dbServerConfig)

		case "et":
			if filter == "" { //Si es primera entrada por defecto se cargan los enrutados por aprobar
				filter = "r"
			}
			tplName = "get_victim_cases_et"
			code, victimCaseRes, count = salvia_ctrl.GetVictimCasesByTownCodeAndEntityBranchWithAttend(s.TownCode, s.EntityBrandICode, filter, page, &db.ConnData{}, dbClientConfig, dbServerConfig)
		case "us":
			tplName = "get_victim_cases_us"
			code, victimCaseRes, count = salvia_ctrl.GetVictimCasesByVictimUser(s.UserICode, page, &db.ConnData{}, dbClientConfig, dbServerConfig)
		case "sv":
			tplName = "get_victim_cases_sv"
			// Se obtienen los operadores generales con rol "op".
			_, operators = security_ctrl.GetGeneralUsersByRole("op", &db.ConnData{}, dbClientConfig, dbServerConfig)
			switch filter {
			case "", "fcv":
				//Es la primera vez y carga por defecto los recontactos, o explícitamente se solicita el primer contacto.
				code, victimContactRes, count = salvia_ctrl.GetVictimContactsWithoutVictimCase(page, "v", &db.ConnData{}, dbClientConfig, dbServerConfig)
			case "fci":
				//Es la primera vez y carga por defecto los recontactos, o explícitamente se solicita el primer contacto.
				code, victimContactRes, count = salvia_ctrl.GetVictimContactsWithoutVictimCase(page, "i", &db.ConnData{}, dbClientConfig, dbServerConfig)
			default:
				code, victimCaseRes, count = salvia_ctrl.GetVictimCaseByAll(filter, page, &db.ConnData{}, dbClientConfig, dbServerConfig)
			}

		case "no":
			tplName = "get_victim_cases_no"
			if filter == "" { //Si es primera entrada por defecto se cargan los entrutados aprobados
				filter = "ra"
			}
			code, victimCaseRes, count = salvia_ctrl.GetVictimCaseByAll(filter, page, &db.ConnData{}, dbClientConfig, dbServerConfig)

		case "fo":
			tplName = "get_victim_cases_fo"

			if filter == "" { //Si es primera entrada por defecto se cargan los entrutados aprobados
				filter = "ra"
			}
			code, victimCaseRes, count = salvia_ctrl.GetVictimCaseByAll(filter, page, &db.ConnData{}, dbClientConfig, dbServerConfig)
		}

	}

	// Se decide el formato de respuesta según el header "Accept" de la solicitud.
	if v, found := c.Request.Header["Accept"]; found && v[0] == "application/json" {
		// Respuesta en formato JSON.
		c.DataFromReader(code, int64(len(victimCaseRes)), gin.MIMEJSON, strings.NewReader(victimCaseRes), nil)
	} else {
		// Se prepara la información necesaria para renderizar la plantilla HTML.
		var menu map[string][]map[string]string
		var contactMenuTools []map[string]string = []map[string]string{}
		var caseMenuTools []map[string]string = []map[string]string{}

		if err == nil {
			menu = s.CurrentMenu
		}
		if v, found := salvia_config.MenuTools[s.Lang][s.CurrentRole]["menu_tool_get_victim_contacts"]; found {
			contactMenuTools = v
		}

		if v, found := salvia_config.MenuTools[s.Lang][s.CurrentRole]["menu_tool_get_victim_cases"]; found {
			caseMenuTools = v
		}

		// Renderiza la plantilla HTML con todos los parámetros necesarios para mostrar el caso de víctima.
		common_facades.RenderTemplate(c, salvia_daos.VictimContactEntityName, "salvia", "victim_case/", salvia_config.HTML_Templates, tplName, utils.GetFullHtmlTemplates(), utils.DEFAULT_VIEW, utils.DEFAULT_PANIC_TEMPLATE,
			map[string]interface{}{
				"windowTitle":               salvia_config.Locale["sp"]["get_victim_case_window_title"],
				"currentUser":               s.Names + " " + s.LastNames,
				"nav_rules":                 salvia_config.TranslateNavigationRule(s.Lang, salvia_config.NAVIGATION_RULES["get_victim_case"]),
				"locale":                    salvia_config.Locale,
				"lang":                      s.Lang,
				"victimContacts":            victimContactRes,
				"victimCases":               victimCaseRes,
				"numPages":                  int(math.Ceil(float64(count) / float64(salvia_config.NUM_ITEMS_PER_PAGE))),
				"gender":                    common_config.GENDER_IDENTITY,
				"docTypeUnified":            mergeDocumentTypes(common_config.DOCUMENT_TYPE_FORM2, common_config.DOCUMENT_TYPE),
				"docType":                   common_config.DOCUMENT_TYPE,
				"docType2":                  common_config.DOCUMENT_TYPE_FORM2,
				"language":                  common_config.LANGUAGE,
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
				"yes_no_enum":               common_config.YES_NO,
				"yes_no_na":                 common_config.YES_NO_NA,
				"riskLevel":                 salvia_config.RISK_LEVEL[s.Lang],
				"actionStatus":              salvia_config.FOLLOW_UP_ENTRY_ACTING_STATUS[s.Lang],

				"victimCaseVictimNationality":                         salvia_config.VICTIM_CASE_VICTIM_NATIONALITY,
				"victimCaseVictimForeignerImmigrationStatus":          salvia_config.VICTIM_CASE_VICTIM_FOREIGNER_IMMIGRATION_STATUS,
				"victimCaseVictimGender":                              salvia_config.VICTIM_CASE_VICTIM_GENDER,
				"victimCaseVictimDependents":                          salvia_config.VICTIM_CASE_VICTIM_DEPENDENTS,
				"victimCaseVictimDeathThreats":                        common_config.YES_NO,
				"victimCaseVictimAggressorHasWeapons":                 common_config.YES_NO,
				"victimCaseExperiencedPhysicalOrSexualViolenceBefore": common_config.YES_NO,
				"victimCaseVictimImminentRisk":                        common_config.YES_NO,
				"victimCaseVictimPreviouslyReportedSituation":         common_config.YES_NO,
				"victimCaseVictimIfPreviouslyReported":                salvia_config.VICTIM_CASE_VICTIM_IF_PREVIOUSLY_REPORTED,
				"followUpStatus":                                      salvia_config.FOLLOW_UP_STATUS[s.Lang],
				"followUpEntryStatus":                                 salvia_config.FOLLOW_UP_ENTRY_STATUS[s.Lang],

				"ifAfro":             salvia_config.AFRO_COMMUNTITIES[s.Lang],
				"ifIndigenous":       salvia_config.COLOMBIAN_INDIGENOUS[s.Lang],
				"ifIndigenousTongue": salvia_config.INDIGENOUS_TONGUES[s.Lang],
				"ifPeasant":          common_config.YES_NO,
				"ifArmedConflict":    common_config.YES_NO,
				"violenceScene":      salvia_config.VIOLENCE_SCENES[s.Lang],

				"menu":                             menu,
				"menuToolsContactTable":            contactMenuTools,
				"menuToolsCaseTable":               caseMenuTools,
				"momentCodes":                      salvia_config.MOMENT[s.Lang],
				"sectors":                          salvia_config.SECTOR[s.Lang],
				"entities":                         entities,
				"operators":                        operators,
				"sectorBarriers":                   sectorBarriers,
				"entitiesByMomentAndSector":        entitiesByMomentAndSector,
				"salviaMomentFormPath":             salvia_config.FormPaths[s.Lang]["MomentPUT"],
				"salviaAlertFormPath":              salvia_config.FormPaths[s.Lang]["AlertGET"],
				"salviaAssignOperatorFormPath":     salvia_config.FormPaths[s.Lang]["AssignOperatorsGET"],
				"salviaReportFormPath":             salvia_config.FormPaths[s.Lang]["VictimCaseReportGET"],
				"salviaVictimCaseFormPath":         salvia_config.FormPaths[s.Lang]["VictimCaseReportPOST"],
				"salviaVictimCaseSearchFormPath":   salvia_config.FormPaths[s.Lang]["VictimCaseGET"],
				"salviaVictimCaseGETFormPath":      salvia_config.FormPaths[s.Lang]["VictimCaseGET"],
				"salviaVictimContactGETFormPath":   salvia_config.FormPaths[s.Lang]["VictimContactGET"],
				"salviaVictimCaseUpdateFormPath":   salvia_config.FormPaths[s.Lang]["VictimCasePUT"],
				"salviaAssignFormPath":             salvia_config.FormPaths[s.Lang]["AssignOperatorsPOST"],
				"salviaCaseLogFormPath":            salvia_config.FormPaths[s.Lang]["CaseLogPOST"],
				"salviaFollowUpPOSTFormPath":       salvia_config.FormPaths[s.Lang]["FollowUpPOST"],
				"salviaFollowUpEntryPUTFormPath":   salvia_config.FormPaths[s.Lang]["FollowUpEntryPUT"],
				"salviaBarrierGETFormPath":         salvia_config.FormPaths[s.Lang]["BarrierGET"],
				"salviaEntityBranchesFormPath":     salvia_config.FormPaths[s.Lang]["EntityBranchGET"],
				"salviaFollowUpActingPOSTFormPath": salvia_config.FormPaths[s.Lang]["FollowUpEntryActingPOST"],
				"salviaFeminicideRiskSetFormPath":  salvia_config.FormPaths[s.Lang]["FeminicideRiskGET"],

				"salviaUserFormPath": security_config.FormPaths[s.Lang]["GeneralUserGET"],
			}, utils.GetFullHtmlFuncMap())
	}
}

// VictimCaseReportGET maneja la solicitud GET para la generación de reportes de casos de víctimas.
// Recupera información necesaria, como los departamentos, y renderiza la plantilla correspondiente.
func VictimCaseReportGET(c *gin.Context) {
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
	common_facades.RenderTemplate(c, salvia_daos.VictimContactEntityName, "salvia", "victim_case/", salvia_config.HTML_Templates, "report_victim_cases", utils.GetFullHtmlTemplates(), utils.REPORT_VIEW, utils.DEFAULT_PANIC_TEMPLATE,
		map[string]interface{}{
			"windowTitle":          salvia_config.Locale["sp"]["report_victim_case_window_title"],
			"currentUser":          s.Names + " " + s.LastNames,
			"nav_rules":            salvia_config.TranslateNavigationRule(s.Lang, salvia_config.NAVIGATION_RULES["report_victim_cases"]),
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
			"salviaReportFormPath": salvia_config.FormPaths[s.Lang]["VictimCaseReportGET"],
		}, utils.GetFullHtmlFuncMap())
}

// VictimCaseReportPOST maneja la solicitud POST para generar reportes de casos de víctimas.
// Verifica permisos, procesa los datos recibidos, y renderiza la respuesta ya sea en JSON o mediante plantilla HTML.
func VictimCaseReportPOST(c *gin.Context) {
	// Se obtiene la sesión y se extrae el ID de sesión.
	session := sessions.Default(c)
	var sessionID string = session.Get("userData").(string)
	s, err := utils.GetCommonSession(sessionID)
	// Establece cabeceras para evitar cache.
	common_facades.SetHeaderNoCache(c)

	var victimCasesRes string
	var code int = http.StatusOK
	var tplName string = "report_victim_cases"
	// Verifica que el usuario tenga el permiso "report_victim_cases".
	if !utils.CheckPermission(salvia_config.PermissionsByRole, "report_victim_cases", s.CurrentRole, c) {
		return
	}
	// Lee el cuerpo de la solicitud.
	buf := new(bytes.Buffer)
	buf.ReadFrom(c.Request.Body)
	// Llama al controlador para generar el reporte de casos de víctimas.
	code, victimCasesRes = salvia_ctrl.ReportVictimCases(buf.String(), &db.ConnData{}, dbClientConfig, dbServerConfig)

	// Según el header "Accept", se devuelve la respuesta en JSON o se renderiza la plantilla.
	if v, found := c.Request.Header["Accept"]; found && v[0] == "application/json" {
		c.DataFromReader(code, int64(len(victimCasesRes)), gin.MIMEJSON, strings.NewReader(victimCasesRes), nil)
	} else {
		var menu map[string][]map[string]string

		if err == nil {
			menu = s.CurrentMenu
		}
		// Se obtienen todas las entidades.
		_, _, entities := salvia_ctrl.GetEntityByAll(&db.ConnData{}, dbClientConfig, dbServerConfig)
		// Renderiza la plantilla HTML para mostrar el reporte de casos de víctimas.
		common_facades.RenderTemplate(c, salvia_daos.VictimContactEntityName, "salvia", "victim_case/", salvia_config.HTML_Templates, tplName, utils.GetFullHtmlTemplates(), utils.REPORT_VIEW, utils.DEFAULT_PANIC_TEMPLATE,
			map[string]interface{}{
				"windowTitle":          salvia_config.Locale["sp"]["report_victim_case_window_title"],
				"currentUser":          s.Names + " " + s.LastNames,
				"nav_rules":            salvia_config.TranslateNavigationRule(s.Lang, salvia_config.NAVIGATION_RULES["report_victim_cases"]),
				"locale":               salvia_config.Locale,
				"lang":                 s.Lang,
				"victimCases":          victimCasesRes,
				"femicideRisk":         common_config.YES_NO,
				"violenceExperienced":  salvia_config.VIOLENCE_EXPERIENCED,
				"menu":                 menu,
				"entities":             entities,
				"salviaMomentFormPath": salvia_config.FormPaths[s.Lang]["MomentPUT"],
				"salviaAlertFormPath":  salvia_config.FormPaths[s.Lang]["AlertGET"],
				"salviaReportFormPath": salvia_config.FormPaths[s.Lang]["VictimCaseReportGET"],
			}, utils.GetFullHtmlFuncMap())
	}
}

// VictimCasePUT_GET maneja la solicitud GET para obtener la información de un caso de víctima a actualizar.
// Recupera datos relacionados con el caso, contacto, ubicación y entidades, y renderiza la plantilla de actualización.
func VictimCasePUT_GET(c *gin.Context) {
	// Se obtiene la sesión y se extrae el ID de sesión.
	session := sessions.Default(c)
	var sessionID string = session.Get("userData").(string)
	s, _ := utils.GetCommonSession(sessionID)

	// Verifica que el usuario tenga el permiso "update_victim_case".
	if !utils.CheckPermission(salvia_config.PermissionsByRole, "update_victim_case", s.CurrentRole, c) {
		return
	}

	// Declaración de variables para almacenar información del caso y contacto.
	var victimContact salvia_daos.VictimContactDTO = salvia_daos.VictimContactDTO{}
	var victimCase salvia_daos.VictimCaseDTO = salvia_daos.VictimCaseDTO{}

	// Variables para almacenar datos geográficos.
	var departments string = "[]"
	var cities string = "[]"
	var towns string = "[]"

	var cities2 string = "[]"
	var towns2 string = "[]"

	var cities3 string = "[]"
	var towns3 string = "[]"

	var violenceCities string = "{}"
	var violenceTowns string = "{}"

	var tplName string

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
		_, _, victimCase = salvia_ctrl.GetVictimCaseByICode(id, &db.ConnData{}, dbClientConfig, dbServerConfig)
		_, _, victimContact = salvia_ctrl.GetVictimContactByICode(strconv.FormatUint(victimCase.VictimCaseVictimContact.VictimContactId, 10), &db.ConnData{}, dbClientConfig, dbServerConfig)

		tplName = "update_victim_case_v2"
		if victimCase.VictimCaseForm1.VictimCaseForm1ICode != "" {
			//Si se diligenció el formulario 1 vamos a su plantilla
			tplName = "update_victim_case_v1"
			// Para el lugar de los hechos, se obtienen las ciudades y pueblos basados en el departamento y ciudad seleccionados.
			if victimCase.VictimCaseForm1.VictimCaseForm1ViolenceDepartment.DepartmentId > 0 {
				_, cities2, _ = security_ctrl.GetCitiesByDeparment(victimCase.VictimCaseForm1.VictimCaseForm1ViolenceDepartment.DepartmentId, &db.ConnData{}, dbClientConfig, dbServerConfig)
			}
			if victimCase.VictimCaseForm1.VictimCaseForm1ViolenceCity.CityId > 0 {
				_, towns2 = security_ctrl.GetTownsByCity(victimCase.VictimCaseForm1.VictimCaseForm1ViolenceCity.CityId, &db.ConnData{}, dbClientConfig, dbServerConfig)
			}

		} else {
			// Para el lugar de los hechos, se obtienen las ciudades y pueblos basados en el departamento y ciudad seleccionados.
			if victimCase.VictimCaseForm2.VictimCaseForm2FactsDepartment.DepartmentId > 0 {
				_, cities2, _ = security_ctrl.GetCitiesByDeparment(victimCase.VictimCaseForm2.VictimCaseForm2FactsDepartment.DepartmentId, &db.ConnData{}, dbClientConfig, dbServerConfig)
			}
			if victimCase.VictimCaseForm2.VictimCaseForm2FactsCity.CityId > 0 {
				_, towns2 = security_ctrl.GetTownsByCity(victimCase.VictimCaseForm2.VictimCaseForm2FactsCity.CityId, &db.ConnData{}, dbClientConfig, dbServerConfig)
			}

			if victimCase.VictimCaseForm2.VictimCaseForm2ResidenceDepartment.DepartmentId > 0 {
				_, cities3, _ = security_ctrl.GetCitiesByDeparment(victimCase.VictimCaseForm2.VictimCaseForm2ResidenceDepartment.DepartmentId, &db.ConnData{}, dbClientConfig, dbServerConfig)
			}
			if victimCase.VictimCaseForm2.VictimCaseForm2ResidenceCity.CityId > 0 {
				_, towns3 = security_ctrl.GetTownsByCity(victimCase.VictimCaseForm2.VictimCaseForm2ResidenceCity.CityId, &db.ConnData{}, dbClientConfig, dbServerConfig)
			}
		}
	}

	// Carga datos adicionales en victimCase, excluyendo información redundante del contacto.
	victimCase.LoadFromWithoutVictimContact(salvia_config.MOMENT[s.Lang], salvia_config.SECTOR[s.Lang], entities)

	// Para el lugar de atención, se obtienen las ciudades y pueblos correspondientes.
	if victimCase.VictimCaseDepartment.DepartmentId > 0 {
		_, cities, _ = security_ctrl.GetCitiesByDeparment(victimCase.VictimCaseDepartment.DepartmentId, &db.ConnData{}, dbClientConfig, dbServerConfig)
	}
	if victimCase.VictimCaseTown.TownId > 0 {
		_, towns = security_ctrl.GetTownsByCity(victimCase.VictimCaseCity.CityId, &db.ConnData{}, dbClientConfig, dbServerConfig)

		// Se obtienen las ramas de entidad asociadas al pueblo para asignar la ruta.
		var branches []salvia_daos.EntityBranchDTO
		_, _, branches = salvia_ctrl.GetEntityBranchesByTownCodeWithMoments(victimCase.VictimCaseTown.TownCode, &db.ConnData{}, dbClientConfig, dbServerConfig)
		/*
			for _, m := range victimCase.VictimCaseMoments {
				if momentsWithSector[m.MomentCode] == nil {
					momentsWithSector[m.MomentCode] = make(map[uint64]salvia_daos.MomentDTO)
				}
				momentsWithSector[m.MomentCode][m.MomentEntityBranch.EntityBranchId] = m
			}*/
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
	common_facades.RenderTemplate(c, salvia_daos.VictimCaseEntityName, "salvia", "victim_case/", salvia_config.HTML_Templates, tplName, utils.GetFullHtmlTemplates(), utils.DEFAULT_VIEW, utils.DEFAULT_PANIC_TEMPLATE,
		map[string]interface{}{
			"windowTitle":   salvia_config.Locale["sp"]["update_victim_case_window_title"],
			"currentUser":   s.Names + " " + s.LastNames,
			"victimCase":    victimCase,
			"victimContact": victimContact,
			"gender":        common_config.GENDER_IDENTITY,
			"language":      common_config.LANGUAGE,
			"locale":        salvia_config.Locale,
			"lang":          s.Lang,
			//Soporte para el formulario v1
			"livingZone":                  salvia_config.LIVING_ZONES,
			"disability":                  salvia_config.DISABILITY,
			"occurrence":                  salvia_config.OCCURRENCE,
			"violenceScope":               salvia_config.VIOLENCE_SCOPE,
			"aggressor":                   salvia_config.AGGRESSOR,
			"femicideRisk":                common_config.YES_NO,
			"violenceExperienced":         salvia_config.VIOLENCE_EXPERIENCED,
			"relationshipWithAggressor":   salvia_config.RELATIONSHIP_WITH_AGGRESSOR,
			"weekDay":                     salvia_config.WEEK_DAY,
			"docType2":                    common_config.DOCUMENT_TYPE_FORM2,
			"victimCaseVictimNationality": salvia_config.VICTIM_CASE_VICTIM_NATIONALITY,
			"victimCaseVictimForeignerImmigrationStatus":          salvia_config.VICTIM_CASE_VICTIM_FOREIGNER_IMMIGRATION_STATUS,
			"victimCaseVictimGender":                              salvia_config.VICTIM_CASE_VICTIM_GENDER,
			"victimCaseVictimDependents":                          salvia_config.VICTIM_CASE_VICTIM_DEPENDENTS,
			"victimCaseVictimDeathThreats":                        common_config.YES_NO,
			"victimCaseVictimAggressorHasWeapons":                 common_config.YES_NO,
			"victimCaseExperiencedPhysicalOrSexualViolenceBefore": common_config.YES_NO,
			"victimCaseVictimImminentRisk":                        common_config.YES_NO,
			"victimCaseVictimPreviouslyReportedSituation":         common_config.YES_NO,
			"victimCaseVictimIfPreviouslyReported":                salvia_config.VICTIM_CASE_VICTIM_IF_PREVIOUSLY_REPORTED,
			"ifAfro":                                              salvia_config.AFRO_COMMUNTITIES[s.Lang],
			"ifIndigenous":                                        salvia_config.COLOMBIAN_INDIGENOUS[s.Lang],
			"ifIndigenousTongue":                                  salvia_config.INDIGENOUS_TONGUES[s.Lang],
			"ifPeasant":                                           common_config.YES_NO,
			"ifArmedConflict":                                     common_config.YES_NO,
			"violenceScene":                                       salvia_config.VIOLENCE_SCENES[s.Lang],
			"sexualOrientationForm1":                              salvia_config.SEXUAL_ORIENTATION,
			"ethnicGroup":                                         salvia_config.ETHNIC_GROUP,
			"maritalStatusForm1":                                  common_config.MARITAL_STATUS,
			"genderIdentityForm1":                                 common_config.GENDER_IDENTITY,
			"origin":                                              salvia_config.ORIGIN_PLACE,
			"occupationForm1":                                     salvia_config.OCCUPATION,

			"yes_no_form1": common_config.YES_NO,

			//-------------------------------

			"yes_no":                              salvia_daos.VictimCaseForm2Enums["yes_no"],
			"docType":                             salvia_daos.VictimCaseForm2Enums["victim_case_form2_victim_doc_type"],
			"factsZone":                           salvia_daos.VictimCaseForm2Enums["victim_case_form2_facts_zone"],
			"scenarioViolence":                    salvia_daos.VictimCaseForm2Enums["victim_case_form2_scenario_violence"],
			"recurrenceAggression":                salvia_daos.VictimCaseForm2Enums["victim_case_form2_recurrence_aggression"],
			"numAgressors":                        salvia_daos.VictimCaseForm2Enums["victim_case_form2_num_agressors"],
			"proximityPrincipalAggressor":         salvia_daos.VictimCaseForm2Enums["victim_case_form2_proximity_principal_aggressor"],
			"relationshipWithPresumedAggressor01": salvia_daos.VictimCaseForm2Enums["victim_case_form2_relationship_with_presumed_aggressor_01"],
			"aggressorOccupation":                 salvia_daos.VictimCaseForm2Enums["victim_case_form2_relationship_with_presumed_aggressor_02"],
			"aggressorGenderIdentity":             salvia_daos.VictimCaseForm2Enums["victim_case_form2_aggressor_gender_identity"],
			"nationality":                         salvia_daos.VictimCaseForm2Enums["victim_case_form2_nationality"],
			"specifiedNationality":                salvia_daos.VictimCaseForm2Enums["victim_case_form2_specified_nationality"],
			"migrationCondition":                  salvia_daos.VictimCaseForm2Enums["victim_case_form2_migration_condition"],
			"genderIdentity":                      salvia_daos.VictimCaseForm2Enums["victim_case_form2_gender_identity"],
			"sexualOrientation":                   salvia_daos.VictimCaseForm2Enums["victim_case_form2_sexual_orientation"],
			"assignedSexAtBirth":                  salvia_daos.VictimCaseForm2Enums["victim_case_form2_assigned_sex_at_birth"],
			"ethnicAffiliation":                   salvia_daos.VictimCaseForm2Enums["victim_case_form2_ethnic_affiliation"],
			"indigenousPeople":                    salvia_daos.VictimCaseForm2Enums["victim_case_form2_indigenous_people"],
			"maritalStatus":                       salvia_daos.VictimCaseForm2Enums["victim_case_form2_marital_status"],
			"lastEducationLevel":                  salvia_daos.VictimCaseForm2Enums["victim_case_form2_last_education_level"],
			"occupation":                          salvia_daos.VictimCaseForm2Enums["victim_case_form2_occupation"],
			"incomeGenerationMethod":              salvia_daos.VictimCaseForm2Enums["victim_case_form2_income_generation_method"],
			"employmentRelationship":              salvia_daos.VictimCaseForm2Enums["victim_case_form2_employment_relationship"],
			"housingTenancyForm":                  salvia_daos.VictimCaseForm2Enums["victim_case_form2_housing_tenancy_form"],
			"housingStratum":                      salvia_daos.VictimCaseForm2Enums["victim_case_form2_housing_stratum"],
			"currentlyPregnant":                   salvia_daos.VictimCaseForm2Enums["victim_case_form2_currently_pregnant"],
			"supportContactKinship":               salvia_daos.VictimCaseForm2Enums["victim_case_form2_support_contact_kinship"],
			//Campos múltiples
			"typeViolenceExperienced":       salvia_daos.VictimCaseForm2Enums["victim_case_form2_type_violence_experienced"],
			"subtypeViolenceExperienced_fi": salvia_daos.VictimCaseForm2Enums["victim_case_form2_subtype_violence_experienced_fi"],
			"subtypeViolenceExperienced_ps": salvia_daos.VictimCaseForm2Enums["victim_case_form2_subtype_violence_experienced_ps"],
			"subtypeViolenceExperienced_se": salvia_daos.VictimCaseForm2Enums["victim_case_form2_subtype_violence_experienced_se"],
			"subtypeViolenceExperienced_po": salvia_daos.VictimCaseForm2Enums["victim_case_form2_subtype_violence_experienced_po"],
			"subtypeViolenceExperienced_pl": salvia_daos.VictimCaseForm2Enums["victim_case_form2_subtype_violence_experienced_pl"],
			"subtypeViolenceExperienced_re": salvia_daos.VictimCaseForm2Enums["victim_case_form2_subtype_violence_experienced_re"],
			"subtypeViolenceExperienced_vi": salvia_daos.VictimCaseForm2Enums["victim_case_form2_subtype_violence_experienced_vi"],
			"scopeOfViolence":               salvia_daos.VictimCaseForm2Enums["victim_case_form2_scope_of_violence"],
			"whoReportTo":                   salvia_daos.VictimCaseForm2Enums["victim_case_form2_who_report_to"],
			"activitiesUnableToPerform":     salvia_daos.VictimCaseForm2Enums["victim_case_form2_activities_unable_to_perform"],
			"adjustmentsGBV":                salvia_daos.VictimCaseForm2Enums["victim_case_form2_adjustments_gbv"],
			"law1996":                       salvia_daos.VictimCaseForm2Enums["victim_case_form2_law_1996"],
			"speciallyProtectedPopulation":  salvia_daos.VictimCaseForm2Enums["victim_case_form2_specially_protected_population"],
			"aspMode":                       salvia_daos.VictimCaseForm2Enums["victim_case_form2_asp_mode"],
			"reasonASP":                     salvia_daos.VictimCaseForm2Enums["victim_case_form2_reason_asp"],
			"hasDependents":                 salvia_daos.VictimCaseForm2Enums["victim_case_form2_has_dependents"],
			"actionPlan":                    salvia_daos.VictimCaseForm2Enums["victim_case_form2_action_plan"],

			//----------

			"moments":                    salvia_config.MOMENT[s.Lang],
			"sectors":                    salvia_config.SECTOR[s.Lang],
			"entities":                   entities,
			"entitiesByMomentAndSector":  entitiesByMomentAndSector,
			"salviaFormPath":             salvia_config.FormPaths[s.Lang]["VictimCasePUT"],
			"securityCityFormPath":       security_config.FormPaths[s.Lang]["CityGET"],
			"salviaEntityBranchFormPath": salvia_config.FormPaths[s.Lang]["EntityBranchGET"],
			"securityTownFormPath":       security_config.FormPaths[s.Lang]["TownGET"],
			"departments":                departments,
			"cities":                     cities,
			"towns":                      towns,
			"cities2":                    cities2,
			"towns2":                     towns2,
			"cities3":                    cities3,
			"towns3":                     towns3,
			"violenceCities":             violenceCities,
			"violenceTowns":              violenceTowns,
			"nav_rules":                  salvia_config.TranslateNavigationRule(s.Lang, salvia_config.NAVIGATION_RULES["update_victim_case"]),
		}, utils.GetFullHtmlFuncMap())
}

// VictimCasePOST_GET maneja la solicitud GET para inicializar la creación de un caso de víctima.
// Obtiene la información de contacto de la víctima, carga datos por defecto para el caso,
// y renderiza la plantilla correspondiente para la creación del caso.
func VictimCasePOST_GET(c *gin.Context) {
	// Se obtiene la sesión y el ID de sesión.
	session := sessions.Default(c)
	var sessionID string = session.Get("userData").(string)
	s, _ := utils.GetCommonSession(sessionID)

	// Verifica que el usuario tenga el permiso "set_victim_case".
	if !utils.CheckPermission(salvia_config.PermissionsByRole, "set_victim_case", s.CurrentRole, c) {
		return
	}

	// Variables para almacenar el caso, contacto y datos geográficos.
	var victimCaseStr string = "{}"
	var victimContact salvia_daos.VictimContactDTO = salvia_daos.VictimContactDTO{}
	var departments string = "{}"
	var cities string = "{}"
	var towns string = "{}"
	var entities []salvia_daos.EntityDTO = []salvia_daos.EntityDTO{}

	// Mapa para almacenar las ramas de entidad agrupadas por momento y sector.
	var entitiesByMomentAndSector map[string]map[string]map[string][]salvia_daos.EntityBranchDTO = make(map[string]map[string]map[string][]salvia_daos.EntityBranchDTO)
	var entitiesByMomentAndSectorStr string = "null"

	// Caso vacío que se llenará con la información del contacto.
	var victimCaseEmpty salvia_daos.VictimCaseDTO = salvia_daos.VictimCaseDTO{}
	// Establece cabeceras para evitar cache.
	common_facades.SetHeaderNoCache(c)

	// Se obtienen los departamentos y las entidades.
	_, departments = security_ctrl.GetDepartmentByAll(&db.ConnData{}, dbClientConfig, dbServerConfig)
	_, _, entities = salvia_ctrl.GetEntityByAll(&db.ConnData{}, dbClientConfig, dbServerConfig)

	// Se obtiene el contacto de la víctima basado en el parámetro "id".
	id := c.Param("id")
	if id != "" {
		_, _, victimContact = salvia_ctrl.GetVictimContactByICode(id, &db.ConnData{}, dbClientConfig, dbServerConfig)
		if victimContact.VictimContactForm1.VictimContactForm1Id > 0 {
			victimCaseEmpty.LoadFromVictimContactForm1(victimContact, salvia_config.MOMENT[s.Lang], salvia_config.SECTOR[s.Lang], entities)
		} else if victimContact.VictimContactForm2.VictimContactForm2Id > 0 {
			victimCaseEmpty.LoadFromVictimContactForm2(victimContact, salvia_config.MOMENT[s.Lang], salvia_config.SECTOR[s.Lang], entities)
		}

	} else {
		//Por defecto asume el formulario 2
		victimCaseEmpty.LoadFromVictimContactForm2(victimContact, salvia_config.MOMENT[s.Lang], salvia_config.SECTOR[s.Lang], entities)
	}

	// Se convierte el caso a JSON de éxito.
	victimCaseStr = utils.CommMsgGetJSONSuccess(victimCaseEmpty)

	// Obtiene las ciudades y pueblos según el departamento y ciudad seleccionados.
	if victimCaseEmpty.VictimCaseDepartment.DepartmentId > 0 {
		_, cities, _ = security_ctrl.GetCitiesByDeparment(victimCaseEmpty.VictimCaseDepartment.DepartmentId, &db.ConnData{}, dbClientConfig, dbServerConfig)
	}
	if victimCaseEmpty.VictimCaseTown.TownId > 0 {
		_, towns = security_ctrl.GetTownsByCity(victimCaseEmpty.VictimCaseCity.CityId, &db.ConnData{}, dbClientConfig, dbServerConfig)

		// Se obtienen las ramas de entidad para el pueblo y se agrupan.
		var branches []salvia_daos.EntityBranchDTO
		_, _, branches = salvia_ctrl.GetEntityBranchesByTownCodeWithMoments(victimCaseEmpty.VictimCaseTown.TownCode, &db.ConnData{}, dbClientConfig, dbServerConfig)

		for _, e := range branches {
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

		// Convierte el mapa de ramas agrupadas a JSON.
		entitiesByMomentAndSectorStr = utils.CommMsgGetJSONSuccess(entitiesByMomentAndSector)
	}

	// Renderiza la plantilla HTML para la creación de un nuevo caso de víctima.
	common_facades.RenderTemplate(c, salvia_daos.VictimCaseEntityName, "salvia", "victim_case/", salvia_config.HTML_Templates, "set_victim_case", utils.GetFullHtmlTemplates(), utils.DEFAULT_VIEW, utils.DEFAULT_PANIC_TEMPLATE,
		map[string]interface{}{
			"windowTitle": salvia_config.Locale["sp"]["set_victim_case_window_title"],
			"currentUser": s.Names + " " + s.LastNames,
			"victimCase":  victimCaseStr,
			"gender":      common_config.GENDER_IDENTITY,
			"language":    common_config.LANGUAGE,
			"locale":      salvia_config.Locale,
			"lang":        s.Lang,
			"livingZone":  salvia_config.LIVING_ZONES,

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

			//New form ------
			"yes_no":                    salvia_daos.VictimCaseForm2Enums["yes_no"],
			"docType":                   salvia_daos.VictimCaseForm2Enums["victim_case_form2_victim_doc_type"],
			"docType2":                  common_config.DOCUMENT_TYPE_FORM2,
			"factsZone":                 salvia_daos.VictimCaseForm2Enums["victim_case_form2_facts_zone"],
			"scenarioViolence":          salvia_daos.VictimCaseForm2Enums["victim_case_form2_scenario_violence"],
			"workplaceSectorOccurrence": salvia_daos.VictimCaseForm2Enums["victim_case_form2_workplace_sector_occurrence"],

			"recurrenceAggression":                salvia_daos.VictimCaseForm2Enums["victim_case_form2_recurrence_aggression"],
			"numAgressors":                        salvia_daos.VictimCaseForm2Enums["victim_case_form2_num_agressors"],
			"proximityPrincipalAggressor":         salvia_daos.VictimCaseForm2Enums["victim_case_form2_proximity_principal_aggressor"],
			"relationshipWithPresumedAggressor01": salvia_daos.VictimCaseForm2Enums["victim_case_form2_relationship_with_presumed_aggressor_01"],
			"aggressorOccupation":                 salvia_daos.VictimCaseForm2Enums["victim_case_form2_relationship_with_presumed_aggressor_02"],
			"aggressorGenderIdentity":             salvia_daos.VictimCaseForm2Enums["victim_case_form2_aggressor_gender_identity"],
			"nationality":                         salvia_daos.VictimCaseForm2Enums["victim_case_form2_nationality"],
			"specifiedNationality":                salvia_daos.VictimCaseForm2Enums["victim_case_form2_specified_nationality"],
			"migrationCondition":                  salvia_daos.VictimCaseForm2Enums["victim_case_form2_migration_condition"],
			"genderIdentity":                      salvia_daos.VictimCaseForm2Enums["victim_case_form2_gender_identity"],
			"sexualOrientation":                   salvia_daos.VictimCaseForm2Enums["victim_case_form2_sexual_orientation"],
			"assignedSexAtBirth":                  salvia_daos.VictimCaseForm2Enums["victim_case_form2_assigned_sex_at_birth"],
			"ethnicAffiliation":                   salvia_daos.VictimCaseForm2Enums["victim_case_form2_ethnic_affiliation"],
			"indigenousPeople":                    salvia_daos.VictimCaseForm2Enums["victim_case_form2_indigenous_people"],
			"maritalStatus":                       salvia_daos.VictimCaseForm2Enums["victim_case_form2_marital_status"],
			"lastEducationLevel":                  salvia_daos.VictimCaseForm2Enums["victim_case_form2_last_education_level"],
			"occupation":                          salvia_daos.VictimCaseForm2Enums["victim_case_form2_occupation"],
			"incomeGenerationMethod":              salvia_daos.VictimCaseForm2Enums["victim_case_form2_income_generation_method"],
			"employmentRelationship":              salvia_daos.VictimCaseForm2Enums["victim_case_form2_employment_relationship"],
			"housingTenancyForm":                  salvia_daos.VictimCaseForm2Enums["victim_case_form2_housing_tenancy_form"],
			"housingStratum":                      salvia_daos.VictimCaseForm2Enums["victim_case_form2_housing_stratum"],
			"currentlyPregnant":                   salvia_daos.VictimCaseForm2Enums["victim_case_form2_currently_pregnant"],
			"supportContactKinship":               salvia_daos.VictimCaseForm2Enums["victim_case_form2_support_contact_kinship"],
			//Campos múltiples
			"typeViolenceExperienced":       salvia_daos.VictimCaseForm2Enums["victim_case_form2_type_violence_experienced"],
			"subtypeViolenceExperienced_fi": salvia_daos.VictimCaseForm2Enums["victim_case_form2_subtype_violence_experienced_fi"],
			"subtypeViolenceExperienced_ps": salvia_daos.VictimCaseForm2Enums["victim_case_form2_subtype_violence_experienced_ps"],
			"subtypeViolenceExperienced_se": salvia_daos.VictimCaseForm2Enums["victim_case_form2_subtype_violence_experienced_se"],
			"subtypeViolenceExperienced_po": salvia_daos.VictimCaseForm2Enums["victim_case_form2_subtype_violence_experienced_po"],
			"subtypeViolenceExperienced_pl": salvia_daos.VictimCaseForm2Enums["victim_case_form2_subtype_violence_experienced_pl"],
			"subtypeViolenceExperienced_re": salvia_daos.VictimCaseForm2Enums["victim_case_form2_subtype_violence_experienced_re"],
			"subtypeViolenceExperienced_vi": salvia_daos.VictimCaseForm2Enums["victim_case_form2_subtype_violence_experienced_vi"],
			"scopeOfViolence":               salvia_daos.VictimCaseForm2Enums["victim_case_form2_scope_of_violence"],
			"whoReportTo":                   salvia_daos.VictimCaseForm2Enums["victim_case_form2_who_report_to"],
			"activitiesUnableToPerform":     salvia_daos.VictimCaseForm2Enums["victim_case_form2_activities_unable_to_perform"],
			"adjustmentsGBV":                salvia_daos.VictimCaseForm2Enums["victim_case_form2_adjustments_gbv"],
			"law1996":                       salvia_daos.VictimCaseForm2Enums["victim_case_form2_law_1996"],
			"speciallyProtectedPopulation":  salvia_daos.VictimCaseForm2Enums["victim_case_form2_specially_protected_population"],
			"aspMode":                       salvia_daos.VictimCaseForm2Enums["victim_case_form2_asp_mode"],
			"reasonASP":                     salvia_daos.VictimCaseForm2Enums["victim_case_form2_reason_asp"],
			"hasDependents":                 salvia_daos.VictimCaseForm2Enums["victim_case_form2_has_dependents"],
			"actionPlan":                    salvia_daos.VictimCaseForm2Enums["victim_case_form2_action_plan"],

			//----------
			"victimCaseVictimNationality":                         salvia_config.VICTIM_CASE_VICTIM_NATIONALITY,
			"victimCaseVictimForeignerImmigrationStatus":          salvia_config.VICTIM_CASE_VICTIM_FOREIGNER_IMMIGRATION_STATUS,
			"victimCaseVictimGender":                              salvia_config.VICTIM_CASE_VICTIM_GENDER,
			"victimCaseVictimDependents":                          salvia_config.VICTIM_CASE_VICTIM_DEPENDENTS,
			"victimCaseVictimDeathThreats":                        common_config.YES_NO,
			"victimCaseVictimAggressorHasWeapons":                 common_config.YES_NO,
			"victimCaseExperiencedPhysicalOrSexualViolenceBefore": common_config.YES_NO,
			"victimCaseVictimImminentRisk":                        common_config.YES_NO,
			"victimCaseVictimPreviouslyReportedSituation":         common_config.YES_NO,
			"victimCaseVictimIfPreviouslyReported":                salvia_config.VICTIM_CASE_VICTIM_IF_PREVIOUSLY_REPORTED,

			"ifAfro":             salvia_config.AFRO_COMMUNTITIES[s.Lang],
			"ifIndigenous":       salvia_config.COLOMBIAN_INDIGENOUS[s.Lang],
			"ifIndigenousTongue": salvia_config.INDIGENOUS_TONGUES[s.Lang],
			"ifPeasant":          common_config.YES_NO,
			"ifArmedConflict":    common_config.YES_NO,
			"violenceScene":      salvia_config.VIOLENCE_SCENES[s.Lang],

			"moments":                    salvia_config.MOMENT[s.Lang],
			"sectors":                    salvia_config.SECTOR[s.Lang],
			"entities":                   entities,
			"entitiesByMomentAndSector":  entitiesByMomentAndSectorStr,
			"salviaFormPath":             salvia_config.FormPaths[s.Lang]["VictimCasePOST_GET"],
			"securityCityFormPath":       security_config.FormPaths[s.Lang]["CityGET"],
			"salviaEntityBranchFormPath": salvia_config.FormPaths[s.Lang]["EntityBranchGET"],
			"securityTownFormPath":       security_config.FormPaths[s.Lang]["TownGET"],
			"departments":                departments,
			"cities":                     cities,
			"towns":                      towns,
			"nav_rules":                  salvia_config.TranslateNavigationRule(s.Lang, salvia_config.NAVIGATION_RULES["set_victim_case"]),
		}, utils.GetFullHtmlFuncMap())
}

// VictimCasePUT maneja la solicitud PUT para actualizar o aprobar un caso de víctima.
// Según el parámetro "by" de la URL, decide si se aprueba el caso o se actualiza la información.
func VictimCasePUT(c *gin.Context) {
	var res string
	var code int
	// Se obtiene la sesión y se extrae el ID de sesión.
	session := sessions.Default(c)
	var sessionID string = session.Get("userData").(string)
	s, _ := utils.GetCommonSession(sessionID)

	// Se obtienen los parámetros "id" y "by" de la URL.
	id := c.Param("id")
	by := c.Param("by")

	form := c.Param("form")

	// Si "by" indica "approve_case" tras la traducción, se procede a aprobar el caso.
	if by != "" && salvia_config.TranslateLocale(by, s.Lang) == "approve_case" {
		// Verifica el permiso "approve_victim_case".
		if !utils.CheckPermission(salvia_config.PermissionsByRole, "approve_victim_case", s.CurrentRole, c) {
			return
		}
		code, res = salvia_ctrl.ApproveVictimCase(id, &db.ConnData{}, dbClientConfig, dbServerConfig)

	} else if by != "" && salvia_config.TranslateLocale(by, s.Lang) == "update" {
		// Si "by" indica "update", se verifica el permiso "update_victim_case".
		if !utils.CheckPermission(salvia_config.PermissionsByRole, "update_victim_case", s.CurrentRole, c) {
			return
		}
		// Se lee el cuerpo de la solicitud y se procesa la actualización.
		buf := new(bytes.Buffer)
		buf.ReadFrom(c.Request.Body)
		//code, res := salvia_ctrl.CtrlSetVictimCase(buf.String())
		if form == "form1" {
			code, res = salvia_ctrl.UpdateVictimCaseForm1(buf.String(), id, *s, dbClientConfig, dbServerConfig)
		} else if form == "form2" {
			code, res = salvia_ctrl.UpdateVictimCaseForm2(buf.String(), id, *s, dbClientConfig, dbServerConfig)
		}

	} else {
		// Si no se reconoce la acción, se responde con un error de solicitud incorrecta.
		code = http.StatusBadRequest
	}

	// Envía la respuesta al cliente en formato JSON.
	c.DataFromReader(code, int64(len(res)), gin.MIMEJSON, strings.NewReader(res), nil)
}

func UpdateVictimCasesOwnersAndRoles() {

	// Se ejecuta a bajo nivel
	salvia_ctrl.UpdateVictimCasesOwnersAndRoles(&db.ConnData{}, dbClientConfig, dbServerConfig)
	println("Fin de ejecución")
}

func mergeDocumentTypes(main, other map[string]string) map[string]string {
	// 1️⃣ Copiamos el mapa FORM2 (evitamos modificar el original)
	result := make(map[string]string, len(main))
	for k, v := range main {
		result[k] = v
	}

	// 2️⃣ Recorremos los demás elementos y sólo añadimos si la clave no está ya.
	for k, v := range other {
		if _, ok := result[k]; !ok { // clave nueva → insertamos
			result[k] = v
		}
	}

	return result
}
