package salvia_config

import "testing"

// Valores de las preguntas "¿A qué institución acudió?" del repeater de barreras,
// tal como están registradas en la tabla option — ver
// DocsMD/Screens/hacer-seguimiento/form-barreras-repeater.md (Q3, Q6, Q9).
var institutionsBySector = map[string][]string{
	FormSectorSalud: {
		"hospital_publico", "hospital_privado", "ips", "eps",
		"puesto_salud", "consultorio_medico", "no_sabe",
	},
	FormSectorJusticia: {
		"cai_policia", "caivas", "casas_justicia", "caviv", "comisaria_familia",
		"estacion_policia", "fiscalia", "inspeccion_policia", "medicina_legal",
		"jueces_civiles", "jueces_penales", "policia_judicial",
		"unidad_reaccion_inmediata", "no_sabe",
	},
	FormSectorProteccion: {
		"comisaria_familia", "fiscalia_proteccion", "icbf", "unidad_nacional_proteccion",
		"jueces_civiles", "inspeccion_policia", "defensoria_familia",
		"ejercito_nacional", "armada_colombiana", "no_sabe",
	},
}

// El catálogo embebido se carga al inicializar el paquete; si estuviera roto,
// cualquier test del paquete moriría en el panic de loadCompetenceKinds.
func TestCompetenceCatalogLoaded(t *testing.T) {
	total := 0
	for _, kind := range CompetenceKinds {
		if kind.Code == "" || kind.Name == "" || kind.Sector == "" {
			t.Fatalf("entidad incompleta en el catálogo: %+v", kind)
		}
		if len(kind.Competences) == 0 {
			t.Errorf("entidad %q sin competencias", kind.Code)
		}
		if _, ok := formSectorByCatalogSector[kind.Sector]; !ok {
			t.Errorf("entidad %q tiene sector %q sin equivalente en el formulario", kind.Code, kind.Sector)
		}
		total += len(kind.Competences)
	}
	if len(CompetenceKinds) != 38 || total != 237 {
		t.Fatalf("catálogo = %d entidades / %d competencias, esperaba 38/237",
			len(CompetenceKinds), total)
	}
}

// Este es el test que evita el modo de falla que importa: un código mal escrito en
// el mapa deja la pregunta sin opciones en producción y nadie se entera.
func TestInstitutionMapPointsToRealKinds(t *testing.T) {
	exists := make(map[string]bool, len(CompetenceKinds))
	for _, kind := range CompetenceKinds {
		exists[kind.Code] = true
	}
	for institution, kindCodes := range competenceKindsByInstitution {
		if len(kindCodes) == 0 {
			t.Errorf("institución %q mapeada a lista vacía", institution)
		}
		for _, code := range kindCodes {
			if !exists[code] {
				t.Errorf("institución %q apunta a la entidad %q, que no está en el catálogo", institution, code)
			}
		}
	}
}

// Ninguna institución que el agente pueda marcar puede dejar la pregunta vacía.
func TestCompetenceOptionsNeverEmptyForFormInstitutions(t *testing.T) {
	for sector, institutions := range institutionsBySector {
		for _, institution := range institutions {
			options := CompetenceOptions(sector, []string{institution})
			if len(options) == 0 {
				t.Errorf("sector %q + institución %q → 0 opciones", sector, institution)
			}
			for _, option := range options {
				if option.Label == "" || option.Value == "" {
					t.Errorf("sector %q + institución %q → opción incompleta %+v", sector, institution, option)
				}
			}
		}
		// Sin institución marcada todavía, el respaldo por sector debe dar opciones.
		if got := CompetenceOptions(sector, nil); len(got) == 0 {
			t.Errorf("sector %q sin institución → 0 opciones, esperaba el respaldo por sector", sector)
		}
	}
}

func TestCompetenceOptions(t *testing.T) {
	// "ips" suma las dos entidades IPS del Excel (15 competencias cada una), y
	// marcar además hospital público no las repite.
	ips := CompetenceOptions(FormSectorSalud, []string{"ips"})
	if len(ips) != 30 {
		t.Errorf("ips → %d opciones, esperaba 30", len(ips))
	}
	if both := CompetenceOptions(FormSectorSalud, []string{"ips", "hospital_publico"}); len(both) != 30 {
		t.Errorf("ips + hospital_publico → %d opciones, esperaba 30 sin repetir", len(both))
	}

	// "policía judicial" junta CTI, SIJIN y DIJIN: 11 + 6 + 5.
	if got := CompetenceOptions(FormSectorJusticia, []string{"policia_judicial"}); len(got) != 22 {
		t.Errorf("policia_judicial → %d opciones, esperaba 22", len(got))
	}

	// "no sabe" no mapea a ninguna entidad: cae al sector completo.
	noSabe := CompetenceOptions(FormSectorProteccion, []string{"no_sabe"})
	sectorOnly := CompetenceOptions(FormSectorProteccion, nil)
	if len(noSabe) != len(sectorOnly) || len(noSabe) != 44 {
		t.Errorf("no_sabe → %d opciones y sector → %d, esperaba 44 en ambos", len(noSabe), len(sectorOnly))
	}

	// Una institución de otro sector manda: el agente marcó Fiscalía dentro de
	// protección y debe recibir las competencias de Fiscalía, no las del sector.
	if got := CompetenceOptions(FormSectorProteccion, []string{"fiscalia_proteccion"}); len(got) != 8 {
		t.Errorf("fiscalia_proteccion → %d opciones, esperaba las 8 de la Fiscalía", len(got))
	}

	// El catálogo no cubre las barreras transversales y el formulario no pregunta
	// institución para ellas: sin opciones, la pregunta queda oculta.
	if got := CompetenceOptions("barrera_transversal", nil); len(got) != 0 {
		t.Errorf("barrera_transversal → %d opciones, esperaba 0", len(got))
	}
	if got := CompetenceOptions("", nil); len(got) != 0 {
		t.Errorf("sector vacío → %d opciones, esperaba 0", len(got))
	}
}
