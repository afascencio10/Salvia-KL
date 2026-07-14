---
okf_version: "1.0"
type: UI_Screen
title: "Pantalla: Detalle de Remisión Psicosocial (host definitivo del card 3x3)"
description: "Pantalla interna definitiva en /salvia/remision-psicosocial/:id (id = psychosocial_support_id) que reemplaza al andamio remision-temporal/:id como host del card psychosocial-contact-card. El card se monta en la pestaña 'Contactos', sustituyendo el placeholder estático del 3x3; convive con la lista existente de sesiones (team_contact). El botón 'Registrar intento' solo es visible para ps/ts (canRegister); los demás roles ven el card en modo lectura."
owner: "@frontend-squad"
status: active
tags: [frontend, ui, psicosocial, 3x3]
dependencies:
  - Component: 3-features/AtencionPsicosocial3x3/frontend/components/PsychosocialContactCard/index.md
  - Endpoint: 3-features/AtencionPsicosocial3x3/backend/endpoints/list-contact-attempts.md
  - Screen: 3-features/AtencionPsicosocial3x3/frontend/screens/RemisionTemporal/index.md
code_refs:
  - "src/salvia/facades/RemisionPsicosocialDetalleFacade.go"
  - "src/salvia/facades/MainRouter.go"
  - "src/salvia/config/Enums.go"
  - "src/frontend/html/salvia/remision-psicosocial/remision_psicosocial_detalle.html"
  - "src/frontend/html/salvia/mis-remisiones-psicosocial/mis_remisiones_psicosocial.html"
  - "src/frontend/js/components/psychosocial-contact-card.js"
last_updated: "2026-07-14"
---

# Pantalla: Detalle de Remisión Psicosocial (host definitivo del card 3x3)

Pantalla **interna definitiva** de una remisión psicosocial. Sustituye al andamio temporal [`remision-temporal/:id`](../RemisionTemporal/index.md) como host del card [`psychosocial-contact-card`](../../components/PsychosocialContactCard/index.md): al existir esta pantalla, la pantalla temporal y su facade se **retiran** (ver su spec).

Se llega aquí desde **Mis Remisiones Psicosocial** (roles `ps`/`ts`) y desde **Historial de Remisiones** (rol `sv`), mediante "Ver remisión".

> [!NOTE]
> **Alcance de este spec (§8 Brownfield).** La pantalla `remision_psicosocial_detalle.html`, su facade `RemisionPsicosocialDetalleFacade.go`, el `psychosocial_detail_controller.go`/`psychosocial_detail_service.go` (endpoints `/api/v1/psychosocial-support/...`) y las pestañas **Info / Tareas / Timeline** y la lista de **sesiones (team_contact)** de la pestaña Contactos **preexisten sin spec OKF** (llegaron en el commit `bb63bec`). Este spec documenta **solo** el montaje del card 3x3 en la pestaña Contactos y los archivos que ese cambio toca. El resto de la pantalla queda como **deuda documental medible** (`make coverage`).

## Ruta y montaje

- **Ruta**: `GET /salvia/remision-psicosocial/:id` — `:id` = **`psychosocial_support_id`**.
- **Registro**: `secRouter.GET("/remision-psicosocial/:id", RemisionPsicosocialDetalleGET)` en `src/salvia/facades/MainRouter.go`.
- **Facade**: `RemisionPsicosocialDetalleGET` en `src/salvia/facades/RemisionPsicosocialDetalleFacade.go`. Valida sesión (redirige a landing si no hay), e inyecta en el template: `remisionId` (= `psychosocial_support_id`), `currentUser`, `currentRole`, `currentUserId`, `userTeam`, `locale`, `lang`. **No** restringe el acceso por rol: cualquier usuario autenticado que llegue por navegación entra (coherente con "card visible para todos"); el control de escritura vive en la visibilidad del botón (`canRegister`) y en el RBAC del backend.
- **Template**: `src/frontend/html/salvia/remision-psicosocial/remision_psicosocial_detalle.html`.
- **Scripts requeridos por el card** (se añaden al final del template, junto a los ya presentes): `dinamic-form.js`, `psychosocial-contact-modal.js`, `psychosocial-contact-card.js`. El card trae su propio `<style>` (clases `psc-*`), por lo que no requiere CSS adicional.

## Navegación de entrada (cambio brownfield)

En `src/frontend/html/salvia/mis-remisiones-psicosocial/mis_remisiones_psicosocial.html`, el handler `onVerRemision` deja de apuntar al andamio y pasa a la pantalla definitiva:

