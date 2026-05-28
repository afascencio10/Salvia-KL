package service_test

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"bitsflow/salvia/service"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// --- Mock FormRepository ---

type MockFormRepo struct{ mock.Mock }

func (m *MockFormRepo) Create(ctx context.Context, e *models.Form) error {
	return m.Called(ctx, e).Error(0)
}
func (m *MockFormRepo) Update(ctx context.Context, e *models.Form) error {
	return m.Called(ctx, e).Error(0)
}
func (m *MockFormRepo) UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error {
	return m.Called(ctx, id, fields).Error(0)
}
func (m *MockFormRepo) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockFormRepo) FindByID(ctx context.Context, id string) (*models.Form, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Form), args.Error(1)
}
func (m *MockFormRepo) FindWithPagination(ctx context.Context, page, size int) (repository.PageResult[models.Form], error) {
	args := m.Called(ctx, page, size)
	return args.Get(0).(repository.PageResult[models.Form]), args.Error(1)
}
func (m *MockFormRepo) FindByStatus(ctx context.Context, status string) ([]models.Form, error) {
	args := m.Called(ctx, status)
	return args.Get(0).([]models.Form), args.Error(1)
}

// --- Tests ---

func TestFormService_GetByID(t *testing.T) {
	ctx := context.Background()

	t.Run("éxito", func(t *testing.T) {
		repo := new(MockFormRepo)
		expected := &models.Form{ID: "uuid-1", Name: "Form A", Status: "active"}
		repo.On("FindByID", ctx, "uuid-1").Return(expected, nil)
		svc := service.NewFormService(service.FormServiceDeps{FormRepo: repo})
		got, err := svc.GetByID(ctx, "uuid-1")
		assert.NoError(t, err)
		assert.Equal(t, expected, got)
		repo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(MockFormRepo)
		repo.On("FindByID", ctx, "x").Return(nil, gorm.ErrRecordNotFound)
		svc := service.NewFormService(service.FormServiceDeps{FormRepo: repo})
		got, err := svc.GetByID(ctx, "x")
		assert.Nil(t, got)
		assert.ErrorIs(t, err, service.ErrFormNotFound)
		repo.AssertExpectations(t)
	})
}

func TestFormService_Create(t *testing.T) {
	ctx := context.Background()

	t.Run("asigna status active por defecto si viene vacío", func(t *testing.T) {
		repo := new(MockFormRepo)
		input := service.CreateFormInput{Name: "Nuevo"}
		
		repo.On("Create", ctx, mock.MatchedBy(func(f *models.Form) bool {
			return f.Name == "Nuevo" && f.Status == "active"
		})).Return(nil)

		svc := service.NewFormService(service.FormServiceDeps{FormRepo: repo})
		got, err := svc.Create(ctx, input)
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, "active", got.Status)
		repo.AssertExpectations(t)
	})
}

func TestFormService_Delete(t *testing.T) {
	ctx := context.Background()

	t.Run("éxito", func(t *testing.T) {
		repo := new(MockFormRepo)
		repo.On("Delete", ctx, "uuid-1").Return(nil)
		svc := service.NewFormService(service.FormServiceDeps{FormRepo: repo})
		assert.NoError(t, svc.Delete(ctx, "uuid-1"))
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(MockFormRepo)
		repo.On("Delete", ctx, "x").Return(gorm.ErrRecordNotFound)
		svc := service.NewFormService(service.FormServiceDeps{FormRepo: repo})
		assert.ErrorIs(t, svc.Delete(ctx, "x"), service.ErrFormNotFound)
	})
}

