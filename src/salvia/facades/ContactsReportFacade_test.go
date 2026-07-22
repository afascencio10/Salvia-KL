package salvia_facades

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func TestDownloadContactsReport_BindingValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.POST("/test-contacts-report", func(ctx *gin.Context) {
		var req DownloadReportRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error_code": "VALIDATION_FAILED", "message": "Fechas requeridas"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Petición con fechas vacías (debe fallar la validación)
	reqBody := map[string]string{
		"start_date": "",
		"end_date":   "",
	}
	bodyBytes, _ := json.Marshal(reqBody)
	c.Request = httptest.NewRequest("POST", "/test-contacts-report", bytes.NewBuffer(bodyBytes))
	c.Request.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, c.Request)

	if w.Code != http.StatusBadRequest {
		t.Errorf("se esperaba status 400 Bad Request, se obtuvo: %d", w.Code)
	}
}

func TestDownloadContactsReport_NoSession(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.POST("/api/v1/reportes/contactos-consolidado", func(ctx *gin.Context) {
		DownloadContactsReport(ctx, nil)
	})

	reqBody := map[string]string{"start_date": "2026-01-01", "end_date": "2026-06-30"}
	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/v1/reportes/contactos-consolidado", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("sin sesión se esperaba 401, se obtuvo: %d", w.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("respuesta no es JSON: %v", err)
	}
	if body["error_code"] != "AUTH_TOKEN_INVALID" {
		t.Errorf("se esperaba error_code AUTH_TOKEN_INVALID, se obtuvo: %q", body["error_code"])
	}
}
