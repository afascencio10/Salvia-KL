// Package service — psychosocial_detail_service.go
// Lógica de negocio para la pantalla de detalle de remisión psicosocial.
package service

import (
	"bitsflow/internal/constants"
	"bitsflow/internal/models"
	salvia_config "bitsflow/salvia/config"
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"
)

// ── Errores de dominio — evento E-01 (cargar pantalla de sesión) ──────────────

// ErrPsicosocialSessionNotFound se retorna cuando la remisión no existe o fue eliminada.
var ErrPsicosocialSessionNotFound = errors.New("psicosocialSession: remisión no encontrada")

// ErrPsicosocialSessionNotAssigned se retorna cuando el agente no es el profesional
// directo ni pertenece a la dupla asignada a la remisión.
var ErrPsicosocialSessionNotAssigned = errors.New("psicosocialSession: sesión no asignada a este profesional")

// ErrPsicosocialContactNotFound se retorna cuando contact_id no existe o no pertenece a la remisión.
var ErrPsicosocialContactNotFound = errors.New("psicosocialSession: contacto de sesión no encontrado")

// ErrPsicosocialContactNoSubmission se retorna al abrir un contacto sin form_submission asociado.
var ErrPsicosocialContactNoSubmission = errors.New("psicosocialSession: la sesión no tiene formulario asociado")

// PsychosocialDetailResponse agrupa los datos para la pantalla de detalle.
type PsychosocialDetailResponse struct {
	// Remisión
	ID             string  `json:"id"`
	CaseID         string  `json:"caseId"`
	Status         string  `json:"status"`
	Type           string  `json:"type"`
	Notes          *string `json:"notes"`
	CreatedAt      string  `json:"createdAt"`
	SessionCount   int     `json:"sessionCount"`
	SubmittedBy    string  `json:"submittedBy"`
	SubmittedByTeam string `json:"submittedByTeam"`

	// Víctima
	VictimNames     string `json:"victimNames"`
	VictimLastNames string `json:"victimLastNames"`
	DocNumber       string `json:"docNumber"`
	VictimPhone     string `json:"victimPhone"`
	Municipality    string `json:"municipality"`
	RiskLevel       *int   `json:"riskLevel"`

	// Profesional asignado
	ProfessionalID   *string `json:"professionalId"`
	ProfessionalName string  `json:"professionalName"`
	ProfessionalRole string  `json:"professionalRole"`

	// Remitente
	SubmittedByName string `json:"submittedByName"`

	// Campos adicionales
	RequiresInterpreter string  `json:"requiresInterpreter"`
	AjustesRazonables   string  `json:"ajustesRazonables"`
	ConsentStatus       string  `json:"consentStatus"`
	ConsentDate         string  `json:"consentDate"`
	SchedulePreference  *string `json:"schedulePreference"`
	ValidationPending   bool    `json:"validationPending"`
	ValidationCriteria  []string `json:"validationCriteria"`

	// Contactos/Sesiones
	Contacts []PsychosocialContactItem `json:"contacts"`
}

// PsychosocialContactItem representa un contacto/sesión individual.
type PsychosocialContactItem struct {
	ID            string  `json:"id"`
	ScheduledDate *string `json:"scheduledDate"`
	ScheduledTime *string `json:"scheduledTime"`
	IsCompleted   bool    `json:"isCompleted"`
	IsPsicoSession bool  `json:"isPsicoSession"`
	Status        *string `json:"status"`
	Summary       *string `json:"summary"`
	CompletedAt   *string `json:"completedAt"`
	SessionType   *string `json:"sessionType"`
	CreatedAt     string  `json:"createdAt"`
}

type PsychosocialDetailService interface {
	GetDetail(ctx context.Context, id string) (*PsychosocialDetailResponse, error)
	ValidateRemision(ctx context.Context, id string, isValid bool, reason string, actorID string) error
	CreateContact(ctx context.Context, psicosocialID, contactType, contactDate, contactTime, summary, sessionType, scheduledDate, scheduledTime string) (*models.TeamContact, error)
	RescheduleContact(ctx context.Context, contactID, scheduledDate, scheduledTime string) error
	CancelContact(ctx context.Context, contactID string) error
	UpdateSchedulePreference(ctx context.Context, psicosocialID, preference string) error
	LoadSession(ctx context.Context, psicosocialID, agentID, contactID string) (*LoadPsicosocialSessionResult, error)
	// CheckSessionAvailability valida si date+time (ventana 2h) está libre para el
	// profesional o la dupla de la remisión. mode: "individual" | "dupla" (vacío = inferir).
	CheckSessionAvailability(ctx context.Context, psicosocialID, date, timeStr, mode string) (available bool, message string, err error)
}

// ── Evento E-01: cargar pantalla "Registrar Sesión Psicosocial" ───────────────

// PsicosocialSessionVictimInfo es la información de la víctima mostrada en la tarjeta
// de la pantalla de sesión (mismos campos que VictimCaseInfo de hacer-seguimiento).
type PsicosocialSessionVictimInfo struct {
	Names             string `json:"Names"`
	LastNames         string `json:"LastNames"`
	Phone             string `json:"Phone"`
	GenderIdentity    string `json:"GenderIdentity"`
	SexualOrientation string `json:"SexualOrientation"`
	ContactPhone      string `json:"ContactPhone"`
	Age               *int64 `json:"Age"`
	TownName          string `json:"TownName"`
	RiskLevel         int    `json:"RiskLevel"`
	CaseICode         string `json:"CaseICode"`
}

