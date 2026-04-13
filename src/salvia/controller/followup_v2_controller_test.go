package controller_test

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"bitsflow/salvia/controller"
	"bitsflow/salvia/service"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func init() { gin.SetMode(gin.TestMode) }

// ── Mock del servicio ─────────────────────────────────────────────────────────

type MockFollowUpV2Service struct{ mock.Mock }

func (m *MockFollowUpV2Service) GetFollowUpByID(ctx context.Context, id string) (*models.FollowUpV2, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.FollowUpV2), args.Error(1)
}

func (m *MockFollowUpV2Service) GetPaginatedFollowUps(ctx context.Context, page, limit int) (repository.PageResult[models.FollowUpV2], error) {
	args := m.Called(ctx, page, limit)
	return args.Get(0).(repository.PageResult[models.FollowUpV2]), args.Error(1)
}

func (m *MockFollowUpV2Service) GetByCaseID(ctx context.Context, caseID string) ([]models.FollowUpV2, error) {
	args := m.Called(ctx, caseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.FollowUpV2), args.Error(1)
}

func (m *MockFollowUpV2Service) GenerateOrRecalculate(ctx context.Context, caseID string, input service.GenerateCalendarInput) ([]models.FollowUpV2, error) {
	args := m.Called(ctx, caseID, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.FollowUpV2), args.Error(1)
}

// ── helper: construye router de test ─────────────────────────────────────────

func newTestRouter(svc service.FollowUpV2Service) *gin.Engine {
	r := gin.New()
	ctrl := controller.NewFollowUpV2Controller(svc)
	ctrl.RegisterRoutes(r.Group("/api/v1"))
	return r
}

// ── Tests GetByID (migrados de net/http a Gin) ────────────────────────────────

func TestGetByIDHandler(t *testing.T) {
	t.Run("200 - registro encontrado", func(t *testing.T) {
		svcMock := new(MockFollowUpV2Service)
		fu := &models.FollowUpV2{ID: "uuid-abc", CaseID: "case-1", Status: "PENDIENTE"}
		svcMock.On("GetFollowUpByID", mock.Anything, "uuid-abc").Return(fu, nil)

		r := newTestRouter(svcMock)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/cases/case-1/follow-ups/by-id?id=uuid-abc", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var body models.FollowUpV2
		_ = json.NewDecoder(rec.Body).Decode(&body)
		assert.Equal(t, "uuid-abc", body.ID)
		svcMock.AssertExpectations(t)
	})

	t.Run("400 - falta parámetro id", func(t *testing.T) {
		svcMock := new(MockFollowUpV2Service)
		r := newTestRouter(svcMock)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/cases/case-1/follow-ups/by-id", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		svcMock.AssertNotCalled(t, "GetFollowUpByID")
	})

	t.Run("404 - registro no encontrado", func(t *testing.T) {
		svcMock := new(MockFollowUpV2Service)
		svcMock.On("GetFollowUpByID", mock.Anything, "no-existe").Return(nil, service.ErrFollowUpNotFound)

		r := newTestRouter(svcMock)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/cases/case-1/follow-ups/by-id?id=no-existe", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		svcMock.AssertExpectations(t)
	})

	t.Run("500 - error de infraestructura", func(t *testing.T) {
		svcMock := new(MockFollowUpV2Service)
		svcMock.On("GetFollowUpByID", mock.Anything, "uuid-err").Return(nil, errors.New("db down"))

		r := newTestRouter(svcMock)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/cases/case-1/follow-ups/by-id?id=uuid-err", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		svcMock.AssertExpectations(t)
	})
}

// ── Tests List (migrados de net/http a Gin) ───────────────────────────────────

