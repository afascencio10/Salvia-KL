// Package salvia_facades — PSWYA-755: Registro de Contacto Víctima V2 (Flujo A).
//
// Handlers públicos, independientes de VictimContactPOST_Public /
// VictimContactPOST_GET_Public (VictimContactFacade.go, V1) — decisión 2A de
// la planeación. No requieren sesión ni permisos, igual que la V1.
package salvia_facades

import (
	"bytes"
	"net/http"
	"strings"

	"bitsflow/common/db"
	common_routers "bitsflow/common/facades"
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	salvia_ctrl "bitsflow/salvia/controllers"
	salvia_daos "bitsflow/salvia/dao"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
)

// VictimContactPOST_Public_V2 maneja la creación pública del reporte V2.
// Igual que la versión V1, no valida sesión ni permisos.
func VictimContactPOST_Public_V2(c *gin.Context) {
	var checkCaptcha bool = true

	buf := new(bytes.Buffer)
	buf.ReadFrom(c.Request.Body)

	if v, found := c.Request.Header["User-Agent"]; found && v[0] == "flutter-client" {
		checkCaptcha = false
	}

	code, res := salvia_ctrl.SetVictimContactV2(buf.String(), checkCaptcha, &db.ConnData{}, dbClientConfig, dbServerConfig)
	c.DataFromReader(code, int64(len(res)), gin.MIMEJSON, strings.NewReader(res), nil)
}

// VictimContactPOST_GET_Public_V2 renderiza el wizard de 4 pasos del reporte
// V2, o devuelve sus catálogos en JSON si el cliente pide application/json.
func VictimContactPOST_GET_Public_V2(c *gin.Context) {
	common_routers.SetHeaderNoCache(c)

	captchaID := captcha.New()

	if v, found := c.Request.Header["Accept"]; found && v[0] == "application/json" {

		type VictimContactV2Fields struct {
			YesNo          []salvia_daos.VictimCaseForm2EnumsDTO `json:"yesNo"`
			ReportType     []salvia_daos.VictimCaseForm2EnumsDTO `json:"reportType"`
			AdjustmentsGBV []salvia_daos.VictimCaseForm2EnumsDTO `json:"adjustmentsGBV"`
			Anonimato      []salvia_daos.VictimCaseForm2EnumsDTO `json:"anonimato"`
			Parentesco     []salvia_daos.VictimCaseForm2EnumsDTO `json:"parentesco"`
			CaptchaID      string                                `json:"captchaID"`
		}

		fields := VictimContactV2Fields{
			YesNo:          salvia_daos.VictimCaseForm2Enums["yes_no"],
			ReportType:     salvia_daos.VictimCaseForm2Enums["victim_case_form2_report_type_v2"],
			AdjustmentsGBV: salvia_daos.VictimCaseForm2Enums["victim_case_form2_adjustments_gbv"],
			Anonimato:      salvia_daos.VictimCaseForm2Enums["victim_case_form2_anonimato"],
			Parentesco:     salvia_daos.VictimCaseForm2Enums["victim_case_form2_parentesco"],
			CaptchaID:      captchaID,
		}
		fieldsStr, _ := utils.DTOtoString(fields)
		c.DataFromReader(http.StatusOK, int64(len(fieldsStr)), gin.MIMEJSON, strings.NewReader(fieldsStr), nil)
		return
	}

	common_routers.RenderTemplate(c, salvia_daos.VictimContactEntityName, "salvia", "victim_contact/", salvia_config.HTML_Templates, "set_victim_contact_v2", utils.GetFullHtmlTemplates(), utils.DEFAULT_VIEW, utils.DEFAULT_PANIC_TEMPLATE,
		map[string]interface{}{
			"windowTitle":    salvia_config.Locale["sp"]["set_victim_contact_window_title"],
			"lang":           "sp",
			"locale":         salvia_config.Locale,
			"yes_no":         salvia_daos.VictimCaseForm2Enums["yes_no"],
			"reportTypeV2":   salvia_daos.VictimCaseForm2Enums["victim_case_form2_report_type_v2"],
			"adjustmentsGBV": salvia_daos.VictimCaseForm2Enums["victim_case_form2_adjustments_gbv"],
			"anonimato":      salvia_daos.VictimCaseForm2Enums["victim_case_form2_anonimato"],
			"parentesco":     salvia_daos.VictimCaseForm2Enums["victim_case_form2_parentesco"],
			"CaptchaID":      captchaID,
			"formPath":       salvia_config.FormPaths["sp"]["VictimContactPOST_Public_V2"],
		},
		utils.GetFullHtmlFuncMap())
}
