package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"

	"bitsflow/internal/models"
)

// ─────────────────────────────────────────────────────────────────────────────
// Entrevista de caracterización — FASE 2b-2
//   1. InitCaracterizacion: crea (o reutiliza) el form_submission del formulario
//      dinámico, precargado desde el Registro de Caso, y lo enlaza a la tarea.
//   2. Medidas: al cerrar la entrevista (rama caracterizado) se reutilizan,
//      asocian o crean las salvia.emergency_measure de la remisión.
// ─────────────────────────────────────────────────────────────────────────────

// ErrCaracterizacionTaskInvalid indica que la tarea no es una
// entrevista_caracterizacion pendiente. El controller lo mapea a 422.
var ErrCaracterizacionTaskInvalid = errors.New("la tarea no es una entrevista de caracterización pendiente")

// ErrCaracterizacionForbidden indica que el agente (am) no es el asignado a la
// tarea. El controller lo mapea a 403.
var ErrCaracterizacionForbidden = errors.New("la tarea no está asignada a este usuario")

// ErrCaracterizacionFormMissing indica que el formulario de caracterización no
// está sembrado en la base de datos. El controller lo mapea a 422.
var ErrCaracterizacionFormMissing = errors.New("el formulario de caracterización no está disponible")

// CaracterizacionInit es la respuesta del init de la entrevista.
type CaracterizacionInit struct {
	SubmissionID string   `json:"submissionId"`
	FormID       string   `json:"formId"`
	Reused       bool     `json:"reused"`
	Prefilled    int      `json:"prefilled"`
	Omitted      []string `json:"omitted"`
}

// submissionIDDeFormData extrae la llave submissionId de un form_data JSON.
func submissionIDDeFormData(raw []byte) string {
	if len(raw) == 0 {
		return ""
	}
	var m map[string]interface{}
	if err := json.Unmarshal(raw, &m); err != nil {
		return ""
	}
	v, _ := m["submissionId"].(string)
	return strings.TrimSpace(v)
}

// submissionCaracterizacionValido indica si el submission existe, no está
// borrado y pertenece al formulario de caracterización.
func submissionCaracterizacionValido(db *gorm.DB, id string) bool {
	id = strings.TrimSpace(id)
	if id == "" {
		return false
	}
	var n int64
	if err := db.Raw(`
		SELECT COUNT(*)
		  FROM salvia.form_submission
		 WHERE id::text = ?
		   AND form_id::text = ?
		   AND deleted_at IS NULL
	`, id, caracterizacionFormID).Scan(&n).Error; err != nil {
		return false
	}
	return n > 0
}

