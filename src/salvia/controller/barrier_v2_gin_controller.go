package controller

import (
	"bitsflow/salvia/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// BarrierV2GinController expone los endpoints HTTP de BarrierV2 usando Gin.
type BarrierV2GinController struct {
	svc service.BarrierV2Service
}

func NewBarrierV2GinController(svc service.BarrierV2Service) *BarrierV2GinController {
	return &BarrierV2GinController{svc: svc}
}

// RegisterRoutes registra las rutas de BarrierV2 en el grupo /api/v1.
//
//	GET /api/v1/barriers-v2?createdById=<id>    → ListByCreatedBy (enriquecido con caso)
//	GET /api/v1/barriers-v2/:id/follow-ups      → ListFollowUps (seguimientos de una barrera)
func (c *BarrierV2GinController) RegisterRoutes(rg *gin.RouterGroup) {
	barriers := rg.Group("/barriers-v2")
	barriers.GET("", c.List)
	barriers.GET("/:id/follow-ups", c.ListFollowUps)
}

// List devuelve las barreras filtradas por createdById enriquecidas con datos del caso.
//
//	GET /api/v1/barriers-v2?createdById=<agentICode>
func (c *BarrierV2GinController) List(ctx *gin.Context) {
	createdByID := ctx.Query("createdById")
	if createdByID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "parámetro 'createdById' requerido"})
		return
	}

	items, err := c.svc.ListByCreatedByIDWithRelations(ctx.Request.Context(), createdByID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}

	ctx.JSON(http.StatusOK, items)
}

// ListFollowUps devuelve los seguimientos de una barrera en orden cronológico.
//
//	GET /api/v1/barriers-v2/:id/follow-ups
func (c *BarrierV2GinController) ListFollowUps(ctx *gin.Context) {
	barrierID := ctx.Param("id")
	if barrierID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id de barrera requerido"})
		return
	}

	items, err := c.svc.ListFollowUps(ctx.Request.Context(), barrierID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}

	ctx.JSON(http.StatusOK, items)
}
