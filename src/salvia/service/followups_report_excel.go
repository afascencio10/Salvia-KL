package service

import (
	"bitsflow/internal/repository"
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"
)

// Header definitions for each sheet
var (
	casesHeaders = []string{
		"Código Caso",
		"Fecha Creación (Bogotá)",
		"Documento Víctima",
		"Nombre Víctima",
		"Apellido Víctima",
		"Edad",
		"Género",
		"Teléfono Víctima",
		"Municipio",
		"Departamento",
		"Equipo Asignado",
		"Estado Caso",
	}

	registroHeaders = []string{
		"Código Caso",
		"Documento Víctima",
		"Tipo Documento",
		"Nombres Víctima",
		"Apellidos Víctima",
		"Edad",
		"Género",
		"Teléfono Víctima",
		"Dirección Víctima",
		"Correo Electrónico (Form 1)",
		"Nombre Identitario (Form 2)",
		"Orientación Sexual (Form 1/2)",
		"Procedencia (Form 1)",
		"Ocupación (Form 1/2)",
		"Grupo Étnico (Form 1/2)",
		"Si es Afrodescendiente (Form 1)",
		"Si es Indígena (Form 1)",
		"Si es Persona Campesina (Form 1)",
		"Número de Hijos (Form 1)",
		"Estado Civil (Form 2)",
		"Discapacidad u Otras Condiciones (Form 2)",
		"Municipio Residencia",
		"Departamento Residencia",
		"Nombres del Contacto (Form 1/2)",
		"Teléfono del Contacto (Form 1/2)",
		"Parentesco del Contacto (Form 1/2)",
		"Descripción de los Hechos (Form 1/2)",
		"Fecha Ocurrencia Hechos (Form 1/2)",
		"Hora Inicio Hechos (Form 1/2)",
		"Hora Fin Hechos (Form 1)",
		"Día de la Semana Hechos (Form 1)",
		"Ocurrencia de los Hechos (Form 1)",
		"Ámbito de la Violencia (Form 1)",
		"Escenario de Violencia (Form 1/2)",
		"Existe Riesgo Feminicida (Form 1)",
		"Nivel de Riesgo (Form 2)",
		"Nombre Agresor (Form 1/2)",
		"Tipo Documento Agresor (Form 1)",
		"Documento Agresor (Form 1/2)",
		"Dirección Agresor (Form 1/2)",
		"Teléfono Agresor (Form 1/2)",
		"Relación con Agresor (Form 1/2)",
		"Violencia Física Aumentada (Último Año) (Form 1)",
		"Separado de Pareja (Último Año) (Form 1)",
		"Amenazado con Arma/Objeto (Form 1)",
		"Amenaza de Muerte a Ella/Hijos (Form 1)",
		"Celoso Violento Constantemente (Form 1)",
		"Cree Capaz de Matarla (Form 1)",
		"Formulario Registro Completado (Form 1)",
		"Formulario Valoración Completado (Form 2)",
	}

	followUpsHeaders = []string{
		"Código Caso",
		"Documento Víctima",
		"ID Seguimiento",
		"Número Secuencia",
		"Fecha Planificada",
		"Fecha Realizado",
		"Estado",
		"Equipo",
		"Nivel Riesgo",
		"Resumen / Notas",
	}

	answersHeaders = []string{
		"Código Caso",
		"Documento Víctima",
		"ID Seguimiento",
		"Fecha Realizado Seguimiento",
		"ID Envío Formulario",
		"Pregunta",
		"Respuesta",
	}

	timelineHeaders = []string{
		"Código Caso",
		"Documento Víctima",
		"Fecha Evento",
		"Categoría",
		"Tipo Evento",
		"Descripción",
		"Actor",
	}
)

// formatGender normaliza los valores a strings legibles
func formatGender(g string) string {
	g = strings.TrimSpace(g)
	switch strings.ToLower(g) {
	case "w", "f", "femenino", "mujer":
		return "Mujer"
	case "m", "masculino", "hombre":
		return "Hombre"
	case "o", "otro":
		return "Otro"
	default:
		if g == "" {
			return "No registra"
		}
		return g
	}
}

