// Package controller — dupla_admin_controller.go
// Endpoints JSON para pantalla Administrar Duplas (E01, E04, E07).
package controller

import (
	"bitsflow/salvia/service"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// DuplaAdminController lecturas y escritura de duplas (Administrar Duplas).
type DuplaAdminController struct {
	svc service.DuplaAdminService
}

func NewDuplaAdminController(svc service.DuplaAdminService) *DuplaAdminController {
	return &DuplaAdminController{svc: svc}
}

func (c *DuplaAdminController) RegisterRoutes(rg *gin.RouterGroup) {
	// Rutas más específicas antes de /duplas (catálogo simple en PsychosocialListController).
	rg.GET("/duplas/profesionales", c.ListProfessionals)
	rg.GET("/duplas/activas", c.ListActiveEnriched)
	rg.POST("/duplas", c.Create)
	rg.PUT("/duplas/:id", c.Update)
	rg.DELETE("/duplas/:id", c.Delete)
}

type duplaAdminSaveRequest struct {
	Name           string `json:"name"`
	PsychologistID string `json:"psychologistId"`
	SocialWorkerID string `json:"socialWorkerId"`
}

func (c *DuplaAdminController) ListProfessionals(ctx *gin.Context) {
	result, err := c.svc.ListProfessionals(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al cargar profesionales"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *DuplaAdminController) ListActiveEnriched(ctx *gin.Context) {
	result, err := c.svc.ListActiveEnriched(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al cargar duplas"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *DuplaAdminController) Create(ctx *gin.Context) {
	var body duplaAdminSaveRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "payload inválido"})
		return
	}

	item, err := c.svc.Create(ctx.Request.Context(), service.DuplaAdminSaveInput{
		Name:           body.Name,
		PsychologistID: body.PsychologistID,
		SocialWorkerID: body.SocialWorkerID,
	})
	if err != nil {
		c.writeSaveError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, item)
}

func (c *DuplaAdminController) Update(ctx *gin.Context) {
	id := strings.TrimSpace(ctx.Param("id"))
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id requerido"})
		return
	}

	var body duplaAdminSaveRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "payload inválido"})
		return
	}

	item, err := c.svc.Update(ctx.Request.Context(), id, service.DuplaAdminSaveInput{
		Name:           body.Name,
		PsychologistID: body.PsychologistID,
		SocialWorkerID: body.SocialWorkerID,
	})
	if err != nil {
		c.writeSaveError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, item)
}

func (c *DuplaAdminController) Delete(ctx *gin.Context) {
	id := strings.TrimSpace(ctx.Param("id"))
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id requerido"})
		return
	}

	if err := c.svc.SoftDelete(ctx.Request.Context(), id); err != nil {
		switch {
		case errors.Is(err, service.ErrDuplaInUse):
			ctx.JSON(http.StatusConflict, gin.H{
				"error":   service.ErrDuplaInUse.Error(),
				"message": service.DuplaInUseMessage,
			})
		case errors.Is(err, service.ErrDuplaNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": service.ErrDuplaNotFound.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al eliminar la dupla"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"ok": true})
}

func (c *DuplaAdminController) writeSaveError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrDuplaNameInUse):
		ctx.JSON(http.StatusConflict, gin.H{"error": service.ErrDuplaNameInUse.Error()})
	case errors.Is(err, service.ErrDuplaMembersConflict):
		ctx.JSON(http.StatusConflict, gin.H{"error": service.ErrDuplaMembersConflict.Error()})
	case errors.Is(err, service.ErrDuplaInvalidMember), errors.Is(err, service.ErrDuplaValidation):
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, service.ErrDuplaNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{"error": service.ErrDuplaNotFound.Error()})
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al guardar la dupla"})
	}
}
