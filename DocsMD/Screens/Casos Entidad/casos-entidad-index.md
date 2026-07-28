# Casos Entidad — Index

Pantalla de listado para el rol **Entidad** (`et`): cada usuario `et` representa una **sede** (`entity_branch`). Ve solo los `entity_case` de su `entity_branch_id`, con filtro por documento.

> **Ruta:** `GET /salvia/casos-entidad`  
> **Fuente de datos:** `salvia.entity_case` filtrado por `entity_branch_id` del usuario logueado → joins a `victim_case` / riesgo.  
> **Relación usuario↔sede:** `security.general_user.entity_branch_id` (campo nuevo, solo relevante para rol `et`).  
> **Detalle del schema:** [related-tables.md](./related-tables.md)

---

## Navegación

| Documento | Link |
|---|---|
| Interfaz | [casos-entidad-interface.md](./casos-entidad-interface.md) |
| Tablas relacionadas | [related-tables.md](./related-tables.md) |

---

## Roles con acceso

| Rol | Código | Alcance |
|---|---|---|
| Entidad | `et` | Solo los `entity_case` de su `general_user.entity_branch_id`. El header muestra la **organización padre** (`entity`) de esa sede. |

---

## Modelo mental

```
usuario et (sesión)
    │  general_user.entity_branch_id
    ▼
entity_branch  (la sede = este usuario)
    │  entity_id                    entity_branch_town_code
    ▼                               ▼
entity (nombre en header)         ciudad implícita de la sede
    │                               (por eso NO hay filtro de ciudad)
    ▼
entity_case WHERE entity_branch_id = sede del usuario
```

---

## Resumen de eventos

| # | Evento | Tipo | Flujo |
|---|---|---|---|
| E01 | Cuando carga la pantalla | Lifecycle | [📄 flow-E01-cuando-carga-pantalla.md](./Flujos/flow-E01-cuando-carga-pantalla.md) |
| E02 | Cuando filtra por documento | User Interaction | [📄 flow-E02-cuando-filtra-por-documento.md](./Flujos/flow-E02-cuando-filtra-por-documento.md) |
| E04 | Cuando cambia de página | User Interaction | [📄 flow-E04-cuando-cambia-de-pagina.md](./Flujos/flow-E04-cuando-cambia-de-pagina.md) |
| E05 | Cuando presiona "Ver caso" | User Interaction | [📄 flow-E05-cuando-presiona-ver-caso.md](./Flujos/flow-E05-cuando-presiona-ver-caso.md) |
| E06 | Cuando presiona "Ver detalle" | User Interaction | [📄 flow-E06-cuando-presiona-ver-detalle.md](./Flujos/flow-E06-cuando-presiona-ver-detalle.md) |

### Eventos eliminados

| # | Evento | Motivo |
|---|---|---|
| ~~E03~~ | Cuando filtra por ciudad | La sede del usuario ya tiene `town_code` / ciudad; no aplica filtrar por ciudad. |
| ~~E07~~ | Cuando cambia entidad | El usuario `et` es una sede fija; no cambia de entidad. |

---

## Inventario de eventos

---

**Nombre del evento:** Cuando carga la pantalla  
**Tipo:** Lifecycle  
**Descripción:** Valida rol `et`, lee `entity_branch_id` del usuario en sesión, resuelve la organización padre (`entity`) para el header (nombre + sector), y carga la primera página de `entity_case` de esa sede (`pageSize=5`, `ORDER BY updated_at DESC`).  
**Requerido:** Sí  

📄 [Ver flujo → flow-E01-cuando-carga-pantalla.md](./Flujos/flow-E01-cuando-carga-pantalla.md)

---

**Nombre del evento:** Cuando filtra por documento  
**Tipo:** User Interaction  
**Descripción:** Tras debounce (~350 ms), resetea página a 0 y recarga el listado con documento parcial (ILIKE), siempre acotado a la sede del usuario.  
**Requerido:** Sí  

📄 [Ver flujo → flow-E02-cuando-filtra-por-documento.md](./Flujos/flow-E02-cuando-filtra-por-documento.md)

---

**Nombre del evento:** Cuando cambia de página  
**Tipo:** User Interaction  
**Descripción:** Actualiza `currentPage`, scroll al inicio y consulta esa página manteniendo el filtro de documento.  
**Requerido:** Sí  

📄 [Ver flujo → flow-E04-cuando-cambia-de-pagina.md](./Flujos/flow-E04-cuando-cambia-de-pagina.md)

---

**Nombre del evento:** Cuando presiona "Ver caso"  
**Tipo:** User Interaction  
**Descripción:** Navega a `/salvia/casos/:caseId/detalle`.  
**Requerido:** Sí  

📄 [Ver flujo → flow-E05-cuando-presiona-ver-caso.md](./Flujos/flow-E05-cuando-presiona-ver-caso.md)

---

**Nombre del evento:** Cuando presiona "Ver detalle"  
**Tipo:** User Interaction  
**Descripción:** Destino planificado `/salvia/entidad/:entityBranchICode` (aún no existe → toast temporal en código).  
**Requerido:** Sí  

📄 [Ver flujo → flow-E06-cuando-presiona-ver-detalle.md](./Flujos/flow-E06-cuando-presiona-ver-detalle.md)

---

## Checklist de validación

- [x] ¿Cada acción del usuario puede ser manejada?
- [x] ¿La carga inicial de datos está cubierta?
- [x] ¿Filtro por documento incluido?
- [x] ¿Paginación cubierta?
- [x] ¿Navegación “Ver caso” / “Ver detalle” incluidas?
- [x] ¿Relación usuario et ↔ sede definida? — `general_user.entity_branch_id`
- [x] ¿Sin selector de entidad ni filtro de ciudad? — E03/E07 eliminados
- [x] ¿Columna `entity_branch_id` en `general_user` creada? — Migración aplicada (pendiente poblar usuarios et concretos)
- [x] ¿Sesión expone `EntityBranchId` al facade/API? — `CommonSession` + login
- [x] ¿API listado filtra por sede de sesión (no por `entityId` del cliente)? — `GET /api/v1/entity-cases`
- [x] ¿UI sin dropdown ni filtro ciudad? — Header con entidad padre + sector
- [ ] ¿Destino “Ver detalle” implementado? — Pendiente
- [ ] ¿Hay actualizaciones en tiempo real? — No aplica

---

## Trabajo pendiente restante

1. **Poblar** `entity_branch_id` en usuarios `et` de prueba/producción (ver `migration-general-user-entity-branch-id.sql`).
2. **Alta/edición de usuarios:** obligar sede al crear rol `et` (UI setup).
3. Pantalla destino E06 `/salvia/entidad/:id`.
4. Reloguearse tras asignar sede (la sesión se llena en login).

---

## Relación con otros módulos

```
usuario et
  └─ general_user.entity_branch_id ──► entity_branch ──► entity (header)
                                              │
                                              ▼
                                       entity_case (listado)
                                              │
                              ┌───────────────┴───────────────┐
                              ▼                               ▼
                     Detalle del Caso (E05)          /salvia/entidad/:id (E06)
```
