# `Administrar Duplas` — Interfaz de la Pantalla

Pantalla para gestionar duplas activas del equipo de Atención Psicosocial: listar profesionales `ps`/`ts`, ver duplas activas, crear, editar y eliminar (lógico).

---

## Archivos relevantes

| Archivo | Rol |
|---|---|
| `src/frontend/html/salvia/duplas/administrar_duplas.html` | Template principal — Vue + layout |
| `src/frontend/css/administrar-duplas.css` | Estilos de la pantalla (prefijo `adp-`) |
| `src/salvia/facades/AdministrarDuplasFacade.go` | Facade GET — permiso `get_administrar_duplas` (solo `sv`) |
| `src/internal/models/dupla.go` | Modelo `salvia.dupla` |
| `src/internal/repository/dupla_repository.go` | Acceso a duplas (hoy solo `ListActive` para reasignación) |

**Ruta:** `GET /salvia/administrar-duplas`  
**Entrada desde:** botón “Administrar duplas” en Historial de Remisiones (`/salvia/historial-remisiones`)

---

## Árbol de interfaz

```
Pantalla: Administrar Duplas
│
├── [v-if loadError] Bloque de error global
│   └── Mensaje + acción opcional "Reintentar" → reloadScreen()
│
├── [v-else]
│   ├── Header
│   │   ├── Título  "Administrar duplas"
│   │   ├── Subtítulo  "Cada profesional puede pertenecer a una sola dupla a la vez"
│   │   └── Botón primaro "+ Nueva dupla"  → openDuplaModal('create')   [E02]
│   │
│   ├── [v-if isLoading] Skeleton / spinner de carga
│   │
│   ├── [v-else] Contenido
│   │   │
│   │   ├── Fila de profesionales (2 columnas)
│   │   │   ├── Card "PSICÓLOGAS"
│   │   │   │   └── ProfessionalRow × N  [v-for psychologists]
│   │   │   │       ├── Avatar (inicial)
│   │   │   │       ├── Nombre completo
│   │   │   │       ├── [si enDupla] Texto secundario "Dupla {name}"
│   │   │   │       └── Badge
│   │   │   │           ├── [si enDupla]  "En dupla"   (.badge--assigned)
│   │   │   │           └── [si !enDupla] "Disponible" (.badge--available)
│   │   │   │
│   │   │   └── Card "TRABAJADORAS SOCIALES"
│   │   │       └── ProfessionalRow × N  [v-for socialWorkers]
│   │   │           ├── Avatar (inicial)
│   │   │           ├── Nombre completo
│   │   │           ├── [si enDupla] Texto secundario "Dupla {name}"
│   │   │           └── Badge "En dupla" | "Disponible"
│   │   │
│   │   └── Card "DUPLAS ACTIVAS ({duplas.length})"
│   │       ├── [v-if duplas.length === 0] Empty state
│   │       │   └── Texto: "No hay duplas activas. Crea la primera con + Nueva dupla."
│   │       │
│   │       └── DuplaRow × N  [v-for duplas]
│   │           ├── Avatar "D"
│   │           ├── Nombre de la dupla  (ej. "Dupla 1")
│   │           ├── Miembros: "{psName} (Psicóloga) + {tsName} (Trabajadora Social)"
│   │           └── Acciones
│   │               ├── Botón "Editar"   → openDuplaModal('edit', dupla)   [E02]
│   │               └── Botón "Eliminar" → openDeleteConfirm(dupla)        [E05]
│   │
│   ├── Modal: Crear / Editar dupla  [v-if modal.visible && modal.kind === 'form']
│   │   ├── Título: [create] "Nueva dupla" | [edit] "Editar dupla"
│   │   ├── Subtítulo: "Cada profesional solo puede pertenecer a una dupla a la vez."
│   │   ├── Campo "Nombre de la dupla"  <input text>  → form.name
│   │   ├── Campo "Psicóloga"           <select>      → form.psychologistId
│   │   │   └── Opciones = availablePsychologists (+ miembro actual si edit)
│   │   ├── Campo "Trabajadora Social"  <select>      → form.socialWorkerId
│   │   │   └── Opciones = availableSocialWorkers (+ miembro actual si edit)
│   │   ├── [v-if warningPs]  Texto aviso naranja (sin psicólogas disponibles)
│   │   ├── [v-if warningTs]  Texto aviso naranja (sin trabajadoras sociales disponibles)
│   │   │   └── Ej.: "Todas las trabajadoras sociales están asignadas a una dupla."
│   │   ├── [v-if saveError]  Mensaje de error de guardado
│   │   └── Footer
│   │       ├── Botón "Cancelar"         → closeFormModal()     [E03]
│   │       └── Botón primario
│   │           ├── [create] "Crear dupla" / "Guardar"
│   │           └── [edit]   "Guardar cambios"
│   │               → saveDuplaModal()   [E04]
│   │               :disabled si isSaving || !canSave
│   │
│   ├── Modal: Confirmar eliminar  [v-if modal.visible && modal.kind === 'delete']
│   │   ├── Título: "Eliminar dupla"
│   │   ├── Texto: "¿Estás seguro de que deseas eliminar {dupla.name}?
│   │   │           Esta acción libera a sus miembros para ser asignados a otra dupla."
│   │   └── Footer
│   │       ├── Botón "Cancelar"  → closeDeleteConfirm()   [E06]
│   │       └── Botón "Eliminar" (destructivo) → confirmDeleteDupla()   [E07]
│   │           :disabled si isDeleting
│   │
│   └── Modal: Error — dupla en uso  [v-if modal.visible && modal.kind === 'deleteBlocked']
│       ├── Título: "No se puede eliminar"
│       ├── Texto: "Esta dupla no se puede eliminar porque se está utilizando en una sesión."
│       └── Footer
│           └── Botón "Entendido" / "Cerrar"  → closeDeleteBlockedModal()   [E08]
```

---

## Estado visual de badges (profesionales)

| Condición | Badge | Texto bajo el nombre |
|---|---|---|
| `duplaId != null` | “En dupla” | `Dupla {duplaName}` |
| `duplaId == null` | “Disponible” | — |

---

## Modo del modal de formulario

| Origen | `modal.mode` | Título | Acción primaria |
|---|---|---|---|
| Botón “+ Nueva dupla” | `create` | Nueva dupla | Crear / Guardar |
| Botón “Editar” en fila | `edit` | Editar dupla | Guardar cambios |
