package controller

import (
	"bitsflow/internal/models"
	"bitsflow/salvia/service"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ─── FormSection ──────────────────────────────────────────────────────────────
// GET    /api/v1/forms/:id/sections
// POST   /api/v1/forms/:id/sections
// GET    /api/v1/form-sections/:id
// PUT    /api/v1/form-sections/:id
// DELETE /api/v1/form-sections/:id

type FormSectionController struct{ svc service.FormSectionService }

func NewFormSectionController(svc service.FormSectionService) *FormSectionController {
	return &FormSectionController{svc: svc}
}

func (c *FormSectionController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/forms/:id/sections", c.ListByForm)
	rg.POST("/forms/:id/sections", c.Create)
	g := rg.Group("/form-sections")
	g.GET("/:id", c.GetByID)
	g.PUT("/:id", c.Update)
	g.DELETE("/:id", c.Delete)
}

func (c *FormSectionController) ListByForm(ctx *gin.Context) {
	items, err := c.svc.ListByFormID(ctx.Request.Context(), ctx.Param("id"))
	if err != nil { ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"}); return }
	ctx.JSON(http.StatusOK, items)
}

func (c *FormSectionController) GetByID(ctx *gin.Context) {
	item, err := c.svc.GetByID(ctx.Request.Context(), ctx.Param("id"))
	if err != nil {
		if errors.Is(err, service.ErrFormSectionNotFound) { ctx.JSON(http.StatusNotFound, gin.H{"error": "sección no encontrada"}); return }
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"}); return
	}
	ctx.JSON(http.StatusOK, item)
}

func (c *FormSectionController) Create(ctx *gin.Context) {
	var body struct {
		Name       string `json:"name"       binding:"required"`
		OrderIndex int    `json:"orderIndex"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil { ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
	fs := &models.FormSection{FormID: ctx.Param("id"), Name: body.Name, OrderIndex: body.OrderIndex}
	if err := c.svc.Create(ctx.Request.Context(), fs); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"}); return
	}
	ctx.JSON(http.StatusCreated, fs)
}

func (c *FormSectionController) Update(ctx *gin.Context) {
	fs, err := c.svc.GetByID(ctx.Request.Context(), ctx.Param("id"))
	if err != nil {
		if errors.Is(err, service.ErrFormSectionNotFound) { ctx.JSON(http.StatusNotFound, gin.H{"error": "sección no encontrada"}); return }
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"}); return
	}
	var body struct {
		Name       *string `json:"name"`
		OrderIndex *int    `json:"orderIndex"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil { ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
	if body.Name != nil       { fs.Name = *body.Name }
	if body.OrderIndex != nil { fs.OrderIndex = *body.OrderIndex }
	if err := c.svc.Update(ctx.Request.Context(), fs); err != nil {
		if errors.Is(err, service.ErrFormSectionNotFound) { ctx.JSON(http.StatusNotFound, gin.H{"error": "sección no encontrada"}); return }
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"}); return
	}
	ctx.JSON(http.StatusOK, fs)
}

func (c *FormSectionController) Delete(ctx *gin.Context) {
	if err := c.svc.Delete(ctx.Request.Context(), ctx.Param("id")); err != nil {
		if errors.Is(err, service.ErrFormSectionNotFound) { ctx.JSON(http.StatusNotFound, gin.H{"error": "sección no encontrada"}); return }
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"}); return
	}
	ctx.JSON(http.StatusNoContent, nil)
}