// PsicosocialStateInfo resume el estado acumulado del proceso, usado por el frontend
// para el badge "Sesión X de 6" y el label del formulario cargado.
type PsicosocialStateInfo struct {
	YaHizoPrimerContacto  bool   `json:"yaHizoPrimerContacto"`
	YaHizoPrimeraAtencion bool   `json:"yaHizoPrimeraAtencion"`
	SessionCount          int    `json:"sessionCount"`
	Status                string `json:"status"`
}

// LoadPsicosocialSessionResult es la respuesta completa de GET /psychosocial-support/:id/load.
type LoadPsicosocialSessionResult struct {
	FormID           string                        `json:"formId"`
	FormType         string                        `json:"formType"`
	SubmissionID     string                        `json:"submissionId"`
	TeamContactID    string                        `json:"teamContactId,omitempty"`
	IsCompleted      bool                          `json:"isCompleted"`
	CanEdit          bool                          `json:"canEdit"`
	VictimInfo       *PsicosocialSessionVictimInfo `json:"victimInfo"`
	PsicosocialState PsicosocialStateInfo          `json:"psicosocialState"`
	FormState        map[string]interface{}        `json:"formState"`
}

type psychosocialDetailService struct {
	db *gorm.DB
}

func NewPsychosocialDetailService(db *gorm.DB) PsychosocialDetailService {
	return &psychosocialDetailService{db: db}
}

