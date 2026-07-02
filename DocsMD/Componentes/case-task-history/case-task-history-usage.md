# `case-task-history` — Guía de uso

Modal de solo lectura que muestra el detalle de una tarea (`case_task`) ya completada: descripción, quién la completó, fecha y el desglose legible de su `formData`. El padre lo controla mediante una referencia Vue (`ref`) y el método público `open(taskId)` — igual que `case-task-modal`.

---

## Props

El componente no recibe props. Se controla exclusivamente mediante el método `open(taskId)`.

---

## Método público

| Método | Parámetros | Descripción |
|---|---|---|
| `open(taskId)` | `taskId: string` — UUID de la `case_task` ya completada | Carga la tarea del API y muestra el modal con su detalle |

---

## Eventos emitidos

Ninguno. Es un componente de presentación pura — no emite eventos al padre.

---

## Qué hace el componente por sí solo

No hace nada hasta que el padre llama a `open(taskId)`. En ese momento carga la tarea del API y la muestra en modo lectura. Al montarse no precarga ningún catálogo (a diferencia de `case-task-modal`, no tiene dropdowns).

---

## Integración

Montado en la pantalla `Detalle del Caso`, junto a `case-task-modal`:

```html
<!-- src/frontend/html/salvia/case_detail/get_case_detail_sv.html -->
<case-task-modal
  ref="taskModal"
  :current-user-id="userICode"
  @completed="onCaseTaskCompleted">
</case-task-modal>

<case-task-history ref="taskHistory"></case-task-history>
```

```javascript
methods: {
  verDetalle(taskId) {
    this.$refs.taskHistory.open(taskId);
  }
}
```

Actualmente el trigger real está expuesto solo vía un panel temporal de desarrollo (pestaña Derivaciones, junto al de `case-task-modal`) — un input de UUID + botón "Ver detalle" que llama `$refs.taskHistory.open(testHistoryTaskId)`. `case-tasks.js` (el componente que lista las tareas del caso) ya separa `pendingTasks`/`completedTasks` internamente pero todavía no renderiza la pestaña de completadas ni emite un evento para abrir este modal — conectar ese trigger real es trabajo pendiente de quien mantiene `case-tasks.js`.

---

## Endpoint que consume

| Acción | Método | URL |
|---|---|---|
| Cargar tarea | GET | `/api/v1/case-tasks/:taskId` |

> ✅ `GetByID` fue actualizado para devolver la tarea enriquecida con `assignedUserName` (antes solo el endpoint de listado lo hacía). Mismo shape que el listado, vía el helper compartido `taskToJSON`.

### Shape real de la respuesta del endpoint

```json
{
  "id": "uuid",
  "caseId": "string",
  "category": "string",
  "type": "gestion_llamada | proyectar_oficio | comite_caso | Corregir oficio",
  "description": "string",
  "status": "Done",
  "assignedUserId": "icode",
  "assignedUserName": "string",
  "barrierId": "uuid | null",
  "entityLetterId": "uuid | null",
  "followUpId": "uuid | null",
  "result": "string | null",
  "completedAt": "timestamp",
  "createdAt": "timestamp",
  "updatedAt": "timestamp",
  "formData": { }
}
```

> Nota de datos reales: `assignedUserName` viene vacío si `assignedUserId` también lo está — pasa en tareas sembradas manualmente para pruebas (`gestion_llamada`, `proyectar_oficio`, `comite_caso` en el caso de prueba E2E). El componente muestra "Sin asignar" en ese caso. `Corregir oficio` sí trae un agente real porque se crea automáticamente desde `crearCaseTaskCorreccion()`, que copia el agente del `entity_letter`.

---

## Ejemplos de `formData` que renderiza (ver detalle completo en `case-task-modal-usage.md`)

#### `gestion_llamada`
```json
{
  "departamentoNombre": "Bogotá D.C.",
  "ciudadNombre":       "Bogotá",
  "municipioNombre":    "Bogotá D.C.",
  "entidadNombre":      "Fundación La Luz",
  "funcionario":        "Coordinadora de atención a víctimas",
  "descripcion":        "Llamada de seguimiento",
  "generaOficio":       false,
  "asunto":             null,
  "rutaKofax":          null
}
```

#### `comite_caso`
```json
{
  "decisionesTexto":              ["Activar Enlace", "Generar Oficio"],
  "observacionesOficio":          "Oficio para seguimiento de ruta institucional",
  "observacionesRecomendaciones": null,
  "nivelMecanismoTexto":          null,
  "observacionesMecanismo":       null
}
```

#### `Corregir oficio`
```json
{}
```
El componente muestra la nota **"El oficio fue corregido por el agente que lo proyectó"** en vez de intentar renderizar campos.
