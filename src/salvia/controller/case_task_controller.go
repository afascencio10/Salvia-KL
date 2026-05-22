package controller

import (
	"bitsflow/salvia/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CaseTaskController expone los endpoints HTTP de CaseTask usando Gin.
type CaseTaskController struct {
	svc service.CaseTaskService
}

func NewCaseTaskController(svc service.CaseTaskService) *CaseTaskController {
	return &CaseTaskController{svc: svc}
}

// RegisterRoutes registra las rutas de CaseTask en el grupo /api/v1.
//
//	GET  /api/v1/case-tasks?assignedUserId=<id>   → tareas con barrera asignadas al usuario
//	POST /api/v1/case-tasks/:id/complete           → marcar tarea como completada
func (c *CaseTaskController) RegisterRoutes(rg *gin.RouterGroup) {
	group := rg.Group("/case-tasks")
	group.GET("", c.ListByAssignedUser)
	group.POST("/:id/complete", c.Complete)
}

// ListByAssignedUser devuelve las tareas con barrierId asignadas al usuario,
// enriquecidas con datos de barrier_v2 y victim_case.
//
//	GET /api/v1/case-tasks?assignedUserId=<userICode>
func (c *CaseTaskController) ListByAssignedUser(ctx *gin.Context) {
	assignedUserID := ctx.Query("assignedUserId")
	if assignedUserID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "parámetro 'assignedUserId' requerido"})
		return
	}

	items, err := c.svc.ListByAssignedUserIDWithRelations(ctx.Request.Context(), assignedUserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}

	ctx.JSON(http.StatusOK, items)
}

// Complete marca una CaseTask como completada registrando su resultado.
//
//	POST /api/v1/case-tasks/:id/complete
//	Body: { "result": "Descripción de lo gestionado" }
func (c *CaseTaskController) Complete(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id requerido"})
		return
	}

	var body struct {
		Result string `json:"result" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "campo 'result' requerido"})
		return
	}

	if err := c.svc.Complete(ctx.Request.Context(), id, body.Result); err != nil {
		if err == service.ErrCaseTaskNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "tarea no encontrada"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "tarea completada exitosamente"})
}
