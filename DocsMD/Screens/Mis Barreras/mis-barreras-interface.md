# `Mis Barreras` — Interfaz de la Pantalla

Listado de barreras activas asignadas al Enlace territorial. Cuatro columnas (Caso, Barrera, Tarea, Acción). Gestionar abre `case-task-modal` en tipo `articulacion_territorial`.

---

## Archivos relevantes

| Archivo | Rol |
|---|---|
| `src/frontend/html/salvia/barriers/barrieris.html` | Template Vue |
| `src/frontend/js/components/barriers.js` | Lógica Vue (E01–E06) |
| `src/frontend/css/barriers.css` | Estilos (`mb-`) |
| `src/frontend/js/components/case-task-modal.js` | Modal Gestionar barrera |
| `src/salvia/facades/MisBarrerasFacade.go` | GET pantalla — `get_mis_barreras` |
| `src/salvia/controller/barrier_v2_gin_controller.go` | `GET /api/v1/barriers-v2/mine` |
| `src/salvia/controller/case_task_controller.go` | `POST /api/v1/case-tasks/articulacion-territorial` |
| `src/internal/repository/barrier_v2_repository.go` | SQL listado mine |
| `src/internal/models/barrier_v2_mine_view.go` | DTO de fila |

**Ruta:** `GET /salvia/mis-barreras`  
**Rol:** `en`

---

## Árbol de interfaz

```
Pantalla: Mis Barreras
│
├── Header  (.mb-header)
│   ├── Título  "Mis Barreras"
│   └── Subtítulo  "Barreras asignadas al enlace territorial"
│
├── [v-if isLoading] Spinner  "Cargando barreras..."
│
├── [v-if successToast] Toast
│
├── [v-if loadError && !isLoading] Error
│
└── [v-if !isLoading && !loadError]
    ├── FilterCard  (.mb-filters.mb-filters-3)
    │   ├── Campo "Número de identidad"  → filters.identidad  [E02]
    │   ├── Campo "Entidad"  → filters.entidad  [E03]
    │   └── Select "Ciudad"  → filters.cityId  [E06]  opciones GET /locations/cities
    │
    ├── Table  (.mb-table-wrapper)
    │   ├── Header  (.mb-grid-4)  Caso | Barrera | Tarea (solicitud) | Acción
    │   ├── [v-if paginatedBarreras.length > 0] Row × N  (.mb-row.mb-grid-4)
    │   │   ├── ColCaso  nombre, "Doc: … · Tel: …", pin + ciudad, "Ver caso →"
    │   │   ├── ColBarrera  sector, descripción, "Ver barrera →"
    │   │   ├── ColTarea  [v-if taskId] tipo + descripción  (vacía si no hay ToDo)
    │   │   └── ColAccion  botón "Gestionar"  → openGestionar()  [E05]
    │   └── [v-else] Empty  "No se encontraron barreras…"
    │
    ├── [v-if numPages > 1] Paginador  [E04]
    │
    └── case-task-modal  ref="taskModal"  @completed="onTaskCompleted"
```
