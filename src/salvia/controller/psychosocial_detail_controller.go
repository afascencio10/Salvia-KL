// Package controller — psychosocial_detail_controller.go
// Endpoint para la pantalla de detalle de remisión psicosocial.
package controller

import (
	"bitsflow/salvia/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PsychosocialDetailController struct {
	svc service.PsychosocialDetailService
}

func NewPsychosocialDetailController(db *gorm.DB) *PsychosocialDetailController {
	return &PsychosocialDetailController{
		svc: service.NewPsychosocialDetailService(db),
	}
}

func (c *PsychosocialDetailController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/psychosocial-support/:id/detail", c.GetDetail)
}

// GetDetail retorna el detalle completo de una remisión psicosocial.
// GET /api/v1/psychosocial-support/:id/detail
func (c *PsychosocialDetailController) GetDetail(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id requerido"})
		return
	}

	detail, err := c.svc.GetDetail(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "remisión no encontrada"})
		return
	}

	ctx.JSON(http.StatusOK, detail)
}
