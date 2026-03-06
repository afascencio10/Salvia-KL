package security_config

import security_daos "bitsflow/security/dao"

var HTML_Templates_folder string = "frontend/html/"

const NUM_ITEMS_PER_PAGE int = 50

var HTML_Templates map[string]string = map[string]string{
	"home":                 "home.html",
	"login":                "login.html",
	"get_general_users":    "get_general_users.html",
	"get_general_users_do": "get_general_users_do.html",
	"get_general_user":     "get_general_user.html",
	"set_general_user":     "set_general_user.html",
	"set_general_user_do":  "set_general_user_do.html",
	"update_general_user":  "update_general_user.html",
	"disable_general_user": "disable_general_user.html",
	"remove_general_user":  "remove_general_user.html",
}

/*
*
Indica el origen y el destino. Pr ejemplo, en el siguiente código indica que del template "get_general_user" para la acción (submit) "any" (Cualquier acción) se redirige al home.
NAVIGATION_RULES["get_general_user"]["any"] = "home"
*/
var NAVIGATION_RULES map[string]map[string]map[string]string = map[string]map[string]map[string]string{
	"home":                {"default": {"module": "", "entity": ""}},
	"get_general_user":    {"default": {"module": "security", "entity": "GeneralUser"}},
	"set_general_user":    {"default": {"module": "security", "entity": "GeneralUser"}, "cancel": {"module": "security", "entity": "GeneralUser"}},
	"set_general_user_do": {"default": {"module": "security", "entity": "GeneralUser"}, "cancel": {"module": "security", "entity": "GeneralUser"}},
	"update_general_user": {"default": {"module": "security", "entity": "GeneralUser"}, "cancel": {"module": "security", "entity": "GeneralUser"}},
	"remove_general_user": {"default": {"module": "security", "entity": "GeneralUser"}, "cancel": {"module": "security", "entity": "GeneralUser"}},
	"load_plain_files":    {"default": {"module": "security", "entity": "GeneralUser"}, "cancel": {"module": "security", "entity": "GeneralUser"}},
	"login":               {"login": {"module": "security", "entity": "GeneralUser"}, "cancel": {"module": "security", "entity": ""}},
}

var USER_STATUS map[string]string = map[string]string{
	"e": "Activo",
	"d": "Inhabilitado",
}

var FormPaths map[string]map[string]string = map[string]map[string]string{
	"sp": {
		//GeneralUser
		"GeneralUserPOST_Public":     "/" + Locale["sp"]["public"] + "/" + Locale["sp"][security_daos.GeneralUserEntityName],
		"GeneralUserPOST_GET_Public": "/" + Locale["sp"]["public"] + "/" + Locale["sp"][security_daos.GeneralUserEntityName] + "/" + Locale["sp"]["new"],
		"GeneralUserPOST_GET":        "/" + Locale["sp"]["security"] + "/" + Locale["sp"][security_daos.GeneralUserEntityName] + "/" + Locale["sp"]["new"],
		"GeneralUserGET":             "/" + Locale["sp"]["security"] + "/" + Locale["sp"][security_daos.GeneralUserEntityName],

		//City
		"CityGET":        "/" + Locale["sp"]["security"] + "/" + Locale["sp"][security_daos.CityEntityName],
		"CityGET_Public": "/" + Locale["sp"]["public"] + "/" + Locale["sp"][security_daos.CityEntityName],
		//Town
		"TownGET":        "/" + Locale["sp"]["security"] + "/" + Locale["sp"][security_daos.TownEntityName],
		"TownGET_Public": "/" + Locale["sp"]["public"] + "/" + Locale["sp"][security_daos.TownEntityName],
	},
}
