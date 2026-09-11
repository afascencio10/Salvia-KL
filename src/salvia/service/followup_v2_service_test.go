package service_test

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"bitsflow/salvia/service"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// ---------------------------------------------------------------------------
// Mock del repositorio
// ---------------------------------------------------------------------------
// MockFollowUpRepo (implementa repository.FollowUpRepository) se declara una
// sola vez para todo el paquete service_test — ver form_service_test.go.

// ---------------------------------------------------------------------------
// Tests de GetFollowUpByID
// ---------------------------------------------------------------------------

func TestGetFollowUpByID(t *testing.T) {
	ctx := context.Background()

	t.Run("éxito - retorna el registro", func(t *testing.T) {
		repoMock := new(MockFollowUpRepo)
		agentID := "agent-1"
		expected := &models.FollowUpV2{ID: "uuid-123", CaseID: "case-1", AgentID: &agentID, Status: "open"}

		repoMock.On("FindByID", ctx, "uuid-123").Return(expected, nil)

		svc := service.NewFollowUpV2Service(repoMock, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
		result, err := svc.GetFollowUpByID(ctx, "uuid-123")

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		repoMock.AssertExpectations(t)
	})

	t.Run("not found - retorna ErrFollowUpNotFound", func(t *testing.T) {
		repoMock := new(MockFollowUpRepo)

		repoMock.On("FindByID", ctx, "no-existe").Return(nil, gorm.ErrRecordNotFound)

		svc := service.NewFollowUpV2Service(repoMock, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
		result, err := svc.GetFollowUpByID(ctx, "no-existe")

		assert.Nil(t, result)
		assert.True(t, errors.Is(err, service.ErrFollowUpNotFound))
		repoMock.AssertExpectations(t)
	})

	t.Run("error de infraestructura - se propaga sin envolver", func(t *testing.T) {
		repoMock := new(MockFollowUpRepo)
		dbErr := errors.New("connection refused")

		repoMock.On("FindByID", ctx, "uuid-err").Return(nil, dbErr)

		svc := service.NewFollowUpV2Service(repoMock, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
		result, err := svc.GetFollowUpByID(ctx, "uuid-err")

		assert.Nil(t, result)
		assert.ErrorIs(t, err, dbErr)
		repoMock.AssertExpectations(t)
	})
}

// ---------------------------------------------------------------------------
// Tests de GetPaginatedFollowUps
// ---------------------------------------------------------------------------

func TestGetPaginatedFollowUps(t *testing.T) {
	ctx := context.Background()

	t.Run("éxito - retorna página con items", func(t *testing.T) {
		repoMock := new(MockFollowUpRepo)
		expected := repository.PageResult[models.FollowUpV2]{
			Items:    []models.FollowUpV2{{ID: "uuid-1"}, {ID: "uuid-2"}},
			Total:    2,
			Page:     0,
			PageSize: 20,
		}

		repoMock.On("FindWithPagination", ctx, 0, 20).Return(expected, nil)

		svc := service.NewFollowUpV2Service(repoMock, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
		result, err := svc.GetPaginatedFollowUps(ctx, 0, 20)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		repoMock.AssertExpectations(t)
	})

	t.Run("error de repositorio - se propaga", func(t *testing.T) {
		repoMock := new(MockFollowUpRepo)
		dbErr := errors.New("timeout")
		empty := repository.PageResult[models.FollowUpV2]{}

		repoMock.On("FindWithPagination", ctx, 0, 20).Return(empty, dbErr)

		svc := service.NewFollowUpV2Service(repoMock, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
		_, err := svc.GetPaginatedFollowUps(ctx, 0, 20)

		assert.ErrorIs(t, err, dbErr)
		repoMock.AssertExpectations(t)
	})
}
