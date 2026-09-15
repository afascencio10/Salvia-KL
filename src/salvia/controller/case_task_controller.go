package controller

import (
	"bitsflow/common/utils"
	"bitsflow/internal/models"
	"bitsflow/salvia/service"
	"errors"
	"net/http"
	"strconv"

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

// misBarrerasRoles — roles con acceso a "Mis Barreras" (barreras asignadas
// personalmente al usuario vía user_territory_link_id). A diferencia de
// enlaceTerritorialRoles, no exige departamento asignado porque "no"
// (Operador territorial nacional) no tiene ámbito departamental.
var misBarrerasRoles = map[string]bool{"en": true, "no": true}

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
	s := requireSessionRole(ctx, enlaceTerritorialRoles)
	if s == nil {
		return nil
	}
	if s.AssignedDepartmentID == "" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "su usuario no tiene un departamento asignado"})
		return nil
	}
	return s
}

// requireMisBarrerasAccess valida la sesión y exige que el rol esté en
// misBarrerasRoles. Usado por ListMine (barrier_v2_gin_controller.go), que
// se filtra por s.UserICode y no depende del departamento asignado.
func requireMisBarrerasAccess(ctx *gin.Context) *utils.CommonSession {
	return requireSessionRole(ctx, misBarrerasRoles)
}

// requireSessionRole resuelve la sesión y exige que el rol actual esté entre los
// permitidos. Devuelve nil si ya respondió con 401/403.
func requireSessionRole(ctx *gin.Context, allowed map[string]bool) *utils.CommonSession {
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
	if !allowed[s.CurrentRole] {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "su rol no tiene permiso para acceder a este recurso"})
		return nil
	}
	return s
}

// misTareasRoles — roles con acceso a la bandeja de "Mis Tareas".
var misTareasRoles = map[string]bool{"op": true, "ro": true, "et": true, "sv": true, "no": true, "do": true, "us": true, "fo": true, "ad": true, "an": true, "ps": true, "ts": true, "ab": true, "mn": true, "te": true}

// equipoTecnicoRoles — único rol que puede agregar comentarios del Equipo
// Técnico en el tab Derivaciones de Detalle del Caso.
var equipoTecnicoRoles = map[string]bool{"te": true}

// RegisterRoutes registra las rutas de CaseTask en el grupo /api/v1.
//
//	GET  /api/v1/case-tasks?assignedUserId=<id>   → tareas con barrera asignadas al usuario
//	POST /api/v1/case-tasks/:id/complete           → marcar tarea como completada (flujo legacy "Gestionar")
//	PUT  /api/v1/case-tasks/:id/complete           → completar tarea con formData (case-task-modal)
//	POST /api/v1/case-tasks/gestion-propia         → Enlace registra gestión propia (crea+completa en un paso)
func (c *CaseTaskController) RegisterRoutes(rg *gin.RouterGroup) {
	group := rg.Group("/case-tasks")
	group.GET("", c.ListByAssignedUser)
	group.GET("/mine", c.ListMine)
	group.GET("/:id", c.GetByID)
	group.POST("/:id/complete", c.Complete)
	group.PUT("/:id/complete", c.CompleteWithFormData)
	group.POST("/:id/reassign", c.Reassign)
	group.POST("/gestion-propia", c.CreateGestionPropia)
	group.POST("/articulacion-territorial", c.SaveArticulacionTerritorial)
	group.POST("/gestion-mecanismo", c.SaveGestionMecanismo)
	group.POST("/comentario-equipo-tecnico", c.SaveComentarioEquipoTecnico)
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
//
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
			ctx.JSON(http.StatusOK, c.enrichTasksWithNames(c.filterTasksForEntityRole(ctx, items)))
			return
		}
		psychosocialID := ctx.Query("psychosocialId")
		if psychosocialID != "" {
			items, err := c.svc.ListByCaseIDAndPsychosocial(ctx.Request.Context(), caseID, psychosocialID)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
				return
			}
			ctx.JSON(http.StatusOK, c.enrichTasksWithNames(c.filterTasksForEntityRole(ctx, items)))
			return
		}
		items, err := c.svc.ListByCaseID(ctx.Request.Context(), caseID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
			return
		}
		ctx.JSON(http.StatusOK, c.enrichTasksWithNames(c.filterTasksForEntityRole(ctx, items)))
		return
	}

	items, err := c.svc.ListByAssignedUserIDWithRelations(ctx.Request.Context(), assignedUserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}

	ctx.JSON(http.StatusOK, items)
}

