---
okf_version: "1.0"
type: UI_Screen
title: "Pantalla: Detalle del Caso (CaseDetail)"
description: "Pantalla de detalle de un caso de víctima. Sirve tanto al módulo de Seguimientos (sv/ro/op, lectura y escritura) como, a partir de este feature, al agente psicosocial (ps/ts, exclusivamente lectura). Server-rendered (Go template) con una única instancia de Vue montada sobre el shell HTML — no hay ruteo SPA (ver ADR-001)."
owner: "@frontend-squad"
status: active
tags: [frontend, ui, psicosocial, case-detail, rbac]
dependencies:
  - Spec: 6-security/rbac-matrix.md
code_refs:
  - "src/frontend/html/salvia/case_detail/get_case_detail_sv.html"
  - "src/salvia/facades/CaseDetailFacade.go"
last_updated: "2026-07-03"
---

# Pantalla: Detalle del Caso (CaseDetail)

## Alcance de Este Spec

Este documento hace *lazy specing* (§8 AGENTS.md) del comportamiento **actual** de la pantalla — construida para `sv`/`ro`/`op` — y añade la regla de negocio **nueva**: acceso de solo consulta para `ps`/`ts`. No rediseña la pantalla ni cambia su comportamiento para `sv`/`ro`/`op`.

**Ruta**: `GET /salvia/casos/:id/detalle` → `CaseDetailFacade.go:CaseDetailGET` (renderiza el shell HTML e inyecta `userRole` desde la sesión). Los datos del caso se cargan client-side desde `GET /api/v1/casos/:id/detalle` (`case_detail_controller.go:GetByID`), fuera del alcance de este spec (no se modifica).

**Entrada para `ps`/`ts`**: desde el componente `remisiones-psicosocial-component.js` (reusado por Remisiones, Mis Remisiones Psicosocial e Historial de Remisiones), evento `@ver-caso` → navega a `window.location.href = '/salvia/casos/' + code + '/detalle'`. Este componente **ya funciona hoy** y no se modifica en este feature.

## Diseño y UX

- **Layout**: header con nombre/datos de la víctima y badges de riesgo/estado; banner de próximo seguimiento; banner de tareas pendientes; panel con 6 tabs (Info General, Derivaciones, Barreras, Seguimientos, Tareas, Timeline).
- **Para `ps`/`ts`**: mismo layout, mismas 6 tabs, mismos datos — sin ningún elemento nuevo. Los botones que ejecutan acciones de escritura no se muestran (ver "Reglas de Autorización" abajo).
- **Accesibilidad**: sin cambios respecto al comportamiento actual.

## Árbol de Interfaz (resumen — solo elementos relevantes a este spec)

```text
Pantalla: CaseDetail
│
├── Header
│   ├── Botón "📋 Info caso" (evento → abrirInfoCaso) — solo lectura, visible para TODOS los roles incl. ps/ts
│   ├── [condición: userRole === 'sv'] Botón "Reasignar" (evento → abrirModalReasignar)
│   └── [condición: userRole === 'sv' OR userRole === 'ro'] Botón "+ Nuevo seguimiento" (evento → abrirModalNuevoSeg)
│
├── Tabs (cd-tabs) — las 6 se muestran para TODOS los roles incl. ps/ts, sin condición
│   ├── Tab "Info General" — solo lectura para todos los roles, sin acciones de escritura
│   ├── Tab "Derivaciones"
│   │   ├── Secciones de solo lectura (Medidas de Emergencia, Psicosocial, Estabilización Socioeconómica)
│   │   ├── [condición NUEVA: userRole NOT IN ('ps','ts')] Widget dev "🧪 Probar case-task-modal" (líneas 265-282)
│   │   └── [condición NUEVA: userRole NOT IN ('ps','ts')] Widget dev "🧪 Probar case-task-history" (líneas 284-301)
│   ├── Tab "Barreras" — lista de solo lectura; al expandir una barrera embebe <case-tasks> (ver spec de componente CaseTasks.md — botón "Gestionar" gateado ahí, no en esta pantalla)
│   ├── Tab "Seguimientos"
│   │   └── Por cada seguimiento pendiente:
│   │       ├── [condición: (userRole === 'ro' OR userRole === 'op') AND seg.agent_id === userICode] Link "▶ Iniciar" (evento → iniciarSeguimiento)
│   │       ├── [condición: userRole === 'sv' OR (userRole === 'ro' AND caso.agentId === userICode)] Botón "✏️ Editar" (evento → editarSeguimiento)
│   │       └── [condición: userRole === 'sv' OR ((userRole === 'ro' OR userRole === 'op') AND seg.agent_id === userICode)] Botón "📋 Reasignar" (evento → reasignarSeguimiento)
│   ├── Tab "Tareas" — embebe <case-tasks> (ver spec de componente CaseTasks.md)
│   └── Tab "Timeline" — embebe <case-timeline>, solo lectura
│
└── Overlays (modales, hermanos del contenido principal)
    ├── Modal "Nuevo seguimiento" [visible si modalNuevoSeg] — solo se abre si el botón que lo dispara está visible
    ├── Modal "Reasignar caso" [visible si modalReasignar] — solo se abre si el botón que lo dispara está visible
    ├── Modal "Reasignar seguimiento" [visible si modalReasignarSeg] — solo se abre si el botón que lo dispara está visible
    ├── Modal "Editar seguimiento" [visible si modalEditarSeg] — solo se abre si el botón que lo dispara está visible
    └── Modal "Info caso" [visible si mostrarInfoCaso] — solo lectura, visible para todos los roles
```

