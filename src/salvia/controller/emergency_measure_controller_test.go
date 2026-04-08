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

type MockEmergencyMeasureService struct{ mock.Mock }

func (m *MockEmergencyMeasureService) GetByID(ctx context.Context, id string) (*models.EmergencyMeasure, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.EmergencyMeasure), args.Error(1)
}
func (m *MockEmergencyMeasureService) List(ctx context.Context, page, limit int) (repository.PageResult[models.EmergencyMeasure], error) {
	args := m.Called(ctx, page, limit)
	return args.Get(0).(repository.PageResult[models.EmergencyMeasure]), args.Error(1)
}
func (m *MockEmergencyMeasureService) Create(ctx context.Context, em *models.EmergencyMeasure) error {
	return m.Called(ctx, em).Error(0)
}
func (m *MockEmergencyMeasureService) Update(ctx context.Context, em *models.EmergencyMeasure) error {
	return m.Called(ctx, em).Error(0)
}
func (m *MockEmergencyMeasureService) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

func TestEmergencyMeasureController_GetByID(t *testing.T) {
	t.Run("200", func(t *testing.T) {
		svc := new(MockEmergencyMeasureService)
		svc.On("GetByID", mock.Anything, "em-1").Return(&models.EmergencyMeasure{ID: "em-1"}, nil)
		ctrl := controller.NewEmergencyMeasureController(svc)
		req := httptest.NewRequest(http.MethodGet, "/emergency-measures?id=em-1", nil)
		rec := httptest.NewRecorder()
		ctrl.GetByIDHandler(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("404", func(t *testing.T) {
		svc := new(MockEmergencyMeasureService)
		svc.On("GetByID", mock.Anything, "x").Return(nil, service.ErrEmergencyMeasureNotFound)
		ctrl := controller.NewEmergencyMeasureController(svc)
		req := httptest.NewRequest(http.MethodGet, "/emergency-measures?id=x", nil)
		rec := httptest.NewRecorder()
		ctrl.GetByIDHandler(rec, req)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestEmergencyMeasureController_Create(t *testing.T) {
	t.Run("201", func(t *testing.T) {
		svc := new(MockEmergencyMeasureService)
		svc.On("Create", mock.Anything, mock.AnythingOfType("*models.EmergencyMeasure")).Return(nil)
		ctrl := controller.NewEmergencyMeasureController(svc)
		body, _ := json.Marshal(map[string]string{"case_id": "c-1", "follow_up_id": "f-1", "type": "PROTECCION"})
		req := httptest.NewRequest(http.MethodPost, "/emergency-measures", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		ctrl.CreateHandler(rec, req)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})
}

func TestEmergencyMeasureController_Delete(t *testing.T) {
	t.Run("204", func(t *testing.T) {
		svc := new(MockEmergencyMeasureService)
		svc.On("Delete", mock.Anything, "em-1").Return(nil)
		ctrl := controller.NewEmergencyMeasureController(svc)
		req := httptest.NewRequest(http.MethodDelete, "/emergency-measures?id=em-1", nil)
		rec := httptest.NewRecorder()
		ctrl.DeleteHandler(rec, req)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})
}
