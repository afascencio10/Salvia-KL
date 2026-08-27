# case-task-modal — Index

## Navegación

| Documento | Link |
|---|---|
| Interfaz | [case-task-modal-interface.md](case-task-modal-interface.md) |
| Uso | [case-task-modal-usage.md](case-task-modal-usage.md) |

**Changelogs:** [changelogAgo2026.md](changelogAgo2026.md)

---

## Resumen de eventos

| # | Evento | Tipo | Flujo |
|---|---|---|---|
| E01 | Cuando carga el componente | Lifecycle | [📄 Ver flujo](Flujos/flow-E01-cuando-carga-el-componente.md) |
| E02 | Cuando se abre con una tarea | User Interaction | [📄 Ver flujo](Flujos/flow-E02-cuando-se-abre-con-una-tarea.md) |
| E03 | Cuando cancela el modal | User Interaction | — |
| E04 | Cuando selecciona departamento | User Interaction | [📄 Ver flujo](Flujos/flow-E04-cuando-selecciona-departamento.md) |
| E05 | Cuando selecciona ciudad | User Interaction | [📄 Ver flujo](Flujos/flow-E05-cuando-selecciona-ciudad.md) |
| E06 | Cuando selecciona municipio | User Interaction | [📄 Ver flujo](Flujos/flow-E06-cuando-selecciona-municipio.md) |
| E07 | Cuando activa "¿Genera oficio?" | User Interaction | [📄 Ver flujo](Flujos/flow-E07-cuando-activa-genera-oficio.md) |
| E08 | Cuando cambia decisiones del comité | User Interaction | [📄 Ver flujo](Flujos/flow-E08-cuando-cambia-decisiones-del-comite.md) |
| E09 | Cuando confirma el formulario | User Interaction | [📄 Ver flujo](Flujos/flow-E09-cuando-confirma-formulario.md) |

---

## Inventario de eventos

📄 [Ver flujo → flow-E01-cuando-carga-el-componente.md](Flujos/flow-E01-cuando-carga-el-componente.md)

**Evento:** Cuando carga el componente
**Tipo:** Lifecycle
**Descripción:** Se ejecuta al montar el componente Vue (`mounted`). Carga la lista de departamentos y ciudades desde el API para pre-poblar los dropdowns de ubicación disponibles en los formularios de tipo `gestion_llamada` y `proyectar_oficio`.
**Requerido:** Sí

---

📄 [Ver flujo → flow-E02-cuando-se-abre-con-una-tarea.md](Flujos/flow-E02-cuando-se-abre-con-una-tarea.md)

**Evento:** Cuando se abre con una tarea
**Tipo:** User Interaction
**Descripción:** El padre llama al método público `open(taskId)`. El componente setea `taskId`, muestra el modal, resetea el formulario y hace `GET /api/v1/case-tasks/:taskId` para cargar la tarea. Una vez cargada, inicializa el formulario según `tarea.type`. Si `tarea.type === 'Corregir oficio'`, además hace `GET /api/v1/entity-letters/:letterId` (usando `tarea.entityLetterId`) para mostrar en solo lectura la razón de corrección y la ruta Kofax.
**Requerido:** Sí

---

**Evento:** Cuando cancela el modal
**Tipo:** User Interaction
**Descripción:** El usuario presiona "Cancelar", el botón "✕" del header o el overlay oscuro. Setea `visible = false`, limpia `tarea`, `form`, `saveError` y `errorTarea`. No emite ningún evento al padre.
**Requerido:** Sí

---

📄 [Ver flujo → flow-E04-cuando-selecciona-departamento.md](Flujos/flow-E04-cuando-selecciona-departamento.md)

**Evento:** Cuando selecciona departamento
**Tipo:** User Interaction
**Descripción:** Aplica únicamente en los formularios `gestion_llamada` y `proyectar_oficio`. Filtra `allCities` en memoria por `departamentoId`. Resetea `form.ciudadId`, `form.municipioId`, `form.entidadId`, `form.entidadNombre` y sus listas de opciones. Sin fetch.
**Requerido:** Sí

