package repository

import (
	internaldb "bitsflow/internal/db"
	"bitsflow/internal/models"
	"context"
	"time"

	"gorm.io/gorm"
)

// PsychosocialSupportReassignRow datos mínimos de una remisión para reasignación (RRM-05).
type PsychosocialSupportReassignRow struct {
	ID     string `gorm:"column:id"`
	Status string `gorm:"column:status"`
}

// PsychosocialReassignRepository lecturas y escrituras para modal reasignar remisiones.
type PsychosocialReassignRepository interface {
	ListActivePsTsProfessionals(ctx context.Context) ([]models.PsTsProfessionalRow, error)
	ListActiveDuplasEnriched(ctx context.Context) ([]models.DuplaReassignRow, error)
	FindPsychosocialForReassign(ctx context.Context, remisionID string) (*PsychosocialSupportReassignRow, error)
	FindProfessionalForReassign(ctx context.Context, icode string) (*models.PsTsProfessionalRow, error)
	FindDuplaForReassign(ctx context.Context, duplaID string) (*models.DuplaReassignRow, error)
	UpdatePsychosocialProfessional(ctx context.Context, remisionID, professionalID string) error
	UpdatePsychosocialDupla(ctx context.Context, remisionID, duplaID string) error
	UpdateTeamContactsProfessional(ctx context.Context, remisionID, professionalID string) (int64, error)
	UpdateTeamContactsDupla(ctx context.Context, remisionID, duplaID string) (int64, error)
}

type psychosocialReassignRepository struct {
	db *gorm.DB
}

func NewPsychosocialReassignRepository(db *gorm.DB) PsychosocialReassignRepository {
	return &psychosocialReassignRepository{db: db}
}

func (r *psychosocialReassignRepository) ListActivePsTsProfessionals(ctx context.Context) ([]models.PsTsProfessionalRow, error) {
	const sql = `
SELECT
	gu.general_user_i_code AS icode,
	TRIM(COALESCE(gup.general_user_profile_names, '') || ' ' || COALESCE(gup.general_user_profile_last_names, '')) AS full_name,
	r.role_code AS role
FROM security.general_user gu
JOIN security.general_user_profile gup
     ON gup.general_user_profile_id = gu.general_user_general_user_profile
JOIN security.rel_role_general_user rrgu
     ON rrgu.general_user_id = gu.general_user_id
JOIN security.role r
     ON r.role_id = rrgu.role_id
WHERE gu.general_user_status = 'e'
  AND r.role_code IN ('ps', 'ts')
ORDER BY r.role_code ASC, full_name ASC`

	var rows []models.PsTsProfessionalRow
	err := internaldb.WithRetry(func() error {
		return r.db.WithContext(ctx).Raw(sql).Scan(&rows).Error
	})
	if err != nil {
		return nil, err
	}
	if rows == nil {
		rows = []models.PsTsProfessionalRow{}
	}
	return rows, nil
}

func (r *psychosocialReassignRepository) ListActiveDuplasEnriched(ctx context.Context) ([]models.DuplaReassignRow, error) {
	const sql = `
SELECT
	d.id,
	d.name,
	TRIM(COALESCE(ps_gup.general_user_profile_names, '') || ' ' || COALESCE(ps_gup.general_user_profile_last_names, '')) AS psychologist_name,
	TRIM(COALESCE(ts_gup.general_user_profile_names, '') || ' ' || COALESCE(ts_gup.general_user_profile_last_names, '')) AS social_worker_name
FROM salvia.dupla d
JOIN security.general_user ps_gu
     ON ps_gu.general_user_i_code = BTRIM(d.psychologist_id::text)
JOIN security.general_user_profile ps_gup
     ON ps_gup.general_user_profile_id = ps_gu.general_user_general_user_profile
JOIN security.general_user ts_gu
     ON ts_gu.general_user_i_code = BTRIM(d.social_worker_id::text)
JOIN security.general_user_profile ts_gup
     ON ts_gup.general_user_profile_id = ts_gu.general_user_general_user_profile
WHERE d.deleted_at IS NULL
  AND ps_gu.general_user_status = 'e'
  AND ts_gu.general_user_status = 'e'
ORDER BY d.name ASC`

	var rows []models.DuplaReassignRow
	err := internaldb.WithRetry(func() error {
		return r.db.WithContext(ctx).Raw(sql).Scan(&rows).Error
	})
	if err != nil {
		return nil, err
	}
	if rows == nil {
		rows = []models.DuplaReassignRow{}
	}
	return rows, nil
}

