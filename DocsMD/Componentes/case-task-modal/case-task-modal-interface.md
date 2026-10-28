# `case-task-modal` — Interfaz del Componente

Modal reutilizable para completar tareas (`case_task`). El padre lo activa vía `open(taskId)`, `openArticulacion(ctx)` (Mis Barreras) o `openMecanismo(ctx)` (Casos Mecanismo). El componente carga la tarea del API (si hay id), renderiza el formulario según `case_task.type` (`gestion_llamada`, `proyectar_oficio`, `comite_caso`, `Corregir oficio`, `justificar_remision`, `solicitud_entidad`, `respuesta_entidad`, `articulacion_territorial`, `gestion_mecanismo`, `decidir_via_riesgo_inminente`, `intentar_contacto_3x3` o `escalar_cruz_roja`) y al confirmar ejecuta el complete correspondiente. `intentar_contacto_3x3` y `escalar_cruz_roja` cierran al pulsar Guardar o Confirmar y no llaman al API.

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
        │   │     'respuesta_entidad'   → "Respuesta institucional"
        │   │     'articulacion_territorial' → "Gestionar barrera"
        │   │     'gestion_mecanismo'        → según scope:
        │   │         municipal     → "Gestión del Comité Municipal"
        │   │         departamental → "Gestión del Comité Departamental"
        │   │         nacional      → "Gestión de la Instancia Nacional"
        │   │     'decidir_via_riesgo_inminente' → "Decidir vía de atención — Riesgo inminente"
        │   │         subtítulo: "Elige una vía y asigna un agente de medidas."
        │   │     'intentar_contacto_3x3' → "Registrar intento de contacto"
        │   │         subtítulo: "Protocolo 3×3 — registra el resultado del intento."
        │   │     'escalar_cruz_roja' → "Escalar a Cruz Roja"
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
            │   [tarea.type === 'respuesta_entidad']
            │   FormRespuestaEntidad  (.ctm-form)
            │   ├── ContextBox  (solo lectura)
            │   │   ├── Campo  "Obligación"  tarea.formData.label
            │   │   └── Campo  "Solicitud"   tarea.formData.descripcionSolicitud
            │   ├── <input> Nº de radicado / NUNC / MP  maxlength=80  (opcional)
            │   ├── <textarea> Ampliar información  maxlength=500  (opcional)
            │   └── <select> Estado institucional*  iniciado | en_proceso | concluido
            │
            │   [tarea.type === 'articulacion_territorial']
            │   FormArticulacion  (.ctm-form)
            │   ├── ContextBox  entidad + descripción de barrera
            │   ├── <input type=date> Fecha de contacto*
            │   ├── <textarea> Compromiso adquirido*
            │   ├── Toggle  ¿El compromiso requiere seguimiento?*  Sí/No
            │   ├── [v-if requiereSeguimiento] <input type=date> Fecha de seguimiento* + nota
            │   └── <textarea> Sugerencias de remisión interna o externa  (opcional)
            │
            │   [tarea.type === 'gestion_mecanismo']
            │   FormMecanismo  (.ctm-form)  // campos según mecanismoCtx.scope
            │   ├── [scope === 'municipal']
            │   │   ├── <select> ¿Abordado en sesión del Comité Municipal?*
            │   │   │     si_ordinaria | si_extraordinaria | no
            │   │   ├── [Si] <input type=date> Fecha de la sesión  (opcional)
            │   │   ├── [Si] <textarea> Compromisos  (opcional)
            │   │   ├── [No] <input type=date> Fecha próxima sesión  (opcional)
            │   │   └── [No] <textarea> Detalle próxima sesión  (opcional)
            │   ├── [scope === 'departamental']
            │   │   ├── <input type=date> Fecha próxima sesión o mesa* 
            │   │   └── <textarea> Detalle próxima sesión o mesa*
            │   ├── <textarea> Reporte de gestión*  (todos los scopes)
            │   ├── [scope !== 'nacional'] <textarea> Remisión a responsables  (opcional)
            │   └── [scope === 'nacional']
            │       ├── <input type=date> Fecha asistencias técnicas*
            │       └── <textarea> Asistencias técnicas*  (si no aplica: «No aplica» + explicación)
            │
            │   [tarea.type === 'decidir_via_riesgo_inminente']
            │   FormDecidirVia
            │   ├── Radio  Vía A — Asignación inmediata a Agente de Medidas
            │   ├── Radio  Vía B — Derivar a Supervisión del Convenio
            │   └── Select Asignar agente*  GET /emergency-measure-referrals/agentes
            │       Confirmar decisión exige una vía y un agente.
            │       PUT /case-tasks/:id/complete guarda la decisión.
            │
            │   [tarea.type === 'intentar_contacto_3x3']
            │   FormContacto
            │   ├── ¿Contestó?  Sí / No
            │   └── Nota
            │       Guardar cierra y no persiste.
            │
            │   [tarea.type === 'escalar_cruz_roja']
            │   Texto  "Confirmar Derivar el caso a Cruz Roja para entrega inmediata de la medida"
            │       Confirmar cierra y no persiste.
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
                    │     'respuesta_entidad'   → "Enviar respuesta"
                    │     'articulacion_territorial' → "Guardar"
                    │     'gestion_mecanismo'        → "Guardar"
                    │     'decidir_via_riesgo_inminente' → "Confirmar decisión"  (vía y agente)
                    │     'intentar_contacto_3x3' → "Guardar"  (cierra, no persiste)
                    │     'escalar_cruz_roja'     → "Confirmar"  (cierra, no persiste)
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
| `solicitud_entidad` | — | Valida ≥1 ítem con descripción. Resuelve exactamente un `general_user` con `entity_branch_id` = sede del `entity_case`. Si hay 0 o N usuarios, **no** marca Done y responde 400. Si hay uno: transacción — marca esta tarea Done + `form_data` y crea N `case_task` `respuesta_entidad` (`ToDo`) asignadas a ese usuario, una por ítem. Tras éxito, crea `case_timeline_event` categoría Entidades, type "Solicitud a entidad" (si el insert falla, WARN; la tarea permanece Done). No cambia el status de `entity_case_obligation`. |
| `respuesta_entidad` | — | Guarda JSON fusionado (contexto de la solicitud + radicado, ampliación, estado) + `Done`. Crea `case_timeline_event` categoría Entidades, type "Respuesta institucional" (si el insert falla, WARN; la tarea permanece Done). No cambia `entity_case_obligation` ni crea más tareas. |
| `articulacion_territorial` | — | `POST /api/v1/case-tasks/articulacion-territorial`. Si hay `taskId` ToDo válida: Done + `form_data`. Si no: crea tarea ya Done (`description=Gestión propia`). Si `requiereSeguimiento`: crea ToDo `compromiso_con_enlace` con `form_data.fechaSeguimiento`. No marca la barrera `Articulada`; si estaba `OPEN` pasa a `En Gestion`. |
| `gestion_mecanismo` | — | `POST /api/v1/case-tasks/gestion-mecanismo`. **Siempre inserta** una `case_task` ya `Done` (`category=Mecanismos`, `case_articulating_mechanism_id`, `form_data`, `result` resumen). Nunca completa una ToDo previa. El `case_articulating_mechanism` permanece `active`. |

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
>
> `respuesta_entidad`: radicado y ampliación opcionales; `estadoInstitucional` required (`iniciado` | `en_proceso` | `concluido`).
>
> `articulacion_territorial`: fecha de contacto y compromiso required; fecha de seguimiento required si Sí; sugerencias opcionales.
>
> `gestion_mecanismo`: `reporteGestion` required en los tres scopes. Municipal: `abordadoEnSesion` required; fecha/compromisos o próxima sesión opcionales según Sí/No. Departamental: `fechaProxima` + `detalleProxima` required. Nacional: `asistenciasFecha` + `asistenciasTexto` required. Remisión a responsables opcional (municipal/departamental).

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
| Nº de radicado / NUNC / MP (`respuesta_entidad`) | 80 |
| Ampliar información (`respuesta_entidad`) | 500 |
| Compromisos / detalle próxima / remisión / asistencias (`gestion_mecanismo`) | 500 |
| Reporte de gestión (`gestion_mecanismo`) | 2000 |
