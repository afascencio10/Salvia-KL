package controller

import (
	"bitsflow/common/utils"
	"bitsflow/salvia/service"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// AssignCaseController maneja POST /api/v1/cases/assign.
type AssignCaseController struct {
	svc service.AssignCaseService
}

func NewAssignCaseController(svc service.AssignCaseService) *AssignCaseController {
	return &AssignCaseController{svc: svc}
}

func (c *AssignCaseController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/cases/assign", c.Assign)
}

func (c *AssignCaseController) Assign(ctx *gin.Context) {
	var body service.AssignCaseInput
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if body.DocNumber == "" || body.AgentLogin == "" || body.Team == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "docNumber, agentLogin y team son requeridos"})
		return
	}

	// Este endpoint no exige sesión, pero si llega una (uso desde el front-end
	// autenticado) se aprovecha para saber quién hizo la reasignación en el log.
	if sessionID, ok := sessions.Default(ctx).Get("userData").(string); ok && sessionID != "" {
		if sess, err := utils.GetCommonSession(sessionID); err == nil {
			body.ActorICode = sess.UserICode
			body.ActorName = strings.TrimSpace(sess.Names + " " + sess.LastNames)
		}
	}

	result, err := c.svc.AssignCase(ctx.Request.Context(), body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}
