package controller

import (
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
//	GET /api/v1/barriers-v2/department          → ListActiveByDepartment (Barreras Departamento)
func (c *BarrierV2GinController) RegisterRoutes(rg *gin.RouterGroup) {
	barriers := rg.Group("/barriers-v2")
	barriers.GET("", c.List)
	barriers.GET("/department", c.ListByDepartment)
	barriers.GET("/:id/detail", c.GetDetail)
	barriers.GET("/:id/follow-ups", c.ListFollowUps)
}

// ListByDepartment devuelve todas las barreras activas de los casos del
// departamento asignado al Enlace Territorial autenticado, sin filtrar por
// asignación de tareas — a diferencia de GET /api/v1/case-tasks (Mis
// Barreras). El departamento se resuelve siempre desde la sesión, nunca
// desde un parámetro del cliente.
//
//	GET /api/v1/barriers-v2/department?docNumber=&entidad=&cityId=
func (c *BarrierV2GinController) ListByDepartment(ctx *gin.Context) {
	s := requireEnlaceAccess(ctx)
	if s == nil {
		return
	}

	filters := repository.DepartmentBarrierFilters{
		DocNumber: ctx.Query("docNumber"),
		Entidad:   ctx.Query("entidad"),
		CityID:    ctx.Query("cityId"),
	}

	items, err := c.svc.ListActiveByDepartment(ctx.Request.Context(), s.AssignedDepartmentID, filters)
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
