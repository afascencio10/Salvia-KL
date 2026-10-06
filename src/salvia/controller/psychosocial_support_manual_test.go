package controller

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"bitsflow/salvia/service"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

type fakeSupportSvc struct {
	service.PsychosocialSupportService
	creada *models.PsychosocialSupport
}

func (f *fakeSupportSvc) Create(_ context.Context, ps *models.PsychosocialSupport) error {
	ps.ID = "ps-1"
	f.creada = ps
	return nil
}

type fakeVictimCaseRepo struct {
	repository.VictimCaseLightRepository
}

func (f *fakeVictimCaseRepo) FindByICode(_ context.Context, iCode string) (*models.VictimCaseLight, error) {
	return &models.VictimCaseLight{VictimICode: iCode}, nil
}

type fakeCaseTaskRepo struct {
	repository.CaseTaskRepository
	creadas []models.CaseTask
	err     error
}

func (f *fakeCaseTaskRepo) Create(_ context.Context, t *models.CaseTask) error {
	if f.err != nil {
		return f.err
	}
	f.creadas = append(f.creadas, *t)
	return nil
}

func postRemisionManual(t *testing.T, tareas *fakeCaseTaskRepo) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)

	ctrl := NewPsychosocialSupportGinController(&fakeSupportSvc{}, &fakeVictimCaseRepo{}, tareas)
	router := gin.New()
	ctrl.RegisterRoutes(router.Group("/api/v1"))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/psychosocial-support/manual",
		strings.NewReader(`{"caseId":"CASO-1"}`))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

// El bug reportado: las remisiones creadas por esta vía nacían sin la tarea de
// validación, y la pantalla de detalle habilita "recepcionar / devolver" solo
// cuando esa tarea existe en ToDo. Sin ella, la remisión no se puede validar nunca.
func TestRemisionManualCreaLaTareaDeValidacion(t *testing.T) {
	tareas := &fakeCaseTaskRepo{}
	rec := postRemisionManual(t, tareas)

	if rec.Code != http.StatusCreated {
		t.Fatalf("se esperaba 201, llegó %d: %s", rec.Code, rec.Body.String())
	}
	if len(tareas.creadas) != 1 {
		t.Fatalf("se esperaba 1 tarea de validación, se crearon %d", len(tareas.creadas))
	}

	tarea := tareas.creadas[0]
	if tarea.Type != "validar_remision" {
		t.Errorf("tipo de tarea: %q", tarea.Type)
	}
	if tarea.Status != models.CaseTaskStatusToDo {
		t.Errorf("la tarea debe nacer en ToDo, está en %q: en cualquier otro estado la pantalla no muestra la validación", tarea.Status)
	}
	if tarea.PsychosocialSupportID == nil || *tarea.PsychosocialSupportID != "ps-1" {
		t.Errorf("la tarea debe apuntar a la remisión creada, apunta a %v", tarea.PsychosocialSupportID)
	}
	if tarea.CaseID != "CASO-1" {
		t.Errorf("la tarea debe quedar en el caso de la remisión, quedó en %q", tarea.CaseID)
	}
}

// Si la tarea falla, la remisión ya está creada: anularla sería peor. Se responde
// 201 igual y queda el WARN en el log para poder repararla.
func TestRemisionManualSobreviveAFalloDeLaTarea(t *testing.T) {
	rec := postRemisionManual(t, &fakeCaseTaskRepo{err: context.DeadlineExceeded})

	if rec.Code != http.StatusCreated {
		t.Fatalf("se esperaba 201 pese al fallo de la tarea, llegó %d", rec.Code)
	}

	var cuerpo map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &cuerpo); err != nil {
		t.Fatalf("respuesta no es JSON: %v", err)
	}
	if cuerpo["id"] != "ps-1" {
		t.Fatalf("la respuesta debe traer la remisión creada, trajo %v", cuerpo["id"])
	}
}