// InitCaracterizacion crea (o reutiliza) el submission de la entrevista de
// caracterización de una tarea, precargado y editable. Idempotente: si la tarea
// ya tiene form_data.submissionId válido, devuelve ese mismo submission.
// role es el rol de sesión (am o sv): el am debe ser el asignado a la tarea.
func (s *caseTaskService) InitCaracterizacion(ctx context.Context, taskID, userID, role string) (*CaracterizacionInit, error) {
	if s.db == nil {
		return nil, fmt.Errorf("InitCaracterizacion: sin base de datos")
	}
	task, err := s.GetByID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	if task.Type != models.CaseTaskTypeEntrevistaCaracterizacion || task.Status != models.CaseTaskStatusToDo {
		return nil, ErrCaracterizacionTaskInvalid
	}
	if role == models.CaseTaskRoleAgenteMedidas && strings.TrimSpace(task.AssignedUserID) != strings.TrimSpace(userID) {
		return nil, ErrCaracterizacionForbidden
	}

	// Camino rápido: ya hay un submission válido enlazado a la tarea.
	if sid := submissionIDDeFormData(task.FormData); submissionCaracterizacionValido(s.db.WithContext(ctx), sid) {
		return &CaracterizacionInit{SubmissionID: sid, FormID: caracterizacionFormID, Reused: true, Omitted: []string{}}, nil
	}

	// Lecturas fuera de la transacción: un SELECT fallido dentro de una
	// transacción de Postgres la aborta. Un fallo de lectura del Registro de
	// Caso no impide el init: el submission queda vacío.
	cat, err := cargarCatalogosPrecarga(ctx, s.db)
	if err != nil {
		return nil, err
	}
	if len(cat.Preguntas) == 0 {
		return nil, ErrCaracterizacionFormMissing
	}
	rc, err := respuestasRegistroCaso(ctx, s.db, task.CaseID)
	if err != nil {
		log.Printf("[InitCaracterizacion] WARN: sin precarga para el caso %s: %v", task.CaseID, err)
		rc = map[string]string{}
	}
	respuestas, omitidas := construirPrecargaCaracterizacion(rc, cat, time.Now())

	var result CaracterizacionInit
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Serializa aperturas simultáneas: bloquea la fila de la tarea.
		var raw string
		res := tx.Raw(`
			SELECT COALESCE(form_data::text, '')
			  FROM salvia.case_task
			 WHERE id::text = ?
			   AND status = ?
			   AND deleted_at IS NULL
			   FOR UPDATE
		`, task.ID, models.CaseTaskStatusToDo).Scan(&raw)
		if res.Error != nil {
			return fmt.Errorf("InitCaracterizacion: bloquear la tarea: %w", res.Error)
		}
		if res.RowsAffected == 0 {
			return ErrCaracterizacionTaskInvalid
		}
		if sid := submissionIDDeFormData([]byte(raw)); submissionCaracterizacionValido(tx, sid) {
			result = CaracterizacionInit{SubmissionID: sid, FormID: caracterizacionFormID, Reused: true, Omitted: []string{}}
			return nil
		}

		sub := &models.FormSubmission{FormID: caracterizacionFormID}
		if err := tx.Create(sub).Error; err != nil {
			return fmt.Errorf("InitCaracterizacion: crear el submission: %w", err)
		}
		// Answers en el orden de la tabla de reglas (determinista).
		answers := make([]models.Answer, 0, len(respuestas))
		for _, rule := range caracterizacionPrefillRules {
			qID := caracterizacionQuestionPrefix + rule.Dest
			if v, ok := respuestas[qID]; ok {
				answers = append(answers, models.Answer{FormSubmissionID: sub.ID, QuestionID: qID, Value: v})
			}
		}
		if len(answers) > 0 {
			if err := tx.Create(&answers).Error; err != nil {
				return fmt.Errorf("InitCaracterizacion: sembrar las respuestas: %w", err)
			}
		}
		// Enlace tarea <-> submission: merge no destructivo en form_data.
		if err := tx.Exec(`
			UPDATE salvia.case_task
			   SET form_data = (CASE WHEN jsonb_typeof(form_data) = 'object' THEN form_data ELSE '{}'::jsonb END)
			                   || jsonb_build_object('submissionId', ?::text),
			       updated_at = NOW()
			 WHERE id::text = ?
		`, sub.ID, task.ID).Error; err != nil {
			return fmt.Errorf("InitCaracterizacion: enlazar el submission a la tarea: %w", err)
		}
		omitted := make([]string, 0, len(omitidas))
		for _, o := range omitidas {
			omitted = append(omitted, o.Label+": "+o.Motivo)
		}
		result = CaracterizacionInit{SubmissionID: sub.ID, FormID: caracterizacionFormID, Prefilled: len(answers), Omitted: omitted}
		return nil
	}); err != nil {
		return nil, err
	}
	return &result, nil
}

