package controller

import (
	"bitsflow/common/utils"
	"bitsflow/salvia/service"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// fakeDirectorySvc cuenta si el handler llegó a ejecutarse. Si un guard falla, el
// contador delata que la escritura pasó.
type fakeDirectorySvc struct {
	service.DirectoryService
	calls int
}

func (f *fakeDirectorySvc) Create(context.Context, service.CreateDirectoryInput) (*service.DirectoryResponse, error) {
	f.calls++
	return &service.DirectoryResponse{}, nil
}

func (f *fakeDirectorySvc) CreateBulk(context.Context, []service.CreateDirectoryInput) (*service.DirectoryBulkCreateResult, error) {
	f.calls++
	return &service.DirectoryBulkCreateResult{Created: 1}, nil
}

func (f *fakeDirectorySvc) UpsertFromCSV(context.Context, io.Reader) (*service.DirectoryUpsertResult, error) {
	f.calls++
	return &service.DirectoryUpsertResult{}, nil
}

func (f *fakeDirectorySvc) Update(context.Context, string, service.UpdateDirectoryInput) (*service.DirectoryResponse, error) {
	f.calls++
	return &service.DirectoryResponse{}, nil
}

func (f *fakeDirectorySvc) Disable(context.Context, string) error { f.calls++; return nil }
func (f *fakeDirectorySvc) Enable(context.Context, string) error  { f.calls++; return nil }

func (f *fakeDirectorySvc) List(context.Context, service.DirectoryListFilters, int, int) (*service.DirectoryListResponse, error) {
	f.calls++
	return &service.DirectoryListResponse{}, nil
}

// newDirectoryRouter arma el router con la sesión real de gin. sessionID vacío
// simula una petición sin autenticar.
func newDirectoryRouter(sessionID string) (*gin.Engine, *fakeDirectorySvc) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(sessions.Sessions("default_session", cookie.NewStore([]byte("test"))))
	router.Use(func(ctx *gin.Context) {
		if sessionID != "" {
			sessions.Default(ctx).Set("userData", sessionID)
		}
		ctx.Next()
	})
	svc := &fakeDirectorySvc{}
	NewDirectoryController(svc).RegisterRoutes(router.Group("/api/v1"))
	return router, svc
}

// escrituras del Directorio con el rol que Menu.go habilita para cada una.
var directoryWrites = []struct {
	method, path, body, contentType string
	allowedRole                     string
}{
	{http.MethodPost, "/api/v1/directories", `{}`, "application/json", "et"},
	{http.MethodPost, "/api/v1/directories/bulk", `[]`, "application/json", "et"},
	{http.MethodPut, "/api/v1/directories/abc", `{}`, "application/json", "sv"},
	{http.MethodPatch, "/api/v1/directories/abc/disable", "", "", "sv"},
	{http.MethodPatch, "/api/v1/directories/abc/enable", "", "", "sv"},
}

func doDirectoryRequest(router *gin.Engine, w struct {
	method, path, body, contentType string
	allowedRole                     string
}) int {
	request := httptest.NewRequest(w.method, w.path, strings.NewReader(w.body))
	if w.contentType != "" {
		request.Header.Set("Content-Type", w.contentType)
	}
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	return recorder.Code
}

// Estas seis rutas estaban abiertas: cualquiera sin sesión podía crear, editar o
// desactivar entidades del Directorio.
func TestDirectoryWritesRejectAnonymous(t *testing.T) {
	for _, w := range directoryWrites {
		t.Run(w.method+" "+w.path, func(t *testing.T) {
			router, svc := newDirectoryRouter("")
			if code := doDirectoryRequest(router, w); code != http.StatusUnauthorized {
				t.Errorf("→ %d, esperaba 401", code)
			}
			if svc.calls != 0 {
				t.Errorf("el servicio se ejecutó %d veces sin sesión", svc.calls)
			}
		})
	}
}

