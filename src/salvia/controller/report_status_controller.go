// Package controller — report_status_controller.go
// Consulta pública del estado de un reporte por su radicado (T4/REP-02).
// PÚBLICO: sin sesión y sin captcha; la protección contra consultas masivas es
// el límite de fallos por IP del servicio. La URI /api/ pasa AuthMiddleware sin
// sesión, como las demás rutas de la API que validan por su cuenta.
package controller

import (
	"bitsflow/salvia/service"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	reportStatusMsgNotFound   = "No encontramos un reporte con ese radicado. Verifique que esté bien escrito."
	reportStatusMsgBadRequest = "Solicitud inválida."
	reportStatusMsgInternal   = "No fue posible consultar en este momento. Intente de nuevo."

	// El cuerpo legítimo es {"radicado":"XXXX-XXXX"}: 1 KiB sobra.
	reportStatusMaxBodyBytes = 1024
)

// ReportStatusController maneja POST /api/v1/public/consulta-reporte.
type ReportStatusController struct {
	svc service.ReportStatusService
}

func NewReportStatusController(svc service.ReportStatusService) *ReportStatusController {
	return &ReportStatusController{svc: svc}
}

// RegisterRoutes registra la consulta pública.
//
//	POST /api/v1/public/consulta-reporte   { "radicado": "K7QM-4R2X" }
//	  200 { radicado, estado, estadoTexto, descripcion, fechaRadicacion }
//	  404 { mensaje }                       formato inválido o radicado inexistente (mismo cuerpo)
//	  429 { mensaje, reintentarEnSegundos } 5 fallos de la IP en 15 min
func (c *ReportStatusController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/public/consulta-reporte", c.Lookup)
}

type reportStatusRequest struct {
	Radicado string `json:"radicado"`
}

// Lookup (CR-E03) consulta el estado público del reporte.
func (c *ReportStatusController) Lookup(ctx *gin.Context) {
	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.Header("Pragma", "no-cache")
	ctx.Header("Expires", "0")

	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, reportStatusMaxBodyBytes)
	var req reportStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"mensaje": reportStatusMsgBadRequest})
		return
	}

	res, err := c.svc.Lookup(ctx.Request.Context(), ctx.ClientIP(), req.Radicado)

	var limited *service.ReportStatusRateLimitedError
	switch {
	case err == nil:
		ctx.JSON(http.StatusOK, res)
	case errors.As(err, &limited):
		seconds := int(math.Ceil(limited.RetryAfter.Seconds()))
		if seconds < 1 {
			seconds = 1
		}
		minutes := (seconds + 59) / 60
		unit := "minutos"
		if minutes == 1 {
			unit = "minuto"
		}
		ctx.Header("Retry-After", strconv.Itoa(seconds))
		ctx.JSON(http.StatusTooManyRequests, gin.H{
			"mensaje":              fmt.Sprintf("Demasiados intentos. Intente de nuevo en %d %s.", minutes, unit),
			"reintentarEnSegundos": seconds,
		})
	case errors.Is(err, service.ErrReportStatusNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{"mensaje": reportStatusMsgNotFound})
	default:
		log.Printf("[consulta-reporte] error consultando el estado: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"mensaje": reportStatusMsgInternal})
	}
}
