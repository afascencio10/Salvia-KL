# `case-oficios` — Interfaz del Componente

> ⚠️ Componente **no implementado aún**. Este MD se escribió antes del código, siguiendo el proceso de set up (`dev-docs-setup.md` → Proceso 2).

Lista los oficios (`entity_letter`) de un caso en formato de cards, con filtros (buscador de texto, chips por tema, filtro por estado) y un modal de detalle completo al hacer click en una card.

Es un componente **nuevo y separado** de `oficios-list.js` (el componente existente que hoy se usa en la pantalla de Notificaciones, en formato tabla). No lo reemplaza ni lo modifica — cubre un caso de uso distinto (vista de cards enfocada en un caso, con filtros que `oficios-list.js` no tiene).

## Archivos relevantes

| Archivo | Rol |
|---|---|
| `src/frontend/js/components/case-oficios.js` | Componente principal — template y lógica Vue (a implementar) |

---

## Supuestos de diseño (a confirmar en revisión)

- Recibe `caseId` como prop (String, requerido) y hace un único fetch al montar: `GET /api/v1/entity-letters?caseId={caseId}`.
- **Todos los filtros son client-side** (computed sobre la lista ya cargada) — el endpoint real (`EntityLetterController.List`) resuelve sus query params en un `switch` mutuamente excluyente: si se manda `state`, ignora `caseId`. No se puede pedir "oficios de este caso Y en este estado" en una sola llamada al backend, así que no tiene sentido intentarlo — se filtra en el cliente sobre los datos ya traídos.
- No pagina — asume que los oficios de un caso son un volumen manejable (igual que `case-task-history`).
- Buscador de texto: coincide contra `entidad`, `url_kofax`, y los campos de asunto disponibles (`subject`, `asunto_radicado`, `asunto_respuesta` — el que esté presente según la etapa del oficio).
- Filtro por estado: `<select>` o chips con los 7 valores de `entity_letter.state` + opción "Todos".
- ⚠️ **GAP de schema — chips por tema:** el mapa de requerimientos pide filtrar por 4 "temas": `barrier_id`, `emergency_measure_id`, `psychosocial_support_id`, `economic_stabilization_id`. Hoy `entity_letter` **solo tiene `barrier_id`** (obligatorio, `not null`). Las otras 3 columnas no existen en el modelo ni en la tabla — son propias de `case_task`, no de `entity_letter`. Para que el chip de tema funcione en los 4 casos, hay que:
  1. Agregar `emergency_measure_id`, `psychosocial_support_id`, `economic_stabilization_id` como columnas nullable a `salvia.entity_letter` (y al modelo Go `EntityLetter`).
  2. Decidir si un oficio puede pertenecer a más de un "tema" a la vez o es siempre exactamente uno (afecta si el chip filter es single-select o multi-select).
  3. Actualizar el flujo de creación de `entity_letter` para poblar el campo correspondiente según el origen (hoy solo `sideEffectsComiteCaso`/`buildEntityLetterFromFormData`/etc. en `case_task_service.go` setean `BarrierID`).
  **Hasta que esto se resuelva, el chip de tema solo puede filtrar por "Barreras" — los otros 3 chips quedarían deshabilitados o no se muestran.**

---

## Árbol de interfaz

