package service

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"errors"
	"testing"
)

type fakeCompetencesQuestionRepo struct {
	repository.QuestionRepository
	questions []models.Question
	err       error
}

func (f *fakeCompetencesQuestionRepo) FindByRepeaterGroupID(context.Context, string) ([]models.Question, error) {
	return f.questions, f.err
}

func question(id string, path *string) models.Question {
	return models.Question{ID: id, StateOptionsPath: path}
}

func TestResolveCompetencesQuestionID(t *testing.T) {
	competencias := competencesStateOptionsPath
	ciudades := "newBarriers.{_entryIndex}.cities"

	cases := []struct {
		name string
		repo *fakeCompetencesQuestionRepo
		want string
	}{
		{
			// La mayoría de preguntas del repeater tiene StateOptionsPath nil; si el
			// resolver lo desreferenciara sin verificar, guardar un seguimiento con
			// barreras entraría en panic.
			name: "la encuentra entre preguntas sin path",
			repo: &fakeCompetencesQuestionRepo{questions: []models.Question{
				question("q-sector", nil),
				question("q-ciudad", &ciudades),
				question("q-gestion", nil),
				question("q-competencias", &competencias),
			}},
			want: "q-competencias",
		},
		{
			// Estado normal antes de crear la pregunta en cada ambiente: la barrera
			// se guarda sin competencias en vez de fallar.
			name: "todavía no existe",
			repo: &fakeCompetencesQuestionRepo{questions: []models.Question{
				question("q-sector", nil),
				question("q-ciudad", &ciudades),
			}},
			want: "",
		},
		{
			name: "repeater vacío",
			repo: &fakeCompetencesQuestionRepo{},
			want: "",
		},
		{
			name: "error de base de datos no tumba el guardado",
			repo: &fakeCompetencesQuestionRepo{err: errors.New("conexión caída")},
			want: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svc := &formService{questionRepo: tc.repo}
			if got := svc.resolveCompetencesQuestionID(context.Background(), "rg-barreras"); got != tc.want {
				t.Errorf("resolveCompetencesQuestionID = %q, esperaba %q", got, tc.want)
			}
		})
	}
}

// El servicio se arma con repos opcionales en varios flujos; sin questionRepo no
// debe reventar.
func TestResolveCompetencesQuestionIDWithoutRepo(t *testing.T) {
	svc := &formService{}
	if got := svc.resolveCompetencesQuestionID(context.Background(), "rg-barreras"); got != "" {
		t.Errorf("resolveCompetencesQuestionID = %q, esperaba vacío", got)
	}
}

// El path tiene que ser exactamente el que usa el frontend en
// hacer_seguimiento.html; si alguien cambia uno de los dos, la pregunta deja de
// resolverse y las competencias se guardarían vacías sin ningún error visible.
func TestCompetencesStateOptionsPathValue(t *testing.T) {
	if competencesStateOptionsPath != "newBarriers.{_entryIndex}.competencias" {
		t.Fatalf("competencesStateOptionsPath = %q", competencesStateOptionsPath)
	}
}
