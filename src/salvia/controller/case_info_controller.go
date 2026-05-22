// Package controller — case_info_controller.go
// Endpoint para el componente <case-info>: información completa del caso.
package controller

import (
	"bitsflow/salvia/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CaseInfoController struct {
	svc service.CaseInfoService
}

func NewCaseInfoController(svc service.CaseInfoService) *CaseInfoController {
	return &CaseInfoController{svc: svc}
}

func (c *CaseInfoController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/casos/:id/info-completa", c.GetFullInfo)
}

// GetFullInfo retorna toda la información del caso combinando victim_case + form1 + form2.
// GET /api/v1/casos/:id/info-completa
func (c *CaseInfoController) GetFullInfo(ctx *gin.Context) {
	caseICode := ctx.Param("id")
	if caseICode == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id requerido"})
		return
	}

	info, err := c.svc.GetFullCaseInfo(ctx.Request.Context(), caseICode)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, info)
}
