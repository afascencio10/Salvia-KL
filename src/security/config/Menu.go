package security_config

//Permisos generales según el rol.
var PermissionsByRole map[string]map[string]bool = map[string]map[string]bool{
	"get_general_users": {
		"ad": true,
		"do": true,
	},
	"get_general_user": {
		"ad": true,
		"do": true,
	},
	"set_general_user": {
		"ad": true,
		"do": true,
	},
	"update_general_user": {
		"ad": true,
		"do": true,
	},
	"update_general_user_disable": {
		"ad": true,
		"do": true,
	},
	"update_general_user_enable": {
		"ad": true,
		"do": true,
	},
	"remove_general_user": {
		"ad": true,
	},
	//City
	"get_city_by_department": {
		"op": true,
		"ue": true,
		"et": true,
		"ad": true,
		"ro": true,
		"do": true,
		"fo": true,
	},
	"get_town_by_city_code": {
		"op": true,
		"ue": true,
		"et": true,
		"ad": true,
		"ro": true,
		"do": true,
		"fo": true,
	},
}

//Menús del sistema. Se filtra por el idioma, luego por el rol y finalmente quedan las secciones con sus respectivos items.
var Menu map[string]map[string]map[string][]map[string]string = map[string]map[string]map[string][]map[string]string{
	"sp": {
		"ad": {
			Locale["sp"]["menu_title_general_user"]: {
				{
					"label": Locale["sp"]["menu_get_general_users"],
					"path":  "/" + Locale["sp"]["security"] + "/" + Locale["sp"]["GeneralUser"],
				},
				{
					"label": Locale["sp"]["set_general_user"],
					"path":  "/" + Locale["sp"]["security"] + "/" + Locale["sp"]["GeneralUser"] + "/" + Locale["sp"]["new"],
				},
			},
		},
		"do": {
			Locale["sp"]["menu_title_general_user"]: {
				{
					"label": Locale["sp"]["menu_get_general_users"],
					"path":  "/" + Locale["sp"]["security"] + "/" + Locale["sp"]["GeneralUser"],
				},
				{
					"label": Locale["sp"]["set_general_user"],
					"path":  "/" + Locale["sp"]["security"] + "/" + Locale["sp"]["GeneralUser"] + "/" + Locale["sp"]["new"],
				},
			},
		},
	},
}

//Menús secundarios del sistema para definir las acciones de edición, eliminación, etc. Se filtra por el idioma, luego por el rol, nombre del menú y finalmente sus respectivos items.
var MenuTools map[string]map[string]map[string][]map[string]string = map[string]map[string]map[string][]map[string]string{
	"sp": {
		"ad": {
			"menu_tool_get_general_users": {
				{
					"label":       Locale["sp"]["menu_tool_get_general_user"],
					"path":        "/" + Locale["sp"]["security"] + "/" + Locale["sp"]["GeneralUser"],
					"action":      "",
					"customClass": "fas fa-eye icoacciones",
				},
				{
					"label":       Locale["sp"]["menu_tool_update_general_user"],
					"path":        "/" + Locale["sp"]["security"] + "/" + Locale["sp"]["GeneralUser"],
					"action":      Locale["sp"]["update"],
					"customClass": "fas fa-edit icoacciones",
				},
				{
					"label":       Locale["sp"]["menu_tool_disable_general_user"],
					"path":        "/" + Locale["sp"]["security"] + "/" + Locale["sp"]["GeneralUser"],
					"action":      Locale["sp"]["router_update_general_user_disable"],
					"customClass": "fas fa-archive unable",
					"assert":      "user.status == 'e'",
				},
				{
					"label":       Locale["sp"]["menu_tool_enable_general_user"],
					"path":        "/" + Locale["sp"]["security"] + "/" + Locale["sp"]["GeneralUser"],
					"action":      Locale["sp"]["router_update_general_user_enable"],
					"customClass": "fas fa-archive unable",
					"assert":      "user.status == 'd'",
				},
			},
		},
		"do": {
			"menu_tool_get_general_users": {
				{
					"label":       Locale["sp"]["menu_tool_get_general_user"],
					"path":        "/" + Locale["sp"]["security"] + "/" + Locale["sp"]["GeneralUser"],
					"action":      "",
					"customClass": "fas fa-eye icoacciones",
				},
				{
					"label":       Locale["sp"]["menu_tool_update_general_user"],
					"path":        "/" + Locale["sp"]["security"] + "/" + Locale["sp"]["GeneralUser"],
					"action":      Locale["sp"]["update"],
					"customClass": "fas fa-edit icoacciones",
				},
				{
					"label":       Locale["sp"]["menu_tool_disable_general_user"],
					"path":        "/" + Locale["sp"]["security"] + "/" + Locale["sp"]["GeneralUser"],
					"action":      Locale["sp"]["router_update_general_user_disable"],
					"customClass": "fas fa-archive unable",
					"assert":      "user.status == 'e'",
				},
				{
					"label":       Locale["sp"]["menu_tool_enable_general_user"],
					"path":        "/" + Locale["sp"]["security"] + "/" + Locale["sp"]["GeneralUser"],
					"action":      Locale["sp"]["router_update_general_user_enable"],
					"customClass": "fas fa-archive unable",
					"assert":      "user.status == 'd'",
				},
			},
		},
	},
}

//Inicializamos los elementos
