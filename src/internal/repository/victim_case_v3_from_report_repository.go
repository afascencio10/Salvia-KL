package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Conversión de un reporte en un borrador de caso de Registro V3 (WP-1 T7/REP-03).
// Corte vertical NUEVO: no modifica CreateDraft ni los archivos compartidos de V3,
// solo reutiliza sus funciones del paquete (CreateDraft, markKeysAnswered,
// resolveFieldUpdates, saveEnumNFields, FieldKeyV3).

var (
	// ErrReportNotConvertibleType: solo Víctima y Tercero se convierten (Familiar y Hombre no tienen destino en V3 todavía).
	ErrReportNotConvertibleType = errors.New("Por ahora no hay procesos de atención para este tipo de reportes")
	// ErrReportNoContact: la conversión exige al menos un intento exitoso.
	ErrReportNoContact = errors.New("aún no hay contacto exitoso")
	// ErrReportHasLinkedCase: el reporte ya tiene un caso vinculado de otro flujo (no es un borrador de Registro V3).
	ErrReportHasLinkedCase = errors.New("el reporte ya tiene un caso vinculado")
)

// ErrReportHasDraft: el reporte ya tiene un borrador de caso (Registro V3 en estado bo o in) —
// se devuelve su icode para ofrecer «Ir al borrador del caso».
type ErrReportHasDraft struct{ CaseICode string }

func (e *ErrReportHasDraft) Error() string { return "el reporte ya tiene un borrador de caso" }

// ReportConversionInput datos de entrada de la conversión. Las credenciales del
// usuario del borrador las genera el servicio (igual que CreateDraft).
type ReportConversionInput struct {
	ReportICode        string
	ActorICode         string // quien convierte (queda en el hito del timeline)
	GeneralUserICode   string
	GeneralUserLogin   string
	GeneralUserPassSha string
}

// ReportConversionRepository crea el borrador de caso a partir de un reporte.
type ReportConversionRepository interface {
	// CheckConvertible valida, sin escribir, que el reporte se pueda convertir; devuelve el
	// error de negocio que corresponda (ErrVictimContactNotFound, ErrVictimContactAlreadyClosed,
	// *ErrReportHasDraft, ErrReportHasLinkedCase, ErrReportNotConvertibleType, ErrReportNoContact).
	CheckConvertible(ctx context.Context, reportICode string) error
	// ConvertReportToDraft hace la conversión en UNA transacción (con el reporte bloqueado FOR UPDATE y
	// las mismas validaciones, para que un doble clic no cree dos borradores): borrador V3 con agent_id
	// vacío, vínculo, prellenado (CI_2, CI_4, HV_1, CI_8) y hito «Borrador de caso iniciado».
	// El reporte NO se cierra. Devuelve el icode del caso.
	ConvertReportToDraft(ctx context.Context, in ReportConversionInput) (caseICode string, err error)
}

type reportConversionRepository struct {
	db *gorm.DB
}

func NewReportConversionRepository(db *gorm.DB) ReportConversionRepository {
	return &reportConversionRepository{db: db}
}

// reportConversionSource lo que se necesita del reporte para validar y prellenar.
type reportConversionSource struct {
	ID             int64  `gorm:"column:id"`
	ICode          string `gorm:"column:icode"`
	Names          string `gorm:"column:names"`
	LastNames      string `gorm:"column:last_names"`
	Closed         bool   `gorm:"column:closed"`
	ReportTypeCode string `gorm:"column:report_type_code"`
	Phone          int64  `gorm:"column:phone"`
	Facts          string `gorm:"column:facts"`
	HasContact     bool   `gorm:"column:has_contact"`
}

type reportLinkedCase struct {
	ICode  string `gorm:"column:victim_case_i_code"`
	Status string `gorm:"column:victim_case_status"`
}

