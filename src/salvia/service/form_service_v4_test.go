package service

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSeguimientoV4ID(t *testing.T) {
	assert.Equal(t, "2bda4735-3927-4634-a666-c7f97949be0a", seguimientoV4ID("repeater:identEntidades"))
	assert.Equal(t, "55faf667-6ea2-40ff-871d-3321cd2e1f31", seguimientoV4ID("q:ident.sector"))
	assert.Equal(t, "a92070e2-e90c-4f1b-9512-75b1d5638250", seguimientoV4ID("q:seg.completadas"))
}

func TestV4ReasignacionDestino(t *testing.T) {
	cases := []struct {
		name     string
		level    int
		extremos string
		equipo   string
		want     v4Reasignacion
	}{
		{"extremo gana al equipo general", 1, "amenazas", "general", v4Reasignacion{apply: true, level: 4, team: "Riesgo alto"}},
		{"extremo en nivel 2", 2, "ninguna,armas", "ninguno", v4Reasignacion{apply: true, level: 4, team: "Riesgo alto"}},
		{"solo ninguna no es extremo", 1, "ninguna", "ninguno", v4Reasignacion{}},
		{"feminicidio no mueve", 1, "", "feminicidio", v4Reasignacion{}},
		{"riesgo sube a alto", 1, "", "riesgo", v4Reasignacion{apply: true, level: 3, team: "Riesgo alto"}},
		{"ya en alto no regenera", 4, "", "riesgo", v4Reasignacion{}},
		{"general baja", 3, "", "general", v4Reasignacion{apply: true, level: 2, team: "Riesgo bajo"}},
		{"ya en bajo no regenera", 2, "", "general", v4Reasignacion{}},
		{"extremo no aplica si ya es alto", 3, "amenazas", "general", v4Reasignacion{apply: true, level: 2, team: "Riesgo bajo"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, v4ReasignacionDestino(tc.level, tc.extremos, tc.equipo))
		})
	}
}

func TestV4PsicoAplica(t *testing.T) {
	assert.True(t, v4PsicoAplica("poblacion,suicida"))
	assert.True(t, v4PsicoAplica("poblacion,temor,territorio,sin_acceso"))
	assert.False(t, v4PsicoAplica("poblacion"))
	assert.False(t, v4PsicoAplica("poblacion,temor"))
	assert.False(t, v4PsicoAplica("suicida"))
}

type v4AnswerFake struct {
	direct  []models.Answer
	byEntry map[string][]models.Answer
}

func (f *v4AnswerFake) Create(context.Context, *models.Answer) error { return nil }
func (f *v4AnswerFake) Update(context.Context, *models.Answer) error { return nil }
func (f *v4AnswerFake) UpdateFields(context.Context, string, map[string]interface{}) error {
	return nil
}
func (f *v4AnswerFake) Delete(context.Context, string) error { return nil }
func (f *v4AnswerFake) FindByID(context.Context, string) (*models.Answer, error) {
	return nil, nil
}
func (f *v4AnswerFake) FindWithPagination(context.Context, int, int) (repository.PageResult[models.Answer], error) {
	return repository.PageResult[models.Answer]{}, nil
}
func (f *v4AnswerFake) FindBySubmissionID(context.Context, string) ([]models.Answer, error) {
	return nil, nil
}
func (f *v4AnswerFake) FindDirectBySubmissionID(context.Context, string) ([]models.Answer, error) {
	return f.direct, nil
}
func (f *v4AnswerFake) FindByRepeaterEntryID(_ context.Context, id string) ([]models.Answer, error) {
	return f.byEntry[id], nil
}
func (f *v4AnswerFake) FindByRepeaterEntryIDs(context.Context, []string) ([]models.Answer, error) {
	return nil, nil
}
func (f *v4AnswerFake) FindByQuestionID(context.Context, string) ([]models.Answer, error) {
	return nil, nil
}

type v4EntryFake struct{ items []models.RepeaterEntry }

