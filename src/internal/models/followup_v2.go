// Package models contiene los structs GORM para la nueva capa de repositorio.
// Estos modelos son independientes de los DTOs del patrón DAO existente.
package models

import (
	"time"

	"gorm.io/gorm"
)

// FollowUpV2 representa la tabla follow_up_v2 en el esquema public (o salvia).
// Usa soft-delete mediante gorm.DeletedAt.
//
// Migración SQL equivalente:
//
//	CREATE TABLE follow_up_v2 (
//	    id          VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid(),
//	    case_id     VARCHAR(36) NOT NULL,
//	    agent_id    VARCHAR(36) NOT NULL,
//	    status      VARCHAR(20) NOT NULL,
//	    created_at  TIMESTAMPTZ,
//	    updated_at  TIMESTAMPTZ,
//	    deleted_at  TIMESTAMPTZ
//	);
//	CREATE INDEX idx_follow_up_v2_case_id ON follow_up_v2(case_id);
//	CREATE INDEX idx_follow_up_v2_deleted_at ON follow_up_v2(deleted_at);
type FollowUpV2 struct {
	ID        string         `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()"`
	CaseID    string         `gorm:"type:varchar(36);index;not null"`
	AgentID   string         `gorm:"type:varchar(36);not null"`
	Status    string         `gorm:"type:varchar(20);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
