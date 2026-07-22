// Package controller — migrate_controller.go
// Endpoints administrativos para migración de datos legacy.
// Protegidos por X-Security-Key en el header.
// Se ejecutan una sola vez desde Postman/curl durante el despliegue.
package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"bitsflow/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

const migrateSecurityKey = "SALVIA_MIGRATE_2026_PROD"

type MigrateController struct {
	db *gorm.DB
}

// Caché en memoria de la estructura de formularios (secciones/preguntas por
// order), usada por MigrateFollowUp — ver comentario en getFormStructure.
type formStructureCacheEntry struct {
	sectionIDByOrder    map[int]string
	questionIDBySection map[string]map[int]string
}

var (
	formStructureCache   = make(map[string]*formStructureCacheEntry)
	formStructureCacheMu sync.RWMutex
)

func NewMigrateController(db *gorm.DB) *MigrateController {
	return &MigrateController{db: db}
}

func (c *MigrateController) RegisterRoutes(rg *gin.RouterGroup) {
	admin := rg.Group("/admin/migrate")
	admin.Use(c.authMiddleware())
	{
		admin.POST("/cases", c.MigrateCases)
		admin.POST("/excel-seguimientos", c.MigrateFromExcel)
		admin.POST("/excel-seguimientos/preview", c.PreviewExcel)
		admin.POST("/follow-up", c.MigrateFollowUp)
		admin.POST("/create-form", c.CreateForm)
	}
}

// authMiddleware valida el header X-Security-Key
func (c *MigrateController) authMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		key := ctx.GetHeader("X-Security-Key")
		if key != migrateSecurityKey {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "clave de seguridad inválida"})
			return
		}
		ctx.Next()
	}
}

