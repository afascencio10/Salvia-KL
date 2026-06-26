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
//
//	GET /api/v1/entity-branches?town_code={townCode}
func (c *EntityBranchAPIController) ListByTownCode(ctx *gin.Context) {
	townCode := ctx.Query("town_code")
	if townCode == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "el parámetro 'town_code' es requerido"})
		return
	}

	log.Printf("[DEBUG] entity-branches: town_code=%q", townCode)

	branches := make([]entityBranchOption, 0)
	err := c.db.WithContext(ctx.Request.Context()).
		Table("salvia.entity_branch").
		Select("entity_branch_id, entity_branch_i_code, entity_branch_name").
		Where("entity_branch_town_code = ?", townCode).
		Order("entity_branch_name ASC").
		Scan(&branches).Error
	if err != nil {
		log.Printf("[ERROR] entity-branches: town_code=%q error=%v", townCode, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}

	log.Printf("[DEBUG] entity-branches: town_code=%q results=%d", townCode, len(branches))
	ctx.JSON(http.StatusOK, branches)
}
