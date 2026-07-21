import { APIRequestContext, expect } from '@playwright/test';

// Credenciales del supervisor de prueba. Nunca se hardcodea la contraseña en el
// repo (CLAUDE.md §6 / data-classification): se toma de la variable de entorno
// SV_PASS. En CI se inyecta como secret; en local: `export SV_PASS='...'`.
export const SV_LOGIN = process.env.SV_LOGIN || 'supervisor.prueba';
export const SV_PASS = process.env.SV_PASS || '';

export function requireCreds(): void {
  if (!SV_PASS) {
    throw new Error(
      'Falta SV_PASS. Exporta la contraseña del supervisor de prueba antes de correr las pruebas:\n' +
        "  export SV_PASS='...'\n" +
        'Ver tests/README.md',
    );
  }
}

// Autentica como supervisor (rol sv) sobre el APIRequestContext dado.
// - En pruebas de API: pasar el fixture `request`.
// - En pruebas E2E: pasar `page.request` (comparte cookies con las navegaciones).
export async function loginSupervisor(request: APIRequestContext): Promise<{ login: string; role: string }> {
  requireCreds();
  const res = await request.post('/seguridad/login', {
    headers: { 'Content-Type': 'application/json' },
    data: { user: { login: SV_LOGIN, pass: SV_PASS } },
  });
  expect(res.status(), 'login del supervisor de prueba debe responder 200').toBe(200);
  const body = await res.json();
  expect(body.role, 'el rol devuelto debe ser sv').toBe('sv');
  return body;
}
