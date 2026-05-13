package salvia_facades

import (
	"net/http"

	"bitsflow/common/db"

	"github.com/gin-gonic/gin"
)

// SyncSequencesPOST sincroniza todas las secuencias del schema salvia.
//
// POST /public/sync/sequences
//
// Sin body, sin autenticación. Solo para desarrollo/pruebas.
func SyncSequencesPOST(c *gin.Context) {
	if err := db.SyncAllSequences(&dbClientConfig, "salvia"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"ok":    false,
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ok":      true,
		"message": "Secuencias del schema 'salvia' sincronizadas correctamente.",
	})
}
