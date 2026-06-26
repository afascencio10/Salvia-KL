package controller

import (
	"bitsflow/salvia/service"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
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
//	GET  /api/v1/case-tasks/:id                    → obtener tarea por ID
//	POST /api/v1/case-tasks/:id/complete           → marcar tarea como completada (flujo legacy)
//	PUT  /api/v1/case-tasks/:id/complete           → completar tarea con formData (case-task-modal)
func (c *CaseTaskController) RegisterRoutes(rg *gin.RouterGroup) {
	group := rg.Group("/case-tasks")
	group.GET("", c.ListByAssignedUser)
	group.GET("/:id", c.GetByID)
	group.POST("/:id/complete", c.Complete)
	group.PUT("/:id/complete", c.CompleteWithFormData)
}

// GetByID devuelve una CaseTask por su UUID.
//
//	GET /api/v1/case-tasks/:id
func (c *CaseTaskController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	task, err := c.svc.GetByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrCaseTaskNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "tarea no encontrada"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}
	ctx.JSON(http.StatusOK, task)
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

// CompleteWithFormData completa una CaseTask guardando el JSON del formulario
// y ejecutando los efectos de lado propios del tipo (proyectar_oficio,
// gestion_llamada, comite_caso).
//
//	PUT /api/v1/case-tasks/:id/complete
//	Body: { "userId": "icode-usuario", "formData": { ... } }
func (c *CaseTaskController) CompleteWithFormData(ctx *gin.Context) {
	id := ctx.Param("id")

	var body struct {
		UserID   string             `json:"userId"   binding:"required"`
		FormData datatypes.JSON     `json:"formData" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := c.svc.CompleteWithFormData(ctx.Request.Context(), id, body.UserID, body.FormData)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrCaseTaskNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": "tarea no encontrada"})
		case errors.Is(err, service.ErrEntityLetterInvalidState):
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusOK, task)
}
