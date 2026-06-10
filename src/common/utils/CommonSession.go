package utils

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type CommonSession struct {
	UserICode         string
	UserLogin         string
	SessionID         string
	Roles             []string
	CurrentRole       string
	EntityBrandICode  string
	TownCode          string
	TownICode         string
	EntityBranchICode string
	Lang              string
	CurrentMenu       map[string][]map[string]string
	Names             string
	LastNames         string
	Emails            []string
	UpdateTime        time.Time
	Team              string // RIESGO_BAJO | RIESGO_ALTO | "" (para ad y otros)
}

var commonSessions map[string]*CommonSession = map[string]*CommonSession{}
var loaded bool = false

func init() {
	loadStoredSessions()
}

func loadStoredSessions() {
	if !loaded {
		if _, err := os.Stat("sessions"); !os.IsNotExist(err) {
			datos, err := os.ReadFile("sessions")
			if err != nil {
				println("Session loadStoredSessions Error: read file")
			} else if err = json.Unmarshal(datos, &commonSessions); err != nil {
				println("Session loadStoredSessions Error: Unmarshal file")
			} else {
				loaded = true
				cleanOutdatedSessions()
			}
		} else {
			//Creamos uno vacío
			storeSessions(true)
		}
	}
}

func cleanOutdatedSessions() {
	for _, s := range commonSessions {
		if time.Since(s.UpdateTime) > 10*time.Hour {
			//No se guarda para poder limpiar todo y guardar una vez al final
			RemoveCommonSession(s.SessionID, false)
		}
	}
	storeSessions(false)
}

func storeSessions(forced bool) {
	if loaded || forced {
		datos, err := json.Marshal(commonSessions)
		if err != nil {
			println("Session storeSessions Error: Marshal file")
		} else if err = os.WriteFile("sessions", datos, 0644); err != nil {
			println("Session storeSessions Error: WriteFile file")
		}
	}
}

func AddCommonSession(sessionId string, commonSession *CommonSession) {
	lock.Lock()
	defer lock.Unlock()
	commonSession.UpdateTime = time.Now()
	commonSessions[sessionId] = commonSession
	storeSessions(false)
}

func GetCommonSession(sessionId string) (*CommonSession, error) {
	lock.RLock()
	defer lock.RUnlock()
	session, found := commonSessions[sessionId]
	if found {
		session.UpdateTime = time.Now()
		return session, nil
	}
	return session, errors.New("Session_not_found")
}

func GetCommonSessionByLogin(login string) (*CommonSession, error) {
	lock.RLock()
	defer lock.RUnlock()
	for _, s := range commonSessions {
		if s.UserLogin == login {
			s.UpdateTime = time.Now()
			return s, nil
		}
	}

	return &CommonSession{}, errors.New("Session_not_found")
}

func RemoveCommonSession(sessionId string, storeSession bool) {
	lock.Lock()
	defer lock.Unlock()
	_, found := commonSessions[sessionId]

	if found {
		delete(commonSessions, sessionId)
	}
	if storeSession {
		storeSessions(false)
	}
}

func HasRole(sessionId string, role string) bool {
	lock.RLock()
	defer lock.RUnlock()
	session, found := commonSessions[sessionId]

	if found {
		for _, v := range session.Roles {
			if v == role {
				return true
			}
		}
	}
	return false
}

func CheckPermission(permissionsMap map[string]map[string]bool, service string, role string, c *gin.Context) bool {
	lock.RLock()
	defer lock.RUnlock()
	if _, found := permissionsMap[service][role]; !found {
		c.DataFromReader(401, int64(len("")), gin.MIMEJSON, strings.NewReader(""), nil)
		return false
	}
	return true
}

// CheckPermissionWithTeam valida el permiso por rol Y verifica que el equipo del usuario
// coincida con el requerido. Si requiredTeam es "" se omite la validación de equipo.
func CheckPermissionWithTeam(permissionsMap map[string]map[string]bool, service string, role string, userTeam string, requiredTeam string, c *gin.Context) bool {
	if role == "ad" {
		return true
	}
	if !CheckPermission(permissionsMap, service, role, c) {
		return false
	}
	if requiredTeam == "" {
		return true
	}
	if userTeam != requiredTeam {
		c.DataFromReader(403, int64(len(`{"error":"acceso denegado para su equipo"}`)), gin.MIMEJSON, strings.NewReader(`{"error":"acceso denegado para su equipo"}`), nil)
		return false
	}
	return true
}
