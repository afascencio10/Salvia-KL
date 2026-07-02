// Package controller — agents_search_controller.go
// Endpoint JSON para autocomplete de persona asignada (evento E-07).
package controller

import (
	"bitsflow/salvia/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AgentsSearchController maneja GET /api/v1/agents/search.
type AgentsSearchController struct {
	svc service.AgentsSearchService
}

func NewAgentsSearchController(svc service.AgentsSearchService) *AgentsSearchController {
	return &AgentsSearchController{svc: svc}
}

// RegisterRoutes registra la ruta de búsqueda de agentes.
//
//	GET /api/v1/agents/search?q=&limit=
func (c *AgentsSearchController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/agents/search", c.Search)
	rg.GET("/agents/search-psicosocial", c.SearchPsicosocial)
}

// Search devuelve agentes cuyo nombre o apellido coincide con q (mín. 3 caracteres).
func (c *AgentsSearchController) Search(ctx *gin.Context) {
	limit := casesListQueryInt(ctx, "limit", 10)

	result, err := c.svc.Search(ctx.Request.Context(), ctx.Query("q"), limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al buscar agentes"})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// SearchPsicosocial devuelve profesionales activos del equipo psicosocial (E-07).
func (c *AgentsSearchController) SearchPsicosocial(ctx *gin.Context) {
	limit := casesListQueryInt(ctx, "limit", 10)

	result, err := c.svc.SearchPsicosocial(ctx.Request.Context(), ctx.Query("q"), limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al buscar profesionales psicosociales"})
		return
	}

	ctx.JSON(http.StatusOK, result)
}