func (s *psychosocialDetailService) GetDetail(ctx context.Context, id string) (*PsychosocialDetailResponse, error) {
	// 1. Cargar remisión
	var ps models.PsychosocialSupport
	if err := s.db.WithContext(ctx).Where("id = ?", id).First(&ps).Error; err != nil {
		return nil, err
	}

	resp := &PsychosocialDetailResponse{
		ID:              ps.ID,
		CaseID:          ps.CaseID,
		Status:          ps.Status,
		Type:            ps.Type,
		Notes:           ps.Notes,
		CreatedAt:       ps.CreatedAt.Format("2006-01-02"),
		SessionCount:    ps.SessionCount,
		SubmittedByTeam: "",
		ProfessionalID:  ps.ProfessionalID,
	}

	if ps.SubmittedBy != nil {
		resp.SubmittedBy = *ps.SubmittedBy
	}
	if ps.SubmittedByTeam != nil {
		resp.SubmittedByTeam = *ps.SubmittedByTeam
	}

	// 2. Datos de la víctima
	type victimRow struct {
		Names     string `gorm:"column:names"`
		LastNames string `gorm:"column:last_names"`
		DocNumber string `gorm:"column:doc_number"`
		Phone     string `gorm:"column:phone"`
		Town      string `gorm:"column:town"`
		RiskLevel *int   `gorm:"column:risk_level"`
	}
	var victim victimRow
	s.db.WithContext(ctx).Raw(`
		SELECT
			COALESCE(vc.victim_case_victim_names, '') AS names,
			COALESCE(vc.victim_case_victim_last_names, '') AS last_names,
			COALESCE(vc.victim_case_victim_doc_number, '') AS doc_number,
			COALESCE(vf2.victim_case_form2_victim_phone::text, '') AS phone,
			COALESCE(t.town_name, '') AS town,
			vf2.victim_case_form2_risk_level AS risk_level
		FROM salvia.victim_case vc
		LEFT JOIN salvia.victim_case_form2 vf2 ON vf2.victim_case_form2_victim_case = vc.victim_case_id
		LEFT JOIN security.town t ON t.town_code = vc.victim_case_victim_town_code
		WHERE vc.victim_case_i_code = ?
		LIMIT 1
	`, ps.CaseID).Scan(&victim)

	resp.VictimNames = victim.Names
	resp.VictimLastNames = victim.LastNames
	resp.DocNumber = victim.DocNumber
	resp.VictimPhone = victim.Phone
	resp.Municipality = victim.Town
	resp.RiskLevel = victim.RiskLevel

	// 3. Nombre del profesional asignado
	if ps.ProfessionalID != nil && *ps.ProfessionalID != "" {
		var profName string
		s.db.WithContext(ctx).Raw(`
			SELECT COALESCE(gup.general_user_profile_names || ' ' || gup.general_user_profile_last_names, '')
			FROM security.general_user gu
			JOIN security.general_user_profile gup ON gup.general_user_profile_id = gu.general_user_general_user_profile
			WHERE gu.general_user_i_code = ?
		`, *ps.ProfessionalID).Scan(&profName)
		resp.ProfessionalName = strings.TrimSpace(profName)

		// Rol del profesional
		var roleName string
		s.db.WithContext(ctx).Raw(`
			SELECT COALESCE(r.role_name, '')
			FROM security.general_user gu
			JOIN security.rel_role_general_user rr ON rr.general_user_id = gu.general_user_id
			JOIN security.role r ON r.role_id = rr.role_id
			WHERE gu.general_user_i_code = ?
			LIMIT 1
		`, *ps.ProfessionalID).Scan(&roleName)
		resp.ProfessionalRole = roleName
	}

	// 4. Nombre del remitente
	if ps.SubmittedBy != nil && *ps.SubmittedBy != "" {
		var submitterName string
		s.db.WithContext(ctx).Raw(`
			SELECT COALESCE(gup.general_user_profile_names || ' ' || gup.general_user_profile_last_names, '')
			FROM security.general_user gu
			JOIN security.general_user_profile gup ON gup.general_user_profile_id = gu.general_user_general_user_profile
			WHERE gu.general_user_i_code = ?
		`, *ps.SubmittedBy).Scan(&submitterName)
		resp.SubmittedByName = strings.TrimSpace(submitterName)
	}

	// 5. Requiere intérprete (desde form2 de la víctima)
	var requiresInterpreter string
	s.db.WithContext(ctx).Raw(`
		SELECT COALESCE(ri.victim_case_form2_enums_name, '')
		FROM salvia.victim_case vc
		LEFT JOIN salvia.victim_case_form2 f2 ON f2.victim_case_form2_victim_case = vc.victim_case_id
		LEFT JOIN salvia.victim_case_form2_enums ri ON ri.victim_case_form2_enums_id = f2.victim_case_form2_require_language_interpreter
		WHERE vc.victim_case_i_code = ?
		LIMIT 1
	`, ps.CaseID).Scan(&requiresInterpreter)
	// Traducir el código del enum (ej: yes_no_n → No)
	if locale, ok := salvia_config.Locale["sp"]; ok {
		if translated, found := locale[requiresInterpreter]; found {
			requiresInterpreter = translated
		}
	}
	resp.RequiresInterpreter = requiresInterpreter

	// 5b. Ajustes razonables (multi-select del form2)
	var ajustesRazonables string
	s.db.WithContext(ctx).Raw(`
		SELECT COALESCE(string_agg(e.victim_case_form2_enums_name, ', '), '')
		FROM salvia.victim_case vc
		JOIN salvia.victim_case_form2 f2 ON f2.victim_case_form2_victim_case = vc.victim_case_id
		JOIN salvia.rel_victim_case_form2_enums_victim_case_form2 rel ON rel.victim_case_form2_id = f2.victim_case_form2_id
		JOIN salvia.victim_case_form2_enums e ON e.victim_case_form2_enums_id = rel.victim_case_form2_enums_id
		WHERE vc.victim_case_i_code = ?
		AND e.victim_case_form2_enums_category LIKE '%adjustments_gbv%'
	`, ps.CaseID).Scan(&ajustesRazonables)
	// Traducir cada valor del multi-select
	if ajustesRazonables != "" {
		parts := strings.Split(ajustesRazonables, ", ")
		var translated []string
		locale := salvia_config.Locale["sp"]
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if t, ok := locale[p]; ok {
				translated = append(translated, t)
			} else {
				translated = append(translated, p)
			}
		}
		ajustesRazonables = strings.Join(translated, ", ")
	}
	resp.AjustesRazonables = ajustesRazonables

	// 6. Consentimiento informado (último intento de contacto con consent_given)
	type consentRow struct {
		ConsentGiven *bool  `gorm:"column:consent_given"`
		AttemptAt    string `gorm:"column:attempt_at"`
	}
	var consent consentRow
	s.db.WithContext(ctx).Raw(`
		SELECT consent_given, TO_CHAR(attempt_at, 'YYYY-MM-DD') AS attempt_at
		FROM salvia.contact_attempts
		WHERE psicosocial_id = ? AND consent_given IS NOT NULL AND deleted_at IS NULL
		ORDER BY attempt_at DESC
		LIMIT 1
	`, ps.ID).Scan(&consent)
	if consent.ConsentGiven != nil {
		if *consent.ConsentGiven {
			resp.ConsentStatus = "Aceptado"
		} else {
			resp.ConsentStatus = "No aceptado"
		}
		resp.ConsentDate = consent.AttemptAt
	}

	// 7. Preferencia de horario
	resp.SchedulePreference = ps.SchedulePreference

	// 8. Contactos/Sesiones
	var contacts []models.TeamContact
	s.db.WithContext(ctx).
		Where("psicosocial_id = ? AND deleted_at IS NULL", ps.ID).
		Order("COALESCE(scheduled_date, created_at) ASC").
		Find(&contacts)

	for _, c := range contacts {
		item := PsychosocialContactItem{
			ID:             c.ID,
			IsCompleted:    c.IsCompleted,
			IsPsicoSession: c.IsPsicoSession,
			Status:         c.Status,
			Summary:        c.Summary,
			CreatedAt:      formatContactDateTimeUTC(c.CreatedAt),
		}
		if c.ScheduledDate != nil {
			d := c.ScheduledDate.Format("2006-01-02")
			item.ScheduledDate = &d
		}
		if c.ScheduledTime != nil {
			item.ScheduledTime = c.ScheduledTime
		}
		if c.CompletedAt != nil {
			d := formatContactDateTimeUTC(*c.CompletedAt)
			item.CompletedAt = &d
		}
		resp.Contacts = append(resp.Contacts, item)
	}

	if resp.Contacts == nil {
		resp.Contacts = []PsychosocialContactItem{}
	}

	// 9. Validación pendiente: verificar si hay tarea validar_remision pendiente
	var validationTaskCount int64
	s.db.WithContext(ctx).Raw(`
		SELECT COUNT(*) FROM salvia.case_task
		WHERE psychosocial_support_id = ? AND type = 'validar_remision' AND status = 'ToDo' AND deleted_at IS NULL
	`, id).Scan(&validationTaskCount)
	resp.ValidationPending = validationTaskCount > 0

	// 10. Criterios de remisión (desde Notes del psychosocial_support)
	if ps.Notes != nil && *ps.Notes != "" {
		criteriosRaw := strings.Split(*ps.Notes, ",")
		criterioLabels := map[string]string{
			"criterio_obligatorio":    "Criterio obligatorio",
			"conducta_suicida":        "Conducta suicida",
			"interseccionalidad":      "Interseccionalidad",
			"sin_ruta":                "Sin ruta de atención",
			"condiciones_territoriales": "Condiciones territoriales",
			"sin_acceso_psico":        "Sin acceso a atención psicosocial",
			"naturalizacion_vbg":      "Naturalización de VBG",
		}
		for _, c := range criteriosRaw {
			c = strings.TrimSpace(c)
			if c == "" || c == "criterio_obligatorio" {
				continue
			}
			if label, ok := criterioLabels[c]; ok {
				resp.ValidationCriteria = append(resp.ValidationCriteria, label)
			} else {
				resp.ValidationCriteria = append(resp.ValidationCriteria, c)
			}
		}
	}
	if resp.ValidationCriteria == nil {
		resp.ValidationCriteria = []string{}
	}

	return resp, nil
}

