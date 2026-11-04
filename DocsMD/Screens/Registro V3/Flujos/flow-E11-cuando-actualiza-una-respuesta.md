━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando actualiza una respuesta del formulario
   Tipo: User Interaction
   Estado: 🟡 Implementado con cambios (subconjunto de lo planeado) — sincronizado con el código 2026-10-31 · actualizado 2026-11-02 (WP-2 · REG-08/09/10/12/34) · actualizado 2026-11-04 (WP-4)
   Funciones: setVal() · setBool() · toggleMarcado() · tzToggle() · guardarTamizaje() · etiqueta() → registro-v3-app.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

> 🔄 **WP-4 (2026-11-04):** (1) **GS_20** «elementos psicosociales» usa `toggleMarcado`; si el campo declara `excluyente` (código 99,
> «Ninguno») marcarlo desmarca las demás y marcar otra lo quita; los demás campos múltiples no cambian (REG-16). (2) **Formulario de
> barrera** (flow-E18): cada respuesta pasa por `setBarrera`; cambiar el sector limpia lo que dejó de aplicar, elegir un departamento
> carga sus municipios, y la perspectiva profesional solo existe con la gestión «alerta». (3) **Formulario de entidad** (flow-E17): la
> cascada Departamento → Ciudad → Municipio → Sector es obligatoria y limpia los niveles inferiores; «Buscar entidad» filtra en el
> navegador con debounce de 250 ms (sin una petición por tecla); marcar una acción como completada la quita de pendientes y viceversa;
> «Ninguna» es excluyente entre los canales. (4) `GS_15` ya no depende de `GS_14`/`GS_16` (siempre visible) y `GS_16` = Sí despliega el
> formulario de activación; con `GS_15` y `GS_16` en «Sí» se piden canales y gestiones.


> Sin `dinamic-form`: cada control escribe en `answers[campoId]` y Vue re-evalúa los bindings. **No hay guardado ni
> llamada al backend por cada cambio** (se guarda por sección, flow-E12, y el Tamizaje con su propio botón). Todo lo que
> sigue ocurre dentro del mismo evento porque es reactividad del framework.

INPUT: {
  campoId:  FieldKey de la pregunta (CI_1, DP_5, HV_9, TZ_3, GS_6…)   → control que disparó el cambio
  valor:    texto | "true"/"false" | icode de enum | CSV de icodes     → input del agente
}

PASO 1 — Guardar el valor en el estado local
  SEGÚN el tipo de control:
    text / textarea / date / time / dropdown → setVal(id, v):  answers[id] = v      // time = <input type="time" step="60"> (HV_3, REG-09)
    boolean (radio Sí/No)             → setBool(id, v): answers[id] = 'true' | 'false'
                                        SI id === 'DP_2' y v es No: answers['DP_51'] = ''   // el tipo de discapacidad (REG-08) queda oculto y se limpia; el servidor lo guarda como NULL
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
    DP_3    visible si DP_2 = Sí             ·  DP_51 (Tipo de discapacidad, 7 opciones, no obligatoria) visible si DP_2 = Sí
    DP_6    visible si DP_5 tiene code 'ex'
    DP_13   visible si la etiqueta de DP_12 contiene "Indígena"
    DP_19   visible si DP_18 tiene code 'em' o 'ta'
    DP_20-22 visibles si DP_17 tiene code 'as' o DP_18 tiene code 'pr' (persona en ASP)
    DP_26   siempre visible (la condición por sexo/identidad planeada devuelve true)
    HV_10   visible si hay tipos marcados en HV_9; sus opciones son los subtipos de los tipos marcados
            (su texto es el de /enums: etiqueta() solo busca textos propios en las categorías del campo; antes un código repetido en otra
             categoría —tr, al, es, ah, ai, pr— daba la etiqueta ajena: «Trata de personas…» salía «Transporte», REG-10)
    HV_12   visible si HV_11 incluye un ámbito con "laboral"
    AG_5    visible si la etiqueta de AG_2 no contiene "desconocid"
  Validaciones visibles al escribir (no bloquean hasta guardar, flow-E12):
    "solo letras" (CI_2, CI_3, CA_1, AG_7): campo en rojo + "No puede ser solo números."
    CI_6: texto "Formato esperado: …" según el tipo de documento, y "No cumple el formato esperado." si no coincide (máscaras en registro-v3-campos.js)
  Si un campo se oculta, su respuesta NO se borra de `answers`.

PASO 4 — Tamizaje por etapas (REG-12 / REG-34): revelar el resto y mostrar el botón
  Las preguntas se presentan en dos tablas «Pregunta | Sí | No» (columnas 70/15/15, un par de radios por fila con aria-label)
  tzToggle(id, v): answers[id] = 'true'|'false'
  tzRevelarResto = alguna de las 5 iniciales (TZ_1, TZ_2, TZ_3, TZ_4, TZ_7) está en 'true'
  tzInicialesRespondidas = las 5 iniciales tienen respuesta ('true' o 'false')            // computed
  SI tzRevelarResto: se muestran TZ_5, TZ_6, TZ_8 (VPP) y la rama de 24 preguntas (segunda tabla)
  SI tzRevelarResto O tzInicialesRespondidas: se muestra el botón "Calcular y guardar Tamizaje"
  SI tzInicialesRespondidas Y NO tzRevelarResto: nota "Las cinco preguntas iniciales están en «No»: no se pregunta el resto del Tamizaje. Ya puedes calcular y guardar el resultado."
  Al abrir un Caso (flow-E07) tzRevelarResto se recalcula con las iniciales guardadas: si alguna es 'true' el resto aparece de una vez.
  Rama: tzWasPartner() = la etiqueta elegida en AG_3 contiene "pareja" → 24 preguntas de pareja (TZP_1…24); si no, 24 de no pareja (TZNP_1…24)

