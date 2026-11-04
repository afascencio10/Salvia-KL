package service

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"

	"bitsflow/internal/repository"
)

// ─────────────────────────────────────────────────────────────────────────────
// Precarga de la Entrevista de Caracterización (FASE 2b-2)
//
// Al abrir la entrevista, el submission del formulario dinámico
// (constants.FormIDCaracterizacion) se siembra con answers copiadas del
// formulario de Registro de Caso del caso. Todo es editable: son answers
// normales de un form_submission.
//
// La tabla caracterizacionPrefillRules es la única fuente de verdad del mapeo
// origen -> destino. Para ajustar la precarga basta con editar esa tabla.
// ─────────────────────────────────────────────────────────────────────────────

// caracterizacionFormID es el form_id del formulario de caracterización
// (debe coincidir con seed.CaracterizacionFormID / constants.FormIDCaracterizacion).
const caracterizacionFormID = "f1a2c3d4-5e6f-4a1b-8c2d-3e4f5a6b7c8d"

// caracterizacionQuestionPrefix es el prefijo común de los question_id sembrados;
// cada regla solo guarda el sufijo de 12 caracteres.
const caracterizacionQuestionPrefix = "c0000001-0000-4000-8000-"

// Preguntas del repeater RG-3 "Medida solicitada" (seed caracterizacion_form.sql).
const (
	caracterizacionQMedida        = "c0000001-0000-4000-8000-000000000921"
	caracterizacionQMedidaSolicit = "c0000001-0000-4000-8000-000000000922"
	caracterizacionQMedidaTiempo  = "c0000001-0000-4000-8000-000000000923"
	caracterizacionQMedidaObs     = "c0000001-0000-4000-8000-000000000924"
)

// Categorías de enum (sin el prefijo "enums.") usadas por las derivaciones.
const (
	enumCatEthnic       = "victim_case_form2_ethnic_affiliation"
	enumCatNationality  = "victim_case_form2_nationality"
	enumCatHasDependent = "victim_case_form2_has_dependents"
)

// Códigos (columna victim_case_form2_enums_code) verificados en Locale.go.
const (
	enumCodeEthnicIndigenous = "in" // ethnic_affiliation_in = Indígena
	enumCodeNationalityForgn = "ex" // nationality_ex = Extranjera
	enumCodeDependentsNobody = "nn" // has_dependents_nn = No, nadie
)

type precargaKind int

const (
	// kindCopy copia el texto tal cual si no está vacío.
	kindCopy precargaKind = iota
	// kindDate copia una fecha válida en formato YYYY-MM-DD.
	kindDate
	// kindAge calcula los años cumplidos desde una fecha de nacimiento.
	kindAge
	// kindNumber copia un número válido.
	kindNumber
	// kindOptionInline valida el valor contra las option inline de la pregunta destino.
	kindOptionInline
	// kindEnum copia un i_code validado contra la categoría del enum destino.
	kindEnum
	// kindEnumMulti copia una lista CSV de i_codes, conservando solo los válidos.
	kindEnumMulti
	// kindBool copia un booleano ("true"/"false").
	kindBool
	// kindBoolFromEnumCode deriva true/false: el código del enum origen == TargetCode.
	kindBoolFromEnumCode
	// kindDependentsBool deriva true si algún código origen != TargetCode ("nn"), false si solo "nn".
	kindDependentsBool
	// kindDependentsWho copia el primer i_code origen cuyo código != TargetCode.
	kindDependentsWho
)

// caracterizacionPrefillRule describe una fila de la tabla de equivalencias.
type caracterizacionPrefillRule struct {
	Dest           string // sufijo (12 chars) del question_id destino
	Label          string // etiqueta legible del destino (para el reporte/logs)
	Source         string // FieldKey del Registro de Caso ("S6Q1")
	Kind           precargaKind
	SourceCategory string // categoría de enum del ORIGEN (solo derivaciones)
	TargetCode     string // código de enum que dispara/excluye la derivación
	Requires       string // sufijo destino que debe haberse precargado con "true"
}

