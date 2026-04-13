package service_test

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"bitsflow/salvia/service"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockQuestionRepo struct{ mock.Mock }

func (m *MockQuestionRepo) Create(ctx context.Context, e *models.Question) error {
	return m.Called(ctx, e).Error(0)
}
func (m *MockQuestionRepo) Update(ctx context.Context, e *models.Question) error {
	return m.Called(ctx, e).Error(0)
}
func (m *MockQuestionRepo) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockQuestionRepo) FindByID(ctx context.Context, id string) (*models.Question, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Question), args.Error(1)
}
func (m *MockQuestionRepo) FindWithPagination(ctx context.Context, page, size int) (repository.PageResult[models.Question], error) {
	args := m.Called(ctx, page, size)
	return args.Get(0).(repository.PageResult[models.Question]), args.Error(1)
}
func (m *MockQuestionRepo) FindByFormID(ctx context.Context, formID string) ([]models.Question, error) {
	args := m.Called(ctx, formID)
	return args.Get(0).([]models.Question), args.Error(1)
}
func (m *MockQuestionRepo) FindBySectionID(ctx context.Context, sectionID string) ([]models.Question, error) {
	args := m.Called(ctx, sectionID)
	return args.Get(0).([]models.Question), args.Error(1)
}

func TestQuestionService_GetByID(t *testing.T) {
	ctx := context.Background()

	t.Run("éxito", func(t *testing.T) {
		repo := new(MockQuestionRepo)
		expected := &models.Question{ID: "q-1", Description: "¿Pregunta?"}
		repo.On("FindByID", ctx, "q-1").Return(expected, nil)
		svc := service.NewQuestionService(repo)
		got, err := svc.GetByID(ctx, "q-1")
		assert.NoError(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(MockQuestionRepo)
		repo.On("FindByID", ctx, "x").Return(nil, gorm.ErrRecordNotFound)
		svc := service.NewQuestionService(repo)
		_, err := svc.GetByID(ctx, "x")
		assert.ErrorIs(t, err, service.ErrQuestionNotFound)
	})
}

func TestQuestionService_Create(t *testing.T) {
	ctx := context.Background()
	repo := new(MockQuestionRepo)
	q := &models.Question{FormID: "f-1", FormSectionID: "s-1", QuestionTypeID: "BOOLEAN", Description: "¿Sí?"}
	repo.On("Create", ctx, q).Return(nil)
	svc := service.NewQuestionService(repo)
	assert.NoError(t, svc.Create(ctx, q))
}

func TestQuestionService_Delete(t *testing.T) {
	ctx := context.Background()

	t.Run("éxito", func(t *testing.T) {
		repo := new(MockQuestionRepo)
		repo.On("Delete", ctx, "q-1").Return(nil)
		svc := service.NewQuestionService(repo)
		assert.NoError(t, svc.Delete(ctx, "q-1"))
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(MockQuestionRepo)
		repo.On("Delete", ctx, "x").Return(gorm.ErrRecordNotFound)
		svc := service.NewQuestionService(repo)
		assert.ErrorIs(t, svc.Delete(ctx, "x"), service.ErrQuestionNotFound)
	})
}