PASO 5 — "Calcular y guardar Tamizaje" (guardarTamizaje)
  wasPartner = tzWasPartner()
  esOsigd    = (DP_9 respondida y su etiqueta no contiene "heterosexual") o (la etiqueta de DP_8 contiene "trans" o "no binaria")
  subtipoFeminicidio / subtipoHomicidioAbyecto = HV_10 incluye "tentativa feminicidio" / "tentativa homicidio"
  tzAnswers (REG-34):
    SI tzRevelarResto:   TZ_1…TZ_8 + las 24 de la rama activa, cada una 'true' o 'false'
    SI NO:               solo TZ_1, TZ_2, TZ_3, TZ_4 y TZ_7 ('true'/'false')   // el resto no se preguntó
  tzGuardando = true; errorGuardado = ''
  [POST] /api/v1/victim-cases-v3/{icode}/tamizaje
  { "answers": tzAnswers, "wasPartner": wasPartner, "esOsigd": esOsigd, "subtipoFeminicidio": …, "subtipoHomicidioAbyecto": … }

  Backend (SaveTamizaje):
    riskType = esOsigd ? "homicidio_abyecto" : "feminicidio"
    SEGÚN subtipos de tentativa:
      ambos      → riskType = "feminicidio" y tieFlag = true      // G-V3-02 resuelto: prioridad a Feminicidio
      solo fem.  → "feminicidio"
      solo hom.  → "homicidio_abyecto"
    score = número de respuestas "true" en `answers` (7 generales + VPP + 24 de la rama; con solo las 5 iniciales, solo esas)
    nivel = primer umbral con score ≤ máximo:
      feminicidio + pareja:      ≤4 bajo · ≤8 moderado · ≤15 alto · más: extremo
      feminicidio + no pareja:   ≤2 bajo · ≤5 moderado · ≤8 alto  · más: extremo
      homicidio abyecto:         ≤5 bajo · ≤11 moderado · ≤16 alto · más: extremo   (misma escala para pareja y no pareja)
    Persiste (SaveTamizajeBranch, UNA transacción — REG-12):
      las preguntas generales presentes (TZ_1…TZ_8 → victim_case_form2_tzv3_g1…g7 y _vpp; antes NO se guardaban y al reabrir todas salían «No»)
      + las columnas de la rama activa + victim_case_form2_tzv3_risk_type / _risk_level / _risk_score / _risk_tie_flag
      SI el guardado no trae nada del resto (TZ_5, TZ_6, TZ_8, rama): esas columnas quedan en «No» —también si antes había respuestas— y NO se marcan como respondidas (REG-34)
      markKeysAnswered al final, dentro de la transacción: solo cuenta como respondido lo que quedó guardado
    Respuesta: { riskType, riskLevel, score, tieFlag }   (rangos "pendientes de ajuste final" según el cliente: no recalibrar)
               un Caso que no existe responde error (500), no 200

  Cliente: SI la respuesta es ok: tzResultado = respuesta
           SI NO (500, sesión vencida, sin conexión): aviso rojo «No se pudo guardar el Tamizaje…» y lo respondido sigue en pantalla (flow-E21)
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
| ~~**Tamizaje sin forma de guardarse si las 5 iniciales son "No":** el botón estaba dentro del bloque `v-show="tzRevelarResto"`…~~ — **resuelto 2026-11-02 (WP-2 · REG-34), verificado en el navegador y en el API:** el botón aparece cuando las 5 iniciales están respondidas, solo se envían esas 5 y `Finalize` acepta el Tamizaje | PASO 4, PASO 5 |
| `guardarYSalir` sigue sin guardar el Tamizaje (se guarda con su propio botón): salir con respuestas del Tamizaje sin calcular las pierde | PASO 5 |
| Planeado y no implementado en el cliente: barra de avance por módulo; deshabilitar opciones según identidad de género (Gay/Lesbiana, sexo asignado al nacer); limpiar las preguntas de escala; regla de gestación por sexo/identidad (DP_26 siempre visible); limpiar subtipos al desmarcar un tipo; modal de motivo al editar el Perfil de un Registro Completo; modal "Registro Completo"; bitácora de trazabilidad del Perfil (`{ts, agente, campo, antes, después, motivo}`) | PASO 2, PASO 3 |
| Planeado y no implementado: cálculo del Tamizaje en cada cambio de DP/HV/TZ, ocultar preguntas OSIGD y de ámbito laboral (las constantes `TAMIZAJE_*_OSIGD_V3` y `*_LABORAL_V3` no se usan: se muestran las 24 de la rama y TZ_8 siempre), y alertas a coordinación / equipo psicosocial (`NotificationService`): hoy solo hay avisos en pantalla | PASO 4, PASO 5 |
| Planeado y no implementado en GS: advertencia al marcar Medidas de Emergencia y Atención Psicosocial a la vez, y presugerencia de equipo según el nivel de riesgo o la ideación | PASO 6 |
| El fan-out del Perfil (CI/DP/CA) no ocurre en este evento sino en el backend al guardar la sección (flow-E12 paso 5) | — |
| `esOsigd` cuenta como OSIGD cualquier orientación sexual respondida distinta de "heterosexual" (incluye "No informa" u "Otro") y compara etiquetas por texto, no por code | PASO 5 |
