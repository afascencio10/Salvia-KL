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

	"github.com/gin-gonic/gin"
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
func (m *MockQuestionService) ListByFormID(ctx context.Context, formID string) ([]models.Question, error) {
	args := m.Called(ctx, formID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Question), args.Error(1)
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

func newQuestionTestRouter(svc service.QuestionService) *gin.Engine {
	r := gin.New()
	ctrl := controller.NewQuestionController(svc)
	ctrl.RegisterRoutes(r.Group("/api/v1"))
	return r
}

func TestQuestionController_GetByID(t *testing.T) {
	t.Run("200", func(t *testing.T) {
		svc := new(MockQuestionService)
		svc.On("GetByID", mock.Anything, "q-1").Return(&models.Question{ID: "q-1"}, nil)
		r := newQuestionTestRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/questions/q-1", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("404", func(t *testing.T) {
		svc := new(MockQuestionService)
		svc.On("GetByID", mock.Anything, "x").Return(nil, service.ErrQuestionNotFound)
		r := newQuestionTestRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/questions/x", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestQuestionController_Create(t *testing.T) {
	t.Run("201", func(t *testing.T) {
		svc := new(MockQuestionService)
		svc.On("Create", mock.Anything, mock.AnythingOfType("*models.Question")).Return(nil)
		r := newQuestionTestRouter(svc)
		body, _ := json.Marshal(map[string]string{
			"formSectionId": "s-1", "questionTypeId": "BOOLEAN", "description": "¿Sí?",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/forms/f-1/questions", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("400 campos faltantes", func(t *testing.T) {
		svc := new(MockQuestionService)
		r := newQuestionTestRouter(svc)
		body, _ := json.Marshal(map[string]string{})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/forms/f-1/questions", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		svc.AssertNotCalled(t, "Create")
	})
}

func TestQuestionController_Delete(t *testing.T) {
	t.Run("204", func(t *testing.T) {
		svc := new(MockQuestionService)
		svc.On("Delete", mock.Anything, "q-1").Return(nil)
		r := newQuestionTestRouter(svc)
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/questions/q-1", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("404", func(t *testing.T) {
		svc := new(MockQuestionService)
		svc.On("Delete", mock.Anything, "x").Return(service.ErrQuestionNotFound)
		r := newQuestionTestRouter(svc)
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/questions/x", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}
