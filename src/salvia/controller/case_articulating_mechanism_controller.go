package controller

import (
	"bitsflow/salvia/service"
	"errors"
	"net/http"
	"strings"

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
		errors.Is(err, service.ErrMechanismCommitmentTooLong):
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, service.ErrMechanismUserNotFound),
		errors.Is(err, service.ErrMechanismNotFound),
		errors.Is(err, gorm.ErrRecordNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, service.ErrMechanismUserOutsideGeo):
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
