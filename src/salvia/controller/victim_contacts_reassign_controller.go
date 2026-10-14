// Package controller — victim_contacts_reassign_controller.go
// Endpoint JSON para reasignación masiva de reportes (Listado Reportes E-06).
// Solo sv — FigJam: "Reasignar un reporte a un Agente integral. Solo el
// supervisor puede hacer esto."
package controller

import (
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	"bitsflow/salvia/service"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type VictimContactsReassignController struct {
	svc service.VictimContactsReassignService
}

func NewVictimContactsReassignController(svc service.VictimContactsReassignService) *VictimContactsReassignController {
	return &VictimContactsReassignController{svc: svc}
}

func (c *VictimContactsReassignController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/victim-contacts/reasignar-bulk", c.ReassignBulk)
}

type victimContactsReassignBulkRequest struct {
	ReportICodes []string `json:"report_icodes" binding:"required"`
	AgentICode   string   `json:"agent_icode" binding:"required"`
}

// ReassignBulk reasigna varios reportes a un agente integral. Solo sv
// (permiso reassign_victim_contacts, Menu.go).
func (c *VictimContactsReassignController) ReassignBulk(ctx *gin.Context) {
	sessionID, ok := sessions.Default(ctx).Get("userData").(string)
	if !ok || sessionID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "sesión inválida"})
		return
	}
	sess, err := utils.GetCommonSession(sessionID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "sesión expirada"})
		return
	}
	if !utils.CheckPermission(salvia_config.PermissionsByRole, "reassign_victim_contacts", sess.CurrentRole, ctx) {
		return // CheckPermission ya escribió el 401
	}

	var body victimContactsReassignBulkRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	supervisorName := strings.TrimSpace(sess.Names + " " + sess.LastNames)

	result, err := c.svc.ReassignBulk(ctx.Request.Context(), service.VictimContactsReassignInput{
		ReportICodes:    body.ReportICodes,
		AgentICode:      body.AgentICode,
		SupervisorICode: sess.UserICode,
		SupervisorName:  supervisorName,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}
