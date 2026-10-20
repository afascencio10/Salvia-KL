// Package controller — victim_contacts_management_controller.go
// Endpoints JSON para gestion-reporte-modal: detalle (M-01), intento de
// contacto (M-05/M-06), cierre (M-02/M-07/M-08/M-09). Ver
// DocsMD/Componentes/Gestion Reporte Modal/gestion-reporte-modal-events.md.
//
// M-03 ("Reporte de prueba") y M-10 ("Apertura de caso") NO se implementan
// en esta fase — ver notas en gestion-reporte-modal-events.md.
package controller

import (
	"bitsflow/common/utils"
	"bitsflow/internal/repository"
	salvia_config "bitsflow/salvia/config"
	"bitsflow/salvia/service"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type VictimContactsManagementController struct {
	svc service.VictimContactsManagementService
}

func NewVictimContactsManagementController(svc service.VictimContactsManagementService) *VictimContactsManagementController {
	return &VictimContactsManagementController{svc: svc}
}

// RegisterRoutes registra las rutas de gestión de un reporte individual.
//
//	GET  /api/v1/victim-contacts/:icode
//	POST /api/v1/victim-contacts/:icode/intento
//	POST /api/v1/victim-contacts/:icode/cerrar
func (c *VictimContactsManagementController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/victim-contacts/:icode", c.GetDetail)
	rg.POST("/victim-contacts/:icode/intento", c.RegisterAttempt)
	rg.POST("/victim-contacts/:icode/cerrar", c.Close)
}

// requireVictimContactsAccess valida sesión + permiso manage_victim_contacts
// (op y sv, GAP-2 confirmado 2026-10-11 — ambos pueden gestionar cualquier
// reporte visible para ellos).
func requireVictimContactsAccess(ctx *gin.Context) *utils.CommonSession {
	sessionID, ok := sessions.Default(ctx).Get("userData").(string)
	if !ok || sessionID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no autenticado"})
		return nil
	}
	s, err := utils.GetCommonSession(sessionID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "sesión expirada"})
		return nil
	}
	if !utils.CheckPermission(salvia_config.PermissionsByRole, "manage_victim_contacts", s.CurrentRole, ctx) {
		return nil // CheckPermission ya escribió el 401
	}
	return s
}

// GetDetail responde con el detalle completo del reporte (M-01).
func (c *VictimContactsManagementController) GetDetail(ctx *gin.Context) {
	s := requireVictimContactsAccess(ctx)
	if s == nil {
		return
	}
	icode := ctx.Param("icode")
	canManage := s.CurrentRole == "op" || s.CurrentRole == "sv"

	detail, err := c.svc.GetDetail(ctx.Request.Context(), icode, canManage)
	if err != nil {
		if errors.Is(err, repository.ErrVictimContactNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "reporte no encontrado"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al cargar el reporte"})
		return
	}
	ctx.JSON(http.StatusOK, detail)
}

// RegisterAttempt registra un intento de contacto (M-05 fallido / M-06 exitoso).
func (c *VictimContactsManagementController) RegisterAttempt(ctx *gin.Context) {
	s := requireVictimContactsAccess(ctx)
	if s == nil {
		return
	}
	icode := ctx.Param("icode")

	var body struct {
		Result string `json:"result" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	intentosHoy, diasCompletos, err := c.svc.RegisterAttempt(ctx.Request.Context(), icode, body.Result, s.UserICode)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrVictimContactNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": "reporte no encontrado"})
		case errors.Is(err, repository.ErrVictimContactAlreadyClosed):
			ctx.JSON(http.StatusConflict, gin.H{"error": "el reporte ya está cerrado"})
		case errors.Is(err, service.ErrInvalidAttemptResult):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al registrar el intento"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"ok":            true,
		"intentosHoy":   intentosHoy,
		"diasCompletos": diasCompletos,
		"quotaReached":  intentosHoy >= 3,
		// Para que el log del modal muestre quién registró el intento sin recargar el detalle.
		"agentName": strings.TrimSpace(s.Names + " " + s.LastNames),
	})
}

// Close cierra el reporte (M-02 "No información mínima", M-07 "Cerrar por no
// contacto", M-08/M-09 post-contacto). M-03 ("Reporte de prueba") y M-10
// ("Apertura de caso") no se implementan en esta fase.
func (c *VictimContactsManagementController) Close(ctx *gin.Context) {
	s := requireVictimContactsAccess(ctx)
	if s == nil {
		return
	}
	icode := ctx.Param("icode")

	var body struct {
		CloseReason      string `json:"close_reason" binding:"required"`
		CloseObservation string `json:"close_observation"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := c.svc.Close(ctx.Request.Context(), icode, body.CloseReason, body.CloseObservation)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrVictimContactNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": "reporte no encontrado"})
		case errors.Is(err, repository.ErrVictimContactAlreadyClosed):
			ctx.JSON(http.StatusConflict, gin.H{"error": "el reporte ya está cerrado"})
		case errors.Is(err, repository.ErrSchemeNotComplete):
			ctx.JSON(http.StatusConflict, gin.H{"error": "el esquema 3x3 no está completo todavía"})
		case errors.Is(err, repository.ErrCloseReasonNotFound):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "motivo de cierre no encontrado"})
		case errors.Is(err, service.ErrInvalidCloseReason), errors.Is(err, service.ErrCloseReasonDeferred):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrObservationRequired), errors.Is(err, service.ErrObservationTooLong):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al cerrar el reporte"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"ok": true, "closeReason": body.CloseReason, "category": category})
}
