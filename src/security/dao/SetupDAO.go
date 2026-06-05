// Package security_daos — funciones DAO exclusivas para los endpoints de setup/administración.
package security_daos

import (
	common_controllers "bitsflow/common/controllers"
	common_dao "bitsflow/common/dao"
	"bitsflow/common/db"
	"context"
	"errors"
	"fmt"
	"time"
)

// UpdateGeneralUserTeamByICode actualiza únicamente el campo team y la fecha de actualización
// de un usuario identificado por su ICode. Usado desde los endpoints de setup.
func UpdateGeneralUserTeamByICode(iCode string, team string, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	var persistenceCtrl common_controllers.PersistenceController = common_controllers.PersistenceController{}

	persistenceCtrl.Setup(connData, clientConfig, serverConfig)
	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Error (Setup team):", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	// SQL directo: UPDATE security.general_user SET general_user_team=$1, general_user_update_date=$2 WHERE general_user_i_code=$3
	query := "UPDATE " + GeneralUserDBScheme + "." + GeneralUserDBName +
		" SET general_user_team=$1, general_user_update_date=$2" +
		" WHERE general_user_i_code=$3"

	persistenceCtrl.Exec(context.Background(), query, team, time.Now().Format("2006-01-02 15:04:05"), iCode)

	if persistenceCtrl.Error != nil {
		fmt.Println("SQL Query:", query)
		fmt.Println("SQL Error:", persistenceCtrl.Error)
		return persistenceCtrl.Error
	}

	if persistenceCtrl.RowsAffected == 0 {
		return errors.New("usuario no encontrado o sin cambios: " + iCode)
	}

	return nil
}

// GetGeneralUserByLogin obtiene un usuario por su login.
// Wrapper de conveniencia sobre GetGeneralUser para los endpoints de setup.
func GetGeneralUserByLogin(login string, user *GeneralUserDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	by := common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"GeneralUserLogin"},
		AttrsValue: []interface{}{login},
	}
	return GetGeneralUser(by, user, connData, clientConfig, serverConfig)
}

// GetGeneralUserByProfileId obtiene el usuario asociado a un perfil dado su ID.
// Retorna error si no existe ningún general_user vinculado al perfil.
func GetGeneralUserByProfileId(profileId uint64, user *GeneralUserDTO, connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) error {
	by := common_controllers.By{
		Operator:   common_dao.SQL_AND,
		AttrsName:  []string{"GeneralUserGeneralUserProfile"},
		AttrsValue: []interface{}{profileId},
	}
	return GetGeneralUser(by, user, connData, clientConfig, serverConfig)
}
