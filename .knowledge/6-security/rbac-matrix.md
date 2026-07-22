---
okf_version: "1.0"
type: Security_Policy
title: "Matriz RBAC y Permisos"
description: "Roles del sistema SALVIA y matriz de permisos por operación, derivada de la tabla security.role (BD) y del mapa PermissionsByRole (código)."
owner: "@security-guild"
status: active
tags: [security, rbac, compliance]
last_updated: "2026-07-17"
code_refs: [src/salvia/config/Menu.go]
---

# Matriz de Control de Acceso (RBAC)

Esta matriz dicta la autorización. Todos los endpoints (y componentes UI) deben alinearse con estas reglas.

## Mecanismo de Enforcement

- La autorización se resuelve en el **backend**, dentro de cada facade, con:
  `utils.CheckPermission(salvia_config.PermissionsByRole, "<permiso>", s.CurrentRole, c)`.
- `PermissionsByRole` (mapa `permiso → { role_code: true }`) vive en [`src/salvia/config/Menu.go`](/.knowledge/1-architecture/container-diagram.md) y es la **fuente de verdad del código**.
- `s.CurrentRole` es el **rol activo de la sesión** (autenticación por sesión server-side, ver [ADR-001](/.knowledge/1-architecture/adrs/adr-001-stack-go-gin-vue-postgres.md)). En la práctica cada usuario tiene un único rol (16 455 asignaciones para 16 454 usuarios en `security.rel_role_general_user`).
- Los mismos códigos filtran el menú que ve el frontend, pero **el menú es solo presentación**: la decisión vinculante es el `CheckPermission` del backend.

> [!WARNING]
> La lógica de Backend **debe** verificar el permiso vía `CheckPermission` en el facade. Jamás confiar en que el Frontend esconde el botón o la opción de menú.

## Catálogo de Roles (`security.role`)

| Código | Rol | Usuarios* |
| :--- | :--- | ---: |
| `ad` | Administrador | 7 |
| `op` | Operador | 153 |
| `sv` | Supervisor | 36 |
| `ro` | Operador de riesgo | 48 |
| `do` | Operador territorial departamental | 4 |
| `no` | Operador territorial nacional | 1 |
| `et` | Entidad | 2 |
| `en` | Enlace territorial | 1 |
| `an` | Agente de notificaciones | 1 |
| `us` | Usuario externo | 16 199 |
| `ps` | Psicólogo | 1 |
| `ts` | Trabajador social | 0 |

\* Conteo de `security.rel_role_general_user` en `develop` (Supabase) al 2026-07-02.

> [!IMPORTANT]
> **Discrepancias código ↔ BD detectadas (deuda a resolver):**
> - El código usa un rol **`fo`** en `PermissionsByRole` que **no existe** en `security.role`. Un usuario nunca tendrá ese rol, así que esas líneas son permisos muertos (o falta crear el rol).
> - Los roles **`ps` (Psicólogo)** y **`ts` (Trabajador social)** existen en BD y desde [Consulta de Caso Psicosocial](/.knowledge/3-features/ConsultaCasoPsicosocial/index.md) tienen **una única** operación habilitada en `PermissionsByRole`: `get_case_detail_sv` (solo lectura). Para el resto de operaciones de este documento siguen sin ninguna entrada.
> - `en` (Enlace territorial) solo habilita `get_mis_barreras`.

## Matriz de Permisos por Operación

Marca ✅ = el rol puede ejecutar la operación; `·` = denegado. Derivada 1:1 de `PermissionsByRole`.

Códigos de columna: `ad`=Administrador, `sv`=Supervisor, `op`=Operador, `ro`=Operador de riesgo, `do`=Op. territorial deptal, `no`=Op. territorial nacional, `et`=Entidad, `en`=Enlace territorial, `an`=Agente notificaciones, `us`=Usuario externo, `fo`=(no está en BD)

### Casos de víctima

> [!NOTE]
> Esta tabla agrega las columnas `ps` y `ts` porque, a la fecha, son las únicas dos columnas de todo este documento donde esos roles tienen alguna entrada en `PermissionsByRole` (ver [Consulta de Caso Psicosocial](/.knowledge/3-features/ConsultaCasoPsicosocial/index.md)). El resto de tablas de este documento mantiene solo las 11 columnas originales porque `ps`/`ts` no aparecen en ningún otro permiso.

