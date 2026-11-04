package service

import (
	"bitsflow/internal/repository"
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"
)

// Estados públicos del reporte (T4/REP-02) — consulta sin sesión.
const (
	ReportStatusReceived   = "recibido"
	ReportStatusInProgress = "en_gestion"
	ReportStatusClosed     = "cerrado"
)

const (
	// Fallos permitidos por IP dentro de la ventana antes de bloquear (429).
	reportStatusMaxFailures   = 5
	reportStatusFailureWindow = 15 * time.Minute
	// Tope de IPs recordadas; al superarlo se barren las entradas vencidas.
	reportStatusMaxTrackedIPs = 10000

	// victim_contact_status que deja el botón «Invalidar» legado (igual que el listado).
	reportStatusLegacyInvalidated = "i"
)

var (
	// ErrReportStatusNotFound: formato inválido o radicado inexistente. Es el
	// mismo error en ambos casos para no revelar cuál fue.
	ErrReportStatusNotFound = errors.New("reporte no encontrado")

	reportStatusRadicadoFormat = regexp.MustCompile(`^[A-HJ-NP-Z2-9]{4}-[A-HJ-NP-Z2-9]{4}$`)
)

// ReportStatusRateLimitedError: la IP superó el límite de fallos. RetryAfter es
// lo que falta para que vuelva a poder consultar.
type ReportStatusRateLimitedError struct {
	RetryAfter time.Duration
}

func (e *ReportStatusRateLimitedError) Error() string {
	return fmt.Sprintf("demasiados intentos; reintentar en %s", e.RetryAfter.Round(time.Second))
}

// ReportStatusResult es TODO lo que la consulta pública devuelve: nada de
// nombres, teléfono, descripción, tipo, agente, equipo, motivo de cierre ni icode.
type ReportStatusResult struct {
	Radicado        string `json:"radicado"`
	Estado          string `json:"estado"`
	EstadoTexto     string `json:"estadoTexto"`
	Descripcion     string `json:"descripcion"`
	FechaRadicacion string `json:"fechaRadicacion"` // YYYY-MM-DD (America/Bogota)
}

// ReportStatusService resuelve «¿en qué estado va mi reporte?» por radicado.
type ReportStatusService interface {
	// Lookup consulta el estado del reporte con ese radicado en nombre de la IP ip.
	// Devuelve ErrReportStatusNotFound (sin distinguir formato inválido de
	// inexistente) o *ReportStatusRateLimitedError.
	Lookup(ctx context.Context, ip, rawRadicado string) (*ReportStatusResult, error)
}

type reportStatusService struct {
	repo    repository.ReportStatusRepository
	limiter *failureLimiter
	now     func() time.Time
}

func NewReportStatusService(repo repository.ReportStatusRepository) ReportStatusService {
	return newReportStatusServiceWithClock(repo, time.Now)
}

func newReportStatusServiceWithClock(repo repository.ReportStatusRepository, now func() time.Time) *reportStatusService {
	return &reportStatusService{
		repo:    repo,
		limiter: newFailureLimiter(reportStatusMaxFailures, reportStatusFailureWindow, reportStatusMaxTrackedIPs),
		now:     now,
	}
}

func (s *reportStatusService) Lookup(ctx context.Context, ip, rawRadicado string) (*ReportStatusResult, error) {
	now := s.now()

	// 1. Límite por IP: si ya está bloqueada no se toca la BD.
	if wait := s.limiter.retryAfter(ip, now); wait > 0 {
		return nil, &ReportStatusRateLimitedError{RetryAfter: wait}
	}

	// 2. Normalizar y validar el formato (un formato inválido ni siquiera consulta).
	radicado, validFormat := NormalizeRadicado(rawRadicado)
	var row *repository.ReportStatusRow
	if validFormat {
		var err error
		if row, err = s.repo.FindByRadicado(ctx, radicado); err != nil {
			return nil, err
		}
	}

	// 3. Formato inválido o inexistente: mismo error, y suma un fallo para la IP.
	if row == nil {
		s.limiter.recordFailure(ip, now)
		return nil, ErrReportStatusNotFound
	}

	return buildReportStatusResult(radicado, row), nil
}

