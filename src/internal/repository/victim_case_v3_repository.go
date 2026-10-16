package repository

import (
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"
)

// VictimCaseV3Repository -- Registro V3 (DocsMD/Screens/Registro V3).
// Implementado SIN dinamic-form (decisión A1): no hay form_submission ni
// BuildAnswersByFieldKey genérico -- el frontend manda directamente el mapa
// de respuestas de cada módulo guardado. Reutiliza el mismo mecanismo de
// "nunca posponer, usar valores dummy" de victim_case_form_repository.go
// (V2) -- ver registro-v3-data.md § D8/D9.
type VictimCaseV3Repository interface {
	// FindByDocumentOrPhone busca Registros/Casos existentes -- D7: solo por
	// documento y/o teléfono, nunca por nombre.
	FindByDocumentOrPhone(ctx context.Context, docType, docNumber, phone string) ([]VictimCaseV3Summary, error)

	// CreateDraft crea el general_user + victim_case ('bo') + victim_case_form2
	// (con TODAS sus columnas NOT NULL -- las 55 heredadas de V2 y las propias
	// de V3 -- rellenas con el dummy "No"/"Pendiente" correspondiente). Nunca
	// se pospone: la fila existe desde el primer momento del Cabezote.
	CreateDraft(ctx context.Context, input CreateDraftV3Input) (icode string, err error)

	// SaveSection actualiza SOLO las columnas del módulo indicado (CI/DP/CA/
	// HV/AG/TZ/GS) -- UPDATE parcial, el resto de columnas no se toca.
	// answersByKey indexado por FieldKeyV3(modulo, orden).
	SaveSection(ctx context.Context, icode, modulo string, answersByKey map[string]string) error

	// SaveTamizajeBranch actualiza la rama pareja o no-pareja completa (24
	// campos) más el resultado (tipo/nivel/score/tie_flag). riskTypeCode y
	// riskLevelCode son los 2-char code de victim_case_form2_enums
	// (categorías victim_case_form2_tzv3_risk_type/_risk_level).
	SaveTamizajeBranch(ctx context.Context, icode string, wasPartner bool, answersByKey map[string]string, riskTypeCode, riskLevelCode string, riskScore int, tieFlag bool) error

	// GetTamizajeResult relee el resultado de Tamizaje guardado por
	// SaveTamizajeBranch (tipo/nivel/score, como codes de 2-char). BUG
	// encontrado en QA (2026-10-15): GetCase nunca repoblaba esto, así que
	// el Cierre mostraba "Nivel de riesgo: —" salvo que Finalizar se
	// disparara en la misma sesión de navegador donde se guardó el
	// Tamizaje (tzResultado era puramente en memoria, nunca se releía del
	// backend al reabrir/recargar el Caso). riskLevelCode == "" si el
	// Tamizaje nunca se ha guardado para este Caso.
	GetTamizajeResult(ctx context.Context, icode string) (riskTypeCode, riskLevelCode string, riskScore int, err error)

	GetAnswers(ctx context.Context, icode string) (map[string]string, error)
	// GetAnsweredKeys -- ver markKeysAnswered. Criterio real de "el agente
	// llegó a responder esto", a diferencia de GetAnswers (valor guardado,
	// dummy o real, indistinguibles para bool_enum/enum1/date/time).
	GetAnsweredKeys(ctx context.Context, icode string) (map[string]bool, error)
	GetCaseSummary(ctx context.Context, icode string) (*VictimCaseV3Summary, error)
	UpdateStatus(ctx context.Context, icode, status string) error
	// UpdateStatusIfCurrent transiciona solo si el estado actual es
	// expectedCurrent -- evita que un recálculo tardío (ej. triadaCompleta
	// reevaluada tras eliminar una respuesta y volver a completarla) revierta
	// un Caso que el agente ya hizo avanzar a Terminado/Cerrado por otra vía.
	UpdateStatusIfCurrent(ctx context.Context, icode, expectedCurrent, newStatus string) error
	MostRecentCaseByDocument(ctx context.Context, docType, docNumber string) (icode string, err error)

	StoreCredentials(ctx context.Context, icode, login, password string) error

	// ListEnumCategories expone el catálogo de opciones (code/icode/label)
	// de todas las categorías que usa Registro V3 -- ver
	// victimCaseV3EnumCategories. Sin dinamic-form (A1) el frontend no tiene
	// otra forma de conocer los icode que debe enviar en SaveSection.
	ListEnumCategories(ctx context.Context) (map[string][]EnumOptionV3, error)
}

type EnumOptionV3 struct {
	Code  string `json:"code"`
	ICode string `json:"icode"`
	Label string `json:"label"`
}

type VictimCaseV3Summary struct {
	ICode     string `json:"icode"`
	Names     string `json:"names"`
	LastNames string `json:"lastNames"`
	DocType   string `json:"docType"`
	DocNumber string `json:"docNumber"`
	Phone     string `json:"phone"`
	Status    string `json:"status"`
	// RegistroCompleto -- D1/D9: calculado por el service (victimCaseV3Service.Search)
	// sobre el Caso más reciente con este documento, no por fila. Todas las
	// filas que comparten documento llevan el mismo valor.
	RegistroCompleto bool `json:"registroCompleto"`
}

type CreateDraftV3Input struct {
	DocType            string
	DocNumber          string
	Phone              string
	ResidenceTownCode  string
	AgentID            string
	Channel            string
	GeneralUserICode   string
	GeneralUserLogin   string
	GeneralUserPassSha string
}

// VictimCaseEntityInput -- "Agregar entidad" (GS § C). Reconstruido sobre el
// modelo real de entity_case/entity_case_obligation (el mismo que usan la
// pestaña "Gestión institucional" y el form de Seguimiento -- ver
// registro-v3-pendientes.md § B): el agente filtra por Departamento/Ciudad/
// Municipio + Sector para encontrar la sede real (entity_branch) en vez de
// escribir el nombre a mano; EntityBranchID es la selección resultante.
type VictimCaseEntityInput struct {
	ICode             string // id real de entity_case (UUID), vacío al crear
	VictimCaseICode   string
	Kind              string // 'acudida' | 'activacion'
	EntityBranchID    int64
	FileNumber        string   // -> entity_case.Radicado
	ActivationChannel string   // code(s) CSV, solo informativo si Kind='activacion' -- ver GAP B-canal en registro-v3-pendientes.md
	CompletedActions  []string // ids de entity_obligation -- solo aplica a Kind='acudida'
	PendingActions    []string // ids de entity_obligation (o "otra")
	OtherInfo         string   // -> entity_case.Objetivo, solo si Kind='acudida'
	Agent             string
}

