# Pruebas WP-4 (API + BD) — scripts de apoyo

Scripts bash/curl para repetir las pruebas de API+BD de WP-4 contra la app local (`./run_dev.sh`, `:9090`).
Son material temporal de la sesión (este directorio está en `DocsMD/Otros/temp/`, ignorado por git: añadir con `git add -f`).

```bash
cd DocsMD/Otros/temp/tests-wp4
node login-cookie.js op                 # inicia sesión como cloud.test.op y guarda /tmp/wp4-cookie-op.txt (Playwright)
bash restart-app.sh                     # reinicia la app (necesario tras cambios en Go o frontend)
bash t9-api-test.sh                     # borrador JSON, cascada, materialización, tareas/oficios, solo lectura (63 casos)
bash t9b-api-test.sh                    # remisiones a otros equipos y casos límite (48 casos)
```
- `wp4-lib.sh`: helpers (`api METHOD ruta [json]`, `SQL "…"`, `check`, `new_case`, `mark_triada`, `ic categoría código`).
- La tríada se completa marcando las claves como respondidas por SQL (`mark_triada`) y guardando un campo de `AG` (el guardado real exige ~35 campos).
