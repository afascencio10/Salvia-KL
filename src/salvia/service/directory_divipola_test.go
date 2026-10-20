package service

import (
	"bitsflow/internal/models"
	"testing"
)

// El recorte a municipio es lo que usa el buscador de sedes del formulario de
// seguimiento (entity_branch_api_controller).
func TestDivipolaMunicipio(t *testing.T) {
	casos := map[string]string{
		"11001000": "11001",
		"11001009": "11001",
		"05034000": "05034",
		"11001":    "11001",
		"110":      "110",
		"":         "",
	}
	for entrada, esperado := range casos {
		if got := models.DivipolaMunicipio(entrada); got != esperado {
			t.Errorf("DivipolaMunicipio(%q) = %q, esperaba %q", entrada, got, esperado)
		}
	}
}

// Las entidades cargadas en masa no traen centro poblado: sin la cabecera del
// municipio se quedaban sin DIVIPOLA que mostrar.
func TestDivipolaDeEntidad(t *testing.T) {
	ptr := func(s string) *string { return &s }

	casos := []struct {
		nombre    string
		townCode  *string
		municipio string
		esperado  string
	}{
		{"con centro poblado manda el de la entidad", ptr("11001009"), "11001", "11001009"},
		{"sin centro poblado cae en la cabecera", nil, "11001", "11001000"},
		{"centro poblado en blanco cuenta como vacío", ptr("   "), "05001", "05001000"},
		{"espacios alrededor no viajan al JSON", ptr(" 05034000 "), "05034", "05034000"},
		{"sin municipio no hay DIVIPOLA", nil, "", ""},
		{"municipio incompleto no se completa", nil, "110", ""},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			if got := divipolaDeEntidad(c.townCode, c.municipio); got != c.esperado {
				t.Fatalf("esperaba %q, obtuve %q", c.esperado, got)
			}
		})
	}
}