// VictimCaseBarrierInput -- "Agregar barrera" (GS § D). Reconstruido sobre el
// modelo real salvia.barrier_v2 (el mismo que usan Seguimiento y la sesión
// psicosocial). EntityBranchID se resuelve con la misma búsqueda por
// Departamento/Ciudad/Municipio + Sector que usa Entidades. Department/City/
// Town se guardan además en la barrera (BarrierV2.DepartmentID/CityID/
// TownID) -- es la ubicación donde se identificó la barrera, no solo un
// filtro de búsqueda.
type VictimCaseBarrierInput struct {
	ICode           string
	VictimCaseICode string
	SectorCode      string // 'salud'|'justicia'|'proteccion'|'otras_instituciones' -- vocabulario real de entity_branch/barrier_v2, NO el enum victim_case_form2_sector
	Department      string
	City            string
	Town            string
	EntityBranchID  int64
	Specific        []string
	OtherDetail     string
	EntityName      string // solo libre cuando SectorCode='otras_instituciones' -- BarrierV2.InstitutionName
	Institutional   []string
	Economic        []string
	Territorial     []string
	Differential    []string
	Date            string
	Official        string
	Description     string
	Management      []string
	Competences     []string // códigos del catálogo salvia_config.CompetenceOptions
	Agent           string
}

type victimCaseV3Repository struct {
	db *gorm.DB
}

func NewVictimCaseV3Repository(db *gorm.DB) VictimCaseV3Repository {
	return &victimCaseV3Repository{db: db}
}

func (r *victimCaseV3Repository) FindByDocumentOrPhone(ctx context.Context, docType, docNumber, phone string) ([]VictimCaseV3Summary, error) {
	if docNumber == "" && phone == "" {
		return nil, fmt.Errorf("findByDocumentOrPhone: se requiere documento o teléfono")
	}
	query := r.db.WithContext(ctx).
		Table("salvia.victim_case vc").
		Joins("LEFT JOIN salvia.victim_case_form2 f2 ON f2.victim_case_form2_victim_case = vc.victim_case_id").
		Select(`vc.victim_case_i_code AS i_code, vc.victim_case_victim_names AS names,
			vc.victim_case_victim_last_names AS last_names, vc.victim_case_victim_doc_type AS doc_type,
			vc.victim_case_victim_doc_number AS doc_number, f2.victim_case_form2_victim_phone AS phone,
			vc.victim_case_status AS status`)

	if docNumber != "" {
		query = query.Where("vc.victim_case_victim_doc_number = ?", docNumber)
		if docType != "" {
			query = query.Where("vc.victim_case_victim_doc_type = ?", docType)
		}
	} else {
		query = query.Where("f2.victim_case_form2_victim_phone = ?", phone)
	}

	var rows []struct {
		ICode     string
		Names     string
		LastNames string
		DocType   string
		DocNumber string
		Phone     *int64
		Status    string
	}
	if err := query.Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("findByDocumentOrPhone: %w", err)
	}

	out := make([]VictimCaseV3Summary, 0, len(rows))
	for _, row := range rows {
		phone := ""
		if row.Phone != nil {
			phone = fmt.Sprintf("%d", *row.Phone)
		}
		out = append(out, VictimCaseV3Summary{
			ICode: row.ICode, Names: row.Names, LastNames: row.LastNames,
			DocType: row.DocType, DocNumber: row.DocNumber, Phone: phone, Status: row.Status,
		})
	}
	return out, nil
}

func (r *victimCaseV3Repository) CreateDraft(ctx context.Context, input CreateDraftV3Input) (string, error) {
	var newICode string
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1-3. general_user_profile + general_user + rol "us" -- idéntico a V2
		// (victim_case_general_user es NOT NULL en el esquema real, no es
		// opcional -- confirma D11).
		var profileID int64
		if err := tx.Raw(`
			INSERT INTO security.general_user_profile
				(general_user_profile_i_code, general_user_profile_creation_date, general_user_profile_update_date,
				 general_user_profile_names, general_user_profile_last_names,
				 general_user_profile_doc_type, general_user_profile_doc_number, general_user_profile_town)
			VALUES (gen_random_uuid()::text, now(), now(), ?, ?, ?, ?, ?)
			RETURNING general_user_profile_id
		`, dummyNamesV3, dummyLastNamesV3, coalesceV3(input.DocType, dummyDocTypeV3), coalesceV3(input.DocNumber, dummyDocNumberV3), coalesceV3(input.ResidenceTownCode, dummyTownCodeV3)).
			Scan(&profileID).Error; err != nil {
			return fmt.Errorf("crear general_user_profile: %w", err)
		}

		var generalUserID int64
		if err := tx.Raw(`
			INSERT INTO security.general_user
				(general_user_i_code, general_user_creation_date, general_user_update_date,
				 general_user_login, general_user_password, general_user_status, general_user_language,
				 general_user_general_user_profile)
			VALUES (?, now(), now(), ?, ?, 'e', 'sp', ?)
			RETURNING general_user_id
		`, input.GeneralUserICode, input.GeneralUserLogin, input.GeneralUserPassSha, profileID).
			Scan(&generalUserID).Error; err != nil {
			return fmt.Errorf("crear general_user: %w", err)
		}

		var roleID int64
		if err := tx.Raw(`SELECT role_id FROM security.role WHERE role_code = 'us' LIMIT 1`).Scan(&roleID).Error; err != nil {
			return fmt.Errorf("buscar rol 'us': %w", err)
		}
		if err := tx.Exec(`
			INSERT INTO security.rel_role_general_user (role_id, general_user_id, rel_role_general_user_creation_date)
			VALUES (?, ?, now())
		`, roleID, generalUserID).Error; err != nil {
			return fmt.Errorf("asignar rol: %w", err)
		}

		// 4. victim_case en Borrador
		newICode = utils.GetUUID()
		if err := tx.Exec(`
			INSERT INTO salvia.victim_case
				(victim_case_i_code, victim_case_creation_date, victim_case_update_date, victim_case_status,
				 victim_case_victim_names, victim_case_victim_last_names, victim_case_victim_doc_type,
				 victim_case_victim_doc_number, victim_case_victim_town_code, victim_case_general_user,
				 victim_case_team, agent_id)
			VALUES (?, now(), now(), 'bo', ?, ?, ?, ?, ?, ?, ?, ?)
		`, newICode, dummyNamesV3, dummyLastNamesV3, coalesceV3(input.DocType, dummyDocTypeV3), coalesceV3(input.DocNumber, dummyDocNumberV3),
			coalesceV3(input.ResidenceTownCode, dummyTownCodeV3), input.GeneralUserICode, "", input.AgentID).Error; err != nil {
			return fmt.Errorf("crear victim_case: %w", err)
		}

		// 5. victim_case_form2 -- fila completa desde el inicio, TODAS las
		// columnas NOT NULL (las 55 heredadas de V2 + las propias de V3 que se
		// marquen como tal) con dummy. Se resuelve aparte porque necesita el
		// victim_case_id numérico recién creado.
		var victimCaseID int64
		if err := tx.Raw(`SELECT victim_case_id FROM salvia.victim_case WHERE victim_case_i_code = ?`, newICode).Scan(&victimCaseID).Error; err != nil {
			return fmt.Errorf("resolver victim_case_id: %w", err)
		}
		if err := insertDummyForm2Row(tx, victimCaseID); err != nil {
			return fmt.Errorf("crear victim_case_form2 con dummies: %w", err)
		}

		// BUG encontrado en QA (2026-10-14): el teléfono del criterio de
		// búsqueda del Cabezote se perdía -- insertDummyForm2Row no conocía
		// el input original y lo dejaba en 0 (dummy). E06 exige "Borrador
		// prediligenciado con los criterios de búsqueda" (D7).
		if input.Phone != "" {
			if err := tx.Exec(`UPDATE salvia.victim_case_form2 SET victim_case_form2_victim_phone = ? WHERE victim_case_form2_victim_case = ?`, input.Phone, victimCaseID).Error; err != nil {
				return fmt.Errorf("prediligenciar teléfono: %w", err)
			}
			if err := markKeysAnswered(tx, newICode, map[string]string{FieldKeyV3("CI", 4): input.Phone}); err != nil {
				return fmt.Errorf("marcar teléfono respondido: %w", err)
			}
		}

		return nil
	})
	return newICode, err
}