**Nota importante (verificada leyendo el código, no asumida):** los `v-if` de "Reasignar" (caso), "+ Nuevo seguimiento", "▶ Iniciar", "✏️ Editar" y "📋 Reasignar" (seguimiento) son listas de permiso **positivas** por rol (`userRole === 'sv'`, `userRole === 'ro'`, `userRole === 'op'`). Como `ps` y `ts` nunca aparecen en esas comparaciones, **estos 5 elementos ya están correctamente ocultos hoy para ps/ts sin ningún cambio de código**. Este spec **no** modifica esas condiciones — solo exige que el plan de pruebas (`testing/consulta-caso-psicosocial-test-plan.md`) lo verifique explícitamente, ya que es un supuesto crítico para el Definition of Done.

## Cambios de Código Requeridos en Este Archivo

Los únicos dos elementos de esta pantalla que **no** tienen gate de rol hoy son los widgets de desarrollador. Deben ganar la condición nueva:

```html
<!-- get_case_detail_sv.html, línea ~266 y ~285 -->
<div v-if="!esConsultaSoloLectura" style="background:#fefce8;...">🧪 Dev — Probar case-task-modal ...</div>
<div v-if="!esConsultaSoloLectura" style="background:#f0fdf4;...">🧪 Dev — Probar case-task-history ...</div>
```

Agregar un computed nuevo en el `data()`/`computed` de la instancia Vue (junto a `puedeReasignar`, línea ~840):

```js
computed: {
  esConsultaSoloLectura() {
    return this.userRole === 'ps' || this.userRole === 'ts';
  },
  // ...resto de computed sin cambios
}
```

Este mismo computed se reutiliza en la spec del componente `CaseTasks.md` (pasado o replicado allí — ver ese documento para el detalle, ya que `case-tasks.js` es un componente separado con su propio `userRole` prop).

## Reglas de Autorización (RBAC)

| Elemento | Rol(es) con acceso hoy | Cambia con este spec? |
| :--- | :--- | :---: |
| Cargar la pantalla (`get_case_detail_sv`) | `op, ro, do, et, us, sv, no, fo` | Sí — se agregan `ps`, `ts` en `Menu.go` (ver [index.md](../../../index.md) del feature) |
| Botón "Reasignar" (caso) | `sv` | No |
| Botón "+ Nuevo seguimiento" | `sv, ro` | No |
| Botón "▶ Iniciar" (seguimiento) | `ro, op` (dueño) | No |
| Botón "✏️ Editar" (seguimiento) | `sv, ro` (dueño) | No |
| Botón "📋 Reasignar" (seguimiento) | `sv, ro, op` (dueño) | No |
| Widgets dev (`case-task-modal`/`case-task-history`) | Todos (sin gate) | Sí — se ocultan para `ps`, `ts` |
| Botón "📋 Info caso" | Todos | No (de solo lectura, se mantiene visible para todos) |

> [!IMPORTANT]
> El backend (`case_detail_controller.go`) no valida ningún rol en sus endpoints mutables. Ocultar estos botones en Vue es la única barrera hoy — ver "Deuda Técnica Reconocida" en el [index.md](../../../index.md) del feature.

## Manejo de Estado

- No usa store global (Redux/Vuex/Pinia) — estado local en `data()` de la instancia Vue montada por página (patrón server-rendered, ver ADR-001).
- `userRole` se inyecta server-side desde la sesión (`CaseDetailFacade.go:35`, `vars["userRole"] = sc.Session.CurrentRole`) y llega al `data()` de Vue como string interpolado (`get_case_detail_sv.html:773`).
- Nuevo: `esConsultaSoloLectura` (computed) deriva de `userRole`, sin nuevo estado ni llamada de red.
