// Package controller — psychosocial_calendar_controller.go
// Endpoints REST para el calendario del módulo Atención Psicosocial.
package controller

import (
	"bitsflow/common/utils"
	"bitsflow/internal/repository"
	"bitsflow/salvia/service"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type PsychosocialCalendarController struct {
	svc service.PsychosocialCalendarService
}

func NewPsychosocialCalendarController(svc service.PsychosocialCalendarService) *PsychosocialCalendarController {
	return &PsychosocialCalendarController{svc: svc}
}

func (c *PsychosocialCalendarController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/psychosocial/calendar/mine", c.Mine)
	rg.GET("/psychosocial/calendar/team", c.Team)
	rg.GET("/psychosocial/calendar/agents", c.Agents)
	rg.GET("/psychosocial/calendar/my-victims", c.MyVictims)
	rg.POST("/psychosocial/sessions", c.CreateSession)
	rg.PUT("/psychosocial/sessions/:id", c.UpdateSession)
	rg.DELETE("/psychosocial/sessions/:id", c.DeleteSession)
	rg.POST("/psychosocial/meetings", c.CreateMeeting)
	rg.GET("/psychosocial/meetings/:id", c.GetMeeting)
	rg.PUT("/psychosocial/meetings/:id", c.UpdateMeeting)
	rg.DELETE("/psychosocial/meetings/:id", c.DeleteMeeting)
}

// ─── Helpers ─────────────────────────────────────────────────────────────────
func (c *PsychosocialCalendarController) getSession(ctx *gin.Context) *utils.CommonSession {
	session := sessions.Default(ctx)
	sid, _ := session.Get("userData").(string)
	s, err := utils.GetCommonSession(sid)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "sesión inválida"})
		return nil
	}
	return s
}

func parseRange(ctx *gin.Context) (time.Time, time.Time) {
	fromStr := ctx.Query("from")
	toStr := ctx.Query("to")
	now := time.Now()
	from := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	to := from.AddDate(0, 1, 0)
	if t, err := time.Parse("2006-01-02", fromStr); err == nil {
		from = t
	}
	if t, err := time.Parse("2006-01-02", toStr); err == nil {
		to = t
	}
	return from, to
}

// ─── Handlers ────────────────────────────────────────────────────────────────
func (c *PsychosocialCalendarController) Mine(ctx *gin.Context) {
	s := c.getSession(ctx)
	if s == nil {
		return
	}
	from, to := parseRange(ctx)
	fmt.Printf("[Calendar.Mine] userICode=%q from=%s to=%s\n", s.UserICode, from.Format("2006-01-02"), to.Format("2006-01-02"))
	events, err := c.svc.MineEvents(ctx.Request.Context(), s.UserICode, from, to)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al cargar calendario"})
		return
	}
	if events == nil {
		events = []repository.CalendarEvent{}
	}
	fmt.Printf("[Calendar.Mine] returning %d events\n", len(events))
	ctx.JSON(http.StatusOK, gin.H{"events": events})
}

func (c *PsychosocialCalendarController) Team(ctx *gin.Context) {
	if s := c.getSession(ctx); s == nil {
		return
	}
	from, to := parseRange(ctx)
	filterAgent := strings.TrimSpace(ctx.Query("agentId"))
	fmt.Printf("[Calendar.Team] from=%s to=%s filterAgent=%q\n", from.Format("2006-01-02"), to.Format("2006-01-02"), filterAgent)
	events, err := c.svc.TeamEvents(ctx.Request.Context(), filterAgent, from, to)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al cargar calendario del equipo"})
		return
	}
	if events == nil {
		events = []repository.CalendarEvent{}
	}
	fmt.Printf("[Calendar.Team] returning %d events\n", len(events))
	ctx.JSON(http.StatusOK, gin.H{"events": events})
}

func (c *PsychosocialCalendarController) Agents(ctx *gin.Context) {
	if s := c.getSession(ctx); s == nil {
		return
	}
	agents, err := c.svc.ListAgents(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al cargar agentes"})
		return
	}
	if agents == nil {
		agents = []repository.CalendarAgentOption{}
	}
	ctx.JSON(http.StatusOK, gin.H{"agents": agents})
}

func (c *PsychosocialCalendarController) MyVictims(ctx *gin.Context) {
	s := c.getSession(ctx)
	if s == nil {
		return
	}
	agentID := s.UserICode
	if q := strings.TrimSpace(ctx.Query("agentId")); q != "" {
		agentID = q
	}
	victims, err := c.svc.MyVictims(ctx.Request.Context(), agentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al cargar víctimas"})
		return
	}
	if victims == nil {
		victims = []repository.CalendarVictimOption{}
	}
	ctx.JSON(http.StatusOK, gin.H{"victims": victims})
}

type createSessionReq struct {
	CaseICode     string  `json:"caseICode"`
	FollowUpID    string  `json:"followUpId"`
	RemissionID   string  `json:"remissionId"`
	AgentID       string  `json:"agentId"`
	ScheduledDate string  `json:"scheduledDate"` // YYYY-MM-DD
	ScheduledTime string  `json:"scheduledTime"` // HH:MM
	Summary       *string `json:"summary"`
}

