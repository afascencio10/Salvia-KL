package controller

import (
	"bitsflow/common/utils"
	"bitsflow/internal/models"
	"bitsflow/salvia/service"
	"errors"
	"net/http"

	"github.com/gin-contrib/sessions"
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

// enlaceTerritorialRoles — únicos roles con acceso a los endpoints
// restringidos por departamento del Enlace Territorial.
var enlaceTerritorialRoles = map[string]bool{"en": true}

// requireEnlaceAccess valida la sesión (sessions.Default + utils.GetCommonSession,
// mismo patrón que entity_case_controller.go), exige rol "en" y que el usuario
// tenga un departamento asignado. No reutiliza salvia_config.PermissionsByRole
// porque la pantalla Detalle de Barrera (BarreraDetalleGET) hoy no tiene ningún
// permiso de rol asociado — el resto de /api/v1 en este repo no valida sesión
// en absoluto (hallazgo reportado aparte); estos endpoints no replican ese
// hueco a sabiendas. Devuelve la sesión resuelta, o nil si ya respondió.
// Usado por CreateGestionPropia (case_task_controller.go) y por el endpoint
// de "Barreras Departamento" (barrier_v2_gin_controller.go).
func requireEnlaceAccess(ctx *gin.Context) *utils.CommonSession {
	sessionID, ok := sessions.Default(ctx).Get("userData").(string)
	if !ok || sessionID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no autenticado"})
		return nil
	}
	s, err := utils.GetCommonSession(sessionID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no autenticado"})
		return nil
	}
	if !enlaceTerritorialRoles[s.CurrentRole] {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "su rol no tiene permiso para acceder a este recurso"})
		return nil
	}
	if s.AssignedDepartmentID == "" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "su usuario no tiene un departamento asignado"})
		return nil
	}
	return s
}

// RegisterRoutes registra las rutas de CaseTask en el grupo /api/v1.
//
//	GET  /api/v1/case-tasks?assignedUserId=<id>   → tareas con barrera asignadas al usuario
//	POST /api/v1/case-tasks/:id/complete           → marcar tarea como completada (flujo legacy "Gestionar")
//	PUT  /api/v1/case-tasks/:id/complete           → completar tarea con formData (case-task-modal)
//	POST /api/v1/case-tasks/gestion-propia         → Enlace registra gestión propia (crea+completa en un paso)
func (c *CaseTaskController) RegisterRoutes(rg *gin.RouterGroup) {
	group := rg.Group("/case-tasks")
	group.GET("", c.ListByAssignedUser)
	group.GET("/:id", c.GetByID)
	group.POST("/:id/complete", c.Complete)
	group.PUT("/:id/complete", c.CompleteWithFormData)
	group.POST("/:id/reassign", c.Reassign)
	group.POST("/gestion-propia", c.CreateGestionPropia)
}

