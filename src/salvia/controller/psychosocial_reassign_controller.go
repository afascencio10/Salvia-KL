// Package controller — psychosocial_reassign_controller.go
// Endpoints JSON para reasignar-remisiones-modal (RRM-03, RRM-05).
package controller

import (
	"bitsflow/common/utils"
	"bitsflow/salvia/service"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// PsychosocialReassignController maneja catálogos y reasignación bulk de remisiones psicosociales.
type PsychosocialReassignController struct {
	svc service.PsychosocialReassignService
}

func NewPsychosocialReassignController(svc service.PsychosocialReassignService) *PsychosocialReassignController {
	return &PsychosocialReassignController{svc: svc}
}

func (c *PsychosocialReassignController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/psychosocial-support/profesionales-reasignacion", c.ProfesionalesReasignacion)
	rg.GET("/duplas/reasignacion", c.DuplasReasignacion)
	rg.POST("/psychosocial-support/reasignar-bulk", c.ReassignBulk)
}

func (c *PsychosocialReassignController) ProfesionalesReasignacion(ctx *gin.Context) {
	result, err := c.svc.ListProfessionalsForReassign(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al cargar profesionales"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *PsychosocialReassignController) DuplasReasignacion(ctx *gin.Context) {
	result, err := c.svc.ListDuplasForReassign(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al cargar duplas"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

type psychosocialReassignBulkRequest struct {
	RemisionIDs    []string `json:"remision_ids" binding:"required"`
	AssignMode     string   `json:"assign_mode" binding:"required"`
	ProfessionalID *string  `json:"professional_id"`
	DuplaID        *string  `json:"dupla_id"`
}

func (c *PsychosocialReassignController) ReassignBulk(ctx *gin.Context) {
	session := sessions.Default(ctx)
	sessionID, ok := session.Get("userData").(string)
	if !ok || sessionID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "sesión inválida"})
		return
	}
	if _, err := utils.GetCommonSession(sessionID); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "sesión expirada"})
		return
	}

	var body psychosocialReassignBulkRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var professionalID, duplaID string
	if body.ProfessionalID != nil {
		professionalID = strings.TrimSpace(*body.ProfessionalID)
	}
	if body.DuplaID != nil {
		duplaID = strings.TrimSpace(*body.DuplaID)
	}

	result, err := c.svc.ReassignBulk(ctx.Request.Context(), service.PsychosocialReassignBulkInput{
		RemisionIDs:    body.RemisionIDs,
		AssignMode:     body.AssignMode,
		ProfessionalID: professionalID,
		DuplaID:        duplaID,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}
