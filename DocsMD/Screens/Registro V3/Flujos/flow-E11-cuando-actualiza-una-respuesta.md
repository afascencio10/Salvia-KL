━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando actualiza una respuesta del formulario
   Tipo: User Interaction
   Estado: 🟡 Implementado con cambios (subconjunto de lo planeado) — sincronizado con el código 2026-10-31
   Funciones: setVal() · setBool() · toggleMarcado() · tzToggle() · guardarTamizaje() → registro-v3-app.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

> Sin `dinamic-form`: cada control escribe en `answers[campoId]` y Vue re-evalúa los bindings. **No hay guardado ni
> llamada al backend por cada cambio** (se guarda por sección, flow-E12, y el Tamizaje con su propio botón). Todo lo que
> sigue ocurre dentro del mismo evento porque es reactividad del framework.

INPUT: {
  campoId:  FieldKey de la pregunta (CI_1, DP_5, HV_9, TZ_3, GS_6…)   → control que disparó el cambio
  valor:    texto | "true"/"false" | icode de enum | CSV de icodes     → input del agente
}

PASO 1 — Guardar el valor en el estado local
  SEGÚN el tipo de control:
    text / textarea / date / dropdown → setVal(id, v):  answers[id] = v
    boolean (radio Sí/No)             → setBool(id, v): answers[id] = 'true' | 'false'
    multiple (checkboxes)             → toggleMarcado(id, icode, max): agrega o quita el icode del CSV;
                                        si ya hay `max` marcados, ignora el nuevo (HV_11 tiene max 2)
    radios del Tamizaje TZ_1…TZ_7     → tzToggle(id, v) (ver PASO 4)

PASO 2 — Efectos por campo (dentro de setVal)
  SI campoId === 'CI_5' (tipo de documento):
    → sugerirNacionalidadPorDocumento():  SI DP_5 sigue vacío (nunca pisa una respuesta manual):
        code ∈ {cc, ti, rc, nc} → DP_5 = "Colombiana" (code 'co')
        code ∈ {ce, ps, pp, pf, pt, dn, vs, vr, sc, mv, no, pn} → DP_5 = "Extranjera" (code 'ex')
        cd / oe / si → no sugiere (ambiguos)
  SI campoId ∈ {'CI_5', 'CI_6'}:
    → verificarDuplicidadDocumento():  espera 500 ms sin cambios (debounce);
        SI hay tipo (code) y número:
          [GET] /api/v1/victim-cases-v3/buscar?tipo_documento=<code>&documento=<número>
          SI hay resultados con icode distinto del actual:
            avisoDuplicidad = "Ya existe un Registro con este documento: <nombres>." (o "Ya existen N Registros …")  → se muestra bajo CI_6
          SI NO: sin aviso
        No bloquea el guardado
  SI campoId === 'CI_1': el encabezado del Formulario recalcula el tratamiento (code la → "Sra.", el → "Sr.", le → "Sre."); también lo usa el párrafo síntesis

PASO 3 — Re-evaluar visibilidad y validaciones en vivo (Vue)
  mostrarCampo(campo) por cada pregunta del módulo; condicionales vigentes en el catálogo:
    CI_6    visible si CI_5 ≠ "Persona sin identificación" (code si)
    CI_8    visible si CI_7 = Sí            ·  DP_4  visible si CI_7 = Sí
    DP_3    visible si DP_2 = Sí
    DP_6    visible si DP_5 tiene code 'ex'
    DP_13   visible si la etiqueta de DP_12 contiene "Indígena"
    DP_19   visible si DP_18 tiene code 'em' o 'ta'
    DP_20-22 visibles si DP_17 tiene code 'as' o DP_18 tiene code 'pr' (persona en ASP)
    DP_26   siempre visible (la condición por sexo/identidad planeada devuelve true)
    HV_10   visible si hay tipos marcados en HV_9; sus opciones son los subtipos de los tipos marcados
    HV_12   visible si HV_11 incluye un ámbito con "laboral"
    AG_5    visible si la etiqueta de AG_2 no contiene "desconocid"
  Validaciones visibles al escribir (no bloquean hasta guardar, flow-E12):
    "solo letras" (CI_2, CI_3, CA_1, AG_7): campo en rojo + "No puede ser solo números."
    CI_6: texto "Formato esperado: …" según el tipo de documento, y "No cumple el formato esperado." si no coincide (máscaras en registro-v3-campos.js)
  Si un campo se oculta, su respuesta NO se borra de `answers`.

