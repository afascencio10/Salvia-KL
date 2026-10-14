━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona "Procede recontacto"
   Tipo: User Interaction
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: clic en el botón primario "Procede recontacto"
               (step 'validacionPrevia', solo si canManage)

Precondición: step === 'validacionPrevia' Y canManage === true

INPUT: {} — no recibe datos, es una transición de vista pura


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — gestion-reporte-modal.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Cambiar de vista (sin backend)

  step = 'esquema3x3'

  // El mockup (línea 479) hace exactamente esto: goTo('recontacto') sin
  // ninguna llamada. El primer registro real en BD ocurre recién en M-05
  // o M-06, cuando el agente registra el primer intento. Hasta entonces,
  // "procede recontacto" es solo una decisión de navegación, no de negocio.

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Variable / decisión                                                              | Paso afectado |
|------------------------------------------------------------------------------------|---------------|
| ¿Si el agente llega hasta aquí y cierra el modal sin registrar ningún intento, la próxima vez que lo abra vuelve a 'validacionPrevia' o sigue en 'esquema3x3'? | PASO 1 — se resuelve en M-01 PASO 7: como no se persiste nada, vuelve a 'validacionPrevia' (0 intentos) |
