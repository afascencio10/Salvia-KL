// Package salvia_ctrl — PSWYA-755: Registro de Contacto Víctima V2 (Flujo A).
//
// SetVictimContactV2 es una función nueva e independiente de SetVictimContact
// (VictimContactController.go). No comparte su lógica de validación ni de
// persistencia de alto nivel — decisión 2A de la planeación (ver
// DocsMD/Screens/Registro de Contacto Víctima V2/related-tables.md). Sí
// reutiliza, sin modificarlas, las funciones DAO existentes que ya saben
// escribir victim_contact y victim_contact_form2, más dos funciones de
// paquete (getAndVerifyVictimContactEnumsMultiple / setVictimContactEnumsMultiple)
// ya definidas en VictimContactController.go — viven en el mismo paquete, así
// que reutilizarlas no requiere tocar ese archivo.
package salvia_ctrl

import (
	"bytes"
	"encoding/json"
	"net/http"

	common_config "bitsflow/common/config"
	common_dao "bitsflow/common/dao"
	"bitsflow/common/db"
	"bitsflow/common/utils"

	salvia_config "bitsflow/salvia/config"
	salvia_daos "bitsflow/salvia/dao"
	security_config "bitsflow/security/config"

	"github.com/dchest/captcha"
)

// victimContactV2Envelope solo captura el objeto crudo "victimContact" — se
// decodifica DOS VECES a partir de los mismos bytes (ver más abajo) en vez de
// componer un struct que embeba VictimContactDTO. VictimContactDTO implementa
// su propio UnmarshalJSON (VictimContactDAO.go); al embeberlo anónimamente
// dentro de otro struct, encoding/json deja de "aplanar" sus campos
// promovidos y el campo hermano (form2V2) nunca se llenaba — confirmado
// depurando: form2V2 llegaba con icode="" pese a que el JSON de entrada sí
// lo traía. Dos decodificaciones independientes evitan el problema por
// completo.
type victimContactV2Envelope struct {
	VictimContact json.RawMessage `json:"victimContact"`
}

type victimContactForm2V2Wrapper struct {
	Form2V2 salvia_daos.VictimContactForm2V2DTO `json:"form2V2"`
}

// Códigos de la clave de enum victim_case_form2_report_type_v2 — ver migración
// migrations/add_victim_contact_v2.sql
const (
	reportTypeV2Victima  string = "v"
	reportTypeV2Tercero  string = "r"
	reportTypeV2Familiar string = "fm"
	reportTypeV2Hombre   string = "ho"
)

