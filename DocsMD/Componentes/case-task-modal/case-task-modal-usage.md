# `case-task-modal` — Guía de uso

Componente modal para completar tareas de tipo `gestion_llamada`, `proyectar_oficio` y `comite_caso`. El padre lo controla mediante una referencia Vue (`ref`) y el método público `open(taskId)`.

---

## Props

El componente no recibe props. Se controla exclusivamente mediante el método `open()` y el evento `@completed`.

---

## Método público

| Método | Parámetros | Descripción |
|---|---|---|
| `open(taskId)` | `taskId: string` — UUID de la `case_task` | Carga la tarea del API, inicializa el formulario según su tipo y muestra el modal |

---

## Eventos emitidos

| Evento | Cuándo se dispara | Payload |
|---|---|---|
| `@completed` | Cuando la tarea se completa exitosamente (status 200) | Objeto `case_task` actualizado con `status: 'Done'` y `form_data` |

---

## Qué hace el componente por sí solo

Al montarse, carga la lista de departamentos y ciudades desde el API para pre-poblar los dropdowns de los formularios `gestion_llamada` y `proyectar_oficio`. No muestra el modal hasta que el padre llame a `open()`.

---

## Integración en el HTML de la pantalla

```html
{{ template "tasks/case_task_modal.html" . }}
```

---

## Tipos de tarea soportados

| `case_task.type` | Título del modal | Formulario que renderiza |
|---|---|---|
| `gestion_llamada` | "Gestión de Llamada" | Ubicación + Funcionario + Descripción + Switch ¿Genera oficio? (condicional: Asunto + Ruta Kofax) |
| `proyectar_oficio` | "Proyectar Oficio" | Ubicación + Funcionario + Asunto + Ruta Kofax |
| `comite_caso` | "Decisiones del Comité" | Checkbox group de decisiones + campos condicionales por decisión |

---

## Ejemplo de uso

```javascript
// En el template de la pantalla padre:
// <case-task-modal ref="taskModal" @completed="onTaskCompleted"></case-task-modal>

methods: {
  abrirTarea(taskId) {
    this.$refs.taskModal.open(taskId);
  },
  onTaskCompleted(task) {
    // Reflejar el nuevo estado en la UI o recargar datos
    this.cargarDatos();
  }
}
```

---

## Endpoint que consume

| Acción | Método | URL |
|---|---|---|
| Cargar tarea | GET | `/api/v1/tasks/:taskId` |
| Completar tarea | PUT | `/api/v1/tasks/:taskId/complete` |
| Cargar departamentos y ciudades | GET | `/api/v1/locations` (al montar) |
| Cargar municipios | GET | `/api/v1/locations/towns?city_id={id}` (al seleccionar ciudad) |
| Cargar sedes de entidades | GET | `/api/v1/entity-branches?town_code={code}` (al seleccionar municipio) |

### Payload de `PUT /api/v1/tasks/:taskId/complete`

```json
{
  "userId":   "icode del usuario en sesión",
  "formData": { }
}
```

#### `formData` para `gestion_llamada`
```json
{
  "departamentoId":  "string",
  "ciudadId":        "string",
  "municipioId":     "string",
  "entidadId":       "integer | 'otra'",
  "entidadNombre":   "string",
  "funcionario":     "string",
  "descripcion":     "string | null",
  "generaOficio":    true,
  "asunto":          "string | null",
  "rutaKofax":       "string | null"
}
```

#### `formData` para `proyectar_oficio`
```json
{
  "departamentoId":  "string",
  "ciudadId":        "string",
  "municipioId":     "string",
  "entidadId":       "integer | 'otra'",
  "entidadNombre":   "string",
  "funcionario":     "string",
  "asunto":          "string",
  "rutaKofax":       "string"
}
```

#### `formData` para `comite_caso`
```json
{
  "decisiones": ["activar_enlace", "oficio", "recomendaciones_agente", "mecanismo_articulador"],
  "observacionesOficio":        "string | null",
  "observacionesRecomendaciones": "string | null",
  "nivelMecanismo":             "Municipal | Departamental | Nacional | null",
  "observacionesMecanismo":     "string | null"
}
```
