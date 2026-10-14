# Lista de casos — Index

Pantalla de listado de `victim_case` que monta [`casos-component`](../../Componentes/Casos%20component/casos-component-events.md). Ruta `GET /salvia/lista-casos`.

Para el rol **Enlace territorial** (`en`) el listado se acota automáticamente al departamento asignado (`security.general_user.general_user_assigned_department` en sesión). Los demás roles (incluido `et`) ven el universo completo. **Casos Entidad** (`/salvia/casos-entidad`) es otra pantalla, solo para `et`: lista `entity_case` por sede.

> **Fuente de datos:** `GET /api/v1/cases/list`  
> **Permiso HTML:** `get_list_cases`  
> **Departamento:** nunca viaja como query param; si el rol es `en`, el backend lo toma de `CommonSession.AssignedDepartmentID`.

---

## Navegación

| Documento | Link |
|---|---|
| Interfaz | [lista-casos-interface.md](./lista-casos-interface.md) |
| Flujo de carga | [Flujos/flow-E01-cuando-carga-pantalla.md](./Flujos/flow-E01-cuando-carga-pantalla.md) |
| Componente | [casos-component-events.md](../../Componentes/Casos%20component/casos-component-events.md) |

---

## Roles con acceso

| Rol | Código | Alcance en esta pantalla |
|---|---|---|
| Enlace territorial | `en` | Solo casos cuyo municipio de víctima (`victim_case_victim_town_code` → ciudad → `department_id`) coincide con el departamento asignado. Sin departamento: mensaje de error, no se monta la tabla. |
| Resto con `get_list_cases` | `sv`, `op`, `ro`, `ad`, `et`, etc. | Todos los casos; filtros del componente sin scope de departamento. |

---

## Modelo mental (rol en)

```
usuario en (sesión)
    │  general_user_assigned_department
    ▼
security.department.department_id
    ▲
    │  city.department_id
    │
victim_case.victim_case_victim_town_code → town → city
```

---

## Resumen de eventos

| # | Evento | Tipo | Flujo |
|---|---|---|---|
| E01 | Cuando carga la pantalla | Lifecycle | [flow-E01-cuando-carga-pantalla.md](./Flujos/flow-E01-cuando-carga-pantalla.md) |

Los filtros, paginación, búsqueda y acciones de fila viven en el componente (E-01, E-03–E-16, E-05).

---

## Relación con otros módulos

- **Mis casos** usa el mismo API: si un `en` lo usara, quedaría acotado a departamento ∩ `agentId`.
- **Casos Entidad** no se modifica (sigue siendo solo `et`).
- El departamento del Enlace ya se configura en el alta de usuario (`set_general_user.html`). Tras asignarlo hay que reloguear.
