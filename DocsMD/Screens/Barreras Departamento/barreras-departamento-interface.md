# `Barreras Departamento` — Interfaz de la Pantalla

Listado de barreras del departamento asignado al Enlace territorial. Layout de dos columnas (CASO | BARRERA) alineado al mockup y a Casos Entidad. Filtros al API.

> **Estado (sep 2026):** implementado en archivos `*_v2` / `barreras_departamento.css`. La UI de 6 columnas anterior se conserva en disco y no se monta.

---

## Archivos relevantes

| Archivo | Rol |
|---|---|
| `src/frontend/html/salvia/barriers/barreras_departamento_v2.html` | Template Vue vigente |
| `src/frontend/js/components/barreras_departamento_v2.js` | Lógica Vue (E01–E08) |
| `src/frontend/css/barreras_departamento.css` | Estilos (`bd-`) |
| `src/salvia/facades/BarrerasDepartamentoFacade.go` | GET pantalla — `get_barreras_departamento` |
| `src/salvia/controller/barrier_v2_gin_controller.go` | `GET /api/v1/barriers-v2/department` |
| `src/internal/repository/barrier_v2_repository.go` | SQL listado + count |
| `src/internal/models/barrier_v2_department_view.go` | DTO de fila |
| `src/frontend/html/salvia/barriers/barreras_departamento.html` | UI anterior (no montada) |
| `src/frontend/js/components/barreras_departamento.js` | JS anterior (no montado) |

**Ruta:** `GET /salvia/barreras-departamento`  
**Rol:** `en`

---

## Árbol de interfaz (objetivo)

```
Pantalla: Barreras Departamento
│
├── [v-if loadError] Bloque de error global
│   └── Mensaje + "Reintentar"  → reloadScreen()
│
└── [v-else]
    ├── welcome-section  (.bd-welcome)
    │   ├── Título  <h2><b>Barreras Departamento</b></h2>
    │   └── Subtítulo  "Todas las barreras identificadas en el departamento"
    │
    ├── FilterCard  (.bd-filters)
    │   ├── Campo "Documento de identidad"
    │   │   └── <input>  placeholder "Ej: 1023456780"
    │   │       → filters.docNumber   [E02] debounce 350 ms
    │   ├── Campo "Entidad"
    │   │   └── <input>  placeholder "Ej: Fiscalía"
    │   │       → filters.entidad     [E03] debounce 350 ms
    │   ├── Campo "Ciudad"
    │   │   └── <select>  primera opción "Todas"
    │   │       ciudades de GET /api/v1/locations/cities?department_id=AssignedDepartmentID
    │   │       → filters.cityId     [E04] inmediato
    │   └── Campo "Estado"
    │       └── <select>  → filters.status   [E05] inmediato
    │           opciones: Todos | Abierta (OPEN) | En gestión (En Gestion)
    │                     | Articulada | Gestionada (MANAGED)
    │
    ├── [v-if isLoading] Spinner  "Cargando barreras…"
    │
    ├── [v-else-if items.length === 0] EmptyState
    │   └── "No se encontraron barreras con los filtros seleccionados."
    │
    └── [v-else] ListCard  (.bd-list)
        ├── ListHeader
        │   ├── Col  "CASO"
        │   └── Col  "BARRERA"
        │
        ├── BarrierRow × N  [v-for items]
        │   ├── ColCaso
        │   │   ├── VictimName
        │   │   ├── Meta  "Doc: {documento} · Tel: {telefono|—}"
        │   │   ├── Location  pin + cityName
        │   │   ├── Badge prioridad  (Bajo/Moderado/Alto/Extremo)
        │   │   └── Link  "Ver caso →"  → goToCase()   [E07]
        │   │
        │   └── ColBarrera
        │       ├── Tags  sector badge + entityBranchName + estado badge
        │       ├── Description
        │       ├── Date  (YYYY-MM-DD o locale)
        │       └── Link  "Ver barrera →"  → goToBarrier()   [E08]
        │
        └── Pagination   [v-if totalPages > 1]  → changePage()   [E06]
```

---

## Notas de UI

1. Departamento **nunca** viaja como query param.
2. Documento y entidad: debounce 350 ms. Ciudad y estado (select): sin debounce.
3. Paginación: 10 por página, `ORDER BY created_at DESC`.
4. Labels de estado: `OPEN` → Abierta; `En Gestion` → En gestión; `Articulada` → Articulada; `MANAGED` → Gestionada.
5. Labels de sector: `salud`, `justicia`, `proteccion`, `otras_instituciones`, `barrera_transversal`.
6. Tipografía de filas alineada a `casos-component` (nombre 0.86rem, meta 0.75rem).
7. Teléfono: form2 del caso, si vacío el del contacto; UI muestra `Tel: —` si no hay dato.
