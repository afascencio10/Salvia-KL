package salvia_config

const NUM_ITEMS_PER_PAGE int = 50

var HTML_Templates_folder string = "frontend/salvia/html/"

var HTML_Templates map[string]string = map[string]string{
	"home":                     "home.html",
	"get_my_follow_ups":        "get_my_follow_ups.html",
	"get_my_cases":             "get_my_cases.html",
	"list_cases":               "list_cases.html",
	"get_victim_cases":         "get_victim_cases.html",
	"get_victim_cases_ro":      "get_victim_cases_ro.html",
	"get_victim_cases_do":      "get_victim_cases_do.html",
	"report_victim_cases":      "report_victim_cases.html",
	"get_victim_cases_et":      "get_victim_cases_et.html",
	"get_victim_cases_us":      "get_victim_cases_us.html",
	"get_victim_cases_sv":      "get_victim_cases_sv.html",
	"get_victim_cases_no":      "get_victim_cases_no.html",
	"get_victim_cases_fo":      "get_victim_cases_fo.html",
	"get_victim_case_fo":       "get_victim_case_fo.html",
	"get_victim_case_fo_v1":    "get_victim_case_fo_v1.html",
	"get_victim_case_op":       "get_victim_case_op.html",
	"get_victim_case_op_v1":    "get_victim_case_op_v1.html",
	"get_victim_case_no_v1":    "get_victim_case_no_v1.html",
	"get_victim_case_ro":       "get_victim_case_ro.html",
	"get_victim_case_ro_v1":    "get_victim_case_ro_v1.html",
	"get_victim_case_sv":       "get_victim_case_sv.html",
	"get_victim_case_sv_v1":    "get_victim_case_sv_v1.html",
	"get_victim_case_no":       "get_victim_case_no.html",
	"get_victim_case_et":       "get_victim_case_et.html",
	"get_victim_case_et_v1":    "get_victim_case_et_v1.html",
	"get_victim_case_us":       "get_victim_case_us.html",
	"get_victim_case_us_v1":    "get_victim_case_us_v1.html",
	"get_victim_case_do":       "get_victim_case_do.html",
	"set_victim_case":          "set_victim_case.html",
	"update_victim_case_v1":    "update_victim_case_v1.html",
	"update_victim_case_v2":    "update_victim_case_v2.html",
	"get_victim_contacts":      "get_victim_contacts.html",
	"get_victim_contact":       "get_victim_contact.html",
	"get_victim_contact_v1":    "get_victim_contact_v1.html",
	"get_victim_contacts_sv":   "get_victim_contacts_sv.html",
	"get_victim_contact_sv_v2": "get_victim_contact_sv_v2.html",
	"get_victim_contact_sv_v1": "get_victim_contact_sv_v1.html",
	"set_victim_contact":       "set_victim_contact.html",
	"get_alerts":               "get_alerts.html",
	"load_plain_files":         "load_plain_files.html",
	"assign_operators":         "assign_operators.html",
	"get_entity_branches":      "get_entity_branches.html",
	"set_feminicide":           "set_feminicide.html",
	"set_feminicide_risk":      "set_feminicide_risk.html",
	"hacer_seguimiento":        "hacer_seguimiento.html",
	"get_follow_up_detail":     "get_follow_up_detail.html",
	"get_case_detail_sv":       "get_case_detail_sv.html",
	"get_feminicide_risks":     "get_feminicide_risks.html",
	"get_feminicide_risk":      "get_feminicide_risk.html",
	"notificaciones":           "notificaciones.html",
	"mis_barreras":             "barrieris.html",
	"barrera_detalle":          "barrera_detalle.html",
	"kpis_dashboard":           "kpis_dashboard.html",
}

