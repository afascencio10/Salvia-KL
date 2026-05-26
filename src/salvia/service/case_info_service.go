// Package service — case_info_service.go
// Servicio para el componente <case-info>: transforma datos crudos en la estructura de respuesta.
package service

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	salvia_config "bitsflow/salvia/config"
	"context"
	"fmt"
)

type CaseInfoService interface {
	GetFullCaseInfo(ctx context.Context, caseICode string) (*models.CaseFullInfo, error)
}

type caseInfoService struct {
	repo repository.CaseInfoRepository
}

func NewCaseInfoService(repo repository.CaseInfoRepository) CaseInfoService {
	return &caseInfoService{repo: repo}
}

func (s *caseInfoService) GetFullCaseInfo(ctx context.Context, caseICode string) (*models.CaseFullInfo, error) {
	if caseICode == "" {
		return nil, fmt.Errorf("caseICode vacío")
	}

	raw, err := s.repo.GetFullInfoByICode(ctx, caseICode)
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, fmt.Errorf("caso no encontrado")
	}

	// Edad: prioridad form2 (calculada desde birth_date), fallback form1
	var edad *int64
	if raw.F2Age != nil && *raw.F2Age > 0 {
		edad = raw.F2Age
	} else if raw.F1Age != nil {
		v := int64(*raw.F1Age)
		edad = &v
	}

	// Nivel de riesgo texto
	nivelTexto := ""
	if raw.F2RiskLevel != nil {
		switch *raw.F2RiskLevel {
		case 1:
			nivelTexto = "Bajo"
		case 2:
			nivelTexto = "Moderado"
		case 3:
			nivelTexto = "Alto"
		case 4:
			nivelTexto = "Extremo"
		}
	} else if raw.F1FemicideRisk == "1" {
		nivelTexto = "Alto (feminicidio)"
	}

	// Resolver enums a texto legible usando el mapa de Locale
	locale := salvia_config.Locale["sp"]
	resolveEnum := func(key string) string {
		if key == "" {
			return ""
		}
		if translated, ok := locale[key]; ok {
			return translated
		}
		return key
	}

	genderIdentity := resolveEnum(raw.F2GenderIdentity)
	sexualOrientation := resolveEnum(raw.F2SexualOrientation)
	contactKinship := resolveEnum(raw.F2ContactKinship)
	maritalStatus := resolveEnum(raw.F2MaritalStatus)
	disability := resolveEnum(raw.F2Disability)
	ethnicAffiliation := resolveEnum(raw.F2EthnicAffiliation)
	indigenousPeople := resolveEnum(raw.F2IndigenousPeople)
	campesino := resolveEnum(raw.F2Campesino)
	scenarioViolence := resolveEnum(raw.F2ScenarioViolence)
	relationshipAggressor := resolveEnum(raw.F2RelationshipAggressor)

	requireInterpreter := resolveEnum(raw.F2RequireInterpreter)
	proximityAggressor := resolveEnum(raw.F2ProximityAggressor)
	aggressorGender := resolveEnum(raw.F2AggressorGender)
	aggressorDocType := resolveEnum(raw.F2AggressorDocType)

	result := &models.CaseFullInfo{
		Encabezado: models.CaseInfoEncabezado{
			Funcionarios:      raw.OwnerDescription,
			FechaCreacion:     raw.CreationDate,
			FechaModificacion: raw.UpdateDate,
			Estado:            raw.Status,
		},
		Victima: models.CaseInfoVictima{
			Nombres:           raw.Nombres,
			Apellidos:         raw.Apellidos,
			NombreIdentitario: raw.F2IdentityName,
			Edad:              edad,
			FechaNacimiento:   raw.F2BirthDate,
			TipoDocumento:     raw.DocType,
			NumeroDocumento:   raw.DocNumber,
			Nacionalidad:      raw.F1Nationality,
			OtraNacionalidad:  raw.F1NationalityOther,
			Municipio:         coalesce(raw.CityName, raw.TownCode),
			CorreoElectronico: raw.F1Email,
			DireccionResidencia: raw.F2ResidenceAddress,
			AjusteRazonable:   requireInterpreter,
		},
		DatosPersonales: models.CaseInfoDatosPersonales{
			CondicionMigratoria:   raw.F1ForeignerStatus,
			Telefono:              raw.F2Phone,
			Genero:                raw.F1Gender,
			IdentidadGenero:       genderIdentity,
			OtraIdentidadGenero:   raw.F1GenderIdentityOther,
			OrientacionSexual:     sexualOrientation,
			OtraOrientacionSexual: raw.F1SexualOrientationOther,
		},
		Etnicos: models.CaseInfoEtnicos{
			GrupoEtnico:      ethnicAffiliation,
			Afrodescendiente: coalesce(raw.F1IfAfro),
			Indigena:         coalesce(indigenousPeople, raw.F1IfIndigenous),
			LenguaIndigena:   raw.F1IfIndigenousTongue,
			Campesino:        coalesce(campesino, raw.F1IfPeasant),
			VictimaConflicto: raw.F1IfArmedConflict,
		},
		Contacto: models.CaseInfoContacto{
			NombreContacto:   raw.F2ContactNames,
			TelefonoContacto: raw.F2ContactPhone,
			Parentesco:       contactKinship,
			PersonasCargo:    raw.F1Dependents,
			NumeroHijos:      raw.F1ChildrenNumber,
			EstadoCivil:      maritalStatus,
			Discapacidad:     disability,
		},
		Ubicacion: models.CaseInfoUbicacion{
			Departamento: raw.DeptName,
			Ciudad:       raw.CityName,
			Municipio:    raw.TownName,
		},
		Hechos: models.CaseInfoHechos{
			Descripcion:        raw.F2FactsDescription,
			FechaHechos:        raw.F2FactsDate,
			Horario:            raw.F2FactsStartTime,
			EscenarioViolencia: coalesce(scenarioViolence, raw.F1ViolenceScene),
			RiesgoFeminicida:   raw.F1FemicideRisk,
			DireccionHechos:    raw.F2FactsAddress,
		},
		Agresor: models.CaseInfoAgresor{
			TipoAgresor:     raw.F1Aggressor,
			Relacion:        coalesce(relationshipAggressor, raw.F1RelationshipAggressor),
			Nombre:          coalesce(raw.F2AggressorNames, raw.F1AggressorName),
			TipoDocumento:   coalesce(aggressorDocType, raw.F1AggressorDocType),
			NumeroDocumento: coalesce(raw.F2AggressorDocNumber, raw.F1AggressorDocNumber),
			Direccion:       coalesce(raw.F2AggressorAddress, raw.F1AggressorAddress),
			Telefono:        coalesce(raw.F2AggressorPhone, raw.F1AggressorPhone),
			NumAgresores:    raw.F2NumAggressors,
			Proximidad:      proximityAggressor,
			GeneroAgresor:   aggressorGender,
		},
		Riesgo: models.CaseInfoRiesgo{
			NivelRiesgo:               raw.F2RiskLevel,
			NivelRiesgoTexto:          nivelTexto,
			AmenazasMuerte:            raw.F1DeathThreats,
			AgresorTieneArmas:         raw.F1AggressorHasWeapons,
			ViolenciaPrevia:           raw.F1ViolenceBefore,
			ViolenciaFisicaIncremento: raw.F1PhysicalIncreased,
			SeparacionUltimoAnio:      raw.F1SeparatedLastYear,
			AmenazoConArma:            raw.F1ThreatenedWeapon,
			AmenazoHijos:              raw.F1ThreatenedChildren,
			CelosoViolento:            raw.F1JealousViolent,
			CreeCapazMatar:            raw.F1CapableKilling,
			RiesgoInminente:           raw.F1ImminentRisk,
			DenunciaPrevia:            raw.F1PreviouslyReported,
			SiDenuncioAntes:           raw.F1IfPreviouslyReported,
		},
	}

	return result, nil
}

// coalesce retorna el primer string no vacío.
func coalesce(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