func TestListHandler(t *testing.T) {
	t.Run("200 - lista paginada con defaults", func(t *testing.T) {
		svcMock := new(MockFollowUpV2Service)
		expected := repository.PageResult[models.FollowUpV2]{
			Items: []models.FollowUpV2{{ID: "uuid-1"}}, Total: 1, Page: 0, PageSize: 20,
		}
		svcMock.On("GetPaginatedFollowUps", mock.Anything, 0, 20).Return(expected, nil)

		r := newTestRouter(svcMock)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/cases/case-1/follow-ups/list", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		svcMock.AssertExpectations(t)
	})

	t.Run("200 - page y limit desde query params", func(t *testing.T) {
		svcMock := new(MockFollowUpV2Service)
		expected := repository.PageResult[models.FollowUpV2]{Items: []models.FollowUpV2{}, Total: 0, Page: 2, PageSize: 5}
		svcMock.On("GetPaginatedFollowUps", mock.Anything, 2, 5).Return(expected, nil)

		r := newTestRouter(svcMock)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/cases/case-1/follow-ups/list?page=2&limit=5", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		svcMock.AssertExpectations(t)
	})

	t.Run("500 - error del servicio", func(t *testing.T) {
		svcMock := new(MockFollowUpV2Service)
		empty := repository.PageResult[models.FollowUpV2]{}
		svcMock.On("GetPaginatedFollowUps", mock.Anything, 0, 20).Return(empty, errors.New("timeout"))

		r := newTestRouter(svcMock)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/cases/case-1/follow-ups/list", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		svcMock.AssertExpectations(t)
	})
}

// ── Tests GetCalendar (nuevos HU-027) ─────────────────────────────────────────

func TestGetCalendar_200_WithResults(t *testing.T) {
	svcMock := new(MockFollowUpV2Service)
	items := []models.FollowUpV2{
		{ID: "fu-1", CaseID: "case-abc"},
		{ID: "fu-2", CaseID: "case-abc"},
	}
	svcMock.On("GetByCaseID", mock.Anything, "case-abc").Return(items, nil)

	r := newTestRouter(svcMock)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/cases/case-abc/follow-ups", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var body []models.FollowUpV2
	_ = json.NewDecoder(rec.Body).Decode(&body)
	assert.Len(t, body, 2)
	svcMock.AssertExpectations(t)
}

func TestGetCalendar_404_Empty(t *testing.T) {
	svcMock := new(MockFollowUpV2Service)
	svcMock.On("GetByCaseID", mock.Anything, "case-vacio").Return(nil, service.ErrFollowUpCaseEmpty)

	r := newTestRouter(svcMock)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/cases/case-vacio/follow-ups", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	svcMock.AssertExpectations(t)
}

// ── Tests GenerateCalendar (nuevos HU-027) ────────────────────────────────────

func TestGenerateCalendar_201_Created(t *testing.T) {
	svcMock := new(MockFollowUpV2Service)
	input := service.GenerateCalendarInput{RiskLevel: 3, AgentID: "agent-1", Team: "Equipo A"}
	today := time.Now().Truncate(24 * time.Hour)
	created := []models.FollowUpV2{
		{ID: "fu-1", CaseID: "case-xyz", ScheduledDate: today.AddDate(0, 0, 1)},
		{ID: "fu-2", CaseID: "case-xyz", ScheduledDate: today.AddDate(0, 0, 3)},
		{ID: "fu-3", CaseID: "case-xyz", ScheduledDate: today.AddDate(0, 0, 15)},
		{ID: "fu-4", CaseID: "case-xyz", ScheduledDate: today.AddDate(0, 0, 30)},
	}
	svcMock.On("GenerateOrRecalculate", mock.Anything, "case-xyz", input).Return(created, nil)

	r := newTestRouter(svcMock)
	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/cases/case-xyz/follow-ups/generate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	var result []models.FollowUpV2
	_ = json.NewDecoder(rec.Body).Decode(&result)
	assert.Len(t, result, 4)
	svcMock.AssertExpectations(t)
}

func TestGenerateCalendar_400_InvalidRiskLevel(t *testing.T) {
	svcMock := new(MockFollowUpV2Service)
	r := newTestRouter(svcMock)

	// risk_level = 5 → falla el binding (max=4)
	body, _ := json.Marshal(map[string]interface{}{"risk_level": 5, "agent_id": "agent-1"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/cases/case-xyz/follow-ups/generate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	svcMock.AssertNotCalled(t, "GenerateOrRecalculate")
}
