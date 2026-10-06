package seed

import (
	"regexp"
	"strings"
	"testing"
)

var soloComentarios = regexp.MustCompile(`(?m)^\s*--.*$`)

func sqlEjecutable() string {
	return strings.TrimSpace(soloComentarios.ReplaceAllString(closureFormSQL, ""))
}

// El seed se aplica con db.Exec, que con el protocolo extendido de Postgres
// rechaza varias sentencias en una sola llamada. Por eso todo el contenido va
// dentro de un único bloque DO: si alguien agrega sentencias sueltas afuera, el
// arranque dejaría de sembrar el formulario y volveríamos al 404 original.
func TestSeedDeCierreEsUnaSolaSentencia(t *testing.T) {
	sql := sqlEjecutable()

	if !strings.HasPrefix(sql, "DO $$") {
		t.Fatal("el seed debe empezar con el bloque DO $$; una sentencia suelta antes rompería db.Exec")
	}
	if !strings.HasSuffix(sql, "END $$;") {
		t.Fatal("el seed debe terminar con END $$; nada puede ir después del bloque")
	}
	if n := strings.Count(sql, "DO $$"); n != 1 {
		t.Fatalf("debe haber exactamente un bloque DO, hay %d", n)
	}
}

// Las escrituras tienen que ser repetibles: el arranque puede encontrarse una base
// a medio sembrar, y un INSERT sin ON CONFLICT la dejaría caer con error de clave
// duplicada en vez de completar lo que falta.
func TestSeedDeCierreEsIdempotente(t *testing.T) {
	sql := sqlEjecutable()

	inserts := strings.Count(sql, "INSERT INTO")
	conflictos := strings.Count(sql, "ON CONFLICT (id) DO NOTHING")
	if inserts != conflictos {
		t.Fatalf("%d INSERT pero %d ON CONFLICT: algún insert no es repetible", inserts, conflictos)
	}
}

// El UUID por el que la aplicación pregunta antes de sembrar tiene que ser el que
// el archivo crea. Si se separan, EnsureClosureForm aplicaría el seed en cada
// arranque y el formulario seguiría sin aparecer.
func TestSeedDeCierreCreaElFormularioQueSeVerifica(t *testing.T) {
	const formID = "fcc7dc8d-835b-4559-9283-2ea8b8e6092b"

	if !strings.Contains(closureFormSQL, formID) {
		t.Fatalf("el seed no crea el formulario %s", formID)
	}
}

// Esta reparación corre en cada arranque contra toda la tabla, así que sus
// exclusiones son lo único que la separa de crear tareas donde no debe. Cada una
// está acá para que quitarla sea una decisión y no un descuido.
func TestReparacionDeValidacionMantieneSusExclusiones(t *testing.T) {
	sql := strings.TrimSpace(soloComentarios.ReplaceAllString(validationTasksSQL, ""))

	exclusiones := map[string]string{
		"'cerrado', 'en_devolucion'": "remisiones cerradas o devueltas: no hay nada que validar",
		"salvia.case_task":           "remisiones que ya tienen la tarea: no se duplica",
		"Remisión Validada":          "remisiones ya validadas sin tarea: no se vuelve a pedir",
		"salvia.team_contact":        "remisiones en gestión: crearles la tarea les bloquea el registro de contactos",
	}
	for fragmento, porque := range exclusiones {
		if !strings.Contains(sql, fragmento) {
			t.Errorf("falta la exclusión de %s (%s)", fragmento, porque)
		}
	}
	if n := strings.Count(sql, "NOT EXISTS"); n != 3 {
		t.Errorf("se esperaban 3 condiciones NOT EXISTS, hay %d", n)
	}
}

// Una sola sentencia: db.Exec la exige, y de su RowsAffected sale el conteo que se
// registra en el log del arranque.
func TestReparacionDeValidacionEsUnaSolaSentencia(t *testing.T) {
	sql := strings.TrimSpace(soloComentarios.ReplaceAllString(validationTasksSQL, ""))

	if !strings.HasPrefix(sql, "INSERT INTO salvia.case_task") {
		t.Fatal("la reparación debe ser un único INSERT")
	}
	if n := strings.Count(sql, ";"); n != 1 || !strings.HasSuffix(sql, ";") {
		t.Fatalf("debe haber un solo ';' y al final, hay %d", n)
	}
}

// La tarea tiene que nacer con el tipo y el estado que la pantalla busca; con
// cualquier otro valor se crearían filas inútiles en cada arranque.
func TestReparacionDeValidacionCreaLaTareaQueLaPantallaBusca(t *testing.T) {
	for _, esperado := range []string{"'validar_remision'", "'ToDo'", "'Psicosocial'"} {
		if !strings.Contains(validationTasksSQL, esperado) {
			t.Errorf("la tarea creada no incluye %s", esperado)
		}
	}
}
