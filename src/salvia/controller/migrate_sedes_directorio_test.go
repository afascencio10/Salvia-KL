package controller

import (
	"strings"
	"testing"

	"bitsflow/internal/models"
	salvia_config "bitsflow/salvia/config"
)

// Un nombre mal escrito en el mapa no falla: la entidad nunca resuelve y todas
// las sedes de ese tipo salen omitidas, con el combo igual de vacío que antes.
func TestTiposDirectorioApuntanAEntidadesDelCatalogo(t *testing.T) {
	delCatalogo := make(map[string]bool, len(salvia_config.EntidadesAcciones))
	for _, e := range salvia_config.EntidadesAcciones {
		delCatalogo[e.Nombre] = true
	}
	for tipo, nombre := range entidadPorTipoDirectorio {
		if !delCatalogo[nombre] {
			t.Errorf("el tipo %q apunta a %q, que no está en el catálogo", tipo, nombre)
		}
	}
}

// Los códigos del mapa tienen que ser tipos de directorio reales: uno inventado
// no trae filas y pasa desapercibido.
func TestTiposDirectorioExisten(t *testing.T) {
	for tipo := range entidadPorTipoDirectorio {
		if _, ok := models.DirectoryTypeLabels[tipo]; !ok {
			t.Errorf("%q no es un tipo del Directorio de Entidades", tipo)
		}
	}
}

// Las alcaldías ya tienen sede en los 1121 municipios: incluirlas acá duplicaría.
func TestTiposDirectorioExcluyeAlcaldias(t *testing.T) {
	if _, ok := entidadPorTipoDirectorio[models.DirectoryTypeAlcaldia]; ok {
		t.Error("AL no debe estar: las alcaldías ya tienen sede en todos los municipios")
	}
}

// entity_branch_name y entity_branch_address son varchar(254).
func TestRecortaTexto(t *testing.T) {
	if got := recortaTexto("  Fiscalía Seccional Leticia  ", maxTextoSede); got != "Fiscalía Seccional Leticia" {
		t.Errorf("un texto corto debería pasar limpio y dio %q", got)
	}
	largo := strings.Repeat("á", 300)
	got := recortaTexto(largo, maxTextoSede)
	if n := len([]rune(got)); n != maxTextoSede {
		t.Errorf("se esperaban %d runas y hay %d", maxTextoSede, n)
	}
}
