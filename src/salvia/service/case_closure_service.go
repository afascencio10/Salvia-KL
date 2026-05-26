package service

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"fmt"
	"log"
	"time"
)

// ─── CasoCierreService ────────────────────────────────────────────────────────
// Servicio modular para cerrar un caso VBG.
// Puede ser llamado desde cualquier controller o servicio que requiera cerrar
// un caso (formulario de seguimiento, detalle del caso, etc.).

const VictimCaseStatusCerrado = "cd"

// CerrarCasoInput reúne todos los datos necesarios para registrar el cierre.
type CerrarCasoInput struct {
	CaseICode               string // victim_case_i_code
	Motivo                  string // valor de la opción seleccionada (e.g. "perdida_contacto")
	OtroMotivo              string // texto libre si Motivo == "otro"
	Descripcion             string // descripción analítica del cierre
	AccionesInstitucionales bool   // si se realizaron acciones institucionales
	ActorID                 string // i_code del usuario que ejecuta el cierre
}

// CasoCierreService expone la lógica de cierre de caso.
type CasoCierreService interface {
	CerrarCaso(ctx context.Context, input CerrarCasoInput) error
}

type casoCierreService struct {
	victimCaseLightRepo repository.VictimCaseLightRepository
	caseTimelineRepo    repository.CaseTimelineEventRepository
}

// NewCasoCierreService crea la instancia del servicio.
func NewCasoCierreService(
	victimCaseLightRepo repository.VictimCaseLightRepository,
	caseTimelineRepo repository.CaseTimelineEventRepository,
) CasoCierreService {
	return &casoCierreService{
		victimCaseLightRepo: victimCaseLightRepo,
		caseTimelineRepo:    caseTimelineRepo,
	}
}

// CerrarCaso actualiza el status del caso a "cd" y registra el evento en el
// timeline. Es idempotente: si el caso ya está cerrado, solo loggea y retorna.
func (s *casoCierreService) CerrarCaso(ctx context.Context, input CerrarCasoInput) error {
	if input.CaseICode == "" {
		return fmt.Errorf("CerrarCaso: CaseICode es requerido")
	}

	// 1. Verificar estado actual del caso
	vc, err := s.victimCaseLightRepo.FindByICode(ctx, input.CaseICode)
	if err != nil {
		return fmt.Errorf("CerrarCaso: buscar caso [%s]: %w", input.CaseICode, err)
	}
	if vc.Status == VictimCaseStatusCerrado {
		log.Printf("[CerrarCaso] caso %s ya está cerrado — se omite la operación", input.CaseICode)
		return nil
	}

	// 2. Actualizar status a "cd"
	if err := s.victimCaseLightRepo.UpdateStatus(ctx, input.CaseICode, VictimCaseStatusCerrado); err != nil {
		return fmt.Errorf("CerrarCaso: actualizar status [%s]: %w", input.CaseICode, err)
	}
	log.Printf("[CerrarCaso] caso %s marcado como cerrado (cd)", input.CaseICode)

	// 3. Registrar evento en el timeline
	if s.caseTimelineRepo == nil {
		log.Printf("⚠️  [CerrarCaso] caseTimelineRepo es nil — evento de cierre NO creado para caso %s", input.CaseICode)
		return nil
	}

	desc := buildCierreDescription(input)
	now := time.Now()
	event := &models.CaseTimelineEvent{
		CaseID:      input.CaseICode,
		Category:    models.TimelineCategoryGeneral,
		Type:        models.TimelineTypeCierreCaso,
		Icon:        models.TimelineIconCierre,
		Color:       models.TimelineColorRed,
		Description: desc,
		EventUserID: input.ActorID,
		Date:        now,
		CreatedAt:   now,
	}
	if err := s.caseTimelineRepo.Create(ctx, event); err != nil {
		// No retornamos error: el cierre ya ocurrió en la BD.
		log.Printf("[CerrarCaso] advertencia: no se pudo crear evento timeline para caso %s: %v", input.CaseICode, err)
	}

	return nil
}

// buildCierreDescription construye el texto descriptivo del evento de cierre
// a partir de los datos del formulario.
func buildCierreDescription(input CerrarCasoInput) string {
	motivo := input.Motivo
	if motivo == "otro" && input.OtroMotivo != "" {
		motivo = input.OtroMotivo
	}

	desc := "Caso cerrado"
	if motivo != "" {
		desc = fmt.Sprintf("Caso cerrado. Motivo: %s", motivo)
	}
	if input.Descripcion != "" {
		desc += fmt.Sprintf(". %s", input.Descripcion)
	}
	return desc
}
