package service

import (
	"bitsflow/internal/repository"
	"context"
	"errors"
	"testing"
	"time"
)

// fakeReportStatusRepo devuelve filas por radicado canónico y cuenta consultas.
type fakeReportStatusRepo struct {
	rows  map[string]*repository.ReportStatusRow
	calls int
	err   error
}

func (f *fakeReportStatusRepo) FindByRadicado(_ context.Context, radicado string) (*repository.ReportStatusRow, error) {
	f.calls++
	if f.err != nil {
		return nil, f.err
	}
	return f.rows[radicado], nil
}

func strp(s string) *string { return &s }

func TestNormalizeRadicado(t *testing.T) {
	cases := []struct {
		in    string
		want  string
		valid bool
	}{
		{"K7QM-4R2X", "K7QM-4R2X", true},
		{"k7qm-4r2x", "K7QM-4R2X", true},
		{"k7qm4r2x", "K7QM-4R2X", true},
		{"  k7qm 4r2x ", "K7QM-4R2X", true},
		{"K7QM–4R2X", "K7QM-4R2X", true}, // guion largo (en dash)
		{"K7QM-4R2", "K7QM4R2", false},   // corto
		{"K7QM-4R2X9", "K7QM4R2X9", false},
		{"IOQM-4R2X", "IOQM-4R2X", false}, // I y O no existen en el alfabeto
		{"K0QM-1R2X", "K0QM-1R2X", false}, // 0 y 1 no existen
		{"' OR 1=1 --", "'OR1=1", false},
		{"", "", false},
	}
	for _, c := range cases {
		got, ok := NormalizeRadicado(c.in)
		if got != c.want || ok != c.valid {
			t.Errorf("NormalizeRadicado(%q) = (%q, %v); esperaba (%q, %v)", c.in, got, ok, c.want, c.valid)
		}
	}
}

func TestBuildReportStatusResult(t *testing.T) {
	cases := []struct {
		name string
		row  repository.ReportStatusRow
		want string
	}{
		{"recién creado", repository.ReportStatusRow{LegacyStatus: "v"}, ReportStatusReceived},
		{"agente asignado sin intentos", repository.ReportStatusRow{LegacyStatus: "v", AgentID: strp("agent-1")}, ReportStatusInProgress},
		{"agente vacío cuenta como sin agente", repository.ReportStatusRow{LegacyStatus: "v", AgentID: strp("")}, ReportStatusReceived},
		{"un intento", repository.ReportStatusRow{LegacyStatus: "v", AttemptCount: 1}, ReportStatusInProgress},
		{"borrador de caso bo", repository.ReportStatusRow{LegacyStatus: "v", LinkedStatus: strp("bo")}, ReportStatusInProgress},
		{"borrador de caso in", repository.ReportStatusRow{LegacyStatus: "v", AgentID: strp("a"), LinkedStatus: strp("in")}, ReportStatusInProgress},
		{"caso terminado te", repository.ReportStatusRow{LegacyStatus: "v", LinkedStatus: strp("te")}, ReportStatusClosed},
		{"caso terminado cd", repository.ReportStatusRow{LegacyStatus: "v", AgentID: strp("a"), AttemptCount: 3, LinkedStatus: strp("cd")}, ReportStatusClosed},
		{"cerrado con motivo ap", repository.ReportStatusRow{LegacyStatus: "v", CloseReasonCode: strp("ap")}, ReportStatusClosed},
		{"cerrado con motivo ni", repository.ReportStatusRow{LegacyStatus: "v", CloseReasonCode: strp("ni"), AgentID: strp("a")}, ReportStatusClosed},
		{"legado invalidado", repository.ReportStatusRow{LegacyStatus: "i"}, ReportStatusClosed},
		{"cierre gana sobre borrador", repository.ReportStatusRow{LegacyStatus: "v", CloseReasonCode: strp("co"), LinkedStatus: strp("bo")}, ReportStatusClosed},
		{"otro estado de caso no cambia nada", repository.ReportStatusRow{LegacyStatus: "v", LinkedStatus: strp("ra")}, ReportStatusReceived},
	}
	for _, c := range cases {
		row := c.row
		row.FechaRadicacion = "2026-11-02"
		got := buildReportStatusResult("K7QM-4R2X", &row)
		if got.Estado != c.want {
			t.Errorf("%s: estado = %s; esperaba %s", c.name, got.Estado, c.want)
		}
		if got.Radicado != "K7QM-4R2X" || got.FechaRadicacion != "2026-11-02" || got.EstadoTexto == "" || got.Descripcion == "" {
			t.Errorf("%s: resultado incompleto: %+v", c.name, got)
		}
	}
}

func newTestReportStatusService(repo *fakeReportStatusRepo, clock *time.Time) *reportStatusService {
	return newReportStatusServiceWithClock(repo, func() time.Time { return *clock })
}

