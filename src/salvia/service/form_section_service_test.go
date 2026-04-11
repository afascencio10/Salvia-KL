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

type MockFormSectionRepo struct{ mock.Mock }

func (m *MockFormSectionRepo) Create(ctx context.Context, e *models.FormSection) error {
	return m.Called(ctx, e).Error(0)
}
func (m *MockFormSectionRepo) Update(ctx context.Context, e *models.FormSection) error {
	return m.Called(ctx, e).Error(0)
}
func (m *MockFormSectionRepo) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockFormSectionRepo) FindByID(ctx context.Context, id string) (*models.FormSection, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.FormSection), args.Error(1)
}
func (m *MockFormSectionRepo) FindWithPagination(ctx context.Context, page, size int) (repository.PageResult[models.FormSection], error) {
	args := m.Called(ctx, page, size)
	return args.Get(0).(repository.PageResult[models.FormSection]), args.Error(1)
}
func (m *MockFormSectionRepo) FindByFormID(ctx context.Context, formID string) ([]models.FormSection, error) {
	args := m.Called(ctx, formID)
	return args.Get(0).([]models.FormSection), args.Error(1)
}

func TestFormSectionService_GetByID(t *testing.T) {
	ctx := context.Background()

	t.Run("éxito", func(t *testing.T) {
		repo := new(MockFormSectionRepo)
		expected := &models.FormSection{ID: "s-1", Name: "Sección A"}
		repo.On("FindByID", ctx, "s-1").Return(expected, nil)
		svc := service.NewFormSectionService(repo)
		got, err := svc.GetByID(ctx, "s-1")
		assert.NoError(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(MockFormSectionRepo)
		repo.On("FindByID", ctx, "x").Return(nil, gorm.ErrRecordNotFound)
		svc := service.NewFormSectionService(repo)
		_, err := svc.GetByID(ctx, "x")
		assert.ErrorIs(t, err, service.ErrFormSectionNotFound)
	})
}

func TestFormSectionService_Create(t *testing.T) {
	ctx := context.Background()
	repo := new(MockFormSectionRepo)
	fs := &models.FormSection{FormID: "f-1", Name: "Nueva"}
	repo.On("Create", ctx, fs).Return(nil)
	svc := service.NewFormSectionService(repo)
	assert.NoError(t, svc.Create(ctx, fs))
}

func TestFormSectionService_Delete(t *testing.T) {
	ctx := context.Background()

	t.Run("éxito", func(t *testing.T) {
		repo := new(MockFormSectionRepo)
		repo.On("Delete", ctx, "s-1").Return(nil)
		svc := service.NewFormSectionService(repo)
		assert.NoError(t, svc.Delete(ctx, "s-1"))
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(MockFormSectionRepo)
		repo.On("Delete", ctx, "x").Return(gorm.ErrRecordNotFound)
		svc := service.NewFormSectionService(repo)
		assert.ErrorIs(t, svc.Delete(ctx, "x"), service.ErrFormSectionNotFound)
	})
}
