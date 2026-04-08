package repository

import (
	"bitsflow/internal/models"
	"context"
	"fmt"

	"gorm.io/gorm"
)

// SubmissionRepository extiende el CRUD genérico con la operación transaccional
// de envío de formulario Kobo (Unit of Work: timeline + submission + answers).
type SubmissionRepository interface {
	Repository[models.KoboSubmission]
	SubmitKoboForm(ctx context.Context, submission *models.KoboSubmission, answers []models.KoboAnswer, timelineEvent *models.CaseTimeline) error
}

type submissionRepository struct {
	repository[models.KoboSubmission]
	db *gorm.DB
}

// NewSubmissionRepository construye un SubmissionRepository listo para usar.
func NewSubmissionRepository(db *gorm.DB) SubmissionRepository {
	return &submissionRepository{
		repository: repository[models.KoboSubmission]{db: db},
		db:         db,
	}
}

// SubmitKoboForm ejecuta el envío completo de un formulario como una única transacción atómica.
func (r *submissionRepository) SubmitKoboForm(
	ctx context.Context,
	submission *models.KoboSubmission,
	answers []models.KoboAnswer,
	timelineEvent *models.CaseTimeline,
) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("submission: no se pudo iniciar la transacción: %w", tx.Error)
	}
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(timelineEvent).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("submission: error al crear timeline event: %w", err)
	}
	submission.TimelineID = timelineEvent.ID
	if err := tx.Create(submission).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("submission: error al crear form submission: %w", err)
	}
	if len(answers) > 0 {
		for i := range answers {
			answers[i].FormSubmissionID = submission.ID
		}
		if err := tx.Create(&answers).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("submission: error al crear answers: %w", err)
		}
	}
	return tx.Commit().Error
}
