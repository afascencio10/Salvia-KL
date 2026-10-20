// Package controller — victim_case_v3_controller.go
// Endpoints JSON para Registro V3 (DocsMD/Screens/Registro V3). Sin
// dinamic-form (decisión A1) — el frontend manda directamente el mapa de
// respuestas de cada módulo.
package controller

import (
	"bitsflow/common/utils"
	"bitsflow/internal/repository"
	salvia_config "bitsflow/salvia/config"
	"bitsflow/salvia/service"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type VictimCaseV3Controller struct {
	svc service.VictimCaseV3Service
}

func NewVictimCaseV3Controller(svc service.VictimCaseV3Service) *VictimCaseV3Controller {
	return &VictimCaseV3Controller{svc: svc}
}

// RegisterRoutes registra las rutas de Registro V3.
//
//	GET    /api/v1/victim-cases-v3/buscar
//	POST   /api/v1/victim-cases-v3/draft
//	GET    /api/v1/victim-cases-v3/:icode
//	POST   /api/v1/victim-cases-v3/:icode/seccion/:modulo
//	POST   /api/v1/victim-cases-v3/:icode/tamizaje
//	POST   /api/v1/victim-cases-v3/:icode/finalizar
//	POST   /api/v1/victim-cases-v3/:icode/entidades
//	DELETE /api/v1/victim-cases-v3/entidades/:entityIcode
//	POST   /api/v1/victim-cases-v3/:icode/barreras
//	DELETE /api/v1/victim-cases-v3/barreras/:barrierIcode
func (c *VictimCaseV3Controller) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/victim-cases-v3/buscar", c.Search)
	rg.GET("/victim-cases-v3/enums", c.GetEnumCatalog)
	rg.POST("/victim-cases-v3/draft", c.CreateDraft)
	rg.GET("/victim-cases-v3/:icode", c.GetCase)
	rg.POST("/victim-cases-v3/:icode/seccion/:modulo", c.SaveSection)
	rg.POST("/victim-cases-v3/:icode/tamizaje", c.SaveTamizaje)
	rg.POST("/victim-cases-v3/:icode/finalizar", c.Finalize)
	rg.POST("/victim-cases-v3/:icode/entidades", c.AddEntity)
	rg.DELETE("/victim-cases-v3/entidades/:entityIcode", c.RemoveEntity)
	rg.POST("/victim-cases-v3/:icode/barreras", c.AddBarrier)
	rg.DELETE("/victim-cases-v3/barreras/:barrierIcode", c.RemoveBarrier)
}

// requireRegistroV3Access valida sesión + permiso manage_registro_v3 (op y
// sv, mismo criterio que Registro de Caso V2 / Gestión de Reportes).
func requireRegistroV3Access(ctx *gin.Context) *utils.CommonSession {
	sessionID, ok := sessions.Default(ctx).Get("userData").(string)
	if !ok || sessionID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no autenticado"})
		return nil
	}
	s, err := utils.GetCommonSession(sessionID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "sesión expirada"})
		return nil
	}
	if !utils.CheckPermission(salvia_config.PermissionsByRole, "manage_registro_v3", s.CurrentRole, ctx) {
		return nil
	}
	return s
}

// Search — Cabezote de búsqueda (E04). D7: solo documento y/o teléfono.
func (c *VictimCaseV3Controller) Search(ctx *gin.Context) {
	if requireRegistroV3Access(ctx) == nil {
		return
	}
	docType := ctx.Query("tipo_documento")
	docNumber := ctx.Query("documento")
	phone := ctx.Query("telefono")
	if docNumber == "" && phone == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "digite al menos un criterio (documento o teléfono)"})
		return
	}
	results, err := c.svc.Search(ctx.Request.Context(), docType, docNumber, phone)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al buscar"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"registros": results})
}