const (
	dummyNamesV3     = "Pendiente"
	dummyLastNamesV3 = "Pendiente"
	dummyDocTypeV3   = "cc"
	dummyDocNumberV3 = "0000000000"
	dummyTownCodeV3  = "00000000"
)

func coalesceV3(v, fallback string) string {
	if v == "" {
		return fallback
	}
	return v
}

// insertDummyForm2Row crea la fila de victim_case_form2 con TODAS sus
// columnas NOT NULL resueltas a un valor dummy neutro según su Kind (D8:
// Text→"", Int/Scale→0, BoolEnum/Enum1→id de "No", Date/Time/Timestamp→
// placeholder) -- tanto las heredadas de V2 (victimCaseFormFields) como las
// propias de V3 (victimCaseV3FormFields). Las columnas nullable se dejan
// NULL.
//
// BUG encontrado en QA (2026-10-14): la versión anterior solo rellenaba
// BoolEnum/Enum1, dejando sin valor columnas NOT NULL reales como
// victim_case_form2_birth_date, _facts_date, _facts_description,
// _residence_town/_address, _saliva_management_explanation -- el INSERT
// fallaba por violación NOT NULL en el primer CreateDraft. Se agrega además
// una red de seguridad que audita en caliente contra information_schema
// cualquier columna NOT NULL que ningún catálogo Go declare (ej. columnas
// propias del motor dinamic-form que V3 no usa pero comparte tabla).
func insertDummyForm2Row(tx *gorm.DB, victimCaseID int64) error {
	var noID int64
	if err := tx.Raw(`SELECT victim_case_form2_enums_id FROM salvia.victim_case_form2_enums WHERE victim_case_form2_enums_category = 'yes_no' AND victim_case_form2_enums_code = 'n' LIMIT 1`).Scan(&noID).Error; err != nil {
		return fmt.Errorf("resolver dummy 'No': %w", err)
	}

	cols := map[string]interface{}{
		"victim_case_form2_i_code":        utils.GetUUID(),
		"victim_case_form2_victim_case":   victimCaseID,
		"victim_case_form2_creation_date": time.Now(),
		"victim_case_form2_update_date":   time.Now(),
	}

	fillDummyByKind := func(f VictimCaseFormField) {
		if f.Column == "" {
			return // KindEnumN -- sin columna propia, va a rel_* aparte
		}
		if _, already := cols[f.Column]; already {
			return
		}
		switch f.Kind {
		case KindBoolEnum, KindEnum1:
			cols[f.Column] = noID
		case KindText:
			cols[f.Column] = ""
		case KindInt, KindScale:
			cols[f.Column] = 0
		case KindDate:
			cols[f.Column] = "1900-01-01"
		case KindTime:
			cols[f.Column] = "00:00:00"
		case KindTimestamp:
			cols[f.Column] = time.Now()
		}
	}
	for _, f := range victimCaseFormFields {
		fillDummyByKind(f)
	}
	for _, f := range victimCaseV3FormFields {
		fillDummyByKind(f)
	}

	// Red de seguridad: columnas NOT NULL SIN DEFAULT que ningún catálogo Go
	// cubra todavía -- se rellenan con un valor neutro según su tipo real.
	// column_default IS NULL excluye columnas con secuencia/identity propia
	// (ej. victim_case_form2_id) -- forzarlas rompería el autoincremento.
	type notNullCol struct {
		Name     string `gorm:"column:column_name"`
		DataType string `gorm:"column:data_type"`
	}
	var missing []notNullCol
	if err := tx.Raw(`
		SELECT column_name, data_type FROM information_schema.columns
		WHERE table_schema = 'salvia' AND table_name = 'victim_case_form2'
		  AND is_nullable = 'NO' AND column_default IS NULL
	`).Scan(&missing).Error; err != nil {
		return fmt.Errorf("auditar columnas NOT NULL de victim_case_form2: %w", err)
	}
	for _, c := range missing {
		if _, already := cols[c.Name]; already {
			continue
		}
		switch {
		case strings.Contains(c.DataType, "timestamp"):
			cols[c.Name] = time.Now()
		case c.DataType == "date":
			cols[c.Name] = "1900-01-01"
		case strings.Contains(c.DataType, "time"):
			cols[c.Name] = "00:00:00"
		case c.DataType == "bigint" || c.DataType == "integer" || c.DataType == "smallint" || c.DataType == "numeric":
			cols[c.Name] = 0
		case c.DataType == "boolean":
			cols[c.Name] = false
		default: // character varying, text, etc.
			cols[c.Name] = ""
		}
	}

	return tx.Table("salvia.victim_case_form2").Create(cols).Error
}

