// Package controller expone los endpoints HTTP de EntityLetter usando Gin.
package controller

import (
	"bitsflow/salvia/service"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// EntityLetterController maneja las rutas REST del recurso EntityLetter (oficios).
type EntityLetterController struct {
	svc service.EntityLetterService
}

func NewEntityLetterController(svc service.EntityLetterService) *EntityLetterController {
	return &EntityLetterController{svc: svc}
}

// RegisterRoutes registra todas las rutas de EntityLetter en el grupo /api/v1.
//
//	GET    /api/v1/entity-letters                      → List (paginado; filtros: state, caseId, barrierId, agentId, notificationUserId)
//	GET    /api/v1/entity-letters/:id                  → GetByID
//	POST   /api/v1/entity-letters                      → Create
//	PUT    /api/v1/entity-letters/:id                  → Update (campos opcionales)
//	DELETE /api/v1/entity-letters/:id                  → Delete (soft delete)
//	PUT    /api/v1/entity-letters/:id/state            → UpdateState (transición de estado)
//	GET    /api/v1/entity-letters/case/:caseId         → ListByCase
//	GET    /api/v1/entity-letters/barrier/:barrierId   → ListByBarrier
func (c *EntityLetterController) RegisterRoutes(rg *gin.RouterGroup) {
	letters := rg.Group("/entity-letters")

	letters.GET("", c.List)
	letters.POST("", c.Create)

	// Rutas de colecciones con prefijo fijo (deben ir ANTES de /:id)
	letters.GET("/case/:caseId", c.ListByCase)
	letters.GET("/barrier/:barrierId", c.ListByBarrier)

	// Rutas con parámetro dinámico
	letters.GET("/:id", c.GetByID)
	letters.PUT("/:id", c.Update)
	letters.DELETE("/:id", c.Delete)
	letters.PUT("/:id/state", c.UpdateState)
	letters.PUT("/:id/action", c.Action)
}

// ─── List ─────────────────────────────────────────────────────────────────────

// List devuelve oficios paginados con filtros opcionales por query string.
//
//	GET /api/v1/entity-letters?page=0&limit=20&state=para_revisar&caseId=...&agentId=...&notificationUserId=...
func (c *EntityLetterController) List(ctx *gin.Context) {
	page  := ginQueryInt(ctx, "page", 0)
	limit := ginQueryInt(ctx, "limit", 20)

	state              := ctx.Query("state")
	caseID             := ctx.Query("caseId")
	barrierID          := ctx.Query("barrierId")
	agentID            := ctx.Query("agentId")
	notificationUserID := ctx.Query("notificationUserId")

	switch {
	case state != "":
		result, err := c.svc.ListByState(ctx.Request.Context(), state, page, limit)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
			return
		}
		ctx.JSON(http.StatusOK, result)

	case caseID != "":
		items, err := c.svc.ListByCase(ctx.Request.Context(), caseID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
			return
		}
		ctx.JSON(http.StatusOK, items)

	case barrierID != "":
		items, err := c.svc.ListByBarrier(ctx.Request.Context(), barrierID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
			return
		}
		ctx.JSON(http.StatusOK, items)

	case agentID != "":
		// Devuelve la lista enriquecida con datos de barrier_v2 y victim_case
		items, err := c.svc.ListByAgentWithRelations(ctx.Request.Context(), agentID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
			return
		}
		ctx.JSON(http.StatusOK, items)

	case notificationUserID != "":
		// Devuelve la lista enriquecida con datos de barrier_v2 y victim_case
		items, err := c.svc.ListByNotificationUserWithRelations(ctx.Request.Context(), notificationUserID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
			return
		}
		ctx.JSON(http.StatusOK, items)

	default:
		result, err := c.svc.List(ctx.Request.Context(), page, limit)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
			return
		}
		ctx.JSON(http.StatusOK, result)
	}
}

// ─── GetByID ──────────────────────────────────────────────────────────────────

