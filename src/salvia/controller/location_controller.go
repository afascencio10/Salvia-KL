package controller

import (
	"bitsflow/internal/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// LocationController expone los endpoints de ubicaciones geográficas.
type LocationController struct {
	repo repository.LocationRepository
}

func NewLocationController(repo repository.LocationRepository) *LocationController {
	return &LocationController{repo: repo}
}

// RegisterRoutes registra las rutas en el grupo /api/v1.
//
//	GET /api/v1/locations/departments         → lista de departamentos { label, value }
//	GET /api/v1/locations/cities              → lista de ciudades { label, value, departmentId }
//	GET /api/v1/locations/towns?city_id=<id>  → municipios de una ciudad { label, value }
func (c *LocationController) RegisterRoutes(rg *gin.RouterGroup) {
	loc := rg.Group("/locations")
	loc.GET("/departments", c.GetDepartments)
	loc.GET("/cities", c.GetCities)
	loc.GET("/towns", c.GetTownsByCityID)
}

func (c *LocationController) GetDepartments(ctx *gin.Context) {
	departments, err := c.repo.GetDepartments(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al cargar departamentos"})
		return
	}
	ctx.JSON(http.StatusOK, departments)
}

func (c *LocationController) GetCities(ctx *gin.Context) {
	cities, err := c.repo.GetCities(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al cargar ciudades"})
		return
	}
	ctx.JSON(http.StatusOK, cities)
}

func (c *LocationController) GetTownsByCityID(ctx *gin.Context) {
	cityIDStr := ctx.Query("city_id")
	if cityIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "parámetro 'city_id' requerido"})
		return
	}

	cityID, err := strconv.ParseUint(cityIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "city_id debe ser un número entero"})
		return
	}

	towns, err := c.repo.GetTownsByCityID(ctx.Request.Context(), cityID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al cargar municipios"})
		return
	}
	ctx.JSON(http.StatusOK, towns)
}
