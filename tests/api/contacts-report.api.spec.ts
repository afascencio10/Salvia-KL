import { test, expect } from '@playwright/test';
import { loginSupervisor } from '../helpers/auth';

// Pruebas de API (Playwright APIRequestContext) del endpoint
// POST /api/v1/reportes/contactos-consolidado.
// Derivan del plan: .knowledge/3-features/ReporteConsolidadoContactos/testing/reporte-contactos-test-plan.md
const ENDPOINT = '/api/v1/reportes/contactos-consolidado';

// Rango amplio con datos reales (últimos ~12 meses).
const withData = { start_date: '2025-07-21', end_date: '2026-07-21' };

test.describe('Reporte consolidado de contactos — API', () => {
  test('TC-RPC-API-08: sin sesión → 401 AUTH_TOKEN_INVALID', async ({ request }) => {
    const res = await request.post(ENDPOINT, { data: withData });
    expect(res.status()).toBe(401);
    const body = await res.json();
    expect(body.error_code).toBe('AUTH_TOKEN_INVALID');
  });

  test('TC-RPC-API-01: rango con datos → 200 y .xlsx descargable', async ({ request }) => {
    await loginSupervisor(request);
    const res = await request.post(ENDPOINT, { data: withData });
    expect(res.status()).toBe(200);
    expect(res.headers()['content-type']).toContain('spreadsheetml.sheet');
    expect(res.headers()['content-disposition']).toContain(
      `filename="reporte-contactos_${withData.start_date}_${withData.end_date}.xlsx"`,
    );
    const buf = await res.body();
    expect(buf.length, 'el .xlsx no debe estar vacío').toBeGreaterThan(1000);
    // Los .xlsx son ZIP → firma "PK".
    expect(buf.subarray(0, 2).toString('latin1')).toBe('PK');
  });

  test('TC-RPC-API-02: rango sin reportes → 400 VALIDATION_FAILED', async ({ request }) => {
    await loginSupervisor(request);
    const res = await request.post(ENDPOINT, { data: { start_date: '1990-01-01', end_date: '1990-12-31' } });
    expect(res.status()).toBe(400);
    const body = await res.json();
    expect(body.error_code).toBe('VALIDATION_FAILED');
    expect(body.message).toContain('No hay reportes');
  });

  test('TC-RPC-API-03: falta end_date → 400 VALIDATION_FAILED', async ({ request }) => {
    await loginSupervisor(request);
    const res = await request.post(ENDPOINT, { data: { start_date: '2026-01-01' } });
    expect(res.status()).toBe(400);
    expect((await res.json()).error_code).toBe('VALIDATION_FAILED');
  });

  test('TC-RPC-API-04: start_date > end_date → 400 VALIDATION_FAILED', async ({ request }) => {
    await loginSupervisor(request);
    const res = await request.post(ENDPOINT, { data: { start_date: '2026-06-30', end_date: '2026-01-01' } });
    expect(res.status()).toBe(400);
    expect((await res.json()).error_code).toBe('VALIDATION_FAILED');
  });

  test('TC-RPC-API-05: rango > 366 días → 400 VALIDATION_FAILED', async ({ request }) => {
    await loginSupervisor(request);
    const res = await request.post(ENDPOINT, { data: { start_date: '2024-01-01', end_date: '2026-01-01' } });
    expect(res.status()).toBe(400);
    expect((await res.json()).message).toContain('366');
  });

  test('TC-RPC-API-06: formato de fecha inválido → 400 VALIDATION_FAILED', async ({ request }) => {
    await loginSupervisor(request);
    const res = await request.post(ENDPOINT, { data: { start_date: '2026/13/40', end_date: '2026-06-30' } });
    expect(res.status()).toBe(400);
    expect((await res.json()).error_code).toBe('VALIDATION_FAILED');
  });
});
