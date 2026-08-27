package service

import (
	"bitsflow/internal/repository"
	salvia_config "bitsflow/salvia/config"
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
		"Utiliza Salud para Chantajear (Form 2)",
		"Amenaza con Revelar Orientación (Form 2)",
		"Dejó de Buscar Ayuda (Form 2)",
		"Forzada a Relaciones Sexuales (Form 2)",
		"Aprovechó Vulnerabilidad Física (Form 2)",
		"Violencia Motivada por Género/Orientación (Form 2)",
		"Formulario Registro Completado (Form 1)",
		"Formulario Valoración Completado (Form 2)",
		// Bloque 1: Autorización y Contacto
		"Autorización Datos Personales",
		"Ajustes Razonables (Form 2)",
		"Intérprete de Idiomas (Form 2)",
		"Correo Contacto Apoyo/Emergencia (Form 2)",
		// Bloque 2: Hechos Victimizantes (ampliación)
		"Departamento Ocurrencia Hechos (Form 2)",
		"Municipio Ocurrencia Hechos (Form 2)",
		"Zona Ocurrencia Hechos (Form 2)",
		"Dirección Ocurrencia Hechos (Form 2)",
		"Tipo de Violencia (Form 2)",
		"Subtipo de Violencia (Form 2)",
		"Ámbito de Violencia (Form 2)",
		"Sector Violencia Laboral (Form 2)",
		"Veces Ocurrencia Agresión (Form 2)",
		// Bloque 3: Características del Agresor
		"Número de Agresores (Form 2)",
		"Cercanía con Persona Agresora (Form 2)",
		"Dependencia Económica del Agresor (Form 2)",
		"Identidad de Género Agresor (Form 2)",
		// Bloque 4: Tamizaje (preguntas adicionales)
		"Persigue o Espía (Form 2)",
		"Agresor Tiene Acceso a Armas (Form 2)",
		"Agresor Desempleado (Form 2)",
		"Otras Denuncias por VIF (Form 2)",
		"Antecedentes Penales (Form 2)",
		"Intento de Estrangulación (Form 2)",
		"Consumo de Drogas (Form 2)",
		"Consumo de Alcohol (Form 2)",
		"Control de Actividades Diarias (Form 2)",
		"Amenaza con Suicidio (Form 2)",
		"Amenaza Daño a Hijos/Familiares (Form 2)",
		"Piensa en Suicidarse (Form 2)",
		"Limita Contacto con Redes de Apoyo (Form 2)",
		"Convivencia con Agresor (Form 2)",
		// Bloque 5: Datos Personales Víctima
		"Fecha de Nacimiento (Form 2)",
		"Dificultad para Realizar Actividades (Form 2)",
		"Nacionalidad (Form 2)",
		"País de Nacionalidad (Form 2)",
		"Condición Migratoria (Form 2)",
		"Identidad de Género Víctima (Form 2)",
		"Sexo Asignado al Nacer (Form 2)",
		"Población de Especial Protección (Form 2)",
		"Nivel de Escolaridad (Form 2)",
		"Forma Generación de Ingresos (Form 2)",
		"Tipo Vinculación Laboral (Form 2)",
		"Modalidad ASP (Form 2)",
		"Fecha Inicio ASP (Form 2)",
		"Motivo Inicio ASP (Form 2)",
		"Tenencia de Vivienda (Form 2)",
		"Estrato de Vivienda (Form 2)",
		"Personas/Animales a Cargo (Form 2)",
		"Estado de Gestación (Form 2)",
		"Zona de Residencia (Form 2)",
		// Bloque 6: Plan de Acción, Denuncia, Operación
		"Plan de Acción (Form 2)",
		"Explicación Gestión Salvia (Form 2)",
		"Permiso Denuncia Fiscalía (Form 2)",
		"Agente Responsable",
		"Historial de Funcionarios",
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
		// Si es un nombre de enum completo (ej: victim_case_form2_gender_identity_mu), traducirlo
		if strings.Contains(g, "gender_identity") {
			return translateSingleEnum(g)
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
	// Headers de "Respuestas de formulario" se escriben dinámicamente más abajo (pivotado)
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
		_ = f.SetCellValue("Registro", fmt.Sprintf("AW%d", rNum), translateEnumKey(c.VictimHealthToBlackmail))
		_ = f.SetCellValue("Registro", fmt.Sprintf("AX%d", rNum), translateEnumKey(c.ThreatenedRevealSexualOrientation))
		_ = f.SetCellValue("Registro", fmt.Sprintf("AY%d", rNum), translateEnumKey(c.StoppedSeekingHelp))
		_ = f.SetCellValue("Registro", fmt.Sprintf("AZ%d", rNum), translateEnumKey(c.AggressorSexuallyHarassment2))
		_ = f.SetCellValue("Registro", fmt.Sprintf("BA%d", rNum), translateEnumKey(c.AggressorTakenAdvantagePhysicalVulnerabil))
		_ = f.SetCellValue("Registro", fmt.Sprintf("BB%d", rNum), translateEnumKey(c.ViolenceMotivatedByGender2))
		_ = f.SetCellValue("Registro", fmt.Sprintf("BC%d", rNum), c.HasForm1)
		_ = f.SetCellValue("Registro", fmt.Sprintf("BD%d", rNum), c.HasForm2)
		// Bloque 1: Autorización y Contacto
		_ = f.SetCellValue("Registro", fmt.Sprintf("BE%d", rNum), translateEnumKey(c.AuthorizationAnswer))
		_ = f.SetCellValue("Registro", fmt.Sprintf("BF%d", rNum), translateEnumKey(c.AdjustmentsGBV))
		_ = f.SetCellValue("Registro", fmt.Sprintf("BG%d", rNum), translateEnumKey(c.RequireInterpreter))
		_ = f.SetCellValue("Registro", fmt.Sprintf("BH%d", rNum), c.SupportContactEmail)
		// Bloque 2: Hechos Victimizantes (ampliación)
		_ = f.SetCellValue("Registro", fmt.Sprintf("BI%d", rNum), c.FactsDeptName)
		_ = f.SetCellValue("Registro", fmt.Sprintf("BJ%d", rNum), c.FactsCityName)
		_ = f.SetCellValue("Registro", fmt.Sprintf("BK%d", rNum), translateEnumKey(c.FactsZone))
		_ = f.SetCellValue("Registro", fmt.Sprintf("BL%d", rNum), c.FactsAddressForm2)
		_ = f.SetCellValue("Registro", fmt.Sprintf("BM%d", rNum), translateEnumKey(c.ViolenceType))
		_ = f.SetCellValue("Registro", fmt.Sprintf("BN%d", rNum), translateEnumKey(c.ViolenceSubtype))
		_ = f.SetCellValue("Registro", fmt.Sprintf("BO%d", rNum), translateEnumKey(c.ViolenceScopeForm2))
		_ = f.SetCellValue("Registro", fmt.Sprintf("BP%d", rNum), translateEnumKey(c.WorkplaceSector))
		_ = f.SetCellValue("Registro", fmt.Sprintf("BQ%d", rNum), translateEnumKey(c.RecurrenceAggression))
		// Bloque 3: Características del Agresor
		_ = f.SetCellValue("Registro", fmt.Sprintf("BR%d", rNum), translateEnumKey(c.NumAggressors))
		_ = f.SetCellValue("Registro", fmt.Sprintf("BS%d", rNum), translateEnumKey(c.ProximityAggressor))
		_ = f.SetCellValue("Registro", fmt.Sprintf("BT%d", rNum), translateEnumKey(c.EconomicallyDependent))
		_ = f.SetCellValue("Registro", fmt.Sprintf("BU%d", rNum), translateEnumKey(c.AggressorGenderIdentity))
		// Bloque 4: Tamizaje
		_ = f.SetCellValue("Registro", fmt.Sprintf("BV%d", rNum), translateEnumKey(c.AggressorPursuesSpies))
		_ = f.SetCellValue("Registro", fmt.Sprintf("BW%d", rNum), translateEnumKey(c.AggressorHasAccessWeapons))
		_ = f.SetCellValue("Registro", fmt.Sprintf("BX%d", rNum), translateEnumKey(c.PartnerUnemployed))
		_ = f.SetCellValue("Registro", fmt.Sprintf("BY%d", rNum), translateEnumKey(c.PartnerOtherDenunciations))
		_ = f.SetCellValue("Registro", fmt.Sprintf("BZ%d", rNum), translateEnumKey(c.AggressorPenalBackground))
		_ = f.SetCellValue("Registro", fmt.Sprintf("CA%d", rNum), translateEnumKey(c.AggressorStrangulation))
		_ = f.SetCellValue("Registro", fmt.Sprintf("CB%d", rNum), translateEnumKey(c.AggressorConsumesDrugs))
		_ = f.SetCellValue("Registro", fmt.Sprintf("CC%d", rNum), translateEnumKey(c.AggressorIsAlcoholic))
		_ = f.SetCellValue("Registro", fmt.Sprintf("CD%d", rNum), translateEnumKey(c.PartnerControls))
		_ = f.SetCellValue("Registro", fmt.Sprintf("CE%d", rNum), translateEnumKey(c.PartnerThreatenedSuicide))
		_ = f.SetCellValue("Registro", fmt.Sprintf("CF%d", rNum), translateEnumKey(c.PartnerThreatenedDamage))
		_ = f.SetCellValue("Registro", fmt.Sprintf("CG%d", rNum), translateEnumKey(c.ThoughtsOfSelfHarm))
		_ = f.SetCellValue("Registro", fmt.Sprintf("CH%d", rNum), translateEnumKey(c.AggressorLimitsContact))
		_ = f.SetCellValue("Registro", fmt.Sprintf("CI%d", rNum), translateEnumKey(c.StillLivesWithAggressor))
		// Bloque 5: Datos Personales Víctima
		_ = f.SetCellValue("Registro", fmt.Sprintf("CJ%d", rNum), c.BirthDate)
		_ = f.SetCellValue("Registro", fmt.Sprintf("CK%d", rNum), translateEnumKey(c.PhysicalDifficulties))
		_ = f.SetCellValue("Registro", fmt.Sprintf("CL%d", rNum), translateEnumKey(c.Nationality))
		_ = f.SetCellValue("Registro", fmt.Sprintf("CM%d", rNum), translateEnumKey(c.SpecifiedNationality))
		_ = f.SetCellValue("Registro", fmt.Sprintf("CN%d", rNum), translateEnumKey(c.MigrationCondition))
		_ = f.SetCellValue("Registro", fmt.Sprintf("CO%d", rNum), translateEnumKey(c.GenderIdentity))
		_ = f.SetCellValue("Registro", fmt.Sprintf("CP%d", rNum), translateEnumKey(c.AssignedSexAtBirth))
		_ = f.SetCellValue("Registro", fmt.Sprintf("CQ%d", rNum), translateEnumKey(c.SpecialProtectedPop))
		_ = f.SetCellValue("Registro", fmt.Sprintf("CR%d", rNum), translateEnumKey(c.LastEducationLevel))
		_ = f.SetCellValue("Registro", fmt.Sprintf("CS%d", rNum), translateEnumKey(c.IncomeGenerationMethod))
		_ = f.SetCellValue("Registro", fmt.Sprintf("CT%d", rNum), translateEnumKey(c.EmploymentRelationship))
		_ = f.SetCellValue("Registro", fmt.Sprintf("CU%d", rNum), translateEnumKey(c.ASPMode))
		_ = f.SetCellValue("Registro", fmt.Sprintf("CV%d", rNum), c.ApproxStartASP)
		_ = f.SetCellValue("Registro", fmt.Sprintf("CW%d", rNum), translateEnumKey(c.ReasonASP))
		_ = f.SetCellValue("Registro", fmt.Sprintf("CX%d", rNum), translateEnumKey(c.HousingTenancy))
		_ = f.SetCellValue("Registro", fmt.Sprintf("CY%d", rNum), translateEnumKey(c.HousingStratum))
		_ = f.SetCellValue("Registro", fmt.Sprintf("CZ%d", rNum), translateEnumKey(c.HasDependents))
		_ = f.SetCellValue("Registro", fmt.Sprintf("DA%d", rNum), translateEnumKey(c.CurrentlyPregnant))
		_ = f.SetCellValue("Registro", fmt.Sprintf("DB%d", rNum), c.ResidenceZone)
		// Bloque 6: Plan de Acción, Denuncia, Operación
		_ = f.SetCellValue("Registro", fmt.Sprintf("DC%d", rNum), translateEnumKey(c.ActionPlan))
		_ = f.SetCellValue("Registro", fmt.Sprintf("DD%d", rNum), c.ManagementExplanation)
		_ = f.SetCellValue("Registro", fmt.Sprintf("DE%d", rNum), translateEnumKey(c.AllowsEasyReport))
		_ = f.SetCellValue("Registro", fmt.Sprintf("DF%d", rNum), c.AgentResponsible)
		_ = f.SetCellValue("Registro", fmt.Sprintf("DG%d", rNum), c.OwnerDescription)
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

	// 4. Llenar hoja "Respuestas de formulario" (formato pivotado: 1 fila por seguimiento)
	{
		sheet := "Respuestas de formulario"

		// 4a. Extraer preguntas únicas en orden de aparición
		var questionOrder []string
		questionSeen := map[string]bool{}
		for _, ans := range answers {
			if !questionSeen[ans.QuestionDesc] {
				questionSeen[ans.QuestionDesc] = true
				questionOrder = append(questionOrder, ans.QuestionDesc)
			}
		}

		// 4b. Construir headers dinámicos: columnas fijas + preguntas como columnas
		dynHeaders := []string{"Código Caso", "Documento Víctima", "ID Seguimiento", "Fecha Realizado Seguimiento"}
		dynHeaders = append(dynHeaders, questionOrder...)
		if err := writeHeaders(f, sheet, dynHeaders, headerStyle); err != nil {
			return nil, err
		}

		// 4c. Agrupar respuestas por followUpID (mantener orden de aparición)
		type followUpRow struct {
			CaseICode      string
			VictimDoc      string
			FollowUpID     string
			FollowUpDate   string
			Answers        map[string]string
		}
		var rowOrder []string
		rowMap := map[string]*followUpRow{}

		for _, ans := range answers {
			row, exists := rowMap[ans.FollowUpID]
			if !exists {
				dateStr := ""
				if ans.FollowUpDate != nil {
					dateStr = ans.FollowUpDate.Format("2006-01-02 15:04:05")
				}
				row = &followUpRow{
					CaseICode:    ans.CaseICode,
					VictimDoc:    ans.VictimDocNumber,
					FollowUpID:   ans.FollowUpID,
					FollowUpDate: dateStr,
					Answers:      map[string]string{},
				}
				rowMap[ans.FollowUpID] = row
				rowOrder = append(rowOrder, ans.FollowUpID)
			}
			// Si ya hay una respuesta para esta pregunta en este seguimiento, concatenar
			if existing, ok := row.Answers[ans.QuestionDesc]; ok && existing != "" {
				row.Answers[ans.QuestionDesc] = existing + "; " + ans.AnswerValue
			} else {
				row.Answers[ans.QuestionDesc] = ans.AnswerValue
			}
		}

		// 4d. Escribir filas pivotadas
		for rowIdx, fuID := range rowOrder {
			rNum := rowIdx + 2
			row := rowMap[fuID]
			_ = f.SetCellValue(sheet, fmt.Sprintf("A%d", rNum), row.CaseICode)
			_ = f.SetCellValue(sheet, fmt.Sprintf("B%d", rNum), row.VictimDoc)
			_ = f.SetCellValue(sheet, fmt.Sprintf("C%d", rNum), row.FollowUpID)
			_ = f.SetCellValue(sheet, fmt.Sprintf("D%d", rNum), row.FollowUpDate)
			// Columnas dinámicas de preguntas (empiezan en E)
			for qIdx, question := range questionOrder {
				colName, _ := excelize.ColumnNumberToName(qIdx + 5) // E=5, F=6, etc.
				val := row.Answers[question]
				// Traducir true/false
				switch strings.ToLower(strings.TrimSpace(val)) {
				case "true":
					val = "Sí"
				case "false":
					val = "No"
				default:
					// Traducir multi-select separados por coma usando labels de option
					if strings.Contains(val, ",") || strings.Contains(val, "_") {
						parts := strings.Split(val, ",")
						var translated []string
						for _, p := range parts {
							p = strings.TrimSpace(p)
							if p == "" {
								continue
							}
							// Humanizar: reemplazar _ por espacios y capitalizar
							readable := strings.ReplaceAll(p, "_", " ")
							readable = strings.ToUpper(readable[:1]) + readable[1:]
							translated = append(translated, readable)
						}
						if len(translated) > 0 {
							val = strings.Join(translated, ", ")
						}
					}
				}
				_ = f.SetCellValue(sheet, fmt.Sprintf("%s%d", colName, rNum), val)
			}
		}
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

	// Si contiene comas, es multi-select: traducir cada valor por separado
	if strings.Contains(key, ",") {
		parts := strings.Split(key, ",")
		var translated []string
		for _, p := range parts {
			t := translateSingleEnum(strings.TrimSpace(p))
			if t != "" {
				translated = append(translated, t)
			}
		}
		if len(translated) > 0 {
			return strings.Join(translated, ", ")
		}
		return key
	}

	return translateSingleEnum(key)
}

func translateSingleEnum(key string) string {
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

	// Identidad de Género (Form 2)
	if strings.HasPrefix(kLower, "victim_case_form2_gender_identity_") {
		suffix := strings.TrimPrefix(kLower, "victim_case_form2_gender_identity_")
		switch suffix {
		case "mu": return "Mujer"
		case "mt": return "Mujer Transgénero"
		case "ho": return "Hombre"
		case "ht": return "Hombre Transgénero"
		case "nb": return "No binaria - asignado masculino al nacer"
		case "nf": return "No binaria - asignado femenino al nacer"
		case "ot": return "Otra"
		case "ni": return "No informa"
		}
	}

	// Ajustes Razonables (Form 2)
	if strings.HasPrefix(kLower, "victim_case_form2_adjustments_gbv_") {
		suffix := strings.TrimPrefix(kLower, "victim_case_form2_adjustments_gbv_")
		switch suffix {
		case "sc": return "Intérprete de lengua de señas colombiana"
		case "gi": return "Guía-intérprete para personas sordociegas"
		case "id": return "Intérprete de idiomas / traducción"
		case "ap": return "Apoyos para la comunicación (pictogramas, tableros, etc.)"
		case "in": return "Información en formatos accesibles (lectura fácil, macrotipo, audio, braille)"
		case "ad": return "Adecuaciones de accesibilidad física"
		case "dt": return "Dispositivos o ayudas técnicas"
		case "pr": return "Presencia de persona de confianza"
		case "tr": return "Transporte seguro o acompañamiento"
		case "ni": return "No informa"
		case "nr": return "No requiere"
		case "ns": return "No requiere apoyos especiales"
		}
	}

	// Zona de ocurrencia hechos (Form 2)
	if strings.HasPrefix(kLower, "victim_case_form2_facts_zone_") {
		suffix := strings.TrimPrefix(kLower, "victim_case_form2_facts_zone_")
		switch suffix {
		case "cm": return "Cabecera municipal"
		case "pr": return "Parte rural (vereda y campo)"
		case "cp": return "Centro poblado (corregimiento, inspección de policía, caserío)"
		}
	}

	// Tipo de violencia experimentada (Form 2)
	if strings.HasPrefix(kLower, "victim_case_form2_type_violence_experienced_") {
		suffix := strings.TrimPrefix(kLower, "victim_case_form2_type_violence_experienced_")
		switch suffix {
		case "fi": return "Física"
		case "ps": return "Psicológica / emocional"
		case "se": return "Sexual"
		case "po": return "Patrimonial o económica"
		case "pl": return "Política"
		case "re": return "Reproductiva"
		case "vi": return "Vicaria"
		}
	}

	// Subtipo de violencia experimentada (Form 2)
	if strings.HasPrefix(kLower, "victim_case_form2_subtype_violence_experienced_") {
		suffix := strings.TrimPrefix(kLower, "victim_case_form2_subtype_violence_experienced_")
		switch suffix {
		case "fi_ac": return "Ataque con agente químico"
		case "fi_gp": return "Golpizas"
		case "fi_qu": return "Quemadura"
		case "fi_so": return "Sofocamiento"
		case "fi_ao": return "Agresión con objetos"
		case "fi_to": return "Tortura"
		case "fi_se": return "Secuestro"
		case "ps_sv": return "Subvaloración, descalificación, humillación"
		case "ps_ia": return "Intimidación y Amenaza"
		case "ps_av": return "Acoso, vigilancia y seguimiento"
		case "ps_ce": return "Control coercitivo por celos"
		case "ps_ai": return "Aislamiento"
		case "ps_al": return "Acoso laboral"
		case "se_ac": return "Acceso carnal violento"
		case "se_ab": return "Abuso sexual"
		case "se_tr": return "Trata de personas con fines de explotación sexual"
		case "se_mu": return "Mutilación de órganos sexuales"
		case "se_ao": return "Acoso sexual"
		case "po_ia": return "Inasistencia alimentaria"
		case "po_vm": return "Vigilancia y manipulación de dinero e ingresos"
		case "po_nr": return "Negación, retención, sustracción o destrucción de bienes"
		case "po_ri": return "Remuneración inequitativa"
		case "pl_aa": return "Abuso de autoridad con fines políticos"
		case "pl_cp": return "Censura y persecución política"
		case "re_ab": return "Aborto forzado"
		case "re_es": return "Esterilización y/o planificación forzada"
		case "re_de": return "Denegación de servicios IVE"
		case "re_ma": return "Matrimonio infantil y/o uniones tempranas"
		case "re_vi": return "Violencia obstétrica"
		case "vi_ah": return "Amenazas contra hijas/os"
		case "vi_pr": return "Privación de libertad y secuestro parental"
		}
	}

	// Ámbito de violencia (Form 2)
	if strings.HasPrefix(kLower, "victim_case_form2_scope_of_violence_") {
		suffix := strings.TrimPrefix(kLower, "victim_case_form2_scope_of_violence_")
		switch suffix {
		case "ae": return "Ámbito educativo"
		case "ac": return "Ámbito de conflicto armado"
		case "as": return "Ámbito de la salud"
		case "ai": return "Ámbito institucional"
		case "ic": return "Ámbito intrafamiliar (conviviente)"
		case "al": return "Ámbito laboral"
		case "ap": return "Ámbito político"
		case "ri": return "Ámbito reclusión intramural"
		case "rn": return "Ámbito relacional (no conviven)"
		case "ad": return "Ámbito digital"
		case "co": return "Ámbito comunitario y/o vecinal"
		case "de": return "Ámbito deportivo"
		case "ce": return "Ámbito cultural, entretenimiento y ocio"
		}
	}

	// Veces ocurrencia agresión (Form 2)
	if strings.HasPrefix(kLower, "victim_case_form2_recurrence_aggression_") {
		suffix := strings.TrimPrefix(kLower, "victim_case_form2_recurrence_aggression_")
		switch suffix {
		case "pr": return "Es la primera vez que la agrede"
		case "se": return "Es la segunda vez que la agrede"
		case "re": return "La agresión ha sido recurrente (más de 2 veces)"
		}
	}

	// Número de agresores (Form 2)
	if strings.HasPrefix(kLower, "victim_case_form2_num_agressors_") {
		suffix := strings.TrimPrefix(kLower, "victim_case_form2_num_agressors_")
		switch suffix {
		case "01": return "1"
		case "02": return "2"
		case "03": return "3"
		case "md": return "Más de 3"
		case "nd": return "No se puede determinar / no informa"
		}
	}

	// Cercanía con persona agresora (Form 2)
	if strings.HasPrefix(kLower, "victim_case_form2_proximity_principal_aggressor_") {
		suffix := strings.TrimPrefix(kLower, "victim_case_form2_proximity_principal_aggressor_")
		switch suffix {
		case "pc": return "Persona conocida con la que convive"
		case "pn": return "Persona conocida con la que no convive"
		case "pd": return "Persona desconocida"
		case "nd": return "No se puede determinar / no informa"
		}
	}

	// Identidad de género agresor (Form 2)
	if strings.HasPrefix(kLower, "victim_case_form2_aggressor_gender_identity_") {
		suffix := strings.TrimPrefix(kLower, "victim_case_form2_aggressor_gender_identity_")
		switch suffix {
		case "ho": return "Hombre"
		case "ht": return "Hombre Trans"
		case "mu": return "Mujer"
		case "mt": return "Mujer Trans"
		case "ot": return "Otro"
		case "ni": return "No identifica"
		}
	}

	// Nacionalidad (Form 2)
	if strings.HasPrefix(kLower, "victim_case_form2_nationality_") {
		suffix := strings.TrimPrefix(kLower, "victim_case_form2_nationality_")
		switch suffix {
		case "co": return "Colombiana"
		case "ex": return "Extranjera"
		case "ap": return "Apátrida"
		}
	}

	// Sexo asignado al nacer (Form 2)
	if strings.HasPrefix(kLower, "victim_case_form2_assigned_sex_at_birth_") {
		suffix := strings.TrimPrefix(kLower, "victim_case_form2_assigned_sex_at_birth_")
		switch suffix {
		case "mu": return "Mujer"
		case "ho": return "Hombre"
		case "in": return "Intersexual"
		}
	}

	// Nivel de escolaridad (Form 2)
	if strings.HasPrefix(kLower, "victim_case_form2_last_education_level_") {
		suffix := strings.TrimPrefix(kLower, "victim_case_form2_last_education_level_")
		switch suffix {
		case "sf": return "Sin educación formal"
		case "pf": return "Primaria finalizada"
		case "sb": return "Secundaria básica finalizada (9°)"
		case "bm": return "Bachillerato / Media finalizada (11°)"
		case "tt": return "Técnico o tecnológico finalizado"
		case "pr": return "Profesional finalizado"
		case "es": return "Especialización / Maestría finalizada"
		case "do": return "Doctorado finalizado"
		case "ni": return "No informa"
		}
	}

	// Forma de generación de ingresos (Form 2)
	if strings.HasPrefix(kLower, "victim_case_form2_income_generation_method_") {
		suffix := strings.TrimPrefix(kLower, "victim_case_form2_income_generation_method_")
		switch suffix {
		case "em": return "Empleada(o) con contrato formal"
		case "es": return "Empleada(o) sin contrato / informal"
		case "de": return "Desempleada(o)"
		case "pe": return "Pensionada(o) o jubilado(a)"
		case "ta": return "Trabajadora independiente"
		case "si": return "Sin forma de generación de ingresos"
		case "pr": return "Actividades sexuales pagas (ASP)"
		case "no": return "No informa / no responde"
		}
	}

	// Tipo vinculación laboral (Form 2)
	if strings.HasPrefix(kLower, "victim_case_form2_employment_relationship_") {
		suffix := strings.TrimPrefix(kLower, "victim_case_form2_employment_relationship_")
		switch suffix {
		case "pl": return "Planta"
		case "pr": return "Provisional"
		case "co": return "Contratista / OPS"
		case "tr": return "Tercerización"
		case "ca": return "Contrato de aprendizaje"
		case "ps": return "Pasantía"
		}
	}

	// Estrato de vivienda (Form 2)
	if strings.HasPrefix(kLower, "victim_case_form2_housing_stratum_") {
		suffix := strings.TrimPrefix(kLower, "victim_case_form2_housing_stratum_")
		switch suffix {
		case "se": return "Sin estrato"
		case "01": return "1"
		case "02": return "2"
		case "03": return "3"
		case "04": return "4"
		case "05": return "5"
		case "06": return "6"
		case "na": return "No aplica / No sabe"
		}
	}

	// Personas/animales a cargo (Form 2)
	if strings.HasPrefix(kLower, "victim_case_form2_has_dependents_") {
		suffix := strings.TrimPrefix(kLower, "victim_case_form2_has_dependents_")
		switch suffix {
		case "nn": return "No, nadie"
		case "si": return "Sí, hijas/os menores de 18 años"
		case "hi": return "Sí, hijas/os mayores de 18 años"
		case "pa": return "Sí, padre y/o madre"
		case "pr": return "Sí, pareja"
		case "of": return "Sí, otros familiares"
		case "op": return "Sí, otras personas no familiares"
		case "an": return "Sí, animales"
		}
	}

	// Plan de acción (Form 2)
	if strings.HasPrefix(kLower, "victim_case_form2_action_plan_") {
		suffix := strings.TrimPrefix(kLower, "victim_case_form2_action_plan_")
		switch suffix {
		case "ei": return "Enrutamiento interinstitucional"
		case "ar": return "Activación de la ruta (oficio)"
		case "ap": return "Acompañamiento psicosocial"
		case "pe": return "Plan de estabilización"
		case "me": return "Medidas de emergencia"
		case "ea": return "Esquema de autocuidado/orientación en derechos"
		case "se": return "Seguimiento"
		case "rm": return "Remisión a Equipo de Masculinidades"
		case "ga": return "Gestión de Barreras y Alertas"
		}
	}

	// Condición migratoria (Form 2) — usar Locale directamente
	if strings.HasPrefix(kLower, "victim_case_form2_migration_condition_") {
		if locale, ok := salvia_config.Locale["sp"]; ok {
			if translated, found := locale[key]; found {
				return translated
			}
		}
	}

	// Tenencia de vivienda (Form 2) — usar Locale directamente
	if strings.HasPrefix(kLower, "victim_case_form2_housing_tenancy_form_") {
		if locale, ok := salvia_config.Locale["sp"]; ok {
			if translated, found := locale[key]; found {
				return translated
			}
		}
	}

	// Población de especial protección (Form 2) — usar Locale directamente
	if strings.HasPrefix(kLower, "victim_case_form2_specially_protected_population_") {
		if locale, ok := salvia_config.Locale["sp"]; ok {
			if translated, found := locale[key]; found {
				return translated
			}
		}
	}

	// País de nacionalidad (Form 2) — lookup dinámico contra Locale
	if strings.HasPrefix(kLower, "victim_case_form2_specified_nationality_") {
		if locale, ok := salvia_config.Locale["sp"]; ok {
			if translated, found := locale[key]; found {
				return translated
			}
		}
	}

	// Fallback general: buscar en Locale["sp"] para cualquier key no traducida
	if locale, ok := salvia_config.Locale["sp"]; ok {
		if translated, found := locale[key]; found && translated != "" {
			return translated
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
