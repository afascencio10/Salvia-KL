// Package controller expone los endpoints HTTP de FollowUpV2 usando Gin.
package controller

import (
	"bitsflow/common/utils"
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"bitsflow/salvia/service"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// FollowUpV2Controller maneja los endpoints del calendario de seguimientos (HU-027).
type FollowUpV2Controller struct {
	svc service.FollowUpV2Service
}

// NewFollowUpV2Controller construye el controlador inyectando el servicio.
func NewFollowUpV2Controller(svc service.FollowUpV2Service) *FollowUpV2Controller {
	return &FollowUpV2Controller{svc: svc}
}

// RegisterRoutes registra todas las rutas de FollowUpV2 en el grupo /api/v1.
//
//		GET  /api/v1/cases/:victim_case_id/follow-ups
//		POST /api/v1/cases/:victim_case_id/follow-ups/generate
//		GET  /api/v1/cases/:victim_case_id/follow-ups/by-id?id=...
//		GET  /api/v1/cases/:victim_case_id/follow-ups/list?page=...&limit=...
//	 GET  /api/v1/follow-ups/my-day
//	 POST /api/v1/follow-ups/:id/attempts
func (c *FollowUpV2Controller) RegisterRoutes(api *gin.RouterGroup) {
	followUps := api.Group("/cases/:victim_case_id/follow-ups")
	{
		followUps.GET("", c.GetCalendar)
		followUps.POST("/generate", c.GenerateCalendar)
		followUps.GET("/by-id", c.GetByID)
		followUps.GET("/list", c.List)
	}

	api.GET("/cases/follow-ups/detail", c.GetDetail)

	// Hacer seguimiento
	api.GET("/follow-ups/:id/load", c.LoadFollowUp)

	// Mis Seguimientos
	myDay := api.Group("/follow-ups")
	{
		myDay.GET("/my-day", c.GetMyDayFollowUps)
		myDay.POST("/:id/attempts", c.RegisterAttempt)
		myDay.PUT("/:id/close", c.CloseCase)
	}
	// Seguimientos Área — rutas para supervisores
	seg := api.Group("/seguimientos")
	{
		seg.GET("", c.ListByArea)
		seg.GET("/carga-agentes", c.AgentWorkload)
		seg.GET("/filtros-opciones", c.FilterOptions)
		seg.PUT("/:id/reagendar", c.Reschedule)
	}
}

// GetCalendar godoc
//
//	@Summary		Consultar calendario de seguimientos de un caso
//	@Description	Retorna todos los seguimientos activos de un caso ordenados por fecha
//	@Tags			Seguimientos
//	@Produce		json
//	@Param			victim_case_id	path		string				true	"ID del caso"
//	@Success		200				{array}		models.FollowUpV2	"Lista de seguimientos"
//	@Failure		404				{object}	map[string]string	"No hay seguimientos para el caso"
//	@Failure		500				{object}	map[string]string	"Error interno"
//	@Router			/cases/{victim_case_id}/follow-ups [get]
func (c *FollowUpV2Controller) GetCalendar(ctx *gin.Context) {
	caseID := ctx.Param("victim_case_id")

	items, err := c.svc.GetByCaseID(ctx.Request.Context(), caseID)
	if err != nil {
		if errors.Is(err, service.ErrFollowUpCaseEmpty) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "no se encontraron seguimientos para el caso"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}
	ctx.JSON(http.StatusOK, items)
}

// GenerateCalendar godoc
//
//	@Summary		Generar o recalcular calendario de seguimientos
//	@Description	Crea los seguimientos según la matriz de riesgo. Si ya existen y cambió el nivel, reprograma los pendientes.
//	@Tags			Seguimientos
//	@Accept			json
//	@Produce		json
//	@Param			victim_case_id	path		string							true	"ID del caso"
//	@Param			body			body		service.GenerateCalendarInput	true	"Datos para generar el calendario"
//	@Success		201				{array}		models.FollowUpV2				"Seguimientos creados"
//	@Failure		400				{object}	map[string]string				"Body inválido o risk_level fuera de rango"
//	@Failure		500				{object}	map[string]string				"Error interno"
//	@Router			/cases/{victim_case_id}/follow-ups/generate [post]
func (c *FollowUpV2Controller) GenerateCalendar(ctx *gin.Context) {
	caseID := ctx.Param("victim_case_id")

	var input service.GenerateCalendarInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created, err := c.svc.GenerateOrRecalculate(ctx.Request.Context(), caseID, input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, created)
}

// GetByID godoc
//
//	@Summary		Obtener seguimiento por ID
//	@Tags			Seguimientos
//	@Produce		json
//	@Param			victim_case_id	path		string				true	"ID del caso"
//	@Param			id				query		string				true	"UUID del seguimiento"
//	@Success		200				{object}	models.FollowUpV2	"Seguimiento encontrado"
//	@Failure		400				{object}	map[string]string	"Parámetro id faltante"
//	@Failure		404				{object}	map[string]string	"No encontrado"
//	@Router			/cases/{victim_case_id}/follow-ups/by-id [get]
func (c *FollowUpV2Controller) GetByID(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "parámetro 'id' requerido"})
		return
	}

	fu, err := c.svc.GetFollowUpByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrFollowUpNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "registro no encontrado"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}
	ctx.JSON(http.StatusOK, fu)
}

