package service

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"
)

var ErrBarrierV2NotFound = errors.New("barrier_v2: registro no encontrado")

// BarrierFollowUpItem es el DTO de respuesta para el endpoint de seguimientos de barrera.
type BarrierFollowUpItem struct {
	models.BarrierFollowUp
	ActorName string `json:"actorName"`
}

// BarrierDetailResponse agrupa toda la data para la pantalla de detalle de barrera.
type BarrierDetailResponse struct {
	Barrier       models.BarrierV2           `json:"barrier"`
	VictimName    string                     `json:"victimName"`
	VictimAge     *int64                     `json:"victimAge"`
	Location      string                     `json:"location"`
	RiskLevel     *int                       `json:"riskLevel"`
	CaseICode     string                     `json:"caseICode"`
	AgentName     string                     `json:"agentName"`
	// Ubicación resuelta de la barrera
	BarrierDepartment string                 `json:"barrierDepartment"`
	BarrierCity       string                 `json:"barrierCity"`
	BarrierTown       string                 `json:"barrierTown"`
	// Quién identificó la barrera
	CreatedByName string                     `json:"createdByName"`
}

type BarrierV2Service interface {
	GetByID(ctx context.Context, id string) (*models.BarrierV2, error)
	GetDetail(ctx context.Context, id string) (*BarrierDetailResponse, error)
	List(ctx context.Context, page, limit int) (repository.PageResult[models.BarrierV2], error)
	Create(ctx context.Context, b *models.BarrierV2) error
	Update(ctx context.Context, b *models.BarrierV2) error
	Delete(ctx context.Context, id string) error
	// ListByCreatedByIDWithRelations devuelve las barreras creadas por el agente
	// enriquecidas con datos de victim_case (nombres, doc, case_code).
	ListByCreatedByIDWithRelations(ctx context.Context, createdByID string) ([]models.BarrierV2WithRelations, error)
	// ListFollowUps devuelve los seguimientos de una barrera en orden cronológico
	// con el nombre del autor resuelto.
	ListFollowUps(ctx context.Context, barrierID string) ([]BarrierFollowUpItem, error)
}

type barrierV2Service struct {
	repo           repository.BarrierV2Repository
	followUpRepo   repository.BarrierFollowUpRepository
	agentLightRepo repository.AgentLightRepository
}

func NewBarrierV2Service(repo repository.BarrierV2Repository, followUpRepo repository.BarrierFollowUpRepository, agentLightRepo repository.AgentLightRepository) BarrierV2Service {
	return &barrierV2Service{repo: repo, followUpRepo: followUpRepo, agentLightRepo: agentLightRepo}
}

