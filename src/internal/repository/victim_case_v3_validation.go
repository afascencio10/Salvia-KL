package repository

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"gorm.io/gorm"
)

// REG-11 — validación y normalización de lo que la persona escribe en Registro V3 ANTES de tocar la
// BD. Un dato con formato o tamaño inválido ya no revienta el guardado con un 500 genérico: se
// devuelve un ErrInvalidAnswerV3 con el campo y el motivo, el controlador responde 400 y la pantalla
// le pide a la persona que lo corrija.

// ErrInvalidAnswerV3 — una respuesta no se puede guardar por su formato, tamaño u opción. La persona
// puede corregirla. Key es el FieldKey ("HV_3"); vacío si no se pudo identificar el campo.
type ErrInvalidAnswerV3 struct {
	Key    string
	Reason string
}

func (e *ErrInvalidAnswerV3) Error() string {
	if e.Key == "" {
		return "respuesta no válida: " + e.Reason
	}
	return fmt.Sprintf("respuesta no válida (%s): %s", e.Key, e.Reason)
}

func invalidAnswerV3(key, format string, args ...interface{}) error {
	return &ErrInvalidAnswerV3{Key: key, Reason: fmt.Sprintf(format, args...)}
}

// asInvalidAnswerV3 — red de seguridad: un error "22xxx" de Postgres (dato con formato o tamaño
// inválido) que la validación previa no atrapó se vuelve ErrInvalidAnswerV3 sin campo. Cualquier otro
// error se deja tal cual.
func asInvalidAnswerV3(err error) error {
	if err == nil {
		return nil
	}
	var inv *ErrInvalidAnswerV3
	if errors.As(err, &inv) {
		return err
	}
	var withState interface{ SQLState() string }
	if errors.As(err, &withState) && strings.HasPrefix(withState.SQLState(), "22") {
		return &ErrInvalidAnswerV3{Reason: "alguno de los datos tiene un formato o un tamaño que no se puede guardar"}
	}
	return err
}

// dummyDateV3 — fecha marcador de "sin respuesta" que insertDummyForm2Row deja en las columnas DATE
// NOT NULL. GetAnswers no la devuelve como respuesta.
const dummyDateV3 = "1900-01-01"

// normalizeDateV3 devuelve la fecha como "YYYY-MM-DD". Acepta ese formato y el crudo que GetAnswers
// devolvía antes ("2026-10-30 00:00:00 +0000 UTC"): se toman los 10 primeros caracteres.
func normalizeDateV3(key, raw string) (string, error) {
	s := strings.TrimSpace(raw)
	if len(s) > 10 && (s[10] == ' ' || s[10] == 'T') {
		s = s[:10]
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return "", invalidAnswerV3(key, "escribe una fecha válida con el formato año-mes-día (por ejemplo 2026-10-30)")
	}
	return t.Format("2006-01-02"), nil
}

// normalizeTimeV3 devuelve la hora como "HH:MM:SS". Acepta "HH:MM" y "HH:MM:SS" (24 horas).
func normalizeTimeV3(key, raw string) (string, error) {
	s := strings.TrimSpace(raw)
	for _, layout := range []string{"15:04", "15:04:05"} {
		if t, err := time.Parse(layout, s); err == nil {
			return t.Format("15:04:05"), nil
		}
	}
	return "", invalidAnswerV3(key, "escribe la hora en formato de 24 horas HH:MM (por ejemplo 19:30)")
}

// parseIntV3 convierte un número entero no negativo; vacío = 0 (valor marcador). Antes Sscanf leía el
// prefijo numérico y descartaba el resto en silencio ("300 123" se guardaba como 300).
func parseIntV3(key, raw string) (int64, error) {
	s := strings.ReplaceAll(strings.TrimSpace(raw), " ", "")
	if s == "" {
		return 0, nil
	}
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil || n < 0 {
		return 0, invalidAnswerV3(key, "escribe solo números, sin letras ni signos")
	}
	return n, nil
}

// Anchos máximos de las columnas de texto (information_schema), con caché corto: así el límite que se
// valida es siempre el real de la BD, aunque una migración lo cambie con la app corriendo.
type columnWidthsV3 struct {
	mu       sync.Mutex
	loadedAt time.Time
	widths   map[string]int // "tabla.columna" -> máximo de caracteres
}

var columnWidthsCacheV3 = &columnWidthsV3{} //nolint:gochecknoglobals

func (c *columnWidthsV3) max(tx *gorm.DB, table, column string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.widths == nil || time.Since(c.loadedAt) > time.Minute {
		var rows []struct {
			TableName  string
			ColumnName string
			MaxLen     int
		}
		err := tx.Raw(`
			SELECT table_name, column_name, character_maximum_length AS max_len
			FROM information_schema.columns
			WHERE table_schema = 'salvia' AND table_name IN ('victim_case_form2', 'victim_case')
			  AND character_maximum_length IS NOT NULL
		`).Scan(&rows).Error
		if err != nil {
			return 0 // sin pre-chequeo: la BD igual lo valida y asInvalidAnswerV3 lo traduce
		}
		m := make(map[string]int, len(rows))
		for _, r := range rows {
			m[r.TableName+"."+r.ColumnName] = r.MaxLen
		}
		c.widths = m
		c.loadedAt = time.Now()
	}
	return c.widths[table+"."+column]
}

// checkTextLengthV3 valida que value quepa en la columna (caracteres, no bytes). what es el sujeto de
// la frase ("Este dato", "Los apellidos"…).
func checkTextLengthV3(tx *gorm.DB, table, column, key, what, value string) error {
	max := columnWidthsCacheV3.max(tx, table, column)
	if max <= 0 {
		return nil
	}
	if n := utf8.RuneCountInString(value); n > max {
		return invalidAnswerV3(key, "%s admite máximo %d caracteres y escribiste %d", what, max, n)
	}
	return nil
}
