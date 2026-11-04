━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando guarda una sección/módulo
   Tipo: User Interaction
   Estado: 🟡 Implementado con cambios — sincronizado con el código 2026-10-31 · actualizado 2026-11-02 (WP-2 · REG-11) · actualizado 2026-11-04 (WP-4)
   Funciones: guardarModulo(moduloId) · enviarSeccion(moduloId) · mensajeErrorGuardado(res) · mostrarErrorGuardado(e) → registro-v3-app.js · SaveSection() → victim_case_v3_controller.go / service / repository (validación: victim_case_v3_validation.go)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

> 🔄 **WP-4 (2026-11-04):** (1) El módulo nuevo **RA** tiene su propio «Guardar sección» (`guardarRutaAtencion`): valida `GS_14` y
> `GS_17` y envía el módulo **GS** con `enviarSeccion('GS')` (esas llaves se guardan con GS); no avanza solo (la barra «Siguiente» sí
> funciona). (2) Si el guardado completa la tríada, el servidor pasa el Caso de `bo` a `in` y **materializa el borrador** de entidades y
> barreras antes del cambio de estado ([sub-flujo](flow-subflujo-materializar-borrador.md)); el cliente, al ver `triadaCompleta` por primera
> vez, recarga entidades, barreras y `entitiesReadOnly` y la sección queda de solo lectura. (3) `GS_20` se guarda con la sección GS (relación
> de enums). (4) Si el vaciado falla, el guardado responde 500 y el Caso sigue en Borrador.


INPUT: {
  moduloId:  módulo visible (CI, DP, CA, HV, AG o GS)   → botón "Guardar sección"
  answers:   mapa FieldKey → valor de todo el formulario   → estado local; se envían solo las keys que empiezan con "{moduloId}_"
  icode:     identificador del Caso                         → estado local
}
(TZ no usa este evento: se guarda con "Calcular y guardar Tamizaje", ver flow-E11.)

Validar formulario (cliente, en este orden; cada falla muestra `mensajeGuardado` y TERMINA):
  • Obligatorias: preguntas VISIBLES del módulo con req 'C1' ([1er contacto]) sin respuesta
      → "Faltan obligatorias de 1er contacto: <etiquetas separadas por ;>"
      (una respuesta "false" cuenta como respondida; solo '' y sin definir cuentan como vacías)
      En GS solo se revisan las de § A (GS_1 y GS_4; `GS_2` se retiró en WP-4 · REG-14); las preguntas dibujadas directamente en el HTML (GS_5…GS_19) no se validan aquí
  • "Solo letras" (CI_2, CI_3, CA_1, AG_7): valor compuesto solo por dígitos
      → "Estos campos no pueden ser solo números: …"
  • CI_6 fuera de la máscara del tipo de documento (solo al guardar CI)
      → "Número de documento inválido para el tipo seleccionado: <formato esperado>"

