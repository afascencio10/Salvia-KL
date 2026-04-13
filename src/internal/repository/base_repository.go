// Package repository provee una capa de acceso a datos genérica basada en GORM.
// Coexiste con el patrón DAO existente (pgx/v4) sin modificarlo.
package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

// PageResult encapsula el resultado de una consulta paginada.
type PageResult[T any] struct {
	Items    []T
	Total    int64
	Page     int
	PageSize int
}

// Repository define el contrato CRUD + paginación para cualquier entidad T.
type Repository[T any] interface {
	Create(ctx context.Context, entity *T) error
	Update(ctx context.Context, entity *T) error
	UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*T, error)
	FindWithPagination(ctx context.Context, page, pageSize int) (PageResult[T], error)
}

// repository es la implementación concreta de Repository[T] usando *gorm.DB.
type repository[T any] struct {
	db *gorm.DB
}

// NewRepository construye un Repository[T] listo para usar.
// db debe ser una instancia válida de *gorm.DB; panics si es nil.
func NewRepository[T any](db *gorm.DB) Repository[T] {
	if db == nil {
		panic("repository: *gorm.DB no puede ser nil")
	}
	return &repository[T]{db: db}
}

// Create inserta entity en la base de datos.
// GORM escanea el ID generado por PostgreSQL (gen_random_uuid()) de vuelta en entity.
func (r *repository[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

// Update guarda todos los campos de entity (debe tener PK asignada).
func (r *repository[T]) Update(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

// UpdateFields actualiza únicamente los campos especificados en el mapa.
// Usa Updates con map para evitar problemas con zero-values y generar
// un UPDATE granular: UPDATE ... SET campo1=v1, campo2=v2 WHERE id=?
// Retorna gorm.ErrRecordNotFound si el registro no existe.
func (r *repository[T]) UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error {
	if len(fields) == 0 {
		return nil
	}
	result := r.db.WithContext(ctx).
		Model(new(T)).
		Where("id = ?", id).
		Updates(fields)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Delete realiza un soft-delete del registro con el id dado.
// Retorna gorm.ErrRecordNotFound si el registro no existe.
func (r *repository[T]) Delete(ctx context.Context, id string) error {
	var entity T
	result := r.db.WithContext(ctx).First(&entity, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	return r.db.WithContext(ctx).Delete(&entity).Error
}

// FindByID busca un registro por su PK (campo id).
// Retorna gorm.ErrRecordNotFound si no existe o fue soft-deleted.
func (r *repository[T]) FindByID(ctx context.Context, id string) (*T, error) {
	var entity T
	result := r.db.WithContext(ctx).First(&entity, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, result.Error
	}
	return &entity, nil
}

// FindWithPagination retorna una página de resultados ordenados por created_at ASC.
// page es base-0 (primera página = 0), pageSize debe ser > 0.
func (r *repository[T]) FindWithPagination(ctx context.Context, page, pageSize int) (PageResult[T], error) {
	if pageSize <= 0 {
		pageSize = 20
	}
	offset := page * pageSize

	var total int64
	var items []T

	if err := r.db.WithContext(ctx).Model(new(T)).Count(&total).Error; err != nil {
		return PageResult[T]{}, err
	}

	if err := r.db.WithContext(ctx).
		Order("created_at ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&items).Error; err != nil {
		return PageResult[T]{}, err
	}

	return PageResult[T]{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}
