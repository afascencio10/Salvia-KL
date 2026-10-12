# `case-gestion-institucional` — Interfaz del Componente

Mini-dashboard con el resumen de gestión institucional de un caso: cantidad de entidades agregadas, % de obligaciones completadas (global y por sector), barreras sin resolver y oficios en gestión. Cada indicador es clickeable y abre un modal con el detalle tabular filtrado. Puerto 1:1 del componente `ResumenGestionInstitucional.jsx` del mockup React (`Salvia/src/components/detalle/ResumenGestionInstitucional.jsx`), adaptado a Vue 3 y a datos reales del backend en vez de mocks.

Se monta dentro del tab **"Gestión institucional"** (`tabActiva === 'barreras'`) de la pantalla Detalle del Caso, **antes** de `<case-entities>` — en la parte superior del tab, igual que en el mockup.

**Permisos:** mismos que el resto del tab — cualquier rol con acceso a `get_case_detail_sv` (todos menos `ad`). Es un componente de solo lectura, no requiere permisos de escritura adicionales.

## Archivos relevantes

| Archivo | Rol |
|---|---|
| `src/frontend/js/components/case-gestion-institucional.js` | Componente principal — template y lógica Vue *(nuevo, por crear)* |
| `src/frontend/html/salvia/case_detail/get_case_detail_sv.html` | Monta `<case-gestion-institucional>` al inicio del tab "Gestión institucional" |
| `src/salvia/controller/entity_case_controller.go` | Reutilizado — `GET /api/v1/casos/:caseId/entidades`, `GET /api/v1/entity-cases/:entityCaseId/obligations` |
| `src/salvia/controller/entity_letter_controller.go` *(a confirmar nombre exacto)* | Reutilizado — `GET /api/v1/entity-letters?caseId=` |
| `src/internal/models/barrier_v2.go` | `EntityBranchID`/`Status`/`Sector` — ya poblados por el registro de barrera (ver sesión donde se conectó `entity_branch_id`) |

---

## Supuestos de diseño (a confirmar en revisión)

- **Nombre del componente/tag**: `case-gestion-institucional` — sigue la convención `case-*` de los hermanos (`case-entities`, `case-oficios`, `case-tasks`) y coincide con el label real del tab ("Gestión institucional"). Ajustable si prefieres otro nombre.
- **Props**: `caseId` (String, requerido) — igual que `case-entities`/`case-oficios` — y `barriers` (Array, requerido) — **reutiliza el array `barriers` que el padre (`get_case_detail_sv.html`) ya carga** para la lista de barreras del tab, igual que hace el mockup (`<ResumenGestionInstitucional casoId barreras />`). Evita duplicar el fetch de barreras.
- **Entidades**: fetch propio a `GET /api/v1/casos/:caseId/entidades` (mismo endpoint que ya usa `case-entities`) — de ahí sale `entidades.length` para el KPI de Entidades agregadas y el detalle del modal `tipo: 'entidades'`.
- **Obligaciones completadas — endpoint agregado nuevo (decisión confirmada):** `GET /api/v1/casos/:caseId/obligaciones-resumen`, que calcula en una sola consulta SQL (join `entity_case` → `entity_branch` → `entity` → `entity_case_obligation` → `entity_obligation`) el agregado global, por sector, y el detalle fila a fila para el modal. Una sola llamada de red, sin importar cuántas entidades tenga el caso. Ver shape completo en `case-gestion-institucional-usage.md`.
- **Oficios**: fetch propio a `GET /api/v1/entity-letters?caseId=` (mismo endpoint que ya usa `case-oficios`) — de ahí sale el bucket "En gestión" (todo estado ≠ `radicado`/`respondido`) vs `Radicado` vs `Respondido`.
- **Barreras por estado**: se calculan client-side sobre la prop `barriers` — `OPEN` → "Abierta", `En Gestion` → "En gestión", cualquier otro (`Articulada`/`MANAGED`/`CLOSED`/`FINISH`/`CLOSE`/`RESOLVED`) → "Resuelta" — mismo mapeo que ya usa `barreraEsResuelta()` en `get_case_detail_sv.html`.
- **Modal de detalle**: confirmado — se incluye desde la v1. Réplica del `ModalDetalleIndicador.jsx` del mockup: lista las filas detrás de cada número clickeado (ver árbol abajo).
- **Umbrales de color (RAG)** para el % de obligaciones: igual que el mockup — `<40%` rojo, `40–79%` ámbar, `≥80%` verde.

---

## Árbol de interfaz