func (s *psychosocialDetailService) CreateContact(ctx context.Context, psicosocialID, contactType, contactDate, contactTime, summary, sessionType, scheduledDate, scheduledTime string) (*models.TeamContact, error) {
	// Obtener la remisión para saber el case_id
	var ps models.PsychosocialSupport
	if err := s.db.WithContext(ctx).Where("id = ?", psicosocialID).First(&ps).Error; err != nil {
		return nil, err
	}

	contact := &models.TeamContact{
		CaseID:        ps.CaseID,
		PsicosocialID: &psicosocialID,
	}

	// Parsear fecha/hora del contacto
	if contactDate != "" && contactTime != "" {
		t, err := time.Parse("2006-01-02 15:04", contactDate+" "+contactTime)
		if err == nil {
			contact.ScheduledDate = &t
			contact.ScheduledTime = &contactTime
		}
	}

	switch contactType {
	case "contacto":
		// Solo contacto — no es sesión psico
		contact.IsPsicoSession = false
		contact.IsCompleted = true
		now := time.Now()
		contact.CompletedAt = &now
		contact.Summary = &summary
		status := "realizado"
		contact.Status = &status

	case "agendar":
		// Agendé sesión — es sesión psico, pendiente
		contact.IsPsicoSession = true
		contact.IsCompleted = false
		status := "agendada"
		contact.Status = &status
		// Parsear fecha/hora de la sesión agendada
		if scheduledDate != "" && scheduledTime != "" {
			t, err := time.Parse("2006-01-02 15:04", scheduledDate+" "+scheduledTime)
			if err == nil {
				contact.ScheduledDate = &t
				contact.ScheduledTime = &scheduledTime
			}
		}

	case "realizar":
		// Realicé sesión — es sesión psico, ya completada
		contact.IsPsicoSession = true
		contact.IsCompleted = true
		now := time.Now()
		contact.CompletedAt = &now
		status := "realizado"
		contact.Status = &status
	}

	if err := s.db.WithContext(ctx).Create(contact).Error; err != nil {
		return nil, err
	}

	// Si es sesión completada, incrementar session_count en psychosocial_support
	if contact.IsPsicoSession && contact.IsCompleted {
		s.db.Exec("UPDATE salvia.psychosocial_support SET session_count = session_count + 1 WHERE id = ?", psicosocialID)
	}

	// Registrar evento en timeline
	now := time.Now()
	var desc string
	var tlType string
	var icon string
	var color string
	switch contactType {
	case "contacto":
		desc = "Contacto registrado: " + summary
		tlType = "Contacto registrado"
		icon = "note-sticky"
		color = models.TimelineColorGray
	case "agendar":
		desc = "Sesión agendada para " + scheduledDate + " " + scheduledTime
		tlType = "Sesión agendada"
		icon = models.TimelineIconSeguimiento
		color = models.TimelineColorGreen
	case "realizar":
		desc = "Sesión realizada"
		tlType = "Sesión realizada"
		icon = "circle-check"
		color = models.TimelineColorGreen
	}
	s.db.WithContext(ctx).Create(&models.CaseTimelineEvent{
		CaseID:                ps.CaseID,
		EventType:             "PSICOSOCIAL_CONTACTO",
		Category:              models.TimelineCategoryPsicosocial,
		Type:                  tlType,
		Icon:                  icon,
		Color:                 color,
		Date:                  now,
		Description:           desc,
		PsychosocialSupportID: psicosocialID,
		CreatedAt:             now,
	})

	return contact, nil
}

func (s *psychosocialDetailService) RescheduleContact(ctx context.Context, contactID, scheduledDate, scheduledTime string) error {
	t, err := time.Parse("2006-01-02 15:04", scheduledDate+" "+scheduledTime)
	if err != nil {
		return err
	}

	// Obtener el contacto para saber case_id y psicosocial_id
	var contact models.TeamContact
	if err := s.db.WithContext(ctx).Where("id = ?", contactID).First(&contact).Error; err != nil {
		return err
	}

	// Actualizar fecha/hora
	if err := s.db.WithContext(ctx).
		Model(&models.TeamContact{}).
		Where("id = ?", contactID).
		Updates(map[string]interface{}{
			"scheduled_date": t,
			"scheduled_time": scheduledTime,
		}).Error; err != nil {
		return err
	}

	// Registrar evento en timeline
	psID := ""
	if contact.PsicosocialID != nil {
		psID = *contact.PsicosocialID
	}
	now := time.Now()
	s.db.WithContext(ctx).Create(&models.CaseTimelineEvent{
		CaseID:                contact.CaseID,
		EventType:             "SESION_REPROGRAMADA",
		Category:              models.TimelineCategoryPsicosocial,
		Type:                  "Sesión reprogramada",
		Icon:                  models.TimelineIconPospuesto,
		Color:                 models.TimelineColorBlue,
		Date:                  now,
		Description:           "Sesión reprogramada para " + scheduledDate + " a las " + scheduledTime,
		PsychosocialSupportID: psID,
		CreatedAt:             now,
	})

	return nil
}

