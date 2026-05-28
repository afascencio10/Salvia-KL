# Flow E-05 — Cuando el formulario emite `answers-updated`

Pantalla: `/salvia/hacer-seguimiento/:id`  
Evento: `answers-updated` emitido por `dinamic-form`  
Función: `onAnswersUpdated(payload)` en `hacer_seguimiento.html`

---

## INPUT

| Campo | Origen |
|---|---|
| `payload.answers.directAnswers` | Emitido por `dinamic-form` al guardar cada sección |
| `Q_PSICO = '71c42c4a-f640-47ad-b2c1-5d4c18480449'` | Constante hardcodeada — ID de la pregunta "Criterios de remisión — Atención Psicosocial" |

---

## Pasos

**1.** Buscar en `payload.answers.directAnswers` la respuesta cuyo `questionId === Q_PSICO`.

**2.** SI no se encuentra respuesta o `answer.value` está vacío:
- Setear `formState.psysocialRemisionState = ''`
- **FIN** (el banner no muestra estado)

**3.** SI se encuentra respuesta con valor:
- Parsear `answer.value` (CSV) → array `selected`

**4.** Evaluar criterio obligatorio:
- SI `selected` incluye `'criterio_obligatorio'` → `tieneObligatorio = true`
- SI NO → `tieneObligatorio = false`

**5.** Calcular puntaje total sumando los puntos de cada valor seleccionado:

| Valor | Puntos |
|---|---|
| `conducta_suicida` | 3 |
| `interseccionalidad` | 2 |
| `sin_ruta` | 1 |
| `condiciones_territoriales` | 1 |
| `sin_acceso_psico` | 1 |
| `naturalizacion_vbg` | 1 |
| `criterio_obligatorio` | 0 (no suma puntos) |

**6.** Evaluar si cumple: `cumple = tieneObligatorio && totalPuntos >= 3`

**7.** SI `cumple`:
- Setear `formState.psysocialRemisionState = 'Remisión: SI cumple'`

**8.** SI NO `cumple`:
- Setear `formState.psysocialRemisionState = 'Remisión: NO cumple'`

**9.** `formState` se pasa como prop `:form-state` a `dinamic-form`, que lo usa para resolver el `render_modification` de tipo REPLACE sobre el banner informativo de remisión psicosocial:
- `search_string = 'Remisión:'`
- Se reemplaza por el valor de `formState.psysocialRemisionState`

---

## Resultado visible

| Criterios seleccionados | Banner muestra |
|---|---|
| Ninguno / sin respuesta | *(vacío — sin estado)* |
| `criterio_obligatorio` solo (0 pts) | `Remisión: NO cumple` |
| `conducta_suicida` solo sin obligatorio (3 pts) | `Remisión: NO cumple` |
| `criterio_obligatorio` + `conducta_suicida` (3 pts) | `Remisión: SI cumple` |
| `criterio_obligatorio` + `interseccionalidad` + `sin_ruta` (3 pts) | `Remisión: SI cumple` |

---

## GAPS

| # | Descripción | Impacto |
|---|---|---|
| G-01 | Si se agrega un criterio nuevo a la pregunta en BD, el puntaje no se actualiza automáticamente — requiere cambio en `onAnswersUpdated` del frontend y en `processFollowUpSubmission` del backend | Medio |
| G-02 | El mismo cálculo está duplicado en frontend (`onAnswersUpdated`) y backend (`processFollowUpSubmission`). Si cambia la regla hay que actualizarlo en ambos lugares | Medio |