func formatYesNo(val string) string {
	val = strings.TrimSpace(val)
	if val == "" {
		return "No registra"
	}
	switch strings.ToLower(val) {
	case "y", "yes", "s", "si", "sí", "true":
		return "Sí"
	case "n", "no", "false":
		return "No"
	default:
		return val
	}
}

func formatDocType(t string) string {
	t = strings.TrimSpace(strings.ToLower(t))
	switch t {
	case "cc":
		return "Cédula de Ciudadanía"
	case "ce":
		return "Cédula de Extranjería"
	case "pa":
		return "Pasaporte"
	case "pe":
		return "Permiso Especial de Permanencia"
	case "pt":
		return "Permiso por Protección Temporal"
	case "rc":
		return "Registro Civil"
	case "ti":
		return "Tarjeta de Identidad"
	default:
		if t == "" {
			return "No registra"
		}
		return strings.ToUpper(t)
	}
}

func formatOrigin(o string) string {
	o = strings.TrimSpace(strings.ToLower(o))
	switch o {
	case "u", "urbano":
		return "Urbano"
	case "r", "rural":
		return "Rural"
	default:
		if o == "" {
			return "No registra"
		}
		return o
	}
}

// BuildEmptyFollowUpsExcel genera un Excel solo con encabezados cuando no hay datos
func BuildEmptyFollowUpsExcel() (*excelize.File, error) {
	f := excelize.NewFile()
	defer f.Close()

	// Configurar hojas
	f.SetSheetName("Sheet1", "Casos")
	f.NewSheet("Registro")
	f.NewSheet("Seguimientos")
	f.NewSheet("Respuestas de formulario")
	f.NewSheet("Timeline")

	// Crear estilos
	headerStyle, err := createHeaderStyle(f)
	if err != nil {
		return nil, err
	}

	// Escribir encabezados en cada hoja
	if err := writeHeaders(f, "Casos", casesHeaders, headerStyle); err != nil {
		return nil, err
	}
	if err := writeHeaders(f, "Registro", registroHeaders, headerStyle); err != nil {
		return nil, err
	}
	if err := writeHeaders(f, "Seguimientos", followUpsHeaders, headerStyle); err != nil {
		return nil, err
	}
	if err := writeHeaders(f, "Respuestas de formulario", answersHeaders, headerStyle); err != nil {
		return nil, err
	}
	if err := writeHeaders(f, "Timeline", timelineHeaders, headerStyle); err != nil {
		return nil, err
	}

	// Auto-ajustar columnas
	autoFitColumns(f)

	return f, nil
}

