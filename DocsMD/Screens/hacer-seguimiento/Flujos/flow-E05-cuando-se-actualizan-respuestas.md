# Flow E-05 — Cuando el formulario emite `answers-updated`

Pantalla: `/salvia/hacer-seguimiento/:id`  
Evento: `answers-updated` emitido por `dinamic-form`  
Función: `onAnswersUpdated(payload)` en `hacer_seguimiento.html`

`onAnswersUpdated` delega a tres métodos privados que corren en secuencia:
1. `_updatePsicosocialState(getVal, splitCSV)`
2. `_updateReasignacionState(getVal, splitCSV)`
3. `_updateBarreraLocationOptions(payload)`

---

## INPUT

| Campo | Origen |
|---|---|
| `payload.answers.directAnswers` | Array plano de `{ questionId, value }` — respuestas directas del submission completo |
| `payload.answers.repeaterEntries` | Array plano de `{ repeaterGroupId, iteration, answers, ... }` — todas las entries de repeaters |
| `this.caseRiskLevel` | Cargado en E-01 desde `data.victimInfo.riskLevel` (entero 1-4) |
| `this.allCities` | Cargado en E-01 por `loadLocations()` — `[{ label, value, departmentId }]` |
| `this._deptCityCache` | `{ [idx]: deptId }` — evita re-filtrar ciudades si el depto no cambió |
| `this._townCityCache` | `{ [idx]: cityId }` — evita re-fetchear municipios si la ciudad no cambió |

### IDs de preguntas evaluadas

| Constante | ID | Descripción |
|---|---|---|
| `Q_PSICO` | `71c42c4a-f640-47ad-b2c1-5d4c18480449` | Criterios remisión — Atención Psicosocial |
| `Q_PROTECTORES` | `a0fdcf67-b05b-4d92-9c19-14753a32bbe3` | Factores protectores |
| `Q_RIESGOS` | `ec5bb242-6f86-4c64-8b9f-afabd5a51878` | Factores de riesgo |
| `Q_EXTREMO` | `65f2d582-a39d-4c93-9519-1a20efbebb03` | Factores de riesgo extremo |
| `BARRIER_GROUP_ID` | `5fd3ecdc-2e5f-4b31-97ef-8a994580586a` | Repeater Sección 4 — Identificación de Barreras |
| `Q_DEPT` | `31c7f8ba-880e-4c9a-89f0-1a6a1e43b7b9` | Q12 Departamento (dentro del repeater) |
| `Q_CITY` | `c6f2d54a-3c61-4ae7-b2bf-50a884031fa4` | Q13 Ciudad (dentro del repeater) |

---

## Pasos

### Bloque 1 — Estado de remisión psicosocial (`_updatePsicosocialState`)

**1.** Parsear `payload.answers.directAnswers` con helpers `getVal(id)` y `splitCSV(v)`.

**2.** `selected = splitCSV(getVal(Q_PSICO))`

**3.** SI `selected` está vacío:
- `formState.psysocialRemisionState = ''`

**4.** SI `selected` tiene valores:
- Evaluar `tieneObligatorio = selected.includes('criterio_obligatorio')`
- Calcular `totalPuntos` sumando:

| Valor | Puntos |
|---|---|
| `conducta_suicida` | 3 |
| `interseccionalidad` | 2 |
| `sin_ruta` | 1 |
| `condiciones_territoriales` | 1 |
| `sin_acceso_psico` | 1 |
| `naturalizacion_vbg` | 1 |
| `criterio_obligatorio` | 0 (no suma) |

- `cumple = tieneObligatorio && totalPuntos >= 3`
- SI `cumple` → `formState.psysocialRemisionState = 'Remisión: SI cumple'`
- SI NO → `formState.psysocialRemisionState = 'Remisión: NO cumple'`

---

### Bloque 2 — Estado de reasignación de caso (`_updateReasignacionState`)

**5.** Leer y filtrar factores (excluir `"ninguno"` antes de contar):
- `protectores = splitCSV(getVal(Q_PROTECTORES)).filter(v => v !== 'ninguno')`
- `riesgos     = splitCSV(getVal(Q_RIESGOS)).filter(v => v !== 'ninguno')`
- `extremos    = splitCSV(getVal(Q_EXTREMO)).filter(v => v !== 'ninguno')`
- `level       = this.caseRiskLevel`

**6.** Reset de todos los valores de reasignación:
```
formState.shouldReassignCase = ''
formState.reassingText       = ''
formState.canReassignHigh    = ''
formState.canReassignLow     = ''
```

**7.** Evaluar según nivel actual del caso:

