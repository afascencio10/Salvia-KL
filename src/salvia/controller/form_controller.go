// Package controller expone los endpoints HTTP de Form y FormSection usando Gin.
package controller

import (
	"bitsflow/common/utils"
	"bitsflow/salvia/service"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
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
//	GET    /api/v1/forms/:id/load?submissionId=<optional>
func (c *FormController) RegisterRoutes(rg *gin.RouterGroup) {
	forms := rg.Group("/forms")
	forms.GET("", c.List)
	forms.GET("/:id", c.GetByID)
	forms.POST("", c.Create)
	forms.PUT("/:id", c.Update)
	forms.DELETE("/:id", c.Delete)
	forms.GET("/:id/load", c.LoadForm)
	forms.POST("/saveSection", c.SaveSection)

	rg.GET("/testEndpoint", c.TestEndpoint)
}

// TestEndpoint es un sandbox de pruebas.
//   ?fn=getFormStructure&id=<formId>
//   ?fn=getFormSubmission&id=<submissionId>
//   ?fn=listSubmissions&id=<formId>
//   ?fn=validateAnswers&submissionId=<submissionId>
//   ?fn=checkVisibility&id=<formId>&submissionId=<submissionId>
//   (sin fn) → default de TestFunction
func (c *FormController) TestEndpoint(ctx *gin.Context) {
	fn := ctx.Query("fn")
	id := ctx.Query("id")
	submissionID := ctx.Query("submissionId")

	var (
		result interface{}
		err    error
	)

	switch fn {
	case "getFormStructure":
		result, err = c.svc.GetFormStructure(ctx.Request.Context(), id)
	case "getFormSubmission":
		result, err = c.svc.GetFormSubmission(ctx.Request.Context(), submissionID)
	default:
		result, err = c.svc.TestFunction(ctx.Request.Context(), fn, id, submissionID)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, result)
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

// LoadForm carga la estructura del form, la submission si existe, y la sección actual.
// GET /api/v1/forms/:id/load?submissionId=<optional>
func (c *FormController) LoadForm(ctx *gin.Context) {
	result, err := c.svc.LoadForm(ctx.Request.Context(), ctx.Param("id"), ctx.Query("submissionId"))
	if err != nil {
		if errors.Is(err, service.ErrFormNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "formulario no encontrado"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

// SaveSection guarda las respuestas de una sección y retorna el LoadFormResult actualizado.
// POST /api/v1/forms/saveSection
func (c *FormController) SaveSection(ctx *gin.Context) {
	var body service.SaveSectionInput
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if body.FormID == "" || body.FormSectionID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "formId y formSectionId son requeridos"})
		return
	}

	// Inyectar actor desde la sesión del servidor
	if sessionID, ok := sessions.Default(ctx).Get("userData").(string); ok && sessionID != "" {
		if s, err := utils.GetCommonSession(sessionID); err == nil {
			body.ActorID = s.UserICode
		}
	}

	result, err := c.svc.SaveSection(ctx.Request.Context(), body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, result)
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
