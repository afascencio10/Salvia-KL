// Package controller — case_detail_controller.go
// Controlador HTTP para la pantalla de detalle de caso (rol sv).
//
// ARQUITECTURA:
//
//	main.go → RegisterRoutes → CaseDetailGET (fachada HTML)
//	main.go → GET /api/v1/casos/:id/detalle → GetByID (API JSON)
//
// La ruta HTML vive en MainRouter.go (sistema de fachadas existente).
// La ruta API JSON se registra aquí con RegisterRoutes.
package controller

import (
	internaldb "bitsflow/internal/db"
	"bitsflow/internal/repository"
	"bitsflow/salvia/service"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CaseDetailController struct {
	svc          service.CaseDetailService
	timelineRepo repository.CaseTimelineEventRepository
}

func NewCaseDetailController(svc service.CaseDetailService, timelineRepo repository.CaseTimelineEventRepository) *CaseDetailController {
	return &CaseDetailController{svc: svc, timelineRepo: timelineRepo}
}

// RegisterRoutes registra el endpoint JSON en /api/v1.
//
//	GET /api/v1/casos/:id/detalle
//	GET /api/v1/operadores
func (c *CaseDetailController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/casos/:id/detalle", c.GetByID)
	rg.GET("/operadores", c.GetOperadores)
	rg.GET("/agentes-ro", c.GetAgentesRO)
	rg.GET("/equipo-operadores", c.GetOperadoresByTeam)
	rg.POST("/casos/:id/seguimiento", c.CrearSeguimiento)
	rg.POST("/casos/:id/timeline", c.AddTimelineEvent)
	rg.GET("/casos/:id/timeline-events", c.GetTimelineEvents)
	rg.POST("/casos/:id/reasignar", c.ReasignarCaso)
	rg.PUT("/seguimiento/:segId/reasignar", c.ReassignFollowUp)
}

// GetByID responde con el detalle completo del caso en JSON.
func (c *CaseDetailController) GetByID(ctx *gin.Context) {
	icode := ctx.Param("id")
	detail, err := c.svc.GetDetail(ctx.Request.Context(), icode)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidICode):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "icode inválido"})
		case errors.Is(err, service.ErrCaseNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": "caso no encontrado"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"})
		}
		return
	}
	ctx.JSON(http.StatusOK, detail)
}

// GetOperadores devuelve la lista de operadores (rol "op") para el dropdown de asignación.
func (c *CaseDetailController) GetOperadores(ctx *gin.Context) {
	team := ctx.Query("team")
	type OperadorResult struct {
		ICode    string `json:"icode" gorm:"column:general_user_i_code"`
		FullName string `json:"fullName" gorm:"column:full_name"`
		Team     string `json:"team" gorm:"column:general_user_team"`
	}
	var results []OperadorResult
	var err error
	err = internaldb.WithRetry(func() error {
		results = nil
		if team != "" {
			return c.svc.GetDB().Raw(`
				SELECT DISTINCT gu.general_user_i_code,
					   gup.general_user_profile_names || ' ' || gup.general_user_profile_last_names as full_name,
					   gu.general_user_team
				FROM security.general_user gu
				JOIN security.general_user_profile gup ON gup.general_user_profile_id = gu.general_user_general_user_profile
				JOIN security.rel_role_general_user rr ON rr.general_user_id = gu.general_user_id
				JOIN security.role r ON r.role_id = rr.role_id
				WHERE r.role_code = 'ro' AND gu.general_user_status = 'e' AND gu.general_user_team = ?
				ORDER BY full_name ASC`, team).Scan(&results).Error
		}
		return c.svc.GetDB().Raw(`
			SELECT DISTINCT gu.general_user_i_code,
				   gup.general_user_profile_names || ' ' || gup.general_user_profile_last_names as full_name,
				   gu.general_user_team
			FROM security.general_user gu
			JOIN security.general_user_profile gup ON gup.general_user_profile_id = gu.general_user_general_user_profile
			JOIN security.rel_role_general_user rr ON rr.general_user_id = gu.general_user_id
			JOIN security.role r ON r.role_id = rr.role_id
			WHERE r.role_code = 'ro' AND gu.general_user_status = 'e'
			ORDER BY full_name ASC`).Scan(&results).Error
	})
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if results == nil {
		results = []OperadorResult{}
	}
	ctx.JSON(200, results)
}

// GetAgentesRO devuelve usuarios con rol "ro" (para nuevo seguimiento y reasignar seguimiento).
func (c *CaseDetailController) GetAgentesRO(ctx *gin.Context) {
	team := ctx.Query("team")
	type OperadorResult struct {
		ICode    string `json:"icode" gorm:"column:general_user_i_code"`
		FullName string `json:"fullName" gorm:"column:full_name"`
		Team     string `json:"team" gorm:"column:general_user_team"`
	}
	var results []OperadorResult
	var err error
	err = internaldb.WithRetry(func() error {
		results = nil
		if team != "" {
			return c.svc.GetDB().Raw(`
				SELECT DISTINCT gu.general_user_i_code,
					   gup.general_user_profile_names || ' ' || gup.general_user_profile_last_names as full_name,
					   gu.general_user_team
				FROM security.general_user gu
				JOIN security.general_user_profile gup ON gup.general_user_profile_id = gu.general_user_general_user_profile
				JOIN security.rel_role_general_user rr ON rr.general_user_id = gu.general_user_id
				JOIN security.role r ON r.role_id = rr.role_id
				WHERE r.role_code = 'ro' AND gu.general_user_status = 'e' AND gu.general_user_team = ?
				ORDER BY full_name ASC`, team).Scan(&results).Error
		}
		return c.svc.GetDB().Raw(`
			SELECT DISTINCT gu.general_user_i_code,
				   gup.general_user_profile_names || ' ' || gup.general_user_profile_last_names as full_name,
				   gu.general_user_team
			FROM security.general_user gu
			JOIN security.general_user_profile gup ON gup.general_user_profile_id = gu.general_user_general_user_profile
			JOIN security.rel_role_general_user rr ON rr.general_user_id = gu.general_user_id
			JOIN security.role r ON r.role_id = rr.role_id
			WHERE r.role_code = 'ro' AND gu.general_user_status = 'e'
			ORDER BY full_name ASC`).Scan(&results).Error
	})
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if results == nil {
		results = []OperadorResult{}
	}
	ctx.JSON(200, results)
}