| Permiso / operación | ad | sv | op | ro | do | no | et | en | an | us | fo | ps | ts |
| :--- | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: |
| `get_victim_case` | · | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | · | · | ✅ | ✅ | · | · |
| `get_victim_cases` | · | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | · | · | ✅ | ✅ | · | · |
| `get_victim_case_by_document` | · | ✅ | ✅ | ✅ | · | ✅ | ✅ | · | · | · | · | · | · |
| `get_victim_cases_by_all` | · | ✅ | ✅ | ✅ | · | ✅ | ✅ | · | · | · | · | · | · |
| `set_victim_case` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | · | ✅ | · | ✅ | · | · |
| `update_victim_case` | · | · | ✅ | ✅ | · | · | · | · | · | · | · | · | · |
| `report_victim_cases` | · | ✅ | ✅ | ✅ | · | ✅ | · | · | · | · | · | · | · |
| `approve_victim_case` | · | · | ✅ | ✅ | · | · | · | · | · | · | · | · | · |
| `get_list_cases` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | · | ✅ | ✅ | ✅ | · | · |
| `get_my_cases` | ✅ | ✅ | ✅ | ✅ | · | ✅ | ✅ | · | ✅ | · | · | · | · |
| `get_case_detail_sv` | · | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | · | · | ✅ | ✅ | ✅ | ✅ |
| `set_case_log` | · | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | · | · | ✅ | · | · | · |

### Atención Psicosocial — Flujo 3x3

> [!NOTE]
> Permisos nuevos del [Flujo 3x3 de Atención Psicosocial](/.knowledge/3-features/AtencionPsicosocial3x3/index.md). Habilitados **solo** para `ps` y `ts` (por eso esta tabla también incluye esas dos columnas). Deben reflejarse 1:1 en `PermissionsByRole` (`src/salvia/config/Menu.go`) en el mismo PR de implementación.

| Permiso / operación | ad | sv | op | ro | do | no | et | en | an | us | fo | ps | ts |
| :--- | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: |
| `get_psychosocial_contact_attempts` | · | · | · | · | · | · | · | · | · | · | · | ✅ | ✅ |
| `register_psychosocial_contact_attempt` | · | · | · | · | · | · | · | · | · | · | · | ✅ | ✅ |
| `set_psychosocial_consent` | · | · | · | · | · | · | · | · | · | · | · | ✅ | ✅ |
| `schedule_psychosocial_session` | · | · | · | · | · | · | · | · | · | · | · | ✅ | ✅ |
| `set_psychosocial_next_attempt` | · | · | · | · | · | · | · | · | · | · | · | ✅ | ✅ |
| `init_psychosocial_closure_form` | · | · | · | · | · | · | · | · | · | · | · | ✅ | ✅ |
| `close_psychosocial_process` | · | · | · | · | · | · | · | · | · | · | · | ✅ | ✅ |

### Contactos de víctima

| Permiso / operación | ad | sv | op | ro | do | no | et | en | an | us | fo |
| :--- | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: |
| `get_victim_contact` | · | ✅ | ✅ | ✅ | · | ✅ | ✅ | · | · | · | · |
| `get_victim_contacts` | · | ✅ | ✅ | ✅ | · | ✅ | ✅ | · | · | · | · |
| `get_victim_contacts_by_all` | · | ✅ | ✅ | ✅ | · | ✅ | ✅ | · | · | · | · |
| `set_victim_contact` | · | · | ✅ | ✅ | · | · | ✅ | · | · | · | · |
| `put_victim_contact` | · | · | ✅ | ✅ | · | · | · | · | · | · | · |

### Entidades y sedes

| Permiso / operación | ad | sv | op | ro | do | no | et | en | an | us | fo |
| :--- | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: |
| `get_entity` | · | · | · | · | ✅ | · | · | · | · | · | · |
| `get_entity_branches` | · | · | · | ✅ | ✅ | · | · | · | · | · | · |
| `get_entity_branches_by_towncode` | ✅ | ✅ | ✅ | ✅ | · | ✅ | ✅ | · | · | · | · |
| `get_entity_branches_by_towncode_and_entity` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | · | · | · | · |
| `get_entity_branches_with_moments` | ✅ | ✅ | ✅ | ✅ | · | ✅ | ✅ | · | · | · | · |
| `set_entity_branch` | · | · | · | · | ✅ | · | · | · | · | · | · |
| `put_entity_branch` | · | · | · | · | ✅ | · | · | · | · | · | · |
| `update_entity_branch` | · | · | · | · | ✅ | · | · | · | · | · | · |
| `update_moment` | · | · | · | · | · | · | ✅ | · | · | · | · |

### Alertas

| Permiso / operación | ad | sv | op | ro | do | no | et | en | an | us | fo |
| :--- | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: |
| `get_alerts_by_victim_case` | · | ✅ | ✅ | ✅ | · | ✅ | ✅ | · | · | · | · |
| `get_alerts_by_town_code` | · | ✅ | ✅ | ✅ | · | ✅ | ✅ | · | · | · | · |
| `get_alerts_in_danger` | · | ✅ | ✅ | ✅ | · | ✅ | · | · | · | · | · |
| `get_alerts_by_all` | · | ✅ | ✅ | ✅ | · | · | ✅ | · | · | · | · |

