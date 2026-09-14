package controller

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
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
//	GET /api/v1/barriers-v2/department          → ListByDepartment (Barreras Departamento)
func (c *BarrierV2GinController) RegisterRoutes(rg *gin.RouterGroup) {
	barriers := rg.Group("/barriers-v2")
	barriers.GET("", c.List)
	barriers.GET("/department", c.ListByDepartment)
	barriers.GET("/mine", c.ListMine)
	barriers.GET("/:id/detail", c.GetDetail)
	barriers.GET("/:id/follow-ups", c.ListFollowUps)
}

// ListByDepartment devuelve las barreras del departamento asignado al Enlace
// Territorial autenticado (todas las no eliminadas). El departamento se
// resuelve siempre desde la sesión, nunca desde un parámetro del cliente.
//
//	GET /api/v1/barriers-v2/department?docNumber=&entidad=&cityId=&status=&page=&pageSize=
func (c *BarrierV2GinController) ListByDepartment(ctx *gin.Context) {
	s := requireEnlaceAccess(ctx)
	if s == nil {
		return
	}

	status := ctx.Query("status")
	if status != "" &&
		status != models.BarrierV2StatusOpen &&
		status != models.BarrierV2StatusEnGestion &&
		status != models.BarrierV2StatusArticulada &&
		status != models.BarrierV2StatusManaged {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "estado inválido"})
		return
	}

	pageSize := ginQueryInt(ctx, "pageSize", 10)
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	page := ginQueryInt(ctx, "page", 0)

	filters := repository.DepartmentBarrierFilters{
		DocNumber: ctx.Query("docNumber"),
		Entidad:   ctx.Query("entidad"),
		CityID:    ctx.Query("cityId"),
		Status:    status,
		Page:      page,
		PageSize:  pageSize,
	}

	items, total, err := c.svc.ListByDepartment(ctx.Request.Context(), s.AssignedDepartmentID, filters)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"items":    items,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// ListMine devuelve las barreras activas asignadas al usuario autenticado
// (user_territory_link_id = sesión) — Enlace Territorial u Operador
// territorial nacional. El usuario se toma de la sesión, nunca de query params.
//
//	GET /api/v1/barriers-v2/mine
func (c *BarrierV2GinController) ListMine(ctx *gin.Context) {
	s := requireMisBarrerasAccess(ctx)
	if s == nil {
		return
	}

	items, err := c.svc.ListMine(ctx.Request.Context(), s.UserICode)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}

	ctx.JSON(http.StatusOK, items)
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

// GetDetail devuelve el detalle completo de una barrera con datos de la víctima.
//
//	GET /api/v1/barriers-v2/:id/detail
func (c *BarrierV2GinController) GetDetail(ctx *gin.Context) {
	barrierID := ctx.Param("id")
	if barrierID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id de barrera requerido"})
		return
	}

	detail, err := c.svc.GetDetail(ctx.Request.Context(), barrierID)
	if err != nil {
		if err.Error() == "barrier_v2: registro no encontrado" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "barrera no encontrada"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}

	ctx.JSON(http.StatusOK, detail)
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