func (s *psychosocialDetailService) CancelContact(ctx context.Context, contactID string) error {
	// Obtener el contacto para saber case_id y psicosocial_id
	var contact models.TeamContact
	if err := s.db.WithContext(ctx).Where("id = ?", contactID).First(&contact).Error; err != nil {
		return err
	}

	cancelado := "cancelada"
	if err := s.db.WithContext(ctx).
		Model(&models.TeamContact{}).
		Where("id = ?", contactID).
		Updates(map[string]interface{}{
			"status": cancelado,
		}).Error; err != nil {
		return err
	}

	// Registrar evento en timeline
	psID := ""
	if contact.PsicosocialID != nil {
		psID = *contact.PsicosocialID
	}
	now := time.Now()
	s.db.WithContext(ctx).Create(&models.CaseTimelineEvent{
		CaseID:                contact.CaseID,
		EventType:             "SESION_CANCELADA",
		Category:              models.TimelineCategoryPsicosocial,
		Type:                  "Sesión cancelada",
		Icon:                  "circle-xmark",
		Color:                 models.TimelineColorRed,
		Date:                  now,
		Description:           "Sesión cancelada",
		PsychosocialSupportID: psID,
		CreatedAt:             now,
	})

	return nil
}

// LoadSession implementa el evento E-01 (carga inicial de la pantalla de sesión):
//  1. Valida que la remisión exista y que el agente tenga acceso (profesional directo
//     o miembro de la dupla asignada).
//  2. Selecciona el formulario psicosocial según el estado del proceso (PASO 7).
//  3. Resuelve (o crea) el team_contact pendiente y fija su form_id la primera vez,
//     reutilizándolo en cargas posteriores para que la sesión no cambie de formulario
//     a medio camino. Un team_contact nuevo solo recibe professional_id O dupla_id
//     (nunca ambos), reflejando el modo de asignación de la remisión padre.
//  4. Resuelve (o crea) el form_submission asociado.
//  5. Carga la información resumida de la víctima.
func (s *psychosocialDetailService) LoadSession(ctx context.Context, psicosocialID, agentID, contactID string) (*LoadPsicosocialSessionResult, error) {
	var ps models.PsychosocialSupport
	if err := s.db.WithContext(ctx).Where("id = ?", psicosocialID).First(&ps).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPsicosocialSessionNotFound
		}
		return nil, err
	}

	if err := s.checkSessionAccess(ctx, ps, agentID); err != nil {
		return nil, err
	}

	var (
		tc           models.TeamContact
		isNewContact bool
		viewByID     = contactID != ""
		err          error
	)

	if viewByID {
		tc, err = s.findSessionContactByID(ctx, psicosocialID, contactID)
		if err != nil {
			return nil, err
		}
		if tc.FormSubmissionID == nil || *tc.FormSubmissionID == "" {
			return nil, ErrPsicosocialContactNoSubmission
		}
	} else {
		tc, isNewContact, err = s.findOrInitPendingContact(ctx, psicosocialID)
		if err != nil {
			return nil, fmt.Errorf("loadSession: resolver team_contact: %w", err)
		}
	}

	formKey, formID := "", ""
	if tc.FormID != nil && *tc.FormID != "" {
		// El formulario ya quedó fijado en una carga anterior — se reutiliza siempre.
		formID = *tc.FormID
		formKey = string(constants.PsicosocialFormKeyByID[formID])
	} else {
		key, id := selectPsicosocialForm(ps)
		formKey, formID = string(key), id
		if !viewByID {
			tc.FormID = &formID
			sessionType := formKey
			tc.SessionType = &sessionType
		}
	}

	if !viewByID {
		if isNewContact {
			if err := s.applyContactAssignment(&tc, ps, agentID); err != nil {
				return nil, fmt.Errorf("loadSession: asignar team_contact: %w", err)
			}
			if err := s.db.WithContext(ctx).Create(&tc).Error; err != nil {
				return nil, fmt.Errorf("loadSession: crear team_contact: %w", err)
			}
		} else {
			if err := s.db.WithContext(ctx).Model(&models.TeamContact{}).Where("id = ?", tc.ID).
				Updates(map[string]interface{}{"form_id": tc.FormID, "session_type": tc.SessionType}).Error; err != nil {
				return nil, fmt.Errorf("loadSession: actualizar team_contact: %w", err)
			}
		}

		if tc.FormSubmissionID == nil || *tc.FormSubmissionID == "" {
			fs := &models.FormSubmission{FormID: formID}
			if err := s.db.WithContext(ctx).Create(fs).Error; err != nil {
				return nil, fmt.Errorf("loadSession: crear form_submission: %w", err)
			}
			if err := s.db.WithContext(ctx).Model(&models.TeamContact{}).Where("id = ?", tc.ID).
				Update("form_submission_id", fs.ID).Error; err != nil {
				return nil, fmt.Errorf("loadSession: asociar form_submission a team_contact: %w", err)
			}
			tc.FormSubmissionID = &fs.ID
		}
	}

	victimInfo, err := s.loadSessionVictimInfo(ctx, ps.CaseID)
	if err != nil {
		log.Printf("[SVC] LoadSession → advertencia: no se pudo cargar victimInfo de caseID=%s: %v", ps.CaseID, err)
	}

	currentBarriers := s.loadActivePsicosocialBarriers(ctx, ps.ID)
	canEdit := !tc.IsCompleted && ps.Status != "cerrado"

	return &LoadPsicosocialSessionResult{
		FormID:        formID,
		FormType:      formKey,
		SubmissionID:  *tc.FormSubmissionID,
		TeamContactID: tc.ID,
		IsCompleted:   tc.IsCompleted,
		CanEdit:       canEdit,
		VictimInfo:    victimInfo,
		PsicosocialState: PsicosocialStateInfo{
			YaHizoPrimerContacto:  ps.YaHizoPrimerContacto,
			YaHizoPrimeraAtencion: ps.YaHizoPrimeraAtencion,
			SessionCount:          ps.SessionCount,
			Status:                ps.Status,
		},
		FormState: map[string]interface{}{
			"currentBarriers": currentBarriers,
		},
	}, nil
}

