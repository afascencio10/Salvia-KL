// Package security_ctrl — lógica para el endpoint que recibe un Excel y ejecuta
// la migración masiva de usuarios directamente en el servidor.
package security_ctrl

import (
	common_controllers "bitsflow/common/controllers"
	common_dao "bitsflow/common/dao"
	"bitsflow/common/db"
	security_daos "bitsflow/security/dao"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"unicode"

	excelize "github.com/xuri/excelize/v2"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// ─────────────────────────────────────────────────────────────────────────────
// Nombre de hoja y columnas del Excel
// ─────────────────────────────────────────────────────────────────────────────

const excelSheetName = "Usuarios en BD"

// Índices de columna esperados en la hoja (base 0)
const (
	colLoginBD     = 2  // Login BD
	colNombreBD    = 3  // Nombre en BD
	colNombreExcel = 0  // Nombre (Excel) — fallback si Nombre BD está vacío
	colRolKreivo   = 6  // RolKreivo  → sv | cualquier otro
	colDocNumero   = 7  // Doc. Número
	colEquipo      = 8  // Equipo
	colEmailBD     = 9  // Email(s) en BD
	colCorreoExcel = 10 // Correo (Excel) — fallback si Email BD está vacío
	colDocTipo     = 13 // Doc. Tipo
)

// ─────────────────────────────────────────────────────────────────────────────
// MigrateExcelResult
// ─────────────────────────────────────────────────────────────────────────────

// MigrateExcelResult resume el resultado de la migración desde Excel.
type MigrateExcelResult struct {
	HojaCalculo       string               `json:"hoja_calculo"`
	TotalFilas        int                  `json:"total_filas"`
	FilasValidas      int                  `json:"filas_validas"`
	FilasOmitidas     int                  `json:"filas_omitidas"`
	CreadosCount      int                  `json:"creados_count"`
	ActualizadosCount int                  `json:"actualizados_count"`
	FallidosCount     int                  `json:"fallidos_count"`
	Creados           []string             `json:"creados"`
	Actualizados      []string             `json:"actualizados"`
	Fallidos          []MigrateSkippedItem `json:"fallidos"`
	ParsedButSkipped  []MigrateSkippedItem `json:"parsed_but_skipped"`
}

// ─────────────────────────────────────────────────────────────────────────────
// MigrateFromExcel
// ─────────────────────────────────────────────────────────────────────────────

// MigrateFromExcel recibe el contenido de un archivo .xlsx como []byte,
// lo parsea, aplica la lógica de negocio y crea o actualiza los usuarios.
//
// Clave de identificación: Doc. Número (cédula). El campo Login BD del Excel se usa
// solo para generar un login al crear un usuario nuevo; no es la clave de búsqueda.
//
// Lógica por fila:
//  1. Si Doc. Número está vacío → skip (no se puede identificar al usuario).
//  2. Si NO existe perfil con esa cédula en BD → CREAR usuario nuevo.
//  3. Si existe perfil:
//     a. Tiene general_user vinculado (credenciales) → actualizar rol + equipo.
//     b. Sin general_user → crear credenciales temporales (test.{login}) y luego actualizar rol + equipo.
//
// Regla de roles:
//   - RolKreivo == "sv"  → ["sv"]  (Supervisor)
//   - cualquier otro     → ["ro"]  (Operador)
func MigrateFromExcel(fileBytes []byte, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {

	// 1. Abrir Excel desde bytes en memoria
	f, err := excelize.OpenReader(bytes.NewReader(fileBytes))
	if err != nil {
		return http.StatusBadRequest, `{"error":"Archivo Excel inválido: ` + err.Error() + `"}`
	}
	defer f.Close()

	// 2. Leer filas de la hoja
	rows, err := f.GetRows(excelSheetName)
	if err != nil {
		sheets := f.GetSheetList()
		return http.StatusBadRequest, fmt.Sprintf(
			`{"error":"Hoja '%s' no encontrada. Hojas disponibles: %s"}`,
			excelSheetName, strings.Join(sheets, ", "),
		)
	}

	if len(rows) < 2 {
		return http.StatusBadRequest, `{"error":"El Excel no tiene datos (solo encabezado o vacío)"}`
	}

	// 3. Parsear filas
	var users []SetupCreateUserRequest
	var parsedSkipped []MigrateSkippedItem

	for i, row := range rows {
		if i == 0 {
			continue // saltar encabezado
		}

		login := strings.TrimSpace(cellAt(row, colLoginBD))

		// Nombre completo → names + lastNames
		fullName := strings.TrimSpace(cellAt(row, colNombreBD))
		if fullName == "" {
			fullName = strings.TrimSpace(cellAt(row, colNombreExcel))
		}
		names, lastNames := splitName(fullName, "")

		// Si el login está vacío o es "Crear", generarlo a partir del nombre
		if login == "" || strings.EqualFold(login, "crear") {
			if fullName == "" {
				parsedSkipped = append(parsedSkipped, MigrateSkippedItem{
					Login:  "(sin login ni nombre)",
					Reason: "fila sin login y sin nombre — no se puede generar login",
				})
				continue
			}
			login = generateLogin(fullName)
		}

		// Si el nombre quedó vacío, usar el login como fallback
		if names == "" {
			names, lastNames = login, login
		}

		// Email: BD primero, Excel como fallback
		email := strings.TrimSpace(cellAt(row, colEmailBD))
		if email == "" {
			email = strings.TrimSpace(cellAt(row, colCorreoExcel))
		}
		if email == "" {
			email = login + "@salvia-temp.com"
		}

		// Documento
		docNumber := strings.TrimSpace(cellAt(row, colDocNumero))
		if docNumber == "" {
			docNumber = "AUTO_" + login
		}
		docType := strings.ToUpper(strings.TrimSpace(cellAt(row, colDocTipo)))
		if docType == "" {
			docType = "CC"
		}

		// Rol: usar el valor del Excel tal cual; si está vacío, defaultea a "ro"
		rolKreivo := strings.TrimSpace(cellAt(row, colRolKreivo))
		var roleCodes []string
		if rolKreivo != "" {
			roleCodes = []string{rolKreivo}
		} else {
			roleCodes = []string{"ro"}
		}

		// Equipo
		teamRaw := strings.TrimSpace(cellAt(row, colEquipo))
		team := normalizeExcelTeam(teamRaw)

		users = append(users, SetupCreateUserRequest{
			Login:     login,
			Password:  "Salvia@Prod2026!", // contraseña por defecto al crear
			Names:     names,
			LastNames: lastNames,
			Email:     email,
			Phone:     "0000000000",
			DocType:   docType,
			DocNumber: docNumber,
			Gender:    "ma", // el Excel no tiene género — default neutro
			Language:  "sp",
			TownCode:  "11001000",
			Team:      team,
			RoleCodes: roleCodes,
		})
	}

	// 4. Ejecutar migración
	result := MigrateExcelResult{
		HojaCalculo:      excelSheetName,
		TotalFilas:       len(rows) - 1,
		FilasValidas:     len(users),
		FilasOmitidas:    len(parsedSkipped),
		Creados:          []string{},
		Actualizados:     []string{},
		Fallidos:         []MigrateSkippedItem{},
		ParsedButSkipped: parsedSkipped,
	}

	total := len(users)
	for idx, u := range users {
		fmt.Printf("[migrate-excel] %d/%d → cédula:%s login:%s ", idx+1, total, u.DocNumber, u.Login)

		// Sin cédula real no podemos identificar al usuario → skip.
		if u.DocNumber == "" || strings.HasPrefix(u.DocNumber, "AUTO_") {
			fmt.Println("⏭ sin cédula válida")
			result.Fallidos = append(result.Fallidos, MigrateSkippedItem{Login: u.Login, Reason: "sin Doc. Número — fila omitida"})
			result.FallidosCount++
			continue
		}

		// ── PASO 1: buscar perfil por cédula ────────────────────────────────
		var existingProfile security_daos.GeneralUserProfileDTO
		var connProfile db.ConnData
		profileFound := security_daos.GetGeneralUserProfile(
			common_controllers.By{
				Operator:   common_dao.SQL_AND,
				AttrsName:  []string{"GeneralUserProfileDocNumber"},
				AttrsValue: []interface{}{u.DocNumber},
			},
			&existingProfile, &connProfile, &dbClientConfig, &dbServerConfig,
		) == nil
		db.ReleaseConnection(&connProfile)

		if !profileFound {
			// ── CASO A: usuario nuevo ────────────────────────────────────────
			// Verificar que el login no esté en uso por otro usuario con distinta cédula
			var loginCheck security_daos.GeneralUserDTO
			var connLoginCheck db.ConnData
			loginExists := security_daos.GetGeneralUserByLogin(u.Login, &loginCheck, &connLoginCheck, &dbClientConfig, &dbServerConfig) == nil
			db.ReleaseConnection(&connLoginCheck)
			if loginExists {
				fmt.Printf("⏭ login '%s' ya está en uso por otro usuario\n", u.Login)
				result.Fallidos = append(result.Fallidos, MigrateSkippedItem{Login: u.Login, Reason: "login ya está en uso por otro usuario con distinta cédula"})
				result.FallidosCount++
				continue
			}

			createBody, _ := json.Marshal(u)
			var connCreate db.ConnData
			_, res := SetupCreateUser(string(createBody), &connCreate, dbClientConfig, dbServerConfig)

			if isErrorResponse(res) {
				fmt.Printf("❌ error creando: %s\n", res)
				result.Fallidos = append(result.Fallidos, MigrateSkippedItem{Login: u.Login, Reason: res})
				result.FallidosCount++
			} else {
				fmt.Println("✅ creado")
				result.Creados = append(result.Creados, u.Login)
				result.CreadosCount++
			}
			continue
		}

		// ── PASO 2: perfil existe — buscar general_user vinculado ───────────
		var existingUser security_daos.GeneralUserDTO
		var connUser db.ConnData
		userFound := security_daos.GetGeneralUserByProfileId(existingProfile.GeneralUserProfileId, &existingUser, &connUser, &dbClientConfig, &dbServerConfig) == nil
		db.ReleaseConnection(&connUser)

		if userFound && existingUser.GeneralUserLogin != "" {
			// ── CASO B: tiene credenciales → solo actualizar rol y equipo ────
			updateBody, _ := json.Marshal(SetupUpdateRoleTeamRequest{
				RoleCodes: u.RoleCodes,
				Team:      u.Team,
			})
			var connUpdate db.ConnData
			_, res := SetupUpdateUserRoleTeam(existingUser.GeneralUserLogin, string(updateBody), &connUpdate, dbClientConfig, dbServerConfig)

			if isErrorResponse(res) {
				fmt.Printf("error actualizando: %s\n", res)
				result.Fallidos = append(result.Fallidos, MigrateSkippedItem{Login: existingUser.GeneralUserLogin, Reason: res})
				result.FallidosCount++
			} else {
				fmt.Printf("actualizado (login:%s)\n", existingUser.GeneralUserLogin)
				result.Actualizados = append(result.Actualizados, existingUser.GeneralUserLogin)
				result.ActualizadosCount++
			}
			continue
		}

		// ── CASO C: tiene perfil pero sin credenciales → crear test.{login} ─
		// Generar login a partir del nombre en el perfil (fuente más confiable que el Excel).
		profileFullName := strings.TrimSpace(existingProfile.GeneralUserProfileNames + " " + existingProfile.GeneralUserProfileLastNames)
		baseLogin := generateLogin(profileFullName)
		if baseLogin == "" || baseLogin == "usuario" {
			baseLogin = u.Login // fallback al login generado del Excel
		}
		testLogin := "test." + baseLogin

		testEmail := testLogin + "@salvia-temp.com"

		testUserReq := SetupCreateUserRequest{
			Login:     baseLogin,
			Password:  u.Password,
			Team:      u.Team,
			RoleCodes: u.RoleCodes,
			Names:     existingProfile.GeneralUserProfileNames,
			LastNames: existingProfile.GeneralUserProfileLastNames,
			DocType:   existingProfile.GeneralUserProfileDocType,
			DocNumber: u.DocNumber,
			Gender:    existingProfile.GeneralUserProfileGender,
			Language:  "sp",
			TownCode:  "11001000",
			Phone:     "0000000000",
			Email:     testEmail,
		}
		if testUserReq.DocType == "" {
			testUserReq.DocType = "CC"
		}
		if testUserReq.Gender == "" {
			testUserReq.Gender = "ma"
		}
		if testUserReq.Names == "" {
			testUserReq.Names = u.Login
		}
		if testUserReq.LastNames == "" {
			testUserReq.LastNames = u.Login
		}

		testBatchBody, _ := json.Marshal(MigrateTestRequest{Users: []SetupCreateUserRequest{testUserReq}})
		var connTest db.ConnData
		_, testRes := MigrateTestUsers(string(testBatchBody), &connTest, dbClientConfig, dbServerConfig)
		db.ReleaseConnection(&connTest)

		if isErrorResponse(testRes) {
			fmt.Printf("error creando credenciales temporales: %s\n", testRes)
			result.Fallidos = append(result.Fallidos, MigrateSkippedItem{Login: testLogin, Reason: "error al crear credenciales temporales: " + testRes})
			result.FallidosCount++
			continue
		}

		// Actualizar rol y equipo del usuario de prueba recién creado.
		updateBody, _ := json.Marshal(SetupUpdateRoleTeamRequest{
			RoleCodes: u.RoleCodes,
			Team:      u.Team,
		})
		var connUpdate db.ConnData
		_, res := SetupUpdateUserRoleTeam(testLogin, string(updateBody), &connUpdate, dbClientConfig, dbServerConfig)

		if isErrorResponse(res) {
			fmt.Printf("credenciales creadas pero error actualizando rol: %s\n", res)
			result.Fallidos = append(result.Fallidos, MigrateSkippedItem{Login: testLogin, Reason: "credenciales creadas, error en rol/equipo: " + res})
			result.FallidosCount++
		} else {
			fmt.Printf("credenciales temporales creadas y actualizadas (login:%s)\n", testLogin)
			result.Actualizados = append(result.Actualizados, testLogin)
			result.ActualizadosCount++
		}
	}
	fmt.Printf("[migrate-excel] Terminado: %d creados, %d actualizados, %d fallidos\n",
		result.CreadosCount, result.ActualizadosCount, result.FallidosCount)

	out, _ := json.Marshal(result)
	return http.StatusOK, string(out)
}

// ─────────────────────────────────────────────────────────────────────────────
// Helpers de parseo
// ─────────────────────────────────────────────────────────────────────────────

// cellAt obtiene el valor de una celda por índice de columna (base 0).
// Retorna "" si la columna no existe en la fila.
func cellAt(row []string, col int) string {
	if col < len(row) {
		return row[col]
	}
	return ""
}

// splitName divide un nombre completo en (names, lastNames).
// Toma las primeras 2 palabras como nombres y el resto como apellidos.
func splitName(fullName, fallback string) (string, string) {
	if fullName == "" {
		return fallback, fallback
	}
	parts := strings.Fields(fullName)
	if len(parts) == 1 {
		return parts[0], parts[0]
	}
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	// 3+ palabras: primeras 2 = nombres, resto = apellidos
	return strings.Join(parts[:2], " "), strings.Join(parts[2:], " ")
}

// generateLogin genera un login a partir del nombre completo.
// Formato: "primer_nombre.primer_apellido" en minúsculas sin tildes.
// Ejemplo: "Isabel Agatón Santander" → "isabel.agaton"
//
//	"Ana María Mojica Quiroz" → "ana.mojica"
func generateLogin(fullName string) string {
	parts := strings.Fields(fullName)
	if len(parts) == 0 {
		return "usuario"
	}

	var first, last string

	if len(parts) == 1 {
		first = parts[0]
		last = parts[0]
	} else if len(parts) == 2 {
		// "Nombre Apellido"
		first = parts[0]
		last = parts[1]
	} else if len(parts) == 3 {
		// "Nombre Apellido1 Apellido2" → nombre.apellido1
		first = parts[0]
		last = parts[1]
	} else {
		// "Nombre1 Nombre2 Apellido1 Apellido2" → nombre1.apellido1
		first = parts[0]
		last = parts[2]
	}

	return removeAccents(strings.ToLower(first)) + "." + removeAccents(strings.ToLower(last))
}

// removeAccents elimina tildes y caracteres especiales de una cadena.
// "ó" → "o", "ñ" → "n", etc.
func removeAccents(s string) string {
	// Normalizar a NFD (descomponer caracteres acentuados)
	t := transform.Chain(norm.NFD, transform.RemoveFunc(func(r rune) bool {
		return unicode.Is(unicode.Mn, r) // eliminar marcas de combinación (tildes)
	}), norm.NFC)
	result, _, _ := transform.String(t, s)

	// Eliminar caracteres no alfanuméricos ni punto (por si acaso)
	re := regexp.MustCompile(`[^a-z0-9._-]`)
	return re.ReplaceAllString(result, "")
}

// normalizeExcelTeam convierte el valor de la columna "Equipo" del Excel
// al valor exacto que acepta el sistema.
func normalizeExcelTeam(raw string) string {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "riesgo alto":
		return "Riesgo alto"
	case "riesgo bajo":
		return "Riesgo bajo"
	case "hombres":
		return "Hombres"
	default:
		return ""
	}
}