func (f *v4EntryFake) Create(context.Context, *models.RepeaterEntry) error { return nil }
func (f *v4EntryFake) Update(context.Context, *models.RepeaterEntry) error { return nil }
func (f *v4EntryFake) UpdateFields(context.Context, string, map[string]interface{}) error {
	return nil
}
func (f *v4EntryFake) Delete(context.Context, string) error { return nil }
func (f *v4EntryFake) FindByID(context.Context, string) (*models.RepeaterEntry, error) {
	return nil, nil
}
func (f *v4EntryFake) FindWithPagination(context.Context, int, int) (repository.PageResult[models.RepeaterEntry], error) {
	return repository.PageResult[models.RepeaterEntry]{}, nil
}
func (f *v4EntryFake) FindBySubmissionID(context.Context, string) ([]models.RepeaterEntry, error) {
	return f.items, nil
}
func (f *v4EntryFake) FindBySubmissionIDOrdered(context.Context, string) ([]models.RepeaterEntry, error) {
	return f.items, nil
}
func (f *v4EntryFake) FindBySubmissionIDAndGroupIDs(_ context.Context, _ string, groupIDs []string) ([]models.RepeaterEntry, error) {
	want := map[string]bool{}
	for _, id := range groupIDs {
		want[id] = true
	}
	var out []models.RepeaterEntry
	for _, e := range f.items {
		if want[e.RepeaterGroupID] {
			out = append(out, e)
		}
	}
	return out, nil
}

type v4TimelineFake struct{ events []models.CaseTimelineEvent }

func (f *v4TimelineFake) Create(_ context.Context, event *models.CaseTimelineEvent) error {
	f.events = append(f.events, *event)
	return nil
}
func (f *v4TimelineFake) GetByCaseID(context.Context, string, string, ...string) ([]repository.TimelineEventRow, error) {
	return nil, nil
}

type v4AttemptFake struct{ n int }

func (f *v4AttemptFake) CreateAttempt(context.Context, *models.FollowUpAttempt) error {
	f.n++
	return nil
}
func (f *v4AttemptFake) GetByFollowUpID(context.Context, string) ([]models.FollowUpAttempt, error) {
	return nil, nil
}
func (f *v4AttemptFake) CountByFollowUpID(context.Context, string) (int64, error) { return 0, nil }

type v4EntityFake struct {
	items   []models.EntityCaseWithRelations
	created []CreateEntityCaseInput
	actions []string
}

func (f *v4EntityFake) ListByCase(context.Context, string) ([]models.EntityCaseWithRelations, error) {
	return f.items, nil
}
func (f *v4EntityFake) Create(_ context.Context, input CreateEntityCaseInput) (*models.EntityCase, error) {
	f.created = append(f.created, input)
	id := fmt.Sprintf("ec-new-%d", len(f.created))
	f.items = append(f.items, models.EntityCaseWithRelations{RelID: id, EntityBranchID: input.EntityBranchID})
	return &models.EntityCase{ID: id, EntityBranchID: input.EntityBranchID}, nil
}
func (f *v4EntityFake) ListByEntity(context.Context, EntityCaseListFilter) (*EntityCaseListResult, error) {
	return nil, nil
}
func (f *v4EntityFake) ListEntities(context.Context) ([]models.EntityCatalogItem, error) {
	return nil, nil
}
func (f *v4EntityFake) ListCitiesByEntity(context.Context, int64) ([]models.EntityCityOption, error) {
	return nil, nil
}
func (f *v4EntityFake) ListObligations(context.Context, string, string) (*models.EntityCaseObligationsResponse, error) {
	return nil, nil
}
func (f *v4EntityFake) UpdateLastAction(_ context.Context, _, lastAction string) error {
	f.actions = append(f.actions, lastAction)
	return nil
}
func (f *v4EntityFake) GetObligacionesResumen(context.Context, string) (*models.ObligacionesResumenResponse, error) {
	return nil, nil
}
func (f *v4EntityFake) GetEntityCaseDetail(context.Context, string) (*models.EntityCaseDetailRow, error) {
	return nil, nil
}

type v4ObligationFake struct {
	created []models.EntityCaseObligation
	updated []string
}