// BuildFollowUpsExcel genera el Excel consolidado con toda la información cargada
func BuildFollowUpsExcel(
	cases []repository.CaseReportDTO,
	followUps []repository.FollowUpReportDTO,
	answers []repository.AnswerDTO,
	events []repository.TimelineReportDTO,
) (*excelize.File, error) {
	f := excelize.NewFile()

	// Configurar hojas
	f.SetSheetName("Sheet1", "Casos")
	f.NewSheet("Registro")
	f.NewSheet("Seguimientos")
	f.NewSheet("Respuestas de formulario")
	f.NewSheet("Timeline")

	// Crear estilos
	headerStyle, err := createHeaderStyle(f)
	if err != nil {
		f.Close()
		return nil, err
	}

	// Escribir encabezados
	if err := writeHeaders(f, "Casos", casesHeaders, headerStyle); err != nil {
		f.Close()
		return nil, err
	}
	if err := writeHeaders(f, "Registro", registroHeaders, headerStyle); err != nil {
		f.Close()
		return nil, err
	}
	if err := writeHeaders(f, "Seguimientos", followUpsHeaders, headerStyle); err != nil {
		f.Close()
		return nil, err
	}
	if err := writeHeaders(f, "Respuestas de formulario", answersHeaders, headerStyle); err != nil {
		f.Close()
		return nil, err
	}
	if err := writeHeaders(f, "Timeline", timelineHeaders, headerStyle); err != nil {
		f.Close()
		return nil, err
	}

	// 1. Llenar hoja "Casos"
	for rowIdx, c := range cases {
		rNum := rowIdx + 2
		_ = f.SetCellValue("Casos", fmt.Sprintf("A%d", rNum), c.VictimCaseICode)
		_ = f.SetCellValue("Casos", fmt.Sprintf("B%d", rNum), c.VictimCaseCreatedAt.Format("2006-01-02 15:04:05"))
		_ = f.SetCellValue("Casos", fmt.Sprintf("C%d", rNum), c.VictimCaseVictimDocNumber)
		_ = f.SetCellValue("Casos", fmt.Sprintf("D%d", rNum), c.VictimCaseVictimName)
		_ = f.SetCellValue("Casos", fmt.Sprintf("E%d", rNum), c.VictimCaseVictimLastName)
		if c.VictimCaseVictimAge != nil {
			_ = f.SetCellValue("Casos", fmt.Sprintf("F%d", rNum), *c.VictimCaseVictimAge)
		} else {
			_ = f.SetCellValue("Casos", fmt.Sprintf("F%d", rNum), "")
		}
		_ = f.SetCellValue("Casos", fmt.Sprintf("G%d", rNum), formatGender(c.VictimCaseVictimGender))
		_ = f.SetCellValue("Casos", fmt.Sprintf("H%d", rNum), c.VictimPhone)
		_ = f.SetCellValue("Casos", fmt.Sprintf("I%d", rNum), c.TownName)
		_ = f.SetCellValue("Casos", fmt.Sprintf("J%d", rNum), c.DeptName)
		_ = f.SetCellValue("Casos", fmt.Sprintf("K%d", rNum), c.VictimCaseTeam)
		_ = f.SetCellValue("Casos", fmt.Sprintf("L%d", rNum), c.VictimCaseStatus)
	}

	// 2. Llenar hoja "Registro" (Form 1 / Form 2)
	for rowIdx, c := range cases {
		rNum := rowIdx + 2
		_ = f.SetCellValue("Registro", fmt.Sprintf("A%d", rNum), c.VictimCaseICode)
		_ = f.SetCellValue("Registro", fmt.Sprintf("B%d", rNum), c.VictimCaseVictimDocNumber)
		_ = f.SetCellValue("Registro", fmt.Sprintf("C%d", rNum), formatDocType(c.VictimDocType))
		_ = f.SetCellValue("Registro", fmt.Sprintf("D%d", rNum), c.VictimCaseVictimName)
		_ = f.SetCellValue("Registro", fmt.Sprintf("E%d", rNum), c.VictimCaseVictimLastName)
		if c.VictimCaseVictimAge != nil {
			_ = f.SetCellValue("Registro", fmt.Sprintf("F%d", rNum), *c.VictimCaseVictimAge)
		} else {
			_ = f.SetCellValue("Registro", fmt.Sprintf("F%d", rNum), "")
		}
		_ = f.SetCellValue("Registro", fmt.Sprintf("G%d", rNum), formatGender(c.VictimCaseVictimGender))
		_ = f.SetCellValue("Registro", fmt.Sprintf("H%d", rNum), c.VictimPhone)
		_ = f.SetCellValue("Registro", fmt.Sprintf("I%d", rNum), c.VictimAddress)
		_ = f.SetCellValue("Registro", fmt.Sprintf("J%d", rNum), c.VictimEmail)
		_ = f.SetCellValue("Registro", fmt.Sprintf("K%d", rNum), c.VictimIdentityName)
		_ = f.SetCellValue("Registro", fmt.Sprintf("L%d", rNum), translateEnumKey(c.VictimSexualOrientation))
		_ = f.SetCellValue("Registro", fmt.Sprintf("M%d", rNum), formatOrigin(c.VictimOrigin))
		_ = f.SetCellValue("Registro", fmt.Sprintf("N%d", rNum), translateEnumKey(c.VictimOccupation))
		_ = f.SetCellValue("Registro", fmt.Sprintf("O%d", rNum), translateEnumKey(c.VictimEthnicGroup))
		_ = f.SetCellValue("Registro", fmt.Sprintf("P%d", rNum), translateEnumKey(c.VictimIfAfro))
		_ = f.SetCellValue("Registro", fmt.Sprintf("Q%d", rNum), translateEnumKey(c.VictimIfIndigenous))
		_ = f.SetCellValue("Registro", fmt.Sprintf("R%d", rNum), translateEnumKey(c.VictimIfPeasant))
		_ = f.SetCellValue("Registro", fmt.Sprintf("S%d", rNum), c.VictimChildrenNumber)
		_ = f.SetCellValue("Registro", fmt.Sprintf("T%d", rNum), translateEnumKey(c.VictimMaritalStatus))
		_ = f.SetCellValue("Registro", fmt.Sprintf("U%d", rNum), translateEnumKey(c.VictimDisability))
		_ = f.SetCellValue("Registro", fmt.Sprintf("V%d", rNum), c.TownName)
		_ = f.SetCellValue("Registro", fmt.Sprintf("W%d", rNum), c.DeptName)
		_ = f.SetCellValue("Registro", fmt.Sprintf("X%d", rNum), c.VictimContactNames)
		_ = f.SetCellValue("Registro", fmt.Sprintf("Y%d", rNum), c.VictimContactPhone)
		_ = f.SetCellValue("Registro", fmt.Sprintf("Z%d", rNum), translateEnumKey(c.VictimContactKinship))
		_ = f.SetCellValue("Registro", fmt.Sprintf("AA%d", rNum), c.FactsDescription)
		_ = f.SetCellValue("Registro", fmt.Sprintf("AB%d", rNum), c.FactsDate)
		_ = f.SetCellValue("Registro", fmt.Sprintf("AC%d", rNum), c.FactsStartTime)
		_ = f.SetCellValue("Registro", fmt.Sprintf("AD%d", rNum), c.FactsEndTime)
		_ = f.SetCellValue("Registro", fmt.Sprintf("AE%d", rNum), formatWeekday(c.FactsWeekday))
		_ = f.SetCellValue("Registro", fmt.Sprintf("AF%d", rNum), c.FactsOccurrence)
		_ = f.SetCellValue("Registro", fmt.Sprintf("AG%d", rNum), c.ViolenceScope)
		_ = f.SetCellValue("Registro", fmt.Sprintf("AH%d", rNum), translateEnumKey(c.ViolenceScene))
		_ = f.SetCellValue("Registro", fmt.Sprintf("AI%d", rNum), translateEnumKey(c.FemicideRisk))
		_ = f.SetCellValue("Registro", fmt.Sprintf("AJ%d", rNum), c.RiskLevel)
		_ = f.SetCellValue("Registro", fmt.Sprintf("AK%d", rNum), c.AggressorName)
		_ = f.SetCellValue("Registro", fmt.Sprintf("AL%d", rNum), formatDocType(c.AggressorDocType))
		_ = f.SetCellValue("Registro", fmt.Sprintf("AM%d", rNum), c.AggressorDocNumber)
		_ = f.SetCellValue("Registro", fmt.Sprintf("AN%d", rNum), c.AggressorAddress)
		_ = f.SetCellValue("Registro", fmt.Sprintf("AO%d", rNum), c.AggressorPhone)
		_ = f.SetCellValue("Registro", fmt.Sprintf("AP%d", rNum), translateEnumKey(c.RelationshipAggressor))
		_ = f.SetCellValue("Registro", fmt.Sprintf("AQ%d", rNum), translateEnumKey(c.PhysicalViolenceIncreased))
		_ = f.SetCellValue("Registro", fmt.Sprintf("AR%d", rNum), translateEnumKey(c.SeparatedFromPartnerLastYear))
		_ = f.SetCellValue("Registro", fmt.Sprintf("AS%d", rNum), translateEnumKey(c.ThreatenedWithWeapon))
		_ = f.SetCellValue("Registro", fmt.Sprintf("AT%d", rNum), translateEnumKey(c.ThreatenedToKillOrHarmChildren))
		_ = f.SetCellValue("Registro", fmt.Sprintf("AU%d", rNum), translateEnumKey(c.JealousAndViolent))
		_ = f.SetCellValue("Registro", fmt.Sprintf("AV%d", rNum), translateEnumKey(c.BelievesCapableOfKilling))
		_ = f.SetCellValue("Registro", fmt.Sprintf("AW%d", rNum), c.HasForm1)
		_ = f.SetCellValue("Registro", fmt.Sprintf("AX%d", rNum), c.HasForm2)
	}

	// 3. Llenar hoja "Seguimientos"
	for rowIdx, fu := range followUps {
		rNum := rowIdx + 2
		_ = f.SetCellValue("Seguimientos", fmt.Sprintf("A%d", rNum), fu.CaseID)
		_ = f.SetCellValue("Seguimientos", fmt.Sprintf("B%d", rNum), fu.VictimDocNumber)
		_ = f.SetCellValue("Seguimientos", fmt.Sprintf("C%d", rNum), fu.ID)
		_ = f.SetCellValue("Seguimientos", fmt.Sprintf("D%d", rNum), fu.SequenceNumber)
		_ = f.SetCellValue("Seguimientos", fmt.Sprintf("E%d", rNum), fu.ScheduledDate.Format("2006-01-02 15:04:05"))

		completedStr := ""
		if fu.CompletedAt != nil {
			completedStr = fu.CompletedAt.Format("2006-01-02 15:04:05")
		}
		_ = f.SetCellValue("Seguimientos", fmt.Sprintf("F%d", rNum), completedStr)
		_ = f.SetCellValue("Seguimientos", fmt.Sprintf("G%d", rNum), fu.Status)
		_ = f.SetCellValue("Seguimientos", fmt.Sprintf("H%d", rNum), fu.Team)

		riskStr := ""
		if fu.RiskStatus != nil {
			riskStr = *fu.RiskStatus
		}
		_ = f.SetCellValue("Seguimientos", fmt.Sprintf("I%d", rNum), riskStr)

		summaryStr := ""
		if fu.Summary != nil {
			summaryStr = *fu.Summary
		}
		_ = f.SetCellValue("Seguimientos", fmt.Sprintf("J%d", rNum), summaryStr)
	}

	// 4. Llenar hoja "Respuestas de formulario"
	for rowIdx, ans := range answers {
		rNum := rowIdx + 2
		_ = f.SetCellValue("Respuestas de formulario", fmt.Sprintf("A%d", rNum), ans.CaseICode)
		_ = f.SetCellValue("Respuestas de formulario", fmt.Sprintf("B%d", rNum), ans.VictimDocNumber)
		_ = f.SetCellValue("Respuestas de formulario", fmt.Sprintf("C%d", rNum), ans.FollowUpID)

		dateStr := ""
		if ans.FollowUpDate != nil {
			dateStr = ans.FollowUpDate.Format("2006-01-02 15:04:05")
		}
		_ = f.SetCellValue("Respuestas de formulario", fmt.Sprintf("D%d", rNum), dateStr)
		_ = f.SetCellValue("Respuestas de formulario", fmt.Sprintf("E%d", rNum), ans.FormSubmissionID)
		_ = f.SetCellValue("Respuestas de formulario", fmt.Sprintf("F%d", rNum), ans.QuestionDesc)
		_ = f.SetCellValue("Respuestas de formulario", fmt.Sprintf("G%d", rNum), ans.AnswerValue)
	}

	// 5. Llenar hoja "Timeline"
	for rowIdx, ev := range events {
		rNum := rowIdx + 2
		_ = f.SetCellValue("Timeline", fmt.Sprintf("A%d", rNum), ev.CaseID)
		_ = f.SetCellValue("Timeline", fmt.Sprintf("B%d", rNum), ev.VictimDocNumber)
		_ = f.SetCellValue("Timeline", fmt.Sprintf("C%d", rNum), ev.Date.Format("2006-01-02 15:04:05"))
		_ = f.SetCellValue("Timeline", fmt.Sprintf("D%d", rNum), ev.Category)
		_ = f.SetCellValue("Timeline", fmt.Sprintf("E%d", rNum), ev.Type)
		_ = f.SetCellValue("Timeline", fmt.Sprintf("F%d", rNum), ev.Description)
		_ = f.SetCellValue("Timeline", fmt.Sprintf("G%d", rNum), ev.ActorName)
	}

	// Auto-ajustar columnas
	autoFitColumns(f)

	return f, nil
}