var FormPaths map[string]map[string]string = map[string]map[string]string{
	"sp": {
		//Vitim contact
		"VictimContactGET":             "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["VictimContact"],
		"VictimContactPOST_Public":     "/" + Locale["sp"]["public"] + "/" + Locale["sp"]["VictimContact"],
		"VictimContactPOST_GET_Public": "/" + Locale["sp"]["public"] + "/" + Locale["sp"]["VictimContact"] + "/" + Locale["sp"]["new"],
		"VictimContactPOST_GET":        "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["VictimContact"] + "/" + Locale["sp"]["new"],
		"VictimContactPUT":             "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["VictimContact"],
		//Vitim case
		"VictimCasePOST_GET":  "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["VictimCase"] + "/" + Locale["sp"]["new"],
		"VictimCaseReportGET": "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["VictimCase"] + "/" + Locale["sp"]["report"],
		"VictimCasePUT":       "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["VictimCase"],
		"VictimCaseGET":       "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["VictimCase"],

		//EntityBranch
		"EntityBranchGET": "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["EntityBranch"],
		//Moment
		"MomentPUT": "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["Moment"],
		//Alert
		"AlertGET": "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["Alert"],
		//CaseLog
		"CaseLogPOST": "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["CaseLog"],
		//PlainFiles
		"PlainFilesPOST": "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["LoadPlainFiles"],
		"PlainFilesGET":  "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["LoadPlainFiles"] + "/" + Locale["sp"]["new"],
		//AssignOperators
		"AssignOperatorsPOST": "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["AssignOperators"],
		"AssignOperatorsGET":  "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["AssignOperators"],
		//FollowUp
		"FollowUpPUT":  "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["FollowUp"],
		"FollowUpPOST": "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["FollowUp"],
		"FollowUpGET":  "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["FollowUp"],

		//FollowUpEntry
		"FollowUpEntryPUT":  "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["FollowUpEntry"],
		"FollowUpEntryPOST": "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["FollowUpEntry"],
		"FollowUpEntryGET":  "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["FollowUpEntry"],

		//FollowUpEntry
		"FollowUpEntryActingPOST": "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["FollowUpEntryActing"],

		//Barrier
		"BarrierGET": "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["Barrier"],

		//Entity
		"EntityGET": "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["Entity"],

		//Feminicide
		"FeminicideGET": "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["Feminicide"],

		//FeminicideRisk
		"FeminicideRiskGET":  "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["FeminicideRisk"],
		"FeminicideRiskPOST": "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["FeminicideRisk"],
	},
}

/*
*
Indica el origen y el destino. Pr ejemplo, en el siguiente código indica que del template "get_general_user" para la acción (submit) "any" (Cualquier acción) se redirige al home.
NAVIGATION_RULES["get_general_user"]["asalvia_config
*/
var NAVIGATION_RULES map[string]map[string]map[string]string = map[string]map[string]map[string]string{
	"home":                {"default": {"module": "", "entity": ""}},
	"get_victim_case":     {"default": {"module": "salvia", "entity": "VictimCase"}},
	"get_feminicide":      {"default": {"module": "salvia", "entity": "Feminicide"}},
	"get_alerts":          {"default": {"module": "salvia", "entity": "Alert"}, "home": {"module": "salvia", "entity": "VictimCase"}},
	"set_victim_case":     {"default": {"module": "salvia", "entity": "VictimCase"}, "cancel": {"module": "salvia", "entity": "VictimCase"}},
	"get_entity_branches": {"default": {"module": "salvia", "entity": "EntityBranch"}, "cancel": {"module": "salvia", "entity": "EntityBranch"}},
	"report_victim_cases": {"default": {"module": "salvia", "entity": "VictimCase"}},
	"update_victim_case":  {"default": {"module": "salvia", "entity": "VictimCase"}, "cancel": {"module": "salvia", "entity": "VictimCase"}},
	"update_moment":       {"default": {"module": "salvia", "entity": "VictimCase"}, "attend": {"module": "salvia", "entity": "VictimCase"}},
	"get_victim_contact":  {"default": {"module": "", "entity": ""}, "invalidate": {"module": "salvia", "entity": "VictimContact"}},
	"set_victim_contact":  {"default": {"module": "salvia", "entity": "VictimContact"}, "cancel": {"module": "salvia", "entity": "VictimContact"}},
	"load_plain_files":    {"default": {"module": "salvia", "entity": "VictimCase"}},
	"assign_operators":    {"default": {"module": "salvia", "entity": "VictimCase"}},
	"login":               {"login": {"module": "salvia", "entity": "VictimCase"}, "cancel": {"module": "salvia", "VictimCase": ""}},
	"get_follow_up":       {"default": {"module": "salvia", "entity": "FollowUp"}},
}