func (f *v4ObligationFake) Create(_ context.Context, e *models.EntityCaseObligation) error {
	f.created = append(f.created, *e)
	return nil
}
func (f *v4ObligationFake) Update(context.Context, *models.EntityCaseObligation) error { return nil }
func (f *v4ObligationFake) UpdateFields(_ context.Context, id string, _ map[string]interface{}) error {
	f.updated = append(f.updated, id)
	return nil
}
func (f *v4ObligationFake) Delete(context.Context, string) error { return nil }
func (f *v4ObligationFake) FindByID(context.Context, string) (*models.EntityCaseObligation, error) {
	return nil, nil
}
func (f *v4ObligationFake) FindWithPagination(context.Context, int, int) (repository.PageResult[models.EntityCaseObligation], error) {
	return repository.PageResult[models.EntityCaseObligation]{}, nil
}
func (f *v4ObligationFake) FindByEntityCaseID(context.Context, string) ([]models.EntityCaseObligation, error) {
	return nil, nil
}
func (f *v4ObligationFake) FindByEntityCaseIDWithLabels(context.Context, string, string) ([]models.EntityCaseObligationItem, error) {
	return nil, nil
}
func (f *v4ObligationFake) ExistsForObligation(context.Context, string, string) (bool, error) {
	return false, nil
}

type v4EntityFollowFake struct{ rows []models.EntityCaseFollowUp }

func (f *v4EntityFollowFake) Create(_ context.Context, e *models.EntityCaseFollowUp) error {
	f.rows = append(f.rows, *e)
	return nil
}
func (f *v4EntityFollowFake) Update(context.Context, *models.EntityCaseFollowUp) error { return nil }
func (f *v4EntityFollowFake) UpdateFields(context.Context, string, map[string]interface{}) error {
	return nil
}
func (f *v4EntityFollowFake) Delete(context.Context, string) error { return nil }
func (f *v4EntityFollowFake) FindByID(context.Context, string) (*models.EntityCaseFollowUp, error) {
	return nil, nil
}
func (f *v4EntityFollowFake) FindWithPagination(context.Context, int, int) (repository.PageResult[models.EntityCaseFollowUp], error) {
	return repository.PageResult[models.EntityCaseFollowUp]{}, nil
}
func (f *v4EntityFollowFake) FindByEntityCaseID(context.Context, string) ([]models.EntityCaseFollowUp, error) {
	return nil, nil
}

type v4BarrierFollowFake struct{ rows []models.BarrierFollowUp }

func (f *v4BarrierFollowFake) Create(_ context.Context, e *models.BarrierFollowUp) error {
	f.rows = append(f.rows, *e)
	return nil
}
func (f *v4BarrierFollowFake) Update(context.Context, *models.BarrierFollowUp) error { return nil }
func (f *v4BarrierFollowFake) UpdateFields(context.Context, string, map[string]interface{}) error {
	return nil
}
func (f *v4BarrierFollowFake) Delete(context.Context, string) error { return nil }
func (f *v4BarrierFollowFake) FindByID(context.Context, string) (*models.BarrierFollowUp, error) {
	return nil, nil
}
func (f *v4BarrierFollowFake) FindWithPagination(context.Context, int, int) (repository.PageResult[models.BarrierFollowUp], error) {
	return repository.PageResult[models.BarrierFollowUp]{}, nil
}
func (f *v4BarrierFollowFake) FindByBarrierID(context.Context, string) ([]models.BarrierFollowUp, error) {
	return nil, nil
}

type v4Calendar struct{ levels []int }

