package repository

import (
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	"context"
	"encoding/json"
	"fmt"
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

	GetAnswers(ctx context.Context, icode string) (map[string]string, error)
	GetCaseSummary(ctx context.Context, icode string) (*VictimCaseV3Summary, error)
	UpdateStatus(ctx context.Context, icode, status string) error
	// UpdateStatusIfCurrent transiciona solo si el estado actual es
	// expectedCurrent -- evita que un recálculo tardío (ej. triadaCompleta
	// reevaluada tras eliminar una respuesta y volver a completarla) revierta
	// un Caso que el agente ya hizo avanzar a Terminado/Cerrado por otra vía.
	UpdateStatusIfCurrent(ctx context.Context, icode, expectedCurrent, newStatus string) error
	IsRegistroComplete(ctx context.Context, docType, docNumber string) (bool, error)

	StoreCredentials(ctx context.Context, icode, login, password string) error

	AddEntity(ctx context.Context, e VictimCaseEntityInput) (string, error)
	RemoveEntity(ctx context.Context, entityICode string) error
	ListEntities(ctx context.Context, icode string) ([]VictimCaseEntityInput, error)

	AddBarrier(ctx context.Context, b VictimCaseBarrierInput) (string, error)
	RemoveBarrier(ctx context.Context, barrierICode string) error
	ListBarriers(ctx context.Context, icode string) ([]VictimCaseBarrierInput, error)

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
}

type CreateDraftV3Input struct {
	DocType           string
	DocNumber         string
	Phone             string
	ResidenceTownCode string
	AgentID           string
	Channel           string
	GeneralUserICode  string
	GeneralUserLogin  string
	GeneralUserPassSha string
}

type VictimCaseEntityInput struct {
	ICode              string
	VictimCaseICode    string
	Kind               string // 'acudida' | 'activacion'
	SectorCode         string // 'sa'|'ju'|'pr'|'ot' -- se resuelve a enum id
	Department         string
	City               string
	Town               string
	Name               string
	FileNumber         string
	ReportDate         string
	ActivationChannel  string // code, solo si Kind='activacion'
	CompletedActions   []string
	PendingActions     []string
	OtherInfo          string
	Agent              string
}

