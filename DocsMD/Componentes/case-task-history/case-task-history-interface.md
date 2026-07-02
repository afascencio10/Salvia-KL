# `case-task-history` — Interfaz del Componente

Modal de solo lectura que muestra el detalle de una tarea (`case_task`) ya completada. El padre —el componente de listado de tareas completadas (`case-tasks.js`)— lo activa vía `open(taskId)`, igual que `case-task-modal`. Este componente carga la tarea del API y renderiza su detalle según `tarea.type`, incluyendo el desglose legible de `formData`.

## Archivos relevantes

| Archivo | Rol |
|---|---|
| `src/frontend/js/components/case-task-history.js` | Componente principal — template y lógica Vue |
| `src/salvia/controller/case_task_controller.go` | Endpoint `GetByID` — enriquece la tarea con `assignedUserName` |

---

## Supuestos de diseño

- Es un **modal**, no una lista inline — se abre y se cierra con un `taskId` distinto cada vez, controlado por el padre vía método público `open(taskId)`. Mismo patrón que `case-task-modal`. **No recibe `caseId` ni ningún prop.**
- El listado de tareas completadas (qué tareas mostrar, paginación, filtros) es responsabilidad de `case-tasks.js` — componente reutilizable ya existente que lista tareas `ToDo`/`Done` de un caso. Ese componente tiene la infraestructura (`pendingTasks`/`completedTasks`/`activeTab`) pero aún no renderiza la pestaña de completadas ni dispara `open()` sobre este modal — **queda pendiente de otra persona conectar el trigger real**.
- Vista de solo lectura: sin campos editables ni botón de confirmar — solo un botón "Cerrar".
- Usa el mismo endpoint que ya consume `case-task-modal` para cargar una tarea individual: `GET /api/v1/case-tasks/:taskId`.
- ✅ **Resuelto:** `GetByID` ahora devuelve la tarea enriquecida con `assignedUserName` (antes solo lo hacía el endpoint de listado). Se refactorizó `enrichTasksWithNames` para compartir el mismo shape (`taskToJSON` + `lookupUserName`) con el nuevo `enrichTaskWithName`, evitando duplicar la lista de campos y sin reintroducir N+1 queries en el listado.

---

## Árbol de interfaz

```
case-task-history
│
├── [v-if !visible]
│   └── (no renderiza nada)
│
└── [v-if visible]  Overlay  (.cth-overlay)
    └── ModalBox  (.cth-modal)
        │
        ├── Header  (.cth-header)
        │   ├── Icono  (.cth-header-icon)  :class="iconoTipo(tarea.type)"
        │   ├── Titulo  (.cth-header-title)   tituloTipo(tarea.type)
        │   └── BtnCerrar  "✕"  → cerrar()
        │
        ├── [v-if cargando]
        │   └── Spinner  "Cargando tarea..."
        │
        ├── [v-else-if error]
        │   └── ErrorMsg  (.cth-error)  error
        │
        └── [v-else-if tarea]
            │
            ├── Body  (.cth-body)
            │   │
            │   ├── CompletadoPor  (.cth-completado-por)  "Completado por: " + tarea.assignedUserName
            │   │
            │   ├── Fecha  (.cth-fecha)  "Completado el " + formatearFecha(tarea.completedAt)
            │   │
            │   ├── [tarea.description]
            │   │   └── Descripcion  (.cth-descripcion)  tarea.description
            │   │
            │   └── DetalleFormData  (.cth-formdata)
            │       │
            │       ├── [tarea.type === 'gestion_llamada']
            │       │   ├── Campo  "Ubicación"          formData.departamentoNombre + " / " + ciudadNombre + " / " + municipioNombre
            │       │   ├── Campo  "Entidad"             formData.entidadNombre
            │       │   ├── Campo  "Funcionario"         formData.funcionario
            │       │   ├── [formData.descripcion]
            │       │   │   └── Campo  "Notas"           formData.descripcion
            │       │   ├── Campo  "¿Generó oficio?"     formData.generaOficio ? "Sí" : "No"
            │       │   └── [formData.generaOficio]
            │       │       ├── Campo  "Asunto"          formData.asunto
            │       │       └── Campo  "Ruta Kofax"      formData.rutaKofax
            │       │
            │       ├── [tarea.type === 'proyectar_oficio']
            │       │   ├── Campo  "Ubicación"           formData.departamentoNombre + " / " + ciudadNombre + " / " + municipioNombre
            │       │   ├── Campo  "Entidad"             formData.entidadNombre
            │       │   ├── Campo  "Funcionario"         formData.funcionario
            │       │   ├── Campo  "Asunto"              formData.asunto
            │       │   └── Campo  "Ruta Kofax"          formData.rutaKofax
            │       │
            │       ├── [tarea.type === 'comite_caso']
            │       │   ├── DecisionChip × N  [v-for d in formData.decisionesTexto]  (.cth-chip)
            │       │   ├── [formData.observacionesOficio]
            │       │   │   └── Campo  "Obs. oficio"             formData.observacionesOficio
            │       │   ├── [formData.observacionesRecomendaciones]
            │       │   │   └── Campo  "Obs. recomendaciones"    formData.observacionesRecomendaciones
            │       │   └── [formData.nivelMecanismo]
            │       │       ├── Campo  "Nivel del mecanismo"     formData.nivelMecanismoTexto
            │       │       └── [formData.observacionesMecanismo]
            │       │           └── Campo  "Obs. mecanismo"      formData.observacionesMecanismo
            │       │
            │       └── [tarea.type === 'Corregir oficio']
            │           └── Nota  (.cth-nota-corregido)  "El oficio fue corregido por el agente que lo proyectó"
            │
            └── Footer  (.cth-footer)
                └── BtnCerrar  "Cerrar"  → cerrar()
```

---

## Presentación de formData por tipo de tarea

| `tarea.type` | Campos que muestra |
|---|---|
| `gestion_llamada` | Ubicación (departamento/ciudad/municipio), entidad, funcionario, notas, ¿generó oficio?, y si aplica: asunto + ruta Kofax |
| `proyectar_oficio` | Ubicación, entidad, funcionario, asunto, ruta Kofax |
| `comite_caso` | Chips de decisiones tomadas + observaciones condicionales por decisión + nivel del mecanismo articulador |
| `Corregir oficio` | Sin campos propios (`formData` es `{}`) — se muestra la nota **"El oficio fue corregido por el agente que lo proyectó"** en vez de un desglose de campos |

> Depende directamente de los campos legibles (`departamentoNombre`, `ciudadNombre`, `municipioNombre`, `decisionesTexto`, `nivelMecanismoTexto`) agregados a `_buildFormData()` en `case-task-modal.js`. Sin esos campos, este componente tendría que re-resolver IDs contra los catálogos de ubicación vigentes, lo cual sería frágil para datos históricos.
