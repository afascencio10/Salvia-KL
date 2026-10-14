package service

import (
	salvia_config "bitsflow/salvia/config"
	"strings"
)

// Etiquetas legibles para los enums de un reporte. La BD guarda el NOMBRE del enum (clave de
// traducción, p. ej. "victim_case_form2_report_type_v2_tercero"); la API de Gestión de reportes
// debe devolver texto para mostrar, no la clave.

// reportTypeLabels — tipo de reporte (V1 y V2) → etiqueta corta de la columna "Tipo" y del badge del modal.
// Los nombres V2 no tienen una etiqueta corta en Locale (allí están las frases largas del wizard).
var reportTypeLabels = map[string]string{
	"victim_case_form2_report_type_v":           "Víctima",
	"victim_case_form2_report_type_r":           "Tercero",
	"victim_case_form2_report_type_v2_victima":  "Víctima",
	"victim_case_form2_report_type_v2_tercero":  "Tercero",
	"victim_case_form2_report_type_v2_familiar": "Familiar",
	"victim_case_form2_report_type_v2_hombre":   "Hombre",
}

// reportTypeLabel devuelve la etiqueta corta, o "" si el nombre no es un tipo de reporte conocido.
func reportTypeLabel(enumName string) string {
	return reportTypeLabels[strings.TrimSpace(enumName)]
}

// localizeEnumName traduce una clave de enum con Locale["sp"]; si no hay traducción devuelve la clave tal cual.
func localizeEnumName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return ""
	}
	if label, ok := salvia_config.Locale["sp"][name]; ok && strings.TrimSpace(label) != "" {
		return label
	}
	return name
}

// localizeEnumNames traduce una lista separada por ", " (p. ej. los ajustes razonables, string_agg en BD).
func localizeEnumNames(csv string) string {
	if strings.TrimSpace(csv) == "" {
		return ""
	}
	parts := strings.Split(csv, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if l := localizeEnumName(p); l != "" {
			out = append(out, l)
		}
	}
	return strings.Join(out, ", ")
}
