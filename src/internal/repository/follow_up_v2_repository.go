package repository

import (
	"bitsflow/internal/models"
	"context"

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
}

type FollowUpV2Repository interface {
	Repository[models.FollowUpV2]
	FindByCaseID(ctx context.Context, caseID string) ([]models.FollowUpV2, error)
	FindPending(ctx context.Context) ([]models.FollowUpV2, error)
	LoadVictimInfoByCaseID(ctx context.Context, caseID string) (*VictimCaseInfo, error)
	UpdateFormSubmissionID(ctx context.Context, id string, fsID string) error
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
	var info VictimCaseInfo
	sql := `
		SELECT
			COALESCE(vc.victim_case_victim_names, '')      AS victim_names,
			COALESCE(vc.victim_case_victim_last_names, '') AS victim_last_names,
			COALESCE(t.town_name, '')                      AS town_name,
			COALESCE(f1.victim_case_form1_victim_phone, '') AS victim_phone,
			COALESCE(f1.victim_case_form1_victim_gender_identity, '') AS victim_gender_identity,
			COALESCE(f1.victim_case_form1_victim_sexual_orientation, '') AS victim_sexual_orientation,
			COALESCE(f1.victim_case_form1_victim_contact_phone, '') AS victim_contact_phone,
			f1.victim_case_form1_age AS age
		FROM salvia.victim_case vc
		LEFT JOIN salvia.victim_case_form1 f1 ON f1.victim_case_form1_victim_case = vc.victim_case_id
		LEFT JOIN security.town t ON t.town_code = vc.victim_case_victim_town_code
		WHERE vc.victim_case_id::text = ?
		LIMIT 1`
	return &info, r.db.WithContext(ctx).Raw(sql, caseID).Scan(&info).Error
}

func (r *followUpV2Repository) UpdateFormSubmissionID(ctx context.Context, id string, fsID string) error {
	return r.db.WithContext(ctx).Model(&models.FollowUpV2{}).
		Where("id = ?", id).
		Update("form_submission_id", fsID).Error
}