// createHeaderStyle configura el estilo premium de los encabezados (Slate Blue oscuro, texto blanco y negrita)
func createHeaderStyle(f *excelize.File) (int, error) {
	return f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#1E293B"}, // Azul Slate oscuro premium
			Pattern: 1,
		},
		Font: &excelize.Font{
			Bold:   true,
			Color:  "FFFFFF",
			Size:   11,
			Family: "Segoe UI",
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   false,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "475569", Style: 1},
			{Type: "top", Color: "475569", Style: 1},
			{Type: "right", Color: "475569", Style: 1},
			{Type: "bottom", Color: "475569", Style: 1},
		},
	})
}

// writeHeaders escribe los encabezados de columna y les aplica la altura y el estilo
func writeHeaders(f *excelize.File, sheet string, headers []string, styleID int) error {
	for i, h := range headers {
		colName, err := excelize.ColumnNumberToName(i + 1)
		if err != nil {
			return err
		}
		cell := fmt.Sprintf("%s1", colName)
		_ = f.SetCellValue(sheet, cell, h)
		_ = f.SetCellStyle(sheet, cell, cell, styleID)
	}
	_ = f.SetRowHeight(sheet, 1, 26) // Altura premium de encabezado
	return nil
}

// autoFitColumns calcula y ajusta el ancho óptimo de cada columna de forma adaptativa
func autoFitColumns(f *excelize.File) {
	for _, sheet := range f.GetSheetList() {
		cols, err := f.GetCols(sheet)
		if err != nil {
			continue
		}
		for i, col := range cols {
			maxLen := 0
			for _, cellVal := range col {
				if len(cellVal) > maxLen {
					maxLen = len(cellVal)
				}
			}
			colName, err := excelize.ColumnNumberToName(i + 1)
			if err == nil {
				width := float64(maxLen) + 4
				if width < 12 {
					width = 12 // Ancho mínimo seguro
				}
				if width > 60 {
					width = 60 // Ancho máximo para evitar columnas monstruosas
				}
				_ = f.SetColWidth(sheet, colName, colName, width)
			}
		}
	}
}