// submissionIDCaracterizacion decide el submissionId a persistir al cerrar la
// entrevista. Prefiere el guardado en la tarea; el que manda el cliente solo
// vale si existe y es de este formulario. Devuelve "" si no hay ninguno válido.
func (s *caseTaskService) submissionIDCaracterizacion(ctx context.Context, task *models.CaseTask, fd map[string]interface{}) string {
	db := s.db.WithContext(ctx)
	if sid := submissionIDDeFormData(task.FormData); submissionCaracterizacionValido(db, sid) {
		return sid
	}
	if v, _ := fd["submissionId"].(string); submissionCaracterizacionValido(db, v) {
		return strings.TrimSpace(v)
	}
	return ""
}

// ── Medidas ──────────────────────────────────────────────────────────────────

// caracterizacionMedidaTipo mapea el valor de la opción del repeater RG-3 al
// tipo de salvia.emergency_measure (decisión G-9).
var caracterizacionMedidaTipo = map[string]string{
	"alojamiento_transitorio": "alojamiento",
	"alojamiento_temporal":    "alojamiento",
	"transporte":              "transporte",
	"alimentacion":            "alimentacion",
	"vestuario":               "vestuario",
	"aseo":                    models.EmergencyMeasureTypeAseo,
	"apoyo_psico":             "apoyo_psico",
	"apoyo_exequial":          models.EmergencyMeasureTypeExequial,
	"accion_simbolica":        models.EmergencyMeasureTypeAccionSimbolica,
}

// caracterizacionMedidaEtiqueta es la etiqueta legible de cada opción RG-3
// (se usa para las Notes de la medida: conserva el subtipo alojamiento).
var caracterizacionMedidaEtiqueta = map[string]string{
	"alojamiento_transitorio": "Alojamiento transitorio (hotel) hasta 15 días",
	"alojamiento_temporal":    "Alojamiento temporal (arriendo) hasta 3 meses",
	"transporte":              "Transporte",
	"alimentacion":            "Alimentación",
	"vestuario":               "Vestuario",
	"aseo":                    "Aseo",
	"apoyo_psico":             "Atención psicosocial",
	"apoyo_exequial":          "Apoyo exequial",
	"accion_simbolica":        "Acción simbólica",
}

// medidaSolicitadaCaracterizacion es un item "se solicita = sí" del repeater RG-3.
type medidaSolicitadaCaracterizacion struct {
	Tipo     string
	Etiqueta string
	Tiempo   string
	Obs      string
}

// leerMedidasSolicitadasCaracterizacion lee del submission los items del
// repeater "Medida solicitada" marcados ¿Se solicita? = true. Debe llamarse
// fuera de la transacción de cierre.
func leerMedidasSolicitadasCaracterizacion(ctx context.Context, db *gorm.DB, submissionID string) ([]medidaSolicitadaCaracterizacion, error) {
	if strings.TrimSpace(submissionID) == "" {
		return nil, nil
	}
	type row struct {
		EntryID    string `gorm:"column:entry_id"`
		QuestionID string `gorm:"column:qid"`
		Value      string `gorm:"column:value"`
	}
	var rows []row
	if err := db.WithContext(ctx).Raw(`
		SELECT a.repeater_entry_id::text AS entry_id,
		       a.question_id::text AS qid,
		       a.value AS value
		  FROM salvia.answer a
		 WHERE a.form_submission_id = ?
		   AND a.repeater_entry_id IS NOT NULL
		   AND a.deleted_at IS NULL
		   AND a.question_id::text IN (?, ?, ?, ?)
		 ORDER BY a.created_at, a.id
	`, submissionID, caracterizacionQMedida, caracterizacionQMedidaSolicit, caracterizacionQMedidaTiempo, caracterizacionQMedidaObs).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("leer medidas solicitadas: %w", err)
	}
	porEntrada := map[string]map[string]string{}
	orden := []string{}
	for _, r := range rows {
		m, ok := porEntrada[r.EntryID]
		if !ok {
			m = map[string]string{}
			porEntrada[r.EntryID] = m
			orden = append(orden, r.EntryID)
		}
		m[r.QuestionID] = strings.TrimSpace(r.Value)
	}
	return medidasSolicitadasDeEntradas(orden, porEntrada), nil
}

