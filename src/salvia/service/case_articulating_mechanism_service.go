package service

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
	"unicode/utf8"

	"gorm.io/gorm"
)

var (
	ErrMechanismInvalidScope      = errors.New("mecanismo: alcance inválido")
	ErrMechanismMissingLocation   = errors.New("mecanismo: ubicación requerida para este alcance")
	ErrMechanismUserRequired      = errors.New("mecanismo: selecciona un usuario")
	ErrMechanismUserNotFound      = errors.New("mecanismo: el usuario no existe o no tiene rol Mecanismo")
	ErrMechanismUserOutsideGeo    = errors.New("mecanismo: el usuario no corresponde a la ubicación seleccionada")
	ErrMechanismNotFound          = errors.New("mecanismo: registro no encontrado")
	ErrMechanismCommentEmpty      = errors.New("mecanismo: el comentario no puede estar vacío")
	ErrMechanismCommitmentTooLong = errors.New("mecanismo: el compromiso supera el máximo de caracteres")
)

const maxSalviaCommitmentRunes = 2000
const maxCommentRunes = 2000

type CreateArticulatingMechanismInput struct {
	CaseID           string
	Scope            string
	AssignedUserID   string
	DepartmentID     string
	CityID           string
	TownCode         string
	SalviaCommitment string
	CreatedByID      string
	ActorName        string
}

type AddMechanismCommentInput struct {
	CaseID      string
	MechanismID string
	AuthorID    string
	AuthorName  string
	Body        string
}

type ListMechanismUsersInput struct {
	Scope        string
	DepartmentID string
	CityID       string
	TownCode     string
}

type CaseArticulatingMechanismService interface {
	ListByCase(ctx context.Context, caseID string) ([]models.CaseArticulatingMechanism, error)
	Create(ctx context.Context, input CreateArticulatingMechanismInput) (*models.CaseArticulatingMechanism, error)
	AddComment(ctx context.Context, input AddMechanismCommentInput) (*models.CaseArticulatingMechanismComment, error)
	ListUsers(ctx context.Context, input ListMechanismUsersInput) ([]models.MechanismUserOption, error)
}

type caseArticulatingMechanismService struct {
	repo         repository.CaseArticulatingMechanismRepository
	timelineRepo repository.CaseTimelineEventRepository
}

func NewCaseArticulatingMechanismService(
	repo repository.CaseArticulatingMechanismRepository,
	timelineRepo repository.CaseTimelineEventRepository,
) CaseArticulatingMechanismService {
	return &caseArticulatingMechanismService{repo: repo, timelineRepo: timelineRepo}
}

func (s *caseArticulatingMechanismService) ListByCase(ctx context.Context, caseID string) ([]models.CaseArticulatingMechanism, error) {
	rows, err := s.repo.FindByCaseID(ctx, caseID)
	if err != nil {
		return nil, err
	}
	if rows == nil {
		return []models.CaseArticulatingMechanism{}, nil
	}
	return rows, nil
}

func (s *caseArticulatingMechanismService) ListUsers(ctx context.Context, input ListMechanismUsersInput) ([]models.MechanismUserOption, error) {
	scope := strings.TrimSpace(strings.ToLower(input.Scope))
	if !validMechanismScope(scope) {
		return nil, ErrMechanismInvalidScope
	}
	if err := requireLocationForScope(scope, input.DepartmentID, input.CityID, input.TownCode, false); err != nil {
		return nil, err
	}
	return s.repo.ListMechanismUsers(ctx, repository.MechanismUserFilter{
		Scope:        scope,
		DepartmentID: strings.TrimSpace(input.DepartmentID),
		CityID:       strings.TrimSpace(input.CityID),
		TownCode:     strings.TrimSpace(input.TownCode),
	})
}