func (c *v4Calendar) GetFollowUpByID(context.Context, string) (*models.FollowUpV2, error) {
	return nil, nil
}
func (c *v4Calendar) GetPaginatedFollowUps(context.Context, int, int) (repository.PageResult[models.FollowUpV2], error) {
	return repository.PageResult[models.FollowUpV2]{}, nil
}
func (c *v4Calendar) GetByCaseID(context.Context, string) ([]models.FollowUpV2, error) {
	return nil, nil
}
func (c *v4Calendar) GenerateOrRecalculate(context.Context, string, GenerateCalendarInput) ([]models.FollowUpV2, error) {
	return nil, nil
}
func (c *v4Calendar) GetFollowUpDetail(context.Context, string, bool) (*models.FollowUpDetailResponse, error) {
	return nil, nil
}
func (c *v4Calendar) GetAgentDayFollowUps(context.Context, string, time.Time) ([]models.FollowUpV2, []models.FollowUpV2, []models.FollowUpV2, error) {
	return nil, nil, nil, nil
}
func (c *v4Calendar) GetMyDayFollowUpsEnriched(context.Context, string, time.Time) (*models.MyDayResponse, error) {
	return nil, nil
}
func (c *v4Calendar) RegisterContactAttempt(context.Context, string, string, bool, *time.Time) (*models.FollowUpV2, error) {
	return nil, nil
}
func (c *v4Calendar) LoadFollowUp(context.Context, string, string, string) (*LoadFollowUpResult, error) {
	return nil, nil
}
func (c *v4Calendar) ReasignarCalendario(_ context.Context, _ string, level int) error {
	c.levels = append(c.levels, level)
	return nil
}
func (c *v4Calendar) GetByTeamPaginated(context.Context, string, repository.FollowUpFilters, int, int) ([]models.FollowUpV2, int64, error) {
	return nil, 0, nil
}
func (c *v4Calendar) GetAgentWorkload(context.Context, string, string) ([]repository.AgentWorkload, error) {
	return nil, nil
}
func (c *v4Calendar) GetFilterOptions(context.Context, string) (FilterOptions, error) {
	return FilterOptions{}, nil
}
func (c *v4Calendar) RescheduleFollowUp(context.Context, string, RescheduleInput) error { return nil }
func (c *v4Calendar) CloseCaseFollowUps(context.Context, string, string) error          { return nil }
func (c *v4Calendar) InitClosureForm(context.Context, string, string) (string, error) {
	return "", nil
}

type v4World struct {
	fu          *models.FollowUpV2
	followUps   *v4FollowFake
	cases       *v4CaseFake
	answers     *v4AnswerFake
	entries     *v4EntryFake
	timeline    *v4TimelineFake
	attempts    *v4AttemptFake
	entities    *v4EntityFake
	obligations *v4ObligationFake
	entFollows  *v4EntityFollowFake
	barFollows  *v4BarrierFollowFake
	calendar    *v4Calendar
	letters     *v4LetterFake
	tasks       *v4TaskFake
	barriers    *v4BarrierFake
	measures    *v4MeasureFake
	referrals   *v4ReferralFake
	psicos      *v4PsicoFake
	svc         *formService
}

func newV4World(level int) *v4World {
	ids := "b-active"
	branch := int64(9)
	w := &v4World{
		fu: &models.FollowUpV2{
			ID: "fu-1", CaseID: "case-1", Status: models.FollowUpStatusPendiente,
			ActiveBarrierIDs: &ids,
		},
		cases:       &v4CaseFake{level: level},
		answers:     &v4AnswerFake{byEntry: map[string][]models.Answer{}},
		entries:     &v4EntryFake{},
		timeline:    &v4TimelineFake{},
		attempts:    &v4AttemptFake{},
		entities:    &v4EntityFake{items: []models.EntityCaseWithRelations{{RelID: "ec-1", EntityBranchID: branch}}},
		obligations: &v4ObligationFake{},
		entFollows:  &v4EntityFollowFake{},
		barFollows:  &v4BarrierFollowFake{},
		calendar:    &v4Calendar{},
		letters:     &v4LetterFake{},
		tasks:       &v4TaskFake{},
		barriers: &v4BarrierFake{active: []models.BarrierV2{{
			ID: "b-active", EntityBranchID: &branch, Status: models.BarrierV2StatusOpen,
		}}},
		measures:  &v4MeasureFake{},
		referrals: &v4ReferralFake{},
		psicos:    &v4PsicoFake{},
	}
	w.followUps = &v4FollowFake{fu: w.fu}
	w.svc = &formService{
		followUpRepo: w.followUps, caseRepo: w.cases, answerRepo: w.answers, repeaterEntryRepo: w.entries,
		caseTimelineRepo: w.timeline, followUpAttemptRepo: w.attempts, entityCaseSvc: w.entities,
		entityCaseObligationRepo: w.obligations, entityCaseFollowUpRepo: w.entFollows,
		barrierFollowUpRepo: w.barFollows, followUpV2Svc: w.calendar, entityLetterRepo: w.letters,
		caseTaskRepo: w.tasks, barrierV2Repo: w.barriers, emergencyMeasureRepo: w.measures,
		emergencyMeasureReferralRepo: w.referrals, psychosocialSupportRepo: w.psicos,
	}
	return w
}

