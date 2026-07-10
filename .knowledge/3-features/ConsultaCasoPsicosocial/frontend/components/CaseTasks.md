---
okf_version: "1.0"
type: UI_Component
title: "Componente: CaseTasks (lista de tareas de un caso/barrera)"
description: "Componente Vue reutilizable que lista las CaseTask de un caso (o de una barrera específica) y expone el botón 'Gestionar' para completarlas. Embebido en la pantalla CaseDetail, en las tabs Tareas y Barreras."
owner: "@frontend-squad"
status: active
tags: [frontend, ui, component, psicosocial, case-detail, rbac]
dependencies:
  - Spec: 6-security/rbac-matrix.md
code_refs:
  - "src/frontend/js/components/case-tasks.js"
last_updated: "2026-07-03"
---

# Componente: CaseTasks

## Alcance de Este Spec

Lazy specing (§8 AGENTS.md) del comportamiento actual de `case-tasks.js` — reutilizado en la pantalla [CaseDetail](../screens/CaseDetail/index.md), tanto en la tab "Tareas" (`:case-id`, sin `barrier-id`) como en la tab "Barreras" al expandir una barrera (`:case-id` + `:barrier-id`). Documenta la regla nueva: el botón "Gestionar" no debe estar disponible para `ps`/`ts`.

## Props

| Prop | Tipo | Requerido | Descripción |
| :--- | :--- | :---: | :--- |
| `case-id` | String | Sí | ID del caso cuyas tareas se listan |
| `barrier-id` | String | No | Si se pasa, filtra tareas de esa barrera específica |
| `user-id` | String | Sí | ICode del usuario actual (usado para ownership) |
| `user-role` | String | Sí | Rol activo de la sesión (`sv`, `ro`, `op`, y a partir de este spec también `ps`/`ts`) |

## Eventos Emitidos

- `open-task-modal` — al hacer click en "Gestionar", pide a la pantalla contenedora abrir `<case-task-modal>` con el ID de la tarea.
- `task-completed` — tras completar una tarea vía el modal.

## Árbol de Interfaz (resumen)

```text
Componente: CaseTasks
│
└── Lista de tareas
    └── Por cada tarea:
        ├── [condición: puedeCompletar(task)] Botón "Gestionar" (evento → abrirModal → emit open-task-modal)
        └── [condición: false — código muerto, nunca se renderiza] Botón "Reasignar" (línea 191, `v-if="false && userRole === 'sv'"`)
```

## Regla de Negocio Actual (`puedeCompletar`, línea 117-120)

```js
puedeCompletar: function(task) {
    if (this.userRole === 'sv') return true;
    return task.assignedUserId === this.userId;
}
```

**Gap identificado:** esta función gatea por **ownership** (`assignedUserId === userId`), no por rol. Si una `CaseTask` llegara a estar asignada al `icode` de un usuario `ps`/`ts` (ej. vía el modelo `dupla`, que empareja psicólogo + trabajador social — `src/internal/models/dupla.go`), el botón "Gestionar" se mostraría incorrectamente para ese rol, violando el requisito de "ninguna acción que modifique el estado del caso" para el agente psicosocial.

## Cambio de Código Requerido

```js
puedeCompletar: function(task) {
    if (this.userRole === 'ps' || this.userRole === 'ts') return false;
    if (this.userRole === 'sv') return true;
    return task.assignedUserId === this.userId;
}
```

La negación por rol debe evaluarse **antes** que la comprobación de ownership, para que ningún dato de asignación futuro pueda habilitar el botón para este rol.

## Reglas de Autorización (RBAC)

| Acción | Rol(es) con acceso hoy | Cambia con este spec? |
| :--- | :--- | :---: |
| Ver lista de tareas (lectura) | Cualquiera que cargue la pantalla contenedora | No — `ps`/`ts` deben poder leer esta lista (Definition of Done exige ver la tab Tareas) |
| Botón "Gestionar" (completar tarea) | `sv` (todas), dueño de la tarea (cualquier rol) | Sí — se excluye explícitamente a `ps`, `ts` sin importar ownership |
| Botón "Reasignar" tarea | Ninguno (código muerto, `v-if="false"`) | No — fuera de alcance, no se activa en este spec |

> [!NOTE]
> Igual que en la pantalla CaseDetail, el backend (`case_task_controller.go: Complete, CompleteWithFormData, Reassign`) no valida ningún rol ni ownership hoy. Este spec solo oculta el botón en Vue; el hardening del API queda como deuda técnica (ver [index.md](../../index.md) del feature).
