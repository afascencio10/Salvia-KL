package salvia_facades

import (
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	"bitsflow/salvia/service"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// DownloadContactsReport maneja la descarga del Excel consolidado de reportes
// (victim_contact) por rango de fechas de creación. Reutiliza DownloadReportRequest.
func DownloadContactsReport(c *gin.Context, svc service.ReportService) {
	session := sessions.Default(c)
	userData := session.Get("userData")
	if userData == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error_code": "AUTH_TOKEN_INVALID", "message": "Sesión no válida o expirada"})
		return
	}
	sessionID := userData.(string)
	s, err := utils.GetCommonSession(sessionID)
	if err != nil || s == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error_code": "AUTH_TOKEN_INVALID", "message": "Sesión no válida o expirada"})
		return
	}

	// 1. Autorización: verificar permiso report_contacts_consolidated (ad, sv)
	if !utils.CheckPermission(salvia_config.PermissionsByRole, "report_contacts_consolidated", s.CurrentRole, c) {
		c.JSON(http.StatusForbidden, gin.H{"error_code": "AUTH_INSUFFICIENT_ROLE", "message": "El rol no tiene permisos para descargar el reporte"})
		return
	}

	// 2. Deserializar request
	var req DownloadReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "VALIDATION_FAILED", "message": "Fechas requeridas o formato incorrecto"})
		return
	}

	// 3. Validar fechas YYYY-MM-DD en America/Bogota
	loc, err := time.LoadLocation("America/Bogota")
	if err != nil {
		loc = time.Local
	}

	start, err := time.ParseInLocation("2006-01-02", req.StartDate, loc)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "VALIDATION_FAILED", "message": "Fecha inicial inválida (debe ser YYYY-MM-DD)"})
		return
	}

	end, err := time.ParseInLocation("2006-01-02", req.EndDate, loc)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "VALIDATION_FAILED", "message": "Fecha final inválida (debe ser YYYY-MM-DD)"})
		return
	}

	if start.After(end) {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "VALIDATION_FAILED", "message": "La fecha inicial no puede ser posterior a la fecha final"})
		return
	}

	// Rango <= 366 días (máximo 1 año)
	diffDays := int(end.Sub(start).Hours() / 24)
	if diffDays > 366 {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "VALIDATION_FAILED", "message": "El rango seleccionado no puede superar los 366 días (1 año)"})
		return
	}

	// Abarcar los días completos del rango
	startRange := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, loc)
	endRange := time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 999999999, loc)

	// 4. Invocar servicio de generación
	file, err := svc.GenerateConsolidatedContactsReport(c.Request.Context(), service.FollowUpsReportRange{
		StartDate: startRange,
		EndDate:   endRange,
	})
	if err != nil {
		if err.Error() == "no_contacts_found" {
			c.JSON(http.StatusBadRequest, gin.H{"error_code": "VALIDATION_FAILED", "message": "No hay reportes en el rango seleccionado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "SYSTEM_INTERNAL_ERROR", "message": "Error interno al generar el reporte Excel"})
		return
	}
	defer file.Close()

	// 5. Headers HTTP para descarga del binario
	filename := fmt.Sprintf("reporte-contactos_%s_%s.xlsx", req.StartDate, req.EndDate)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")

	if _, err := file.WriteTo(c.Writer); err != nil {
		c.Error(err)
	}
}
