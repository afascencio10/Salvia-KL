// Package controller expone los endpoints HTTP de Form y FormSection usando Gin.
package controller

import (
	"bitsflow/salvia/service"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ─── Form ─────────────────────────────────────────────────────────────────────

type FormController struct {
	svc service.FormService
}

func NewFormController(svc service.FormService) *FormController {
	return &FormController{svc: svc}
}

// RegisterRoutes registra las rutas de Form en el grupo /api/v1.
//
//	GET    /api/v1/forms
//	GET    /api/v1/forms/:id
//	POST   /api/v1/forms
//	PUT    /api/v1/forms/:id
//	DELETE /api/v1/forms/:id
func (c *FormController) RegisterRoutes(rg *gin.RouterGroup) {
	forms := rg.Group("/forms")
	forms.GET("", c.List)
	forms.GET("/:id", c.GetByID)
	forms.POST("", c.Create)
	forms.PUT("/:id", c.Update)
	forms.DELETE("/:id", c.Delete)
}

func (c *FormController) List(ctx *gin.Context) {
	page := ginQueryInt(ctx, "page", 0)
	limit := ginQueryInt(ctx, "limit", 20)

	result, err := c.svc.List(ctx.Request.Context(), page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *FormController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	form, err := c.svc.GetByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrFormNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "formulario no encontrado"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}
	ctx.JSON(http.StatusOK, form)
}

func (c *FormController) Create(ctx *gin.Context) {
	var body struct {
		Name        string `json:"name"        binding:"required"`
		Description string `json:"description" binding:"required"`
		Status      string `json:"status"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	form, err := c.svc.Create(ctx.Request.Context(), service.CreateFormInput{
		Name:        body.Name,
		Description: body.Description,
		Status:      body.Status,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}
	ctx.JSON(http.StatusCreated, form)
}

func (c *FormController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var body struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
		Status      *string `json:"status"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	form, err := c.svc.Update(ctx.Request.Context(), id, service.UpdateFormInput{
		Name:        body.Name,
		Description: body.Description,
		Status:      body.Status,
	})
	if err != nil {
		if errors.Is(err, service.ErrFormNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "formulario no encontrado"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}
	ctx.JSON(http.StatusOK, form)
}

func (c *FormController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.svc.Delete(ctx.Request.Context(), id); err != nil {
		if errors.Is(err, service.ErrFormNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "formulario no encontrado"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}

// ─── FormSection ──────────────────────────────────────────────────────────────

type FormSectionController struct {
	svc service.FormSectionService
}

func NewFormSectionController(svc service.FormSectionService) *FormSectionController {
	return &FormSectionController{svc: svc}
}

// RegisterRoutes registra las rutas de FormSection en el grupo /api/v1.
//
//	GET    /api/v1/forms/:formId/sections
//	GET    /api/v1/form-sections/:id
//	POST   /api/v1/forms/:formId/sections
//	PUT    /api/v1/form-sections/:id
//	DELETE /api/v1/form-sections/:id
func (c *FormSectionController) RegisterRoutes(rg *gin.RouterGroup) {
	// Rutas anidadas bajo /forms/:id — Gin exige el mismo nombre de wildcard que la ruta padre
	rg.GET("/forms/:id/sections", c.ListByFormID)
	rg.POST("/forms/:id/sections", c.Create)

	// Rutas standalone por ID de sección
	sections := rg.Group("/form-sections")
	sections.GET("/:id", c.GetByID)
	sections.PUT("/:id", c.Update)
	sections.DELETE("/:id", c.Delete)
}

func (c *FormSectionController) ListByFormID(ctx *gin.Context) {
	formID := ctx.Param("id")
	sections, err := c.svc.ListByFormID(ctx.Request.Context(), formID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}
	ctx.JSON(http.StatusOK, sections)
}

func (c *FormSectionController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	section, err := c.svc.GetByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrFormSectionNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "sección no encontrada"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}
	ctx.JSON(http.StatusOK, section)
}

func (c *FormSectionController) Create(ctx *gin.Context) {
	formID := ctx.Param("id")
	var body struct {
		Name        string  `json:"name"        binding:"required"`
		Description *string `json:"description"`
		Order       int     `json:"order"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	section, err := c.svc.Create(ctx.Request.Context(), service.CreateFormSectionInput{
		FormID:      formID,
		Name:        body.Name,
		Description: body.Description,
		Order:       body.Order,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}
	ctx.JSON(http.StatusCreated, section)
}

func (c *FormSectionController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var body struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
		Order       *int    `json:"order"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	section, err := c.svc.Update(ctx.Request.Context(), id, service.UpdateFormSectionInput{
		Name:        body.Name,
		Description: body.Description,
		Order:       body.Order,
	})
	if err != nil {
		if errors.Is(err, service.ErrFormSectionNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "sección no encontrada"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}
	ctx.JSON(http.StatusOK, section)
}

func (c *FormSectionController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.svc.Delete(ctx.Request.Context(), id); err != nil {
		if errors.Is(err, service.ErrFormSectionNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "sección no encontrada"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}

// ─── helpers ──────────────────────────────────────────────────────────────────

func ginQueryInt(ctx *gin.Context, key string, defaultVal int) int {
	raw := ctx.Query(key)
	if raw == "" {
		return defaultVal
	}
	v, err := strconv.Atoi(raw)
	if err != nil || v < 0 {
		return defaultVal
	}
	return v
}
