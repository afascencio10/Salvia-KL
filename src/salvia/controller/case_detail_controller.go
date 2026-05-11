// Package controller — case_detail_controller.go
// Controlador HTTP para la pantalla de detalle de caso (rol sv).
//
// ARQUITECTURA:
//
//	main.go → RegisterRoutes → CaseDetailGET (fachada HTML)
//	main.go → GET /api/v1/casos/:id/detalle → GetByID (API JSON)
//
// La ruta HTML vive en MainRouter.go (sistema de fachadas existente).
// La ruta API JSON se registra aquí con RegisterRoutes.
package controller

import (
	"bitsflow/common/db"
	"bitsflow/common/utils"
	"bitsflow/salvia/service"
	security_ctrl "bitsflow/security/controllers"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CaseDetailController struct {
	svc service.CaseDetailService
}

func NewCaseDetailController(svc service.CaseDetailService) *CaseDetailController {
	return &CaseDetailController{svc: svc}
}

// RegisterRoutes registra el endpoint JSON en /api/v1.
//
//	GET /api/v1/casos/:id/detalle
//	GET /api/v1/operadores
func (c *CaseDetailController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/casos/:id/detalle", c.GetByID)
	rg.GET("/operadores", c.GetOperadores)
	rg.GET("/agentes-ro", c.GetAgentesRO)
	rg.POST("/casos/:id/seguimiento", c.CrearSeguimiento)
	rg.POST("/casos/:id/timeline", c.AddTimelineEvent)
	rg.PUT("/seguimiento/:segId/reasignar", c.ReassignFollowUp)
}

// GetByID responde con el detalle completo del caso en JSON.
// Útil para consumir desde el frontend sin recargar la página.
func (c *CaseDetailController) GetByID(ctx *gin.Context) {
	icode := ctx.Param("id")
	detail, err := c.svc.GetDetail(ctx.Request.Context(), icode)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidICode):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "icode inválido"})
		case errors.Is(err, service.ErrCaseNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": "caso no encontrado"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"})
		}
		return
	}
	ctx.JSON(http.StatusOK, detail)
}

// GetOperadores devuelve la lista de operadores (rol "op") para el dropdown de asignación.
func (c *CaseDetailController) GetOperadores(ctx *gin.Context) {
	dbCfg := utils.LoadDBCLientConfig()
	dbSrv := db.DBServerConfig{PoolSize: 80}
	code, res := security_ctrl.GetGeneralUsersByRole("op", &db.ConnData{}, dbCfg, dbSrv)
	ctx.Data(code, "application/json", []byte(res))
}

// GetAgentesRO devuelve la lista de usuarios con rol "ro" para asignar seguimientos.
func (c *CaseDetailController) GetAgentesRO(ctx *gin.Context) {
	//No crear conexiones a la DB en controllers, usar repositorios
	dbCfg := utils.LoadDBCLientConfig()
	dbSrv := db.DBServerConfig{PoolSize: 80}
	code, res := security_ctrl.GetGeneralUsersByRole("ro", &db.ConnData{}, dbCfg, dbSrv)
	ctx.Data(code, "application/json", []byte(res))
}

// CrearSeguimiento crea un registro básico en follow_up_v2.
// POST /api/v1/casos/:id/seguimiento
func (c *CaseDetailController) CrearSeguimiento(ctx *gin.Context) {
	caseICode := ctx.Param("id")

	var body struct {
		AgentID       string `json:"agent_id"       binding:"required"`
		ScheduledDate string `json:"scheduled_date"  binding:"required"`
		Notas         string `json:"notas"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	followUp, err := c.svc.CreateFollowUp(ctx.Request.Context(), caseICode, body.AgentID, body.ScheduledDate, body.Notas)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, followUp)
}

// AddTimelineEvent registra un evento en el timeline del caso.
// POST /api/v1/casos/:id/timeline
func (c *CaseDetailController) AddTimelineEvent(ctx *gin.Context) {
	caseICode := ctx.Param("id")
	var body struct {
		EventType   string `json:"event_type"   binding:"required"`
		Description string `json:"description"  binding:"required"`
		ActorName   string `json:"actor_name"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.svc.AddTimelineEvent(ctx.Request.Context(), caseICode, body.EventType, body.Description, "", body.ActorName); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"ok": true})
}

// ReassignFollowUp actualiza el agent_id de un seguimiento.
// PUT /api/v1/seguimiento/:segId/reasignar
func (c *CaseDetailController) ReassignFollowUp(ctx *gin.Context) {
	segID := ctx.Param("segId")
	var body struct {
		AgentID string `json:"agent_id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.svc.ReassignFollowUp(ctx.Request.Context(), segID, body.AgentID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"ok": true})
}
