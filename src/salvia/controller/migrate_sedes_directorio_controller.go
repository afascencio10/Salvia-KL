package controller

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Tipos del Directorio de Entidades contra el nombre de la entidad en la hoja 4.
// Solo los que corresponden sin ambigüedad.
//
// Queda fuera AL (alcaldías): ya tienen sede en los 1121 municipios y el
// directorio solo agregaría duplicados. También quedan fuera RU y LE, que no son
// entidades, y los que el equipo no ha confirmado todavía: hospitales, centros
// de salud, IPS/EPS y clínicas contra "IPS"; centros zonales contra ICBF o
// Defensorías de Familia; CAIVAS y CAVIF como sedes de Fiscalía; y Personería,
// Defensoría, Procuraduría y UARIV, que no están en la hoja 4.
var entidadPorTipoDirectorio = map[string]string{
	"FI": "Fiscalía General de la Nación (FGN)",
	"CF": "Comisarías de Familia",
	"PN": "Policía Nacional de Colombia",
	"IP": "Inspecciones de Policía",
	"JZ": "Jueces de Control de Garantías/Civil Municipal o Promiscuo/Familia/Tutela",
	"IC": "Instituto Colombiano de Bienestar Familiar (ICBF)",
	"SS": "Secretarías de Salud (Dptal./Distrital/Municipal)",
	"SE": "Secretarías de Educación",
	"SG": "Secretarias de la Mujer",
	"GO": "Gobernaciones",
}

const maxTextoSede = 254

// recortaTexto deja el texto dentro del ancho de la columna, cortando por runas
// para no partir una tilde por la mitad.
func recortaTexto(texto string, max int) string {
	texto = strings.TrimSpace(texto)
	r := []rune(texto)
	if len(r) <= max {
		return texto
	}
	return string(r[:max])
}

type sedeOmitida struct {
	Nombre    string `json:"nombre"`
	Municipio string `json:"municipio"`
	Motivo    string `json:"motivo"`
}

type sedesDirectorioResult struct {
	SedesCreadas    int           `json:"sedesCreadas"`
	SedesYaExistian int           `json:"sedesYaExistian"`
	Omitidas        []sedeOmitida `json:"omitidas,omitempty"`
	Errores         []string      `json:"errores,omitempty"`
	DryRun          bool          `json:"dryRun"`
}

