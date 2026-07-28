# `Barrera Detalle` — Interfaz de la Pantalla

Vista interna de una barrera: información general, tareas (pendientes/completadas, incluyendo gestión propia del Enlace Territorial) y timeline de actuaciones. Carga sus datos desde `GET /api/v1/barriers-v2/:id/detail`.

## Archivos relevantes

| Archivo | Rol |
|---|---|
| `src/frontend/html/salvia/barriers/barrera_detalle.html` | Template principal — lógica Vue y estructura completa |
| `src/frontend/css/barrera_detalle.css` | Estilos de layout y estados visuales |
| `src/salvia/facades/BarreraDetalleFacade.go` | Fachada Go que renderiza el template e inyecta sesión (rol, departamento asignado) |
| `src/frontend/js/components/case-tasks.js` | Componente reutilizable — tab "Tareas": listado pendientes/completadas, modal de gestión propia |
| `src/frontend/js/components/case-task-modal.js` | Modal para completar una tarea pendiente (montado dentro de `case-tasks`) |
| `src/frontend/js/components/case-task-history.js` | Modal de solo lectura para ver una tarea completada (montado dentro de `case-tasks`) |
| `src/frontend/js/components/case-timeline.js` | Componente reutilizable — tab "Timeline" |
| `src/frontend/js/components/barrier-follow-up-timeline.js` | Timeline de seguimientos — tab "Información General" |

---

## Árbol de interfaz

```
BarreraDetalle  (.brd-container)
│
├── [v-if cargando]
│   └── Spinner  "Cargando información de la barrera..."
│
├── [v-if error && !cargando]
│   └── ErrorBox  "No se pudo cargar la barrera"  ${ error }  → BtnVolver "← Volver"
│
└── [v-if barrera && !cargando]
    │
    ├── Breadcrumb  (.brd-breadcrumb)
    │   ├── BtnBack  (.brd-back-link)  "← {victimName}"  → history.back()
    │   └── Sector  barrera.sector
    │
    ├── HeaderCard  (.brd-header-card)
    │   ├── Avatar  (.brd-victim-avatar)  inicial de victimName
    │   └── VictimInfo  (.brd-victim-info)
    │       ├── NameRow: Name  victimName  ·  BadgeRiesgo  labelRiesgo(riskLevel)
    │       └── Meta  "{victimAge} años · {location}"
    │
    ├── Tabs  (.brd-tabs-wrapper)
    │   ├── TabBtn "⚠️ Información General"  :class active si activeTab === 'info'
    │   ├── TabBtn "📋 Tareas"  :class active si activeTab === 'tareas'
    │   └── TabBtn "🕐 Timeline"  :class active si activeTab === 'timeline'
    │
    ├── [v-if activeTab === 'info']
    │   TabContentInfo  (.brd-tab-content-card)
    │   ├── InfoRow: FechaIdentificacion "📅 Identificada el {createdAt}"  ·  TagStatus  labelBarreraStatus(status)
    │   ├── BarrierDetailRow  sector, entidad(es), dependencia, identificada por, agente responsable,
    │   │   fecha, [v-if enlaceActivado] "Enlace territorial ✓ Activado", barreras identificadas, descripción...
    │   ├── [v-if ubicación]  Ubicación de la barrera
    │   ├── [v-if barreras estructurales]  Chips por categoría (institucional/económica/territorial/diferencial)
    │   ├── SeguimientosBlock  (.brd-tasks-block)
    │   │   ├── BtnVerRegistroCompleto "📋 Ver registro completo"  → mostrarRegistroCompleto = true
    │   │   └── <barrier-follow-up-timeline>  // src/frontend/js/components/barrier-follow-up-timeline.js
    │   └── [v-if mostrarRegistroCompleto]
    │       └── ModalRegistroCompleto  — desglose de las 22 preguntas del registro de la barrera (solo lectura)
    │
    ├── [v-if activeTab === 'tareas']
    │   TabContentTareas  (.brd-tab-content-card)
    │   └── <case-tasks :case-id :barrier-id :user-id :user-role :puede-gestion-propia>
    │       // src/frontend/js/components/case-tasks.js
    │       ├── Tabs internas "Pendientes (N)" / "Completadas (N)"
    │       ├── [pendingTask] BtnGestionar  → abre <case-task-modal>
    │       ├── [completedTask] BtnVerDetalle  → abre <case-task-history>
    │       └── [puedeGestionPropia]  ← rol "en" Y barrera.departmentId === currentUserAssignedDepartment
    │           ├── BtnRegistrarGestionPropia  "+ Registrar gestión propia"  → abrirModalGestionPropia()
    │           └── [modalGestionPropiaAbierto]
    │               ModalGestionPropia  (.ct-modal-backdrop)
    │               ├── <select> Tipo de gestión*  [Llamada | Visita presencial | Oficio a entidad | Otra gestión]
    │               ├── <textarea> Descripción de la gestión*
    │               ├── [errorGestionPropia]  ErrorMsg
    │               ├── BtnCancelar  → cerrarModalGestionPropia()
    │               └── BtnGuardar  "Registrar como completada"  :disabled si !descripcion.trim()
    │                   → confirmarGestionPropia()  →  POST /api/v1/case-tasks/gestion-propia
    │
    └── [v-if activeTab === 'timeline']
        TabContentTimeline  (.brd-tab-content-card)
        └── <case-timeline :case-id :barrier-id :show-filters="false">
            // src/frontend/js/components/case-timeline.js
```

---

## Acceso por rol y departamento

El botón "+ Registrar gestión propia" (tab Tareas) solo se muestra si:
- `currentRole === 'en'` (Enlace Territorial), **y**
- `barrera.departmentId === currentUserAssignedDepartment` — el departamento asignado al Enlace (`security.general_user.general_user_assigned_department`, independiente de su municipio de residencia).

El backend revalida ambas condiciones en `POST /api/v1/case-tasks/gestion-propia` — el chequeo del frontend es solo para UX, no seguridad.
