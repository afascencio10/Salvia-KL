# Flow E-02 — Cuando el formulario emite `answers-updated`

Pantalla: `/salvia/casos/nuevo-v2`
Evento: `answers-updated` emitido por `dinamic-form`
Función: `onAnswersUpdated(payload)` en la pantalla anfitriona (análoga a `hacer_seguimiento.html`)

`onAnswersUpdated` delega a cuatro bloques que corren en secuencia:
1. `_updateAggressorRelationState` — wasPartner / partnerKnown / tamizajeAggressorCheck
2. `_updateTamizajeScore` — riskScore / riskLevel / hasRisk / riskBadgeText
3. `_updateViolenceState` — hasWorkplaceScope / hasSelectedViolenceTypes / violenceSubtypes
4. `_updateLocationChains` — las 3 cadenas geográficas (residencia, hechos, atención)

---

## INPUT

| Campo | Origen |
|---|---|
| `payload.answers.directAnswers` | Array plano `{ questionId, value }` — respuestas directas del submission completo |
| `this.allCities` | Cargado en E-01 — `[{ label, value, departmentId }]` |
| `this._deptCityCache` / `this._townCityCache` | `{ [chainName]: id }` — evita re-filtrar/re-fetchear si no cambió, generalizado de índice numérico (Barreras) a nombre de cadena (residencia/hechos/atencion) |

---

## Bloque 1 — Estado de relación con el agresor (`_updateAggressorRelationState`)

**1.** `proximity = getVal(qProximityAggressor)` · `relationship = getVal(qRelationshipAggressor)`

**2.** `wasPartner = relationship en {'pi', 'ex'}` (pareja íntima / ex-pareja — mismo criterio que `relationWithAggressorChanged` hoy)

**3.** `partnerKnown = proximity en {'pc', 'pn'} OR wasPartner`  (persona conocida / pareja — mismo criterio que `proximityWithAggressorChanged` hoy)

**4.** Escribir en `formState`:
```
formState.wasPartner   = wasPartner ? 'true' : 'false'
formState.partnerKnown = partnerKnown ? 'true' : 'false'
```

> Estos dos flags controlan, vía `visibility_condition` con `trigger_state_path`, el bloque de tamizaje (pareja/no-pareja) y la pregunta "¿Depende económicamente?" — igual que `shouldReassignCase` controla el banner de reasignación en Seguimiento.

---

## Bloque 2 — Puntaje de tamizaje (`_updateTamizajeScore`)

**5.** Recolectar las respuestas de tamizaje visibles:
- 6 comunes (siempre)
- 18 de pareja SI `wasPartner == true`, o 14 de no-pareja SI `wasPartner == false`

**6.** `totalSi = cantidad de 'true' entre las preguntas recolectadas`

**7.** Determinar `riskLevel` según la tabla de umbrales (idéntica a la actual — ver `registro-caso-interface.md` § "Lógica del tamizaje"):

| `wasPartner` | Bajo | Moderado | Alto | Extremo |
|---|---|---|---|---|
| `true` (24 preguntas máx) | 0-4 | 5-8 | 9-15 | 16-24 |
| `false` (20 preguntas máx) | 0-2 | 3-5 | 6-8 | 9-20 |

**8.** Escribir en `formState`:
```
formState.hasRisk       = (riskLevel > 0) ? 'true' : 'false'
formState.riskBadgeText = `Puntaje de riesgo: ${totalSi} — Nivel: ${nombreNivel(riskLevel)}`
```

> ⚠️ Esta fórmula está **duplicada** entre este bloque (frontend, para el banner en vivo) y `getRiskScore()` en el backend (E-03, PASO 3) — mismo patrón de duplicación ya señalado como gap en Seguimiento.

---

## Bloque 3 — Estado de violencia (`_updateViolenceState`)

**9.** `types = splitCSV(getVal(qTypeViolence))`

**10.** `formState.hasSelectedViolenceTypes = (types.length > 0) ? 'true' : 'false'`

**11.** `formState.violenceSubtypes = types.flatMap(code => enums['victim_case_form2_subtype_violence_experienced_' + code] || [])`
— concatena las opciones de cada tipo seleccionado, igual que `violenceTypeChanged` hoy (los 7 sufijos de categoría: `fi, pl, po, ps, re, se, vi` — ver conteos en la tabla de enums).

**12.** `scope = splitCSV(getVal(qScopeOfViolence))`

**13.** `formState.hasWorkplaceScope = scope.includes('al') ? 'true' : 'false'`
— código `'al'` (ámbito laboral) tomado de `registro-caso-index.md`; confirmar contra el `code` real en `victim_case_form2_enums` (ver GAPS de `form-registro-caso-v2-data.md`, G-03).

---

## Bloque 4 — Cadenas geográficas (`_updateLocationChains`)

**14.** Para cada cadena `chain` en `['residencia', 'hechos', 'atencion']`:

  a. Leer `deptId = getVal(qDept[chain])` · `cityId = getVal(qCity[chain])`

  b. **Actualizar ciudades** (solo si el departamento cambió):
  ```
  SI deptId !== _deptCityCache[chain]:
    _deptCityCache[chain] = deptId
    formState.geo[chain].cities = allCities
      .filter(c => c.departmentId === deptId)
      .map(c => ({ label, value }))
  ```

  c. **Actualizar municipios** (solo si la ciudad cambió):
  ```
  SI cityId !== _townCityCache[chain]:
    _townCityCache[chain] = cityId
    SI cityId:
      GET /api/v1/locations/towns?city_id={cityId}
      → formState.geo[chain].towns = response
    SI NO:
      formState.geo[chain].towns = []
  ```

> Mismo patrón exacto que `_updateBarreraLocationOptions` en Seguimiento (flow-E05 de hacer-seguimiento), generalizado de "por índice de entry en un repeater" a "por nombre de cadena fija" — las 3 cadenas no son repetibles, son 3 preguntas nombradas.

---

## Resultado visible en el formulario

| Sección | Efecto |
|---|---|
| Agresor | "¿Depende económicamente?" aparece solo si `partnerKnown` |
| Tamizaje | Bloque de 18 (pareja) o 14 (no-pareja) preguntas según `wasPartner`; banner de riesgo aparece cuando `hasRisk` |
| Hechos | "Subtipo de violencia" muestra opciones según los tipos seleccionados; "Sector laboral" aparece solo si se marcó ámbito laboral |
| Datos de la Víctima / Hechos / Lugar de Atención | Ciudad se filtra por departamento elegido; Municipio se carga vía API al elegir ciudad — por cadena, independientes entre sí |

---

## GAPS

| # | Descripción | Impacto |
|---|---|---|
| G-08 | Confirmar código exacto de "ámbito laboral" en `victim_case_form2_scope_of_violence` (se asumió `'al'`) | Medio |
| G-09 | Confirmar sufijos de categoría usados por `subtype_violence_experienced_{code}` contra los 7 valores reales de `victim_case_form2_type_violence_experienced` | Medio |
| G-10 | Igual que en Seguimiento: si `loadLocations`/enums termina después del primer `answers-updated` en la carga inicial, las opciones quedan vacías en el primer render | Bajo |