func (r *victimCaseV3Repository) SaveSection(ctx context.Context, icode, modulo string, answersByKey map[string]string) error {
	var victimCaseID int64
	if err := r.db.WithContext(ctx).Raw(`SELECT victim_case_id FROM salvia.victim_case WHERE victim_case_i_code = ?`, icode).Scan(&victimCaseID).Error; err != nil {
		return fmt.Errorf("saveSection: resolver victim_case_id: %w", err)
	}
	if victimCaseID == 0 {
		return fmt.Errorf("saveSection: no existe victim_case con icode=%s", icode)
	}

	if err := markKeysAnswered(r.db.WithContext(ctx), icode, answersByKey); err != nil {
		return fmt.Errorf("saveSection[%s]: %w", modulo, err)
	}

	updates, err := resolveFieldUpdates(r.db.WithContext(ctx), victimCaseV3FormFields, answersByKey)
	if err != nil {
		return fmt.Errorf("saveSection[%s]: %w", modulo, err)
	}
	if len(updates) > 0 {
		updates["victim_case_form2_update_date"] = time.Now()
		if err := r.db.WithContext(ctx).Table("salvia.victim_case_form2").
			Where("victim_case_form2_victim_case = ?", victimCaseID).
			Updates(updates).Error; err != nil {
			return fmt.Errorf("saveSection[%s]: actualizar form2: %w", modulo, err)
		}
	}

	var form2ID int64
	if err := r.db.WithContext(ctx).Raw(`SELECT victim_case_form2_id FROM salvia.victim_case_form2 WHERE victim_case_form2_victim_case = ?`, victimCaseID).Scan(&form2ID).Error; err != nil {
		return fmt.Errorf("saveSection[%s]: resolver form2_id: %w", modulo, err)
	}
	if err := saveEnumNFields(r.db.WithContext(ctx), form2ID, answersByKey); err != nil {
		return fmt.Errorf("saveSection[%s]: %w", modulo, err)
	}

	// Sincroniza los campos núcleo de victim_case si vinieron en este módulo
	// (CI trae nombres/documento; DP trae municipio de residencia).
	coreUpdates := map[string]interface{}{}
	if v, ok := answersByKey[FieldKeyV3("CI", 2)]; ok {
		parts := splitNames(v)
		coreUpdates["victim_case_victim_names"] = parts[0]
		coreUpdates["victim_case_victim_last_names"] = parts[1]
	}
	if v, ok := answersByKey[FieldKeyV3("CI", 5)]; ok && v != "" {
		// BUG encontrado en QA (2026-10-14): CI_5 es KindEnum1 -- el frontend
		// manda el icode (uuid), pero victim_case_victim_doc_type es
		// varchar(2) y espera el code corto (ej. "cc"). Resolver antes de
		// escribir (antes reventaba con "value too long for type character
		// varying(2)").
		code, err := resolveEnumCodeByICode(r.db.WithContext(ctx), v)
		if err != nil {
			return fmt.Errorf("saveSection[%s]: resolver code de CI_5: %w", modulo, err)
		}
		if code != "" {
			coreUpdates["victim_case_victim_doc_type"] = code
		}
	}
	if v, ok := answersByKey[FieldKeyV3("CI", 6)]; ok {
		coreUpdates["victim_case_victim_doc_number"] = v
	}
	if v, ok := answersByKey[FieldKeyV3("DP", 28)]; ok {
		coreUpdates["victim_case_victim_town_code"] = v
	}
	if len(coreUpdates) > 0 {
		coreUpdates["victim_case_update_date"] = time.Now()
		if err := r.db.WithContext(ctx).Table("salvia.victim_case").
			Where("victim_case_i_code = ?", icode).Updates(coreUpdates).Error; err != nil {
			return fmt.Errorf("saveSection[%s]: sincronizar victim_case: %w", modulo, err)
		}
	}

	// D8 (flow-E11, "fan-out al editar el Perfil"): si el módulo editado es
	// parte del Perfil (CI/DP/CA, D1), replica el MISMO cambio en los demás
	// victim_case que comparten Tipo+Número de documento. Best-effort: un
	// error acá se loggea pero no revierte ni falla el guardado de ESTE
	// Caso (ya quedó guardado correctamente arriba) -- mismo criterio que
	// ya usa este archivo para Entidades/Barreras y que V2-real usa para su
	// propio OnSectionUpdate. Sin control de concurrencia entre agentes
	// editando el mismo Perfil a la vez (GAP G-V3-01 ya documentado).
	if modulo == "CI" || modulo == "DP" || modulo == "CA" {
		if err := r.fanOutPerfilV3(ctx, icode, answersByKey); err != nil {
			log.Printf("[SaveSection] WARN: fan-out de Perfil falló para icode=%s módulo=%s: %v", icode, modulo, err)
		}
	}
	return nil
}

// fanOutPerfilV3 -- ver nota D8 en SaveSection. Resuelve el documento del
// Caso actual, busca los demás victim_case con el mismo Tipo+Número, y les
// aplica el mismo UPDATE de columnas escalares + relaciones EnumN que ya se
// le aplicó al Caso actual en esta misma llamada a SaveSection.
func (r *victimCaseV3Repository) fanOutPerfilV3(ctx context.Context, icode string, answersByKey map[string]string) error {
	var doc struct {
		DocType   string
		DocNumber string
	}
	if err := r.db.WithContext(ctx).Table("salvia.victim_case").
		Select("victim_case_victim_doc_type AS doc_type, victim_case_victim_doc_number AS doc_number").
		Where("victim_case_i_code = ?", icode).Take(&doc).Error; err != nil {
		return fmt.Errorf("fanOutPerfilV3: resolver documento: %w", err)
	}
	if doc.DocType == "" || doc.DocNumber == "" {
		return nil // todavía no hay documento real (dummy D8) -- nada que sincronizar
	}

	var siblings []struct {
		VictimCaseID int64
		ICode        string
	}
	if err := r.db.WithContext(ctx).Table("salvia.victim_case").
		Select("victim_case_id AS victim_case_id, victim_case_i_code AS i_code").
		Where("victim_case_victim_doc_type = ? AND victim_case_victim_doc_number = ? AND victim_case_i_code != ?", doc.DocType, doc.DocNumber, icode).
		Find(&siblings).Error; err != nil {
		return fmt.Errorf("fanOutPerfilV3: buscar hermanos: %w", err)
	}
	if len(siblings) == 0 {
		return nil
	}

	for _, sib := range siblings {
		updates, err := resolveFieldUpdates(r.db.WithContext(ctx), victimCaseV3FormFields, answersByKey)
		if err != nil {
			return fmt.Errorf("fanOutPerfilV3: resolver updates para %s: %w", sib.ICode, err)
		}
		if len(updates) > 0 {
			updates["victim_case_form2_update_date"] = time.Now()
			if err := r.db.WithContext(ctx).Table("salvia.victim_case_form2").
				Where("victim_case_form2_victim_case = ?", sib.VictimCaseID).Updates(updates).Error; err != nil {
				return fmt.Errorf("fanOutPerfilV3: actualizar form2 de %s: %w", sib.ICode, err)
			}
		}
		var sibForm2ID int64
		if err := r.db.WithContext(ctx).Raw(`SELECT victim_case_form2_id FROM salvia.victim_case_form2 WHERE victim_case_form2_victim_case = ?`, sib.VictimCaseID).Scan(&sibForm2ID).Error; err != nil {
			return fmt.Errorf("fanOutPerfilV3: resolver form2_id de %s: %w", sib.ICode, err)
		}
		if err := saveEnumNFields(r.db.WithContext(ctx), sibForm2ID, answersByKey); err != nil {
			return fmt.Errorf("fanOutPerfilV3: enumN de %s: %w", sib.ICode, err)
		}
		if err := markKeysAnswered(r.db.WithContext(ctx), sib.ICode, answersByKey); err != nil {
			return fmt.Errorf("fanOutPerfilV3: marcar respondidas en %s: %w", sib.ICode, err)
		}
	}
	return nil
}

