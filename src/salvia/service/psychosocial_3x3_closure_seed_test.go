package service

import (
	"os"
	"strings"
	"testing"
)

// El bug reportado ("no permite cierre por imposibilidad de contacto") fue que el
// código apuntaba por UUID a un formulario que ningún archivo del repositorio
// creaba: el seed existía solo en la máquina de quien lo escribió. Estas pruebas
// atan el código al seed, para que mover un id en un lado y no en el otro falle acá
// y no en la cara de una profesional a mitad de una gestión.
func leerSeedDeCierre(t *testing.T) string {
	t.Helper()
	// Es el mismo archivo que la aplicación embebe y aplica al arrancar
	// (src/internal/seed/seed.go).
	ruta := "../../internal/seed/closure_form.sql"
	contenido, err := os.ReadFile(ruta)
	if err != nil {
		t.Fatalf("falta el seed del formulario de cierre (%s): sin él, el cierre por imposibilidad de contacto responde 404", ruta)
	}
	return string(contenido)
}

func TestSeedDeCierreCreaElFormularioQueUsaElCodigo(t *testing.T) {
	seed := leerSeedDeCierre(t)

	if !strings.Contains(seed, ClosureFormID) {
		t.Fatalf("el seed no crea el formulario %s, que es el que abre el flujo 3x3", ClosureFormID)
	}
	if !strings.Contains(seed, ClosureQMotivo) {
		t.Fatalf("el seed no crea la pregunta 'Motivo de cierre' %s, que el backend prerresponde y luego lee para cerrar", ClosureQMotivo)
	}
}

// Cada motivo que el backend sabe traducir a un estado tiene que existir como
// opción: si falta, la profesional llega al formulario y no encuentra el motivo con
// el que venía, o lo guarda y el cierre no transiciona.
func TestSeedDeCierreTieneTodosLosMotivos(t *testing.T) {
	seed := leerSeedDeCierre(t)

	for motivo := range closureReasonToStatus {
		if !strings.Contains(seed, "'"+motivo+"'") {
			t.Errorf("el motivo %q no está como opción en el seed", motivo)
		}
	}
}

// Los dos disparadores del cierre son los que abren el formulario desde el flujo,
// incluido el del reporte (imposibilidad de contacto).
func TestSeedDeCierreCubreLosDisparadores(t *testing.T) {
	seed := leerSeedDeCierre(t)

	for motivo := range validClosureReasons {
		if !strings.Contains(seed, "'"+motivo+"'") {
			t.Errorf("el disparador de cierre %q no tiene opción en el seed", motivo)
		}
		if _, ok := closureReasonToStatus[motivo]; !ok {
			t.Errorf("el disparador %q no tiene estado resultante: el cierre quedaría a medias", motivo)
		}
	}
}
