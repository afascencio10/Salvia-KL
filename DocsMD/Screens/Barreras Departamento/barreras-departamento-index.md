# Barreras Departamento — Index

Pantalla de listado para el rol **Enlace territorial** (`en`): ve las barreras de `barrier_v2` cuyo `department_id` coincide con el departamento asignado al usuario (`general_user.general_user_assigned_department` → sesión `AssignedDepartmentID`). Incluye todos los estados; el estado es un filtro opcional.

No confundir con **Mis Barreras** (`/salvia/mis-barreras`): esa pantalla lista barreras vía `case_task` asignadas al agente. Barreras Departamento no filtra por asignación de tareas.

> **Ruta:** `GET /salvia/barreras-departamento`  
> **Fuente de datos:** `GET /api/v1/barriers-v2/department`  
> **Permiso:** `get_barreras_departamento` (solo `en`)  
> **UI vigente:** `barreras_departamento_v2.html` + `barreras_departamento_v2.js` + `barreras_departamento.css`  
> **UI anterior (conservada, no montada):** `barreras_departamento.html` + `barreras_departamento.js`

---

## Navegación

| Documento | Link |
|---|---|
| Interfaz | [barreras-departamento-interface.md](./barreras-departamento-interface.md) |
| Tablas relacionadas | [related-tables.md](./related-tables.md) |

**Changelog:** [changelogSep2026.md](./changelogSep2026.md)

---

## Roles con acceso

| Rol | Código | Alcance |
|---|---|---|
| Enlace territorial | `en` | Barreras con `barrier_v2.department_id` = departamento asignado en sesión. Si el usuario no tiene departamento, el API responde 403. |

---

## Modelo mental

```
usuario en (sesión)
    │  AssignedDepartmentID  ← general_user.general_user_assigned_department
    ▼
barrier_v2 WHERE department_id = :AssignedDepartmentID
              AND deleted_at IS NULL
    │
    ├─ victim_case (nombre, documento)
    ├─ victim_case_form2 (riesgo + teléfono; fallback teléfono del contacto)
    ├─ GET /api/v1/locations/cities (select de filtro)
    ├─ entity_branch (nombre de sede)
    └─ security.city (ciudad de la barrera; JOIN city_id::text)
```

---

## Resumen de eventos

| # | Evento | Tipo | Flujo |
|---|---|---|---|
| E01 | Cuando carga la pantalla | Lifecycle | [📄 flow-E01-cuando-carga-pantalla.md](./Flujos/flow-E01-cuando-carga-pantalla.md) |
| E02 | Cuando filtra por documento | User Interaction | [📄 flow-E02-cuando-filtra-por-documento.md](./Flujos/flow-E02-cuando-filtra-por-documento.md) |
| E03 | Cuando filtra por entidad | User Interaction | [📄 flow-E03-cuando-filtra-por-entidad.md](./Flujos/flow-E03-cuando-filtra-por-entidad.md) |
| E04 | Cuando filtra por ciudad | User Interaction | [📄 flow-E04-cuando-filtra-por-ciudad.md](./Flujos/flow-E04-cuando-filtra-por-ciudad.md) |
| E05 | Cuando filtra por estado | User Interaction | [📄 flow-E05-cuando-filtra-por-estado.md](./Flujos/flow-E05-cuando-filtra-por-estado.md) |
| E06 | Cuando cambia de página | User Interaction | [📄 flow-E06-cuando-cambia-de-pagina.md](./Flujos/flow-E06-cuando-cambia-de-pagina.md) |
| E07 | Cuando presiona "Ver caso" | User Interaction | [📄 flow-E07-cuando-presiona-ver-caso.md](./Flujos/flow-E07-cuando-presiona-ver-caso.md) |
| E08 | Cuando presiona "Ver barrera" | User Interaction | [📄 flow-E08-cuando-presiona-ver-barrera.md](./Flujos/flow-E08-cuando-presiona-ver-barrera.md) |

---

## Inventario de eventos

**Nombre del evento:** Cuando carga la pantalla  
**Tipo:** Lifecycle  
**Descripción:** Valida rol `en`, muestra `#app` y carga la primera página (`pageSize=10`) sin filtros. El departamento lo aplica el servidor desde la sesión.  
**Requerido:** Sí  

📄 [Ver flujo → flow-E01-cuando-carga-pantalla.md](./Flujos/flow-E01-cuando-carga-pantalla.md)

---

**Nombre del evento:** Cuando filtra por documento  
**Tipo:** User Interaction  
**Descripción:** Tras debounce 350 ms, resetea página a 0 y consulta el API con `docNumber`.  
**Requerido:** Sí  

📄 [Ver flujo → flow-E02-cuando-filtra-por-documento.md](./Flujos/flow-E02-cuando-filtra-por-documento.md)

---

**Nombre del evento:** Cuando filtra por entidad  
**Tipo:** User Interaction  
**Descripción:** Tras debounce 350 ms, resetea página a 0 y consulta el API con `entidad` (ILIKE sobre `entity_branch_name`).  
**Requerido:** Sí  

📄 [Ver flujo → flow-E03-cuando-filtra-por-entidad.md](./Flujos/flow-E03-cuando-filtra-por-entidad.md)

---

**Nombre del evento:** Cuando filtra por ciudad  
**Tipo:** User Interaction  
**Descripción:** Select de ciudades del departamento (catálogo completo). Sin debounce; envía `cityId` exacto.  
**Requerido:** Sí  

📄 [Ver flujo → flow-E04-cuando-filtra-por-ciudad.md](./Flujos/flow-E04-cuando-filtra-por-ciudad.md)

---

**Nombre del evento:** Cuando filtra por estado  
**Tipo:** User Interaction  
**Descripción:** Select inmediato (sin debounce): envía `status` de BD o vacío (todos). Resetea página a 0.  
**Requerido:** Sí  

📄 [Ver flujo → flow-E05-cuando-filtra-por-estado.md](./Flujos/flow-E05-cuando-filtra-por-estado.md)

---

**Nombre del evento:** Cuando cambia de página  
**Tipo:** User Interaction  
**Descripción:** Actualiza `currentPage` y vuelve a consultar el API manteniendo filtros.  
**Requerido:** Sí  

📄 [Ver flujo → flow-E06-cuando-cambia-de-pagina.md](./Flujos/flow-E06-cuando-cambia-de-pagina.md)

---

**Nombre del evento:** Cuando presiona "Ver caso"  
**Tipo:** User Interaction  
**Descripción:** Navega a `/salvia/casos/:caseId/detalle`.  
**Requerido:** Sí  

📄 [Ver flujo → flow-E07-cuando-presiona-ver-caso.md](./Flujos/flow-E07-cuando-presiona-ver-caso.md)

---

**Nombre del evento:** Cuando presiona "Ver barrera"  
**Tipo:** User Interaction  
**Descripción:** Navega a `/salvia/barreras/:id`.  
**Requerido:** Sí  

📄 [Ver flujo → flow-E08-cuando-presiona-ver-barrera.md](./Flujos/flow-E08-cuando-presiona-ver-barrera.md)

---

## Checklist de validación

- [x] ¿Carga inicial cubierta?
- [x] ¿Filtros documento / entidad / ciudad / estado pegan al API?
- [x] ¿Paginación servidor `pageSize=10`?
- [x] ¿Departamento solo desde sesión?
- [x] ¿Todas las barreras (no solo `!= MANAGED`)?
- [x] ¿UI 2 columnas en archivos v2 sin borrar la UI anterior?
- [x] ¿Navegación Ver caso / Ver barrera?