// optionalSession lee la sesión si existe; no responde 401 cuando falta
// (GET /case-tasks por caseId sigue abierto para otros roles).
func optionalSession(ctx *gin.Context) *utils.CommonSession {
	sessionID, ok := sessions.Default(ctx).Get("userData").(string)
	if !ok || sessionID == "" {
		return nil
	}
	s, err := utils.GetCommonSession(sessionID)
	if err != nil {
		return nil
	}
	return s
}

// filterTasksForEntityRole deja solo las asignadas al usuario si el rol de
// sesión es et. Otros roles y ausencia de sesión no filtran.
func (c *CaseTaskController) filterTasksForEntityRole(ctx *gin.Context, items []models.CaseTask) []models.CaseTask {
	s := optionalSession(ctx)
	if s == nil || s.CurrentRole != "et" {
		return items
	}
	out := make([]models.CaseTask, 0, len(items))
	for _, t := range items {
		if t.AssignedUserID == s.UserICode {
			out = append(out, t)
		}
	}
	return out
}

// ListMine devuelve, paginadas, las tareas del agente autenticado — pantalla
// "Mis Tareas". El usuario se resuelve siempre desde la sesión y nunca desde un
// parámetro del cliente, a diferencia de GET /api/v1/case-tasks?assignedUserId=.
//
//	GET /api/v1/case-tasks/mine?status=ToDo&page=1&pageSize=20
func (c *CaseTaskController) ListMine(ctx *gin.Context) {
	s := requireSessionRole(ctx, misTareasRoles)
	if s == nil {
		return
	}

	status := ctx.Query("status")
	if status != "" && status != models.CaseTaskStatusToDo && status != models.CaseTaskStatusDone {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "estado inválido: " + status})
		return
	}

	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("pageSize"))

	res, err := c.svc.ListByAssignedUserIDPaged(ctx.Request.Context(), s.UserICode, status, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"items":    res.Items,
		"total":    res.Total,
		"page":     res.Page,
		"pageSize": res.PageSize,
	})
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
		case errors.Is(err, service.ErrEntityUserNotFound), errors.Is(err, service.ErrEntityUserNotUnique):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
	"Llamada":           true,
	"Visita presencial": true,
	"Oficio a entidad":  true,
	"Otra gestión":      true,
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

