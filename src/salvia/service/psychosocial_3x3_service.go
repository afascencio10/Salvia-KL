// Package service — Psychosocial3x3Service: lógica de dominio del flujo 3x3 de
// Atención Psicosocial (intentos de contacto, umbral diario, tope 50, cierre,
// consentimiento, agendamiento de sesión y próximo intento).
package service

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

// Errores de dominio del flujo 3x3.
var (
	ErrPsychosocial3x3NotFound = errors.New("psychosocial process not found")
	ErrContactAttemptNotFound  = errors.New("contact attempt not found")
	ErrMaxAttemptsReached      = errors.New("max contact attempts reached")
	ErrConsentNotAccepted      = errors.New("informed consent not accepted")
	ErrInvalidAttempt          = errors.New("invalid contact attempt")
	ErrInvalidReason           = errors.New("invalid closure reason")
)

// Formulario de Cierre de proceso psicosocial (ver closure-form-model.md).
const (
	ClosureFormID  = "fcc7dc8d-835b-4559-9283-2ea8b8e6092b"
	ClosureQMotivo = "1c6d97ac-dc83-4e73-b747-14cd93ad422d" // pregunta "Motivo de cierre"
)

// validClosureReasons: disparadores permitidos hoy para iniciar el cierre.
var validClosureReasons = map[string]bool{
	"no_consentimiento":          true,
	"imposibilidad_contacto_3x3": true,
}

// closureReasonToStatus mapea el Motivo de cierre al estado resultante (RN-12).
var closureReasonToStatus = map[string]string{
	"imposibilidad_contacto_3x3": models.PsychosocialSupportStatusCerrado,
	"cumplimiento_objetivos":     models.PsychosocialSupportStatusCerrado,
	"cumplimiento_esquema":       models.PsychosocialSupportStatusCerrado,
	"no_consentimiento":          models.PsychosocialSupportStatusEnDevolucion,
	"desistimiento_proceso":      models.PsychosocialSupportStatusEnDevolucion,
}

// ClosureFormInit es la respuesta de InitClosureForm.
type ClosureFormInit struct {
	SubmissionID      string `json:"submission_id"`
	FormID            string `json:"form_id"`
	PreselectedMotivo string `json:"preselectedMotivo"`
}

// Umbrales del flujo 3x3.
const (
	DailyThreshold   = 3  // intentos fallidos por día que abren el modal de acciones
	MaxAttempts      = 50 // tope total de intentos por proceso
	ClosureThreshold = 9  // intentos mínimos para elegibilidad de cierre
	ClosureMinDays   = 3  // días distintos mínimos para elegibilidad de cierre
)

// Counters agrupa los contadores que gobiernan la UI del flujo 3x3.
type Counters struct {
	DailyFailedCount      int64 `json:"dailyFailedCount"`
	TotalCount            int64 `json:"totalCount"`
	DistinctDaysCount     int64 `json:"distinctDaysCount"`
	DailyThresholdReached bool  `json:"dailyThresholdReached"`
	MaxAttemptsReached    bool  `json:"maxAttemptsReached"`
	ClosureEligible       bool  `json:"closureEligible"`
}

// AttemptHistoryItem es una fila del historial que consume el card.
type AttemptHistoryItem struct {
	ID             string    `json:"id"`
	SequenceNumber int       `json:"sequenceNumber"`
	WasAnswered    bool      `json:"wasAnswered"`
	Note           *string   `json:"note"`
	ConsentGiven   *bool     `json:"consentGiven"`
	AttemptAt      time.Time `json:"attemptAt"`
}

// HistoryResult es la respuesta del GET del historial.
type HistoryResult struct {
	Attempts             []AttemptHistoryItem `json:"attempts"`
	Counters             Counters             `json:"counters"`
	LastAttemptAt        *time.Time           `json:"lastAttemptAt"`
	NextContactAttemptAt *time.Time           `json:"nextContactAttemptAt"`
	ProcessStatus        string               `json:"processStatus"`
	CaseID               string               `json:"caseId"`
	CaseName             string               `json:"caseName"`
}

// SessionResult es la respuesta del agendamiento de sesión.
type SessionResult struct {
	Session     *models.TeamContact `json:"session"`
	RedirectURL *string             `json:"redirectUrl,omitempty"`
}

