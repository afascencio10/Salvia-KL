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

type MockQuestionService struct{ mock.Mock }

func (m *MockQuestionService) GetByID(ctx context.Context, id string) (*models.Question, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Question), args.Error(1)
}
func (m *MockQuestionService) List(ctx context.Context, page, limit int) (repository.PageResult[models.Question], error) {
	args := m.Called(ctx, page, limit)
	return args.Get(0).(repository.PageResult[models.Question]), args.Error(1)
}
func (m *MockQuestionService) Create(ctx context.Context, q *models.Question) error {
	return m.Called(ctx, q).Error(0)
}
func (m *MockQuestionService) Update(ctx context.Context, q *models.Question) error {
	return m.Called(ctx, q).Error(0)
}
func (m *MockQuestionService) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

func TestQuestionController_GetByID(t *testing.T) {
	t.Run("200", func(t *testing.T) {
		svc := new(MockQuestionService)
		svc.On("GetByID", mock.Anything, "q-1").Return(&models.Question{ID: "q-1"}, nil)
		ctrl := controller.NewQuestionController(svc)
		req := httptest.NewRequest(http.MethodGet, "/questions?id=q-1", nil)
		rec := httptest.NewRecorder()
		ctrl.GetByIDHandler(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("404", func(t *testing.T) {
		svc := new(MockQuestionService)
		svc.On("GetByID", mock.Anything, "x").Return(nil, service.ErrQuestionNotFound)
		ctrl := controller.NewQuestionController(svc)
		req := httptest.NewRequest(http.MethodGet, "/questions?id=x", nil)
		rec := httptest.NewRecorder()
		ctrl.GetByIDHandler(rec, req)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestQuestionController_Create(t *testing.T) {
	t.Run("201", func(t *testing.T) {
		svc := new(MockQuestionService)
		svc.On("Create", mock.Anything, mock.AnythingOfType("*models.Question")).Return(nil)
		ctrl := controller.NewQuestionController(svc)
		body, _ := json.Marshal(map[string]string{
			"form_id": "f-1", "form_section_id": "s-1",
			"question_type_id": "BOOLEAN", "description": "¿Sí?",
		})
		req := httptest.NewRequest(http.MethodPost, "/questions", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		ctrl.CreateHandler(rec, req)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("400 campos faltantes", func(t *testing.T) {
		svc := new(MockQuestionService)
		ctrl := controller.NewQuestionController(svc)
		body, _ := json.Marshal(map[string]string{"form_id": "f-1"})
		req := httptest.NewRequest(http.MethodPost, "/questions", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		ctrl.CreateHandler(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestQuestionController_Delete(t *testing.T) {
	t.Run("204", func(t *testing.T) {
		svc := new(MockQuestionService)
		svc.On("Delete", mock.Anything, "q-1").Return(nil)
		ctrl := controller.NewQuestionController(svc)
		req := httptest.NewRequest(http.MethodDelete, "/questions?id=q-1", nil)
		rec := httptest.NewRecorder()
		ctrl.DeleteHandler(rec, req)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})
}