func (r *psychosocialReassignRepository) FindPsychosocialForReassign(ctx context.Context, remisionID string) (*PsychosocialSupportReassignRow, error) {
	var row PsychosocialSupportReassignRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT id, COALESCE(status, '') AS status
		FROM salvia.psychosocial_support
		WHERE id = ?
		  AND deleted_at IS NULL
		LIMIT 1
	`, remisionID).Scan(&row).Error
	if err != nil {
		return nil, err
	}
	if row.ID == "" {
		return nil, gorm.ErrRecordNotFound
	}
	return &row, nil
}

func (r *psychosocialReassignRepository) FindProfessionalForReassign(ctx context.Context, icode string) (*models.PsTsProfessionalRow, error) {
	var row models.PsTsProfessionalRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			gu.general_user_i_code AS icode,
			TRIM(COALESCE(gup.general_user_profile_names, '') || ' ' || COALESCE(gup.general_user_profile_last_names, '')) AS full_name,
			r.role_code AS role
		FROM security.general_user gu
		JOIN security.general_user_profile gup
		     ON gup.general_user_profile_id = gu.general_user_general_user_profile
		JOIN security.rel_role_general_user rrgu
		     ON rrgu.general_user_id = gu.general_user_id
		JOIN security.role r
		     ON r.role_id = rrgu.role_id
		WHERE gu.general_user_i_code = ?
		  AND gu.general_user_status = 'e'
		  AND r.role_code IN ('ps', 'ts')
		LIMIT 1
	`, icode).Scan(&row).Error
	if err != nil {
		return nil, err
	}
	if row.ICode == "" {
		return nil, gorm.ErrRecordNotFound
	}
	return &row, nil
}

func (r *psychosocialReassignRepository) FindDuplaForReassign(ctx context.Context, duplaID string) (*models.DuplaReassignRow, error) {
	var row models.DuplaReassignRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			d.id,
			d.name,
			TRIM(COALESCE(ps_gup.general_user_profile_names, '') || ' ' || COALESCE(ps_gup.general_user_profile_last_names, '')) AS psychologist_name,
			TRIM(COALESCE(ts_gup.general_user_profile_names, '') || ' ' || COALESCE(ts_gup.general_user_profile_last_names, '')) AS social_worker_name
		FROM salvia.dupla d
		JOIN security.general_user ps_gu
		     ON ps_gu.general_user_i_code = BTRIM(d.psychologist_id::text)
		JOIN security.general_user_profile ps_gup
		     ON ps_gup.general_user_profile_id = ps_gu.general_user_general_user_profile
		JOIN security.general_user ts_gu
		     ON ts_gu.general_user_i_code = BTRIM(d.social_worker_id::text)
		JOIN security.general_user_profile ts_gup
		     ON ts_gup.general_user_profile_id = ts_gu.general_user_general_user_profile
		WHERE d.id = ?
		  AND d.deleted_at IS NULL
		  AND ps_gu.general_user_status = 'e'
		  AND ts_gu.general_user_status = 'e'
		LIMIT 1
	`, duplaID).Scan(&row).Error
	if err != nil {
		return nil, err
	}
	if row.ID == "" {
		return nil, gorm.ErrRecordNotFound
	}
	return &row, nil
}

func (r *psychosocialReassignRepository) UpdatePsychosocialProfessional(ctx context.Context, remisionID, professionalID string) error {
	return r.db.WithContext(ctx).Exec(`
		UPDATE salvia.psychosocial_support
		SET professional_id = ?,
		    dupla_id = NULL,
		    updated_at = ?
		WHERE id = ?
		  AND deleted_at IS NULL
		  AND status <> ?
	`, professionalID, time.Now(), remisionID, models.PsychosocialSupportStatusCerrado).Error
}

func (r *psychosocialReassignRepository) UpdatePsychosocialDupla(ctx context.Context, remisionID, duplaID string) error {
	return r.db.WithContext(ctx).Exec(`
		UPDATE salvia.psychosocial_support
		SET dupla_id = ?,
		    professional_id = NULL,
		    updated_at = ?
		WHERE id = ?
		  AND deleted_at IS NULL
		  AND status <> ?
	`, duplaID, time.Now(), remisionID, models.PsychosocialSupportStatusCerrado).Error
}

func (r *psychosocialReassignRepository) UpdateTeamContactsProfessional(ctx context.Context, remisionID, professionalID string) (int64, error) {
	result := r.db.WithContext(ctx).Exec(`
		UPDATE salvia.team_contact
		SET professional_id = ?,
		    dupla_id = NULL,
		    updated_at = ?
		WHERE BTRIM(psicosocial_id::text) = BTRIM(?)
		  AND is_completed = false
		  AND deleted_at IS NULL
	`, professionalID, time.Now(), remisionID)
	return result.RowsAffected, result.Error
}

func (r *psychosocialReassignRepository) UpdateTeamContactsDupla(ctx context.Context, remisionID, duplaID string) (int64, error) {
	result := r.db.WithContext(ctx).Exec(`
		UPDATE salvia.team_contact
		SET dupla_id = ?,
		    professional_id = NULL,
		    updated_at = ?
		WHERE BTRIM(psicosocial_id::text) = BTRIM(?)
		  AND is_completed = false
		  AND deleted_at IS NULL
	`, duplaID, time.Now(), remisionID)
	return result.RowsAffected, result.Error
}
