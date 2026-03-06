package utils

import (
	common_config "bitsflow/common/config"
	"bitsflow/common/db"
	"encoding/json"
	"errors"
	"html/template"
	"io/fs"
	"math/rand/v2"
	"strconv"
	"unicode"

	"strings"
	"sync"

	"github.com/google/uuid"
)

const DEFAULT_PANIC_TEMPLATE string = "views/error.html"
const DEFAULT_VIEW string = "views/index.html"
const REPORT_VIEW string = "views/index_report.html"
const LOGIN_VIEW string = "views/login.html"
const DEFAULT_SCRIPTS string = "frontend/templates/layouts/standard_scripts.html"
const HTML_HELPER_SCRIPTS string = "frontend/templates/layouts/html_helper_scripts.html"
const HTML_LOGIN_SCRIPTS string = "frontend/templates/layouts/login_scripts.html"
const INPUT_TEMPLATE string = "frontend/templates/layouts_helpers/labeled_input.html"
const DATE_TIME_TEMPLATE string = "frontend/templates/layouts_helpers/date_time.html"
const TEXT_TEMPLATE string = "frontend/templates/layouts_helpers/text_translated.html"
const TEXT_AREA_TEMPLATE string = "frontend/templates/layouts_helpers/labeled_textarea.html"
const OVERLAY_TEMPLATE string = "frontend/templates/layouts/overlay.html"
const SELECT_ENUMS_TEMPLATE string = "frontend/templates/layouts_helpers/labeled_select_enums.html"
const STAR_RATING_TEMPLATE string = "frontend/templates/layouts_helpers/labeled_star_rating.html"
const DATALIST_ENUMS_TEMPLATE string = "frontend/templates/layouts_helpers/labeled_datalist_enums.html"
const SELECT_ENTITY_TEMPLATE string = "frontend/templates/layouts_helpers/labeled_select_entity.html"
const BUTTON_TEMPLATE string = "frontend/templates/layouts_helpers/button.html"

var lock = sync.RWMutex{}

/*
La separación de los tags viene con " | " de la siguiente forma:

	bf:"type:datetime | format:default | required:true | unique:false | default_value:auto

Pero podría venir así:

	bf:"type:datetime |format:default | required:true| unique:false|default_value:auto

La idea es que esta función normalice el tag para que quede así:

	bf:"type:datetime|format:default|required:true|unique:false|default_value:auto
*/
func NormalizeBFStructTag(tag string) string {
	re := strings.NewReplacer(" |", "|", "  |", "|", "| ", "|", "|  ", "|")
	return re.Replace(tag)
}

func MergeMaps[M1 ~map[K]V, M2 ~map[K]V, K comparable, V any](dst M1, src M2) {
	for k, v := range src {
		dst[k] = v
	}
}

func GetRandomString(size int) string {
	const (
		letterBytes   = "abcdefghijklmnopqrstuvwxyz1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%&*_-+="
		letterIdxBits = 6                    // 6 bits para representar un índice de letra
		letterIdxMask = 1<<letterIdxBits - 1 // Todos los bits de 1 para los índices de letra
	)
	sb := strings.Builder{}
	sb.Grow(size)

	// iterate over the indices
	for i := 0; i < size; i++ {
		// Obtener un índice aleatorio seguro del conjunto de caracteres
		idx := rand.IntN(len(letterBytes))
		sb.WriteByte(letterBytes[idx])
	}

	return sb.String()
}

func GetRandomInt(min, max int) int {
	return rand.IntN(max-min) + min
}

