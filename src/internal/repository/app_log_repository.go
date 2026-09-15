package repository

import (
	"bitsflow/internal/models"
	"context"
	"encoding/json"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// AppLogRepository persiste logs estructurados de funciones del backend.
// Es append-only: solo se insertan registros, nunca se modifican ni eliminan.
type AppLogRepository interface {
	// Log inserta un registro con una key para filtrar (ej. i_code de un caso),
	// un mensaje en lenguaje humano y datos adicionales que se serializan a JSON.
	// data puede ser nil si el registro no necesita payload.
	Log(ctx context.Context, key, message string, data interface{}) error
}

type appLogRepository struct {
	db *gorm.DB
}

func NewAppLogRepository(db *gorm.DB) AppLogRepository {
	return &appLogRepository{db: db}
}

func (r *appLogRepository) Log(ctx context.Context, key, message string, data interface{}) error {
	entry := &models.AppLog{Key: key, Message: message}
	if data != nil {
		raw, err := json.Marshal(data)
		if err != nil {
			return err
		}
		entry.Data = datatypes.JSON(raw)
	}
	return r.db.WithContext(ctx).Create(entry).Error
}