// GetOperadoresByTeam devuelve operadores (op + ro) filtrados por team.
func (c *CaseDetailController) GetOperadoresByTeam(ctx *gin.Context) {
	team := ctx.Query("team")

	type OperadorResult struct {
		ICode     string `json:"icode" gorm:"column:general_user_i_code"`
		FullName  string `json:"fullName" gorm:"column:full_name"`
		Team      string `json:"team" gorm:"column:general_user_team"`
		Names     string `json:"names" gorm:"column:general_user_profile_names"`
		LastNames string `json:"lastNames" gorm:"column:general_user_profile_last_names"`
	}

	var results []OperadorResult
	query := `
		SELECT DISTINCT gu.general_user_i_code,
			   gup.general_user_profile_names || ' ' || gup.general_user_profile_last_names as full_name,
			   gu.general_user_team,
			   gup.general_user_profile_names,
			   gup.general_user_profile_last_names
		FROM security.general_user gu
		JOIN security.general_user_profile gup ON gup.general_user_profile_id = gu.general_user_general_user_profile
		JOIN security.rel_role_general_user rr ON rr.general_user_id = gu.general_user_id
		JOIN security.role r ON r.role_id = rr.role_id
		WHERE r.role_code = 'ro'
		  AND gu.general_user_status = 'e'
	`
	args := []interface{}{}
	if team != "" {
		query += " AND gu.general_user_team = ?"
		args = append(args, team)
	}
	query += " ORDER BY full_name ASC"

	internaldb.WithRetry(func() error {
		results = nil
		return c.svc.GetDB().Raw(query, args...).Scan(&results).Error
	})

	if results == nil {
		results = []OperadorResult{}
	}

	type FrontendUser struct {
		ICode    string `json:"icode"`
		FullName string `json:"fullName"`
		Team     string `json:"team"`
	}
	var output []FrontendUser
	for _, r := range results {
		output = append(output, FrontendUser{
			ICode:    r.ICode,
			FullName: r.FullName,
			Team:     r.Team,
		})
	}
	if output == nil {
		output = []FrontendUser{}
	}
	ctx.JSON(200, output)
}

// CrearSeguimiento crea un registro básico en follow_up_v2.
// POST /api/v1/casos/:id/seguimiento
func (c *CaseDetailController) CrearSeguimiento(ctx *gin.Context) {
	caseICode := ctx.Param("id")

	var body struct {
		AgentID       string `json:"agent_id"       binding:"required"`
		ScheduledDate string `json:"scheduled_date"  binding:"required"`
		Notas         string `json:"notas"`
		CreatedBy     string `json:"created_by"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	followUp, err := c.svc.CreateFollowUp(ctx.Request.Context(), caseICode, body.AgentID, body.ScheduledDate, body.Notas, body.CreatedBy)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, followUp)
}

// AddTimelineEvent registra un evento en el timeline del caso.
// POST /api/v1/casos/:id/timeline
func (c *CaseDetailController) AddTimelineEvent(ctx *gin.Context) {
	caseICode := ctx.Param("id")
	var body struct {
		EventType   string `json:"event_type"   binding:"required"`
		Description string `json:"description"  binding:"required"`
		ActorName   string `json:"actor_name"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.svc.AddTimelineEvent(ctx.Request.Context(), caseICode, body.EventType, body.Description, "", body.ActorName); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"ok": true})
}

// GetTimelineEvents retorna los eventos del timeline de un caso con el nombre del actor resuelto.
// GET /api/v1/casos/:id/timeline-events?barrierId=xxx
func (c *CaseDetailController) GetTimelineEvents(ctx *gin.Context) {
	caseID    := ctx.Param("id")
	barrierID := ctx.Query("barrierId")

	events, err := c.timelineRepo.GetByCaseID(ctx.Request.Context(), caseID, barrierID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al cargar eventos del timeline"})
		return
	}
	ctx.JSON(http.StatusOK, events)
}

// ReassignFollowUp actualiza el agent_id de un seguimiento.
// PUT /api/v1/seguimiento/:segId/reasignar
func (c *CaseDetailController) ReassignFollowUp(ctx *gin.Context) {
	segID := ctx.Param("segId")
	var body struct {
		AgentID string `json:"agent_id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.svc.ReassignFollowUp(ctx.Request.Context(), segID, body.AgentID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"ok": true})
}

// ReasignarCaso reasigna un caso a un nuevo operador usando GORM directamente.
// POST /api/v1/casos/:id/reasignar
func (c *CaseDetailController) ReasignarCaso(ctx *gin.Context) {
	caseICode := ctx.Param("id")
	var body struct {
		OperadorICode string `json:"operador_icode" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.svc.ReasignarCaso(ctx.Request.Context(), caseICode, body.OperadorICode); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"ok": true})
}
