package controller

import (
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	"bitsflow/salvia/service"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// EmergencyMeasureReferralController expone la bandeja de Medidas de Emergencia.
type EmergencyMeasureReferralController struct {
	svc service.EmergencyMeasureReferralService
}

func NewEmergencyMeasureReferralController(svc service.EmergencyMeasureReferralService) *EmergencyMeasureReferralController {
	return &EmergencyMeasureReferralController{svc: svc}
}

func (c *EmergencyMeasureReferralController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/emergency-measure-referrals/agentes", c.Agents)
	rg.GET("/emergency-measure-referrals/profesionales-revision", c.Reviewers)
	rg.GET("/emergency-measure-referrals", c.List)
	rg.GET("/emergency-measure-referrals/:id", c.Detail)
	rg.POST("/emergency-measure-referrals/reasignar", c.Reassign)
	rg.POST("/emergency-measure-referrals/:id/revision", c.Review)
}

func (c *EmergencyMeasureReferralController) List(ctx *gin.Context) {
	s := referralSession(ctx)
	if s == nil {
		return
	}
	vista, _, ok := salvia_config.VistaRemisionesMedidas(s.CurrentRole)
	if !ok {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "este rol no tiene bandeja de medidas"})
		return
	}
	perm := "get_remisiones_medidas"
	if vista == salvia_config.VistaMedidasAgente {
		perm = "get_mis_remisiones_medidas"
	}
	if !utils.CheckPermission(salvia_config.PermissionsByRole, perm, s.CurrentRole, ctx) {
		return
	}
	queryAgent := strings.TrimSpace(ctx.Query("agent_id"))
	if vista == salvia_config.VistaMedidasAgente && queryAgent != "" && queryAgent != s.UserICode {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "el agente no coincide con la sesión"})
		return
	}
	status := ""
	if vista == salvia_config.VistaMedidasSupervisor || vista == salvia_config.VistaMedidasAgente {
		status = strings.TrimSpace(ctx.Query("status"))
	}
	agentFilter := ""
	if vista == salvia_config.VistaMedidasSupervisor {
		agentFilter = queryAgent
	}
	result, err := c.svc.List(ctx.Request.Context(), service.EmergencyReferralListInput{
		Vista:          vista,
		SessionAgentID: s.UserICode,
		Q:              ctx.Query("q"),
		Status:         status,
		AgentID:        agentFilter,
		Urgency:        strings.TrimSpace(ctx.Query("urgency")),
	})
	if err != nil {
		writeReferralErr(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *EmergencyMeasureReferralController) Detail(ctx *gin.Context) {
	s := referralSession(ctx)
	if s == nil {
		return
	}
	if !utils.CheckPermission(salvia_config.PermissionsByRole, "get_remision_medida_detalle", s.CurrentRole, ctx) {
		return
	}
	vista, _, ok := salvia_config.VistaRemisionesMedidas(s.CurrentRole)
	if !ok {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "este rol no tiene bandeja de medidas"})
		return
	}
	item, err := c.svc.Detail(ctx.Request.Context(), vista, s.UserICode, ctx.Param("id"))
	if err != nil {
		writeReferralErr(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, item)
}

func (c *EmergencyMeasureReferralController) Agents(ctx *gin.Context) {
	s := referralSession(ctx)
	if s == nil {
		return
	}
	if !utils.CheckPermission(salvia_config.PermissionsByRole, "get_remisiones_medidas", s.CurrentRole, ctx) {
		return
	}
	result, err := c.svc.Agents(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al cargar los agentes"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *EmergencyMeasureReferralController) Reviewers(ctx *gin.Context) {
	s := referralSession(ctx)
	if s == nil {
		return
	}
	if !utils.CheckPermission(salvia_config.PermissionsByRole, "review_remisiones_medidas", s.CurrentRole, ctx) {
		return
	}
	result, err := c.svc.Reviewers(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al cargar los profesionales"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *EmergencyMeasureReferralController) Reassign(ctx *gin.Context) {
	s := referralSession(ctx)
	if s == nil {
		return
	}
	if !utils.CheckPermission(salvia_config.PermissionsByRole, "reassign_remisiones_medidas", s.CurrentRole, ctx) {
		return
	}
	var body struct {
		IDs     []string `json:"ids"`
		AgentID string   `json:"agentId"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "body inválido"})
		return
	}
	result, err := c.svc.Reassign(ctx.Request.Context(), body.IDs, body.AgentID)
	if err != nil {
		writeReferralErr(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *EmergencyMeasureReferralController) Review(ctx *gin.Context) {
	s := referralSession(ctx)
	if s == nil {
		return
	}
	if !utils.CheckPermission(salvia_config.PermissionsByRole, "review_remisiones_medidas", s.CurrentRole, ctx) {
		return
	}
	var body struct {
		ProfessionalID string `json:"professionalId"`
		Concepto       string `json:"concepto"`
		Resultado      string `json:"resultado"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "body inválido"})
		return
	}
	item, err := c.svc.Review(ctx.Request.Context(), ctx.Param("id"), body.ProfessionalID, body.Concepto, body.Resultado)
	if err != nil {
		writeReferralErr(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, item)
}

func referralSession(ctx *gin.Context) *utils.CommonSession {
	sessionID, ok := sessions.Default(ctx).Get("userData").(string)
	if !ok || sessionID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no autenticado"})
		return nil
	}
	s, err := utils.GetCommonSession(sessionID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no autenticado"})
		return nil
	}
	return s
}

func writeReferralErr(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrReferralInvalidFilter), errors.Is(err, service.ErrReferralInvalidAgent), errors.Is(err, service.ErrReferralInvalidReview):
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, service.ErrReferralNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, service.ErrReferralConflict):
		ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al procesar la remisión"})
	}
}
