# `dinamic-form` — Guía de uso

Componente Vue 3 reutilizable que renderiza un formulario dinámico completo a partir de un `formId` y un `submissionId`. Gestiona internamente la navegación por secciones, el guardado por sección, la validación, los repeaters y las condiciones de visibilidad.

Archivo principal: `src/frontend/html/salvia/components/dinamic_form.html`

---

## Props

| Prop | Tipo | Requerido | Default | Descripción |
|---|---|---|---|---|
| `formId` | String | Sí | — | UUID del formulario a renderizar |
| `submissionId` | String | No | `''` | UUID del submission existente. Vacío = crear nuevo al guardar la primera sección |
| `canEdit` | Boolean | No | `true` | Si `false`, el formulario se muestra en modo solo lectura completo |
| `formState` | Object | No | `{}` | Estado externo del padre. Se pasa a cada sección para resolver `visibility_condition` de tipo `formState` y `render_modification` de tipo REPLACE/SET |

---

## Eventos emitidos

| Evento | Cuándo se dispara | Payload |
|---|---|---|
| `form-completed` | Cuando el usuario guarda la última sección y todas las secciones visibles quedan respondidas | `{ submissionId: string }` |
| `answers-updated` | Cada vez que el usuario guarda una sección (incluyendo secciones intermedias) | `{ submissionId: string, answers: SubmissionStructure }` |

---

## Ejemplos de uso

### Caso 1 — Formulario nuevo (sin submission previo)

```html
<dinamic-form
  :form-id="followUp.formId"
  @form-completed="onFormCompleted"
></dinamic-form>
```

### Caso 2 — Continuar submission existente

```html
<dinamic-form
  :form-id="followUp.formId"
  :submission-id="followUp.formSubmissionId"
  @form-completed="onFormCompleted"
></dinamic-form>
```

### Caso 3 — Solo lectura (canEdit = false)

```html
<dinamic-form
  :form-id="followUp.formId"
  :submission-id="followUp.formSubmissionId"
  :can-edit="false"
></dinamic-form>
```

### Caso 4 — Con formState y answers-updated (para render_modification reactivo)

```html
<dinamic-form
  :form-id="followUp.formId"
  :submission-id="followUp.formSubmissionId"
  :can-edit="canEdit"
  :form-state="formState"
  @form-completed="onFormCompleted"
  @answers-updated="onAnswersUpdated"
></dinamic-form>
```

```js
data() {
  return {
    formState: { psysocialRemisionState: '' }
  }
},
methods: {
  onAnswersUpdated(payload) {
    // Computar estado reactivo a partir de las respuestas
    // y actualizar formState para que dinamic-form lo consuma
    // en sus render_modifications y visibility_conditions
    const answer = payload.answers?.directAnswers?.find(
      a => a.questionId === Q_PSICO
    )
    this.formState.psysocialRemisionState = computeState(answer)
  }
}
```

---

## Qué hace el componente por sí solo

- Carga la estructura del formulario (`GET /api/v1/forms/:id/load`) con el `submissionId` y el `formState` actuales.
- Determina la `currentSection` (primera sección visible no respondida).
- Renderiza la sección activa con sus preguntas, repeaters y condiciones de visibilidad.
- Aplica `render_modification` sobre labels y descripciones usando `formState`.
- Al guardar una sección llama `POST /api/v1/forms/saveSection` y actualiza el estado interno.
- Emite `answers-updated` después de cada guardado exitoso.
- Cuando todas las secciones visibles quedan respondidas, emite `form-completed`.
- Si `canEdit = false`, deshabilita todos los inputs y oculta los botones de guardado.
