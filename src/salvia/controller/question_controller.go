package controller

import (
	"bitsflow/internal/models"
	"bitsflow/salvia/service"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ─── Question ─────────────────────────────────────────────────────────────────
// GET    /api/v1/forms/:id/questions
// POST   /api/v1/forms/:id/questions
// GET    /api/v1/questions/:id
// PUT    /api/v1/questions/:id
// DELETE /api/v1/questions/:id

type QuestionController struct{ svc service.QuestionService }

func NewQuestionController(svc service.QuestionService) *QuestionController {
	return &QuestionController{svc: svc}
}

func (c *QuestionController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/forms/:id/questions", c.ListByForm)
	rg.POST("/forms/:id/questions", c.Create)
	g := rg.Group("/questions")
	g.GET("/:id", c.GetByID)
	g.PUT("/:id", c.Update)
	g.DELETE("/:id", c.Delete)
}

func (c *QuestionController) ListByForm(ctx *gin.Context) {
	items, err := c.svc.ListByFormID(ctx.Request.Context(), ctx.Param("id"))
	if err != nil { ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"}); return }
	ctx.JSON(http.StatusOK, items)
}

func (c *QuestionController) GetByID(ctx *gin.Context) {
	item, err := c.svc.GetByID(ctx.Request.Context(), ctx.Param("id"))
	if err != nil {
		if errors.Is(err, service.ErrQuestionNotFound) { ctx.JSON(http.StatusNotFound, gin.H{"error": "pregunta no encontrada"}); return }
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"}); return
	}
	ctx.JSON(http.StatusOK, item)
}

func (c *QuestionController) Create(ctx *gin.Context) {
	var body struct {
		FormSectionID   string  `json:"formSectionId"  binding:"required"`
		RepeaterGroupID *string `json:"repeaterGroupId"`
		QuestionTypeID  string  `json:"questionTypeId" binding:"required"`
		Description     string  `json:"description"    binding:"required"`
		Required        bool    `json:"required"`
		Order           int     `json:"order"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil { ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
	q := &models.Question{
		FormID:          ctx.Param("id"),
		FormSectionID:   body.FormSectionID,
		RepeaterGroupID: body.RepeaterGroupID,
		QuestionTypeID:  body.QuestionTypeID,
		Description:     body.Description,
		Required:        body.Required,
		Order:           body.Order,
	}
	if err := c.svc.Create(ctx.Request.Context(), q); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"}); return
	}
	ctx.JSON(http.StatusCreated, q)
}

func (c *QuestionController) Update(ctx *gin.Context) {
	q, err := c.svc.GetByID(ctx.Request.Context(), ctx.Param("id"))
	if err != nil {
		if errors.Is(err, service.ErrQuestionNotFound) { ctx.JSON(http.StatusNotFound, gin.H{"error": "pregunta no encontrada"}); return }
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"}); return
	}
	var body struct {
		QuestionTypeID *string `json:"questionTypeId"`
		Description    *string `json:"description"`
		Required       *bool   `json:"required"`
		Order          *int    `json:"order"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil { ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
	if body.QuestionTypeID != nil { q.QuestionTypeID = *body.QuestionTypeID }
	if body.Description != nil    { q.Description = *body.Description }
	if body.Required != nil       { q.Required = *body.Required }
	if body.Order != nil          { q.Order = *body.Order }
	if err := c.svc.Update(ctx.Request.Context(), q); err != nil {
		if errors.Is(err, service.ErrQuestionNotFound) { ctx.JSON(http.StatusNotFound, gin.H{"error": "pregunta no encontrada"}); return }
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"}); return
	}
	ctx.JSON(http.StatusOK, q)
}

func (c *QuestionController) Delete(ctx *gin.Context) {
	if err := c.svc.Delete(ctx.Request.Context(), ctx.Param("id")); err != nil {
		if errors.Is(err, service.ErrQuestionNotFound) { ctx.JSON(http.StatusNotFound, gin.H{"error": "pregunta no encontrada"}); return }
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"}); return
	}
	ctx.JSON(http.StatusNoContent, nil)
}
