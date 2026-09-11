package controller

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"bitsflow/salvia/service"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PsychosocialSupportGinController expone la creación manual de remisiones a
// Atención Psicosocial (fuera del flujo de seguimiento), usando Gin.
type PsychosocialSupportGinController struct {
	svc            service.PsychosocialSupportService
	victimCaseRepo repository.VictimCaseLightRepository
}

func NewPsychosocialSupportGinController(svc service.PsychosocialSupportService, victimCaseRepo repository.VictimCaseLightRepository) *PsychosocialSupportGinController {
	return &PsychosocialSupportGinController{svc: svc, victimCaseRepo: victimCaseRepo}
}

// RegisterRoutes registra las rutas de creación manual en el grupo /api/v1.
//
//	POST /api/v1/psychosocial-support/manual
func (c *PsychosocialSupportGinController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/psychosocial-support/manual", c.CreateManual)
}

type createManualPsychosocialSupportRequest struct {
	CaseID string `json:"caseId"`
	Cedula string `json:"cedula"`
}

// CreateManual crea una remisión a Atención Psicosocial directamente para un caso,
// sin pasar por el flujo de seguimiento (form-based) ni sus criterios de puntaje.
// Prioriza caseId sobre cedula cuando llegan ambos. No requiere autenticación.
//
//	POST /api/v1/psychosocial-support/manual
//	Body: { "caseId": "..." } o { "cedula": "..." }
func (c *PsychosocialSupportGinController) CreateManual(ctx *gin.Context) {
	var body createManualPsychosocialSupportRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "body inválido"})
		return
	}
	if body.CaseID == "" && body.Cedula == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "se requiere 'caseId' o 'cedula'"})
		return
	}

	var victimCase *models.VictimCaseLight
	var err error
	if body.CaseID != "" {
		victimCase, err = c.victimCaseRepo.FindByICode(ctx.Request.Context(), body.CaseID)
	} else {
		victimCase, err = c.victimCaseRepo.FindOpenByDocNumber(ctx.Request.Context(), body.Cedula)
	}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "no se encontró un caso para el caseId/cédula indicado"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}

	team := "Supervisor"
	ps := &models.PsychosocialSupport{
		CaseID:          victimCase.VictimICode,
		Type:            "derivacion",
		Status:          models.PsychosocialSupportStatusAbierto,
		SubmittedByTeam: &team,
	}
	if err := c.svc.Create(ctx.Request.Context(), ps); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}

	ctx.JSON(http.StatusCreated, ps)
}