// NormalizeRadicado lleva la entrada a la forma canónica «XXXX-XXXX»: sin
// espacios ni guiones, en mayúsculas. El segundo valor indica si el resultado
// cumple el formato de un radicado (alfabeto sin I, O, 0 ni 1).
func NormalizeRadicado(raw string) (string, bool) {
	var b strings.Builder
	for _, r := range strings.ToUpper(strings.TrimSpace(raw)) {
		switch r {
		case ' ', '\t', '-', '‐', '‑', '‒', '–', '—', '−':
			continue
		}
		b.WriteRune(r)
	}
	s := b.String()
	if len(s) == 8 {
		s = s[:4] + "-" + s[4:]
	}
	return s, reportStatusRadicadoFormat.MatchString(s)
}

// buildReportStatusResult traduce la fila al estado público. Reglas (ver
// «Estado derivado del reporte», S_EST, y la tabla de T4):
//   - cerrado:    con motivo de cierre, invalidado en el flujo legado, o con el caso vinculado terminado (te/cd)
//   - en gestión: sin cierre y con agente asignado, con al menos un intento, o con borrador de caso (bo/in)
//   - recibido:   el resto
func buildReportStatusResult(radicado string, row *repository.ReportStatusRow) *ReportStatusResult {
	linked := ""
	if row.LinkedStatus != nil {
		linked = *row.LinkedStatus
	}
	hasCloseReason := row.CloseReasonCode != nil && *row.CloseReasonCode != ""
	hasAgent := row.AgentID != nil && *row.AgentID != ""

	estado := ReportStatusReceived
	switch {
	case hasCloseReason || row.LegacyStatus == reportStatusLegacyInvalidated || linked == "te" || linked == "cd":
		estado = ReportStatusClosed
	case hasAgent || row.AttemptCount > 0 || linked == "bo" || linked == "in":
		estado = ReportStatusInProgress
	}

	texto, descripcion := reportStatusTexts(estado)
	return &ReportStatusResult{
		Radicado:        radicado,
		Estado:          estado,
		EstadoTexto:     texto,
		Descripcion:     descripcion,
		FechaRadicacion: row.FechaRadicacion,
	}
}

func reportStatusTexts(estado string) (string, string) {
	switch estado {
	case ReportStatusClosed:
		return "Cerrado", "La gestión de su reporte finalizó. Gracias por comunicarse con Salvia."
	case ReportStatusInProgress:
		return "En gestión", "Su reporte está siendo gestionado por el equipo de Salvia."
	default:
		return "Recibido", "Su reporte fue recibido. Un agente integral lo revisará y se comunicará con usted."
	}
}

// failureLimiter cuenta fallos por clave (IP) en una ventana deslizante, en
// memoria: se reinicia con el servicio y no se comparte entre réplicas.
type failureLimiter struct {
	mu       sync.Mutex
	failures map[string][]time.Time
	max      int
	window   time.Duration
	maxKeys  int
}

func newFailureLimiter(max int, window time.Duration, maxKeys int) *failureLimiter {
	return &failureLimiter{failures: map[string][]time.Time{}, max: max, window: window, maxKeys: maxKeys}
}

// prune descarta los fallos fuera de la ventana. Debe llamarse con l.mu tomado.
func (l *failureLimiter) prune(key string, now time.Time) []time.Time {
	list := l.failures[key]
	cut := 0
	for cut < len(list) && !list[cut].Add(l.window).After(now) {
		cut++
	}
	if cut > 0 {
		list = list[cut:]
		if len(list) == 0 {
			delete(l.failures, key)
			return nil
		}
		l.failures[key] = list
	}
	return list
}

// retryAfter devuelve 0 si la clave puede consultar, o el tiempo que falta
// para que baje del límite de fallos.
func (l *failureLimiter) retryAfter(key string, now time.Time) time.Duration {
	l.mu.Lock()
	defer l.mu.Unlock()
	list := l.prune(key, now)
	if len(list) < l.max {
		return 0
	}
	return list[len(list)-l.max].Add(l.window).Sub(now)
}

func (l *failureLimiter) recordFailure(key string, now time.Time) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if len(l.failures) >= l.maxKeys {
		for k := range l.failures {
			l.prune(k, now)
		}
	}
	l.prune(key, now)
	l.failures[key] = append(l.failures[key], now)
}
