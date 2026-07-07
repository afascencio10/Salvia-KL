package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// entityBranchOption es el DTO mínimo para el dropdown de entidades en el modal proyectar.
type entityBranchOption struct {
	ID    int64  `gorm:"column:entity_branch_id"     json:"id"`
	ICode string `gorm:"column:entity_branch_i_code" json:"icode"`
	Name  string `gorm:"column:entity_branch_name"   json:"name"`
}

// EntityBranchAPIController expone entity_branch como API REST en /api/v1.
type EntityBranchAPIController struct {
	db *gorm.DB
}

func NewEntityBranchAPIController(db *gorm.DB) *EntityBranchAPIController {
	return &EntityBranchAPIController{db: db}
}

func (c *EntityBranchAPIController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/entity-branches", c.ListByTownCode)
}

// ListByTownCode devuelve las sedes del municipio indicado, ordenadas por nombre.
// Soporta filtro opcional por sector de barrera (se mapea al código de entity_sector).
//
//	GET /api/v1/entity-branches?town_code={townCode}&sector={sector}
func (c *EntityBranchAPIController) ListByTownCode(ctx *gin.Context) {
	townCode := ctx.Query("town_code")
	if townCode == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "el parámetro 'town_code' es requerido"})
		return
	}

	// Mapear sector de barrera a código de entity_sector
	barrierSector := ctx.Query("sector")
	entitySectorCode := ""
	if barrierSector != "" {
		sectorMap := map[string]string{
			"salud":                "he",
			"proteccion":          "pt",
			"justicia":            "js",
			"otras_instituciones": "os",
		}
		entitySectorCode = sectorMap[barrierSector]
	}

	log.Printf("[DEBUG] entity-branches: town_code=%q barrierSector=%q entitySectorCode=%q", townCode, barrierSector, entitySectorCode)

	branches := make([]entityBranchOption, 0)
	var err error

	if entitySectorCode != "" {
		err = c.db.WithContext(ctx.Request.Context()).
			Raw(`SELECT DISTINCT ON (eb.entity_branch_name) eb.entity_branch_id, eb.entity_branch_i_code, eb.entity_branch_name
				FROM salvia.entity_branch eb
				JOIN salvia.entity e ON e.entity_id = eb.entity_id
				WHERE eb.entity_branch_town_code = ? AND e.entity_sector = ?
				ORDER BY eb.entity_branch_name ASC, eb.entity_branch_id ASC`, townCode, entitySectorCode).
			Scan(&branches).Error
	} else {
		err = c.db.WithContext(ctx.Request.Context()).
			Raw(`SELECT DISTINCT ON (entity_branch_name) entity_branch_id, entity_branch_i_code, entity_branch_name
				FROM salvia.entity_branch
				WHERE entity_branch_town_code = ?
				ORDER BY entity_branch_name ASC, entity_branch_id ASC`, townCode).
			Scan(&branches).Error
	}

	if err != nil {
		log.Printf("[ERROR] entity-branches: town_code=%q error=%v", townCode, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}

	log.Printf("[DEBUG] entity-branches: town_code=%q sector=%q results=%d", townCode, barrierSector, len(branches))
	ctx.JSON(http.StatusOK, branches)
}
