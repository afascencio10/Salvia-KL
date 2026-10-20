package controller

import (
	"net/http"
	"strings"

	"bitsflow/internal/models"
	salvia_config "bitsflow/salvia/config"

	"github.com/gin-gonic/gin"
)

const (
	entidadAlcaldias      = "Alcaldías"
	sectorAlcaldias       = "os"
	direccionSinRegistrar = "Sin dirección registrada"
	maxEtiquetaObligacion = 255
)

// etiquetaObligacion recorta la acción al ancho de entity_obligation.label
// (varchar 255), cortando en un espacio y dejando visible el recorte.
func etiquetaObligacion(accion string) string {
	accion = strings.TrimSpace(accion)
	r := []rune(accion)
	if len(r) <= maxEtiquetaObligacion {
		return accion
	}
	r = r[:maxEtiquetaObligacion-1]
	if i := strings.LastIndex(string(r), " "); i > 0 {
		return string([]rune(string(r)[:i])) + "…"
	}
	return string(r) + "…"
}

type alcaldiasResult struct {
	EntidadID           int64    `json:"entidadId"`
	EntidadCreada       bool     `json:"entidadCreada"`
	ObligacionesCreadas int      `json:"obligacionesCreadas"`
	ObligacionesExisten int      `json:"obligacionesYaExistian"`
	SedesCreadas        int      `json:"sedesCreadas"`
	SedesYaExistian     int      `json:"sedesYaExistian"`
	MunicipiosSinTown   int      `json:"municipiosSinCentroPoblado"`
	Errores             []string `json:"errores,omitempty"`
	DryRun              bool     `json:"dryRun"`
}

// MigrateAlcaldias crea la entidad "Alcaldías", sus acciones institucionales y una
// sede por municipio. Es idempotente: lo que ya existe se cuenta y no se duplica.
//
//	POST /api/v1/admin/migrate/alcaldias?dry_run=true
func (c *MigrateController) MigrateAlcaldias(ctx *gin.Context) {
	if ctx.Query("inspect") == "true" {
		c.inspeccionarTablasAlcaldias(ctx)
		return
	}

	dryRun := ctx.Query("dry_run") == "true"
	res := alcaldiasResult{DryRun: dryRun}

	acciones := accionesDeAlcaldias()
	if len(acciones) == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "el catálogo embebido no trae la entidad Alcaldías"})
		return
	}

	db := c.db.WithContext(ctx.Request.Context())

	var entityID int64
	err := db.Raw(`SELECT entity_id FROM salvia.entity WHERE entity_name = ?`, entidadAlcaldias).Scan(&entityID).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error consultando la entidad: " + err.Error()})
		return
	}
	if entityID == 0 {
		res.EntidadCreada = true
		if !dryRun {
			err = db.Raw(`
				INSERT INTO salvia.entity
					(entity_i_code, entity_creation_date, entity_update_date, entity_name,
					 entity_is_interoperable, entity_interoperability_code, entity_response_time, entity_sector)
				VALUES (gen_random_uuid()::text, now(), now(), ?, 'n', 'ALCALDIAS', 0, ?)
				RETURNING entity_id`, entidadAlcaldias, sectorAlcaldias).Scan(&entityID).Error
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error creando la entidad: " + err.Error()})
				return
			}
		}
	}
	res.EntidadID = entityID

	for i, accion := range acciones {
		label := etiquetaObligacion(accion)
		var existe int64
		if entityID != 0 {
			db.Raw(`SELECT COUNT(*) FROM salvia.entity_obligation
			        WHERE entity_id = ? AND label = ? AND deleted_at IS NULL`, entityID, label).Scan(&existe)
		}
		if existe > 0 {
			res.ObligacionesExisten++
			continue
		}
		res.ObligacionesCreadas++
		if dryRun || entityID == 0 {
			continue
		}
		ob := models.EntityObligation{EntityID: entityID, Label: label, Order: i + 1}
		if err := db.Create(&ob).Error; err != nil {
			res.Errores = append(res.Errores, "obligación "+label[:40]+": "+err.Error())
		}
	}

	// Una sede por municipio. El town_code sale del menor centro poblado del
	// municipio: la cabecera termina en 000 y es el primero en orden.
	type municipio struct {
		CityName string `gorm:"column:city_name"`
		TownCode string `gorm:"column:town_code"`
	}
	var municipios []municipio
	if err := db.Raw(`
		SELECT c.city_name AS city_name, MIN(t.town_code) AS town_code
		FROM security.city c
		JOIN security.town t ON t.city_id = c.city_id
		GROUP BY c.city_id, c.city_name
		ORDER BY c.city_name`).Scan(&municipios).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error listando municipios: " + err.Error()})
		return
	}

	for _, m := range municipios {
		if strings.TrimSpace(m.TownCode) == "" {
			res.MunicipiosSinTown++
			continue
		}
		icode := "Alcaldia" + m.TownCode
		var existe int64
		db.Raw(`SELECT COUNT(*) FROM salvia.entity_branch WHERE entity_branch_i_code = ?`, icode).Scan(&existe)
		if existe > 0 {
			res.SedesYaExistian++
			continue
		}
		res.SedesCreadas++
		if dryRun || entityID == 0 {
			continue
		}
		nombre := "Alcaldía de " + m.CityName
		err := db.Exec(`
			INSERT INTO salvia.entity_branch
				(entity_branch_i_code, entity_branch_creation_date, entity_branch_update_date,
				 entity_branch_name, entity_branch_description, entity_branch_address,
				 entity_branch_town_code, entity_id, entity_branch_source)
			VALUES (?, now(), now(), ?, ?, ?, ?, ?, 's')`,
			icode, nombre, nombre, direccionSinRegistrar, m.TownCode, entityID).Error
		if err != nil {
			res.SedesCreadas--
			res.Errores = append(res.Errores, icode+": "+err.Error())
		}
	}

	ctx.JSON(http.StatusOK, res)
}