PASO 4 — Tamizaje: revelar el resto (solo TZ_1…TZ_7)
  tzToggle(id, v): answers[id] = 'true'|'false'
  tzRevelarResto = alguna de las 5 iniciales (TZ_1, TZ_2, TZ_3, TZ_4, TZ_7) está en 'true'
  SI tzRevelarResto: se muestran TZ_5, TZ_6, TZ_8 (VPP), la rama de 24 preguntas y el botón "Calcular y guardar Tamizaje"
  SI NO: solo se ven las 5 iniciales (+ TZ_5/TZ_6 ocultas); no se muestra ningún botón de guardado (ver GAP)
  Rama: tzWasPartner() = la etiqueta elegida en AG_3 contiene "pareja" → 24 preguntas de pareja (TZP_1…24); si no, 24 de no pareja (TZNP_1…24)

PASO 5 — "Calcular y guardar Tamizaje" (guardarTamizaje)
  wasPartner = tzWasPartner()
  esOsigd    = (DP_9 respondida y su etiqueta no contiene "heterosexual") o (la etiqueta de DP_8 contiene "trans" o "no binaria")
  subtipoFeminicidio / subtipoHomicidioAbyecto = HV_10 incluye "tentativa feminicidio" / "tentativa homicidio"
  tzAnswers = TZ_1…TZ_8 + las 24 de la rama activa, cada una 'true' o 'false'
  [POST] /api/v1/victim-cases-v3/{icode}/tamizaje
  { "answers": tzAnswers, "wasPartner": wasPartner, "esOsigd": esOsigd, "subtipoFeminicidio": …, "subtipoHomicidioAbyecto": … }

  Backend (SaveTamizaje):
    riskType = esOsigd ? "homicidio_abyecto" : "feminicidio"
    SEGÚN subtipos de tentativa:
      ambos      → riskType = "feminicidio" y tieFlag = true      // G-V3-02 resuelto: prioridad a Feminicidio
      solo fem.  → "feminicidio"
      solo hom.  → "homicidio_abyecto"
    score = número de respuestas "true" en `answers` (7 generales + VPP + 24 de la rama)
    nivel = primer umbral con score ≤ máximo:
      feminicidio + pareja:      ≤4 bajo · ≤8 moderado · ≤15 alto · más: extremo
      feminicidio + no pareja:   ≤2 bajo · ≤5 moderado · ≤8 alto  · más: extremo
      homicidio abyecto:         ≤5 bajo · ≤11 moderado · ≤16 alto · más: extremo   (misma escala para pareja y no pareja)
    Persiste: marca las keys como respondidas, actualiza las columnas de la rama activa y victim_case_form2_tzv3_risk_type / _risk_level / _risk_score / _risk_tie_flag
    Respuesta: { riskType, riskLevel, score, tieFlag }   (rangos "pendientes de ajuste final" según el cliente: no recalibrar)

  Cliente: tzResultado = respuesta
    Vue muestra "Tipo de riesgo: … — Nivel: … — Puntaje: N" (rojo si alto/extremo, con "⚠️ Nivel de riesgo alto — considere activar la Alerta Inminente en Gestión Salvia.") y, si tieFlag, la nota de empate
    SI TZP_24 o TZNP_24 = 'true': "⚠️ Se marcó ideación suicida. Debe abordarse con remisión psicosocial (Gestión Salvia § B, Atención Psicosocial)."