// MigrateCases realiza la carga de seguimientos para un caso específico.
// POST /api/v1/admin/migrate/cases
// Header: X-Security-Key: SALVIA_MIGRATE_2026_PROD
// Body JSON:
//
//	{
//	  "cedula": "1099000005",
//	  "seguimientos": [
//	    { "fecha": "2026-06-01", "hora": "09:00", "estado": "PENDIENTE" },
//	    { "fecha": "2026-06-15", "hora": "", "estado": "PENDIENTE" },
//	    { "fecha": "2026-05-10", "hora": "", "estado": "REALIZADO" }
//	  ]
//	}
func (c *MigrateController) MigrateCases(ctx *gin.Context) {
	var body struct {
		Cedula       string `json:"cedula" binding:"required"`
		Seguimientos []struct {
			Fecha  string `json:"fecha" binding:"required"`
			Hora   string `json:"hora"`
			Estado string `json:"estado"`
		} `json:"seguimientos" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[MIGRATE] Buscando caso para cédula: %s", body.Cedula)

	// 1. Buscar el caso por número de documento
	type casoRow struct {
		VictimCaseId    int64  `gorm:"column:victim_case_id"`
		VictimCaseICode string `gorm:"column:victim_case_i_code"`
		VictimCaseTeam  string `gorm:"column:victim_case_team"`
	}
	var caso casoRow
	err := c.db.Raw(`
		SELECT victim_case_id, victim_case_i_code, COALESCE(victim_case_team, '') AS victim_case_team
		FROM salvia.victim_case
		WHERE victim_case_victim_doc_number = ?
		  AND victim_case_status IN ('ra', 'is')
		ORDER BY victim_case_creation_date DESC
		LIMIT 1
	`, body.Cedula).Scan(&caso).Error
	if err != nil || caso.VictimCaseICode == "" {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No se encontró un caso activo para la cédula: " + body.Cedula})
		return
	}

	log.Printf("[MIGRATE] Caso encontrado: id=%d icode=%s team=%s", caso.VictimCaseId, caso.VictimCaseICode, caso.VictimCaseTeam)

	// 2. Obtener el último agente asignado al caso
	var agentID string
	c.db.Raw(`
		SELECT co.case_owner_general_user
		FROM salvia.rel_case_owner_victim_case rcov
		JOIN salvia.case_owner co ON co.case_owner_id = rcov.case_owner_id
		WHERE rcov.victim_case_id = ?
		ORDER BY rcov.rel_case_owner_victim_case_creation_date DESC
		LIMIT 1
	`, caso.VictimCaseId).Scan(&agentID)

	log.Printf("[MIGRATE] Agente asignado: %s", agentID)

	// 3. Obtener nivel de riesgo para risk_status
	var riskLevel *int
	c.db.Raw(`SELECT victim_case_form2_risk_level FROM salvia.victim_case_form2 WHERE victim_case_form2_victim_case = ?`, caso.VictimCaseId).Scan(&riskLevel)

	riskStatus := "SIN_EVALUAR"
	if riskLevel != nil {
		switch *riskLevel {
		case 1:
			riskStatus = "BAJO"
		case 2:
			riskStatus = "MODERADO"
		case 3:
			riskStatus = "ALTO"
		case 4:
			riskStatus = "EXTREMO"
		}
	}

	// Team: usar el del caso, o calcular por riesgo
	team := caso.VictimCaseTeam
	if team == "" {
		if riskLevel != nil && *riskLevel >= 3 {
			team = "Riesgo alto"
		} else {
			team = "Riesgo bajo"
		}
		// Insertar el team calculado en la tabla victim_case
		c.db.Exec(`UPDATE salvia.victim_case SET victim_case_team = ? WHERE victim_case_i_code = ?`, team, caso.VictimCaseICode)
		log.Printf("[MIGRATE] Team asignado al caso %s: %s", caso.VictimCaseICode, team)
	}

	// 4. Crear los seguimientos en follow_up_v2
	creados := 0
	type segCreado struct {
		ID     string `json:"id"`
		Fecha  string `json:"fecha"`
		Estado string `json:"estado"`
		Seq    int    `json:"sequenceNumber"`
	}
	var seguimientosCreados []segCreado

	for i, seg := range body.Seguimientos {
		estado := seg.Estado
		if estado == "" {
			estado = "PENDIENTE"
		}

		// Parsear fecha en formato DD/MM/YYYY o YYYY-MM-DD
		fechaStr := seg.Fecha
		if len(fechaStr) == 10 && fechaStr[2] == '/' {
			fechaStr = fechaStr[6:10] + "-" + fechaStr[3:5] + "-" + fechaStr[0:2]
		}

		seqNum := i + 1
		hora := seg.Hora

		// Generar UUID y crear el seguimiento
		var newID string
		err := c.db.Raw(`
			INSERT INTO salvia.follow_up_v2 (
				id, case_id, agent_id, status, team, risk_status,
				scheduled_date, scheduled_time, sequence_number,
				attempts, is_priority, created_at, updated_at
			) VALUES (
				gen_random_uuid(), ?, ?, ?, ?, ?,
				(?::timestamp AT TIME ZONE 'America/Bogota'), ?, ?,
				0, false, NOW(), NOW()
			) RETURNING id
		`,
			caso.VictimCaseICode, agentID, estado, team, riskStatus,
			fechaStr, hora, seqNum,
		).Scan(&newID).Error

		if err != nil {
			log.Printf("[MIGRATE] Error creando seguimiento #%d: %v", seqNum, err)
		} else {
			creados++
			seguimientosCreados = append(seguimientosCreados, segCreado{
				ID:     newID,
				Fecha:  fechaStr,
				Estado: estado,
				Seq:    seqNum,
			})
		}
	}

	log.Printf("[MIGRATE] Migración completada para cédula %s: %d/%d seguimientos creados", body.Cedula, creados, len(body.Seguimientos))

	ctx.JSON(http.StatusOK, gin.H{
		"success":      true,
		"message":      fmt.Sprintf("Caso %s: %d seguimientos creados de %d solicitados", caso.VictimCaseICode, creados, len(body.Seguimientos)),
		"caso":         caso.VictimCaseICode,
		"agente":       agentID,
		"team":         team,
		"riesgo":       riskStatus,
		"creados":      creados,
		"total":        len(body.Seguimientos),
		"seguimientos": seguimientosCreados,
	})
}

// MigrateFromExcel procesa un archivo Excel con la hoja "Riesgo nuevo".
// POST /api/v1/admin/migrate/excel
// Header: X-Security-Key: SALVIA_MIGRATE_2026_PROD
// Body: multipart/form-data con campo "file" (archivo .xlsx)
//
// Estructura del Excel (hoja "Riesgo nuevo"):
//
//	Columnas de fecha: "Fecha primer seguimiento", "Fecha segundo seguimiento", etc.
//	Columna "Alerta pendientes": indica cuál seguimiento es el próximo PENDIENTE (ej: "Segundo seguimiento")
//	Columna "Caso cerrado": si dice "si" se omite la fila
//	Columna "Documento": cédula de la víctima
//
// Lógica:
//   - Seguimientos con fecha ANTES del indicado en "Alerta pendientes" → REALIZADO
//   - El seguimiento indicado en "Alerta pendientes" → PENDIENTE (usa su fecha si tiene)
//   - Seguimientos posteriores con fecha → PENDIENTE
func (c *MigrateController) MigrateFromExcel(ctx *gin.Context) {
	file, _, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Se requiere un archivo Excel en el campo 'file': " + err.Error()})
		return
	}
	defer file.Close()

	f, err := excelize.OpenReader(file)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No se pudo leer el archivo Excel: " + err.Error()})
		return
	}
	defer f.Close()

	// Buscar la hoja "Riesgo nuevo"
	sheetName := ""
	for _, name := range f.GetSheetList() {
		if strings.EqualFold(strings.TrimSpace(name), "riesgo nuevo") {
			sheetName = name
			break
		}
	}
	if sheetName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No se encontró la hoja 'Riesgo nuevo' en el archivo"})
		return
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error leyendo filas: " + err.Error()})
		return
	}
	if len(rows) < 2 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "La hoja no tiene datos (solo encabezado o vacía)"})
		return
	}

	// Mapear encabezados a índices (case-insensitive, trimmed, sin saltos de línea)
	header := rows[0]
	colIndex := make(map[string]int)
	for i, h := range header {
		// Normalizar: quitar saltos de línea, espacios extra, lowercase
		normalized := strings.ToLower(strings.Join(strings.Fields(h), " "))
		colIndex[normalized] = i
	}

	// Buscar columna de documento/cédula
	colCedula := findColIndex(colIndex, []string{"documento", "cedula", "cédula", "tipo id"})
	if colCedula == -1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No se encontró columna 'Documento'", "encabezados": header})
		return
	}

	// Buscar columna "Alerta pendientes"
	colAlerta := findColIndex(colIndex, []string{"alerta pendientes", "alerta pendiente", "alertapendientes"})
	if colAlerta == -1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No se encontró columna 'Alerta pendientes'", "encabezados": header})
		return
	}

	// Buscar columna "Caso cerrado"
	colCerrado := findColIndex(colIndex, []string{"caso cerrado", "caso_cerrado"})
	if colCerrado == -1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No se encontró columna 'Caso cerrado'", "encabezados": header})
		return
	}

	// Mapear columnas de fecha de seguimiento por número ordinal
	// "fecha primer seguimiento" → 1, "fecha segundo seguimiento" → 2, etc.
	ordinalMap := map[string]int{
		"primer": 1, "primero": 1, "1": 1,
		"segundo": 2, "2": 2,
		"tercer": 3, "tercero": 3, "3": 3,
		"cuarto": 4, "4": 4,
		"quinto": 5, "5": 5,
		"sexto": 6, "6": 6,
		"septimo": 7, "séptimo": 7, "7": 7,
		"octavo": 8, "8": 8,
	}

	type fechaCol struct {
		ColIndex int
		SeqNum   int
	}
	var fechaCols []fechaCol

	for normalized, idx := range colIndex {
		if strings.Contains(normalized, "fecha") && strings.Contains(normalized, "seguimiento") {
			// Extraer el ordinal de la columna
			for ordinal, num := range ordinalMap {
				if strings.Contains(normalized, ordinal) {
					fechaCols = append(fechaCols, fechaCol{ColIndex: idx, SeqNum: num})
					break
				}
			}
		}
	}

	// Ordenar por SeqNum
	for i := 0; i < len(fechaCols)-1; i++ {
		for j := i + 1; j < len(fechaCols); j++ {
			if fechaCols[i].SeqNum > fechaCols[j].SeqNum {
				fechaCols[i], fechaCols[j] = fechaCols[j], fechaCols[i]
			}
		}
	}

	// Mapear columnas "¿Se hizo seguimiento?" por número
	// "¿se hizo seguimiento?" sin número → seg 1
	// "¿se hizo seguimiento?2" → seg 2, etc.
	seHizoColMap := make(map[int]int)
	for normalized, idx := range colIndex {
		if strings.Contains(normalized, "se hizo") && strings.Contains(normalized, "seguimiento") {
			num := extractNumber(normalized)
			if num == 0 {
				// Sin número = seguimiento 1 (solo si no contiene dígito alguno)
				num = 1
			}
			seHizoColMap[num] = idx
			log.Printf("[MIGRATE-EXCEL]   seHizo #%d → columna %d (%s)", num, idx, header[idx])
		}
	}
	// Fallback: si no encontró por "se hizo", buscar por "hizo" solamente
	if len(seHizoColMap) == 0 {
		for normalized, idx := range colIndex {
			if strings.Contains(normalized, "hizo") && strings.Contains(normalized, "seguimiento") {
				num := extractNumber(normalized)
				if num == 0 {
					num = 1
				}
				seHizoColMap[num] = idx
				log.Printf("[MIGRATE-EXCEL]   seHizo(fallback) #%d → columna %d (%s)", num, idx, header[idx])
			}
		}
	}

	// Buscar columna "Login" para identificar al operador asignado al caso
	colLogin := findColIndex(colIndex, []string{"login", "login operador", "usuario"})

	log.Printf("[MIGRATE-EXCEL] Hoja: %s | Filas: %d | Columnas fecha: %d | Columnas seHizo: %d | Col login: %d", sheetName, len(rows)-1, len(fechaCols), len(seHizoColMap), colLogin)
	for _, fc := range fechaCols {
		log.Printf("[MIGRATE-EXCEL]   Seg #%d → columna %d (%s)", fc.SeqNum, fc.ColIndex, header[fc.ColIndex])
	}

	// Procesar cada fila
	type resultRow struct {
		Fila    int    `json:"fila"`
		Cedula  string `json:"cedula"`
		Caso    string `json:"caso,omitempty"`
		Creados int    `json:"creados"`
		Omitido bool   `json:"omitido,omitempty"`
		Motivo  string `json:"motivo,omitempty"`
		Error   string `json:"error,omitempty"`
	}
	var resultados []resultRow
	totalCreados := 0
	totalOmitidos := 0
	totalErrores := 0

	for rowIdx := 1; rowIdx < len(rows); rowIdx++ {
		row := rows[rowIdx]
		fila := rowIdx + 1

		// Obtener cédula
		cedula := getCellValue(row, colCedula)
		if cedula == "" {
			continue
		}

		// Verificar si caso cerrado
		casoCerrado := strings.ToLower(getCellValue(row, colCerrado))
		if casoCerrado == "si" || casoCerrado == "sí" || casoCerrado == "x" {
			totalOmitidos++
			resultados = append(resultados, resultRow{Fila: fila, Cedula: cedula, Omitido: true, Motivo: "caso cerrado"})
			continue
		}

		// Leer "Alerta pendientes" — determina cuál seguimiento es PENDIENTE
		alertaTexto := strings.ToLower(getCellValue(row, colAlerta))
		numAlerta := resolveOrdinal(alertaTexto, ordinalMap)

		// Buscar caso por cédula
		type casoRow struct {
			VictimCaseId    int64  `gorm:"column:victim_case_id"`
			VictimCaseICode string `gorm:"column:victim_case_i_code"`
			VictimCaseTeam  string `gorm:"column:victim_case_team"`
		}
		var caso casoRow
		err := c.db.Raw(`
			SELECT victim_case_id, victim_case_i_code, COALESCE(victim_case_team, '') AS victim_case_team
			FROM salvia.victim_case
			WHERE victim_case_victim_doc_number = ?
			  AND victim_case_status IN ('ra', 'is')
			ORDER BY victim_case_creation_date DESC
			LIMIT 1
		`, cedula).Scan(&caso).Error
		if err != nil || caso.VictimCaseICode == "" {
			totalErrores++
			resultados = append(resultados, resultRow{Fila: fila, Cedula: cedula, Error: "caso no encontrado"})
			continue
		}

		// Buscar agente por login del Excel (columna "Login") — solo si existe la columna y tiene valor
		var agentID string
		if colLogin >= 0 {
			loginValue := getCellValue(row, colLogin)
			if loginValue != "" && strings.ToLower(loginValue) != "na" {
				agentID = c.findAgentByName(loginValue)
			}
		}
		// Si no hay login o no se encontró, NO asignar a nadie (dejar vacío)
		// Ya no se usa el fallback del owner

		// Si encontramos agente por login, actualizar victim_case.agent_id y victim_case_team
		team := caso.VictimCaseTeam
		if agentID != "" {
			agentTeam := c.getAgentTeam(agentID)
			if agentTeam != "" {
				c.db.Exec(`UPDATE salvia.victim_case SET agent_id = ?, victim_case_team = ? WHERE victim_case_i_code = ?`,
					agentID, agentTeam, caso.VictimCaseICode)
				team = agentTeam
			} else {
				c.db.Exec(`UPDATE salvia.victim_case SET agent_id = ? WHERE victim_case_i_code = ?`,
					agentID, caso.VictimCaseICode)
			}
		}

		// Obtener nivel de riesgo
		var riskLevel *int
		c.db.Raw(`SELECT victim_case_form2_risk_level FROM salvia.victim_case_form2 WHERE victim_case_form2_victim_case = ?`, caso.VictimCaseId).Scan(&riskLevel)

		riskStatus := "SIN_EVALUAR"
		if riskLevel != nil {
			switch *riskLevel {
			case 1:
				riskStatus = "BAJO"
			case 2:
				riskStatus = "MODERADO"
			case 3:
				riskStatus = "ALTO"
			case 4:
				riskStatus = "EXTREMO"
			}
		}

		// Team: si no se asignó del agente, calcular por riesgo
		if team == "" {
			if riskLevel != nil && *riskLevel >= 3 {
				team = "Riesgo alto"
			} else {
				team = "Riesgo bajo"
			}
			c.db.Exec(`UPDATE salvia.victim_case SET victim_case_team = ? WHERE victim_case_i_code = ?`, team, caso.VictimCaseICode)
		}

		// Determinar el verdadero seguimiento pendiente:
		// Desde numAlerta, verificar "¿Se hizo seguimiento?" — si dice "Sí", avanzar al siguiente
		realPendiente := numAlerta
		if realPendiente > 0 {
			for seq := numAlerta; seq <= 8; seq++ {
				if seHizoIdx, ok := seHizoColMap[seq]; ok {
					valor := strings.ToLower(getCellValue(row, seHizoIdx))
					if valor == "sí" || valor == "si" || valor == "s" {
						realPendiente = seq + 1 // ya se hizo, avanzar
					} else {
						break // este es el verdadero pendiente
					}
				} else {
					break // no hay columna para este, es el pendiente
				}
			}
		}

		// Caso especial: "Seguimientos al día" → crear 1 solo seguimiento con fecha de hoy
		if numAlerta == 0 && strings.Contains(alertaTexto, "al d") {
			// Verificar si ya existe algún seguimiento para este caso
			var existCount int64
			c.db.Raw(`SELECT COUNT(*) FROM salvia.follow_up_v2 WHERE case_id = ?`, caso.VictimCaseICode).Scan(&existCount)

			if existCount == 0 {
				loc, _ := time.LoadLocation("America/Bogota")
				hoy := time.Now().In(loc).Format("2006-01-02")
				var newID string
				c.db.Raw(`
					INSERT INTO salvia.follow_up_v2 (
						id, case_id, agent_id, status, team, risk_status,
						scheduled_date, scheduled_time, sequence_number,
						attempts, is_priority, created_at, updated_at
					) VALUES (
						gen_random_uuid(), ?, ?, 'PENDIENTE', ?, ?,
						(?::date + interval '12 hours'), '', 1,
						0, false, NOW(), NOW()
					) RETURNING id
				`, caso.VictimCaseICode, agentID, team, riskStatus, hoy).Scan(&newID)
				if newID != "" {
					totalCreados++
					resultados = append(resultados, resultRow{Fila: fila, Cedula: cedula, Caso: caso.VictimCaseICode, Creados: 1})
				}
			} else {
				resultados = append(resultados, resultRow{Fila: fila, Cedula: cedula, Caso: caso.VictimCaseICode, Creados: 0})
			}
			continue
		}

		// Solo crear seguimientos PENDIENTES (desde realPendiente en adelante)
		creados := 0
		for _, fc := range fechaCols {
			// Ignorar seguimientos anteriores al pendiente real
			if fc.SeqNum < realPendiente {
				continue
			}

			fechaStr := getCellValue(row, fc.ColIndex)
			if fechaStr == "" {
				continue
			}

			// Parsear la fecha y hora
			fechaParsed, horaParsed := parseDateTimeFromExcel(fechaStr)
			if fechaParsed == "" {
				log.Printf("[MIGRATE-EXCEL] Fila %d, seg %d: fecha no parseable: '%s'", fila, fc.SeqNum, fechaStr)
				continue
			}

			// Todos los que creamos son PENDIENTE
			estado := "PENDIENTE"
			segAgentID := agentID

			// Verificar si ya existe un seguimiento con el mismo case_id y sequence_number
			var existingCount int64
			c.db.Raw(`SELECT COUNT(*) FROM salvia.follow_up_v2 WHERE case_id = ? AND sequence_number = ?`,
				caso.VictimCaseICode, fc.SeqNum).Scan(&existingCount)

			if existingCount > 0 {
				if segAgentID != "" {
					c.db.Exec(`UPDATE salvia.follow_up_v2 SET agent_id = ? WHERE case_id = ? AND sequence_number = ?`,
						segAgentID, caso.VictimCaseICode, fc.SeqNum)
				}
				log.Printf("[MIGRATE-EXCEL] Fila %d, seg %d: ya existe, omitiendo creación", fila, fc.SeqNum)
				continue
			}

			// Guardar la fecha con la hora del Excel en zona Colombia
			scheduledTime := horaParsed
			dateSQL := fechaParsed + " 12:00:00"
			if horaParsed != "" {
				dateSQL = fechaParsed + " " + horaParsed + ":00"
			}

			var newID string
			insertErr := c.db.Raw(`
				INSERT INTO salvia.follow_up_v2 (
					id, case_id, agent_id, status, team, risk_status,
					scheduled_date, scheduled_time, sequence_number,
					attempts, is_priority, created_at, updated_at
				) VALUES (
					gen_random_uuid(), ?, ?, ?, ?, ?,
					(?::timestamp AT TIME ZONE 'America/Bogota'), ?, ?,
					0, false, NOW(), NOW()
				) RETURNING id
			`,
				caso.VictimCaseICode, segAgentID, estado, team, riskStatus,
				dateSQL, scheduledTime, fc.SeqNum,
			).Scan(&newID).Error

			if insertErr != nil {
				log.Printf("[MIGRATE-EXCEL] Fila %d, seg %d: error INSERT: %v", fila, fc.SeqNum, insertErr)
			} else {
				creados++
			}
		}

		totalCreados += creados
		resultados = append(resultados, resultRow{Fila: fila, Cedula: cedula, Caso: caso.VictimCaseICode, Creados: creados})
	}

	log.Printf("[MIGRATE-EXCEL] Completado: %d creados, %d omitidos, %d errores", totalCreados, totalOmitidos, totalErrores)

	ctx.JSON(http.StatusOK, gin.H{
		"success":        true,
		"total_filas":    len(rows) - 1,
		"total_creados":  totalCreados,
		"total_omitidos": totalOmitidos,
		"total_errores":  totalErrores,
		"detalle":        resultados,
	})
}

// findAgentByName busca un usuario por login en general_user.
// Retorna el general_user_i_code o vacío si no lo encuentra.
func (c *MigrateController) findAgentByName(login string) string {
	login = strings.TrimSpace(login)
	if login == "" || strings.ToLower(login) == "na" {
		return ""
	}

	var userICode string
	c.db.Raw(`
		SELECT general_user_i_code
		FROM security.general_user
		WHERE LOWER(general_user_login) = LOWER(?)
		  AND general_user_status = 'e'
		LIMIT 1
	`, login).Scan(&userICode)

	if userICode != "" {
		log.Printf("[MIGRATE-EXCEL] Operador encontrado por login '%s' → %s", login, userICode)
	}

	return userICode
}

// getAgentTeam obtiene el team del agente desde general_user.
func (c *MigrateController) getAgentTeam(agentICode string) string {
	if agentICode == "" {
		return ""
	}
	var team string
	c.db.Raw(`SELECT COALESCE(general_user_team, '') FROM security.general_user WHERE general_user_i_code = ?`, agentICode).Scan(&team)
	return team
}

// extractNumber extrae el primer número encontrado en un string.
func extractNumber(s string) int {
	num := 0
	found := false
	for _, ch := range s {
		if ch >= '0' && ch <= '9' {
			num = num*10 + int(ch-'0')
			found = true
		} else if found {
			break
		}
	}
	return num
}

// findColIndex busca en el mapa de columnas normalizadas usando una lista de aliases.
func findColIndex(colIndex map[string]int, aliases []string) int {
	for _, alias := range aliases {
		if idx, ok := colIndex[alias]; ok {
			return idx
		}
	}
	// Búsqueda parcial: si algún encabezado CONTIENE el alias
	for _, alias := range aliases {
		for normalized, idx := range colIndex {
			if strings.Contains(normalized, alias) {
				return idx
			}
		}
	}
	return -1
}

// getCellValue obtiene el valor de una celda de forma segura.
func getCellValue(row []string, colIdx int) string {
	if colIdx < 0 || colIdx >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[colIdx])
}

// resolveOrdinal convierte texto como "segundo seguimiento" al número 2.
func resolveOrdinal(texto string, ordinalMap map[string]int) int {
	texto = strings.ToLower(strings.TrimSpace(texto))
	if texto == "" {
		return 0
	}
	// Buscar directamente en el mapa
	for ordinal, num := range ordinalMap {
		if strings.Contains(texto, ordinal) {
			return num
		}
	}
	// Intentar extraer número directo
	return extractNumber(texto)
}

// parseDateTimeFromExcel convierte una fecha del Excel a formato YYYY-MM-DD y extrae la hora si existe.
// Retorna (fecha YYYY-MM-DD, hora HH:MM o vacío si es 00:00)
func parseDateTimeFromExcel(s string) (string, string) {
	s = strings.TrimSpace(s)
	if s == "" || strings.ToLower(s) == "na" || strings.ToLower(s) == "n/a" {
		return "", ""
	}

	// Separar fecha de hora
	datePart := s
	timePart := ""
	if idx := strings.Index(s, " "); idx > 0 {
		datePart = s[:idx]
		timePart = strings.TrimSpace(s[idx+1:])
	}

	// Parsear la parte de fecha
	fecha := parseDateOnly(datePart)
	if fecha == "" {
		return "", ""
	}

	// Parsear la hora: si es "00:00" o "0:00" → no guardar hora
	hora := ""
	if timePart != "" && timePart != "00:00" && timePart != "0:00" {
		// Normalizar hora a HH:MM
		hParts := strings.Split(timePart, ":")
		if len(hParts) >= 2 {
			h := hParts[0]
			m := hParts[1]
			if len(h) == 1 {
				h = "0" + h
			}
			if len(m) == 1 {
				m = "0" + m
			}
			hora = h + ":" + m
		}
	}

	return fecha, hora
}

// parseDateOnly convierte solo la parte de fecha a YYYY-MM-DD.
func parseDateOnly(datePart string) string {
	datePart = strings.TrimSpace(datePart)
	if datePart == "" {
		return ""
	}

	// Formato DD/MM/YYYY (10 chars, día y mes con 2 dígitos)
	if len(datePart) == 10 && datePart[2] == '/' && datePart[5] == '/' {
		return datePart[6:10] + "-" + datePart[3:5] + "-" + datePart[0:2]
	}

	// Formato YYYY-MM-DD
	if len(datePart) == 10 && datePart[4] == '-' && datePart[7] == '-' {
		return datePart
	}

	// Formato DD-MM-YYYY
	if len(datePart) == 10 && datePart[2] == '-' && datePart[5] == '-' {
		return datePart[6:10] + "-" + datePart[3:5] + "-" + datePart[0:2]
	}

	// Formato M/D/YY o MM/DD/YY (formato corto americano con año 2 dígitos)
	if strings.Contains(datePart, "/") {
		parts := strings.Split(datePart, "/")
		if len(parts) == 3 {
			month := parts[0]
			day := parts[1]
			year := parts[2]

			if len(year) == 2 {
				yearNum := 0
				for _, ch := range year {
					yearNum = yearNum*10 + int(ch-'0')
				}
				year = fmt.Sprintf("20%02d", yearNum)
			}

			if len(month) == 1 {
				month = "0" + month
			}
			if len(day) == 1 {
				day = "0" + day
			}

			return year + "-" + month + "-" + day
		}
	}

	// Formato numérico de Excel (serial date number)
	num := 0
	isNumeric := true
	for _, ch := range datePart {
		if ch >= '0' && ch <= '9' {
			num = num*10 + int(ch-'0')
		} else if ch == '.' {
			break
		} else {
			isNumeric = false
			break
		}
	}
	if isNumeric && num > 0 {
		if num > 59 {
			num--
		}
		days := num - 1
		month := 1
		year := 1900
		totalDays := 1 + days
		daysInMonth := func(y, m int) int {
			switch m {
			case 1, 3, 5, 7, 8, 10, 12:
				return 31
			case 4, 6, 9, 11:
				return 30
			case 2:
				if (y%4 == 0 && y%100 != 0) || y%400 == 0 {
					return 29
				}
				return 28
			}
			return 30
		}
		for totalDays > daysInMonth(year, month) {
			totalDays -= daysInMonth(year, month)
			month++
			if month > 12 {
				month = 1
				year++
			}
		}
		return fmt.Sprintf("%04d-%02d-%02d", year, month, totalDays)
	}

	return ""
}

// parseDateFromExcel wrapper de compatibilidad — solo retorna la fecha.
func parseDateFromExcel(s string) string {
	fecha, _ := parseDateTimeFromExcel(s)
	return fecha
}

// PreviewExcel lee el Excel y devuelve los encabezados y las primeras filas para diagnóstico.
// POST /api/v1/admin/migrate/excel/preview
func (c *MigrateController) PreviewExcel(ctx *gin.Context) {
	file, _, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Se requiere un archivo Excel en el campo 'file': " + err.Error()})
		return
	}
	defer file.Close()

	f, err := excelize.OpenReader(file)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No se pudo leer el archivo Excel: " + err.Error()})
		return
	}
	defer f.Close()

	// Buscar la hoja "Riesgo nuevo"
	sheetName := ""
	for _, name := range f.GetSheetList() {
		if strings.EqualFold(strings.TrimSpace(name), "riesgo nuevo") {
			sheetName = name
			break
		}
	}
	if sheetName == "" {
		ctx.JSON(http.StatusOK, gin.H{"hojas": f.GetSheetList(), "error": "No se encontró 'Riesgo nuevo'"})
		return
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error leyendo filas: " + err.Error()})
		return
	}

	// Devolver encabezados y primeras 3 filas de datos
	preview := gin.H{
		"hoja":        sheetName,
		"total_filas": len(rows) - 1,
		"encabezados": rows[0],
	}

	var sampleRows []map[string]string
	limit := 3
	if len(rows)-1 < limit {
		limit = len(rows) - 1
	}
	for i := 1; i <= limit; i++ {
		rowMap := make(map[string]string)
		for j, cell := range rows[i] {
			colName := ""
			if j < len(rows[0]) {
				colName = rows[0][j]
			} else {
				colName = fmt.Sprintf("col_%d", j)
			}
			rowMap[colName] = cell
		}
		sampleRows = append(sampleRows, rowMap)
	}
	preview["muestra_filas"] = sampleRows

	ctx.JSON(http.StatusOK, preview)
}

// followUpExistente es lo mínimo que hace falta devolver cuando un external_ref
// ya tenía un follow_up_v2 creado — ver MigrateFollowUp.
type followUpExistente struct {
	ID               string `gorm:"column:id"`
	FormSubmissionID string `gorm:"column:form_submission_id"`
}

// buscarFollowUpPorExternalRef busca un follow_up_v2 ya creado para (form_id,
// external_ref). external_ref vive dentro de la columna JSONB kobo_metadata
// (clave "external_ref") en vez de una columna propia — se reusa el mismo campo
// que ya existía para metadata de origen, sea o no una migración de Kobo. Único
// por (form_id, external_ref) vía un índice parcial (ver mapeo-columnas.md /
// migrarBasesSalvia y DocsMD — creado a mano, no por AutoMigrate, porque
// requiere permisos de owner que el rol de la app no tiene en esta BD).
func (c *MigrateController) buscarFollowUpPorExternalRef(formID, externalRef string) (*followUpExistente, error) {
	var existente followUpExistente
	err := c.db.Raw(`
		SELECT id, form_submission_id FROM salvia.follow_up_v2
		WHERE form_id = ? AND kobo_metadata->>'external_ref' = ?
		LIMIT 1
	`, formID, externalRef).Scan(&existente).Error
	if err != nil {
		return nil, err
	}
	if existente.ID == "" {
		return nil, nil
	}
	return &existente, nil
}

// esErrorDeUnicidad detecta una violación del índice único de (form_id,
// external_ref) — código de error 23505 de Postgres. Detección por texto (igual
// que IsConnectionError en internal/db/gorm_connection.go) para no acoplarse al
// tipo de error concreto del driver.
func esErrorDeUnicidad(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "duplicate key") || strings.Contains(msg, "23505") || strings.Contains(msg, "already exists")
}

// getFormStructure resuelve secciones (por order) y preguntas (por sectionId+order)
// de un form, cacheadas en memoria por form_id tras la primera resolución.
//
// MigrateFollowUp se usa para migraciones masivas: miles de llamadas seguidas con
// el MISMO form_id, donde esta estructura nunca cambia dentro de la corrida. Antes
// se volvía a pedir a la BD en cada llamada (form_exists + form_section + question
// = 3 round-trips evitables, ~600-900ms de los ~2.5s totales medidos por llamada).
// Cachear indefinidamente por el tiempo de vida del proceso es seguro para este
// endpoint (admin, protegido por clave de seguridad, uso exclusivo de migraciones
// puntuales contra forms que no se editan mientras la migración corre) — si el
// form cambiara en caliente durante la corrida no se reflejaría hasta reiniciar
// el proceso, trade-off aceptado para este caso de uso.
func (c *MigrateController) getFormStructure(formID string) (map[int]string, map[string]map[int]string, error) {
	formStructureCacheMu.RLock()
	entry, ok := formStructureCache[formID]
	formStructureCacheMu.RUnlock()
	if ok {
		return entry.sectionIDByOrder, entry.questionIDBySection, nil
	}

	formStructureCacheMu.Lock()
	defer formStructureCacheMu.Unlock()
	// Revalidar tras tomar el lock de escritura: otra goroutine pudo haber
	// poblado el caché mientras esperábamos (llamadas concurrentes al mismo form_id).
	if entry, ok := formStructureCache[formID]; ok {
		return entry.sectionIDByOrder, entry.questionIDBySection, nil
	}

	var formExists int64
	c.db.Raw(`SELECT COUNT(*) FROM salvia.form WHERE id = ?`, formID).Scan(&formExists)
	if formExists == 0 {
		return nil, nil, fmt.Errorf("no existe un form con form_id: %s", formID)
	}

	type sectionRow struct {
		ID    string `gorm:"column:id"`
		Order int    `gorm:"column:order"`
	}
	var sections []sectionRow
	c.db.Raw(`SELECT id::text AS id, "order" FROM salvia.form_section WHERE form_id = ? AND deleted_at IS NULL`, formID).Scan(&sections)
	sectionIDByOrder := make(map[int]string, len(sections))
	for _, s := range sections {
		sectionIDByOrder[s.Order] = s.ID
	}

	type questionRow struct {
		ID            string `gorm:"column:id"`
		FormSectionID string `gorm:"column:form_section_id"`
		Order         int    `gorm:"column:order"`
	}
	var questions []questionRow
	c.db.Raw(`SELECT id, form_section_id, "order" FROM salvia.question WHERE form_id = ? AND deleted_at IS NULL`, formID).Scan(&questions)
	questionIDBySection := make(map[string]map[int]string, len(sections))
	for _, q := range questions {
		if questionIDBySection[q.FormSectionID] == nil {
			questionIDBySection[q.FormSectionID] = make(map[int]string)
		}
		questionIDBySection[q.FormSectionID][q.Order] = q.ID
	}

	formStructureCache[formID] = &formStructureCacheEntry{
		sectionIDByOrder:    sectionIDByOrder,
		questionIDBySection: questionIDBySection,
	}
	return sectionIDByOrder, questionIDBySection, nil
}

// MigrateFollowUp crea un follow_up_v2 + form_submission + answers en una sola llamada,
// pensado para el script de migración del Excel de KoBoToolbox (ver migrarKobos/).
//
// POST /api/v1/admin/migrate/follow-up
// Header: X-Security-Key: SALVIA_MIGRATE_2026_PROD
// Body JSON:
//
//	{
//	  "followUp": {
//	    "case_id": "019ea7f6-df8f-7bb4-a768-c9c9bf538945",
//	    "form_id": "f132614c-bd11-4871-9b3d-5a85bfab9abd",
//	    "agent_id": "",
//	    "status": "REALIZADO",
//	    "team": "Riesgo bajo",
//	    "risk_status": "BAJO",
//	    "scheduled_date": "2026-07-06",
//	    "scheduled_time": "10:30",
//	    "completed_at": "2026-07-06T10:45:00Z",
//	    "sequence_number": 1,
//	    "summary": "",
//	    "kobo_metadata": { "_id": "753500613", "_uuid": "b8262a05-...", "_submission_time": "2026-03-14T09:20:11" }
//	  },
//	  "formSubmission": {
//	    "section1": { "question1": "respuesta 1", "question2": "respuesta 2" },
//	    "section2": { "question1": "respuesta 1", "question3": "respuesta 2" }
//	  }
//	}
//
// "sectionN"/"questionM" refieren al `order` (1-indexado) de form_section/question dentro
// de followUp.form_id — NO al UUID real — así el script de migración solo necesita conocer
// la posición de cada columna, no los IDs de la BD. Preguntas con valor vacío se omiten.
// Claves que no resuelven a ninguna sección/pregunta real NO abortan la fila completa:
// se reportan en "warnings" y el resto de respuestas válidas se crea igual.
//
// "kobo_metadata" es opcional y se guarda tal cual (JSONB) en follow_up_v2.kobo_metadata —
// pensado para conservar la metadata cruda del sistema de origen (_id, _uuid,
// _submission_time, _index, etc. de KoBoToolbox) con fines de trazabilidad/auditoría.
func (c *MigrateController) MigrateFollowUp(ctx *gin.Context) {
	var body struct {
		FollowUp struct {
			CaseID         string                 `json:"case_id" binding:"required"`
			FormID         string                 `json:"form_id" binding:"required"`
			AgentID        string                 `json:"agent_id"`
			Status         string                 `json:"status"`
			Team           string                 `json:"team"`
			RiskStatus     string                 `json:"risk_status"`
			ScheduledDate  string                 `json:"scheduled_date"`
			ScheduledTime  string                 `json:"scheduled_time"`
			CompletedAt    string                 `json:"completed_at"`
			SequenceNumber int                    `json:"sequence_number"`
			Summary        string                 `json:"summary"`
			KoboMetadata   map[string]interface{} `json:"kobo_metadata"`
		} `json:"followUp" binding:"required"`
		FormSubmission map[string]map[string]string `json:"formSubmission"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fu := body.FollowUp

	// Idempotencia: si el llamador manda kobo_metadata.external_ref, es un
	// identificador estable que ÉL mismo construyó de forma determinística a
	// partir del origen (ej. el _uuid de un envío de KoBoToolbox, o
	// "archivo::hoja::filaN::bloqueM" para un Excel legacy sin ID propio) — el
	// mismo campo sirve para cualquier migración futura, no es específico de
	// Kobo pese al nombre de la columna que lo aloja.
	//
	// Si ya existe un follow_up_v2 con el mismo (form_id, external_ref), NO se
	// crea uno nuevo — se devuelve el existente con alreadyExisted:true. Esto
	// es lo que hace que reintentar/reanudar una migración interrumpida a la
	// mitad (kill -9, crash, corte de luz — cualquier cosa que el cliente no
	// pueda registrar a tiempo) nunca duplique un seguimiento: la garantía vive
	// en la BD, no en que el script haya alcanzado a guardar su log a tiempo.
	externalRef, _ := fu.KoboMetadata["external_ref"].(string)
	if externalRef != "" {
		if existente, err := c.buscarFollowUpPorExternalRef(fu.FormID, externalRef); err == nil && existente != nil {
			log.Printf("[MIGRATE-FOLLOWUP] external_ref=%s ya existía (followUpId=%s) — no se crea duplicado", externalRef, existente.ID)
			ctx.JSON(http.StatusOK, gin.H{
				"success":          true,
				"followUpId":       existente.ID,
				"formSubmissionId": existente.FormSubmissionID,
				"alreadyExisted":   true,
			})
			return
		}
	}

	// Validar que el caso exista (evita follow_up_v2/form_submission huérfanos) y
	// de paso traer su fecha de creación — se usa como fallback de scheduled_date
	// más abajo en vez de time.Now(), ver comentario ahí.
	type caseRow struct {
		CreationDate time.Time `gorm:"column:victim_case_creation_date"`
	}
	var caso caseRow
	c.db.Raw(`SELECT victim_case_creation_date FROM salvia.victim_case WHERE victim_case_i_code = ?`, fu.CaseID).Scan(&caso)
	if caso.CreationDate.IsZero() {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "no existe un victim_case con case_id: " + fu.CaseID})
		return
	}

	// Estructura del form (existencia + secciones/preguntas por order) — cacheada
	// en memoria por form_id, ver getFormStructure.
	sectionIDByOrder, questionIDBySection, err := c.getFormStructure(fu.FormID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Resolver sectionN.questionM → IDs reales
	type answerToCreate struct {
		QuestionID string
		Value      string
	}
	var answers []answerToCreate
	var warnings []string

	for sectionKey, preguntas := range body.FormSubmission {
		if !strings.HasPrefix(strings.ToLower(sectionKey), "section") {
			warnings = append(warnings, fmt.Sprintf("clave de sección inválida (se esperaba 'sectionN'): %s", sectionKey))
			continue
		}
		secOrder := extractNumber(sectionKey)
		sectionID, ok := sectionIDByOrder[secOrder]
		if !ok {
			warnings = append(warnings, fmt.Sprintf("no existe sección con order=%d en form_id=%s (clave: %s)", secOrder, fu.FormID, sectionKey))
			continue
		}

		for questionKey, value := range preguntas {
			if strings.TrimSpace(value) == "" {
				continue
			}
			if !strings.HasPrefix(strings.ToLower(questionKey), "question") {
				warnings = append(warnings, fmt.Sprintf("clave de pregunta inválida (se esperaba 'questionN'): %s.%s", sectionKey, questionKey))
				continue
			}
			qOrder := extractNumber(questionKey)
			questionID, ok := questionIDBySection[sectionID][qOrder]
			if !ok {
				warnings = append(warnings, fmt.Sprintf("no existe pregunta con order=%d en %s (clave: %s.%s)", qOrder, sectionKey, sectionKey, questionKey))
				continue
			}
			answers = append(answers, answerToCreate{QuestionID: questionID, Value: value})
		}
	}

	// Fechas: scheduled_date es NOT NULL en DB. Si no llega o no parsea (común en
	// datos legacy migrados: la celda de fecha del origen vino vacía o corrupta,
	// pero el resto del bloque sí tiene contenido real), se usa la fecha de
	// creación del caso como fallback — NO time.Now(), que fabricaba una fecha
	// falsa (la del momento de la migración) y hacía parecer que el seguimiento
	// ocurrió el día en que se corrió el script, en vez de dejarlo aproximado a
	// algún punto real del historial del caso.
	scheduledDate := parseFlexibleDateTime(fu.ScheduledDate)
	if scheduledDate.IsZero() {
		scheduledDate = caso.CreationDate
	}
	var completedAt *time.Time
	if fu.CompletedAt != "" {
		if t := parseFlexibleDateTime(fu.CompletedAt); !t.IsZero() {
			completedAt = &t
		}
	}

	status := fu.Status
	if status == "" {
		status = "REALIZADO"
	}

	// Metadata cruda del sistema de origen (ej. _id/_uuid/_submission_time de KoBoToolbox).
	// Se guarda tal cual llegó, sin validar su forma — es solo para trazabilidad/auditoría.
	var koboMetadataJSON []byte
	if len(fu.KoboMetadata) > 0 {
		var err error
		koboMetadataJSON, err = json.Marshal(fu.KoboMetadata)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "kobo_metadata inválido: " + err.Error()})
			return
		}
	}

	// Crear todo en una transacción: form_submission → follow_up_v2 → answers
	var followUpID, formSubmissionID string
	txErr := c.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Raw(`
			INSERT INTO salvia.form_submission (id, form_id, created_at, updated_at)
			VALUES (gen_random_uuid(), ?, NOW(), NOW())
			RETURNING id
		`, fu.FormID).Scan(&formSubmissionID).Error; err != nil {
			return fmt.Errorf("crear form_submission: %w", err)
		}

		var koboMetadataParam interface{}
		if koboMetadataJSON != nil {
			koboMetadataParam = string(koboMetadataJSON)
		}

		if err := tx.Raw(`
			INSERT INTO salvia.follow_up_v2 (
				id, case_id, form_submission_id, form_id, agent_id, status, team, risk_status,
				scheduled_date, scheduled_time, completed_at, sequence_number, summary,
				kobo_metadata, attempts, is_priority, created_at, updated_at
			) VALUES (
				gen_random_uuid(), ?, ?, ?, ?, ?, ?, ?,
				?, ?, ?, ?, ?,
				?::jsonb, 0, false, NOW(), NOW()
			) RETURNING id
		`,
			fu.CaseID, formSubmissionID, fu.FormID, fu.AgentID, status, fu.Team, fu.RiskStatus,
			scheduledDate, fu.ScheduledTime, completedAt, fu.SequenceNumber, fu.Summary,
			koboMetadataParam,
		).Scan(&followUpID).Error; err != nil {
			return fmt.Errorf("crear follow_up_v2: %w", err)
		}

		// Un solo INSERT multi-fila para todas las respuestas en vez de uno por
		// respuesta — con ~7-8 respuestas típicas por llamada, esto reemplaza
		// 7-8 round-trips secuenciales al pooler remoto por 1 solo.
		if len(answers) > 0 {
			valuePlaceholders := make([]string, 0, len(answers))
			args := make([]interface{}, 0, len(answers)*3)
			for _, a := range answers {
				valuePlaceholders = append(valuePlaceholders, "(gen_random_uuid(), ?, ?, ?, NOW(), NOW())")
				args = append(args, formSubmissionID, a.QuestionID, a.Value)
			}
			insertSQL := fmt.Sprintf(`
				INSERT INTO salvia.answer (id, form_submission_id, question_id, value, created_at, updated_at)
				VALUES %s
			`, strings.Join(valuePlaceholders, ", "))
			if err := tx.Exec(insertSQL, args...).Error; err != nil {
				return fmt.Errorf("crear answers (batch de %d): %w", len(answers), err)
			}
		}

		return nil
	})

	if txErr != nil {
		// Carrera rarísima: dos requests con el mismo external_ref llegaron casi
		// simultáneas y ambas pasaron el chequeo de arriba antes de que la otra
		// terminara de insertar. El índice único de (form_id, external_ref) la
		// atrapa acá — se resuelve igual que el caso normal (se busca el que sí
		// se creó y se devuelve como "ya existía") en vez de fallar la petición.
		if externalRef != "" && esErrorDeUnicidad(txErr) {
			if existente, err := c.buscarFollowUpPorExternalRef(fu.FormID, externalRef); err == nil && existente != nil {
				log.Printf("[MIGRATE-FOLLOWUP] carrera detectada en external_ref=%s — devolviendo el ya creado (followUpId=%s)", externalRef, existente.ID)
				ctx.JSON(http.StatusOK, gin.H{
					"success":          true,
					"followUpId":       existente.ID,
					"formSubmissionId": existente.FormSubmissionID,
					"alreadyExisted":   true,
				})
				return
			}
		}
		log.Printf("[MIGRATE-FOLLOWUP] error: %v", txErr)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": txErr.Error()})
		return
	}

	log.Printf("[MIGRATE-FOLLOWUP] followUpId=%s formSubmissionId=%s answers=%d warnings=%d",
		followUpID, formSubmissionID, len(answers), len(warnings))

	ctx.JSON(http.StatusOK, gin.H{
		"success":          true,
		"followUpId":       followUpID,
		"formSubmissionId": formSubmissionID,
		"answersCreated":   len(answers),
		"warnings":         warnings,
	})
}