func formatWeekday(w string) string {
	w = strings.TrimSpace(w)
	switch w {
	case "1":
		return "Lunes"
	case "2":
		return "Martes"
	case "3":
		return "Miércoles"
	case "4":
		return "Jueves"
	case "5":
		return "Viernes"
	case "6":
		return "Sábado"
	case "7":
		return "Domingo"
	default:
		return w
	}
}

func translateEnumKey(key string) string {
	key = strings.TrimSpace(key)
	if key == "" {
		return "No registra"
	}
	kLower := strings.ToLower(key)

	// Traducciones comunes de sí/no
	switch kLower {
	case "yes_no_y", "si", "sí", "y", "yes", "true":
		return "Sí"
	case "yes_no_n", "no", "n", "no_", "false":
		return "No"
	case "yes_no_nr", "nr", "no registra", "n/r":
		return "No registra"
	}

	// Orientación Sexual
	if strings.HasPrefix(kLower, "victim_case_form2_sexual_orientation_") {
		suffix := strings.TrimPrefix(kLower, "victim_case_form2_sexual_orientation_")
		switch suffix {
		case "he": return "Heterosexual"
		case "ga": return "Gay"
		case "le": return "Lesbiana"
		case "bi": return "Bisexual"
		case "ot": return "Otro"
		case "ns": return "No estoy segura/o/e"
		case "ni": return "No informa"
		}
	}

	// Grupo Étnico
	if strings.HasPrefix(kLower, "victim_case_form2_ethnic_affiliation_") {
		suffix := strings.TrimPrefix(kLower, "victim_case_form2_ethnic_affiliation_")
		switch suffix {
		case "af": return "Afrocolombiana(o)"
		case "ne": return "Negra(o)"
		case "pa": return "Palenquera(o)"
		case "ar": return "Raizal"
		case "in": return "Indígena"
		case "pr": return "Pueblo Rrom / Gitano"
		case "ns": return "No se autoreconoce/No informa"
		case "ni": return "Ninguno"
		}
	}

	// Estado Civil
	if strings.HasPrefix(kLower, "victim_case_form2_marital_status_") {
		suffix := strings.TrimPrefix(kLower, "victim_case_form2_marital_status_")
		switch suffix {
		case "so": return "Soltera(o)"
		case "ul": return "Unión libre / Unión marital de hecho"
		case "ca": return "Casada(o)"
		case "se": return "Separada(o)"
		case "di": return "Divorciada(o)"
		case "vi": return "Viuda(o)"
		case "ni": return "No informa"
		}
	}

	// Discapacidad
	if strings.HasPrefix(kLower, "feminicide_risk_form1_victim_disability_type_") {
		suffix := strings.TrimPrefix(kLower, "feminicide_risk_form1_victim_disability_type_")
		switch suffix {
		case "df": return "Discapacidad Física"
		case "ds": return "Discapacidad Sensorial"
		case "da": return "Discapacidad Auditiva"
		case "dv": return "Discapacidad Visual"
		case "so": return "Sordoceguera"
		case "di": return "Discapacidad Intelectual/Cognitiva"
		case "dm": return "Discapacidad Mental/Psicosocisocial"
		case "mu": return "Múltiple"
		}
	}

	// Ocupación (Form 2)
	if strings.HasPrefix(kLower, "victim_case_form2_occupation_") {
		suffix := strings.TrimPrefix(kLower, "victim_case_form2_occupation_")
		switch suffix {
		case "sg": return "Servicios generales y aseo"
		case "ca": return "Cocina y alimentos"
		case "bc": return "Belleza y cuidado personal"
		case "ac": return "Ama de casa (Trabajo doméstico no remunerado)"
		case "cu": return "Cuidado (cuidadora, niñera(o), etc.)"
		case "td": return "Trabajo doméstico remunerado"
		case "cv": return "Comercio y ventas"
		case "cc": return "Atención al cliente / call center / recepción"
		case "se": return "Seguridad"
		case "tm": return "Transporte y mensajería"
		case "co": return "Construcción y oficios de obra"
		case "mp": return "Manufactura y producción"
		case "ao": return "Administrativo y oficina"
		case "sc": return "Salud y cuidado"
		case "ed": return "Educación"
		case "ar": return "Agro / rural"
		case "tp": return "Tecnología / plataformas"
		case "es": return "Estudiante"
		case "as": return "Persona en Actividades Sexuales Pagas (ASP)"
		case "dh": return "Defensora de derechos humanos"
		case "fp": return "Funcionaria(o) pública(o)"
		case "pe": return "Periodista"
		case "vt": return "Turista/visitante temporal"
		case "em": return "Empresaria(o)"
		case "ds": return "Dirigente sindical"
		case "af": return "Afiliada(o) sindical"
		case "ni": return "No informa"
		case "ot": return "Otra"
		}
	}

	// Parentesco del Contacto (Form 2)
	if strings.HasPrefix(kLower, "victim_case_form2_support_contact_kinship_") {
		suffix := strings.TrimPrefix(kLower, "victim_case_form2_support_contact_kinship_")
		switch suffix {
		case "ma": return "Madre / padre"
		case "he": return "Hermano(a)"
		case "hi": return "Hijo(a)"
		case "pa": return "Pareja o expareja"
		case "ot": return "Otro familiar"
		case "am": return "Amigo(a) o conocido(a)"
		case "je": return "Jefe(a) o colega de trabajo"
		case "cu": return "Otro"
		case "no": return "No informa"
		}
	}

	// Escenario de Violencia (Form 2)
	if strings.HasPrefix(kLower, "victim_case_form2_scenario_violence_") {
		suffix := strings.TrimPrefix(kLower, "victim_case_form2_scenario_violence_")
		switch suffix {
		case "am": return "Ambulancia - transporte sanitario"
		case "ar": return "Áreas deportivas y/o recreativas"
		case "ca": return "Calle"
		case "cr": return "Carretera"
		case "ce": return "Centro de atención médica"
		case "co": return "Centros de reclusión"
		case "ed": return "Centros educativos"
		case "ac": return "Espacios acuáticos al aire libre"
		case "te": return "Espacios terrestres al aire libre"
		case "cm": return "Establecimiento comercial"
		case "in": return "Establecimiento industrial"
		case "ex": return "Establecimientos de expendio de comidas"
		case "ad": return "Establecimientos dedicados a la administración pública"
		case "fi": return "Establecimientos financieros"
		case "es": return "Estaciones de servicio"
		case "gu": return "Guarniciones militares y/o de Policía"
		case "cu": return "Lugares de actividades culturales"
		case "ci": return "Lugares de cuidado de personas"
		case "al": return "Lugares de esparcimiento con expendio de alcohol"
		case "mi": return "Lugares de explotación de minas"
		case "ho": return "Lugares de hospedaje"
		case "of": return "Oficinas de empresas privadas"
		case "pa": return "Parqueaderos"
		case "si": return "Sitio de culto"
		case "ta": return "Terminales de pasajeros"
		case "tr": return "Transporte masivo"
		case "sp": return "Vehículo de servicio particular"
		case "vp": return "Vehículo servicio público"
		case "vi": return "Vivienda"
		case "za": return "Zonas de actividades agropecuarias"
		case "ot": return "Otro"
		}
	}

	// Relación con el Agresor (Form 2)
	if strings.HasPrefix(kLower, "victim_case_form2_relationship_with_presumed_aggressor_") {
		suffix := strings.TrimPrefix(kLower, "victim_case_form2_relationship_with_presumed_aggressor_")
		// Quitar posibles prefijos numéricos intermedios como "01_" o "02_"
		suffix = strings.TrimPrefix(suffix, "01_")
		suffix = strings.TrimPrefix(suffix, "02_")
		switch suffix {
		case "am": return "Amigo(a)"
		case "en": return "Encargado del cuidado"
		case "he": return "Hermano(a)"
		case "ma": return "Madre"
		case "of": return "Otros familiares"
		case "pa": return "Padre"
		case "pi": return "Pareja íntima"
		case "ex": return "Expareja o antigua pareja"
		case "rp", "rc": return "Representante de la entidad"
		case "su": return "Supervisor(a) / empleador(a) / jefe(a)"
		case "co": return "Compañero(a) de trabajo"
		case "pr", "di": return "Profesor(a) / instructor(a) / coach"
		case "ce": return "Compañero(a) de estudios"
		case "ps": return "Proveedor(a) de servicios"
		case "ci", "cp": return "Coinquilino(a)"
		case "af": return "Amigo de la familia"
		case "ve": return "Vecino(a)"
		case "cl": return "Clientes"
		case "sr": return "Sin relación"
		case "dc": return "Delincuencia común"
		case "de": return "Persona desmovilizada"
		case "mg": return "Miembro de grupos alzados al margen de la ley"
		case "ms": return "Personal de seguridad privada"
		case "pc": return "Personal de custodia"
		case "sp", "fp": return "Servidor(a) público(a)"
		case "lr": return "Líder(esa) religioso(a) o espiritual"
		case "pn": return "Miembro de Policía Nacional"
		case "an": return "Miembro de la Armada Nacional"
		case "fa": return "Miembro de la Fuerza Aérea"
		case "is": return "Personal de institución de salud"
		case "ls": return "Líder(esa) social"
		case "as": return "Afiliado(a) sindical"
		case "ds": return "Dirigente sindical"
		case "ni": return "No informa"
		case "ot": return "Otro"
		}
	}

	// Mapeo genérico por sufijo
	if strings.HasSuffix(kLower, "_he") { return "Heterosexual" }
	if strings.HasSuffix(kLower, "_ga") { return "Gay" }
	if strings.HasSuffix(kLower, "_le") { return "Lesbiana" }
	if strings.HasSuffix(kLower, "_bi") { return "Bisexual" }
	if strings.HasSuffix(kLower, "_ot") { return "Otro" }
	if strings.HasSuffix(kLower, "_ns") { return "No informa" }
	if strings.HasSuffix(kLower, "_ni") { return "No informa" }
	
	if strings.HasSuffix(kLower, "_af") { return "Afrocolombiana(o)" }
	if strings.HasSuffix(kLower, "_ne") { return "Negra(o)" }
	if strings.HasSuffix(kLower, "_pa") { return "Palenquera(o)" }
	if strings.HasSuffix(kLower, "_ar") { return "Raizal" }
	if strings.HasSuffix(kLower, "_in") { return "Indígena" }
	if strings.HasSuffix(kLower, "_pr") { return "Pueblo Rrom / Gitano" }
	
	if strings.HasSuffix(kLower, "_so") { return "Soltera(o)" }
	if strings.HasSuffix(kLower, "_ul") { return "Unión libre / Unión marital de hecho" }
	if strings.HasSuffix(kLower, "_ca") { return "Casada(o)" }
	if strings.HasSuffix(kLower, "_se") { return "Separada(o)" }
	if strings.HasSuffix(kLower, "_di") { return "Divorciada(o)" }
	if strings.HasSuffix(kLower, "_vi") { return "Viuda(o)" }

	return key
}
