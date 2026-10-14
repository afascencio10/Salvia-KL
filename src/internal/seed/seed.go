// Package seed aplica al arrancar los datos sin los cuales la aplicación no
// funciona: definiciones de formularios dinámicos que el código abre por un UUID
// fijo. No son datos de ejemplo ni de prueba — si faltan, la pantalla que los usa
// se rompe.
//
// Por qué va en el arranque y no en un script aparte: el formulario de cierre
// psicosocial se aplicaba a mano y en algún ambiente nadie lo corrió, así que el
// cierre por imposibilidad de contacto fallaba con un 404. Un paso manual que hay
// que recordar en cada base termina olvidándose.
package seed

import (
	_ "embed"
	"fmt"

	"gorm.io/gorm"
)

//go:embed closure_form.sql
var closureFormSQL string

//go:embed validation_tasks.sql
var validationTasksSQL string

// EnsureValidationTasks crea la case_task de validación en las remisiones
// psicosociales activas que no la tengan — sin ella, la pantalla de detalle no
// ofrece recepcionar ni devolver. Devuelve cuántas creó.
//
// Corre en cada arranque a propósito: en la vía del seguimiento la tarea se crea
// "best effort" y un fallo solo deja un WARN, así que esto es la red que recoge
// esos casos. Repara la condición, no una lista de remisiones; el archivo SQL
// explica qué deja fuera y por qué.
func EnsureValidationTasks(db *gorm.DB) (int64, error) {
	resultado := db.Exec(validationTasksSQL)
	if resultado.Error != nil {
		return 0, fmt.Errorf("reparando tareas de validación de remisiones: %w", resultado.Error)
	}
	return resultado.RowsAffected, nil
}

// EnsureClosureForm crea el formulario "Cierre de proceso psicosocial" si no está
// en la base, o lo revive si quedó borrado lógicamente. Devuelve true cuando tuvo
// que aplicarlo.
//
// En un arranque normal no escribe nada: comprueba primero y, si el formulario ya
// está activo, no toca la base. El SQL es idempotente de todos modos.
//
// formID lo pasa quien llama (service.ClosureFormID) para no repetir el UUID ni
// hacer que este paquete dependa de la capa de dominio.
func EnsureClosureForm(db *gorm.DB, formID string) (bool, error) {
	var activo int64
	if err := db.Raw(
		`SELECT COUNT(*) FROM salvia.form WHERE id = ? AND deleted_at IS NULL`, formID,
	).Scan(&activo).Error; err != nil {
		return false, fmt.Errorf("verificando el formulario de cierre: %w", err)
	}
	if activo > 0 {
		return false, nil
	}

	if err := db.Exec(closureFormSQL).Error; err != nil {
		return false, fmt.Errorf("aplicando el seed del formulario de cierre: %w", err)
	}

	// Se vuelve a comprobar: si el seed corrió pero el formulario sigue sin estar,
	// es que el UUID del código y el del archivo se separaron, y es mejor decirlo
	// al arrancar que dejar que falle cuando una profesional intente cerrar.
	if err := db.Raw(
		`SELECT COUNT(*) FROM salvia.form WHERE id = ? AND deleted_at IS NULL`, formID,
	).Scan(&activo).Error; err != nil {
		return false, fmt.Errorf("verificando el formulario de cierre tras aplicarlo: %w", err)
	}
	if activo == 0 {
		return false, fmt.Errorf("el seed se aplicó pero el formulario %s sigue sin existir", formID)
	}

	return true, nil
}
