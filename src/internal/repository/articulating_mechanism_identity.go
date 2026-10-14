package repository

import (
	"bitsflow/internal/models"
	"context"
	"strings"

	"gorm.io/gorm"
)

// ArticulatingMechanismIdentityStore persiste el registro 1:1 del usuario mc.
type ArticulatingMechanismIdentityStore interface {
	Create(ctx context.Context, row *models.ArticulatingMechanism) error
	Update(ctx context.Context, row *models.ArticulatingMechanism) error
	FindByID(ctx context.Context, id string) (*models.ArticulatingMechanism, error)
	ResolveDepartmentName(ctx context.Context, departmentID string) (string, error)
	ResolveCityName(ctx context.Context, cityID string) (string, error)
	ResolveTownName(ctx context.Context, townCode string) (string, error)
}

type articulatingMechanismIdentityStore struct {
	db *gorm.DB
}

var articulatingMechanismIdentity ArticulatingMechanismIdentityStore

func InitArticulatingMechanismIdentity(db *gorm.DB) {
	articulatingMechanismIdentity = &articulatingMechanismIdentityStore{db: db}
}

func ArticulatingMechanismIdentity() ArticulatingMechanismIdentityStore {
	return articulatingMechanismIdentity
}

func (s *articulatingMechanismIdentityStore) Create(ctx context.Context, row *models.ArticulatingMechanism) error {
	return s.db.WithContext(ctx).Create(row).Error
}

func (s *articulatingMechanismIdentityStore) Update(ctx context.Context, row *models.ArticulatingMechanism) error {
	return s.db.WithContext(ctx).Save(row).Error
}

func (s *articulatingMechanismIdentityStore) FindByID(ctx context.Context, id string) (*models.ArticulatingMechanism, error) {
	var row models.ArticulatingMechanism
	err := s.db.WithContext(ctx).Where("id = ?", id).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (s *articulatingMechanismIdentityStore) ResolveDepartmentName(ctx context.Context, departmentID string) (string, error) {
	var name string
	err := s.db.WithContext(ctx).
		Raw("SELECT department_name FROM security.department WHERE department_id::text = ? LIMIT 1", departmentID).
		Scan(&name).Error
	return name, err
}

func (s *articulatingMechanismIdentityStore) ResolveCityName(ctx context.Context, cityID string) (string, error) {
	var name string
	err := s.db.WithContext(ctx).
		Raw("SELECT city_name FROM security.city WHERE city_id::text = ? LIMIT 1", cityID).
		Scan(&name).Error
	return name, err
}

func (s *articulatingMechanismIdentityStore) ResolveTownName(ctx context.Context, townCode string) (string, error) {
	var name string
	err := s.db.WithContext(ctx).
		Raw("SELECT town_name FROM security.town WHERE town_code = ? LIMIT 1", townCode).
		Scan(&name).Error
	return name, err
}

func strPtrOrNil(s string) *string {
	v := strings.TrimSpace(s)
	if v == "" {
		return nil
	}
	return &v
}

// FillMechanismLocationNames resuelve y asigna los nombres de ubicación.
func FillMechanismLocationNames(ctx context.Context, store ArticulatingMechanismIdentityStore, row *models.ArticulatingMechanism) error {
	row.DepartmentName = nil
	row.CityName = nil
	row.TownName = nil
	if row.DepartmentID != nil && strings.TrimSpace(*row.DepartmentID) != "" {
		name, err := store.ResolveDepartmentName(ctx, *row.DepartmentID)
		if err != nil {
			return err
		}
		if strings.TrimSpace(name) != "" {
			row.DepartmentName = strPtrOrNil(name)
		}
	}
	if row.CityID != nil && strings.TrimSpace(*row.CityID) != "" {
		name, err := store.ResolveCityName(ctx, *row.CityID)
		if err != nil {
			return err
		}
		if strings.TrimSpace(name) != "" {
			row.CityName = strPtrOrNil(name)
		}
	}
	if row.TownCode != nil && strings.TrimSpace(*row.TownCode) != "" {
		name, err := store.ResolveTownName(ctx, *row.TownCode)
		if err != nil {
			return err
		}
		if strings.TrimSpace(name) != "" {
			row.TownName = strPtrOrNil(name)
		}
	}
	return nil
}
