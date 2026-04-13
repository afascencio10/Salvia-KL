package controller

import (
	"bitsflow/salvia/service"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ─── RepeaterGroup ────────────────────────────────────────────────────────────
// GET    /api/v1/form-sections/:id/repeater-groups
// POST   /api/v1/form-sections/:id/repeater-groups
// GET    /api/v1/repeater-groups/:id
// PUT    /api/v1/repeater-groups/:id
// DELETE /api/v1/repeater-groups/:id

type RepeaterGroupController struct{ svc service.RepeaterGroupService }

func NewRepeaterGroupController(svc service.RepeaterGroupService) *RepeaterGroupController {
	return &RepeaterGroupController{svc: svc}
}

func (c *RepeaterGroupController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/form-sections/:id/repeater-groups", c.List)
	rg.POST("/form-sections/:id/repeater-groups", c.Create)
	g := rg.Group("/repeater-groups")
	g.GET("/:id", c.GetByID)
	g.PUT("/:id", c.Update)
	g.DELETE("/:id", c.Delete)
}

func (c *RepeaterGroupController) List(ctx *gin.Context) {
	items, err := c.svc.ListByFormSectionID(ctx.Request.Context(), ctx.Param("id"))
	if err != nil { ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"}); return }
	ctx.JSON(http.StatusOK, items)
}

func (c *RepeaterGroupController) GetByID(ctx *gin.Context) {
	item, err := c.svc.GetByID(ctx.Request.Context(), ctx.Param("id"))
	if err != nil {
		if errors.Is(err, service.ErrRepeaterGroupNotFound) { ctx.JSON(http.StatusNotFound, gin.H{"error": "no encontrado"}); return }
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"}); return
	}
	ctx.JSON(http.StatusOK, item)
}

func (c *RepeaterGroupController) Create(ctx *gin.Context) {
	var body struct {
		Name           string  `json:"name"           binding:"required"`
		ItemName       *string `json:"itemName"`
		Order          int     `json:"order"`
		MinRepetitions int     `json:"minRepetitions"`
		MaxRepetitions *int    `json:"maxRepetitions"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil { ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
	item, err := c.svc.Create(ctx.Request.Context(), service.CreateRepeaterGroupInput{
		FormSectionID:  ctx.Param("id"),
		Name:           body.Name,
		ItemName:       body.ItemName,
		Order:          body.Order,
		MinRepetitions: body.MinRepetitions,
		MaxRepetitions: body.MaxRepetitions,
	})
	if err != nil { ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"}); return }
	ctx.JSON(http.StatusCreated, item)
}

func (c *RepeaterGroupController) Update(ctx *gin.Context) {
	var body struct {
		Name           *string `json:"name"`
		ItemName       *string `json:"itemName"`
		Order          *int    `json:"order"`
		MinRepetitions *int    `json:"minRepetitions"`
		MaxRepetitions *int    `json:"maxRepetitions"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil { ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
	item, err := c.svc.Update(ctx.Request.Context(), ctx.Param("id"), service.UpdateRepeaterGroupInput{
		Name:           body.Name,
		ItemName:       body.ItemName,
		Order:          body.Order,
		MinRepetitions: body.MinRepetitions,
		MaxRepetitions: body.MaxRepetitions,
	})
	if err != nil {
		if errors.Is(err, service.ErrRepeaterGroupNotFound) { ctx.JSON(http.StatusNotFound, gin.H{"error": "no encontrado"}); return }
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"}); return
	}
	ctx.JSON(http.StatusOK, item)
}

func (c *RepeaterGroupController) Delete(ctx *gin.Context) {
	if err := c.svc.Delete(ctx.Request.Context(), ctx.Param("id")); err != nil {
		if errors.Is(err, service.ErrRepeaterGroupNotFound) { ctx.JSON(http.StatusNotFound, gin.H{"error": "no encontrado"}); return }
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"}); return
	}
	ctx.JSON(http.StatusNoContent, nil)
}

// ─── VisibilityCondition ──────────────────────────────────────────────────────
// GET    /api/v1/visibility-conditions?targetId=xxx
// POST   /api/v1/visibility-conditions
// GET    /api/v1/visibility-conditions/:id
// PUT    /api/v1/visibility-conditions/:id
// DELETE /api/v1/visibility-conditions/:id

type VisibilityConditionController struct{ svc service.VisibilityConditionService }

func NewVisibilityConditionController(svc service.VisibilityConditionService) *VisibilityConditionController {
	return &VisibilityConditionController{svc: svc}
}

func (c *VisibilityConditionController) RegisterRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/visibility-conditions")
	g.GET("", c.List)
	g.POST("", c.Create)
	g.GET("/:id", c.GetByID)
	g.PUT("/:id", c.Update)
	g.DELETE("/:id", c.Delete)
}

func (c *VisibilityConditionController) List(ctx *gin.Context) {
	targetID := ctx.Query("targetId")
	if targetID == "" { ctx.JSON(http.StatusBadRequest, gin.H{"error": "query param 'targetId' requerido"}); return }
	items, err := c.svc.ListByTargetID(ctx.Request.Context(), targetID)
	if err != nil { ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"}); return }
	ctx.JSON(http.StatusOK, items)
}

func (c *VisibilityConditionController) GetByID(ctx *gin.Context) {
	item, err := c.svc.GetByID(ctx.Request.Context(), ctx.Param("id"))
	if err != nil {
		if errors.Is(err, service.ErrVisibilityConditionNotFound) { ctx.JSON(http.StatusNotFound, gin.H{"error": "no encontrado"}); return }
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"}); return
	}
	ctx.JSON(http.StatusOK, item)
}

func (c *VisibilityConditionController) Create(ctx *gin.Context) {
	var body struct {
		TargetType        string  `json:"targetType"        binding:"required"`
		TargetID          string  `json:"targetId"          binding:"required"`
		TriggerQuestionID string  `json:"triggerQuestionId" binding:"required"`
		TriggerValue      *string `json:"triggerValue"`
		Operator          string  `json:"operator"          binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil { ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
	item, err := c.svc.Create(ctx.Request.Context(), service.CreateVisibilityConditionInput{
		TargetType:        body.TargetType,
		TargetID:          body.TargetID,
		TriggerQuestionID: body.TriggerQuestionID,
		TriggerValue:      body.TriggerValue,
		Operator:          body.Operator,
	})
	if err != nil { ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"}); return }
	ctx.JSON(http.StatusCreated, item)
}

func (c *VisibilityConditionController) Update(ctx *gin.Context) {
	var body struct {
		TargetType        *string `json:"targetType"`
		TargetID          *string `json:"targetId"`
		TriggerQuestionID *string `json:"triggerQuestionId"`
		TriggerValue      *string `json:"triggerValue"`
		Operator          *string `json:"operator"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil { ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
	item, err := c.svc.Update(ctx.Request.Context(), ctx.Param("id"), service.UpdateVisibilityConditionInput{
		TargetType:        body.TargetType,
		TargetID:          body.TargetID,
		TriggerQuestionID: body.TriggerQuestionID,
		TriggerValue:      body.TriggerValue,
		Operator:          body.Operator,
	})
	if err != nil {
		if errors.Is(err, service.ErrVisibilityConditionNotFound) { ctx.JSON(http.StatusNotFound, gin.H{"error": "no encontrado"}); return }
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"}); return
	}
	ctx.JSON(http.StatusOK, item)
}

func (c *VisibilityConditionController) Delete(ctx *gin.Context) {
	if err := c.svc.Delete(ctx.Request.Context(), ctx.Param("id")); err != nil {
		if errors.Is(err, service.ErrVisibilityConditionNotFound) { ctx.JSON(http.StatusNotFound, gin.H{"error": "no encontrado"}); return }
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"}); return
	}
	ctx.JSON(http.StatusNoContent, nil)
}
