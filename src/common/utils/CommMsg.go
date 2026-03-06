package utils

import (
	common_config "bitsflow/common/config"
	"encoding/json"
	"strings"
)

/*func getMap(status bool, message string, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"status":  status,
		"message": message,
		"data":    data,
	}
}*/

func CommMsgGetJSONSuccess(dataStr interface{}) string {
	var resData string = ""
	data, error := json.Marshal(dataStr)
	if error != nil {
		resData = `{"` + common_config.Enums.GLOBAL_ERROR + `":"Error interno"}`
	} else {
		resData = string(data)
		if resData == "null" {
			resData = "{}"
		}
	}
	return resData
}

func CommMsgGetJSONSuccessMultiple(attrNAmes []string, attrValues []interface{}, defaultNullValues []string) string {
	var resData []string = []string{}
	for i := 0; i < len(attrNAmes); i++ {
		data, error := json.Marshal(attrValues[i])
		if error != nil {
			resData = append(resData, `{"`+common_config.Enums.GLOBAL_ERROR+`":"Error interno"}`)
		} else {
			var dataTmp string = string(data)
			if dataTmp == "null" {
				dataTmp = defaultNullValues[i]
			}
			resData = append(resData, `"`+attrNAmes[i]+`":`+dataTmp)
		}
	}

	return "{" + strings.Join(resData, ",") + "}"
}

func CommMsgGetJSONErrors(m map[string]map[string]string) string {
	var resData string = ""
	data, error := json.Marshal(m)
	if error != nil {
		resData = `{"` + common_config.Enums.GLOBAL_ERROR + `":"Error interno"}`
	} else {
		resData = string(data)
	}
	return resData
}
