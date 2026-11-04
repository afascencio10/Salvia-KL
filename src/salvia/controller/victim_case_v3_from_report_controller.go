// Package controller — victim_case_v3_from_report_controller.go
// «Registro de caso»: crea un borrador de Registro V3 a partir de un reporte
// (WP-1 T7/REP-03). Corte vertical nuevo: no toca victim_case_v3_controller.go.
package controller

import (
	"bitsflow/internal/repository"
	"bitsflow/salvia/service"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReportConversionController struct {
	svc service.ReportConversionService
}

func NewReportConversionController(svc service.ReportConversionService) *ReportConversionController {
	return &ReportConversionController{svc: svc}
}

// RegisterRoutes registra la conversión.
//
//	POST /api/v1/victim-cases-v3/desde-reporte/:icode
//	  200 { caseIcode, redirect }
//	  409 { error }                  reporte cerrado, tipo no soportado o sin contacto exitoso
//	  409 { error, caseIcode }       el reporte ya tiene un borrador (para «Ir al borrador del caso»)
func (c *ReportConversionController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/victim-cases-v3/desde-reporte/:icode", c.Convert)
}

// Convert (M-10 «Cuando confirma «Registro de caso»») — op y sv (permisos
// manage_victim_contacts y manage_registro_v3).
func (c *ReportConversionController) Convert(ctx *gin.Context) {
	s := requireVictimContactsAccess(ctx)
	if s == nil {
		return
	}
	if requireRegistroV3Access(ctx) == nil {
		return
	}

	caseICode, err := c.svc.ConvertToDraft(ctx.Request.Context(), ctx.Param("icode"), s.UserICode)
	if err != nil {
		var hasDraft *repository.ErrReportHasDraft
		switch {
		case errors.Is(err, repository.ErrVictimContactNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": "reporte no encontrado"})
		case errors.As(err, &hasDraft):
			ctx.JSON(http.StatusConflict, gin.H{"error": hasDraft.Error(), "caseIcode": hasDraft.CaseICode})
		case errors.Is(err, repository.ErrVictimContactAlreadyClosed):
			ctx.JSON(http.StatusConflict, gin.H{"error": "el reporte ya está cerrado"})
		case errors.Is(err, repository.ErrReportHasLinkedCase),
			errors.Is(err, repository.ErrReportNotConvertibleType),
			errors.Is(err, repository.ErrReportNoContact):
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			log.Printf("[reportConversion] error convirtiendo el reporte %s: %v", ctx.Param("icode"), err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al crear el registro"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"caseIcode": caseICode,
		"redirect":  "/salvia/registro-v3?caso=" + caseICode,
	})
}
