package controller

import (
	"bitsflow/salvia/service"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ─── Option ───────────────────────────────────────────────────────────────────
// GET    /api/v1/questions/:id/options
// POST   /api/v1/questions/:id/options
// GET    /api/v1/options/:id
// PUT    /api/v1/options/:id
// DELETE /api/v1/options/:id

type OptionController struct{ svc service.OptionService }

func NewOptionController(svc service.OptionService) *OptionController {
	return &OptionController{svc: svc}
}

func (c *OptionController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/questions/:id/options", c.ListByQuestion)
	rg.POST("/questions/:id/options", c.Create)
	g := rg.Group("/options")
	g.GET("/:id", c.GetByID)
	g.PUT("/:id", c.Update)
	g.DELETE("/:id", c.Delete)
}

func (c *OptionController) ListByQuestion(ctx *gin.Context) {
	items, err := c.svc.ListByQuestionID(ctx.Request.Context(), ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}
	ctx.JSON(http.StatusOK, items)
}

func (c *OptionController) GetByID(ctx *gin.Context) {
	item, err := c.svc.GetByID(ctx.Request.Context(), ctx.Param("id"))
	if err != nil {
		if errors.Is(err, service.ErrOptionNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "opción no encontrada"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}
	ctx.JSON(http.StatusOK, item)
}

func (c *OptionController) Create(ctx *gin.Context) {
	var body struct {
		Label string `json:"label" binding:"required"`
		Value string `json:"value" binding:"required"`
		Order int    `json:"order"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	o, err := c.svc.Create(ctx.Request.Context(), service.CreateOptionInput{
		QuestionID: ctx.Param("id"),
		Label:      body.Label,
		Value:      body.Value,
		Order:      body.Order,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}
	ctx.JSON(http.StatusCreated, o)
}

func (c *OptionController) Update(ctx *gin.Context) {
	var body struct {
		Label *string `json:"label"`
		Value *string `json:"value"`
		Order *int    `json:"order"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	o, err := c.svc.Update(ctx.Request.Context(), ctx.Param("id"), service.UpdateOptionInput{
		Label: body.Label,
		Value: body.Value,
		Order: body.Order,
	})
	if err != nil {
		if errors.Is(err, service.ErrOptionNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "opción no encontrada"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}
	ctx.JSON(http.StatusOK, o)
}

func (c *OptionController) Delete(ctx *gin.Context) {
	if err := c.svc.Delete(ctx.Request.Context(), ctx.Param("id")); err != nil {
		if errors.Is(err, service.ErrOptionNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "opción no encontrada"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}