PASO 6 — Reglas de Gestión Salvia (módulo GS, solo visibilidad y avisos locales)
  GS_5 = Sí  → se listan los 6 equipos (GS_6)
  equipo marcado con criterios → checkboxes de criterios (GS_9 Psicosocial, GS_10 Hombres, GS_11 Discapacidad, GS_12 Estabilización)
    cumpleCriterioEquipo():  Psicosocial y Hombres exigen el criterio obligatorio 'ob'; Psicosocial y Estabilización suman puntos
      (Psicosocial: conducta suicida 3, interseccionalidad 2, y 1 cada una de las otras 4; mínimo 3 · Estabilización: 1 cada una, mínimo 1);
      Discapacidad pide al menos una opción. Se muestra "✓ Cumple criterios de remisión" o "✗ Aún no cumple los criterios mínimos" (informativo, no bloquea)
  equipo 'Medidas de Emergencia' marcado → checkboxes de medidas (GS_7)
  GS_14 = Sí o GS_16 = Sí → aparece GS_15 · GS_17 = Sí → habilita "+ Agregar barrera a esta entidad" · GS_18 = Sí → aparece GS_19
  Las entidades y barreras no pasan por aquí: ver flow-E17 y flow-E18

PASO 7 — Párrafo síntesis
  parrafoSintesis() = sintesisA (CI) + sintesisB (DP) + sintesisC (HV) + sintesisD (AG), cada una con sus reglas de omisión
  Se calcula en cada render y se muestra al entrar a GS y en el Cierre; no se persiste (ver registro-v3-parrafo-sintesis.md)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión | Paso afectado |
|---|---|
| **Tamizaje sin forma de guardarse si las 5 iniciales son "No":** el botón "Calcular y guardar Tamizaje" está dentro del bloque `v-show="tzRevelarResto"`, `guardarYSalir` no guarda TZ y `Finalize` exige TZ_1, 2, 3, 4 y 7 respondidas (solo se marcan al guardar el Tamizaje). Por lectura de código, un Caso con las 5 iniciales en "No" no podría finalizar. No verificado en el navegador (el QA TZ1 solo comprobó que el resto queda oculto) | PASO 4, PASO 5 |
| Planeado y no implementado en el cliente: barra de avance por módulo; deshabilitar opciones según identidad de género (Gay/Lesbiana, sexo asignado al nacer); limpiar las preguntas de escala; regla de gestación por sexo/identidad (DP_26 siempre visible); limpiar subtipos al desmarcar un tipo; modal de motivo al editar el Perfil de un Registro Completo; modal "Registro Completo"; bitácora de trazabilidad del Perfil (`{ts, agente, campo, antes, después, motivo}`) | PASO 2, PASO 3 |
| Planeado y no implementado: cálculo del Tamizaje en cada cambio de DP/HV/TZ, ocultar preguntas OSIGD y de ámbito laboral (las constantes `TAMIZAJE_*_OSIGD_V3` y `*_LABORAL_V3` no se usan: se muestran las 24 de la rama y TZ_8 siempre), y alertas a coordinación / equipo psicosocial (`NotificationService`): hoy solo hay avisos en pantalla | PASO 4, PASO 5 |
| Planeado y no implementado en GS: advertencia al marcar Medidas de Emergencia y Atención Psicosocial a la vez, y presugerencia de equipo según el nivel de riesgo o la ideación | PASO 6 |
| El fan-out del Perfil (CI/DP/CA) no ocurre en este evento sino en el backend al guardar la sección (flow-E12 paso 5) | — |
| `esOsigd` cuenta como OSIGD cualquier orientación sexual respondida distinta de "heterosexual" (incluye "No informa" u "Otro") y compara etiquetas por texto, no por code | PASO 5 |
