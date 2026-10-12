package service

import (
	"reflect"
	"testing"
)

// La traducción vivía dentro del bloque que solo corría para remisiones con
// seguimiento, que es lo que dejaba en blanco la pantalla de validación.
func TestEtiquetasCriteriosPsico(t *testing.T) {
	casos := []struct {
		nombre   string
		valor    string
		esperado []string
	}{
		{
			nombre:   "traduce los códigos y descarta la casilla habilitadora",
			valor:    "criterio_obligatorio,conducta_suicida,sin_ruta",
			esperado: []string{"Conducta suicida", "Sin ruta de atención"},
		},
		{
			nombre:   "un código desconocido se muestra tal cual",
			valor:    "interseccionalidad,criterio_nuevo",
			esperado: []string{"Interseccionalidad", "criterio_nuevo"},
		},
		{
			nombre:   "tolera espacios y comas sueltas",
			valor:    " sin_acceso_psico , , naturalizacion_vbg ",
			esperado: []string{"Sin acceso a atención psicosocial", "Naturalización de VBG"},
		},
		{
			nombre:   "sin criterios devuelve lista vacía, no nula",
			valor:    "",
			esperado: []string{},
		},
		{
			nombre:   "solo la casilla habilitadora no es un criterio",
			valor:    "criterio_obligatorio",
			esperado: []string{},
		},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			got := etiquetasCriteriosPsico(c.valor)
			if got == nil {
				t.Fatal("devolvió nil; el JSON debe llevar una lista vacía")
			}
			if !reflect.DeepEqual(got, c.esperado) {
				t.Fatalf("esperaba %v, obtuve %v", c.esperado, got)
			}
		})
	}
}
