package controller

import (
	salvia_facades "bitsflow/salvia/facades"
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
	rg.POST("/reportes/seguimientos-consolidado", c.DownloadConsolidatedReport)
	rg.POST("/reportes/contactos-consolidado", c.DownloadContactsReport)
}

func (c *ReportController) DownloadConsolidatedReport(ctx *gin.Context) {
	salvia_facades.DownloadConsolidatedReport(ctx, c.svc)
}

func (c *ReportController) DownloadContactsReport(ctx *gin.Context) {
	salvia_facades.DownloadContactsReport(ctx, c.svc)
}
