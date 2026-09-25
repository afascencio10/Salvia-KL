package salvia_config

// Catálogo de competencias (acciones institucionales) en VBG por entidad.
//
// Origen: hoja "4.EntidadesAcciones" del consolidado de rutas que entregó
// Coordinación de Rutas Salvia — 38 entidades, 237 competencias, cada una con su
// tiempo de respuesta esperado y su sustento normativo.
//
// El catálogo va embebido y no en una tabla a propósito: es normativa, cambia
// cuando cambia la ley y no cuando opera el sistema. Para actualizarlo se
// regenera competences.json con scripts/extract_competences.py y se versiona el
// resultado.
//
// La competencia pertenece a la ENTIDAD (la institución), nunca a una sede: el
// Excel lista tipos de institución ("IPS", "Comisarías de Familia"), no sedes
// puntuales, así que una obligación legal no puede quedar atada a una dirección.

import (
	"embed"
	"encoding/json"
	"fmt"
)

//go:embed competences.json
var competencesFS embed.FS

// Competence es una acción institucional que una entidad está obligada a cumplir.
type Competence struct {
	Code         string `json:"code"`
	Action       string `json:"action"`
	Scope        string `json:"scope"`
	Summary      string `json:"summary"`
	ResponseTime string `json:"responseTime"`
	LegalBasis   string `json:"legalBasis"`
}

// CompetenceKind es un tipo de entidad del catálogo con sus competencias.
type CompetenceKind struct {
	Code        string       `json:"code"`
	Sector      string       `json:"sector"`
	Name        string       `json:"name"`
	Competences []Competence `json:"competences"`
}

// CompetenceOption es el formato { label, value } que consume stateOptionsPath en
// dinamic-form (igual que models.LocationOption para departamentos y ciudades).
type CompetenceOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// CompetenceKinds son las entidades del catálogo en el orden del Excel.
var CompetenceKinds = loadCompetenceKinds()

// Sectores del formulario de seguimiento (valores de la pregunta Q1 del repeater
// de barreras). No coinciden con los sectores del Excel de competencias ni con
// los del Directorio de Entidades: son tres catálogos distintos y este mapa es el
// puente, para no tocar ninguno de los dos.
const (
	FormSectorSalud              = "salud"
	FormSectorJusticia           = "justicia"
	FormSectorProteccion         = "proteccion"
	FormSectorOtrasInstituciones = "otras_instituciones"
)

// formSectorByCatalogSector traduce el sector del Excel al sector del formulario.
// Los siete sectores que el formulario no tiene (entes territoriales, interior,
// educación, trabajo, defensa, relaciones exteriores y agricultura) caen en
// "otras instituciones", que es justamente la opción que el formulario ofrece
// para lo que no es salud, justicia ni protección.
var formSectorByCatalogSector = map[string]string{
	"SALUD":                 FormSectorSalud,
	"JUSTICIA":              FormSectorJusticia,
	"PROTECCIÓN":            FormSectorProteccion,
	"ENTES TERRITORIALES":   FormSectorOtrasInstituciones,
	"INTERIOR":              FormSectorOtrasInstituciones,
	"EDUCACIÓN":             FormSectorOtrasInstituciones,
	"TRABAJO":               FormSectorOtrasInstituciones,
	"DEFENSA":               FormSectorOtrasInstituciones,
	"RELACIONES EXTERIORES": FormSectorOtrasInstituciones,
	"AGRICULTURA":           FormSectorOtrasInstituciones,
}