// ConsentResult es la respuesta del registro de consentimiento.
type ConsentResult struct {
	ID                  string `json:"id"`
	ConsentGiven        bool   `json:"consentGiven"`
	RequiresClosureForm bool   `json:"requiresClosureForm"`
	CanScheduleSession  bool   `json:"canScheduleSession"`
}

// Psychosocial3x3Service expone las operaciones del flujo 3x3.
type Psychosocial3x3Service interface {
	GetHistory(ctx context.Context, psicosocialID string) (*HistoryResult, error)
	RegisterAttempt(ctx context.Context, psicosocialID string, wasAnswered bool, note *string, attemptAt *time.Time, professionalID, team string) (*models.ContactAttempt, *Counters, string, error)
	SetConsent(ctx context.Context, attemptID string, consentGiven bool) (*ConsentResult, error)
	ScheduleSession(ctx context.Context, psicosocialID string, immediate bool, scheduledAt *time.Time, scheduledTime *string, professionalID, team string) (*SessionResult, error)
	SetNextAttempt(ctx context.Context, psicosocialID string, at time.Time) (*time.Time, error)
	InitClosureForm(ctx context.Context, psicosocialID, reason string) (*ClosureFormInit, error)
	CloseProcess(ctx context.Context, psicosocialID, submissionID string) (status, motivo string, err error)
}

type psychosocial3x3Service struct {
	repo repository.ContactAttemptRepository
	db   *gorm.DB
}

// NewPsychosocial3x3Service construye el servicio.
func NewPsychosocial3x3Service(repo repository.ContactAttemptRepository, db *gorm.DB) Psychosocial3x3Service {
	return &psychosocial3x3Service{repo: repo, db: db}
}

func (s *psychosocial3x3Service) loadProcess(ctx context.Context, psicosocialID string) (*models.PsychosocialSupport, error) {
	var p models.PsychosocialSupport
	err := s.db.WithContext(ctx).Where("id = ?", psicosocialID).First(&p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPsychosocial3x3NotFound
		}
		return nil, err
	}
	return &p, nil
}

func (s *psychosocial3x3Service) buildCounters(raw repository.ContactAttemptCounters) Counters {
	return Counters{
		DailyFailedCount:      raw.DailyFailedCount,
		TotalCount:            raw.TotalCount,
		DistinctDaysCount:     raw.DistinctDaysCount,
		DailyThresholdReached: raw.DailyFailedCount >= DailyThreshold,
		MaxAttemptsReached:    raw.TotalCount >= MaxAttempts,
		ClosureEligible:       raw.TotalCount >= ClosureThreshold && raw.DistinctDaysCount >= ClosureMinDays,
	}
}

func (s *psychosocial3x3Service) caseName(ctx context.Context, caseICode string) string {
	var name string
	s.db.WithContext(ctx).Raw(
		`SELECT TRIM(COALESCE(victim_case_victim_names,'') || ' ' || COALESCE(victim_case_victim_last_names,''))
		 FROM salvia.victim_case WHERE victim_case_i_code = ? LIMIT 1`, caseICode).Scan(&name)
	return name
}

func (s *psychosocial3x3Service) GetHistory(ctx context.Context, psicosocialID string) (*HistoryResult, error) {
	process, err := s.loadProcess(ctx, psicosocialID)
	if err != nil {
		return nil, err
	}

	attempts, err := s.repo.GetByPsicosocialID(ctx, psicosocialID)
	if err != nil {
		return nil, err
	}

	items := make([]AttemptHistoryItem, 0, len(attempts))
	var lastAttemptAt *time.Time
	for i, a := range attempts {
		items = append(items, AttemptHistoryItem{
			ID:             a.ID,
			SequenceNumber: i + 1,
			WasAnswered:    a.WasAnswered,
			Note:           a.Note,
			ConsentGiven:   a.ConsentGiven,
			AttemptAt:      a.AttemptAt,
		})
		at := a.AttemptAt
		if lastAttemptAt == nil || at.After(*lastAttemptAt) {
			cp := at
			lastAttemptAt = &cp
		}
	}

	// Orden de presentación: más nuevo → más viejo. El sequenceNumber permanece
	// cronológico (1 = intento más antiguo) aunque se muestre de último.
	for l, r := 0, len(items)-1; l < r; l, r = l+1, r-1 {
		items[l], items[r] = items[r], items[l]
	}

	rawCounters, err := s.repo.Counters(ctx, psicosocialID)
	if err != nil {
		return nil, err
	}

	return &HistoryResult{
		Attempts:             items,
		Counters:             s.buildCounters(rawCounters),
		LastAttemptAt:        lastAttemptAt,
		NextContactAttemptAt: process.NextContactAttemptAt,
		ProcessStatus:        process.Status,
		CaseID:               process.CaseID,
		CaseName:             s.caseName(ctx, process.CaseID),
	}, nil
}

