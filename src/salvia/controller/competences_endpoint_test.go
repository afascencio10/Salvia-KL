package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// ListCompetences solo lee el catálogo embebido en salvia_config, así que no
// necesita servicios ni base de datos: el controlador se arma vacío.
func newCompetencesRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	NewFormController(nil, nil).RegisterRoutes(router.Group("/api/v1"))
	return router
}

func getCompetences(t *testing.T, query string) []map[string]string {
	t.Helper()
	recorder := httptest.NewRecorder()
	newCompetencesRouter().ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/competences"+query, nil))

	if recorder.Code != http.StatusOK {
		t.Fatalf("GET %s → %d, esperaba 200; cuerpo: %s", query, recorder.Code, recorder.Body.String())
	}
	var options []map[string]string
	if err := json.Unmarshal(recorder.Body.Bytes(), &options); err != nil {
		t.Fatalf("GET %s → respuesta no parseable: %v", query, err)
	}
	return options
}

func TestListCompetences(t *testing.T) {
	cases := []struct {
		name  string
		query string
		want  int
	}{
		{"institución marcada", "?sector=salud&institutions=ips", 30},
		{"varias instituciones", "?sector=salud&institutions=ips,eps", 39},
		{"sin institución cae al sector", "?sector=proteccion", 44},
		{"institución de otro sector manda", "?sector=proteccion&institutions=fiscalia_proteccion", 8},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			options := getCompetences(t, tc.query)
			if len(options) != tc.want {
				t.Fatalf("%s → %d opciones, esperaba %d", tc.query, len(options), tc.want)
			}
			for _, option := range options {
				if option["label"] == "" || option["value"] == "" {
					t.Errorf("opción incompleta: %+v", option)
				}
			}
		})
	}
}

// Los valores llegan del querystring tal como los arme el frontend: comas
// sobrantes y espacios no deben inventar una institución vacía, que resolvería a
// "sin institución" y devolvería el sector entero en vez de las 39 esperadas.
func TestListCompetencesIgnoresBlankInstitutions(t *testing.T) {
	if options := getCompetences(t, "?sector=salud&institutions=,ips,%20eps%20,"); len(options) != 39 {
		t.Fatalf("→ %d opciones, esperaba 39", len(options))
	}
}

// Sin opciones el cuerpo tiene que ser [] y no null: resolveQuestionOptions en
// dinamic-form espera un arreglo.
func TestListCompetencesEmptyIsArray(t *testing.T) {
	recorder := httptest.NewRecorder()
	newCompetencesRouter().ServeHTTP(recorder,
		httptest.NewRequest(http.MethodGet, "/api/v1/competences?sector=barrera_transversal", nil))

	if body := strings.TrimSpace(recorder.Body.String()); body != "[]" {
		t.Fatalf("cuerpo = %q, esperaba []", body)
	}
}