// List godoc
//
//	@Summary		Listar seguimientos paginados
//	@Tags			Seguimientos
//	@Produce		json
//	@Param			victim_case_id	path		string									true	"ID del caso"
//	@Param			page			query		int										false	"Página (base 0)"
//	@Param			limit			query		int										false	"Registros por página (default 20)"
//	@Success		200				{object}	repository.PageResult[models.FollowUpV2]	"Página de resultados"
//	@Failure		500				{object}	map[string]string						"Error interno"
//	@Router			/cases/{victim_case_id}/follow-ups/list [get]
func (c *FollowUpV2Controller) List(ctx *gin.Context) {
	page := ginQueryInt(ctx, "page", 0)
	limit := ginQueryInt(ctx, "limit", 20)

	result, err := c.svc.GetPaginatedFollowUps(ctx.Request.Context(), page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

// GetDetail godoc
//
//	@Summary		Obtener el detalle estructurado de un seguimiento
//	@Tags			Seguimientos
//	@Produce		json
//	@Param			victim_case_id	path		string				true	"ID del caso"
//	@Param			id				query		string				true	"UUID del seguimiento"
//	@Success		200				{object}	models.FollowUpDetailResponse	"Detalle del seguimiento"
//	@Failure		400				{object}	map[string]string	"Parámetro id faltante"
//	@Failure		404				{object}	map[string]string	"No encontrado"
//	@Router			/cases/{victim_case_id}/follow-ups/detail [get]
func (c *FollowUpV2Controller) GetDetail(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "parámetro 'id' requerido"})
		return
	}

	// Por ahora asignamos isSupervisor = true. En el futuro, integrar middleware de roles (ej. GetCommonSession)
	isSupervisor := true

	detail, err := c.svc.GetFollowUpDetail(ctx.Request.Context(), id, isSupervisor)
	if err != nil {
		if errors.Is(err, service.ErrFollowUpNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "seguimiento no encontrado"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno al cargar detalle"})
		return
	}
	ctx.JSON(http.StatusOK, detail)
}

// GetMyDayFollowUps godoc
//
//	@Summary		Obtener seguimientos del día para un agente
//	@Description	Retorna los seguimientos pendientes, priorizados y realizados del agente logueado para el día actual
//	@Tags			Seguimientos
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}	"Seguimientos del día"
//	@Failure		500	{object}	map[string]string	"Error interno"
//	@Router			/follow-ups/my-day [get]
func (c *FollowUpV2Controller) GetMyDayFollowUps(ctx *gin.Context) {
	// Obtener agentID de la sesión del usuario autenticado
	session := sessions.Default(ctx)
	sessionID, ok := session.Get("userData").(string)
	if !ok || sessionID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "sesión inválida"})
		return
	}

	s, err := utils.GetCommonSession(sessionID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "sesión expirada"})
		return
	}

	// Usar el UserICode de la sesión como agentID
	agentID := s.UserICode	
	//log agent id
	fmt.Println("agentID", agentID)
	
	if agentID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "usuario sin identificador"})
		return
	}

	date := time.Now()

	// Obtener seguimientos enriquecidos con datos del caso
	response, err := c.svc.GetMyDayFollowUpsEnriched(ctx.Request.Context(), agentID, date)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno al obtener seguimientos"})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// RegisterAttempt godoc