// SaveArticulacionTerritorial completa o crea+completa una tarea de
// Articulación territorial (Mis Barreras y «Registrar gestión propia» en Detalle).
//
//	POST /api/v1/case-tasks/articulacion-territorial
//	Body: { "barrierId": "...", "taskId": "...", "formData": { ... } }
func (c *CaseTaskController) SaveArticulacionTerritorial(ctx *gin.Context) {
	s := requireEnlaceAccess(ctx)
	if s == nil {
		return
	}

	var body struct {
		BarrierID string         `json:"barrierId" binding:"required"`
		TaskID    string         `json:"taskId"`
		FormData  datatypes.JSON `json:"formData" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "campos 'barrierId' y 'formData' requeridos"})
		return
	}

	task, err := c.svc.SaveArticulacionTerritorial(ctx.Request.Context(), s.UserICode, s.AssignedDepartmentID, body.BarrierID, body.TaskID, body.FormData)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrBarrierNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": "barrera no encontrada"})
		case errors.Is(err, service.ErrBarrierLinkMismatch):
			ctx.JSON(http.StatusForbidden, gin.H{"error": "esta barrera no está asignada a tu usuario"})
		case errors.Is(err, service.ErrBarrierDepartmentMismatch):
			ctx.JSON(http.StatusForbidden, gin.H{"error": "esta barrera no pertenece a tu departamento asignado"})
		case errors.Is(err, service.ErrCaseTaskNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": "tarea no encontrada"})
		case errors.Is(err, service.ErrArticulacionInvalid):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "completa fecha de contacto, compromiso y seguimiento si aplica"})
		case errors.Is(err, service.ErrArticulacionTaskMismatch):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "la tarea no corresponde a esta barrera"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		}
		return
	}

	ctx.JSON(http.StatusOK, c.enrichTaskWithName(*task))
}

// SaveGestionMecanismo crea una case_task Done con el formulario del mecanismo.
//
//	POST /api/v1/case-tasks/gestion-mecanismo
//	Body: { "mechanismId": "...", "formData": { ... } }
func (c *CaseTaskController) SaveGestionMecanismo(ctx *gin.Context) {
	s := requireMcAccess(ctx)
	if s == nil {
		return
	}

	var body struct {
		MechanismID string         `json:"mechanismId" binding:"required"`
		FormData    datatypes.JSON `json:"formData" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "campos 'mechanismId' y 'formData' requeridos"})
		return
	}

	task, err := c.svc.SaveGestionMecanismo(ctx.Request.Context(), s.UserICode, body.MechanismID, body.FormData)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrMechanismNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": "mecanismo no encontrado"})
		case errors.Is(err, service.ErrMechanismNotAssigned):
			ctx.JSON(http.StatusForbidden, gin.H{"error": "este caso no está asignado a tu usuario"})
		case errors.Is(err, service.ErrMechanismNotActive):
			ctx.JSON(http.StatusConflict, gin.H{"error": "este escalamiento ya no está activo"})
		case errors.Is(err, service.ErrGestionMecanismoInvalid),
			errors.Is(err, service.ErrMechanismInvalidScope):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "completa los campos requeridos del formulario"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		}
		return
	}

	ctx.JSON(http.StatusOK, c.enrichTaskWithName(*task))
}

