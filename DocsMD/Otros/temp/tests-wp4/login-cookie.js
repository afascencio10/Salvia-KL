// Inicia sesión como cloud.test.<rol> y escribe la cookie en un archivo para curl.
const { chromium } = require('playwright');
const fs = require('fs');
const role = process.argv[2] || 'op';
const out = process.argv[3] || `/tmp/wp4-cookie-${role}.txt`;
(async () => {
  const browser = await chromium.launch();
  const ctx = await browser.newContext();
  const page = await ctx.newPage();
  await page.route('**/seguridad/login', (route) => {
    const req = route.request();
    if (req.method() !== 'POST') return route.continue();
    return route.continue({ headers: { ...req.headers(), 'user-agent': 'flutter-client' } });
  });
  await page.goto('http://127.0.0.1:9090/seguridad/login');
  await page.locator('input#login').fill(`cloud.test.${role}`);
  await page.locator('input#pass').fill('CloudDev!Test2026');
  await page.locator('input#captcha').fill('0000');
  await Promise.all([
    page.waitForResponse((r) => r.url().endsWith('/seguridad/login') && r.request().method() === 'POST'),
    page.locator('button#login').click(),
  ]);
  await page.waitForURL((u) => !u.pathname.startsWith('/seguridad/login'), { timeout: 15000 });
  const cookies = await ctx.cookies();
  const header = cookies.map((c) => `${c.name}=${c.value}`).join('; ');
  fs.writeFileSync(out, header);
  console.log('cookie guardada en', out, '(', cookies.map((c) => c.name).join(','), ')');
  await browser.close();
})().catch((e) => { console.error(e); process.exit(1); });
