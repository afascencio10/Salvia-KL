package service

import (
	"bitsflow/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
)

func TestTaskHasSolicitudChecklist(t *testing.T) {
	assert.False(t, taskHasSolicitudChecklist(nil))
	assert.False(t, taskHasSolicitudChecklist(datatypes.JSON(`{"solicitudes":[]}`)))
	assert.False(t, taskHasSolicitudChecklist(datatypes.JSON(`{"checklist":null}`)))
	assert.True(t, taskHasSolicitudChecklist(datatypes.JSON(`{"checklist":[]}`)))
	assert.True(t, taskHasSolicitudChecklist(datatypes.JSON(`{"checklist":[{"obligationId":"a","label":"A"}]}`)))
}

func TestObligacionesParaSolicitud(t *testing.T) {
	followUpID := "fu-1"
	rows := obligacionesParaSolicitud("ec-1", "agent-1", &followUpID, []solicitudItem{
		{ObligationID: "obl-1", Label: "Visita"},
		{ObligationID: "obl-1", Label: "Visita"},
		{ObligationID: "obl-2", Label: "Oficio"},
		{ObligationID: "otra", Label: "Acompañar"},
		{ObligationID: "otra", Label: "  "},
	}, map[string]bool{"obl-2": true})

	if assert.Len(t, rows, 3) {
		assert.Equal(t, "obl-1", *rows[0].EntityObligationID)
		assert.Equal(t, models.EntityCaseObligationStatusPendiente, rows[0].Status)
		assert.Equal(t, "ec-1", rows[0].EntityCaseID)
		assert.Equal(t, "agent-1", rows[0].CreatedByID)
		assert.Nil(t, rows[1].EntityObligationID)
		assert.Equal(t, "Acompañar", *rows[1].CustomLabel)
		assert.Nil(t, rows[2].EntityObligationID)
		assert.Equal(t, "Otra", *rows[2].CustomLabel)
	}

	assert.Empty(t, obligacionesParaSolicitud("ec-1", "agent-1", nil, nil, nil))
}
