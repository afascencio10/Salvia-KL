package controller

import (
	"net/http"
	"strings"

	"bitsflow/internal/models"
	salvia_config "bitsflow/salvia/config"

	"github.com/gin-gonic/gin"
)

// Sectores del Excel contra entity_sector. Los que no tienen código propio caen
// en "os" (otros sectores), igual que quedaron las alcaldías.
var sectorPorExcel = map[string]string{
	"SALUD":      "he",
	"JUSTICIA":   "js",
	"PROTECCIÓN": "pt",
}

// El Excel renombró entidades que llevan años en la base. Sin este puente la
// migración crearía un duplicado de cada una.
var nombreEnBase = map[string]string{
	"Comisarías de Familia":                                                                                        "Comisaría de Familia",
	"Fiscalía General de la Nación (FGN)":                                                                          "Fiscalía General de la Nación",
	"Instituto Colombiano de Bienestar Familiar (ICBF)":                                                            "Instituto Colombiano de Bienestar Familiar",
	"Instituto Nacional de Medicina Legal y Ciencias Forenses (INMLCF)":                                            "Instituto Nacional de Medicina Legal y Ciencias Forenses",
	"Instituciones Prestadoras de Servicios de Salud (IPS)":                                                        "IPS",
	"Instituciones Prestadoras de Servicios de Salud (IPS) (Clínicas, Hospitales, Consultorios, Centros de Salud)": "IPS",
	"Policía Nacional de Colombia":                                                                                 "Policía Nacional",
}

var sinTilde = strings.NewReplacer(
	"á", "a", "é", "e", "í", "i", "ó", "o", "ú", "u", "ü", "u", "ñ", "n",
	"Á", "a", "É", "e", "Í", "i", "Ó", "o", "Ú", "u", "Ü", "u", "Ñ", "n",
)

