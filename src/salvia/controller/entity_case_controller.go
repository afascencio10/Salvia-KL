// Package controller expone los endpoints HTTP de EntityCase (componente case-entities).
package controller

import (
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	"bitsflow/salvia/service"
	"errors"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// EntityCaseController maneja las rutas REST de EntityCase — entidades
// relacionadas manualmente con un caso ("Gestión institucional").
type EntityCaseController struct {
	svc service.EntityCaseService
}

func NewEntityCaseController(svc service.EntityCaseService) *EntityCaseController {
	return &EntityCaseController{svc: svc}
}

// entityCaseWriteRoles — roles que pueden agregar una entidad al caso.
// Más restrictivo que el acceso de lectura: cualquiera que vea la pantalla de
// Detalle del Caso puede ver/filtrar entidades, pero solo sv/op/ro pueden agregar
// (ver DocsMD/Componentes/case-entities/case-entities-interface.md).
var entityCaseWriteRoles = map[string]bool{"sv": true, "op": true, "ro": true}

// requireCaseDetailAccess valida la sesión igual que las fachadas HTML
// (sessions.Default + utils.GetCommonSession) y exige el mismo permiso que ya
// protege la pantalla de Detalle del Caso ("get_case_detail_sv" en
// salvia_config.PermissionsByRole) — no una lista de roles propia, para que
// ambos se mantengan sincronizados automáticamente. Responde 401 JSON en vez
// de redirigir — el resto de /api/v1 en este repo no hace ningún chequeo de
// sesión (hallazgo reportado aparte); estos dos endpoints nuevos no replican
// ese hueco a sabiendas. Devuelve la sesión resuelta, o nil si ya respondió.
func requireCaseDetailAccess(ctx *gin.Context) *utils.CommonSession {
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
	if !utils.CheckPermission(salvia_config.PermissionsByRole, "get_case_detail_sv", s.CurrentRole, ctx) {
		return nil // CheckPermission ya escribió el 401
	}
	return s
}

// RegisterRoutes registra las rutas de EntityCase en el grupo /api/v1.
//
// Usa ":id" (no ":caseId") para el segmento de caso porque Gin exige que el
// nombre del wildcard sea el mismo en todas las rutas que comparten la misma
// posición del árbol — CaseDetailController ya registró "/casos/:id/..." en
// este mismo grupo; un nombre distinto aquí provoca panic en el arranque.
//
//	GET  /api/v1/casos/:id/entidades   → List
//	POST /api/v1/casos/:id/entidades   → Create
//	GET  /api/v1/entity-cases          → ListByEntity (Casos Entidad)
func (c *EntityCaseController) RegisterRoutes(rg *gin.RouterGroup) {
	casos := rg.Group("/casos/:id/entidades")
	casos.GET("", c.List)
	casos.POST("", c.Create)

	rg.GET("/entity-cases", c.ListByEntity)
}

// ListByEntity — GET /api/v1/entity-cases?document=&page=&pageSize=
// Listado paginado de Casos Entidad (rol et). Filtra por entity_branch_id de sesión.
func (c *EntityCaseController) ListByEntity(ctx *gin.Context) {
	s := requireCasosEntidadAccess(ctx)
	if s == nil {
		return
	}
	if s.EntityBranchId <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "el usuario no tiene sede (entity_branch_id) asignada"})
		return
	}

	page := ginQueryInt(ctx, "page", 0)
	pageSize := ginQueryInt(ctx, "pageSize", 5)
	if pageSize <= 0 {
		pageSize = 5
	}

	result, err := c.svc.ListByEntity(ctx.Request.Context(), service.EntityCaseListFilter{
		EntityBranchID: s.EntityBranchId,
		Document:       ctx.Query("document"),
		Page:           page,
		PageSize:       pageSize,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "la sede asignada no existe o no tiene entidad padre"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

// List devuelve las entidades relacionadas con un caso.
//
//	GET /api/v1/casos/:id/entidades
func (c *EntityCaseController) List(ctx *gin.Context) {
	if requireCaseDetailAccess(ctx) == nil {
		return
	}
	caseID := ctx.Param("id")
	items, err := c.svc.ListByCase(ctx.Request.Context(), caseID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}
	ctx.JSON(http.StatusOK, items)
}

// Create asocia una sede de entidad (entity_branch) al caso.
//
//	POST /api/v1/casos/:id/entidades
//	Body: { "entityBranchId": 12, "objetivo": "...", "createdById": "icode" }
func (c *EntityCaseController) Create(ctx *gin.Context) {
	s := requireCaseDetailAccess(ctx)
	if s == nil {
		return
	}
	if !entityCaseWriteRoles[s.CurrentRole] {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "su rol no tiene permiso para agregar entidades"})
		return
	}
	caseID := ctx.Param("id")

	var body struct {
		EntityBranchID int64   `json:"entityBranchId" binding:"required"`
		Objetivo       *string `json:"objetivo"`
		CreatedByID    string  `json:"createdById" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ec, err := c.svc.Create(ctx.Request.Context(), service.CreateEntityCaseInput{
		CaseID:         caseID,
		EntityBranchID: body.EntityBranchID,
		Objetivo:       body.Objetivo,
		CreatedByID:    body.CreatedByID,
	})
	if err != nil {
		if errors.Is(err, service.ErrEntityCaseDuplicate) {
			ctx.JSON(http.StatusConflict, gin.H{"error": "Esta entidad ya está asociada a este caso"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}
	ctx.JSON(http.StatusCreated, ec)
}
