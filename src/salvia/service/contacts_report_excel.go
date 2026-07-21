package service

import (
	"bitsflow/internal/repository"
	common_config "bitsflow/common/config"
	salvia_config "bitsflow/salvia/config"
	"fmt"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

const contactsSheetName = "Reportes"

// contactsHeaders define las 34 columnas de la hoja única del reporte consolidado
// de contactos, en el orden exigido por el spec del servicio.
var contactsHeaders = []string{
	"contacto_id",
	"contacto_icode",
	"contacto_fecha_creacion",
	"contacto_fecha_actualizacion",
	"contacto_estado",
	"contacto_nombres",
	"contacto_apellidos",
	"contacto_latitud",
	"contacto_longitud",
	"tipo_formulario_reporte",
	"f1_apodo_nombre_identitario",
	"f1_tipo_documento",
	"f1_numero_documento",
	"f1_fecha_nacimiento",
	"f1_codigo_municipio",
	"f1_municipio",
	"f1_direccion",
	"f1_telefono",
	"f1_identidad_genero",
	"f1_orientacion_sexual",
	"f1_procedencia_zona",
	"f1_ocupacion",
	"f1_otra_ocupacion",
	"f1_descripcion_hechos",
	"f2_nombre_reportante",
	"f2_telefono_reportante",
	"f2_telefono_contacto_victima",
	"f2_descripcion_hechos",
	"f2_mejor_hora_contacto",
	"f2_recibira_llamada",
	"f2_rol_cuidado",
	"f2_victima_enterada",
	"f2_tipo_reporte",
	"f2_ajustes_gbv",
}

// translateCode resuelve un código corto de formulario 1 (ej. "cc", "fe") a su
// etiqueta en español usando los catálogos de configuración; si el código no
// existe en el catálogo se exporta tal cual (defensa) y vacío queda vacío.
func translateCode(catalog map[string]string, code string) string {
	if code == "" {
		return ""
	}
	if label, found := catalog[code]; found {
		return label
	}
	return code
}

// translateForm2Enum resuelve la clave de un enum de Form2 (ej. "yes_no_y",
// "victim_case_form2_report_type_v") a su texto en español vía Locale; si la
// clave no existe en Locale cae al patrón translateEnumKey y, en última
// instancia, exporta la clave cruda (defensa).
func translateForm2Enum(key string) string {
	key = strings.TrimSpace(key)
	if key == "" {
		return ""
	}
	if label, found := salvia_config.Locale["sp"][key]; found {
		return label
	}
	return translateEnumKey(key)
}

// translateForm2EnumList traduce una lista de claves separadas por coma
// (f2_ajustes_gbv) manteniendo el separador ", ".
func translateForm2EnumList(list string) string {
	if strings.TrimSpace(list) == "" {
		return ""
	}
	parts := strings.Split(list, ", ")
	for i, p := range parts {
		parts[i] = translateForm2Enum(p)
	}
	return strings.Join(parts, ", ")
}

func formatContactFormType(formType string) string {
	switch formType {
	case "f1":
		return "Reporte propio (Formulario 1)"
	case "f2":
		return "Reporte de tercero (Formulario 2)"
	default:
		return ""
	}
}

func formatContactDateTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("02/01/2006 15:04")
}

// BuildContactsExcel genera el Excel de una sola hoja plana con una fila por
// reporte (victim_contact), incluyendo la información de Form1 o Form2.
func BuildContactsExcel(contacts []repository.ContactReportDTO) (*excelize.File, error) {
	f := excelize.NewFile()
	f.SetSheetName("Sheet1", contactsSheetName)

	headerStyle, err := createHeaderStyle(f)
	if err != nil {
		f.Close()
		return nil, err
	}
	if err := writeHeaders(f, contactsSheetName, contactsHeaders, headerStyle); err != nil {
		f.Close()
		return nil, err
	}

	for i, c := range contacts {
		rNum := i + 2

		values := []interface{}{
			c.VictimContactId,
			c.VictimContactICode,
			formatContactDateTime(c.CreationDate),
			formatContactDateTime(c.UpdateDate),
			c.StatusDescription,
			c.Names,
			c.LastNames,
			contactCoordinate(c.Latitude),
			contactCoordinate(c.Longitude),
			formatContactFormType(c.FormType),
			c.Form1Nick,
			translateCode(common_config.DOCUMENT_TYPE, c.Form1DocType),
			c.Form1DocNumber,
			c.Form1BirthDate,
			c.Form1TownCode,
			c.Form1TownName,
			c.Form1Address,
			c.Form1Phone,
			translateCode(common_config.GENDER_IDENTITY, c.Form1GenderIdentity),
			translateCode(salvia_config.SEXUAL_ORIENTATION, c.Form1SexualOrientation),
			translateCode(salvia_config.ORIGIN_PLACE, c.Form1Origin),
			translateCode(salvia_config.OCCUPATION, c.Form1Occupation),
			c.Form1OccupationOther,
			c.Form1FactsDescription,
			c.Form2ReporterNames,
			c.Form2ReporterPhone,
			c.Form2VictimColPhone,
			c.Form2FactsDescription,
			c.Form2BestContactTime,
			translateForm2Enum(c.Form2WillReceiveCall),
			translateForm2Enum(c.Form2HasCareRole),
			translateForm2Enum(c.Form2VictimAware),
			translateForm2Enum(c.Form2ReportType),
			translateForm2EnumList(c.Form2AdjustmentsGBV),
		}

		for col, val := range values {
			colName, err := excelize.ColumnNumberToName(col + 1)
			if err != nil {
				f.Close()
				return nil, err
			}
			cell := fmt.Sprintf("%s%d", colName, rNum)
			if s, ok := val.(string); ok {
				// Los strings (teléfonos, documentos, códigos) se escriben como texto
				// para evitar notación científica o pérdida de ceros a la izquierda.
				_ = f.SetCellStr(contactsSheetName, cell, s)
			} else {
				_ = f.SetCellValue(contactsSheetName, cell, val)
			}
		}
	}

	autoFitColumns(f)
	return f, nil
}

// contactCoordinate deja la celda vacía cuando la coordenada no fue registrada (0).
func contactCoordinate(v float64) interface{} {
	if v == 0 {
		return ""
	}
	return v
}