// loadActivePsicosocialBarriers resuelve las "barreras activas" de una remisión psicosocial
// para poblar formState.currentBarriers (consumido por el repeater "Seguimiento a Barreras" vía
// stateItems, igual patrón que hacer_seguimiento.html/followUpV2Service). A diferencia de
// hacer_seguimiento (que filtra por case_id + guarda un CSV de IDs en el propio follow_up), aquí
// se filtra por team_contact_id: se buscan todos los team_contact de este psychosocial_support y
// luego las barrier_v2 (no MANAGED) cuyo team_contact_id esté en ese conjunto — decisión de los
// líderes (Jul 2026) para no depender de un campo adicional en psychosocial_support.
// Reutiliza ActiveBarrierInfo/buildBarrierName definidos en followup_v2_service.go (mismo paquete).
func (s *psychosocialDetailService) loadActivePsicosocialBarriers(ctx context.Context, psicosocialID string) []ActiveBarrierInfo {
	var contacts []models.TeamContact
	if err := s.db.WithContext(ctx).
		Where("psicosocial_id = ? AND deleted_at IS NULL", psicosocialID).
		Find(&contacts).Error; err != nil {
		log.Printf("[SVC] loadActivePsicosocialBarriers → error leyendo team_contact de psicosocialId=%s: %v", psicosocialID, err)
		return []ActiveBarrierInfo{}
	}
	if len(contacts) == 0 {
		return []ActiveBarrierInfo{}
	}
	contactIDs := make([]string, len(contacts))
	for i, c := range contacts {
		contactIDs[i] = c.ID
	}

	var barriers []models.BarrierV2
	if err := s.db.WithContext(ctx).
		Where("team_contact_id IN ? AND status != ? AND deleted_at IS NULL", contactIDs, models.BarrierV2StatusManaged).
		Order("created_at ASC").
		Find(&barriers).Error; err != nil {
		log.Printf("[SVC] loadActivePsicosocialBarriers → error leyendo barrier_v2 de psicosocialId=%s: %v", psicosocialID, err)
		return []ActiveBarrierInfo{}
	}

	result := make([]ActiveBarrierInfo, len(barriers))
	for i, b := range barriers {
		result[i] = ActiveBarrierInfo{ID: b.ID, BarrierName: buildBarrierName(b)}
	}
	return result
}

// checkSessionAccess valida que agentID sea el profesional directo asignado a la
// remisión, o que pertenezca a la dupla asignada (psicóloga o trabajador social).
func (s *psychosocialDetailService) checkSessionAccess(ctx context.Context, ps models.PsychosocialSupport, agentID string) error {
	if ps.ProfessionalID != nil && *ps.ProfessionalID == agentID {
		return nil
	}
	if ps.DuplaID != nil && *ps.DuplaID != "" {
		var dupla models.Dupla
		if err := s.db.WithContext(ctx).Where("id = ?", *ps.DuplaID).First(&dupla).Error; err == nil {
			if dupla.PsychologistID == agentID || dupla.SocialWorkerID == agentID {
				return nil
			}
		}
	}
	return ErrPsicosocialSessionNotAssigned
}

// findSessionContactByID carga un team_contact por id que pertenece a la remisión
// (flujo "Ver sesión"). No crea registros.
func (s *psychosocialDetailService) findSessionContactByID(ctx context.Context, psicosocialID, contactID string) (models.TeamContact, error) {
	var tc models.TeamContact
	err := s.db.WithContext(ctx).
		Where("id = ? AND psicosocial_id = ? AND deleted_at IS NULL", contactID, psicosocialID).
		First(&tc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.TeamContact{}, ErrPsicosocialContactNotFound
		}
		return models.TeamContact{}, fmt.Errorf("loadSession: cargar team_contact: %w", err)
	}
	return tc, nil
}

// findOrInitPendingContact busca un team_contact pendiente (no completado, sesión
// psico) para la remisión. Si no existe, retorna un struct vacío listo para crear
// (isNewContact = true); el llamador decide asignación y formulario antes de crearlo.
func (s *psychosocialDetailService) findOrInitPendingContact(ctx context.Context, psicosocialID string) (models.TeamContact, bool, error) {
	var tc models.TeamContact
	err := s.db.WithContext(ctx).
		Where("psicosocial_id = ? AND is_completed = false AND is_psico_session = true AND deleted_at IS NULL", psicosocialID).
		Order("created_at DESC").
		First(&tc).Error

	if err == nil {
		return tc, false, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.TeamContact{
			CaseID:         "", // se completa en LoadSession antes de crear
			PsicosocialID:  &psicosocialID,
			IsCompleted:    false,
			IsPsicoSession: true,
		}, true, nil
	}
	return models.TeamContact{}, false, err
}

