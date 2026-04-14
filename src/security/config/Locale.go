package security_config

var Locale = map[string]map[string]string{
	//Validaciones GeneralUser
	"sp": {
		"public":                               "public",
		"security_general_user_password_error": "La contraseña debe tener una o más letras en minúsculas, mayúsculas, uno o más números y se sugiere que también contenga uno o más caracteres especiales como: ! @ # $ % & * ( ) _ - = + / \\ ? < > . , ;",
		"security_general_user_password_do_not_match":   "La contraseña y su confirmación no coinciden",
		"security_general_user_password_internal_error": "Se produjo un error interno al cifrar la contraseña, por favor contacte a un administrador",
		"security_general_user_unique":                  "Ya existe un usuario con este dato, por favor escriba uno diferente",
		"security_email_unique":                         "Ya existe un usuario con este correo, por favor escriba uno diferente",
		"security_phone_unique":                         "Ya existe un usuario con este número telefónico, por favor escriba uno diferente",
		"security_role_not_found":                       "Por favor seleccione un rol válido",
		"login_fail":                                    "Usuario o contraseña incorrectos",
		"security":                                      "seguridad",
		"GeneralUser":                                   "usuarios",
		"City":                                          "ciudad",
		"Town":                                          "centro_poblado",
		"login":                                         "ingreso",
		"new":                                           "nuevo",
		"update":                                        "actualizar",
		"remove":                                        "eliminar",
		"disable":                                       "deshabilitar",
		"general_user_save":                             "Guardar",
		"general_user_remove":                           "Eliminar",
		"general_user_disable":                          "Inhabilitar",
		"general_user_enable":                           "Habilitar",
		"general_user_cancel":                           "Cancelar",
		"general_user_back":                             "Regresar sin guardar",
		"security_general_user_captcha_fail":            "CAPTCHA incorrecto, por favor intenta escribir nuevamente los números de la imagen",

		"security_general_user_reset_token_not_found": "No se encontró el token",
		"security_general_user_reset_mail_not_found":  "El usuario no tiene correo de notificación. Por favor solicite al administrador que actualice los datos del usuario.",

		"login_window_title":               "Iniciar sesión",
		"get_general_user_window_title":    "Usuarios",
		"update_general_user_window_title": "Modificar usuario",
		"remove_general_user_window_title": "Eliminar usuario",
		"set_general_user_window_title":    "Crear usuario",

		"global_error":   "Se produjo un error general",
		"internal_error": "Se produjo un error interno desconocido",

		"router_get_city_by_department": "departamento",
		"router_get_town_by_city_code":  "codigo_ciudad",
		//Menú
		"menu_title_general_user":       "Usuarios",
		"menu_get_general_users":        "Listar usuarios",
		"menu_get_general_users_by_all": "Listar todos los usuarios",
		"set_general_user":              "Crear nuevo usuario",
		//Login
		"user_name":                "Nombre de usuario",
		"user_password":            "Contraseña",
		"user_name_h":              "Escriba un nombre de usuario",
		"user_password_h":          "Escriba una contraseña",
		"login_forgot_pass_title":  "Restablecer contraseña",
		"login_forgot_pass_msg":    "Si usted es funcionario, por favor escriba el nombre de usuario asignado en la plataforma.",
		"login_forgot_pass_back":   "Cerrar",
		"login_forgot_pass_save":   "Restablecer contraseña",
		"login_forgot_pass_text":   "!Olvidé mi contraseña!",
		"login_forgot_success_msg": "Se ha enviado un correo electrónico con las instrucciones para restablecer su contraseña.",

		"login_reset_success_msg_btn_back": "Cerrar",

		"login_reset_pass_title": "Nueva contraseña",
		"login_reset_pass_msg":   "Por favor escriba la nueva contraseña para su cuenta",
		"login_reset_pass_back":  "Cerrar",
		"login_reset_pass_save":  "Guardar",

		"login_reset_pass_no_email": "El usuario no tiene correo de notificación",

		"captcha_solution":   "Verificación",
		"captcha_solution_h": "Escribe los caracteres que ves en la imagen",

		//GeneralUser
		"general_user_profile_email":     "Correo electrónico para notificaciones",
		"general_user_profile_email_h":   "Escriba un correo electrónico válido",
		"general_user_creation_date":     "Fecha de creación",
		"general_user_creation_date_h":   "Fecha en la que se creó el usuario",
		"general_user_update_date":       "Fecha de última actualización",
		"general_user_update_date_h":     "Fecha de la última actualización del usuario",
		"general_user_login":             "Nombre de usuario",
		"general_user_login_h":           "Escriba un Nombre de usuario",
		"general_user_password":          "Contraseña",
		"general_user_password_h":        "Mínimo 8 dígitos con caracteres variados",
		"general_user_password_repeat":   "Confirmación de contraseña",
		"general_user_password_repeat_h": "Mínimo 8 dígitos con caracteres variados",
		"general_user_status":            "Estado",
		"general_user_status_h":          "Seleccione un estado",
		"general_user_language":          "Idioma",
		"general_user_language_h":        "Seleccione un Idioma",
		"general_user_roles":             "Roles de usuario",
		"general_user_roles_h":           "Seleccione los roles de usuario",
		"general_user_select_empty":      "",

		"menu_tool_get_general_user":         "Ver",
		"menu_tool_update_general_user":      "Actualizar",
		"menu_tool_remove_general_user":      "Eliminar",
		"menu_tool_disable_general_user":     "Inhabilitar",
		"menu_tool_enable_general_user":      "Habilitar",
		"router_update_general_user_disable": "deshabilitar_usuario",
		"router_update_general_user_enable":  "habilitar_usuario",
		"get_department_users_empty_town":    "Se produjo un error interno. El usuario no tiene asociada una ubicación georgráfica.",
		//GeneralProfile
		"general_user_profile_creation_date":     "Fecha de creación",
		"general_user_profile_creation_date_h":   "Fecha en la que se creó el perfíl de usuario",
		"general_user_profile_update_date":       "Fecha de última actualización",
		"general_user_profile_update_date_h":     "EFecha de la última actualización del perfíl de usuario",
		"general_user_profile_gender":            "Género",
		"general_user_profile_gender_h":          "Seleccione un Género",
		"general_user_profile_nick":              "Nick / Nombre identitario",
		"general_user_profile_nick_h":            "Escriba una Nick / Nombre identitario",
		"general_user_profile_names":             "Nombres",
		"general_user_profile_names_h":           "Escriba nombres completos",
		"general_user_profile_last_names":        "Apellidos",
		"general_user_profile_last_names_h":      "Escriba un apellidos completos",
		"general_user_profile_doc_type":          "Tipo de documento",
		"general_user_profile_doc_type_h":        "Seleccione un tipo de documento",
		"general_user_profile_doc_number":        "Número de documento",
		"general_user_profile_doc_number_h":      "Escriba un número de documento",
		"general_user_profile_phones":            "Números de teléfono",
		"general_user_profile_phones_h":          "Escriba los números de teléfono",
		"general_user_profile_mails":             "Correos electrónicos",
		"general_user_profile_mails_h":           "Escriba los correos electrónicos",
		"general_user_profile_description":       "Descripción / nota adicional",
		"general_user_profile_description_h":     "Escriba una descripción / nota adicional",
		"security_general_user_profile_town":     "Por favor seleccione un municipio",
		"security_general_user_profile_branch":   "Por favor seleccione una entidad",
		"security_general_user_branch_not_found": "No se encontró la entidad",
	},
}

func TranslateLocale(targetValue string, language string) string {

	for key, value := range Locale[language] {
		if value == targetValue {
			return key
		}
	}
	return ""
}

/*
*
Traduce y retorna una regla de navegación y la convierte en algo como:

	nav["save"]["/usuarios/perfil"]
*/
func TranslateNavigationRule(lang string, navRule map[string]map[string]string) map[string]string {

	var resNav map[string]string = map[string]string{}

	for submit, params := range navRule {

		locale, found := Locale[lang]
		if found {
			mod, found := locale[params["module"]]
			if found {
				ent, found := locale[params["entity"]]
				if found {
					resNav[submit] = "/" + mod + "/" + ent
				}
			}
		}
	}

	return resNav
}
