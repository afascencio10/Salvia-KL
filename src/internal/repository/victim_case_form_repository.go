package repository

import (
	"bitsflow/common/utils"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

// VictimCaseFormFieldKind describe cómo se persiste la respuesta de una pregunta
// del formulario "Registro de Caso" (dinamic-form) sobre victim_case_form2.
type VictimCaseFormFieldKind string

const (
	KindText      VictimCaseFormFieldKind = "text"      // varchar/text plano
	KindInt       VictimCaseFormFieldKind = "int"        // entero plano (teléfonos, escalas)
	KindDate      VictimCaseFormFieldKind = "date"        // date
	KindTime      VictimCaseFormFieldKind = "time"        // time without time zone
	KindTimestamp VictimCaseFormFieldKind = "timestamp"   // timestamp
	KindEnum1     VictimCaseFormFieldKind = "enum1"        // FK simple -> victim_case_form2_enums_id
	KindEnumN     VictimCaseFormFieldKind = "enumN"        // multi-valor -> rel_victim_case_form2_enums_victim_case_form2
	KindBoolEnum  VictimCaseFormFieldKind = "bool_enum"    // 'true'/'false' -> FK a categoría yes_no
	KindScale     VictimCaseFormFieldKind = "scale"        // escala de dificultad -> entero 0-3
)

// VictimCaseFormField mapea una pregunta (por FieldKey "S{seccion}Q{orden}") a su columna en victim_case_form2.
type VictimCaseFormField struct {
	Key    string
	Column string
	Kind   VictimCaseFormFieldKind
}

// scaleValues mapea los value de las opciones fijas de escala de dificultad a un entero 0-3.
var scaleValues = map[string]int{
	"sin_dificultad":    0,
	"alguna_dificultad": 1,
	"mucha_dificultad":  2,
	"no_puede":          3,
}

// VictimCaseFormRepository gestiona la creación progresiva (Borrador -> Activo) del
// victim_case a partir del formulario dinámico "Registro de Caso".
type VictimCaseFormRepository interface {
	// FindICodeBySubmissionID retorna el i_code del victim_case ya vinculado a este
	// form_submission, si existe (idempotencia de E-03/E-04).
	FindICodeBySubmissionID(ctx context.Context, submissionID string) (string, bool, error)

	// CreateDraft crea el general_user y el victim_case en estado 'bo' (Borrador)
	// con los datos mínimos de la Sección 1. Retorna el i_code del caso nuevo.
	CreateDraft(ctx context.Context, input CreateDraftInput) (string, error)

	// UpdateDraftCoreFields sincroniza nombres/documento/municipio en victim_case
	// en cada guardado de sección posterior a la creación del draft -- el
	// llamador ya resolvió valores dummy para lo que falte, por lo que este
	// update nunca se pospone.
	UpdateDraftCoreFields(ctx context.Context, iCode, names, lastNames, docType, docNumber, townCode string) error

	// UpsertForm2FromAnswers proyecta las respuestas disponibles del submission
	// sobre victim_case_form2 -- lo crea si no existe, o lo actualiza si ya
	// existe. Se llama en CADA guardado de sección (no solo al completar), por
	// lo que muchas respuestas pueden faltar aún: las columnas NOT NULL sin
	// responder (enum1/bool_enum) reciben el i_code "No" como valor dummy
	// (ver KindBoolEnum/KindEnum1 en la implementación), así el upsert nunca
	// se pospone -- se sobrescribe con la respuesta real en cuanto llega.
	// answersByKey está indexado por FieldKey (ver victim_case_form_fields.go),
	// no por questionID ni descripción -- así se evita ambigüedad entre
	// preguntas de Tamizaje pareja/no-pareja que comparten texto idéntico.
	UpsertForm2FromAnswers(ctx context.Context, iCode string, answersByKey map[string]string, riskScore, riskLevel int) (created bool, err error)

	// MarkActive transiciona el caso de 'bo' (Borrador) a 'ra' (activo),
	// refrescando nombres/documento/municipio de atención si cambiaron.
	// Asume que UpsertForm2FromAnswers ya creó victim_case_form2 exitosamente.
	MarkActive(ctx context.Context, iCode string, answersByKey map[string]string) error

	// StoreCredentials guarda login/password en texto plano en victim_case_new_user_credentials.
	StoreCredentials(ctx context.Context, iCode, login, password string) error

	// GetCredentials lee de vuelta las credenciales guardadas (para el GET de resultado).
	GetCredentials(ctx context.Context, iCode string) (login, password string, err error)

	// ResolveEnumCode resuelve un i_code de victim_case_form2_enums a su "code"
	// corto (ej. 'pi', 'ex') -- usado por el servicio para lógica de negocio
	// (wasPartner, etc.) que compara contra códigos, no contra i_codes.
	ResolveEnumCode(ctx context.Context, icode string) (string, error)

	// BuildAnswersByFieldKey junta answer+question+form_section para devolver
	// las respuestas de un submission indexadas por FieldKey ("S{n}Q{n}"),
	// listas para pasar a ActivateFromAnswers o para lógica de negocio del
	// servicio (wasPartner, riskScore, etc.).
	BuildAnswersByFieldKey(ctx context.Context, formID, submissionID string) (map[string]string, error)
}

type CreateDraftInput struct {
	SubmissionID       string
	Names              string
	LastNames          string
	DocType            string
	DocNumber          string
	ResidenceTownCode  string // fallback para victim_case_victim_town_code (igual que SetVictimCase hoy)
	GeneralUserICode   string
	GeneralUserLogin   string
	GeneralUserPassSha string // password ya hasheada (bcrypt)
	GeneralUserGender  string // code, ej. "mu"/"ho" -- puede venir vacío en la Sección 1
}

type victimCaseFormRepository struct {
	db *gorm.DB
}

func NewVictimCaseFormRepository(db *gorm.DB) VictimCaseFormRepository {
	return &victimCaseFormRepository{db: db}
}

func (r *victimCaseFormRepository) FindICodeBySubmissionID(ctx context.Context, submissionID string) (string, bool, error) {
	var iCode string
	err := r.db.WithContext(ctx).
		Table("salvia.victim_case").
		Select("victim_case_i_code").
		Where("victim_case_form_submission_id = ?", submissionID).
		Limit(1).
		Scan(&iCode).Error
	if err != nil {
		return "", false, err
	}
	return iCode, iCode != "", nil
}

func (r *victimCaseFormRepository) CreateDraft(ctx context.Context, input CreateDraftInput) (string, error) {
	return r.createDraftTx(ctx, input)
}

func (r *victimCaseFormRepository) createDraftTx(ctx context.Context, input CreateDraftInput) (string, error) {
	var newICode string

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Perfil de usuario
		var profileID int64
		err := tx.Raw(`
			INSERT INTO security.general_user_profile
				(general_user_profile_i_code, general_user_profile_creation_date, general_user_profile_update_date,
				 general_user_profile_gender, general_user_profile_names, general_user_profile_last_names,
				 general_user_profile_doc_type, general_user_profile_doc_number, general_user_profile_town)
			VALUES (gen_random_uuid()::text, now(), now(), ?, ?, ?, ?, ?, ?)
			RETURNING general_user_profile_id
		`, input.GeneralUserGender, input.Names, input.LastNames, input.DocType, input.DocNumber, input.ResidenceTownCode).
			Scan(&profileID).Error
		if err != nil {
			return fmt.Errorf("crear general_user_profile: %w", err)
		}

		// 2. Usuario
		var generalUserID int64
		err = tx.Raw(`
			INSERT INTO security.general_user
				(general_user_i_code, general_user_creation_date, general_user_update_date,
				 general_user_login, general_user_password, general_user_status, general_user_language,
				 general_user_general_user_profile)
			VALUES (?, now(), now(), ?, ?, 'e', 'sp', ?)
			RETURNING general_user_id
		`, input.GeneralUserICode, input.GeneralUserLogin, input.GeneralUserPassSha, profileID).
			Scan(&generalUserID).Error
		if err != nil {
			return fmt.Errorf("crear general_user: %w", err)
		}

		// 3. Rol "us" (Usuario externo)
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
		townCode := input.ResidenceTownCode // mismo fallback que SetVictimCase (residencia si no hay atención)
		if err := tx.Exec(`
			INSERT INTO salvia.victim_case
				(victim_case_i_code, victim_case_creation_date, victim_case_update_date, victim_case_status,
				 victim_case_victim_names, victim_case_victim_last_names, victim_case_victim_doc_type,
				 victim_case_victim_doc_number, victim_case_victim_town_code, victim_case_general_user,
				 victim_case_form_submission_id)
			VALUES (?, now(), now(), 'bo', ?, ?, ?, ?, ?, ?, ?)
		`, newICode, input.Names, input.LastNames, input.DocType, input.DocNumber, townCode, input.GeneralUserICode, input.SubmissionID).Error; err != nil {
			return fmt.Errorf("crear victim_case: %w", err)
		}

		return nil
	})

	return newICode, err
}

