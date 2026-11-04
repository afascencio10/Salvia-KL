package repository

import (
	"errors"
	"fmt"
	"testing"
)

// REG-11 — validadores puros de lo que la persona escribe en Registro V3.

func TestNormalizeDateV3(t *testing.T) {
	cases := []struct {
		in      string
		want    string
		wantErr bool
	}{
		{"2026-10-30", "2026-10-30", false},
		{" 2026-10-30 ", "2026-10-30", false},
		{"2026-10-30 00:00:00 +0000 UTC", "2026-10-30", false}, // formato crudo que GetAnswers devolvía antes
		{"2026-10-30T05:00:00Z", "2026-10-30", false},
		{"2024-02-29", "2024-02-29", false}, // bisiesto
		{"2026-02-30", "", true},            // día imposible
		{"2026-13-01", "", true},
		{"30/10/2026", "", true},
		{"", "", true},
		{"275760-09-13", "", true}, // el máximo de <input type=date> no cabe en un DATE de Postgres
	}
	for _, c := range cases {
		got, err := normalizeDateV3("HV_2", c.in)
		if (err != nil) != c.wantErr || got != c.want {
			t.Errorf("normalizeDateV3(%q) = %q, err=%v; quería %q, error=%v", c.in, got, err, c.want, c.wantErr)
		}
		var inv *ErrInvalidAnswerV3
		if err != nil && (!errors.As(err, &inv) || inv.Key != "HV_2" || inv.Reason == "") {
			t.Errorf("normalizeDateV3(%q): el error debe ser *ErrInvalidAnswerV3 con campo y motivo, fue %#v", c.in, err)
		}
	}
}

func TestNormalizeTimeV3(t *testing.T) {
	cases := []struct {
		in      string
		want    string
		wantErr bool
	}{
		{"19:30", "19:30:00", false},
		{"08:05:30", "08:05:30", false},
		{"00:00", "00:00:00", false},
		{"23:59", "23:59:00", false},
		{" 07:15 ", "07:15:00", false},
		{"7pm", "", true},
		{"24:00", "", true},
		{"19:60", "", true},
		{"7:30", "07:30:00", false}, // la hora de un solo dígito se acepta y se normaliza
		{"", "", true},
	}
	for _, c := range cases {
		got, err := normalizeTimeV3("HV_3", c.in)
		if (err != nil) != c.wantErr || got != c.want {
			t.Errorf("normalizeTimeV3(%q) = %q, err=%v; quería %q, error=%v", c.in, got, err, c.want, c.wantErr)
		}
	}
}

func TestParseIntV3(t *testing.T) {
	cases := []struct {
		in      string
		want    int64
		wantErr bool
	}{
		{"3001234567", 3001234567, false},
		{"300 123 4567", 3001234567, false}, // antes se guardaba 300 sin avisar
		{"", 0, false},                      // vacío = valor marcador
		{"  ", 0, false},
		{"abc", 0, true},
		{"300abc", 0, true},
		{"-5", 0, true},
		{"99999999999999999999", 0, true}, // no cabe en un bigint
	}
	for _, c := range cases {
		got, err := parseIntV3("CA_2", c.in)
		if (err != nil) != c.wantErr || got != c.want {
			t.Errorf("parseIntV3(%q) = %d, err=%v; quería %d, error=%v", c.in, got, err, c.want, c.wantErr)
		}
	}
}

type fakeSQLStateErr struct{ state string }

func (e fakeSQLStateErr) Error() string    { return "sql error " + e.state }
func (e fakeSQLStateErr) SQLState() string { return e.state }

func TestAsInvalidAnswerV3(t *testing.T) {
	// 22xxx (dato con formato o tamaño inválido) → se vuelve respuesta no válida sin campo.
	err := asInvalidAnswerV3(fmt.Errorf("saveSection[HV]: actualizar form2: %w", fakeSQLStateErr{"22001"}))
	var inv *ErrInvalidAnswerV3
	if !errors.As(err, &inv) || inv.Key != "" || inv.Reason == "" {
		t.Errorf("22001 debía volverse *ErrInvalidAnswerV3 sin campo, fue %#v", err)
	}
	// Otros errores de Postgres (p. ej. llave foránea 23503) siguen siendo falla del servidor.
	other := fmt.Errorf("x: %w", fakeSQLStateErr{"23503"})
	got := asInvalidAnswerV3(other)
	var inv2 *ErrInvalidAnswerV3
	if errors.As(got, &inv2) || !errors.Is(got, other) {
		t.Errorf("23503 no debía transformarse en respuesta no válida, fue %#v", got)
	}
	// Un ErrInvalidAnswerV3 con campo pasa intacto; nil sigue siendo nil.
	orig := &ErrInvalidAnswerV3{Key: "HV_3", Reason: "x"}
	if got := asInvalidAnswerV3(fmt.Errorf("w: %w", orig)); !errors.Is(got, orig) {
		t.Errorf("un ErrInvalidAnswerV3 existente debía conservarse")
	}
	if asInvalidAnswerV3(nil) != nil {
		t.Errorf("nil debía seguir siendo nil")
	}
}
