# `case-entities` — Interfaz del Componente

Lista las entidades institucionales (`entity_branch`) relacionadas manualmente con un caso, en formato de cards, con filtros (buscador de texto, sector, y "¿con barreras activas?"). Incluye un modal para agregar una nueva entidad al caso. No tiene funcionalidad de quitar entidades en esta versión.

Montado dentro del tab **"Gestión institucional"** (antes "Barreras") de la pantalla Detalle del Caso, antes de la sección de Barreras — mismo lugar que ocupa `SeccionEntidades` en el mockup React, y hermano de `case-oficios` (que va después de la lista de barreras).

**Permisos:**
- **Ver y filtrar entidades:** cualquier rol con acceso a la pantalla Detalle del Caso — es decir, cualquiera de los roles listados en `salvia_config.PermissionsByRole["get_case_detail_sv"]` (`op`, `ro`, `do`, `et`, `us`, `sv`, `no`, `fo`, `ps`, `ts`, `an` — todos menos `ad`). El backend reutiliza ese mismo permiso en vez de mantener una lista de roles propia, para que ambos permisos no se desincronicen.
- **Agregar una entidad:** solo `sv`, `op`, `ro` — el botón "+ Agregar entidad" no se renderiza para el resto (`v-if="puedeAgregar"` en el frontend) y el backend rechaza la petición con 403 aunque se llame directamente.

## Archivos relevantes

| Archivo | Rol |
|---|---|
| `src/frontend/js/components/case-entities.js` | Componente principal — template y lógica Vue |
| `src/frontend/html/salvia/case_detail/get_case_detail_sv.html` | Monta `<case-entities>` dentro del tab "Gestión institucional" |
| `src/salvia/controller/entity_case_controller.go` | Endpoints REST + guard de permisos (`requireCaseDetailAccess`, `entityCaseWriteRoles`) |
| `src/salvia/service/entity_case_service.go` | Lógica de negocio (duplicados, etc.) |
| `src/internal/repository/entity_case_repository.go` | Queries a `entity_case` + joins |
| `src/internal/models/entity_case.go` | Modelo — ver `related-tables.md` |
| `src/internal/models/barrier_v2.go` | Incluye `EntityBranchID *int64` nullable — ver `related-tables.md` |
| `src/salvia/dao/EntityBranchDAO.go` | DAO existente reutilizado para resolver datos de la sede |

---

## Supuestos de diseño (a confirmar en revisión)

- Recibe `caseId` como prop (String, requerido) y hace un fetch al montar: `GET /api/v1/casos/:caseId/entidades`.
- La ubicación de cada entidad viene de `entity_branch.entity_branch_town_code` (resuelta a municipio/ciudad/departamento por el backend), **no** de campos de texto libre como en el mockup original.
- El modal "+ Agregar entidad" no usa un catálogo hardcodeado de nombres (como el mockup): reutiliza el endpoint ya existente `GET /api/v1/entity-branches?town_code=&sector=` para que el agente busque sedes reales ya registradas en el sistema. Requiere elegir primero ubicación (departamento → ciudad → municipio, mismo selector en cascada que ya usa el registro de caso) **y sector — ambos obligatorios** — antes de ver resultados.
- El campo "Oficios con esta entidad" se calcula en el backend contando `entity_letter` por `case_id` + `entity_branch_id`.
- El campo "Barreras activas de esta entidad" se calcula en el backend contando `barrier_v2` por `case_id` + `entity_branch_id` (status `OPEN` o `En Gestion`) — habilitado gracias al campo `entity_branch_id` agregado a `barrier_v2` (ver `related-tables.md`). El filtro "¿Con barreras activas?" ya no está bloqueado.
- El campo "Última acción" (`lastAction`) es un texto libre almacenado directamente en `entity_case.last_action` — no se calcula. Queda pendiente definir en un siguiente ciclo el mecanismo exacto para poblarlo/editarlo (ver GAP en `related-tables.md` y en este archivo).
- **No existe funcionalidad de "quitar entidad" en esta versión** — una vez agregada, la entidad permanece asociada al caso. Si más adelante se necesita, se reactivaría usando el soft-delete (`deleted_at`) que ya trae el modelo `entity_case`.
- "Ver Detalle" navega a `/salvia/entidad/:id` (ruta nueva — no existe hoy en el código real, ver GAP en `case-entities-index.md`).

---

## Árbol de interfaz