var LIVING_ZONES map[string]string = map[string]string{
	"cm": "Cabecera municipal",
	"cp": "Centro poblado",
	"rd": "Rural disperso",
}

var SEXUAL_ORIENTATION map[string]string = map[string]string{
	"he": "Heterosexual",
	"lb": "Lesbiana",
	"gy": "Gay",
	"bs": "Bisexual",
	"ot": "Otro ¿Cuál?",
}

var ORIGIN_PLACE map[string]string = map[string]string{
	"r": "Rural",
	"u": "Urbano",
}

var OCCUPATION map[string]string = map[string]string{
	"st": "Estudiante",
	"hs": "Trabajadora sector salud",
	"es": "Trabajadora sector educación",
	"ps": "Trabajadora sector público / funcionaria pública",
	"ls": "Lideresa social",
	"pe": "Periodista",
	"dr": "Trabajo doméstico remunerado",
	"dp": "Trabajo doméstico no remunerado",
	"ue": "Desempleada",
	"rt": "Jubilada",
	"pt": "Campesino",
	"ts": "Trabajadora sexual",
	"ot": "Otra ¿Cuál?",
}

var ETHNIC_GROUP map[string]string = map[string]string{
	"in": "Indígena",
	"af": "Afrodescendientes",
	"rg": "ROM",
	"no": "Ninguno",
}

var SECTOR map[string][]map[string]string = map[string][]map[string]string{
	"sp": {
		{"code": "he", "label": Locale["sp"]["sector_he"], "img": "fas fa-stethoscope fa-2x"},
		{"code": "pt", "label": Locale["sp"]["sector_pt"], "img": "fa-regular fa-folder fa-2x"},
		{"code": "js", "label": Locale["sp"]["sector_js"], "img": "fa-regular fa-hand-back-fist fa-2x"},
		{"code": "pm", "label": Locale["sp"]["sector_pm"], "img": "fa-regular fa-hand-back-fist fa-2x"},
		{"code": "os", "label": Locale["sp"]["sector_os"], "img": "fa-regular fa-book-open-reader fa-2x"},
	},
}

var BARRIER_SECTOR map[string][]map[string]string = map[string][]map[string]string{
	"sp": {
		{"code": "he", "label": Locale["sp"]["sector_he"], "img": "fas fa-stethoscope fa-2x"},
		{"code": "pt", "label": Locale["sp"]["sector_pt"], "img": "fa-regular fa-folder fa-2x"},
		{"code": "js", "label": Locale["sp"]["sector_js"], "img": "fa-regular fa-hand-back-fist fa-2x"},
		{"code": "pm", "label": Locale["sp"]["sector_pm"], "img": "fa-regular fa-hand-back-fist fa-2x"},
		{"code": "os", "label": Locale["sp"]["sector_os"], "img": "fa-regular fa-book-open-reader fa-2x"},
	},
}

var MOMENT map[string][]map[string]string = map[string][]map[string]string{
	"sp": {
		{"code": "01", "label": Locale["sp"]["moment_01"], "img": "fas fa-eye fa-2x"},
		{"code": "02", "label": Locale["sp"]["moment_02"], "img": "fas fa-recycle fa-2x"},
		{"code": "03", "label": Locale["sp"]["moment_03"], "img": "fas fa-people-arrows fa-2x"},
	},
}

