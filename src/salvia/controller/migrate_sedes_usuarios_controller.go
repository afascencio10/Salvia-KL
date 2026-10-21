package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type sedesUsuariosResult struct {
	UsuariosActualizados int64 `json:"usuariosActualizados"`
	DryRun               bool  `json:"dryRun"`
}

// Condiciones compartidas por el conteo del dry run y el UPDATE real. El
// IS DISTINCT FROM es lo que hace la migración idempotente.
const sedesUsuariosWhere = `
	WHERE co.case_owner_general_user = u.general_user_i_code
	  AND co.entity_branch_id IS NOT NULL
	  AND u.entity_branch_id IS DISTINCT FROM co.entity_branch_id`

// MigrateSedesUsuarios copia la sede de salvia.case_owner a
// security.general_user.entity_branch_id, que es donde la sesión la busca.
//
//	POST /api/v1/admin/migrate/sedes-usuarios?dry_run=true
func (c *MigrateController) MigrateSedesUsuarios(ctx *gin.Context) {
	dryRun := ctx.Query("dry_run") == "true"
	res := sedesUsuariosResult{DryRun: dryRun}
	db := c.db.WithContext(ctx.Request.Context())

	if dryRun {
		err := db.Raw(`SELECT COUNT(*) FROM security.general_user u, salvia.case_owner co` + sedesUsuariosWhere).
			Scan(&res.UsuariosActualizados).Error
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error contando usuarios: " + err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, res)
		return
	}

	tx := db.Exec(`UPDATE security.general_user u
		SET entity_branch_id = co.entity_branch_id
		FROM salvia.case_owner co` + sedesUsuariosWhere)
	if tx.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error actualizando usuarios: " + tx.Error.Error()})
		return
	}
	res.UsuariosActualizados = tx.RowsAffected

	ctx.JSON(http.StatusOK, res)
}
