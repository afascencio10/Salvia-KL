package common_dao

import (
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	"fmt"
	"strconv"
	"strings"
)

const (
	SQL_SELECT             string = "SELECT"
	SQL_UPDATE             string = "UPDATE"
	SQL_DELETE             string = "DELETE"
	SQL_INSERT             string = "INSERT"
	SQL_SELECT_FIELDS_ONLY string = "SELECT_ONLY_FIELDS"
	SQL_UPDATE_FIELDS_ONLY string = "UPDATE_ONLY_FIELDS"
	SQL_SELECT_WHERE_ONLY  string = "SELECT_WHERE_ONLY"
	SQL_AND                string = "AND"
	SQL_OR                 string = "OR"
)

/*
Ayuda a formar SQLs simples.

@action Especifica la acción que se quiere realizar para formar el SQL:

	SQL_SELECT: Un select completo
	SQL_UPDATE: Un update completo
	SQL_DELETE: Un delete completo
	SQL_INSERT: Un insert completo
	SQL_SELECT_FIELDS_ONLY: Sólo genera la lista de campos del select para cuando el Where es generado externamente
	SQL_UPDATE_FIELDS_ONLY: Sólo genera la lista de campos con su asignación para cuando el Where es generado externamente

@fields Tiene la lista de campos que se quieren usar dentro del select. Estos son los nombres de atributos de la clase. Él lo traduce internamente a los nombres en la BD
@fieldsAlias En el mismo orden que fields, pero con los alias para la consulta
@from nombre de la tabla de BD
@where lista de atributos que participarán del where
@returning lista de atributos que se espera como respuesta del query

@whereLogic operador lógico que se usará en el where

	SQL_AND: Se operará todos los atributos en where con AND
	SQL_OR: Se operará todos los atributos en where con OR

@scheme Esquema de BD
@definitions mapa con las definiciones de la clase. Allí estará el mapeo del nombre del atributo con los nombres en BD
@fieldsSbsolutePath true si se quiere poner en los atributos toda la ruta de scheme + el nombre de la tabla + el atributo
*/
func GetSQL(action string, fields []string, fieldsAlias []string, from string, where []string, equalSign []string, returning []string, whereLogic string, scheme string, definitions map[string]utils.FieldDefinition, absolutePath bool) string {

	var query string
	var fieldSlice []string
	var returningSlice []string

	var whereStr string
	var valuesStr string

	switch action {
	case SQL_SELECT:

		for i, v := range fields {
			var f string = definitions[v].DBName
			if absolutePath {
				f = scheme + "." + from + "." + f
			}
			if len(fieldsAlias) > 0 && fieldsAlias[i] != "" {
				f += " AS " + fieldsAlias[i]
			} else if definitions[v].Alias != "" {
				f += " AS " + definitions[v].Alias
			}
			fieldSlice = append(fieldSlice, f)
		}
		if whereLogic != "" {
			for i, v := range where {
				var f string = definitions[v].DBName
				if absolutePath {
					f = scheme + "." + from + "." + f
				}
				if len(equalSign) == 0 {
					whereStr += f + " = $" + fmt.Sprint(i+1) + " "
				} else {
					var value string = " $" + fmt.Sprint(i+1)

					switch equalSign[i] {
					case "%LIKE%":
						value = "%" + value + "%"
					case "%LIKE":
						value = "%" + value
					case "LIKE%":
						value = value + "%"
					}
					whereStr += f + " " + equalSign[i] + value + " "
				}

				if i < len(where)-1 {
					if whereLogic == SQL_AND {
						whereStr += " AND "
					} else {
						whereStr += " OR "
					}
				}
			}
		} else {
			whereStr = " TRUE "
		}

		query = `SELECT ` + strings.Join(fieldSlice, ", ") + ` FROM ` + scheme + `.` + from + `
				WHERE ` + whereStr

	case SQL_INSERT:

		for i, v := range fields {
			var f string = definitions[v].DBName
			if absolutePath {
				f = scheme + "." + from + "." + f
			}
			fieldSlice = append(fieldSlice, f)

			valuesStr += "$" + fmt.Sprint(i+1) + " "
			if i < len(fields)-1 {
				valuesStr += ", "
			}
		}
		for _, v := range returning {
			var f string = definitions[v].DBName
			if absolutePath {
				f = scheme + "." + from + "." + f
			}
			returningSlice = append(returningSlice, f)
		}
		var returningStr string = ``
		if len(returning) > 0 {
			returningStr += `RETURNING ` + strings.Join(returningSlice, ",")
		}

		query = `INSERT INTO ` + scheme + `.` + from + ` (` + strings.Join(fieldSlice, ", ") + `) VALUES (` + valuesStr + `) ` + returningStr

	case SQL_UPDATE:
		var i int
		var v string
		if whereLogic != "" {
			for i, v = range where {
				var f string = definitions[v].DBName
				if absolutePath {
					f = scheme + "." + from + "." + f
				}
				whereStr += f + " = $" + fmt.Sprint(i+1) + " "
				if i < len(where)-1 {
					if whereLogic == SQL_AND {
						whereStr += " AND "
					} else {
						whereStr += " OR "
					}
				}
			}
			i++
		} else {
			whereStr = " FALSE "
		}

		for j, v := range fields {
			var f string = definitions[v].DBName
			if absolutePath {
				f = scheme + "." + from + "." + f
			}
			valuesStr += f + " = $" + fmt.Sprint(i+1) + " "
			if j < len(fields)-1 {
				valuesStr += ", "
			}
			i++
		}

		query = `UPDATE ` + scheme + `.` + from + ` SET ` + valuesStr + ` WHERE ` + whereStr

	case SQL_DELETE:

		if whereLogic != "" {
			for i, v := range where {
				var f string = definitions[v].DBName
				if absolutePath {
					f = scheme + "." + from + "." + f
				}
				whereStr += f + " = $" + fmt.Sprint(i+1) + " "
				if i < len(where)-1 {
					if whereLogic == SQL_AND {
						whereStr += " AND "
					} else {
						whereStr += " OR "
					}
				}
			}
		} else {
			whereStr = " TRUE "
		}

		query = `DELETE FROM ` + scheme + `.` + from + `
				WHERE ` + whereStr

	case SQL_SELECT_FIELDS_ONLY:

		for i, v := range fields {
			var f string = definitions[v].DBName
			if absolutePath {
				f = scheme + "." + from + "." + f
			}
			if len(fieldsAlias) > 0 && fieldsAlias[i] != "" {
				f += " AS " + fieldsAlias[i]
			} else if definitions[v].Alias != "" {
				f += " AS " + definitions[v].Alias
			}
			fieldSlice = append(fieldSlice, f)
		}
		query = strings.Join(fieldSlice, ", ")

	case SQL_UPDATE_FIELDS_ONLY:

		for i, v := range fields {
			var f string = definitions[v].DBName
			if absolutePath {
				f = scheme + "." + from + "." + f
			}
			valuesStr += f + " = $" + fmt.Sprint(i) + " "
			if i < len(fields)-1 {
				valuesStr += ", "
			}
		}

		query = valuesStr

	case SQL_SELECT_WHERE_ONLY:

		if whereLogic != "" {
			for i, v := range where {
				var f string = definitions[v].DBName
				if absolutePath {
					f = scheme + "." + from + "." + f
				}

				if len(equalSign) == 0 {
					whereStr += f + " = $" + fmt.Sprint(i+1) + " "
				} else {
					var value string = " $" + fmt.Sprint(i+1)

					switch equalSign[i] {
					case "%LIKE%":
						value = "%" + value + "%"
					case "%LIKE":
						value = "%" + value
					case "LIKE%":
						value = value + "%"
					}
					whereStr += f + " " + equalSign[i] + value + " "
				}

				if i < len(where)-1 {
					if whereLogic == SQL_AND {
						whereStr += " AND "
					} else {
						whereStr += " OR "
					}
				}
			}
		} else {
			whereStr = " FALSE "
		}

		query = `WHERE ` + whereStr

	}
	return " " + query + " "
}

func GetOffsetQuery(page int) string {
	var pageQuery string = ""
	if page >= 0 {
		var offset int = page * salvia_config.NUM_ITEMS_PER_PAGE

		pageQuery = ` LIMIT ` + strconv.Itoa(salvia_config.NUM_ITEMS_PER_PAGE) + ` OFFSET ` + strconv.Itoa(offset) + ` `
	}

	return pageQuery
}
