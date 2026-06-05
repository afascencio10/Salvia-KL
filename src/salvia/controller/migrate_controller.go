// Package controller — migrate_controller.go
// Endpoints administrativos para migración de datos legacy.
// Protegidos por X-Security-Key en el header.
// Se ejecutan una sola vez desde Postman/curl durante el despliegue.
package controller

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

const migrateSecurityKey = "SALVIA_MIGRATE_2026_PROD"

type MigrateController struct {
	db *gorm.DB
}

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
		"success":       true,
		"message":       fmt.Sprintf("Caso %s: %d seguimientos creados de %d solicitados", caso.VictimCaseICode, creados, len(body.Seguimientos)),
		"caso":          caso.VictimCaseICode,
		"agente":        agentID,
		"team":          team,
		"riesgo":        riskStatus,
		"creados":       creados,
		"total":         len(body.Seguimientos),
		"seguimientos":  seguimientosCreados,
	})
}

// MigrateFromExcel procesa un archivo Excel con la hoja "Riesgo nuevo".
// POST /api/v1/admin/migrate/excel
// Header: X-Security-Key: SALVIA_MIGRATE_2026_PROD
// Body: multipart/form-data con campo "file" (archivo .xlsx)
//
// Estructura del Excel (hoja "Riesgo nuevo"):
//   Columnas de fecha: "Fecha primer seguimiento", "Fecha segundo seguimiento", etc.
//   Columna "Alerta pendientes": indica cuál seguimiento es el próximo PENDIENTE (ej: "Segundo seguimiento")
//   Columna "Caso cerrado": si dice "si" se omite la fila
//   Columna "Documento": cédula de la víctima
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
		"hoja":         sheetName,
		"total_filas":  len(rows) - 1,
		"encabezados":  rows[0],
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
