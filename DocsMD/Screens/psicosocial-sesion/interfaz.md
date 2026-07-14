# `/salvia/psicosocial/registrar/:psicosocial_id` — Interfaz de la Pantalla

Pantalla para que el agente psicosocial registre una sesión. Muestra la tarjeta de información del caso y un formulario dinámico seleccionado automáticamente según el estado del proceso psicosocial (`psychosocial_support`).

**Archivo fuente:** `src/frontend/html/salvia/psicosocial/registrar_sesion.html`  
**Componente central:** `src/frontend/js/components/dinamic-form.js` (compartido con `hacer-seguimiento`)  
**Entry point:** `GET /salvia/psicosocial/registrar/:psicosocial_id` → `PsicosocialSesionFacade.go`

---

## Árbol de interfaz

```
Page  (.container.mt-4)
│
├── WelcomeSection  (.welcome-section.mb-4)
│   └── <h2>  "Registrar Sesión Psicosocial"
│
└── #psicosocial-app  (Vue mount)
    │
    ├── [v-if victimInfo]
    │   PsicosocialCard  (.ps-card)
    │   └── .card-body
    │       ├── Header  (.ps-card-header)
    │       │   ├── Avatar  (.ps-card-avatar)
    │       │   │   └── <i fa-user>
    │       │   └── NameBlock  (.ps-card-name-block)
    │       │       ├── <h5.ps-card-name>  victimInfo.Names + LastNames
    │       │       └── Location  (.ps-card-location)
    │       │           └── <i fa-map-marker-alt>  victimInfo.TownName
    │       ├── <hr.ps-divider>
    │       ├── Fields  (.ps-card-fields)
    │       │   ├── Field  (.ps-field)  "Teléfono"
    │       │   ├── Field  (.ps-field)  "Identidad de género"
    │       │   ├── Field  (.ps-field)  "Nivel de riesgo"
    │       │   └── Field  (.ps-field)  "Código de caso"
    │       └── SessionStatus  (.ps-session-status)
    │           ├── Badge  "Sesión X de 6"  (session_count + 1 de 6)
    │           ├── Badge  status label  (abierto | en gestión | en devolución | cerrado)
    │           └── FormTypeBadge  "Formulario: {nombre del form cargado}"
    │               (Primer Contacto | Primera Atención | Seguimiento | Cierre)
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
    │           │       └── [type === 'question' && isVisible]
    │           │           DirectQuestion  (.df-question)
    │           │           ├── <label>  description
    │           │           │   └── [v-if required]  <span.df-required>  "*"
    │           │           ├── [single]    → chips de selección única
    │           │           ├── [boolean]   → botones Sí / No
    │           │           ├── [multiple]  → chips de selección múltiple
    │           │           ├── [text]      → <textarea>
    │           │           ├── [date]      → <input type="date">
    │           │           └── [v-if questionErrors[id]]  <span.df-error>
    │           │
    │           └── NavigationBar  (.df-nav)
    │               ├── PrevBtn  (.df-btn-prev)  "← Anterior"
    │               │   :disabled si isFirstSection
    │               │   → goPrev()
    │               └── [v-if !isLastSection || canEdit]
    │                   NextBtn  (.df-btn-next)
    │                   ├── [isLastSection && canEdit]  "Guardar sesión ✓"  → saveSection()
    │                   └── [v-else]  "Siguiente →"  → goNext() | saveSection()
    │
    │       ValidationPopup  [v-if showValidationError]  (.df-popup-backdrop)
    │       └── .df-popup
    │           ├── Icon  (.df-popup-icon)  "⚠️"
    │           ├── <p.df-popup-title>  "Algunas respuestas no son válidas, revisa el formulario"
    │           └── CloseBtn  (.df-popup-btn)  "Entendido"
    │               → showValidationError = false
    │
    └── [v-if sessionCompleted]
        CompletedOverlay  (.ps-completed-backdrop)
        └── .ps-completed-card
            ├── Icon  (.ps-completed-icon)  SVG checkmark verde
            ├── <h4.ps-completed-title>  "Sesión registrada"
            ├── <p.ps-completed-text>
            │   [si continuar_primera_atencion = true]
            │     "La sesión de Primer Contacto fue guardada. Continúa con la Primera Atención."
            │     → ContinuarBtn  "Iniciar Primera Atención →"  (recarga la pantalla)
            │   [si solo_contacto = true]
            │     "El contacto fue registrado. La próxima sesión queda agendada para {fecha_nueva}."
            │   [cualquier otro caso]
            │     "La sesión fue registrada exitosamente."
            └── VerRemisionBtn  (.ps-completed-btn)  "Ver remisión"
                → goToRemision()  [redirige a /salvia/remisiones/:psicosocial_id]
```

---

## Tipos de input soportados (igual que DinamicForm en hacer-seguimiento)

| `questionTypeId` | Elemento renderizado | Clase CSS | Modo read-only |
|---|---|---|---|
| `single` | Grupo de `<button>` chip-selección única | `.df-btn-group` / `.df-opt-btn` | `:disabled` si `!canEdit` |
| `boolean` | Dos `<button>` "Sí" / "No" | `.df-bool-group` / `.df-bool-btn` | `:disabled` si `!canEdit` |
| `multiple` | Grupo de `<button>` chip-selección múltiple | `.df-btn-group` / `.df-opt-btn` | `:disabled` si `!canEdit` |
| `text` | `<textarea>` | `.df-textarea` | `:disabled` si `!canEdit` |
| `date` | `<input type="date">` | `.df-input` | `:disabled` si `!canEdit` |

---

## Estados visuales del SectionItem (sidebar)

| Estado | Condición | Clase CSS | Badge |
|---|---|---|---|
| `active` | Es la sección actualmente visible en MainPanel | `.df-section-item.active` | `.df-badge.active` — número (violeta) |
| `completed` | `section.isAnswered === true` | `.df-section-item.completed` | `.df-badge.completed` — "✓" (verde) |
| `enabled` | Primera sección sin responder (no activa) | `.df-section-item.enabled` | `.df-badge.enabled` — número (gris) |
| `pending` | Secciones posteriores no habilitadas | `.df-section-item.pending` | `.df-badge.pending` — número (gris claro) |

`:disabled` en pending — no es clickeable.

---

## Estados del modo canEdit

`canEdit` se calcula en `loadPsicosocial()` del app raíz:

| Condición | canEdit | Efecto visual |
|---|---|---|
| `psychosocial_support.status === 'cerrado'` | `false` | Badge "🔒 Solo lectura", inputs `:disabled` |
| Cualquier otro status | `true` | Formulario editable, botón "Guardar sesión ✓" |

---

## Overlay de sesión completada — variantes

| Escenario | Texto principal | Botón extra |
|---|---|---|
| Primer Contacto con Continuar = Sí | "Primer Contacto guardado. Continúa con la Primera Atención." | "Iniciar Primera Atención →" (recarga pantalla) |
| Solo Contacto registrado | "El contacto fue registrado. Próxima sesión: {fecha_nueva}." | — |
| Sesión normal completada | "La sesión fue registrada exitosamente." | — |
| Cierre completado | "El proceso psicosocial fue cerrado exitosamente." | — |
