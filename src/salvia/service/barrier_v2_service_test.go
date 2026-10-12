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
func (m *MockBarrierV2Repo) UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error {
	return m.Called(ctx, id, fields).Error(0)
}
func (m *MockBarrierV2Repo) FindActiveByCaseID(ctx context.Context, caseID string) ([]models.BarrierV2, error) {
	args := m.Called(ctx, caseID)
	return args.Get(0).([]models.BarrierV2), args.Error(1)
}
func (m *MockBarrierV2Repo) FindByIDs(ctx context.Context, ids []string) ([]models.BarrierV2, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]models.BarrierV2), args.Error(1)
}
func (m *MockBarrierV2Repo) FindActiveByTeamContactIDs(ctx context.Context, teamContactIDs []string) ([]models.BarrierV2, error) {
	args := m.Called(ctx, teamContactIDs)
	return args.Get(0).([]models.BarrierV2), args.Error(1)
}
func (m *MockBarrierV2Repo) UpdateStatus(ctx context.Context, id string, status string) error {
	return m.Called(ctx, id, status).Error(0)
}
func (m *MockBarrierV2Repo) FindByCreatedByIDWithRelations(ctx context.Context, createdByID string) ([]models.BarrierV2WithRelations, error) {
	args := m.Called(ctx, createdByID)
	return args.Get(0).([]models.BarrierV2WithRelations), args.Error(1)
}
func (m *MockBarrierV2Repo) FindByDepartmentWithRelations(ctx context.Context, departmentID string, filters repository.DepartmentBarrierFilters) ([]models.BarrierV2DepartmentItem, int64, error) {
	args := m.Called(ctx, departmentID, filters)
	return args.Get(0).([]models.BarrierV2DepartmentItem), args.Get(1).(int64), args.Error(2)
}
func (m *MockBarrierV2Repo) FindMineWithRelations(ctx context.Context, userTerritoryLinkID string) ([]models.BarrierV2MineItem, error) {
	args := m.Called(ctx, userTerritoryLinkID)
	return args.Get(0).([]models.BarrierV2MineItem), args.Error(1)
}
func (m *MockBarrierV2Repo) GetDB() *gorm.DB {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*gorm.DB)
}

func TestBarrierV2Service_GetByID(t *testing.T) {
	ctx := context.Background()

	t.Run("éxito", func(t *testing.T) {
		repo := new(MockBarrierV2Repo)
		expected := &models.BarrierV2{ID: "b-1", Status: "OPEN"}
		repo.On("FindByID", ctx, "b-1").Return(expected, nil)
		svc := service.NewBarrierV2Service(repo, nil, nil)
		got, err := svc.GetByID(ctx, "b-1")
		assert.NoError(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(MockBarrierV2Repo)
		repo.On("FindByID", ctx, "x").Return(nil, gorm.ErrRecordNotFound)
		svc := service.NewBarrierV2Service(repo, nil, nil)
		_, err := svc.GetByID(ctx, "x")
		assert.ErrorIs(t, err, service.ErrBarrierV2NotFound)
	})
}

func TestBarrierV2Service_Create(t *testing.T) {
	ctx := context.Background()
	repo := new(MockBarrierV2Repo)
	followUpID := "f-1"
	b := &models.BarrierV2{CaseID: "c-1", FollowUpID: &followUpID, Sector: "JUSTICIA", Description: "Desc"}
	repo.On("Create", ctx, b).Return(nil)
	svc := service.NewBarrierV2Service(repo, nil, nil)
	err := svc.Create(ctx, b)
	assert.NoError(t, err)
	assert.Equal(t, "OPEN", b.Status)
}