func TestFormService_List(t *testing.T) {
	ctx := context.Background()
	repo := new(MockFormRepo)
	expected := repository.PageResult[models.Form]{Items: []models.Form{{ID: "1"}}, Total: 1}
	repo.On("FindWithPagination", ctx, 0, 20).Return(expected, nil)
	svc := service.NewFormService(service.FormServiceDeps{FormRepo: repo})
	got, err := svc.List(ctx, 0, 20)
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestFormService_Update(t *testing.T) {
	ctx := context.Background()

	t.Run("éxito", func(t *testing.T) {
		repo := new(MockFormRepo)
		name := "Updated"
		input := service.UpdateFormInput{Name: &name}
		
		repo.On("UpdateFields", ctx, "uuid-1", map[string]interface{}{"name": "Updated"}).Return(nil)
		
		expected := &models.Form{ID: "uuid-1", Name: "Updated"}
		repo.On("FindByID", ctx, "uuid-1").Return(expected, nil)

		svc := service.NewFormService(service.FormServiceDeps{FormRepo: repo})
		got, err := svc.Update(ctx, "uuid-1", input)
		assert.NoError(t, err)
		assert.Equal(t, expected, got)
		repo.AssertExpectations(t)
	})

	t.Run("error infraestructura", func(t *testing.T) {
		repo := new(MockFormRepo)
		name := "Updated"
		input := service.UpdateFormInput{Name: &name}
		dbErr := errors.New("db error")
		repo.On("UpdateFields", ctx, "uuid-err", map[string]interface{}{"name": "Updated"}).Return(dbErr)
		svc := service.NewFormService(service.FormServiceDeps{FormRepo: repo})
		_, err := svc.Update(ctx, "uuid-err", input)
		assert.ErrorIs(t, err, dbErr)
	})
}

func TestFormService_OnEndFormSubmission_CaseClosure_Critical(t *testing.T) {
	ctx := context.Background()

	followUpRepo := new(MockFollowUpRepo)
	answerRepo := new(MockAnswerRepo)
	caseRepo := new(MockVictimCaseLightRepo)
	entityLetterRepo := new(MockEntityLetterRepo)
	caseTaskRepo := new(MockCaseTaskRepo)

	fu := &models.FollowUpV2{
		ID:     "fu-456",
		CaseID: "case-123",
		Status: "pendiente",
	}

	submissionID := "sub-789"
	actorID := "agent-999"
	closureFormID := "da8423ab-1a8c-47db-96b7-d10496df571a"

	followUpRepo.On("FindByFormSubmissionID", ctx, submissionID).Return(fu, nil)
	followUpRepo.On("UpdateStatus", ctx, "fu-456", models.FollowUpStatusRealizado).Return(nil)
	followUpRepo.On("CloseCaseFollowUps", ctx, "fu-456").Return(nil)

	answers := []models.Answer{
		{QuestionID: "d2c6e1af-651f-43b4-8e61-aecafd07443d", Value: "perdida_contacto"},
		{QuestionID: "f4b162fd-adf5-4341-ab4f-162fdadf5341", Value: "false"},
	}
	answerRepo.On("FindDirectBySubmissionID", ctx, submissionID).Return(answers, nil)
	caseRepo.On("UpdateStatus", ctx, "case-123", "cd").Return(nil)

	entityLetterRepo.On("Create", ctx, mock.MatchedBy(func(el *models.EntityLetter) bool {
		return el.CaseID == "case-123" &&
			el.State == models.EntityLetterStatePorProyectar &&
			el.Priority == "normal" &&
			*el.AgentID == actorID &&
			*el.RegisterBy == actorID
	})).Run(func(args mock.Arguments) {
		el := args.Get(1).(*models.EntityLetter)
		el.ID = "letter-888"
	}).Return(nil)

	caseTaskRepo.On("Create", ctx, mock.MatchedBy(func(ct *models.CaseTask) bool {
		return ct.CaseID == "case-123" &&
			ct.Category == "Oficios" &&
			ct.Type == "Escribir oficio" &&
			ct.AssignedUserID == actorID &&
			ct.Status == models.CaseTaskStatusToDo &&
			*ct.FollowUpID == "fu-456" &&
			*ct.EntityLetterID == "letter-888"
	})).Return(nil)

	svc := service.NewFormService(service.FormServiceDeps{
		FollowUpRepo:     followUpRepo,
		AnswerRepo:       answerRepo,
		CaseRepo:         caseRepo,
		EntityLetterRepo: entityLetterRepo,
		CaseTaskRepo:     caseTaskRepo,
	})

	err := svc.OnEndFormSubmission(ctx, closureFormID, submissionID, actorID)
	assert.NoError(t, err)

	followUpRepo.AssertExpectations(t)
	answerRepo.AssertExpectations(t)
	caseRepo.AssertExpectations(t)
	entityLetterRepo.AssertExpectations(t)
	caseTaskRepo.AssertExpectations(t)
}

// --- Mock FollowUpRepository ---
type MockFollowUpRepo struct { mock.Mock }
func (m *MockFollowUpRepo) Create(ctx context.Context, e *models.FollowUpV2) error {
	return m.Called(ctx, e).Error(0)
}
func (m *MockFollowUpRepo) Update(ctx context.Context, e *models.FollowUpV2) error {
	return m.Called(ctx, e).Error(0)
}
func (m *MockFollowUpRepo) UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error {
	return m.Called(ctx, id, fields).Error(0)
}
func (m *MockFollowUpRepo) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockFollowUpRepo) FindByID(ctx context.Context, id string) (*models.FollowUpV2, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.FollowUpV2), args.Error(1)
}
func (m *MockFollowUpRepo) FindWithPagination(ctx context.Context, page, size int) (repository.PageResult[models.FollowUpV2], error) {
	args := m.Called(ctx, page, size)
	return args.Get(0).(repository.PageResult[models.FollowUpV2]), args.Error(1)
}
func (m *MockFollowUpRepo) FindByCaseID(ctx context.Context, caseID string) ([]models.FollowUpV2, error) {
	args := m.Called(ctx, caseID)
	return args.Get(0).([]models.FollowUpV2), args.Error(1)
}
func (m *MockFollowUpRepo) FindByCaseIDOrdered(ctx context.Context, caseID string) ([]models.FollowUpV2, error) {
	args := m.Called(ctx, caseID)
	return args.Get(0).([]models.FollowUpV2), args.Error(1)
}
func (m *MockFollowUpRepo) FindPendingByCaseID(ctx context.Context, caseID string) ([]models.FollowUpV2, error) {
	args := m.Called(ctx, caseID)
	return args.Get(0).([]models.FollowUpV2), args.Error(1)
}
func (m *MockFollowUpRepo) FindCompletedByCaseID(ctx context.Context, caseID string) ([]models.FollowUpV2, error) {
	args := m.Called(ctx, caseID)
	return args.Get(0).([]models.FollowUpV2), args.Error(1)
}
func (m *MockFollowUpRepo) BulkCreate(ctx context.Context, tx *gorm.DB, followUps []models.FollowUpV2) error {
	return m.Called(ctx, tx, followUps).Error(0)
}
func (m *MockFollowUpRepo) SoftDeleteAndReprogramPending(ctx context.Context, tx *gorm.DB, caseID string) error {
	return m.Called(ctx, tx, caseID).Error(0)
}
func (m *MockFollowUpRepo) RunInTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return m.Called(ctx, fn).Error(0)
}
func (m *MockFollowUpRepo) FindByAgentAndDate(ctx context.Context, agentID string, date time.Time) ([]models.FollowUpV2, error) {
	args := m.Called(ctx, agentID, date)
	return args.Get(0).([]models.FollowUpV2), args.Error(1)
}
func (m *MockFollowUpRepo) FindRealizedTodayByAgent(ctx context.Context, agentID string, date time.Time) ([]models.FollowUpV2, error) {
	args := m.Called(ctx, agentID, date)
	return args.Get(0).([]models.FollowUpV2), args.Error(1)
}
func (m *MockFollowUpRepo) IncrementAttempt(ctx context.Context, followUpID string) error {
	return m.Called(ctx, followUpID).Error(0)
}
func (m *MockFollowUpRepo) FindByTeamPaginated(ctx context.Context, team string, filters repository.FollowUpFilters, page, limit int) ([]models.FollowUpV2, int64, error) {
	args := m.Called(ctx, team, filters, page, limit)
	return args.Get(0).([]models.FollowUpV2), args.Get(1).(int64), args.Error(2)
}
func (m *MockFollowUpRepo) FindPendingByTeamGroupedByAgent(ctx context.Context, team string, fecha string) ([]repository.AgentWorkload, error) {
	args := m.Called(ctx, team, fecha)
	return args.Get(0).([]repository.AgentWorkload), args.Error(1)
}
func (m *MockFollowUpRepo) FindAgentsByTeam(ctx context.Context, team string) ([]repository.AgentOption, error) {
	args := m.Called(ctx, team)
	return args.Get(0).([]repository.AgentOption), args.Error(1)
}
func (m *MockFollowUpRepo) Reschedule(ctx context.Context, id string, fields map[string]interface{}) error {
	return m.Called(ctx, id, fields).Error(0)
}
func (m *MockFollowUpRepo) FindWorkloadByDates(ctx context.Context, team string, dates []time.Time) ([]repository.AgentDateWorkload, error) {
	args := m.Called(ctx, team, dates)
	return args.Get(0).([]repository.AgentDateWorkload), args.Error(1)
}
func (m *MockFollowUpRepo) FindGlobalWorkloadByTeam(ctx context.Context, team string) ([]repository.AgentWorkload, error) {
	args := m.Called(ctx, team)
	return args.Get(0).([]repository.AgentWorkload), args.Error(1)
}
func (m *MockFollowUpRepo) CloseCaseFollowUps(ctx context.Context, followUpID string) error {
	return m.Called(ctx, followUpID).Error(0)
}
func (m *MockFollowUpRepo) LoadVictimInfoByCaseID(ctx context.Context, caseID string) (*repository.VictimCaseInfo, error) {
	args := m.Called(ctx, caseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.VictimCaseInfo), args.Error(1)
}
func (m *MockFollowUpRepo) UpdateFormSubmissionID(ctx context.Context, id string, fsID string) error {
	return m.Called(ctx, id, fsID).Error(0)
}
func (m *MockFollowUpRepo) UpdateFormIDAndSubmissionID(ctx context.Context, id string, formID string, fsID string) error {
	return m.Called(ctx, id, formID, fsID).Error(0)
}
func (m *MockFollowUpRepo) FindByFormSubmissionID(ctx context.Context, formSubmissionID string) (*models.FollowUpV2, error) {
	args := m.Called(ctx, formSubmissionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.FollowUpV2), args.Error(1)
}
func (m *MockFollowUpRepo) UpdateStatus(ctx context.Context, id string, status string) error {
	return m.Called(ctx, id, status).Error(0)
}
func (m *MockFollowUpRepo) CreateTimelineEvent(ctx context.Context, event *models.CaseTimelineEvent) error {
	return m.Called(ctx, event).Error(0)
}

