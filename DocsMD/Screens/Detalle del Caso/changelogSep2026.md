# Changelog — Detalle del Caso — septiembre 2026

## 2026-09-24 — Restringir "Reabrir caso" a rol supervisor

El endpoint `POST /api/v1/casos/:id/reabrir` no validaba rol en el backend — solo el botón se ocultaba en frontend según reglas de UI, por lo que cualquier rol autenticado podía llamarlo directo. Ahora exige sesión válida y rol `sv` (403 en otro caso) vía `requireCaseDetailAccess`, el mismo helper que ya protege `case-entities`. De paso, `actor_id` del evento de timeline se toma de la sesión (`s.UserICode`) en vez de confiar en el body enviado por el cliente. En frontend, el botón "Reabrir caso" ahora depende solo de `puedeReasignar` (antes también lo veían `ro` y el `op` dueño del caso).

### Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/salvia/controller/case_detail_controller.go` | `ReabrirCaso` usa `requireCaseDetailAccess` + chequeo `CurrentRole == "sv"`; `actor_id` viene de sesión |
| `src/frontend/html/salvia/case_detail/get_case_detail_sv.html` | Botón "Reabrir caso": `v-if` reducido a `puedeReasignar && status === 'cd'` |
| `DocsMD/Screens/Detalle del Caso/detalle-caso-interface.md` | Condición del botón y tabla "Controles por rol" actualizadas |
| `DocsMD/Screens/Detalle del Caso/detalle-caso-index.md` | Descripciones de E21 y E22 actualizadas |
| `DocsMD/Screens/Detalle del Caso/Flujos/flow-E22-cuando-confirma-reapertura.md` | Nuevo PASO B0 (validación de sesión y rol) |

## 2026-09-10 — Tipo de mecanismo en el usuario mc

El tipo Municipal / Departamental / Nacional se registra 1:1 en `salvia.articulating_mechanism` al crear o editar el usuario `mc`. El modal del caso filtra por ese tipo y por la ubicación. Se retiró el input de comentarios de la card.

**Archivos**

| Archivo | Cambio |
|---|---|
| `src/internal/models/articulating_mechanism.go` | Tabla `salvia.articulating_mechanism` |
| `src/security/dao/GeneralUserDAO.go` | Columna `articulating_mechanism_id` |
| `src/frontend/html/security/general_user/set_general_user.html` | Input tipo + ubicación condicional |
| `src/frontend/js/components/case-mecanismos.js` | Filtro por tipo; sin input de comentarios |

## 2026-09-09 — Tarjeta Mecanismos en Gestión institucional

Nueva tarjeta para registrar la cadena de escalamiento del caso ante mecanismos articuladores. El modal asigna un usuario con rol `mc` filtrado por alcance y por el municipio del perfil (`town → city → department`). Cada registro admite comentarios (autor + fecha) y la creación escribe un evento en el Timeline.

**Archivos**

| Archivo | Cambio |
|---|---|
| `src/frontend/js/components/case-mecanismos.js` | Card + modal |
| `src/frontend/html/salvia/case_detail/get_case_detail_sv.html` | Montaje en el tab Gestión institucional |
| `src/internal/models/case_articulating_mechanism.go` | Tablas `case_articulating_mechanism` y comentarios |
| `src/salvia/controller/case_articulating_mechanism_controller.go` | REST |
| `src/main.go` | Seed rol `mc` + AutoMigrate + wiring |

## 2026-09-02 — Enlace Territorial puede abrir el detalle

El rol `en` obtiene `get_case_detail_sv` para que «Ver caso» (Mis barreras / Detalle de barrera) no responda 401. Las acciones de reasignar, seguimientos y similares siguen ocultas: el HTML las limita a `sv`/`op`/`ro`.

**Archivos**

| Archivo | Cambio |
|---|---|
| `src/salvia/config/Menu.go` | `"en": true` en `get_case_detail_sv` |
