# `Casos Entidad` — Interfaz de la Pantalla

Listado de casos referidos a la **sede** (`entity_branch`) del usuario `et`. El header muestra la organización padre (`entity`). Solo filtro por documento + paginación.

> **Estado (jul 2026):** implementado — sede desde `general_user.entity_branch_id` / sesión; sin selector de entidad ni filtro de ciudad. Header con entidad padre + sector vía respuesta de listado.

---

## Archivos relevantes

| Archivo | Rol |
|---|---|
| `src/frontend/html/salvia/entity_case/casos_entidad.html` | Template Vue |
| `src/frontend/js/components/casos-entidad.js` | Lógica Vue (E01, E02, E04, E05, E06) |
| `src/frontend/css/casos-entidad.css` | Estilos (`ce-`) |
| `src/salvia/facades/CasosEntidadFacade.go` | Facade GET — `get_casos_entidad` |
| `src/salvia/controller/entity_case_controller.go` | `GET /api/v1/entity-cases` (ajustar a sede de sesión) |
| `src/security/dao/GeneralUserDAO.go` | Campo nuevo `entity_branch_id` |
| `src/common/utils/CommonSession.go` | Exponer `EntityBranchId` en sesión |
| `src/internal/models/entity_case.go` | Modelo + DTOs |
| `src/internal/repository/entity_case_repository.go` | Listado por `entity_branch_id` |

**Ruta:** `GET /salvia/casos-entidad`  
**Rol:** `et`

---

## Árbol de interfaz (objetivo)

```
Pantalla: Casos Entidad
│
├── [v-if loadError] Bloque de error global
│   └── Mensaje + "Reintentar"  → reloadScreen()
│
└── [v-else]
    ├── welcome-section  (.ce-welcome)
    │   ├── Título  <h2><b>Casos Entidad</b></h2>
    │   ├── Subtítulo  (.text-muted)  "Casos con víctimas referidas a tu entidad"
    │   ├── EntityName  (.ce-entity-name)   parentEntity.name   // organización padre
    │   └── SectorLabel  (.ce-sector)       parentEntity.sectorName   (ej. "Justicia")
    │       // SIN selector de entidad
    │
    ├── FilterCard  (.ce-filters)
    │   └── Campo "Documento de identidad"   // único filtro
    │       └── <input type="text">  (.ce-input)
    │           placeholder "Ej: 1023456780"
    │           → filters.document   [E02]
    │       // SIN filtro de ciudad
    │
    ├── [v-if isLoading] Spinner
    │
    ├── [v-else-if items.length === 0] EmptyState
    │   └── "No hay casos referidos a tu entidad con los filtros actuales."
    │
    └── [v-else] ListCard  (.ce-list)
        ├── ListHeader
        │   ├── Col  "CASO"
        │   └── Col  "RELACIÓN CON ENTIDAD"
        │
        ├── CaseRow × N  [v-for items]
        │   ├── ColCaso
        │   │   ├── VictimName
        │   │   ├── Identifiers  "{caseCode} · {document}"
        │   │   ├── Location  (ciudad del caso — solo display)
        │   │   ├── Badges  risk + status
        │   │   └── Link  "Ver caso →"  → goToCase()   [E05]
        │   │
        │   └── ColRelacion
        │       ├── RelationDesc  item.objetivo
        │       ├── LastAction  item.lastAction | "Sin acciones registradas"
        │       └── Link  "Ver detalle →"  → goToRelationDetail()   [E06]
        │
        └── Pagination   [v-if totalPages > 1]  → changePage()   [E04]
```

---

## Notas de UI

1. **Sin** dropdown de entidad: la sede sale de `session.EntityBranchId`.
2. **Sin** filtro de ciudad: la sede ya está ligada a un `town_code`.
3. Header: nombre de la **entidad padre** (`entity.entity_name`) + sector.
4. Documento: debounce ~350 ms, match parcial ILIKE.
5. Paginación: 5 por página, `ORDER BY updated_at DESC`.
