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

// --- Mock FormService ---

type MockFormService struct{ mock.Mock }

func (m *MockFormService) GetByID(ctx context.Context, id string) (*models.Form, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Form), args.Error(1)
}
func (m *MockFormService) List(ctx context.Context, page, limit int) (repository.PageResult[models.Form], error) {
	args := m.Called(ctx, page, limit)
	return args.Get(0).(repository.PageResult[models.Form]), args.Error(1)
}
func (m *MockFormService) Create(ctx context.Context, f *models.Form) error {
	return m.Called(ctx, f).Error(0)
}
func (m *MockFormService) Update(ctx context.Context, f *models.Form) error {
	return m.Called(ctx, f).Error(0)
}
func (m *MockFormService) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

// --- Tests GetByIDHandler ---

func TestFormController_GetByID(t *testing.T) {
	t.Run("200", func(t *testing.T) {
		svc := new(MockFormService)
		svc.On("GetByID", mock.Anything, "uuid-1").Return(&models.Form{ID: "uuid-1", Name: "F"}, nil)
		ctrl := controller.NewFormController(svc)
		req := httptest.NewRequest(http.MethodGet, "/forms?id=uuid-1", nil)
		rec := httptest.NewRecorder()
		ctrl.GetByIDHandler(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("400 sin id", func(t *testing.T) {
		svc := new(MockFormService)
		ctrl := controller.NewFormController(svc)
		req := httptest.NewRequest(http.MethodGet, "/forms", nil)
		rec := httptest.NewRecorder()
		ctrl.GetByIDHandler(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("404", func(t *testing.T) {
		svc := new(MockFormService)
		svc.On("GetByID", mock.Anything, "x").Return(nil, service.ErrFormNotFound)
		ctrl := controller.NewFormController(svc)
		req := httptest.NewRequest(http.MethodGet, "/forms?id=x", nil)
		rec := httptest.NewRecorder()
		ctrl.GetByIDHandler(rec, req)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

// --- Tests ListHandler ---

func TestFormController_List(t *testing.T) {
	t.Run("200 defaults", func(t *testing.T) {
		svc := new(MockFormService)
		svc.On("List", mock.Anything, 0, 20).Return(repository.PageResult[models.Form]{}, nil)
		ctrl := controller.NewFormController(svc)
		req := httptest.NewRequest(http.MethodGet, "/forms/list", nil)
		rec := httptest.NewRecorder()
		ctrl.ListHandler(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

// --- Tests CreateHandler ---

func TestFormController_Create(t *testing.T) {
	t.Run("201", func(t *testing.T) {
		svc := new(MockFormService)
		svc.On("Create", mock.Anything, mock.AnythingOfType("*models.Form")).Return(nil)
		ctrl := controller.NewFormController(svc)
		body, _ := json.Marshal(map[string]string{"name": "Nuevo Form"})
		req := httptest.NewRequest(http.MethodPost, "/forms", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		ctrl.CreateHandler(rec, req)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("400 sin name", func(t *testing.T) {
		svc := new(MockFormService)
		ctrl := controller.NewFormController(svc)
		body, _ := json.Marshal(map[string]string{})
		req := httptest.NewRequest(http.MethodPost, "/forms", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		ctrl.CreateHandler(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

// --- Tests DeleteHandler ---

func TestFormController_Delete(t *testing.T) {
	t.Run("204", func(t *testing.T) {
		svc := new(MockFormService)
		svc.On("Delete", mock.Anything, "uuid-1").Return(nil)
		ctrl := controller.NewFormController(svc)
		req := httptest.NewRequest(http.MethodDelete, "/forms?id=uuid-1", nil)
		rec := httptest.NewRecorder()
		ctrl.DeleteHandler(rec, req)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("404", func(t *testing.T) {
		svc := new(MockFormService)
		svc.On("Delete", mock.Anything, "x").Return(service.ErrFormNotFound)
		ctrl := controller.NewFormController(svc)
		req := httptest.NewRequest(http.MethodDelete, "/forms?id=x", nil)
		rec := httptest.NewRecorder()
		ctrl.DeleteHandler(rec, req)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}
