package repository

import (
	"bitsflow/internal/models"
	"context"
	"fmt"

	"gorm.io/gorm"
)

// SubmissionRepository extiende el CRUD genérico con la operación transaccional
// de envío de formulario (Unit of Work: timeline + submission + answers).
type SubmissionRepository interface {
	Repository[models.FormSubmission]
	SubmitKoboForm(ctx context.Context, submission *models.FormSubmission, answers []models.Answer, timelineEvent *models.CaseTimeline) error
}

type submissionRepository struct {
	repository[models.FormSubmission]
	db *gorm.DB
}

// NewSubmissionRepository construye un SubmissionRepository listo para usar.
func NewSubmissionRepository(db *gorm.DB) SubmissionRepository {
	return &submissionRepository{
		repository: repository[models.FormSubmission]{db: db},
		db:         db,
	}
}

// SubmitKoboForm ejecuta el envío completo de un formulario como una única transacción atómica.
//
// Orden de operaciones (Unit of Work):
//  1. Inserta el CaseTimeline (ancla forense).
//  2. Asigna el TimelineID al FormSubmission e inserta el submission.
//  3. Asigna el FormSubmissionID a cada Answer y las inserta en batch.
//
// Si cualquier paso falla, se hace rollback completo.
func (r *submissionRepository) SubmitKoboForm(
	ctx context.Context,
	submission *models.FormSubmission,
	answers []models.Answer,
	timelineEvent *models.CaseTimeline,
) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("submission: no se pudo iniciar la transacción: %w", tx.Error)
	}

	// Rollback automático ante panic inesperado.
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	// Paso A: insertar el evento de auditoría forense.
	if err := tx.Create(timelineEvent).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("submission: error al crear timeline event: %w", err)
	}

	// Paso B: anclar el submission al timeline e insertarlo.
	submission.TimelineID = timelineEvent.ID
	if err := tx.Create(submission).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("submission: error al crear form submission: %w", err)
	}

	// Paso C: asignar el submission ID a cada respuesta e insertar en batch.
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