// markKeysAnswered acumula en victim_case_v3_answered_keys (CSV) cada
// FieldKey presente en answersByKey -- sin importar si el valor es vacío.
// Necesario porque los valores dummy (D8) son indistinguibles de una
// respuesta real una vez guardados en victim_case_form2 (ej. dummy "No" en
// bool_enum == respuesta real "No") -- la única forma confiable de saber
// "¿el agente llegó a responder esto?" es rastrear qué keys llegaron
// alguna vez en un guardado real, no inferirlo del valor almacenado. Ver
// SaveSection / GetAnsweredKeys.
func markKeysAnswered(tx *gorm.DB, icode string, answersByKey map[string]string) error {
	if len(answersByKey) == 0 {
		return nil
	}
	var existing string
	if err := tx.Raw(`SELECT COALESCE(victim_case_v3_answered_keys, '') FROM salvia.victim_case WHERE victim_case_i_code = ?`, icode).Scan(&existing).Error; err != nil {
		return fmt.Errorf("markKeysAnswered: leer existentes: %w", err)
	}
	set := map[string]bool{}
	for _, k := range splitCSVV3(existing) {
		set[k] = true
	}
	for k := range answersByKey {
		set[k] = true
	}
	merged := make([]string, 0, len(set))
	for k := range set {
		merged = append(merged, k)
	}
	return tx.Exec(`UPDATE salvia.victim_case SET victim_case_v3_answered_keys = ? WHERE victim_case_i_code = ?`, strings.Join(merged, ","), icode).Error
}

// GetAnsweredKeys devuelve el set de FieldKey que el agente llegó a
// enviar alguna vez (ver markKeysAnswered) -- es el criterio real de
// "respondido", a diferencia de GetAnswers (que refleja lo que hay
// guardado en BD, dummy o real, indistinguibles para varios Kind).
func (r *victimCaseV3Repository) GetAnsweredKeys(ctx context.Context, icode string) (map[string]bool, error) {
	var raw string
	if err := r.db.WithContext(ctx).Raw(`SELECT COALESCE(victim_case_v3_answered_keys, '') FROM salvia.victim_case WHERE victim_case_i_code = ?`, icode).Scan(&raw).Error; err != nil {
		return nil, fmt.Errorf("getAnsweredKeys: %w", err)
	}
	set := map[string]bool{}
	for _, k := range splitCSVV3(raw) {
		set[k] = true
	}
	return set, nil
}

// splitNames parte "Nombres Apellidos" en [nombres, apellidos] de forma
// simple (primer espacio) -- el formulario V3 pide nombres y apellidos
// como un solo campo de texto (secciones.md § CI Q2), victim_case los
// guarda separados por herencia del esquema de V2.
func splitNames(full string) [2]string {
	for i, c := range full {
		if c == ' ' {
			return [2]string{full[:i], full[i+1:]}
		}
	}
	return [2]string{full, ""}
}

// resolveFieldUpdates resuelve un subconjunto de victimCaseV3FormFields
// contra answersByKey, devolviendo solo las columnas que SÍ vinieron en
// este guardado (no toca el resto -- UPDATE parcial, no upsert completo).
func resolveFieldUpdates(tx *gorm.DB, fields []VictimCaseFormField, answersByKey map[string]string) (map[string]interface{}, error) {
	updates := map[string]interface{}{}
	for _, f := range fields {
		raw, ok := answersByKey[f.Key]
		if !ok || f.Column == "" {
			continue // enumN (sin columna propia) se maneja aparte por el llamador vía rel_* -- ver GAP si se requiere
		}
		switch f.Kind {
		case KindText:
			updates[f.Column] = raw
		case KindInt:
			var n int64
			fmt.Sscanf(raw, "%d", &n)
			updates[f.Column] = n
		case KindDate:
			updates[f.Column] = raw
		case KindTime:
			updates[f.Column] = raw
		case KindTimestamp:
			updates[f.Column] = raw
		case KindBoolEnum:
			code := "n"
			if raw == "true" {
				code = "y"
			}
			id, err := resolveEnumIDByCategory(tx, "yes_no", code)
			if err != nil {
				return nil, err
			}
			updates[f.Column] = id
		case KindEnum1:
			id, err := resolveEnumIDByICode(tx, raw)
			if err != nil {
				return nil, err
			}
			updates[f.Column] = id
		case KindScale:
			updates[f.Column] = raw
		}
	}
	return updates, nil
}

