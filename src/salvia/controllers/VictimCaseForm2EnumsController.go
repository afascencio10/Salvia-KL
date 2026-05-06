// Package salvia_ctrl contiene los controladores para manejar operaciones relacionadas con
// "VictimCaseForm2Enums" en el módulo Salvia, incluyendo la creación, consulta e invalidación de registros.
package salvia_ctrl

import (
	"net/http"

	"bitsflow/common/db"
	"bitsflow/common/utils"

	salvia_config "bitsflow/salvia/config"
	salvia_daos "bitsflow/salvia/dao"
)

type VictimCaseForm2EnumsRequest struct {
	VContact salvia_daos.VictimCaseForm2EnumsDTO `json:"victimContact"`
}

func GetVictimCaseForm2EnumsByAll(connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	var resData string = ""
	var resCode int = http.StatusInternalServerError

	// Libera la conexión si no se proporcionó un ID de conexión.
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// Consulta todos los registros de VictimCaseForm2Enums.
	enums, err := salvia_daos.GetAllVictimCaseForm2Enums(connData, &dbClientConfig, &dbServerConfig)

	println("VictimCaseForm2Enums Loaded!!")

	if err != nil {
		// Retorna error interno en caso de fallo.
		resData = err.Error()
		resCode = http.StatusInternalServerError
	} else {
		// Se carga un mapa para consulta offline
		for _, e := range enums {

			//TODO: Traducimos el nombre, aunque a futuro habrá que traducirlo cuando el usuario lo utilice en idiomas diferentes

			e.VictimCaseForm2EnumsName = salvia_config.Locale["sp"][e.VictimCaseForm2EnumsCategory+"_"+e.VictimCaseForm2EnumsCode]
			salvia_daos.VictimCaseForm2Enums[e.VictimCaseForm2EnumsCategory] = append(salvia_daos.VictimCaseForm2Enums[e.VictimCaseForm2EnumsCategory], e)
		}
		// Retorna la lista completa de contactos en formato JSON con éxito.
		resCode = http.StatusOK
		resData = utils.CommMsgGetJSONSuccess(enums)
	}
	return resCode, resData
}
