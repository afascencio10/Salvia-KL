// Package service — case_info_service.go
// Servicio para el componente <case-info>: transforma datos crudos en la estructura de respuesta.
package service

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	common_config "bitsflow/common/config"
	salvia_config "bitsflow/salvia/config"
	"context"
	"fmt"
	"strings"
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
			AgenteAsignado:    raw.AgentName,
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
			TipoDocumento:     resolveForm1Code(raw.DocType, common_config.DOCUMENT_TYPE),
			NumeroDocumento:   raw.DocNumber,
			Nacionalidad:      resolveForm1Code(raw.F1Nationality, salvia_config.VICTIM_CASE_VICTIM_NATIONALITY),
			OtraNacionalidad:  raw.F1NationalityOther,
			Municipio:         coalesce(raw.CityName, raw.TownCode),
			CorreoElectronico: raw.F1Email,
			DireccionResidencia: raw.F2ResidenceAddress,
			AjusteRazonable:   requireInterpreter,
		},
		DatosPersonales: models.CaseInfoDatosPersonales{
			CondicionMigratoria:   resolveForm1Code(raw.F1ForeignerStatus, salvia_config.VICTIM_CASE_VICTIM_FOREIGNER_IMMIGRATION_STATUS),
			Telefono:              raw.F2Phone,
			Genero:                resolveForm1Code(raw.F1Gender, salvia_config.VICTIM_CASE_VICTIM_GENDER),
			IdentidadGenero:       genderIdentity,
			OtraIdentidadGenero:   raw.F1GenderIdentityOther,
			OrientacionSexual:     sexualOrientation,
			OtraOrientacionSexual: raw.F1SexualOrientationOther,
		},
		Etnicos: models.CaseInfoEtnicos{
			GrupoEtnico:      ethnicAffiliation,
			Afrodescendiente: resolveForm1Code(coalesce(raw.F1IfAfro), nil),
			Indigena:         coalesce(indigenousPeople, resolveForm1Code(raw.F1IfIndigenous, nil)),
			LenguaIndigena:   raw.F1IfIndigenousTongue,
			Campesino:        coalesce(campesino, resolveForm1Code(raw.F1IfPeasant, nil)),
			VictimaConflicto: resolveForm1Code(raw.F1IfArmedConflict, nil),
		},
		Contacto: models.CaseInfoContacto{
			NombreContacto:   raw.F2ContactNames,
			TelefonoContacto: raw.F2ContactPhone,
			Parentesco:       contactKinship,
			PersonasCargo:    resolveForm1Code(raw.F1Dependents, salvia_config.VICTIM_CASE_VICTIM_DEPENDENTS),
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
			Descripcion:           coalesce(raw.F2FactsDescription, raw.F1FactsDescription),
			Ocurrencia:            resolveForm1Code(raw.F1FactsOccurrence, salvia_config.OCCURRENCE),
			Horario:               coalesce(raw.F2FactsStartTime, formatTimeRange(raw.F1FactsStartTime, raw.F1FactsEndTime)),
			DiaSemana:             resolveForm1Code(raw.F1FactsWeekday, salvia_config.WEEK_DAY),
			FechaHechos:           coalesce(raw.F2FactsDate, raw.F1FactsDate),
			ViolenciaExperimentada: resolveForm1Code(raw.F1ViolenceExperienced, salvia_config.VIOLENCE_EXPERIENCED),
			OtroTipoViolencia:     raw.F1ViolenceExperiencedOther,
			AmbitoViolencia:       resolveForm1Code(raw.F1ViolenceScope, salvia_config.VIOLENCE_SCOPE),
			EscenarioViolencia:    coalesce(scenarioViolence, resolveForm1Code(raw.F1ViolenceScene, salvia_config.VIOLENCE_SCENES["sp"])),
			RiesgoFeminicida:      resolveForm1Code(raw.F1FemicideRisk, nil),
			DireccionHechos:       raw.F2FactsAddress,
		},
		Agresor: models.CaseInfoAgresor{
			TipoAgresor:     resolveForm1Code(raw.F1Aggressor, salvia_config.AGGRESSOR),
			Relacion:        coalesce(relationshipAggressor, resolveForm1Code(raw.F1RelationshipAggressor, salvia_config.RELATIONSHIP_WITH_AGGRESSOR)),
			Nombre:          coalesce(raw.F2AggressorNames, raw.F1AggressorName),
			TipoDocumento:   coalesce(aggressorDocType, resolveForm1Code(raw.F1AggressorDocType, common_config.DOCUMENT_TYPE)),
			NumeroDocumento: coalesce(raw.F2AggressorDocNumber, raw.F1AggressorDocNumber),
			Direccion:       coalesce(raw.F2AggressorAddress, raw.F1AggressorAddress),
			Telefono:        coalesce(raw.F2AggressorPhone, raw.F1AggressorPhone),
			NumAgresores:    resolveEnum(raw.F2NumAggressors),
			Proximidad:      proximityAggressor,
			GeneroAgresor:   aggressorGender,
		},
		Riesgo: models.CaseInfoRiesgo{
			NivelRiesgo:               raw.F2RiskLevel,
			NivelRiesgoTexto:          nivelTexto,
			AmenazasMuerte:            resolveForm1Code(raw.F1DeathThreats, nil),
			AgresorTieneArmas:         resolveForm1Code(raw.F1AggressorHasWeapons, nil),
			ViolenciaPrevia:           resolveForm1Code(raw.F1ViolenceBefore, nil),
			ViolenciaFisicaIncremento: resolveForm1Code(raw.F1PhysicalIncreased, nil),
			SeparacionUltimoAnio:      resolveForm1Code(raw.F1SeparatedLastYear, nil),
			AmenazoConArma:            resolveForm1Code(raw.F1ThreatenedWeapon, nil),
			AmenazoHijos:              resolveForm1Code(raw.F1ThreatenedChildren, nil),
			CelosoViolento:            resolveForm1Code(raw.F1JealousViolent, nil),
			CreeCapazMatar:            resolveForm1Code(raw.F1CapableKilling, nil),
			RiesgoInminente:           resolveForm1Code(raw.F1ImminentRisk, nil),
			DenunciaPrevia:            resolveForm1Code(raw.F1PreviouslyReported, nil),
			SiDenuncioAntes:           raw.F1IfPreviouslyReported,
		},
	}

	// Cargar plan de atención (enums multi-select)
	planAtencion, _ := s.repo.GetPlanAtencionByICode(ctx, caseICode)
	if len(planAtencion) > 0 {
		for i, p := range planAtencion {
			planAtencion[i] = resolveEnum(p)
		}
		result.Hechos.PlanAtencion = planAtencion
	}

	return result, nil
}

// resolveForm1Code traduce un código corto de form1 usando los mapas legacy.
func resolveForm1Code(code string, m map[string]string) string {
	if code == "" {
		return ""
	}
	if m != nil {
		if translated, ok := m[code]; ok {
			return translated
		}
	}
	// Intentar sí/no genérico
	switch code {
	case "y", "1", "true", "Si", "si":
		return "Sí"
	case "n", "0", "false", "No", "no":
		return "No"
	}
	return code
}

// formatTimeRange forma un rango de horas "HH:MM - HH:MM" si ambas existen.
func formatTimeRange(start, end string) string {
	start = strings.TrimSpace(start)
	end = strings.TrimSpace(end)
	if start == "" && end == "" {
		return ""
	}
	// Truncar a HH:MM si viene HH:MM:SS
	if len(start) > 5 {
		start = start[:5]
	}
	if len(end) > 5 {
		end = end[:5]
	}
	if start != "" && end != "" {
		return start + " - " + end
	}
	return coalesce(start, end)
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