// UpdateDraftCoreFields sincroniza nombres/documento/municipio en victim_case
// en CADA guardado de sección (no solo al crear el draft) -- el llamador ya
// resuelve valores dummy para lo que aún no se ha respondido, así el update
// nunca se pospone. Los valores reales (cuando lleguen) sobrescriben el dummy.
func (r *victimCaseFormRepository) UpdateDraftCoreFields(ctx context.Context, iCode, names, lastNames, docType, docNumber, townCode string) error {
	return r.db.WithContext(ctx).
		Table("salvia.victim_case").
		Where("victim_case_i_code = ?", iCode).
		Updates(map[string]interface{}{
			"victim_case_victim_names":      names,
			"victim_case_victim_last_names": lastNames,
			"victim_case_victim_doc_type":   docType,
			"victim_case_victim_doc_number": docNumber,
			"victim_case_victim_town_code":  townCode,
		}).Error
}

func (r *victimCaseFormRepository) StoreCredentials(ctx context.Context, iCode, login, password string) error {
	payload, err := json.Marshal(map[string]string{"login": login, "pass": password})
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).
		Table("salvia.victim_case").
		Where("victim_case_i_code = ?", iCode).
		Update("victim_case_new_user_credentials", string(payload)).Error
}

func (r *victimCaseFormRepository) GetCredentials(ctx context.Context, iCode string) (string, string, error) {
	var raw string
	err := r.db.WithContext(ctx).
		Table("salvia.victim_case").
		Select("victim_case_new_user_credentials").
		Where("victim_case_i_code = ?", iCode).
		Scan(&raw).Error
	if err != nil {
		return "", "", err
	}
	if raw == "" {
		return "", "", nil
	}
	var parsed struct {
		Login string `json:"login"`
		Pass  string `json:"pass"`
	}
	if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
		return "", "", err
	}
	return parsed.Login, parsed.Pass, nil
}

