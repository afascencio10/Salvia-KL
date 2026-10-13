package repository

import (
	"bitsflow/internal/models"
	"context"
	"testing"
)

// El catálogo de estados es la fuente de verdad del select del modal: cada código
// tiene que tener etiqueta, o la UI mostraría el código crudo.
func TestEntityCaseObligationEstados(t *testing.T) {
	codigos := []string{
		models.EntityCaseObligationEstadoCompletada,
		models.EntityCaseObligationEstadoEnTramite,
		models.EntityCaseObligationEstadoFueraCompetencia,
		models.EntityCaseObligationEstadoRequiereInfo,
	}
	for _, codigo := range codigos {
		if !models.IsValidEntityCaseObligationEstado(codigo) {
			t.Errorf("el estado %q no está en el catálogo", codigo)
		}
		if models.EntityCaseObligationEstadoLabels[codigo] == "" {
			t.Errorf("el estado %q no tiene etiqueta", codigo)
		}
	}
	if len(models.EntityCaseObligationEstadoLabels) != len(codigos) {
		t.Errorf("el catálogo tiene %d entradas y hay %d constantes declaradas",
			len(models.EntityCaseObligationEstadoLabels), len(codigos))
	}
	if models.IsValidEntityCaseObligationEstado("cualquier_cosa") {
		t.Error("un estado inventado pasó la validación")
	}
}

// FindByEntityCaseIDs con la lista vacía debe cortar antes de tocar la base:
// "WHERE entity_case_id IN ()" no es SQL válido en Postgres. El repositorio se
// construye con db nil a propósito — si el guard no estuviera, esto entraría en
// panic en vez de devolver una lista vacía.
func TestFindByEntityCaseIDsSinIDs(t *testing.T) {
	repo := NewEntityCaseObligationResponseRepository(nil)

	items, err := repo.FindByEntityCaseIDs(context.Background(), nil)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if items == nil {
		t.Error("devolvió nil; se espera una lista vacía para que el JSON sea [] y no null")
	}
	if len(items) != 0 {
		t.Errorf("devolvió %d elementos, esperaba 0", len(items))
	}

	if items, err = repo.FindByEntityCaseIDs(context.Background(), []string{}); err != nil || len(items) != 0 {
		t.Errorf("con slice vacío: %d elementos, err=%v", len(items), err)
	}
}
