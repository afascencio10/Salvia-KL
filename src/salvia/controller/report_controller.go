package controller

import (
	"bitsflow/salvia/service"

	"github.com/gin-gonic/gin"
)

type ReportController struct {
	svc service.ReportService
}

func NewReportController(svc service.ReportService) *ReportController {
	return &ReportController{svc: svc}
}

func (c *ReportController) RegisterRoutes(rg *gin.RouterGroup) {
	// TODO: registrar rutas de reportes
}
