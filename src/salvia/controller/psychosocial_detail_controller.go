// Package controller — psychosocial_detail_controller.go
// Endpoint para la pantalla de detalle de remisión psicosocial.
package controller

import (
	"bitsflow/salvia/service"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PsychosocialDetailController struct {
	svc service.PsychosocialDetailService
}

func NewPsychosocialDetailController(db *gorm.DB) *PsychosocialDetailController {
	return &PsychosocialDetailController{
		svc: service.NewPsychosocialDetailService(db),
	}
}

func (c *PsychosocialDetailController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/psychosocial-support/:id/detail", c.GetDetail)
	rg.POST("/psychosocial-support/:id/contacts", c.CreateContact)
	rg.PUT("/psychosocial-support/contacts/:contactId/reschedule", c.RescheduleContact)
	rg.PUT("/psychosocial-support/contacts/:contactId/cancel", c.CancelContact)
	rg.PUT("/psychosocial-support/:id/schedule-preference", c.UpdateSchedulePreference)
	rg.GET("/psychosocial-support/:id/load", c.LoadSession)
}

// GetDetail retorna el detalle completo de una remisión psicosocial.
// GET /api/v1/psychosocial-support/:id/detail
func (c *PsychosocialDetailController) GetDetail(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id requerido"})
		return
	}

	detail, err := c.svc.GetDetail(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "remisión no encontrada"})
		return
	}

	ctx.JSON(http.StatusOK, detail)
}

// CreateContact registra un nuevo contacto/sesión para la remisión.
// POST /api/v1/psychosocial-support/:id/contacts
func (c *PsychosocialDetailController) CreateContact(ctx *gin.Context) {
	psicosocialID := ctx.Param("id")
	if psicosocialID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id requerido"})
		return
	}

	var body struct {
		ContactDate   string `json:"contactDate" binding:"required"`
		ContactTime   string `json:"contactTime" binding:"required"`
		Type          string `json:"type" binding:"required"`
		Summary       string `json:"summary"`
		SessionType   string `json:"sessionType"`
		ScheduledDate string `json:"scheduledDate"`
		ScheduledTime string `json:"scheduledTime"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contact, err := c.svc.CreateContact(ctx.Request.Context(), psicosocialID, body.Type, body.ContactDate, body.ContactTime, body.Summary, body.SessionType, body.ScheduledDate, body.ScheduledTime)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, contact)
}

// RescheduleContact actualiza la fecha/hora de una sesión agendada sin duplicar.
// PUT /api/v1/psychosocial-support/contacts/:contactId/reschedule
func (c *PsychosocialDetailController) RescheduleContact(ctx *gin.Context) {
	contactID := ctx.Param("contactId")
	var body struct {
		ScheduledDate string `json:"scheduledDate" binding:"required"`
		ScheduledTime string `json:"scheduledTime" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.svc.RescheduleContact(ctx.Request.Context(), contactID, body.ScheduledDate, body.ScheduledTime); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"ok": true})
}

// CancelContact marca una sesión agendada como cancelada sin eliminarla.
// PUT /api/v1/psychosocial-support/contacts/:contactId/cancel
func (c *PsychosocialDetailController) CancelContact(ctx *gin.Context) {
	contactID := ctx.Param("contactId")

	if err := c.svc.CancelContact(ctx.Request.Context(), contactID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"ok": true})
}

// LoadSession carga la pantalla "Registrar Sesión Psicosocial" — evento E-01:
// selecciona el formulario según el estado del proceso, resuelve el team_contact
// y form_submission activos, y retorna la info de la víctima.
// GET /api/v1/psychosocial-support/:id/load?agent_id=...
func (c *PsychosocialDetailController) LoadSession(ctx *gin.Context) {
	id := ctx.Param("id")
	agentID := ctx.Query("agent_id")
	if id == "" || agentID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id y agent_id son requeridos"})
		return
	}

	result, err := c.svc.LoadSession(ctx.Request.Context(), id, agentID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrPsicosocialSessionNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": "remisión psicosocial no encontrada"})
		case errors.Is(err, service.ErrPsicosocialSessionNotAssigned):
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Hubo un error al obtener la información del servidor, por favor verifique su conexión a internet y vuelva a intentarlo"})
		}
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// UpdateSchedulePreference actualiza la preferencia de horario de la paciente.
// PUT /api/v1/psychosocial-support/:id/schedule-preference
func (c *PsychosocialDetailController) UpdateSchedulePreference(ctx *gin.Context) {
	id := ctx.Param("id")
	var body struct {
		Preference string `json:"preference" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.svc.UpdateSchedulePreference(ctx.Request.Context(), id, body.Preference); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"ok": true})
}