func (s *caseArticulatingMechanismService) Create(ctx context.Context, input CreateArticulatingMechanismInput) (*models.CaseArticulatingMechanism, error) {
	scope := strings.TrimSpace(strings.ToLower(input.Scope))
	if !validMechanismScope(scope) {
		return nil, ErrMechanismInvalidScope
	}
	userID := strings.TrimSpace(input.AssignedUserID)
	if userID == "" {
		return nil, ErrMechanismUserRequired
	}
	deptID := strings.TrimSpace(input.DepartmentID)
	cityID := strings.TrimSpace(input.CityID)
	townCode := strings.TrimSpace(input.TownCode)
	if err := requireLocationForScope(scope, deptID, cityID, townCode, true); err != nil {
		return nil, err
	}

	commitment := strings.TrimSpace(input.SalviaCommitment)
	if utf8.RuneCountInString(commitment) > maxSalviaCommitmentRunes {
		return nil, ErrMechanismCommitmentTooLong
	}

	user, err := s.repo.FindMechanismUser(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMechanismUserNotFound
		}
		return nil, err
	}
	candidates, err := s.repo.ListMechanismUsers(ctx, repository.MechanismUserFilter{
		Scope:        scope,
		DepartmentID: deptID,
		CityID:       cityID,
		TownCode:     townCode,
	})
	if err != nil {
		return nil, err
	}
	if !mechanismUserInList(candidates, user.ICode) {
		return nil, ErrMechanismUserOutsideGeo
	}

	row := &models.CaseArticulatingMechanism{
		CaseID:           strings.TrimSpace(input.CaseID),
		Scope:            scope,
		Status:           models.CaseArticulatingMechanismStatusActive,
		AssignedUserID:   user.ICode,
		AssignedUserName: strings.TrimSpace(user.FullName),
		CreatedByID:      strings.TrimSpace(input.CreatedByID),
	}
	if commitment != "" {
		row.SalviaCommitment = &commitment
	}

	if scope != models.ArticulatingMechanismScopeNacional {
		if deptID != "" {
			name, err := s.repo.ResolveDepartmentName(ctx, deptID)
			if err != nil {
				return nil, err
			}
			row.DepartmentID = strPtr(deptID)
			if name != "" {
				row.DepartmentName = &name
			}
		}
		if cityID != "" {
			name, err := s.repo.ResolveCityName(ctx, cityID)
			if err != nil {
				return nil, err
			}
			row.CityID = strPtr(cityID)
			if name != "" {
				row.CityName = &name
			}
		}
		if scope == models.ArticulatingMechanismScopeMunicipal && townCode != "" {
			name, err := s.repo.ResolveTownName(ctx, townCode)
			if err != nil {
				return nil, err
			}
			row.TownCode = strPtr(townCode)
			if name != "" {
				row.TownName = &name
			}
		}
	}

	if err := s.repo.Create(ctx, row); err != nil {
		return nil, err
	}

	if s.timelineRepo != nil {
		now := time.Now()
		event := &models.CaseTimelineEvent{
			CaseID:                  row.CaseID,
			EventType:               models.TimelineEventMecanismo,
			Description:             timelineDescription(row),
			ActorID:                 input.CreatedByID,
			ActorName:               strings.TrimSpace(input.ActorName),
			Category:                models.TimelineCategoryMecanismos,
			Type:                    models.TimelineTypeMecanismoRegistrado,
			Icon:                    models.TimelineIconMecanismo,
			Date:                    now,
			EventUserID:             input.CreatedByID,
			Color:                   models.TimelineColorPurple,
			ArticulatingMechanismID: row.ID,
			CreatedAt:               now,
		}
		if err := s.timelineRepo.Create(ctx, event); err != nil {
			log.Printf("[CaseArticulatingMechanism] WARN: no se pudo registrar timeline del caso %s: %v", row.CaseID, err)
		}
	}

	if row.Comments == nil {
		row.Comments = []models.CaseArticulatingMechanismComment{}
	}
	return row, nil
}

func (s *caseArticulatingMechanismService) AddComment(ctx context.Context, input AddMechanismCommentInput) (*models.CaseArticulatingMechanismComment, error) {
	body := strings.TrimSpace(input.Body)
	if body == "" {
		return nil, ErrMechanismCommentEmpty
	}
	if utf8.RuneCountInString(body) > maxCommentRunes {
		return nil, ErrMechanismCommentEmpty
	}
	if _, err := s.repo.FindByIDAndCaseID(ctx, input.MechanismID, input.CaseID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMechanismNotFound
		}
		return nil, err
	}
	authorName := strings.TrimSpace(input.AuthorName)
	if authorName == "" {
		if n, err := s.repo.UserFullName(ctx, input.AuthorID); err == nil && strings.TrimSpace(n) != "" {
			authorName = n
		} else {
			authorName = input.AuthorID
		}
	}
	c := &models.CaseArticulatingMechanismComment{
		MechanismID: input.MechanismID,
		AuthorID:    input.AuthorID,
		AuthorName:  authorName,
		Body:        body,
	}
	if err := s.repo.CreateComment(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func validMechanismScope(scope string) bool {
	return scope == models.ArticulatingMechanismScopeMunicipal ||
		scope == models.ArticulatingMechanismScopeDepartamental ||
		scope == models.ArticulatingMechanismScopeNacional
}

// requireLocationForScope: forCreate=true exige ciudad en municipal; en listado
// departamental la ciudad es opcional.
func requireLocationForScope(scope, deptID, cityID, townCode string, _ bool) error {
	deptID = strings.TrimSpace(deptID)
	cityID = strings.TrimSpace(cityID)
	townCode = strings.TrimSpace(townCode)
	switch scope {
	case models.ArticulatingMechanismScopeMunicipal:
		if deptID == "" || cityID == "" || townCode == "" {
			return ErrMechanismMissingLocation
		}
	case models.ArticulatingMechanismScopeDepartamental:
		if deptID == "" {
			return ErrMechanismMissingLocation
		}
	}
	return nil
}

func mechanismUserInList(users []models.MechanismUserOption, userICode string) bool {
	for _, u := range users {
		if u.ICode == userICode {
			return true
		}
	}
	return false
}

func timelineDescription(row *models.CaseArticulatingMechanism) string {
	scopeLabel := map[string]string{
		models.ArticulatingMechanismScopeMunicipal:     "Municipal",
		models.ArticulatingMechanismScopeDepartamental: "Departamental",
		models.ArticulatingMechanismScopeNacional:      "Nacional",
	}[row.Scope]
	loc := locationLabel(row)
	if loc != "" {
		return fmt.Sprintf("Se registró escalamiento %s a %s (%s)", scopeLabel, row.AssignedUserName, loc)
	}
	return fmt.Sprintf("Se registró escalamiento %s a %s", scopeLabel, row.AssignedUserName)
}

func locationLabel(row *models.CaseArticulatingMechanism) string {
	parts := make([]string, 0, 3)
	if row.TownName != nil && *row.TownName != "" {
		parts = append(parts, *row.TownName)
	}
	if row.CityName != nil && *row.CityName != "" {
		parts = append(parts, *row.CityName)
	}
	if row.DepartmentName != nil && *row.DepartmentName != "" {
		parts = append(parts, *row.DepartmentName)
	}
	if len(parts) == 0 {
		return ""
	}
	if len(parts) >= 3 {
		return fmt.Sprintf("%s, %s (%s)", parts[0], parts[1], parts[2])
	}
	return strings.Join(parts, ", ")
}

func strPtr(s string) *string {
	v := s
	return &v
}
