package service

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"errors"
	"testing"
)

type fakeObligationResponseRepo struct {
	repository.EntityCaseObligationResponseRepository
	pertenece        bool
	perteneceErr     error
	guardado         *models.EntityCaseObligationResponse
	filtroUsado      repository.EntityObligationListFilter
	filtroContadores repository.EntityObligationListFilter
	filas            []repository.EntityObligationRow
	checklistItems   []repository.EntityObligationChecklistItem
	idsChecklist     []string
	respuestas       []models.EntityCaseObligationResponse
	caseIDPedido     string
}

func (f *fakeObligationResponseRepo) FindByCaseID(_ context.Context, caseID string) ([]models.EntityCaseObligationResponse, error) {
	f.caseIDPedido = caseID
	return f.respuestas, nil
}

func (f *fakeObligationResponseRepo) BelongsToEntityBranch(_ context.Context, _ int64, _, _ string) (bool, error) {
	return f.pertenece, f.perteneceErr
}

func (f *fakeObligationResponseRepo) Upsert(_ context.Context, r *models.EntityCaseObligationResponse) error {
	f.guardado = r
	return nil
}

func (f *fakeObligationResponseRepo) ListForEntityBranch(_ context.Context, filter repository.EntityObligationListFilter) ([]repository.EntityObligationRow, int64, error) {
	f.filtroUsado = filter
	if f.filas != nil {
		return f.filas, int64(len(f.filas)), nil
	}
	return []repository.EntityObligationRow{{EntityCaseID: "ec-1"}}, 1, nil
}

func (f *fakeObligationResponseRepo) CountForEntityBranch(_ context.Context, filter repository.EntityObligationListFilter) (repository.EntityObligationCounts, error) {
	f.filtroContadores = filter
	return repository.EntityObligationCounts{Pendientes: 3, Completadas: 2}, nil
}

func (f *fakeObligationResponseRepo) ListCaseCitiesForEntityBranch(context.Context, int64) ([]models.EntityCityOption, error) {
	return []models.EntityCityOption{{ID: "c1", Name: "Leticia"}}, nil
}

func (f *fakeObligationResponseRepo) FindChecklistsByEntityCaseIDs(_ context.Context, ids []string) ([]repository.EntityObligationChecklistItem, error) {
	f.idsChecklist = ids
	return f.checklistItems, nil
}

type fakeEntityCaseMetaRepo struct {
	repository.EntityCaseRepository
	err error
}

func (f *fakeEntityCaseMetaRepo) GetEntityMetaByBranchID(context.Context, int64) (*models.EntityBranchMeta, error) {
	if f.err != nil {
		return nil, f.err
	}
	return &models.EntityBranchMeta{EntityName: "ESE Hospital San Rafael"}, nil
}

func newObligationSvc(pertenece bool) (EntityObligationResponseService, *fakeObligationResponseRepo) {
	repo := &fakeObligationResponseRepo{pertenece: pertenece}
	return NewEntityObligationResponseService(repo, &fakeEntityCaseMetaRepo{}), repo
}

// El guard de escritura: una entidad no puede responder una acción que no tiene
// asignada, aunque mande ids válidos de otra entidad.
func TestSaveResponseRechazaAccionAjena(t *testing.T) {
	svc, repo := newObligationSvc(false)

	err := svc.SaveResponse(context.Background(), SaveEntityObligationResponseInput{
		EntityBranchID: 7, EntityCaseID: "ec-ajeno", EntityObligationID: "eo-1",
		Estado: models.EntityCaseObligationEstadoCompletada,
	})
	if !errors.Is(err, ErrEntityObligationNoAsignada) {
		t.Fatalf("err = %v, esperaba ErrEntityObligationNoAsignada", err)
	}
	if repo.guardado != nil {
		t.Error("guardó la respuesta de una acción ajena")
	}
}

func TestSaveResponseValidaEstado(t *testing.T) {
	svc, repo := newObligationSvc(true)

	err := svc.SaveResponse(context.Background(), SaveEntityObligationResponseInput{
		EntityBranchID: 7, EntityCaseID: "ec-1", EntityObligationID: "eo-1",
		Estado: "inventado",
	})
	if !errors.Is(err, ErrEntityObligationEstadoInvalido) {
		t.Fatalf("err = %v, esperaba ErrEntityObligationEstadoInvalido", err)
	}
	if repo.guardado != nil {
		t.Error("guardó con un estado fuera del catálogo")
	}
}

