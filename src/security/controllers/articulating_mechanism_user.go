package security_ctrl

import (
	"bitsflow/common/utils"
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	security_config "bitsflow/security/config"
	security_daos "bitsflow/security/dao"
	"context"
	"errors"
	"strings"
)

func validateMechanismProfile(profile *security_daos.GeneralUserProfileDTO, collectedErrors map[string]map[string]string) {
	tipo := strings.TrimSpace(strings.ToLower(profile.GeneralUserProfileMechanismType))
	if tipo != models.ArticulatingMechanismScopeMunicipal &&
		tipo != models.ArticulatingMechanismScopeDepartamental &&
		tipo != models.ArticulatingMechanismScopeNacional {
		utils.SetError(collectedErrors, security_daos.GeneralUserProfileJSONName, "mechanismType", "security_general_user_mechanism_type", "common_global_error", security_config.Locale)
		return
	}
	profile.GeneralUserProfileMechanismType = tipo

	dept := strings.TrimSpace(profile.GeneralUserProfileMechanismDepartmentId)
	city := strings.TrimSpace(profile.GeneralUserProfileMechanismCityId)
	town := strings.TrimSpace(profile.GeneralUserProfileTownCode)

	switch tipo {
	case models.ArticulatingMechanismScopeMunicipal:
		if dept == "" {
			utils.SetError(collectedErrors, security_daos.GeneralUserProfileJSONName, "mechanismDepartmentId", "security_general_user_mechanism_department", "common_global_error", security_config.Locale)
		}
		if city == "" {
			utils.SetError(collectedErrors, security_daos.GeneralUserProfileJSONName, "mechanismCityId", "security_general_user_mechanism_city", "common_global_error", security_config.Locale)
		}
		if town == "" {
			utils.SetError(collectedErrors, security_daos.GeneralUserProfileJSONName, utils.GetTag(profile, "GeneralUserProfileTownCode", "json"), "security_general_user_profile_town", "common_global_error", security_config.Locale)
		}
	case models.ArticulatingMechanismScopeDepartamental:
		if dept == "" {
			utils.SetError(collectedErrors, security_daos.GeneralUserProfileJSONName, "mechanismDepartmentId", "security_general_user_mechanism_department", "common_global_error", security_config.Locale)
		}
	}
}

func buildMechanismFromProfile(profile *security_daos.GeneralUserProfileDTO) *models.ArticulatingMechanism {
	tipo := strings.TrimSpace(strings.ToLower(profile.GeneralUserProfileMechanismType))
	row := &models.ArticulatingMechanism{Type: tipo}
	if tipo == models.ArticulatingMechanismScopeNacional {
		return row
	}
	if dept := strings.TrimSpace(profile.GeneralUserProfileMechanismDepartmentId); dept != "" {
		row.DepartmentID = &dept
	}
	if tipo == models.ArticulatingMechanismScopeMunicipal {
		if city := strings.TrimSpace(profile.GeneralUserProfileMechanismCityId); city != "" {
			row.CityID = &city
		}
		if town := strings.TrimSpace(profile.GeneralUserProfileTownCode); town != "" {
			row.TownCode = &town
		}
	}
	return row
}

func persistUserMechanism(existingID string, profile *security_daos.GeneralUserProfileDTO) (string, error) {
	store := repository.ArticulatingMechanismIdentity()
	if store == nil {
		return "", errors.New("mecanismo: store no inicializado")
	}
	ctx := context.Background()
	row := buildMechanismFromProfile(profile)
	if err := repository.FillMechanismLocationNames(ctx, store, row); err != nil {
		return "", err
	}
	if strings.TrimSpace(existingID) != "" {
		existing, err := store.FindByID(ctx, existingID)
		if err == nil && existing != nil {
			existing.Type = row.Type
			existing.DepartmentID = row.DepartmentID
			existing.CityID = row.CityID
			existing.TownCode = row.TownCode
			existing.DepartmentName = row.DepartmentName
			existing.CityName = row.CityName
			existing.TownName = row.TownName
			if err := store.Update(ctx, existing); err != nil {
				return "", err
			}
			return existing.ID, nil
		}
	}
	if err := store.Create(ctx, row); err != nil {
		return "", err
	}
	return row.ID, nil
}

func hydrateProfileFromMechanism(usr *security_daos.GeneralUserDTO) {
	id := strings.TrimSpace(usr.GeneralUserArticulatingMechanismId)
	if id == "" {
		return
	}
	store := repository.ArticulatingMechanismIdentity()
	if store == nil {
		return
	}
	row, err := store.FindByID(context.Background(), id)
	if err != nil || row == nil {
		return
	}
	p := &usr.GeneralUserGeneralUserProfile
	p.GeneralUserProfileMechanismType = row.Type
	if row.DepartmentID != nil {
		p.GeneralUserProfileMechanismDepartmentId = *row.DepartmentID
	}
	if row.CityID != nil {
		p.GeneralUserProfileMechanismCityId = *row.CityID
	}
	if row.TownCode != nil && strings.TrimSpace(p.GeneralUserProfileTownCode) == "" {
		p.GeneralUserProfileTownCode = *row.TownCode
	}
}