// saveEnumNFields reescribe las relaciones multi-valor (rel_victim_case_
// form2_enums_victim_case_form2) de los campos KindEnumN presentes en
// answersByKey (CSV de icode). Acotado por categoría -- NO se borra toda la
// fila, solo las categorías de los campos que vinieron en este guardado, así
// un módulo no pisa las selecciones multi-valor de otro (ver
// victimCaseV3EnumNCategories).
func saveEnumNFields(tx *gorm.DB, form2ID int64, answersByKey map[string]string) error {
	for key, categories := range victimCaseV3EnumNCategories {
		raw, ok := answersByKey[key]
		if !ok {
			continue
		}
		if err := tx.Exec(`
			DELETE FROM salvia.rel_victim_case_form2_enums_victim_case_form2 r
			USING salvia.victim_case_form2_enums e
			WHERE r.victim_case_form2_enums_id = e.victim_case_form2_enums_id
			  AND r.victim_case_form2_id = ?
			  AND e.victim_case_form2_enums_category = ANY(string_to_array(?, ','))
		`, form2ID, strings.Join(categories, ",")).Error; err != nil {
			return fmt.Errorf("limpiar relaciones multi-valor (%s): %w", key, err)
		}
		if raw == "" {
			continue
		}
		for _, icode := range splitCSVV3(raw) {
			enumID, err := resolveEnumIDByICode(tx, icode)
			if err != nil || enumID == 0 {
				continue // icode inválido/huérfano, igual criterio que V2
			}
			if err := tx.Exec(`
				INSERT INTO salvia.rel_victim_case_form2_enums_victim_case_form2
					(victim_case_form2_enums_id, victim_case_form2_id)
				VALUES (?, ?)
			`, enumID, form2ID).Error; err != nil {
				return fmt.Errorf("crear relación multi-valor (%s): %w", key, err)
			}
		}
	}
	return nil
}

func splitCSVV3(s string) []string {
	var out []string
	start := 0
	for i, c := range s {
		if c == ',' {
			if v := s[start:i]; v != "" {
				out = append(out, v)
			}
			start = i + 1
		}
	}
	if v := s[start:]; v != "" {
		out = append(out, v)
	}
	return out
}

func resolveEnumIDByCategory(tx *gorm.DB, category, code string) (int64, error) {
	var id int64
	err := tx.Raw(`SELECT victim_case_form2_enums_id FROM salvia.victim_case_form2_enums WHERE victim_case_form2_enums_category = ? AND victim_case_form2_enums_code = ? LIMIT 1`, category, code).Scan(&id).Error
	return id, err
}

func resolveEnumIDByICode(tx *gorm.DB, icode string) (int64, error) {
	var id int64
	err := tx.Raw(`SELECT victim_case_form2_enums_id FROM salvia.victim_case_form2_enums WHERE victim_case_form2_enums_i_code = ?`, icode).Scan(&id).Error
	return id, err
}

func (r *victimCaseV3Repository) SaveTamizajeBranch(ctx context.Context, icode string, wasPartner bool, answersByKey map[string]string, riskTypeCode, riskLevelCode string, riskScore int, tieFlag bool) error {
	var victimCaseID int64
	if err := r.db.WithContext(ctx).Raw(`SELECT victim_case_id FROM salvia.victim_case WHERE victim_case_i_code = ?`, icode).Scan(&victimCaseID).Error; err != nil {
		return fmt.Errorf("saveTamizajeBranch: %w", err)
	}

	// BUG encontrado en QA (2026-10-15): SaveTamizajeBranch es una ruta de
	// guardado separada de SaveSection (POST /tamizaje, no POST /seccion/TZ)
	// y nunca llamaba a markKeysAnswered -- TZ_1..8 y la rama pareja/no-
	// pareja completa nunca quedaban en victim_case_v3_answered_keys, así
	// que Finalize (que exige TZ_1/2/3/4/7 respondidas) rechazaba
	// cualquier Caso aunque el agente sí hubiera llenado el Tamizaje.
	// answersByKey aquí ya trae TZ_1..8 + las 24 de la rama (ver
	// SaveTamizaje en el service), aunque resolveFieldUpdates de abajo solo
	// use las de la rama activa para el UPDATE de columnas.
	if err := markKeysAnswered(r.db.WithContext(ctx), icode, answersByKey); err != nil {
		return fmt.Errorf("saveTamizajeBranch: %w", err)
	}

	fields := victimCaseV3TamizajeNoPareja
	if wasPartner {
		fields = victimCaseV3TamizajePareja
	}
	updates, err := resolveFieldUpdates(r.db.WithContext(ctx), fields, answersByKey)
	if err != nil {
		return fmt.Errorf("saveTamizajeBranch: %w", err)
	}

	noID, _ := resolveEnumIDByCategory(r.db.WithContext(ctx), "yes_no", "n")
	yesID, _ := resolveEnumIDByCategory(r.db.WithContext(ctx), "yes_no", "y")
	if tieFlag {
		updates["victim_case_form2_tzv3_risk_tie_flag"] = yesID
	} else {
		updates["victim_case_form2_tzv3_risk_tie_flag"] = noID
	}

	if riskTypeID, err := resolveEnumIDByCategory(r.db.WithContext(ctx), "victim_case_form2_tzv3_risk_type", riskTypeCode); err == nil && riskTypeID != 0 {
		updates["victim_case_form2_tzv3_risk_type"] = riskTypeID
	}
	if riskLevelID, err := resolveEnumIDByCategory(r.db.WithContext(ctx), "victim_case_form2_tzv3_risk_level", riskLevelCode); err == nil && riskLevelID != 0 {
		updates["victim_case_form2_tzv3_risk_level"] = riskLevelID
	}
	updates["victim_case_form2_tzv3_risk_score"] = riskScore
	updates["victim_case_form2_update_date"] = time.Now()

	return r.db.WithContext(ctx).Table("salvia.victim_case_form2").
		Where("victim_case_form2_victim_case = ?", victimCaseID).
		Updates(updates).Error
}

func (r *victimCaseV3Repository) GetTamizajeResult(ctx context.Context, icode string) (string, string, int, error) {
	row := struct {
		RiskType  sql.NullInt64
		RiskLevel sql.NullInt64
		RiskScore sql.NullInt64
	}{}
	err := r.db.WithContext(ctx).Table("salvia.victim_case_form2 f").
		Joins("JOIN salvia.victim_case vc ON vc.victim_case_id = f.victim_case_form2_victim_case").
		Where("vc.victim_case_i_code = ?", icode).
		Select("f.victim_case_form2_tzv3_risk_type AS risk_type, f.victim_case_form2_tzv3_risk_level AS risk_level, f.victim_case_form2_tzv3_risk_score AS risk_score").
		Take(&row).Error
	if err != nil {
		return "", "", 0, fmt.Errorf("getTamizajeResult: %w", err)
	}
	if !row.RiskLevel.Valid || row.RiskLevel.Int64 == 0 {
		return "", "", 0, nil
	}
	tipoCode, err := enumCodeByID(r.db.WithContext(ctx), row.RiskType.Int64)
	if err != nil {
		return "", "", 0, fmt.Errorf("getTamizajeResult: tipo: %w", err)
	}
	nivelCode, err := enumCodeByID(r.db.WithContext(ctx), row.RiskLevel.Int64)
	if err != nil {
		return "", "", 0, fmt.Errorf("getTamizajeResult: nivel: %w", err)
	}
	return tipoCode, nivelCode, int(row.RiskScore.Int64), nil
}