PASO 1 — Persistir el módulo
  guardando = true
  [POST] /api/v1/victim-cases-v3/{icode}/seccion/{moduloId}
  // payload:
  { "answers": { "<FieldKey>": "<valor>", ... } }   // valores como texto: icode de enum, "true"/"false", CSV de icodes en multi-selección, fecha ISO

  Backend (SaveSection) — UNA transacción: o se guarda todo el módulo o no se guarda nada (REG-11):
    0. Validar y normalizar ANTES de escribir (resolveFieldUpdates y los campos núcleo):
         fecha → "AAAA-MM-DD" (acepta también el formato crudo anterior "2026-10-30 00:00:00 +0000 UTC") · hora → "HH:MM" o "HH:MM:SS"
         número → solo dígitos (los espacios se quitan: "300 123 4567" → 3001234567) · texto → máximo el ancho real de la columna
         (HV_7 150, CI_6 32, nombres 32…; el ancho se lee de information_schema) · enum → el icode debe existir
         valor vacío de fecha/hora/lista → no se escribe (la columna conserva su valor; NULL solo si la columna admite NULL, p. ej. DP_51)
       SI algún dato no es válido: ROLLBACK (ni valores ni marcas) → 400 { "error": "valor no válido", "campo": "<FieldKey>", "detalle": "<motivo>" } → TERMINAR
    1. UPDATE parcial de victim_case_form2 solo con las columnas de las keys recibidas (bool → id del "Sí"/"No"; enum → id por icode; fecha/hora/número normalizados)
    2. Relaciones multi-valor (KindEnumN): borra y reinserta las del campo, acotado a sus categorías de enum
    3. Campos núcleo en victim_case: CI_2 → names/last_names (se parte en el primer espacio) · CI_5 → doc_type (code corto) · CI_6 → doc_number · DP_28 → town_code
    4. markKeysAnswered (al final, dentro de la transacción): acumula en victim_case.victim_case_v3_answered_keys cada key recibida (aunque el valor sea vacío)
    COMMIT
    5. SI módulo ∈ CI/DP/CA: sincroniza el mismo cambio en los demás `victim_case` con el mismo tipo+número de documento (fan-out DESPUÉS del commit, best-effort: un error solo se loguea)
    6. SI módulo ∈ CI/DP/CA/HV/AG: triadaCompleta = ¿todas las 32 obligatorias de la tríada respondidas alguna vez?
         SI sí: UPDATE victim_case SET status='in' WHERE status='bo'   // nunca revierte un Caso ya Terminado/Cerrado
       SI módulo es TZ o GS: responde triadaCompleta=false sin evaluar nada
  Respuesta: { "ok": true, "triadaCompleta": <bool> }
             400 { "error": "valor no válido", "campo", "detalle" }   (dato inválido; el log del servidor deja «[RegistroV3] SaveSection rechazado …»)
             500 { "error": "error al guardar la sección" }          (cualquier otra falla; el log deja «[RegistroV3] SaveSection falló …»)

PASO 2 — Actualizar el estado del cliente
  SI la respuesta es ok:
    triadaCompleta = !!data.triadaCompleta
    mensajeGuardado = "Sección guardada."; errorGuardado = ''
  SI NO (400, 401/403, 500 o sin conexión): se lanza un error con el texto para la persona → mostrarErrorGuardado()  (ver flow-E21):
    errorGuardado = aviso rojo con la pregunta a corregir y el motivo (400) o el error del servidor / de sesión / de conexión
    mensajeGuardado = "No se guardó."   // nunca "Sección guardada." si no se guardó
  guardando = false

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━CAMBIOS WP-3 (2026-11-02):
  Validar formulario (REG-05): cada falla (obligatorias, "solo letras", máscara de CI_6) deja mensajeGuardadoEsError = true: el mensaje sale en rojo
    (clase v3-msg-error, #b02a37, contraste 6,5:1). El texto no cambia. guardarModulo devuelve una promesa true/false.
  PASO 1 · SaveSection (REG-35): la tríada se evalúa con CUALQUIER módulo (también TZ y GS) y se devuelve el valor real; Borrador a Iniciado sigue
    siendo solo de la tríada (CI/DP/CA/HV/AG). Cliente: triadaCompleta = triadaCompleta || !!data.triadaCompleta (nunca baja de true a false).
  PASO 2 (REG-36, ya cubierto por REG-11/E21): una respuesta no exitosa muestra el aviso rojo de flow-E21 y no dice "Sección guardada.".
  PASO 3 (REG-03/04): tras guardar bien se vuelve a leer el resumen del Caso (GET /victim-cases-v3/:icode): encabezado con nombre, documento y estado.
  PASO 4 (REG-01): si se guardó y el usuario sigue en ese módulo, pasa al siguiente módulo no bloqueado (salvo en el último) con
    mensajeAvance = "✓ Sección guardada. Pasaste de «X» a «Y».". Ver flow-E10 / flow-E22.
  Supuesto (comentario de Felipe, 2026-11-02): por ahora un solo usuario edita el registro a la vez; no hay control de concurrencia.

━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión | Paso afectado |
|---|---|
| Vaciar una pregunta opcional de fecha/hora cuyo campo es NOT NULL (HV_3, DP_21…) y guardar **no** la borra: el valor anterior se conserva, porque la columna no admite NULL (decisión de BD pendiente) | PASO 1 (0) |
| No hay toast ni marca de error por campo: solo el texto de `mensajeGuardado` junto al botón | Validar formulario |