// caracterizacionPrefillRules es la tabla de equivalencias origen -> destino.
// No se precargan (sin fuente cierta): I.7 Régimen, I.8 EPS, I.13 Corregimiento/Vereda,
// fecha de entrevista, "Otra ¿cuál?", textos de ayuda, repeaters y secciones II.1/II.5 a VIII.
var caracterizacionPrefillRules = []caracterizacionPrefillRule{
	{Dest: "000000000101", Label: "I.1 Nombres", Source: "S1Q1", Kind: kindCopy},
	{Dest: "000000000102", Label: "I.2 Apellidos", Source: "S1Q2", Kind: kindCopy},
	{Dest: "000000000103", Label: "I.3 Nombre identitario", Source: "S1Q3", Kind: kindCopy},
	{Dest: "000000000104", Label: "I.4 Tipo de documento", Source: "S1Q5", Kind: kindOptionInline},
	{Dest: "000000000105", Label: "I.5 Fecha de nacimiento", Source: "S6Q1", Kind: kindDate},
	{Dest: "000000000106", Label: "I.6 Edad (calculada)", Source: "S6Q1", Kind: kindAge},
	{Dest: "000000000109", Label: "I.9 Dirección", Source: "S1Q11", Kind: kindCopy},
	{Dest: "000000000110", Label: "I.10 Zona", Source: "S1Q7", Kind: kindEnum},
	{Dest: "000000000111", Label: "I.11 Departamento", Source: "S1Q8", Kind: kindCopy},
	{Dest: "000000000112", Label: "I.12 Ciudad", Source: "S1Q9", Kind: kindCopy},
	{Dest: "000000000114", Label: "I.14 Teléfono", Source: "S1Q4", Kind: kindNumber},
	{Dest: "000000000115", Label: "I.15 Contacto de apoyo", Source: "S2Q1", Kind: kindCopy},
	{Dest: "000000000116", Label: "I.16 Teléfono del contacto", Source: "S2Q2", Kind: kindNumber},
	{Dest: "000000000117", Label: "I.17 Estado civil", Source: "S6Q23", Kind: kindEnum},
	{Dest: "000000000118", Label: "I.18 Sexo", Source: "S6Q18", Kind: kindEnum},
	{Dest: "000000000119", Label: "I.19 Identidad de género", Source: "S6Q16", Kind: kindEnum},
	{Dest: "000000000120", Label: "I.20 Orientación sexual", Source: "S6Q17", Kind: kindEnum},
	{Dest: "000000000121", Label: "I.21 Pertenencia étnica", Source: "S6Q20", Kind: kindEnum},
	{Dest: "000000000122", Label: "I.22 ¿Pueblo indígena? (derivado)", Source: "S6Q20", Kind: kindBoolFromEnumCode, SourceCategory: enumCatEthnic, TargetCode: enumCodeEthnicIndigenous},
	{Dest: "00000000012b", Label: "I.22b ¿Cuál pueblo?", Source: "S6Q21", Kind: kindEnum, Requires: "000000000122"},
	{Dest: "000000000123", Label: "I.23 ¿Migrante? (derivado)", Source: "S6Q13", Kind: kindBoolFromEnumCode, SourceCategory: enumCatNationality, TargetCode: enumCodeNationalityForgn},
	{Dest: "000000000124", Label: "I.24 Situación migratoria", Source: "S6Q15", Kind: kindEnum, Requires: "000000000123"},
	{Dest: "000000000125", Label: "I.25 Nivel educativo", Source: "S6Q24", Kind: kindEnum},
	{Dest: "000000000126", Label: "I.26 Población especial", Source: "S6Q19", Kind: kindEnumMulti},
	{Dest: "000000000127", Label: "I.27 ¿Discapacidad?", Source: "S6Q2", Kind: kindBool},
	{Dest: "000000000128", Label: "I.28 Ajustes razonables", Source: "S1Q14", Kind: kindEnumMulti},
	{Dest: "000000000129", Label: "I.29 Ocupación", Source: "S6Q25", Kind: kindEnum},
	{Dest: "000000000130", Label: "I.30 Ingresos", Source: "S6Q26", Kind: kindEnum},
	{Dest: "000000000131", Label: "I.31 Hechos victimizantes", Source: "S3Q1", Kind: kindCopy},
	{Dest: "000000000203", Label: "II.3 ¿Personas a cargo? (derivado)", Source: "S6Q33", Kind: kindDependentsBool, SourceCategory: enumCatHasDependent, TargetCode: enumCodeDependentsNobody},
	{Dest: "000000000204", Label: "II.4 ¿Quiénes?", Source: "S6Q33", Kind: kindDependentsWho, SourceCategory: enumCatHasDependent, TargetCode: enumCodeDependentsNobody, Requires: "000000000203"},
}

