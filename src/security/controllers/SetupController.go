// Package security_ctrl — controladores para los endpoints de setup/administración.
// Estos endpoints están protegidos por API Key y no requieren sesión de usuario.
package security_ctrl

import (
	common_controllers "bitsflow/common/controllers"
	common_dao "bitsflow/common/dao"
	"bitsflow/common/db"
	"bitsflow/common/utils"
	salvia_daos "bitsflow/salvia/dao"
	security_daos "bitsflow/security/dao"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// ─────────────────────────────────────────────────────────────────────────────
// Teams válidos del sistema
// ─────────────────────────────────────────────────────────────────────────────

// ValidTeams define los equipos reconocidos por el sistema.
// El team del usuario debe coincidir con el team de los seguimientos (follow_up_v2.team)
// para que el enrutamiento de casos funcione correctamente.
//
//	"Riesgo alto"  → agentes que atienden casos con nivel de riesgo 3 o 4
//	"Riesgo bajo"  → agentes que atienden casos con nivel de riesgo 1 o 2
//	"SIN_EQUIPO"   → sin equipo asignado (valor por defecto del sistema)
var ValidTeams = map[string]bool{
	"Riesgo alto": true,
	"Riesgo bajo": true,
	"SIN_EQUIPO":  true,
	"Hombres":     true, // Equipo del programa de atención a hombres
}

// isValidTeam retorna true si el team enviado es uno de los valores reconocidos.
// Permite cadena vacía (el usuario no tiene equipo aún).
func isValidTeam(team string) bool {
	if team == "" {
		return true
	}
	return ValidTeams[team]
}

// ─────────────────────────────────────────────────────────────────────────────
// DTOs de request
// ─────────────────────────────────────────────────────────────────────────────

// SetupCreateUserRequest contiene los datos necesarios para crear un nuevo usuario
// a través del endpoint de setup.
type SetupCreateUserRequest struct {
	Login             string   `json:"login"`
	Password          string   `json:"pass"`
	Team              string   `json:"team"`
	RoleCodes         []string `json:"roleCodes"`
	Names             string   `json:"names"`
	LastNames         string   `json:"lastNames"`
	Email             string   `json:"email"`
	Phone             string   `json:"phone"`
	DocType           string   `json:"docType"`
	DocNumber         string   `json:"docNumber"`
	Gender            string   `json:"gender"`
	Language          string   `json:"lang"`
	EntityBranchICode string   `json:"entityBranchICode"`
	TownCode          string   `json:"townCode"`
}

// SetupUpdateRoleTeamRequest contiene los datos para cambiar rol y equipo de un usuario.
type SetupUpdateRoleTeamRequest struct {
	RoleCodes []string `json:"roleCodes"`
	Team      string   `json:"team"`
}

// ─────────────────────────────────────────────────────────────────────────────
// SetupCreateUser
// ─────────────────────────────────────────────────────────────────────────────

// SetupCreateUser crea un nuevo usuario en el sistema a partir de los datos recibidos.
// No requiere sesión activa — está protegido por API Key en la capa de facade.
// Retorna código HTTP y JSON de respuesta.
func SetupCreateUser(dataInput string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	// 1. Parsear request
	var req SetupCreateUserRequest
	if err := json.Unmarshal([]byte(dataInput), &req); err != nil {
		return http.StatusBadRequest, `{"error":"JSON inválido: ` + err.Error() + `"}`
	}

	// 2. Validar campos obligatorios
	if req.Login == "" || req.Password == "" || req.Names == "" || req.LastNames == "" ||
		req.Email == "" || req.DocType == "" || req.DocNumber == "" || req.Gender == "" || len(req.RoleCodes) == 0 {
		return http.StatusBadRequest, `{"error":"Campos obligatorios faltantes: login, pass, names, lastNames, email, docType, docNumber, gender, roleCodes"}`
	}

	// 3. Verificar fortaleza de contraseña
	if passErrs := utils.CheckPasswordStrength(req.Password); len(passErrs) > 0 {
		return http.StatusBadRequest, `{"error":"Contraseña inválida: ` + passErrs[0] + `"}`
	}

	// 3b. Validar team (debe ser uno de los valores del sistema o vacío)
	if !isValidTeam(req.Team) {
		return http.StatusBadRequest, `{"error":"Team inválido: '` + req.Team + `'. Valores permitidos: \"Riesgo alto\", \"Riesgo bajo\", \"SIN_EQUIPO\""}`
	}

	// 4. Valores por defecto opcionales
	if req.Language == "" {
		req.Language = "sp"
	}
	if req.TownCode == "" {
		req.TownCode = "11001000" // Bogotá por defecto
	}
	if req.Phone == "" {
		req.Phone = "0000000000"
	}

	// 5. Hashear contraseña con bcrypt (cost=10, igual que el flujo estándar)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return http.StatusInternalServerError, `{"error":"Error al hashear la contraseña"}`
	}

	// 6. Resolver roles por código → obtener DTOs con IDs reales
	var roles []security_daos.RoleDTO
	for _, code := range req.RoleCodes {
		var role security_daos.RoleDTO
		by := common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"RoleCode"},
			AttrsValue: []interface{}{code},
		}
		if err = security_daos.GetRole(by, &role, connData, &dbClientConfig, &dbServerConfig); err != nil {
			return http.StatusBadRequest, `{"error":"Rol no encontrado: ` + code + `"}`
		}
		roles = append(roles, role)
	}

	// 7. Construir DTOs de usuario y perfil con valores por defecto
	var userDTO security_daos.GeneralUserDTO
	security_daos.SetGeneralUserDefaults(&userDTO, common_dao.SQL_INSERT)
	userDTO.GeneralUserLogin = req.Login
	userDTO.GeneralUserPassword = string(hashedPassword)
	userDTO.GeneralUserLanguage = req.Language
	userDTO.GeneralUserTeam = req.Team

	var profileDTO security_daos.GeneralUserProfileDTO
	security_daos.SetGeneralUserProfileDefaults(&profileDTO, common_dao.SQL_INSERT)
	profileDTO.GeneralUserProfileNames = req.Names
	profileDTO.GeneralUserProfileLastNames = req.LastNames
	profileDTO.GeneralUserProfileDocType = req.DocType
	profileDTO.GeneralUserProfileDocNumber = req.DocNumber
	profileDTO.GeneralUserProfileGender = req.Gender
	profileDTO.GeneralUserProfileTownCode = req.TownCode

	var mailDTO security_daos.EMailDTO
	security_daos.SetEMailDefaults(&mailDTO, common_dao.SQL_INSERT)
	mailDTO.EMailData = req.Email

	var phoneDTO security_daos.PhoneNumberDTO
	security_daos.SetPhoneNumberDefaults(&phoneDTO, common_dao.SQL_INSERT)
	phoneDTO.PhoneNumberData = req.Phone

	// 8. Verificar unicidad de login
	{
		var existing security_daos.GeneralUserDTO
		by := common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"GeneralUserLogin"},
			AttrsValue: []interface{}{req.Login},
		}
		if security_daos.GetGeneralUser(by, &existing, connData, &dbClientConfig, &dbServerConfig) == nil {
			return http.StatusConflict, `{"error":"El login ya está en uso"}`
		}
	}

	// 9. Verificar unicidad de email
	{
		var existing security_daos.EMailDTO
		by := common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"EMailData"},
			AttrsValue: []interface{}{req.Email},
		}
		if security_daos.GetEMail(by, &existing, connData, &dbClientConfig, &dbServerConfig) == nil {
			return http.StatusConflict, `{"error":"El email ya está registrado"}`
		}
	}

	// 10. Verificar unicidad de documento
	{
		var existing security_daos.GeneralUserProfileDTO
		by := common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"GeneralUserProfileDocType", "GeneralUserProfileDocNumber"},
			AttrsValue: []interface{}{req.DocType, req.DocNumber},
		}
		if security_daos.GetGeneralUserProfile(by, &existing, connData, &dbClientConfig, &dbServerConfig) == nil {
			return http.StatusConflict, `{"error":"El documento ya está registrado"}`
		}
	}

	// 11. Iniciar transacción
	if _, err = db.StartTransaction(connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, `{"error":"` + err.Error() + `"}`
	}

	// 12. Crear perfil
	if err = security_daos.SetGeneralUserProfile(&profileDTO, connData, &dbClientConfig, &dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, `{"error":"Error creando perfil: ` + err.Error() + `"}`
	}
	userDTO.GeneralUserGeneralUserProfile = profileDTO

	// 13. Crear email
	mailDTO.EMailGeneralUserProfile = profileDTO
	if err = security_daos.SetEMail(&mailDTO, connData, &dbClientConfig, &dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, `{"error":"Error creando email: ` + err.Error() + `"}`
	}

	// 14. Crear teléfono
	phoneDTO.PhoneNumberGeneralUserProfile = profileDTO
	if err = security_daos.SetPhoneNumber(&phoneDTO, connData, &dbClientConfig, &dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, `{"error":"Error creando teléfono: ` + err.Error() + `"}`
	}

	// 15. Crear CaseOwner (siempre se crea; se asocia a sede si se envió entityBranchICode)
	var cowner salvia_daos.CaseOwnerDTO
	salvia_daos.SetCaseOwnerDefaults(&cowner, common_dao.SQL_INSERT)
	if req.EntityBranchICode != "" {
		var branch salvia_daos.EntityBranchDTO
		by := common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"EntityBranchICode"},
			AttrsValue: []interface{}{req.EntityBranchICode},
		}
		if err = salvia_daos.GetEntityBranch(by, &branch, connData, &dbClientConfig, &dbServerConfig); err != nil {
			db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
			return http.StatusBadRequest, `{"error":"Sede no encontrada: ` + req.EntityBranchICode + `"}`
		}
		cowner.EntityBranch = branch
	}
	cowner.CaseOwnerGeneralUser = userDTO.GeneralUserICode
	if err = salvia_daos.SetCaseOwner(&cowner, connData, &dbClientConfig, &dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, `{"error":"Error creando case owner: ` + err.Error() + `"}`
	}

	// 16. Crear usuario
	if err = security_daos.SetGeneralUser(&userDTO, connData, &dbClientConfig, &dbServerConfig); err != nil {
		db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
		return http.StatusInternalServerError, `{"error":"Error creando usuario: ` + err.Error() + `"}`
	}

	// 17. Asignar roles
	for _, r := range roles {
		var rel security_daos.RelRoleGeneralUserDTO
		security_daos.SetRelRoleGeneralUserDefaults(&rel, common_dao.SQL_INSERT)
		rel.RelRoleGeneralUserGeneralUser = userDTO
		rel.RelRoleGeneralUserRole = r
		if err = security_daos.SetRelRoleGeneralUser(&rel, connData, &dbClientConfig, &dbServerConfig); err != nil {
			db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
			return http.StatusInternalServerError, `{"error":"Error asignando rol: ` + err.Error() + `"}`
		}
	}

	// 18. Confirmar transacción
	if _, err = db.CommitTransaction(connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, `{"error":"` + err.Error() + `"}`
	}

	// 19. Actualizar team (operación separada — el DAO de SetGeneralUser no incluye el campo team)
	if req.Team != "" {
		var connTeam db.ConnData
		defer db.ReleaseConnection(&connTeam)
		if teamErr := security_daos.UpdateGeneralUserTeamByICode(userDTO.GeneralUserICode, req.Team, &connTeam, &dbClientConfig, &dbServerConfig); teamErr != nil {
			// No es fatal, el usuario fue creado correctamente; solo logeamos el warning
			// El team queda vacío y puede actualizarse con el segundo endpoint
			_ = teamErr
		}
	}

	return http.StatusCreated, `{"success":true,"icode":"` + userDTO.GeneralUserICode + `","login":"` + userDTO.GeneralUserLogin + `"}`
}

