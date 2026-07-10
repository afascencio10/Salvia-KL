---
okf_version: "1.0"
type: UI_Screen
title: "Pantalla: Remisión Temporal (host del card 3x3)"
description: "Pantalla puente temporal en /salvia/remision-temporal/:id (id = psychosocial_support_id) a la que se llega desde 'Ver remisión' en Mis Remisiones Psicosocial. Su única función es montar el card psychosocial-contact-card (que compone el flujo de modales) mientras no exista la pantalla interna definitiva de Remisión."
owner: "@frontend-squad"
status: active
tags: [frontend, ui, psicosocial, 3x3, temporal]
dependencies:
  - Component: 3-features/AtencionPsicosocial3x3/frontend/components/PsychosocialContactCard/index.md
  - Endpoint: 3-features/AtencionPsicosocial3x3/backend/endpoints/list-contact-attempts.md
code_refs:
  - "src/salvia/facades/RemisionTemporalFacade.go"
  - "src/salvia/facades/MainRouter.go"
  - "src/salvia/config/Enums.go"
  - "src/frontend/html/salvia/remision-temporal/remision_temporal.html"
  - "src/frontend/html/salvia/mis-remisiones-psicosocial/mis_remisiones_psicosocial.html"
last_updated: "2026-07-09"
---

# Pantalla: Remisión Temporal (host del card 3x3)

Pantalla **puente temporal**. No existe todavía la pantalla interna definitiva de Remisión Psicosocial; mientras tanto, "Ver remisión" en **Mis Remisiones Psicosocial** lleva a esta pantalla, cuya única responsabilidad es **montar el card** [`psychosocial-contact-card`](../../components/PsychosocialContactCard/index.md) (que a su vez compone el flujo de modales [`psychosocial-contact-modal`](../../components/PsychosocialContactModal/index.md)).

## Ruta y montaje

- **Ruta**: `GET /salvia/remision-temporal/:id` — `:id` = **`psychosocial_support_id`**.
- **Registro**: `secRouter.GET("/remision-temporal/:id", RemisionTemporalGET)` en `src/salvia/facades/MainRouter.go` (mismo patrón que `hacer-seguimiento/:id`).
- **Facade**: `RemisionTemporalGET` en `src/salvia/facades/RemisionTemporalFacade.go`. Valida sesión con `CheckAndGetSession(c, "get_psychosocial_contact_attempts")` (misma puerta que usa el card para leer el historial; habilitada para `ps`/`ts`), inyecta `psicosocialId`, datos de la profesional (`currentUser`, `currentUserId`, `userTeam`) y `caseName` en el template.
- **Template**: `src/frontend/html/salvia/remision-temporal/remision_temporal.html`.

## Navegación de entrada (cambio brownfield)

En `src/frontend/html/salvia/mis-remisiones-psicosocial/mis_remisiones_psicosocial.html`, el handler `onVerRemision` deja de ir a `hacer-seguimiento` y pasa a:

```js
onVerRemision: function(payload) {
    var psicosocialId = payload && payload.remisionId; // remisionId = psychosocial_support_id (PsychosocialListItem.ID)
    if (psicosocialId) {
        window.location.href = '/salvia/remision-temporal/' + encodeURIComponent(psicosocialId);
    }
}
```

> [!NOTE]
> El componente `remisiones-psicosocial-component` ya emite `ver-remision` con `{ remisionId: r.id, followUpId, remision }`, donde `r.id` es el `psychosocial_support_id`. **Solo cambia el handler de esta pantalla**; el componente de listado no se toca. El mismo handler existe en `historial-remisiones` pero **no** se modifica (fuera de alcance: solo "Mis Remisiones Psicosocial").

## Diseño y UX

- **Layout**: cabecera mínima (título "Remisión — {caseName}") + el card 3x3 centrado en el contenedor estándar de `salvia-main`. Sin acciones adicionales; toda la interacción vive en el card y sus modales.
- **Accesibilidad**: heredada del card (ver su spec). La pantalla no añade inputs propios.

## Árbol de Interfaz

```text
Pantalla: Remisión Temporal
│
├── main.salvia-main > .container
│   ├── Cabecera "Remisión — ${caseName}"
│   └── <psychosocial-contact-card
│           :psicosocial-id="psicosocialId"
│           :case-id="caseId"
│           :case-name="caseName"
│           :current-user="currentUser"
│           :current-user-id="currentUserId"
│           :user-team="userTeam" />
│         └── (internamente compone <psychosocial-contact-modal /> — ver spec del card)
│
└── Overlays estándar (success/fail) provistos por overlay.html
```

## Inventario de Eventos

1. **`mounted` (pantalla)** — fija `document.title`, muestra `#app`. No hace fetch propio: el card carga su historial vía `GET /api/v1/psychosocial/{psicosocialId}/contact-attempts`.
2. Toda otra interacción (registrar intento, consentimiento, agendar, próximo intento) la maneja el card/modal — ver sus specs.

## Manejo de Estado

- Sin store global: la pantalla pasa props al card por atributos inyectados desde el facade (Vue embebido, patrón del proyecto). El card mantiene su propio estado local y refresca vía API.

## Naturaleza temporal (deuda reconocida)

> [!WARNING]
> Esta pantalla es un **andamio temporal** para poder ejercitar el flujo 3x3 antes de que exista la pantalla interna definitiva de Remisión Psicosocial. Cuando esa pantalla se especifique, el card se montará allí y esta ruta (`remision-temporal/:id`) y su facade se retiran, restaurando o reapuntando `onVerRemision`.