var DISABILITY map[string]string = map[string]string{
	"no": "Ninguna",
	"ad": "Auditiva",
	"vi": "Visual",
	"mo": "Motora",
	"co": "Cognitiva",
	"at": "Autismo",
	"mu": "Múltiple",
}

var OCCURRENCE map[string]string = map[string]string{
	"f": "Primera vez",
	"r": "Recurrente",
}

var WEEK_DAY map[string]string = map[string]string{
	"1": "Lunes",
	"2": "Martes",
	"3": "Miércoles",
	"4": "Jueves",
	"5": "Viernes",
	"6": "Sábado",
	"7": "Domingo",
}

var VIOLENCE_SCOPE map[string]string = map[string]string{
	"fl": "Familiar convivivente",
	"fn": "Familiar no convivivente",
	"pn": "Pareja",
	"p2": "Expareja",
	"fs": "Amistad",
	"cm": "Comunitario",
	"ht": "Salud",
	"la": "Laboral",
	"it": "Institucional",
	"jl": "Reclusión Intramural",
	"pi": "Instituciones de protección",
	"pt": "Transporte Público",
	"ed": "Educativo",
	"ps": "Espacio público",
	"dg": "Digital",
	"co": "Establecimientos de comercio",
	"cb": "Cibernético",
	"po": "Político",
	"nr": "Sin relación",
	"ch": "Cohabitación, no familiar ni expareja",
	"ac": "Conflicto armado",
}

var VIOLENCE_EXPERIENCED map[string]string = map[string]string{
	"ph": "Física",
	"ps": "Psicológica",
	"sx": "Sexual",
	"pm": "Patrimonial o económica",
	"vc": "Vicaria",
	"re": "Reproductiva",
	"po": "Política",
	"ot": "Otras",
}

var AGGRESSOR map[string]string = map[string]string{
	"k": "Conocido",
	"u": "Desconocido",
}
var RELATIONSHIP_WITH_AGGRESSOR map[string]string = map[string]string{
	"fl": "Familiar convivivente",
	"fn": "Familiar no convivivente",
	"fr": "Amigo",
	"ex": "Expareja",
	"pn": "Pareja permanente/temporal",
	"wp": "Compañero de trabajo",
	"ng": "Vecino",
	"sp": "Compañero de estudios",
	"pf": "Amigo de la pareja",
	"fm": "Familiar",
	"bs": "Jefe",
	"bf": "Novio",
	"fd": "Familiar diferente a la pareja",
	"ot": "Otro",
	"na": "Ninguno",
}

var MOMENT_APPROVAL_SOURCE map[string]string = map[string]string{
	"AUTOMATIC": "a",
	"MANUAL":    "m",
}

var VICTIM_CASE_STATUS map[string]map[string]string = map[string]map[string]string{
	"sp": {
		"fc": "primer contacto",
		"r":  "enrutado por aprobar",
		"ra": "enrutado aprobado",
		"c":  "completado",
		"ex": "vencido",
		"is": "con novedades",
		"cd": "cerrado",
		"iv": "inválido",
	},
}

var VICTIM_CONTACT_STATUS map[string]map[string]string = map[string]map[string]string{
	"sp": {
		"v": "Válido",
		"i": "Inválido",
	},
}

var CASE_OWNER_STATUS map[string]map[string]string = map[string]map[string]string{
	"sp": {
		"a": "Activo",
		"i": "Inactivo",
	},
}

var ALERT_TYPE map[string]string = map[string]string{
	"CASE_LOG":        "case_log",
	"MOMENT_LOG":      "moment_log",
	"CASE_EXPIRED":    "case_expired",
	"FEMINICIDE_RISK": "feminicide_risk",
}

var ALERT_TYPE_LOCALE map[string]map[string]string = map[string]map[string]string{
	"sp": {
		"case_log":        "Novedad del caso",
		"moment_log":      "Novedad del momento del caso",
		"case_expired":    "Caso expiró",
		"feminicide_risk": "Riesgo feminicida",
	},
}

