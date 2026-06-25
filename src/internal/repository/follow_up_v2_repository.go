package repository

import (
	"bitsflow/internal/models"
	"context"
	"log"
	"time"

	"gorm.io/gorm"
)

// VictimCaseInfo contiene la información resumida del caso para la pantalla hacer-seguimiento.
type VictimCaseInfo struct {
	Names             string `json:"Names"`
	LastNames         string `json:"LastNames"`
	TownName          string `json:"TownName"`
	Phone             string `json:"Phone"`
	GenderIdentity    string `json:"GenderIdentity"`
	SexualOrientation string `json:"SexualOrientation"`
	ContactPhone      string `json:"ContactPhone"`
	Age               *int64 `json:"Age"`
	RiskLevel         int    `json:"riskLevel"`
}

type FollowUpV2Repository interface {
	Repository[models.FollowUpV2]
	FindByCaseID(ctx context.Context, caseID string) ([]models.FollowUpV2, error)
	FindPending(ctx context.Context) ([]models.FollowUpV2, error)
	FindPendingByCaseID(ctx context.Context, caseID string) ([]models.FollowUpV2, error)
	DeletePendingByCaseID(ctx context.Context, caseID string) error
	UpdateAgentForPendingByCaseID(ctx context.Context, caseID string, agentID string) error
	LoadVictimInfoByCaseID(ctx context.Context, caseID string) (*VictimCaseInfo, error)
	UpdateFormSubmissionID(ctx context.Context, id string, fsID string) error
	ExistsForCaseOnDate(ctx context.Context, caseID string, date time.Time) (bool, error)
}

type followUpV2Repository struct {
	repository[models.FollowUpV2]
	db *gorm.DB
}

func NewFollowUpV2Repository(db *gorm.DB) FollowUpV2Repository {
	return &followUpV2Repository{
		repository: repository[models.FollowUpV2]{db: db},
		db:         db,
	}
}

func (r *followUpV2Repository) FindByCaseID(ctx context.Context, caseID string) ([]models.FollowUpV2, error) {
	var items []models.FollowUpV2
	return items, r.db.WithContext(ctx).Where("case_id = ?", caseID).Find(&items).Error
}

func (r *followUpV2Repository) FindPending(ctx context.Context) ([]models.FollowUpV2, error) {
	var items []models.FollowUpV2
	return items, r.db.WithContext(ctx).Where("status = ?", models.FollowUpStatusPendiente).Find(&items).Error
}

func (r *followUpV2Repository) LoadVictimInfoByCaseID(ctx context.Context, caseID string) (*VictimCaseInfo, error) {
	log.Printf("[REPO] LoadVictimInfoByCaseID → caseID=%s", caseID)
	var info VictimCaseInfo
	sql := `
		SELECT
			COALESCE(vc.victim_case_victim_names, '')                           AS names,
			COALESCE(vc.victim_case_victim_last_names, '')                      AS last_names,
			COALESCE(t.town_name, '')                                           AS town_name,
			COALESCE(f2.victim_case_form2_victim_phone::text, '')               AS phone,
			COALESCE(gi.victim_case_form2_enums_name, '')                       AS gender_identity,
			COALESCE(so.victim_case_form2_enums_name, '')                       AS sexual_orientation,
			COALESCE(f2.victim_case_form2_support_contact_phone::text, '')      AS contact_phone,
			EXTRACT(YEAR FROM AGE(NOW(), f2.victim_case_form2_birth_date))::int AS age,
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
		LIMIT 1`
	err := r.db.WithContext(ctx).Raw(sql, caseID).Scan(&info).Error
	log.Printf("[REPO] LoadVictimInfoByCaseID → resultado: err=%v info=%+v", err, info)
	return &info, err
}

func (r *followUpV2Repository) FindPendingByCaseID(ctx context.Context, caseID string) ([]models.FollowUpV2, error) {
	var items []models.FollowUpV2
	err := r.db.WithContext(ctx).
		Where("case_id = ? AND status = ?", caseID, models.FollowUpStatusPendiente).
		Find(&items).Error
	return items, err
}

func (r *followUpV2Repository) DeletePendingByCaseID(ctx context.Context, caseID string) error {
	return r.db.WithContext(ctx).
		Where("case_id = ? AND status = ?", caseID, models.FollowUpStatusPendiente).
		Delete(&models.FollowUpV2{}).Error
}

func (r *followUpV2Repository) UpdateAgentForPendingByCaseID(ctx context.Context, caseID string, agentID string) error {
	return r.db.WithContext(ctx).
		Model(&models.FollowUpV2{}).
		Where("case_id = ? AND status = ?", caseID, models.FollowUpStatusPendiente).
		Update("agent_id", agentID).Error
}

func (r *followUpV2Repository) UpdateFormSubmissionID(ctx context.Context, id string, fsID string) error {
	return r.db.WithContext(ctx).Model(&models.FollowUpV2{}).
		Where("id = ?", id).
		Update("form_submission_id", fsID).Error
}

// ExistsForCaseOnDate devuelve true si ya existe un seguimiento (no eliminado) para el
// caso en el mismo día calendario que date, independientemente de la hora.
func (r *followUpV2Repository) ExistsForCaseOnDate(ctx context.Context, caseID string, date time.Time) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.FollowUpV2{}).
		Where("case_id = ? AND DATE(scheduled_date) = DATE(?)", caseID, date).
		Count(&count).Error
	return count > 0, err
}

