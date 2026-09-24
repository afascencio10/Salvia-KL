package repository

import (
	"context"

	"gorm.io/gorm"
)

// CaseProfessionalRow es una fila de la tabla "Profesionales de atención"
// (pestaña Derivaciones del detalle del caso): una persona/usuario que ha
// intervenido en el caso, con todos sus roles concatenados.
//
// Antes traía además una nota que concatenaba todas las descripciones del
// timeline de esa persona. Se quitó junto con la columna de la tabla: era
// ilegible en una sola celda y el detalle ya está en la pestaña Timeline.
type CaseProfessionalRow struct {
	Usuario string `gorm:"column:usuario" json:"usuario"`
	Rol     string `gorm:"column:rol"     json:"rol"`
}

// CaseProfessionalRepository expone la consulta consolidada de intervinientes
// de un caso a partir de la tabla append-only salvia.case_timeline_event.
type CaseProfessionalRepository interface {
	// FindByCaseID devuelve una fila por persona que intervino en el caso,
	// ordenadas cronológicamente por su primera intervención.
	FindByCaseID(ctx context.Context, caseID string) ([]CaseProfessionalRow, error)
}

type caseProfessionalRepository struct {
	db *gorm.DB
}

func NewCaseProfessionalRepository(db *gorm.DB) CaseProfessionalRepository {
	return &caseProfessionalRepository{db: db}
}

// caseProfessionalsSQL consolida los intervinientes del caso desde
// salvia.case_timeline_event (fuente ya consolidada por diseño: todos los
// módulos registran ahí un evento cuando un actor actúa sobre el caso).
//
// Se agrupa por usuario (una fila por persona). Para evitar el producto
// cartesiano entre "eventos del usuario" y "roles del usuario", se resuelven en
// dos CTEs separadas y se unen al final:
//   - ev: agrega los eventos del usuario (nombre resuelto/snapshot y fecha de la
//     primera intervención para ordenar).
//   - roles: concatena los role_name distintos de ese usuario.
//
// Solo se consideran eventos con event_user_id no vacío (los eventos de sistema
// sin actor identificable no aportan una "persona" a la tabla). Rendimiento:
// case_timeline_event tiene índice en case_id; el filtro por caso acota el set
// antes de agregar.
const caseProfessionalsSQL = `
WITH ev AS (
	SELECT
		cte.event_user_id AS user_icode,
		COALESCE(
			NULLIF(TRIM(gup.general_user_profile_names || ' ' || gup.general_user_profile_last_names), ''),
			NULLIF(cte.actor_name, ''),
			cte.event_user_id
		) AS usuario,
		MIN(COALESCE(cte.date, cte.created_at)) AS primera_intervencion
	FROM salvia.case_timeline_event cte
	LEFT JOIN security.general_user gu
		ON gu.general_user_i_code = cte.event_user_id
	LEFT JOIN security.general_user_profile gup
		ON gup.general_user_profile_id = gu.general_user_general_user_profile
	WHERE cte.case_id = ?
	  AND cte.deleted_at IS NULL
	  AND COALESCE(NULLIF(TRIM(cte.event_user_id), ''), '') <> ''
	GROUP BY cte.event_user_id,
		COALESCE(
			NULLIF(TRIM(gup.general_user_profile_names || ' ' || gup.general_user_profile_last_names), ''),
			NULLIF(cte.actor_name, ''),
			cte.event_user_id
		)
),
roles AS (
	-- Solo se calculan los roles de los usuarios que efectivamente
	-- intervinieron en el caso (subquery contra ev), evitando recorrer todas
	-- las asignaciones de rol del sistema.
	SELECT
		gu.general_user_i_code AS user_icode,
		string_agg(DISTINCT r.role_name, ', ') AS rol
	FROM security.general_user gu
	JOIN security.rel_role_general_user rr ON rr.general_user_id = gu.general_user_id
	JOIN security.role r ON r.role_id = rr.role_id
	WHERE gu.general_user_i_code IN (SELECT user_icode FROM ev)
	GROUP BY gu.general_user_i_code
)
SELECT
	ev.usuario                           AS usuario,
	COALESCE(NULLIF(roles.rol, ''), '—') AS rol
FROM ev
LEFT JOIN roles ON roles.user_icode = ev.user_icode
ORDER BY ev.primera_intervencion ASC
`

func (r *caseProfessionalRepository) FindByCaseID(ctx context.Context, caseID string) ([]CaseProfessionalRow, error) {
	rows := make([]CaseProfessionalRow, 0)
	err := r.db.WithContext(ctx).Raw(caseProfessionalsSQL, caseID).Scan(&rows).Error
	return rows, err
}
