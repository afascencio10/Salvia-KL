package repository

import (
	internaldb "bitsflow/internal/db"
	"bitsflow/internal/models"
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DuplaRepository acceso a salvia.dupla.
type DuplaRepository interface {
	ListActive(ctx context.Context) ([]models.DuplaOption, error)
	ListActiveEnriched(ctx context.Context) ([]models.DuplaAdminItem, error)
	FindEnrichedByID(ctx context.Context, id string) (*models.DuplaAdminItem, error)
	ExistsActiveByName(ctx context.Context, name string, excludeID string) (bool, error)
	ExistsActivePsychologistConflict(ctx context.Context, psychologistID, excludeID string) (bool, error)
	Create(ctx context.Context, name, psychologistID, socialWorkerID string) (*models.Dupla, error)
	UpdateActive(ctx context.Context, id, name, psychologistID, socialWorkerID string) error
	ExistsActiveByID(ctx context.Context, id string) (bool, error)
	IsInUseNonClosed(ctx context.Context, duplaID string) (bool, error)
	SoftDelete(ctx context.Context, id string) error
}

type duplaRepository struct {
	db *gorm.DB
}

func NewDuplaRepository(db *gorm.DB) DuplaRepository {
	return &duplaRepository{db: db}
}

func (r *duplaRepository) ListActive(ctx context.Context) ([]models.DuplaOption, error) {
	var items []models.DuplaOption
	err := r.db.WithContext(ctx).
		Table("salvia.dupla").
		Select("id, name").
		Where("deleted_at IS NULL").
		Order("name ASC").
		Scan(&items).Error
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []models.DuplaOption{}
	}
	return items, nil
}

const duplaEnrichedSelectSQL = `
SELECT
	d.id,
	d.name,
	BTRIM(d.psychologist_id::text) AS psychologist_id,
	TRIM(COALESCE(ps_gup.general_user_profile_names, '') || ' ' || COALESCE(ps_gup.general_user_profile_last_names, '')) AS psychologist_name,
	BTRIM(d.social_worker_id::text) AS social_worker_id,
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
WHERE d.deleted_at IS NULL`

// ListActiveEnriched lista duplas activas con ids y nombres de integrantes (Administrar Duplas E01).
func (r *duplaRepository) ListActiveEnriched(ctx context.Context) ([]models.DuplaAdminItem, error) {
	sql := duplaEnrichedSelectSQL + `
ORDER BY d.name ASC`

	var items []models.DuplaAdminItem
	err := internaldb.WithRetry(func() error {
		return r.db.WithContext(ctx).Raw(sql).Scan(&items).Error
	})
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []models.DuplaAdminItem{}
	}
	return items, nil
}

func (r *duplaRepository) FindEnrichedByID(ctx context.Context, id string) (*models.DuplaAdminItem, error) {
	sql := duplaEnrichedSelectSQL + `
  AND d.id = ?
LIMIT 1`

	var item models.DuplaAdminItem
	err := r.db.WithContext(ctx).Raw(sql, id).Scan(&item).Error
	if err != nil {
		return nil, err
	}
	if item.ID == "" {
		return nil, gorm.ErrRecordNotFound
	}
	return &item, nil
}

func (r *duplaRepository) ExistsActiveByName(ctx context.Context, name string, excludeID string) (bool, error) {
	var id string
	q := r.db.WithContext(ctx).
		Table("salvia.dupla").
		Select("id").
		Where("deleted_at IS NULL").
		Where("name = ?", name)
	if excludeID != "" {
		q = q.Where("id <> ?", excludeID)
	}
	err := q.Limit(1).Scan(&id).Error
	if err != nil {
		return false, err
	}
	return id != "", nil
}

func (r *duplaRepository) ExistsActivePsychologistConflict(ctx context.Context, psychologistID, excludeID string) (bool, error) {
	var id string
	q := r.db.WithContext(ctx).Raw(`
SELECT id
FROM salvia.dupla
WHERE deleted_at IS NULL
  AND BTRIM(psychologist_id::text) = BTRIM(?)
  AND (? = '' OR id <> ?)
LIMIT 1
`, psychologistID, excludeID, excludeID)
	err := q.Scan(&id).Error
	if err != nil {
		return false, err
	}
	return id != "", nil
}

func (r *duplaRepository) Create(ctx context.Context, name, psychologistID, socialWorkerID string) (*models.Dupla, error) {
	now := time.Now()
	dupla := &models.Dupla{
		ID:             uuid.NewString(),
		Name:           name,
		PsychologistID: psychologistID,
		SocialWorkerID: socialWorkerID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	err := r.db.WithContext(ctx).Create(dupla).Error
	if err != nil {
		return nil, err
	}
	return dupla, nil
}

func (r *duplaRepository) UpdateActive(ctx context.Context, id, name, psychologistID, socialWorkerID string) error {
	result := r.db.WithContext(ctx).Exec(`
UPDATE salvia.dupla
SET name = ?,
    psychologist_id = ?,
    social_worker_id = ?,
    updated_at = ?
WHERE id = ?
  AND deleted_at IS NULL
`, name, psychologistID, socialWorkerID, time.Now(), id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *duplaRepository) ExistsActiveByID(ctx context.Context, id string) (bool, error) {
	var found string
	err := r.db.WithContext(ctx).
		Table("salvia.dupla").
		Select("id").
		Where("id = ? AND deleted_at IS NULL", id).
		Limit(1).
		Scan(&found).Error
	if err != nil {
		return false, err
	}
	return found != "", nil
}

// IsInUseNonClosed true si la dupla está referenciada por remisión o sesión
// cuya psychosocial_support no está en estado cerrado (E07).
func (r *duplaRepository) IsInUseNonClosed(ctx context.Context, duplaID string) (bool, error) {
	const sql = `
SELECT 1 AS ok
WHERE EXISTS (
	SELECT 1
	FROM salvia.psychosocial_support ps
	WHERE BTRIM(ps.dupla_id::text) = BTRIM(?)
	  AND ps.deleted_at IS NULL
	  AND COALESCE(ps.status, '') <> ?
)
OR EXISTS (
	SELECT 1
	FROM salvia.team_contact tc
	JOIN salvia.psychosocial_support ps
	  ON BTRIM(ps.id::text) = BTRIM(tc.psicosocial_id::text)
	 AND ps.deleted_at IS NULL
	WHERE BTRIM(tc.dupla_id::text) = BTRIM(?)
	  AND tc.deleted_at IS NULL
	  AND COALESCE(ps.status, '') <> ?
)
LIMIT 1`

	var ok int
	err := r.db.WithContext(ctx).Raw(
		sql,
		duplaID, models.PsychosocialSupportStatusCerrado,
		duplaID, models.PsychosocialSupportStatusCerrado,
	).Scan(&ok).Error
	if err != nil {
		return false, err
	}
	return ok == 1, nil
}

func (r *duplaRepository) SoftDelete(ctx context.Context, id string) error {
	now := time.Now()
	result := r.db.WithContext(ctx).Exec(`
UPDATE salvia.dupla
SET deleted_at = ?,
    updated_at = ?
WHERE id = ?
  AND deleted_at IS NULL
`, now, now, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
