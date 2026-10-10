# `/salvia/hacer-seguimiento/:id` — Interfaz de la Pantalla

Pantalla para ejecutar un seguimiento: muestra la tarjeta de la víctima y un formulario dinámico por secciones que el profesional completa y guarda progresivamente.
Archivo fuente: `src/frontend/html/salvia/follow_up_v2/hacer_seguimiento.html`
Componente central: `src/frontend/js/components/dinamic-form.js`
Entry point: `GET /salvia/hacer-seguimiento/:id` → `HacerSeguimientoFacade.go`

---

## Árbol de interfaz

```
Page  (.container.mt-4)
│
├── WelcomeSection  (.welcome-section.mb-4)
│   └── <h2>  "Hacer Seguimiento"
│
└── #seguimiento-app  (Vue mount)
    │
    ├── [v-if !canEdit]
    │   ReadonlyBanner  (.ms-readonly-banner)  🔒 "Este seguimiento no se puede editar" + noEditReason
    │
    ├── [v-if victimInfo]
    │   VictimCard  (.hs-victim-card)
    │   └── .card-body
    │       ├── Header  (.hs-victim-header)
    │       │   ├── Avatar  (.hs-victim-avatar)
    │       │   │   └── <i fa-user>
    │       │   └── NameBlock  (.hs-victim-name-block)
    │       │       ├── <h5.hs-victim-name>  victimInfo.Names + LastNames
    │       │       └── Location  (.hs-victim-location)
    │       │           └── <i fa-map-marker-alt>  victimInfo.TownName
    │       ├── <hr.hs-divider>
    │       └── Fields  (.hs-victim-fields)
    │           ├── Field  (.hs-field)  "Teléfono"
    │           ├── Field  (.hs-field)  "Identidad de género"
    │           ├── Field  (.hs-field)  "Orientación sexual"
    │           ├── Field  (.hs-field)  "Teléfono de contacto"
    │           └── Field  (.hs-field)  "Edad"
    │
    ├── [v-if loadError]
    │   ErrorAlert  (.alert.alert-danger)  loadError
    │
    ├── [v-else-if formId && submissionId]
    │   DinamicForm  (componente Vue)
    │   │
    │   ├── [v-if loading]
    │   │   LoadingSpinner  — SVG spin + "Cargando formulario..."
    │   │
    │   ├── [v-else-if error]
    │   │   ErrorPanel  — error message
    │   │
    │   └── [v-else]  (formulario cargado)
    │       .df-wrapper
    │       │
    │       ├── Sidebar  (.df-sidebar)
    │       │   ├── Header  (.df-sidebar-header)
    │       │   │   └── <p.df-sidebar-title>  "SECCIONES"
    │       │   ├── Nav  (.df-sidebar-nav)
    │       │   │   └── SectionItem × N  [v-for visibleSections]  (.df-section-item)
    │       │   │       │   State: active | completed | enabled | pending
    │       │   │       │   :disabled si !isSectionEnabled(section)
    │       │   │       │   @click → goToSectionById(section)
    │       │   │       ├── Badge  (.df-badge)
    │       │   │       │   ├── [v-if section.isAnswered]  "✓"
    │       │   │       │   └── [v-else]  section.order
    │       │   │       └── section.name
    │       │   └── ProgressWrap  (.df-progress-wrap)
    │       │       ├── Label  (.df-progress-label)  "Progreso" + "N / total"
    │       │       └── BarBg  (.df-progress-bar-bg)
    │       │           └── BarFill  (.df-progress-bar-fill)  :width=progressPercent%
    │       │
    │       └── MainPanel  (.df-main)
    │           │
    │           ├── [v-if saving]
    │           │   SavingOverlay  (.df-saving-overlay)  — SVG spin + "Guardando..."
    │           │
    │           ├── SectionHeader  (.df-section-header)
    │           │   ├── Row  (.df-section-header-row)
    │           │   │   ├── Number  (.df-section-number)  currentSection.order
    │           │   │   ├── <h2.df-section-title>  currentSection.name
    │           │   │   └── [v-if !canEdit]
    │           │   │       ReadonlyBadge  (.df-readonly-badge)  "🔒 Solo lectura"
    │           │   └── [v-if currentSection.description]
    │           │       <p.df-section-desc>  currentSection.description
    │           │
    │           ├── QuestionsArea  (.df-questions)  [v-if currentSectionRender]
    │           │   └── FormItem × N  [v-for currentSectionRender.formItems]
    │           │       │
    │           │       ├── [type === 'question' && isVisible]
    │           │       │   DirectQuestion  (.df-question)
    │           │       │   ├── <label>  description
    │           │       │   │   └── [v-if required]  <span.df-required>  "*"
    │           │       │   ├── [questionTypeId === 'single']  → ver tabla de tipos
    │           │       │   ├── [questionTypeId === 'dropdown']
    │           │       │   ├── [questionTypeId === 'boolean']
    │           │       │   ├── [questionTypeId === 'multiple']
    │           │       │   ├── [questionTypeId === 'text']
    │           │       │   ├── [questionTypeId === 'number']
    │           │       │   ├── [questionTypeId === 'date']
    │           │       │   ├── [questionTypeId === 'datetime']
    │           │       │   └── [v-if questionErrors[id]]  <span.df-error>
    │           │       │
    │           │       └── [type === 'repeater' && isVisible]
    │           │           RepeaterGroup  (.df-question)
    │           │           ├── <label>  repeater.name
    │           │           │   └── [v-if minRepetitions > 0]  <span.df-required>  "*"
    │           │           ├── .df-repeater
    │           │           │   ├── RepeaterEntry × N  [v-for entries]  (.df-repeater-item)
    │           │           │   │   ├── Header  (.df-repeater-item-header)
    │           │           │   │   │   ├── Title  (.df-repeater-item-title)  "itemName #iteration"
    │           │           │   │   │   └── [v-if canEdit]
    │           │           │   │   │       DeleteBtn  (.df-repeater-delete)  "✕ Eliminar"
    │           │           │   │   │       → removeRepeaterEntry()
    │           │           │   │   └── RepeaterQuestion × N  [v-for questions, v-if isVisible]
    │           │           │   │       (.df-question)
    │           │           │   │       ├── <label>  description
    │           │           │   │       ├── [por questionTypeId]  → mismos tipos que DirectQuestion
    │           │           │   │       └── [v-if questionErrors[key]]  <span.df-error>
    │           │           │   └── [v-if canEdit]
    │           │           │       AddBtn  (.df-repeater-add)  "+ Agregar {itemName}"
    │           │           │       → addRepeaterEntry()
    │           │           └── [v-if repeaterErrors[id]]  <span.df-error>
    │           │
    │           └── NavigationBar  (.df-nav)
    │               ├── PrevBtn  (.df-btn-prev)  "← Anterior"
    │               │   :disabled si isFirstSection
    │               │   → goPrev()
    │               └── [v-if !isLastSection || canEdit]
    │                   NextBtn  (.df-btn-next)
    │                   ├── [isLastSection && canEdit]  "Guardar seguimiento ✓"  → saveSection()
    │                   └── [v-else]  "Siguiente →"  → goNext() | saveSection()
    │
    │       ValidationPopup  [v-if showValidationError]  (.df-popup-backdrop)
    │       └── .df-popup
    │           ├── Icon  (.df-popup-icon)  "⚠️"
    │           ├── <p.df-popup-title>  "Algunas respuestas no son válidas, revisa el formulario"
    │           └── CloseBtn  (.df-popup-btn)  "Entendido"
    │               → showValidationError = false
    │
    └── [v-if formCompleted]
        CompletedOverlay  (.hs-completed-backdrop)
        └── .hs-completed-card
            ├── Icon  (.hs-completed-icon)  SVG checkmark verde
            ├── <h4.hs-completed-title>  "Seguimiento registrado"
            ├── <p.hs-completed-text>  "El seguimiento fue completado exitosamente."
            └── VerCasoBtn  (.hs-completed-btn)  "Ver caso"
                → goToCase()  [redirige a /salvia/casos/:caseId/detalle]
```

