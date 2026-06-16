package controller

import (
	"bitsflow/internal/repository"
	"bitsflow/salvia/service"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DirectoryController struct {
	svc service.DirectoryService
}

func NewDirectoryController(svc service.DirectoryService) *DirectoryController {
	return &DirectoryController{svc: svc}
}

// RegisterRoutes registra las rutas en el grupo /api/v1.
//
//	POST /api/v1/directories           → Create
//	POST /api/v1/directories/bulk      → CreateBulk (array JSON)
//	GET  /api/v1/directories           → List (?type=... requerido; ?cityName=... opcional; ?page=...&limit=...)
//	GET  /api/v1/directories/:id       → GetByID
func (c *DirectoryController) RegisterRoutes(rg *gin.RouterGroup) {
	dirs := rg.Group("/directories")
	dirs.POST("/bulk", c.CreateBulk)
	dirs.POST("", c.Create)
	dirs.GET("", c.List)
	dirs.GET("/:id", c.GetByID)
}

func (c *DirectoryController) Create(ctx *gin.Context) {
	var body service.CreateDirectoryInput
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "body inválido"})
		return
	}

	result, err := c.svc.Create(ctx.Request.Context(), body)
	if err != nil {
		writeDirectoryError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, result)
}

func (c *DirectoryController) CreateBulk(ctx *gin.Context) {
	var body []service.CreateDirectoryInput
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "body inválido: se espera un array JSON de directorios"})
		return
	}

	result, err := c.svc.CreateBulk(ctx.Request.Context(), body)
	if err != nil {
		writeDirectoryError(ctx, err)
		return
	}

	if result.Created == 0 {
		ctx.JSON(http.StatusBadRequest, result)
		return
	}
	ctx.JSON(http.StatusCreated, result)
}

func (c *DirectoryController) List(ctx *gin.Context) {
	dirType := ctx.Query("type")
	if dirType == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "parámetro 'type' requerido"})
		return
	}

	cityName := ctx.Query("cityName")
	page := ginQueryInt(ctx, "page", 0)
	limit := ginQueryInt(ctx, "limit", 20)

	result, err := c.svc.List(ctx.Request.Context(), cityName, dirType, page, limit)
	if err != nil {
		writeDirectoryError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *DirectoryController) GetByID(ctx *gin.Context) {
	result, err := c.svc.GetByID(ctx.Request.Context(), ctx.Param("id"))
	if err != nil {
		writeDirectoryError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func writeDirectoryError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, repository.ErrCityNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, service.ErrDirectoryNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{"error": "directorio no encontrado"})
	case errors.Is(err, service.ErrDirectoryInvalidType):
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "tipo de directorio inválido"})
	case errors.Is(err, service.ErrDirectoryInvalidInput):
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"})
	}
}