type VictimCaseBarrierInput struct {
	ICode           string
	VictimCaseICode string
	SectorCode      string
	Specific        []string
	OtherDetail     string
	EntityName      string
	Institutional   []string
	Economic        []string
	Territorial     []string
	Differential    []string
	Date            string
	Official        string
	Description     string
	Management      []string
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
// columnas NOT NULL resueltas al dummy "No" (categoría yes_no, code='n') --
// tanto las 55 heredadas de V2 (victimCaseFormFields) como cualquier
// columna NOT NULL propia de V3. Las columnas nullable se dejan NULL.
func insertDummyForm2Row(tx *gorm.DB, victimCaseID int64) error {
	var noID int64
	if err := tx.Raw(`SELECT victim_case_form2_enums_id FROM salvia.victim_case_form2_enums WHERE victim_case_form2_enums_category = 'yes_no' AND victim_case_form2_enums_code = 'n' LIMIT 1`).Scan(&noID).Error; err != nil {
		return fmt.Errorf("resolver dummy 'No': %w", err)
	}

	cols := map[string]interface{}{
		"victim_case_form2_victim_case":   victimCaseID,
		"victim_case_form2_creation_date": time.Now(),
		"victim_case_form2_update_date":   time.Now(),
	}
	// Todas las columnas NOT NULL heredadas de V2 (enum1/bool_enum) → dummy "No".
	for _, f := range victimCaseFormFields {
		if f.Column != "" && (f.Kind == KindBoolEnum || f.Kind == KindEnum1) {
			cols[f.Column] = noID
		}
	}
	// Columnas propias de V3 que son NOT NULL en el esquema -- ninguna lo es
	// hoy (todas nullable, ver migrations/add_registro_v3.sql), se dejan sin
	// setear (quedan NULL hasta que el módulo correspondiente se guarde).

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
	if v, ok := answersByKey[FieldKeyV3("CI", 5)]; ok {
		coreUpdates["victim_case_victim_doc_type"] = v
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
	return nil
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
		if core.DocType != "" && core.DocType != dummyDocTypeV3 {
			answers[FieldKeyV3("CI", 5)] = core.DocType
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

// IsRegistroComplete -- D1: "Registro completo" es calculado, no una
// columna. Evalúa sobre el Caso MÁS RECIENTE con ese documento si las
// obligatorias [1er contacto] de CI+DP+CA están todas respondidas.
// Implementación completa de la lista de obligatorias pendiente -- ver
// GAP G-V3-22.
func (r *victimCaseV3Repository) IsRegistroComplete(ctx context.Context, docType, docNumber string) (bool, error) {
	return false, nil
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

// ─── Entidades (GS § C) ──────────────────────────────────────────────────

func (r *victimCaseV3Repository) AddEntity(ctx context.Context, e VictimCaseEntityInput) (string, error) {
	var victimCaseID int64
	if err := r.db.WithContext(ctx).Raw(`SELECT victim_case_id FROM salvia.victim_case WHERE victim_case_i_code = ?`, e.VictimCaseICode).Scan(&victimCaseID).Error; err != nil {
		return "", fmt.Errorf("addEntity: %w", err)
	}
	sectorID, err := resolveEnumIDByCategory(r.db.WithContext(ctx), "victim_case_form2_sector", e.SectorCode)
	if err != nil {
		return "", fmt.Errorf("addEntity: resolver sector: %w", err)
	}
	var channelID *int64
	if e.ActivationChannel != "" {
		id, err := resolveEnumIDByCategory(r.db.WithContext(ctx), "victim_case_form2_activation_channel", e.ActivationChannel)
		if err != nil {
			return "", fmt.Errorf("addEntity: resolver canal: %w", err)
		}
		channelID = &id
	}
	completedJSON, _ := json.Marshal(e.CompletedActions)
	pendingJSON, _ := json.Marshal(e.PendingActions)

	newICode := utils.GetUUID()
	err = r.db.WithContext(ctx).Exec(`
		INSERT INTO salvia.victim_case_entity
			(victim_case_entity_i_code, victim_case_entity_victim_case, victim_case_entity_kind,
			 victim_case_entity_sector, victim_case_entity_department, victim_case_entity_city,
			 victim_case_entity_town, victim_case_entity_name, victim_case_entity_file_number,
			 victim_case_entity_report_date, victim_case_entity_activation_channel,
			 victim_case_entity_completed_actions, victim_case_entity_pending_actions,
			 victim_case_entity_other_info, victim_case_entity_agent,
			 victim_case_entity_creation_date, victim_case_entity_update_date)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, now(), now())
	`, newICode, victimCaseID, e.Kind, sectorID, nullIfEmpty(e.Department), nullIfEmpty(e.City), nullIfEmpty(e.Town),
		nullIfEmpty(e.Name), nullIfEmpty(e.FileNumber), nullIfEmptyDate(e.ReportDate), channelID,
		string(completedJSON), string(pendingJSON), nullIfEmpty(e.OtherInfo), e.Agent).Error
	if err != nil {
		return "", fmt.Errorf("addEntity: insertar: %w", err)
	}
	return newICode, nil
}

func (r *victimCaseV3Repository) RemoveEntity(ctx context.Context, entityICode string) error {
	return r.db.WithContext(ctx).Exec(`UPDATE salvia.victim_case_entity SET victim_case_entity_deleted_at = now() WHERE victim_case_entity_i_code = ?`, entityICode).Error
}

func (r *victimCaseV3Repository) ListEntities(ctx context.Context, icode string) ([]VictimCaseEntityInput, error) {
	var out []VictimCaseEntityInput
	err := r.db.WithContext(ctx).Raw(`
		SELECT e.victim_case_entity_i_code AS i_code, e.victim_case_entity_kind AS kind,
			e.victim_case_entity_name AS name, e.victim_case_entity_file_number AS file_number
		FROM salvia.victim_case_entity e
		JOIN salvia.victim_case vc ON vc.victim_case_id = e.victim_case_entity_victim_case
		WHERE vc.victim_case_i_code = ? AND e.victim_case_entity_deleted_at IS NULL
		ORDER BY e.victim_case_entity_creation_date
	`, icode).Scan(&out).Error
	return out, err
}

// ─── Barreras (GS § D) ───────────────────────────────────────────────────

func (r *victimCaseV3Repository) AddBarrier(ctx context.Context, b VictimCaseBarrierInput) (string, error) {
	var victimCaseID int64
	if err := r.db.WithContext(ctx).Raw(`SELECT victim_case_id FROM salvia.victim_case WHERE victim_case_i_code = ?`, b.VictimCaseICode).Scan(&victimCaseID).Error; err != nil {
		return "", fmt.Errorf("addBarrier: %w", err)
	}
	sectorID, err := resolveEnumIDByCategory(r.db.WithContext(ctx), "victim_case_form2_sector", b.SectorCode)
	if err != nil {
		return "", fmt.Errorf("addBarrier: resolver sector: %w", err)
	}
	specificJSON, _ := json.Marshal(b.Specific)
	instJSON, _ := json.Marshal(b.Institutional)
	econJSON, _ := json.Marshal(b.Economic)
	terrJSON, _ := json.Marshal(b.Territorial)
	diffJSON, _ := json.Marshal(b.Differential)
	mgmtJSON, _ := json.Marshal(b.Management)

	newICode := utils.GetUUID()
	err = r.db.WithContext(ctx).Exec(`
		INSERT INTO salvia.victim_case_barrier
			(victim_case_barrier_i_code, victim_case_barrier_victim_case, victim_case_barrier_sector,
			 victim_case_barrier_specific, victim_case_barrier_other_detail, victim_case_barrier_entity_name,
			 victim_case_barrier_institutional, victim_case_barrier_economic, victim_case_barrier_territorial,
			 victim_case_barrier_differential, victim_case_barrier_date, victim_case_barrier_official,
			 victim_case_barrier_description, victim_case_barrier_management, victim_case_barrier_agent,
			 victim_case_barrier_creation_date, victim_case_barrier_update_date)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, now(), now())
	`, newICode, victimCaseID, sectorID, string(specificJSON), nullIfEmpty(b.OtherDetail), nullIfEmpty(b.EntityName),
		string(instJSON), string(econJSON), string(terrJSON), string(diffJSON), b.Date, b.Official, b.Description,
		string(mgmtJSON), b.Agent).Error
	if err != nil {
		return "", fmt.Errorf("addBarrier: insertar: %w", err)
	}
	return newICode, nil
}

func (r *victimCaseV3Repository) RemoveBarrier(ctx context.Context, barrierICode string) error {
	return r.db.WithContext(ctx).Exec(`UPDATE salvia.victim_case_barrier SET victim_case_barrier_deleted_at = now() WHERE victim_case_barrier_i_code = ?`, barrierICode).Error
}

func (r *victimCaseV3Repository) ListBarriers(ctx context.Context, icode string) ([]VictimCaseBarrierInput, error) {
	var out []VictimCaseBarrierInput
	err := r.db.WithContext(ctx).Raw(`
		SELECT b.victim_case_barrier_i_code AS i_code, b.victim_case_barrier_description AS description
		FROM salvia.victim_case_barrier b
		JOIN salvia.victim_case vc ON vc.victim_case_id = b.victim_case_barrier_victim_case
		WHERE vc.victim_case_i_code = ? AND b.victim_case_barrier_deleted_at IS NULL
		ORDER BY b.victim_case_barrier_creation_date
	`, icode).Scan(&out).Error
	return out, err
}

func nullIfEmpty(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

func nullIfEmptyDate(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