```
case-oficios
│
├── [v-if cargando]
│   └── Spinner  "Cargando oficios..."
│
├── [v-else-if error]
│   └── ErrorMsg  (.co-error)  error  +  BtnReintentar  → cargar()
│
└── [v-else]
    │
    ├── Filtros  (.co-filtros)
    │   ├── <input> Buscador  "Buscar por entidad, ruta o asunto..."  (.co-buscador)
    │   │
    │   ├── ChipsTema  (.co-chips-tema)
    │   │   ├── Chip "Barreras"                    :active="temaActivo === 'barrera'"
    │   │   ├── Chip "Medidas de Emergencia"        :disabled  ⚠️ GAP — sin columna en entity_letter
    │   │   ├── Chip "Apoyo Psicosocial"            :disabled  ⚠️ GAP — sin columna en entity_letter
    │   │   └── Chip "Estabilización Económica"     :disabled  ⚠️ GAP — sin columna en entity_letter
    │   │
    │   └── <select> Estado  (.co-select-estado)
    │       opciones: Todos | Por proyectar | Para revisar | En corrección |
    │                 Aprobación jurídica | Para radicar | Radicado | Respondido
    │
    ├── [v-if oficiosFiltrados.length === 0]
    │   └── EmptyState  (.co-empty)  "No hay oficios que coincidan con los filtros."
    │
    └── [v-else]
        └── Grid  (.co-grid)
            └── OficioCard × N  [v-for oficio in oficiosFiltrados]  (.co-card)  @click="abrirDetalle(oficio)"
                ├── Header
                │   ├── EntidadNombre  (.co-card-entidad)   oficio.entidad
                │   └── BadgeEstado  :class="stateClass(oficio.state)"   statusLabel(oficio.state)
                ├── Badges  (.co-card-badges)
                │   ├── [oficio.nivel] BadgeNivel   nivelLabel(oficio.nivel)
                │   └── BadgePrioridad   priorityLabel(oficio.priority)
                ├── [oficio.subjectResuelto] Asunto  (.co-card-asunto)  truncado
                ├── [oficio.url_kofax] DocLink  (.co-card-doc)  <i fa-file-pdf/>  truncateUrl(oficio.url_kofax)
                └── Fecha  (.co-card-fecha)  formatDate(oficio.created_at)

─── MODAL ─────────────────────────────────────────────────────────────────

[v-if oficioSeleccionado]  Overlay  (.co-modal-overlay)  @click.self="cerrarDetalle"
  └── ModalBox  (.co-modal)
      ├── Header
      │   ├── Titulo "Detalle del Oficio"
      │   ├── Subtitulo  truncateUrl(oficio.url_kofax || oficio.id, 55)
      │   └── BtnCerrar "×"  → cerrarDetalle()
      │
      └── Body  (.co-modal-body)
          │  — mismo contenido que el modal ya existente en oficios-list.js —
          ├── Sección "Oficio"          (id, estado, prioridad, entidad+nivel, correo entidad,
          │                              documento/url_kofax, creado, actualizado, registrado por)
          ├── [si tiene datos] Sección "Radicación y Respuesta"
          │                              (número radicado, asunto radicado, correo remitente,
          │                               asunto respuesta, fecha respuesta, revisado por)
          ├── [si reason_correction] Sección "Corrección requerida"  (fondo warning)
          ├── [si barrier_id] Sección "Barrera relacionada"  → Link "Ver barrera →"
          ├── [si case_id] Sección "Caso relacionado"  → Link "Ver caso →"
          └── Acciones
              ├── BtnDescargarPDF  → downloadPDF(oficio)
              ├── BtnImprimir      → printOficio()
              └── BtnCerrar        → cerrarDetalle()
```

---

## Filtros — comportamiento

| Filtro | Tipo | Campo(s) que evalúa | Estado |
|---|---|---|---|
| Buscador | texto libre, client-side | `entidad`, `url_kofax`, `subject` \| `asunto_radicado` \| `asunto_respuesta` | ✅ Implementable ya |
| Chip "Barreras" | toggle, client-side | `barrier_id != null` | ✅ Implementable ya |
| Chip "Medidas de Emergencia" | toggle, client-side | `emergency_measure_id != null` | ⚠️ GAP — columna no existe |
| Chip "Apoyo Psicosocial" | toggle, client-side | `psychosocial_support_id != null` | ⚠️ GAP — columna no existe |
| Chip "Estabilización Económica" | toggle, client-side | `economic_stabilization_id != null` | ⚠️ GAP — columna no existe |
| Select Estado | single-select, client-side | `state === valor` | ✅ Implementable ya |

> Los 3 filtros marcados como GAP dependen de la migración de schema descrita arriba en "Supuestos de diseño".
