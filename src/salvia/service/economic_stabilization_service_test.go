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

type MockEconomicStabilizationRepo struct{ mock.Mock }

func (m *MockEconomicStabilizationRepo) Create(ctx context.Context, e *models.EconomicStabilization) error {
	return m.Called(ctx, e).Error(0)
}
func (m *MockEconomicStabilizationRepo) Update(ctx context.Context, e *models.EconomicStabilization) error {
	return m.Called(ctx, e).Error(0)
}
func (m *MockEconomicStabilizationRepo) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockEconomicStabilizationRepo) FindByID(ctx context.Context, id string) (*models.EconomicStabilization, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.EconomicStabilization), args.Error(1)
}
func (m *MockEconomicStabilizationRepo) FindWithPagination(ctx context.Context, page, size int) (repository.PageResult[models.EconomicStabilization], error) {
	args := m.Called(ctx, page, size)
	return args.Get(0).(repository.PageResult[models.EconomicStabilization]), args.Error(1)
}
func (m *MockEconomicStabilizationRepo) FindByCaseID(ctx context.Context, caseID string) ([]models.EconomicStabilization, error) {
	args := m.Called(ctx, caseID)
	return args.Get(0).([]models.EconomicStabilization), args.Error(1)
}
func (m *MockEconomicStabilizationRepo) FindByFollowUpID(ctx context.Context, fid string) ([]models.EconomicStabilization, error) {
	args := m.Called(ctx, fid)
	return args.Get(0).([]models.EconomicStabilization), args.Error(1)
}

func TestEconomicStabilizationService_GetByID(t *testing.T) {
	ctx := context.Background()

	t.Run("éxito", func(t *testing.T) {
		repo := new(MockEconomicStabilizationRepo)
		expected := &models.EconomicStabilization{ID: "es-1", Type: "SUBSIDIO"}
		repo.On("FindByID", ctx, "es-1").Return(expected, nil)
		svc := service.NewEconomicStabilizationService(repo)
		got, err := svc.GetByID(ctx, "es-1")
		assert.NoError(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(MockEconomicStabilizationRepo)
		repo.On("FindByID", ctx, "x").Return(nil, gorm.ErrRecordNotFound)
		svc := service.NewEconomicStabilizationService(repo)
		_, err := svc.GetByID(ctx, "x")
		assert.ErrorIs(t, err, service.ErrEconomicStabilizationNotFound)
	})
}

func TestEconomicStabilizationService_Create(t *testing.T) {
	ctx := context.Background()
	repo := new(MockEconomicStabilizationRepo)
	es := &models.EconomicStabilization{CaseID: "c-1", FollowUpID: "f-1", Type: "SUBSIDIO"}
	repo.On("Create", ctx, es).Return(nil)
	svc := service.NewEconomicStabilizationService(repo)
	err := svc.Create(ctx, es)
	assert.NoError(t, err)
	assert.Equal(t, "ACTIVE", es.Status)
}
