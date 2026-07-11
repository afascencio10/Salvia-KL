// Package service — psychosocial_detail_service.go
// Lógica de negocio para la pantalla de detalle de remisión psicosocial.
package service

import (
	"bitsflow/internal/models"
	"context"
	"strings"

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
