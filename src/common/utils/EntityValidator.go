package utils

import (
	common_config "bitsflow/common/config"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type FieldDefinition struct {
	Name      string
	DBName    string
	Alias     string
	ModelType string
	MinSize   int64
	MaxSize   int64
	Required  bool
}

/*
Esta función valida los campos y retorna un DTO en caso de éxito. ECC, retorna el DTO en blanco y un mapa con los errores detectados
Se espera que jsonStr tenga la siguiente estructura:

	"objectName":{
		"EntityClassName":{
			"json_attr1_name":"value",
			"json_attr2_name":5,
			...
		}
	}

Pero también puede ser de la siguiente forma:

	{
		"objectName":{
				"EntityClassName":{
				"json_attr1_name":"value",
				"json_attr2_name":5,
				...
				"EntitySubClassName":{
					"json_attr1_name":"value",
					"json_attr2_name":5,
					...

				}
			}
		}
	}

Para lo cual se puede especificar en jsonFieldNameContainer así: "EntityClassName.EntitySubClassName". Sólo admite 1 nivel de profundidad

Actualiza el slice de errores encontrados para retornar al usuario. También retorna un error global en caso de que jsonFieldNameContainer no se encuentre en el json

Retorna true si encuentra que el dto a procesar vino vacío o con los valores por defecto. False ECC
*/
func ValidateJSONInput(dto interface{}, dtoMap map[string]interface{}, jsonFieldNameContainer string, fields map[string]FieldDefinition, checkFields map[string]bool, config map[string]map[string]string, dateTimeFormat string, dateFormat string, timeFormat string, collectedErrors map[string]map[string]string, canBeEmpty bool) bool {

	var ancestry []string = strings.Split(jsonFieldNameContainer, ".")
	var container interface{}
	var mapContainer map[string]interface{}
	var found bool
	var isEmpty bool = true
	var currentErrors map[string]map[string]string = map[string]map[string]string{}

	if len(ancestry) > 1 {

		container, found = dtoMap[ancestry[0]]
		if found {
			switch container.(type) {
			case map[string]interface{}:
				var c map[string]interface{} = container.(map[string]interface{})
				var c2 interface{}
				c2, found = c[ancestry[1]]
				switch c2.(type) {
				case map[string]interface{}:
					mapContainer, found = c2.(map[string]interface{})
					jsonFieldNameContainer = ancestry[1]
				default:
					SetError(collectedErrors, jsonFieldNameContainer, "default", common_config.Enums.GLOBAL_ERROR, "", config)
				}
			default:
				SetError(collectedErrors, jsonFieldNameContainer, "default", common_config.Enums.GLOBAL_ERROR, "", config)
			}
		}
	} else {
		container, found = dtoMap[jsonFieldNameContainer]
		switch container.(type) {
		case map[string]interface{}:
			mapContainer, found = container.(map[string]interface{})
		default:
			SetError(collectedErrors, jsonFieldNameContainer, "default", common_config.Enums.GLOBAL_ERROR, "", config)
		}
	}

	if !found {
		SetError(currentErrors, jsonFieldNameContainer, jsonFieldNameContainer, "common_validation_field_json_field_not_found", "common_global_error", config)
		return false
	}

	isEmpty = processFields(dto, fields, checkFields, currentErrors, mapContainer, jsonFieldNameContainer, config, dateTimeFormat, dateFormat, timeFormat)

	if !isEmpty || isEmpty && !canBeEmpty {
		//Mezclamos los errores obtenidos ocn los que vienen
		MergeMaps(collectedErrors, currentErrors)
	}
	return isEmpty
}

/*
*

	Igual que ValidateJSONInput pero recibe un slice de elementos:

	[
		{
			"json_attr1_name":"value",
			"json_attr2_name":5,
			...
		},
		{
			"json_attr1_name":"value2",
			"json_attr2_name":8,
			...
		}
	]
*/
func ValidateJSONMultipleInput(dtos []interface{}, dtoMap []interface{}, jsonFieldNameContainer string, fields map[string]FieldDefinition, checkFields map[string]bool, config map[string]map[string]string, dateTimeFormat string, dateFormat string, timeFormat string, collectedErrors map[string]map[string]string, canBeEmpty bool) bool {

	var mapContainer map[string]interface{}

	var isEmpty bool = true
	var currentErrors map[string]map[string]string = map[string]map[string]string{}

	for i, c := range dtoMap {

		switch c := c.(type) {
		case map[string]interface{}:
			mapContainer = c

			if processFields(dtos[i], fields, checkFields, currentErrors, mapContainer, jsonFieldNameContainer, config, dateTimeFormat, dateFormat, timeFormat) {
				isEmpty = false
			}
		default:
			SetError(collectedErrors, jsonFieldNameContainer, "default", common_config.Enums.GLOBAL_ERROR, "", config)
		}

	}

	if !isEmpty || isEmpty && !canBeEmpty {
		//Mezclamos los errores obtenidos ocn los que vienen
		MergeMaps(collectedErrors, currentErrors)
	}
	return isEmpty
}

func processFields(dto interface{}, fields map[string]FieldDefinition, checkFields map[string]bool, currentErrors map[string]map[string]string, container map[string]interface{}, jsonFieldNameContainer string, config map[string]map[string]string, dateTimeFormat string, dateFormat string, timeFormat string) bool {
	var isEmpty bool = true
	/*
		Hacemos un primer barrido de los campos que deben venir que se requieren para la base de datos
	*/

	for fname := range fields {
		isRequired, found := checkFields[fname]
		if len(checkFields) == 0 || found {
			var jsonTag = GetTag(dto, fname, "json")

			value, found := container[jsonTag]
			if found {
				switch fields[fname].ModelType {
				case "email":
					var tmp string
					switch value := value.(type) {
					case string:
						tmp = value
					case int:
						tmp = fmt.Sprintf("%v", value)
					case int8:
						tmp = fmt.Sprintf("%v", value)
					case int16:
						tmp = fmt.Sprintf("%v", value)
					case int32:
						tmp = fmt.Sprintf("%v", value)
					case int64:
						tmp = fmt.Sprintf("%v", value)
					case uint:
						tmp = fmt.Sprintf("%v", value)
					case uint8:
						tmp = fmt.Sprintf("%v", value)
					case uint16:
						tmp = fmt.Sprintf("%v", value)
					case uint32:
						tmp = fmt.Sprintf("%v", value)
					case uint64:
						tmp = fmt.Sprintf("%v", value)
					case float32:
						tmp = fmt.Sprintf("%v", value)
					case float64:
						tmp = fmt.Sprintf("%v", value)

					default:
						SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_format_error", "common_global_error", config)
					}
					//Si es realmente un string, se adiciona al dto. ECC significa que puede ser una fecha u otro campo que se validará y adicionará después
					if validateStringSize(fields, fname, tmp, isRequired, currentErrors, jsonFieldNameContainer, jsonTag, config) {
						var e error
						if isRequired && IsEmailValid(tmp) {
							e = SetField(dto, fname, tmp)
						} else if isRequired {
							e = errors.New("")
						}

						if e != nil {
							SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_email_error", "common_global_error", config)
						}
					}

					//Revisamos si el valor es diferente al de por defecto
					var empty string
					if empty != tmp {
						isEmpty = false
					}
				case "string":
					var tmp string
					switch value := value.(type) {
					case string:
						tmp = value
					case int:
						tmp = fmt.Sprintf("%v", value)
					case int8:
						tmp = fmt.Sprintf("%v", value)
					case int16:
						tmp = fmt.Sprintf("%v", value)
					case int32:
						tmp = fmt.Sprintf("%v", value)
					case int64:
						tmp = fmt.Sprintf("%v", value)
					case uint:
						tmp = fmt.Sprintf("%v", value)
					case uint8:
						tmp = fmt.Sprintf("%v", value)
					case uint16:
						tmp = fmt.Sprintf("%v", value)
					case uint32:
						tmp = fmt.Sprintf("%v", value)
					case uint64:
						tmp = fmt.Sprintf("%v", value)
					case float32:
						tmp = fmt.Sprintf("%v", value)
					case float64:
						tmp = fmt.Sprintf("%v", value)

					default:
						SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_format_error", "common_global_error", config)
					}
					//Si es realmente un string, se adiciona al dto. ECC significa que puede ser una fecha u otro campo que se validará y adicionará después
					if validateStringSize(fields, fname, tmp, isRequired, currentErrors, jsonFieldNameContainer, jsonTag, config) {
						e := SetField(dto, fname, tmp)
						if e != nil {
							SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_string_error", "common_global_error", config)
						}
					}

					//Revisamos si el valor es diferente al de por defecto
					var empty string
					if empty != tmp {
						isEmpty = false
					}
				case "datetime":
					switch value := value.(type) {
					case string:
						t, error := time.Parse(dateTimeFormat, value)
						var hasError bool = false
						if error != nil {
							//Significa que puede venir por defecto en zero
							t, error = time.Parse(time.RFC3339, value)
							if error != nil {
								hasError = true
								SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_date_error", "common_global_error", config)
							}
						}

						if !hasError && !t.IsZero() {
							//Se adiciona al dto
							e := SetField(dto, fname, t)
							if e != nil {
								SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_string_error", "common_global_error", config)
							}
						}

						//Revisamos si el valor es diferente al de por defecto
						var tmp string
						if tmp != value {
							isEmpty = false
						}
					case time.Time:
						//Se adiciona al dto
						e := SetField(dto, fname, value)
						if e != nil {
							SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_string_error", "common_global_error", config)
						}

						//Revisamos si el valor es diferente al de por defecto
						var empty time.Time
						if empty != value {
							isEmpty = false
						}

					default:
						SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_format_error", "common_global_error", config)
					}
				case "time":
					switch value := value.(type) {
					case nil:
						if isRequired {
							SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_date_error", "common_global_error", config)
						}
					case string:
						t, error := time.Parse(timeFormat, value)
						var hasError bool = false
						if error != nil {
							//Significa que puede venir por defecto en zero
							t, error = time.Parse(time.RFC3339, value)
							if error != nil {
								hasError = true
								SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_date_error", "common_global_error", config)
							}
						} else if isRequired && value == common_config.DateTime.TIME_ZERO_VALUE {
							hasError = true
							SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_date_error", "common_global_error", config)
						}

						if !hasError && !t.IsZero() {
							//Se adiciona al dto
							e := SetField(dto, fname, t)
							if e != nil {
								SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_string_error", "common_global_error", config)
							}
						}

						//Revisamos si el valor es diferente al de por defecto
						var tmp string
						if tmp != value {
							isEmpty = false
						}
					case time.Time:
						//Se adiciona al dto
						e := SetField(dto, fname, value)
						if e != nil {
							SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_string_error", "common_global_error", config)
						}

						//Revisamos si el valor es diferente al de por defecto
						var empty time.Time
						if empty != value {
							isEmpty = false
						}

					default:
						SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_format_error", "common_global_error", config)
					}
				case "date":
					switch value := value.(type) {
					case nil:
						if isRequired {
							SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_date_error", "common_global_error", config)
						}
					case string:
						t, error := time.Parse(dateFormat, value)
						var hasError bool = false
						if error != nil {
							//Significa que puede venir por defecto en zero
							t, error = time.Parse(time.RFC3339, value)
							if error != nil {
								hasError = true
								SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_date_error", "common_global_error", config)
							}
						} else if isRequired && value == common_config.DateTime.DATE_ZERO_VALUE {
							hasError = true
							SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_date_error", "common_global_error", config)
						}

						if !hasError && !t.IsZero() {
							//Se adiciona al dto
							e := SetField(dto, fname, t)
							if e != nil {
								SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_date_error", "common_global_error", config)
							}
						}

						//Revisamos si el valor es diferente al de por defecto
						var tmp string
						if tmp != value {
							isEmpty = false
						}
					case time.Time:
						//Se adiciona al dto
						e := SetField(dto, fname, value)
						if e != nil {
							SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_date_error", "common_global_error", config)
						}

						//Revisamos si el valor es diferente al de por defecto
						var empty time.Time
						if empty != value {
							isEmpty = false
						}

					default:
						SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_format_error", "common_global_error", config)
					}
				case "bool":
					switch value := value.(type) {
					case bool:
						e := SetField(dto, fname, value)
						if e != nil {
							SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_bool_error", "common_global_error", config)
						}

						var empty bool
						if empty != value {
							isEmpty = false
						}
					case string:
						var tmp string = strings.ToLower(value)
						var tmpBool bool = tmp == "true"

						if tmp != "true" && tmp != "false" {
							SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_bool_error", "common_global_error", config)
						} else {
							e := SetField(dto, fname, tmpBool)
							if e != nil {
								SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_bool_error", "common_global_error", config)
							}
						}
						isEmpty = false
					default:
						SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_format_error", "common_global_error", config)
					}
				case "uint":
					switch value := value.(type) {
					case map[string]interface{}:
						//Significa que es un objeto complejo.

						if _, found := value["icode"]; !found && isRequired {
							SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_uuid_error", "common_global_error", config)
						}
						if _, ok := value["icode"].(string); !ok && isRequired {
							SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_uuid_error", "common_global_error", config)
						}

						var icode string
						if value["icode"] != nil {
							icode = value["icode"].(string)
						}

						if len(icode) > 0 && len(icode) < 36 {
							SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_uuid_error", "common_global_error", config)
						}

						isEmpty = false
					case uint:
						tmp := uint64(value)
						if validateUintSize(fields, fname, tmp, isRequired, currentErrors, jsonFieldNameContainer, jsonTag, config) {
							e := SetField(dto, fname, tmp)
							if e != nil {
								SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_number_error", "common_global_error", config)
							}
						}

						isEmpty = false
					case uint16:
						tmp := uint64(value)
						if validateUintSize(fields, fname, tmp, isRequired, currentErrors, jsonFieldNameContainer, jsonTag, config) {
							e := SetField(dto, fname, tmp)
							if e != nil {
								SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_number_error", "common_global_error", config)
							}
						}

						isEmpty = false
					case uint32:
						tmp := uint64(value)
						if validateUintSize(fields, fname, tmp, isRequired, currentErrors, jsonFieldNameContainer, jsonTag, config) {
							e := SetField(dto, fname, tmp)
							if e != nil {
								SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_number_error", "common_global_error", config)
							}
						}

						isEmpty = false
					case uint64:
						tmp := value
						if validateUintSize(fields, fname, tmp, isRequired, currentErrors, jsonFieldNameContainer, jsonTag, config) {
							e := SetField(dto, fname, tmp)
							if e != nil {
								SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_number_error", "common_global_error", config)
							}
						}

						isEmpty = false
					case int:
						tmp := uint64(value)
						if validateUintSize(fields, fname, tmp, isRequired, currentErrors, jsonFieldNameContainer, jsonTag, config) {
							e := SetField(dto, fname, tmp)
							if e != nil {
								SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_number_error", "common_global_error", config)
							}
						}

						isEmpty = false
					case int16:
						tmp := uint64(value)
						if validateUintSize(fields, fname, tmp, isRequired, currentErrors, jsonFieldNameContainer, jsonTag, config) {
							e := SetField(dto, fname, tmp)
							if e != nil {
								SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_number_error", "common_global_error", config)
							}
						}

						isEmpty = false
					case int32:
						tmp := uint64(value)
						if validateUintSize(fields, fname, tmp, isRequired, currentErrors, jsonFieldNameContainer, jsonTag, config) {
							e := SetField(dto, fname, tmp)
							if e != nil {
								SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_number_error", "common_global_error", config)
							}
						}

						isEmpty = false
					case int64:
						tmp := uint64(value)
						if validateUintSize(fields, fname, tmp, isRequired, currentErrors, jsonFieldNameContainer, jsonTag, config) {
							e := SetField(dto, fname, tmp)
							if e != nil {
								SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_number_error", "common_global_error", config)
							}
						}

						isEmpty = false
					case float64:
						tmp := uint(value)
						tmp64 := uint64(tmp)

						if validateUintSize(fields, fname, tmp64, isRequired, currentErrors, jsonFieldNameContainer, jsonTag, config) {
							e := SetField(dto, fname, tmp64)
							if e != nil {
								SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_number_error", "common_global_error", config)
							}
						}

						isEmpty = false
					case string:
						if value == "" {
							value = "0"
						}
						tmp, err := strconv.ParseUint(value, 10, 64)

						if err != nil {
							SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_number_error", "common_global_error", config)
						} else {
							if validateUintSize(fields, fname, tmp, isRequired, currentErrors, jsonFieldNameContainer, jsonTag, config) {
								e := SetField(dto, fname, tmp)
								if e != nil {
									SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_number_error", "common_global_error", config)
								}
							}
						}
						isEmpty = false
					default:
						SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_format_error", "common_global_error", config)
					}
				case "int":
					switch value := value.(type) {
					case int:
						tmp := int64(value)
						if validateIntSize(fields, fname, tmp, isRequired, currentErrors, jsonFieldNameContainer, jsonTag, config) {
							e := SetField(dto, fname, tmp)
							if e != nil {
								SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_number_error", "common_global_error", config)
							}
						}

						isEmpty = false
					case int16:
						tmp := int64(value)
						if validateIntSize(fields, fname, tmp, isRequired, currentErrors, jsonFieldNameContainer, jsonTag, config) {
							e := SetField(dto, fname, tmp)
							if e != nil {
								SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_number_error", "common_global_error", config)
							}
						}

						isEmpty = false
					case int32:
						tmp := int64(value)
						if validateIntSize(fields, fname, tmp, isRequired, currentErrors, jsonFieldNameContainer, jsonTag, config) {
							e := SetField(dto, fname, tmp)
							if e != nil {
								SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_number_error", "common_global_error", config)
							}
						}

						isEmpty = false
					case int64:
						tmp := value
						if validateIntSize(fields, fname, tmp, isRequired, currentErrors, jsonFieldNameContainer, jsonTag, config) {
							e := SetField(dto, fname, tmp)
							if e != nil {
								SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_number_error", "common_global_error", config)
							}
						}

						isEmpty = false
					case uint:
						tmp := int64(value)
						if validateIntSize(fields, fname, tmp, isRequired, currentErrors, jsonFieldNameContainer, jsonTag, config) {
							e := SetField(dto, fname, tmp)
							if e != nil {
								SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_number_error", "common_global_error", config)
							}
						}

						isEmpty = false
					case uint16:
						tmp := int64(value)
						if validateIntSize(fields, fname, tmp, isRequired, currentErrors, jsonFieldNameContainer, jsonTag, config) {
							e := SetField(dto, fname, tmp)
							if e != nil {
								SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_number_error", "common_global_error", config)
							}
						}

						isEmpty = false
					case uint32:
						tmp := int64(value)
						if validateIntSize(fields, fname, tmp, isRequired, currentErrors, jsonFieldNameContainer, jsonTag, config) {
							e := SetField(dto, fname, tmp)
							if e != nil {
								SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_number_error", "common_global_error", config)
							}
						}

						isEmpty = false
					case uint64:
						tmp := int64(value)
						if validateIntSize(fields, fname, tmp, isRequired, currentErrors, jsonFieldNameContainer, jsonTag, config) {
							e := SetField(dto, fname, tmp)
							if e != nil {
								SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_number_error", "common_global_error", config)
							}
						}

						isEmpty = false
					case float64:
						tmp := int(value)
						tmp64 := int64(tmp)

						if validateIntSize(fields, fname, tmp64, isRequired, currentErrors, jsonFieldNameContainer, jsonTag, config) {
							e := SetField(dto, fname, tmp64)
							if e != nil {
								SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_number_error", "common_global_error", config)
							}
						}

						isEmpty = false
					case string:
						tmp, err := strconv.ParseInt(value, 10, 64)

						if err != nil {
							SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_number_error", "common_global_error", config)
						} else {
							if validateIntSize(fields, fname, tmp, isRequired, currentErrors, jsonFieldNameContainer, jsonTag, config) {
								e := SetField(dto, fname, tmp)
								if e != nil {
									SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_number_error", "common_global_error", config)
								}
							}
						}
						isEmpty = false
					default:
						SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_format_error", "common_global_error", config)
					}
				case "float":
					switch value := value.(type) {
					case float32:
						tmp := float64(value)
						if validateFloatSize(fields, fname, tmp, isRequired, currentErrors, jsonFieldNameContainer, jsonTag, config) {
							e := SetField(dto, fname, tmp)
							if e != nil {
								SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_number_error", "common_global_error", config)
							}
						}

						isEmpty = false
					case float64:
						tmp := value
						if validateFloatSize(fields, fname, tmp, isRequired, currentErrors, jsonFieldNameContainer, jsonTag, config) {
							e := SetField(dto, fname, tmp)
							if e != nil {
								SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_number_error", "common_global_error", config)
							}
						}
						isEmpty = false
					case int:
						tmp := float64(value)
						if validateFloatSize(fields, fname, tmp, isRequired, currentErrors, jsonFieldNameContainer, jsonTag, config) {
							e := SetField(dto, fname, tmp)
							if e != nil {
								SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_number_error", "common_global_error", config)
							}
						}

						isEmpty = false
					case int16:
						tmp := float64(value)
						if validateFloatSize(fields, fname, tmp, isRequired, currentErrors, jsonFieldNameContainer, jsonTag, config) {
							e := SetField(dto, fname, tmp)
							if e != nil {
								SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_number_error", "common_global_error", config)
							}
						}

						isEmpty = false
					case int32:
						tmp := float64(value)
						if validateFloatSize(fields, fname, tmp, isRequired, currentErrors, jsonFieldNameContainer, jsonTag, config) {
							e := SetField(dto, fname, tmp)
							if e != nil {
								SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_number_error", "common_global_error", config)
							}
						}

						isEmpty = false
					case int64:
						tmp := float64(value)
						if validateFloatSize(fields, fname, tmp, isRequired, currentErrors, jsonFieldNameContainer, jsonTag, config) {
							e := SetField(dto, fname, tmp)
							if e != nil {
								SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_number_error", "common_global_error", config)
							}
						}

						isEmpty = false
					case uint:
						tmp := float64(value)
						if validateFloatSize(fields, fname, tmp, isRequired, currentErrors, jsonFieldNameContainer, jsonTag, config) {
							e := SetField(dto, fname, tmp)
							if e != nil {
								SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_number_error", "common_global_error", config)
							}
						}

						isEmpty = false
					case uint16:
						tmp := float64(value)
						if validateFloatSize(fields, fname, tmp, isRequired, currentErrors, jsonFieldNameContainer, jsonTag, config) {
							e := SetField(dto, fname, tmp)
							if e != nil {
								SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_number_error", "common_global_error", config)
							}
						}

						isEmpty = false
					case uint32:
						tmp := float64(value)
						if validateFloatSize(fields, fname, tmp, isRequired, currentErrors, jsonFieldNameContainer, jsonTag, config) {
							e := SetField(dto, fname, tmp)
							if e != nil {
								SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_number_error", "common_global_error", config)
							}
						}

						isEmpty = false
					case uint64:
						tmp := float64(value)
						if validateFloatSize(fields, fname, tmp, isRequired, currentErrors, jsonFieldNameContainer, jsonTag, config) {
							e := SetField(dto, fname, tmp)
							if e != nil {
								SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_number_error", "common_global_error", config)
							}
						}

						isEmpty = false
					case string:
						tmp, err := strconv.ParseFloat(value, 64)

						if err != nil {
							SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_number_error", "common_global_error", config)
						} else {
							if validateFloatSize(fields, fname, tmp, isRequired, currentErrors, jsonFieldNameContainer, jsonTag, config) {
								e := SetField(dto, fname, tmp)
								if e != nil {
									SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_number_error", "common_global_error", config)
								}
							}
						}
						isEmpty = false
					default:
						SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_format_error", "common_global_error", config)
					}
				default:
					SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_format_error", "common_global_error", config)
				}

			} else if fields[fname].Required || isRequired {
				SetError(currentErrors, jsonFieldNameContainer, jsonTag, "common_validation_field_required_error", "common_global_error", config)
			}

		}
	}
	return isEmpty
}

func JSONToStruct(jsonStr string, target interface{}) {
	// Convertir el JSON string a bytes y intentar deserializar
	// Ignora los errores ya que se asume que habrá una validación especializada complementaria
	error := json.Unmarshal([]byte(jsonStr), target)
	println(error)
}

func GetDTOMap(jsonStr string, jsonFieldNameContainer string, locale map[string]map[string]string, collectedErrors map[string]map[string]string) interface{} {
	var dtoMap interface{}

	var hasError bool = CheckDataInput(&dtoMap, jsonStr, jsonFieldNameContainer, common_config.Enums.GLOBAL_ERROR, locale, collectedErrors)

	if hasError {
		SetError(collectedErrors, jsonFieldNameContainer, common_config.Enums.GLOBAL_ERROR, common_config.Enums.GLOBAL_ERROR, "", locale)
	}
	return dtoMap
}

func ValidateId(id string, entity string, fieldName string, config map[string]map[string]string) (map[string]map[string]string, error) {
	var collectedErrors map[string]map[string]string = map[string]map[string]string{}

	_, err := strconv.Atoi(id)
	if err != nil {
		SetError(collectedErrors, entity, fieldName, "common_validation_field_number_error", "common_global_error", config)
	}

	return collectedErrors, nil
}

func ValidateICode(id string, entity string, fieldName string, config map[string]map[string]string) (map[string]map[string]string, error) {
	var collectedErrors map[string]map[string]string = map[string]map[string]string{}

	_, err := uuid.Parse(id)

	if err != nil {
		SetError(collectedErrors, entity, fieldName, "common_validation_field_uuid_error", "common_global_error", config)
	}

	return collectedErrors, nil
}

/*
Asigna un valor de error al mapa. En caso de que no se encuentre el código, no hace nada.
Esto se hizo porque si el código del error no se encuentra, entonces podría generarse un balor en nil, cosa que no le gusta a go.
*/
func SetError(errorMap map[string]map[string]string, entity string, attr string, code string, generalMsgCode string, config map[string]map[string]string) {
	v, found := config["sp"][code]
	if found {
		_, found = errorMap[entity]
		if !found {
			errorMap[entity] = map[string]string{}
		}
		if generalMsgCode != "" {
			if _, found = errorMap["default"]; !found {
				errorMap["default"] = map[string]string{}
			}
			errorMap["default"][common_config.Enums.GLOBAL_MSG] = config["sp"][generalMsgCode]
		}
		errorMap[entity][attr] = v
	}
}

// Se utiliza para indicar que la entidad vino vacía, es decir, con los valores por defecto de los atributos. Es un string, pero el hecho de que exista ya muestra que está vacío. SI no existe es porque no está vacío.
func SetEmptyError(errorMap map[string]map[string]string, entity string, config map[string]map[string]string) {

	_, found := errorMap[entity]
	if !found {
		errorMap[entity] = map[string]string{}
	}
	errorMap[entity]["empty"] = "true"

}
func CheckDataInput(dtoMap *interface{}, dataInput string, JSONName string, globalError string, locale map[string]map[string]string, collectedErrors map[string]map[string]string) bool {
	var err error = json.Unmarshal([]byte(dataInput), &dtoMap)
	if err != nil {
		SetError(collectedErrors, JSONName, globalError, globalError, "", locale)
		return true
	}
	return false
}
func SetField(v interface{}, name string, value interface{}) error {
	// v must be a pointer to a struct
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
		return errors.New("v must be pointer to struct")
	}

	// Dereference pointer
	rv = rv.Elem()

	// Lookup field by name
	fv := rv.FieldByName(name)
	if !fv.IsValid() {
		return fmt.Errorf("not a field name: %s", name)
	}

	// Field must be exported
	if !fv.CanSet() {
		return fmt.Errorf("cannot set field %s", name)
	}

	/*// We expect a string field
	if fv.Kind() != reflect.String {
		return fmt.Errorf("%s is not a string field", name)
	}
	*/
	// Set the value
	switch fv.Kind() {
	case reflect.String:
		fv.SetString(value.(string))
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		fv.SetInt(value.(int64))
	case reflect.Bool:
		fv.SetBool(value.(bool))
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fv.SetUint(value.(uint64))
	case reflect.Float32, reflect.Float64:
		fv.SetFloat(value.(float64))
	case reflect.Struct:
		fv.Set(reflect.ValueOf(value))
	default:
		return fmt.Errorf("%s is not a known field", name)
	}

	return nil
}

