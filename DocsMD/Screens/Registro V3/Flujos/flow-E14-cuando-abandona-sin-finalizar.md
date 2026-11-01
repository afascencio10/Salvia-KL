━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando abandona la pantalla sin finalizar
   Tipo: User Interaction (primario) + Lifecycle (respaldo)
   Estado: ✅ Implementado — sincronizado con el código 2026-10-31
   Funciones: guardarYSalir() · volverACabezote() · listener 'beforeunload' de mounted() → registro-v3-app.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  etapa, moduloActivo, icode, answers   → estado local
}

CAMINO PRIMARIO — botón "Finalizar llamada (guardar y salir)" (solo visible en etapa 'formulario')

  SI moduloActivo ∉ CI/DP/CA/HV/AG/GS (es TZ):
    → volverACabezote() sin guardar nada (el Tamizaje se guarda con su propio botón)
    → TERMINAR ejecución

  SI NO:
    PASO 1 — Guardar el módulo activo sin exigir obligatorias
      guardando = true
      enviarSeccion(moduloActivo)   → [POST] /api/v1/victim-cases-v3/{icode}/seccion/{moduloActivo}  (ver flow-E12; sin validaciones)
      SI la petición falla: se ignora el error (best-effort)
    PASO 2 — Salir
      guardando = false
      volverACabezote():  etapa = 'cabezote'; icode = null; answers = {}; entidades = []; barreras = []
      (el Caso conserva su estado: Borrador o Iniciado; las entidades y barreras ya se guardaron al agregarlas)

CAMINO DE RESPALDO — cierre o recarga de la pestaña (window 'beforeunload')

  SI etapa !== 'formulario' o moduloActivo ∉ CI/DP/CA/HV/AG/GS:
    → no hace nada
  SI NO:
    → payload = { answers: solo las keys "{moduloActivo}_*" }
    → navigator.sendBeacon("/api/v1/victim-cases-v3/{icode}/seccion/{moduloActivo}", Blob JSON)   // fire-and-forget

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión | Paso afectado |
|---|---|
| Solo se guarda el módulo activo: lo editado en otros módulos sin guardar se pierde al salir | Ambos caminos |
| `sendBeacon` envía un Blob `application/json`; si el navegador lo rechaza (política de tipo de contenido) el respaldo no persiste. No verificado | Camino de respaldo |
| Por lo dicho en flow-E12/E10, guardar estando en GS deja `triadaCompleta` en false en el cliente, pero como se vuelve al Cabezote no tiene efecto visible | Camino primario |
