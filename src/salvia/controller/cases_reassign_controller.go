// Package controller — cases_reassign_controller.go
// Endpoint JSON para reasignación masiva de casos (M-05).
package controller

import (
	"bitsflow/common/utils"
	"bitsflow/salvia/service"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// CasesReassignController maneja POST /api/v1/cases/reasignar-bulk.
type CasesReassignController struct {
	svc service.CasesReassignService
}

func NewCasesReassignController(svc service.CasesReassignService) *CasesReassignController {
	return &CasesReassignController{svc: svc}
}

func (c *CasesReassignController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/cases/reasignar-bulk", c.ReassignBulk)
}

type reassignBulkRequest struct {
	CaseICodes []string `json:"case_icodes" binding:"required"`
	AgentICode string   `json:"agent_icode" binding:"required"`
}

// ReassignBulk reasigna varios casos a un agente (M-05).
func (c *CasesReassignController) ReassignBulk(ctx *gin.Context) {
	session := sessions.Default(ctx)
	sessionID, ok := session.Get("userData").(string)
	if !ok || sessionID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "sesión inválida"})
		return
	}

	sess, err := utils.GetCommonSession(sessionID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "sesión expirada"})
		return
	}

	var body reassignBulkRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	supervisorName := strings.TrimSpace(sess.Names + " " + sess.LastNames)

	result, err := c.svc.ReassignBulk(ctx.Request.Context(), service.CasesReassignInput{
		CaseICodes:      body.CaseICodes,
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
