package service_test

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"bitsflow/salvia/service"
	"context"
	"errors"
	"testing"

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
func (m *MockFormRepo) FindByCampaignID(ctx context.Context, cid string) ([]models.Form, error) {
	args := m.Called(ctx, cid)
	return args.Get(0).([]models.Form), args.Error(1)
}

// --- Tests ---

func TestFormService_GetByID(t *testing.T) {
	ctx := context.Background()

	t.Run("éxito", func(t *testing.T) {
		repo := new(MockFormRepo)
		expected := &models.Form{ID: "uuid-1", Name: "Form A", Status: "ACTIVE"}
		repo.On("FindByID", ctx, "uuid-1").Return(expected, nil)
		svc := service.NewFormService(repo)
		got, err := svc.GetByID(ctx, "uuid-1")
		assert.NoError(t, err)
		assert.Equal(t, expected, got)
		repo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(MockFormRepo)
		repo.On("FindByID", ctx, "x").Return(nil, gorm.ErrRecordNotFound)
		svc := service.NewFormService(repo)
		got, err := svc.GetByID(ctx, "x")
		assert.Nil(t, got)
		assert.ErrorIs(t, err, service.ErrFormNotFound)
		repo.AssertExpectations(t)
	})
}

func TestFormService_Create(t *testing.T) {
	ctx := context.Background()

	t.Run("asigna status ACTIVE por defecto", func(t *testing.T) {
		repo := new(MockFormRepo)
		f := &models.Form{Name: "Nuevo"}
		repo.On("Create", ctx, f).Return(nil)
		svc := service.NewFormService(repo)
		err := svc.Create(ctx, f)
		assert.NoError(t, err)
		assert.Equal(t, "ACTIVE", f.Status)
		repo.AssertExpectations(t)
	})
}

func TestFormService_Delete(t *testing.T) {
	ctx := context.Background()

	t.Run("éxito", func(t *testing.T) {
		repo := new(MockFormRepo)
		repo.On("Delete", ctx, "uuid-1").Return(nil)
		svc := service.NewFormService(repo)
		assert.NoError(t, svc.Delete(ctx, "uuid-1"))
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(MockFormRepo)
		repo.On("Delete", ctx, "x").Return(gorm.ErrRecordNotFound)
		svc := service.NewFormService(repo)
		assert.ErrorIs(t, svc.Delete(ctx, "x"), service.ErrFormNotFound)
	})
}

func TestFormService_List(t *testing.T) {
	ctx := context.Background()
	repo := new(MockFormRepo)
	expected := repository.PageResult[models.Form]{Items: []models.Form{{ID: "1"}}, Total: 1}
	repo.On("FindWithPagination", ctx, 0, 20).Return(expected, nil)
	svc := service.NewFormService(repo)
	got, err := svc.List(ctx, 0, 20)
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestFormService_Update(t *testing.T) {
	ctx := context.Background()

	t.Run("éxito", func(t *testing.T) {
		repo := new(MockFormRepo)
		f := &models.Form{ID: "uuid-1", Name: "Updated"}
		repo.On("Update", ctx, f).Return(nil)
		svc := service.NewFormService(repo)
		assert.NoError(t, svc.Update(ctx, f))
	})

	t.Run("error infraestructura", func(t *testing.T) {
		repo := new(MockFormRepo)
		f := &models.Form{ID: "uuid-err"}
		dbErr := errors.New("db error")
		repo.On("Update", ctx, f).Return(dbErr)
		svc := service.NewFormService(repo)
		assert.ErrorIs(t, svc.Update(ctx, f), dbErr)
	})
}
