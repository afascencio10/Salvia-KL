package service

import (
	"bitsflow/common/utils"
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// RegistroCasoFormID es el form_id fijo del formulario dinámico "Registro de Caso"
// (creado por src/cmd/seed/seed_registro_caso.sql). Análogo a seguimientoFormID.
const RegistroCasoFormID = "0a24ab30-3cfc-4861-b74d-65d21524bc00"

// Valores dummy para los campos obligatorios de victim_case que aún no se han
// respondido -- así el draft nunca se pospone, y se sobrescriben con la
// respuesta real en cuanto la Sección 1 se guarde (ver UpdateDraftCoreFields).
const (
	dummyNames     = "Pendiente"
	dummyLastNames = "Pendiente"
	dummyDocType   = "cc"
	dummyDocNumber = "0000000000"
	dummyTownCode  = "00000000"
)

// Umbrales de riesgo del tamizaje -- idénticos a los usados hoy en
// set_victim_case.html (updateTamizajeScore) y en getRiskScore (Go).
var riskThresholdsPartner = []struct{ max, level int }{
	{4, 1}, {8, 2}, {15, 3}, {24, 4},
}
var riskThresholdsNonPartner = []struct{ max, level int }{
	{2, 1}, {5, 2}, {8, 3}, {20, 4},
}

// VictimCaseFormResult es lo que el frontend recupera vía GET tras 'form-completed'.
type VictimCaseFormResult struct {
	Ready  bool   `json:"ready"`
	CaseID string `json:"caseId,omitempty"`
	Login  string `json:"login,omitempty"`
	Pass   string `json:"pass,omitempty"`
}

// VictimCaseFormService orquesta la creación progresiva (Borrador -> Activo)
// del formulario dinámico "Registro de Caso" -- ver DocsMD/Screens/Registro de Caso V2.
type VictimCaseFormService interface {
	// UpdateCaseDraft se llama en CADA guardado de sección (vía el hook
	// genérico OnSectionUpdate, no solo la primera). Crea el victim_case en
	// 'bo' (Borrador) desde el primer saveSection -- nunca se pospone: los
	// campos obligatorios de Sección 1 que aún no tengan respuesta se rellenan
	// con valores dummy (ver dummyNames/dummyDocType/etc.). En cada llamada
	// posterior sincroniza esos mismos campos con la respuesta real ya
	// disponible y proyecta todas las respuestas sobre victim_case_form2 (que
	// también usa dummies en sus columnas NOT NULL aún sin responder) -- el
	// "borrador" refleja siempre el progreso real, sin esperas.
	UpdateCaseDraft(ctx context.Context, submissionID, actorID string) error

	// Activate se llama solo al completar el formulario (OnEndFormSubmission).
	// victim_case_form2 ya debería existir (lo dejó UpdateCaseDraft en esta
	// misma request, justo antes) -- aquí solo se confirma, se activa el caso
	// (bo -> ra), se genera el calendario y se asigna equipo/agente.
	Activate(ctx context.Context, submissionID, actorID string) error

	// GetResult retorna las credenciales + caseId una vez el caso está activo.
	GetResult(ctx context.Context, submissionID string) (*VictimCaseFormResult, error)
}

type victimCaseFormService struct {
	formRepo         repository.VictimCaseFormRepository
	caseRepo         repository.VictimCaseLightRepository
	caseTimelineRepo repository.CaseTimelineEventRepository
	followUpV2Svc    FollowUpV2Service
}

func NewVictimCaseFormService(
	formRepo repository.VictimCaseFormRepository,
	caseRepo repository.VictimCaseLightRepository,
	caseTimelineRepo repository.CaseTimelineEventRepository,
	followUpV2Svc FollowUpV2Service,
) VictimCaseFormService {
	return &victimCaseFormService{
		formRepo:         formRepo,
		caseRepo:         caseRepo,
		caseTimelineRepo: caseTimelineRepo,
		followUpV2Svc:    followUpV2Svc,
	}
}

// ─── onSectionUpdate → UpdateCaseDraft (cada sección guardada) ─────────────────

func (s *victimCaseFormService) UpdateCaseDraft(ctx context.Context, submissionID, actorID string) error {
	answers, err := s.formRepo.BuildAnswersByFieldKey(ctx, RegistroCasoFormID, submissionID)
	if err != nil {
		return fmt.Errorf("updateCaseDraft: leer respuestas: %w", err)
	}

	iCode, found, err := s.formRepo.FindICodeBySubmissionID(ctx, submissionID)
	if err != nil {
		return fmt.Errorf("updateCaseDraft: verificar existencia del caso: %w", err)
	}

	if !found {
		newICode, err := s.createDraftShell(ctx, submissionID, answers)
		if err != nil {
			return err
		}
		iCode = newICode
	} else if err := s.formRepo.UpdateDraftCoreFields(ctx, iCode,
		coalesce(answers[repository.FieldKey(1, 1)], dummyNames),
		coalesce(answers[repository.FieldKey(1, 2)], dummyLastNames),
		coalesce(answers[repository.FieldKey(1, 5)], dummyDocType),
		coalesce(answers[repository.FieldKey(1, 6)], dummyDocNumber),
		coalesce(answers[repository.FieldKey(1, 10)], dummyTownCode),
	); err != nil {
		return fmt.Errorf("updateCaseDraft: sincronizar datos base del caso: %w", err)
	}

	wasPartner, err := s.resolveWasPartner(ctx, answers)
	if err != nil {
		return fmt.Errorf("updateCaseDraft: resolver wasPartner: %w", err)
	}
	riskScore, riskLevel := computeRiskScore(answers, wasPartner)

	created, err := s.formRepo.UpsertForm2FromAnswers(ctx, iCode, answers, riskScore, riskLevel)
	if err != nil {
		return fmt.Errorf("updateCaseDraft: proyectar respuestas sobre form2: %w", err)
	}
	if created {
		log.Printf("[victimCaseForm] victim_case_form2 creado para caso=%s", iCode)
	}

	return nil
}

// createDraftShell crea el general_user + victim_case en 'bo' de inmediato, en
// el primer saveSection que llegue -- nunca se pospone. Los campos mínimos de
// Sección 1 que aún no se hayan respondido se rellenan con valores dummy
// (ver dummyNames/dummyLastNames/etc.); UpdateCaseDraft los sincroniza con los
// valores reales en cada guardado de sección posterior.
func (s *victimCaseFormService) createDraftShell(ctx context.Context, submissionID string, answers map[string]string) (string, error) {
	names := coalesce(answers[repository.FieldKey(1, 1)], dummyNames)
	lastNames := coalesce(answers[repository.FieldKey(1, 2)], dummyLastNames)
	docType := coalesce(answers[repository.FieldKey(1, 5)], dummyDocType)
	docNumber := coalesce(answers[repository.FieldKey(1, 6)], dummyDocNumber)
	residenceTown := coalesce(answers[repository.FieldKey(1, 10)], dummyTownCode)

	login := randomLogin(6)
	plainPassword := randomPassword(4)
	hashed, err := bcrypt.GenerateFromPassword([]byte(plainPassword), 10)
	if err != nil {
		return "", fmt.Errorf("createDraftShell: generar hash de password: %w", err)
	}

	newICode, err := s.formRepo.CreateDraft(ctx, repository.CreateDraftInput{
		SubmissionID:       submissionID,
		Names:              names,
		LastNames:          lastNames,
		DocType:            docType,
		DocNumber:          docNumber,
		ResidenceTownCode:  residenceTown,
		GeneralUserICode:   utils.GetUUID(),
		GeneralUserLogin:   login,
		GeneralUserPassSha: string(hashed),
	})
	if err != nil {
		return "", fmt.Errorf("createDraftShell: crear victim_case borrador: %w", err)
	}

	if err := s.formRepo.StoreCredentials(ctx, newICode, login, plainPassword); err != nil {
		return "", fmt.Errorf("createDraftShell: guardar credenciales: %w", err)
	}

	log.Printf("[victimCaseForm] draft creado: caso=%s submission=%s login=%s", newICode, submissionID, login)
	return newICode, nil
}

// ─── OnEndFormSubmission → Activate (al completar el formulario) ──────────────

func (s *victimCaseFormService) Activate(ctx context.Context, submissionID, actorID string) error {
	iCode, found, err := s.formRepo.FindICodeBySubmissionID(ctx, submissionID)
	if err != nil {
		return fmt.Errorf("activate: buscar victim_case: %w", err)
	}
	if !found {
		return fmt.Errorf("activate: no existe victim_case para submission=%s (UpdateCaseDraft nunca llegó a crearlo)", submissionID)
	}

	answers, err := s.formRepo.BuildAnswersByFieldKey(ctx, RegistroCasoFormID, submissionID)
	if err != nil {
		return fmt.Errorf("activate: leer respuestas: %w", err)
	}

	wasPartner, err := s.resolveWasPartner(ctx, answers)
	if err != nil {
		return fmt.Errorf("activate: resolver wasPartner: %w", err)
	}
	riskScore, riskLevel := computeRiskScore(answers, wasPartner)

	// Defensivo: UpdateCaseDraft (vía OnSectionUpdate) ya debería haber
	// dejado form2 creado en esta misma request, justo antes de que
	// allAnswered diera true -- se repite aquí por si acaso, es idempotente.
	if _, err := s.formRepo.UpsertForm2FromAnswers(ctx, iCode, answers, riskScore, riskLevel); err != nil {
		return fmt.Errorf("activate: proyectar respuestas: %w", err)
	}

	if err := s.formRepo.MarkActive(ctx, iCode, answers); err != nil {
		return fmt.Errorf("activate: marcar caso activo: %w", err)
	}
	log.Printf("[victimCaseForm] caso=%s activado — wasPartner=%v riskScore=%d riskLevel=%d", iCode, wasPartner, riskScore, riskLevel)

	if s.followUpV2Svc != nil {
		followUps, calErr := s.followUpV2Svc.GenerateOrRecalculate(ctx, iCode, GenerateCalendarInput{
			RiskLevel: riskLevel,
			AgentID:   "",
			Team:      "",
		})
		if calErr != nil {
			log.Printf("[WARN] [victimCaseForm] error generando calendario para caso %s: %v", iCode, calErr)
		} else if len(followUps) > 0 {
			team := followUps[0].Team
			agentID := ""
			if followUps[0].AgentID != nil {
				agentID = *followUps[0].AgentID
			}
			if err := s.caseRepo.UpdateTeamAndAgent(ctx, iCode, team, agentID); err != nil {
				log.Printf("[WARN] [victimCaseForm] no se pudo actualizar team/agent en caso %s: %v", iCode, err)
			}
		}
	}

	if s.caseTimelineRepo != nil {
		event := &models.CaseTimelineEvent{
			CaseID:      iCode,
			Category:    "General",
			Type:        "Caso Activado",
			Icon:        "folder-plus",
			Color:       "#22c55e",
			EventUserID: actorID,
			Date:        time.Now(),
		}
		if err := s.caseTimelineRepo.Create(ctx, event); err != nil {
			log.Printf("[WARN] [victimCaseForm] no se pudo crear evento de timeline para caso %s: %v", iCode, err)
		}
	}

	return nil
}

func (s *victimCaseFormService) GetResult(ctx context.Context, submissionID string) (*VictimCaseFormResult, error) {
	iCode, found, err := s.formRepo.FindICodeBySubmissionID(ctx, submissionID)
	if err != nil {
		return nil, err
	}
	if !found {
		return &VictimCaseFormResult{Ready: false}, nil
	}

	login, pass, err := s.formRepo.GetCredentials(ctx, iCode)
	if err != nil {
		return nil, err
	}
	if login == "" {
		return &VictimCaseFormResult{Ready: false}, nil
	}

	vcase, err := s.caseRepo.FindByICode(ctx, iCode)
	if err != nil {
		return nil, err
	}
	if vcase.Status != "ra" {
		// El draft y las credenciales ya existen, pero el caso aún no se activó
		// (E-04 todavía corriendo, o el formulario no se ha completado).
		return &VictimCaseFormResult{Ready: false}, nil
	}

	return &VictimCaseFormResult{Ready: true, CaseID: iCode, Login: login, Pass: pass}, nil
}

// ─── Helpers ────────────────────────────────────────────────────────────────────

func (s *victimCaseFormService) resolveWasPartner(ctx context.Context, answers map[string]string) (bool, error) {
	icode := answers[repository.FieldKey(4, 3)] // Relación con el presunto agresor
	if icode == "" {
		return false, nil
	}
	code, err := s.formRepo.ResolveEnumCode(ctx, icode)
	if err != nil {
		return false, err
	}
	return code == "pi" || code == "ex", nil
}

// computeRiskScore replica la fórmula de tamizaje: suma de "Sí" entre las 6
// preguntas comunes + las 18 (pareja) o 14 (no-pareja) específicas, y determina
// el nivel según la tabla de umbrales -- ver registro-caso-interface.md.
func computeRiskScore(answers map[string]string, wasPartner bool) (score int, level int) {
	for q := 1; q <= 6; q++ {
		if answers[repository.FieldKey(5, q)] == "true" {
			score++
		}
	}

	start, end := 25, 38 // no-pareja
	if wasPartner {
		start, end = 7, 24
	}
	for q := start; q <= end; q++ {
		if answers[repository.FieldKey(5, q)] == "true" {
			score++
		}
	}

	thresholds := riskThresholdsNonPartner
	if wasPartner {
		thresholds = riskThresholdsPartner
	}
	level = len(thresholds) // por si supera todos los umbrales, queda en el más alto
	for _, t := range thresholds {
		if score <= t.max {
			level = t.level
			break
		}
	}
	return score, level
}

func randomLogin(size int) string {
	rand.Seed(time.Now().UnixNano())
	const letters = "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, size)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b) + strconv.Itoa(rand.Intn(1000))
}

func randomPassword(size int) string {
	rand.Seed(time.Now().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyz1234567890"
	b := make([]byte, size)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}