// GetAnswers relee victim_case_form2 (+ victim_case para los campos núcleo
// CI_2/CI_5/CI_6/DP_28) y lo resuelve a FieldKey->valor crudo, en el mismo
// formato que espera SaveSection/el frontend. Permite reabrir un Borrador y
// repoblar el formulario. Los campos KindEnumN se resuelven aparte (no están
// en victimCaseV3FormFields) consultando rel_victim_case_form2_enums_*.
func (r *victimCaseV3Repository) GetAnswers(ctx context.Context, icode string) (map[string]string, error) {
	var form2ID int64
	if err := r.db.WithContext(ctx).Raw(`
		SELECT f.victim_case_form2_id FROM salvia.victim_case_form2 f
		JOIN salvia.victim_case vc ON vc.victim_case_id = f.victim_case_form2_victim_case
		WHERE vc.victim_case_i_code = ?
	`, icode).Scan(&form2ID).Error; err != nil || form2ID == 0 {
		return nil, fmt.Errorf("getAnswers: no existe victim_case_form2 para icode=%s", icode)
	}

	row := map[string]interface{}{}
	if err := r.db.WithContext(ctx).Table("salvia.victim_case_form2").
		Where("victim_case_form2_id = ?", form2ID).Take(&row).Error; err != nil {
		return nil, fmt.Errorf("getAnswers: leer form2: %w", err)
	}

	answers := map[string]string{}
	allFields := append(append([]VictimCaseFormField{}, victimCaseV3FormFields...), victimCaseV3TamizajePareja...)
	allFields = append(allFields, victimCaseV3TamizajeNoPareja...)
	for _, f := range allFields {
		if f.Column == "" {
			continue
		}
		v, ok := row[f.Column]
		if !ok || v == nil {
			continue
		}
		switch f.Kind {
		case KindBoolEnum:
			id, _ := toInt64V3(v)
			code, err := enumCodeByID(r.db.WithContext(ctx), id)
			if err == nil {
				answers[f.Key] = fmt.Sprintf("%v", code == "y")
			}
		case KindEnum1:
			id, _ := toInt64V3(v)
			ic, err := enumICodeByID(r.db.WithContext(ctx), id)
			if err == nil && ic != "" {
				answers[f.Key] = ic
			}
		default:
			answers[f.Key] = fmt.Sprintf("%v", v)
		}
	}

	// KindEnumN -- icodes seleccionados por categoría, vía la relación
	// compartida (acotado por categoría, igual criterio que saveEnumNFields).
	for key, categories := range victimCaseV3EnumNCategories {
		var icodes []string
		err := r.db.WithContext(ctx).Raw(`
			SELECT e.victim_case_form2_enums_i_code
			FROM salvia.rel_victim_case_form2_enums_victim_case_form2 r
			JOIN salvia.victim_case_form2_enums e ON e.victim_case_form2_enums_id = r.victim_case_form2_enums_id
			WHERE r.victim_case_form2_id = ? AND e.victim_case_form2_enums_category = ANY(string_to_array(?, ','))
		`, form2ID, strings.Join(categories, ",")).Scan(&icodes).Error
		if err == nil && len(icodes) > 0 {
			answers[key] = strings.Join(icodes, ",")
		}
	}

	// Campos núcleo de victim_case (CI_1/CI_2/CI_5/CI_6/DP_28 viven en esa
	// tabla, no en form2 -- ver SaveSection).
	var core struct {
		Names, LastNames, DocType, DocNumber, TownCode string
	}
	if err := r.db.WithContext(ctx).Table("salvia.victim_case").
		Select(`victim_case_victim_names AS names, victim_case_victim_last_names AS last_names,
			victim_case_victim_doc_type AS doc_type, victim_case_victim_doc_number AS doc_number,
			victim_case_victim_town_code AS town_code`).
		Where("victim_case_i_code = ?", icode).Scan(&core).Error; err == nil {
		if core.Names != "" && core.Names != dummyNamesV3 {
			full := core.Names
			if core.LastNames != "" && core.LastNames != dummyLastNamesV3 {
				full = full + " " + core.LastNames
			}
			answers[FieldKeyV3("CI", 2)] = full
		}
		if core.DocType != "" {
			// BUG encontrado en QA (2026-10-15): NO se puede comparar contra
			// dummyDocTypeV3 aquí -- el dummy es "cc" (Cédula de Ciudadanía)
			// y es, a la vez, el tipo de documento real más común en
			// Colombia. Cualquier víctima con cédula real quedaba excluida
			// de 'answers' (se confundía con "nunca respondido"). El filtro
			// correcto de "¿lo respondió el agente?" ya lo hace
			// GetAnsweredKeys/answeredKeys en el frontend (H3) -- aquí basta
			// con que el dato exista.
			//
			// CI_5 es KindEnum1 -- el frontend espera icode, no el code
			// corto que guarda victim_case_victim_doc_type (ver SaveSection).
			if ic, err := resolveEnumICodeByCategoryAndCode(r.db.WithContext(ctx), "victim_case_form2_victim_doc_type", core.DocType); err == nil && ic != "" {
				answers[FieldKeyV3("CI", 5)] = ic
			}
		}
		if core.DocNumber != "" && core.DocNumber != dummyDocNumberV3 {
			answers[FieldKeyV3("CI", 6)] = core.DocNumber
		}
		if core.TownCode != "" && core.TownCode != dummyTownCodeV3 {
			answers[FieldKeyV3("DP", 28)] = core.TownCode
		}
	}

	return answers, nil
}

func toInt64V3(v interface{}) (int64, error) {
	switch n := v.(type) {
	case int64:
		return n, nil
	case int32:
		return int64(n), nil
	case float64:
		return int64(n), nil
	default:
		var out int64
		_, err := fmt.Sscanf(fmt.Sprintf("%v", v), "%d", &out)
		return out, err
	}
}

func enumCodeByID(tx *gorm.DB, id int64) (string, error) {
	var code string
	err := tx.Raw(`SELECT victim_case_form2_enums_code FROM salvia.victim_case_form2_enums WHERE victim_case_form2_enums_id = ?`, id).Scan(&code).Error
	return code, err
}

