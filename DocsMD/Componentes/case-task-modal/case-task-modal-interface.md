# `case-task-modal` — Interfaz del Componente

Modal reutilizable para completar tareas (`case_task`). El padre lo activa vía `open(taskId)`, el componente carga la tarea del API, renderiza el formulario específico según `case_task.type` (`gestion_llamada`, `proyectar_oficio`, `comite_caso`, `Corregir oficio`, `justificar_remision` o `solicitud_entidad`) y al confirmar la marca como completada ejecutando los efectos de lado correspondientes en el backend.

## Archivos relevantes

| Archivo | Rol |
|---|---|
| `src/frontend/js/components/case-task-modal.js` | Componente principal — template y lógica Vue |

---

## Árbol de interfaz

```
case-task-modal
│
├── [v-if !visible]
│   └── (no renderiza nada)
│
└── [v-if visible]  Overlay  (.ctm-overlay)
    └── ModalBox  (.ctm-modal)
        │
        ├── Header  (.ctm-header)
        │   ├── Title  (.ctm-title)
        │   │   SEGÚN tarea.type:
        │   │     'gestion_llamada'     → "Gestión de Llamada"
        │   │     'proyectar_oficio'    → "Proyectar Oficio"
        │   │     'comite_caso'         → "Decisiones del Comité"
        │   │     'Corregir oficio'     → "Corregir Oficio"
        │   │     'justificar_remision' → "Justificar Remisión"
        │   │     'solicitud_entidad'   → "Solicitudes a {entityName}"
        │   └── BtnCerrar  "✕"  → cancelar()
        │
        ├── [v-if cargandoTarea]
        │   └── Spinner  "Cargando tarea..."
        │
        ├── [v-if errorTarea && !cargandoTarea]
        │   └── ErrorMsg  (.ctm-error-load)  errorTarea
        │
        └── [v-if tarea && !cargandoTarea]
            │
            ├── Body  (.ctm-body)
            │   │
            │   ├── [tarea.type === 'gestion_llamada']
            │   │   FormLlamada  (.ctm-form)
            │   │   ├── <select> Departamento*  (.ctm-select)
            │   │   ├── <select> Ciudad*  (.ctm-select)  :disabled si !form.departamentoId
            │   │   ├── <select> Municipio*  (.ctm-select)  :disabled si !form.ciudadId
            │   │   ├── <select> Entidad*  (.ctm-select)  :disabled si !form.municipioId
            │   │   │   └── opción fija "Otra entidad"  value="otra"
            │   │   ├── [v-if form.entidadId === 'otra']
            │   │   │   <input> Nombre de la entidad*  (.ctm-input)
            │   │   ├── <input> Funcionario*  (.ctm-input)
            │   │   ├── <input> Descripción  (.ctm-input)
            │   │   ├── Switch  "¿Genera oficio?"  → form.generaOficio
            │   │   └── [v-if form.generaOficio]
            │   │       ├── <input> Asunto*  (.ctm-input)
            │   │       └── <input> Ruta al archivo en Kofax*  (.ctm-input)
            │   │
            │   ├── [tarea.type === 'proyectar_oficio']
            │   │   FormOficio  (.ctm-form)
            │   │   ├── <select> Departamento*  (.ctm-select)
            │   │   ├── <select> Ciudad*  (.ctm-select)  :disabled si !form.departamentoId
            │   │   ├── <select> Municipio*  (.ctm-select)  :disabled si !form.ciudadId
            │   │   ├── <select> Entidad*  (.ctm-select)  :disabled si !form.municipioId
            │   │   │   └── opción fija "Otra entidad"  value="otra"
            │   │   ├── [v-if form.entidadId === 'otra']
            │   │   │   <input> Nombre de la entidad*  (.ctm-input)
            │   │   ├── <input> Funcionario*  (.ctm-input)
            │   │   ├── <input> Asunto*  (.ctm-input)
            │   │   └── <input> Ruta al archivo en Kofax*  (.ctm-input)
            │   │
            │   └── [tarea.type === 'comite_caso']
            │       FormComite  (.ctm-form)
            │       ├── CheckboxGroup  "Decisiones del comité"*  [min 1 opción]
            │       │   opciones:
            │       │     activar_enlace          "Activar Enlace"
            │       │     oficio                  "Oficio"
            │       │     recomendaciones_agente  "Recomendaciones al agente"
            │       │     mecanismo_articulador   "Mecanismo articulador"
            │       │
            │       ├── [form.decisiones.includes('oficio')]
            │       │   <textarea> Observaciones del oficio  (.ctm-textarea)
            │       │
            │       ├── [form.decisiones.includes('recomendaciones_agente')]
            │       │   <textarea> Observaciones de recomendaciones  (.ctm-textarea)
            │       │
            │       └── [form.decisiones.includes('mecanismo_articulador')]
            │           ├── <select> Nivel*  (.ctm-select)
            │           │   opciones: Municipal | Departamental | Nacional
            │           └── <textarea> Observaciones del mecanismo  (.ctm-textarea)
            │
            │   [tarea.type === 'solicitud_entidad']
            │   FormSolicitudEntidad  (.ctm-form)
            │   ├── Intro  "Selecciona los ítems pendientes..."
            │   ├── [v-if cargandoObligaciones]
            │   │   └── Spinner  "Cargando obligaciones..."
            │   ├── [v-else-if errorObligaciones]
            │   │   └── ErrorMsg  errorObligaciones
            │   └── [v-else]
            │       ├── Label  "Checklist entidad pendientes"
            │       ├── [v-if !form.items.length]
            │       │   └── Empty  "No hay obligaciones pendientes"
            │       └── Item × N  [v-for item in form.items]
            │           ├── Chip  item.label  → toggleObligacion(item.id)
            │           │   selected → fondo púrpura + check
            │           └── [v-if seleccionado]
            │               └── <textarea> Describe la solicitud*  maxlength=500
            │
            │   [tarea.type === 'justificar_remision']
            │   FormJustificarRemision  (.ctm-form)
            │   ├── MotivoDevolucion  (fondo rojo claro)
            │   │   ├── Label  "⚠ Motivo de la devolución"
            │   │   └── Texto  tarea.description   // "Justificar remisión devuelta — Motivo: {motivo}"
            │   └── <textarea> Tu respuesta*  (.ctm-textarea)  maxlength=500
            │
            │   [tarea.type === 'Corregir oficio']
            │   FormCorregirOficio  (.ctm-form)  // solo lectura, sin campos editables
            │   │
            │   ├── [v-if cargandoOficio]
            │   │   └── Spinner  "Cargando datos del oficio..."
            │   │
            │   ├── [v-else-if oficioVinculado]
            │   │   ├── [v-if oficioVinculado.reasonCorrection]
            │   │   │   RazonCorreccion  (.ctm-razon-box)
            │   │   │   ├── Label  "⚠ Razón de corrección"
            │   │   │   └── Texto  oficioVinculado.reasonCorrection
            │   │   ├── [v-if oficioVinculado.urlKofax]
            │   │   │   KofaxBox  (.ctm-kofax-box)
            │   │   │   ├── Label  "Oficio en Kofax"
            │   │   │   └── Texto  oficioVinculado.urlKofax
            │   │   └── Instrucciones  "Revisa el oficio en el Kofax, aplica las correcciones indicadas y presiona 'Marcar como corregido'..."
            │   │
            │   └── [v-else]
            │       └── ErrorMsg  "No se pudo cargar la información del oficio."
            │
            ├── [v-if saveError]
            │   ErrorMsg  (.ctm-save-error)  saveError
            │
            └── Footer  (.ctm-footer)
                ├── BtnCancelar  "Cancelar"  → cancelar()
                └── BtnConfirmar
                    │   SEGÚN tarea.type:
                    │     'Corregir oficio'     → "Marcar como corregido"
                    │     'solicitud_entidad'   → "Enviar solicitudes"
                    │     default               → "Completar tarea"
                    :disabled si !formularioValido || guardando || cargandoTarea
                    → confirmar()
```

