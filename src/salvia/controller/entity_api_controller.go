// Package controller — EntityAPIController: catálogo de organizaciones y ciudades por entidad.
package controller

import (
	"bitsflow/common/utils"
	"bitsflow/salvia/service"
	salvia_config "bitsflow/salvia/config"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// EntityAPIController expone catálogos de salvia.entity para Casos Entidad.
type EntityAPIController struct {
	svc service.EntityCaseService
}

func NewEntityAPIController(svc service.EntityCaseService) *EntityAPIController {
	return &EntityAPIController{svc: svc}
}

func (c *EntityAPIController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/entities", c.List)
	rg.GET("/entities/:entityId/cities", c.ListCities)
}

// requireCasosEntidadAccess valida sesión + permiso get_casos_entidad (rol et).
func requireCasosEntidadAccess(ctx *gin.Context) *utils.CommonSession {
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
	if !utils.CheckPermission(salvia_config.PermissionsByRole, "get_casos_entidad", s.CurrentRole, ctx) {
		return nil
	}
	return s
}

// List — GET /api/v1/entities
func (c *EntityAPIController) List(ctx *gin.Context) {
	if requireCasosEntidadAccess(ctx) == nil {
		return
	}
	items, err := c.svc.ListEntities(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}
	ctx.JSON(http.StatusOK, items)
}

// ListCities — GET /api/v1/entities/:entityId/cities
func (c *EntityAPIController) ListCities(ctx *gin.Context) {
	if requireCasosEntidadAccess(ctx) == nil {
		return
	}
	entityID, err := strconv.ParseInt(ctx.Param("entityId"), 10, 64)
	if err != nil || entityID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "entityId inválido"})
		return
	}
	items, err := c.svc.ListCitiesByEntity(ctx.Request.Context(), entityID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}
	ctx.JSON(http.StatusOK, items)
}