func (w *v4World) direct(pairs ...string) {
	for i := 0; i < len(pairs); i += 2 {
		w.answers.direct = append(w.answers.direct, models.Answer{QuestionID: v4Key(pairs[i]), Value: pairs[i+1]})
	}
}

func (w *v4World) entry(group, id string, iteration int, parent string, pairs ...string) {
	var parentID *string
	if parent != "" {
		parentID = &parent
	}
	w.entries.items = append(w.entries.items, models.RepeaterEntry{
		ID: id, FormSubmissionID: "sub-1", RepeaterGroupID: v4Key(group),
		Iteration: iteration, ParentEntryID: parentID,
	})
	var rows []models.Answer
	for i := 0; i < len(pairs); i += 2 {
		rows = append(rows, models.Answer{QuestionID: v4Key(pairs[i]), Value: pairs[i+1]})
	}
	w.answers.byEntry[id] = rows
}

func (w *v4World) run(t *testing.T) {
	t.Helper()
	assert.NoError(t, w.svc.completeFollowUpV4(context.Background(), "sub-1", "agent-1"))
}

func taskTypes(tasks []models.CaseTask) []string {
	out := make([]string, len(tasks))
	for i, task := range tasks {
		out[i] = task.Type
	}
	return out
}

func TestCompleteFollowUpV4ChannelsAndOtra(t *testing.T) {
	w := newV4World(1)
	w.direct("q:ident.conocimiento", "true")
	w.entry("repeater:identEntidades", "ident-1", 1, "",
		"q:ident.entidad", "42",
		"q:ident.acuerdo", "true",
		"q:ident.debida", "true",
		"q:ident.canales", "notificacion,articulacion,denuncia,misalvia,ninguna",
		"q:ident.radicado", "NUNC-1",
		"q:ident.otraRuta", "ruta",
	)
	w.run(t)

	assert.Len(t, w.entities.created, 1)
	assert.Equal(t, int64(42), w.entities.created[0].EntityBranchID)
	assert.Len(t, w.letters.rows, 1)
	assert.Empty(t, w.letters.rows[0].BarrierID)
	assert.Equal(t, []string{"proyectar_oficio", "gestion_llamada", "gestion_llamada", "solicitud_entidad"}, taskTypes(w.tasks.rows))
	assert.Equal(t, "Articulación institucional telefónica o por mensajería", w.tasks.rows[1].Description)
	assert.Equal(t, "Denuncia fácil", w.tasks.rows[2].Description)
	assert.Equal(t, 1, w.attempts.n)

	otra := newV4World(1)
	otra.direct("q:ident.conocimiento", "true")
	otra.entry("repeater:identEntidades", "ident-otra", 1, "",
		"q:ident.entidad", "otra",
		"q:ident.acuerdo", "true",
		"q:ident.debida", "true",
		"q:ident.canales", "notificacion",
	)
	otra.run(t)
	assert.Empty(t, otra.entities.created)
	assert.Empty(t, otra.letters.rows)
	assert.Empty(t, otra.tasks.rows)
}