---

## Tipos de input soportados en DirectQuestion y RepeaterQuestion

| `questionTypeId` | Elemento renderizado | Clase CSS | Modo read-only |
|---|---|---|---|
| `single` | Grupo de `<button>` chip-selección única | `.df-btn-group` / `.df-opt-btn` | `:disabled` si `!canEdit` |
| `dropdown` | `<select>` con `<option> × N` | `.df-select` | `:disabled` si `!canEdit` |
| `boolean` | Dos `<button>` "Sí" / "No" | `.df-bool-group` / `.df-bool-btn` | `:disabled` si `!canEdit` |
| `multiple` | Grupo de `<button>` chip-selección múltiple | `.df-btn-group` / `.df-opt-btn` | `:disabled` si `!canEdit` |
| `text` | `<textarea>` | `.df-textarea` | `:disabled` si `!canEdit` |
| `number` | `<input type="number">` | `.df-input` | `:disabled` si `!canEdit` |
| `date` | `<input type="date">` | `.df-input` | `:disabled` si `!canEdit` |
| `datetime` | `<input type="datetime-local">` | `.df-input` | `:disabled` si `!canEdit` |

---

## Estados visuales del SectionItem (sidebar)

| Estado | Condición | Clase CSS | Badge |
|---|---|---|---|
| `active` | Es la sección actualmente visible en MainPanel | `.df-section-item.active` | `.df-badge.active` — número (violeta) |
| `completed` | `section.isAnswered === true` | `.df-section-item.completed` | `.df-badge.completed` — "✓" (verde) |
| `enabled` | Es la primera sección sin responder (y no está activa) | `.df-section-item.enabled` | `.df-badge.enabled` — número (gris) |
| `pending` | Secciones posteriores aún no habilitadas | `.df-section-item.pending` | `.df-badge.pending` — número (gris claro) |

