package controller

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"bitsflow/salvia/service"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PsychosocialSupportGinController expone la creación manual de remisiones a
// Atención Psicosocial (fuera del flujo de seguimiento), usando Gin.
type PsychosocialSupportGinController struct {
	svc            service.PsychosocialSupportService
	victimCaseRepo repository.VictimCaseLightRepository
	caseTaskRepo   repository.CaseTaskRepository
}

func NewPsychosocialSupportGinController(
	svc service.PsychosocialSupportService,
	victimCaseRepo repository.VictimCaseLightRepository,
	caseTaskRepo repository.CaseTaskRepository,
) *PsychosocialSupportGinController {
	return &PsychosocialSupportGinController{svc: svc, victimCaseRepo: victimCaseRepo, caseTaskRepo: caseTaskRepo}
}

// RegisterRoutes registra las rutas de creación manual en el grupo /api/v1.
//
//	POST /api/v1/psychosocial-support/manual
func (c *PsychosocialSupportGinController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/psychosocial-support/manual", c.CreateManual)
}

type createManualPsychosocialSupportRequest struct {
	CaseID    string   `json:"caseId"`
	Cedula    string   `json:"cedula"`
	Criterios []string `json:"criterios"`
	Evidencia string   `json:"evidencia"`
}

// CreateManual crea una remisión a Atención Psicosocial directamente para un caso,
// sin pasar por el flujo de seguimiento (form-based). Prioriza caseId sobre cedula
// cuando llegan ambos.
//
// Criterios y evidencia son obligatorios: sin ellos el equipo psicosocial no tiene
// con qué validar la remisión. No se aplica la regla de puntaje del formulario de
// seguimiento — esta vía es para casos ya en proceso y la procedencia la decide
// quien valida.
//
//	POST /api/v1/psychosocial-support/manual
//	Body: { "caseId": "...", "criterios": ["..."], "evidencia": "..." }
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

	criterios := make([]string, 0, len(body.Criterios))
	for _, c := range body.Criterios {
		if c = strings.TrimSpace(c); c != "" {
			criterios = append(criterios, c)
		}
	}
	evidencia := strings.TrimSpace(body.Evidencia)
	if len(criterios) == 0 || evidencia == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "se requieren los criterios de la remisión y la evidencia"})
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
	criteriosGuardados := strings.Join(criterios, ",")
	ps := &models.PsychosocialSupport{
		CaseID:          victimCase.VictimICode,
		Type:            "derivacion",
		Status:          models.PsychosocialSupportStatusAbierto,
		SubmittedByTeam: &team,
		// Mismo formato que la respuesta del formulario, para que el detalle los lea igual.
		ReferralCriteria: &criteriosGuardados,
		ReferralEvidence: &evidencia,
	}
	if err := c.svc.Create(ctx.Request.Context(), ps); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}

	// La remisión manual también necesita su tarea de validación. La pantalla de
	// detalle habilita "recepcionar / devolver" solo cuando existe una case_task
	// validar_remision en ToDo (ver validationPending en
	// psychosocial_detail_service.go), así que sin esta tarea la remisión nace sin
	// forma de validarse — que es el bug que esto corrige. La vía del seguimiento
	// ya la creaba; esta se había quedado atrás.
	if err := c.crearTareaValidacion(ctx, ps); err != nil {
		// No se devuelve error: la remisión ya existe y anularla sería peor. Se
		// registra para que la reparación sea rastreable.
		log.Printf("[PsychosocialSupportGin] WARN: remisión %s creada sin case_task validar_remision: %v", ps.ID, err)
	}

	ctx.JSON(http.StatusCreated, ps)
}

// crearTareaValidacion crea la case_task con la que el equipo psicosocial
// recepciona o devuelve la remisión. Mismos campos que la que crea el flujo de
// seguimiento (ver processFollowUpSubmission), sin follow_up_id porque esta
// remisión no nace de un seguimiento. Queda sin asignar: todavía no hay dupla ni
// profesional, igual que en la otra vía.
func (c *PsychosocialSupportGinController) crearTareaValidacion(ctx *gin.Context, ps *models.PsychosocialSupport) error {
	if c.caseTaskRepo == nil {
		return errors.New("caseTaskRepo no inyectado")
	}
	psID := ps.ID
	return c.caseTaskRepo.Create(ctx.Request.Context(), &models.CaseTask{
		Category:              "Psicosocial",
		Type:                  "validar_remision",
		Description:           "Validar que remision a psicosocial es valida",
		Status:                models.CaseTaskStatusToDo,
		CaseID:                ps.CaseID,
		PsychosocialSupportID: &psID,
	})
}
