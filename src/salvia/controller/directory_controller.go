package controller

import (
	"bitsflow/internal/repository"
	salvia_config "bitsflow/salvia/config"
	"bitsflow/salvia/service"
	"errors"
	"net/http"
	"strconv"
	"strings"

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
//	POST   /api/v1/directories/bulk-upload   → BulkUpload (CSV del "Mapeo de actores", multipart campo "file")
//	GET    /api/v1/directories               → List (filtros opcionales)
//	GET    /api/v1/directories/:id           → GetByID
//	PUT    /api/v1/directories/:id           → Update
//	PATCH  /api/v1/directories/:id/disable   → Disable (soft)
//	PATCH  /api/v1/directories/:id/enable    → Enable
//
// Las lecturas quedan abiertas: las consume la app móvil sin sesión de la web.
// Toda escritura pasa por requireDirectoryPermission con el permiso que ya define
// Menu.go para esa acción — son tres distintos, no uno: crear lo puede hacer
// entidad y supervisor, actualizar y desactivar solo supervisor.
func (c *DirectoryController) RegisterRoutes(rg *gin.RouterGroup) {
	dirs := rg.Group("/directories")
	dirs.GET("", c.List)
	dirs.GET("/:id", c.GetByID)

	// El guard va como middleware del grupo y no repetido dentro de cada handler
	// para que una ruta de escritura nueva no pueda quedar sin protección por
	// olvido — que es exactamente cómo quedaron abiertas estas seis.
	crear := dirs.Group("", requireDirectoryPermission("set_directorio_entidad"))
	crear.POST("", c.Create)
	crear.POST("/bulk", c.CreateBulk)

	actualizar := dirs.Group("", requireDirectoryPermission("update_directorio_entidad"))
	actualizar.PUT("/:id", c.Update)
	// La carga masiva actualiza lo existente además de crear, así que exige el
	// permiso de actualizar, el más estricto de los dos.
	actualizar.POST("/bulk-upload", c.BulkUpload)

	desactivar := dirs.Group("", requireDirectoryPermission("disable_directorio_entidad"))
	desactivar.PATCH("/:id/disable", c.Disable)
	desactivar.PATCH("/:id/enable", c.Enable)
}

// requireDirectoryPermission corta la petición con 401 o 403 si la sesión no tiene
// el permiso indicado en salvia_config.PermissionsByRole.
func requireDirectoryPermission(permission string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// requireSessionRole ya escribió la respuesta; falta cortar la cadena
		// porque usa ctx.JSON y no ctx.AbortWithStatusJSON.
		if requireSessionRole(ctx, salvia_config.PermissionsByRole[permission]) == nil {
			ctx.Abort()
			return
		}
		ctx.Next()
	}
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

// BulkUpload recibe rutas_atencion_territorial.csv (paquete "BD consolidada mapeo de
// actores" de SALVIA): actualiza las entidades existentes y crea las nuevas. El
// permiso lo exige el middleware del grupo en RegisterRoutes.
func (c *DirectoryController) BulkUpload(ctx *gin.Context) {
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "adjunta el archivo en el campo 'file' (multipart/form-data)"})
		return
	}
	if !strings.HasSuffix(strings.ToLower(fileHeader.Filename), ".csv") {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "el archivo debe ser .csv"})
		return
	}
	file, err := fileHeader.Open()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "no se pudo abrir el archivo"})
		return
	}
	defer file.Close()

	result, err := c.svc.UpsertFromCSV(ctx.Request.Context(), file)
	if err != nil {
		writeDirectoryError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, result)
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