```js
onVerRemision: function(payload) {
    var psicosocialId = payload && payload.remisionId; // remisionId = psychosocial_support_id
    if (psicosocialId) {
        window.location.href = '/salvia/remision-psicosocial/' + encodeURIComponent(psicosocialId);
    }
}
```

> [!NOTE]
> `Historial de Remisiones` (`historial_remisiones.html`) **ya** navega a `/salvia/remision-psicosocial/:id` desde el commit `bb63bec` — no se toca. El único cambio de navegación es el de Mis Remisiones (que iba al andamio temporal).

## Visibilidad del botón "Registrar intento" (RBAC en UI)

- El card recibe una prop **`canRegister`** (ver [spec del card](../../components/PsychosocialContactCard/index.md)). El template la calcula desde `currentRole`:

  ```js
  canRegister = (currentRole === 'ps' || currentRole === 'ts')
  ```

- Solo `ps`/`ts` ven el botón "Registrar intento" y pueden abrir el flujo de modales, en línea con el permiso **`register_psychosocial_contact_attempt`** ([RBAC](/.knowledge/6-security/rbac-matrix.md), solo `ps`/`ts`). El supervisor (`sv`) y cualquier otro rol ven el card en **modo lectura** (contador, tiles e historial), sin el botón.
- Esta es una defensa de UI. La escritura real la protege el RBAC del backend en el endpoint `POST /contact-attempts`. Ver "Deuda / Riesgo conocido".

## Diseño y UX

- **Layout**: sin cambios estructurales. Dentro de la pestaña **Contactos** del detalle, el bloque **placeholder estático del 3x3** (los tres círculos `0/3`, el botón muerto y el texto "Sin intentos registrados.") se **reemplaza** por `<psychosocial-contact-card>`. La sección inferior **"Contactos (N)"** (lista de sesiones `team_contact` con Registrar/Reprogramar/Cancelar) se **mantiene sin cambios**.
- **Accesibilidad**: heredada del card (cabecera de historial con `aria-expanded`, badges con texto accesible, tiles con `aria-label`). La pantalla no añade inputs propios al montar el card.

## Árbol de Interfaz (solo la porción afectada: pestaña Contactos)

```text
Pantalla: Detalle de Remisión Psicosocial
│
└── main > .container.brd-container
    ├── Header (víctima, badges riesgo/estado, profesional)   [preexistente]
    ├── Tabs [ Info | Contactos | Tareas | Timeline ]         [preexistente]
    └── [tab activo: 'contactos']
        ├── <psychosocial-contact-card                        ← REEMPLAZA el placeholder estático
        │      :psicosocial-id="remisionID"
        │      :current-user="currentUser"
        │      :current-user-id="currentUserId"
        │      :user-team="userTeam"
        │      :can-register="canRegister" />
        │     └── (internamente compone <psychosocial-contact-modal /> — ver spec del card)
        └── Sección "Contactos (N)" (lista de sesiones team_contact)  [preexistente, sin cambios]
```

## Inventario de Eventos

1. **`mounted` (pantalla)** — comportamiento preexistente: muestra `#app` y carga el detalle vía `GET /api/v1/psychosocial-support/:id/detail`. Este spec no lo modifica.
2. **Card montado** — el card hace su **propio** fetch de historial vía `GET /api/v1/psychosocial/{psicosocialId}/contact-attempts` al montarse; la pantalla no orquesta ese fetch. Toda la interacción 3x3 (registrar intento, consentimiento, agendar, próximo intento, cierre) la maneja el card/modal — ver sus specs.

## Manejo de Estado

- Sin store global. La pantalla pasa props al card por atributos inyectados desde el facade (patrón Vue embebido del proyecto). El card mantiene su estado local y refresca vía API tras cada `completed` del modal.

## Deuda / Riesgo conocido

> [!WARNING]
> **Lectura del historial sin RBAC de backend.** El endpoint `GET /api/v1/psychosocial/:id/contact-attempts` **no** está protegido por middleware RBAC (el grupo `/api/v1` no lo aplica). Al montar el card para "todos los roles" que abran el detalle, cualquiera con la pantalla abierta puede leer el historial. Se acepta como comportamiento actual (el andamio temporal ya exponía el mismo GET). Endurecer este endpoint (y reflejar la lectura de `sv` en la matriz RBAC) queda como deuda del feature 3x3, fuera del alcance de esta pantalla.