// MigrateSedesDirectorio crea una sede en salvia.entity_branch por cada registro
// del Directorio de Entidades cuyo tipo corresponda a una entidad de la hoja 4.
// El icode de la sede es el id del registro del directorio, así que volver a
// correrla no duplica nada.
//
// No crea sedes en municipios homónimos. El directorio resolvió varias filas al
// municipio equivocado — la Fiscalía Seccional de Armenia, Quindío quedó con el
// DIVIPOLA de Armenia, Antioquia — y ahí el dato de origen está mal, no la
// migración. Esas filas salen en "omitidas" para corregirlas en el directorio.
//
//	POST /api/v1/admin/migrate/sedes-directorio?dry_run=true
func (c *MigrateController) MigrateSedesDirectorio(ctx *gin.Context) {
	dryRun := ctx.Query("dry_run") == "true"
	res := sedesDirectorioResult{DryRun: dryRun}
	db := c.db.WithContext(ctx.Request.Context())

	type entidadFila struct {
		EntityID   int64  `gorm:"column:entity_id"`
		EntityName string `gorm:"column:entity_name"`
	}
	var entidades []entidadFila
	if err := db.Raw(`SELECT entity_id, entity_name FROM salvia.entity`).Scan(&entidades).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error listando entidades: " + err.Error()})
		return
	}
	idPorClave := make(map[string]int64, len(entidades))
	for _, e := range entidades {
		idPorClave[claveNombre(e.EntityName)] = e.EntityID
	}

	// Un mismo nombre de municipio en varios departamentos: el directorio no
	// distingue cuál es, así que no se puede confiar en su DIVIPOLA.
	var homonimos []string
	db.Raw(`SELECT city_name FROM security.city GROUP BY city_name HAVING COUNT(*) > 1`).Scan(&homonimos)
	esHomonimo := make(map[string]bool, len(homonimos))
	for _, n := range homonimos {
		esHomonimo[claveNombre(n)] = true
	}

	tipos := make([]string, 0, len(entidadPorTipoDirectorio))
	for tipo := range entidadPorTipoDirectorio {
		tipos = append(tipos, tipo)
	}

	// town_i_code guarda el DIVIPOLA cuando viene, igual que lo lee el Directorio
	// de Entidades; si no, se cae al menor centro poblado del municipio.
	type sedeFila struct {
		ID        string   `gorm:"column:id"`
		Nombre    string   `gorm:"column:name"`
		Direccion string   `gorm:"column:address"`
		Tipo      string   `gorm:"column:type"`
		Municipio string   `gorm:"column:city_name"`
		TownCode  string   `gorm:"column:town_code"`
		Latitud   *float64 `gorm:"column:latitude"`
		Longitud  *float64 `gorm:"column:longitude"`
	}
	var filas []sedeFila
	err := db.Raw(`
		SELECT d.id, d.name, d.address, d.type, c.city_name,
		       COALESCE(NULLIF(TRIM(d.town_i_code), ''),
		                (SELECT MIN(t.town_code) FROM security.town t WHERE t.city_id = c.city_id)) AS town_code,
		       d.latitude, d.longitude
		FROM salvia.directories d
		JOIN security.city c ON c.city_i_code = d.city_id
		WHERE d.active = true AND d.type IN ?
		ORDER BY d.name`, tipos).Scan(&filas).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error listando el directorio: " + err.Error()})
		return
	}

	for _, fila := range filas {
		omite := func(motivo string) {
			res.Omitidas = append(res.Omitidas, sedeOmitida{
				Nombre: fila.Nombre, Municipio: fila.Municipio, Motivo: motivo,
			})
		}

		nombreEntidad := entidadPorTipoDirectorio[fila.Tipo]
		clave := claveNombre(nombreEntidad)
		if alias, ok := nombreEnBase[nombreEntidad]; ok {
			clave = claveNombre(alias)
		}
		entityID, hayEntidad := idPorClave[clave]
		if !hayEntidad {
			omite("la entidad " + nombreEntidad + " todavía no existe: corre antes entidades-acciones")
			continue
		}
		if strings.TrimSpace(fila.TownCode) == "" {
			omite("sin municipio resuelto")
			continue
		}
		if esHomonimo[claveNombre(fila.Municipio)] {
			omite("el municipio " + fila.Municipio + " existe en más de un departamento y el directorio no distingue cuál")
			continue
		}

		var existe int64
		db.Raw(`SELECT COUNT(*) FROM salvia.entity_branch WHERE entity_branch_i_code = ?`, fila.ID).Scan(&existe)
		if existe > 0 {
			res.SedesYaExistian++
			continue
		}
		res.SedesCreadas++
		if dryRun {
			continue
		}

		nombre := recortaTexto(fila.Nombre, maxTextoSede)
		err := db.Exec(`
			INSERT INTO salvia.entity_branch
				(entity_branch_i_code, entity_branch_creation_date, entity_branch_update_date,
				 entity_branch_name, entity_branch_description, entity_branch_address,
				 entity_branch_latitude, entity_branch_longitude,
				 entity_branch_town_code, entity_id, entity_branch_source)
			VALUES (?, now(), now(), ?, ?, ?, ?, ?, ?, ?, 's')`,
			fila.ID, nombre, nombre, recortaTexto(fila.Direccion, maxTextoSede),
			fila.Latitud, fila.Longitud, fila.TownCode, entityID).Error
		if err != nil {
			res.SedesCreadas--
			res.Errores = append(res.Errores, fila.Nombre+": "+err.Error())
		}
	}

	ctx.JSON(http.StatusOK, res)
}
