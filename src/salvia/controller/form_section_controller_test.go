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

type MockFormSectionService struct{ mock.Mock }

func (m *MockFormSectionService) GetByID(ctx context.Context, id string) (*models.FormSection, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.FormSection), args.Error(1)
}
func (m *MockFormSectionService) List(ctx context.Context, page, limit int) (repository.PageResult[models.FormSection], error) {
	args := m.Called(ctx, page, limit)
	return args.Get(0).(repository.PageResult[models.FormSection]), args.Error(1)
}
func (m *MockFormSectionService) Create(ctx context.Context, s *models.FormSection) error {
	return m.Called(ctx, s).Error(0)
}
func (m *MockFormSectionService) Update(ctx context.Context, s *models.FormSection) error {
	return m.Called(ctx, s).Error(0)
}
func (m *MockFormSectionService) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

func TestFormSectionController_GetByID(t *testing.T) {
	t.Run("200", func(t *testing.T) {
		svc := new(MockFormSectionService)
		svc.On("GetByID", mock.Anything, "s-1").Return(&models.FormSection{ID: "s-1"}, nil)
		ctrl := controller.NewFormSectionController(svc)
		req := httptest.NewRequest(http.MethodGet, "/form-sections?id=s-1", nil)
		rec := httptest.NewRecorder()
		ctrl.GetByIDHandler(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("404", func(t *testing.T) {
		svc := new(MockFormSectionService)
		svc.On("GetByID", mock.Anything, "x").Return(nil, service.ErrFormSectionNotFound)
		ctrl := controller.NewFormSectionController(svc)
		req := httptest.NewRequest(http.MethodGet, "/form-sections?id=x", nil)
		rec := httptest.NewRecorder()
		ctrl.GetByIDHandler(rec, req)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestFormSectionController_Create(t *testing.T) {
	t.Run("201", func(t *testing.T) {
		svc := new(MockFormSectionService)
		svc.On("Create", mock.Anything, mock.AnythingOfType("*models.FormSection")).Return(nil)
		ctrl := controller.NewFormSectionController(svc)
		body, _ := json.Marshal(map[string]string{"name": "Sec A", "form_id": "f-1"})
		req := httptest.NewRequest(http.MethodPost, "/form-sections", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		ctrl.CreateHandler(rec, req)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("400 sin campos requeridos", func(t *testing.T) {
		svc := new(MockFormSectionService)
		ctrl := controller.NewFormSectionController(svc)
		body, _ := json.Marshal(map[string]string{})
		req := httptest.NewRequest(http.MethodPost, "/form-sections", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		ctrl.CreateHandler(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestFormSectionController_Delete(t *testing.T) {
	t.Run("204", func(t *testing.T) {
		svc := new(MockFormSectionService)
		svc.On("Delete", mock.Anything, "s-1").Return(nil)
		ctrl := controller.NewFormSectionController(svc)
		req := httptest.NewRequest(http.MethodDelete, "/form-sections?id=s-1", nil)
		rec := httptest.NewRecorder()
		ctrl.DeleteHandler(rec, req)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})
}