---

## Efectos de lado por tipo de tarea (ejecutados por el backend)

| Tipo | Condición | Efecto en el backend |
|---|---|---|
| `gestion_llamada` | `generaOficio = false` | Solo guarda JSON en `case_task.form_data` + marca `status = 'Done'` |
| `gestion_llamada` | `generaOficio = true` | Crea `entity_letter` en estado `para_revisar` con datos del form + guarda JSON + Done |
| `proyectar_oficio` | — | Actualiza la `entity_letter` vinculada (`case_task.entity_letter_id`) con los datos del form → pasa a `para_revisar` + guarda JSON + Done |
| `comite_caso` | `decisiones.includes('activar_enlace')` | Activa flag booleano de enlace territorial en el registro + guarda JSON + Done |
| `comite_caso` | `decisiones.includes('oficio')` | Crea `entity_letter` (estado `por_proyectar`) + nueva `case_task` de tipo `proyectar_oficio` asignada al agente del caso + guarda JSON + Done |
| `Corregir oficio` | — | Transiciona el `entity_letter` vinculado (`case_task.entity_letter_id`) de `en_correccion` → `para_revisar` + marca la tarea `Done`. No hay `form_data` que guardar — es solo confirmación. |
| `justificar_remision` | — | Guarda JSON (`respuesta`) + Done. Reabre la `psychosocial_support` vinculada (`status → 'abierto'`) y crea una nueva `case_task` `validar_remision` (`ToDo`) para que el psicólogo/a vuelva a evaluarla — asignada al `psychologist_id` de la `dupla` si la remisión tiene una, si no al `professional_id` individual, si no queda sin asignar (igual que la remisión). Crea `case_timeline_event` "Remisión Justificada". Ver `DocsMD/Screens/Remision Psicosocial Detalle/QA/validar-remision-justificar-qa-cases.md`. |
| `solicitud_entidad` | — | Valida ≥1 ítem con descripción. Resuelve exactamente un `general_user` con `entity_branch_id` = sede del `entity_case`. Si hay 0 o N usuarios, **no** marca Done y responde 400. Si hay uno: transacción — marca esta tarea Done + `form_data` y crea N `case_task` `respuesta_entidad` (`ToDo`) asignadas a ese usuario, una por ítem. No cambia el status de `entity_case_obligation`. |

