import { test, expect, Page, TestInfo } from '@playwright/test';
import { loginSupervisor } from '../helpers/auth';

// Pruebas E2E/UI (Playwright, navegador real) de la descarga del reporte
// consolidado de contactos desde la pantalla "Reportes Salvia" del supervisor.
// Plan: .knowledge/3-features/ReporteConsolidadoContactos/testing/reporte-contactos-test-plan.md

// Abre la pantalla de reportes (rol sv) y espera a que Vue monte el contenido.
async function openReportsScreen(page: Page) {
  await page.goto('/salvia/casos', { waitUntil: 'networkidle' });
  await expect(page.getByText('Reportes Salvia')).toBeVisible();
}

const btn = (page: Page) =>
  page.locator('button.btn-secondary', { hasText: 'Descargar reportes' });

// Guarda una captura como artefacto en el output dir de la prueba (test-results/).
async function shot(page: Page, testInfo: TestInfo, name: string) {
  const p = testInfo.outputPath(name);
  await page.screenshot({ path: p });
  await testInfo.attach(name, { path: p, contentType: 'image/png' });
}

test.describe('Reporte consolidado de contactos — UI (rol sv)', () => {
  test.beforeEach(async ({ page }) => {
    await loginSupervisor(page.request);
    await openReportsScreen(page);
  });

  test('TC-RPC-UI-01: el supervisor ve el botón "Descargar reportes"', async ({ page }, testInfo) => {
    await expect(btn(page)).toBeVisible();
    await shot(page, testInfo, 'ui-01-boton-visible.png');
  });

  test('TC-RPC-UI-03: abrir modal, rango válido, descarga .xlsx con nombre correcto', async ({ page }, testInfo) => {
    await btn(page).click();
    await expect(page.locator('.crm-dialog')).toBeVisible();
    await expect(page.getByText('Descarga de Reportes (Excel)')).toBeVisible();

    const from = '2025-07-21';
    const to = '2026-07-21';
    await page.locator('#crm-date-from').fill(from);
    await page.locator('#crm-date-to').fill(to);
    await shot(page, testInfo, 'ui-03-modal-abierto.png');

    const downloadPromise = page.waitForEvent('download', { timeout: 60_000 });
    await page.locator('.crm-btn--primary').click();
    const download = await downloadPromise;

    expect(download.suggestedFilename()).toBe(`reporte-contactos_${from}_${to}.xlsx`);
    const dest = testInfo.outputPath(download.suggestedFilename());
    await download.saveAs(dest);
    const fs = await import('fs');
    const bytes = fs.readFileSync(dest);
    expect(bytes.length, 'el .xlsx descargado no debe estar vacío').toBeGreaterThan(1000);
    expect(bytes.subarray(0, 2).toString('latin1')).toBe('PK'); // firma ZIP
  });

  test('TC-RPC-UI-04: rango > 1 año muestra error en el DOM y no dispara la petición', async ({ page }, testInfo) => {
    await btn(page).click();
    await expect(page.locator('.crm-dialog')).toBeVisible();

    let requestFired = false;
    page.on('request', (r) => {
      if (r.url().includes('/api/v1/reportes/contactos-consolidado')) requestFired = true;
    });

    await page.locator('#crm-date-from').fill('2024-01-01');
    await page.locator('#crm-date-to').fill('2026-01-01'); // > 366 días
    await page.locator('.crm-btn--primary').click();

    const banner = page.locator('.crm-error-banner');
    await expect(banner).toBeVisible();
    await expect(banner).toContainText('1 año');
    await shot(page, testInfo, 'ui-04-error-rango.png');
    await page.waitForTimeout(500);
    expect(requestFired, 'no debe dispararse la petición al backend').toBe(false);
  });

  test('TC-RPC-UI-05: rango sin reportes muestra el mensaje del backend', async ({ page }, testInfo) => {
    await btn(page).click();
    await expect(page.locator('.crm-dialog')).toBeVisible();

    await page.locator('#crm-date-from').fill('1990-01-01');
    await page.locator('#crm-date-to').fill('1990-12-31');
    await page.locator('.crm-btn--primary').click();

    const banner = page.locator('.crm-error-banner');
    await expect(banner).toBeVisible({ timeout: 30_000 });
    await expect(banner).toContainText('No hay reportes');
    await shot(page, testInfo, 'ui-05-sin-reportes.png');
  });

  test('TC-RPC-UI-06: error 500 del backend muestra mensaje de error en el modal', async ({ page }, testInfo) => {
    await btn(page).click();
    await expect(page.locator('.crm-dialog')).toBeVisible();

    // Se intercepta la llamada y se fuerza un 500 para verificar el manejo de error en la UI.
    await page.route('**/api/v1/reportes/contactos-consolidado', (route) =>
      route.fulfill({
        status: 500,
        contentType: 'application/json',
        body: JSON.stringify({ error_code: 'SYSTEM_INTERNAL_ERROR', message: 'Error interno al generar el reporte.' }),
      }),
    );

    await page.locator('#crm-date-from').fill('2025-07-21');
    await page.locator('#crm-date-to').fill('2026-07-21');
    await page.locator('.crm-btn--primary').click();

    const banner = page.locator('.crm-error-banner');
    await expect(banner).toBeVisible({ timeout: 20_000 });
    await expect(banner).toContainText('Error interno');
    await shot(page, testInfo, 'ui-06-error-500.png');
  });
});