### Seguimientos

| Permiso / operación | ad | sv | op | ro | do | no | et | en | an | us | fo |
| :--- | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: |
| `set_follow_up` | · | · | · | ✅ | · | · | · | · | · | · | · |
| `set_follow_up_entry` | · | · | · | ✅ | · | · | · | · | · | · | · |
| `update_follow_up_entry` | · | · | · | ✅ | · | · | · | · | · | · | · |
| `set_follow_up_entry_acting` | · | · | · | · | ✅ | · | · | · | · | · | · |
| `get_hacer_seguimiento` | ✅ | · | ✅ | ✅ | · | · | · | · | · | · | · |
| `generate_calendario_seguimiento` | ✅ | ✅ | ✅ | · | · | · | · | · | · | · | · |
| `get_seguimiento_detalle` | ✅ | ✅ | ✅ | · | · | · | · | · | · | · | · |
| `get_seguimiento_detalle_caso` | ✅ | ✅ | ✅ | · | · | · | · | · | · | · | · |
| `get_seguimiento_formulario` | · | · | ✅ | · | · | · | · | · | · | · | · |
| `get_seguimientos_area` | ✅ | ✅ | · | · | · | · | · | · | · | · | · |
| `get_mis_seguimientos_dia` | · | · | ✅ | · | · | · | · | · | · | · | · |

### Barreras

| Permiso / operación | ad | sv | op | ro | do | no | et | en | an | us | fo |
| :--- | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: |
| `get_barrier` | · | · | · | ✅ | · | · | · | · | · | · | · |
| `get_mis_barreras` | · | · | · | · | · | · | · | ✅ | · | · | · |

### Feminicidio y riesgo

| Permiso / operación | ad | sv | op | ro | do | no | et | en | an | us | fo |
| :--- | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: |
| `get_feminicide` | · | · | · | ✅ | · | · | · | · | · | · | ✅ |
| `set_feminicide` | · | · | · | ✅ | · | · | · | · | · | · | ✅ |
| `get_feminicides` | · | · | · | ✅ | · | · | · | · | · | · | ✅ |
| `get_feminicide_risk` | · | · | · | ✅ | · | · | · | · | · | · | ✅ |
| `set_feminicide_risk` | · | · | · | ✅ | · | · | · | · | · | · | ✅ |
| `get_feminicide_risks` | · | · | · | ✅ | · | · | · | · | · | · | ✅ |

### Remisiones, KPIs y notificaciones

| Permiso / operación | ad | sv | op | ro | do | no | et | en | an | us | fo |
| :--- | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: |
| `get_historial_remisiones` | · | ✅ | · | · | · | · | · | · | · | · | · |
| `get_kpis` | ✅ | ✅ | · | · | · | · | · | · | · | · | · |
| `get_notificaciones` | ✅ | ✅ | ✅ | ✅ | · | ✅ | ✅ | · | ✅ | · | · |
| `assign_operators` | · | ✅ | · | · | · | · | · | · | · | · | · |
| `load_plain_files` | ✅ | · | · | · | · | · | · | · | · | · | · |
| `report_followups_consolidated` | ✅ | ✅ | · | · | · | · | · | · | · | · | · |
| `report_contacts_consolidated` | ✅ | ✅ | · | · | · | · | · | · | · | · | · |

> [!NOTE]
> `report_followups_consolidated` pertenece al feature [Reporte consolidado de seguimientos](/.knowledge/3-features/ReporteConsolidadoSeguimientos/index.md) y ya está implementado en `PermissionsByRole` (`Menu.go`) y en el endpoint `POST /api/v1/reportes/seguimientos-consolidado`. Habilita descargar el Excel consolidado (solo `ad`, `sv`).

> [!NOTE]
> `report_contacts_consolidated` pertenece al feature [Reporte consolidado de contactos](/.knowledge/3-features/ReporteConsolidadoContactos/index.md) (spec `active`, **implementado** en `Menu.go`). Habilita descargar el Excel consolidado de reportes vía `POST /api/v1/reportes/contactos-consolidado` (solo `ad`, `sv`).

## Cómo mantener esta matriz

Esta matriz es un **espejo** de `PermissionsByRole` en `src/salvia/config/Menu.go` (declarado en `code_refs`). Al agregar o cambiar un permiso en ese mapa, actualiza aquí la fila correspondiente en el mismo PR — el hook Spec-First exige que el cambio de código y su spec viajen juntos.
