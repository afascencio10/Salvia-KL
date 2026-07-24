# Casos Entidad — Index

Pantalla de listado para el rol **Entidad** (`et`): el usuario elige una organización del catálogo y ve los `entity_case` cuyas sedes (`entity_branch`) pertenecen a esa entidad, con filtros por documento y ciudad.

> **Ruta:** `GET /salvia/casos-entidad`  
> **Fuente de datos:** `salvia.entity_case` → `entity_branch` → `entity` + datos de `victim_case` / contacto. UI aún con mock hasta existir el endpoint de listado.  
> **Detalle del schema:** [related-tables.md](./related-tables.md) · origen del modelo: [case-entities/related-tables.md](../../Componentes/case-entities/related-tables.md)  
> **Temporal:** la relación usuario logueado ↔ entidad **aún no está definida**. E01 toma la **primera** entidad del catálogo para inicializar; el selector (E07) permite cambiar a cualquier otra.

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
| Entidad | `et` | Accede a la pantalla. Por ahora puede consultar cualquier entidad del catálogo vía selector (E07). Alcance real por usuario: **pendiente**. |

---

## Resumen de eventos

| # | Evento | Tipo | Flujo |
|---|---|---|---|
| E01 | Cuando carga la pantalla | Lifecycle | [📄 flow-E01-cuando-carga-pantalla.md](./Flujos/flow-E01-cuando-carga-pantalla.md) |
| E02 | Cuando filtra por documento | User Interaction | [📄 flow-E02-cuando-filtra-por-documento.md](./Flujos/flow-E02-cuando-filtra-por-documento.md) |
| E03 | Cuando filtra por ciudad | User Interaction | [📄 flow-E03-cuando-filtra-por-ciudad.md](./Flujos/flow-E03-cuando-filtra-por-ciudad.md) |
| E04 | Cuando cambia de página | User Interaction | [📄 flow-E04-cuando-cambia-de-pagina.md](./Flujos/flow-E04-cuando-cambia-de-pagina.md) |
| E05 | Cuando presiona "Ver caso" | User Interaction | [📄 flow-E05-cuando-presiona-ver-caso.md](./Flujos/flow-E05-cuando-presiona-ver-caso.md) |
| E06 | Cuando presiona "Ver detalle" | User Interaction | [📄 flow-E06-cuando-presiona-ver-detalle.md](./Flujos/flow-E06-cuando-presiona-ver-detalle.md) |
| E07 | Cuando cambia entidad | User Interaction | [📄 flow-E07-cuando-cambia-entidad.md](./Flujos/flow-E07-cuando-cambia-entidad.md) |

---

## Inventario de eventos

---

**Nombre del evento:** Cuando carga la pantalla  
**Tipo:** Lifecycle  
**Descripción:** Valida rol `et`, carga `GET /api/v1/entities`, preselecciona la **primera** entidad del catálogo (temporal), carga ciudades de sus sedes y la primera página de casos (`pageSize=5`, `ORDER BY updated_at DESC`).  
**Requerido:** Sí  

📄 [Ver flujo → flow-E01-cuando-carga-pantalla.md](./Flujos/flow-E01-cuando-carga-pantalla.md)

---

**Nombre del evento:** Cuando filtra por documento  
**Tipo:** User Interaction  
**Descripción:** Tras debounce (~350 ms), resetea página a 0 y recarga el listado con el documento, manteniendo entidad activa y ciudad. Si no hay entidad seleccionada, no llama al API.  
**Requerido:** Sí  

📄 [Ver flujo → flow-E02-cuando-filtra-por-documento.md](./Flujos/flow-E02-cuando-filtra-por-documento.md)

---

**Nombre del evento:** Cuando filtra por ciudad  
**Tipo:** User Interaction  
**Descripción:** Tras debounce (~350 ms), resetea página a 0 y recarga el listado con la ciudad, manteniendo entidad activa y documento. Si no hay entidad seleccionada, no llama al API.  
**Requerido:** Sí  

📄 [Ver flujo → flow-E03-cuando-filtra-por-ciudad.md](./Flujos/flow-E03-cuando-filtra-por-ciudad.md)

---

**Nombre del evento:** Cuando cambia de página  
**Tipo:** User Interaction  
**Descripción:** Actualiza `currentPage`, hace scroll al inicio y consulta esa página manteniendo entidad y filtros.  
**Requerido:** Sí  

📄 [Ver flujo → flow-E04-cuando-cambia-de-pagina.md](./Flujos/flow-E04-cuando-cambia-de-pagina.md)

---

**Nombre del evento:** Cuando presiona "Ver caso"  
**Tipo:** User Interaction  
**Descripción:** Navega al detalle del caso con `caseId` (`entity_case.case_id`).  
**Requerido:** Sí  

📄 [Ver flujo → flow-E05-cuando-presiona-ver-caso.md](./Flujos/flow-E05-cuando-presiona-ver-caso.md)

---

**Nombre del evento:** Cuando presiona "Ver detalle"  
**Tipo:** User Interaction  
**Descripción:** Navega a `/salvia/entidad/:entityBranchICode` (pantalla planificada en case-entities).  
**Requerido:** Sí  

📄 [Ver flujo → flow-E06-cuando-presiona-ver-detalle.md](./Flujos/flow-E06-cuando-presiona-ver-detalle.md)

---

**Nombre del evento:** Cuando cambia entidad  
**Tipo:** User Interaction  
**Descripción:** El usuario elige una organización del catálogo completo. Actualiza entidad activa + sector del header, resetea página a 0 y carga los `entity_case` de todas las sedes de esa entidad. Conserva filtros de documento/ciudad.  
**Requerido:** Sí  

📄 [Ver flujo → flow-E07-cuando-cambia-entidad.md](./Flujos/flow-E07-cuando-cambia-entidad.md)

---

## Checklist de validación

- [x] ¿Cada acción del usuario puede ser manejada?
- [x] ¿La carga inicial de datos está cubierta?
- [x] ¿Filtros por documento y ciudad están incluidos?
- [x] ¿Paginación está cubierta?
- [x] ¿Navegación “Ver caso” y “Ver detalle” están incluidas?
- [x] ¿Cambio de entidad está incluido? — E07 (catálogo completo, selección libre)
- [x] ¿Schema de `entity_case` definido? — Sí
- [x] ¿Endpoint de listado et definido? — `GET /api/v1/entity-cases`
- [x] ¿Endpoint catálogo de entities definido? — `GET /api/v1/entities`
- [x] ¿Ciudades por entidad? — `GET /api/v1/entities/:id/cities`
- [ ] ¿Relación usuario et ↔ entidad definida? — Pendiente (hoy: primera del catálogo)
- [ ] ¿Destino de “Ver detalle” implementado? — Pendiente (E06 muestra toast hasta existir `/salvia/entidad/:id`)
- [ ] ¿Hay actualizaciones en tiempo real? — No aplica

---

## Relación con otros módulos

```
Casos Entidad (listado)
        │  E01: carga TODAS las entities → selector
        │  E07: usuario elige entity_id
        │  lee entity_case JOIN entity_branch WHERE entity_id = elegida
        │  last_action = entity_case.last_action (texto)
        ▼
Detalle del Caso              ← E05
/salvia/entidad/:branchICode  ← E06 (planificado)

case-entities (tab del detalle de caso)
        │  mismo entity_case, vista por caseId
```