// ListByAssignedUser devuelve las tareas con barrierId asignadas al usuario,
// enriquecidas con datos de barrier_v2 y victim_case.
// También soporta filtro por caseId.
//
// GetByID retorna una tarea por su ID, enriquecida con el nombre del usuario
// asignado (assignedUserName) — necesario para vistas de solo lectura como
// case-task-history, que muestran "Completado por" sin re-consultar aparte.
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
	ctx.JSON(http.StatusOK, c.enrichTaskWithName(*task))
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
		psychosocialID := ctx.Query("psychosocialId")
		if psychosocialID != "" {
			items, err := c.svc.ListByCaseIDAndPsychosocial(ctx.Request.Context(), caseID, psychosocialID)
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

// CompleteWithFormData completa una CaseTask guardando el JSON del formulario
// y ejecutando los efectos de lado propios del tipo (proyectar_oficio,
// gestion_llamada, comite_caso, Corregir oficio).
//
//	PUT /api/v1/case-tasks/:id/complete
//	Body: { "userId": "icode-usuario", "formData": { ... } }
func (c *CaseTaskController) CompleteWithFormData(ctx *gin.Context) {
	id := ctx.Param("id")

	var body struct {
		UserID   string         `json:"userId"   binding:"required"`
		FormData datatypes.JSON `json:"formData" binding:"required"`
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

// tiposGestionPropiaValidos — opciones válidas para el campo "tipo" del
// endpoint de gestión propia (mismas 4 opciones del modal ModalRegistrarGestion).
var tiposGestionPropiaValidos = map[string]bool{
	"Llamada":            true,
	"Visita presencial":  true,
	"Oficio a entidad":   true,
	"Otra gestión":       true,
}

// CreateGestionPropia crea una CaseTask ya completada para una gestión que el
// Enlace Territorial registró por iniciativa propia sobre una barrera — sin
// que exista una tarea "ToDo" previa. Solo el rol "en" puede llamarlo, y solo
// sobre barreras del departamento asignado al Enlace de la sesión.
//
//	POST /api/v1/case-tasks/gestion-propia
//	Body: { "caseId": "...", "barrierId": "...", "tipo": "Llamada", "descripcion": "..." }
func (c *CaseTaskController) CreateGestionPropia(ctx *gin.Context) {
	s := requireEnlaceAccess(ctx)
	if s == nil {
		return
	}

	var body struct {
		CaseID      string `json:"caseId" binding:"required"`
		BarrierID   string `json:"barrierId" binding:"required"`
		Tipo        string `json:"tipo" binding:"required"`
		Descripcion string `json:"descripcion" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "campos 'caseId', 'barrierId', 'tipo' y 'descripcion' requeridos"})
		return
	}
	if !tiposGestionPropiaValidos[body.Tipo] {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "tipo de gestión inválido: " + body.Tipo})
		return
	}

	task, err := c.svc.CreateGestionPropia(ctx.Request.Context(), body.CaseID, body.BarrierID, body.Tipo, body.Descripcion, s.UserICode, s.AssignedDepartmentID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrBarrierNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": "barrera no encontrada"})
		case errors.Is(err, service.ErrBarrierDepartmentMismatch):
			ctx.JSON(http.StatusForbidden, gin.H{"error": "esta barrera no pertenece a tu departamento asignado"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		}
		return
	}

	ctx.JSON(http.StatusOK, c.enrichTaskWithName(*task))
}

// lookupUserName resuelve el nombre completo de un usuario a partir de su icode.
// Devuelve "" si el icode está vacío o no se encuentra.
func (c *CaseTaskController) lookupUserName(userID string) string {
	if userID == "" {
		return ""
	}
	var name string
	c.svc.GetDB().Raw(`
		SELECT COALESCE(gup.general_user_profile_names, '') || ' ' || COALESCE(gup.general_user_profile_last_names, '')
		FROM security.general_user gu
		JOIN security.general_user_profile gup ON gup.general_user_profile_id = gu.general_user_general_user_profile
		WHERE gu.general_user_i_code = ?
	`, userID).Scan(&name)
	return name
}

// taskToJSON arma el gin.H completo de una tarea, incluyendo el nombre del
// usuario asignado ya resuelto. Centraliza el shape para que GetByID y el
// listado devuelvan siempre el mismo conjunto de campos.
func (c *CaseTaskController) taskToJSON(t models.CaseTask, assignedUserName string) gin.H {
	result := gin.H{
		"id":               t.ID,
		"caseId":           t.CaseID,
		"category":         t.Category,
		"type":             t.Type,
		"description":      t.Description,
		"status":           t.Status,
		"assignedUserId":   t.AssignedUserID,
		"assignedUserName": assignedUserName,
		"barrierId":        t.BarrierID,
		"entityLetterId":   t.EntityLetterID,
		"followUpId":       t.FollowUpID,
		"result":           t.Result,
		"formData":         t.FormData,
		"completedAt":      t.CompletedAt,
		"createdAt":        t.CreatedAt,
		"updatedAt":        t.UpdatedAt,
	}
	// Resolver sector de la barrera si tiene barrierId
	if t.BarrierID != nil && *t.BarrierID != "" {
		var sector string
		c.svc.GetDB().Raw("SELECT COALESCE(sector, '') FROM salvia.barrier_v2 WHERE id = ? LIMIT 1", *t.BarrierID).Scan(&sector)
		result["barrierSector"] = sector
	}
	return result
}

// enrichTaskWithName enriquece una tarea individual con el nombre del usuario
// asignado. Usado por GetByID — vistas de detalle como case-task-history
// necesitan "todos los datos" de una sola tarea, incluyendo el nombre legible.
func (c *CaseTaskController) enrichTaskWithName(t models.CaseTask) gin.H {
	return c.taskToJSON(t, c.lookupUserName(t.AssignedUserID))
}

// enrichTasksWithNames agrega el nombre del usuario asignado a cada tarea.
// Resuelve los nombres en batch (una consulta por icode único) para evitar
// N+1 queries cuando el listado tiene muchas tareas del mismo agente.
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
		if name := c.lookupUserName(id); name != "" {
			namesMap[id] = name
		}
	}

	// Construir response enriquecido
	var result []gin.H
	for _, t := range tasks {
		result = append(result, c.taskToJSON(t, namesMap[t.AssignedUserID]))
	}
	if result == nil {
		result = []gin.H{}
	}
	return result
}
