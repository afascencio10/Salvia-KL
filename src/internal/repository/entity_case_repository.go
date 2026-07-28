package repository

import (
	"bitsflow/internal/models"
	"context"
	"strings"

	"gorm.io/gorm"
)

// EntityCaseRepository extiende el contrato CRUD genérico con las consultas
// específicas del componente case-entities y de la pantalla Casos Entidad.
type EntityCaseRepository interface {
	Repository[models.EntityCase]

	// FindByCaseIDWithRelations devuelve las entidades relacionadas con un caso,
	// enriquecidas con datos de entity_branch/entity, la cadena de ubicación
	// town/city/department, y los conteos de oficios y barreras activas.
	FindByCaseIDWithRelations(ctx context.Context, caseID string) ([]models.EntityCaseWithRelations, error)

	// ExistsActive indica si ya existe una relación activa (no eliminada) entre
	// el caso y la sede dados.
	ExistsActive(ctx context.Context, caseID string, entityBranchID int64) (bool, error)

	// FindByEntityIDPaginated lista entity_case de todas las sedes de una
	// organización, con filtros opcionales de documento (parcial) y ciudad del caso.
	FindByEntityIDPaginated(ctx context.Context, entityID int64, document, city string, page, pageSize int) ([]models.EntityCaseListItem, int64, error)

	// ListEntities cataloga salvia.entity ordenado por nombre.
	ListEntities(ctx context.Context) ([]models.EntityCatalogItem, error)

	// ListCitiesByEntityID ciudades distintas de las sedes de una organización.
	ListCitiesByEntityID(ctx context.Context, entityID int64) ([]models.EntityCityOption, error)
}

type entityCaseRepository struct {
	repository[models.EntityCase]
	db *gorm.DB
}

func NewEntityCaseRepository(db *gorm.DB) EntityCaseRepository {
	return &entityCaseRepository{
		repository: repository[models.EntityCase]{db: db},
		db:         db,
	}
}

func (r *entityCaseRepository) ExistsActive(ctx context.Context, caseID string, entityBranchID int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.EntityCase{}).
		Where("case_id = ? AND entity_branch_id = ?", caseID, entityBranchID).
		Count(&count).Error
	return count > 0, err
}

// entityCaseWithRelationsSQL enriquece entity_case con:
//   - entity_branch + entity (nombre, sector, dirección)
//   - town → city → department (ubicación completa)
//   - COUNT(entity_letter) por case_id + entity_branch_id → oficiosCount
//   - COUNT(barrier_v2 activas) por case_id + entity_branch_id → barrerasActivasCount
const entityCaseWithRelationsSQL = `
SELECT
    ec.id                       AS rel_id,
    ec.entity_branch_id         AS entity_branch_id,
    eb.entity_branch_i_code     AS entity_branch_icode,
    eb.entity_branch_name       AS entity_branch_name,
    e.entity_sector             AS sector,
    eb.entity_branch_address    AS address,
    COALESCE(d.department_name, '') AS department_name,
    COALESCE(c.city_name, '')       AS city_name,
    COALESCE(t.town_name, '')       AS town_name,
    ec.objetivo                 AS objetivo,
    ec.last_action              AS last_action,
    ec.created_by_id            AS created_by_id,
    ec.created_at               AS created_at,
    (SELECT COUNT(*) FROM salvia.entity_letter el
       WHERE el.case_id = ec.case_id AND el.entity_branch_id = ec.entity_branch_id AND el.deleted_at IS NULL
    ) AS oficios_count,
    (SELECT COUNT(*) FROM salvia.barrier_v2 b
       WHERE b.case_id = ec.case_id AND b.entity_branch_id = ec.entity_branch_id
         AND b.status IN ('OPEN', 'En Gestion') AND b.deleted_at IS NULL
    ) AS barreras_activas_count
FROM salvia.entity_case ec
JOIN salvia.entity_branch eb ON eb.entity_branch_id = ec.entity_branch_id
JOIN salvia.entity e ON e.entity_id = eb.entity_id
LEFT JOIN security.town t ON t.town_code = eb.entity_branch_town_code
LEFT JOIN security.city c ON c.city_id = t.city_id
LEFT JOIN security.department d ON d.department_id = c.department_id
WHERE ec.case_id = ? AND ec.deleted_at IS NULL
ORDER BY ec.created_at ASC
`

func (r *entityCaseRepository) FindByCaseIDWithRelations(ctx context.Context, caseID string) ([]models.EntityCaseWithRelations, error) {
	items := make([]models.EntityCaseWithRelations, 0)
	err := r.db.WithContext(ctx).Raw(entityCaseWithRelationsSQL, caseID).Scan(&items).Error
	return items, err
}

