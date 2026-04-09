<<<<<<< HEAD
package controller

import (
	"bitsflow/internal/models"
	"bitsflow/salvia/service"
	"encoding/json"
	"errors"
	"net/http"
)

// FormController maneja los endpoints HTTP para Form.
=======
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

>>>>>>> e10fe9a63702dc7f74e0c51178b6ec39a2464fcb
type FormController struct {
	svc service.FormService
}

func NewFormController(svc service.FormService) *FormController {
	return &FormController{svc: svc}
}

<<<<<<< HEAD
// GET /forms?id={uuid}
func (c *FormController) GetByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "parámetro 'id' requerido"})
		return
	}
	f, err := c.svc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrFormNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "formulario no encontrado"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "error interno"})
		return
	}
	writeJSON(w, http.StatusOK, f)
}

// GET /forms/list?page={int}&limit={int}
func (c *FormController) ListHandler(w http.ResponseWriter, r *http.Request) {
	page := queryInt(r, "page", 0)
	limit := queryInt(r, "limit", 20)
	result, err := c.svc.List(r.Context(), page, limit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "error interno"})
		return
	}
	writeJSON(w, http.StatusOK, result)
}

// POST /forms
func (c *FormController) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var f models.Form
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "body inválido"})
		return
	}
	if f.Name == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "campo 'name' requerido"})
		return
	}
	if err := c.svc.Create(r.Context(), &f); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "error interno"})
		return
	}
	writeJSON(w, http.StatusCreated, f)
}

// PUT /forms?id={uuid}
func (c *FormController) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "parámetro 'id' requerido"})
		return
	}
	var f models.Form
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "body inválido"})
		return
	}
	f.ID = id
	if err := c.svc.Update(r.Context(), &f); err != nil {
		if errors.Is(err, service.ErrFormNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "formulario no encontrado"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "error interno"})
		return
	}
	writeJSON(w, http.StatusOK, f)
}

// DELETE /forms?id={uuid}
func (c *FormController) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "parámetro 'id' requerido"})
		return
	}
	if err := c.svc.Delete(r.Context(), id); err != nil {
		if errors.Is(err, service.ErrFormNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "formulario no encontrado"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "error interno"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
=======
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
>>>>>>> e10fe9a63702dc7f74e0c51178b6ec39a2464fcb
}
