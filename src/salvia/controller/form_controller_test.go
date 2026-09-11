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
func (m *MockFormService) Create(ctx context.Context, input service.CreateFormInput) (*models.Form, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Form), args.Error(1)
}
func (m *MockFormService) Update(ctx context.Context, id string, input service.UpdateFormInput) (*models.Form, error) {
	args := m.Called(ctx, id, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Form), args.Error(1)
}
func (m *MockFormService) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockFormService) GetFormStructure(ctx context.Context, formID string) (*service.FormStructure, error) {
	args := m.Called(ctx, formID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.FormStructure), args.Error(1)
}
func (m *MockFormService) GetFormSubmission(ctx context.Context, submissionID string) (*service.SubmissionStructure, error) {
	args := m.Called(ctx, submissionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.SubmissionStructure), args.Error(1)
}
func (m *MockFormService) GetFormSubmissionStructured(ctx context.Context, formID, submissionID string) (*service.FormSubmissionStructured, error) {
	args := m.Called(ctx, formID, submissionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.FormSubmissionStructured), args.Error(1)
}
func (m *MockFormService) LoadForm(ctx context.Context, formID, submissionID string, formState map[string]interface{}) (*service.LoadFormResult, error) {
	args := m.Called(ctx, formID, submissionID, formState)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.LoadFormResult), args.Error(1)
}
func (m *MockFormService) SaveSection(ctx context.Context, input service.SaveSectionInput) (*service.LoadFormResult, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.LoadFormResult), args.Error(1)
}
func (m *MockFormService) OnEndFormSubmission(ctx context.Context, formID, submissionID, actorID string) error {
	return m.Called(ctx, formID, submissionID, actorID).Error(0)
}
func (m *MockFormService) OnSectionUpdate(ctx context.Context, formID, submissionID, actorID string) error {
	return m.Called(ctx, formID, submissionID, actorID).Error(0)
}
func (m *MockFormService) TestFunction(ctx context.Context, fn, id, submissionID string) (interface{}, error) {
	args := m.Called(ctx, fn, id, submissionID)
	return args.Get(0), args.Error(1)
}
func (m *MockFormService) ProcessVictimCaseEntidadEntries(ctx context.Context, caseID, actorID string, acudidas, activacion []service.VictimCaseEntidadEntryInput) error {
	return m.Called(ctx, caseID, actorID, acudidas, activacion).Error(0)
}

// ── helper: construye router de test ─────────────────────────────────────────

func newFormTestRouter(svc service.FormService) *gin.Engine {
	r := gin.New()
	ctrl := controller.NewFormController(svc, nil)
	ctrl.RegisterRoutes(r.Group("/api/v1"))
	return r
}

// --- Tests GetByID ---

func TestFormController_GetByID(t *testing.T) {
	t.Run("200", func(t *testing.T) {
		svc := new(MockFormService)
		svc.On("GetByID", mock.Anything, "uuid-1").Return(&models.Form{ID: "uuid-1", Name: "F"}, nil)
		r := newFormTestRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/forms/uuid-1", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("404", func(t *testing.T) {
		svc := new(MockFormService)
		svc.On("GetByID", mock.Anything, "x").Return(nil, service.ErrFormNotFound)
		r := newFormTestRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/forms/x", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

// --- Tests List ---

func TestFormController_List(t *testing.T) {
	t.Run("200 defaults", func(t *testing.T) {
		svc := new(MockFormService)
		svc.On("List", mock.Anything, 0, 20).Return(repository.PageResult[models.Form]{}, nil)
		r := newFormTestRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/forms", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

// --- Tests Create ---

func TestFormController_Create(t *testing.T) {
	t.Run("201", func(t *testing.T) {
		svc := new(MockFormService)
		svc.On("Create", mock.Anything, service.CreateFormInput{Name: "Nuevo Form", Description: "Desc"}).
			Return(&models.Form{ID: "uuid-new", Name: "Nuevo Form"}, nil)
		r := newFormTestRouter(svc)
		body, _ := json.Marshal(map[string]string{"name": "Nuevo Form", "description": "Desc"})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/forms", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("400 sin description", func(t *testing.T) {
		svc := new(MockFormService)
		r := newFormTestRouter(svc)
		body, _ := json.Marshal(map[string]string{"name": "Nuevo Form"})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/forms", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		svc.AssertNotCalled(t, "Create")
	})
}

// --- Tests Delete ---

func TestFormController_Delete(t *testing.T) {
	t.Run("204", func(t *testing.T) {
		svc := new(MockFormService)
		svc.On("Delete", mock.Anything, "uuid-1").Return(nil)
		r := newFormTestRouter(svc)
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/forms/uuid-1", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("404", func(t *testing.T) {
		svc := new(MockFormService)
		svc.On("Delete", mock.Anything, "x").Return(service.ErrFormNotFound)
		r := newFormTestRouter(svc)
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/forms/x", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}
