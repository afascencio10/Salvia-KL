// Package controller — endpoints HTTP (Gin) del flujo 3x3 de Atención Psicosocial.
package controller

import (
	"bitsflow/common/utils"
	"bitsflow/salvia/service"
	"errors"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// PsychosocialContactController maneja los endpoints del flujo 3x3.
type PsychosocialContactController struct {
	svc service.Psychosocial3x3Service
}

// NewPsychosocialContactController construye el controlador.
func NewPsychosocialContactController(svc service.Psychosocial3x3Service) *PsychosocialContactController {
	return &PsychosocialContactController{svc: svc}
}

// RegisterRoutes registra las rutas del flujo 3x3 en /api/v1.
//
//	GET   /api/v1/psychosocial/:psicosocialId/contact-attempts
//	POST  /api/v1/psychosocial/:psicosocialId/contact-attempts
//	POST  /api/v1/psychosocial/:psicosocialId/sessions
//	PUT   /api/v1/psychosocial/:psicosocialId/next-attempt
//	PATCH /api/v1/contact-attempts/:id/consent
func (c *PsychosocialContactController) RegisterRoutes(api *gin.RouterGroup) {
	grp := api.Group("/psychosocial/:psicosocialId")
	{
		grp.GET("/contact-attempts", c.GetHistory)
		grp.POST("/contact-attempts", c.RegisterAttempt)
		grp.POST("/sessions", c.ScheduleSession)
		grp.PUT("/next-attempt", c.SetNextAttempt)
		grp.POST("/init-closure-form", c.InitClosureForm)
		grp.POST("/close", c.CloseProcess)
	}
	api.PATCH("/contact-attempts/:id/consent", c.SetConsent)
}

// professionalContext extrae (best-effort) el id y equipo de la profesional logueada.
func professionalContext(ctx *gin.Context) (professionalID, team string) {
	defer func() { _ = recover() }()
	session := sessions.Default(ctx)
	sessionID, ok := session.Get("userData").(string)
	if !ok || sessionID == "" {
		return "", ""
	}
	s, err := utils.GetCommonSession(sessionID)
	if err != nil || s == nil {
		return "", ""
	}
	return s.UserICode, ""
}

func parseRFC3339(v *string) *time.Time {
	if v == nil || *v == "" {
		return nil
	}
	parsed, err := time.Parse(time.RFC3339, *v)
	if err != nil {
		return nil
	}
	utc := parsed.UTC()
	return &utc
}

func (c *PsychosocialContactController) GetHistory(ctx *gin.Context) {
	psicosocialID := ctx.Param("psicosocialId")
	result, err := c.svc.GetHistory(ctx.Request.Context(), psicosocialID)
	if err != nil {
		if errors.Is(err, service.ErrPsychosocial3x3NotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "proceso psicosocial no encontrado"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *PsychosocialContactController) RegisterAttempt(ctx *gin.Context) {
	psicosocialID := ctx.Param("psicosocialId")

	var body struct {
		WasAnswered *bool   `json:"was_answered"`
		Note        *string `json:"note"`
		AttemptAt   *string `json:"attempt_at"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil || body.WasAnswered == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "was_answered es requerido"})
		return
	}

	profID, team := professionalContext(ctx)
	attempt, counters, processStatus, err := c.svc.RegisterAttempt(
		ctx.Request.Context(), psicosocialID, *body.WasAnswered, body.Note, parseRFC3339(body.AttemptAt),
		profID, team,
	)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrPsychosocial3x3NotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": "proceso psicosocial no encontrado"})
		case errors.Is(err, service.ErrMaxAttemptsReached):
			ctx.JSON(http.StatusConflict, gin.H{"error": "tope de 50 intentos alcanzado"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"attempt":       attempt,
		"counters":      counters,
		"processStatus": processStatus,
	})
}

func (c *PsychosocialContactController) SetConsent(ctx *gin.Context) {
	attemptID := ctx.Param("id")

	var body struct {
		ConsentGiven *bool `json:"consent_given"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil || body.ConsentGiven == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "consent_given es requerido"})
		return
	}

	result, err := c.svc.SetConsent(ctx.Request.Context(), attemptID, *body.ConsentGiven)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrContactAttemptNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": "intento no encontrado"})
		case errors.Is(err, service.ErrInvalidAttempt):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "el consentimiento solo aplica a un contacto exitoso"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *PsychosocialContactController) ScheduleSession(ctx *gin.Context) {
	psicosocialID := ctx.Param("psicosocialId")

	var body struct {
		Immediate     *bool   `json:"immediate"`
		ScheduledDate *string `json:"scheduled_date"`
		ScheduledTime *string `json:"scheduled_time"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil || body.Immediate == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "immediate es requerido"})
		return
	}
	if !*body.Immediate && (body.ScheduledDate == nil || body.ScheduledTime == nil) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "scheduled_date y scheduled_time son requeridos al agendar"})
		return
	}

	var scheduledAt *time.Time
	if !*body.Immediate && body.ScheduledDate != nil {
		combined := *body.ScheduledDate + "T" + firstNonEmptyStr(body.ScheduledTime, "00:00") + ":00Z"
		scheduledAt = parseRFC3339(&combined)
	}

	profID, team := professionalContext(ctx)
	result, err := c.svc.ScheduleSession(ctx.Request.Context(), psicosocialID, *body.Immediate, scheduledAt, body.ScheduledTime, profID, team)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrPsychosocial3x3NotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": "proceso psicosocial no encontrado"})
		case errors.Is(err, service.ErrConsentNotAccepted):
			ctx.JSON(http.StatusConflict, gin.H{"error": "consentimiento no aceptado"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusCreated, result)
}

func (c *PsychosocialContactController) SetNextAttempt(ctx *gin.Context) {
	psicosocialID := ctx.Param("psicosocialId")

	var body struct {
		NextContactAttemptAt *string `json:"next_contact_attempt_at"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil || body.NextContactAttemptAt == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "next_contact_attempt_at es requerido"})
		return
	}
	at := parseRFC3339(body.NextContactAttemptAt)
	if at == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "fecha/hora inválida"})
		return
	}

	saved, err := c.svc.SetNextAttempt(ctx.Request.Context(), psicosocialID, *at)
	if err != nil {
		if errors.Is(err, service.ErrPsychosocial3x3NotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "proceso psicosocial no encontrado"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"psicosocialId":        psicosocialID,
		"nextContactAttemptAt": saved,
	})
}

