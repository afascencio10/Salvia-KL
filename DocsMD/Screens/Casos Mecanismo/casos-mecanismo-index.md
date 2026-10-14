# Casos Mecanismo — Index

Pantalla de listado exclusiva del rol **Mecanismo articulador** (`mc`). Cada usuario `mc` ve solo los `case_articulating_mechanism` **activos** asignados a su `assigned_user_id` de sesión.

Al login, el rol `mc` entra en esta ruta. El menú lateral solo muestra **Casos Mecanismo**.

> **Ruta:** `GET /salvia/casos-mecanismo`  
> **Fuente de datos:** `salvia.case_articulating_mechanism` filtrado por `assigned_user_id` de sesión y `status = 'active'` → joins a `victim_case` / riesgo / ciudad, y última `case_task` Done de ese mecanismo.  
> **Permiso:** `get_casos_mecanismo` (solo `mc`).  
> **Detalle del schema:** [related-tables.md](./related-tables.md)

---

## Navegación

| Documento | Link |
|---|---|
| Interfaz | [casos-mecanismo-interface.md](./casos-mecanismo-interface.md) |
| Tablas relacionadas | [related-tables.md](./related-tables.md) |

---

## Roles con acceso

| Rol | Código | Alcance |
|---|---|---|
| Mecanismo articulador | `mc` | Solo los escalamientos `active` con `assigned_user_id` = usuario en sesión. |

---

## Modelo mental

```
usuario mc (sesión)
    │  general_user_i_code
    ▼
case_articulating_mechanism
  WHERE assigned_user_id = sesión
    AND status = 'active'
    │  case_id (= victim_case_i_code)
    ▼
victim_case (+ form2 riesgo, town→city)
    │
    └── última case_task Done
          WHERE case_articulating_mechanism_id = cam.id
          → result = “Último comentario del mecanismo”
```

Cada **Gestionar** crea una `case_task` nueva (`type = gestion_mecanismo`, `status = Done`) y un evento de timeline (categoría `Mecanismos`). El registro del mecanismo **permanece `active`**: puede haber varios reportes en el tiempo.

---

## Resumen de eventos

| # | Evento | Tipo | Flujo |
|---|---|---|---|
| E01 | Cuando carga la pantalla | Lifecycle | [📄 flow-E01-cuando-carga-pantalla.md](./Flujos/flow-E01-cuando-carga-pantalla.md) |
| E02 | Cuando filtra por documento | User Interaction | [📄 flow-E02-cuando-filtra-por-documento.md](./Flujos/flow-E02-cuando-filtra-por-documento.md) |
| E03 | Cuando filtra por ciudad | User Interaction | [📄 flow-E03-cuando-filtra-por-ciudad.md](./Flujos/flow-E03-cuando-filtra-por-ciudad.md) |
| E04 | Cuando cambia de página | User Interaction | [📄 flow-E04-cuando-cambia-de-pagina.md](./Flujos/flow-E04-cuando-cambia-de-pagina.md) |
| E05 | Cuando presiona "Ver caso" | User Interaction | [📄 flow-E05-cuando-presiona-ver-caso.md](./Flujos/flow-E05-cuando-presiona-ver-caso.md) |
| E06 | Cuando presiona "Gestionar" | User Interaction | [📄 flow-E06-cuando-presiona-gestionar.md](./Flujos/flow-E06-cuando-presiona-gestionar.md) |

---

## Inventario de eventos

---

**Nombre del evento:** Cuando carga la pantalla  
**Tipo:** Lifecycle  
**Descripción:** Valida rol `mc` y carga la primera página de mecanismos activos del usuario (`pageSize=5`, `ORDER BY created_at DESC`).  
**Requerido:** Sí  

📄 [Ver flujo → flow-E01-cuando-carga-pantalla.md](./Flujos/flow-E01-cuando-carga-pantalla.md)

---

**Nombre del evento:** Cuando filtra por documento  
**Tipo:** User Interaction  
**Descripción:** Tras debounce (~350 ms), resetea página a 0 y recarga con documento parcial (ILIKE), siempre acotado al usuario de sesión.  
**Requerido:** Sí  

📄 [Ver flujo → flow-E02-cuando-filtra-por-documento.md](./Flujos/flow-E02-cuando-filtra-por-documento.md)

---

**Nombre del evento:** Cuando filtra por ciudad  
**Tipo:** User Interaction  
**Descripción:** Tras debounce (~350 ms), resetea página a 0 y recarga con ciudad parcial (ILIKE sobre el nombre de ciudad de la víctima).  
**Requerido:** Sí  

📄 [Ver flujo → flow-E03-cuando-filtra-por-ciudad.md](./Flujos/flow-E03-cuando-filtra-por-ciudad.md)

---

**Nombre del evento:** Cuando cambia de página  
**Tipo:** User Interaction  
**Descripción:** Actualiza `currentPage`, scroll al inicio y consulta esa página manteniendo filtros.  
**Requerido:** Sí  

📄 [Ver flujo → flow-E04-cuando-cambia-de-pagina.md](./Flujos/flow-E04-cuando-cambia-de-pagina.md)

---

**Nombre del evento:** Cuando presiona "Ver caso"  
**Tipo:** User Interaction  
**Descripción:** Navega a `/salvia/casos/:i_code/detalle`. Requiere `get_case_detail_sv` (habilitado para `mc`).  
**Requerido:** Sí  

📄 [Ver flujo → flow-E05-cuando-presiona-ver-caso.md](./Flujos/flow-E05-cuando-presiona-ver-caso.md)

---

**Nombre del evento:** Cuando presiona "Gestionar"  
**Tipo:** User Interaction  
**Descripción:** Abre `case-task-modal` vía `openMecanismo({ mechanismId, caseId, scope })`. Al guardar, `POST /api/v1/case-tasks/gestion-mecanismo` crea una `case_task` Done; toast y recarga del listado.  
**Requerido:** Sí  

📄 [Ver flujo → flow-E06-cuando-presiona-gestionar.md](./Flujos/flow-E06-cuando-presiona-gestionar.md)

---

## Checklist de validación

- [x] ¿Cada acción del usuario puede ser manejada?
- [x] ¿La carga inicial de datos está cubierta?
- [x] ¿Filtro por documento incluido?
- [x] ¿Filtro por ciudad incluido?
- [x] ¿Paginación cubierta?
- [x] ¿Navegación “Ver caso” incluida?
- [x] ¿Gestionar crea una tarea Done nueva (no actualiza una previa)?
- [x] ¿API listado filtra por `assigned_user_id` de sesión (no por userId del cliente)?
- [x] ¿Login y “Volver” caen en esta pantalla para `mc`?
- [ ] ¿Hay actualizaciones en tiempo real? — No aplica

---

## Fuera de esta entrega

- Imprimir el formulario.
- Escalar / remitir a otro nivel de comité (el campo “Remisión a responsables” se guarda, no dispara escalamiento).
- Dashboard.
- Sincronizar la card `case-mecanismos` del Detalle del Caso (sigue con comentarios viejos).

---

## Relación con otros módulos

```
usuario mc
  └─ GET /salvia/casos-mecanismo
          │
          ▼
     GET /api/v1/mecanismos/mine
          │
          ├── Ver caso → Detalle del Caso (get_case_detail_sv)
          └── Gestionar → case-task-modal.openMecanismo
                            └── POST /api/v1/case-tasks/gestion-mecanismo
                                  └── case_task Done + result
```
