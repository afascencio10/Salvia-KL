// Package controller — admin_users_controller.go
// Endpoint de búsqueda de usuarios para la pantalla de administración.
package controller

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminUsersController struct {
	db *gorm.DB
}

func NewAdminUsersController(db *gorm.DB) *AdminUsersController {
	return &AdminUsersController{db: db}
}

func (c *AdminUsersController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/admin/users/search", c.Search)
}

// Search busca usuarios por nombre, documento o login con paginación del servidor.
// GET /api/v1/admin/users/search?name=...&doc=...&login=...&status=e|d&page=1&page_size=20
func (c *AdminUsersController) Search(ctx *gin.Context) {
	name := strings.TrimSpace(ctx.Query("name"))
	doc := strings.TrimSpace(ctx.Query("doc"))
	login := strings.TrimSpace(ctx.Query("login"))
	status := strings.TrimSpace(ctx.Query("status"))

	// Paginación
	page := 1
	pageSize := 20
	if p := ctx.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if ps := ctx.Query("page_size"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && v > 0 && v <= 100 {
			pageSize = v
		}
	}

	type UserResult struct {
		ICode     string `gorm:"column:icode"`
		Login     string `gorm:"column:login"`
		Status    string `gorm:"column:status"`
		Names     string `gorm:"column:names"`
		LastNames string `gorm:"column:last_names"`
		DocType   string `gorm:"column:doc_type"`
		DocNumber string `gorm:"column:doc_number"`
		RoleName  string `gorm:"column:role_name"`
	}

	// Base WHERE conditions
	whereClause := " WHERE 1=1"
	var args []interface{}

	if status != "" {
		whereClause += " AND gu.general_user_status = ?"
		args = append(args, status)
	}
	if name != "" {
		whereClause += " AND (LOWER(COALESCE(gup.general_user_profile_names,'') || ' ' || COALESCE(gup.general_user_profile_last_names,'')) LIKE ?)"
		args = append(args, "%"+strings.ToLower(name)+"%")
	}
	if doc != "" {
		whereClause += " AND gup.general_user_profile_doc_number LIKE ?"
		args = append(args, "%"+doc+"%")
	}
	if login != "" {
		whereClause += " AND LOWER(gu.general_user_login) LIKE ?"
		args = append(args, "%"+strings.ToLower(login)+"%")
	}

	baseFrom := `
		FROM security.general_user gu
		JOIN security.general_user_profile gup ON gup.general_user_profile_id = gu.general_user_general_user_profile
		LEFT JOIN security.rel_role_general_user rru ON rru.general_user_id = gu.general_user_id
		LEFT JOIN security.role r ON r.role_id = rru.role_id
	`

	// Count total
	countQuery := "SELECT COUNT(DISTINCT gu.general_user_id) " + baseFrom + whereClause
	var total int64
	c.db.Raw(countQuery, args...).Scan(&total)

	// Fetch page
	offset := (page - 1) * pageSize
	dataQuery := `
		SELECT DISTINCT
			gu.general_user_i_code AS icode,
			gu.general_user_login AS login,
			gu.general_user_status AS status,
			COALESCE(gup.general_user_profile_names, '') AS names,
			COALESCE(gup.general_user_profile_last_names, '') AS last_names,
			COALESCE(gup.general_user_profile_doc_type, '') AS doc_type,
			COALESCE(gup.general_user_profile_doc_number, '') AS doc_number,
			COALESCE(r.role_name, '') AS role_name
	` + baseFrom + whereClause + " ORDER BY names ASC LIMIT ? OFFSET ?"
	dataArgs := append(args, pageSize, offset)

	var results []UserResult
	c.db.Raw(dataQuery, dataArgs...).Scan(&results)

	// Transformar al formato que espera el frontend
	type ProfileDTO struct {
		Names     string `json:"names"`
		LastNames string `json:"lastNames"`
		DocType   string `json:"docType"`
		DocNumber string `json:"docNumber"`
	}
	type RoleDTO struct {
		Name string `json:"name"`
	}
	type UserDTO struct {
		ICode   string     `json:"icode"`
		Login   string     `json:"login"`
		Status  string     `json:"status"`
		Profile ProfileDTO `json:"profile"`
		Roles   []RoleDTO  `json:"roles"`
	}

	var users []UserDTO
	for _, r := range results {
		users = append(users, UserDTO{
			ICode:  r.ICode,
			Login:  r.Login,
			Status: r.Status,
			Profile: ProfileDTO{
				Names:     r.Names,
				LastNames: r.LastNames,
				DocType:   r.DocType,
				DocNumber: r.DocNumber,
			},
			Roles: []RoleDTO{{Name: r.RoleName}},
		})
	}
	if users == nil {
		users = []UserDTO{}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"users": users,
		"total": total,
		"page":  page,
		"pageSize": pageSize,
	})
}
