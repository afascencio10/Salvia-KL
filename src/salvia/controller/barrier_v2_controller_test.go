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

type MockBarrierV2Service struct{ mock.Mock }

func (m *MockBarrierV2Service) GetByID(ctx context.Context, id string) (*models.BarrierV2, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BarrierV2), args.Error(1)
}
func (m *MockBarrierV2Service) List(ctx context.Context, page, limit int) (repository.PageResult[models.BarrierV2], error) {
	args := m.Called(ctx, page, limit)
	return args.Get(0).(repository.PageResult[models.BarrierV2]), args.Error(1)
}
func (m *MockBarrierV2Service) Create(ctx context.Context, b *models.BarrierV2) error {
	return m.Called(ctx, b).Error(0)
}
func (m *MockBarrierV2Service) Update(ctx context.Context, b *models.BarrierV2) error {
	return m.Called(ctx, b).Error(0)
}
func (m *MockBarrierV2Service) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

func TestBarrierV2Controller_GetByID(t *testing.T) {
	t.Run("200", func(t *testing.T) {
		svc := new(MockBarrierV2Service)
		svc.On("GetByID", mock.Anything, "b-1").Return(&models.BarrierV2{ID: "b-1"}, nil)
		ctrl := controller.NewBarrierV2Controller(svc)
		req := httptest.NewRequest(http.MethodGet, "/barriers?id=b-1", nil)
		rec := httptest.NewRecorder()
		ctrl.GetByIDHandler(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("404", func(t *testing.T) {
		svc := new(MockBarrierV2Service)
		svc.On("GetByID", mock.Anything, "x").Return(nil, service.ErrBarrierV2NotFound)
		ctrl := controller.NewBarrierV2Controller(svc)
		req := httptest.NewRequest(http.MethodGet, "/barriers?id=x", nil)
		rec := httptest.NewRecorder()
		ctrl.GetByIDHandler(rec, req)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestBarrierV2Controller_Create(t *testing.T) {
	t.Run("201", func(t *testing.T) {
		svc := new(MockBarrierV2Service)
		svc.On("Create", mock.Anything, mock.AnythingOfType("*models.BarrierV2")).Return(nil)
		ctrl := controller.NewBarrierV2Controller(svc)
		body, _ := json.Marshal(map[string]string{
			"case_id": "c-1", "follow_up_id": "f-1", "sector": "JUSTICIA", "description": "Desc",
		})
		req := httptest.NewRequest(http.MethodPost, "/barriers", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		ctrl.CreateHandler(rec, req)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})
}

func TestBarrierV2Controller_Delete(t *testing.T) {
	t.Run("204", func(t *testing.T) {
		svc := new(MockBarrierV2Service)
		svc.On("Delete", mock.Anything, "b-1").Return(nil)
		ctrl := controller.NewBarrierV2Controller(svc)
		req := httptest.NewRequest(http.MethodDelete, "/barriers?id=b-1", nil)
		rec := httptest.NewRecorder()
		ctrl.DeleteHandler(rec, req)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})
}
