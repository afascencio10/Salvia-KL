━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando guarda la sección «Ruta de atención integral»
   Tipo: User Interaction
   Estado: ✅ Implementado 2026-11-04 (WP-4 · REG-17)
   Funciones: guardarRutaAtencion() · enviarSeccion('GS') · obligatoriasFaltantesModulo('RA') → registro-v3-app.js · SaveSection() → victim_case_v3_service.go
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

> `GS_14` y `GS_17` se contestan aquí pero se guardan con la sección **GS** (misma llave `GS_`): por eso el guardado reutiliza `enviarSeccion('GS')`. Las entidades y barreras ya se guardaron al agregarlas (RA-E03, RA-E06).

PASO 1 — «Guardar sección»: `obligatoriasFaltantesModulo('RA')` (`GS_14` y `GS_17` visibles y sin responder, `req: C1`); SI faltan: mensaje en rojo «Faltan obligatorias de 1er contacto: …» y TERMINAR
PASO 2 — [POST] /api/v1/victim-cases-v3/{icode}/seccion/GS  { answers: todas las `GS_*` } → `{ ok, triadaCompleta }`
  Si el guardado completa la tríada, el servidor materializa el borrador antes de pasar el Caso a Iniciado ([sub-flujo](../../Registro%20V3/Flujos/flow-subflujo-materializar-borrador.md)) y el cliente recarga entidades, barreras y `entitiesReadOnly`
PASO 3 — Mensaje «Sección guardada.»; una falla se avisa con el mismo aviso rojo de los demás módulos (`mostrarErrorGuardado`). **No avanza** al siguiente módulo: la barra «Siguiente →» sí funciona

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión | Paso afectado |
|---|---|
| `guardarYSalir` y `beforeunload` (WP-3) solo guardan CI, DP, CA, HV, AG y GS: agregar `RA` y enviar `GS_14` y `GS_17` («Para integración») | PASO 2 |
