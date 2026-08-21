// Package controller — report_async_controller.go
// Generación asíncrona de reportes para evitar timeout del gateway.
// Endpoints:
//   POST /api/v1/reportes/iniciar          → inicia generación en background
//   GET  /api/v1/reportes/estado/:jobId    → consulta estado del job
//   GET  /api/v1/reportes/descargar/:jobId → descarga el archivo generado
package controller

import (
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	"bitsflow/salvia/service"
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

// ReportJob representa un trabajo de generación de reporte en background.
type ReportJob struct {
	ID        string
	Status    string // "procesando" | "listo" | "error"
	Progress  string
	File      *excelize.File
	Error     string
	CreatedAt time.Time
	UserID    string
}

// ReportAsyncController maneja la generación asíncrona de reportes.
type ReportAsyncController struct {
	svc  service.ReportService
	jobs map[string]*ReportJob
	mu   sync.RWMutex
}

func NewReportAsyncController(svc service.ReportService) *ReportAsyncController {
	ctrl := &ReportAsyncController{
		svc:  svc,
		jobs: make(map[string]*ReportJob),
	}
	// Limpiar jobs expirados cada 10 minutos
	go ctrl.cleanupLoop()
	return ctrl
}

func (c *ReportAsyncController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/reportes/iniciar", c.IniciarReporte)
	rg.GET("/reportes/estado/:jobId", c.EstadoReporte)
	rg.GET("/reportes/descargar/:jobId", c.DescargarReporte)
}

// IniciarReporte inicia la generación del reporte en background.
// POST /api/v1/reportes/iniciar { start_date, end_date }
func (c *ReportAsyncController) IniciarReporte(ctx *gin.Context) {
	// Auth
	session := sessions.Default(ctx)
	userData := session.Get("userData")
	if userData == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no autenticado"})
		return
	}
	sessionID := userData.(string)
	s, err := utils.GetCommonSession(sessionID)
	if err != nil || s == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "sesión inválida"})
		return
	}
	if !utils.CheckPermission(salvia_config.PermissionsByRole, "report_followups_consolidated", s.CurrentRole, ctx) {
		return
	}

	// Parse body
	var req struct {
		StartDate string `json:"start_date" binding:"required"`
		EndDate   string `json:"end_date" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "fechas requeridas"})
		return
	}

	// Validar fechas
	loc, _ := time.LoadLocation("America/Bogota")
	start, err := time.ParseInLocation("2006-01-02", req.StartDate, loc)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "fecha inicial inválida"})
		return
	}
	end, err := time.ParseInLocation("2006-01-02", req.EndDate, loc)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "fecha final inválida"})
		return
	}
	if start.After(end) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "fecha inicial posterior a fecha final"})
		return
	}
	if int(end.Sub(start).Hours()/24) > 366 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "rango máximo 366 días"})
		return
	}

	// Crear job
	jobID := fmt.Sprintf("%s_%d", s.UserICode, time.Now().UnixNano())
	job := &ReportJob{
		ID:        jobID,
		Status:    "procesando",
		Progress:  "Iniciando generación...",
		CreatedAt: time.Now(),
		UserID:    s.UserICode,
	}

	c.mu.Lock()
	c.jobs[jobID] = job
	c.mu.Unlock()

	// Lanzar generación en background
	startRange := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, loc)
	endRange := time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 999999999, loc)

	go func() {
		job.Progress = "Consultando casos..."
		file, err := c.svc.GenerateConsolidatedFollowUpsReport(context.Background(), service.FollowUpsReportRange{
			StartDate: startRange,
			EndDate:   endRange,
		})
		if err != nil {
			job.Status = "error"
			job.Error = err.Error()
			if err.Error() == "no_cases_found" {
				job.Error = "No se encontraron casos en el rango seleccionado"
			}
			log.Printf("[ReportAsync] Job %s error: %v", jobID, err)
			return
		}
		job.File = file
		job.Status = "listo"
		job.Progress = "Reporte generado exitosamente"
		log.Printf("[ReportAsync] Job %s completado", jobID)
	}()

	ctx.JSON(http.StatusOK, gin.H{"jobId": jobID, "status": "procesando"})
}

// EstadoReporte consulta el estado de un job.
// GET /api/v1/reportes/estado/:jobId
func (c *ReportAsyncController) EstadoReporte(ctx *gin.Context) {
	jobID := ctx.Param("jobId")

	c.mu.RLock()
	job, exists := c.jobs[jobID]
	c.mu.RUnlock()

	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "reporte no encontrado", "status": "not_found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"jobId":    job.ID,
		"status":   job.Status,
		"progress": job.Progress,
		"error":    job.Error,
	})
}

// DescargarReporte descarga el archivo Excel generado.
// GET /api/v1/reportes/descargar/:jobId
func (c *ReportAsyncController) DescargarReporte(ctx *gin.Context) {
	jobID := ctx.Param("jobId")

	c.mu.RLock()
	job, exists := c.jobs[jobID]
	c.mu.RUnlock()

	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "reporte no encontrado"})
		return
	}
	if job.Status != "listo" || job.File == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "reporte no está listo aún"})
		return
	}

	filename := fmt.Sprintf("reporte-casos_%s.xlsx", time.Now().Format("2006-01-02"))
	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")

	if _, err := job.File.WriteTo(ctx.Writer); err != nil {
		log.Printf("[ReportAsync] Error escribiendo archivo: %v", err)
	}

	// Limpiar job después de descargar
	c.mu.Lock()
	delete(c.jobs, jobID)
	c.mu.Unlock()
	job.File.Close()
}

// cleanupLoop elimina jobs expirados (>30 min) cada 10 minutos.
func (c *ReportAsyncController) cleanupLoop() {
	ticker := time.NewTicker(10 * time.Minute)
	for range ticker.C {
		c.mu.Lock()
		for id, job := range c.jobs {
			if time.Since(job.CreatedAt) > 30*time.Minute {
				if job.File != nil {
					job.File.Close()
				}
				delete(c.jobs, id)
				log.Printf("[ReportAsync] Job %s expirado y limpiado", id)
			}
		}
		c.mu.Unlock()
	}
}
