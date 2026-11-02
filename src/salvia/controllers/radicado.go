// T3 (REP-01): radicado corto y legible de los reportes V2 («K7QM-4R2X»).
//
// El icode (UUID) sigue siendo el identificador interno del reporte; el
// radicado es solo para que la persona lo guarde y lo consulte (T4).
package salvia_ctrl

import (
	"crypto/rand"
	"errors"
	"io"
	"time"
)

const (
	// 32 símbolos: A–Z sin I ni O, y 2–9 (sin 0 ni 1), para no confundirlos al leerlos o dictarlos.
	radicadoAlphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	// Intentos máximos de guardar un radicado antes de rendirse (colisión con uno ya existente).
	radicadoMaxAttempts = 5
	// Símbolos aleatorios del radicado (los otros 2 salen de la fecha).
	radicadoRandomSymbols = 6
)

var (
	// Colombia no tiene horario de verano: un huso fijo UTC-5 equivale a America/Bogota y no depende de tzdata.
	radicadoZone = time.FixedZone("COT", -5*60*60)
	// Día 0 de la parte «fecha» del radicado.
	radicadoEpoch = time.Date(2026, time.January, 1, 0, 0, 0, 0, radicadoZone)

	// ErrRadicadoExhausted indica que se agotaron los intentos sin lograr guardar un radicado único.
	ErrRadicadoExhausted = errors.New("no se pudo asignar un radicado único")
)

// radicadoDatePrefix devuelve los 2 primeros símbolos del radicado: los días
// transcurridos desde radicadoEpoch (fecha de Colombia) mezclados con una
// fórmula fija e inversible, (día×619 + 331) mod 1024, escrita en base 32.
// Dos reportes de días distintos no pueden chocar (cambia el prefijo) y la
// fecha no se ve literal. El prefijo se repite cada 1024 días (~2,8 años);
// desde ahí vuelve a decidir el índice único.
func radicadoDatePrefix(now time.Time) string {
	y, m, d := now.In(radicadoZone).Date()
	day := int(time.Date(y, m, d, 0, 0, 0, 0, radicadoZone).Sub(radicadoEpoch) / (24 * time.Hour))
	mix := ((day*619+331)%1024 + 1024) % 1024
	return string([]byte{radicadoAlphabet[mix>>5], radicadoAlphabet[mix&31]})
}

// radicadoRandomSymbolsFrom lee n símbolos del alfabeto desde r, descartando
// los bytes que sesgarían el módulo.
func radicadoRandomSymbolsFrom(r io.Reader, n int) (string, error) {
	size := len(radicadoAlphabet)
	limit := 256 - (256 % size)
	out := make([]byte, 0, n)
	buf := make([]byte, n*2)
	for len(out) < n {
		if _, err := io.ReadFull(r, buf); err != nil {
			return "", err
		}
		for _, b := range buf {
			if int(b) >= limit {
				continue
			}
			out = append(out, radicadoAlphabet[int(b)%size])
			if len(out) == n {
				break
			}
		}
	}
	return string(out), nil
}

// generateRadicadoAt arma un radicado «XXXX-XXXX»: 2 símbolos de la fecha
// mezclada + 6 aleatorios (≈10⁹ combinaciones por día).
func generateRadicadoAt(now time.Time, r io.Reader) (string, error) {
	random, err := radicadoRandomSymbolsFrom(r, radicadoRandomSymbols)
	if err != nil {
		return "", err
	}
	s := radicadoDatePrefix(now) + random
	return s[:4] + "-" + s[4:], nil
}

// generateRadicado genera un radicado con la hora actual y crypto/rand.
func generateRadicado() (string, error) {
	return generateRadicadoAt(time.Now(), rand.Reader)
}

// assignRadicado genera radicados y los guarda con assign hasta que uno quede
// asignado al reporte icode. assign devuelve false (sin error) si el valor ya
// lo tenía otro reporte: entonces se genera otro, hasta radicadoMaxAttempts.
// Cualquier error de assign (incluida una clave duplicada por carrera entre
// dos transacciones) se devuelve tal cual para que el llamador haga rollback.
func assignRadicado(icode string, gen func() (string, error), assign func(icode, radicado string) (bool, error)) (string, error) {
	for attempt := 0; attempt < radicadoMaxAttempts; attempt++ {
		radicado, err := gen()
		if err != nil {
			return "", err
		}
		assigned, err := assign(icode, radicado)
		if err != nil {
			return "", err
		}
		if assigned {
			return radicado, nil
		}
	}
	return "", ErrRadicadoExhausted
}