// ─────────────────────────────────────────────────────────────────────────────
// SetupUpdateUserRoleTeam
// ─────────────────────────────────────────────────────────────────────────────

// SetupUpdateUserRoleTeam cambia el rol y el equipo de un usuario identificado por su login.
// Reemplaza todos los roles actuales por los nuevos enviados en la petición.
// No requiere sesión activa — está protegido por API Key en la capa de facade.
// Retorna código HTTP y JSON de respuesta.
func SetupUpdateUserRoleTeam(login string, dataInput string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	if login == "" {
		return http.StatusBadRequest, `{"error":"El parámetro login es obligatorio"}`
	}

	// 1. Parsear request
	var req SetupUpdateRoleTeamRequest
	if err := json.Unmarshal([]byte(dataInput), &req); err != nil {
		return http.StatusBadRequest, `{"error":"JSON inválido: ` + err.Error() + `"}`
	}

	if len(req.RoleCodes) == 0 && req.Team == "" {
		return http.StatusBadRequest, `{"error":"Debe enviar al menos roleCodes o team para actualizar"}`
	}

	// Validar team (debe ser uno de los valores del sistema o vacío)
	if !isValidTeam(req.Team) {
		return http.StatusBadRequest, `{"error":"Team inválido: '` + req.Team + `'. Valores permitidos: \"Riesgo alto\", \"Riesgo bajo\", \"SIN_EQUIPO\""}`
	}

	// 2. Obtener usuario por login
	var user security_daos.GeneralUserDTO
	if err := security_daos.GetGeneralUserByLogin(login, &user, connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusNotFound, `{"error":"Usuario no encontrado: ` + login + `"}`
	}

	// 3. Si vienen roles, resolver los nuevos roles por código
	var newRoles []security_daos.RoleDTO
	for _, code := range req.RoleCodes {
		var role security_daos.RoleDTO
		by := common_controllers.By{
			Operator:   common_dao.SQL_AND,
			AttrsName:  []string{"RoleCode"},
			AttrsValue: []interface{}{code},
		}
		if err := security_daos.GetRole(by, &role, connData, &dbClientConfig, &dbServerConfig); err != nil {
			return http.StatusBadRequest, `{"error":"Rol no encontrado: ` + code + `"}`
		}
		newRoles = append(newRoles, role)
	}

	// 4. Obtener roles actuales del usuario (para eliminarlos)
	currentRoles, err := security_daos.GetRolesByGeneralUser(&user, connData, &dbClientConfig, &dbServerConfig)
	if err != nil {
		return http.StatusInternalServerError, `{"error":"Error obteniendo roles actuales: ` + err.Error() + `"}`
	}

	// 5. Iniciar transacción
	if _, err = db.StartTransaction(connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, `{"error":"` + err.Error() + `"}`
	}

	// 6. Eliminar todos los roles actuales (solo si se envían nuevos roles)
	if len(req.RoleCodes) > 0 {
		for _, r := range currentRoles {
			by := common_controllers.By{
				Operator:   common_dao.SQL_AND,
				AttrsName:  []string{"RelRoleGeneralUserRole", "RelRoleGeneralUserGeneralUser"},
				AttrsValue: []interface{}{r.RoleId, user.GeneralUserId},
			}
			if err = security_daos.RemoveRelRoleGeneralUsers(by, connData, &dbClientConfig, &dbServerConfig); err != nil {
				db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
				return http.StatusInternalServerError, `{"error":"Error eliminando roles actuales: ` + err.Error() + `"}`
			}
		}

		// 7. Insertar nuevos roles
		for _, r := range newRoles {
			var rel security_daos.RelRoleGeneralUserDTO
			security_daos.SetRelRoleGeneralUserDefaults(&rel, common_dao.SQL_INSERT)
			rel.RelRoleGeneralUserGeneralUser = user
			rel.RelRoleGeneralUserRole = r
			if err = security_daos.SetRelRoleGeneralUser(&rel, connData, &dbClientConfig, &dbServerConfig); err != nil {
				db.RollbackTransaction(connData, &dbClientConfig, &dbServerConfig)
				return http.StatusInternalServerError, `{"error":"Error asignando nuevo rol: ` + err.Error() + `"}`
			}
		}
	}

	// 8. Confirmar transacción
	if _, err = db.CommitTransaction(connData, &dbClientConfig, &dbServerConfig); err != nil {
		return http.StatusInternalServerError, `{"error":"` + err.Error() + `"}`
	}

	// 9. Actualizar team (fuera de la transacción — operación independiente)
	if req.Team != "" {
		var connTeam db.ConnData
		defer db.ReleaseConnection(&connTeam)
		if teamErr := security_daos.UpdateGeneralUserTeamByICode(user.GeneralUserICode, req.Team, &connTeam, &dbClientConfig, &dbServerConfig); teamErr != nil {
			return http.StatusInternalServerError, `{"error":"Roles actualizados pero falló el team: ` + teamErr.Error() + `"}`
		}
	}

	return http.StatusOK, `{"success":true,"login":"` + login + `","icode":"` + user.GeneralUserICode + `"}`
}