// inspeccionarTablasAlcaldias devuelve las columnas reales de entity y
// entity_branch más una fila de ejemplo. Solo lectura: los DAO de Go no conocen
// todas las columnas de estas tablas y alguna es NOT NULL.
func (c *MigrateController) inspeccionarTablasAlcaldias(ctx *gin.Context) {
	db := c.db.WithContext(ctx.Request.Context())

	type columna struct {
		Nombre   string  `gorm:"column:column_name"    json:"columna"`
		Tipo     string  `gorm:"column:data_type"      json:"tipo"`
		Nullable string  `gorm:"column:is_nullable"    json:"aceptaNulo"`
		Default  *string `gorm:"column:column_default" json:"default,omitempty"`
	}

	salida := gin.H{}
	for _, tabla := range []string{"entity", "entity_branch", "entity_obligation"} {
		var cols []columna
		if err := db.Raw(`
			SELECT column_name, data_type, is_nullable, column_default
			FROM information_schema.columns
			WHERE table_schema = 'salvia' AND table_name = ?
			ORDER BY ordinal_position`, tabla).Scan(&cols).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error leyendo columnas de " + tabla + ": " + err.Error()})
			return
		}

		var ejemplo *string
		db.Raw(`SELECT row_to_json(t)::text FROM salvia.` + tabla + ` t LIMIT 1`).Scan(&ejemplo)

		salida[tabla] = gin.H{"columnas": cols, "ejemplo": ejemplo}
	}

	ctx.JSON(http.StatusOK, salida)
}

// accionesDeAlcaldias saca las acciones del catálogo embebido (hoja
// 4.EntidadesAcciones) en vez de repetirlas acá: si se regenera el JSON, la
// migración sigue el archivo.
func accionesDeAlcaldias() []string {
	for _, kind := range salvia_config.CompetenceKinds {
		if !strings.HasPrefix(strings.ToLower(kind.Name), "alcald") {
			continue
		}
		out := make([]string, 0, len(kind.Competences))
		for _, comp := range kind.Competences {
			if accion := strings.TrimSpace(comp.Action); accion != "" {
				out = append(out, accion)
			}
		}
		return out
	}
	return nil
}
