package repository

import (
	"bitsflow/internal/models"
	"context"

	"gorm.io/gorm"
)

type AgentLightRepository interface {
	FindByICode(ctx context.Context, iCode string) (*models.AgentLight, error)
}

type agentLightRepository struct {
	db *gorm.DB
}

func NewAgentLightRepository(db *gorm.DB) AgentLightRepository {
	return &agentLightRepository{db: db}
}

func (r *agentLightRepository) FindByICode(ctx context.Context, iCode string) (*models.AgentLight, error) {
	var agent models.AgentLight
	err := r.db.WithContext(ctx).
		Table("security.general_user").
		Select("security.general_user.general_user_i_code, security.general_user_profile.general_user_profile_names, security.general_user_profile.general_user_profile_last_names").
		Joins("JOIN security.general_user_profile ON security.general_user.general_user_general_user_profile = security.general_user_profile.general_user_profile_id").
		Where("security.general_user.general_user_i_code = ?", iCode).
		Take(&agent).Error

	if err != nil {
		return nil, err
	}
	return &agent, nil
}
