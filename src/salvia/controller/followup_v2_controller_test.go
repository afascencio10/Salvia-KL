package controller_test

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"bitsflow/salvia/controller"
	"bitsflow/salvia/service"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ---------------------------------------------------------------------------
// Mock del servicio
// ---------------------------------------------------------------------------

type MockFollowUpV2Service struct {
	mock.Mock
}

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

// ---------------------------------------------------------------------------
// Tests de GetByIDHandler
// ---------------------------------------------------------------------------

func TestGetByIDHandler(t *testing.T) {
	t.Run("200 - registro encontrado", func(t *testing.T) {
		svcMock := new(MockFollowUpV2Service)
		fu := &models.FollowUpV2{ID: "uuid-abc", CaseID: "case-1", Status: "open"}
		svcMock.On("GetFollowUpByID", mock.Anything, "uuid-abc").Return(fu, nil)

		ctrl := controller.NewFollowUpV2Controller(svcMock)
		req := httptest.NewRequest(http.MethodGet, "/followups/v2?id=uuid-abc", nil)
		rec := httptest.NewRecorder()

		ctrl.GetByIDHandler(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var body models.FollowUpV2
		_ = json.NewDecoder(rec.Body).Decode(&body)
		assert.Equal(t, "uuid-abc", body.ID)
		svcMock.AssertExpectations(t)
	})

	t.Run("400 - falta parámetro id", func(t *testing.T) {
		svcMock := new(MockFollowUpV2Service)
		ctrl := controller.NewFollowUpV2Controller(svcMock)

		req := httptest.NewRequest(http.MethodGet, "/followups/v2", nil)
		rec := httptest.NewRecorder()

		ctrl.GetByIDHandler(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		svcMock.AssertNotCalled(t, "GetFollowUpByID")
	})

	t.Run("404 - registro no encontrado", func(t *testing.T) {
		svcMock := new(MockFollowUpV2Service)
		svcMock.On("GetFollowUpByID", mock.Anything, "no-existe").Return(nil, service.ErrFollowUpNotFound)

		ctrl := controller.NewFollowUpV2Controller(svcMock)
		req := httptest.NewRequest(http.MethodGet, "/followups/v2?id=no-existe", nil)
		rec := httptest.NewRecorder()

		ctrl.GetByIDHandler(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		svcMock.AssertExpectations(t)
	})

	t.Run("500 - error de infraestructura", func(t *testing.T) {
		svcMock := new(MockFollowUpV2Service)
		svcMock.On("GetFollowUpByID", mock.Anything, "uuid-err").Return(nil, errors.New("db down"))

		ctrl := controller.NewFollowUpV2Controller(svcMock)
		req := httptest.NewRequest(http.MethodGet, "/followups/v2?id=uuid-err", nil)
		rec := httptest.NewRecorder()

		ctrl.GetByIDHandler(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		svcMock.AssertExpectations(t)
	})
}

// ---------------------------------------------------------------------------
// Tests de ListHandler
// ---------------------------------------------------------------------------

func TestListHandler(t *testing.T) {
	t.Run("200 - lista paginada con defaults", func(t *testing.T) {
		svcMock := new(MockFollowUpV2Service)
		expected := repository.PageResult[models.FollowUpV2]{
			Items: []models.FollowUpV2{{ID: "uuid-1"}}, Total: 1, Page: 0, PageSize: 20,
		}
		svcMock.On("GetPaginatedFollowUps", mock.Anything, 0, 20).Return(expected, nil)

		ctrl := controller.NewFollowUpV2Controller(svcMock)
		req := httptest.NewRequest(http.MethodGet, "/followups/v2/list", nil)
		rec := httptest.NewRecorder()

		ctrl.ListHandler(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		svcMock.AssertExpectations(t)
	})

	t.Run("200 - page y limit desde query params", func(t *testing.T) {
		svcMock := new(MockFollowUpV2Service)
		expected := repository.PageResult[models.FollowUpV2]{Items: []models.FollowUpV2{}, Total: 0, Page: 2, PageSize: 5}
		svcMock.On("GetPaginatedFollowUps", mock.Anything, 2, 5).Return(expected, nil)

		ctrl := controller.NewFollowUpV2Controller(svcMock)
		req := httptest.NewRequest(http.MethodGet, "/followups/v2/list?page=2&limit=5", nil)
		rec := httptest.NewRecorder()

		ctrl.ListHandler(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		svcMock.AssertExpectations(t)
	})

	t.Run("500 - error del servicio", func(t *testing.T) {
		svcMock := new(MockFollowUpV2Service)
		empty := repository.PageResult[models.FollowUpV2]{}
		svcMock.On("GetPaginatedFollowUps", mock.Anything, 0, 20).Return(empty, errors.New("timeout"))

		ctrl := controller.NewFollowUpV2Controller(svcMock)
		req := httptest.NewRequest(http.MethodGet, "/followups/v2/list", nil)
		rec := httptest.NewRecorder()

		ctrl.ListHandler(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		svcMock.AssertExpectations(t)
	})
}
