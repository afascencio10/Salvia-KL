# Mis Reportes — Inventario de Eventos

Pantalla: `GET /salvia/mis-reportes` (rol `op`; `sv` y los demás roles → 401)
Template: `get_my_reportes.html` (nuevo)
API de listado: `GET /api/v1/victim-contacts/list` (mismo endpoint que `Listado Reportes`, con `filter_key=persona_asignada` fijo)

---

## Eventos identificados

### E-01 — Cuando carga la pantalla

📄 [Ver flujo → flow-E01-cuando-carga-pantalla.md](./Flujos/flow-E01-cuando-carga-pantalla.md)

```
Evento:       Cuando carga la pantalla
Tipo:         Lifecycle
Descripción:  MyReportsGET valida el permiso (solo op; otro rol → 401) y
              renderiza get_my_reportes.html. Monta casos-component con
              :agent-id="currentUserId" (fuerza persona_asignada = usuario de
              la sesión, no editable) y reasignacion=false.
Requerido:    Sí
```

---

### E-02 — Cuando presiona "Ver Reporte" en una fila

📄 [Ver flujo → flow-E02-cuando-presiona-ver-reporte.md](./Flujos/flow-E02-cuando-presiona-ver-reporte.md)

```
Evento:       Cuando presiona "Ver Reporte" en una fila
Tipo:         User Interaction
Descripción:  Idéntico a E-04 de Listado Reportes: casos-component emite
              action-clicked { buttonId:'ver', case }, la pantalla abre
              gestion-reporte-modal con los mismos botones de acción
              (op y sv pueden gestionar cualquier reporte, GAP-2).
Requerido:    Sí
```

---

## Checklist de completitud

- [x] ¿Carga inicial? → E-01
- [x] ¿Acción de ver detalle? → E-02
- [x] ¿Hay eventos propios de filtros/paginación/orden? → No, heredados de `casos-component` sin cambios
- [x] ¿Hay reasignación? → No aplica (`reasignacion=false`)

---

## Resumen

| # | Evento | Tipo | Persiste en backend |
|---|---|---|---|
| E-01 | Cuando carga la pantalla | Lifecycle | No (solo lee) |
| E-02 | Cuando presiona "Ver Reporte" | User Interaction | No (apertura de modal) |

**Total: 2 eventos propios de la pantalla.**