var VICTIM_CASE_VICTIM_NATIONALITY map[string]string = map[string]string{
	"c": "Colombiana",
	"f": "Extranjera ¿Cuál?",
}

var VICTIM_CASE_VICTIM_FOREIGNER_IMMIGRATION_STATUS map[string]string = map[string]string{
	"r": "Regular",
	"i": "Irregular",
}

var VICTIM_CASE_VICTIM_GENDER map[string]string = map[string]string{
	"m": "Hombre",
	"w": "Mujer",
	"i": "Intersexual",
}

var VICTIM_CASE_VICTIM_DEPENDENTS map[string]string = map[string]string{
	"m": "Madre",
	"f": "Padre",
	"s": "Hermanos",
	"n": "Ninguna",
	"o": "Otro ¿Cuál?",
}

var VICTIM_CASE_VICTIM_IF_PREVIOUSLY_REPORTED map[string]string = map[string]string{
	"cf": "Comisaría de Familia",
	"fi": "Fiscalía",
	"ml": "Medicina Legal",
	"is": "IPS",
	"ip": "Inspección de policía",
	"ot": "Otro",
	"dr": "No recuerdo",
}

var AFRO_COMMUNTITIES map[string]map[string]string = map[string]map[string]string{
	"sp": {
		"b": "Negro",
		"p": "Palenquero",
		"r": "Raizal",
	},
}

var INDIGENOUS_TONGUES map[string]map[string]string = map[string]map[string]string{
	"sp": {
		"cb": "CUBEO",
		"cu": "CUIBA",
		"yu": "YUHUP",
		"si": "SIRIANO",
		"ac": "ACHAGUA",
		"wa": "WAUNANA",
		"ec": "EMBERA CHAMI",
		"gu": "GUAYABERO",
		"wi": "WIWA",
		"tk": "TUKANO",
		"ty": "TUYUCA",
		"ch": "CHIMILA",
		"sa": "SÁLIBA",
		"de": "DESANO",
		"ed": "EMBERA DOBIDA",
		"ti": "TIKUNA",
		"gr": "Gitanos-Rom",
		"yc": "YUCUNA",
		"ma": "MAKUNA",
		"pu": "PUINAVE",
		"na": "NASAYUWE",
		"hi": "HITNU",
		"sk": "SIKUANI",
		"ek": "EMBERA KATIO",
		"ka": "KAMSÁ",
		"pi": "PISAMIRA",
		"ta": "TATUYO",
		"wn": "WANANO",
		"mi": "MIRAÑA",
		"no": "NONUYA",
		"tn": "TANIMUKA",
		"yr": "YURUTI",
		"ya": "YARURO",
		"yg": "YAGUA",
		"co": "COCAMA",
		"wy": "WAYUUNAIKI",
		"ku": "KURRIPACO",
		"ba": "BARASANA",
		"br": "BARI",
		"tr": "TARIANO",
		"bo": "BORA",
		"pa": "PALENQUE",
		"in": "INGA",
		"so": "SIONA",
		"ko": "KOREGUAJE",
		"ca": "CARIJONA",
		"un": "KUNA",
		"oc": "OCAINA",
		"tw": "TAIWANO",
		"pr": "PIAROA",
		"pt": "PIRATAPUYO",
		"ia": "IKA ARHUACO",
		"cf": "COFÁN",
		"aw": "AWAPIT",
		"ar": "BARÁ",
		"kg": "KOGUI",
		"ui": "UITOTO",
		"tg": "TINIGUA",
		"uw": "U'WA",
		"an": "ANDOQUE",
		"nu": "NUKAK",
		"nm": "NAMTRIK",
		"cc": "CACUA",
		"cr": "CREOLE",
		"ci": "CABIYARI",
		"mu": "MUINANE",
		"pp": "PIAPOCO",
		"kr": "KARAPANA",
	},
}

