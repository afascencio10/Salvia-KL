package controller

import (
	"bitsflow/internal/repository"
	"bitsflow/salvia/service"
	"errors"
	"net/http"
	"strconv"

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
//	POST   /api/v1/directories               → Create
//	POST   /api/v1/directories/bulk          → CreateBulk (array JSON)
//	GET    /api/v1/directories               → List (filtros opcionales)
//	GET    /api/v1/directories/:id           → GetByID
//	PUT    /api/v1/directories/:id           → Update
//	PATCH  /api/v1/directories/:id/disable   → Disable (soft)
//	PATCH  /api/v1/directories/:id/enable    → Enable
func (c *DirectoryController) RegisterRoutes(rg *gin.RouterGroup) {
	dirs := rg.Group("/directories")
	dirs.POST("/bulk", c.CreateBulk)
	dirs.POST("", c.Create)
	dirs.GET("", c.List)
	dirs.GET("/:id", c.GetByID)
	dirs.PUT("/:id", c.Update)
	dirs.PATCH("/:id/disable", c.Disable)
	dirs.PATCH("/:id/enable", c.Enable)
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
	filters := service.DirectoryListFilters{
		CityName:  ctx.Query("cityName"),
		CityICode: ctx.Query("cityICode"),
		Type:      ctx.Query("type"),
		Sector:    ctx.Query("sector"),
		NameQuery: ctx.Query("name"),
	}
	if v := ctx.Query("departmentId"); v != "" {
		id, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "departmentId debe ser un número entero"})
			return
		}
		filters.DepartmentID = &id
	}
	if ctx.Query("includeInactive") == "true" {
		filters.IncludeInactive = true
	}

	page := ginQueryInt(ctx, "page", 0)
	limit := ginQueryInt(ctx, "limit", 200)

	result, err := c.svc.List(ctx.Request.Context(), filters, page, limit)
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

func (c *DirectoryController) Update(ctx *gin.Context) {
	var body service.UpdateDirectoryInput
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "body inválido"})
		return
	}

	result, err := c.svc.Update(ctx.Request.Context(), ctx.Param("id"), body)
	if err != nil {
		writeDirectoryError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *DirectoryController) Disable(ctx *gin.Context) {
	if err := c.svc.Disable(ctx.Request.Context(), ctx.Param("id")); err != nil {
		writeDirectoryError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"ok": true})
}

func (c *DirectoryController) Enable(ctx *gin.Context) {
	if err := c.svc.Enable(ctx.Request.Context(), ctx.Param("id")); err != nil {
		writeDirectoryError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"ok": true})
}

func writeDirectoryError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, repository.ErrCityNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, service.ErrDirectoryNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{"error": "directorio no encontrado"})
	case errors.Is(err, service.ErrDirectoryInvalidType):
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "tipo de directorio inválido"})
	case errors.Is(err, service.ErrDirectorySectorType):
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "sector inválido"})
	case errors.Is(err, service.ErrDirectoryInvalidInput):
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"})
	}
}