func GetFullHtmlTemplates() []string {
	return []string{
		BUTTON_TEMPLATE, OVERLAY_TEMPLATE, SELECT_ENTITY_TEMPLATE,
		SELECT_ENUMS_TEMPLATE, DATALIST_ENUMS_TEMPLATE, TEXT_TEMPLATE,
		INPUT_TEMPLATE, DATE_TIME_TEMPLATE, TEXT_AREA_TEMPLATE, DEFAULT_SCRIPTS,
		HTML_HELPER_SCRIPTS, DATALIST_ENUMS_TEMPLATE, STAR_RATING_TEMPLATE}
}
func GetFullHtmlFuncMap() template.FuncMap {

	var input HTMLInput = HTMLInput{}
	var textArea HTMLTextArea = HTMLTextArea{}
	var sel HTMLSelect = HTMLSelect{}
	var dl HTMLDatalist = HTMLDatalist{}
	var btn HTMLButton = HTMLButton{}
	var text HTMLText = HTMLText{}
	var starRating HTMLStarRating = HTMLStarRating{}
	var dateTime HTMLDateTime = HTMLDateTime{}

	funcMap := template.FuncMap{
		"Input":      input.New,
		"TextArea":   textArea.New,
		"Select":     sel.New,
		"DataList":   dl.New,
		"Button":     btn.New,
		"Text":       text.New,
		"StarRating": starRating.New,
		"DateTime":   dateTime.New,
	}
	return funcMap
}

func GetLoginHtmlTemplates() []string {
	return []string{BUTTON_TEMPLATE, OVERLAY_TEMPLATE, INPUT_TEMPLATE, DATE_TIME_TEMPLATE, HTML_LOGIN_SCRIPTS, TEXT_TEMPLATE}
}

func GetLoginHtmlFuncMap() template.FuncMap {

	var input HTMLInput = HTMLInput{}
	var btn HTMLButton = HTMLButton{}
	var text HTMLText = HTMLText{}
	var sel HTMLSelect = HTMLSelect{}

	funcMap := template.FuncMap{
		"Input":  input.New,
		"Button": btn.New,
		"Text":   text.New,
		"Select": sel.New,
	}
	return funcMap
}

func GetUUID() string {
	lock.RLock()
	defer lock.RUnlock()

	uuid7, err := uuid.NewV7()
	var uuidStr string
	if err != nil {
		uuidStr = uuid.New().String()
	} else {
		uuidStr = uuid7.String()
	}
	return uuidStr
}

func GeneratePassword(length int, language string) (string, error) {
	if length < 8 {
		return "", errors.New(common_config.Locale[language]["common_pass_num_chars_fail"])
	}

	// Caracteres para cada tipo
	lower := "abcdefghijklmnopqrstuvwxyz"
	upper := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits := "0123456789"
	specials := "!@#$%^&*()-_=+<>?{}[]|~"

	allChars := lower + upper + digits + specials
	password := make([]byte, length)

	// Aseguramos que tenga al menos un carácter de cada tipo
	password[0] = lower[GetRandomInt(0, len(lower))]
	password[1] = upper[GetRandomInt(0, len(upper))]
	password[2] = digits[GetRandomInt(0, len(digits))]
	password[3] = specials[GetRandomInt(0, len(specials))]

	// Llenamos el resto de la contraseña
	for i := 4; i < length; i++ {
		password[i] = allChars[GetRandomInt(0, len(allChars))]
	}

	// Mezclamos los caracteres de la contraseña
	for i := range password {
		j := GetRandomInt(0, len(password))
		password[i], password[j] = password[j], password[i]
	}

	return string(password), nil
}

// Función para verificar la seguridad de una contraseña y reportar sus debilidades
func CheckPasswordStrength(password string) []string {
	var issues []string

	if len(password) < 8 {
		issues = append(issues, "common_pass_num_chars_fail")
	}

	hasLower := false
	hasUpper := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasLower {
		issues = append(issues, "common_pass_lower_case_missing")
	}
	if !hasUpper {
		issues = append(issues, "common_pass_upper_case_missing")
	}
	if !hasDigit {
		issues = append(issues, "common_pass_number_missing")
	}
	if !hasSpecial {
		issues = append(issues, "common_pass_special_char_missing")
	}

	return issues
}

func LoadDBCLientConfig() db.DBClientConfig {
	data, err := fs.ReadFile(ConfigAssets, "config/db_config.json")
	var config db.DBClientConfig = db.DBClientConfig{}
	if err == nil {
		json.Unmarshal([]byte(data), &config)
	}
	return config
}

func ParseUint64(u string) uint64 {
	v, err := strconv.ParseUint(u, 10, 64)
	if err != nil {
		panic(err)
	}
	return v
}

func NilIfZero[T comparable](v T) interface{} {
	var zero T
	if v == zero {
		return nil
	}
	return v
}