// SaveComentarioEquipoTecnico crea una case_task Done con un comentario del
// Equipo Técnico sobre un caso (tab Derivaciones → Comentarios del Equipo
// Técnico). Solo accesible para el rol 'te'.
//
//	POST /api/v1/case-tasks/comentario-equipo-tecnico
//	Body: { "caseId": "...", "texto": "..." }
func (c *CaseTaskController) SaveComentarioEquipoTecnico(ctx *gin.Context) {
	s := requireSessionRole(ctx, equipoTecnicoRoles)
	if s == nil {
		return
	}

	var body struct {
		CaseID string `json:"caseId" binding:"required"`
		Texto  string `json:"texto" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "campos 'caseId' y 'texto' requeridos"})
		return
	}

	task, err := c.svc.SaveComentarioEquipoTecnico(ctx.Request.Context(), s.UserICode, body.CaseID, body.Texto)
	if err != nil {
		if errors.Is(err, service.ErrComentarioEquipoTecnicoInvalid) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "el comentario no puede estar vacío"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}

	ctx.JSON(http.StatusCreated, c.enrichTaskWithName(*task))
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
func (c *CaseTaskController) taskToJSON(t models.CaseTask, assignedUserName, entityName string) gin.H {
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
		"entityCaseId":     t.EntityCaseID,
		"followUpId":       t.FollowUpID,
		"result":           t.Result,
		"formData":         t.FormData,
		"completedAt":      t.CompletedAt,
		"createdAt":        t.CreatedAt,
		"updatedAt":        t.UpdatedAt,
	}
	if entityName != "" {
		result["entityName"] = entityName
	}
	// Resolver sector y entidad de la barrera si tiene barrierId
	if t.BarrierID != nil && *t.BarrierID != "" {
		var barrierRow struct {
			Sector           string `gorm:"column:sector"`
			EntityBranchID   *int64 `gorm:"column:entity_branch_id"`
			EntityBranchName string `gorm:"column:entity_branch_name"`
		}
		c.svc.GetDB().Raw(`
			SELECT COALESCE(b.sector, '') AS sector,
			       b.entity_branch_id,
			       COALESCE(eb.entity_branch_name, '') AS entity_branch_name
			FROM salvia.barrier_v2 b
			LEFT JOIN salvia.entity_branch eb ON eb.entity_branch_id = b.entity_branch_id
			WHERE b.id = ? LIMIT 1`, *t.BarrierID).Scan(&barrierRow)
		result["barrierSector"] = barrierRow.Sector
		result["barrierEntityBranchId"] = barrierRow.EntityBranchID
		result["barrierEntityName"] = barrierRow.EntityBranchName
	}
	return result
}

func (c *CaseTaskController) lookupEntityNames(entityCaseIDs []string) map[string]string {
	out := make(map[string]string)
	if len(entityCaseIDs) == 0 {
		return out
	}
	type row struct {
		ID   string `gorm:"column:id"`
		Name string `gorm:"column:entity_name"`
	}
	var rows []row
	c.svc.GetDB().Raw(`
		SELECT ec.id, COALESCE(e.entity_name, '') AS entity_name
		FROM salvia.entity_case ec
		JOIN salvia.entity_branch eb ON eb.entity_branch_id = ec.entity_branch_id
		JOIN salvia.entity e ON e.entity_id = eb.entity_id
		WHERE ec.id IN ?
	`, entityCaseIDs).Scan(&rows)
	for _, r := range rows {
		out[r.ID] = r.Name
	}
	return out
}

func (c *CaseTaskController) entityNameForTask(t models.CaseTask) string {
	if t.EntityCaseID == nil || *t.EntityCaseID == "" {
		return ""
	}
	return c.lookupEntityNames([]string{*t.EntityCaseID})[*t.EntityCaseID]
}

// enrichTaskWithName enriquece una tarea individual con el nombre del usuario
// asignado. Usado por GetByID — vistas de detalle como case-task-history
// necesitan "todos los datos" de una sola tarea, incluyendo el nombre legible.
func (c *CaseTaskController) enrichTaskWithName(t models.CaseTask) gin.H {
	return c.taskToJSON(t, c.lookupUserName(t.AssignedUserID), c.entityNameForTask(t))
}

// enrichTasksWithNames agrega el nombre del usuario asignado a cada tarea.
// Resuelve los nombres en batch (una consulta por icode único) para evitar
// N+1 queries cuando el listado tiene muchas tareas del mismo agente.
func (c *CaseTaskController) enrichTasksWithNames(tasks []models.CaseTask) []gin.H {
	// Recolectar IDs únicos
	idsMap := make(map[string]bool)
	ecIDsMap := make(map[string]bool)
	for _, t := range tasks {
		if t.AssignedUserID != "" {
			idsMap[t.AssignedUserID] = true
		}
		if t.EntityCaseID != nil && *t.EntityCaseID != "" {
			ecIDsMap[*t.EntityCaseID] = true
		}
	}

	// Resolver nombres
	namesMap := make(map[string]string)
	for id := range idsMap {
		if name := c.lookupUserName(id); name != "" {
			namesMap[id] = name
		}
	}

	ecIDs := make([]string, 0, len(ecIDsMap))
	for id := range ecIDsMap {
		ecIDs = append(ecIDs, id)
	}
	entityNames := c.lookupEntityNames(ecIDs)

	// Construir response enriquecido
	var result []gin.H
	for _, t := range tasks {
		entityName := ""
		if t.EntityCaseID != nil {
			entityName = entityNames[*t.EntityCaseID]
		}
		result = append(result, c.taskToJSON(t, namesMap[t.AssignedUserID], entityName))
	}
	if result == nil {
		result = []gin.H{}
	}
	return result
}
