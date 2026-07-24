# `Casos Entidad` — Interfaz de la Pantalla

Listado de casos con víctimas referidas a la entidad del usuario. Cada fila muestra datos del caso (víctima, documento, ciudad, riesgo, estado) y la relación con la entidad (descripción de la actuación + última acción del timeline).

> **Estado:** UI conectada a API para E01 (catálogo + primera entidad + listado). Filtros documento/ciudad y paginación usan el mismo endpoint de listado.

---

## Archivos relevantes

| Archivo | Rol |
|---|---|
| `src/frontend/html/salvia/entity_case/casos_entidad.html` | Template principal — markup Vue |
| `src/frontend/js/components/casos-entidad.js` | Lógica Vue (E01–E07) |
| `src/frontend/css/casos-entidad.css` | Estilos (prefijo `ce-`) |
| `src/salvia/facades/CasosEntidadFacade.go` | Facade GET — permiso `get_casos_entidad` |
| `src/salvia/controller/entity_api_controller.go` | `GET /api/v1/entities`, `GET /api/v1/entities/:id/cities` |
| `src/salvia/controller/entity_case_controller.go` | `GET /api/v1/entity-cases` (+ rutas por caso) |
| `src/internal/models/entity_case.go` | Modelo + DTOs de listado/catálogo |
| `src/internal/repository/entity_case_repository.go` | Queries listado / entities / cities |

**Ruta:** `GET /salvia/casos-entidad`  
**Rol:** `et` (Entidad) — único rol con acceso.

---

## Árbol de interfaz

```
Pantalla: Casos Entidad
│
├── [v-if loadError] Bloque de error global
│   └── Mensaje + acción opcional "Reintentar"  → reloadScreen()
│
└── [v-else]
    ├── welcome-section  (.ce-welcome)   // mismo patrón que Historial de Remisiones
    │   ├── Título  <h2><b>Casos Entidad</b></h2>
    │   ├── Subtítulo  (.text-muted)  "Casos con víctimas referidas a tu entidad"
    │   ├── SectorLabel  (.ce-sector)   entity.sectorName   (ej. "Justicia")
    │   └── EntitySelector  (.ce-entity-selector)
    │       ├── Label  "Entidad:"
    │       └── <select>  → filters.entityId   [E07]
    │           └── <option> × N  [v-for entities]
    │
    ├── FilterCard  (.ce-filters)
    │   ├── Campo "Documento de identidad"
    │   │   └── <input type="text">  (.ce-input)
    │   │       placeholder "Ej: 1023456780"
    │   │       → filters.document   [E02]
    │   └── Campo "Ciudad"
    │       └── <input type="text">  (.ce-input)
    │           placeholder "Ej: San Salvador"
    │           → filters.city   [E03]
    │
    ├── [v-if isLoading] Skeleton / spinner de carga
    │
    ├── [v-else-if items.length === 0] EmptyState  (.ce-empty)
    │   └── Texto: "No hay casos referidos a tu entidad con los filtros actuales."
    │
    └── [v-else] ListCard  (.ce-list)
        ├── ListHeader  (.ce-list-header)
        │   ├── Col  "CASO"
        │   └── Col  "RELACIÓN CON ENTIDAD"
        │
        ├── CaseRow × N  [v-for items]  (.ce-row)
        │   │
        │   ├── ColCaso  (.ce-col-caso)
        │   │   ├── VictimName  (.ce-name)   item.victimFullName
        │   │   ├── Identifiers  (.ce-ids)   "{caseCode} · {document}"
        │   │   ├── Location  (.ce-city)
        │   │   │   ├── PinIcon
        │   │   │   └── CityName  item.city
        │   │   ├── Badges  (.ce-badges)
        │   │   │   ├── RiskBadge  (.ce-badge--risk)   item.riskLevel
        │   │   │   │   State: Bajo | Medio | Alto | Extremo
        │   │   │   └── StatusBadge  (.ce-badge--status)   item.caseStatus
        │   │   │       State: Activo | En seguimiento | … (según catálogo)
        │   │   └── Link  "Ver caso →"  → goToCase(item)   [E05]
        │   │
        │   └── ColRelacion  (.ce-col-relacion)
        │       ├── RelationDesc  (.ce-relation)   item.objetivo   // entity_case.objetivo
        │       ├── LastActionBlock  (.ce-last-action)
        │       │   ├── Label  "Última acción"
        │       │   ├── [v-if item.lastAction]
        │       │   │   └── LastActionText  item.lastAction   // entity_case.last_action (texto)
        │       │   └── [v-else]
        │       │       └── EmptyAction  "Sin acciones registradas"  (.ce-last-action--empty)
        │       └── Link  "Ver detalle →"  → goToRelationDetail(item)   [E06]
        │
        └── Pagination  (.ce-pagination)   [v-if totalPages > 1]
            ├── BtnPrev  "Anterior"  → changePage(currentPage - 1)   [E04]
            │   :disabled si currentPage === 0
            ├── PageBtn × N  → changePage(n)   [E04]
            └── BtnNext  "Siguiente"  → changePage(currentPage + 1)   [E04]
                :disabled si currentPage >= totalPages - 1
```

---

## Estados visuales de badges

| Campo | Valores esperados (ejemplo UI) | Clase sugerida |
|---|---|---|
| Riesgo (`riskLevel`) | Bajo / Medio / Alto / Extremo | `.ce-badge--risk` + modificador por nivel |
| Estado (`caseStatus`) | Activo / En seguimiento / … | `.ce-badge--status` + modificador por estado |

---

## Notas de UI

1. E01 preselecciona la **primera** entidad del catálogo; el selector permite cambiar (E07).
2. Ciudad: `<select>` con ciudades de las sedes de la entidad; el filtro aplica sobre la **ciudad del caso**.
3. Documento: input con debounce 350 ms e ILIKE parcial en backend.
4. “Relación con entidad” = `objetivo`; “Última acción” = `last_action` (texto libre).
5. Paginación server-side, 5 por página, `ORDER BY updated_at DESC`.