func GetFieldValue(v interface{}, name string) interface{} {
	// v must be a pointer to a struct
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
		return errors.New("v must be pointer to struct")
	}

	// Dereference pointer
	rv = rv.Elem()

	// Lookup field by name
	fv := rv.FieldByName(name)
	if !fv.IsValid() {
		return nil
	}

	return fv.Interface()
}

/*
Retorna el valor del tag de nombre "tagName" correspondiente al atributo "attrName" en la estructura "dto"
*/
func GetTag(dto interface{}, attrName string, tagName string) string {
	field, ok := reflect.TypeOf(dto).Elem().FieldByName(attrName)
	if ok {
		return field.Tag.Get(tagName)
	}
	return ""
}

/*
Compara dos arreglos de tipo interface{} y devuelve true si son iguales
*/
func CompareInterfaceSlices(slice1 []interface{}, slice2 []interface{}) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	for i := 0; i < len(slice1); i++ {
		if slice1[i] != slice2[i] {
			return false
		}
	}
	return true
}

func validateStringSize(fields map[string]FieldDefinition, fieldName string, value string, isRequired bool, errors map[string]map[string]string, entity string, jsonName string, config map[string]map[string]string) bool {
	var empty string

	if empty == value && !isRequired {
		return true
	}

	if empty == value {
		SetError(errors, entity, jsonName, "common_validation_field_required_error", "common_global_error", config)
		return false
	}
	//Significa que no hay contorl de tamaño
	if fields[fieldName].MinSize == 0 && fields[fieldName].MaxSize == 0 {
		return true
	}
	if int64(len(value)) < fields[fieldName].MinSize {
		SetError(errors, entity, jsonName, "common_validation_field_string_min_size_error", "common_global_error", config)
		return false
	}
	if fields[fieldName].MaxSize > 0 && int64(len(value)) > fields[fieldName].MaxSize {
		SetError(errors, entity, jsonName, "common_validation_field_string_max_size_error", "common_global_error", config)
		return false
	}
	return true
}

