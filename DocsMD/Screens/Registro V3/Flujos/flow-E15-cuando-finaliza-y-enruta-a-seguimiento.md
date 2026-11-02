━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona "Finalizar registro y enrutar a Seguimiento"
   Tipo: User Interaction
   Estado: 🟡 Implementado con cambios — sincronizado con el código 2026-10-31
   Funciones: finalizarRegistro() · confirmacionFinalizar() → registro-v3-app.js · Finalize() → victim_case_v3_controller.go / service
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  icode:  identificador del Caso   → estado local
}
(el botón está siempre habilitado, solo en el módulo GS; no lleva cuerpo)

PASO 1 — Finalizar en el backend
  errorFinalizar = ''
  [POST] /api/v1/victim-cases-v3/{icode}/finalizar

  Backend (Finalize):
    answered = keys respondidas alguna vez (victim_case_v3_answered_keys)
    Validar obligatorias — deben estar TODAS en `answered`:
      • Tríada [1er contacto]: CI_1-7, CI_9 · DP_1, 2, 5, 8, 9, 11, 12, 14, 17, 27, 28, 29 · HV_1, 2, 4, 6, 8, 9, 11, 13, 14 · AG_1, 2, 3
      • [máx. 2º contacto] incondicionales: DP_10, DP_15 · CA_1, CA_2, CA_3, CA_4
      • Tamizaje: TZ_1, TZ_2, TZ_3, TZ_4, TZ_7 (las 5 iniciales)
      • Gestión Salvia siempre visibles: GS_1, GS_2, GS_4, GS_5, GS_14, GS_16, GS_17, GS_18
      Las condicionales (DP_3, DP_6, DP_13, DP_19-22, GS_3, GS_6-13, GS_15, GS_19, HV_12, AG_5…) y las [Deseable] no se exigen.
    SI falta alguna:
      → 400 { "error": "No se puede finalizar — hay módulos incompletos", "modulosIncompletos": ["DP", "TZ", …] }   // prefijos de las keys faltantes, ordenados
      → TERMINAR ejecución (el Caso no cambia)
    SI NO:
      denunciaFacilSolicitada = (GS_1 === "true")
      equiposRemisionCSV = GS_6 (CSV de icodes) si GS_5 === "true"; si no, ""
      UPDATE victim_case SET status = 'te' WHERE status = 'in'      // Terminado; si ya estaba 'te' no cambia (idempotente)
      → 200 { "ok": true, "denunciaFacilSolicitada": …, "equiposRemisionCSV": … }
    Otros errores de servicio: 500 { "error": "error al finalizar el registro" }

PASO 2 — Resultado en el cliente
  SI respuesta ok:
    → mensajeCierre = confirmacionFinalizar(data):  "Se solicitó la Denuncia fácil ante la FGN." (si aplica) + "Se confirma remisión a: <equipos>." (nombres resueltos con GS_6)
    → etapa = 'cierre'  (ver interface: 4 tarjetas, párrafo síntesis, tabla de actuaciones, bitácora, "Volver al Cabezote")
  SI NO:
    → errorFinalizar = <error> + " — falta diligenciar: <etiquetas de los módulos faltantes>" (si hay modulosIncompletos)
    → el agente permanece en GS

> No hay traspaso a Seguimiento que ejecutar: Registro V3 escribe directo sobre `victim_case`/`victim_case_form2`, las
> mismas tablas que lee Seguimiento; `te` es la señal. La transición a Cerrado (`cd`) es de Seguimiento.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━CAMBIOS WP-3 (2026-11-02 · REG-37, REG-21):
  PASO 0: si ya hay un guardado o cierre en curso (guardando/finalizando) se ignora el clic (sin doble envío).
  PASO 1: guardarModulo('GS', sinAvanzar) — mismas validaciones de "Guardar sección" (obligatorias visibles del § A de GS) y POST /seccion/GS. Si no se
    guarda, NO se llama a /finalizar y el recordatorio se repite junto al botón (errorFinalizar).
  PASO 2: POST /finalizar. Backend: el 400 incluye camposFaltantes (llaves sin responder, ordenadas) además de modulosIncompletos.
  PASO 3 (respuesta): 200 → Cierre · 400 → "<error> — falta diligenciar: <módulos>. Gestión Salvia, falta: <etiquetas>" (catálogo o mapa local para
    GS_5, 14, 16, 17 y 18) · 401/403 → "tu sesión venció o no tienes permiso. Inicia sesión de nuevo en otra pestaña y vuelve a pulsar «Finalizar
    registro»: lo que diligenciaste sigue en esta pantalla" (la respuesta se lee sin lanzar: un 401 puede venir sin cuerpo) · otro → aviso de error del servidor.

━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión | Paso afectado |
|---|---|
| No hay integración con la FGN ni un workflow de remisión: la Denuncia fácil y las remisiones solo se **confirman en pantalla** a partir de GS_1 y GS_6 (ver registro-v3-pendientes.md C2/C3). Tampoco se crean entradas en un módulo de Seguimiento, ni se notifica a los equipos | PASO 1 |
| `POST /entidades/registrar-lote` y `/barreras/registrar-lote` de la planeación no existen: entidades y barreras se guardan una a una al agregarlas (flow-E17, flow-E18) | — |
| ~~Con las 5 iniciales del Tamizaje en "No" el Tamizaje no llega a guardarse y el Caso no puede finalizar~~ — **resuelto 2026-11-02 (WP-2 · REG-34):** con las 5 iniciales respondidas el botón «Calcular y guardar Tamizaje» aparece, se guardan TZ_1, 2, 3, 4 y 7 (el resto queda en «No», sin marcarse) y `Finalize` ya acepta el Tamizaje (ver flow-E11) | PASO 1 |
| Solo se exige `Terminado` desde `in`: si el Caso estuviera en `bo` (nunca ocurre con la tríada completa) respondería ok sin cambiar el estado | PASO 1 |
