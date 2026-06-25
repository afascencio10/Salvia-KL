// Package controller — cases_list_controller.go
// Endpoint JSON para el componente casos-component (evento E-01 y filtros).
package controller

import (
	"bitsflow/salvia/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CasesListController maneja GET /api/v1/cases/list.
type CasesListController struct {
	svc service.CasesListService
}

func NewCasesListController(svc service.CasesListService) *CasesListController {
	return &CasesListController{svc: svc}
}

// RegisterRoutes registra la ruta del listado de casos.
//
//	GET /api/v1/cases/list?filter_key=&filter_value=&chip_filter=
//	    &filter_riesgo=&filter_equipo=&filter_seguimientos_ejecutados=&filter_estado_caso=
//	    &filter_barreras_activas=
//	    &dropdown_filter_key=&dropdown_filter_value=&search=&sort=&order=&page=&page_size=
func (c *CasesListController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/cases/list", c.List)
}

// List devuelve casos paginados con filtros para casos-component.
func (c *CasesListController) List(ctx *gin.Context) {
	page := casesListQueryInt(ctx, "page", 1)
	pageSize := casesListQueryInt(ctx, "page_size", 20)

	result, err := c.svc.List(ctx.Request.Context(), service.CasesListInput{
		FilterKey:           ctx.Query("filter_key"),
		FilterValue:         ctx.Query("filter_value"),
		ChipFilter:          ctx.Query("chip_filter"),
		DropdownFilterKey:            ctx.Query("dropdown_filter_key"),
		DropdownFilterValue:          ctx.Query("dropdown_filter_value"),
		FilterRiesgo:                 ctx.Query("filter_riesgo"),
		FilterEquipo:                 ctx.Query("filter_equipo"),
		FilterSeguimientosEjecutados: ctx.Query("filter_seguimientos_ejecutados"),
		FilterEstadoCaso:             ctx.Query("filter_estado_caso"),
		FilterBarrerasActivas:        ctx.Query("filter_barreras_activas"),
		Search:              ctx.Query("search"),
		Sort:        ctx.DefaultQuery("sort", "registration_date"),
		Order:       ctx.DefaultQuery("order", "desc"),
		Page:        page,
		PageSize:    pageSize,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al cargar los casos"})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func casesListQueryInt(ctx *gin.Context, key string, defaultVal int) int {
	raw := ctx.Query(key)
	if raw == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(raw)
	if err != nil {
		return defaultVal
	}
	return val
}