func enumICodeByID(tx *gorm.DB, id int64) (string, error) {
	var ic string
	err := tx.Raw(`SELECT victim_case_form2_enums_i_code FROM salvia.victim_case_form2_enums WHERE victim_case_form2_enums_id = ?`, id).Scan(&ic).Error
	return ic, err
}

// resolveEnumCodeByICode -- ver CI_5 en SaveSection/GetAnswers: victim_case
// guarda el code corto (varchar2) en victim_case_victim_doc_type, pero el
// frontend siempre manda/espera icode para campos KindEnum1.
func resolveEnumCodeByICode(tx *gorm.DB, icode string) (string, error) {
	var code string
	err := tx.Raw(`SELECT victim_case_form2_enums_code FROM salvia.victim_case_form2_enums WHERE victim_case_form2_enums_i_code = ?`, icode).Scan(&code).Error
	return code, err
}

// resolveEnumICodeByCategoryAndCode -- inverso de resolveEnumCodeByICode,
// para reconstruir el icode esperado por el frontend a partir del code
// corto guardado en victim_case (ver GetAnswers, CI_5).
func resolveEnumICodeByCategoryAndCode(tx *gorm.DB, category, code string) (string, error) {
	var ic string
	err := tx.Raw(`SELECT victim_case_form2_enums_i_code FROM salvia.victim_case_form2_enums WHERE victim_case_form2_enums_category = ? AND victim_case_form2_enums_code = ? LIMIT 1`, category, code).Scan(&ic).Error
	return ic, err
}

// ListEnumCategories -- ver victimCaseV3EnumCategories. label se resuelve
// contra el diccionario Locale["sp"] (mismo que usa el resto del sistema
// para traducir victim_case_form2_enums_name, ej. "victim_case_form2_
// gender_identity_ho" -> "Hombre"); si no hay traducción (categorías
// propias de V3, el frontend ya tiene su propio texto en español y solo
// usa 'code'/'icode' de aquí) se devuelve el name crudo como fallback.
func (r *victimCaseV3Repository) ListEnumCategories(ctx context.Context) (map[string][]EnumOptionV3, error) {
	type row struct {
		Category string
		Code     string
		ICode    string
		Name     string
	}
	var rows []row
	err := r.db.WithContext(ctx).Raw(`
		SELECT victim_case_form2_enums_category AS category, victim_case_form2_enums_code AS code,
			victim_case_form2_enums_i_code AS i_code, victim_case_form2_enums_name AS name
		FROM salvia.victim_case_form2_enums
		WHERE victim_case_form2_enums_category = ANY(string_to_array(?, ','))
		ORDER BY victim_case_form2_enums_category, victim_case_form2_enums_code
	`, strings.Join(victimCaseV3EnumCategories, ",")).Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("listEnumCategories: %w", err)
	}

	out := map[string][]EnumOptionV3{}
	for _, rr := range rows {
		label := salvia_config.Locale["sp"][rr.Name]
		if label == "" {
			label = rr.Name
		}
		out[rr.Category] = append(out[rr.Category], EnumOptionV3{Code: rr.Code, ICode: rr.ICode, Label: label})
	}
	return out, nil
}

func (r *victimCaseV3Repository) GetCaseSummary(ctx context.Context, icode string) (*VictimCaseV3Summary, error) {
	var row VictimCaseV3Summary
	err := r.db.WithContext(ctx).Table("salvia.victim_case").
		Select(`victim_case_i_code AS i_code, victim_case_victim_names AS names,
			victim_case_victim_last_names AS last_names, victim_case_victim_doc_type AS doc_type,
			victim_case_victim_doc_number AS doc_number, victim_case_status AS status`).
		Where("victim_case_i_code = ?", icode).Scan(&row).Error
	return &row, err
}

func (r *victimCaseV3Repository) UpdateStatus(ctx context.Context, icode, status string) error {
	return r.db.WithContext(ctx).Table("salvia.victim_case").
		Where("victim_case_i_code = ?", icode).
		Updates(map[string]interface{}{"victim_case_status": status, "victim_case_update_date": time.Now()}).Error
}

func (r *victimCaseV3Repository) UpdateStatusIfCurrent(ctx context.Context, icode, expectedCurrent, newStatus string) error {
	return r.db.WithContext(ctx).Table("salvia.victim_case").
		Where("victim_case_i_code = ? AND victim_case_status = ?", icode, expectedCurrent).
		Updates(map[string]interface{}{"victim_case_status": newStatus, "victim_case_update_date": time.Now()}).Error
}

// MostRecentCaseByDocument -- D1/D9 (badge "Registro Completo/Incompleto"):
// resuelve el icode del Caso MÁS RECIENTE (por fecha de creación) con ese
// Tipo+Número de documento -- "Registro completo" se evalúa sobre ESE Caso,
// nunca sobre uno cualquiera de los que comparten documento. El llamador
// (service) decide qué campos cuentan como obligatorias del Perfil (CI+DP+CA)
// vía GetAnsweredKeys -- acá solo se resuelve CUÁL Caso mirar.
func (r *victimCaseV3Repository) MostRecentCaseByDocument(ctx context.Context, docType, docNumber string) (string, error) {
	if docType == "" || docNumber == "" {
		return "", nil
	}
	var icode string
	err := r.db.WithContext(ctx).Table("salvia.victim_case").
		Select("victim_case_i_code").
		Where("victim_case_victim_doc_type = ? AND victim_case_victim_doc_number = ?", docType, docNumber).
		Order("victim_case_creation_date DESC").
		Limit(1).
		Scan(&icode).Error
	return icode, err
}

func (r *victimCaseV3Repository) StoreCredentials(ctx context.Context, icode, login, password string) error {
	payload, err := json.Marshal(map[string]string{"login": login, "pass": password})
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).Table("salvia.victim_case").
		Where("victim_case_i_code = ?", icode).
		Update("victim_case_new_user_credentials", string(payload)).Error
}

// Entidades y Barreras (GS § C/D) -- AddEntity/ListEntities/AddBarrier/
// ListBarriers y sus RemoveX ya NO viven acá: se reconstruyeron sobre el
// modelo real (entity_case/entity_case_obligation/barrier_v2) directamente
// en el service layer (salvia/service/victim_case_v3_service.go), que puede
// importar internal/models y reusar EntityCaseRepository/BarrierV2Repository
// sin crear un ciclo de imports. Ver registro-v3-pendientes.md § B.
