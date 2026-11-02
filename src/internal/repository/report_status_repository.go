package repository

import (
	"context"

	"gorm.io/gorm"
)

// ReportStatusRow — lo mínimo que la consulta pública necesita de un reporte
// (T4/REP-02). No incluye nombres, teléfono, descripción, tipo, agente (solo si
// existe), equipo, motivo de cierre (solo si existe) ni icode: el servicio
// los reduce a recibido / en gestión / cerrado.
type ReportStatusRow struct {
	// Día de radicación en hora de Colombia (YYYY-MM-DD). victim_contact_creation_date
	// se guarda con la hora local del servidor, que main.go fija en America/Bogota.
	FechaRadicacion string  `gorm:"column:fecha_radicacion"`
	CloseReasonCode *string `gorm:"column:close_reason_code"`
	LegacyStatus    string  `gorm:"column:legacy_status"`
	AgentID         *string `gorm:"column:agent_id"`
	AttemptCount    int     `gorm:"column:attempt_count"`
	// Estado del caso vinculado más reciente de Registro V3 (bo, in, te o cd);
	// nil si no hay (un caso de otro flujo no cuenta). Ver «Estado derivado del reporte».
	LinkedStatus *string `gorm:"column:linked_status"`
}

// ReportStatusRepository consulta el estado de un reporte por su radicado.
type ReportStatusRepository interface {
	// FindByRadicado devuelve la fila del reporte con ese radicado exacto, o
	// nil (sin error) si no existe.
	FindByRadicado(ctx context.Context, radicado string) (*ReportStatusRow, error)
}

type reportStatusRepository struct {
	db *gorm.DB
}

func NewReportStatusRepository(db *gorm.DB) ReportStatusRepository {
	return &reportStatusRepository{db: db}
}

// reportStatusFindSQL: SELECT parametrizado (el radicado nunca se concatena).
// Los intentos cuentan cualquier resultado no eliminado; el caso vinculado se
// limita a los estados de Registro V3.
const reportStatusFindSQL = `
SELECT TO_CHAR(vc.victim_contact_creation_date, 'YYYY-MM-DD') AS fecha_radicacion,
       closeEnum.victim_case_form2_enums_code                AS close_reason_code,
       vc.victim_contact_status                              AS legacy_status,
       vc.victim_contact_agent_id                            AS agent_id,
       (SELECT COUNT(*)::int
          FROM salvia.victim_contact_attempt vca
         WHERE vca.victim_contact_attempt_victim_contact = vc.victim_contact_i_code
           AND vca.victim_contact_attempt_deleted_at IS NULL) AS attempt_count,
       (SELECT c.victim_case_status
          FROM salvia.victim_case c
         WHERE c.victim_case_victim_contact = vc.victim_contact_id
           AND c.victim_case_status IN ('bo', 'in', 'te', 'cd')
         ORDER BY c.victim_case_creation_date DESC, c.victim_case_id DESC
         LIMIT 1)                                            AS linked_status
  FROM salvia.victim_contact vc
  LEFT JOIN salvia.victim_case_form2_enums closeEnum
         ON closeEnum.victim_case_form2_enums_id = vc.victim_contact_close_reason
 WHERE vc.victim_contact_radicado = ?
 LIMIT 1`

func (r *reportStatusRepository) FindByRadicado(ctx context.Context, radicado string) (*ReportStatusRow, error) {
	var rows []ReportStatusRow
	if err := r.db.WithContext(ctx).Raw(reportStatusFindSQL, radicado).Scan(&rows).Error; err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, nil
	}
	return &rows[0], nil
}