```
case-gestion-institucional  (.cgi-card)
│
├── Header  (.cgi-header)
│   ├── Icono  "📊"
│   └── Titulo  "Resumen de Gestión Institucional"
│
├── [v-if cargando]
│   └── Spinner  "Cargando resumen..."
│
├── [v-else-if error]
│   └── ErrorMsg  (.cgi-error)  error  +  BtnReintentar  → cargar()
│
└── [v-else]
    │
    ├── KPIs  (.cgi-kpis, grid 2×2 / 4×1)
    │   ├── StatTile  (.cgi-tile.violet)  "🏛️ Entidades agregadas"  entidades.length
    │   │   → abrirModal('entidades')
    │   ├── StatTile  (.cgi-tile)  "📋 Obligaciones completadas"  "{completadas}/{total} · {pct}%"
    │   │   :tone="toneObligaciones"  (verde/ámbar/rojo/gris según pct)
    │   │   [v-if total === 0]  "Sin obligaciones"
    │   │   → abrirModal('obligaciones')
    │   ├── StatTile  (.cgi-tile)  "⚠️ Barreras sin resolver"  barrerasSinResolver
    │   │   :tone="barrerasSinResolver > 0 ? 'red' : 'green'"
    │   │   → abrirModal('barreras')
    │   └── StatTile  (.cgi-tile)  "📄 Oficios en gestión"  oficiosPorBucket['En gestión']
    │       :tone="oficiosEnGestion > 0 ? 'amber' : 'green'"
    │       → abrirModal('oficios', 'En gestión')
    │
    ├── SeccionSectores  (.cgi-sectores)
    │   ├── Label  "% OBLIGACIONES COMPLETADAS POR SECTOR"
    │   └── Grid  (.cgi-sectores-grid)
    │       └── SectorItem × 3  [v-for ['Salud','Justicia','Protección']]  (.cgi-sector-item)
    │           ├── CircularProgress  (SVG, :pct :tone)
    │           │   Centro: "{pct}%" | "—" si sin obligaciones
    │           ├── NombreSector
    │           ├── CantidadEntidades  "{n} entidad(es)"
    │           └── DetalleObligaciones  "{completadas}/{total} obligaciones" | "Sin obligaciones"
    │           → abrirModal('obligaciones', sector)
    │
    └── SeccionEstados  (.cgi-estados, grid 2 columnas, border-top)
        ├── ColBarreras
        │   ├── Label  "BARRERAS POR ESTADO"
        │   └── Badges  (.cgi-badges)
        │       └── Badge × 3  [v-for ['Abierta','En gestión','Resuelta']]  (.cgi-badge)
        │           "{estado}: {count}"
        │           → abrirModal('barreras', estado)
        └── ColOficios
            ├── Label  "OFICIOS POR ESTADO"
            └── Badges  (.cgi-badges)
                └── Badge × 3  [v-for ['En gestión','Radicado','Respondido']]  (.cgi-badge)
                    "{bucket}: {count}"
                    → abrirModal('oficios', bucket)

─── MODAL: DETALLE DE INDICADOR ────────────────────────────────────────────────

[v-if modal]  Overlay  (.cgi-modal-overlay)  @click.self="cerrarModal"
  └── ModalBox  (.cgi-modal)
      ├── Header
      │   ├── Titulo  modal.titulo
      │   ├── Subtitulo  "{filas.length} resultado(s)"
      │   └── BtnCerrar  "×"  → cerrarModal()
      │
      ├── Body  (.cgi-modal-body, scroll)
      │   ├── [v-if filas.length === 0]
      │   │   EmptyState  "No hay registros para este indicador"
      │   │
      │   └── [v-else]
      │       Tabla  (.cgi-modal-table)
      │       ├── Thead  — columnas según modal.tipo
      │       │   [tipo === 'entidades']    Sector | Entidad | Ubicación | Oficios | Barreras activas
      │       │   [tipo === 'obligaciones'] Entidad | Sector | Obligación | Estado
      │       │   [tipo === 'barreras']     Sector | Estado | Descripción | Fecha
      │       │   [tipo === 'oficios']      Entidad | Estado | Fecha
      │       └── Tbody
      │           └── Fila × N  [v-for filas]  (contenido según modal.tipo)
      │
      └── Footer
          └── BtnCerrar  "Cerrar"  → cerrarModal()
```

---

## Estados visuales de los tiles (StatTile)

| Tone | Uso |
|---|---|
| `violet` | Entidades agregadas (siempre) |
| `green` | % ≥ 80, barreras sin resolver = 0, oficios en gestión = 0 |
| `amber` | % entre 40 y 79 |
| `red` | % < 40, barreras sin resolver > 0 |
| `gray` | Sin datos (0 obligaciones totales) |

## Estados de barrera → bucket visual

| `barrier_v2.status` | Bucket mostrado |
|---|---|
| `OPEN` | Abierta |
| `En Gestion` | En gestión |
| `Articulada` / `MANAGED` / `CLOSED` / `FINISH` / `CLOSE` / `RESOLVED` | Resuelta |

## Estados de oficio → bucket visual

| `entity_letter.state` | Bucket mostrado |
|---|---|
| `radicado` | Radicado |
| `respondido` | Respondido |
| cualquier otro (`por_proyectar`, `para_revisar`, `en_correccion`, `aprobacion_juridica`, `para_radicar`) | En gestión |