// applyContactAssignment fija case_id y professional_id/dupla_id en un team_contact
// nuevo, reflejando exactamente el modo de asignación de la remisión padre: nunca
// se guardan ambos campos a la vez.
func (s *psychosocialDetailService) applyContactAssignment(tc *models.TeamContact, ps models.PsychosocialSupport, agentID string) error {
	tc.CaseID = ps.CaseID
	switch {
	case ps.DuplaID != nil && *ps.DuplaID != "":
		tc.DuplaID = ps.DuplaID
		tc.ProfessionalID = nil
	case ps.ProfessionalID != nil && *ps.ProfessionalID != "":
		tc.ProfessionalID = ps.ProfessionalID
		tc.DuplaID = nil
	default:
		// La remisión no tiene asignación explícita todavía (caso raro) — se asigna
		// directamente al agente que abrió la sesión.
		tc.ProfessionalID = &agentID
		tc.DuplaID = nil
	}
	return nil
}

// selectPsicosocialForm implementa el PASO 7 del flujo E-01: elige, en orden, el
// primer formulario cuyo criterio se cumpla según el estado acumulado del proceso.
func selectPsicosocialForm(ps models.PsychosocialSupport) (constants.PsicosocialFormKey, string) {
	switch {
	case !ps.YaHizoPrimerContacto:
		return constants.PsicosocialFormPrimerContacto, constants.FormIDPrimerContacto
	case !ps.YaHizoPrimeraAtencion:
		return constants.PsicosocialFormPrimeraAtencion, constants.FormIDPrimeraAtencion
	case ps.SessionCount < 3:
		return constants.PsicosocialFormAtencionPsicosocial, constants.FormIDAtencionPsicosocial
	default:
		return constants.PsicosocialFormCierre, constants.FormIDCierre
	}
}

// loadSessionVictimInfo carga los datos resumidos de la víctima para la tarjeta de
// la pantalla de sesión psicosocial (mismos campos que LoadVictimInfoByCaseID).
func (s *psychosocialDetailService) loadSessionVictimInfo(ctx context.Context, caseID string) (*PsicosocialSessionVictimInfo, error) {
	type victimRow struct {
		Names             string `gorm:"column:names"`
		LastNames         string `gorm:"column:last_names"`
		Phone             string `gorm:"column:phone"`
		GenderIdentity    string `gorm:"column:gender_identity"`
		SexualOrientation string `gorm:"column:sexual_orientation"`
		ContactPhone      string `gorm:"column:contact_phone"`
		Age               *int64 `gorm:"column:age"`
		TownName          string `gorm:"column:town_name"`
		RiskLevel         int    `gorm:"column:risk_level"`
	}
	var row victimRow
	err := s.db.WithContext(ctx).Raw(`
		SELECT
			COALESCE(vc.victim_case_victim_names, '')                           AS names,
			COALESCE(vc.victim_case_victim_last_names, '')                      AS last_names,
			COALESCE(f2.victim_case_form2_victim_phone::text, '')               AS phone,
			COALESCE(gi.victim_case_form2_enums_name, '')                       AS gender_identity,
			COALESCE(so.victim_case_form2_enums_name, '')                       AS sexual_orientation,
			COALESCE(f2.victim_case_form2_support_contact_phone::text, '')      AS contact_phone,
			EXTRACT(YEAR FROM AGE(NOW(), f2.victim_case_form2_birth_date))::int AS age,
			COALESCE(t.town_name, '')                                           AS town_name,
			COALESCE(f2.victim_case_form2_risk_level, 0)                        AS risk_level
		FROM salvia.victim_case vc
		LEFT JOIN salvia.victim_case_form2 f2
			ON f2.victim_case_form2_victim_case = vc.victim_case_id
		LEFT JOIN security.town t
			ON t.town_code = vc.victim_case_victim_town_code
		LEFT JOIN salvia.victim_case_form2_enums gi
			ON gi.victim_case_form2_enums_id = f2.victim_case_form2_gender_identity
		LEFT JOIN salvia.victim_case_form2_enums so
			ON so.victim_case_form2_enums_id = f2.victim_case_form2_sexual_orientation
		WHERE vc.victim_case_i_code = ?
		LIMIT 1
	`, caseID).Scan(&row).Error
	if err != nil {
		return nil, err
	}

	locale := salvia_config.Locale["sp"]
	return &PsicosocialSessionVictimInfo{
		Names:             row.Names,
		LastNames:         row.LastNames,
		Phone:             row.Phone,
		GenderIdentity:    locale[row.GenderIdentity],
		SexualOrientation: locale[row.SexualOrientation],
		ContactPhone:      row.ContactPhone,
		Age:               row.Age,
		TownName:          row.TownName,
		RiskLevel:         row.RiskLevel,
		CaseICode:         caseID,
	}, nil
}

func (s *psychosocialDetailService) UpdateSchedulePreference(ctx context.Context, psicosocialID, preference string) error {
	return s.db.WithContext(ctx).
		Model(&models.PsychosocialSupport{}).
		Where("id = ?", psicosocialID).
		Update("schedule_preference", preference).Error
}

