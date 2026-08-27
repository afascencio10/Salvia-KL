package repository

import (
	"bitsflow/internal/models"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

// TeamContactRepository gestiona la persistencia de salvia.team_contact.
type TeamContactRepository interface {
	Repository[models.TeamContact]
	// FindByFormSubmissionID busca el team_contact asociado a un form_submission —
	// usado por el procesamiento de guardado de los formularios psicosociales (evento E-02),
	// que recibe el submissionId desde OnEndFormSubmission y necesita el contexto
	// (case_id, psicosocial_id, professional/dupla) para actualizar el estado.
	FindByFormSubmissionID(ctx context.Context, submissionID string) (*models.TeamContact, error)
	// FindByPsicosocialID retorna todos los team_contact (completados o no) de una remisión
	// psicosocial — usado para resolver las "barreras activas" de la sección "Seguimiento a
	// Barreras" (se buscan las barrier_v2 cuyo team_contact_id esté entre estos IDs).
	FindByPsicosocialID(ctx context.Context, psicosocialID string) ([]models.TeamContact, error)
	// FindScheduledOnDate retorna team_contact de sesión psicosocial agendados en la fecha
	// indicada (is_psico_session, no eliminados) — usado para validar disponibilidad (ventana 2h).
	FindScheduledOnDate(ctx context.Context, date time.Time) ([]models.TeamContact, error)
}

type teamContactRepository struct {
	repository[models.TeamContact]
	db *gorm.DB
}

func NewTeamContactRepository(db *gorm.DB) TeamContactRepository {
	return &teamContactRepository{
		repository: repository[models.TeamContact]{db: db},
		db:         db,
	}
}

func (r *teamContactRepository) FindByFormSubmissionID(ctx context.Context, submissionID string) (*models.TeamContact, error) {
	var tc models.TeamContact
	err := r.db.WithContext(ctx).
		Where("form_submission_id = ? AND deleted_at IS NULL", submissionID).
		Order("created_at DESC").
		First(&tc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &tc, nil
}

func (r *teamContactRepository) FindByPsicosocialID(ctx context.Context, psicosocialID string) ([]models.TeamContact, error) {
	var items []models.TeamContact
	err := r.db.WithContext(ctx).
		Where("psicosocial_id = ? AND deleted_at IS NULL", psicosocialID).
		Order("created_at ASC").
		Find(&items).Error
	return items, err
}

func (r *teamContactRepository) FindScheduledOnDate(ctx context.Context, date time.Time) ([]models.TeamContact, error) {
	var items []models.TeamContact
	day := date.Format("2006-01-02")
	err := r.db.WithContext(ctx).
		Where("scheduled_date IS NOT NULL AND (scheduled_date AT TIME ZONE 'UTC')::date = ?::date AND is_psico_session = true AND deleted_at IS NULL", day).
		Find(&items).Error
	return items, err
}