`:disabled` en pending — no es clickeable.

---

## Estados del modo canEdit

`canEdit` se calcula en `loadFollowUp()` del app raíz (no del componente), en cascada por prioridad (la primera condición que aplica gana). Cuando queda en `false`, `noEditReason` guarda el motivo exacto que muestra el `ReadonlyBanner` al inicio de la pantalla.

| # | Condición | canEdit | noEditReason |
|---|---|---|---|
| 1 | `formId !== FORM_ID_DEFAULT` (seguimiento migrado — Kobo, Bases Salvia, u otra fuente) | `false` permanente | "Fue migrado desde un formulario anterior (Kobo o Bases Salvia) y quedó fijo como solo lectura." |
| 2 | `caseStatus === 'cd'` (caso cerrado) | `false` permanente | "El caso asociado está cerrado." |
| 3 | `status === 'REALIZADO'` y `completed_at` hace > 5 días | `false` | "Se completó hace más de 5 días." |
| 3b | `status === 'REALIZADO'` y `completed_at` hace ≤ 5 días | `true` | — |
| 4 | Cualquier otro estado | `true` | — |
| 5 | **Override final** — `userRole === 'ps'` o `'ts'`: se evalúa después de las anteriores y pisa cualquier resultado, incluso si alguna de 1-4 dio `true` | `false` | "Tu rol solo tiene acceso de lectura a los seguimientos." |

Efecto visual cuando `canEdit = false`: `ReadonlyBanner` (🔒 + `noEditReason`) al inicio de la pantalla, badge "🔒 Solo lectura" en el header de la sección activa, todos los inputs `:disabled`, y el botón final de navegación solo avanza sin guardar.