// CreateDraft — "Nuevo Registro" (E06).
func (c *VictimCaseV3Controller) CreateDraft(ctx *gin.Context) {
	s := requireRegistroV3Access(ctx)
	if s == nil {
		return
	}
	var body struct {
		DocType       string `json:"tipo_documento"`
		DocNumber     string `json:"documento"`
		Phone         string `json:"telefono"`
		ResidenceTown string `json:"municipio_residencia"`
		Channel       string `json:"canal"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	icode, err := c.svc.CreateDraft(ctx.Request.Context(), body.DocType, body.DocNumber, body.Phone, body.ResidenceTown, s.UserICode, body.Channel)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al crear el registro"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"icode": icode})
}

// GetCase — estado completo del Caso (header + respuestas + entidades +
// barreras), para repoblar el Formulario (ej. reabrir un Borrador).
func (c *VictimCaseV3Controller) GetCase(ctx *gin.Context) {
	if requireRegistroV3Access(ctx) == nil {
		return
	}
	icode := ctx.Param("icode")
	state, err := c.svc.GetCase(ctx.Request.Context(), icode)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al cargar el caso"})
		return
	}
	ctx.JSON(http.StatusOK, state)
}

// GetEnumCatalog — catálogo de opciones (code/icode/label) que necesita el
// Formulario para resolver sus dropdowns/checkboxes a icode (sin
// dinamic-form, A1, no hay otra forma de conocerlos). Público dentro del
// mismo criterio de acceso que el resto de Registro V3 -- es normativa, no
// dato de caso.
func (c *VictimCaseV3Controller) GetEnumCatalog(ctx *gin.Context) {
	if requireRegistroV3Access(ctx) == nil {
		return
	}
	catalog, err := c.svc.EnumCatalog(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al cargar el catálogo"})
		return
	}
	ctx.JSON(http.StatusOK, catalog)
}

// SaveSection — Cuando guarda una sección/módulo (E12). modulo = CI/DP/CA/HV/AG/GS
// (TZ usa SaveTamizaje aparte, por la lógica de ramas/puntaje).
func (c *VictimCaseV3Controller) SaveSection(ctx *gin.Context) {
	if requireRegistroV3Access(ctx) == nil {
		return
	}
	icode := ctx.Param("icode")
	modulo := strings.ToUpper(ctx.Param("modulo"))

	var body struct {
		Answers map[string]string `json:"answers" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	triadaCompleta, err := c.svc.SaveSection(ctx.Request.Context(), icode, modulo, body.Answers)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al guardar la sección"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"ok": true, "triadaCompleta": triadaCompleta})
}

// SaveTamizaje — guarda la rama pareja/no-pareja activa y calcula tipo/nivel/score.
func (c *VictimCaseV3Controller) SaveTamizaje(ctx *gin.Context) {
	if requireRegistroV3Access(ctx) == nil {
		return
	}
	icode := ctx.Param("icode")

	var body struct {
		Answers                 map[string]string `json:"answers" binding:"required"`
		WasPartner              bool              `json:"wasPartner"`
		EsOsigd                 bool              `json:"esOsigd"`
		SubtipoFeminicidio      bool              `json:"subtipoFeminicidio"`
		SubtipoHomicidioAbyecto bool              `json:"subtipoHomicidioAbyecto"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := c.svc.SaveTamizaje(ctx.Request.Context(), icode, body.Answers, body.WasPartner, body.EsOsigd, body.SubtipoFeminicidio, body.SubtipoHomicidioAbyecto)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al guardar el tamizaje"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

// Finalize — "Finalizar registro y enrutar a Seguimiento" (E15). Valida
// las obligatorias antes de transicionar a Terminado (ver flow-E15 §
// "Validar formulario" y el fix de QA del 2026-10-15 en
// victim_case_v3_service.go).
func (c *VictimCaseV3Controller) Finalize(ctx *gin.Context) {
	if requireRegistroV3Access(ctx) == nil {
		return
	}
	icode := ctx.Param("icode")
	result, err := c.svc.Finalize(ctx.Request.Context(), icode)
	if err != nil {
		var incompletos *service.ErrModulosIncompletos
		if errors.As(err, &incompletos) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":              "No se puede finalizar — hay módulos incompletos",
				"modulosIncompletos": incompletos.Modulos,
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al finalizar el registro"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"ok": true, "denunciaFacilSolicitada": result.DenunciaFacilSolicitada, "equiposRemisionCSV": result.EquiposRemisionCSV,
	})
}

// AddEntity — "+ Siguiente entidad" (E17).
func (c *VictimCaseV3Controller) AddEntity(ctx *gin.Context) {
	s := requireRegistroV3Access(ctx)
	if s == nil {
		return
	}
	icode := ctx.Param("icode")

	var body repository.VictimCaseEntityInput
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	body.VictimCaseICode = icode
	body.Agent = s.UserICode

	entityICode, err := c.svc.AddEntity(ctx.Request.Context(), body)
	if err != nil {
		if errors.Is(err, service.ErrEntityCaseDuplicate) {
			ctx.JSON(http.StatusConflict, gin.H{"error": "Esta entidad ya está asociada a este caso"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al registrar la entidad"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"icode": entityICode})
}

// RemoveEntity — "✕ Eliminar" sobre una entidad (E17).
func (c *VictimCaseV3Controller) RemoveEntity(ctx *gin.Context) {
	if requireRegistroV3Access(ctx) == nil {
		return
	}
	entityICode := ctx.Param("entityIcode")
	if err := c.svc.RemoveEntity(ctx.Request.Context(), entityICode); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al eliminar la entidad"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"ok": true})
}

// AddBarrier — "+ Siguiente barrera" (E18).
func (c *VictimCaseV3Controller) AddBarrier(ctx *gin.Context) {
	s := requireRegistroV3Access(ctx)
	if s == nil {
		return
	}
	icode := ctx.Param("icode")

	var body repository.VictimCaseBarrierInput
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	body.VictimCaseICode = icode
	body.Agent = s.UserICode

	barrierICode, err := c.svc.AddBarrier(ctx.Request.Context(), body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al registrar la barrera"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"icode": barrierICode})
}

// RemoveBarrier — "✕ Eliminar" sobre una barrera (E18).
func (c *VictimCaseV3Controller) RemoveBarrier(ctx *gin.Context) {
	if requireRegistroV3Access(ctx) == nil {
		return
	}
	barrierICode := ctx.Param("barrierIcode")
	if err := c.svc.RemoveBarrier(ctx.Request.Context(), barrierICode); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al eliminar la barrera"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"ok": true})
}
