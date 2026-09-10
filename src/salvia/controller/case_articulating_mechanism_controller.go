package controller

import (
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	"bitsflow/salvia/service"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CaseArticulatingMechanismController struct {
	svc service.CaseArticulatingMechanismService
}

func NewCaseArticulatingMechanismController(svc service.CaseArticulatingMechanismService) *CaseArticulatingMechanismController {
	return &CaseArticulatingMechanismController{svc: svc}
}

var mechanismWriteRoles = map[string]bool{"sv": true, "op": true, "ro": true}

func (c *CaseArticulatingMechanismController) RegisterRoutes(rg *gin.RouterGroup) {
	casos := rg.Group("/casos/:id/mecanismos")
	casos.GET("", c.List)
	casos.POST("", c.Create)
	casos.POST("/:mid/comentarios", c.AddComment)

	rg.GET("/mecanismos/mine", c.ListMine)
	rg.GET("/mecanismos/usuarios", c.ListUsers)
}

func (c *CaseArticulatingMechanismController) List(ctx *gin.Context) {
	if requireCaseDetailAccess(ctx) == nil {
		return
	}
	items, err := c.svc.ListByCase(ctx.Request.Context(), ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}
	ctx.JSON(http.StatusOK, items)
}

func requireMcAccess(ctx *gin.Context) *utils.CommonSession {
	sessionID, ok := sessions.Default(ctx).Get("userData").(string)
	if !ok || sessionID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no autenticado"})
		return nil
	}
	s, err := utils.GetCommonSession(sessionID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no autenticado"})
		return nil
	}
	if !utils.CheckPermission(salvia_config.PermissionsByRole, "get_casos_mecanismo", s.CurrentRole, ctx) {
		return nil
	}
	return s
}

// ListMine — GET /api/v1/mecanismos/mine?document=&city=&page=&pageSize=
func (c *CaseArticulatingMechanismController) ListMine(ctx *gin.Context) {
	s := requireMcAccess(ctx)
	if s == nil {
		return
	}
	page := ginQueryInt(ctx, "page", 0)
	pageSize := ginQueryInt(ctx, "pageSize", 5)
	if pageSize <= 0 {
		pageSize = 5
	}
	result, err := c.svc.ListMine(ctx.Request.Context(), service.ListMineFilter{
		AssignedUserID: s.UserICode,
		Document:       ctx.Query("document"),
		City:           ctx.Query("city"),
		Page:           page,
		PageSize:       pageSize,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *CaseArticulatingMechanismController) Create(ctx *gin.Context) {
	s := requireCaseDetailAccess(ctx)
	if s == nil {
		return
	}
	if !mechanismWriteRoles[s.CurrentRole] {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "su rol no tiene permiso para agregar mecanismos"})
		return
	}

	var body struct {
		Scope            string `json:"scope"`
		AssignedUserID   string `json:"assignedUserId"`
		DepartmentID     string `json:"departmentId"`
		CityID           string `json:"cityId"`
		TownCode         string `json:"townCode"`
		SalviaCommitment string `json:"salviaCommitment"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	row, err := c.svc.Create(ctx.Request.Context(), service.CreateArticulatingMechanismInput{
		CaseID:           ctx.Param("id"),
		Scope:            body.Scope,
		AssignedUserID:   body.AssignedUserID,
		DepartmentID:     body.DepartmentID,
		CityID:           body.CityID,
		TownCode:         body.TownCode,
		SalviaCommitment: body.SalviaCommitment,
		CreatedByID:      s.UserICode,
		ActorName:        strings.TrimSpace(s.Names + " " + s.LastNames),
	})
	if err != nil {
		writeMechanismError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, row)
}

func (c *CaseArticulatingMechanismController) AddComment(ctx *gin.Context) {
	s := requireCaseDetailAccess(ctx)
	if s == nil {
		return
	}
	if !mechanismWriteRoles[s.CurrentRole] {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "su rol no tiene permiso para comentar mecanismos"})
		return
	}

	var body struct {
		Body string `json:"body"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment, err := c.svc.AddComment(ctx.Request.Context(), service.AddMechanismCommentInput{
		CaseID:      ctx.Param("id"),
		MechanismID: ctx.Param("mid"),
		AuthorID:    s.UserICode,
		AuthorName:  strings.TrimSpace(s.Names + " " + s.LastNames),
		Body:        body.Body,
	})
	if err != nil {
		writeMechanismError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, comment)
}

func (c *CaseArticulatingMechanismController) ListUsers(ctx *gin.Context) {
	if requireCaseDetailAccess(ctx) == nil {
		return
	}
	items, err := c.svc.ListUsers(ctx.Request.Context(), service.ListMechanismUsersInput{
		Scope:        ctx.Query("scope"),
		DepartmentID: ctx.Query("department_id"),
		CityID:       ctx.Query("city_id"),
		TownCode:     firstNonEmpty(ctx.Query("town_id"), ctx.Query("town_code")),
	})
	if err != nil {
		writeMechanismError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, items)
}

func writeMechanismError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrMechanismInvalidScope),
		errors.Is(err, service.ErrMechanismMissingLocation),
		errors.Is(err, service.ErrMechanismUserRequired),
		errors.Is(err, service.ErrMechanismCommentEmpty),
		errors.Is(err, service.ErrMechanismCommitmentTooLong),
		errors.Is(err, service.ErrGestionMecanismoInvalid):
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, service.ErrMechanismUserNotFound),
		errors.Is(err, service.ErrMechanismNotFound),
		errors.Is(err, gorm.ErrRecordNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, service.ErrMechanismUserOutsideGeo),
		errors.Is(err, service.ErrMechanismNotAssigned),
		errors.Is(err, service.ErrMechanismNotActive):
		ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
	}
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}