// ValidateRemision procesa la validación de una remisión psicosocial.
// Si isValid=true: completa la tarea y deja la remisión habilitada para gestión.
// Si isValid=false: completa la tarea y pone la remisión en estado "en_devolucion".
func (s *psychosocialDetailService) ValidateRemision(ctx context.Context, id string, isValid bool, reason string, actorID string) error {
	// 1. Buscar la remisión
	var ps models.PsychosocialSupport
	if err := s.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&ps).Error; err != nil {
		return err
	}

	// 2. Completar la tarea de validación
	s.db.WithContext(ctx).Exec(`
		UPDATE salvia.case_task SET status = 'Done', updated_at = NOW()
		WHERE psychosocial_support_id = ? AND type = 'validar_remision' AND status = 'ToDo' AND deleted_at IS NULL
	`, id)

	// 3. Resolver nombre del actor
	var actorName string
	s.db.WithContext(ctx).Raw(`
		SELECT gup.general_user_profile_names || ' ' || gup.general_user_profile_last_names
		FROM security.general_user gu
		JOIN security.general_user_profile gup ON gup.general_user_profile_id = gu.general_user_general_user_profile
		WHERE gu.general_user_i_code = ?
	`, actorID).Scan(&actorName)
	if actorName == "" {
		actorName = actorID
	}

	now := time.Now()

	if isValid {
		// 4a. Remisión validada — no cambia de status (queda "abierto" para ser gestionada)
		s.db.WithContext(ctx).Exec(`INSERT INTO salvia.case_timeline_event
			(case_id, category, type, icon, color, description, event_user_id, actor_name, psychosocial_support_id, date, created_at)
			VALUES (?, 'Psicosocial', 'Remisión Validada', 'circle-check', '#10b981', ?, ?, ?, ?, ?, ?)`,
			ps.CaseID, "Remisión validada por "+actorName, actorID, actorName, id, now, now)
	} else {
		// 4b. Remisión devuelta — cambiar status a en_devolucion
		s.db.WithContext(ctx).Exec(`UPDATE salvia.psychosocial_support SET status = 'en_devolucion', updated_at = NOW() WHERE id = ?`, id)

		descripcion := "Remisión devuelta — Motivo: " + reason
		s.db.WithContext(ctx).Exec(`INSERT INTO salvia.case_timeline_event
			(case_id, category, type, icon, color, description, event_user_id, actor_name, psychosocial_support_id, date, created_at)
			VALUES (?, 'Psicosocial', 'Remisión Devuelta', 'arrow-rotate-left', '#dc2626', ?, ?, ?, ?, ?, ?)`,
			ps.CaseID, descripcion, actorID, actorName, id, now, now)

		// 5. Tarea "Justificar Remisión" para el agente dueño del caso — solo
		// cuando la remisión se devuelve (no aplica si se valida como correcta,
		// no hay nada que justificar en ese caso).
		var caseAgentID string
		s.db.WithContext(ctx).Raw(`SELECT COALESCE(agent_id, '') FROM salvia.victim_case WHERE victim_case_i_code = ?`, ps.CaseID).Scan(&caseAgentID)
		if caseAgentID == "" {
			log.Printf("[ValidateRemision] WARN: caso %s sin agente asignado — no se creó tarea Justificar Remisión", ps.CaseID)
		} else {
			s.db.WithContext(ctx).Exec(`
				INSERT INTO salvia.case_task
					(category, type, description, assigned_user_id, status, case_id, follow_up_id, psychosocial_support_id, created_at, updated_at)
				VALUES ('Psicosocial', 'justificar_remision', ?, ?, 'ToDo', ?, ?, ?, NOW(), NOW())
			`, "Justificar remisión devuelta — Motivo: "+reason, caseAgentID, ps.CaseID, ps.FollowUpID, id)
		}
	}

	return nil
}

// CheckSessionAvailability GET /psychosocial-support/:id/availability
// Valida si date (YYYY-MM-DD) + time (HH:MM) está libre en ventana de 2 horas
// para el profesional (individual) o los miembros de la dupla.
func (s *psychosocialDetailService) CheckSessionAvailability(ctx context.Context, psicosocialID, date, timeStr, mode string) (bool, string, error) {
	var ps models.PsychosocialSupport
	if err := s.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", psicosocialID).First(&ps).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, "", ErrPsicosocialSessionNotFound
		}
		return false, "", err
	}

	fecha, err := parseAgendaDate(date)
	if err != nil {
		return false, "Fecha inválida; use formato YYYY-MM-DD", nil
	}
	hora := normalizeScheduledTime(timeStr)
	if _, ok := parseTimeToMinutes(hora); !ok {
		return false, "Hora inválida; use formato HH:MM", nil
	}

	day := fecha.Format("2006-01-02")
	var contacts []models.TeamContact
	if err := s.db.WithContext(ctx).
		Where("scheduled_date IS NOT NULL AND (scheduled_date AT TIME ZONE 'UTC')::date = ?::date AND is_psico_session = true AND deleted_at IS NULL", day).
		Find(&contacts).Error; err != nil {
		return false, "", err
	}

	mode = strings.ToLower(strings.TrimSpace(mode))
	if mode == "" {
		if ps.DuplaID != nil && *ps.DuplaID != "" {
			mode = "dupla"
		} else {
			mode = "individual"
		}
	}

	var psychID, swID string
	if mode == "dupla" && ps.DuplaID != nil && *ps.DuplaID != "" {
		var dupla models.Dupla
		if err := s.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", *ps.DuplaID).First(&dupla).Error; err == nil {
			psychID = dupla.PsychologistID
			swID = dupla.SocialWorkerID
		}
	}

	ok, msg := evaluatePsicosocialAvailability(contacts, &ps, hora, mode, psychID, swID)
	return ok, msg, nil
}

// formatContactDateTimeUTC serializa un instante como RFC3339 en UTC para que el
// frontend lo muestre en America/Bogota sin ambigüedad de zona.
func formatContactDateTimeUTC(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}
