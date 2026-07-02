# `case-task-modal` — Interfaz del Componente

Modal reutilizable para completar tareas (`case_task`). El padre lo activa vía `open(taskId)`, el componente carga la tarea del API, renderiza el formulario específico según `case_task.type` (`gestion_llamada`, `proyectar_oficio`, `comite_caso` o `Corregir oficio`) y al confirmar la marca como completada ejecutando los efectos de lado correspondientes en el backend.

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
        │   │     'gestion_llamada'  → "Gestión de Llamada"
        │   │     'proyectar_oficio' → "Proyectar Oficio"
        │   │     'comite_caso'      → "Decisiones del Comité"
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
                    │     'Corregir oficio'  → "Marcar como corregido"
                    │     default             → "Completar tarea"
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

---

## Validaciones por tipo de tarea

| Campo | `gestion_llamada` | `proyectar_oficio` | `comite_caso` | `Corregir oficio` |
|---|---|---|---|---|
| Departamento | required | required | — | — |
| Ciudad | required | required | — | — |
| Municipio | required | required | — | — |
| Entidad | required | required | — | — |
| Nombre entidad (si "otra") | required | required | — | — |
| Funcionario | required | required | — | — |
| Descripción | opcional | — | — | — |
| Asunto | required si `generaOficio` | required | — | — |
| Ruta al archivo en Kofax | required si `generaOficio` | required | — | — |
| Decisiones del comité | — | — | min 1 requerida | — |
| Nivel (mecanismo articulador) | — | — | required si decisión activa | — |

> `Corregir oficio` no tiene campos editables. `formularioValido` solo exige que `cargandoOficio === false` (que haya terminado de cargar el `entity_letter` vinculado en modo lectura).
