package controller

import (
	"testing"

	salvia_config "bitsflow/salvia/config"
)

// Una llave de nombreEnBase que no coincida exactamente con el nombre del Excel
// no falla: el alias nunca se aplica y la migración crea una entidad duplicada
// de una que ya existe. Es el error más caro de este archivo y no lo ve el
// compilador.
func TestAliasApuntanAEntidadesDelCatalogo(t *testing.T) {
	delCatalogo := make(map[string]bool, len(salvia_config.EntidadesAcciones))
	for _, e := range salvia_config.EntidadesAcciones {
		delCatalogo[e.Nombre] = true
	}
	for nombreExcel := range nombreEnBase {
		if !delCatalogo[nombreExcel] {
			t.Errorf("el alias %q no corresponde a ninguna entidad del catálogo", nombreExcel)
		}
	}
}

// El catálogo embebido es el insumo completo de la migración: una entidad sin
// nombre o sin acciones es una fila mal extraída del Excel.
func TestCatalogoCompleto(t *testing.T) {
	if len(salvia_config.EntidadesAcciones) != 38 {
		t.Fatalf("se esperaban 38 entidades y hay %d", len(salvia_config.EntidadesAcciones))
	}
	total := 0
	for _, e := range salvia_config.EntidadesAcciones {
		if e.Nombre == "" {
			t.Error("hay una entidad sin nombre")
		}
		if len(e.Acciones) == 0 {
			t.Errorf("%q no trae acciones", e.Nombre)
		}
		total += len(e.Acciones)
	}
	if total != 235 {
		t.Errorf("se esperaban 235 acciones y hay %d", total)
	}
}

func TestClaveNombre(t *testing.T) {
	casos := []struct{ entra, sale string }{
		{"Comisarías de Familia", "comisariasdefamilia"},
		{"Fiscalía General de la Nación (FGN)", "fiscaliageneraldelanacionfgn"},
		{"  IPS  ", "ips"},
		{"Policía Nacional", "policianacional"},
	}
	for _, caso := range casos {
		if got := claveNombre(caso.entra); got != caso.sale {
			t.Errorf("claveNombre(%q) = %q, se esperaba %q", caso.entra, got, caso.sale)
		}
	}
}

// entity_interoperability_code es varchar(32): un nombre largo tiene que caber.
func TestCodigoInteroperabilidadCabe(t *testing.T) {
	for _, e := range salvia_config.EntidadesAcciones {
		if n := len(codigoInteroperabilidad(e.Nombre)); n == 0 || n > 32 {
			t.Errorf("%q genera un código de %d caracteres", e.Nombre, n)
		}
	}
}

// Todo sector del Excel tiene que caer en un código de dos letras.
func TestSectorDeEntidad(t *testing.T) {
	for _, e := range salvia_config.EntidadesAcciones {
		if got := sectorDeEntidad(e.SectorExcel); len(got) != 2 {
			t.Errorf("sector %q de %q da %q", e.SectorExcel, e.Nombre, got)
		}
	}
	if got := sectorDeEntidad("JUSTICIA"); got != "js" {
		t.Errorf("JUSTICIA debería dar js y dio %q", got)
	}
	if got := sectorDeEntidad("AGRICULTURA"); got != sectorAlcaldias {
		t.Errorf("un sector sin código propio debería caer en %q y dio %q", sectorAlcaldias, got)
	}
}
