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

type MockBarrierV2Repo struct{ mock.Mock }

func (m *MockBarrierV2Repo) Create(ctx context.Context, e *models.BarrierV2) error {
	return m.Called(ctx, e).Error(0)
}
func (m *MockBarrierV2Repo) Update(ctx context.Context, e *models.BarrierV2) error {
	return m.Called(ctx, e).Error(0)
}
func (m *MockBarrierV2Repo) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockBarrierV2Repo) FindByID(ctx context.Context, id string) (*models.BarrierV2, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BarrierV2), args.Error(1)
}
func (m *MockBarrierV2Repo) FindWithPagination(ctx context.Context, page, size int) (repository.PageResult[models.BarrierV2], error) {
	args := m.Called(ctx, page, size)
	return args.Get(0).(repository.PageResult[models.BarrierV2]), args.Error(1)
}
func (m *MockBarrierV2Repo) FindByCaseID(ctx context.Context, caseID string) ([]models.BarrierV2, error) {
	args := m.Called(ctx, caseID)
	return args.Get(0).([]models.BarrierV2), args.Error(1)
}
func (m *MockBarrierV2Repo) FindByFollowUpID(ctx context.Context, fid string) ([]models.BarrierV2, error) {
	args := m.Called(ctx, fid)
	return args.Get(0).([]models.BarrierV2), args.Error(1)
}

func TestBarrierV2Service_GetByID(t *testing.T) {
	ctx := context.Background()

	t.Run("éxito", func(t *testing.T) {
		repo := new(MockBarrierV2Repo)
		expected := &models.BarrierV2{ID: "b-1", Status: "OPEN"}
		repo.On("FindByID", ctx, "b-1").Return(expected, nil)
		svc := service.NewBarrierV2Service(repo)
		got, err := svc.GetByID(ctx, "b-1")
		assert.NoError(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(MockBarrierV2Repo)
		repo.On("FindByID", ctx, "x").Return(nil, gorm.ErrRecordNotFound)
		svc := service.NewBarrierV2Service(repo)
		_, err := svc.GetByID(ctx, "x")
		assert.ErrorIs(t, err, service.ErrBarrierV2NotFound)
	})
}

func TestBarrierV2Service_Create(t *testing.T) {
	ctx := context.Background()
	repo := new(MockBarrierV2Repo)
	b := &models.BarrierV2{CaseID: "c-1", FollowUpID: "f-1", Sector: "JUSTICIA", Description: "Desc"}
	repo.On("Create", ctx, b).Return(nil)
	svc := service.NewBarrierV2Service(repo)
	err := svc.Create(ctx, b)
	assert.NoError(t, err)
	assert.Equal(t, "OPEN", b.Status)
}
