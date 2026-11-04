━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona "Copiar radicado"
   Tipo: User Interaction
   Estado: ✅ Implementado y probado — WP-1 T3 (2026-11-02) (ver QA/E10-cuando-envia-reporte-qa-cases.md, casos C31 y C34)
   Función: copiarRadicado() — botón del paso 'confirmacion'
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: { radicado: this.radicado   → devuelto por E-10 }

PASO 1 — copiaError = false
PASO 2 — Copiar
  SI navigator.clipboard.writeText existe: writeText(radicado)
       → éxito: PASO 3 ; si lo rechaza (permiso denegado, contexto no seguro): respaldo
  SI NO: respaldo
  Respaldo: <textarea> temporal + document.execCommand('copy')
       → éxito: PASO 3 ; si falla: copiado = false ; copiaError = true (aviso «No fue posible copiarlo
         automáticamente: seleccione el número y cópielo manualmente.») → TERMINAR ejecución
PASO 3 — copiado = true  (el botón dice «Copiado»); un temporizador lo devuelve a «Copiar radicado» a los 3 s

No hay llamada a backend.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Ninguno.
