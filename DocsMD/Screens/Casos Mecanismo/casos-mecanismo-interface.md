# `Casos Mecanismo` — Interfaz de la Pantalla

Listado de casos escalados al mecanismo articulador del usuario `mc`. Filtros por documento y ciudad (inputs de texto). Cuatro columnas: Caso, Compromiso de Salvia, Último comentario del mecanismo, Acción.

> **Estado (sep 2026):** implementado — recorte por `assigned_user_id` de sesión + `status = active`. Login `mc` → esta ruta.

---

## Archivos relevantes

| Archivo | Rol |
|---|---|
| `src/frontend/html/salvia/casos-mecanismo/casos_mecanismo.html` | Template Vue |
| `src/frontend/js/components/casos-mecanismo.js` | Lógica Vue (E01–E06) |
| `src/frontend/css/casos-mecanismo.css` | Estilos (`cme-`) |
| `src/frontend/js/components/case-task-modal.js` | Modal Gestionar (`openMecanismo`) |
| `src/salvia/facades/CasosMecanismoFacade.go` | Facade GET — `get_casos_mecanismo` |
| `src/salvia/controller/case_articulating_mechanism_controller.go` | `GET /api/v1/mecanismos/mine` |
| `src/salvia/controller/case_task_controller.go` | `POST /api/v1/case-tasks/gestion-mecanismo` |
| `src/internal/models/case_articulating_mechanism.go` | Modelo + DTO `CaseMechanismMineItem` |
| `src/internal/repository/case_articulating_mechanism_repository.go` | `FindMinePaginated` |

**Ruta:** `GET /salvia/casos-mecanismo`  
**Rol:** `mc`

---

## Árbol de interfaz

```
Pantalla: Casos Mecanismo
│
├── [v-if loadError] Bloque de error global
│   └── Mensaje + "Reintentar"  → reloadScreen()
│
└── [v-else]
    ├── welcome-section  (.cme-welcome)
    │   ├── Título  <h2><b>Casos Mecanismo</b></h2>
    │   └── Subtítulo  "Casos escalados a tu instancia del mecanismo articulador"
    │
    ├── FilterCard  (.cme-filters)
    │   ├── Campo "Documento de identidad"
    │   │   └── <input type="text">  (.cme-input)
    │   │       placeholder "Ej: 1023456780"
    │   │       → filters.document   [E02]
    │   └── Campo "Ciudad"
    │       └── <input type="text">  (.cme-input)
    │           placeholder "Ej: San Salvador"
    │           → filters.city       [E03]
    │
    ├── [v-if isLoading] Spinner
    │
    ├── [v-else-if items.length === 0] EmptyState
    │   └── "No hay casos asignados a tu mecanismo con los filtros actuales."
    │
    └── [v-else] ListCard  (.cme-list)
        ├── ListHeader
        │   ├── Col  "CASO"
        │   ├── Col  "COMPROMISO DE SALVIA"
        │   ├── Col  "ÚLTIMO COMENTARIO DEL MECANISMO"
        │   └── Col  "ACCIÓN"
        │
        ├── CaseRow × N  [v-for items]
        │   ├── ColCaso
        │   │   ├── VictimName
        │   │   ├── Identifiers  "Doc: {document} · Tel: {phone}"
        │   │   ├── Location  (ciudad de la víctima)
        │   │   ├── Badges  risk + status + scope (Nacional / Departamental / Municipal)
        │   │   └── Link  "Ver caso →"  → goToCase()   [E05]
        │   │
        │   ├── ColCompromiso
        │   │   └── salviaCommitment | "—"
        │   │
        │   ├── ColComentario
        │   │   ├── lastComment | "Sin comentarios aún"
        │   │   └── lastCommentAt  (fecha debajo, si hay)
        │   │
        │   └── ColAccion
        │       └── Btn  "Gestionar"  → openGestionar()   [E06]
        │
        └── Pagination   [v-if totalPages > 1]  → changePage()   [E04]
```

---

## Notas de UI

1. Sin dropdowns de filtro: documento y ciudad son texto libre (ILIKE).
2. Debounce ~350 ms en ambos filtros.
3. Paginación: 5 por página, `ORDER BY cam.created_at DESC`.
4. Badge de alcance según `scope`: `municipal` / `departamental` / `nacional`.
5. “Último comentario” = `case_task.result` de la última tarea Done ligada al mecanismo.
6. El i_code del caso no se muestra; sigue en el JSON (`caseCode`) para **Ver caso**.
7. El modal de gestión es `case-task-modal` (tipo `gestion_mecanismo`); ver [case-task-modal-interface.md](../../Componentes/case-task-modal/case-task-modal-interface.md).
