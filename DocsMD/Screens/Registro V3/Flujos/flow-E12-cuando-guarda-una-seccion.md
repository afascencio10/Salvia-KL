━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando guarda una sección/módulo
   Tipo: User Interaction
   Estado: 🟡 Implementado con cambios — sincronizado con el código 2026-10-31
   Funciones: guardarModulo(moduloId) · enviarSeccion(moduloId) → registro-v3-app.js · SaveSection() → victim_case_v3_controller.go / service / repository
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

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
      En GS solo se revisan las de § A (GS_1, GS_2, GS_4); las preguntas dibujadas directamente en el HTML (GS_5…GS_19) no se validan aquí
  • "Solo letras" (CI_2, CI_3, CA_1, AG_7): valor compuesto solo por dígitos
      → "Estos campos no pueden ser solo números: …"
  • CI_6 fuera de la máscara del tipo de documento (solo al guardar CI)
      → "Número de documento inválido para el tipo seleccionado: <formato esperado>"

PASO 1 — Persistir el módulo
  guardando = true
  [POST] /api/v1/victim-cases-v3/{icode}/seccion/{moduloId}
  // payload:
  { "answers": { "<FieldKey>": "<valor>", ... } }   // valores como texto: icode de enum, "true"/"false", CSV de icodes en multi-selección, fecha ISO

  Backend (SaveSection):
    1. markKeysAnswered: acumula en victim_case.victim_case_v3_answered_keys cada key recibida (aunque el valor sea vacío)
    2. UPDATE parcial de victim_case_form2 solo con las columnas de las keys recibidas (bool → id del "Sí"/"No"; enum → id por icode; fecha/hora/número tal cual)
    3. Relaciones multi-valor (KindEnumN): borra y reinserta las del campo, acotado a sus categorías de enum
    4. Campos núcleo en victim_case: CI_2 → names/last_names (se parte en el primer espacio) · CI_5 → doc_type (code corto) · CI_6 → doc_number · DP_28 → town_code
    5. SI módulo ∈ CI/DP/CA: sincroniza el mismo cambio en los demás `victim_case` con el mismo tipo+número de documento (fan-out, best-effort: un error solo se loguea)
    6. SI módulo ∈ CI/DP/CA/HV/AG: triadaCompleta = ¿todas las 32 obligatorias de la tríada respondidas alguna vez?
         SI sí: UPDATE victim_case SET status='in' WHERE status='bo'   // nunca revierte un Caso ya Terminado/Cerrado
       SI módulo es TZ o GS: responde triadaCompleta=false sin evaluar nada
  Respuesta: { "ok": true, "triadaCompleta": <bool> }   (error de servicio: 500 { "error": "error al guardar la sección" })

PASO 2 — Actualizar el estado del cliente
  triadaCompleta = !!data.triadaCompleta
  mensajeGuardado = "Sección guardada."
  SI la petición lanza una excepción de red: mensajeGuardado = "Error al guardar la sección."
  guardando = false

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión | Paso afectado |
|---|---|
| `enviarSeccion` no revisa `r.ok`: una respuesta 500 se lee como JSON sin `triadaCompleta` y aun así muestra "Sección guardada." (y deja `triadaCompleta` en false). Solo una caída de red muestra el error | PASO 2 |
| Al guardar GS (o TZ) `triadaCompleta` queda en false en el cliente (ver GAP en flow-E10) | PASO 2 |
| No hay toast ni marca de error por campo: solo el texto de `mensajeGuardado` junto al botón | Validar formulario |
| Sin control de concurrencia en el fan-out del Perfil entre agentes que editan el mismo documento (G-V3-01) | PASO 1 (5) |