// riskLevelSelect mapea victim_case_form2_risk_level (1-4) a etiqueta UI.
const riskLevelSelect = `
    CASE vf2.victim_case_form2_risk_level
        WHEN 1 THEN 'Bajo'
        WHEN 2 THEN 'Medio'
        WHEN 3 THEN 'Alto'
        WHEN 4 THEN 'Extremo'
        ELSE ''
    END
`

const entityCaseListSelectSQL = `
SELECT
    ec.id AS entity_case_id,
    ec.case_id AS case_id,
    vc.victim_case_i_code AS case_code,
    ec.entity_branch_id AS entity_branch_id,
    eb.entity_branch_i_code AS entity_branch_icode,
    TRIM(CONCAT(COALESCE(vc.victim_case_victim_names, ''), ' ', COALESCE(vc.victim_case_victim_last_names, ''))) AS victim_full_name,
    COALESCE(vc.victim_case_victim_doc_number, '') AS document,
    COALESCE(case_city.city_name, '') AS city,
    (` + riskLevelSelect + `) AS risk_level,
    COALESCE(vc.victim_case_status, '') AS case_status,
    ec.objetivo AS objetivo,
    ec.last_action AS last_action,
    ec.updated_at AS updated_at,
    COALESCE(e.entity_sector, '') AS sector
FROM salvia.entity_case ec
JOIN salvia.entity_branch eb ON eb.entity_branch_id = ec.entity_branch_id
JOIN salvia.entity e ON e.entity_id = eb.entity_id
JOIN salvia.victim_case vc ON vc.victim_case_i_code = ec.case_id
LEFT JOIN salvia.victim_case_form2 vf2 ON vf2.victim_case_form2_victim_case = vc.victim_case_id
LEFT JOIN security.town case_town ON case_town.town_code = vc.victim_case_victim_town_code
LEFT JOIN security.city case_city ON case_city.city_id = case_town.city_id
WHERE ec.deleted_at IS NULL
  AND eb.entity_id = ?
`

func (r *entityCaseRepository) FindByEntityIDPaginated(ctx context.Context, entityID int64, document, city string, page, pageSize int) ([]models.EntityCaseListItem, int64, error) {
	if page < 0 {
		page = 0
	}
	if pageSize <= 0 {
		pageSize = 5
	}

	args := []interface{}{entityID}
	whereExtra := ""

	doc := strings.TrimSpace(document)
	if doc != "" {
		whereExtra += " AND vc.victim_case_victim_doc_number ILIKE ?"
		args = append(args, "%"+doc+"%")
	}

	cityFilter := strings.TrimSpace(city)
	if cityFilter != "" {
		whereExtra += " AND case_city.city_name ILIKE ?"
		args = append(args, cityFilter)
	}

	countSQL := `
SELECT COUNT(*)
FROM salvia.entity_case ec
JOIN salvia.entity_branch eb ON eb.entity_branch_id = ec.entity_branch_id
JOIN salvia.victim_case vc ON vc.victim_case_i_code = ec.case_id
LEFT JOIN security.town case_town ON case_town.town_code = vc.victim_case_victim_town_code
LEFT JOIN security.city case_city ON case_city.city_id = case_town.city_id
WHERE ec.deleted_at IS NULL
  AND eb.entity_id = ?` + whereExtra

	var total int64
	if err := r.db.WithContext(ctx).Raw(countSQL, args...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	listSQL := entityCaseListSelectSQL + whereExtra + `
ORDER BY ec.updated_at DESC
LIMIT ? OFFSET ?`
	listArgs := append(append([]interface{}{}, args...), pageSize, page*pageSize)

	items := make([]models.EntityCaseListItem, 0)
	if err := r.db.WithContext(ctx).Raw(listSQL, listArgs...).Scan(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *entityCaseRepository) ListEntities(ctx context.Context) ([]models.EntityCatalogItem, error) {
	items := make([]models.EntityCatalogItem, 0)
	err := r.db.WithContext(ctx).Raw(`
SELECT entity_id, entity_i_code, entity_name, COALESCE(entity_sector, '') AS entity_sector
FROM salvia.entity
ORDER BY entity_name ASC
`).Scan(&items).Error
	return items, err
}

func (r *entityCaseRepository) ListCitiesByEntityID(ctx context.Context, entityID int64) ([]models.EntityCityOption, error) {
	items := make([]models.EntityCityOption, 0)
	err := r.db.WithContext(ctx).Raw(`
SELECT DISTINCT c.city_id, c.city_name
FROM salvia.entity_branch eb
JOIN security.town t ON t.town_code = eb.entity_branch_town_code
JOIN security.city c ON c.city_id = t.city_id
WHERE eb.entity_id = ?
  AND c.city_name IS NOT NULL
  AND TRIM(c.city_name) <> ''
ORDER BY c.city_name ASC
`, entityID).Scan(&items).Error
	return items, err
}
