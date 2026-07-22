# Pruebas Playwright (API + E2E)

Suite oficial de pruebas end-to-end de SALVIA. Herramienta: **Playwright**
(ver [testing-strategy.md](../.knowledge/5-standards/testing-strategy.md)).

```
tests/
├── api/     # Pruebas de contrato/API (Playwright APIRequestContext)
├── e2e/     # Pruebas de UI en navegador real
└── helpers/ # Utilidades compartidas (login, etc.)
```

## Requisitos

- Node LTS + `npm ci` (instala `@playwright/test`).
- Navegadores: `npx playwright install --with-deps`.
- Un servidor SALVIA accesible **con base de datos**:
  - **Local**: no hace falta levantarlo a mano — `playwright.config.ts` arranca
    `go run main.go` (dir `src/`, puerto `8090`) y lo reutiliza si ya está corriendo.
    Requiere Go y acceso a la BD de `config/db_config.json`.
  - **Contra un entorno desplegado**: exporta `BASE_URL` (p. ej.
    `export BASE_URL=https://qa.tu-host`) y no se levantará server local.

## Credenciales

Las pruebas se autentican como un supervisor de prueba (rol `sv`). La contraseña
**no** se versiona (CLAUDE.md §6): se toma de variables de entorno.

```bash
export SV_LOGIN='supervisor.prueba'   # opcional, este es el valor por defecto
export SV_PASS='********'              # obligatorio
```

En CI se inyectan como secrets del pipeline.

> El login programático omite el captcha porque el contexto usa el
> `User-Agent: flutter-client`, que el backend trata como cliente de la app.

## Ejecutar

```bash
npm test            # toda la suite (api + e2e)
npm run test:api    # solo tests/api/
npm run test:e2e    # solo tests/e2e/
npm run report      # abre el último reporte HTML
```

Artefactos (capturas, .xlsx descargados, trazas) quedan en `test-results/` y el
reporte HTML en `playwright-report/` (ambos ignorados por git).

## Cobertura actual

- `tests/api/contacts-report.api.spec.ts` — `TC-RPC-API-01/02/03/04/05/06/08`.
- `tests/e2e/contacts-report.spec.ts` — `TC-RPC-UI-01/03/04/05/06`.

Pendientes (requieren un usuario con rol distinto de `sv`/`ad`): `TC-RPC-API-07`
y `TC-RPC-UI-02`.

## Nota sobre CI

`.github/workflows/playwright.yml` corre `tests/api/` y `tests/e2e/` en cada
push/PR a `develop`/`qa`/`main`, pero **no** aprovisiona el servidor Go ni la BD:
para que el pipeline sea verde hay que apuntar `BASE_URL` a un entorno desplegado
(y definir `SV_PASS`) o extender el workflow para levantar app + BD. Hoy corre en
verde solo en local.