---

## Validaciones por tipo de tarea

| Campo | `gestion_llamada` | `proyectar_oficio` | `comite_caso` | `Corregir oficio` | `justificar_remision` | `solicitud_entidad` |
|---|---|---|---|---|---|---|
| Departamento | required | required | — | — | — | — |
| Ciudad | required | required | — | — | — | — |
| Municipio | required | required | — | — | — | — |
| Entidad | required | required | — | — | — | — |
| Nombre entidad (si "otra") | required | required | — | — | — | — |
| Funcionario | required | required | — | — | — | — |
| Descripción | opcional | — | — | — | — | — |
| Asunto | required si `generaOficio` | required | — | — | — | — |
| Ruta al archivo en Kofax | required si `generaOficio` | required | — | — | — | — |
| Decisiones del comité | — | — | min 1 requerida | — | — | — |
| Nivel (mecanismo articulador) | — | — | required si decisión activa | — | — | — |
| Tu respuesta | — | — | — | — | required | — |
| Ítems del checklist | — | — | — | — | — | min 1 seleccionado |
| Descripción por ítem | — | — | — | — | — | required en cada seleccionado |

> `Corregir oficio` no tiene campos editables. `formularioValido` solo exige que `cargandoOficio === false` (que haya terminado de cargar el `entity_letter` vinculado en modo lectura).

---

## Límite de caracteres por campo

Todos los `<input>`/`<textarea>` de texto libre tienen `maxlength` (atributo HTML nativo — el navegador bloquea escribir/pegar más allá del límite). Los `<select>` no aplican.

| Campo | `maxlength` |
|---|---|
| Funcionario contactado / Dependencia del funcionario | 50 |
| Asunto (del oficio) | 80 |
| Nombre de la entidad (si "otra") | 150 |
| Ruta Kofax / URL del archivo | 255 |
| Descripción / Notas (`gestion_llamada`) | 500 |
| Observaciones del oficio / para el agente / del mecanismo (`comite_caso`) | 500 |
| Tu respuesta (`justificar_remision`) | 500 |
| Descripción por ítem (`solicitud_entidad`) | 500 |