var COLOMBIAN_INDIGENOUS map[string]map[string]string = map[string]map[string]string{
	"sp": {
		"ac": "ACHAGUA (Achagua,ajagua,xagua)",
		"am": "AMORÚA (Wipiwe,Siripu,Mariposa)",
		"wi": "WIPIJIWI0WAÜPIJIWI",
		"ya": "YAMALERO (Yamalero)",
		"yr": "YARURO (Pume)",
		"an": "ANDOQUE (Andoque-andoque,cha’oie,businka)",
		"ar": "ARHUACO (Arhuacos,ika,iku,Ijku-Ijka)",
		"ww": "WIWA (Wiwa, Arzario,Guamaca,Malayo,Sanjá,Dumana)",
		"ba": "BARÁ ( waimaja,posanga0mira. Barasana del Norte)",
		"br": "BARASANO (Barasano del sur,eduria,yebá0masã,yepa0mahsã,yepá0matsó,hanerã (o janena),paneroa,komea,teiuana (o taiwano),banera yae,hanera oka)",
		"bi": "BARÍ (Motilone barí,motilón,Barís,barira,dobocubi,cunausaya)",
		"be": "BETOYE (Jirarre, Betoi,Jirara,Guahibos,Betoyes)",
		"bo": "BORA (Meamuyna)",
		"ka": "KAWIYARI (Kawiyarí,Kawiarí,Kabiyarí,Cabiyari)",
		"ca": "CARAPANA (Ucomaja-Karapana-Moxdoa,Muxtea,Mitea,Mochda,Karapaná,Karapano,Carapana-tapuya,Mextä,Mehta)",
		"kr": "KARIJONA (Carijona,Carifuna,Hianacoto0umaua,Kaliohona)",
		"ce": "CHIMILAS-ETTE ENEKA (Chimila-ett E'neka-simiza,chimile y shimizya)",
		"ch": "CHIRICOA (Chiricoa Guahibo)",
		"co": "COCAMA (Cocama-kokama,cocama,ucayali,xibitaoan,huallaga,pampadeque,pandequebo,omagua)",
		"ko": "KOKONUCO (Coconuco,Puracé)",
		"ke": "KOREGUAJE (Koreguaje-Coreguaje,korebaju,coreguaxe o Koré pâín)",
		"pi": "PIJAO (Pijao-Coyaima,Natagaima. Pixaos,Pyjaos y Pinaos)",
		"aw": "AWÁ (cuaiquer,kwaiker)",
		"ku": "KUBEO (Pamíva,0Kubeo,Paniwa,Cobewa,Hipnwa,Kaniwa)",
		"cu": "CUIBA",
		"gu": "GUANADULE0TULE0CUNA (Guanadule,Tula,Cuna0Tule,cuna,kuna,tacarcuna,cerracuna,darienes)",
		"cr": "CURRIPAKO (Kurripaco,Baniva,Waquenia,Kurrupaku)",
		"bn": "BANIVA (Kurripaco,Baniva,Waquenia,Kurrupaku)",
		"ga": "GUARIQUEMA",
		"de": "DESANO (Desana,Uina,Winá,Uira,Wirá boleka,Oregua,Kusibi,Wirá,Kotedia,Dessana)",
		"ta": "TAMA DUJO (Tamas0Dujo,dujos. )",
		"em": "EMBERA",
		"eb": "EMBERA EYABIDA-EMBERA KATÍO (Embera Katío-Katío,catío,katio,embena,eyabida.)",
		"ec": "EMBERA CHAMI (Embera Chamí )",
		"ep": "EPERARA SIAPIDARA (Eperara,saija,epená saija,epea pedée,cholo)",
		"ed": "EMBERA DOBIDÁ (Emberá Dobidá0Mokaná,macaná,Dóvida)",
		"nu": "NUTABE (Nutabes)",
		"mi": "MISAK (Misak)",
		"ab": "AMBALÓ (Ambaló )",
		"qi": "QUIZGÓ (Quizgó)",
		"gn": "GUANACO (Guanaca)",
		"go": "GUANANO",
		"ji": "JIW0GUAYABERO (guayabero,piapoco,bisanigua,cunimía,mitúa,mítiwa)",
		"cm": "CAÑAMOMO",
		"in": "INGA (Ingano)",
		"km": "KAMËNTSA (kamsá,camsá,sibundoy-gache Camentsa )",
		"kf": "KOFÁN (Kofán-Cofan,Kofane)",
		"kg": "KOGUI (Kogui,Kaggabba,Kogui,Cogui)",
		"le": "LETUAMA (Lituana,Detuama,Wejeñeme maja)",
		"ma": "MAKAGUAJE (Makaguaje0Macaguaje,Macaguaxe,Airubain)",
		"hi": "HITNÜ,MACAGUÁN (Hitnü,Macaguán-Hitnü,macaguane,jitnu,macaguán,hutnun)",
		"mk": "MAKUNA (Macuna,Sara,Ide masa,Buhagana,Siroa,Tsoloa)",
		"ju": "JUDPA-JUJUPDA",
		"kk": "KAKUA (kakua,kakwa,cacua)",
		"hu": "HUPDË-HUPDAH-HUPDU",
		"jh": "JUHUP-YUJU",
		"jd": "JUDPA",
		"nk": "NUKAK",
		"mb": "MAIBÉN MASIWARE-PODIPODI (Masiware0Maibén.)",
		"mt": "MATAPÍ (Matapí,jupichiya,upichia)",
		"je": "JE'ERURIWA",
		"mr": "MIRAÑA (Miraña0 Miranas wacho Améjimínaa. mirnha,miraya,wacho améjimínaa)",
		"mu": "MUISCA (Muisca-Chibcha)",
		"no": "NONUYA (Nonuya,Nononota)",
		"ok": "OKAINA (Dyo'xaiya0o0Ivo'tsa)",
		"na": "NASA (Naza Paéz0Nasa Yuwe)",
		"po": "POLINDARA",
		"pa": "PIAPOCO (Dzase,dejá,kuipaco,wenéwika,enegua,yapoco,amarizado)",
		"pr": "PIAROA (Wotiheh,Uhothuha, Wóthuha o Dearwa)",
		"pt": "PIRATAPUYO (Piratapuyo-piratapuyo,piratapuya,parata puya,uaikama,wai kana,waikhana,urubú-tapuy)",
		"ps": "PISAMIRA (Pisamira0Papiwa,pasatapuyo,wasona,wasina)",
		"pu": "PUINAVE (Puinave Guaipunare,puinabe,uapi,wantyinht)",
		"as": "PASTO (Pastos)",
		"qu": "QUILLACINGA",
		"sa": "SÁLIBA (Salivá)",
		"si": "SIKUANI",
		"mp": "MAPAYERRI",
		"so": "SIONA (Siona-Ganteyabain,ganteya,ceona,zeona,kokakanú o Katucha0Pai)",
		"tu": "TUBÚ-SIRIANO (Siriano-Sura masa,Cirnga,Chiranga,Si0Ra - Ciriana,Ciriano. )",
		"ti": "TAIWANO-EDURIA (Taiwano)",
		"tn": "TANIMUCA (Tanimuka ufania,tanimuca,taniboka,ohañara,opaima)",
		"tg": "TANIGUA (Tinigwa,Tinigua)",
		"tr": "TARIANO (Tariana,Taliaseri)",
		"tt": "TATUYO (Tatuyo,Juna maja,Pamoa,Tatutapuyo,Sina.)",
		"to": "TOTORÓ (Tontotuna0Totoroe)",
		"tk": "TIKUNA",
		"ts": "TSIRIPU (Tsiripo0Mariposo)",
		"uk": "TUKANO (Tukano0Ye'pámahsa)",
		"uw": "U'WA (U'wa 0 uwkuwa.)",
		"ty": "TUYUCA (Tatuyo,Juna maja,Pamoa,Tatutapuyo,Sina.)",
		"wo": "WOUNAAN (Waunana,Wounaan,Noanamá,Waumeu)",
		"wy": "WAYUÚ (Wayú,Guajiros)",
		"ui": "UITOTO (Muina Murui-Witotos-Huitoto,Witoto,Murui,Muinane,Mi0ka,Huitoto,Mi0pode. Wuitotos-Uitotos)",
		"mm": "MUINA MURUI (Muina Murui-Witotos-Huitoto,Witoto,Murui,Muinane,Mi0ka,Huitoto,Mi0pode. Wuitotos-Uitotos)",
		"yi": "YARÍ (Yari )",
		"yg": "YAGUA (Ñihamwo,Mishara.)",
		"yk": "YANAKONA-YANAKUNA (Yanacona,Yanakuna,Yanacuna.)",
		"yn": "YAUNA (Yahuna)",
		"un": "YUKUNA (Yucuna,Yukuna Matapí, Kamejeya)",
		"kp": "YUKPA (Yuko,Yuco,Yukpa.)",
		"yt": "YURUTI (otsoca,Wadyana,Wadzana,Wai jiara masa-wadzana,Waikana)",
		"ze": "ZENÚ (Zenú,senú)",
		"ua": "GUANE",
		"kn": "MOKANÁ",
		"ot": "OTAVALEÑO",
		"qc": "QUICHUA (Quichua 0 Kechua,Quechua, Kichwa)",
		"k1": "KANKUAMO (Kankankuamos,Kankuamos)",
		"ay": "TAYRONAS",
		"it": "CHITARERO",
		"qb": "QUIMBAYA",
		"cl": "CALIMA",
		"pn": "PANCHES",
		"ie": "INDIGENAS-ECUADOR (DIFERENTE DE OTAVALEÑOS)",
		"ip": "INDIGENAS-PERU",
		"iv": "INDIGENAS-VENEZUELA",
		"im": "INDIGENAS-MÉXICO",
		"ib": "INDIGENAS-BRASIL",
		"ye": "YERAL",
		"ii": "INDIGENAS-PANAMA",
		"nb": "INDIGENAS-BOLIVIA",
		"gm": "MAYA-(GUATEMALA)",
		"ni": "INDÍGENA-SIN-INFORMACIÓN",
	},
}

