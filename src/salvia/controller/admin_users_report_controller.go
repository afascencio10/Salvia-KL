// Package controller — admin_users_report_controller.go
// Genera un reporte Excel con todos los usuarios del sistema.
// Disponible solo para el rol administrador (ad).
package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type AdminUsersReportController struct {
	db *gorm.DB
}

func NewAdminUsersReportController(db *gorm.DB) *AdminUsersReportController {
	return &AdminUsersReportController{db: db}
}

func (c *AdminUsersReportController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/admin/users/report", c.DownloadReport)
}

// DownloadReport genera y descarga un Excel con todos los usuarios del sistema.
// GET /api/v1/admin/users/report
func (c *AdminUsersReportController) DownloadReport(ctx *gin.Context) {
	type UserRow struct {
		Login          string `gorm:"column:login"`
		Status         string `gorm:"column:status"`
		Team           string `gorm:"column:team"`
		Names          string `gorm:"column:names"`
		LastNames      string `gorm:"column:last_names"`
		DocType        string `gorm:"column:doc_type"`
		DocNumber      string `gorm:"column:doc_number"`
		Gender         string `gorm:"column:gender"`
		TownName       string `gorm:"column:town_name"`
		CityName       string `gorm:"column:city_name"`
		DepartmentName string `gorm:"column:department_name"`
		RoleName       string `gorm:"column:role_name"`
		RoleCode       string `gorm:"column:role_code"`
		Email          string `gorm:"column:email"`
		Phone          string `gorm:"column:phone"`
		CreatedAt      string `gorm:"column:created_at"`
		UpdatedAt      string `gorm:"column:updated_at"`
		AssignedDept   string `gorm:"column:assigned_dept"`
	}

	query := `
		SELECT
			gu.general_user_login AS login,
			gu.general_user_status AS status,
			COALESCE(gu.general_user_team, '') AS team,
			COALESCE(gup.general_user_profile_names, '') AS names,
			COALESCE(gup.general_user_profile_last_names, '') AS last_names,
			COALESCE(gup.general_user_profile_doc_type, '') AS doc_type,
			COALESCE(gup.general_user_profile_doc_number, '') AS doc_number,
			COALESCE(gup.general_user_profile_gender, '') AS gender,
			COALESCE(t.town_name, '') AS town_name,
			COALESCE(c.city_name, '') AS city_name,
			COALESCE(d.department_name, '') AS department_name,
			COALESCE(r.role_name, '') AS role_name,
			COALESCE(r.role_code, '') AS role_code,
			COALESCE((SELECT e.e_mail_data FROM security.email e WHERE e.e_mail_general_user_profile = gup.general_user_profile_id LIMIT 1), '') AS email,
			COALESCE((SELECT p.phone_number_data FROM security.phone_number p WHERE p.phone_number_general_user_profile = gup.general_user_profile_id LIMIT 1), '') AS phone,
			TO_CHAR(gu.general_user_creation_date, 'YYYY-MM-DD') AS created_at,
			TO_CHAR(gu.general_user_update_date, 'YYYY-MM-DD') AS updated_at,
			COALESCE(gu.general_user_assigned_department, '') AS assigned_dept
		FROM security.general_user gu
		JOIN security.general_user_profile gup ON gup.general_user_profile_id = gu.general_user_general_user_profile
		LEFT JOIN security.town t ON t.town_code = gup.general_user_profile_town
		LEFT JOIN security.city c ON c.city_id = t.city_id
		LEFT JOIN security.department d ON d.department_id = c.department_id
		LEFT JOIN security.rel_role_general_user rru ON rru.general_user_id = gu.general_user_id
		LEFT JOIN security.role r ON r.role_id = rru.role_id
		ORDER BY gup.general_user_profile_names ASC
	`

	var rows []UserRow
	c.db.Raw(query).Scan(&rows)

	// Generar Excel
	f := excelize.NewFile()
	sheet := "Usuarios"
	f.SetSheetName("Sheet1", sheet)

	// Headers
	headers := []string{
		"Usuario (Login)", "Estado", "Rol", "Código Rol", "Equipo", "Nombres", "Apellidos",
		"Tipo Documento", "Número Documento", "Género", "Correo", "Teléfono",
		"Municipio", "Ciudad", "Departamento", "Departamento Asignado",
		"Fecha Creación", "Última Actualización",
	}

	// Estilo de header
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"5106A7"}, Pattern: 1},
		Font: &excelize.Font{Bold: true, Color: "FFFFFF", Size: 11},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "374151", Style: 1},
			{Type: "top", Color: "374151", Style: 1},
			{Type: "right", Color: "374151", Style: 1},
			{Type: "bottom", Color: "374151", Style: 1},
		},
	})

	for i, h := range headers {
		colName, _ := excelize.ColumnNumberToName(i + 1)
		cell := colName + "1"
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}

	// Data rows
	genderMap := map[string]string{
		"m": "Hombre", "w": "Mujer", "i": "Intersexual",
		"ma": "Hombre", "fe": "Mujer",
		"ho": "Hombre", "mu": "Mujer",
		"ht": "Hombre Transgénero", "mt": "Mujer Transgénero",
		"tf": "Mujer Transgénero", "tm": "Hombre Transgénero",
		"nb": "No binaria", "nf": "No binaria",
		"n2": "No binaria",
		"ot": "Otra", "no": "No registra",
		"cl": "No binaria", "ci": "Cisgénero",
	}

	for rowIdx, r := range rows {
		row := rowIdx + 2
		statusLabel := "Inhabilitado"
		if r.Status == "e" {
			statusLabel = "Activo"
		}

		// Traducir género
		genero := r.Gender
		if translated, ok := genderMap[r.Gender]; ok {
			genero = translated
		}

		values := []interface{}{
			r.Login, statusLabel, r.RoleName, r.RoleCode, r.Team, r.Names, r.LastNames,
			r.DocType, r.DocNumber, genero, r.Email, r.Phone,
			r.TownName, r.CityName, r.DepartmentName, r.AssignedDept,
			r.CreatedAt, r.UpdatedAt,
		}

		for col, val := range values {
			colName, _ := excelize.ColumnNumberToName(col + 1)
			cell := fmt.Sprintf("%s%d", colName, row)
			f.SetCellValue(sheet, cell, val)
		}
	}

	// Auto-fit column widths
	for i := range headers {
		colName, _ := excelize.ColumnNumberToName(i + 1)
		f.SetColWidth(sheet, colName, colName, 18)
	}

	// Response
	filename := fmt.Sprintf("usuarios-salvia_%s.xlsx", time.Now().Format("2006-01-02"))
	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")

	if err := f.Write(ctx.Writer); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error generando el archivo Excel"})
		return
	}
	f.Close()
}