// caracterizacionPregunta es lo que se necesita saber de una pregunta destino
// para validar el valor a precargar.
type caracterizacionPregunta struct {
	Type         string
	EnumCategory string              // categoría (sin "enums.") si state_options_path es enums.*
	Options      map[string]struct{} // valores de las option inline
}

// caracterizacionCatalogos agrupa lo que se lee de la BD para validar.
type caracterizacionCatalogos struct {
	Preguntas map[string]caracterizacionPregunta // por question_id destino completo
	// Enums[categoria][i_code] = code
	Enums map[string]map[string]string
}

// caracterizacionOmitida describe una regla que no se pudo precargar.
type caracterizacionOmitida struct {
	Label  string
	Motivo string
}

var reFechaISO = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}`)

// calcularEdadCaracterizacion devuelve los años cumplidos a la fecha hoy.
// ok=false si la fecha es inválida, futura, o la edad supera 120 años.
func calcularEdadCaracterizacion(nacimiento string, hoy time.Time) (int, bool) {
	nac, ok := normalizarFechaPrecarga(nacimiento)
	if !ok {
		return 0, false
	}
	t, err := time.Parse("2006-01-02", nac)
	if err != nil {
		return 0, false
	}
	if t.After(hoy) {
		return 0, false
	}
	edad := hoy.Year() - t.Year()
	if hoy.Month() < t.Month() || (hoy.Month() == t.Month() && hoy.Day() < t.Day()) {
		edad--
	}
	if edad < 0 || edad > 120 {
		return 0, false
	}
	return edad, true
}

// normalizarFechaPrecarga devuelve la fecha como YYYY-MM-DD si es válida.
func normalizarFechaPrecarga(v string) (string, bool) {
	v = strings.TrimSpace(v)
	if !reFechaISO.MatchString(v) {
		return "", false
	}
	d := v[:10]
	if _, err := time.Parse("2006-01-02", d); err != nil {
		return "", false
	}
	return d, true
}

func splitCSVPrecarga(v string) []string {
	out := []string{}
	for _, p := range strings.Split(v, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

// enumValido indica si el i_code existe en la categoría.
func (c *caracterizacionCatalogos) enumValido(categoria, icode string) bool {
	m, ok := c.Enums[categoria]
	if !ok {
		return false
	}
	_, ok = m[icode]
	return ok
}

// enumCode devuelve el code del i_code dentro de la categoría ("" si no existe).
func (c *caracterizacionCatalogos) enumCode(categoria, icode string) string {
	return c.Enums[categoria][icode]
}

// validarValorPrecarga aplica la regla de un solo campo. Devuelve el valor a
// insertar; ok=false con motivo si hay que omitirlo (nunca se inserta basura).
func validarValorPrecarga(rule caracterizacionPrefillRule, rc map[string]string, cat *caracterizacionCatalogos, hoy time.Time) (valor string, ok bool, motivo string) {
	raw := strings.TrimSpace(rc[rule.Source])
	qID := caracterizacionQuestionPrefix + rule.Dest
	preg, hayPregunta := cat.Preguntas[qID]
	if !hayPregunta {
		return "", false, "la pregunta destino no existe en el formulario"
	}
	if raw == "" {
		return "", false, "sin dato en el Registro de Caso"
	}
	switch rule.Kind {
	case kindCopy:
		return raw, true, ""
	case kindDate:
		d, valid := normalizarFechaPrecarga(raw)
		if !valid {
			return "", false, "fecha inválida"
		}
		return d, true, ""
	case kindAge:
		edad, valid := calcularEdadCaracterizacion(raw, hoy)
		if !valid {
			return "", false, "fecha de nacimiento inválida, futura o edad fuera de rango"
		}
		return strconv.Itoa(edad), true, ""
	case kindNumber:
		if _, err := strconv.ParseFloat(raw, 64); err != nil {
			return "", false, "no es un número válido"
		}
		return raw, true, ""
	case kindOptionInline:
		if _, found := preg.Options[raw]; !found {
			return "", false, "valor no existe entre las opciones de la pregunta destino"
		}
		return raw, true, ""
	case kindEnum:
		if preg.EnumCategory == "" || !cat.enumValido(preg.EnumCategory, raw) {
			return "", false, "valor no pertenece a la categoría de enum destino"
		}
		return raw, true, ""
	case kindEnumMulti:
		if preg.EnumCategory == "" {
			return "", false, "la pregunta destino no usa un enum"
		}
		kept := []string{}
		for _, v := range splitCSVPrecarga(raw) {
			if cat.enumValido(preg.EnumCategory, v) {
				kept = append(kept, v)
			}
		}
		if len(kept) == 0 {
			return "", false, "ningún valor es válido para la categoría destino"
		}
		return strings.Join(kept, ","), true, ""
	case kindBool:
		switch strings.ToLower(raw) {
		case "true":
			return "true", true, ""
		case "false":
			return "false", true, ""
		}
		return "", false, "no es booleano"
	case kindBoolFromEnumCode:
		if !cat.enumValido(rule.SourceCategory, raw) {
			return "", false, "valor origen no pertenece a su categoría"
		}
		if cat.enumCode(rule.SourceCategory, raw) == rule.TargetCode {
			return "true", true, ""
		}
		return "false", true, ""
	case kindDependentsBool, kindDependentsWho:
		vals := []string{}
		for _, v := range splitCSVPrecarga(raw) {
			if cat.enumValido(rule.SourceCategory, v) {
				vals = append(vals, v)
			}
		}
		if len(vals) == 0 {
			return "", false, "ningún valor origen es válido"
		}
		primero := ""
		for _, v := range vals {
			if cat.enumCode(rule.SourceCategory, v) != rule.TargetCode {
				primero = v
				break
			}
		}
		if rule.Kind == kindDependentsBool {
			if primero != "" {
				return "true", true, ""
			}
			return "false", true, ""
		}
		if primero == "" {
			return "", false, "no tiene personas a cargo"
		}
		if preg.EnumCategory != "" && !cat.enumValido(preg.EnumCategory, primero) {
			return "", false, "valor no pertenece a la categoría de enum destino"
		}
		return primero, true, ""
	}
	return "", false, "tipo de regla desconocido"
}

// construirPrecargaCaracterizacion aplica la tabla de reglas sobre las
// respuestas del Registro de Caso. Devuelve las answers por question_id destino
// y la lista de reglas omitidas con su motivo. Función pura (testeable).
func construirPrecargaCaracterizacion(rc map[string]string, cat *caracterizacionCatalogos, hoy time.Time) (map[string]string, []caracterizacionOmitida) {
	respuestas := map[string]string{}
	omitidas := []caracterizacionOmitida{}
	if cat == nil {
		return respuestas, omitidas
	}
	for _, rule := range caracterizacionPrefillRules {
		if rule.Requires != "" {
			if respuestas[caracterizacionQuestionPrefix+rule.Requires] != "true" {
				omitidas = append(omitidas, caracterizacionOmitida{Label: rule.Label, Motivo: "depende de una pregunta que no quedó en Sí"})
				continue
			}
		}
		valor, ok, motivo := validarValorPrecarga(rule, rc, cat, hoy)
		if !ok {
			omitidas = append(omitidas, caracterizacionOmitida{Label: rule.Label, Motivo: motivo})
			continue
		}
		respuestas[caracterizacionQuestionPrefix+rule.Dest] = valor
	}
	return respuestas, omitidas
}

// cargarCatalogosPrecarga lee las preguntas destino (tipo, ruta de opciones,
// option inline) y el catálogo de enums de victim_case_form2_enums.
func cargarCatalogosPrecarga(ctx context.Context, db *gorm.DB) (*caracterizacionCatalogos, error) {
	cat := &caracterizacionCatalogos{
		Preguntas: map[string]caracterizacionPregunta{},
		Enums:     map[string]map[string]string{},
	}

	type qRow struct {
		ID   string `gorm:"column:id"`
		Type string `gorm:"column:qtype"`
		Path string `gorm:"column:path"`
	}
	var qs []qRow
	if err := db.WithContext(ctx).Raw(`
		SELECT id::text AS id, question_type AS qtype, COALESCE(state_options_path, '') AS path
		  FROM salvia.question
		 WHERE form_id::text = ?
		   AND deleted_at IS NULL
	`, caracterizacionFormID).Scan(&qs).Error; err != nil {
		return nil, fmt.Errorf("precarga: leer preguntas: %w", err)
	}
	for _, q := range qs {
		p := caracterizacionPregunta{Type: q.Type, Options: map[string]struct{}{}}
		if strings.HasPrefix(q.Path, "enums.") {
			p.EnumCategory = strings.TrimPrefix(q.Path, "enums.")
		}
		cat.Preguntas[q.ID] = p
	}

	type oRow struct {
		QID   string `gorm:"column:qid"`
		Value string `gorm:"column:value"`
	}
	var os []oRow
	if err := db.WithContext(ctx).Raw(`
		SELECT o.question_id::text AS qid, o.value AS value
		  FROM salvia.option o
		  JOIN salvia.question q ON q.id::text = o.question_id::text
		 WHERE q.form_id::text = ?
		   AND o.deleted_at IS NULL
	`, caracterizacionFormID).Scan(&os).Error; err != nil {
		return nil, fmt.Errorf("precarga: leer opciones: %w", err)
	}
	for _, o := range os {
		if p, ok := cat.Preguntas[o.QID]; ok {
			p.Options[o.Value] = struct{}{}
		}
	}

	type eRow struct {
		Category string `gorm:"column:category"`
		ICode    string `gorm:"column:icode"`
		Code     string `gorm:"column:code"`
	}
	var es []eRow
	if err := db.WithContext(ctx).Raw(`
		SELECT victim_case_form2_enums_category AS category,
		       victim_case_form2_enums_i_code::text AS icode,
		       victim_case_form2_enums_code AS code
		  FROM salvia.victim_case_form2_enums
	`).Scan(&es).Error; err != nil {
		return nil, fmt.Errorf("precarga: leer enums: %w", err)
	}
	for _, e := range es {
		m, ok := cat.Enums[e.Category]
		if !ok {
			m = map[string]string{}
			cat.Enums[e.Category] = m
		}
		m[e.ICode] = e.Code
	}
	return cat, nil
}

// respuestasRegistroCaso devuelve las answers del Registro de Caso del caso,
// indexadas por FieldKey ("S1Q1"). Vía de enlace: salvia.victim_case.
// victim_case_form_submission_id (victim_case_i_code == case_task.case_id).
// Casos legacy sin submission: fallback mínimo con Nombres/Apellidos de las
// columnas de victim_case (omitiendo el valor dummy "Pendiente").
func respuestasRegistroCaso(ctx context.Context, db *gorm.DB, caseID string) (map[string]string, error) {
	out := map[string]string{}
	caseID = strings.TrimSpace(caseID)
	if caseID == "" {
		return out, nil
	}

	var vc struct {
		SubmissionID string `gorm:"column:submission_id"`
		Names        string `gorm:"column:names"`
		LastNames    string `gorm:"column:last_names"`
	}
	if err := db.WithContext(ctx).Raw(`
		SELECT COALESCE(victim_case_form_submission_id::text, '') AS submission_id,
		       COALESCE(victim_case_victim_names, '') AS names,
		       COALESCE(victim_case_victim_last_names, '') AS last_names
		  FROM salvia.victim_case
		 WHERE victim_case_i_code = ?
		 LIMIT 1
	`, caseID).Scan(&vc).Error; err != nil {
		return out, fmt.Errorf("precarga: leer victim_case: %w", err)
	}

	if strings.TrimSpace(vc.SubmissionID) != "" {
		type row struct {
			SectionOrder  int
			QuestionOrder int
			Value         string
		}
		var rows []row
		// Misma consulta que BuildAnswersByFieldKey (victim_case_form_repository.go).
		if err := db.WithContext(ctx).Raw(`
			SELECT fs."order" AS section_order, q."order" AS question_order, a.value AS value
			FROM salvia.answer a
			JOIN salvia.question q ON q.id = a.question_id
			JOIN salvia.form_section fs ON fs.id::varchar = q.form_section_id::varchar
			WHERE a.form_submission_id = ? AND q.form_id = ?
		`, vc.SubmissionID, RegistroCasoFormID).Scan(&rows).Error; err != nil {
			return out, fmt.Errorf("precarga: leer respuestas del Registro de Caso: %w", err)
		}
		for _, r := range rows {
			out[repository.FieldKey(r.SectionOrder, r.QuestionOrder)] = r.Value
		}
		return out, nil
	}

	// Legacy: sin submission de Registro de Caso.
	if n := strings.TrimSpace(vc.Names); n != "" && n != dummyNames {
		out["S1Q1"] = n
	}
	if n := strings.TrimSpace(vc.LastNames); n != "" && n != dummyLastNames {
		out["S1Q2"] = n
	}
	return out, nil
}
