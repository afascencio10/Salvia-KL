package repository

import (
	"bitsflow/internal/models"
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// EntityObligationRow es una fila de "Mis Acciones Institucionales": el cruce de
// un caso con una acción institucional asignada a la sede del usuario, más la
// respuesta de la entidad si ya la diligenció.
type EntityObligationRow struct {
	EntityCaseID string `gorm:"column:entity_case_id" json:"entityCaseId"`

	// Datos del caso — columna "Caso" de la tabla.
	CaseID         string `gorm:"column:case_id"          json:"caseId"`
	VictimFullName string `gorm:"column:victim_full_name" json:"victimFullName"`
	Document       string `gorm:"column:document"         json:"document"`
	CityID         string `gorm:"column:city_id"          json:"cityId"`
	CityName       string `gorm:"column:city_name"        json:"cityName"`
	RiskLevel      string `gorm:"column:risk_level"       json:"riskLevel"`

	// Datos de la acción — columna "Acción Institucional".
	EntityObligationID string  `gorm:"column:entity_obligation_id" json:"entityObligationId"`
	Label              string  `gorm:"column:label"                json:"label"`
	Description        *string `gorm:"column:description"          json:"description,omitempty"`

	// Respuesta de la entidad. Nulos mientras la acción esté pendiente.
	ResponseID  *string    `gorm:"column:response_id" json:"responseId,omitempty"`
	Radicado    *string    `gorm:"column:radicado"    json:"radicado,omitempty"`
	Ampliacion  *string    `gorm:"column:ampliacion"  json:"ampliacion,omitempty"`
	Estado      *string    `gorm:"column:estado"      json:"estado,omitempty"`
	RespondedAt *time.Time `gorm:"column:responded_at" json:"respondedAt,omitempty"`
}

// EntityObligationChecklistItem es una acción del checklist "Relación con
// entidad": todas las acciones asignadas a la entidad para ese caso, marcadas
// según si la entidad ya respondió.
type EntityObligationChecklistItem struct {
	EntityCaseID       string `gorm:"column:entity_case_id"       json:"entityCaseId"`
	EntityObligationID string `gorm:"column:entity_obligation_id" json:"entityObligationId"`
	Label              string `gorm:"column:label"                json:"label"`
	Completada         bool   `gorm:"column:completada"           json:"completada"`
}

// EntityObligationCounts alimenta los contadores de las pestañas.
type EntityObligationCounts struct {
	Pendientes  int64 `json:"pendientes"`
	Completadas int64 `json:"completadas"`
}

// EntityObligationListFilter acota el listado. EntityBranchID sale siempre de la
// sesión, nunca del querystring: es lo que garantiza que la entidad solo vea sus
// propias acciones.
type EntityObligationListFilter struct {
	EntityBranchID int64
	Status         string // "pendiente" | "completada"
	Document       string
	CityID         string
	Page           int
	PageSize       int
}

// EntityCaseObligationResponseRepository maneja la gestión que la entidad reporta
// sobre sus acciones institucionales. Es la fuente única que consumen tanto la
// pantalla "Mis Acciones Institucionales" como la tarjeta de entidad en Gestión
// Institucional.
type EntityCaseObligationResponseRepository interface {
	// ListForEntityBranch pagina las acciones institucionales de la sede, según el
	// estado pedido, y devuelve además el total de esa pestaña.
	ListForEntityBranch(ctx context.Context, filter EntityObligationListFilter) ([]EntityObligationRow, int64, error)

	// CountForEntityBranch cuenta pendientes y completadas con los mismos filtros
	// de documento y ciudad, para que los contadores de las pestañas concuerden
	// con lo que el usuario está viendo.
	CountForEntityBranch(ctx context.Context, filter EntityObligationListFilter) (EntityObligationCounts, error)

	// ListCaseCitiesForEntityBranch devuelve las ciudades de los casos de la sede,
	// que son las opciones reales del filtro de ciudad. No se usan las ciudades de
	// las sedes de la organización: el filtro es sobre la ubicación de la víctima.
	ListCaseCitiesForEntityBranch(ctx context.Context, entityBranchID int64) ([]models.EntityCityOption, error)

	// FindChecklistsByEntityCaseIDs devuelve, para los entity_case indicados, todas
	// sus acciones institucionales marcadas según si la entidad ya respondió. Es lo
	// que alimenta el panel "Relación con entidad" del listado: se pide una sola vez
	// para toda la página, no una consulta por fila.
	FindChecklistsByEntityCaseIDs(ctx context.Context, entityCaseIDs []string) ([]EntityObligationChecklistItem, error)

	// BelongsToEntityBranch indica si esa acción de ese caso está realmente
	// asignada a la sede indicada. Es el guard de escritura: sin él, una entidad
	// podría responder la acción de otra enviando ids ajenos al endpoint.
	BelongsToEntityBranch(ctx context.Context, entityBranchID int64, entityCaseID, entityObligationID string) (bool, error)

	// Upsert crea la respuesta o actualiza la existente. La entidad puede corregir
	// lo que reportó, y el índice único (entity_case_id, entity_obligation_id)
	// garantiza una sola fila por acción.
	Upsert(ctx context.Context, response *models.EntityCaseObligationResponse) error

	// FindByEntityCaseIDs devuelve las respuestas de varios entity_case en una
	// sola consulta, para no hacer una por fila del listado.
	FindByEntityCaseIDs(ctx context.Context, entityCaseIDs []string) ([]models.EntityCaseObligationResponse, error)

	// FindByCaseID devuelve las respuestas de todas las entidades de un caso —
	// lo que necesita la tarjeta de entidad en Gestión Institucional.
	FindByCaseID(ctx context.Context, caseID string) ([]models.EntityCaseObligationResponse, error)
}

type entityCaseObligationResponseRepository struct {
	db *gorm.DB
}

func NewEntityCaseObligationResponseRepository(db *gorm.DB) EntityCaseObligationResponseRepository {
	return &entityCaseObligationResponseRepository{db: db}
}

// entityObligationFromSQL es el FROM compartido por el listado, los contadores y
// las opciones de ciudad, para que las tres vean exactamente el mismo universo.
//
// Los joins del caso (víctima, documento, ciudad, nivel de riesgo) son los mismos
// que usa entityCaseListSelectSQL en entity_case_repository.go — misma pantalla de
// entidad, mismos datos del caso.
//
// Las acciones asignadas salen de entity_case_obligation, que es append-only: una
// misma acción aparece una vez por cada seguimiento que tocó el checklist, así que
// se agrupa con DISTINCT. Se descartan las filas sin entity_obligation_id, que son
// las acciones "Otra" escritas a mano y no tienen catálogo al que responder.
const entityObligationFromSQL = `
FROM salvia.entity_case ec
JOIN salvia.victim_case vc ON vc.victim_case_i_code = ec.case_id
LEFT JOIN salvia.victim_case_form2 vf2 ON vf2.victim_case_form2_victim_case = vc.victim_case_id
LEFT JOIN security.town case_town ON case_town.town_code = vc.victim_case_victim_town_code
LEFT JOIN security.city case_city ON case_city.city_id = case_town.city_id
JOIN (
    SELECT DISTINCT entity_case_id, entity_obligation_id
    FROM salvia.entity_case_obligation
    WHERE entity_obligation_id IS NOT NULL
) acc ON acc.entity_case_id = ec.id
JOIN salvia.entity_obligation eo ON eo.id = acc.entity_obligation_id AND eo.deleted_at IS NULL
LEFT JOIN salvia.entity_case_obligation_response resp
    ON resp.entity_case_id = acc.entity_case_id
   AND resp.entity_obligation_id = acc.entity_obligation_id
WHERE ec.deleted_at IS NULL
  AND ec.entity_branch_id = ?
`

// appendEntityObligationFilters agrega los filtros opcionales al WHERE y sus
// argumentos en el mismo orden, que es lo que exige el placeholder posicional.
func appendEntityObligationFilters(sql string, args []interface{}, filter EntityObligationListFilter) (string, []interface{}) {
	if doc := filter.Document; doc != "" {
		sql += " AND COALESCE(vc.victim_case_victim_doc_number, '') ILIKE ?"
		args = append(args, "%"+doc+"%")
	}
	if filter.CityID != "" {
		sql += " AND case_city.city_id = ?"
		args = append(args, filter.CityID)
	}
	return sql, args
}

// statusCondition traduce la pestaña a una condición sobre la respuesta: una
// acción está completada cuando la entidad ya la diligenció, y pendiente cuando
// todavía no existe respuesta.
func statusCondition(status string) string {
	if status == models.EntityCaseObligationStatusCompletada {
		return " AND resp.id IS NOT NULL"
	}
	return " AND resp.id IS NULL"
}

func (r *entityCaseObligationResponseRepository) ListForEntityBranch(ctx context.Context, filter EntityObligationListFilter) ([]EntityObligationRow, int64, error) {
	rows := make([]EntityObligationRow, 0)
	if filter.EntityBranchID <= 0 {
		return rows, 0, nil
	}
	pageSize := filter.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	page := filter.Page
	if page < 0 {
		page = 0
	}

	where, args := appendEntityObligationFilters(entityObligationFromSQL, []interface{}{filter.EntityBranchID}, filter)
	where += statusCondition(filter.Status)

	var total int64
	if err := r.db.WithContext(ctx).Raw("SELECT COUNT(*) "+where, args...).Scan(&total).Error; err != nil {
		return rows, 0, err
	}

	const selectSQL = `
SELECT
    ec.id                        AS entity_case_id,
    ec.case_id                   AS case_id,
    TRIM(CONCAT(COALESCE(vc.victim_case_victim_names, ''), ' ', COALESCE(vc.victim_case_victim_last_names, ''))) AS victim_full_name,
    COALESCE(vc.victim_case_victim_doc_number, '') AS document,
    COALESCE(case_city.city_id, '')   AS city_id,
    COALESCE(case_city.city_name, '') AS city_name,
    (` + riskLevelSelect + `)         AS risk_level,
    eo.id                        AS entity_obligation_id,
    eo.label                     AS label,
    eo.description               AS description,
    resp.id                      AS response_id,
    resp.radicado                AS radicado,
    resp.ampliacion              AS ampliacion,
    resp.estado                  AS estado,
    resp.updated_at              AS responded_at
`
	// Orden estable: primero los casos tocados más recientemente y, dentro de cada
	// caso, las acciones en el orden del catálogo. Se agrega el id como último
	// criterio para que la paginación no repita ni se salte filas cuando dos
	// acciones empatan.
	paged := selectSQL + where + ` ORDER BY ec.updated_at DESC, eo."order" ASC, eo.id ASC LIMIT ? OFFSET ?`
	pagedArgs := append(append([]interface{}{}, args...), pageSize, page*pageSize)

	err := r.db.WithContext(ctx).Raw(paged, pagedArgs...).Scan(&rows).Error
	return rows, total, err
}

func (r *entityCaseObligationResponseRepository) CountForEntityBranch(ctx context.Context, filter EntityObligationListFilter) (EntityObligationCounts, error) {
	var counts EntityObligationCounts
	if filter.EntityBranchID <= 0 {
		return counts, nil
	}
	where, args := appendEntityObligationFilters(entityObligationFromSQL, []interface{}{filter.EntityBranchID}, filter)

	// Una sola pasada para las dos pestañas.
	err := r.db.WithContext(ctx).Raw(`
SELECT
    COUNT(*) FILTER (WHERE resp.id IS NULL)     AS pendientes,
    COUNT(*) FILTER (WHERE resp.id IS NOT NULL) AS completadas
`+where, args...).Scan(&counts).Error
	return counts, err
}

func (r *entityCaseObligationResponseRepository) ListCaseCitiesForEntityBranch(ctx context.Context, entityBranchID int64) ([]models.EntityCityOption, error) {
	items := make([]models.EntityCityOption, 0)
	if entityBranchID <= 0 {
		return items, nil
	}
	err := r.db.WithContext(ctx).Raw(`
SELECT DISTINCT case_city.city_id AS city_id, case_city.city_name AS city_name
`+entityObligationFromSQL+`
  AND case_city.city_id IS NOT NULL
ORDER BY case_city.city_name ASC
`, entityBranchID).Scan(&items).Error
	return items, err
}

func (r *entityCaseObligationResponseRepository) FindChecklistsByEntityCaseIDs(ctx context.Context, entityCaseIDs []string) ([]EntityObligationChecklistItem, error) {
	items := make([]EntityObligationChecklistItem, 0)
	if len(entityCaseIDs) == 0 {
		return items, nil
	}
	// Mismo criterio que el listado: DISTINCT sobre entity_case_obligation porque es
	// append-only, y solo acciones del catálogo. "Completada" es que exista respuesta
	// de la entidad, igual que en las pestañas, para que el conteo del panel y el de
	// las pestañas nunca digan cosas distintas.
	err := r.db.WithContext(ctx).Raw(`
SELECT
    acc.entity_case_id                   AS entity_case_id,
    eo.id                                AS entity_obligation_id,
    eo.label                             AS label,
    (resp.id IS NOT NULL)                AS completada
FROM (
    SELECT DISTINCT entity_case_id, entity_obligation_id
    FROM salvia.entity_case_obligation
    WHERE entity_obligation_id IS NOT NULL
      AND entity_case_id IN ?
) acc
JOIN salvia.entity_obligation eo ON eo.id = acc.entity_obligation_id AND eo.deleted_at IS NULL
LEFT JOIN salvia.entity_case_obligation_response resp
    ON resp.entity_case_id = acc.entity_case_id
   AND resp.entity_obligation_id = acc.entity_obligation_id
ORDER BY acc.entity_case_id, eo."order" ASC, eo.id ASC
`, entityCaseIDs).Scan(&items).Error
	return items, err
}

func (r *entityCaseObligationResponseRepository) BelongsToEntityBranch(ctx context.Context, entityBranchID int64, entityCaseID, entityObligationID string) (bool, error) {
	if entityBranchID <= 0 || entityCaseID == "" || entityObligationID == "" {
		return false, nil
	}
	var count int64
	err := r.db.WithContext(ctx).Raw(`
SELECT COUNT(*)
FROM salvia.entity_case ec
JOIN salvia.entity_case_obligation eco
    ON eco.entity_case_id = ec.id
   AND eco.entity_obligation_id = ?
WHERE ec.deleted_at IS NULL
  AND ec.id = ?
  AND ec.entity_branch_id = ?
`, entityObligationID, entityCaseID, entityBranchID).Scan(&count).Error
	return count > 0, err
}

func (r *entityCaseObligationResponseRepository) Upsert(ctx context.Context, response *models.EntityCaseObligationResponse) error {
	// ON CONFLICT sobre el índice único: si la entidad ya había respondido esta
	// acción, se actualizan los campos diligenciados y queda constancia de quién
	// hizo la última corrección. created_at no se toca.
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "entity_case_id"}, {Name: "entity_obligation_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"radicado", "ampliacion", "estado", "responded_by_id", "updated_at"}),
	}).Create(response).Error
}

func (r *entityCaseObligationResponseRepository) FindByEntityCaseIDs(ctx context.Context, entityCaseIDs []string) ([]models.EntityCaseObligationResponse, error) {
	items := make([]models.EntityCaseObligationResponse, 0)
	if len(entityCaseIDs) == 0 {
		return items, nil
	}
	err := r.db.WithContext(ctx).
		Where("entity_case_id IN ?", entityCaseIDs).
		Find(&items).Error
	return items, err
}

func (r *entityCaseObligationResponseRepository) FindByCaseID(ctx context.Context, caseID string) ([]models.EntityCaseObligationResponse, error) {
	items := make([]models.EntityCaseObligationResponse, 0)
	err := r.db.WithContext(ctx).Raw(`
SELECT ecor.*
FROM salvia.entity_case_obligation_response ecor
JOIN salvia.entity_case ec ON ec.id = ecor.entity_case_id
WHERE ec.case_id = ? AND ec.deleted_at IS NULL
ORDER BY ecor.updated_at DESC
`, caseID).Scan(&items).Error
	return items, err
}