func (c *PsychosocialContactController) InitClosureForm(ctx *gin.Context) {
	psicosocialID := ctx.Param("psicosocialId")
	reason := ctx.Query("reason")

	init, err := c.svc.InitClosureForm(ctx.Request.Context(), psicosocialID, reason)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidReason):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "reason inválido"})
		case errors.Is(err, service.ErrPsychosocial3x3NotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": "proceso psicosocial no encontrado"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"submission_id":     init.SubmissionID,
		"form_id":           init.FormID,
		"preselectedMotivo": init.PreselectedMotivo,
	})
}

func (c *PsychosocialContactController) CloseProcess(ctx *gin.Context) {
	psicosocialID := ctx.Param("psicosocialId")

	var body struct {
		SubmissionID *string `json:"submission_id"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil || body.SubmissionID == nil || *body.SubmissionID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "submission_id es requerido"})
		return
	}

	status, motivo, err := c.svc.CloseProcess(ctx.Request.Context(), psicosocialID, *body.SubmissionID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidReason):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "la submission no tiene Motivo de cierre respondido"})
		case errors.Is(err, service.ErrPsychosocial3x3NotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": "proceso psicosocial no encontrado"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"psicosocialId": psicosocialID,
		"status":        status,
		"motivo":        motivo,
	})
}

func firstNonEmptyStr(v *string, def string) string {
	if v == nil || *v == "" {
		return def
	}
	return *v
}