// claveNombre deja solo letras y dígitos sin tildes, para comparar el nombre del
// Excel contra el que ya está en la base pese a mayúsculas, acentos o puntuación.
func claveNombre(nombre string) string {
	var b strings.Builder
	for _, r := range sinTilde.Replace(strings.ToLower(strings.TrimSpace(nombre))) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// codigoInteroperabilidad arma el valor de entity_interoperability_code, que es
// varchar(32), a partir del nombre.
func codigoInteroperabilidad(nombre string) string {
	clave := claveNombre(nombre)
	if clave == "" {
		return "ENTIDAD"
	}
	if len(clave) > 32 {
		clave = clave[:32]
	}
	return strings.ToUpper(clave)
}

// sectorDeEntidad traduce el sector del Excel al código de entity_sector.
func sectorDeEntidad(sectorExcel string) string {
	if codigo, ok := sectorPorExcel[strings.TrimSpace(sectorExcel)]; ok {
		return codigo
	}
	return sectorAlcaldias
}

type sobranteObligacion struct {
	Entidad  string `json:"entidad"`
	Etiqueta string `json:"etiqueta"`
}

type entidadesAccionesResult struct {
	EntidadesCreadas   int                  `json:"entidadesCreadas"`
	EntidadesYaExistan int                  `json:"entidadesYaExistian"`
	AccionesCreadas    int                  `json:"accionesCreadas"`
	AccionesYaExistan  int                  `json:"accionesYaExistian"`
	Sobrantes          []sobranteObligacion `json:"sobrantes,omitempty"`
	Errores            []string             `json:"errores,omitempty"`
	DryRun             bool                 `json:"dryRun"`
}

// MigrateEntidadesAcciones crea las 38 entidades de la hoja 4 y carga sus
// acciones institucionales en salvia.entity_obligation. Es idempotente: lo que
// ya existe se cuenta y no se duplica.
//
// No borra nada. Las obligaciones que están en la base y no en el Excel se
// devuelven en "sobrantes" para que el equipo decida: varias se cargaron a mano.
//
//	POST /api/v1/admin/migrate/entidades-acciones?dry_run=true
func (c *MigrateController) MigrateEntidadesAcciones(ctx *gin.Context) {
	dryRun := ctx.Query("dry_run") == "true"
	res := entidadesAccionesResult{DryRun: dryRun}
	db := c.db.WithContext(ctx.Request.Context())

	type entidadFila struct {
		EntityID   int64  `gorm:"column:entity_id"`
		EntityName string `gorm:"column:entity_name"`
	}
	var existentes []entidadFila
	if err := db.Raw(`SELECT entity_id, entity_name FROM salvia.entity`).Scan(&existentes).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error listando entidades: " + err.Error()})
		return
	}
	idPorClave := make(map[string]int64, len(existentes))
	for _, e := range existentes {
		idPorClave[claveNombre(e.EntityName)] = e.EntityID
	}

	// Acumulan, por entidad, las etiquetas del catálogo y el nombre a reportar:
	// varias filas del Excel pueden caer en la misma entidad.
	catalogoPorEntidad := map[int64]map[string]bool{}
	nombrePorEntidad := map[int64]string{}

	// entity_city_code es NOT NULL y no existe en el DAO: se copia el de una
	// entidad ya cargada en vez de inventar la convención.
	var cityCode string
	db.Raw(`SELECT entity_city_code FROM salvia.entity
	        WHERE entity_city_code IS NOT NULL LIMIT 1`).Scan(&cityCode)

	for _, entidad := range salvia_config.EntidadesAcciones {
		clave := claveNombre(entidad.Nombre)
		if alias, ok := nombreEnBase[entidad.Nombre]; ok {
			clave = claveNombre(alias)
		}

		entityID, existe := idPorClave[clave]
		if existe {
			res.EntidadesYaExistan++
		} else {
			res.EntidadesCreadas++
			if !dryRun {
				err := db.Raw(`
					INSERT INTO salvia.entity
						(entity_i_code, entity_creation_date, entity_update_date, entity_city_code,
						 entity_name, entity_description, entity_is_interoperable,
						 entity_interoperability_code, entity_response_time, entity_sector)
					VALUES (gen_random_uuid()::text, now(), now(), ?, ?, ?, 'n', ?, 0, ?)
					RETURNING entity_id`,
					cityCode, entidad.Nombre, entidad.Nombre,
					codigoInteroperabilidad(entidad.Nombre), sectorDeEntidad(entidad.SectorExcel),
				).Scan(&entityID).Error
				if err != nil {
					res.Errores = append(res.Errores, entidad.Nombre+": "+err.Error())
					continue
				}
				// Dos filas del Excel apuntan a la misma entidad (las dos de IPS):
				// la segunda debe encontrar la que acaba de crear la primera.
				idPorClave[clave] = entityID
			}
		}

		if entityID == 0 {
			// Dry run sobre una entidad que todavía no existe: sus acciones se
			// cuentan como nuevas sin consultar la base.
			res.AccionesCreadas += len(entidad.Acciones)
			continue
		}

		var enBase []string
		db.Raw(`SELECT label FROM salvia.entity_obligation
		        WHERE entity_id = ? AND deleted_at IS NULL`, entityID).Scan(&enBase)
		yaEsta := make(map[string]bool, len(enBase))
		for _, label := range enBase {
			yaEsta[claveNombre(label)] = true
		}

		if catalogoPorEntidad[entityID] == nil {
			catalogoPorEntidad[entityID] = map[string]bool{}
			nombrePorEntidad[entityID] = entidad.Nombre
		}
		delCatalogo := catalogoPorEntidad[entityID]
		for i, accion := range entidad.Acciones {
			label := etiquetaObligacion(accion)
			delCatalogo[claveNombre(label)] = true
			if yaEsta[claveNombre(label)] {
				res.AccionesYaExistan++
				continue
			}
			res.AccionesCreadas++
			if dryRun {
				continue
			}
			ob := models.EntityObligation{EntityID: entityID, Label: label, Order: i + 1}
			if err := db.Create(&ob).Error; err != nil {
				res.Errores = append(res.Errores, entidad.Nombre+" / obligación: "+err.Error())
			}
		}

	}

	// Los sobrantes se calculan al final y por entidad, no por fila del Excel:
	// las dos filas de IPS caen en la misma entidad, y comparar contra una sola
	// marcaría como sobrantes las acciones que acaba de cargar la otra.
	for entityID, delCatalogo := range catalogoPorEntidad {
		var enBase []string
		db.Raw(`SELECT label FROM salvia.entity_obligation
		        WHERE entity_id = ? AND deleted_at IS NULL`, entityID).Scan(&enBase)
		for _, label := range enBase {
			if !delCatalogo[claveNombre(label)] {
				res.Sobrantes = append(res.Sobrantes, sobranteObligacion{
					Entidad: nombrePorEntidad[entityID], Etiqueta: label,
				})
			}
		}
	}

	ctx.JSON(http.StatusOK, res)
}