func TestCompleteFollowUpV4Barriers(t *testing.T) {
	w := newV4World(1)
	w.direct("q:ident.conocimiento", "true")
	w.entry("repeater:identEntidades", "ident-1", 1, "",
		"q:ident.entidad", "42",
		"q:ident.nuevaBarrera", "true",
		"q:ident.acuerdo", "false",
	)
	w.entry("repeater:identBarreras", "ident-b", 1, "ident-1",
		"q:ident.barrera.sector", "salud",
		"q:ident.barrera.gestion", "activacion",
		"q:ident.barrera.descripcion", "demora",
	)
	w.entry("repeater:segEntidades", "seg-1", 1, "",
		"q:seg.actualizacion", "hubo cambios",
		"q:seg.completadas", "obl-1",
		"q:seg.pendientes", "obl-1,otra",
		"q:seg.otraPendiente", "Visita",
		"q:seg.nuevaBarrera", "false",
		"q:seg.acuerdo", "false",
	)
	w.entry("repeater:segBarrerasActivas", "seg-a", 1, "seg-1",
		"q:seg.barreraActiva.persisten", "true",
		"q:seg.barreraActiva.cierre", "true",
		"q:seg.barreraActiva.motivoCierre", "superada",
		"q:seg.barreraActiva.gestion", "activacion",
	)
	w.run(t)

	assert.Len(t, w.barriers.rows, 1)
	assert.Equal(t, "salud", w.barriers.rows[0].Sector)
	assert.NotNil(t, w.barriers.rows[0].EntityBranchID)
	assert.Equal(t, int64(42), *w.barriers.rows[0].EntityBranchID)
	assert.Len(t, w.letters.rows, 1)
	assert.Equal(t, w.barriers.rows[0].ID, w.letters.rows[0].BarrierID)
	assert.Equal(t, []string{"proyectar_oficio"}, taskTypes(w.tasks.rows))
	assert.Equal(t, "Barreras", w.tasks.rows[0].Category)

	assert.Len(t, w.barFollows.rows, 1)
	assert.Equal(t, "b-active", w.barFollows.rows[0].BarrierID)
	assert.True(t, w.barFollows.rows[0].ClosesBarrier)
	assert.Equal(t, "activacion", w.barFollows.rows[0].ManagementActions)
	assert.Equal(t, models.BarrierV2StatusManaged, w.barriers.status["b-active"])

	assert.Equal(t, []string{"obl-1"}, w.obligations.updated)
	assert.Len(t, w.obligations.created, 1)
	assert.Equal(t, "Visita", *w.obligations.created[0].CustomLabel)
	assert.Len(t, w.entFollows.rows, 1)
	assert.True(t, w.entFollows.rows[0].RutaActualizada)
}

func TestCompleteFollowUpV4MedidasApaganPsico(t *testing.T) {
	w := newV4World(1)
	w.direct(
		"q:plan.remision", "true",
		"q:plan.equipos", "medidas,psico,osignh",
		"q:plan.medidasTipos", "alojamiento,transporte",
		"q:plan.psico", "poblacion,suicida",
		"q:plan.osignh", "obligatorio",
	)
	w.run(t)
	if assert.Len(t, w.referrals.referrals, 1) {
		assert.Equal(t, models.EmergencyReferralStatusPendienteDeRevision, w.referrals.referrals[0].Status)
		assert.Equal(t, []string{"alojamiento", "transporte"}, w.referrals.types[0])
	}
	assert.Empty(t, w.measures.rows)
	assert.Empty(t, w.psicos.rows)

	ok := newV4World(1)
	ok.direct(
		"q:plan.remision", "true",
		"q:plan.equipos", "psico,osignh",
		"q:plan.psico", "poblacion,suicida",
		"q:plan.osignh", "obligatorio",
	)
	ok.run(t)
	assert.Empty(t, ok.measures.rows)
	if assert.Len(t, ok.psicos.rows, 2) {
		assert.Equal(t, "derivacion", ok.psicos.rows[0].Type)
		assert.Equal(t, "osignh", ok.psicos.rows[1].Type)
	}
	assert.Equal(t, []string{"validar_remision", "validar_remision"}, taskTypes(ok.tasks.rows))
}

