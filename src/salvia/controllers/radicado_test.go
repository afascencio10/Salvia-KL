package salvia_ctrl

import (
	"bytes"
	"errors"
	"regexp"
	"strings"
	"testing"
	"time"
)

var radicadoFormat = regexp.MustCompile(`^[A-HJ-NP-Z2-9]{4}-[A-HJ-NP-Z2-9]{4}$`)

func bogota(y int, m time.Month, d, h, mi int) time.Time {
	return time.Date(y, m, d, h, mi, 0, 0, radicadoZone)
}

func TestRadicadoAlphabetHas32UnambiguousSymbols(t *testing.T) {
	if len(radicadoAlphabet) != 32 {
		t.Fatalf("el alfabeto debe tener 32 símbolos, tiene %d", len(radicadoAlphabet))
	}
	if strings.ContainsAny(radicadoAlphabet, "IO01") {
		t.Fatalf("el alfabeto no debe contener I, O, 0 ni 1: %s", radicadoAlphabet)
	}
}

func TestGenerateRadicadoFormatAndUniqueness(t *testing.T) {
	seen := map[string]bool{}
	for i := 0; i < 500; i++ {
		r, err := generateRadicado()
		if err != nil {
			t.Fatal(err)
		}
		if !radicadoFormat.MatchString(r) {
			t.Fatalf("formato inválido: %q", r)
		}
		if seen[r] {
			t.Fatalf("radicado repetido en 500 generados: %q", r)
		}
		seen[r] = true
	}
}

// C33: el prefijo es igual todo el día (hora de Colombia), cambia con la fecha y
// no es la fecha literal.
func TestRadicadoDatePrefix(t *testing.T) {
	inicio := radicadoDatePrefix(bogota(2026, time.November, 2, 0, 0))
	fin := radicadoDatePrefix(bogota(2026, time.November, 2, 23, 59))
	if inicio != fin {
		t.Fatalf("el prefijo debe ser el mismo todo el día: %s vs %s", inicio, fin)
	}
	// 2 de noviembre de 2026 = día 305 desde 2026-01-01 → (305×619+331) mod 1024 = 710 → «YG».
	if inicio != "YG" {
		t.Fatalf("prefijo esperado YG para 2026-11-02, obtuve %s", inicio)
	}
	if siguiente := radicadoDatePrefix(bogota(2026, time.November, 3, 0, 0)); siguiente == inicio {
		t.Fatalf("el prefijo debe cambiar con la fecha: %s", siguiente)
	}
	// En UTC, las 23:00 de Bogotá del 2 ya son el 3 de noviembre: debe mandar la fecha de Colombia.
	if got := radicadoDatePrefix(time.Date(2026, time.November, 3, 4, 0, 0, 0, time.UTC)); got != inicio {
		t.Fatalf("23:00 en Bogotá (04:00 UTC del día siguiente) sigue siendo el día 2: %s vs %s", got, inicio)
	}
	// La fórmula es inversible: 1024 días seguidos dan 1024 prefijos distintos.
	vistos := map[string]bool{}
	for d := 0; d < 1024; d++ {
		p := radicadoDatePrefix(bogota(2026, time.January, 1, 12, 0).AddDate(0, 0, d))
		if vistos[p] {
			t.Fatalf("prefijo repetido dentro de 1024 días (día %d): %s", d, p)
		}
		vistos[p] = true
	}
}

func TestGenerateRadicadoAtWithFixedRandomness(t *testing.T) {
	// bytes 0..5 → símbolos A,B,C,D,E,F; el prefijo de fecha es «YG».
	r, err := generateRadicadoAt(bogota(2026, time.November, 2, 12, 0), bytes.NewReader([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}))
	if err != nil {
		t.Fatal(err)
	}
	if r != "YGAB-CDEF" {
		t.Fatalf("esperaba YGAB-CDEF, obtuve %s", r)
	}
}

func TestGenerateRadicadoAtFailsWithoutRandomness(t *testing.T) {
	if _, err := generateRadicadoAt(time.Now(), bytes.NewReader(nil)); err == nil {
		t.Fatal("esperaba error cuando no hay aleatoriedad disponible")
	}
}

// C29: colisión simulada → el generador reintenta y guarda otro.
func TestAssignRadicadoRetriesOnCollision(t *testing.T) {
	candidatos := []string{"AAAA-AAAA", "BBBB-BBBB", "CCCC-CCCC"}
	generados, guardados := 0, []string{}
	gen := func() (string, error) { c := candidatos[generados]; generados++; return c, nil }
	assign := func(icode, radicado string) (bool, error) {
		guardados = append(guardados, radicado)
		return radicado == "CCCC-CCCC", nil // los dos primeros «ya existen»
	}
	r, err := assignRadicado("icode-1", gen, assign)
	if err != nil || r != "CCCC-CCCC" {
		t.Fatalf("esperaba CCCC-CCCC sin error, obtuve %q, %v", r, err)
	}
	if generados != 3 || len(guardados) != 3 {
		t.Fatalf("esperaba 3 intentos, hubo %d generados y %d guardados", generados, len(guardados))
	}
}

func TestAssignRadicadoGivesUpAfterFiveAttempts(t *testing.T) {
	intentos := 0
	_, err := assignRadicado("icode-1",
		func() (string, error) { return "AAAA-AAAA", nil },
		func(string, string) (bool, error) { intentos++; return false, nil })
	if !errors.Is(err, ErrRadicadoExhausted) {
		t.Fatalf("esperaba ErrRadicadoExhausted, obtuve %v", err)
	}
	if intentos != radicadoMaxAttempts {
		t.Fatalf("esperaba %d intentos, hubo %d", radicadoMaxAttempts, intentos)
	}
}

func TestAssignRadicadoStopsOnStoreError(t *testing.T) {
	boom := errors.New("23505 duplicate key")
	intentos := 0
	_, err := assignRadicado("icode-1",
		func() (string, error) { return "AAAA-AAAA", nil },
		func(string, string) (bool, error) { intentos++; return false, boom })
	if !errors.Is(err, boom) || intentos != 1 {
		t.Fatalf("un error de guardado debe cortar de inmediato: err=%v intentos=%d", err, intentos)
	}
}
