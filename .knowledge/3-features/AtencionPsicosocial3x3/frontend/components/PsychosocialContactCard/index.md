---
okf_version: "1.0"
type: UI_Component
title: "Componente: psychosocial-contact-card"
description: "Card reutilizable 'Intentos de contacto 3x3' que muestra contador de intentos, tiles de último/próximo intento, historial de contacto colapsable con badges (NO CONTESTÓ / SÍ CONTESTÓ / PROGRAMADO) y un botón 'Registrar intento' que dispara el flujo de modales (psychosocial-contact-modal). Autosuficiente: carga su propio historial vía GET y se monta en una o varias pantallas."
owner: "@frontend-squad"
status: active
tags: [frontend, component, vue, psicosocial, 3x3, card]
dependencies:
  - Endpoint: 3-features/AtencionPsicosocial3x3/backend/endpoints/list-contact-attempts.md
  - Component: 3-features/AtencionPsicosocial3x3/frontend/components/PsychosocialContactModal/index.md
code_refs:
  - "src/frontend/js/components/psychosocial-contact-card.js"
last_updated: "2026-07-09"
---

# Componente: `<psychosocial-contact-card />`

Card Vue **reutilizable** que resume y gestiona el flujo 3x3 de contacto psicosocial de un proceso `psychosocial_support`. Es **la unidad que se monta en las pantallas** (una o varias); **compone** internamente el flujo de modales `psychosocial-contact-modal` (invocado por `ref` desde el botón "Registrar intento").

Es autosuficiente: al montarse carga su propio historial vía `GET /api/v1/psychosocial/{psicosocialId}/contact-attempts`.

## Props (API Pública)

```typescript
interface PsychosocialContactCardProps {
  psicosocialId: string;   // proceso psychosocial_support
  caseId: string;
  caseName: string;
  currentUser: string;
  currentUserId: string;   // profesional ps/ts
  userTeam: string;
  closureThreshold?: number; // denominador del contador "X de N"; default 9 (umbral de elegibilidad de cierre)
}
```

## Fuentes de datos

- `GET .../contact-attempts` → `{ attempts[], counters, lastAttemptAt, nextContactAttemptAt, processStatus }`.
- Refresco tras cada evento `completed` emitido por el modal (ver Eventos).

## Anatomía visual (según mockup aprobado)

```text
┌─ Card "Intentos de contacto 3x3" ────────────────────────────────────┐
│  Título: "Intentos de contacto 3x3"        [ Registrar intento ]      │
│  Subtítulo: "${counters.totalCount} de ${closureThreshold} intentos"  │
│  Ayuda: "Al llegar a ${closureThreshold} intentos en ≥3 días distintos │
│          se habilita el cierre por imposibilidad de contacto."         │
│                                                                        │
│  ┌─ Tile ────────────────┐   ┌─ Tile ─────────────────┐               │
│  │ 🕐 Último intento      │   │ 📅 Próximo intento      │               │
│  │    realizado          │   │    ${nextContactAttemptAt│               │
│  │    ${lastAttemptAt}    │   │      | '—' si null}      │               │
│  └───────────────────────┘   └────────────────────────┘               │
│                                                                        │
│  ┌─ Sección colapsable "HISTORIAL DE CONTACTO"  [chevron ▲/▼] ───────┐ │
│  │  (timeline vertical, MÁS NUEVO → MÁS VIEJO)                       │ │
│  │  ○ Próximo      🕐 ${nextContactAttemptAt}         [PROGRAMADO]   │ │
│  │  ● Intento #2   🕐 ${attemptAt}                    [SÍ CONTESTÓ]  │ │
│  │      ${note || 'Sí contestó'}                                     │ │
│  │  ● Intento #1   🕐 ${attemptAt}                    [NO CONTESTÓ]  │ │
│  │      ${note || 'Sin nota'}                                        │ │
│  └──────────────────────────────────────────────────────────────────┘ │
│                                                                        │
│  <psychosocial-contact-modal ref="modal" ... />  (oculto hasta abrir) │
└────────────────────────────────────────────────────────────────────────┘
```

### Badges de estado por fila del historial

| Badge | Condición | Estilo (referencia) |
| :--- | :--- | :--- |
| `NO CONTESTÓ` | `attempt.wasAnswered === false` | rojo (fondo red-100, texto red-700) |
| `SÍ CONTESTÓ` | `attempt.wasAnswered === true` | verde (fondo green-100, texto green-700) |
| `PROGRAMADO` | fila sintética derivada de `nextContactAttemptAt` (no es un `contact_attempts` real) | violeta (fondo violet-100, texto violet-700) |