func TestCompleteFollowUpV4Reasignacion(t *testing.T) {
	w := newV4World(1)
	w.direct("q:risk.extremos", "amenazas", "q:plan.equipoReasigna", "general")
	w.run(t)
	assert.Equal(t, []int{4}, w.calendar.levels)
	assert.True(t, w.cases.levelSet)
	assert.Equal(t, 4, w.cases.newLevel)
	assert.Equal(t, "Riesgo alto", w.cases.team)
	assert.Equal(t, "ag-1", w.cases.agent)

	quiet := newV4World(1)
	quiet.direct("q:risk.extremos", "ninguna", "q:plan.equipoReasigna", "feminicidio")
	quiet.run(t)
	assert.Empty(t, quiet.calendar.levels)
	assert.False(t, quiet.cases.levelSet)
}

type v4CierreFake struct {
	called bool
	input  CerrarCasoInput
}

func (f *v4CierreFake) CerrarCaso(_ context.Context, input CerrarCasoInput) error {
	f.called = true
	f.input = input
	return nil
}
func (f *v4CierreFake) ReabrirCaso(context.Context, string, string) error { return nil }

func TestCompleteFollowUpV4Conclusion(t *testing.T) {
	w := newV4World(1)
	cierre := &v4CierreFake{}
	w.svc.casoCierreService = cierre
	w.direct(
		"q:con.concluye", "true",
		"q:con.motivo", "otro",
		"q:con.otro", "se mudó",
		"q:con.descripcion", "texto",
		"q:con.accion", "true",
		"q:con.entidad", "15",
		"q:con.canales", "denuncia",
	)
	w.run(t)
	assert.True(t, cierre.called)
	assert.Equal(t, "otro", cierre.input.Motivo)
	assert.Equal(t, "se mudó", cierre.input.OtroMotivo)
	assert.True(t, cierre.input.AccionesInstitucionales)
	assert.Len(t, w.entities.created, 1)
	assert.Equal(t, int64(15), w.entities.created[0].EntityBranchID)
	assert.Equal(t, []string{"gestion_llamada"}, taskTypes(w.tasks.rows))
	assert.Empty(t, w.letters.rows)
}

func TestCompleteFollowUpV4AlertaDeBarreraActiva(t *testing.T) {
	w := newV4World(1)
	w.entry("repeater:segEntidades", "seg-1", 1, "",
		"q:seg.actualizacion", "sin cambios",
		"q:seg.acuerdo", "false",
	)
	w.entry("repeater:segBarrerasActivas", "seg-a", 1, "seg-1",
		"q:seg.barreraActiva.persisten", "true",
		"q:seg.barreraActiva.cierre", "false",
		"q:seg.barreraActiva.gestion", "alerta",
	)
	w.run(t)
	assert.Empty(t, w.letters.rows)
	assert.Equal(t, []string{"comite_caso"}, taskTypes(w.tasks.rows))
	assert.Empty(t, w.tasks.rows[0].AssignedUserID)
	assert.Equal(t, "b-active", *w.tasks.rows[0].BarrierID)
}

func TestCompleteFollowUpV4SecondSaveDoesNotRepeat(t *testing.T) {
	w := newV4World(1)
	w.direct("q:ident.conocimiento", "true")
	w.entry("repeater:identEntidades", "ident-1", 1, "",
		"q:ident.entidad", "42",
		"q:ident.acuerdo", "true",
		"q:ident.debida", "true",
		"q:ident.canales", "notificacion",
	)
	w.run(t)
	letters := len(w.letters.rows)
	tasks := len(w.tasks.rows)
	w.fu.Status = models.FollowUpStatusRealizado
	assert.NoError(t, w.svc.completeFollowUpV4(context.Background(), "sub-1", "agent-1"))
	assert.Equal(t, letters, len(w.letters.rows))
	assert.Equal(t, tasks, len(w.tasks.rows))
	assert.Equal(t, 1, w.followUps.updates)
	var edited bool
	for _, event := range w.timeline.events {
		if event.Type == models.TimelineTypeSeguimientoEditado {
			edited = true
		}
	}
	assert.True(t, edited)
}
