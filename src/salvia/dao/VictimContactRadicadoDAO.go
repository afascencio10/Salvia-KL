// T3 (REP-01): acceso a victim_contact.victim_contact_radicado. Archivo nuevo
// para no tocar VictimContactDAO.go (V1): la columna se crea por AutoMigrate
// (VictimContactSync en main.go) y V1 no la conoce.
package salvia_daos

import (
	"context"

	common_controllers "bitsflow/common/controllers"
	"bitsflow/common/db"
)

// AssignVictimContactRadicado guarda el radicado del reporte icode solo si el
// reporte aún no tiene uno y nadie más usa ese valor. Devuelve true si lo
// asignó, false si el valor ya existía (el llamador genera otro). Debe correr
// dentro de la transacción que insertó el reporte.
//
// reporterRole (WP-1 T8) es el rol de sesión de quien reporta ("" si es anónimo →
// NULL); se guarda en la misma sentencia (columna victim_contact_reporter_role,
// creada por AutoMigrate junto con el radicado).
func AssignVictimContactRadicado(icode, radicado, reporterRole string, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) (bool, error) {
	var persistenceCtrl common_controllers.PersistenceController
	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		return false, persistenceCtrl.Error
	}

	var role interface{}
	if reporterRole != "" {
		role = reporterRole
	}

	persistenceCtrl.Exec(context.Background(),
		`UPDATE salvia.victim_contact
		    SET victim_contact_radicado = $1,
		        victim_contact_reporter_role = $3
		  WHERE victim_contact_i_code = $2
		    AND victim_contact_radicado IS NULL
		    AND NOT EXISTS (SELECT 1 FROM salvia.victim_contact WHERE victim_contact_radicado = $1)`,
		radicado, icode, role)
	if persistenceCtrl.Error != nil {
		return false, persistenceCtrl.Error
	}
	return persistenceCtrl.RowsAffected == 1, nil
}