func (s *barrierV2Service) GetDetail(ctx context.Context, id string) (*BarrierDetailResponse, error) {
	b, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBarrierV2NotFound
		}
		return nil, err
	}

	resp := &BarrierDetailResponse{
		Barrier:   *b,
		CaseICode: b.CaseID,
	}

	// Cargar datos de la víctima desde victim_case
	db := s.repo.GetDB()
	if db != nil {
		// Si la barrera tiene tareas pendientes y está en OPEN, actualizar a "En Gestion" en BD
		var pendingCount int64
		db.WithContext(ctx).Raw(`
			SELECT COUNT(*) FROM salvia.case_task
			WHERE barrier_id = ? AND status = 'ToDo' AND deleted_at IS NULL
		`, b.ID).Scan(&pendingCount)
		if pendingCount > 0 && b.Status == models.BarrierV2StatusOpen {
			resp.Barrier.Status = models.BarrierV2StatusEnGestion
			db.Exec(`UPDATE salvia.barrier_v2 SET status = ? WHERE id = ?`, models.BarrierV2StatusEnGestion, b.ID)
		}
		type victimRow struct {
			Names     string `gorm:"column:names"`
			LastNames string `gorm:"column:last_names"`
			Age       *int64 `gorm:"column:age"`
			CityName  string `gorm:"column:city_name"`
			DeptName  string `gorm:"column:dept_name"`
			RiskLevel *int   `gorm:"column:risk_level"`
			AgentName string `gorm:"column:agent_name"`
		}
		var row victimRow
		db.WithContext(ctx).Raw(`
			SELECT
				COALESCE(vc.victim_case_victim_names, '') AS names,
				COALESCE(vc.victim_case_victim_last_names, '') AS last_names,
				EXTRACT(YEAR FROM AGE(NOW(), f2.victim_case_form2_birth_date))::int AS age,
				COALESCE(c.city_name, '') AS city_name,
				COALESCE(d.department_name, '') AS dept_name,
				f2.victim_case_form2_risk_level AS risk_level,
				COALESCE((
					SELECT gup.general_user_profile_names || ' ' || gup.general_user_profile_last_names
					FROM security.general_user gu
					JOIN security.general_user_profile gup ON gup.general_user_profile_id = gu.general_user_general_user_profile
					WHERE gu.general_user_i_code = vc.agent_id LIMIT 1
				), '') AS agent_name
			FROM salvia.victim_case vc
			LEFT JOIN salvia.victim_case_form2 f2 ON f2.victim_case_form2_victim_case = vc.victim_case_id
			LEFT JOIN security.town t ON t.town_code = vc.victim_case_victim_town_code
			LEFT JOIN security.city c ON c.city_id = t.city_id
			LEFT JOIN security.department d ON d.department_id = c.department_id
			WHERE vc.victim_case_i_code = ?
			LIMIT 1
		`, b.CaseID).Scan(&row)

		resp.VictimName = strings.TrimSpace(row.Names + " " + row.LastNames)
		resp.VictimAge = row.Age
		resp.RiskLevel = row.RiskLevel
		resp.AgentName = strings.TrimSpace(row.AgentName)
		if row.CityName != "" && row.DeptName != "" {
			resp.Location = row.CityName + ", " + row.DeptName
		} else if row.CityName != "" {
			resp.Location = row.CityName
		} else {
			resp.Location = row.DeptName
		}

		// Resolver ubicación geográfica de la barrera
		if b.TownID != "" {
			type barrierLocRow struct {
				DeptName string `gorm:"column:dept_name"`
				CityName string `gorm:"column:city_name"`
				TownName string `gorm:"column:town_name"`
			}
			var loc barrierLocRow
			db.WithContext(ctx).Raw(`
				SELECT COALESCE(d.department_name, '') AS dept_name,
				       COALESCE(c.city_name, '') AS city_name,
				       COALESCE(t.town_name, '') AS town_name
				FROM security.town t
				LEFT JOIN security.city c ON c.city_id = t.city_id
				LEFT JOIN security.department d ON d.department_id = c.department_id
				WHERE t.town_code = ?
				LIMIT 1
			`, b.TownID).Scan(&loc)
			resp.BarrierDepartment = loc.DeptName
			resp.BarrierCity = loc.CityName
			resp.BarrierTown = loc.TownName
		} else if b.CityID != "" {
			type cityLocRow struct {
				DeptName string `gorm:"column:dept_name"`
				CityName string `gorm:"column:city_name"`
			}
			var loc cityLocRow
			db.WithContext(ctx).Raw(`
				SELECT COALESCE(d.department_name, '') AS dept_name,
				       COALESCE(c.city_name, '') AS city_name
				FROM security.city c
				LEFT JOIN security.department d ON d.department_id = c.department_id
				WHERE c.city_id = ?
				LIMIT 1
			`, b.CityID).Scan(&loc)
			resp.BarrierDepartment = loc.DeptName
			resp.BarrierCity = loc.CityName
		} else if b.DepartmentID != "" {
			var deptName string
			db.WithContext(ctx).Raw(`SELECT COALESCE(department_name, '') FROM security.department WHERE department_id = ? LIMIT 1`, b.DepartmentID).Scan(&deptName)
			resp.BarrierDepartment = deptName
		}

		// Resolver nombre del creador de la barrera
		if b.CreatedByID != "" {
			var creatorName string
			db.WithContext(ctx).Raw(`
				SELECT COALESCE(gup.general_user_profile_names || ' ' || gup.general_user_profile_last_names, '')
				FROM security.general_user gu
				JOIN security.general_user_profile gup ON gup.general_user_profile_id = gu.general_user_general_user_profile
				WHERE gu.general_user_i_code = ?
			`, b.CreatedByID).Scan(&creatorName)
			resp.CreatedByName = strings.TrimSpace(creatorName)
		}
	}

	return resp, nil
}

func (s *barrierV2Service) GetByID(ctx context.Context, id string) (*models.BarrierV2, error) {
	b, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBarrierV2NotFound
		}
		return nil, err
	}
	return b, nil
}

func (s *barrierV2Service) List(ctx context.Context, page, limit int) (repository.PageResult[models.BarrierV2], error) {
	return s.repo.FindWithPagination(ctx, page, limit)
}

func (s *barrierV2Service) Create(ctx context.Context, b *models.BarrierV2) error {
	if b.Status == "" {
		b.Status = "OPEN"
	}
	return s.repo.Create(ctx, b)
}

func (s *barrierV2Service) Update(ctx context.Context, b *models.BarrierV2) error {
	return s.repo.Update(ctx, b)
}

func (s *barrierV2Service) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrBarrierV2NotFound
	}
	return err
}

func (s *barrierV2Service) ListByCreatedByIDWithRelations(ctx context.Context, createdByID string) ([]models.BarrierV2WithRelations, error) {
	return s.repo.FindByCreatedByIDWithRelations(ctx, createdByID)
}

func (s *barrierV2Service) ListFollowUps(ctx context.Context, barrierID string) ([]BarrierFollowUpItem, error) {
	records, err := s.followUpRepo.FindByBarrierID(ctx, barrierID)
	if err != nil {
		return nil, err
	}
	items := make([]BarrierFollowUpItem, len(records))
	for i, r := range records {
		item := BarrierFollowUpItem{BarrierFollowUp: r}
		if s.agentLightRepo != nil && r.CreatedByID != "" {
			if agent, err := s.agentLightRepo.FindByICode(ctx, r.CreatedByID); err == nil && agent != nil {
				item.ActorName = strings.TrimSpace(agent.Names + " " + agent.LastNames)
			}
		}
		items[i] = item
	}
	return items, nil
}