// Los campos opcionales en blanco se guardan como nulos, para distinguir "sin
// radicado" de "radicado vacío" cuando la tarjeta de Gestión Institucional los lea.
func TestSaveResponseNormalizaCamposOpcionales(t *testing.T) {
	svc, repo := newObligationSvc(true)

	err := svc.SaveResponse(context.Background(), SaveEntityObligationResponseInput{
		EntityBranchID: 7, EntityCaseID: "ec-1", EntityObligationID: "eo-1",
		Radicado: "   ", Ampliacion: "  Se realizó la valoración.  ",
		Estado: models.EntityCaseObligationEstadoEnTramite, RespondedByID: "u-1",
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	g := repo.guardado
	if g == nil {
		t.Fatal("no guardó nada")
	}
	if g.Radicado != nil {
		t.Errorf("radicado en blanco quedó como %q, esperaba nulo", *g.Radicado)
	}
	if g.Ampliacion == nil || *g.Ampliacion != "Se realizó la valoración." {
		t.Errorf("ampliación = %v, esperaba el texto sin espacios sobrantes", g.Ampliacion)
	}
	if g.RespondedByID != "u-1" || g.EntityCaseID != "ec-1" || g.EntityObligationID != "eo-1" {
		t.Errorf("guardado incompleto: %+v", g)
	}
}

// La pestaña por defecto es pendientes: cualquier valor raro en el querystring no
// debe colar un filtro de estado inesperado.
func TestListNormalizaStatus(t *testing.T) {
	for _, entrada := range []string{"", "pendiente", "cualquier_cosa", "COMPLETADA"} {
		svc, repo := newObligationSvc(true)
		if _, err := svc.ListForEntityBranch(context.Background(), repository.EntityObligationListFilter{
			EntityBranchID: 7, Status: entrada,
		}); err != nil {
			t.Fatalf("status %q: error %v", entrada, err)
		}
		if repo.filtroUsado.Status != models.EntityCaseObligationStatusPendiente {
			t.Errorf("status %q → filtro %q, esperaba pendiente", entrada, repo.filtroUsado.Status)
		}
	}

	svc, repo := newObligationSvc(true)
	if _, err := svc.ListForEntityBranch(context.Background(), repository.EntityObligationListFilter{
		EntityBranchID: 7, Status: models.EntityCaseObligationStatusCompletada,
	}); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if repo.filtroUsado.Status != models.EntityCaseObligationStatusCompletada {
		t.Errorf("filtro = %q, esperaba completada", repo.filtroUsado.Status)
	}
}

// Los contadores de las dos pestañas tienen que respetar los filtros de documento
// y ciudad, o mostrarían un número que no concuerda con la tabla.
func TestListContadoresRespetanFiltros(t *testing.T) {
	svc, repo := newObligationSvc(true)

	result, err := svc.ListForEntityBranch(context.Background(), repository.EntityObligationListFilter{
		EntityBranchID: 7, Document: "1099", CityID: "c1",
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if repo.filtroContadores.Document != "1099" || repo.filtroContadores.CityID != "c1" {
		t.Errorf("los contadores se pidieron con %+v, sin los filtros del listado", repo.filtroContadores)
	}
	if result.Counts.Pendientes != 3 || result.Counts.Completadas != 2 {
		t.Errorf("contadores = %+v", result.Counts)
	}
	if result.EntityName != "ESE Hospital San Rafael" {
		t.Errorf("entityName = %q", result.EntityName)
	}
	if len(result.CityOptions) != 1 {
		t.Errorf("cityOptions = %d, esperaba 1", len(result.CityOptions))
	}
}

// El panel "Relación con entidad" se pide una sola vez por caso, aunque el caso
// aparezca en varias filas: si se pidiera por fila, una página con 10 acciones del
// mismo caso haría 10 consultas iguales.
func TestChecklistsAgrupaYDeduplicaCasos(t *testing.T) {
	repo := &fakeObligationResponseRepo{
		pertenece: true,
		filas: []repository.EntityObligationRow{
			{EntityCaseID: "ec-1", EntityObligationID: "eo-1"},
			{EntityCaseID: "ec-1", EntityObligationID: "eo-2"},
			{EntityCaseID: "ec-2", EntityObligationID: "eo-9"},
			{EntityCaseID: "", EntityObligationID: "eo-x"}, // fila defectuosa
		},
		checklistItems: []repository.EntityObligationChecklistItem{
			{EntityCaseID: "ec-1", EntityObligationID: "eo-1", Label: "Valorar", Completada: true},
			{EntityCaseID: "ec-1", EntityObligationID: "eo-2", Label: "Remitir", Completada: false},
			{EntityCaseID: "ec-1", EntityObligationID: "eo-3", Label: "Notificar", Completada: true},
		},
	}
	svc := NewEntityObligationResponseService(repo, &fakeEntityCaseMetaRepo{})

	result, err := svc.ListForEntityBranch(context.Background(), repository.EntityObligationListFilter{EntityBranchID: 7})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}

	if len(repo.idsChecklist) != 2 {
		t.Fatalf("se pidieron %d casos (%v), esperaba 2 sin repetir y sin el id vacío",
			len(repo.idsChecklist), repo.idsChecklist)
	}

	ec1 := result.Checklists["ec-1"]
	if ec1.Total != 3 || ec1.Completadas != 2 || len(ec1.Items) != 3 {
		t.Errorf("ec-1: total=%d completadas=%d items=%d, esperaba 3/2/3",
			ec1.Total, ec1.Completadas, len(ec1.Items))
	}

	// ec-2 no trajo acciones, pero debe existir con lista vacía: si no, el frontend
	// tendría que distinguir entre "no consultado" y "sin acciones".
	ec2, ok := result.Checklists["ec-2"]
	if !ok {
		t.Fatal("ec-2 quedó sin entrada en checklists")
	}
	if ec2.Total != 0 || ec2.Items == nil {
		t.Errorf("ec-2: total=%d items=%v, esperaba 0 y una lista vacía", ec2.Total, ec2.Items)
	}
}

// Sin filas no se consulta el checklist.
func TestChecklistsSinFilasNoConsulta(t *testing.T) {
	repo := &fakeObligationResponseRepo{pertenece: true, filas: []repository.EntityObligationRow{}}
	svc := NewEntityObligationResponseService(repo, &fakeEntityCaseMetaRepo{})

	result, err := svc.ListForEntityBranch(context.Background(), repository.EntityObligationListFilter{EntityBranchID: 7})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if repo.idsChecklist != nil {
		t.Errorf("consultó el checklist con %v", repo.idsChecklist)
	}
	if result.Checklists == nil {
		t.Error("checklists quedó nulo; se espera un mapa vacío para que el JSON sea {}")
	}
}

// El catálogo del select va en orden fijo: recorrer el mapa lo cambiaría en cada
// carga y el usuario vería las opciones bailando.
func TestEstadoOptionsOrdenFijo(t *testing.T) {
	esperado := []string{
		models.EntityCaseObligationEstadoCompletada,
		models.EntityCaseObligationEstadoEnTramite,
		models.EntityCaseObligationEstadoFueraCompetencia,
		models.EntityCaseObligationEstadoRequiereInfo,
	}
	for i := 0; i < 20; i++ {
		options := entityObligationEstadoOptions()
		if len(options) != len(esperado) {
			t.Fatalf("%d opciones, esperaba %d", len(options), len(esperado))
		}
		for j, code := range esperado {
			if options[j].Value != code {
				t.Fatalf("posición %d = %q, esperaba %q", j, options[j].Value, code)
			}
			if options[j].Label == "" {
				t.Errorf("la opción %q quedó sin etiqueta", code)
			}
		}
	}
}

// Si falla resolver el nombre de la organización, la pantalla debe cargar igual:
// es un dato de encabezado, no el contenido.
func TestListToleraFalloDelNombreDeEntidad(t *testing.T) {
	repo := &fakeObligationResponseRepo{pertenece: true}
	svc := NewEntityObligationResponseService(repo, &fakeEntityCaseMetaRepo{err: errors.New("sin sede")})

	result, err := svc.ListForEntityBranch(context.Background(), repository.EntityObligationListFilter{EntityBranchID: 7})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if result.EntityName != "" {
		t.Errorf("entityName = %q, esperaba vacío", result.EntityName)
	}
	if len(result.Items) != 1 {
		t.Errorf("items = %d, esperaba que el listado cargara igual", len(result.Items))
	}
}