// ─────────────────────────────────────────────────────────────────────────────
// MigrateBatchUsers  —  POST /api/v1/admin/migrate/users
// ─────────────────────────────────────────────────────────────────────────────

// MigrateBatchRequest es el body del endpoint de migración masiva.
type MigrateBatchRequest struct {
	Users []SetupCreateUserRequest `json:"users"`
}

// MigrateBatchResult resume el resultado de la operación masiva.
type MigrateBatchResult struct {
	TotalProcessed int                  `json:"total_processed"`
	CreatedCount   int                  `json:"created_count"`
	UpdatedCount   int                  `json:"updated_count"`
	SkippedCount   int                  `json:"skipped_count"`
	Created        []string             `json:"created"`
	Updated        []string             `json:"updated"`
	Skipped        []MigrateSkippedItem `json:"skipped"`
}

// MigrateSkippedItem representa un usuario que no pudo procesarse.
type MigrateSkippedItem struct {
	Login  string `json:"login"`
	Reason string `json:"reason"`
}

// MigrateBatchUsers crea o actualiza masivamente usuarios a partir del array enviado.
//
// Regla de roles: si roleCodes contiene "sv" → se asigna ["sv"], en cualquier otro caso → ["ro"].
//
// Para cada usuario:
//   - Si ya existe por login → actualiza rol y team (PATCH).
//   - Si no existe → lo crea con todos los datos (POST).
func MigrateBatchUsers(dataInput string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	var req MigrateBatchRequest
	if err := json.Unmarshal([]byte(dataInput), &req); err != nil {
		return http.StatusBadRequest, `{"error":"JSON inválido: ` + err.Error() + `"}`
	}
	if len(req.Users) == 0 {
		return http.StatusBadRequest, `{"error":"El array 'users' está vacío"}`
	}

	result := MigrateBatchResult{
		Created: []string{},
		Updated: []string{},
		Skipped: []MigrateSkippedItem{},
	}

	total := len(req.Users)
	for idx, u := range req.Users {
		fmt.Printf("[migrate-users] %d/%d → cédula:%s login:%s ", idx+1, total, u.DocNumber, u.Login)

		// Rol: usar el valor del request tal cual; si está vacío, defaultea a "ro"
		if len(u.RoleCodes) == 0 || (len(u.RoleCodes) == 1 && u.RoleCodes[0] == "") {
			u.RoleCodes = []string{"ro"}
		}

		if !isValidTeam(u.Team) {
			fmt.Printf("⏭ team inválido: %s\n", u.Team)
			result.Skipped = append(result.Skipped, MigrateSkippedItem{Login: u.Login, Reason: "team inválido: " + u.Team})
			result.SkippedCount++
			continue
		}

		// Sin cédula no podemos identificar al usuario → skip.
		if u.DocNumber == "" {
			fmt.Println("⏭ sin cédula válida")
			result.Skipped = append(result.Skipped, MigrateSkippedItem{Login: u.Login, Reason: "sin docNumber — omitido"})
			result.SkippedCount++
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
			// Si no viene login, generarlo del nombre. Si tampoco hay nombre, skip.
			if u.Login == "" {
				fullName := strings.TrimSpace(u.Names + " " + u.LastNames)
				if fullName == "" {
					fmt.Println("⏭ usuario nuevo sin login ni nombre")
					result.Skipped = append(result.Skipped, MigrateSkippedItem{Login: "(sin login)", Reason: "usuario no existe en BD y no hay login ni nombre para generarlo"})
					result.SkippedCount++
					continue
				}
				u.Login = generateLogin(fullName)
			}

			// Verificar que el login no esté en uso por otro usuario
			var loginCheck security_daos.GeneralUserDTO
			var connLoginCheck db.ConnData
			loginExists := security_daos.GetGeneralUserByLogin(u.Login, &loginCheck, &connLoginCheck, &dbClientConfig, &dbServerConfig) == nil
			db.ReleaseConnection(&connLoginCheck)
			if loginExists {
				fmt.Printf("⏭ login '%s' ya está en uso por otro usuario\n", u.Login)
				result.Skipped = append(result.Skipped, MigrateSkippedItem{Login: u.Login, Reason: "login ya está en uso por otro usuario con distinta cédula"})
				result.SkippedCount++
				continue
			}

			createBody, _ := json.Marshal(u)
			var connCreate db.ConnData
			_, res := SetupCreateUser(string(createBody), &connCreate, dbClientConfig, dbServerConfig)

			if isErrorResponse(res) {
				fmt.Printf("❌ %s\n", res)
				result.Skipped = append(result.Skipped, MigrateSkippedItem{Login: u.Login, Reason: "error al crear: " + res})
				result.SkippedCount++
			} else {
				fmt.Println("✅ creado")
				result.Created = append(result.Created, u.Login)
				result.CreatedCount++
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
				fmt.Printf("❌ %s\n", res)
				result.Skipped = append(result.Skipped, MigrateSkippedItem{Login: existingUser.GeneralUserLogin, Reason: "error al actualizar: " + res})
				result.SkippedCount++
			} else {
				fmt.Printf("✅ actualizado (login:%s)\n", existingUser.GeneralUserLogin)
				result.Updated = append(result.Updated, existingUser.GeneralUserLogin)
				result.UpdatedCount++
			}
			continue
		}

		// ── CASO C: tiene perfil pero sin credenciales → crear test.{login} ─
		profileFullName := fmt.Sprintf("%s %s", existingProfile.GeneralUserProfileNames, existingProfile.GeneralUserProfileLastNames)
		testUserReq := SetupCreateUserRequest{
			Login:     u.Login,
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
			Email:     u.Login + "@salvia-temp.com",
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
		_ = profileFullName

		testBatchBody, _ := json.Marshal(MigrateTestRequest{Users: []SetupCreateUserRequest{testUserReq}})
		var connTest db.ConnData
		_, testRes := MigrateTestUsers(string(testBatchBody), &connTest, dbClientConfig, dbServerConfig)
		db.ReleaseConnection(&connTest)

		if isErrorResponse(testRes) {
			fmt.Printf("❌ error creando credenciales temporales: %s\n", testRes)
			result.Skipped = append(result.Skipped, MigrateSkippedItem{Login: "test." + u.Login, Reason: "error al crear credenciales temporales: " + testRes})
			result.SkippedCount++
			continue
		}

		updateBody, _ := json.Marshal(SetupUpdateRoleTeamRequest{
			RoleCodes: u.RoleCodes,
			Team:      u.Team,
		})
		var connUpdate db.ConnData
		_, res := SetupUpdateUserRoleTeam("test."+u.Login, string(updateBody), &connUpdate, dbClientConfig, dbServerConfig)

		if isErrorResponse(res) {
			fmt.Printf("❌ credenciales creadas pero error actualizando rol: %s\n", res)
			result.Skipped = append(result.Skipped, MigrateSkippedItem{Login: "test." + u.Login, Reason: "credenciales creadas, error en rol/equipo: " + res})
			result.SkippedCount++
		} else {
			fmt.Printf("✅ credenciales temporales creadas y actualizadas (login:test.%s)\n", u.Login)
			result.Updated = append(result.Updated, "test."+u.Login)
			result.UpdatedCount++
		}
	}

	result.TotalProcessed = result.CreatedCount + result.UpdatedCount + result.SkippedCount
	out, _ := json.Marshal(result)
	return http.StatusOK, string(out)
}

// ─────────────────────────────────────────────────────────────────────────────
// MigrateTestUsers  —  POST /api/v1/admin/migrate/test-users
// ─────────────────────────────────────────────────────────────────────────────

// MigrateTestRequest es el body del endpoint de credenciales de prueba.
type MigrateTestRequest struct {
	Users        []SetupCreateUserRequest `json:"users"`
	TestPassword string                   `json:"testPassword"` // opcional — default: Salvia@Test2026!
}

const defaultTestPassword = "Salvia@Test2026!"
const testLoginPrefix = "test."

// MigrateTestUsers crea credenciales temporales de prueba para verificación en producción.
//
// Para cada usuario en el array:
//   - login de prueba: "test." + login original (ej: test.ana.garcia)
//   - contraseña: testPassword o "Salvia@Test2026!" por defecto
//   - rol y team: iguales al usuario original (aplicando regla sv/ro)
//   - docNumber de prueba: "TEST_" + docNumber original
//   - email de prueba: "test." + email original (o generado)
//
// Si el usuario de prueba ya existe → se omite (no se sobreescribe).
func MigrateTestUsers(dataInput string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {
	var req MigrateTestRequest
	if err := json.Unmarshal([]byte(dataInput), &req); err != nil {
		return http.StatusBadRequest, `{"error":"JSON inválido: ` + err.Error() + `"}`
	}
	if len(req.Users) == 0 {
		return http.StatusBadRequest, `{"error":"El array 'users' está vacío"}`
	}

	testPassword := req.TestPassword
	if testPassword == "" {
		testPassword = defaultTestPassword
	}

	result := MigrateBatchResult{
		Created: []string{},
		Updated: []string{},
		Skipped: []MigrateSkippedItem{},
	}

	for _, u := range req.Users {
		if u.Login == "" {
			result.Skipped = append(result.Skipped, MigrateSkippedItem{Login: "(sin login)", Reason: "login vacío — omitido"})
			result.SkippedCount++
			continue
		}

		testLogin := testLoginPrefix + u.Login

		// Si ya existe el usuario de prueba → omitir
		var existing security_daos.GeneralUserDTO
		var connCheck db.ConnData
		alreadyExists := security_daos.GetGeneralUserByLogin(testLogin, &existing, &connCheck, &dbClientConfig, &dbServerConfig) == nil
		db.ReleaseConnection(&connCheck)

		if alreadyExists {
			result.Skipped = append(result.Skipped, MigrateSkippedItem{Login: testLogin, Reason: "ya existe — omitido"})
			result.SkippedCount++
			continue
		}

		// Aplicar regla de roles
		u.RoleCodes = resolveRoleCodes(u.RoleCodes)

		// Construir usuario de prueba
		testUser := SetupCreateUserRequest{
			Login:     testLogin,
			Password:  testPassword,
			Team:      u.Team,
			RoleCodes: u.RoleCodes,
			Names:     "TEST " + u.Names,
			LastNames: u.LastNames,
			DocType:   "CC",
			DocNumber: "TEST_" + u.DocNumber,
			Gender:    u.Gender,
			Language:  "sp",
			TownCode:  "11001000",
			Phone:     "0000000000",
		}

		// Email: prefija con "test." si se provee, si no genera uno
		if u.Email != "" {
			testUser.Email = "test." + u.Email
		} else {
			testUser.Email = testLogin + "@salvia-test.com"
		}

		// Team: si no es válido, usar SIN_EQUIPO
		if !isValidTeam(testUser.Team) {
			testUser.Team = "SIN_EQUIPO"
		}

		if testUser.Gender == "" {
			testUser.Gender = "ma"
		}

		createBody, _ := json.Marshal(testUser)
		var connCreate db.ConnData
		_, res := SetupCreateUser(string(createBody), &connCreate, dbClientConfig, dbServerConfig)
		db.ReleaseConnection(&connCreate)

		if isErrorResponse(res) {
			result.Skipped = append(result.Skipped, MigrateSkippedItem{Login: testLogin, Reason: "error al crear: " + res})
			result.SkippedCount++
		} else {
			result.Created = append(result.Created, testLogin)
			result.CreatedCount++
		}
	}

	result.TotalProcessed = result.CreatedCount + result.UpdatedCount + result.SkippedCount
	out, _ := json.Marshal(result)
	return http.StatusOK, string(out)
}

// ─────────────────────────────────────────────────────────────────────────────
// Helpers internos
// ─────────────────────────────────────────────────────────────────────────────

// resolveRoleCodes aplica la regla de negocio:
// si el usuario tiene rol "sv" → conserva ["sv"]
// en cualquier otro caso → fuerza ["ro"]
func resolveRoleCodes(codes []string) []string {
	for _, c := range codes {
		if c == "sv" {
			return []string{"sv"}
		}
	}
	return []string{"ro"}
}

// isErrorResponse detecta si la respuesta JSON contiene un campo "error".
func isErrorResponse(res string) bool {
	return len(res) > 0 && res[0] == '{' && len(res) > 8 && res[1:8] == `"error"`
}
