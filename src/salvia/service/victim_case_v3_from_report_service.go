package service

import (
	"bitsflow/common/utils"
	"bitsflow/internal/repository"
	"context"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// ReportConversionService convierte un reporte (Víctima o Tercero, con contacto
// exitoso) en un borrador de caso de Registro V3 (WP-1 T7/REP-03). El reporte no
// se cierra: queda en «Borrador iniciado» (estado derivado del caso vinculado).
type ReportConversionService interface {
	// ConvertToDraft devuelve el icode del borrador creado. Los errores de negocio son
	// los de repository.ReportConversionRepository.CheckConvertible.
	ConvertToDraft(ctx context.Context, reportICode, actorICode string) (caseICode string, err error)
}

type reportConversionService struct {
	repo   repository.ReportConversionRepository
	v3Repo repository.VictimCaseV3Repository
}

func NewReportConversionService(repo repository.ReportConversionRepository, v3Repo repository.VictimCaseV3Repository) ReportConversionService {
	return &reportConversionService{repo: repo, v3Repo: v3Repo}
}

func (s *reportConversionService) ConvertToDraft(ctx context.Context, reportICode, actorICode string) (string, error) {
	// Falla rápido (sin generar credenciales) si el reporte no se puede convertir; la
	// transacción vuelve a validar con el reporte bloqueado para cubrir un doble clic.
	if err := s.repo.CheckConvertible(ctx, reportICode); err != nil {
		return "", err
	}

	// Mismas credenciales aleatorias que CreateDraft (helpers del paquete).
	login := randomLogin(6)
	plainPassword := randomPassword(4)
	hashed, err := bcrypt.GenerateFromPassword([]byte(plainPassword), 10)
	if err != nil {
		return "", fmt.Errorf("convertir reporte: generar hash: %w", err)
	}

	caseICode, err := s.repo.ConvertReportToDraft(ctx, repository.ReportConversionInput{
		ReportICode:        reportICode,
		ActorICode:         actorICode,
		GeneralUserICode:   utils.GetUUID(),
		GeneralUserLogin:   login,
		GeneralUserPassSha: string(hashed),
	})
	if err != nil {
		return "", err
	}

	// Después del commit, como CreateDraft: si falla, el caso queda creado y se registra el error.
	if err := s.v3Repo.StoreCredentials(ctx, caseICode, login, plainPassword); err != nil {
		log.Printf("[WARN] convertir reporte %s: no se pudieron guardar las credenciales del caso %s: %v", reportICode, caseICode, err)
	}
	log.Printf("[reportConversion] borrador creado: caso=%s reporte=%s actor=%s", caseICode, reportICode, actorICode)
	return caseICode, nil
}