// SetVictimContactV2 crea un reporte ciudadano V2: valida según el tipo de
// reportante (Víctima/Tercero/Familiar/Hombre), inserta victim_contact +
// victim_contact_form2 (tablas de V1, sin modificarlas) y victim_contact_form2_v2
// (tabla de extensión, exclusiva de V2). Devuelve el icode del VictimContact
// creado como número de radicado.
func SetVictimContactV2(dataInput string, checkCaptcha bool, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	var err error

	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	var collectedErrors map[string]map[string]string = map[string]map[string]string{}

	var envelope victimContactV2Envelope
	utils.JSONToStruct(dataInput, &envelope)

	var vc salvia_daos.VictimContactDTO
	utils.JSONToStruct(string(envelope.VictimContact), &vc)

	var form2v2Wrapper victimContactForm2V2Wrapper
	utils.JSONToStruct(string(envelope.VictimContact), &form2v2Wrapper)

	form2 := &vc.VictimContactForm2
	form2v2 := &form2v2Wrapper.Form2V2

	salvia_daos.SetVictimContactDefaults(&vc, common_dao.SQL_INSERT)
	salvia_daos.SetVictimContactForm2Defaults(form2, common_dao.SQL_INSERT)
	salvia_daos.SetVictimContactForm2V2Defaults(form2v2, common_dao.SQL_INSERT)

	// 1. Captcha — igual que V1 (bypass sólo para el cliente Flutter).
	captchaPassed := !checkCaptcha || captcha.VerifyString(vc.VictimContactCaptchaID, vc.VictimContactCaptchaSolution)
	if !captchaPassed {
		utils.SetError(collectedErrors, salvia_daos.VictimContactJSONName, "captchaSolution", "security_general_user_captcha_fail", "common_global_error", security_config.Locale)
	}

	// 2. Autorización de datos — debe existir y resolver a "Sí".
	if err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&form2v2.VictimContactForm2V2AuthorizationAnswer); err != nil {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2V2JSONName, "authorizationAnswer", "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
	} else if form2v2.VictimContactForm2V2AuthorizationAnswer.VictimCaseForm2EnumsCode != "y" {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2V2JSONName, "authorizationAnswer", "common_validation_field_empty_error", "common_global_error", common_config.Locale)
	}

	// 3. Tipo de reporte — debe resolver a una fila de la clave V2.
	if err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&form2.VictimContactForm2ReportType); err != nil {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, "reportType", "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
	}
	reportTypeCode := form2.VictimContactForm2ReportType.VictimCaseForm2EnumsCode

	// 4. Ajuste razonable — al menos una opción, todas deben resolver.
	// (getAndVerifyVictimContactEnumsMultiple ya registra el error en collectedErrors)
	getAndVerifyVictimContactEnumsMultiple(form2, collectedErrors)

	// 5. Descripción de los hechos — obligatoria para todos los tipos.
	if len(form2.VictimContactForm2FactsDescription) < 10 {
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, "factsDescription", "common_validation_field_empty_error", "common_global_error", common_config.Locale)
	}

	// 6. Validaciones específicas por tipo de reportante.
	switch reportTypeCode {

	case reportTypeV2Victima, reportTypeV2Tercero:
		// Datos de la persona afectada (ella misma, o descrita por el tercero).
		if vc.VictimContactNames == "" {
			utils.SetError(collectedErrors, salvia_daos.VictimContactJSONName, "names", "common_validation_field_empty_error", "common_global_error", common_config.Locale)
		}
		if vc.VictimContactLastNames == "" {
			utils.SetError(collectedErrors, salvia_daos.VictimContactJSONName, "lastNames", "common_validation_field_empty_error", "common_global_error", common_config.Locale)
		}
		if form2.VictimContactForm2VictimColPhone == 0 {
			utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, "victimColPhone", "common_validation_field_number_error", "common_global_error", common_config.Locale)
		}
		if err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&form2.VictimContactForm2WillReceiveCall); err != nil {
			utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, "willReceiveCall", "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		}

		if reportTypeCode == reportTypeV2Tercero {
			// Anonimato del reportante.
			if err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&form2v2.VictimContactForm2V2Anonimato); err != nil {
				utils.SetError(collectedErrors, salvia_daos.VictimContactForm2V2JSONName, "anonimato", "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
			} else if form2v2.VictimContactForm2V2Anonimato.VictimCaseForm2EnumsCode == "pr" {
				if form2.VictimContactForm2ReporterNames == "" {
					utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, "reporterNames", "common_validation_field_empty_error", "common_global_error", common_config.Locale)
				}
				if form2v2.VictimContactForm2V2ReporterLastNames == "" {
					utils.SetError(collectedErrors, salvia_daos.VictimContactForm2V2JSONName, "reporterLastNames", "common_validation_field_empty_error", "common_global_error", common_config.Locale)
				}
			}
			if err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&form2.VictimContactForm2HasCareRole); err != nil {
				utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, "hasCareRole", "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
			}
			if err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&form2.VictimContactForm2VictimAwareOfReport); err != nil {
				utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, "victimAwareOfReport", "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
			}
		}

	case reportTypeV2Familiar, reportTypeV2Hombre:
		// Datos de quien reporta — siempre identificado (no hay opción de anonimato).
		if form2.VictimContactForm2ReporterNames == "" {
			utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, "reporterNames", "common_validation_field_empty_error", "common_global_error", common_config.Locale)
		}
		if form2v2.VictimContactForm2V2ReporterLastNames == "" {
			utils.SetError(collectedErrors, salvia_daos.VictimContactForm2V2JSONName, "reporterLastNames", "common_validation_field_empty_error", "common_global_error", common_config.Locale)
		}
		if form2.VictimContactForm2ReporterPhone == 0 {
			utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, "reporterPhone", "common_validation_field_number_error", "common_global_error", common_config.Locale)
		}
		if err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&form2v2.VictimContactForm2V2ReporterWillReceiveCall); err != nil {
			utils.SetError(collectedErrors, salvia_daos.VictimContactForm2V2JSONName, "reporterWillReceiveCall", "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
		} else if form2v2.VictimContactForm2V2ReporterWillReceiveCall.VictimCaseForm2EnumsCode == "y" && form2v2.VictimContactForm2V2ReporterBestContactTime.IsZero() {
			utils.SetError(collectedErrors, salvia_daos.VictimContactForm2V2JSONName, "reporterBestContactTime", "common_validation_field_empty_error", "common_global_error", common_config.Locale)
		}

		if reportTypeCode == reportTypeV2Familiar {
			// Datos de la víctima (fallecida) — solo nombre, sin teléfono ni disposición.
			if vc.VictimContactNames == "" {
				utils.SetError(collectedErrors, salvia_daos.VictimContactJSONName, "names", "common_validation_field_empty_error", "common_global_error", common_config.Locale)
			}
			if vc.VictimContactLastNames == "" {
				utils.SetError(collectedErrors, salvia_daos.VictimContactJSONName, "lastNames", "common_validation_field_empty_error", "common_global_error", common_config.Locale)
			}
			if err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&form2v2.VictimContactForm2V2Parentesco); err != nil {
				utils.SetError(collectedErrors, salvia_daos.VictimContactForm2V2JSONName, "parentesco", "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
			}
			if err = salvia_daos.GetLocalVictimCaseForm2EnumsByICode(&form2.VictimContactForm2HasCareRole); err != nil {
				utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, "hasCareRole", "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
			}
		} else {
			// Hombre: no hay "persona afectada" nombrada — decisión 7A, se deja
			// pendiente para que el agente lo complete tras el recontacto.
			vc.VictimContactNames = ""
			vc.VictimContactLastNames = ""
		}
		// victim_contact_form2_victim_col_phone / best_contact_time son NOT NULL
		// en la tabla existente (no se puede alterar esa restricción sin tocar
		// V1) — para Familiar/Hombre, donde no aplica, se guardan en 0 / vacío.
		form2.VictimContactForm2VictimColPhone = 0

	default:
		utils.SetError(collectedErrors, salvia_daos.VictimContactForm2JSONName, "reportType", "victim_case_form2_enums_not_found", "common_global_error", salvia_config.Locale)
	}

	if len(collectedErrors) > 0 {
		collectedErrors["default"]["captchaID"] = captcha.New()
		return http.StatusBadRequest, utils.CommMsgGetJSONErrors(collectedErrors)
	}

	// --- Persistencia ---------------------------------------------------

	if _, err = db.StartTransaction(connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	if err = salvia_daos.SetVictimContact(&vc, connData, &dbClientConfig, &dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, err.Error()
	}

	form2.VictimContactForm2VictimContact = vc

	if err = salvia_daos.SetVictimContactForm2(form2, connData, &dbClientConfig, &dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, err.Error()
	}

	if err = setVictimContactEnumsMultiple(form2, connData, dbClientConfig, dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, err.Error()
	}

	form2v2.VictimContactForm2V2VictimContactForm2 = *form2

	if err = salvia_daos.SetVictimContactForm2V2(form2v2, connData, &dbClientConfig, &dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, err.Error()
	}

	if _, err = db.CommitTransaction(connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// A diferencia de V1 (que responde con cuerpo vacío), V2 sí devuelve el
	// icode: el wizard lo muestra como número de radicado en la confirmación.
	var buf bytes.Buffer
	buf.WriteString(`{"icode":"`)
	buf.WriteString(vc.VictimContactICode)
	buf.WriteString(`"}`)

	return http.StatusOK, buf.String()
}
