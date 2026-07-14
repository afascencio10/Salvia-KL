// Package service — psychosocial_detail_service.go
// Lógica de negocio para la pantalla de detalle de remisión psicosocial.
package service

import (
	"bitsflow/internal/models"
	"context"
	"strings"
	"time"

	"gorm.io/gorm"
)

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
	CreateContact(ctx context.Context, psicosocialID, contactType, contactDate, contactTime, summary, sessionType, scheduledDate, scheduledTime string) (*models.TeamContact, error)
	RescheduleContact(ctx context.Context, contactID, scheduledDate, scheduledTime string) error
	CancelContact(ctx context.Context, contactID string) error
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

	// 5. Contactos/Sesiones
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
			CreatedAt:      c.CreatedAt.Format("2006-01-02 15:04"),
		}
		if c.ScheduledDate != nil {
			d := c.ScheduledDate.Format("2006-01-02")
			item.ScheduledDate = &d
		}
		if c.ScheduledTime != nil {
			item.ScheduledTime = c.ScheduledTime
		}
		if c.CompletedAt != nil {
			d := c.CompletedAt.Format("2006-01-02 15:04")
			item.CompletedAt = &d
		}
		resp.Contacts = append(resp.Contacts, item)
	}

	if resp.Contacts == nil {
		resp.Contacts = []PsychosocialContactItem{}
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