func validateIntSize(fields map[string]FieldDefinition, fieldName string, value int64, isRequired bool, errors map[string]map[string]string, entity string, jsonName string, config map[string]map[string]string) bool {
	var empty int64

	if empty == value && (!isRequired || value == fields[fieldName].MinSize) {
		return true
	}

	if empty == value {
		SetError(errors, entity, jsonName, "common_validation_field_required_error", "common_global_error", config)
		return false
	}

	if fields[fieldName].MinSize == 0 && fields[fieldName].MaxSize == 0 {
		return true
	}

	if value < fields[fieldName].MinSize {
		SetError(errors, entity, jsonName, "common_validation_field_number_min_size_error", "common_global_error", config)
		return false
	}
	if value > fields[fieldName].MaxSize {
		SetError(errors, entity, jsonName, "common_validation_field_number_max_size_error", "common_global_error", config)
		return false
	}
	return true
}

func validateFloatSize(fields map[string]FieldDefinition, fieldName string, value float64, isRequired bool, errors map[string]map[string]string, entity string, jsonName string, config map[string]map[string]string) bool {
	var empty float64

	if empty == value && (!isRequired || value == float64(fields[fieldName].MinSize)) {
		return true
	}

	if empty == value {
		SetError(errors, entity, jsonName, "common_validation_field_required_error", "common_global_error", config)
		return false
	}

	if fields[fieldName].MinSize == 0 && fields[fieldName].MaxSize == 0 {
		return true
	}
	if value < float64(fields[fieldName].MinSize) {
		SetError(errors, entity, jsonName, "common_validation_field_number_min_size_error", "common_global_error", config)
		return false
	}
	if value > float64(fields[fieldName].MaxSize) {
		SetError(errors, entity, jsonName, "common_validation_field_number_max_size_error", "common_global_error", config)
		return false
	}
	return true
}

