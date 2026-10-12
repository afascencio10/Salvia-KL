package repository

import (
	"bitsflow/internal/models"
	"context"
	"strings"

	"gorm.io/gorm"
)

// DirectoryFilters agrupa los criterios opcionales del listado de entidades.
type DirectoryFilters struct {
	CityICode    string
	Type         string
	Sector       string
	DepartmentID *uint64
	NameQuery    string
	IncludeInactive bool
}

type DirectoryRepository interface {
	Repository[models.Directory]
	FindByFilters(ctx context.Context, filters DirectoryFilters, page, limit int) ([]models.Directory, int64, error)
	SetActive(ctx context.Context, id string, active bool) error
	// FindAll devuelve todos los directorios, activos e inactivos.
	FindAll(ctx context.Context) ([]models.Directory, error)
}

type directoryRepository struct {
	repository[models.Directory]
	db *gorm.DB
}

func NewDirectoryRepository(db *gorm.DB) DirectoryRepository {
	return &directoryRepository{
		repository: repository[models.Directory]{db: db},
		db:         db,
	}
}

func (r *directoryRepository) FindByFilters(ctx context.Context, f DirectoryFilters, page, limit int) ([]models.Directory, int64, error) {
	if limit <= 0 {
		limit = 20
	}
	if page < 0 {
		page = 0
	}

	query := r.db.WithContext(ctx).Model(&models.Directory{})
	if !f.IncludeInactive {
		query = query.Where("active = ?", true)
	}
	if f.CityICode != "" {
		query = query.Where("city_id = ?", f.CityICode)
	}
	if f.Type != "" {
		query = query.Where("type = ?", f.Type)
	}
	if f.Sector != "" {
		query = query.Where("sector = ?", f.Sector)
	}
	if f.DepartmentID != nil {
		cond, args := departmentFilterSQL(*f.DepartmentID)
		query = query.Where(cond, args...)
	}
	if q := strings.TrimSpace(f.NameQuery); q != "" {
		query = query.Where("name ILIKE ?", "%"+q+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []models.Directory
	err := query.
		Order("name ASC").
		Offset(page * limit).
		Limit(limit).
		Find(&items).Error
	return items, total, err
}

// departmentFilterSQL arma la condición de "filtrar por departamento" del Directorio.
//
// directories.department_id está denormalizado y llegó después: solo lo llenan la
// pantalla y la carga masiva. Las entidades que cargó SOG antes lo tienen en NULL,
// así que comparar contra la columna devolvía cero resultados para departamentos
// enteros aunque el listado sí mostrara esas entidades — la respuesta ya deduce el
// departamento del municipio cuando la columna viene vacía (ver directory_service).
// Acá se aplica el mismo criterio al filtro: la columna manda si está, y si no, el
// departamento sale de security.city a través de city_id, que guarda city_i_code.
//
// La condición va entre paréntesis a propósito: lleva un OR y se encadena con AND
// al resto de los filtros, así que sin ellos el OR se los tragaría.
func departmentFilterSQL(departmentID uint64) (string, []interface{}) {
	return "(department_id = ? OR (department_id IS NULL AND city_id IN " +
			"(SELECT city_i_code FROM security.city WHERE department_id = ?)))",
		[]interface{}{departmentID, departmentID}
}

func (r *directoryRepository) FindAll(ctx context.Context) ([]models.Directory, error) {
	var items []models.Directory
	err := r.db.WithContext(ctx).Find(&items).Error
	return items, err
}

func (r *directoryRepository) SetActive(ctx context.Context, id string, active bool) error {
	return r.db.WithContext(ctx).
		Model(&models.Directory{}).
		Where("id = ?", id).
		Update("active", active).Error
}
