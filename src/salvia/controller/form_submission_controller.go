package controller

import (
	"bitsflow/salvia/service"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ─── FormSubmission ───────────────────────────────────────────────────────────
// GET    /api/v1/forms/:id/submissions
// POST   /api/v1/forms/:id/submissions
// GET    /api/v1/form-submissions/:id
// PUT    /api/v1/form-submissions/:id
// DELETE /api/v1/form-submissions/:id

type FormSubmissionController struct{ svc service.FormSubmissionService }

func NewFormSubmissionController(svc service.FormSubmissionService) *FormSubmissionController {
	return &FormSubmissionController{svc: svc}
}

func (c *FormSubmissionController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/forms/:id/submissions", c.ListByForm)
	rg.POST("/forms/:id/submissions", c.Create)
	g := rg.Group("/form-submissions")
	g.GET("/:id", c.GetByID)
	g.PUT("/:id", c.Update)
	g.DELETE("/:id", c.Delete)
}

func (c *FormSubmissionController) ListByForm(ctx *gin.Context) {
	items, err := c.svc.ListByFormID(ctx.Request.Context(), ctx.Param("id"))
	if err != nil { ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"}); return }
	ctx.JSON(http.StatusOK, items)
}

func (c *FormSubmissionController) GetByID(ctx *gin.Context) {
	item, err := c.svc.GetByID(ctx.Request.Context(), ctx.Param("id"))
	if err != nil {
		if errors.Is(err, service.ErrFormSubmissionNotFound) { ctx.JSON(http.StatusNotFound, gin.H{"error": "no encontrado"}); return }
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"}); return
	}
	ctx.JSON(http.StatusOK, item)
}

func (c *FormSubmissionController) Create(ctx *gin.Context) {
	item, err := c.svc.Create(ctx.Request.Context(), ctx.Param("id"))
	if err != nil { ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"}); return }
	ctx.JSON(http.StatusCreated, item)
}

func (c *FormSubmissionController) Update(ctx *gin.Context) {
	var body struct {
		FormID *string `json:"formId"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil { ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
	item, err := c.svc.Update(ctx.Request.Context(), ctx.Param("id"), service.UpdateFormSubmissionInput{FormID: body.FormID})
	if err != nil {
		if errors.Is(err, service.ErrFormSubmissionNotFound) { ctx.JSON(http.StatusNotFound, gin.H{"error": "no encontrado"}); return }
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"}); return
	}
	ctx.JSON(http.StatusOK, item)
}

func (c *FormSubmissionController) Delete(ctx *gin.Context) {
	if err := c.svc.Delete(ctx.Request.Context(), ctx.Param("id")); err != nil {
		if errors.Is(err, service.ErrFormSubmissionNotFound) { ctx.JSON(http.StatusNotFound, gin.H{"error": "no encontrado"}); return }
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"}); return
	}
	ctx.JSON(http.StatusNoContent, nil)
}

// ─── RepeaterEntry ────────────────────────────────────────────────────────────
// GET    /api/v1/form-submissions/:id/repeater-entries
// POST   /api/v1/form-submissions/:id/repeater-entries
// GET    /api/v1/repeater-entries/:id
// PUT    /api/v1/repeater-entries/:id
// DELETE /api/v1/repeater-entries/:id

type RepeaterEntryController struct{ svc service.RepeaterEntryService }

func NewRepeaterEntryController(svc service.RepeaterEntryService) *RepeaterEntryController {
	return &RepeaterEntryController{svc: svc}
}

func (c *RepeaterEntryController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/form-submissions/:id/repeater-entries", c.List)
	rg.POST("/form-submissions/:id/repeater-entries", c.Create)
	g := rg.Group("/repeater-entries")
	g.GET("/:id", c.GetByID)
	g.PUT("/:id", c.Update)
	g.DELETE("/:id", c.Delete)
}

func (c *RepeaterEntryController) List(ctx *gin.Context) {
	items, err := c.svc.ListBySubmissionID(ctx.Request.Context(), ctx.Param("id"))
	if err != nil { ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"}); return }
	ctx.JSON(http.StatusOK, items)
}

func (c *RepeaterEntryController) GetByID(ctx *gin.Context) {
	item, err := c.svc.GetByID(ctx.Request.Context(), ctx.Param("id"))
	if err != nil {
		if errors.Is(err, service.ErrRepeaterEntryNotFound) { ctx.JSON(http.StatusNotFound, gin.H{"error": "no encontrado"}); return }
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"}); return
	}
	ctx.JSON(http.StatusOK, item)
}

func (c *RepeaterEntryController) Create(ctx *gin.Context) {
	var body struct {
		RepeaterGroupID string `json:"repeaterGroupId" binding:"required"`
		Iteration       int    `json:"iteration"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil { ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
	item, err := c.svc.Create(ctx.Request.Context(), service.CreateRepeaterEntryInput{
		FormSubmissionID: ctx.Param("id"),
		RepeaterGroupID:  body.RepeaterGroupID,
		Iteration:        body.Iteration,
	})
	if err != nil { ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"}); return }
	ctx.JSON(http.StatusCreated, item)
}

func (c *RepeaterEntryController) Update(ctx *gin.Context) {
	var body struct {
		Iteration *int `json:"iteration"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil { ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
	item, err := c.svc.Update(ctx.Request.Context(), ctx.Param("id"), service.UpdateRepeaterEntryInput{Iteration: body.Iteration})
	if err != nil {
		if errors.Is(err, service.ErrRepeaterEntryNotFound) { ctx.JSON(http.StatusNotFound, gin.H{"error": "no encontrado"}); return }
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"}); return
	}
	ctx.JSON(http.StatusOK, item)
}

func (c *RepeaterEntryController) Delete(ctx *gin.Context) {
	if err := c.svc.Delete(ctx.Request.Context(), ctx.Param("id")); err != nil {
		if errors.Is(err, service.ErrRepeaterEntryNotFound) { ctx.JSON(http.StatusNotFound, gin.H{"error": "no encontrado"}); return }
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"}); return
	}
	ctx.JSON(http.StatusNoContent, nil)
}

// ─── Answer ───────────────────────────────────────────────────────────────────
// GET    /api/v1/form-submissions/:id/answers
// POST   /api/v1/answers
// GET    /api/v1/answers/:id
// PUT    /api/v1/answers/:id
// DELETE /api/v1/answers/:id

type AnswerController struct{ svc service.AnswerService }

func NewAnswerController(svc service.AnswerService) *AnswerController {
	return &AnswerController{svc: svc}
}

func (c *AnswerController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/form-submissions/:id/answers", c.ListBySubmission)
	g := rg.Group("/answers")
	g.POST("", c.Create)
	g.GET("/:id", c.GetByID)
	g.PUT("/:id", c.Update)
	g.DELETE("/:id", c.Delete)
}

func (c *AnswerController) ListBySubmission(ctx *gin.Context) {
	items, err := c.svc.ListBySubmissionID(ctx.Request.Context(), ctx.Param("id"))
	if err != nil { ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"}); return }
	ctx.JSON(http.StatusOK, items)
}

func (c *AnswerController) GetByID(ctx *gin.Context) {
	item, err := c.svc.GetByID(ctx.Request.Context(), ctx.Param("id"))
	if err != nil {
		if errors.Is(err, service.ErrAnswerNotFound) { ctx.JSON(http.StatusNotFound, gin.H{"error": "no encontrado"}); return }
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"}); return
	}
	ctx.JSON(http.StatusOK, item)
}

func (c *AnswerController) Create(ctx *gin.Context) {
	var body struct {
		FormSubmissionID string  `json:"formSubmissionId" binding:"required"`
		QuestionID       string  `json:"questionId"       binding:"required"`
		RepeaterEntryID  *string `json:"repeaterEntryId"`
		Value            *string `json:"value"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil { ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
	item, err := c.svc.Create(ctx.Request.Context(), service.CreateAnswerInput{
		FormSubmissionID: body.FormSubmissionID,
		QuestionID:       body.QuestionID,
		RepeaterEntryID:  body.RepeaterEntryID,
		Value:            body.Value,
	})
	if err != nil { ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"}); return }
	ctx.JSON(http.StatusCreated, item)
}

func (c *AnswerController) Update(ctx *gin.Context) {
	var body struct {
		Value *string `json:"value"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil { ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
	item, err := c.svc.Update(ctx.Request.Context(), ctx.Param("id"), service.UpdateAnswerInput{Value: body.Value})
	if err != nil {
		if errors.Is(err, service.ErrAnswerNotFound) { ctx.JSON(http.StatusNotFound, gin.H{"error": "no encontrado"}); return }
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"}); return
	}
	ctx.JSON(http.StatusOK, item)
}

func (c *AnswerController) Delete(ctx *gin.Context) {
	if err := c.svc.Delete(ctx.Request.Context(), ctx.Param("id")); err != nil {
		if errors.Is(err, service.ErrAnswerNotFound) { ctx.JSON(http.StatusNotFound, gin.H{"error": "no encontrado"}); return }
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"}); return
	}
	ctx.JSON(http.StatusNoContent, nil)
}
