## 2026-09-27 — Implementación de la pantalla (PSWYA-755)

Se construyó `/salvia/public-v2` como segunda versión del formulario público de reporte, en ruta y tablas propias, sin tocar la versión actual. Wizard de 4 pasos con 4 tipos de reportante (antes 2). QA con Playwright y `curl` contra el backend real encontró y corrigió 7 bugs: pantalla invisible por `id="app"` duplicado, decodificación JSON rota por embedding de structs, dos variantes del mismo bug de Locale mal referenciado (captcha y validaciones de campo requerido), una validación de longitud siempre falsa, permisos de BD faltantes para el rol de runtime, y un popup de éxito compartido que bloqueaba un botón. Migración corrida contra la BD real.

### Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/frontend/html/salvia/victim_contact/set_victim_contact_v2.html` | Nuevo — wizard de 4 pasos |
| `src/salvia/facades/VictimContactFacadeV2.go` | Nuevo — handlers GET/POST públicos |
| `src/salvia/controllers/VictimContactControllerV2.go` | Nuevo — validación y orquestación de persistencia |
| `src/salvia/dao/VictimContactForm2V2DAO.go` | Nuevo — DTO y persistencia de la tabla de extensión |
| `src/migrations/add_victim_contact_v2.sql` | Nuevo — tabla `victim_contact_form2_v2` + 9 filas de enum + grants (no está en git) |
| `src/salvia/config/Locale.go` | Solo adiciones — labels de los enums/campos nuevos de V2 |
| `src/salvia/config/Enums.go` | Solo adiciones — template y form paths de V2 |
| `src/salvia/facades/MainRouter.go` | Solo adiciones — registro de rutas `/salvia/public-v2` |
