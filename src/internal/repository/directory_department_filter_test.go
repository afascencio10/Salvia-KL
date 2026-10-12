package repository

import (
	"strings"
	"testing"
)

// El filtro por departamento del Directorio devolvía cero resultados porque
// comparaba contra directories.department_id, que está en NULL en las entidades
// cargadas antes de que existiera la columna. La condición ahora cae en
// security.city cuando la columna viene vacía.
func TestDepartmentFilterSQL(t *testing.T) {
	cond, args := departmentFilterSQL(11)

	// Sin los paréntesis el OR se comería los demás filtros (active, type,
	// sector, name), que se encadenan con AND: "sector = 'X' AND a = 1 OR b"
	// devuelve todo lo que cumpla b, sin importar el sector.
	if !strings.HasPrefix(cond, "(") || !strings.HasSuffix(cond, ")") {
		t.Fatalf("la condición debe ir entre paréntesis para encadenarse con AND, obtuve: %s", cond)
	}

	// Un placeholder de más o de menos no falla al compilar: revienta en Postgres
	// o, peor, ata el argumento a la columna equivocada.
	if got := strings.Count(cond, "?"); got != len(args) {
		t.Fatalf("placeholders (%d) y argumentos (%d) no coinciden: %s", got, len(args), cond)
	}

	for i, a := range args {
		if a != uint64(11) {
			t.Fatalf("argumento %d: esperaba el id del departamento, obtuve %v", i, a)
		}
	}

	// La razón de ser del cambio: el respaldo por municipio.
	if !strings.Contains(cond, "department_id IS NULL") ||
		!strings.Contains(cond, "security.city") ||
		!strings.Contains(cond, "city_i_code") {
		t.Fatalf("falta el respaldo que deduce el departamento del municipio: %s", cond)
	}
}