func validateUintSize(fields map[string]FieldDefinition, fieldName string, value uint64, isRequired bool, errors map[string]map[string]string, entity string, jsonName string, config map[string]map[string]string) bool {
	var empty uint64

	if empty == value && (!isRequired || value == uint64(fields[fieldName].MinSize)) {
		return true
	}

	if empty == value {
		SetError(errors, entity, jsonName, "common_validation_field_required_error", "common_global_error", config)
		return false
	}

	if fields[fieldName].MinSize == 0 && fields[fieldName].MaxSize == 0 {
		return true
	}
	if value < uint64(fields[fieldName].MinSize) {
		SetError(errors, entity, jsonName, "common_validation_field_number_min_size_error", "common_global_error", config)
		return false
	}
	if value > uint64(fields[fieldName].MaxSize) {
		SetError(errors, entity, jsonName, "common_validation_field_number_max_size_error", "common_global_error", config)
		return false
	}
	return true
}

func DTOtoString(dto interface{}) (string, error) {
	data, error := json.Marshal(dto)
	if error != nil {
		return "", errors.New(`{"` + common_config.Enums.GLOBAL_ERROR + `":"Error interno"}`)
	}
	return string(data), nil
}

func IsEmailValid(email string) bool {
	// Expresión regular para validar el email
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}