const reportConversionSourceSQL = `
	SELECT vc.victim_contact_id                      AS id,
	       vc.victim_contact_i_code                  AS icode,
	       COALESCE(vc.victim_contact_names, '')     AS names,
	       COALESCE(vc.victim_contact_last_names, '') AS last_names,
	       (vc.victim_contact_close_reason IS NOT NULL OR vc.victim_contact_status = '` + victimContactLegacyInvalidated + `') AS closed,
	       COALESCE((SELECT e.victim_case_form2_enums_code
	                   FROM salvia.victim_contact_form2 vf
	                   JOIN salvia.victim_case_form2_enums e ON e.victim_case_form2_enums_id = vf.victim_contact_form2_report_type
	                  WHERE vf.victim_contact_form2_victim_contact = vc.victim_contact_id
	                    AND e.victim_case_form2_enums_category = 'victim_case_form2_report_type_v2'
	                  ORDER BY vf.victim_contact_form2_id DESC LIMIT 1), '') AS report_type_code,
	       COALESCE((SELECT vf.victim_contact_form2_victim_col_phone FROM salvia.victim_contact_form2 vf
	                  WHERE vf.victim_contact_form2_victim_contact = vc.victim_contact_id
	                  ORDER BY vf.victim_contact_form2_id DESC LIMIT 1), 0) AS phone,
	       COALESCE((SELECT vf.victim_contact_form2_facts_description FROM salvia.victim_contact_form2 vf
	                  WHERE vf.victim_contact_form2_victim_contact = vc.victim_contact_id
	                  ORDER BY vf.victim_contact_form2_id DESC LIMIT 1), '') AS facts,
	       EXISTS (SELECT 1 FROM salvia.victim_contact_attempt a
	                WHERE a.victim_contact_attempt_victim_contact = vc.victim_contact_i_code
	                  AND a.victim_contact_attempt_result = 'exitoso'
	                  AND a.victim_contact_attempt_deleted_at IS NULL) AS has_contact
	  FROM salvia.victim_contact vc
	 WHERE vc.victim_contact_i_code = ?`

// loadConversionSource lee el reporte y sus casos vinculados. Con forUpdate bloquea la fila del
// reporte (debe correr dentro de una transacción).
func loadConversionSource(tx *gorm.DB, reportICode string, forUpdate bool) (*reportConversionSource, []reportLinkedCase, error) {
	q := reportConversionSourceSQL
	if forUpdate {
		q += ` FOR UPDATE OF vc`
	}
	var rows []reportConversionSource
	if err := tx.Raw(q, reportICode).Scan(&rows).Error; err != nil {
		return nil, nil, err
	}
	if len(rows) == 0 {
		return nil, nil, ErrVictimContactNotFound
	}
	var linked []reportLinkedCase
	if err := tx.Raw(`
		SELECT victim_case_i_code, victim_case_status
		  FROM salvia.victim_case
		 WHERE victim_case_victim_contact = ?
		 ORDER BY victim_case_creation_date DESC, victim_case_id DESC`, rows[0].ID).Scan(&linked).Error; err != nil {
		return nil, nil, err
	}
	return &rows[0], linked, nil
}

// validateConversion aplica, en este orden, las reglas de «Registro de caso» (flow M-10, PASO 5).
func validateConversion(src *reportConversionSource, linked []reportLinkedCase) error {
	for _, c := range linked { // el caso vinculado más reciente de Registro V3 decide
		switch c.Status {
		case "te", "cd":
			return ErrVictimContactAlreadyClosed // el reporte ya se lee como «Cerrado — válido»
		case "bo", "in":
			return &ErrReportHasDraft{CaseICode: c.ICode}
		}
	}
	if src.Closed {
		return ErrVictimContactAlreadyClosed
	}
	if len(linked) > 0 {
		return ErrReportHasLinkedCase
	}
	if src.ReportTypeCode != "v" && src.ReportTypeCode != "r" {
		return ErrReportNotConvertibleType
	}
	if !src.HasContact {
		return ErrReportNoContact
	}
	return nil
}

func (r *reportConversionRepository) CheckConvertible(ctx context.Context, reportICode string) error {
	src, linked, err := loadConversionSource(r.db.WithContext(ctx), reportICode, false)
	if err != nil {
		return err
	}
	return validateConversion(src, linked)
}

