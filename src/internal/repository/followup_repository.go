package repository

import (
	"bitsflow/internal/models"
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// FollowUpRepository extiende el CRUD genérico con consultas específicas de follow_up_v2.
// Create, FindByID, Delete, Update, UpdateFields y FindWithPagination
// son heredados de base_repository — NO se reimplementan aquí.
type FollowUpRepository interface {
	Repository[models.FollowUpV2]

	// Existente — conservado
	FindByCaseID(ctx context.Context, caseID string) ([]models.FollowUpV2, error)

	// Nuevos HU-027
	FindByCaseIDOrdered(ctx context.Context, caseID string) ([]models.FollowUpV2, error)
	FindPendingByCaseID(ctx context.Context, caseID string) ([]models.FollowUpV2, error)
	FindCompletedByCaseID(ctx context.Context, caseID string) ([]models.FollowUpV2, error)
	BulkCreate(ctx context.Context, tx *gorm.DB, followUps []models.FollowUpV2) error
	SoftDeleteAndReprogramPending(ctx context.Context, tx *gorm.DB, caseID string) error
	RunInTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error

	FindByAgentAndDate(ctx context.Context, agentID string, date time.Time) ([]models.FollowUpV2, error)
	FindRealizedTodayByAgent(ctx context.Context, agentID string, date time.Time) ([]models.FollowUpV2, error)
	IncrementAttempt(ctx context.Context, followUpID string) error
	// Seguimientos Área
	FindByTeamPaginated(ctx context.Context, team string, filters FollowUpFilters, page, limit int) ([]models.FollowUpV2, int64, error)
	FindPendingByTeamGroupedByAgent(ctx context.Context, team string, fecha string) ([]AgentWorkload, error)
	FindAgentsByTeam(ctx context.Context, team string) ([]AgentOption, error)
	Reschedule(ctx context.Context, id string, fields map[string]interface{}) error
	
	// De backupDevelop
	CloseCaseFollowUps(ctx context.Context, followUpID string) error

	// Hacer seguimiento (De develop)
	LoadVictimInfoByCaseID(ctx context.Context, caseID string) (*VictimCaseInfo, error)
	UpdateFormSubmissionID(ctx context.Context, id string, fsID string) error
	FindByFormSubmissionID(ctx context.Context, formSubmissionID string) (*models.FollowUpV2, error)
	UpdateStatus(ctx context.Context, id string, status string) error
}

// FollowUpFilters contiene los filtros dinámicos para la consulta paginada.
type FollowUpFilters struct {
	Tab         string // "pendientes", "realizados", "vencidos", "todos"
	FechaInicio string // formato YYYY-MM-DD
	FechaFin    string // formato YYYY-MM-DD
	Estado      string
	AgentID     string
}

// AgentWorkload agrupa la carga de seguimientos por agente.
type AgentWorkload struct {
	AgentID string `json:"agent_id"`
	Total   int64  `json:"total"`
}

// AgentOption representa un agente disponible para filtros.
type AgentOption struct {
	AgentID string `json:"agent_id"`
}

type followUpRepository struct {
	repository[models.FollowUpV2]
	db *gorm.DB
}

func NewFollowUpRepository(db *gorm.DB) FollowUpRepository {
	return &followUpRepository{
		repository: repository[models.FollowUpV2]{db: db},
		db:         db,
	}
}

// FindByCaseID retorna los seguimientos de un caso (sin orden garantizado).
// Conservado del patrón original.
func (r *followUpRepository) FindByCaseID(ctx context.Context, caseID string) ([]models.FollowUpV2, error) {
	var items []models.FollowUpV2
	err := r.db.WithContext(ctx).
		Where("case_id = ?", caseID).
		Find(&items).Error
	return items, err
}

// FindByCaseIDOrdered retorna todos los seguimientos del caso ordenados por scheduled_date ASC.
func (r *followUpRepository) FindByCaseIDOrdered(ctx context.Context, caseID string) ([]models.FollowUpV2, error) {
	var items []models.FollowUpV2
	err := r.db.WithContext(ctx).
		Where("case_id = ?", caseID).
		Order("scheduled_date ASC").
		Find(&items).Error
	return items, err
}

// FindPendingByCaseID retorna los seguimientos con status = 'PENDIENTE'.
func (r *followUpRepository) FindPendingByCaseID(ctx context.Context, caseID string) ([]models.FollowUpV2, error) {
	var items []models.FollowUpV2
	err := r.db.WithContext(ctx).
		Where("case_id = ? AND status = ?", caseID, models.FollowUpStatusPendiente).
		Order("scheduled_date ASC").
		Find(&items).Error
	return items, err
}

// FindCompletedByCaseID retorna los seguimientos con status REALIZADO o VENCIDO.
func (r *followUpRepository) FindCompletedByCaseID(ctx context.Context, caseID string) ([]models.FollowUpV2, error) {
	var items []models.FollowUpV2
	err := r.db.WithContext(ctx).
		Where("case_id = ? AND status IN ?", caseID,
			[]string{models.FollowUpStatusRealizado, models.FollowUpStatusVencido}).
		Order("scheduled_date ASC").
		Find(&items).Error
	return items, err
}

// BulkCreate inserta múltiples seguimientos en una sola operación usando la tx provista.
func (r *followUpRepository) BulkCreate(ctx context.Context, tx *gorm.DB, followUps []models.FollowUpV2) error {
	return tx.WithContext(ctx).Create(&followUps).Error
}

// SoftDeleteAndReprogramPending marca los PENDIENTES como REPROGRAMADO y aplica soft-delete.
// Sigue el patrón de submission_repository.go: primero Updates, luego Delete sobre el mismo Where.
func (r *followUpRepository) SoftDeleteAndReprogramPending(ctx context.Context, tx *gorm.DB, caseID string) error {
	return tx.WithContext(ctx).
		Where("case_id = ? AND status = ?", caseID, models.FollowUpStatusPendiente).
		Updates(map[string]interface{}{"status": models.FollowUpStatusReprogramado}).
		Delete(&models.FollowUpV2{}).Error
}

// RunInTransaction ejecuta fn dentro de una transacción GORM.
// Sigue exactamente el patrón de submission_repository.go.
func (r *followUpRepository) RunInTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("followup: no se pudo iniciar la transacción: %w", tx.Error)
	}

	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// FindByAgentAndDate retorna los seguimientos pendientes de un agente para una fecha específica.
