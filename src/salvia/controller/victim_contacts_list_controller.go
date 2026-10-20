// Package controller — victim_contacts_list_controller.go
// Endpoint JSON para casos-component mostrando reportes (Listado Reportes /
// Mis Reportes). Ver DocsMD/Screens/Listado Reportes/listado-reportes-events.md.
package controller

import (
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	"bitsflow/salvia/service"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// VictimContactsListController maneja GET /api/v1/victim-contacts/list.
type VictimContactsListController struct {
	svc service.VictimContactsListService
}

func NewVictimContactsListController(svc service.VictimContactsListService) *VictimContactsListController {
	return &VictimContactsListController{svc: svc}
}

// RegisterRoutes registra la ruta del listado de reportes.
//
//	GET /api/v1/victim-contacts/list?filter_key=&filter_value=&filter_estado_reporte=
//	    &search=&sort=&order=&page=&page_size=
func (c *VictimContactsListController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/victim-contacts/list", c.List)
}

// List devuelve reportes paginados con filtros para casos-component.
// Solo op/sv (permiso list_victim_contacts, Menu.go) — FigJam: "ve TODOS los
// reportes; visible para op y sv".
func (c *VictimContactsListController) List(ctx *gin.Context) {
	sessionID, ok := sessions.Default(ctx).Get("userData").(string)
	if !ok || sessionID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no autenticado"})
		return
	}
	sess, err := utils.GetCommonSession(sessionID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "sesión expirada"})
		return
	}
	if !utils.CheckPermission(salvia_config.PermissionsByRole, "list_victim_contacts", sess.CurrentRole, ctx) {
		return // CheckPermission ya escribió el 401
	}

	page := victimContactsQueryInt(ctx, "page", 1)
	pageSize := victimContactsQueryInt(ctx, "page_size", 20)

	result, err := c.svc.List(ctx.Request.Context(), service.VictimContactsListInput{
		FilterKey:           ctx.Query("filter_key"),
		FilterValue:         ctx.Query("filter_value"),
		FilterEstadoReporte: ctx.Query("filter_estado_reporte"),
		Search:              ctx.Query("search"),
		Sort:                ctx.DefaultQuery("sort", "registration_date"),
		Order:               ctx.DefaultQuery("order", "desc"),
		Page:                page,
		PageSize:            pageSize,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al cargar los reportes"})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func victimContactsQueryInt(ctx *gin.Context, key string, defaultVal int) int {
	raw := ctx.Query(key)
	if raw == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(raw)
	if err != nil {
		return defaultVal
	}
	return val
}
