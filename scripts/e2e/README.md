# scripts/e2e

Scripts E2E de requerimientos concretos: cada uno prueba un requerimiento de
punta a punta (API, base de datos y pantalla) contra la app corriendo en el
sandbox cloud de Claude Code, con Playwright y `psql`. Hoy están los de WP-2
(REG-07, 08, 09, 10, 11, 12, 34 y REP-12): 92 comprobaciones en total, 93 si le
pasas `APP_LOG` (el log de la app) a REG-11.

Estaban en `cloudDevelopment/e2e/`. Salieron de ahí porque `cloudDevelopment/` es
**solo** el setup y los datos del sandbox (ver «Alcance de esta carpeta» en
[`../../cloudDevelopment/README.md`](../../cloudDevelopment/README.md)): un script
de QA no va allí, y este es su lugar.

## Cómo correrlos

Requisitos: el sandbox arriba (`./run_dev.sh`; la app en
`http://127.0.0.1:9090`, o la dirección que pongas en `BASE_URL`), las
migraciones de `src/migrations/` aplicadas (el setup las aplica solo en cada
arranque) y `psql` en el PATH. Usan el Playwright que ya trae el sandbox y los
usuarios `cloud.test.*` (clave `CloudDev!Test2026`).

```bash
node scripts/e2e/wp2-rep12-fecha-reporte.js                   # uno
for s in scripts/e2e/wp2-re*.js; do node "$s"; done           # todos los de WP-2 (~3 min)
```

Cada script imprime una línea por comprobación y termina con código 0 si todas
pasan, o 1 si alguna falla. Los datos que crean (reportes, casos, sus usuarios y
los hitos del historial) los borran al terminar; no asumen un sandbox vacío, así
que conviven con los reportes de ejemplo de `seed_test_data/`.

## Convenciones

- Un script por requerimiento, `<paquete>-<id>-<slug>.js` (p. ej.
  `wp2-reg09-hora-hecho.js`); lo que comparten va en `wp2-lib.js`.
- Los casos de QA que implementan están en `DocsMD/Screens/<pantalla>/QA/`; cada
  caso cita su script.
- Son de QA, no del sandbox: no los guardes en `cloudDevelopment/`.
