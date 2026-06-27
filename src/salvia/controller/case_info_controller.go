// Package controller — case_info_controller.go
// Endpoint para el componente <case-info>: información completa del caso.
package controller

import (
	"bitsflow/salvia/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CaseInfoController struct {
	svc service.CaseInfoService
	db  *gorm.DB
}

func NewCaseInfoController(svc service.CaseInfoService, db ...*gorm.DB) *CaseInfoController {
	ctrl := &CaseInfoController{svc: svc}
	if len(db) > 0 {
		ctrl.db = db[0]
	}
	return ctrl
}

func (c *CaseInfoController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/casos/:id/info-completa", c.GetFullInfo)
	rg.GET("/casos/:id/tamizaje", c.GetTamizaje)
}

// GetFullInfo retorna toda la información del caso combinando victim_case + form1 + form2.
// GET /api/v1/casos/:id/info-completa
func (c *CaseInfoController) GetFullInfo(ctx *gin.Context) {
	caseICode := ctx.Param("id")
	if caseICode == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id requerido"})
		return
	}

	info, err := c.svc.GetFullCaseInfo(ctx.Request.Context(), caseICode)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, info)
}

// GetTamizaje retorna las preguntas del tamizaje de riesgo con sus respuestas.
// Soporta form2 (enums) y form1 (campos directos y/n).
// GET /api/v1/casos/:id/tamizaje
func (c *CaseInfoController) GetTamizaje(ctx *gin.Context) {
	caseICode := ctx.Param("id")
	if caseICode == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id requerido"})
		return
	}
	if c.db == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "servicio no disponible"})
		return
	}

	type preguntaResp struct {
		Pregunta  string `json:"pregunta"`
		Respuesta string `json:"respuesta"`
	}

	// ── Form2: query con sub-selects a enums ──
	type tamRow struct {
		Q1  string `gorm:"column:q1"`
		Q2  string `gorm:"column:q2"`
		Q3  string `gorm:"column:q3"`
		Q4  string `gorm:"column:q4"`
		Q5  string `gorm:"column:q5"`
		Q6  string `gorm:"column:q6"`
		Q7  string `gorm:"column:q7"`
		Q8  string `gorm:"column:q8"`
		Q9  string `gorm:"column:q9"`
		Q10 string `gorm:"column:q10"`
		Q11 string `gorm:"column:q11"`
		Q12 string `gorm:"column:q12"`
		Q13 string `gorm:"column:q13"`
		Q14 string `gorm:"column:q14"`
		Q15 string `gorm:"column:q15"`
		Q16 string `gorm:"column:q16"`
		Q17 string `gorm:"column:q17"`
		Q18 string `gorm:"column:q18"`
		Q19 string `gorm:"column:q19"`
		Q20 string `gorm:"column:q20"`
		Q21 string `gorm:"column:q21"`
		RiskLevel *int `gorm:"column:risk_level"`
		RiskScore *int `gorm:"column:risk_score"`
	}
	var row tamRow
	c.db.Raw(`
		SELECT
			COALESCE((SELECT e.victim_case_form2_enums_name FROM salvia.victim_case_form2_enums e WHERE e.victim_case_form2_enums_id = f2.victim_case_form2_aggressor_violence_physical_increase), '') AS q1,
			COALESCE((SELECT e.victim_case_form2_enums_name FROM salvia.victim_case_form2_enums e WHERE e.victim_case_form2_enums_id = f2.victim_case_form2_aggressor_weapon_used), '') AS q2,
			COALESCE((SELECT e.victim_case_form2_enums_name FROM salvia.victim_case_form2_enums e WHERE e.victim_case_form2_enums_id = f2.victim_case_form2_aggressor_threat_kill), '') AS q3,
			COALESCE((SELECT e.victim_case_form2_enums_name FROM salvia.victim_case_form2_enums e WHERE e.victim_case_form2_enums_id = f2.victim_case_form2_aggressor_pursues_spies_destroys), '') AS q4,
			COALESCE((SELECT e.victim_case_form2_enums_name FROM salvia.victim_case_form2_enums e WHERE e.victim_case_form2_enums_id = f2.victim_case_form2_aggressor_capable_of_killing), '') AS q5,
			COALESCE((SELECT e.victim_case_form2_enums_name FROM salvia.victim_case_form2_enums e WHERE e.victim_case_form2_enums_id = f2.victim_case_form2_aggressor_has_access_to_weapons), '') AS q6,
			COALESCE((SELECT e.victim_case_form2_enums_name FROM salvia.victim_case_form2_enums e WHERE e.victim_case_form2_enums_id = f2.victim_case_form2_partner_unemployed), '') AS q7,
			COALESCE((SELECT e.victim_case_form2_enums_name FROM salvia.victim_case_form2_enums e WHERE e.victim_case_form2_enums_id = f2.victim_case_form2_partner_other_denunciations), '') AS q8,
			COALESCE((SELECT e.victim_case_form2_enums_name FROM salvia.victim_case_form2_enums e WHERE e.victim_case_form2_enums_id = f2.victim_case_form2_aggressor_has_penal_background), '') AS q9,
			COALESCE((SELECT e.victim_case_form2_enums_name FROM salvia.victim_case_form2_enums e WHERE e.victim_case_form2_enums_id = f2.victim_case_form2_aggressor_forced_sex), '') AS q10,
			COALESCE((SELECT e.victim_case_form2_enums_name FROM salvia.victim_case_form2_enums e WHERE e.victim_case_form2_enums_id = f2.victim_case_form2_aggressor_attempted_strangulation), '') AS q11,
			COALESCE((SELECT e.victim_case_form2_enums_name FROM salvia.victim_case_form2_enums e WHERE e.victim_case_form2_enums_id = f2.victim_case_form2_aggressor_consumes_drugs), '') AS q12,
			COALESCE((SELECT e.victim_case_form2_enums_name FROM salvia.victim_case_form2_enums e WHERE e.victim_case_form2_enums_id = f2.victim_case_form2_aggressor_is_alcoholic), '') AS q13,
			COALESCE((SELECT e.victim_case_form2_enums_name FROM salvia.victim_case_form2_enums e WHERE e.victim_case_form2_enums_id = f2.victim_case_form2_partner_controls), '') AS q14,
			COALESCE((SELECT e.victim_case_form2_enums_name FROM salvia.victim_case_form2_enums e WHERE e.victim_case_form2_enums_id = f2.victim_case_form2_aggressor_had_hit_in_vulnerability), '') AS q15,
			COALESCE((SELECT e.victim_case_form2_enums_name FROM salvia.victim_case_form2_enums e WHERE e.victim_case_form2_enums_id = f2.victim_case_form2_partner_threatened_suicide), '') AS q16,
			COALESCE((SELECT e.victim_case_form2_enums_name FROM salvia.victim_case_form2_enums e WHERE e.victim_case_form2_enums_id = f2.victim_case_form2_partner_threatened_damage_members), '') AS q17,
			COALESCE((SELECT e.victim_case_form2_enums_name FROM salvia.victim_case_form2_enums e WHERE e.victim_case_form2_enums_id = f2.victim_case_form2_thoughts_of_self_harm), '') AS q18,
			COALESCE((SELECT e.victim_case_form2_enums_name FROM salvia.victim_case_form2_enums e WHERE e.victim_case_form2_enums_id = f2.victim_case_form2_aggressor_limits_contact_support_networks), '') AS q19,
			COALESCE((SELECT e.victim_case_form2_enums_name FROM salvia.victim_case_form2_enums e WHERE e.victim_case_form2_enums_id = f2.victim_case_form2_still_lives_with_aggressor), '') AS q20,
			COALESCE((SELECT e.victim_case_form2_enums_name FROM salvia.victim_case_form2_enums e WHERE e.victim_case_form2_enums_id = f2.victim_case_form2_aggressor_violently_jealous), '') AS q21,
			f2.victim_case_form2_risk_level AS risk_level,
			f2.victim_case_form2_risk_score AS risk_score
		FROM salvia.victim_case_form2 f2
		JOIN salvia.victim_case vc ON vc.victim_case_id = f2.victim_case_form2_victim_case
		WHERE vc.victim_case_i_code = ?
		LIMIT 1
	`, caseICode).Scan(&row)

	resolveYN := func(val string) string {
		switch val {
		case "yes_no_y", "yes_no_01":
			return "Sí"
		case "yes_no_n", "yes_no_02":
			return "No"
		default:
			if val != "" {
				return val
			}
			return ""
		}
	}

	preguntas := []preguntaResp{
		{"¿Ha aumentado la violencia física en severidad o frecuencia en los últimos 60 días?", resolveYN(row.Q1)},
		{"¿La ha agredido o amenazado utilizando algún arma u objeto en su contra?", resolveYN(row.Q2)},
		{"¿Ha amenazado con matarla?", resolveYN(row.Q3)},
		{"¿La persigue o espía, le deja notas amenazantes o mensajes, destruye sus cosas, o le llama cuando no quiere?", resolveYN(row.Q4)},
		{"¿Usted cree que es capaz de matarla?", resolveYN(row.Q5)},
		{"¿La persona agresora tiene acceso a armas?", resolveYN(row.Q6)},
		{"¿Está su pareja / expareja desempleado(a) o sin trabajo actualmente?", resolveYN(row.Q7)},
		{"¿Cuenta con otras denuncias por violencia intrafamiliar o incumplimiento de medidas?", resolveYN(row.Q8)},
		{"¿Tiene antecedentes penales u otras denuncias por otro delito?", resolveYN(row.Q9)},
		{"¿Su pareja / expareja la ha forzado a mantener relaciones sexuales?", resolveYN(row.Q10)},
		{"¿Ha intentado alguna vez estrangularla?", resolveYN(row.Q11)},
		{"¿Su pareja / expareja consume drogas?", resolveYN(row.Q12)},
		{"¿Su pareja / expareja es alcohólica o tiene problemas con el alcohol?", resolveYN(row.Q13)},
		{"¿Su pareja / expareja controla la mayoría de sus actividades diarias?", resolveYN(row.Q14)},
		{"¿La ha golpeado en estado de vulnerabilidad física o embarazo?", resolveYN(row.Q15)},
		{"¿Su pareja/expareja la ha amenazado con suicidarse o lo ha intentado?", resolveYN(row.Q16)},
		{"¿Ha amenazado con hacer daño a sus hijos, mascotas u otros familiares?", resolveYN(row.Q17)},
		{"¿Ha llegado a pensar en hacerse daño o en suicidarse?", resolveYN(row.Q18)},
		{"¿La persona agresora limita su contacto con familiares o redes de apoyo?", resolveYN(row.Q19)},
		{"¿Aún convive con la persona agresora?", resolveYN(row.Q20)},
		{"¿Es celoso(a) con usted constante o violentamente?", resolveYN(row.Q21)},
	}

	var resultado []preguntaResp
	for _, p := range preguntas {
		if p.Respuesta != "" {
			resultado = append(resultado, p)
		}
	}

	// ── Fallback Form1: campos directos y/n ──
	nivelTexto := ""
	if len(resultado) == 0 {
		type f1Row struct {
			DeathThreats       string `gorm:"column:q1"`
			HasWeapons         string `gorm:"column:q2"`
			ViolenceBefore     string `gorm:"column:q3"`
			PhysicalIncreased  string `gorm:"column:q4"`
			SeparatedLastYear  string `gorm:"column:q5"`
			ThreatenedWeapon   string `gorm:"column:q6"`
			ThreatenedChildren string `gorm:"column:q7"`
			JealousViolent     string `gorm:"column:q8"`
			CapableKilling     string `gorm:"column:q9"`
			ImminentRisk       string `gorm:"column:q10"`
			PreviouslyReported string `gorm:"column:q11"`
			FemicideRisk       string `gorm:"column:q12"`
		}
		var f1 f1Row
		c.db.Raw(`
			SELECT
				COALESCE(f1.victim_case_form1_victim_death_threats, '') AS q1,
				COALESCE(f1.victim_case_form1_victim_aggressor_has_weapons, '') AS q2,
				COALESCE(f1.victim_case_form1_experienced_physical_or_sexual_violence_befor, '') AS q3,
				COALESCE(f1.victim_case_form1_physical_violence_increased, '') AS q4,
				COALESCE(f1.victim_case_form1_separated_from_partner_last_year, '') AS q5,
				COALESCE(f1.victim_case_form1_threatened_with_weapon, '') AS q6,
				COALESCE(f1.victim_case_form1_threatened_to_kill_or_harm_children, '') AS q7,
				COALESCE(f1.victim_case_form1_jealous_and_violent, '') AS q8,
				COALESCE(f1.victim_case_form1_believes_capable_of_killing, '') AS q9,
				COALESCE(f1.victim_case_form1_victim_imminent_risk, '') AS q10,
				COALESCE(f1.victim_case_form1_victim_previously_reported_situation, '') AS q11,
				COALESCE(f1.victim_case_form1_victim_femicide_risk, '') AS q12
			FROM salvia.victim_case_form1 f1
			JOIN salvia.victim_case vc ON vc.victim_case_id = f1.victim_case_form1_victim_case
			WHERE vc.victim_case_i_code = ?
			LIMIT 1
		`, caseICode).Scan(&f1)

		resolveF1 := func(val string) string {
			switch val {
			case "y", "1":
				return "Sí"
			case "n", "0":
				return "No"
			}
			return ""
		}

		preguntasF1 := []preguntaResp{
			{"¿Ha amenazado con matarla?", resolveF1(f1.DeathThreats)},
			{"¿La persona agresora tiene acceso a armas?", resolveF1(f1.HasWeapons)},
			{"¿Ha experimentado violencia física o sexual anteriormente?", resolveF1(f1.ViolenceBefore)},
			{"¿Ha aumentado la violencia física en severidad o frecuencia?", resolveF1(f1.PhysicalIncreased)},
			{"¿Se ha separado de su pareja en el último año?", resolveF1(f1.SeparatedLastYear)},
			{"¿Le ha amenazado con arma?", resolveF1(f1.ThreatenedWeapon)},
			{"¿Ha amenazado con hacer daño a sus hijos o familiares?", resolveF1(f1.ThreatenedChildren)},
			{"¿Es celoso(a) y violento(a)?", resolveF1(f1.JealousViolent)},
			{"¿Usted cree que es capaz de matarla?", resolveF1(f1.CapableKilling)},
			{"¿Existe riesgo inminente?", resolveF1(f1.ImminentRisk)},
			{"¿Ha realizado denuncia previa?", resolveF1(f1.PreviouslyReported)},
			{"¿Existe riesgo feminicida?", resolveF1(f1.FemicideRisk)},
		}
		for _, p := range preguntasF1 {
			if p.Respuesta != "" {
				resultado = append(resultado, p)
			}
		}
		if f1.FemicideRisk == "1" || f1.FemicideRisk == "y" {
			nivelTexto = "Alto (feminicidio)"
		}
	}

	// Nivel de riesgo (form2)
	if row.RiskLevel != nil {
		switch *row.RiskLevel {
		case 1:
			nivelTexto = "Bajo"
		case 2:
			nivelTexto = "Moderado"
		case 3:
			nivelTexto = "Alto"
		case 4:
			nivelTexto = "Extremo"
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"preguntas":  resultado,
		"riskLevel":  row.RiskLevel,
		"riskScore":  row.RiskScore,
		"nivelTexto": nivelTexto,
	})
}