//
//	@Summary		Registrar un intento fallido de contacto
//	@Description	Registra un intento fallido con su justificación. Incrementa el contador de intentos.
//	@Tags			Seguimientos
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"UUID del seguimiento"
//	@Param			body	body		RegisterAttemptInput	true	"Razón del intento fallido"
//	@Success		200		{object}	models.FollowUpV2	"Seguimiento actualizado"
//	@Failure		400		{object}	map[string]string	"Body inválido o razón vacía"
//	@Failure		404		{object}	map[string]string	"Seguimiento no encontrado"
//	@Failure		500		{object}	map[string]string	"Error interno o máximo de intentos alcanzado"
//	@Router			/follow-ups/{id}/attempts [post]
func (c *FollowUpV2Controller) RegisterAttempt(ctx *gin.Context) {
	followUpID := ctx.Param("id")

	var body struct {
		Reason      string `json:"reason" binding:"required"`
		WasAnswered bool   `json:"was_answered"`
	}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "razón es requerida"})
		return
	}

	fu, err := c.svc.RegisterContactAttempt(ctx.Request.Context(), followUpID, body.Reason, body.WasAnswered)
	if err != nil {
		if errors.Is(err, service.ErrFollowUpNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "seguimiento no encontrado"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Retornar información adicional sobre el estado
	response := gin.H{
		"success":                 true,
		"message":                 "Intento registrado con éxito",
		"followUp":                fu,
		"attempts":                fu.Attempts,
		"maxAttemptsReached":      fu.Attempts >= 3,
		"criticalAttemptsReached": fu.Attempts >= 9,
	}

	ctx.JSON(http.StatusOK, response)
}

// RegisterAttemptInput es el body de entrada para registrar un intento
type RegisterAttemptInput struct {
	Reason      string `json:"reason" binding:"required"`
	WasAnswered bool   `json:"was_answered"`
}

// ginQueryInt está definido en form_controller.go (mismo package controller).
// models importado para las anotaciones Swagger.
var _ models.FollowUpV2

// ── Seguimientos Área — Handlers ──────────────────────────────────────────────

// ListByArea godoc
//
//	@Summary		Listar seguimientos del área (supervisores)
//	@Description	Retorna seguimientos filtrados por el team del supervisor en sesión
//	@Tags			Seguimientos Área
//	@Produce		json
//	@Param			tab				query	string	false	"Tab: pendientes, realizados, vencidos, todos"
//	@Param			fecha_inicio	query	string	false	"Fecha inicio YYYY-MM-DD"
//	@Param			fecha_fin		query	string	false	"Fecha fin YYYY-MM-DD"
//	@Param			estado			query	string	false	"Estado del seguimiento"
//	@Param			agente_id		query	string	false	"ID del agente"
//	@Param			page			query	int		false	"Página (base 0)"
//	@Param			limit			query	int		false	"Registros por página"
//	@Success		200	{object}	map[string]interface{}
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/seguimientos [get]
func (c *FollowUpV2Controller) ListByArea(ctx *gin.Context) {
	team := ctx.GetHeader("X-User-Team")
	if team == "" {
		team = ctx.Query("team")
	}
	if team == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "team requerido (header X-User-Team o query param)"})
		return
	}

	filters := repository.FollowUpFilters{
		Tab:         ctx.Query("tab"),
		FechaInicio: ctx.Query("fecha_inicio"),
		FechaFin:    ctx.Query("fecha_fin"),
		Estado:      ctx.Query("estado"),
		AgentID:     ctx.Query("agente_id"),
	}
	page := ginQueryInt(ctx, "page", 0)
	limit := ginQueryInt(ctx, "limit", 20)

	items, total, err := c.svc.GetByTeamPaginated(ctx.Request.Context(), team, filters, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"items": items,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// AgentWorkload godoc
//
//	@Summary		Carga de seguimientos por agente
//	@Description	Retorna seguimientos pendientes agrupados por agente para una fecha
//	@Tags			Seguimientos Área
//	@Produce		json
//	@Param			fecha	query	string	false	"Fecha YYYY-MM-DD"
//	@Success		200		{array}	repository.AgentWorkload
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/seguimientos/carga-agentes [get]
func (c *FollowUpV2Controller) AgentWorkload(ctx *gin.Context) {
	team := ctx.GetHeader("X-User-Team")
	if team == "" {
		team = ctx.Query("team")
	}
	if team == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "team requerido"})
		return
	}

	fecha := ctx.Query("fecha")
	results, err := c.svc.GetAgentWorkload(ctx.Request.Context(), team, fecha)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}
	ctx.JSON(http.StatusOK, results)
}

