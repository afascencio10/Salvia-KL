// Package psychosocial3x3_test ejercita los endpoints HTTP del flujo 3x3 con
// httptest y un stub del servicio (sin base de datos). Vive en un paquete
// separado para no arrastrar tests preexistentes del paquete controller.
package psychosocial3x3_test

import (
	ctrl "bitsflow/salvia/controller"
	"bitsflow/internal/models"
	"bitsflow/salvia/service"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

type stub3x3 struct{ notFound bool }

func (s *stub3x3) GetHistory(_ context.Context, psicosocialID string) (*service.HistoryResult, error) {
	if s.notFound {
		return nil, service.ErrPsychosocial3x3NotFound
	}
	note := "Sí contestó, pide que lo llamen otro día."
	yes := true
	return &service.HistoryResult{
		Attempts:      []service.AttemptHistoryItem{{ID: "a1", SequenceNumber: 1, WasAnswered: true, Note: &note, ConsentGiven: &yes, AttemptAt: time.Now()}},
		Counters:      service.Counters{TotalCount: 1, DistinctDaysCount: 1},
		ProcessStatus: "en_gestion",
		CaseID:        "c-123",
		CaseName:      "Camila Pérez",
	}, nil
}
func (s *stub3x3) RegisterAttempt(_ context.Context, psicosocialID string, wasAnswered bool, note *string, attemptAt *time.Time, professionalID, team string) (*models.ContactAttempt, *service.Counters, string, error) {
	return &models.ContactAttempt{ID: "att-1", PsicosocialID: psicosocialID, CaseID: "c-123", WasAnswered: wasAnswered, Note: note},
		&service.Counters{DailyFailedCount: 3, TotalCount: 3, DistinctDaysCount: 1, DailyThresholdReached: true}, "abierto", nil
}
func (s *stub3x3) SetConsent(_ context.Context, attemptID string, consentGiven bool) (*service.ConsentResult, error) {
	return &service.ConsentResult{ID: attemptID, ConsentGiven: consentGiven, RequiresClosureForm: !consentGiven, CanScheduleSession: consentGiven}, nil
}
func (s *stub3x3) ScheduleSession(_ context.Context, psicosocialID string, immediate bool, scheduledAt *time.Time, scheduledTime *string, professionalID, team string) (*service.SessionResult, error) {
	url := "/salvia/psicosocial/sesion/t-1"
	res := &service.SessionResult{Session: &models.TeamContact{ID: "t-1", CaseID: "c-123", IsPsicoSession: true}}
	if immediate {
		res.RedirectURL = &url
	}
	return res, nil
}
func (s *stub3x3) SetNextAttempt(_ context.Context, psicosocialID string, at time.Time) (*time.Time, error) {
	return &at, nil
}

func router(svc service.Psychosocial3x3Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	ctrl.NewPsychosocialContactController(svc).RegisterRoutes(r.Group("/api/v1"))
	return r
}

func req(r *gin.Engine, method, path, body string) *httptest.ResponseRecorder {
	httpReq := httptest.NewRequest(method, path, bytes.NewReader([]byte(body)))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httpReq)
	return w
}

func TestRoutesRegisterWithoutConflict(t *testing.T) { _ = router(&stub3x3{}) }

func TestGetHistory(t *testing.T) {
	w := req(router(&stub3x3{}), http.MethodGet, "/api/v1/psychosocial/ps-1/contact-attempts", "")
	if w.Code != http.StatusOK {
		t.Fatalf("esperaba 200, obtuvo %d (%s)", w.Code, w.Body.String())
	}
	var out map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil || out["caseName"] != "Camila Pérez" {
		t.Fatalf("respuesta inesperada: %s", w.Body.String())
	}
}

func TestGetHistoryNotFound(t *testing.T) {
	w := req(router(&stub3x3{notFound: true}), http.MethodGet, "/api/v1/psychosocial/x/contact-attempts", "")
	if w.Code != http.StatusNotFound {
		t.Fatalf("esperaba 404, obtuvo %d", w.Code)
	}
}

func TestRegisterAttempt(t *testing.T) {
	w := req(router(&stub3x3{}), http.MethodPost, "/api/v1/psychosocial/ps-1/contact-attempts", `{"was_answered":false,"note":"No contestó, celular apagado."}`)
	if w.Code != http.StatusCreated {
		t.Fatalf("esperaba 201, obtuvo %d (%s)", w.Code, w.Body.String())
	}
}

func TestRegisterAttemptMissingField(t *testing.T) {
	w := req(router(&stub3x3{}), http.MethodPost, "/api/v1/psychosocial/ps-1/contact-attempts", `{"note":"sin was_answered"}`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("esperaba 400, obtuvo %d", w.Code)
	}
}

func TestSetConsent(t *testing.T) {
	w := req(router(&stub3x3{}), http.MethodPatch, "/api/v1/contact-attempts/att-1/consent", `{"consent_given":false}`)
	if w.Code != http.StatusOK {
		t.Fatalf("esperaba 200, obtuvo %d (%s)", w.Code, w.Body.String())
	}
	var out map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &out)
	if out["requiresClosureForm"] != true {
		t.Fatalf("requiresClosureForm esperado true")
	}
}

func TestScheduleSessionImmediate(t *testing.T) {
	w := req(router(&stub3x3{}), http.MethodPost, "/api/v1/psychosocial/ps-1/sessions", `{"immediate":true}`)
	if w.Code != http.StatusCreated {
		t.Fatalf("esperaba 201, obtuvo %d (%s)", w.Code, w.Body.String())
	}
}

func TestScheduleSessionMissingDate(t *testing.T) {
	w := req(router(&stub3x3{}), http.MethodPost, "/api/v1/psychosocial/ps-1/sessions", `{"immediate":false}`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("esperaba 400 por falta de fecha/hora, obtuvo %d", w.Code)
	}
}

func TestSetNextAttempt(t *testing.T) {
	w := req(router(&stub3x3{}), http.MethodPut, "/api/v1/psychosocial/ps-1/next-attempt", `{"next_contact_attempt_at":"2026-07-10T14:00:00Z"}`)
	if w.Code != http.StatusOK {
		t.Fatalf("esperaba 200, obtuvo %d (%s)", w.Code, w.Body.String())
	}
}