// --- Mock AnswerRepository ---
type MockAnswerRepo struct { mock.Mock }
func (m *MockAnswerRepo) Create(ctx context.Context, e *models.Answer) error {
	return m.Called(ctx, e).Error(0)
}
func (m *MockAnswerRepo) Update(ctx context.Context, e *models.Answer) error {
	return m.Called(ctx, e).Error(0)
}
func (m *MockAnswerRepo) UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error {
	return m.Called(ctx, id, fields).Error(0)
}
func (m *MockAnswerRepo) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockAnswerRepo) FindByID(ctx context.Context, id string) (*models.Answer, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Answer), args.Error(1)
}
func (m *MockAnswerRepo) FindWithPagination(ctx context.Context, page, size int) (repository.PageResult[models.Answer], error) {
	args := m.Called(ctx, page, size)
	return args.Get(0).(repository.PageResult[models.Answer]), args.Error(1)
}
func (m *MockAnswerRepo) FindBySubmissionID(ctx context.Context, submissionID string) ([]models.Answer, error) {
	args := m.Called(ctx, submissionID)
	return args.Get(0).([]models.Answer), args.Error(1)
}
func (m *MockAnswerRepo) FindDirectBySubmissionID(ctx context.Context, submissionID string) ([]models.Answer, error) {
	args := m.Called(ctx, submissionID)
	return args.Get(0).([]models.Answer), args.Error(1)
}
func (m *MockAnswerRepo) FindByRepeaterEntryID(ctx context.Context, entryID string) ([]models.Answer, error) {
	args := m.Called(ctx, entryID)
	return args.Get(0).([]models.Answer), args.Error(1)
}
func (m *MockAnswerRepo) FindByRepeaterEntryIDs(ctx context.Context, entryIDs []string) ([]models.Answer, error) {
	args := m.Called(ctx, entryIDs)
	return args.Get(0).([]models.Answer), args.Error(1)
}
func (m *MockAnswerRepo) FindByQuestionID(ctx context.Context, questionID string) ([]models.Answer, error) {
	args := m.Called(ctx, questionID)
	return args.Get(0).([]models.Answer), args.Error(1)
}