// CreateForm crea un form con secciones y preguntas de solo tipo texto — pensado
// para crear en producción, de forma repetible y sin necesitar acceso directo a
// la BD, los forms que usan los scripts de migración (Kobo, Bases Salvia, y
// cualquier migración futura con el mismo patrón "1 fila → N respuestas de
// texto"). El endpoint /admin/migrate/follow-up ya resuelve sectionN.questionM
// por `order`, así que basta con que las secciones y preguntas queden creadas en
// el mismo orden en que vienen en el body.
//
// POST /api/v1/admin/migrate/create-form?name=...&description=...&status=inactive
// Header: X-Security-Key: SALVIA_MIGRATE_2026_PROD
// Query params:
//   - name (requerido): nombre del form.
//   - description (opcional): descripción del form.
//   - status (opcional, default "inactive"): forms de migración no deberían
//     quedar seleccionables para diligenciamiento en vivo hasta confirmarse.
//
// Body JSON: un array de secciones, en el orden en que deben quedar (order
// 1-indexado por posición) — NO un objeto envolvente:
//
//	[
//	  { "sectionName": "Nombre Sección 1", "questions": ["Pregunta 1", "Pregunta 2"] },
//	  { "sectionName": "Nombre Sección 2", "questions": ["Pregunta 1", "Pregunta 2"] }
//	]
//
// Las preguntas dentro de cada sección también quedan con order 1-indexado por
// posición. Todas se crean con question_type="text" y required=false.
func (c *MigrateController) CreateForm(ctx *gin.Context) {
	name := strings.TrimSpace(ctx.Query("name"))
	if name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "falta el query param 'name' (nombre del form)"})
		return
	}
	description := ctx.Query("description")
	status := ctx.Query("status")
	if status == "" {
		status = "inactive"
	}

	type sectionInput struct {
		SectionName string   `json:"sectionName" binding:"required"`
		Questions   []string `json:"questions"`
	}
	var body []sectionInput
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(body) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "el body debe ser un array con al menos una sección"})
		return
	}

	type questionOut struct {
		ID          string `json:"id"`
		Description string `json:"description"`
		Order       int    `json:"order"`
	}
	type sectionOut struct {
		ID          string        `json:"id"`
		SectionName string        `json:"sectionName"`
		Order       int           `json:"order"`
		Questions   []questionOut `json:"questions"`
	}

	form := models.Form{Name: name, Description: description, Status: status}
	var sectionsOut []sectionOut

	txErr := c.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&form).Error; err != nil {
			return fmt.Errorf("crear form: %w", err)
		}

		for i, sec := range body {
			section := models.FormSection{
				FormID: form.ID,
				Name:   sec.SectionName,
				Order:  i + 1,
			}
			if err := tx.Create(&section).Error; err != nil {
				return fmt.Errorf("crear form_section %q: %w", sec.SectionName, err)
			}

			secOut := sectionOut{ID: section.ID, SectionName: sec.SectionName, Order: section.Order}
			for j, qText := range sec.Questions {
				question := models.Question{
					FormID:         form.ID,
					FormSectionID:  section.ID,
					QuestionTypeID: "text",
					Description:    qText,
					Required:       false,
					Order:          j + 1,
				}
				if err := tx.Create(&question).Error; err != nil {
					return fmt.Errorf("crear question %q (sección %q): %w", qText, sec.SectionName, err)
				}
				secOut.Questions = append(secOut.Questions, questionOut{ID: question.ID, Description: qText, Order: question.Order})
			}
			sectionsOut = append(sectionsOut, secOut)
		}
		return nil
	})

	if txErr != nil {
		log.Printf("[CREATE-FORM] error: %v", txErr)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": txErr.Error()})
		return
	}

	log.Printf("[CREATE-FORM] formId=%s name=%q secciones=%d", form.ID, name, len(body))

	ctx.JSON(http.StatusOK, gin.H{
		"success":  true,
		"formId":   form.ID,
		"name":     name,
		"status":   status,
		"sections": sectionsOut,
	})
}

// parseFlexibleDateTime intenta parsear varios formatos comunes de fecha/fecha-hora.
// Retorna time.Time{} (zero value) si ninguno matchea.
func parseFlexibleDateTime(s string) time.Time {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}
	}
	formatos := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}
	for _, layout := range formatos {
		if t, err := time.Parse(layout, s); err == nil {
			return t
		}
	}
	return time.Time{}
}