```
case-entities
│
├── [v-if cargando]
│   └── Spinner  "Cargando entidades..."
│
├── [v-else-if error]
│   └── ErrorMsg  (.ce-error)  error  +  BtnReintentar  → cargar()
│
└── [v-else]
    │
    ├── Filtros  (.ce-filtros)
    │   ├── <select> Sector  (.ce-select-sector)
    │   │   opciones: Todos | Salud | Justicia | Protección
    │   │
    │   ├── <input> Buscador  "Buscar por nombre de entidad..."  (.ce-buscador)
    │   │
    │   ├── <select> ¿Con barreras activas?  (.ce-select-barreras)
    │   │   opciones: Todos | Sí | No
    │   │
    │   └── [v-if puedeAgregar]  BtnAgregar  "+ Agregar entidad"  → modalAgregarAbierto = true
    │       puedeAgregar: userRole === 'sv' || 'op' || 'ro'
    │
    ├── [v-if entidadesFiltradas.length === 0]
    │   └── EmptyState  (.ce-empty)
    │       [entidades.length === 0]  "No hay entidades registradas para este caso"
    │       [v-else]  "No hay entidades que coincidan con los filtros"
    │
    └── [v-else]
        └── Grid  (.ce-grid)
            └── EntidadCard × N  [v-for entidadesFiltradas]  (.ce-card)
                ├── Header
                │   ├── BadgeSector  (.ce-badge-sector)  entidad.sector
                │   └── Nombre  (.ce-nombre)  entidad.entityBranchName
                ├── Ubicacion  (.ce-ubicacion)  📍 departamento · ciudad · municipio  (resueltos vía town_code)
                ├── [v-if entidad.address]  Direccion  (.ce-direccion)  entidad.address
                ├── [v-if entidad.objetivo]  Objetivo  (.ce-objetivo)  texto truncado (line-clamp-2)
                ├── StatsRow  (.ce-stats)
                │   ├── OficiosCount  "Oficios: N"
                │   ├── BarrerasActivasCount  "Barreras activas: N"
                │   └── UltimaAccion  [v-if entidad.lastAction]  "Última acción: {lastAction}"  |  [v-else]  "Sin acciones registradas"
                └── Actions  (.ce-actions)
                    └── BtnVerDetalle  "Ver Detalle →"  → onVerDetalleEntidad(entidad)

─── MODAL: AGREGAR ENTIDAD ────────────────────────────────────────────────────

[v-if modalAgregarAbierto]  Overlay  (.ce-modal-overlay)  @click.self="cerrarModalAgregar"
  └── ModalBox  (.ce-modal)
      ├── Header
      │   ├── Titulo  "Agregar entidad al caso"
      │   └── BtnCerrar  "×"  → cerrarModalAgregar()
      │
      ├── Body  (.ce-modal-body)
      │   ├── <select> Departamento*  (.ce-select-depto)  → onCambiaDepartamento()
      │   ├── <select> Ciudad*  (.ce-select-ciudad)  :disabled si !departamento  → onCambiaCiudad()
      │   ├── <select> Municipio*  (.ce-select-municipio)  :disabled si !ciudad  → onCambiaMunicipio()
      │   ├── <select> Sector*  (.ce-select-sector-modal)  :disabled si !municipio  → onCambiaSector()
      │   │   (obligatorio — ya no es opcional; sin sector no se buscan entidades)
      │   │
      │   ├── [v-if cargandoEntidades]  Spinner  "Buscando sedes disponibles..."
      │   ├── [v-if !cargandoEntidades && sector && entidadesDisponibles.length === 0]
      │   │   EmptyState  "No hay sedes registradas en este municipio y sector. Puedes crear una desde Sedes."
      │   │
      │   ├── [v-if entidadesDisponibles.length > 0]
      │   │   <select> Entidad*  (.ce-select-entidad)  [v-for entidadesDisponibles]
      │   │
      │   ├── [v-if errores.entidad] ErrorMsg  "Selecciona una entidad"
      │   │
      │   └── <textarea> Objetivo con la entidad*  (.ce-textarea-objetivo)
      │       [v-if errores.objetivo] ErrorMsg  "Describe el objetivo con la entidad"
      │
      └── Footer
          ├── BtnCancelar  → cerrarModalAgregar()
          └── BtnGuardar  "Agregar entidad" | "Guardando..."  :disabled si !formValido || guardando  → guardarEntidad()
```
