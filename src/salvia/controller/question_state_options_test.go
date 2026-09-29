package controller

import (
	"bitsflow/internal/models"
	"bitsflow/salvia/service"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// Solo se usan Create, GetByID y Update; el resto de la interfaz va embebida.
type fakeQuestionSvc struct {
	service.QuestionService
	existing *models.Question
	saved    *models.Question
}

func (f *fakeQuestionSvc) Create(_ context.Context, q *models.Question) error {
	f.saved = q
	return nil
}

func (f *fakeQuestionSvc) GetByID(context.Context, string) (*models.Question, error) {
	return f.existing, nil
}

func (f *fakeQuestionSvc) Update(_ context.Context, q *models.Question) error {
	f.saved = q
	return nil
}

func callQuestions(t *testing.T, svc *fakeQuestionSvc, method, path, body string) int {
	t.Helper()
	gin.SetMode(gin.TestMode)
	router := gin.New()
	NewQuestionController(svc).RegisterRoutes(router.Group("/api/v1"))

	request := httptest.NewRequest(method, path, strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	return recorder.Code
}

// Registrar la pregunta de competencias necesita stateOptionsPath; antes el body
// lo ignoraba y la pregunta quedaba leyendo opciones de la tabla option, donde no
// hay ninguna.
func TestCreateQuestionKeepsStateOptionsPath(t *testing.T) {
	svc := &fakeQuestionSvc{}
	code := callQuestions(t, svc, http.MethodPost, "/api/v1/forms/form-1/questions", `{
		"formSectionId":"sec-1","repeaterGroupId":"rg-1","questionTypeId":"multiple",
		"description":"Competencias de la entidad","required":false,"order":23,
		"stateOptionsPath":"newBarriers.{_entryIndex}.competencias"}`)

	if code != http.StatusCreated {
		t.Fatalf("POST → %d, esperaba 201", code)
	}
	if svc.saved.StateOptionsPath == nil {
		t.Fatal("stateOptionsPath se perdió: llegó nil al servicio")
	}
	if got := *svc.saved.StateOptionsPath; got != "newBarriers.{_entryIndex}.competencias" {
		t.Errorf("stateOptionsPath = %q", got)
	}
	if svc.saved.FormID != "form-1" || svc.saved.QuestionTypeID != "multiple" || svc.saved.Order != 23 {
		t.Errorf("resto del cuerpo mal mapeado: %+v", svc.saved)
	}
}

// Las preguntas normales siguen sin path.
func TestCreateQuestionWithoutStateOptionsPath(t *testing.T) {
	svc := &fakeQuestionSvc{}
	callQuestions(t, svc, http.MethodPost, "/api/v1/forms/form-1/questions",
		`{"formSectionId":"sec-1","questionTypeId":"text","description":"Descripción"}`)

	if svc.saved.StateOptionsPath != nil {
		t.Errorf("stateOptionsPath = %v, esperaba nil", *svc.saved.StateOptionsPath)
	}
}

func TestUpdateQuestionStateOptionsPath(t *testing.T) {
	previous := "newBarriers.{_entryIndex}.competenicas" // con el typo
	cases := []struct {
		name string
		body string
		want *string
	}{
		{"corrige el path", `{"stateOptionsPath":"newBarriers.{_entryIndex}.competencias"}`,
			ptr("newBarriers.{_entryIndex}.competencias")},
		{"cadena vacía lo borra", `{"stateOptionsPath":""}`, nil},
		{"si no se envía no se toca", `{"required":true}`, &previous},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path := previous
			svc := &fakeQuestionSvc{existing: &models.Question{ID: "q-1", StateOptionsPath: &path}}
			if code := callQuestions(t, svc, http.MethodPut, "/api/v1/questions/q-1", tc.body); code != http.StatusOK {
				t.Fatalf("PUT → %d, esperaba 200", code)
			}
			switch {
			case tc.want == nil && svc.saved.StateOptionsPath != nil:
				t.Errorf("stateOptionsPath = %q, esperaba nil", *svc.saved.StateOptionsPath)
			case tc.want != nil && svc.saved.StateOptionsPath == nil:
				t.Errorf("stateOptionsPath = nil, esperaba %q", *tc.want)
			case tc.want != nil && *svc.saved.StateOptionsPath != *tc.want:
				t.Errorf("stateOptionsPath = %q, esperaba %q", *svc.saved.StateOptionsPath, *tc.want)
			}
		})
	}
}

func ptr(s string) *string { return &s }
