## 2026-11-02 — Regreso al reporte V2 tras iniciar sesión (WP-1 T8, REP-07)

El evento «Ingresar» lee `?next=` y, si vale exactamente `/salvia/public-v2/nuevo`, redirige ahí en vez de a la pantalla por rol; cualquier otro valor se ignora (sin redirecciones abiertas). Lo usa el paso «¿Tiene usuario?» del reporte V2 (enlace territorial, agente integral y entidad).

| Archivo | Cambio |
|---|---|
| `src/frontend/html/security/login.html` | lectura de `next` con lista blanca exacta en el éxito del login (`login.html` no figuraba en la lista de archivos del Miro: se anota para integración) |