func TestLookupFoundAndNotFoundAreIndistinguishable(t *testing.T) {
	repo := &fakeReportStatusRepo{rows: map[string]*repository.ReportStatusRow{
		"K7QM-4R2X": {FechaRadicacion: "2026-11-02", LegacyStatus: "v"},
	}}
	now := time.Date(2026, 11, 2, 10, 0, 0, 0, time.UTC)
	svc := newTestReportStatusService(repo, &now)

	res, err := svc.Lookup(context.Background(), "1.1.1.1", " k7qm4r2x ")
	if err != nil || res.Estado != ReportStatusReceived {
		t.Fatalf("esperaba recibido, obtuve %+v, %v", res, err)
	}

	_, errInexistente := svc.Lookup(context.Background(), "1.1.1.1", "ZZZZ-ZZZZ")
	_, errFormato := svc.Lookup(context.Background(), "1.1.1.1", "hola")
	if !errors.Is(errInexistente, ErrReportStatusNotFound) || !errors.Is(errFormato, ErrReportStatusNotFound) {
		t.Fatalf("ambos deben dar ErrReportStatusNotFound: %v / %v", errInexistente, errFormato)
	}
	if repo.calls != 2 { // el formato inválido no consulta la BD
		t.Fatalf("esperaba 2 consultas a la BD (válida e inexistente), hubo %d", repo.calls)
	}
}

// C8/C9: 5 fallos bloquean al sexto sin consultar la BD; los aciertos no cuentan;
// pasada la ventana vuelve a funcionar; cada IP tiene su propio contador.
func TestLookupRateLimitByIP(t *testing.T) {
	repo := &fakeReportStatusRepo{rows: map[string]*repository.ReportStatusRow{
		"K7QM-4R2X": {FechaRadicacion: "2026-11-02", LegacyStatus: "v"},
	}}
	now := time.Date(2026, 11, 2, 10, 0, 0, 0, time.UTC)
	svc := newTestReportStatusService(repo, &now)
	ctx := context.Background()

	for i := 0; i < 20; i++ { // C9: consultas válidas repetidas no bloquean
		if _, err := svc.Lookup(ctx, "9.9.9.9", "K7QM-4R2X"); err != nil {
			t.Fatalf("un acierto no debe fallar ni sumar fallos (intento %d): %v", i, err)
		}
	}

	for i := 0; i < 5; i++ {
		now = now.Add(time.Minute)
		if _, err := svc.Lookup(ctx, "9.9.9.9", "ZZZZ-ZZZZ"); !errors.Is(err, ErrReportStatusNotFound) {
			t.Fatalf("fallo %d: esperaba no encontrado, obtuve %v", i+1, err)
		}
	}
	callsBefore := repo.calls
	_, err := svc.Lookup(ctx, "9.9.9.9", "K7QM-4R2X") // sexto intento, aunque sea válido
	var limited *ReportStatusRateLimitedError
	if !errors.As(err, &limited) {
		t.Fatalf("el sexto intento debe quedar limitado, obtuve %v", err)
	}
	if repo.calls != callsBefore {
		t.Fatalf("limitada no debe consultar la BD")
	}
	// El primer fallo fue a las 10:01 y la ventana es de 15 min: a las 10:05 faltan 11 min.
	if limited.RetryAfter != 11*time.Minute {
		t.Fatalf("RetryAfter = %v; esperaba 11m", limited.RetryAfter)
	}

	if _, err := svc.Lookup(ctx, "8.8.8.8", "K7QM-4R2X"); err != nil { // otra IP no se ve afectada
		t.Fatalf("otra IP debe poder consultar: %v", err)
	}

	now = now.Add(limited.RetryAfter + time.Second) // vence el primer fallo
	if _, err := svc.Lookup(ctx, "9.9.9.9", "K7QM-4R2X"); err != nil {
		t.Fatalf("pasada la ventana debe volver a funcionar: %v", err)
	}
}

func TestLookupPropagatesRepositoryErrorWithoutCountingFailure(t *testing.T) {
	boom := errors.New("db caída")
	repo := &fakeReportStatusRepo{err: boom}
	now := time.Now()
	svc := newTestReportStatusService(repo, &now)
	for i := 0; i < 10; i++ {
		if _, err := svc.Lookup(context.Background(), "1.1.1.1", "K7QM-4R2X"); !errors.Is(err, boom) {
			t.Fatalf("esperaba el error de la BD, obtuve %v", err)
		}
	}
	repo.err = nil
	repo.rows = map[string]*repository.ReportStatusRow{"K7QM-4R2X": {LegacyStatus: "v"}}
	if _, err := svc.Lookup(context.Background(), "1.1.1.1", "K7QM-4R2X"); err != nil {
		t.Fatalf("un error de BD no debe sumar fallos: %v", err)
	}
}

func TestFailureLimiterSweepsExpiredKeys(t *testing.T) {
	l := newFailureLimiter(5, time.Minute, 3)
	t0 := time.Now()
	for _, ip := range []string{"a", "b", "c"} {
		l.recordFailure(ip, t0)
	}
	l.recordFailure("d", t0.Add(2*time.Minute)) // lleva a maxKeys: barre las vencidas
	if len(l.failures) != 1 {
		t.Fatalf("esperaba 1 clave tras el barrido, hay %d", len(l.failures))
	}
}
