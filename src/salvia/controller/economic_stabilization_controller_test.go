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

type MockEconomicStabilizationService struct{ mock.Mock }

func (m *MockEconomicStabilizationService) GetByID(ctx context.Context, id string) (*models.EconomicStabilization, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.EconomicStabilization), args.Error(1)
}
func (m *MockEconomicStabilizationService) List(ctx context.Context, page, limit int) (repository.PageResult[models.EconomicStabilization], error) {
	args := m.Called(ctx, page, limit)
	return args.Get(0).(repository.PageResult[models.EconomicStabilization]), args.Error(1)
}
func (m *MockEconomicStabilizationService) Create(ctx context.Context, es *models.EconomicStabilization) error {
	return m.Called(ctx, es).Error(0)
}
func (m *MockEconomicStabilizationService) Update(ctx context.Context, es *models.EconomicStabilization) error {
	return m.Called(ctx, es).Error(0)
}
func (m *MockEconomicStabilizationService) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

func TestEconomicStabilizationController_GetByID(t *testing.T) {
	t.Run("200", func(t *testing.T) {
		svc := new(MockEconomicStabilizationService)
		svc.On("GetByID", mock.Anything, "es-1").Return(&models.EconomicStabilization{ID: "es-1"}, nil)
		ctrl := controller.NewEconomicStabilizationController(svc)
		req := httptest.NewRequest(http.MethodGet, "/economic?id=es-1", nil)
		rec := httptest.NewRecorder()
		ctrl.GetByIDHandler(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("404", func(t *testing.T) {
		svc := new(MockEconomicStabilizationService)
		svc.On("GetByID", mock.Anything, "x").Return(nil, service.ErrEconomicStabilizationNotFound)
		ctrl := controller.NewEconomicStabilizationController(svc)
		req := httptest.NewRequest(http.MethodGet, "/economic?id=x", nil)
		rec := httptest.NewRecorder()
		ctrl.GetByIDHandler(rec, req)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestEconomicStabilizationController_Create(t *testing.T) {
	t.Run("201", func(t *testing.T) {
		svc := new(MockEconomicStabilizationService)
		svc.On("Create", mock.Anything, mock.AnythingOfType("*models.EconomicStabilization")).Return(nil)
		ctrl := controller.NewEconomicStabilizationController(svc)
		body, _ := json.Marshal(map[string]string{"case_id": "c-1", "follow_up_id": "f-1", "type": "SUBSIDIO"})
		req := httptest.NewRequest(http.MethodPost, "/economic", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		ctrl.CreateHandler(rec, req)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})
}

func TestEconomicStabilizationController_Delete(t *testing.T) {
	t.Run("204", func(t *testing.T) {
		svc := new(MockEconomicStabilizationService)
		svc.On("Delete", mock.Anything, "es-1").Return(nil)
		ctrl := controller.NewEconomicStabilizationController(svc)
		req := httptest.NewRequest(http.MethodDelete, "/economic?id=es-1", nil)
		rec := httptest.NewRecorder()
		ctrl.DeleteHandler(rec, req)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})
}