// --- Mock VictimCaseLightRepository ---
type MockVictimCaseLightRepo struct { mock.Mock }
func (m *MockVictimCaseLightRepo) FindByID(ctx context.Context, caseId string) (*models.VictimCaseLight, error) {
	args := m.Called(ctx, caseId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.VictimCaseLight), args.Error(1)
}
func (m *MockVictimCaseLightRepo) FindByICode(ctx context.Context, iCode string) (*models.VictimCaseLight, error) {
	args := m.Called(ctx, iCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.VictimCaseLight), args.Error(1)
}
func (m *MockVictimCaseLightRepo) UpdateStatus(ctx context.Context, caseID string, status string) error {
	return m.Called(ctx, caseID, status).Error(0)
}

// --- Mock EntityLetterRepository ---
type MockEntityLetterRepo struct { mock.Mock }
func (m *MockEntityLetterRepo) Create(ctx context.Context, e *models.EntityLetter) error {
	return m.Called(ctx, e).Error(0)
}
func (m *MockEntityLetterRepo) Update(ctx context.Context, e *models.EntityLetter) error {
	return m.Called(ctx, e).Error(0)
}
func (m *MockEntityLetterRepo) UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error {
	return m.Called(ctx, id, fields).Error(0)
}
func (m *MockEntityLetterRepo) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockEntityLetterRepo) FindByID(ctx context.Context, id string) (*models.EntityLetter, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.EntityLetter), args.Error(1)
}
func (m *MockEntityLetterRepo) FindWithPagination(ctx context.Context, page, size int) (repository.PageResult[models.EntityLetter], error) {
	args := m.Called(ctx, page, size)
	return args.Get(0).(repository.PageResult[models.EntityLetter]), args.Error(1)
}
func (m *MockEntityLetterRepo) FindByCaseID(ctx context.Context, caseID string) ([]models.EntityLetter, error) {
	args := m.Called(ctx, caseID)
	return args.Get(0).([]models.EntityLetter), args.Error(1)
}
func (m *MockEntityLetterRepo) FindByBarrierID(ctx context.Context, barrierID string) ([]models.EntityLetter, error) {
	args := m.Called(ctx, barrierID)
	return args.Get(0).([]models.EntityLetter), args.Error(1)
}
func (m *MockEntityLetterRepo) FindByState(ctx context.Context, state string, page, pageSize int) (repository.PageResult[models.EntityLetter], error) {
	args := m.Called(ctx, state, page, pageSize)
	return args.Get(0).(repository.PageResult[models.EntityLetter]), args.Error(1)
}
func (m *MockEntityLetterRepo) FindByAgentID(ctx context.Context, agentID string) ([]models.EntityLetter, error) {
	args := m.Called(ctx, agentID)
	return args.Get(0).([]models.EntityLetter), args.Error(1)
}
func (m *MockEntityLetterRepo) FindByNotificationUserID(ctx context.Context, notificationUserID string) ([]models.EntityLetter, error) {
	args := m.Called(ctx, notificationUserID)
	return args.Get(0).([]models.EntityLetter), args.Error(1)
}
func (m *MockEntityLetterRepo) FindByAgentIDWithRelations(ctx context.Context, agentID string) ([]models.EntityLetterWithRelations, error) {
	args := m.Called(ctx, agentID)
	return args.Get(0).([]models.EntityLetterWithRelations), args.Error(1)
}
func (m *MockEntityLetterRepo) FindByNotificationUserIDWithRelations(ctx context.Context, notifUserID string) ([]models.EntityLetterWithRelations, error) {
	args := m.Called(ctx, notifUserID)
	return args.Get(0).([]models.EntityLetterWithRelations), args.Error(1)
}
func (m *MockEntityLetterRepo) UpdateState(ctx context.Context, id, state string) error {
	return m.Called(ctx, id, state).Error(0)
}