// medidasSolicitadasDeEntradas convierte las entradas del repeater en items
// solicitados con tipo válido. Función pura (testeable).
func medidasSolicitadasDeEntradas(orden []string, porEntrada map[string]map[string]string) []medidaSolicitadaCaracterizacion {
	out := []medidaSolicitadaCaracterizacion{}
	for _, id := range orden {
		m := porEntrada[id]
		if strings.ToLower(m[caracterizacionQMedidaSolicit]) != "true" {
			continue
		}
		valor := m[caracterizacionQMedida]
		tipo, ok := caracterizacionMedidaTipo[valor]
		if !ok {
			continue
		}
		out = append(out, medidaSolicitadaCaracterizacion{
			Tipo:     tipo,
			Etiqueta: caracterizacionMedidaEtiqueta[valor],
			Tiempo:   m[caracterizacionQMedidaTiempo],
			Obs:      m[caracterizacionQMedidaObs],
		})
	}
	return out
}

// tiposMedidasDeChips convierte fd["medidasRequeridas"] (labels de los chips)
// en códigos de tipo, ignorando los labels desconocidos.
func tiposMedidasDeChips(fd map[string]interface{}) []string {
	out := []string{}
	lista, _ := fd["medidasRequeridas"].([]interface{})
	for _, it := range lista {
		label, _ := it.(string)
		if code := models.EmergencyMeasureTypeFromLabel(label); code != "" {
			out = append(out, code)
		}
	}
	return out
}

// unirTiposMedidas une sin duplicados (chips primero, luego el repeater) y
// arma las líneas de Notes por tipo (solo vienen del repeater). Pura.
func unirTiposMedidas(chips []string, formulario []medidaSolicitadaCaracterizacion) ([]string, map[string][]string) {
	tipos := []string{}
	visto := map[string]bool{}
	agregar := func(t string) {
		if t != "" && !visto[t] {
			visto[t] = true
			tipos = append(tipos, t)
		}
	}
	for _, t := range chips {
		agregar(t)
	}
	notas := map[string][]string{}
	for _, m := range formulario {
		agregar(m.Tipo)
		linea := m.Etiqueta
		if m.Tiempo != "" {
			linea += " — tiempo: " + m.Tiempo
		}
		if m.Obs != "" {
			linea += " — obs: " + m.Obs
		}
		notas[m.Tipo] = append(notas[m.Tipo], linea)
	}
	return tipos, notas
}

const prefijoNotasCaracterizacion = "Caracterización:"

// resumenMedidasCaracterizacion cuenta lo hecho por tipo (para el timeline).
type resumenMedidasCaracterizacion struct {
	Creadas    []string
	Asociadas  []string
	Reutilizad []string
}