// GetByID devuelve un oficio por su UUID.
//
//	GET /api/v1/entity-letters/:id
func (c *EntityLetterController) GetByID(ctx *gin.Context) {
	letter, err := c.svc.GetByID(ctx.Request.Context(), ctx.Param("id"))
	if err != nil {
		if errors.Is(err, service.ErrEntityLetterNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "oficio no encontrado"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}
	ctx.JSON(http.StatusOK, letter)
}

// ─── Create ───────────────────────────────────────────────────────────────────

// Create crea un nuevo oficio en estado inicial "por_proyectar".
//
//	POST /api/v1/entity-letters
//	Body: { "barrierId": "...", "caseId": "...", "agentId": "...", "notificationUserId": "..." }
func (c *EntityLetterController) Create(ctx *gin.Context) {
	var body struct {
		BarrierID          string  `json:"barrierId"          binding:"required"`
		CaseID             string  `json:"caseId"             binding:"required"`
		AgentID            *string `json:"agentId"`
		NotificationUserID *string `json:"notificationUserId"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	letter, err := c.svc.Create(ctx.Request.Context(), service.CreateEntityLetterInput{
		BarrierID:          body.BarrierID,
		CaseID:             body.CaseID,
		AgentID:            body.AgentID,
		NotificationUserID: body.NotificationUserID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}
	ctx.JSON(http.StatusCreated, letter)
}

// ─── Update ───────────────────────────────────────────────────────────────────

// Update actualiza los campos opcionales de un oficio (sin cambiar el estado).
// Solo se actualizan los campos que vienen en el body (patch semántico).
//
//	PUT /api/v1/entity-letters/:id
//	Body: { "agentId": "...", "notificationUserId": "...", "reviewBy": "...", "radicadoBy": "...", "registerBy": "..." }
func (c *EntityLetterController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var body struct {
		AgentID            *string `json:"agentId"`
		NotificationUserID *string `json:"notificationUserId"`
		ReviewBy           *string `json:"reviewBy"`
		RadicadoBy         *string `json:"radicadoBy"`
		RegisterBy         *string `json:"registerBy"`
		Entidad            *string `json:"entidad"`
		Nivel              *string `json:"nivel"`
		UrlKofax           *string `json:"urlKofax"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	letter, err := c.svc.Update(ctx.Request.Context(), id, service.UpdateEntityLetterInput{
		AgentID:            body.AgentID,
		NotificationUserID: body.NotificationUserID,
		ReviewBy:           body.ReviewBy,
		RadicadoBy:         body.RadicadoBy,
		RegisterBy:         body.RegisterBy,
		Entidad:            body.Entidad,
		Nivel:              body.Nivel,
		UrlKofax:           body.UrlKofax,
	})
	if err != nil {
		if errors.Is(err, service.ErrEntityLetterNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "oficio no encontrado"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}
	ctx.JSON(http.StatusOK, letter)
}

// ─── Delete ───────────────────────────────────────────────────────────────────

// Delete realiza un soft-delete del oficio.
//
//	DELETE /api/v1/entity-letters/:id
func (c *EntityLetterController) Delete(ctx *gin.Context) {
	if err := c.svc.Delete(ctx.Request.Context(), ctx.Param("id")); err != nil {
		if errors.Is(err, service.ErrEntityLetterNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "oficio no encontrado"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}

// ─── UpdateState ──────────────────────────────────────────────────────────────

// UpdateState ejecuta una transición de estado validada.
// Si la transición no está permitida devuelve 422 con el error de dominio.
//
//	PUT /api/v1/entity-letters/:id/state
//	Body: { "state": "para_revisar", "userId": "uuid-del-agente" }
func (c *EntityLetterController) UpdateState(ctx *gin.Context) {
	id := ctx.Param("id")

	var body struct {
		State  string `json:"state"  binding:"required"`
		UserID string `json:"userId"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	letter, err := c.svc.UpdateState(ctx.Request.Context(), id, service.UpdateStateInput{
		State:  body.State,
		UserID: body.UserID,
	})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEntityLetterNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": "oficio no encontrado"})
		case errors.Is(err, service.ErrEntityLetterInvalidState):
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		}
		return
	}
	ctx.JSON(http.StatusOK, letter)
}

// ─── ListByCase ───────────────────────────────────────────────────────────────

// ListByCase devuelve todos los oficios de un caso específico.
//
//	GET /api/v1/entity-letters/case/:caseId
func (c *EntityLetterController) ListByCase(ctx *gin.Context) {
	items, err := c.svc.ListByCase(ctx.Request.Context(), ctx.Param("caseId"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}
	ctx.JSON(http.StatusOK, items)
}

// ─── ListByBarrier ────────────────────────────────────────────────────────────

// ListByBarrier devuelve todos los oficios de una barrera específica.
//
//	GET /api/v1/entity-letters/barrier/:barrierId
func (c *EntityLetterController) ListByBarrier(ctx *gin.Context) {
	items, err := c.svc.ListByBarrier(ctx.Request.Context(), ctx.Param("barrierId"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}
	ctx.JSON(http.StatusOK, items)
}

// ─── Action ───────────────────────────────────────────────────────────────────

// Action ejecuta la acción de un modal de gestión (actualiza campos + transición de estado).
// Es el endpoint unificado para todos los modales; el campo "action" determina qué lógica aplicar.
//
//	PUT /api/v1/entity-letters/:id/action
//	Body proyectar: { "action": "proyectar", "userId": "...", "nivel": "...", "entidad": "...", "urlKofax": "..." }
func (c *EntityLetterController) Action(ctx *gin.Context) {
	id := ctx.Param("id")

	var body struct {
		Action             string  `json:"action"              binding:"required"`
		UserID             string  `json:"userId"`
		Nivel              *string `json:"nivel"`
		Entidad            *string `json:"entidad"`
		UrlKofax           *string `json:"urlKofax"`
		Priority           *string `json:"priority"`
		// Campos nuevos del modal proyectar v2
		EntityBranchID     *int64  `json:"entityBranchId"`
		EntityName         *string `json:"entityName"`
		DepartmentID       *string `json:"departmentId"`
		CityID             *string `json:"cityId"`
		TownID             *string `json:"townId"`
		OfficialDependency *string `json:"officialDependency"`
		Subject            *string `json:"subject"`
		// Campos radicar / revisar / respuesta
		AsuntoRadicado     *string `json:"asuntoRadicado"`
		CorreoEntidad      *string `json:"correoEntidad"`
		NumeroRadicado     *string `json:"numeroRadicado"`
		ReasonCorrection   *string `json:"reasonCorrection"`
		ResponseDate       *string `json:"responseDate"`
		CorreoRemitente    *string `json:"correoRemitente"`
		AsuntoRespuesta    *string `json:"asuntoRespuesta"`
		ResponseReviewBy   *string `json:"responseReviewBy"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	letter, err := c.svc.PerformAction(ctx.Request.Context(), id, service.ActionInput{
		Action:             body.Action,
		UserID:             body.UserID,
		Nivel:              body.Nivel,
		Entidad:            body.Entidad,
		UrlKofax:           body.UrlKofax,
		Priority:           body.Priority,
		EntityBranchID:     body.EntityBranchID,
		EntityName:         body.EntityName,
		DepartmentID:       body.DepartmentID,
		CityID:             body.CityID,
		TownID:             body.TownID,
		OfficialDependency: body.OfficialDependency,
		Subject:            body.Subject,
		AsuntoRadicado:     body.AsuntoRadicado,
		CorreoEntidad:      body.CorreoEntidad,
		NumeroRadicado:     body.NumeroRadicado,
		ReasonCorrection:   body.ReasonCorrection,
		ResponseDate:       body.ResponseDate,
		CorreoRemitente:    body.CorreoRemitente,
		AsuntoRespuesta:    body.AsuntoRespuesta,
		ResponseReviewBy:   body.ResponseReviewBy,
	})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEntityLetterNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": "oficio no encontrado"})
		case errors.Is(err, service.ErrEntityLetterInvalidState):
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusOK, letter)
}
