package common_config

type date_time struct {
	DATE_TIME_FORMAT              string
	DB_DATE_TIME_FORMAT           string
	DATE_FORMAT                   string
	DATE_ZERO_VALUE               string
	DB_DATE_FORMAT                string
	TIME_FORMAT                   string
	TIME_ZERO_VALUE               string
	TIME_WITH_SECONDS_FORMAT      string
	TIME_WITH_MILLISECONDS_FORMAT string
}

var DateTime date_time = date_time{
	DATE_TIME_FORMAT:              "02/01/2006 15:04:05",
	DB_DATE_TIME_FORMAT:           "2006/01/02 15:04:05",
	DATE_FORMAT:                   "02/01/2006",
	DB_DATE_FORMAT:                "2006/01/02",
	DATE_ZERO_VALUE:               "01/01/0001",
	TIME_FORMAT:                   "15:04",
	TIME_ZERO_VALUE:               "00:00",
	TIME_WITH_SECONDS_FORMAT:      "15:04:05",
	TIME_WITH_MILLISECONDS_FORMAT: "15:04:05.000000",
}

type enums struct {
	INTERNAL_ERROR string
	GLOBAL_ERROR   string
	GLOBAL_MSG     string
}

var Enums enums = enums{
	INTERNAL_ERROR: "internal_error",
	GLOBAL_ERROR:   "global_error",
	GLOBAL_MSG:     "global_msg",
}

var GENDER_IDENTITY map[string]string = map[string]string{
	"nb": "No binaria - asignado femenino al nacer",
	"n2": "No binaria - asignado masculino al nacer",
	"ma": "Hombre",
	"fe": "Mujer",
	"tf": "Mujer Transgénero",
	"tm": "Hombre Transgénero",
	"ot": "Otra ¿Cuál?",
}

var GENDER map[string]string = map[string]string{
	"m": "Hombre",
	"w": "Mujer",
	"i": "Intersexual",
}

var YES_NO map[string]string = map[string]string{
	"y": "Sí",
	"n": "No",
}

var YES_NO_NA map[string]string = map[string]string{
	"y": "Sí",
	"n": "No",
	"a": "No Aplica",
}

var MARITAL_STATUS map[string]string = map[string]string{
	"si": "Soltera",
	"ma": "Casada",
	"fu": "Unión libre / Unión marital de hecho",
	"wd": "Viuda",
	"fv": "Divorciada",
	"ot": "Otro",
}

var DOCUMENT_TYPE map[string]string = map[string]string{
	"cc": "Cédula de Ciudadanía",
	"ce": "Cédula de Extranjería",
	"ti": "Tarjeta de identidad",
	"te": "Tarjeta de Extranjería",
	"ps": "Pasaporte",
	"ni": "NIT",
	"dn": "DNI - documento identidad extranjero",
	"np": "NIT de otro país",
	"rc": "Registro civil",
	"sc": "Salvoconducto refugiados SVC",
	"pp": "Permiso especial de permanencia - PEP",
	"pt": "Permiso de protección temporal",
	"sn": "Sin información SN",
}

var DOCUMENT_TYPE_FORM2 map[string]string = map[string]string{
	"cd": "Carné Diplomático",
	"cc": "Cédula de Ciudadanía",
	"ce": "Cédula de Extranjería",
	"no": "Certificado de Nacido Vivo del país de origen",
	"nc": "Certificado de Nacido Vivo en Colombia",
	"mv": "Certificado de Registro Administrativo de Migrantes Venezolanos - RAMV",
	"dn": "Documento Extranjero o Documento Nacional de Identidad (del país de origen -DNI )",
	"oe": "Otro no especificado",
	"pn": "Partida de Nacimiento de su país de origen",
	"ps": "Pasaporte",
	"pp": "Permiso Especial de Permanencia (PEP)",
	"pf": "Permiso Especial de Permanencia para el Fomento de la Formalización (PEP FF)",
	"pt": "Permiso por Protección Temporal",
	"si": "Persona sin identificación",
	"rc": "Registro Civil de Nacimiento",
	"sc": "Salvoconducto",
	"ti": "Tarjeta Identidad",
	"vs": "Visa",
	"vr": "Visa de Refugiado",
	"np": "NIT de otro país",
}

var LANGUAGE map[string]string = map[string]string{
	"sp": "Español",
}
