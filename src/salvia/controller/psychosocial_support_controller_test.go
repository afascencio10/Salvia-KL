package controller_test

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"bitsflow/salvia/controller"
	"bitsflow/salvia/service"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPsychosocialSupportService struct{ mock.Mock }

func (m *MockPsychosocialSupportService) GetByID(ctx context.Context, id string) (*models.PsychosocialSupport, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PsychosocialSupport), args.Error(1)
}
func (m *MockPsychosocialSupportService) List(ctx context.Context, page, limit int) (repository.PageResult[models.PsychosocialSupport], error) {
	args := m.Called(ctx, page, limit)
	return args.Get(0).(repository.PageResult[models.PsychosocialSupport]), args.Error(1)
}
func (m *MockPsychosocialSupportService) Create(ctx context.Context, ps *models.PsychosocialSupport) error {
	return m.Called(ctx, ps).Error(0)
}
func (m *MockPsychosocialSupportService) Update(ctx context.Context, ps *models.PsychosocialSupport) error {
	return m.Called(ctx, ps).Error(0)
}
func (m *MockPsychosocialSupportService) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

func TestPsychosocialSupportController_GetByID(t *testing.T) {
	t.Run("200", func(t *testing.T) {
		svc := new(MockPsychosocialSupportService)
		svc.On("GetByID", mock.Anything, "ps-1").Return(&models.PsychosocialSupport{ID: "ps-1"}, nil)
		ctrl := controller.NewPsychosocialSupportController(svc)
		req := httptest.NewRequest(http.MethodGet, "/psychosocial?id=ps-1", nil)
		rec := httptest.NewRecorder()
		ctrl.GetByIDHandler(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("404", func(t *testing.T) {
		svc := new(MockPsychosocialSupportService)
		svc.On("GetByID", mock.Anything, "x").Return(nil, service.ErrPsychosocialSupportNotFound)
		ctrl := controller.NewPsychosocialSupportController(svc)
		req := httptest.NewRequest(http.MethodGet, "/psychosocial?id=x", nil)
		rec := httptest.NewRecorder()
		ctrl.GetByIDHandler(rec, req)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestPsychosocialSupportController_Create(t *testing.T) {
	t.Run("201", func(t *testing.T) {
		svc := new(MockPsychosocialSupportService)
		svc.On("Create", mock.Anything, mock.AnythingOfType("*models.PsychosocialSupport")).Return(nil)
		ctrl := controller.NewPsychosocialSupportController(svc)
		body, _ := json.Marshal(map[string]string{"case_id": "c-1", "follow_up_id": "f-1", "type": "TERAPIA"})
		req := httptest.NewRequest(http.MethodPost, "/psychosocial", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		ctrl.CreateHandler(rec, req)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})
}

func TestPsychosocialSupportController_Delete(t *testing.T) {
	t.Run("204", func(t *testing.T) {
		svc := new(MockPsychosocialSupportService)
		svc.On("Delete", mock.Anything, "ps-1").Return(nil)
		ctrl := controller.NewPsychosocialSupportController(svc)
		req := httptest.NewRequest(http.MethodDelete, "/psychosocial?id=ps-1", nil)
		rec := httptest.NewRecorder()
		ctrl.DeleteHandler(rec, req)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})
}
