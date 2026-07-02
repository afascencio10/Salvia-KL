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

func TestDownloadConsolidatedReport_CompileAndVerify(t *testing.T) {
	// Configurar Gin en modo test
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	// Configurar middleware de sesión mock
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.POST("/test-report", func(ctx *gin.Context) {
		// Test básico de binding del payload
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
	c.Request = httptest.NewRequest("POST", "/test-report", bytes.NewBuffer(bodyBytes))
	c.Request.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, c.Request)

	if w.Code != http.StatusBadRequest {
		t.Errorf("se esperaba status 400 Bad Request, se obtuvo: %d", w.Code)
	}
}
