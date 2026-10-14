package controller

import (
	"bitsflow/internal/repository"
	salvia_config "bitsflow/salvia/config"
	"bitsflow/salvia/service"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// EntityObligationResponseController expone la pantalla "Mis Acciones
// Institucionales" del rol entidad y el registro de su gestión.
type EntityObligationResponseController struct {
	svc service.EntityObligationResponseService
}

func NewEntityObligationResponseController(svc service.EntityObligationResponseService) *EntityObligationResponseController {
	return &EntityObligationResponseController{svc: svc}
}

// RegisterRoutes registra las rutas en el grupo /api/v1.
//
//	GET  /api/v1/mis-acciones-institucionales?status=&document=&cityId=&page=&pageSize=
//	POST /api/v1/mis-acciones-institucionales/respuesta
func (c *EntityObligationResponseController) RegisterRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/mis-acciones-institucionales")
	g.GET("", c.List)
	g.POST("/respuesta", c.SaveResponse)
}

// List devuelve las acciones institucionales de la sede del usuario, separadas por
// pestaña, con los contadores de ambas y las opciones de los filtros.
//
// La sede sale de la sesión y nunca del querystring: es lo que garantiza que la
// entidad solo vea sus propias acciones, como exige el DoD.
func (c *EntityObligationResponseController) List(ctx *gin.Context) {
	s := requireSessionRole(ctx, salvia_config.PermissionsByRole["get_mis_acciones_institucionales"])
	if s == nil {
		return
	}
	if s.EntityBranchId <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "el usuario no tiene sede (entity_branch_id) asignada"})
		return
	}

	result, err := c.svc.ListForEntityBranch(ctx.Request.Context(), repository.EntityObligationListFilter{
		EntityBranchID: s.EntityBranchId,
		Status:         ctx.Query("status"),
		Document:       ctx.Query("document"),
		CityID:         ctx.Query("cityId"),
		Page:           ginQueryInt(ctx, "page", 0),
		PageSize:       ginQueryInt(ctx, "pageSize", 10),
	})
	if err != nil {
		// El detalle va al log y no a la respuesta: al usuario no se le muestran
		// errores de base de datos, pero sin esta línea un 500 acá no deja rastro
		// de qué falló y no hay forma de diagnosticarlo desde QA.
		log.Printf("[MisAccionesInstitucionales] List falló (branch=%d status=%q document=%q cityId=%q): %v",
			s.EntityBranchId, ctx.Query("status"), ctx.Query("document"), ctx.Query("cityId"), err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

// SaveResponse guarda lo diligenciado en el modal "Completar acción institucional".
// Si la entidad vuelve a guardar la misma acción, corrige su respuesta anterior.
//
//	POST /api/v1/mis-acciones-institucionales/respuesta
//	Body: { "entityCaseId": "...", "entityObligationId": "...",
//	        "radicado": "...", "ampliacion": "...", "estado": "en_tramite" }
func (c *EntityObligationResponseController) SaveResponse(ctx *gin.Context) {
	s := requireSessionRole(ctx, salvia_config.PermissionsByRole["get_mis_acciones_institucionales"])
	if s == nil {
		return
	}
	if s.EntityBranchId <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "el usuario no tiene sede (entity_branch_id) asignada"})
		return
	}

	var body struct {
		EntityCaseID       string `json:"entityCaseId"       binding:"required"`
		EntityObligationID string `json:"entityObligationId" binding:"required"`
		Radicado           string `json:"radicado"`
		Ampliacion         string `json:"ampliacion"`
		Estado             string `json:"estado"             binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.svc.SaveResponse(ctx.Request.Context(), service.SaveEntityObligationResponseInput{
		EntityBranchID:     s.EntityBranchId,
		EntityCaseID:       body.EntityCaseID,
		EntityObligationID: body.EntityObligationID,
		Radicado:           body.Radicado,
		Ampliacion:         body.Ampliacion,
		Estado:             body.Estado,
		RespondedByID:      s.UserICode,
	})
	switch {
	case err == nil:
		ctx.JSON(http.StatusOK, gin.H{"ok": true})
	case errors.Is(err, service.ErrEntityObligationEstadoInvalido):
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "el estado de la respuesta no es válido"})
	case errors.Is(err, service.ErrEntityObligationNoAsignada):
		// 404 y no 403: si la acción es de otra entidad, no se confirma que exista.
		ctx.JSON(http.StatusNotFound, gin.H{"error": "la acción institucional no está asignada a esta entidad"})
	default:
		log.Printf("[MisAccionesInstitucionales] SaveResponse falló (branch=%d caso=%s accion=%s): %v",
			s.EntityBranchId, body.EntityCaseID, body.EntityObligationID, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
	}
}
