package common_config

var Locale = map[string]map[string]string{
	"sp": {
		Enums.GLOBAL_ERROR:             "Se ha producido un error inesperado",
		"common_op_success":            "Operación realizada exitosamente",
		"common_db_connection_success": "Conectado a base de datos OK",
		"common_db_max_conn_reached":   "Se alcanzó el número máximo de conexiones disponibles",
		"common_db_get_conn_fail":      "Error al obtener una nueva conexión: ",
		"common_db_get_conn_fail_2":    "Error al obtener una conexión existente ",
		"common_db_realease_conn_fail": "Error al liberar una conexión ",
		"common_db_no_rows_affected":   "Error al realizar la operación ",

		//Validación genérica de atributos
		"common_validation_field_format_error":          "Campo con formato inesperado",
		"common_validation_field_empty_error":           "Este campo es obligatorio",
		"common_validation_field_string_min_size_error": "El número de cacarteres debe ser mayor",
		"common_validation_field_string_max_size_error": "Ha excedido el número máximo de cacarteres",
		"common_validation_field_number_min_size_error": "El número debe ser mayor",
		"common_validation_field_number_max_size_error": "El número debe ser menor",
		"common_validation_field_email_error":           "Digite un correo válido",
		"common_validation_field_number_error":          "Digite un número válido",
		"common_validation_field_date_error":            "Digite una fecha válida",
		"common_validation_field_date_future_error":     "Esta fecha no puede establecerse en el futuro",
		"common_validation_field_string_error":          "Digite una secuencia de caracteres válida",
		"common_validation_field_bool_error":            "Digite un booleano válido",
		"common_validation_field_required_error":        "Este campo es requerido",
		"common_validation_field_uuid_error":            "Digite un identificador válido",
		"common_validation_field_json_field_not_found":  "No se encontró el campo especificado en el JSON",
		"common_global_error":                           "Se presentaron los siguientes errores",
		"common_pass_num_chars_fail":                    "La contraseña debe tener al menos 8 caracteres.",
		"common_pass_lower_case_missing":                "La contraseña debe contener al menos una letra minúscula.",
		"common_pass_upper_case_missing":                "La contraseña debe contener al menos una letra mayúscula.",
		"common_pass_number_missing":                    "La contraseña debe contener al menos un número.",
		"common_pass_special_char_missing":              "La contraseña debe contener al menos un caracter especial como estos: $%&*!@.-_(){}[]",
	},
}
