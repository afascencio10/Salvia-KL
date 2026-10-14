package service

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"errors"
	"strings"
	"unicode/utf8"
)

var ErrInvalidAttemptResult = errors.New("resultado de intento inválido — use 'fallido' o 'exitoso'")
var ErrInvalidCloseReason = errors.New("motivo de cierre inválido")
var ErrObservationRequired = errors.New("la observación es obligatoria para este motivo")
var ErrObservationTooLong = errors.New("la observación no puede superar 500 caracteres")
var ErrCloseReasonDeferred = errors.New("el motivo de cierre 'ap' (apertura de caso) no está disponible en esta fase")

// maxObservationLength — victim_contact_close_observation es varchar(500); sin esta validación un
// texto más largo llegaba a la BD y respondía 500 en vez de 400.
const maxObservationLength = 500

// VictimContactsManagementService lógica de negocio del modal gestion-reporte-modal
// (M-01 detalle, M-05/M-06 intentos, M-02/M-07/M-08/M-09 cierre).
type VictimContactsManagementService interface {
	GetDetail(ctx context.Context, icode string, canManage bool) (*models.ReportDetail, error)
	RegisterAttempt(ctx context.Context, icode, result, agentICode string) (intentosHoy int, diasCompletos int, err error)
	Close(ctx context.Context, icode, closeReasonCode, observation string) (category string, err error)
}

type victimContactsManagementService struct {
	repo repository.VictimContactsManagementRepository
}

func NewVictimContactsManagementService(repo repository.VictimContactsManagementRepository) VictimContactsManagementService {
	return &victimContactsManagementService{repo: repo}
}

// GetDetail delega al repositorio y sobreescribe canManage con el valor
// resuelto por el controller a partir del rol de sesión (GAP-2: op y sv
// pueden gestionar cualquier reporte, confirmado 2026-10-11).
func (s *victimContactsManagementService) GetDetail(ctx context.Context, icode string, canManage bool) (*models.ReportDetail, error) {
	detail, err := s.repo.GetDetail(ctx, icode)
	if err != nil {
		return nil, err
	}
	// La BD guarda claves de enum; la API devuelve texto para mostrar.
	detail.ReportType = reportTypeLabel(detail.ReportType)
	detail.Parentesco = localizeEnumName(detail.Parentesco)
	detail.AjusteRazonable = localizeEnumNames(detail.AjusteRazonable)
	detail.CanManage = canManage && detail.CloseReasonCode == ""
	return detail, nil
}

var validCloseReasonCodes = map[string]bool{"ni": true, "nc": true, "co": true, "dp": true, "ap": true}

// deferredCloseReasonCodes — motivos que existen en el catálogo (migración) pero que NINGÚN endpoint puede
// escribir todavía: 'ap' (Apertura de caso) depende de M-10, diferido. Sin este bloqueo, una llamada directa a
// POST /cerrar dejaba el reporte "Cerrado - válido" sin abrir ningún caso (confirmado con el usuario, 2026-10-13).
// Cuando se construya M-10, quitar 'ap' de este mapa.
var deferredCloseReasonCodes = map[string]bool{"ap": true}

// validCloseReasonsRequiringObservation — 'co' (no competencia) y 'dp'
// (decisión de la persona) exigen observación obligatoria (M-08/M-09).
var closeReasonsRequiringObservation = map[string]bool{"co": true, "dp": true}

func (s *victimContactsManagementService) RegisterAttempt(ctx context.Context, icode, result, agentICode string) (int, int, error) {
	result = strings.TrimSpace(strings.ToLower(result))
	if result != "fallido" && result != "exitoso" {
		return 0, 0, ErrInvalidAttemptResult
	}
	return s.repo.InsertAttempt(ctx, icode, result, agentICode)
}

func (s *victimContactsManagementService) Close(ctx context.Context, icode, closeReasonCode, observation string) (string, error) {
	code := strings.TrimSpace(strings.ToLower(closeReasonCode))
	if !validCloseReasonCodes[code] {
		return "", ErrInvalidCloseReason
	}
	if deferredCloseReasonCodes[code] {
		return "", ErrCloseReasonDeferred
	}
	if closeReasonsRequiringObservation[code] && strings.TrimSpace(observation) == "" {
		return "", ErrObservationRequired
	}
	if utf8.RuneCountInString(strings.TrimSpace(observation)) > maxObservationLength {
		return "", ErrObservationTooLong
	}
	return s.repo.Close(ctx, icode, code, observation)
}
