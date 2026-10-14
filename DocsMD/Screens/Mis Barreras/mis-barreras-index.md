# Mis Barreras — Index

Pantalla de listado para el rol **Enlace territorial** (`en`): ve las barreras de `barrier_v2` cuyo `user_territory_link_id` coincide con el usuario autenticado, en estado activo (`OPEN` o `En Gestion`). Incluye el caso y, si existe, una tarea ToDo de tipo `articulacion_territorial` asignada al mismo usuario.

No confundir con **Barreras Departamento** (`/salvia/barreras-departamento`): esa pantalla lista barreras del departamento asignado, todos los estados, sin recorte por `user_territory_link_id`.

> **Ruta:** `GET /salvia/mis-barreras`  
> **Fuente de datos:** `GET /api/v1/barriers-v2/mine`  
> **Guardar gestión:** `POST /api/v1/case-tasks/articulacion-territorial`  
> **Permiso:** `get_mis_barreras` (solo `en`)  
> **UI:** `barrieris.html` + `barriers.js` + `barriers.css` + `case-task-modal.js`

---

## Navegación

| Documento | Link |
|---|---|
| Interfaz | [mis-barreras-interface.md](./mis-barreras-interface.md) |
| Tablas relacionadas | [related-tables.md](./related-tables.md) |

**Changelog:** [changelogSep2026.md](./changelogSep2026.md)

---

## Roles con acceso

| Rol | Código | Alcance |
|---|---|---|
| Enlace territorial | `en` | Barreras con `user_territory_link_id` = `UserICode` de sesión y status `OPEN` o `En Gestion`. Si no hay departamento asignado el API responde 403 (`requireEnlaceAccess`). |

---

## Modelo mental

```
usuario en (sesión)
    │  UserICode
    ▼
barrier_v2 WHERE user_territory_link_id = :UserICode
              AND status IN ('OPEN', 'En Gestion')
              AND deleted_at IS NULL
    │
    ├─ victim_case (nombre, documento)
    ├─ victim_case_form2 / victim_contact_form1 (teléfono)
    ├─ security.city (nombre; city_id::text = barrier_v2.city_id)
    ├─ entity_branch (nombre de sede)
    └─ case_task ToDo type=articulacion_territorial asignada al usuario (opcional)
```

---

## Resumen de eventos

| # | Evento | Tipo | Flujo |
|---|---|---|---|
| E01 | Cuando carga la pantalla | Lifecycle | [📄 flow-E01-cuando-carga-pantalla.md](./Flujos/flow-E01-cuando-carga-pantalla.md) |
| E02 | Cuando filtra por documento | User Interaction | [📄 flow-E02-cuando-filtra-por-documento.md](./Flujos/flow-E02-cuando-filtra-por-documento.md) |
| E03 | Cuando filtra por entidad | User Interaction | [📄 flow-E03-cuando-filtra-por-entidad.md](./Flujos/flow-E03-cuando-filtra-por-entidad.md) |
| E04 | Cuando cambia de página | User Interaction | [📄 flow-E04-cuando-cambia-de-pagina.md](./Flujos/flow-E04-cuando-cambia-de-pagina.md) |
| E05 | Cuando presiona Gestionar | User Interaction | [📄 flow-E05-cuando-presiona-gestionar.md](./Flujos/flow-E05-cuando-presiona-gestionar.md) |
| E06 | Cuando filtra por ciudad | User Interaction | [📄 flow-E06-cuando-filtra-por-ciudad.md](./Flujos/flow-E06-cuando-filtra-por-ciudad.md) |

El Guardar del modal se documenta en `case-task-modal` (E09).

---

## Inventario de eventos

**Nombre del evento:** Cuando carga la pantalla  
**Tipo:** Lifecycle  
**Descripción:** Muestra `#app`, carga ciudades del departamento asignado y barreras activas del enlace desde `GET /api/v1/barriers-v2/mine`.  
**Requerido:** Sí  

📄 [Ver flujo → flow-E01-cuando-carga-pantalla.md](./Flujos/flow-E01-cuando-carga-pantalla.md)

---

**Nombre del evento:** Cuando filtra por documento  
**Tipo:** User Interaction  
**Descripción:** Filtra en cliente por `docNumber` (includes). Resetea página a 0.  
**Requerido:** Sí  

📄 [Ver flujo → flow-E02-cuando-filtra-por-documento.md](./Flujos/flow-E02-cuando-filtra-por-documento.md)

---

**Nombre del evento:** Cuando filtra por entidad  
**Tipo:** User Interaction  
**Descripción:** Filtra en cliente por sector o nombre de sede (ILIKE-like includes). Resetea página a 0.  
**Requerido:** Sí  

📄 [Ver flujo → flow-E03-cuando-filtra-por-entidad.md](./Flujos/flow-E03-cuando-filtra-por-entidad.md)

---

**Nombre del evento:** Cuando cambia de página  
**Tipo:** User Interaction  
**Descripción:** Pagina en cliente (`itemsPerPage=5`) el arreglo ya filtrado.  
**Requerido:** Sí  

📄 [Ver flujo → flow-E04-cuando-cambia-de-pagina.md](./Flujos/flow-E04-cuando-cambia-de-pagina.md)

---

**Nombre del evento:** Cuando presiona Gestionar  
**Tipo:** User Interaction  
**Descripción:** Abre `case-task-modal.openArticulacion` con o sin `taskId`. Al `@completed` recarga el listado.  
**Requerido:** Sí  

📄 [Ver flujo → flow-E05-cuando-presiona-gestionar.md](./Flujos/flow-E05-cuando-presiona-gestionar.md)

---

**Nombre del evento:** Cuando filtra por ciudad  
**Tipo:** User Interaction  
**Descripción:** Select de ciudades del departamento asignado (opción Todas). Filtra en cliente por `cityId` exacto. Resetea página a 0.  
**Requerido:** Sí  

📄 [Ver flujo → flow-E06-cuando-filtra-por-ciudad.md](./Flujos/flow-E06-cuando-filtra-por-ciudad.md)

---

## Checklist de validación

- [x] ¿Carga inicial cubierta?
- [x] ¿Solo barreras activas del enlace autenticado?
- [x] ¿Sin tabs Por articular / Articuladas?
- [x] ¿Cuatro columnas (sin Gestión Realizada)?
- [x] ¿Gestionar abre el mismo modal con o sin tarea?
- [x] ¿Refresh tras guardar?
- [x] ¿Columna Caso con Doc, Tel y ciudad (sin i_code)?
- [x] ¿Filtro ciudad del departamento asignado?