// Ordena por:
// 1. nivel Riesgo DESC (Alto > Moderado > Bajo) usando CASE
// 2. lastAttemptBy null primero (NULLS FIRST)
// 3. lastAttemptBy ASC (intentos más antiguos primero)
func (r *followUpRepository) FindByAgentAndDate(ctx context.Context, agentID string, date time.Time) ([]models.FollowUpV2, error) {
	var items []models.FollowUpV2

	// Obtener solo la parte de la fecha (YYYY-MM-DD)
	dateOnly := date.Format("2006-01-02")

	err := r.db.WithContext(ctx).
		Where("agent_id = ? AND status IN (?, ?) AND (scheduled_date::date) = ?",
			agentID,
			models.FollowUpStatusPendiente,
			models.FollowUpStatusReprogramado,
			dateOnly).
		Order(`
			CASE UPPER(risk_status)
				WHEN 'EXTREMO' THEN 1 
				WHEN 'CRÍTICO' THEN 2 
				WHEN 'ALTO' THEN 3
				WHEN 'MODERADO' THEN 4
				WHEN 'BAJO' THEN 5
				ELSE 5 
			END ASC,
			CASE 
				WHEN (scheduled_time IS NULL OR scheduled_time = '00:00:00') THEN 1 
				ELSE 0 
			END ASC,
			scheduled_time ASC,			
			last_attempt_at ASC NULLS FIRST
		`).
		Find(&items).Error

	return items, err
}

// FindRealizedTodayByAgent retorna los seguimientos marcados como REALIZADO
// que tuvieron su último intento el día de hoy.
func (r *followUpRepository) FindRealizedTodayByAgent(ctx context.Context, agentID string, date time.Time) ([]models.FollowUpV2, error) {
	var items []models.FollowUpV2
	dateOnly := date.Format("2006-01-02")

	err := r.db.WithContext(ctx).
		Where("agent_id = ? AND status = ? AND (scheduled_date::date) = ?",
			agentID,
			models.FollowUpStatusRealizado,
			dateOnly).
		Order("scheduled_time ASC").
		Find(&items).Error

	return items, err
}

// IncrementAttempt incrementa el contador de intentos de un seguimiento y actualiza last_attempt_at.
func (r *followUpRepository) IncrementAttempt(ctx context.Context, followUpID string) error {
	return r.db.WithContext(ctx).
		Model(&models.FollowUpV2{}).
		Where("id = ?", followUpID).
		Updates(map[string]interface{}{
			"attempts":        gorm.Expr("attempts + 1"),
			"last_attempt_at": time.Now(),
		}).Error
}