// --- Mock CaseTaskRepository ---
type MockCaseTaskRepo struct { mock.Mock }
func (m *MockCaseTaskRepo) Create(ctx context.Context, e *models.CaseTask) error {
	return m.Called(ctx, e).Error(0)
}
func (m *MockCaseTaskRepo) Update(ctx context.Context, e *models.CaseTask) error {
	return m.Called(ctx, e).Error(0)
}
func (m *MockCaseTaskRepo) UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error {
	return m.Called(ctx, id, fields).Error(0)
}
func (m *MockCaseTaskRepo) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockCaseTaskRepo) FindByID(ctx context.Context, id string) (*models.CaseTask, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CaseTask), args.Error(1)
}
func (m *MockCaseTaskRepo) FindWithPagination(ctx context.Context, page, size int) (repository.PageResult[models.CaseTask], error) {
	args := m.Called(ctx, page, size)
	return args.Get(0).(repository.PageResult[models.CaseTask]), args.Error(1)
}
func (m *MockCaseTaskRepo) FindByAssignedUserID(ctx context.Context, assignedUserID string) ([]models.CaseTask, error) {
	args := m.Called(ctx, assignedUserID)
	return args.Get(0).([]models.CaseTask), args.Error(1)
}
func (m *MockCaseTaskRepo) FindByAssignedUserIDWithRelations(ctx context.Context, assignedUserID string) ([]models.CaseTaskWithRelations, error) {
	args := m.Called(ctx, assignedUserID)
	return args.Get(0).([]models.CaseTaskWithRelations), args.Error(1)
}
func (m *MockCaseTaskRepo) FindByCaseID(ctx context.Context, caseID string) ([]models.CaseTask, error) {
	args := m.Called(ctx, caseID)
	return args.Get(0).([]models.CaseTask), args.Error(1)
}
func (m *MockCaseTaskRepo) FindByBarrierID(ctx context.Context, barrierID string) ([]models.CaseTask, error) {
	args := m.Called(ctx, barrierID)
	return args.Get(0).([]models.CaseTask), args.Error(1)
}
func (m *MockCaseTaskRepo) FindTodoByEntityLetterID(ctx context.Context, entityLetterID string) (*models.CaseTask, error) {
	args := m.Called(ctx, entityLetterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CaseTask), args.Error(1)
}