// ─── Proyección incremental de respuestas (llamada en cada sección) ────────────

// isNotNullViolation detecta SQLSTATE 23502 (not_null_violation) en el texto
// del error -- Postgres/pgx no siempre exponen un tipo tipado fácil de
// inspeccionar a través de GORM, así que se detecta por código en el mensaje.
func isNotNullViolation(err error) bool {
	return err != nil && strings.Contains(err.Error(), "23502")
}

func (r *victimCaseFormRepository) UpsertForm2FromAnswers(ctx context.Context, iCode string, answers map[string]string, riskScore, riskLevel int) (bool, error) {
	created := false
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var victimCaseID int64
		if err := tx.Raw(`SELECT victim_case_id FROM salvia.victim_case WHERE victim_case_i_code = ?`, iCode).Scan(&victimCaseID).Error; err != nil {
			return fmt.Errorf("buscar victim_case_id: %w", err)
		}
		if victimCaseID == 0 {
			return fmt.Errorf("victim_case no encontrado para i_code=%s", iCode)
		}

		var existingForm2ID int64
		if err := tx.Raw(`SELECT victim_case_form2_id FROM salvia.victim_case_form2 WHERE victim_case_form2_victim_case = ?`, victimCaseID).Scan(&existingForm2ID).Error; err != nil {
			return fmt.Errorf("buscar victim_case_form2 existente: %w", err)
		}

		yesID, noID, err := r.resolveYesNoIDs(tx)
		if err != nil {
			return err
		}

		cols := map[string]interface{}{}
		var multiValueInserts []struct {
			column string
			icodes []string
		}

		for _, f := range victimCaseFormFields {
			raw := answers[f.Key] // zero value "" si no se ha respondido todavía
			switch f.Kind {
			case KindText:
				cols[f.Column] = raw

			case KindInt:
				if raw == "" {
					cols[f.Column] = 0
				} else {
					n, _ := strconv.ParseInt(raw, 10, 64)
					cols[f.Column] = n
				}

			case KindDate:
				cols[f.Column] = nullableOrNow(raw, "2006-01-02")

			case KindTime:
				// El input es datetime-local ("2026-06-01T10:00") pero la
				// columna es time without time zone -- se extrae solo la hora.
				cols[f.Column] = extractTimePart(raw)

			case KindTimestamp:
				cols[f.Column] = nullableTimestampOrNow(raw)

			case KindBoolEnum:
				if raw == "true" {
					cols[f.Column] = yesID
				} else if raw == "false" {
					cols[f.Column] = noID
				} else {
					// Aún sin responder -- se usa "No" como valor dummy para que
					// la columna NOT NULL nunca bloquee el upsert. Se sobrescribe
					// con la respuesta real en cuanto la sección correspondiente
					// se guarde.
					cols[f.Column] = noID
				}

			case KindEnum1:
				if raw == "" {
					// Aún sin responder -- mismo criterio que KindBoolEnum: se usa
					// el i_code "No" como placeholder (cualquier fila válida de
					// victim_case_form2_enums sirve para satisfacer la FK/NOT NULL).
					cols[f.Column] = noID
					continue
				}
				id, err := r.resolveEnumID(tx, raw)
				if err != nil {
					return fmt.Errorf("resolver enum %q (%s): %w", raw, f.Key, err)
				}
				cols[f.Column] = id

			case KindEnumN:
				if raw == "" {
					continue
				}
				icodes := splitCSV(raw)
				multiValueInserts = append(multiValueInserts, struct {
					column string
					icodes []string
				}{f.Column, icodes})

			case KindScale:
				v, ok := scaleValues[raw]
				if !ok {
					v = 0
				}
				cols[f.Column] = v
			}
		}

		cols["victim_case_form2_update_date"] = time.Now()
		cols["victim_case_form2_risk_score"] = riskScore
		cols["victim_case_form2_risk_level"] = riskLevel

		var form2ID int64

		if existingForm2ID == 0 {
			// Primera vez que hay suficientes datos -- intentar crear.
			cols["victim_case_form2_i_code"] = utils.GetUUID()
			cols["victim_case_form2_creation_date"] = time.Now()
			cols["victim_case_form2_victim_case"] = victimCaseID

			insertCols, insertVals, placeholders := buildInsert(cols)
			sql := fmt.Sprintf(
				`INSERT INTO salvia.victim_case_form2 (%s) VALUES (%s) RETURNING victim_case_form2_id`,
				strings.Join(insertCols, ", "), strings.Join(placeholders, ", "),
			)
			if err := tx.Raw(sql, insertVals...).Scan(&form2ID).Error; err != nil {
				if isNotNullViolation(err) {
					// No debería ocurrir: todo NOT NULL de victim_case_form2 tiene
					// un valor dummy de respaldo en victimCaseFormFields. Si esto
					// dispara, es una columna NOT NULL sin mapear -- error real.
					return fmt.Errorf("crear victim_case_form2: columna NOT NULL sin valor dummy de respaldo: %w", err)
				}
				return fmt.Errorf("crear victim_case_form2: %w", err)
			}
			created = true
		} else {
			// Ya existe -- actualizar con los valores más recientes.
			form2ID = existingForm2ID
			if err := tx.Table("salvia.victim_case_form2").Where("victim_case_form2_id = ?", form2ID).Updates(cols).Error; err != nil {
				return fmt.Errorf("actualizar victim_case_form2: %w", err)
			}
			// Limpiar relaciones multi-valor previas antes de reinsertar (idempotente).
			if err := tx.Exec(`DELETE FROM salvia.rel_victim_case_form2_enums_victim_case_form2 WHERE victim_case_form2_id = ?`, form2ID).Error; err != nil {
				return fmt.Errorf("limpiar relaciones multi-valor: %w", err)
			}
		}

		for _, mv := range multiValueInserts {
			for _, icode := range mv.icodes {
				enumID, err := r.resolveEnumID(tx, icode)
				if err != nil {
					// icode inválido/huérfano -- se ignora, igual que filterAnswers en dinamic-form
					continue
				}
				if err := tx.Exec(`
					INSERT INTO salvia.rel_victim_case_form2_enums_victim_case_form2
						(victim_case_form2_enums_id, victim_case_form2_id)
					VALUES (?, ?)
				`, enumID, form2ID).Error; err != nil {
					return fmt.Errorf("crear relación multi-valor (%s): %w", mv.column, err)
				}
			}
		}

		return nil
	})

	return created, err
}

