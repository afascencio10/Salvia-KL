━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando actualiza una respuesta del formulario
   Tipo: User Interaction
   Función: handleChangeValor(campoId, valor) → dinamic-form.js (answers-updated)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  campoId: id de la pregunta modificada   → FieldRenderer que disparó el cambio
  valor:   nuevo valor (string | boolean | string[] | número)   → input del agente
}

PASO 1 — Guardar el valor en el estado local
  respuestas[campoId] = valor

PASO 2 — Recalcular condicionales de visibilidad
  Para cada pregunta con `visibility_condition` que depende de campoId:
    Evaluar EQUALS/CONTAINS contra el nuevo valor (motor existente de dinamic-form, sin cambios)
  → Angular/Vue re-evalúa los `[v-if]` del árbol → preguntas se muestran/ocultan según corresponda
  (ver SUB-FLUJO: Recalcular visibilidad específica de Registro V3)

┌─────────────────────────────────────────┐
│  SUB-FLUJO: Recalcular visibilidad específica de Registro V3  │
└─────────────────────────────────────────┘

  S0. SI campoId === 'pronombres':
        → Recalcula `tratamiento`: Ella → "Señora", Él → "Señor", Elle → "Señore"
        → Se refleja de inmediato en NombreIdentitarioPrevalente (header) y en el párrafo síntesis § A (ver S10)
  S0b. SI campoId ∈ {'nombres', 'nombre_identitario', 'contacto_nombres', 'nombres_agresor'} (campos de nombre):
        → Validar "solo letras": si `valor` coincide con `/^\d+$/` (solo dígitos) → marcar el campo como inválido, mostrar error inline "Este campo no admite solo números" (obs. #13 del cliente)
        → No bloquea el resto del flujo — solo impide que ese campo se cuente como respondido para la barra de avance (PASO 3) hasta que se corrija
  S1. SI campoId === 'tipo_documento':
        → Auto-sugerir Nacionalidad: si tipo_documento ∈ {colombianos} → respuestas.nacionalidad = 'Colombiana' (editable después); si no → 'Extranjero(a)'
        → Si Nacionalidad = Colombiana(o): oculta/inhabilita "condición migratoria"
        → Aplica la máscara de validación de "número de documento" correspondiente (ver registro-v3-secciones.md § CI)
  S2. SI campoId === 'nacionalidad' Y valor !== 'Extranjero(a)':
        → Oculta "¿De qué país es su nacionalidad?" y "condición migratoria"
  S3. SI campoId === 'identidad_genero':
        → Inhabilita "Gay" en orientación sexual si valor ∈ {Mujer, Mujer Trans}; inhabilita "Lesbiana" si valor ∈ {Hombre, Hombre Trans}
        → Inhabilita "Masculino" en sexo asignado al nacer si valor === 'No binaria-asignado femenino al nacer'; inhabilita "Femenino" si valor === 'No binaria-asignado masculino al nacer'
        → Recalcula `esOsigd` (ver S7)
  S4. SI campoId === 'orientacion_sexual' O campoId === 'sexo_asignado_nacer':
        → Recalcula `esOsigd` (ver S7)
  S5. SI campoId === 'tiene_discapacidad' (Q2 de DP) Y valor === false:
        → Limpia y oculta las 9 preguntas de escala de dificultad
  S5b. SI campoId === 'se_reconoce_discapacidad' (CI.Q7) Y valor === false:
        → Limpia y oculta "Ley 1996" (DP.Q4) — **corregido por obs. #7 del cliente**: Ley 1996 depende de CI.Q7
          ("¿se reconoce con discapacidad?"), no de DP.Q2 ("dificultades físicas/mentales/sensoriales") como se
          había asumido de la redacción del Excel
  S6. SI campoId === 'tipo_violencia':
        → Recalcula las opciones disponibles de "subtipo_violencia" según la tabla de subtipos por tipo (registro-v3-secciones.md § HV)
        → Si se desmarca un tipo, se eliminan de la respuesta los subtipos que ya no aplican
  S7. Recalcular `esOsigd`:
        esOsigd = identidad_genero ∈ {Mujer Trans, Hombre Trans, No binaria-*, Otro}
                  O orientacion_sexual ∈ {Gay, Lesbiana, Bisexual, Otro, No estoy segura/o/e}
                  O sexo_asignado_nacer === 'Intersexual o indeterminado'
        → Si esOsigd cambia de false a true: se revela la pregunta fuerte VPP en TZ (si el módulo ya está desbloqueado)
        → Si cambia de true a false: se oculta y se limpia esa respuesta
  S8. SI campoId === 'relacion_agresor' (AG):
        → wasPartner = valor ∈ {'Pareja íntima', 'Expareja o antigua pareja'}
        → Si wasPartner cambia: cambia la rama de Tamizaje visible (pareja ↔ no pareja) — las respuestas de la rama anterior se conservan en memoria pero no se envían al guardar
  S9. SI campoId pertenece al módulo DP y es "¿Cuál es tu condición migratoria?" (Q7) O cualquier campo con el flag `soloSiAmbitoLaboral`:
        → Visible solo si 'Laboral' ∈ ambito_violencia (ver también sub-flujo Tamizaje más abajo para las preguntas de jerarquía/autoridad)
  S9b. Recalcular visibilidad de "¿Actualmente en embarazo?" (DP.Q26) — **condición corregida por obs. #8 del cliente**:
        visible = (sexo_asignado_nacer === 'Femenino') O (identidad_genero ∈ {'Mujer', 'Hombre Trans', 'No binaria-asignado femenino al nacer'})
        → Se re-evalúa cada vez que cambia `sexo_asignado_nacer` o `identidad_genero`, no solo una de las dos
  → FIN SUB-FLUJO → CONTINÚA FLUJO GENERAL

PASO 3 — Recalcular barra de avance del módulo del campo modificado
  progresoModulo(moduloDelCampo) = (preguntas obligatorias respondidas / preguntas obligatorias visibles) × 100
  → Vue re-evalúa la BarraProgreso del ModuloItem correspondiente en el sidebar

PASO 4 — Recalcular habilitación secuencial

┌─────────────────────────────────────────┐
│  SUB-FLUJO: Recalcular habilitación secuencial (triadaCompleta)  │
└─────────────────────────────────────────┘

  S1. triadaCompleta = todas las obligatorias **[1er contacto]** (no se exige [máx. 2º contacto]) de CI, DP, CA, HV y AG están respondidas
      — unificado con el umbral de Borrador→Iniciado de `flow-E12` (antes este sub-flujo exigía también [máx. 2º contacto], inconsistente con esa transición — ver registro-v3-data.md)
      — CA no aporta preguntas a este cálculo: sus 4 preguntas son todas [máx. 2º contacto], puede quedar vacío sin bloquear la tríada
  S2. SI triadaCompleta cambia de false a true:
        → Se desbloquean TZ y GS en el sidebar (se quita el 🔒)
  S3. SI triadaCompleta cambia de true a false (el agente borró una respuesta obligatoria ya diligenciada):
        → Se vuelven a bloquear TZ y GS
        → SI moduloActual ∈ ['TZ', 'GS']: → moduloActual = último módulo de la tríada, navega automáticamente fuera de TZ/GS
  → FIN SUB-FLUJO → CONTINÚA FLUJO GENERAL

PASO 5 — Verificar duplicidad (solo si campoId === 'numero_documento')

SI campoId === 'numero_documento' Y valor tiene longitud válida para el tipo de documento actual:
  → SUB-FLUJO: Verificar duplicidad
    S1. VictimCaseService.buscarPorDocumento({ tipoDocumento, numeroDocumento: valor })
    S2. SI existe otro `victim_case` con ese documento Y su icode !== caso.icode actual:
          → Mostrar alerta emergente "Ya existe un Registro con este documento"
          → Ofrecer botón "Ir al Expediente de víctima" (deshabilitado/placeholder en esta ronda — Expediente diferido, ver registro-v3-data.md D2) o "Continuar igual" (el agente decide)
          → TERMINAR sub-flujo, no bloquea el guardado
    S3. SI NO existe:
          → CONTINÚA sin alerta
  → FIN SUB-FLUJO → CONTINÚA FLUJO GENERAL

SI NO (campoId !== 'numero_documento'):
  → CONTINÚA FLUJO GENERAL

PASO 5b — Exigir motivo si se edita un campo de Perfil ya respondido (agregado del mockup)

SI campoId pertenece a CI, DP o CA Y ya existe un valor previo para ese campo (no es la primera vez que se responde) Y el Registro ya está en estado Completo:
  → Abrir modal "Editar Perfil de la víctima" — el agente elige un motivo:
    'Solicitud de la víctima' | 'Pertinencia de actualización identificada por el agente' | 'Corrección de digitación'
  → Sin motivo seleccionado, el cambio no se aplica — TERMINAR ejecución
  → Con motivo seleccionado: respuestas[campoId] = valor, CONTINÚA FLUJO GENERAL, y el motivo viaja junto con el campo al PASO 6 (fan-out) y a la trazabilidad

SI NO (primera respuesta a ese campo, o Registro aún Incompleto):
  → No se exige motivo, se guarda directo (ya cubierto por PASO 1)
  → CONTINÚA FLUJO GENERAL

PASO 6 — Sincronizar Perfil con otros Casos del mismo documento (fan-out)

SI campoId pertenece a los módulos CI, DP o CA:
  → SUB-FLUJO: Sincronización fan-out del Perfil
    S1. VictimCaseService.listarPorDocumento({ tipoDocumento, numeroDocumento })
    S2. SI hay más de un `victim_case` con el mismo documento (existen otros Casos del mismo Registro):
          PARA CADA otroCaso EN esosVictimCases EXCEPTO el actual:
            [PUT] /api/v1/victim-cases/{otroCaso.icode}/sincronizar-perfil
            // payload:
            {
              "campo":  campoId,                  // qué campo cambió
              "valor":  valor,                     // nuevo valor
              "motivo": motivoEdicion || null       // del PASO 5b, si aplicó
            }
          → Cada otroCaso queda con el mismo valor para ese campo de Perfil
    S3. SI NO hay otros Casos:
          → No hace nada
    S4. Registrar trazabilidad (todo campo de Perfil editado, tenga o no otros Casos — formato del mockup):
          log.push({
            ts:      now(),
            agente:  agentId,
            entidad: 'Perfil',          // o 'Hecho'/'Agresor' según el módulo, fuera de Perfil
            campo:   campoId,
            antes:   valorAnterior,
            despues: valor,
            motivo:  motivoEdicion || null
          })
          → Visible luego en el log de Trazabilidad del Expediente de víctima (diferido, D2) y en el inspector interno si existe
  → FIN SUB-FLUJO → CONTINÚA FLUJO GENERAL

SI NO (campoId pertenece a HV, AG, TZ o GS):
  → CONTINÚA FLUJO GENERAL (esos campos son del Caso, no del Perfil — no se sincronizan)

PASO 6b — Recalcular estado del Registro (Incompleto → Completo)

SI campoId pertenece a CI, DP o CA:
  → registro.estado = (todas las obligatorias [1er contacto] y [máx. 2º contacto] de CI+DP+CA respondidas) ? 'Completo' : 'Incompleto'
  SI registro.estado cambia de 'Incompleto' a 'Completo' en esta ejecución:
    → Mostrar modal de confirmación "Registro de víctima Completo" (agregado del mockup): muestra el `idRegistro` generado,
      aclara que de aquí en adelante toda modificación a estos campos exige motivo (ver PASO 5b) y queda en el log,
      y que el Registro no se cierra — permanece disponible en el futuro Expediente de víctima (diferido, D2) agrupando todos los Casos
  → CONTINÚA FLUJO GENERAL

PASO 7 — Recalcular resultado de Tamizaje (si el módulo del campo es DP, HV o TZ — el tipo ahora depende de DP/HV, no solo de TZ)

SI moduloDelCampo ∈ {'DP', 'HV', 'TZ'}:
  → Ver SUB-FLUJO: Calcular tipo y nivel de riesgo (abajo)

┌─────────────────────────────────────────┐
│  SUB-FLUJO: Calcular tipo y nivel de riesgo  │
└─────────────────────────────────────────┘

  > Corregido tras releer `Observaciones y ajustes al MockUp de registro
  > 2.0.pdf` (obs. #5 y #15) — este PDF tiene prioridad sobre el Guion
  > Funcional de Drive. Ver `registro-v3-data.md` § D3/D4 para el porqué.

  S1. Determinar tipo de riesgo por defecto — **directo de identidad/orientación (D3), ya no de la pregunta VPP**:
      SI esOsigd (calculado en PASO 2, S7, a partir de identidad_genero/orientacion_sexual/sexo_asignado_nacer de DP):
        → tipoRiesgoDefault = 'Homicidio abyecto o fútil'
      SI NO:
        → tipoRiesgoDefault = 'Feminicidio'
      (la pregunta VPP de TZ ya no se consulta para este cálculo — solo suma punto en S4)

  S2. Aplicar override por subtipo de violencia (D3 — ver registro-v3-data.md):
      SI 'Tentativa Feminicidio' ∈ respuestas.subtipo_violencia Y 'Tentativa Homicidio abyecto o fútil' ∈ respuestas.subtipo_violencia:
        → GAP G-V3-02: ambos marcados — tipoRiesgo = (tipoRiesgoDefault si coincide con alguno de los dos) o 'Feminicidio' por defecto
        → Marcar caso.bandera_revision_tipo_riesgo = true
      SI SOLO 'Tentativa Feminicidio' ∈ respuestas.subtipo_violencia:
        → tipoRiesgo = 'Feminicidio'  (override, ignora tipoRiesgoDefault)
      SI SOLO 'Tentativa Homicidio abyecto o fútil' ∈ respuestas.subtipo_violencia:
        → tipoRiesgo = 'Homicidio abyecto o fútil'  (override)
      SI NINGUNO:
        → tipoRiesgo = tipoRiesgoDefault

  S3. Evaluar el "gate" de las 5 preguntas iniciales (D4 — obs. #15 del cliente, propuesta GAP G-V3-08):
      lasCinco = [pregunta 1, 2, 3, 4, 7 de las generales]  // propuesta, ver GAP G-V3-08
      SI TODAS las lasCinco están respondidas Y TODAS === false (No):
        → tamizajeAbierto = false
        → Ocultar/no mostrar: preguntas 5-6 generales, pregunta VPP, rama pareja/no-pareja completa
        → score = 0, nivelRiesgo = 'Bajo' (el más bajo de la escala de tipoRiesgo)
        → IR A S5 (persistir) — se salta S4
      SI alguna de lasCinco === true (Sí), o aún no están todas respondidas:
        → tamizajeAbierto = true
        → Se revela/mantiene visible el resto (preguntas 5-6, VPP si esOsigd, rama pareja/no-pareja)
        → CONTINÚA a S4

  S4. Calcular el score (solo si tamizajeAbierto):
      score = cuenta de respuestas "Sí" entre: las 7 preguntas generales + la pregunta VPP (si esOsigd, ahora solo aporta puntaje) + todas las preguntas visibles de la rama activa (pareja si wasPartner, no pareja si no)

  S4b. Determinar nivel de riesgo según tipoRiesgo (solo si tamizajeAbierto):
      SEGÚN tipoRiesgo:
        CASO 'Feminicidio' Y wasPartner:        → umbrales 0-4 Bajo / 5-8 Moderado / 9-15 Alto / 16-24 Extremo
        CASO 'Feminicidio' Y !wasPartner:       → umbrales 0-2 Bajo / 3-5 Moderado / 6-8 Alto / 9-20 Extremo
        CASO 'Homicidio abyecto o fútil' (pareja o no pareja, comparten escala): → umbrales 0-5 Bajo / 6-11 Moderado / 12-16 Alto / 17-28 Extremo
      → nivelRiesgo = umbral donde cae `score`
      ⚠️ Estos rangos siguen marcados "pendientes de ajustes finales" por el propio cliente (obs. #3) — no recalibrar sin instrucción nueva

  S5. Persistir resultado (no solo en memoria):
      victim_case_form2.tamizaje_tipo_riesgo = tipoRiesgo
      victim_case_form2.tamizaje_nivel_riesgo = nivelRiesgo
      → BannerResultadoTamizaje muestra "{tipoRiesgo}: nivel {nivelRiesgo}"

  S6. SI nivelRiesgo ∈ {Alto, Extremo}:
        → SUB-FLUJO: Alerta por riesgo (ver abajo) — se dispara inmediatamente, sin esperar a guardar la sección
      SI nivelRiesgo === 'Moderado':
        → Sugerir "activación de rutas pertinentes" en GS (no es una alerta a coordinación, es una sugerencia de acción — ver PASO 8)

  S7. SI campoId === pregunta de ideación suicida (pareja ítem 24 / no pareja ítem 24) Y valor === true:
        → SUB-FLUJO: Alerta por ideación suicida (ver abajo)

  → FIN SUB-FLUJO → CONTINÚA FLUJO GENERAL

┌─────────────────────────────────────────┐
│  SUB-FLUJO: Alerta por riesgo  │
└─────────────────────────────────────────┘
  S1. NotificationService.alertarCoordinacion({
        icode: caso.icode,
        motivo: 'tamizaje_alto_extremo',
        nivelRiesgo
      })
  S2. Mostrar banner inmediato al agente: "Alerta enviada a coordinación — nivel {nivelRiesgo}"
  → No espera a que el formulario esté completo (regla explícita del Guion Funcional)
  → FIN SUB-FLUJO

┌─────────────────────────────────────────┐
│  SUB-FLUJO: Alerta por ideación suicida  │
└─────────────────────────────────────────┘
  S1. Mostrar al agente un letrero: "Verifique si se trató de ideación o de intento"
  S2. SI el agente marca "Ideación": → sugerirá Psicosocial en GS (ver PASO 8)
      SI el agente marca "Intento": → sugerirá Medidas de Emergencia en GS Y activa de inmediato la ruta de salud mental
      NotificationService.alertarEquipoPsicosocial({ icode: caso.icode, tipo: 'ideacion' | 'intento' })
  → FIN SUB-FLUJO

PASO 8 — Reglas de Gestión Salvia (solo si el módulo del campo es GS)

> Reestructurado sobre el modelo real de Seguimiento (D9) — reemplaza la
> versión anterior de 4 "procesos" genéricos y el campo `medio_gestion`
> a nivel de GS (el canal de activación ahora vive por entidad, dentro del
> repeater de Registrar Entidades — ver SUB-FLUJO más abajo).

SI moduloDelCampo === 'GS':

  SI campoId === 'equipos_remision' (GS pregunta 6):
    → Revela los bloques de criterios de cada equipo marcado (medidas de emergencia, psicosocial, hombres, discapacidad, estabilización)
    → Oculta y limpia los bloques de los equipos que se desmarcaron
    → SI valor incluye 'medidas_emergencia' Y valor incluye 'atencion_psico':
        → Mostrar BannerAdvertencia "No se puede remitir a estos dos equipos simultáneamente" (no bloquea, solo advierte — igual que en Seguimiento real)
    → CONTINÚA FLUJO GENERAL

  SI campoId pertenece a los criterios de Atención Psicosocial (GS pregunta 9):
    → Recalcular puntaje: suma 3 (conducta suicida) + 2 (interseccionalidad) + 1 por cada uno de los 4 restantes marcados
    → SI el criterio obligatorio NO está marcado O el puntaje < 3:
        → Mostrar advertencia inline "Para que la remisión sea válida se debe marcar el criterio obligatorio y sumar mínimo 3 puntos" (informativo, igual que el banner real — no bloquea el guardado)
    → CONTINÚA FLUJO GENERAL

  SI campoId === 'requiere_remision' (GS pregunta 5) Y valor === true:
    → Se revela "¿Cuáles equipos?"
    → Presugerir equipo: si tamizaje_nivel_riesgo ∈ {Alto, Extremo} → sugerir 'medidas_emergencia'; si hubo ideación suicida → sugerir 'atencion_psico' (sugerencia, no obliga, el agente puede cambiarla)
    → CONTINÚA FLUJO GENERAL

┌─────────────────────────────────────────┐
│  SUB-FLUJO: Registrar/editar una Entidad (repeater)  │
└─────────────────────────────────────────┘
  S1. SI campoId === 'reported_to_entity' (GS pregunta 14) Y valor === true:
        → Revela el repeater "Entidades Acudidas"
  S2. SI campoId === 'salvia_activates_routes' (GS pregunta 16) Y valor === true:
        → Revela el repeater "Entidades Activación"
  S3. SI campoId pertenece a un ítem del repeater Y es el campo "Sector de la entidad":
        → Recalcula las opciones del campo "Entidad" (cascada dependiente del sector, catálogo `entidadesDisponibles` — mecanismo dinámico, no enum estático, ver GAP G-V3-21)
  → FIN SUB-FLUJO → CONTINÚA FLUJO GENERAL

┌─────────────────────────────────────────┐
│  SUB-FLUJO: Registrar/editar una Barrera (repeater)  │
└─────────────────────────────────────────┘
  S1. SI campoId === 'identifies_barriers' (GS pregunta 17) Y valor === true:
        → Revela el repeater de Barreras
  S2. SI campoId pertenece a un ítem del repeater Y es el campo "Sector de la barrera":
        → Revela el bloque de multiselect de barreras correspondiente a ese sector (Salud/Justicia/Protección) y oculta los otros 2
  → FIN SUB-FLUJO → CONTINÚA FLUJO GENERAL

SI campoId === 'wants_imminent_alert' (GS pregunta 18, botón):
  → Ver flujo: Cuando presiona "Alerta inminente" (mismo flujo que el botón del header, origen='GS')

PASO 8b — Recalcular el párrafo síntesis (si el campo pertenece a CI, DP, HV o AG)

SI moduloDelCampo ∈ {'CI', 'DP', 'HV', 'AG'}:
  → Recalcular las 4 secciones (A-D) del párrafo síntesis aplicando las reglas de
    `registro-v3-parrafo-sintesis.md` — cada sección se omite si no tiene datos suficientes
  → Se ve reflejado de inmediato en el módulo GS, aunque el agente no esté viéndolo en ese momento
  → No se persiste — se recalcula on-demand, igual que el resultado de Tamizaje
  → CONTINÚA FLUJO GENERAL

PASO 9 — Re-render final
  Vue re-evalúa todos los bindings afectados por los pasos anteriores

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                                                     | Paso afectado |
|-----------------------------------------------------------------------------------------------|---------------|
| Catálogo de criterios de remisión por proceso (G-V3-04) — bloquea completar el detalle de PASO 8 | PASO 8 |
| Confirmar los códigos EQUALS reales de `visibility_condition` contra `victim_case_form2_enums` antes de implementar (G-V3-06, mismo gap heredado de V2-real) | PASO 2, S1-S9 |
| Control de concurrencia del fan-out (G-V3-01) si dos agentes editan Perfil de Casos distintos del mismo documento simultáneamente | PASO 6 |
| Umbral exacto de "longitud válida" antes de disparar la verificación de duplicidad en PASO 5 (¿validar formato completo primero, o cualquier longitud mínima?) | PASO 5 |
| G-V3-08 — cuáles son exactamente las 5 preguntas obligatorias iniciales del Tamizaje (usamos una propuesta no confirmada: generales 1,2,3,4,7) | SUB-FLUJO Calcular tipo y nivel de riesgo, S3 |
| G-V3-02 — sin regla de prioridad si se marcan ambos subtipos de tentativa a la vez | SUB-FLUJO Calcular tipo y nivel de riesgo, S2 |
