# `Barrera Detalle` — Interfaz de la Pantalla

Vista interna de una barrera: información general, estado, tareas pendientes/completadas y timeline de actuaciones.

> ⚠️ **Pantalla en desarrollo temprano.** Actualmente renderiza datos mock (hardcoded en Vue data). No hay llamadas a API desde el frontend.

## Archivos relevantes

| Archivo | Rol |
|---|---|
| `src/frontend/html/salvia/barriers/barrera_detalle.html` | Template principal — lógica Vue y estructura completa |
| `src/frontend/css/barrera_detalle.css` | Estilos de layout y estados visuales |
| `src/salvia/facades/BarreraDetalleFacade.go` | Fachada Go que renderiza el template e inyecta barrierICode |

---

## Árbol de interfaz

```
BarreraDetalle  (.brd-container)
│
├── Breadcrumb  (.brd-breadcrumb)
│   └── BtnBack  (.brd-back-link)  "← {barrera.victimName}"  → history.back()
│
├── HeaderCard  (.brd-header-card)
│   ├── Avatar  (.brd-victim-avatar)  inicial del nombre de la víctima
│   └── VictimInfo  (.brd-victim-info)
│       ├── NameRow  (.brd-victim-name-row)
│       │   ├── Name  (.brd-victim-name)  barrera.victimName
│       │   └── TagPrioridad  (.brd-prioridad-{slug})  barrera.priority
│       └── Meta  (.brd-victim-meta)  "{barrera.age} años · {barrera.location}"
│
├── Tabs  (.brd-tabs-wrapper)
│   ├── TabBtn "⚠️ Información General"  :class active si activeTab === 'info'
│   └── TabBtn "📋 Timeline"  :class active si activeTab === 'timeline'
│
├── [v-if activeTab === 'info']
│   TabContentInfo  (.brd-tab-content-card)
│   │
│   ├── SectionHeader "Información general"
│   ├── Divider
│   │
│   ├── ToggleEstado  (.brd-info-subsection)
│   │   └── Toggle  (.brd-toggle-switch)  :class on si barrera.active
│   │       @click → barrera.active = !barrera.active
│   │       Label: "Activa" | "Inactiva"
│   │
│   ├── Divider
│   │
│   ├── InfoRow  (.brd-info-row)
│   │   ├── FechaIdentificacion "📅 Identificada el {barrera.identifiedAt}"
│   │   └── TagStatus  (.brd-status-{slug})  barrera.status
│   │
│   ├── Divider
│   │
│   ├── BarrierDetailRow  (.brd-barrier-detail-row)
│   │   ├── FieldSector: TagSector  (.brd-sector-{slug})  barrera.sector
│   │   ├── FieldEntidad: barrera.org
│   │   └── FieldDescripcion (.brd-detail-full): barrera.description
│   │
│   ├── Divider
│   │
│   ├── TareasPendientes  (.brd-tasks-block)
│   │   ├── [v-if pendingTasks.length > 0] TaskItem × N  [v-for pendingTasks]
│   │   │   └── icon dot pending · task.label
│   │   └── [v-else] "No hay tareas pendientes"
│   │
│   └── TareasCompletadas  (.brd-tasks-block)
│       ├── [v-if completedTasks.length > 0] TaskItem × N  [v-for completedTasks]
│       │   └── icon check-circle completed · task.label
│       └── [v-else] "No hay tareas completadas aún"
│
└── [v-if activeTab === 'timeline']
    TabContentTimeline  (.brd-tab-content-card)
    └── EmptyState  (.brd-empty-timeline)
        "No hay eventos registrados en el timeline de esta barrera."
```

---

## Datos actuales (mock)

La pantalla usa datos hardcoded en `data()`. No hay carga desde API.

| Campo | Valor mock |
|---|---|
| victimName | 'María González' |
| age | 32 |
| location | 'San Salvador, San Salvador' |
| priority | 'Alto' |
| sector | 'Justicia' |
| org | 'Policía Nacional' |
| description | 'Negación de recibir denuncia formal' |
| status | 'Por articular' |
| identifiedAt | '15 may. 2026' |
| active | true |
| pendingTasks / completedTasks | [] (vacíos) |
