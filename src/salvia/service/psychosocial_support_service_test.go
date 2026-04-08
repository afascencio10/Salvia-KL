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

type MockPsychosocialSupportRepo struct{ mock.Mock }

func (m *MockPsychosocialSupportRepo) Create(ctx context.Context, e *models.PsychosocialSupport) error {
	return m.Called(ctx, e).Error(0)
}
func (m *MockPsychosocialSupportRepo) Update(ctx context.Context, e *models.PsychosocialSupport) error {
	return m.Called(ctx, e).Error(0)
}
func (m *MockPsychosocialSupportRepo) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockPsychosocialSupportRepo) FindByID(ctx context.Context, id string) (*models.PsychosocialSupport, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PsychosocialSupport), args.Error(1)
}
func (m *MockPsychosocialSupportRepo) FindWithPagination(ctx context.Context, page, size int) (repository.PageResult[models.PsychosocialSupport], error) {
	args := m.Called(ctx, page, size)
	return args.Get(0).(repository.PageResult[models.PsychosocialSupport]), args.Error(1)
}
func (m *MockPsychosocialSupportRepo) FindByCaseID(ctx context.Context, caseID string) ([]models.PsychosocialSupport, error) {
	args := m.Called(ctx, caseID)
	return args.Get(0).([]models.PsychosocialSupport), args.Error(1)
}
func (m *MockPsychosocialSupportRepo) FindByFollowUpID(ctx context.Context, fid string) ([]models.PsychosocialSupport, error) {
	args := m.Called(ctx, fid)
	return args.Get(0).([]models.PsychosocialSupport), args.Error(1)
}

func TestPsychosocialSupportService_GetByID(t *testing.T) {
	ctx := context.Background()

	t.Run("éxito", func(t *testing.T) {
		repo := new(MockPsychosocialSupportRepo)
		expected := &models.PsychosocialSupport{ID: "ps-1", Type: "TERAPIA"}
		repo.On("FindByID", ctx, "ps-1").Return(expected, nil)
		svc := service.NewPsychosocialSupportService(repo)
		got, err := svc.GetByID(ctx, "ps-1")
		assert.NoError(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(MockPsychosocialSupportRepo)
		repo.On("FindByID", ctx, "x").Return(nil, gorm.ErrRecordNotFound)
		svc := service.NewPsychosocialSupportService(repo)
		_, err := svc.GetByID(ctx, "x")
		assert.ErrorIs(t, err, service.ErrPsychosocialSupportNotFound)
	})
}

func TestPsychosocialSupportService_Create(t *testing.T) {
	ctx := context.Background()
	repo := new(MockPsychosocialSupportRepo)
	ps := &models.PsychosocialSupport{CaseID: "c-1", FollowUpID: "f-1", Type: "TERAPIA"}
	repo.On("Create", ctx, ps).Return(nil)
	svc := service.NewPsychosocialSupportService(repo)
	err := svc.Create(ctx, ps)
	assert.NoError(t, err)
	assert.Equal(t, "ACTIVE", ps.Status)
}