// El rol de psicología no administra el Directorio: ni crea, ni edita, ni desactiva.
func TestDirectoryWritesRejectWrongRole(t *testing.T) {
	utils.AddCommonSession("sid-ps", &utils.CommonSession{UserICode: "u-ps", CurrentRole: "ps"})
	for _, w := range directoryWrites {
		t.Run(w.method+" "+w.path, func(t *testing.T) {
			router, svc := newDirectoryRouter("sid-ps")
			if code := doDirectoryRequest(router, w); code != http.StatusForbidden {
				t.Errorf("→ %d, esperaba 403", code)
			}
			if svc.calls != 0 {
				t.Errorf("el servicio se ejecutó %d veces con rol ps", svc.calls)
			}
		})
	}
}

// Con el rol habilitado la petición llega al handler. Importa porque los permisos
// son tres distintos: entidad crea pero no actualiza.
func TestDirectoryWritesAllowPermittedRole(t *testing.T) {
	utils.AddCommonSession("sid-sv", &utils.CommonSession{UserICode: "u-sv", CurrentRole: "sv"})
	utils.AddCommonSession("sid-et", &utils.CommonSession{UserICode: "u-et", CurrentRole: "et"})

	for _, w := range directoryWrites {
		t.Run(w.method+" "+w.path+" ("+w.allowedRole+")", func(t *testing.T) {
			router, svc := newDirectoryRouter("sid-" + w.allowedRole)
			code := doDirectoryRequest(router, w)
			if code == http.StatusUnauthorized || code == http.StatusForbidden {
				t.Fatalf("→ %d, el rol %q debería tener permiso", code, w.allowedRole)
			}
			if svc.calls != 1 {
				t.Errorf("el servicio se ejecutó %d veces, esperaba 1", svc.calls)
			}
		})
	}
}

// Entidad puede crear pero no actualizar ni desactivar: si el guard usara un solo
// permiso para todas las escrituras, esto pasaría desapercibido en un sentido o en
// el otro.
func TestDirectoryUpdateRejectsEntidadRole(t *testing.T) {
	utils.AddCommonSession("sid-et2", &utils.CommonSession{UserICode: "u-et", CurrentRole: "et"})
	router, svc := newDirectoryRouter("sid-et2")

	request := httptest.NewRequest(http.MethodPut, "/api/v1/directories/abc", strings.NewReader(`{}`))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusForbidden {
		t.Errorf("PUT con rol et → %d, esperaba 403", recorder.Code)
	}
	if svc.calls != 0 {
		t.Errorf("el servicio se ejecutó %d veces con rol et", svc.calls)
	}
}

// La carga masiva va en el mismo grupo que PUT, pero se prueba aparte: si alguien
// la saca de ahí, el test de PUT seguiría pasando y este endpoint quedaría abierto
// a que cualquiera reescriba las 3.400 entidades del Directorio. El guard corre
// antes de leer el archivo, así que no hace falta armar el multipart.
func TestDirectoryBulkUploadRejectsAnonymous(t *testing.T) {
	router, svc := newDirectoryRouter("")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/api/v1/directories/bulk-upload", nil))

	if recorder.Code != http.StatusUnauthorized {
		t.Errorf("POST bulk-upload sin sesión → %d, esperaba 401", recorder.Code)
	}
	if svc.calls != 0 {
		t.Errorf("el servicio se ejecutó %d veces sin sesión", svc.calls)
	}
}

// Las lecturas siguen abiertas: las consume la app móvil, que no tiene la sesión
// de la web.
func TestDirectoryReadsStayPublic(t *testing.T) {
	router, svc := newDirectoryRouter("")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/directories?type=FI", nil))

	if recorder.Code != http.StatusOK {
		t.Errorf("GET sin sesión → %d, esperaba 200", recorder.Code)
	}
	if svc.calls != 1 {
		t.Errorf("el servicio se ejecutó %d veces, esperaba 1", svc.calls)
	}
}