---

📄 [Ver flujo → flow-E05-cuando-selecciona-ciudad.md](Flujos/flow-E05-cuando-selecciona-ciudad.md)

**Evento:** Cuando selecciona ciudad
**Tipo:** User Interaction
**Descripción:** Aplica únicamente en los formularios `gestion_llamada` y `proyectar_oficio`. Hace `GET /api/v1/locations/towns?city_id={ciudadId}` para cargar los municipios. Resetea `form.municipioId`, `form.entidadId`, `form.entidadNombre` y sus listas.
**Requerido:** Sí

---

📄 [Ver flujo → flow-E06-cuando-selecciona-municipio.md](Flujos/flow-E06-cuando-selecciona-municipio.md)

**Evento:** Cuando selecciona municipio
**Tipo:** User Interaction
**Descripción:** Aplica únicamente en los formularios `gestion_llamada` y `proyectar_oficio`. Hace `GET /api/v1/entity-branches?town_code={municipioId}` para cargar las sedes de entidades del municipio. La respuesta incluye `{id, icode, name}`. La opción "Otra entidad" (`value="otra"`) siempre se agrega al final. Resetea `form.entidadId` y `form.entidadNombre`.
**Requerido:** Sí

---

📄 [Ver flujo → flow-E07-cuando-activa-genera-oficio.md](Flujos/flow-E07-cuando-activa-genera-oficio.md)

**Evento:** Cuando activa "¿Genera oficio?"
**Tipo:** User Interaction
**Descripción:** Aplica únicamente en el formulario `gestion_llamada`. El usuario activa o desactiva el switch. Actualiza `form.generaOficio`. Vue reactivamente muestra u oculta los campos de Asunto y Ruta al archivo. Si se desactiva, limpia esos campos.
**Requerido:** Sí

---

📄 [Ver flujo → flow-E08-cuando-cambia-decisiones-del-comite.md](Flujos/flow-E08-cuando-cambia-decisiones-del-comite.md)

**Evento:** Cuando cambia decisiones del comité
**Tipo:** User Interaction
**Descripción:** Aplica únicamente en el formulario `comite_caso`. El usuario marca o desmarca una opción del checkbox group. Actualiza `form.decisiones`. Vue reactivamente muestra u oculta los bloques de campos condicionales: Observaciones del oficio, Observaciones de recomendaciones, y el bloque de Nivel + Observaciones del mecanismo. Si se desmarca una opción, limpia los campos de su bloque.
**Requerido:** Sí

---

📄 [Ver flujo → flow-E09-cuando-confirma-formulario.md](Flujos/flow-E09-cuando-confirma-formulario.md)

**Evento:** Cuando confirma el formulario
**Tipo:** User Interaction
**Descripción:** El usuario presiona "Completar tarea" (o "Marcar como corregido" si `tarea.type === 'Corregir oficio'`). Valida los campos requeridos según `tarea.type` (ver tabla en interface) — `Corregir oficio` no tiene campos, solo requiere que haya terminado de cargar el oficio vinculado. Si válido, construye el payload con `formData` como JSON (`{}` para `Corregir oficio`) y llama `PUT /api/v1/case-tasks/:taskId/complete`. Si responde 200, cierra el modal y emite `@completed` con la tarea actualizada. Si error, muestra `saveError` sin cerrar el modal.
**Requerido:** Sí

---

## Checklist de completitud

- [x] ¿Se cubre la carga inicial de datos?
- [x] ¿Todos los campos interactivos del formulario tienen evento?
- [x] ¿Las validaciones de formulario están cubiertas?
- [x] ¿Los estados de error de carga y guardado están cubiertos?
- [x] ¿El callback al padre está cubierto?
