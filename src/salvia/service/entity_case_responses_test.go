package service

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"errors"
	"testing"
	"time"
)

type fakeEntityCaseRelationsRepo struct {
	repository.EntityCaseRepository
	items []models.EntityCaseWithRelations
	err   error
}

func (f *fakeEntityCaseRelationsRepo) FindByCaseIDWithRelations(context.Context, string) ([]models.EntityCaseWithRelations, error) {
	return f.items, f.err
}

func ptr(s string) *string { return &s }

// La tarjeta muestra el conteo del checklist y los datos de la respuesta más
// reciente. "Más reciente" se decide por fecha: las filas llegan en cualquier
// orden y confiar en el ORDER BY del SQL es lo que rompería este caso.
func TestListByCaseAgregaRespuestaMasReciente(t *testing.T) {
	vieja := time.Date(2026, 3, 1, 10, 0, 0, 0, time.UTC)
	nueva := time.Date(2026, 6, 15, 8, 30, 0, 0, time.UTC)

	casos := &fakeEntityCaseRelationsRepo{items: []models.EntityCaseWithRelations{
		{RelID: "ec-1", EntityBranchName: "Fiscalía Leticia"},
	}}
	respuestas := &fakeObligationResponseRepo{
		checklistItems: []repository.EntityObligationChecklistItem{
			{EntityCaseID: "ec-1", EntityObligationID: "ob-1", Label: "Recepción de denuncia", Completada: true},
			{EntityCaseID: "ec-1", EntityObligationID: "ob-2", Label: "Medidas de protección", Completada: true},
			{EntityCaseID: "ec-1", EntityObligationID: "ob-3", Label: "Valoración de riesgo", Completada: false},
		},
		respuestas: []models.EntityCaseObligationResponse{
			{EntityCaseID: "ec-1", EntityObligationID: "ob-1", Estado: models.EntityCaseObligationEstadoEnTramite, UpdatedAt: vieja},
			{EntityCaseID: "ec-1", EntityObligationID: "ob-2", Estado: models.EntityCaseObligationEstadoCompletada,
				Radicado: ptr("NUNC-123"), Ampliacion: ptr("Se otorgaron medidas."), UpdatedAt: nueva},
		},
	}

	items, err := NewEntityCaseService(casos, nil, respuestas).ListByCase(context.Background(), "caso-9")
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}

	resumen := items[0].RespuestaEntidad
	if resumen == nil {
		t.Fatal("la tarjeta debería traer el resumen de la respuesta de la entidad")
	}
	if resumen.Total != 3 || resumen.Respondidas != 2 {
		t.Fatalf("conteos del checklist: total=%d respondidas=%d, se esperaba 3 y 2", resumen.Total, resumen.Respondidas)
	}
	if resumen.AccionLabel != "Medidas de protección" {
		t.Fatalf("la acción de la última respuesta debería ser 'Medidas de protección', fue %q", resumen.AccionLabel)
	}
	if resumen.Estado != models.EntityCaseObligationEstadoCompletada {
		t.Fatalf("estado de la última respuesta: %q", resumen.Estado)
	}
	if resumen.EstadoLabel != models.EntityCaseObligationEstadoLabels[models.EntityCaseObligationEstadoCompletada] {
		t.Fatalf("la etiqueta debería resolverse contra el catálogo, fue %q", resumen.EstadoLabel)
	}
	if resumen.Radicado != "NUNC-123" || resumen.Ampliacion != "Se otorgaron medidas." {
		t.Fatalf("radicado/ampliación no se desreferenciaron: %q / %q", resumen.Radicado, resumen.Ampliacion)
	}
	if resumen.RespondedAt == nil || !resumen.RespondedAt.Equal(nueva) {
		t.Fatalf("respondedAt debería ser la fecha de la respuesta más reciente, fue %v", resumen.RespondedAt)
	}
	if respuestas.caseIDPedido != "caso-9" {
		t.Fatalf("las respuestas deben buscarse por el caso consultado, se pidió %q", respuestas.caseIDPedido)
	}
}

// Una entidad con acciones asignadas pero sin responder sí muestra el bloque (con
// el contador en 0); una entidad sin acciones institucionales no lo muestra.
func TestListByCaseDistingueSinRespuestaDeSinAcciones(t *testing.T) {
	casos := &fakeEntityCaseRelationsRepo{items: []models.EntityCaseWithRelations{
		{RelID: "ec-pendiente"},
		{RelID: "ec-sin-acciones"},
	}}
	respuestas := &fakeObligationResponseRepo{
		checklistItems: []repository.EntityObligationChecklistItem{
			{EntityCaseID: "ec-pendiente", EntityObligationID: "ob-1", Label: "Atención en salud"},
		},
	}

	items, err := NewEntityCaseService(casos, nil, respuestas).ListByCase(context.Background(), "caso-9")
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}

	pendiente := items[0].RespuestaEntidad
	if pendiente == nil {
		t.Fatal("con acciones asignadas la tarjeta debe mostrar el bloque, aunque no haya respuesta")
	}
	if pendiente.Total != 1 || pendiente.Respondidas != 0 || pendiente.Estado != "" {
		t.Fatalf("resumen sin respuesta mal armado: %+v", pendiente)
	}
	if items[1].RespuestaEntidad != nil {
		t.Fatalf("sin acciones institucionales no debería haber bloque: %+v", items[1].RespuestaEntidad)
	}
}

// Si la consulta de entidades falla no hay nada que enriquecer: el error sube tal
// cual, porque sin tarjetas no hay pestaña.
func TestListByCasePropagaErrorDeEntidades(t *testing.T) {
	fallo := errors.New("db caída")
	casos := &fakeEntityCaseRelationsRepo{err: fallo}

	if _, err := NewEntityCaseService(casos, nil, &fakeObligationResponseRepo{}).ListByCase(context.Background(), "caso-9"); !errors.Is(err, fallo) {
		t.Fatalf("se esperaba el error del repositorio, llegó %v", err)
	}
}
