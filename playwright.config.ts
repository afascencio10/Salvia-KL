import { defineConfig } from '@playwright/test';

// Config Playwright para SALVIA (API + E2E). El CI (.github/workflows/playwright.yml)
// corre `tests/api/` y `tests/e2e/` por separado.
export default defineConfig({
  testDir: './tests',
  timeout: 90_000,
  fullyParallel: false,
  workers: 1,
  retries: process.env.CI ? 1 : 0,
  reporter: [['list'], ['html', { open: 'never' }]],
  use: {
    // En local levanta el server Go (ver webServer). En CI/entorno desplegado,
    // apunta BASE_URL al host correspondiente.
    baseURL: process.env.BASE_URL || 'http://localhost:8090',
    // El backend omite la validación de captcha cuando el User-Agent es
    // 'flutter-client' (security_facades.GeneralUserLOGIN_POST), lo que permite
    // autenticar de forma programática en las pruebas.
    userAgent: 'flutter-client',
    acceptDownloads: true,
    ignoreHTTPSErrors: true,
    screenshot: 'only-on-failure',
    trace: 'on-first-retry',
  },
  // Solo en local (cuando BASE_URL no está definido): arranca el server Go y lo
  // reutiliza si ya está corriendo. Requiere Go y acceso a la BD (config/db_config.json).
  webServer: process.env.BASE_URL
    ? undefined
    : {
        command: 'go run main.go',
        cwd: 'src',
        env: { PORT: '8090' },
        url: 'http://localhost:8090',
        timeout: 120_000,
        reuseExistingServer: !process.env.CI,
        stdout: 'ignore',
        stderr: 'pipe',
      },
});