// competenceKindsByInstitution traduce el valor de la pregunta "¿A qué institución
// acudió?" del repeater de barreras (Q3 Salud, Q6 Justicia, Q9 Protección) a las
// entidades del catálogo que le corresponden.
//
// Es uno a varios porque el formulario pregunta por la puerta de entrada y el
// catálogo por la institución obligada:
//   - hospital, puesto de salud y consultorio son todos IPS, y el Excel trae IPS
//     escrita de dos formas distintas, así que suma las dos;
//   - CAIVAS, CAVIF y la URI son dependencias de la Fiscalía;
//   - "policía judicial" cubre CTI, SIJIN y DIJIN.
//
// Un valor sin entrada aquí (hoy solo "no sabe / no responde") no es un error: se
// resuelve por sector en CompetenceOptions.
var competenceKindsByInstitution = map[string][]string{
	// Q3 — Salud
	"hospital_publico":   {"instituciones_prestadoras_de_servicios_d_5694", "instituciones_prestadoras_de_servicios_d_357c"},
	"hospital_privado":   {"instituciones_prestadoras_de_servicios_d_5694", "instituciones_prestadoras_de_servicios_d_357c"},
	"ips":                {"instituciones_prestadoras_de_servicios_d_5694", "instituciones_prestadoras_de_servicios_d_357c"},
	"puesto_salud":       {"instituciones_prestadoras_de_servicios_d_5694", "instituciones_prestadoras_de_servicios_d_357c"},
	"consultorio_medico": {"instituciones_prestadoras_de_servicios_d_5694", "instituciones_prestadoras_de_servicios_d_357c"},
	"eps":                {"empresas_administradoras_de_planes_de_be_e2ed"},

	// Q6 — Justicia
	"fiscalia":                  {"fiscalia_general_de_la_nacion_fgn_6eb1"},
	"caivas":                    {"fiscalia_general_de_la_nacion_fgn_6eb1"},
	"caviv":                     {"fiscalia_general_de_la_nacion_fgn_6eb1"},
	"unidad_reaccion_inmediata": {"fiscalia_general_de_la_nacion_fgn_6eb1"},
	"medicina_legal":            {"instituto_nacional_de_medicina_legal_y_c_4ad6"},
	"policia_judicial": {
		"policia_judicial_cti_cuerpo_tecnico_de_i_0091",
		"policia_judicial_sijin_seccional_de_inve_62c3",
		"policia_judicial_dijin_direccion_de_inve_856a",
	},
	"jueces_civiles":   {"jueces_de_control_de_garantias_civil_mun_20e4"},
	"jueces_penales":   {"jueces_de_control_de_garantias_civil_mun_20e4"},
	"casas_justicia":   {"ministerio_de_justicia_y_del_derecho_916b"},
	"cai_policia":      {"policia_nacional_de_colombia_4ec0"},
	"estacion_policia": {"policia_nacional_de_colombia_4ec0"},

	// Q9 — Protección
	"comisaria_familia":          {"comisarias_de_familia_b249"},
	"inspeccion_policia":         {"inspecciones_de_policia_e4bb"},
	"icbf":                       {"instituto_colombiano_de_bienestar_famili_0e5a"},
	"defensoria_familia":         {"defensorias_de_familia_icbf_90d0"},
	"unidad_nacional_proteccion": {"unidad_nacional_de_proteccion_e984"},
	"fiscalia_proteccion":        {"fiscalia_general_de_la_nacion_fgn_6eb1"},
	"ejercito_nacional":          {"ministerio_de_defensa_809e"},
	"armada_colombiana":          {"ministerio_de_defensa_809e"},
}

// CompetenceOptions devuelve las competencias a mostrar como opciones en el
// formulario de seguimiento, en el orden del catálogo y sin repetidas.
//
// institutions son los valores marcados en la pregunta de institución de la
// entrada de barrera. Si ninguno resuelve a una entidad del catálogo (el agente
// eligió "no sabe", o todavía no marcó institución), se devuelven las
// competencias de todo el sector: la idea es que la pregunta nunca quede sin
// opciones cuando el sector sí se conoce.
//
// Un sector que el catálogo no cubre — "barrera_transversal", que en el
// formulario ni siquiera tiene pregunta de institución — devuelve vacío.
func CompetenceOptions(formSector string, institutions []string) []CompetenceOption {
	selected := make(map[string]bool)
	for _, institution := range institutions {
		for _, kindCode := range competenceKindsByInstitution[institution] {
			selected[kindCode] = true
		}
	}

	options := make([]CompetenceOption, 0, 16)
	seen := make(map[string]bool)
	for _, kind := range CompetenceKinds {
		if len(selected) > 0 {
			if !selected[kind.Code] {
				continue
			}
		} else if formSector == "" || formSectorByCatalogSector[kind.Sector] != formSector {
			continue
		}
		for _, competence := range kind.Competences {
			if seen[competence.Code] {
				continue
			}
			seen[competence.Code] = true
			options = append(options, CompetenceOption{Label: competence.Action, Value: competence.Code})
		}
	}
	return options
}

// loadCompetenceKinds parsea el catálogo embebido. Un JSON inválido o incompleto
// es un error de build, no de operación: revienta al arrancar en vez de dejar el
// formulario sin opciones en producción.
func loadCompetenceKinds() []CompetenceKind {
	raw, err := competencesFS.ReadFile("competences.json")
	if err != nil {
		panic(fmt.Sprintf("salvia_config: leer competences.json: %v", err))
	}
	var catalog struct {
		Kinds []CompetenceKind `json:"kinds"`
	}
	if err := json.Unmarshal(raw, &catalog); err != nil {
		panic(fmt.Sprintf("salvia_config: competences.json inválido: %v", err))
	}
	if len(catalog.Kinds) == 0 {
		panic("salvia_config: competences.json sin entidades")
	}
	return catalog.Kinds
}
