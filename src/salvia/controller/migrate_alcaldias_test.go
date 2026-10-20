package controller

import (
	"strings"
	"testing"
)

// entity_obligation.label es varchar(255) y una de las acciones de la hoja 4
// mide 269 caracteres: sin recorte el insert falla.
func TestEtiquetaObligacion(t *testing.T) {
	largo := strings.Repeat("palabra ", 40) // 320 caracteres
	casos := []struct {
		nombre  string
		entrada string
		revisar func(t *testing.T, got string)
	}{
		{
			nombre:  "una acción corta pasa intacta",
			entrada: "Garantizar cupos en casa de refugio",
			revisar: func(t *testing.T, got string) {
				if got != "Garantizar cupos en casa de refugio" {
					t.Fatalf("la cambió: %q", got)
				}
			},
		},
		{
			nombre:  "una larga se recorta y queda marcada",
			entrada: largo,
			revisar: func(t *testing.T, got string) {
				if n := len([]rune(got)); n > maxEtiquetaObligacion {
					t.Fatalf("quedó en %d caracteres, no cabe en varchar(255)", n)
				}
				if !strings.HasSuffix(got, "…") {
					t.Fatalf("sin marca de recorte: %q", got[len(got)-20:])
				}
				if strings.HasSuffix(strings.TrimSuffix(got, "…"), " ") {
					t.Fatal("cortó dejando el espacio antes de la marca")
				}
			},
		},
		{
			nombre:  "una sola palabra larguísima igual se recorta",
			entrada: strings.Repeat("a", 400),
			revisar: func(t *testing.T, got string) {
				if n := len([]rune(got)); n > maxEtiquetaObligacion {
					t.Fatalf("quedó en %d caracteres", n)
				}
			},
		},
		{
			nombre:  "los espacios de los bordes no cuentan",
			entrada: "   Otorgar subsidios   ",
			revisar: func(t *testing.T, got string) {
				if got != "Otorgar subsidios" {
					t.Fatalf("no la limpió: %q", got)
				}
			},
		},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) { c.revisar(t, etiquetaObligacion(c.entrada)) })
	}
}

// Las acciones salen del catálogo embebido, no de texto repetido en el código.
func TestAccionesDeAlcaldias(t *testing.T) {
	acciones := accionesDeAlcaldias()
	if len(acciones) != 5 {
		t.Fatalf("se esperaban las 5 acciones de la hoja 4, llegaron %d", len(acciones))
	}
	for _, a := range acciones {
		if strings.TrimSpace(a) == "" {
			t.Fatal("llegó una acción vacía")
		}
		if n := len([]rune(etiquetaObligacion(a))); n > maxEtiquetaObligacion {
			t.Fatalf("una acción no cabe ni recortada: %d caracteres", n)
		}
	}
}
