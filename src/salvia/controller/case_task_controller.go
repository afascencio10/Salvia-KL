package controller

import (
	"bitsflow/internal/models"
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
	group.GET("/:id", c.GetByID)
	group.POST("/:id/complete", c.Complete)
	group.POST("/:id/reassign", c.Reassign)
}

// ListByAssignedUser devuelve las tareas con barrierId asignadas al usuario,
// enriquecidas con datos de barrier_v2 y victim_case.
// También soporta filtro por caseId.
//
// GetByID retorna una tarea por su ID.
// GET /api/v1/case-tasks/:id
func (c *CaseTaskController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id requerido"})
		return
	}
	task, err := c.svc.GetByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "tarea no encontrada"})
		return
	}
	ctx.JSON(http.StatusOK, task)
}

// ListByAssignedUser devuelve las tareas asignadas a un usuario o las de un caso.
//	GET /api/v1/case-tasks?assignedUserId=<userICode>
//	GET /api/v1/case-tasks?caseId=<caseICode>
func (c *CaseTaskController) ListByAssignedUser(ctx *gin.Context) {
	assignedUserID := ctx.Query("assignedUserId")
	caseID := ctx.Query("caseId")

	if assignedUserID == "" && caseID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "parámetro 'assignedUserId' o 'caseId' requerido"})
		return
	}

	if caseID != "" {
		barrierID := ctx.Query("barrierId")
		if barrierID != "" {
			items, err := c.svc.ListByCaseIDAndBarrier(ctx.Request.Context(), caseID, barrierID)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
				return
			}
			ctx.JSON(http.StatusOK, c.enrichTasksWithNames(items))
			return
		}
		items, err := c.svc.ListByCaseID(ctx.Request.Context(), caseID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
			return
		}
		ctx.JSON(http.StatusOK, c.enrichTasksWithNames(items))
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
// Si la tarea tiene EntityLetter asociado, también actualiza el oficio.
//
//	POST /api/v1/case-tasks/:id/complete
//	Body: { "result": "...", "nivel": "municipal", "entidad": "...", "kofaxPath": "...", "prioridad": "normal" }
func (c *CaseTaskController) Complete(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id requerido"})
		return
	}

	var body struct {
		Result    string `json:"result" binding:"required"`
		Nivel     string `json:"nivel"`
		Entidad   string `json:"entidad"`
		KofaxPath string `json:"kofaxPath"`
		Prioridad string `json:"prioridad"`
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

	// Si la tarea tenía EntityLetter y se enviaron datos de oficio, actualizar el EntityLetter
	if body.KofaxPath != "" || body.Entidad != "" {
		task, _ := c.svc.GetByID(ctx.Request.Context(), id)
		if task != nil && task.EntityLetterID != nil && *task.EntityLetterID != "" {
			updates := map[string]interface{}{
				"state": "para_revisar",
			}
			if body.KofaxPath != "" {
				updates["url_kofax"] = body.KofaxPath
			}
			if body.Entidad != "" {
				updates["entidad"] = body.Entidad
			}
			if body.Nivel != "" {
				updates["nivel"] = body.Nivel
			}
			if body.Prioridad != "" {
				updates["priority"] = body.Prioridad
			}
			c.svc.UpdateEntityLetter(ctx.Request.Context(), *task.EntityLetterID, updates)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "tarea completada exitosamente"})
}


// Reassign reasigna una CaseTask a otro usuario.
// POST /api/v1/case-tasks/:id/reassign
// Body: { "assignedUserId": "<icode>" }
func (c *CaseTaskController) Reassign(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id requerido"})
		return
	}

	var body struct {
		AssignedUserID string `json:"assignedUserId" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "campo 'assignedUserId' requerido"})
		return
	}

	if err := c.svc.Reassign(ctx.Request.Context(), id, body.AssignedUserID); err != nil {
		if err == service.ErrCaseTaskNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "tarea no encontrada"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "tarea reasignada exitosamente"})
}


// enrichTasksWithNames agrega el nombre del usuario asignado a cada tarea.
func (c *CaseTaskController) enrichTasksWithNames(tasks []models.CaseTask) []gin.H {
	// Recolectar IDs únicos
	idsMap := make(map[string]bool)
	for _, t := range tasks {
		if t.AssignedUserID != "" {
			idsMap[t.AssignedUserID] = true
		}
	}

	// Resolver nombres
	namesMap := make(map[string]string)
	for id := range idsMap {
		var name string
		c.svc.GetDB().Raw(`
			SELECT COALESCE(gup.general_user_profile_names, '') || ' ' || COALESCE(gup.general_user_profile_last_names, '')
			FROM security.general_user gu
			JOIN security.general_user_profile gup ON gup.general_user_profile_id = gu.general_user_general_user_profile
			WHERE gu.general_user_i_code = ?
		`, id).Scan(&name)
		if name != "" {
			namesMap[id] = name
		}
	}

	// Construir response enriquecido
	var result []gin.H
	for _, t := range tasks {
		item := gin.H{
			"id":               t.ID,
			"caseId":           t.CaseID,
			"category":         t.Category,
			"type":             t.Type,
			"description":      t.Description,
			"status":           t.Status,
			"assignedUserId":   t.AssignedUserID,
			"assignedUserName": namesMap[t.AssignedUserID],
			"barrierId":        t.BarrierID,
			"entityLetterId":   t.EntityLetterID,
			"followUpId":       t.FollowUpID,
			"result":           t.Result,
			"completedAt":      t.CompletedAt,
			"createdAt":        t.CreatedAt,
		}
		result = append(result, item)
	}
	if result == nil {
		result = []gin.H{}
	}
	return result
}
