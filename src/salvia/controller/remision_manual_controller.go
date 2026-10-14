// Package controller — remision_manual_controller.go
// Endpoint para el modal "Remisión a equipos" del Equipo Técnico en Detalle
// del Caso → Derivaciones. Replica las preguntas de orden 5 a 13 de la
// sección Seguimiento de Caso (form de seguimiento) fuera de dinamic-form.
package controller

import (
	"bitsflow/salvia/service"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RemisionManualController struct {
	svc service.FormService
}

func NewRemisionManualController(svc service.FormService) *RemisionManualController {
	return &RemisionManualController{svc: svc}
}

// RegisterRoutes registra el endpoint en /api/v1.
//
//	POST /api/v1/casos/:id/remision-equipos
func (c *RemisionManualController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/casos/:id/remision-equipos", c.SaveRemisionManual)
}

// SaveRemisionManual crea derivaciones a equipos SALVIA (y medidas de
// emergencia) directamente para el caso, sin pasar por un seguimiento.
// Accesible para el rol 'te' (Equipo Técnico) y 'sv' (Supervisor).
//
//	POST /api/v1/casos/:id/remision-equipos
//	Body: { "answers": { "<questionId>": "valor1,valor2", ... } }
func (c *RemisionManualController) SaveRemisionManual(ctx *gin.Context) {
	s := requireSessionRole(ctx, map[string]bool{"te": true, "sv": true})
	if s == nil {
		return
	}

	caseID := ctx.Param("id")
	if caseID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id de caso requerido"})
		return
	}

	var body struct {
		Answers map[string]string `json:"answers" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "campo 'answers' requerido"})
		return
	}

	count, err := c.svc.ProcessManualRemision(ctx.Request.Context(), caseID, s.UserICode, body.Answers)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrRemisionManualInvalid):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "selecciona al menos un equipo"})
		case errors.Is(err, service.ErrCaseNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": "caso no encontrado"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		}
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"remisionesCreadas": count})
}
