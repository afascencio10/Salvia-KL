package controller

import (
	"bitsflow/salvia/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CaseProfessionalController expone el listado de "Profesionales de atención"
// de un caso (pestaña Derivaciones del detalle del caso).
type CaseProfessionalController struct {
	svc service.CaseProfessionalService
}

func NewCaseProfessionalController(svc service.CaseProfessionalService) *CaseProfessionalController {
	return &CaseProfessionalController{svc: svc}
}

// RegisterRoutes registra la ruta en el grupo /api/v1.
//
// Usa ":id" (no ":caseId") para el segmento de caso: Gin exige el mismo nombre
// de wildcard en todas las rutas que comparten la posición del árbol, y
// CaseDetailController ya registró "/casos/:id/..." en este mismo grupo.
//
//	GET /api/v1/casos/:id/profesionales-atencion → List
func (c *CaseProfessionalController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/casos/:id/profesionales-atencion", c.List)
}

// List devuelve los profesionales/personas que han intervenido en el caso,
// una fila por persona, con sus roles concatenados y sus intervenciones.
//
//	GET /api/v1/casos/:id/profesionales-atencion
func (c *CaseProfessionalController) List(ctx *gin.Context) {
	if requireCaseDetailAccess(ctx) == nil {
		return
	}
	caseID := ctx.Param("id")
	items, err := c.svc.ListByCase(ctx.Request.Context(), caseID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}
	ctx.JSON(http.StatusOK, items)
}