func (r *reportConversionRepository) ConvertReportToDraft(ctx context.Context, in ReportConversionInput) (string, error) {
	var caseICode string
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		src, linked, err := loadConversionSource(tx, in.ReportICode, true)
		if err != nil {
			return err
		}
		if err := validateConversion(src, linked); err != nil {
			return err
		}

		// a. Borrador con los mismos pasos de CreateDraft (en un SAVEPOINT de esta transacción),
		// con agent_id VACÍO: el agente del equipo de seguimiento se asigna después. Si el reporte
		// trae teléfono, CreateDraft lo guarda y marca CI_4 como respondida.
		phone := ""
		if src.Phone > 0 {
			phone = strconv.FormatInt(src.Phone, 10)
		}
		caseICode, err = NewVictimCaseV3Repository(tx).CreateDraft(ctx, CreateDraftV3Input{
			Phone:              phone,
			AgentID:            "",
			GeneralUserICode:   in.GeneralUserICode,
			GeneralUserLogin:   in.GeneralUserLogin,
			GeneralUserPassSha: in.GeneralUserPassSha,
		})
		if err != nil {
			return fmt.Errorf("crear borrador: %w", err)
		}

		// b. Vincular el reporte y prellenar los nombres de la persona afectada (CI_2).
		if err := tx.Exec(`
			UPDATE salvia.victim_case
			   SET victim_case_victim_contact = ?,
			       victim_case_victim_names = ?,
			       victim_case_victim_last_names = ?,
			       victim_case_update_date = now()
			 WHERE victim_case_i_code = ?
		`, src.ID, src.Names, src.LastNames, caseICode).Error; err != nil {
			return fmt.Errorf("vincular el reporte: %w", err)
		}

		// c. Prellenar hechos (HV_1) y ajuste razonable (CI_8) y marcar CI_2, HV_1 y CI_8 como respondidas.
		answers := map[string]string{
			FieldKeyV3("CI", 2): strings.TrimSpace(src.Names + " " + src.LastNames),
			FieldKeyV3("HV", 1): src.Facts,
		}
		var adjustmentICodes []string
		if err := tx.Raw(`
			SELECT e.victim_case_form2_enums_i_code
			  FROM salvia.victim_contact_form2 vf
			  JOIN salvia.rel_victim_case_form2_enums_victim_contact_form2 rel ON rel.victim_contact_form2_id = vf.victim_contact_form2_id
			  JOIN salvia.victim_case_form2_enums e ON e.victim_case_form2_enums_id = rel.victim_case_form2_enums_id
			 WHERE vf.victim_contact_form2_victim_contact = ?
			   AND e.victim_case_form2_enums_category = 'victim_case_form2_adjustments_gbv'
		`, src.ID).Scan(&adjustmentICodes).Error; err != nil {
			return fmt.Errorf("leer ajuste razonable del reporte: %w", err)
		}
		if len(adjustmentICodes) > 0 {
			answers[FieldKeyV3("CI", 8)] = strings.Join(adjustmentICodes, ",")
		}

		if err := markKeysAnswered(tx, caseICode, answers); err != nil {
			return fmt.Errorf("marcar respondidas: %w", err)
		}
		updates, err := resolveFieldUpdates(tx, victimCaseV3FormFields, answers)
		if err != nil {
			return fmt.Errorf("resolver campos: %w", err)
		}
		var form2ID int64
		if err := tx.Raw(`
			SELECT f.victim_case_form2_id FROM salvia.victim_case_form2 f
			  JOIN salvia.victim_case c ON c.victim_case_id = f.victim_case_form2_victim_case
			 WHERE c.victim_case_i_code = ?`, caseICode).Scan(&form2ID).Error; err != nil || form2ID == 0 {
			return fmt.Errorf("resolver form2 del borrador: %v", err)
		}
		if len(updates) > 0 {
			updates["victim_case_form2_update_date"] = time.Now()
			if err := tx.Table("salvia.victim_case_form2").Where("victim_case_form2_id = ?", form2ID).Updates(updates).Error; err != nil {
				return fmt.Errorf("prellenar el formulario: %w", err)
			}
		}
		if err := saveEnumNFields(tx, form2ID, answers); err != nil {
			return fmt.Errorf("prellenar el ajuste razonable: %w", err)
		}

		// d. Hito «Borrador de caso iniciado» (lleva reporte y caso), en la misma transacción.
		now := time.Now().UTC()
		description := "Se inició el borrador del caso " + caseICode + " a partir del reporte."
		if err := NewCaseTimelineEventRepository(tx).Create(ctx,
			NewReportMilestoneEvent(in.ReportICode, ReportMilestoneDraftStarted, description, in.ActorICode, caseICode, now)); err != nil {
			return fmt.Errorf("registrar el hito: %w", err)
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return caseICode, nil
}