// FilterOptions godoc
//
//	@Summary		Opciones de filtros para seguimientos del área
//	@Description	Retorna catálogos de agentes del área y estados disponibles
//	@Tags			Seguimientos Área
//	@Produce		json
//	@Success		200	{object}	service.FilterOptions
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/seguimientos/filtros-opciones [get]
func (c *FollowUpV2Controller) FilterOptions(ctx *gin.Context) {
	team := ctx.GetHeader("X-User-Team")
	if team == "" {
		team = ctx.Query("team")
	}
	if team == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "team requerido"})
		return
	}

	options, err := c.svc.GetFilterOptions(ctx.Request.Context(), team)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}
	ctx.JSON(http.StatusOK, options)
}

// Reschedule godoc
//
//	@Summary		Reagendar un seguimiento
//	@Description	Actualiza la fecha, hora, prioridad y motivo de un seguimiento
//	@Tags			Seguimientos Área
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string					true	"UUID del seguimiento"
//	@Param			body	body	service.RescheduleInput	true	"Datos de reagendamiento"
//	@Success		200		{object}	map[string]string
//	@Failure		400		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/seguimientos/{id}/reagendar [put]
func (c *FollowUpV2Controller) Reschedule(ctx *gin.Context) {
	id := ctx.Param("id")

	var input service.RescheduleInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.svc.RescheduleFollowUp(ctx.Request.Context(), id, input)
	if err != nil {
		if errors.Is(err, service.ErrFollowUpNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "seguimiento no encontrado"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "seguimiento reagendado exitosamente"})
}

// CloseCase godoc
//
//	@Summary		Cerrar seguimiento y futuros del mismo caso
//	@Tags			Seguimientos
//	@Produce		json
//	@Param			id	path		string				true	"UUID del seguimiento"
//	@Success		200	{object}	map[string]string	"Cierre exitoso"
//	@Failure		404	{object}	map[string]string	"No encontrado"
//	@Failure		500	{object}	map[string]string	"Error interno"
//	@Router			/follow-ups/{id}/close [put]
func (c *FollowUpV2Controller) CloseCase(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.svc.CloseCaseFollowUps(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "seguimiento no encontrado"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno al cerrar seguimientos"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "seguimientos cerrados exitosamente"})
}

// LoadFollowUp godoc
//
//	@Summary		Cargar pantalla hacer-seguimiento
//	@Description	Valida acceso, carga info del caso y crea FormSubmission si no existe
//	@Tags			Seguimientos
//	@Produce		json
//	@Param			id			path	string	true	"UUID del seguimiento"
//	@Param			agent_id	query	string	true	"ICode del agente"
//	@Param			form_id		query	string	true	"ID del formulario"
//	@Success		200	{object}	service.LoadFollowUpResult
//	@Failure		403	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/follow-ups/{id}/load [get]
func (c *FollowUpV2Controller) LoadFollowUp(ctx *gin.Context) {
	id := ctx.Param("id")
	agentID := ctx.Query("agent_id")
	formID := ctx.Query("form_id")

	if agentID == "" || formID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "agent_id y form_id son requeridos"})
		return
	}

	result, err := c.svc.LoadFollowUp(ctx.Request.Context(), id, agentID, formID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrFollowUpNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": "seguimiento no encontrado"})
		case errors.Is(err, service.ErrFollowUpNotAssigned):
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrFollowUpNotYetDue):
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		}
		return
	}
	ctx.JSON(http.StatusOK, result)
}