// asociarMedidasCaracterizacion aplica, dentro de la transacción de cierre,
// la regla reutilizar-o-crear por tipo:
//  1. medida del mismo referral_id y tipo -> se reutiliza;
//  2. medida ACTIVE del mismo caso y tipo sin remisión -> se asocia a la remisión;
//  3. si no, se crea (ACTIVE, con CaseID, FollowUpID y ReferralID).
//
// Nunca borra ni cierra medidas existentes. Notes solo se escribe si está vacío
// o ya empieza por "Caracterización:" (idempotente; no pisa notas ajenas).
func asociarMedidasCaracterizacion(tx *gorm.DB, task *models.CaseTask, refID string, tipos []string, notas map[string][]string) (resumenMedidasCaracterizacion, error) {
	var resumen resumenMedidasCaracterizacion
	type medidaRow struct {
		ID    string
		Notes *string
	}
	for _, tipo := range tipos {
		label := models.EmergencyMeasureTypeLabel(tipo)
		notaNueva := ""
		if lineas := notas[tipo]; len(lineas) > 0 {
			notaNueva = prefijoNotasCaracterizacion + " " + strings.Join(lineas, "; ")
		}

		var existentes []medidaRow
		if err := tx.Raw(`
			SELECT id::text AS id, notes
			  FROM salvia.emergency_measure
			 WHERE referral_id::text = ?
			   AND type = ?
			   AND deleted_at IS NULL
			 ORDER BY created_at
			 LIMIT 1
		`, refID, tipo).Scan(&existentes).Error; err != nil {
			return resumen, fmt.Errorf("buscar medida %s de la remisión: %w", tipo, err)
		}
		accion := "reutilizada"
		if len(existentes) == 0 {
			if err := tx.Raw(`
				SELECT id::text AS id, notes
				  FROM salvia.emergency_measure
				 WHERE case_id = ?
				   AND type = ?
				   AND status = ?
				   AND referral_id IS NULL
				   AND deleted_at IS NULL
				 ORDER BY created_at
				 LIMIT 1
			`, task.CaseID, tipo, models.EmergencyMeasureStatusActive).Scan(&existentes).Error; err != nil {
				return resumen, fmt.Errorf("buscar medida %s del caso: %w", tipo, err)
			}
			if len(existentes) > 0 {
				if err := tx.Exec(`
					UPDATE salvia.emergency_measure
					   SET referral_id = ?::uuid,
					       updated_at = NOW()
					 WHERE id::text = ?
				`, refID, existentes[0].ID).Error; err != nil {
					return resumen, fmt.Errorf("asociar medida %s a la remisión: %w", tipo, err)
				}
				accion = "asociada"
			}
		}

		if len(existentes) == 0 {
			rid := refID
			m := &models.EmergencyMeasure{
				CaseID:     task.CaseID,
				FollowUpID: task.FollowUpID,
				Type:       tipo,
				Status:     models.EmergencyMeasureStatusActive,
				ReferralID: &rid,
			}
			if notaNueva != "" {
				m.Notes = &notaNueva
			}
			if err := tx.Create(m).Error; err != nil {
				return resumen, fmt.Errorf("crear medida %s: %w", tipo, err)
			}
			resumen.Creadas = append(resumen.Creadas, label)
			continue
		}

		// Reutilizada o asociada: actualiza Notes solo si es seguro.
		if notaNueva != "" {
			actual := ""
			if existentes[0].Notes != nil {
				actual = strings.TrimSpace(*existentes[0].Notes)
			}
			if actual == "" || strings.HasPrefix(actual, prefijoNotasCaracterizacion) {
				if err := tx.Exec(`
					UPDATE salvia.emergency_measure
					   SET notes = ?,
					       updated_at = NOW()
					 WHERE id::text = ?
				`, notaNueva, existentes[0].ID).Error; err != nil {
					return resumen, fmt.Errorf("actualizar notas de la medida %s: %w", tipo, err)
				}
			}
		}
		if accion == "asociada" {
			resumen.Asociadas = append(resumen.Asociadas, label)
		} else {
			resumen.Reutilizad = append(resumen.Reutilizad, label)
		}
	}
	return resumen, nil
}

// descripcionMedidasCaracterizacion arma el texto extra del timeline.
func descripcionMedidasCaracterizacion(r resumenMedidasCaracterizacion) string {
	partes := []string{}
	if len(r.Creadas) > 0 {
		partes = append(partes, "creadas: "+strings.Join(r.Creadas, ", "))
	}
	if len(r.Asociadas) > 0 {
		partes = append(partes, "asociadas: "+strings.Join(r.Asociadas, ", "))
	}
	if len(r.Reutilizad) > 0 {
		partes = append(partes, "ya existentes: "+strings.Join(r.Reutilizad, ", "))
	}
	if len(partes) == 0 {
		return ""
	}
	return " Medidas " + strings.Join(partes, "; ") + "."
}