// FindByTeamPaginated retorna seguimientos filtrados por team + filtros dinámicos con paginación.
func (r *followUpRepository) FindByTeamPaginated(ctx context.Context, team string, filters FollowUpFilters, page, limit int) ([]models.FollowUpV2, int64, error) {
	query := r.db.WithContext(ctx).Model(&models.FollowUpV2{})

	// Solo filtrar por team si viene con valor (vacío = todos)
	if team != "" {
		query = query.Where("team = ?", team)
	}

	// Filtros dinámicos
	switch filters.Tab {
	case "pendientes":
		query = query.Where("status = ?", models.FollowUpStatusPendiente)
	case "realizados":
		query = query.Where("status = ?", models.FollowUpStatusRealizado)
	case "vencidos":
		query = query.Where("status = ?", models.FollowUpStatusVencido)
	}
	if filters.Estado != "" && filters.Tab == "" {
		query = query.Where("status = ?", filters.Estado)
	}
	if filters.FechaInicio != "" {
		query = query.Where("scheduled_date >= ?", filters.FechaInicio)
	}
	if filters.FechaFin != "" {
		query = query.Where("scheduled_date <= ?", filters.FechaFin)
	}
	if filters.AgentID != "" {
		query = query.Where("agent_id = ?", filters.AgentID)
	}

	var total int64
	if err := query.Model(&models.FollowUpV2{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []models.FollowUpV2
	offset := page * limit
	if err := query.Order("scheduled_date ASC").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// FindPendingByTeamGroupedByAgent retorna la carga de seguimientos pendientes por agente para una fecha.
func (r *followUpRepository) FindPendingByTeamGroupedByAgent(ctx context.Context, team string, fecha string) ([]AgentWorkload, error) {
	var results []AgentWorkload
	query := r.db.WithContext(ctx).
		Model(&models.FollowUpV2{}).
		Select("agent_id, COUNT(*) as total").
		Where("status = ?", models.FollowUpStatusPendiente)

	if team != "" {
		query = query.Where("team = ?", team)
	}

	if fecha != "" {
		query = query.Where("DATE(scheduled_date) = ?", fecha)
	}

	err := query.Group("agent_id").Order("total DESC").Scan(&results).Error
	return results, err
}

// FindAgentsByTeam retorna los agentes distintos que tienen seguimientos en un team.
func (r *followUpRepository) FindAgentsByTeam(ctx context.Context, team string) ([]AgentOption, error) {
	var results []AgentOption
	query := r.db.WithContext(ctx).
		Model(&models.FollowUpV2{}).
		Select("DISTINCT agent_id")

	if team != "" {
		query = query.Where("team = ?", team)
	}

	err := query.Scan(&results).Error
	return results, err
}

// Reschedule actualiza los campos de reagendamiento de un seguimiento.
func (r *followUpRepository) Reschedule(ctx context.Context, id string, fields map[string]interface{}) error {
	result := r.db.WithContext(ctx).
		Model(&models.FollowUpV2{}).
		Where("id = ?", id).
		Updates(fields)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// CloseCaseFollowUps busca el caseID del seguimiento indicado y cierra todos los 
// seguimientos pendientes o reprogramados del mismo caso.
func (r *followUpRepository) CloseCaseFollowUps(ctx context.Context, followUpID string) error {
	var fu models.FollowUpV2
	if err := r.db.WithContext(ctx).Where("id = ?", followUpID).First(&fu).Error; err != nil {
		return err
	}

	return r.db.WithContext(ctx).
		Model(&models.FollowUpV2{}).
		Where("case_id = ? AND status IN ?", fu.CaseID,
			[]string{models.FollowUpStatusPendiente, models.FollowUpStatusReprogramado}).
		Update("status", models.FollowUpStatusCerrado).Error
}

// LoadVictimInfoByCaseID obtiene la información resumida del caso para la pantalla hacer-seguimiento.
func (r *followUpRepository) LoadVictimInfoByCaseID(ctx context.Context, caseID string) (*VictimCaseInfo, error) {
	var info VictimCaseInfo
	sql := `
		SELECT
			COALESCE(vc.victim_case_victim_names, '')                   AS names,
			COALESCE(vc.victim_case_victim_last_names, '')              AS last_names,
			COALESCE(t.town_name, '')                                   AS town_name,
			COALESCE(f1.victim_case_form1_victim_phone, '')             AS phone,
			COALESCE(f1.victim_case_form1_victim_gender_identity, '')   AS gender_identity,
			COALESCE(f1.victim_case_form1_victim_sexual_orientation, '') AS sexual_orientation,
			COALESCE(f1.victim_case_form1_victim_contact_phone, '')     AS contact_phone,
			f1.victim_case_form1_age                                    AS age
		FROM salvia.victim_case vc
		LEFT JOIN salvia.victim_case_form1 f1 ON f1.victim_case_form1_victim_case = vc.victim_case_id
		LEFT JOIN security.town t ON t.town_code = vc.victim_case_victim_town_code
		WHERE vc.victim_case_id::text = ?
		LIMIT 1`
	return &info, r.db.WithContext(ctx).Raw(sql, caseID).Scan(&info).Error
}

// UpdateFormSubmissionID asigna un formSubmissionId a un seguimiento.
func (r *followUpRepository) UpdateFormSubmissionID(ctx context.Context, id string, fsID string) error {
	return r.db.WithContext(ctx).
		Model(&models.FollowUpV2{}).
		Where("id = ?", id).
		Update("form_submission_id", fsID).Error
}

// FindByFormSubmissionID busca el seguimiento asociado a un formSubmissionId.
func (r *followUpRepository) FindByFormSubmissionID(ctx context.Context, formSubmissionID string) (*models.FollowUpV2, error) {
	var fu models.FollowUpV2
	err := r.db.WithContext(ctx).
		Where("form_submission_id = ?", formSubmissionID).
		First(&fu).Error
	return &fu, err
}

// UpdateStatus actualiza el status de un seguimiento.
// Si el status es REALIZADO también registra el completed_at.
func (r *followUpRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	fields := map[string]interface{}{"status": status}
	if status == models.FollowUpStatusRealizado {
		fields["completed_at"] = time.Now()
	}
	return r.db.WithContext(ctx).
		Model(&models.FollowUpV2{}).
		Where("id = ?", id).
		Updates(fields).Error
}