| Condición | shouldReassignCase | reassingText | canReassignHigh | canReassignLow |
|---|---|---|---|---|
| `level ∈ [1,2]` + `extremos.length > 0` | `'true'` | `'Se identificó un factor de riesgo extremo. El caso será reasignado automáticamente al equipo de Riesgo Alto.'` | `''` | `''` |
| `level ∈ [1,2]` + `riesgos.length >= 4` (sin extremos) | `'true'` | `'El caso califica para reasignación al equipo de Riesgo Alto. Confirma para continuar.'` | `'true'` | `''` |
| `level >= 3` + `extremos.length == 0` + `protectores.length >= 3` | `'true'` | `'El caso califica para reasignación al equipo de Riesgo Bajo. Confirma para continuar.'` | `''` | `'true'` |
| Cualquier otro caso | `''` | `''` | `''` | `''` |

**8.** `formState` se pasa como prop `:form-state` a `dinamic-form`, que lo usa para:
- **`psysocialRemisionState`**: `render_modification` SET en el banner psicosocial
- **`shouldReassignCase`**: `visibility_condition` del banner de reasignación (order 9, `504cdad6`)
- **`canReassignHigh`**: `visibility_condition` de la pregunta "¿Confirmar reasignación a riesgo alto?" (order 10, `df7a0293`)
- **`canReassignLow`**: `visibility_condition` de la pregunta "¿Confirmar reasignación a riesgo bajo?" (order 11, `7ec8d66d`)
- **`reassingText`**: `render_modification` REPLACE en el banner de reasignación

---

### Bloque 3 — Opciones de ubicación por barrera (`_updateBarreraLocationOptions`)

**9.** Filtrar entries del repeater de Sección 4:
```
entries = payload.answers.repeaterEntries
  .filter(e => e.repeaterGroupId === BARRIER_GROUP_ID)
```

**10.** Asegurar que `formState.newBarriers` tiene al menos `entries.length` posiciones:
- Por cada posición faltante → push `{ cities: [], towns: [] }`

**11.** Para cada `entry` en `entries` (por índice `idx`):

  a. Leer `deptId = entry.answers.find(Q_DEPT)?.value`  
  b. Leer `cityId = entry.answers.find(Q_CITY)?.value`

  c. **Actualizar ciudades** (solo si el departamento cambió):
  ```
  SI deptId !== _deptCityCache[idx]:
    _deptCityCache[idx] = deptId
    formState.newBarriers[idx].cities = allCities
      .filter(c => c.departmentId === deptId)
      .map(c => ({ label, value }))
  ```

  d. **Actualizar municipios** (solo si la ciudad cambió):
  ```
  SI cityId !== _townCityCache[idx]:
    _townCityCache[idx] = cityId
    SI cityId:
      GET /api/v1/locations/towns?city_id={cityId}
      → formState.newBarriers[idx].towns = response
    SI NO:
      formState.newBarriers[idx].towns = []
  ```

> Si `allCities` aún no está cargado (request de `loadLocations` pendiente), `cities` queda `[]`. Se populará cuando `answers-updated` vuelva a dispararse.

---

## Resultado visible en el formulario

| Nivel caso | Factores | Banner reasignación | Pregunta confirmación |
|---|---|---|---|
| Bajo (1-2) | Factor extremo seleccionado | Aviso reasignación automática | No aparece |
| Bajo (1-2) | ≥ 4 factores de riesgo (sin extremo) | Aviso reasignación + confirmar | Sí: "¿Confirmar a riesgo alto?" |
| Alto (3-4) | ≥ 3 protectores + sin extremo | Aviso reasignación + confirmar | Sí: "¿Confirmar a riesgo bajo?" |
| Cualquier otro | — | No aparece | No aparece |

| Sección 4 — Barrera | Efecto |
|---|---|
| Usuario selecciona departamento (Q12) | Q13 muestra ciudades de ese departamento |
| Usuario selecciona ciudad (Q13) | Q14 hace fetch y muestra municipios de esa ciudad |
| Departamento no cambia entre saves | No re-filtra ciudades (cache) |
| Ciudad no cambia entre saves | No re-fetchea municipios (cache) |

---

## GAPS

| # | Descripción | Impacto |
|---|---|---|
| G-01 | Si se agrega un criterio nuevo a `Q_PSICO` en BD, el puntaje no se actualiza automáticamente — requiere cambio en `_updatePsicosocialState` y en `processFollowUpSubmission` | Medio |
| G-02 | El cálculo de criterios psicosociales está duplicado en frontend y backend | Medio |
| G-03 | `caseRiskLevel` se lee solo una vez en E-01. Si el nivel del caso cambia en otra sesión, el valor en pantalla queda desactualizado hasta que se recargue | Bajo |
| G-04 | Si `loadLocations` termina después de que `answers-updated` se emite en carga inicial, las opciones de Q12 quedan vacías en el primer render | Bajo |