var VIOLENCE_SCENES map[string]map[string]string = map[string]map[string]string{
	"sp": {
		"ed": "Establecimiento educativo",
		"is": "Institución de salud",
		"ad": "Áreas de deporte y recreación",
		"cc": "Calles y carreteras",
		"cs": "Comercio y áreas de servicios",
		"ai": "Áreas industriales o de construcción",
		"gf": "Granjas o fincas",
		"oe": "Otros espacios abiertos",
		"oc": "Otros espacios cerrados: bares, restaurantes",
		"di": "Digital",
		"ip": "Institución de protección",
		"ri": "Resguardo indígena",
		"co": "Consejo comunitario",
	},
}

var RISK_LEVEL map[string]map[string]string = map[string]map[string]string{
	"sp": {
		"h": "Alto",
		"m": "Medio",
		"l": "Bajo",
	},
}

var FOLLOW_UP_STATUS map[string]map[string]string = map[string]map[string]string{
	"sp": {
		"p": "En progreso",
		"d": "Realizado",
	},
}

var FOLLOW_UP_ENTRY_STATUS map[string]map[string]string = map[string]map[string]string{
	"sp": {
		"p": "En progreso",
		"d": "Realizado",
		"u": "Sin iniciar",
		"x": "Expirado",
	},
}

var FOLLOW_UP_ENTRY_ACTING_STATUS map[string]map[string]string = map[string]map[string]string{
	"sp": {
		"p": "En progreso",
		"d": "Realizado",
		"x": "Expirado",
		"e": "Escalado",
	},
}