func (r *victimCaseFormRepository) MarkActive(ctx context.Context, iCode string, answers map[string]string) error {
	updates := map[string]interface{}{"victim_case_status": "ra"}
	if v, ok := answers[FieldKey(1, 1)]; ok && v != "" {
		updates["victim_case_victim_names"] = v
	}
	if v, ok := answers[FieldKey(1, 2)]; ok && v != "" {
		updates["victim_case_victim_last_names"] = v
	}
	if v, ok := answers[FieldKey(1, 5)]; ok && v != "" {
		updates["victim_case_victim_doc_type"] = v
	}
	if v, ok := answers[FieldKey(1, 6)]; ok && v != "" {
		updates["victim_case_victim_doc_number"] = v
	}
	if atencionTown, ok := answers[FieldKey(9, 3)]; ok && atencionTown != "" {
		updates["victim_case_victim_town_code"] = atencionTown
	}
	return r.db.WithContext(ctx).
		Table("salvia.victim_case").
		Where("victim_case_i_code = ?", iCode).
		Updates(updates).Error
}

func (r *victimCaseFormRepository) BuildAnswersByFieldKey(ctx context.Context, formID, submissionID string) (map[string]string, error) {
	type row struct {
		SectionOrder  int
		QuestionOrder int
		Value         string
	}
	var rows []row
	err := r.db.WithContext(ctx).Raw(`
		SELECT fs."order" AS section_order, q."order" AS question_order, a.value AS value
		FROM salvia.answer a
		JOIN salvia.question q ON q.id = a.question_id
		JOIN salvia.form_section fs ON fs.id::varchar = q.form_section_id::varchar
		WHERE a.form_submission_id = ? AND q.form_id = ?
	`, submissionID, formID).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make(map[string]string, len(rows))
	for _, rw := range rows {
		out[FieldKey(rw.SectionOrder, rw.QuestionOrder)] = rw.Value
	}
	return out, nil
}

