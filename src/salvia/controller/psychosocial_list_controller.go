// Package controller — psychosocial_list_controller.go
// Endpoints JSON para remisiones-psicosocial-component (evento E-01).
package controller

import (
	"bitsflow/salvia/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// PsychosocialListController maneja listado, stats, duplas y equipos remitentes.
type PsychosocialListController struct {
	svc service.PsychosocialListService
}

func NewPsychosocialListController(svc service.PsychosocialListService) *PsychosocialListController {
	return &PsychosocialListController{svc: svc}
}

// RegisterRoutes registra rutas del componente remisiones psicosocial.
func (c *PsychosocialListController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/psychosocial-support/list", c.List)
	rg.GET("/psychosocial-support/stats", c.Stats)
	rg.GET("/psychosocial-support/equipos-remitentes", c.EquiposRemitentes)
	rg.GET("/duplas", c.Duplas)
}

func (c *PsychosocialListController) List(ctx *gin.Context) {
	input := c.inputFromQuery(ctx)
	result, err := c.svc.List(ctx.Request.Context(), input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al cargar las remisiones"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *PsychosocialListController) Stats(ctx *gin.Context) {
	input := c.inputFromQuery(ctx)
	result, err := c.svc.Stats(ctx.Request.Context(), input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al cargar estadísticas"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *PsychosocialListController) EquiposRemitentes(ctx *gin.Context) {
	result, err := c.svc.ListEquiposRemitentes(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al cargar equipos remitentes"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *PsychosocialListController) Duplas(ctx *gin.Context) {
	result, err := c.svc.ListDuplas(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al cargar duplas"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *PsychosocialListController) inputFromQuery(ctx *gin.Context) service.PsychosocialListInput {
	return service.PsychosocialListInput{
		FilterProfessionalID:      ctx.Query("filter_professional_id"),
		FilterDuplaID:             ctx.Query("filter_dupla_id"),
		FilterEstadoRemision:      ctx.Query("filter_estado_remision"),
		FilterSesionesCompletadas: ctx.Query("filter_sesiones_completadas"),
		FilterEquipoRemitente:     ctx.Query("filter_equipo_remitente"),
		FilterNivelRiesgo:         ctx.Query("filter_nivel_riesgo"),
		FilterNumeroIdentidad:     ctx.Query("filter_numero_identidad"),
		FilterTelefono:            ctx.Query("filter_telefono"),
		Sort:        ctx.DefaultQuery("sort", "created_at"),
		Order:       ctx.DefaultQuery("order", "desc"),
		Page:        casesListQueryInt(ctx, "page", 1),
		PageSize:    casesListQueryInt(ctx, "page_size", 20),
	}
}
