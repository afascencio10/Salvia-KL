---
okf_version: "1.0"
type: Project
title: "Feature: Consulta de Caso para Agente Psicosocial (Ver Caso)"
description: "El agente psicosocial (roles ps y ts) puede consultar el detalle completo de un caso, en modo de solo lectura, desde el botón 'Ver caso' en las pantallas Remisiones y Mis Remisiones. Las acciones operativas del módulo de seguimientos (reasignar caso, nuevo seguimiento, editar/reasignar/iniciar seguimiento, gestionar tarea) permanecen ocultas para este rol; para el resto de roles (sv, ro, op) la pantalla no cambia."
owner: "@frontend-squad"
status: active
tags: [feature, psicosocial, case-detail, rbac, brownfield]
dependencies:
  - Spec: 6-security/rbac-matrix.md
last_updated: "2026-07-03"
---

# Feature: Consulta de Caso para Agente Psicosocial (Ver Caso)

## Descripción

La pantalla de detalle del caso (`get_case_detail_sv.html`) fue construida originalmente para el módulo de Seguimientos (roles `sv`, `ro`, `op`) e incluye acciones operativas que modifican el estado del caso. El agente psicosocial (roles `ps` = Psicólogo, `ts` = Trabajador social) necesita acceder a esta misma pantalla desde el botón "Ver caso" de "Remisiones" y "Mis Remisiones Psicosocial", pero exclusivamente en modo de consulta: ninguna acción que modifique el estado del caso debe estar disponible para este rol.

Este feature **no crea una pantalla nueva** — reutiliza la existente y ajusta la visibilidad condicional por rol, siguiendo el sistema de permisos ya existente (RBAC por rol, ver [Matriz RBAC](/.knowledge/6-security/rbac-matrix.md)).

**Es brownfield (§8 AGENTS.md):** ningún archivo de código tocado por este feature tenía spec OKF antes. Este spec documenta el comportamiento actual (lazy specing) de los archivos modificados, además de la regla nueva.

**Hallazgo clave de la investigación de código:** la mayoría de los botones operativos de esta pantalla (`Reasignar` de caso, `+ Nuevo seguimiento`, `Iniciar`/`Editar`/`Reasignar` de seguimiento) ya están ocultos hoy para cualquier rol que no sea `sv`/`ro`/`op`, porque sus `v-if` son listas de permiso positivas (ej. `userRole === 'sv'`), no negativas. Es decir, **ya funcionan correctamente para `ps`/`ts` sin cambio de código** — solo falta verificarlo con test. Los dos únicos puntos que si requieren cambio de código son: (1) el botón "Gestionar" tarea, cuyo gate depende de *ownership* (`assignedUserId === userId`) y no de rol, y (2) dos widgets de prueba de desarrollador que no tienen ningún gate de rol.

## Documentos Relacionados

### Seguridad
- [Matriz RBAC](/.knowledge/6-security/rbac-matrix.md) — Matriz de control de acceso; el permiso `get_case_detail_sv` debe extenderse a `ps` y `ts`

### Frontend — Screens
| Pantalla | Spec | Estado |
| :--- | :--- | :---: |
| CaseDetail (adaptación de solo lectura para ps/ts) | [index.md](./frontend/screens/CaseDetail/index.md) | 🔜 |

### Frontend — Componentes
| Componente | Spec | Estado |
| :--- | :--- | :---: |
| CaseTasks (botón "Gestionar") | [CaseTasks.md](./frontend/components/CaseTasks.md) | 🔜 |

### Testing
| Plan | Spec | Estado |
| :--- | :--- | :---: |
| Test Plan | [consulta-caso-psicosocial-test-plan.md](./testing/consulta-caso-psicosocial-test-plan.md) | 🔜 |

## Historias de Usuario

| ID | Historia | Estado |
| :--- | :--- | :---: |
| US-CCP-01 | Como agente psicosocial (`ps`/`ts`), quiero consultar el detalle completo de un caso desde "Ver caso" en Remisiones/Mis Remisiones, para conocer su contexto sin poder modificar su estado. | 🔜 |
| US-CCP-02 | Como supervisor/operador/operador de riesgo (`sv`/`op`/`ro`), quiero seguir teniendo exactamente las mismas acciones de hoy en esta pantalla, sin ninguna regresión por este cambio. | 🔜 |

## Checklist de Implementación (cambios de código no cubiertos por un spec nuevo)

Estos cambios modifican un archivo que **ya tiene spec propio** (`rbac-matrix.md`, activo), por lo que no requieren un spec de endpoint nuevo — solo deben mantenerlo sincronizado en el mismo PR:

1. En `src/salvia/config/Menu.go`, agregar `"ps": true, "ts": true` a la entrada `PermissionsByRole["get_case_detail_sv"]`. Esta es la única puerta de autorización real hoy (la valida `CaseDetailFacade.go:20` vía `CheckAndGetSession(c, "get_case_detail_sv")` al cargar la pantalla HTML).
2. En `.knowledge/6-security/rbac-matrix.md`, en la tabla "Casos de víctima": (a) agregar columnas `ps` y `ts` a la tabla completa (hoy ninguna fila las incluye — sería el primer permiso real que reciben estos dos roles), y (b) marcar ✅ en la fila `get_case_detail_sv` para ambas columnas nuevas.

## Deuda Técnica Reconocida (Fuera de Alcance de Este Spec)

> [!NOTE]
> Decisión explícita del usuario (no es un gap silencioso): el API JSON de mutaciones (`case_detail_controller.go`: `CrearSeguimiento`, `AddTimelineEvent`, `ReasignarCaso`, `ReassignFollowUp`, `EditFollowUp`; `case_task_controller.go`: `Complete`, `CompleteWithFormData`, `Reassign`) no tiene ningún `CheckPermission` hoy — ni antes ni después de este feature. Cualquier sesión autenticada podría invocarlos directamente sin pasar por la UI. Este spec **no** lo resuelve. Queda pendiente para un spec de seguridad futuro dedicado a endpoint hardening (bloqueo por rol + validación de ownership).