func (s *psychosocial3x3Service) RegisterAttempt(ctx context.Context, psicosocialID string, wasAnswered bool, note *string, attemptAt *time.Time, professionalID, team string) (*models.ContactAttempt, *Counters, string, error) {
	process, err := s.loadProcess(ctx, psicosocialID)
	if err != nil {
		return nil, nil, "", err
	}

	rawCounters, err := s.repo.Counters(ctx, psicosocialID)
	if err != nil {
		return nil, nil, "", err
	}
	if rawCounters.TotalCount >= MaxAttempts {
		return nil, nil, "", ErrMaxAttemptsReached
	}

	when := time.Now().UTC()
	if attemptAt != nil {
		when = attemptAt.UTC()
	}

	attempt := &models.ContactAttempt{
		PsicosocialID: psicosocialID,
		CaseID:        process.CaseID,
		WasAnswered:   wasAnswered,
		Note:          note,
		AttemptAt:     when,
	}
	if professionalID != "" {
		attempt.ProfessionalID = &professionalID
	}
	if team != "" {
		attempt.Team = &team
	}
	if err := s.repo.Create(ctx, attempt); err != nil {
		return nil, nil, "", err
	}

	// Contacto exitoso → el proceso pasa a en_gestion (si estaba abierto).
	if wasAnswered && process.Status == models.PsychosocialSupportStatusAbierto {
		if err := s.db.WithContext(ctx).
			Model(&models.PsychosocialSupport{}).
			Where("id = ?", psicosocialID).
			Update("status", models.PsychosocialSupportStatusEnGestion).Error; err != nil {
			return nil, nil, "", err
		}
		process.Status = models.PsychosocialSupportStatusEnGestion
	}

	// Si había un próximo intento programado (next_contact_attempt_at) cuya fecha
	// es igual o anterior al día del intento recién registrado, se limpia: el
	// intento planeado ya se ejecutó (sin importar si contestó o no; cubre el caso
	// de que el agente lo hiciera el mismo día o días después de lo programado).
	// Solo se conserva si lo programado es a futuro respecto al intento.
	if process.NextContactAttemptAt != nil {
		if err := s.db.WithContext(ctx).Exec(
			`UPDATE salvia.psychosocial_support SET next_contact_attempt_at = NULL, updated_at = now()
			 WHERE id = ? AND next_contact_attempt_at IS NOT NULL AND next_contact_attempt_at::date <= ?::date`,
			psicosocialID, when).Error; err != nil {
			return nil, nil, "", err
		}
	}

	newRaw, err := s.repo.Counters(ctx, psicosocialID)
	if err != nil {
		return nil, nil, "", err
	}
	counters := s.buildCounters(newRaw)
	return attempt, &counters, process.Status, nil
}

func (s *psychosocial3x3Service) SetConsent(ctx context.Context, attemptID string, consentGiven bool) (*ConsentResult, error) {
	attempt, err := s.repo.GetByID(ctx, attemptID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrContactAttemptNotFound
		}
		return nil, err
	}
	if !attempt.WasAnswered {
		return nil, ErrInvalidAttempt
	}
	if err := s.repo.UpdateConsent(ctx, attemptID, consentGiven); err != nil {
		return nil, err
	}
	return &ConsentResult{
		ID:                  attemptID,
		ConsentGiven:        consentGiven,
		RequiresClosureForm: !consentGiven,
		CanScheduleSession:  consentGiven,
	}, nil
}

