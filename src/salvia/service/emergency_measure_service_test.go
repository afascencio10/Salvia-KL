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

type MockEmergencyMeasureRepo struct{ mock.Mock }

func (m *MockEmergencyMeasureRepo) Create(ctx context.Context, e *models.EmergencyMeasure) error {
	return m.Called(ctx, e).Error(0)
}
func (m *MockEmergencyMeasureRepo) Update(ctx context.Context, e *models.EmergencyMeasure) error {
	return m.Called(ctx, e).Error(0)
}
func (m *MockEmergencyMeasureRepo) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockEmergencyMeasureRepo) FindByID(ctx context.Context, id string) (*models.EmergencyMeasure, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.EmergencyMeasure), args.Error(1)
}
func (m *MockEmergencyMeasureRepo) FindWithPagination(ctx context.Context, page, size int) (repository.PageResult[models.EmergencyMeasure], error) {
	args := m.Called(ctx, page, size)
	return args.Get(0).(repository.PageResult[models.EmergencyMeasure]), args.Error(1)
}
func (m *MockEmergencyMeasureRepo) FindByCaseID(ctx context.Context, caseID string) ([]models.EmergencyMeasure, error) {
	args := m.Called(ctx, caseID)
	return args.Get(0).([]models.EmergencyMeasure), args.Error(1)
}
func (m *MockEmergencyMeasureRepo) FindByFollowUpID(ctx context.Context, fid string) ([]models.EmergencyMeasure, error) {
	args := m.Called(ctx, fid)
	return args.Get(0).([]models.EmergencyMeasure), args.Error(1)
}

func TestEmergencyMeasureService_GetByID(t *testing.T) {
	ctx := context.Background()

	t.Run("éxito", func(t *testing.T) {
		repo := new(MockEmergencyMeasureRepo)
		expected := &models.EmergencyMeasure{ID: "em-1", Type: "PROTECCION"}
		repo.On("FindByID", ctx, "em-1").Return(expected, nil)
		svc := service.NewEmergencyMeasureService(repo)
		got, err := svc.GetByID(ctx, "em-1")
		assert.NoError(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(MockEmergencyMeasureRepo)
		repo.On("FindByID", ctx, "x").Return(nil, gorm.ErrRecordNotFound)
		svc := service.NewEmergencyMeasureService(repo)
		_, err := svc.GetByID(ctx, "x")
		assert.ErrorIs(t, err, service.ErrEmergencyMeasureNotFound)
	})
}

func TestEmergencyMeasureService_Create(t *testing.T) {
	ctx := context.Background()
	repo := new(MockEmergencyMeasureRepo)
	em := &models.EmergencyMeasure{CaseID: "c-1", FollowUpID: "f-1", Type: "PROTECCION"}
	repo.On("Create", ctx, em).Return(nil)
	svc := service.NewEmergencyMeasureService(repo)
	err := svc.Create(ctx, em)
	assert.NoError(t, err)
	assert.Equal(t, "ACTIVE", em.Status)
}