func (c *PsychosocialCalendarController) CreateSession(ctx *gin.Context) {
	s := c.getSession(ctx)
	if s == nil {
		return
	}
	var req createSessionReq
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "body inválido"})
		return
	}
	date, err := time.Parse("2006-01-02", req.ScheduledDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "scheduledDate inválido (YYYY-MM-DD)"})
		return
	}
	agentID := strings.TrimSpace(req.AgentID)
	agentName := s.Names + " " + s.LastNames
	if agentID == "" {
		agentID = s.UserICode
	} else if agentID != s.UserICode {
		agentName = ""
	}
	fmt.Printf("[Calendar.CreateSession] agentID=%q caseICode=%q remissionID=%q date=%s time=%s\n", agentID, req.CaseICode, req.RemissionID, req.ScheduledDate, req.ScheduledTime)
	id, err := c.svc.CreateSession(ctx.Request.Context(), repository.CreateSessionInput{
		AgentID:       agentID,
		AgentName:     agentName,
		CaseICode:     req.CaseICode,
		FollowUpID:    req.FollowUpID,
		RemissionID:   req.RemissionID,
		ScheduledDate: date,
		ScheduledTime: req.ScheduledTime,
		Summary:       req.Summary,
	})
	if err != nil {
		fmt.Printf("[Calendar.CreateSession] ERROR: %v\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("[Calendar.CreateSession] created id=%s\n", id)
	ctx.JSON(http.StatusCreated, gin.H{"id": id})
}

func (c *PsychosocialCalendarController) UpdateSession(ctx *gin.Context) {
	s := c.getSession(ctx)
	if s == nil {
		return
	}
	var req createSessionReq
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "body inválido"})
		return
	}
	date, err := time.Parse("2006-01-02", req.ScheduledDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "scheduledDate inválido (YYYY-MM-DD)"})
		return
	}
	agentID := strings.TrimSpace(req.AgentID)
	if agentID == "" {
		agentID = s.UserICode
	}
	err = c.svc.UpdateSession(ctx.Request.Context(), ctx.Param("id"), repository.CreateSessionInput{
		AgentID:       agentID,
		ScheduledDate: date,
		ScheduledTime: req.ScheduledTime,
		Summary:       req.Summary,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"ok": true})
}

func (c *PsychosocialCalendarController) DeleteSession(ctx *gin.Context) {
	if s := c.getSession(ctx); s == nil {
		return
	}
	if err := c.svc.DeleteSession(ctx.Request.Context(), ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"ok": true})
}

type createMeetingReq struct {
	Title         string   `json:"title"`
	Description   *string  `json:"description"`
	ScheduledDate string   `json:"scheduledDate"`
	ScheduledTime string   `json:"scheduledTime"`
	DurationMin   int      `json:"durationMin"`
	AgentIDs      []string `json:"agentIds"`
}

func (c *PsychosocialCalendarController) CreateMeeting(ctx *gin.Context) {
	s := c.getSession(ctx)
	if s == nil {
		return
	}
	var req createMeetingReq
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "body inválido"})
		return
	}
	date, err := time.Parse("2006-01-02", req.ScheduledDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "scheduledDate inválido (YYYY-MM-DD)"})
		return
	}
	id, err := c.svc.CreateMeeting(ctx.Request.Context(), repository.CreateMeetingInput{
		Title:         req.Title,
		Description:   req.Description,
		ScheduledDate: date,
		ScheduledTime: req.ScheduledTime,
		DurationMin:   req.DurationMin,
		CreatedBy:     s.UserICode,
		CreatedByName: s.Names + " " + s.LastNames,
		AgentIDs:      req.AgentIDs,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"id": id})
}

func (c *PsychosocialCalendarController) GetMeeting(ctx *gin.Context) {
	if s := c.getSession(ctx); s == nil {
		return
	}
	m, err := c.svc.GetMeeting(ctx.Request.Context(), ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "reunión no encontrada"})
		return
	}
	ctx.JSON(http.StatusOK, m)
}

func (c *PsychosocialCalendarController) UpdateMeeting(ctx *gin.Context) {
	if s := c.getSession(ctx); s == nil {
		return
	}
	var req createMeetingReq
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "body inválido"})
		return
	}
	date, err := time.Parse("2006-01-02", req.ScheduledDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "scheduledDate inválido (YYYY-MM-DD)"})
		return
	}
	err = c.svc.UpdateMeeting(ctx.Request.Context(), ctx.Param("id"), repository.CreateMeetingInput{
		Title:         req.Title,
		Description:   req.Description,
		ScheduledDate: date,
		ScheduledTime: req.ScheduledTime,
		DurationMin:   req.DurationMin,
		AgentIDs:      req.AgentIDs,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"ok": true})
}

func (c *PsychosocialCalendarController) DeleteMeeting(ctx *gin.Context) {
	if s := c.getSession(ctx); s == nil {
		return
	}
	if err := c.svc.DeleteMeeting(ctx.Request.Context(), ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"ok": true})
}
