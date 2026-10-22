package salvia_config

// Volcado de las columnas B (Entidad) y E (Acciones Institucionales) de la hoja
// "4.EntidadesAcciones" — 38 entidades, 235 acciones.
//
// Va aparte de competences.json y no lo reemplaza: los códigos de ese catálogo
// llevan un hash del texto de la acción y quedan guardados en los casos ya
// registrados, así que regenerarlo dejaría esas competencias huérfanas. Este
// archivo alimenta salvia.entity_obligation, que guarda texto libre.
//
// Para actualizarlo se regenera con scripts/extract_entidades_acciones.py y se
// versiona el resultado.

import (
	"embed"
	"encoding/json"
	"fmt"
)

//go:embed entidades_acciones.json
var entidadesAccionesFS embed.FS

// EntidadAcciones es una entidad del Excel con sus acciones institucionales.
type EntidadAcciones struct {
	Nombre      string   `json:"nombre"`
	SectorExcel string   `json:"sectorExcel"`
	Acciones    []string `json:"acciones"`
}

// EntidadesAcciones son las entidades en el orden del Excel.
var EntidadesAcciones = loadEntidadesAcciones()

// Un JSON inválido revienta al arrancar, no en mitad de una migración.
func loadEntidadesAcciones() []EntidadAcciones {
	raw, err := entidadesAccionesFS.ReadFile("entidades_acciones.json")
	if err != nil {
		panic(fmt.Sprintf("salvia_config: leer entidades_acciones.json: %v", err))
	}
	var catalogo struct {
		Entidades []EntidadAcciones `json:"entidades"`
	}
	if err := json.Unmarshal(raw, &catalogo); err != nil {
		panic(fmt.Sprintf("salvia_config: entidades_acciones.json inválido: %v", err))
	}
	if len(catalogo.Entidades) == 0 {
		panic("salvia_config: entidades_acciones.json sin entidades")
	}
	return catalogo.Entidades
}