func (r *victimCaseFormRepository) ResolveEnumCode(ctx context.Context, icode string) (string, error) {
	if icode == "" {
		return "", nil
	}
	var code string
	err := r.db.WithContext(ctx).
		Raw(`SELECT victim_case_form2_enums_code FROM salvia.victim_case_form2_enums WHERE victim_case_form2_enums_i_code = ?`, icode).
		Scan(&code).Error
	return code, err
}

func (r *victimCaseFormRepository) resolveYesNoIDs(tx *gorm.DB) (int64, int64, error) {
	var yesID, noID int64
	if err := tx.Raw(`SELECT victim_case_form2_enums_id FROM salvia.victim_case_form2_enums WHERE victim_case_form2_enums_category = 'yes_no' AND victim_case_form2_enums_code = 'y' LIMIT 1`).Scan(&yesID).Error; err != nil {
		return 0, 0, fmt.Errorf("resolver enum yes_no/y: %w", err)
	}
	if err := tx.Raw(`SELECT victim_case_form2_enums_id FROM salvia.victim_case_form2_enums WHERE victim_case_form2_enums_category = 'yes_no' AND victim_case_form2_enums_code = 'n' LIMIT 1`).Scan(&noID).Error; err != nil {
		return 0, 0, fmt.Errorf("resolver enum yes_no/n: %w", err)
	}
	return yesID, noID, nil
}

// resolveEnumID resuelve un i_code de victim_case_form2_enums a su ID interno.
// Las respuestas de preguntas con state_options_path=enums.* llegan como i_code
// (mismo valor que el frontend recibe en formState.enums.<categoria>).
func (r *victimCaseFormRepository) resolveEnumID(tx *gorm.DB, icode string) (int64, error) {
	var id int64
	err := tx.Raw(`SELECT victim_case_form2_enums_id FROM salvia.victim_case_form2_enums WHERE victim_case_form2_enums_i_code = ?`, icode).Scan(&id).Error
	if err != nil {
		return 0, err
	}
	if id == 0 {
		return 0, fmt.Errorf("i_code no encontrado: %s", icode)
	}
	return id, nil
}

func splitCSV(v string) []string {
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func nullableOrNow(raw, layout string) interface{} {
	if raw == "" {
		return time.Now()
	}
	t, err := time.Parse(layout, raw)
	if err != nil {
		return time.Now()
	}
	return t
}

// extractTimePart obtiene "HH:MM:SS" de un valor datetime-local ("2026-06-01T10:00")
// u hora simple ("10:00"). Devuelve "00:00:00" si no se puede interpretar.
func extractTimePart(raw string) string {
	if raw == "" {
		return "00:00:00"
	}
	s := raw
	if idx := strings.Index(s, "T"); idx != -1 {
		s = s[idx+1:]
	}
	if t, err := time.Parse("15:04:05", s); err == nil {
		return t.Format("15:04:05")
	}
	if t, err := time.Parse("15:04", s); err == nil {
		return t.Format("15:04:05")
	}
	return "00:00:00"
}

func nullableTimestampOrNow(raw string) interface{} {
	if raw == "" {
		return time.Now()
	}
	if t, err := time.Parse("2006-01-02T15:04", raw); err == nil {
		return t
	}
	if t, err := time.Parse(time.RFC3339, raw); err == nil {
		return t
	}
	return time.Now()
}

func buildInsert(cols map[string]interface{}) ([]string, []interface{}, []string) {
	names := make([]string, 0, len(cols))
	vals := make([]interface{}, 0, len(cols))
	placeholders := make([]string, 0, len(cols))
	for k, v := range cols {
		names = append(names, k)
		vals = append(vals, v)
		placeholders = append(placeholders, "?")
	}
	return names, vals, placeholders
}