- Título de cada fila real: `Intento #${sequenceNumber}` (numeración cronológica, 1 = más antiguo).
- Orden de presentación: **más nuevo → más viejo** (el API entrega los intentos ya en ese orden; ver [list-contact-attempts](../../../backend/endpoints/list-contact-attempts.md)).
- Texto de la fila: `note` (nota libre) si existe; si no, un fallback (`"Sin nota"` en fallidos, `"Sí contestó"` en exitosos).
- La fila `PROGRAMADO` se renderiza **solo** si `nextContactAttemptAt != null`, siempre **al inicio** del timeline (es un evento futuro).

## Estados del Componente

| Estado | Descripción |
| :--- | :--- |
| `loading` | Cargando historial **en la carga inicial** (GET en vuelo). Muestra skeleton. En los **refrescos** (tras un `completed` del modal) NO se activa `loading`: el re-fetch es silencioso para no desmontar el bloque del card ni el `<psychosocial-contact-modal>` que está en pleno flujo (evita perder, p. ej., el modal de Consentimiento). |
| `loaded` | Card renderizado con datos. |
| `historyCollapsed` | Sección "HISTORIAL DE CONTACTO" plegada/desplegada (toggle chevron). |
| `error` | Fallo de carga; mensaje con reintento. |

## Árbol de Renderizado

```text
<psychosocial-contact-card />
│
├── [loading]  → skeleton
├── [error]    → mensaje + botón "Reintentar"
└── [loaded]
      ├── Header (título + subtítulo "X de N intentos" + botón "Registrar intento")
      │      └── click "Registrar intento" → $refs.modal.start(processContext)
      ├── Tiles (Último intento realizado | Próximo intento)
      ├── Historial colapsable (más nuevo → más viejo)
      │      ├── fila sintética "Próximo" con badge PROGRAMADO (si nextContactAttemptAt) — arriba
      │      └── v-for attempts (desc) → fila con badge NO CONTESTÓ / SÍ CONTESTÓ
      └── <psychosocial-contact-modal @completed="onModalCompleted" />
```

## Interacción con el modal

- El botón **"Registrar intento"** llama a `$refs.modal.start(ctx)` con `ctx = { id: psicosocialId, caseId, caseName, dailyFailedCount, totalCount, distinctDaysCount }`.
- El modal decide (según `dailyFailedCount >= 3`) si abre "¿Contestó?" o directamente el modal de acciones (ver [PsychosocialContactModal](../PsychosocialContactModal/index.md)).
- Tras cualquier `completed`, el card refresca el historial (re-fetch **silencioso**, sin activar `loading`) para reflejar el nuevo intento / próximo intento / transición de estado, **sin desmontar el modal** que pueda seguir abierto (Consentimiento, Sesión, Acciones, Cierre).

## Accesibilidad (a11y)

- La cabecera de "HISTORIAL DE CONTACTO" es un botón con `aria-expanded` que refleja `historyCollapsed`.
- Cada badge incluye texto accesible (no solo color).
- Tiles con `aria-label` ("Último intento realizado: …", "Próximo intento: …").

## Eventos Emitidos (re-emitidos desde el modal)

| Evento | Trigger | Payload |
| :--- | :--- | :--- |
| `refreshed` | Historial recargado tras un cambio | `{ psicosocialId, counters }` |
| `contacted` | Contacto exitoso + paso a en_gestion | `{ psicosocialId }` |
| `closure-triggered` | Se disparó cierre (no acepta / imposibilidad) | `{ psicosocialId, reason }` |

## Montaje en pantallas

Este card se monta en **una o varias pantallas** del módulo psicosocial (p.ej. detalle del proceso / gestión del caso). Las pantallas anfitrionas concretas se listan a medida que se especifican; cada una solo debe pasar las props (`psicosocialId`, `caseId`, `caseName`, datos de la profesional). El card no asume una pantalla específica.

> [!NOTE]
> **Confirmado:** el denominador "de 9" del contador es el **umbral de elegibilidad de cierre** (9 intentos). La interfaz lo comunica en el subtítulo ("X de 9 intentos") y en el texto de ayuda. Se expone como prop `closureThreshold` (default 9). El cierre por imposibilidad requiere además ≥3 días distintos (RN-06); por eso el texto de ayuda menciona ambas condiciones.