func (s *psychosocial3x3Service) ScheduleSession(ctx context.Context, psicosocialID string, immediate bool, scheduledAt *time.Time, scheduledTime *string, professionalID, team string) (*SessionResult, error) {
	process, err := s.loadProcess(ctx, psicosocialID)
	if err != nil {
		return nil, err
	}

	// Requiere consentimiento aceptado (algún intento exitoso con consent_given = true).
	var consented int64
	if err := s.db.WithContext(ctx).
		Model(&models.ContactAttempt{}).
		Where("psicosocial_id = ? AND was_answered = true AND consent_given = true", psicosocialID).
		Count(&consented).Error; err != nil {
		return nil, err
	}
	if consented == 0 {
		return nil, ErrConsentNotAccepted
	}

	when := time.Now().UTC()
	if !immediate && scheduledAt != nil {
		when = scheduledAt.UTC()
	}

	isPsico := true
	session := &models.TeamContact{
		CaseID:         process.CaseID,
		PsicosocialID:  &psicosocialID,
		ScheduledDate:  &when,
		ScheduledTime:  scheduledTime,
		IsCompleted:    false,
		IsPsicoSession: isPsico,
	}
	if professionalID != "" {
		session.ProfessionalID = &professionalID
	}
	if team != "" {
		session.Team = &team
	}
	if process.DuplaID != nil {
		session.DuplaID = process.DuplaID
	}
	if err := s.db.WithContext(ctx).Create(session).Error; err != nil {
		return nil, err
	}

	result := &SessionResult{Session: session}
	if immediate {
		url := "/salvia/psicosocial/sesion/" + session.ID
		result.RedirectURL = &url
	}
	return result, nil
}

func (s *psychosocial3x3Service) SetNextAttempt(ctx context.Context, psicosocialID string, at time.Time) (*time.Time, error) {
	if _, err := s.loadProcess(ctx, psicosocialID); err != nil {
		return nil, err
	}
	utc := at.UTC()
	if err := s.db.WithContext(ctx).
		Model(&models.PsychosocialSupport{}).
		Where("id = ?", psicosocialID).
		Update("next_contact_attempt_at", utc).Error; err != nil {
		return nil, err
	}
	return &utc, nil
}

// InitClosureForm crea la form_submission del formulario de cierre y pre-responde
// el Motivo según el disparador (queda preseleccionado pero editable).
func (s *psychosocial3x3Service) InitClosureForm(ctx context.Context, psicosocialID, reason string) (*ClosureFormInit, error) {
	if !validClosureReasons[reason] {
		return nil, ErrInvalidReason
	}
	if _, err := s.loadProcess(ctx, psicosocialID); err != nil {
		return nil, err
	}

	submission := &models.FormSubmission{FormID: ClosureFormID}
	if err := s.db.WithContext(ctx).Create(submission).Error; err != nil {
		return nil, err
	}
	// Pre-respuesta del Motivo (editable por la profesional en el dinamic-form).
	answer := &models.Answer{FormSubmissionID: submission.ID, QuestionID: ClosureQMotivo, Value: reason}
	if err := s.db.WithContext(ctx).Create(answer).Error; err != nil {
		return nil, err
	}

	return &ClosureFormInit{
		SubmissionID:      submission.ID,
		FormID:            ClosureFormID,
		PreselectedMotivo: reason,
	}, nil
}

// CloseProcess lee el Motivo de la submission y transiciona el proceso a
// cerrado / en_devolucion. Afecta solo al proceso psicosocial.
func (s *psychosocial3x3Service) CloseProcess(ctx context.Context, psicosocialID, submissionID string) (string, string, error) {
	if _, err := s.loadProcess(ctx, psicosocialID); err != nil {
		return "", "", err
	}

	var motivos []string
	if err := s.db.WithContext(ctx).
		Model(&models.Answer{}).
		Where("form_submission_id = ? AND question_id = ?", submissionID, ClosureQMotivo).
		Order("updated_at DESC").
		Limit(1).
		Pluck("value", &motivos).Error; err != nil {
		return "", "", err
	}
	if len(motivos) == 0 || motivos[0] == "" {
		return "", "", ErrInvalidReason
	}
	motivo := motivos[0]

	status, ok := closureReasonToStatus[motivo]
	if !ok {
		status = models.PsychosocialSupportStatusCerrado
	}

	if err := s.db.WithContext(ctx).
		Model(&models.PsychosocialSupport{}).
		Where("id = ?", psicosocialID).
		Update("status", status).Error; err != nil {
		return "", "", err
	}
	return status, motivo, nil
}
